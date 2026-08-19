# AGENTS.md

httpSMS turns an Android phone into an SMS gateway via a Go HTTP API. This is a fork of `NdoleStudio/httpsms` with an added integration-test suite and Render self-hosting blueprint.

## Layout

- `api/` — Go API. Module `github.com/NdoleStudio/httpsms` (Go 1.25), Fiber v2, GORM (Postgres/CockroachDB) + Redis. Entrypoint `api/main.go`; all logic under `api/pkg/` (handlers, services, repositories, entities, events, listeners, di). Deps wired in `api/pkg/di`; OpenAPI comments are inline swagger annotations.
- `web/` — Nuxt 4 + Vuetify 4 SPA (`ssr: false`), static host. Uses **pnpm** (lockfile `web/pnpm-lock.yaml`, CI uses pnpm 11 / Node 24). Package manager is pnpm, not npm/yarn.
- `android/` — Kotlin Android client. Requires Firebase `google-services.json` to build; `android/app/google-services.json` exists but is a placeholder.
- `tests/` — Standalone Go module (`github.com/NdoleStudio/httpsms/tests`) for E2E integration tests. Docker-based, see below.
- `docs/superpowers/specs/` — design specs/plans (e.g. the Nuxt4/Vuetify4 migration).

## Commands

API (run from `api/`):
- `go test ./...` — unit tests
- `swag init --requiredByDefault --parseDependency --parseInternal` — regenerate `api/docs/` (swagger.json) after changing route handlers. Run inside the Docker build, so usually only needed locally to refresh swagger.json for the web client types.
- `go build -o http-sms .` then `./http-sms` — manual run. `main.go` only loads `.env` when run with **no CLI args**; the Docker image runs `--dotenv=false` (env from container). All config via env vars; see `api/.env.docker` for the full list.

Web (run from `web/`):
- `pnpm install`
- `pnpm dev` — dev server on `:3000`, proxies to `API_BASE_URL` (default `http://localhost:8000`)
- `pnpm lint` — eslint + stylelint + prettier (all three; CI runs this)
- `pnpm lintfix` — autofix
- `pnpm api:models` — regenerates `web/shared/types/api.ts` from `../api/docs/swagger.json`. **If you change API responses, regenerate this or the web build/tests break.**
- `pnpm run generate` — static build (used by CI and both Dockerfiles); the build **requires** `.env` to exist (CI does `mv .env.production .env`).

## Full local stack (root `docker-compose.yml`)

`docker compose up --build` runs postgres, redis, api (:8000), web (:3000). Requires `api/.env` (copy from `api/.env.docker`) and `web/.env` (copy from `web/.env.docker`). Postgres on host port **5435** (not 5432). After first boot, insert the events-queue system user into `users` matching `EVENTS_QUEUE_USER_ID`/`EVENTS_QUEUE_USER_API_KEY` or async event processing fails; restart the api container after changing those vars.

## Integration tests (`tests/`)

Full-stack E2E. Stack = API + CockroachDB (v22.2, in-memory) + Redis + MongoDB (heartbeat backend, `HEARTBEAT_DB_BACKEND=mongodb`) + **WireMock** (mocks FCM push and webhooks via `tests/wiremock/mappings/`; `FCM_ENDPOINT=http://wiremock:8080`). No real Firebase/emulator needed — fake credentials are generated locally.

Prereqs: Docker, Go, `jq`, `openssl`.

Run from `tests/`:
```bash
bash generate-firebase-credentials.sh
export FIREBASE_CREDENTIALS=$(jq -c . firebase-credentials.json)
docker compose up -d --build --wait
docker compose wait seed && sleep 2   # seed must run after API migrations
go test -v -timeout 120s ./...
docker compose down -v
```

Notes:
- The seed container must finish after the API healthcheck (GORM migrations) — `docker compose wait seed` handles ordering.
- **Root `docker-compose.yml` and `tests/docker-compose.yml` both map host ports 8000/6379** — don't run both stacks at once.
- `tests/README.md` is stale: it describes an "emulator" service that no longer exists; the actual stack uses WireMock (see `tests/docker-compose.yml`).
- API changes should add tests to `integration_test.go`; helpers: `doRequest()`, `pollMessageStatus()`, `waitForFCMPush()`, `waitForWebhookEvents()`. Seed data lives in `tests/seed.sql` (test user `test-user-api-key`, phone `pk_test-phone-api-key`).

