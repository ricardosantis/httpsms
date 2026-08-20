<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiSpeedometer,
  mdiAlertCircle,
  mdiShieldCheck,
  mdiCog,
  mdiArrowRight,
  mdiKey,
  mdiWebhook,
  mdiClockOutline,
  mdiTrayFull,
  mdiCodeBraces,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('featureControlRate.title')} - httpSMS`),
  description: computed(() => t('featureControlRate.subtitle')),
  ogTitle: computed(() => `${t('featureControlRate.title')} | httpSMS`),
  ogDescription: computed(() => t('featureControlRate.subtitle')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer class="pt-8 pb-16">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <VCard class="pa-6 pa-md-10" elevation="2">
          <!-- Header Tag & Title -->
          <div class="d-flex align-center mb-3">
            <VChip color="success" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiSpeedometer" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              /features/control-sms-send-rate
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('featureControlRate.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('featureControlRate.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Why Control Send Rate -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('featureControlRate.why.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureControlRate.why.desc1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureControlRate.why.desc2') }}
          </p>

          <VAlert
            color="warning"
            variant="tonal"
            class="mb-6"
            :icon="mdiAlertCircle"
          >
            {{ $t('featureControlRate.why.carrierAlert') }}
          </VAlert>

          <!-- How to Configure -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="primary" :icon="mdiCog" class="mr-1 mt-n1" />
            {{ $t('featureControlRate.config.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureControlRate.config.desc') }}
          </p>

          <VRow class="mb-6">
            <VCol cols="12" sm="6">
              <VCard variant="outlined" class="pa-4 h-100">
                <h3
                  class="text-title-medium font-weight-bold text-primary mb-2"
                >
                  {{ $t('featureControlRate.config.defaultLimitTitle') }}
                </h3>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('featureControlRate.config.defaultLimitDesc') }}
                </p>
              </VCard>
            </VCol>
            <VCol cols="12" sm="6">
              <VCard variant="outlined" class="pa-4 h-100">
                <h3
                  class="text-title-medium font-weight-bold text-success mb-2"
                >
                  {{ $t('featureControlRate.config.maxLimitTitle') }}
                </h3>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('featureControlRate.config.maxLimitDesc') }}
                </p>
              </VCard>
            </VCol>
          </VRow>

          <!-- Queue Delay Calculation -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="success" :icon="mdiShieldCheck" class="mr-1 mt-n1" />
            {{ $t('featureControlRate.formula.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('featureControlRate.formula.desc') }}
          </p>

          <pre
            class="pa-4 bg-surface rounded mb-6"
          ><code class="language-text">Delay (em segundos) = Índice da Mensagem × (60 / Taxa por Minuto)

Exemplo com taxa de 10 mensagens/min:
- Mensagem 0: 0s (imediato)
- Mensagem 1: 1 × (60 / 10) = 6s
- Mensagem 2: 2 × (60 / 10) = 12s
- Mensagem 3: 3 × (60 / 10) = 18s</code></pre>

          <VBtn color="primary" variant="flat" size="large" to="/settings">
            {{ $t('featureControlRate.config.settingsBtn') }}
            <VIcon end :icon="mdiArrowRight" size="small" />
          </VBtn>
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
              active
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
