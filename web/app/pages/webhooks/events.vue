<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiWebhook,
  mdiMessageText,
  mdiCheckCircle,
  mdiAlertCircle,
  mdiHeartPulse,
  mdiInformation,
  mdiKey,
  mdiSpeedometer,
  mdiClockOutline,
  mdiTrayFull,
  mdiCodeBraces,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('webhookEvents.title')} - httpSMS`),
  description: computed(() => t('webhookEvents.subtitle')),
  ogTitle: computed(() => `${t('webhookEvents.title')} | httpSMS`),
  ogDescription: computed(() => t('webhookEvents.subtitle')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})

const selectedEventTab = ref('received')
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
              /webhooks/events
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('webhookEvents.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('webhookEvents.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- CloudEvents Standard Info -->
          <VAlert
            color="info"
            variant="tonal"
            class="mb-6"
            :icon="mdiInformation"
          >
            <p class="mb-0">
              {{ $t('webhookEvents.cloudEventsInfo') }}
            </p>
          </VAlert>

          <!-- Events Navigation Tabs -->
          <VTabs
            v-model="selectedEventTab"
            color="primary"
            bg-color="#212121"
            show-arrows
            class="rounded-t"
          >
            <VTab value="received">
              <VIcon color="success" class="mr-2" :icon="mdiMessageText" />
              message.phone.received
            </VTab>
            <VTab value="sent">
              <VIcon color="info" class="mr-2" :icon="mdiCheckCircle" />
              message.phone.sent
            </VTab>
            <VTab value="failed">
              <VIcon color="error" class="mr-2" :icon="mdiAlertCircle" />
              message.phone.failed
            </VTab>
            <VTab value="delivered">
              <VIcon color="primary" class="mr-2" :icon="mdiCheckCircle" />
              message.phone.delivered
            </VTab>
            <VTab value="heartbeatReceived">
              <VIcon color="success" class="mr-2" :icon="mdiHeartPulse" />
              phone.heartbeat.received
            </VTab>
            <VTab value="heartbeatMissed">
              <VIcon color="warning" class="mr-2" :icon="mdiAlertCircle" />
              phone.heartbeat.missed
            </VTab>
          </VTabs>

          <VTabsWindow v-model="selectedEventTab" v-highlight class="mb-8">
            <!-- 1. message.phone.received -->
            <VTabsWindowItem value="received">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.received') }}
                </p>
                <pre><code class="language-json">{
  "id": "e8a2b537-8e6d-4957-8fb1-2a6288ea56f9",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "message.phone.received",
  "data": {
    "id": "01HGB7M7E0V105S8D36G292C23",
    "owner": "1c07f822b4ab1731",
    "user_id": "1c07f822b4ab1731",
    "content": "Olá! Gostaria de saber mais sobre os planos do produto.",
    "from": "+5511988888888",
    "to": "+5511999999999",
    "type": "received",
    "status": "RECEIVED",
    "order_date": "2026-08-20T18:00:00Z",
    "created_at": "2026-08-20T18:00:00Z",
    "updated_at": "2026-08-20T18:00:00Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>

            <!-- 2. message.phone.sent -->
            <VTabsWindowItem value="sent">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.sent') }}
                </p>
                <pre><code class="language-json">{
  "id": "f912c448-9f7e-5068-9ac2-3b7399fb67a0",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "message.phone.sent",
  "data": {
    "id": "01HGB7M7E0V105S8D36G292C23",
    "owner": "1c07f822b4ab1731",
    "user_id": "1c07f822b4ab1731",
    "content": "Seu código de confirmação é 849201.",
    "from": "+5511999999999",
    "to": "+5511988888888",
    "type": "sent",
    "status": "SENT",
    "order_date": "2026-08-20T18:00:05Z",
    "created_at": "2026-08-20T18:00:00Z",
    "updated_at": "2026-08-20T18:00:05Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>

            <!-- 3. message.phone.failed -->
            <VTabsWindowItem value="failed">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.failed') }}
                </p>
                <pre><code class="language-json">{
  "id": "a023d559-0a8f-6179-abde-4c8400ac78b1",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "message.phone.failed",
  "data": {
    "id": "01HGB7M7E0V105S8D36G292C23",
    "owner": "1c07f822b4ab1731",
    "user_id": "1c07f822b4ab1731",
    "content": "Seu código de confirmação é 849201.",
    "from": "+5511999999999",
    "to": "+5511988888888",
    "type": "sent",
    "status": "FAILED",
    "failure_reason": "RESULT_ERROR_GENERIC_FAILURE",
    "order_date": "2026-08-20T18:00:05Z",
    "created_at": "2026-08-20T18:00:00Z",
    "updated_at": "2026-08-20T18:00:05Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>

            <!-- 4. message.phone.delivered -->
            <VTabsWindowItem value="delivered">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.delivered') }}
                </p>
                <pre><code class="language-json">{
  "id": "b134e660-1b90-7280-bcef-5d9511bd89c2",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "message.phone.delivered",
  "data": {
    "id": "01HGB7M7E0V105S8D36G292C23",
    "owner": "1c07f822b4ab1731",
    "user_id": "1c07f822b4ab1731",
    "content": "Seu código de confirmação é 849201.",
    "from": "+5511999999999",
    "to": "+5511988888888",
    "type": "sent",
    "status": "DELIVERED",
    "order_date": "2026-08-20T18:00:08Z",
    "created_at": "2026-08-20T18:00:00Z",
    "updated_at": "2026-08-20T18:00:08Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>

            <!-- 5. phone.heartbeat.received -->
            <VTabsWindowItem value="heartbeatReceived">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.heartbeatReceived') }}
                </p>
                <pre><code class="language-json">{
  "id": "c245f771-2c01-8391-cdfa-6ea622ce90d3",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "phone.heartbeat.received",
  "data": {
    "id": "01HGB8A1F1W216T9E47H303D34",
    "phone_number": "+5511999999999",
    "user_id": "1c07f822b4ab1731",
    "status": "ONLINE",
    "timestamp": "2026-08-20T18:00:00Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>

            <!-- 6. phone.heartbeat.missed -->
            <VTabsWindowItem value="heartbeatMissed">
              <div class="pa-4 bg-surface rounded-b">
                <p class="text-body-medium mb-3">
                  {{ $t('webhookEvents.types.heartbeatMissed') }}
                </p>
                <pre><code class="language-json">{
  "id": "d3560882-3d12-9402-de0b-7fb733df01e4",
  "source": "httpsms.com/users/1c07f822b4ab1731",
  "specversion": "1.0",
  "type": "phone.heartbeat.missed",
  "data": {
    "id": "01HGB8A1F1W216T9E47H303D34",
    "phone_number": "+5511999999999",
    "user_id": "1c07f822b4ab1731",
    "status": "OFFLINE",
    "last_seen_at": "2026-08-20T17:45:00Z",
    "timestamp": "2026-08-20T18:00:00Z"
  }
}</code></pre>
              </div>
            </VTabsWindowItem>
          </VTabsWindow>
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
              active
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