## Production deployment (this fork — httpsms-prod)

Fork deployed as an SMS gateway for **mesaquevende.com.br** (SaaS in `/home/ricardosantis/projetos/mesaquevende.com.br`). Stack:

- **API**: GCP Cloud Run `http-sms-api` (us-east1), image built via Cloud Build (kaniko). Domain `api.mesaquevende.com.br`.
- **Web**: Firebase Hosting site `httpsms-prod`, SPA built with Nuxt 4 `pnpm run generate`. Domain `sms.mesaquevende.com.br`.
- **DB**: Cloud SQL Postgres 15 `httpsms-postgres`, **private IP only** (`10.13.0.3`), connected **directly** from Cloud Run via VPC connector `httpsms-connector` (private-ranges-only egress). Do NOT use `--add-cloudsql-instances` (the built-in proxy sidecar does not get injected in this setup).
- **Redis**: Memorystore `httpsms-redis` private IP `10.192.49.99:6379` (also reached via the VPC connector).
- **Events queue**: Cloud Tasks queue `httpsms-events` (us-east1), full path `projects/479893056966/locations/us-east1/queues/httpsms-events`. `EVENTS_QUEUE_TYPE=cloud-task`, `EVENTS_QUEUE_ENDPOINT=https://api.mesaquevende.com.br/v1/events`.
- **Auth**: Firebase Auth (email/password only — Google/GitHub OAuth providers were NOT configured), project `httpsms-prod` (billing account Principal/Blaze 0143B4). API validates Firebase ID tokens (`bearer_auth_middleware.go`).
- **CI**: GitHub Actions on `main`. `api.yml` runs integration tests then `gcloud builds submit . --config=api/cloudbuild.yaml` via Workload Identity Federation (secrets `GCP_WORKLOAD_IDENTITY_PROVIDER`, `GCP_SERVICE_ACCOUNT=httpsms-api-sa@httpsms-prod.iam.gserviceaccount.com`). `web.yml` runs lint+generate then Firebase Hosting deploy (secret `FIREBASE_SERVICE_ACCOUNT_HTTPSMS_PROD`).

Key prod env vars for the Cloud Run service (all on the service, secrets via Secret Manager `--set-secrets`): `ENV=production`, `GCP_PROJECT_ID=httpsms-prod`, `APP_HOST=0.0.0.0`, `APP_PORT=8000`, `SWAGGER_HOST=api.mesaquevende.com.br`, `APP_URL=https://sms.mesaquevende.com.br`, `EVENTS_QUEUE_NAME=<full tasks path>`, `FIREBASE_CREDENTIALS`, `DATABASE_URL`/`DATABASE_URL_DEDICATED` (point at `10.13.0.3:5432`), `REDIS_URL`, `EVENTS_QUEUE_USER_ID`/`EVENTS_QUEUE_USER_API_KEY`, `SMTP_*` (mailtrap placeholder; optional for SMS core), `GCS_BUCKET_NAME=` (empty → in-memory MMS).

### System user (events queue)

Async event processing (Cloud Tasks → `POST /v1/events`) requires a row in `users` whose `id` == `EVENTS_QUEUE_USER_ID` and `api_key` == `EVENTS_QUEUE_USER_API_KEY`. The events handler 403s unless `userIDFromContext(c) == EVENTS_QUEUE_USER_ID`. Upstream has **no auto-seed** — insert manually after first deploy (GORM migrations create the tables). Note the `users` column is `notification_webhook_enabled` (not `webhook_enabled`). Current system user: id `1c07f822b4ab1731`, api key `6f6d88d038754b5a1e6eba477f7576d120726cd0a15db591` (also in Secret Manager). Local psql admin access requires temporarily re-enabling the public IP + authorized network, or a compute VM inside the VPC.

