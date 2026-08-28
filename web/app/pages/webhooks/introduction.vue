<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiWebhook,
  mdiShieldCheck,
  mdiCodeBraces,
  mdiArrowRight,
  mdiInformation,
  mdiLockCheck,
  mdiCog,
  mdiKey,
  mdiSpeedometer,
  mdiClockOutline,
  mdiTrayFull,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('webhooksIntro.title')} - httpSMS`),
  description: computed(() => t('webhooksIntro.subtitle')),
  ogTitle: computed(() => `${t('webhooksIntro.title')} | httpSMS`),
  ogDescription: computed(() => t('webhooksIntro.subtitle')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})

const selectedTab = ref('node')
</script>

<template>
  <VContainer class="pt-8 pb-16">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <VCard class="pa-6 pa-md-10" elevation="2">
          <!-- Header Tag & Title -->
          <div class="d-flex align-center mb-3">
            <VChip color="info" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiWebhook" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              /webhooks/introduction
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('webhooksIntro.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('webhooksIntro.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- What are Webhooks -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('webhooksIntro.whatIs.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('webhooksIntro.whatIs.desc') }}
          </p>

          <VAlert
            color="info"
            variant="tonal"
            class="mb-6"
            :icon="mdiInformation"
          >
            {{ $t('webhooksIntro.whatIs.alert') }}
          </VAlert>

          <!-- Step 1: Configuration -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('webhooksIntro.config.title') }}
          </h2>
          <p class="text-body-medium text-medium-emphasis mb-4">
            {{ $t('webhooksIntro.config.desc') }}
          </p>

          <VCard variant="outlined" class="pa-4 mb-6">
            <div class="d-flex align-center mb-2">
              <VIcon color="primary" :icon="mdiCog" class="mr-2" />
              <h3 class="text-title-medium font-weight-bold">
                {{ $t('webhooksIntro.config.stepTitle') }}
              </h3>
            </div>
            <p class="text-body-medium text-medium-emphasis mb-3">
              {{ $t('webhooksIntro.config.stepDesc') }}
            </p>
            <VBtn color="primary" variant="tonal" size="small" to="/settings">
              {{ $t('webhooksIntro.config.goToSettings') }}
              <VIcon end :icon="mdiArrowRight" size="small" />
            </VBtn>
          </VCard>

          <!-- Step 2: Signature Verification -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="success" :icon="mdiShieldCheck" class="mr-1 mt-n1" />
            {{ $t('webhooksIntro.security.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('webhooksIntro.security.desc') }}
          </p>

          <VAlert
            color="success"
            variant="tonal"
            class="mb-6"
            :icon="mdiLockCheck"
          >
            {{ $t('webhooksIntro.security.secretNote') }}
          </VAlert>

          <!-- Code tabs for signature verification -->
          <VTabs
            v-model="selectedTab"
            color="primary"
            bg-color="#212121"
            show-arrows
            class="rounded-t"
          >
            <VTab value="node">Node.js (Express)</VTab>
            <VTab value="python">Python (Flask)</VTab>
            <VTab value="php">PHP</VTab>
            <VTab value="go">Go</VTab>
          </VTabs>

          <VTabsWindow v-model="selectedTab" v-highlight class="mb-8">
            <VTabsWindowItem value="node">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-javascript">const crypto = require('crypto');
const express = require('express');
const app = express();

// Use express.raw() to get raw request body buffer for HMAC verification
app.post('/webhook', express.raw({ type: 'application/json' }), (req, res) => {
  const signature = req.headers['x-signature'];
  const secret = 'SEU_WEBHOOK_SIGNING_SECRET';

  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(req.body)
    .digest('hex');

  if (signature !== expectedSignature) {
    return res.status(401).send('Invalid signature');
  }

  const event = JSON.parse(req.body.toString());
  console.log('Recebido evento:', event.type, event.data);

  res.status(200).send('OK');
});

app.listen(3000, () => console.log('Webhook server running on port 3000'));</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="python">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-python">import hmac
import hashlib
from flask import Flask, request, abort

app = Flask(__name__)

SECRET = b'SEU_WEBHOOK_SIGNING_SECRET'

@app.route('/webhook', methods=['POST'])
def webhook():
    signature = request.headers.get('x-signature')
    if not signature:
        abort(401)

    expected = hmac.new(SECRET, request.data, hashlib.sha256).hexdigest()
    if not hmac.compare_digest(signature, expected):
        abort(401)

    event = request.get_json()
    print(f"Recebido evento: {event['type']}")
    return 'OK', 200

if __name__ == '__main__':
    app.run(port=3000)</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="php">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-php">&lt;?php
$secret = 'SEU_WEBHOOK_SIGNING_SECRET';
$rawBody = file_get_contents('php://input');
$signature = $_SERVER['HTTP_X_SIGNATURE'] ?? '';

$expectedSignature = hash_hmac('sha256', $rawBody, $secret);

if (!hash_equals($expectedSignature, $signature)) {
    http_response_code(401);
    die('Invalid signature');
}

$event = json_decode($rawBody, true);
// Processar evento...
http_response_code(200);
echo 'OK';</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="go">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-go">package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "io"
    "net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    secret := []byte("SEU_WEBHOOK_SIGNING_SECRET")
    signature := r.Header.Get("x-signature")

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    mac := hmac.New(sha256.New, secret)
    mac.Write(body)
    expected := hex.EncodeToString(mac.Sum(nil))

    if signature != expected {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    http.HandleFunc("/webhook", webhookHandler)
    http.ListenAndServe(":3000", nil)
}</code></pre>
            </VTabsWindowItem>
          </VTabsWindow>

          <!-- Link to events list -->
          <VCard color="info" variant="tonal" class="pa-6">
            <div
              class="d-flex flex-column flex-sm-row align-sm-center justify-space-between"
            >
              <div>
                <h3 class="text-title-large font-weight-bold mb-1">
                  {{ $t('webhooksIntro.eventsLink.title') }}
                </h3>
                <p class="text-body-medium mb-2 mb-sm-0">
                  {{ $t('webhooksIntro.eventsLink.desc') }}
                </p>
              </div>
              <VBtn
                color="info"
                variant="flat"
                to="/webhooks/events"
                class="mt-2 mt-sm-0"
              >
                {{ $t('webhooksIntro.eventsLink.btn') }}
                <VIcon end :icon="mdiArrowRight" size="small" />
              </VBtn>
            </div>
          </VCard>
        </VCard>
      </VCol>

      <VCol v-if="mdAndUp" md="3">
        <VCard class="pa-6" elevation="2">
          <h3 class="text-title-medium font-weight-bold mb-3">
            {{ $t('docs.sidebar.guidesTitle') }}
          </h3>
          <VList density="compact" nav class="pa-0">
            <VListItem
              to="/docs"
              :prepend-icon="mdiCodeBraces"
              :title="$t('docs.sections.quickstart.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/webhooks/introduction"
              :prepend-icon="mdiWebhook"
              :title="$t('docs.sections.webhooks.title')"
              class="mb-1 rounded"
              active
            />
            <VListItem
              to="/webhooks/events"
              :prepend-icon="mdiWebhook"
              :title="$t('docs.sections.webhookEvents.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/phone-api-keys"
              :prepend-icon="mdiKey"
              :title="$t('docs.sections.phoneApiKeys.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/control-sms-send-rate"
              :prepend-icon="mdiSpeedometer"
              :title="$t('docs.sections.rateLimit.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/scheduling-sms-messages"
              :prepend-icon="mdiClockOutline"
              :title="$t('docs.sections.scheduling.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/outgoing-message-queue"
              :prepend-icon="mdiTrayFull"
              :title="$t('docs.sections.queue.title')"
              class="mb-1 rounded"
            />
          </VList>
        </VCard>
      </VCol>
    </VRow>
  </VContainer>
</template>