### Local-only fixes in this fork (diff vs upstream)

- **Axiom removed**: `api/pkg/di/container.go` — `logDriver()` falls back to console logger when `AXIOM_TOKEN` is empty; `InitializeTraceProvider()` falls back to Google trace when `AXIOM_TOKEN` is empty. Run prod with `ENV=production` and no Axiom vars.
- **Turnstile optional**: `turnstile_token_validator.go` (`if v.secretKey == "" { return true }`) and `message_handler_validator.go` (only requires `token` when the secret key is configured). `CLOUDFLARE_TURNSTILE_SECRET_KEY`/`SITE_KEY` are absent in prod.
- The web `.env.production` was rewritten for prod (`API_BASE_URL=https://api.mesaquevende.com.br`, `APP_URL=https://sms.mesaquevende.com.br`, Firebase config for `httpsms-prod`); upstream CHECKOUT/PUSHER/TURNSTILE vars removed (nuxt.config uses `|| ''` fallbacks).

## Deployment learnings (GCP / Cloudflare gotchas)

- **Cloud Run + Cloud SQL**: the built-in connector annotation (`run.googleapis.com/cloudsql-instances`) did NOT inject the proxy sidecar; simplest reliable path was connecting directly to the instance private IP over a VPC connector. Set `--vpc-connector` + `--set-env-vars`/`--set-secrets` on `gcloud run deploy`; secrets pinned to a version (`secretName:3`).
- **Secret Manager values must have no trailing newline** — `echo` appends `\n` and Postgres DSNs fail with `invalid control character in URL`. Use `printf` (or `--data-file`).
- **Cloud Build kaniko**: cloudbuild.yaml uses `dir: api` + `--context=.`, so source must be the **repo root** (`gcloud builds submit .`), never `api/`. `$SHORT_SHA` is empty on manual submits (only set on triggers) — the deploy step must reference `:latest`.
- **Cloud Run custom domains require Search Console (Webmaster Central) verification** of the domain with the owning Google account (`ricardo.santis@gmail.com`). `gcloud` tokens lack the `siteverification` scope; do it interactively: add property at search.google.com/search-console → add the `google-site-verification` TXT to Cloudflare → click Verify. After verification, delete and recreate the DomainMapping (old one keeps `PermissionDenied`), then add `CNAME api → ghs.googlehosted.com`.
- **Firebase Hosting custom domains**: there is no CLI command — use `POST https://firebasehosting.googleapis.com/v1beta1/projects/{p}/sites/{s}/customDomains?customDomainId=<host>` with `Authorization: Bearer` + **`X-Goog-User-Project: <project>`** (required or you get 403 SERVICE_DISABLED). Add the returned CNAME + `_acme-challenge` TXT to Cloudflare.
- **Cloudflare global API key** (`cfk_…`): use `X-Auth-Email` + `X-Auth-Key` headers (NOT `Authorization: Bearer`, which returns 9109). Wrangler OAuth scopes do NOT include DNS write, so zone record changes must go through the API v4 (`zones/{id}/dns_records`).
- **Firebase Auth initialize**: `POST https://identitytoolkit.googleapis.com/v2/projects/{p}/identityPlatform:initializeAuth` (empty body, `X-Goog-User-Project` header). Enable email/password via `PATCH …/admin/v2/projects/{p}/config?updateMask=signIn.email`. Authorized domains must be patched at the **root level** (`authorizedDomains`), not under `signIn`.
- **gh auth accounts**: the CLI has 3 accounts (ricardosantisinc active, webtechnegocios, ricardosantis). Use `gh auth switch --user ricardosantis` before `gh secret set` on this fork's repo.

## Repo conventions

- Pre-commit (`.pre-commit-config.yaml`) enforces on Go: `gofumpt`, `go mod tidy` (check), `goimports`, `go-lint`; prettier excludes `web/` (web has its own prettier via husky/lint-staged).
- Web commits use conventional commits (`@commitlint/config-conventional`) via husky.
- No existing instruction files (CLAUDE.md/.cursor rules) in this repo.