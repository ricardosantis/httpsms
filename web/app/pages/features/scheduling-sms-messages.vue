<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiClockOutline,
  mdiCalendarClock,
  mdiMicrosoftExcel,
  mdiInformation,
  mdiKey,
  mdiWebhook,
  mdiSpeedometer,
  mdiTrayFull,
  mdiCodeBraces,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('featureScheduling.title')} - httpSMS`),
  description: computed(() => t('featureScheduling.subtitle')),
  ogTitle: computed(() => `${t('featureScheduling.title')} | httpSMS`),
  ogDescription: computed(() => t('featureScheduling.subtitle')),
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
            <VChip color="purple" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiClockOutline" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              /features/scheduling-sms-messages
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('featureScheduling.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('featureScheduling.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Introduction -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('featureScheduling.howItWorks.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureScheduling.howItWorks.desc') }}
          </p>

          <VAlert
            color="info"
            variant="tonal"
            class="mb-6"
            :icon="mdiInformation"
          >
            {{ $t('featureScheduling.howItWorks.alertMaxDays') }}
          </VAlert>

          <!-- API Scheduling with send_at -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon
              color="primary"
              :icon="mdiCalendarClock"
              class="mr-1 mt-n1"
            />
            {{ $t('featureScheduling.api.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('featureScheduling.api.desc') }}
          </p>

          <pre
            class="pa-4 bg-surface rounded mb-6"
          ><code class="language-bash">curl --location --request POST '{{ appStore.appData.apiBaseUrl }}/v1/messages/send' \
--header 'x-api-key: SUA_CHAVE_DE_API' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "+5511999999999",
    "to": "+5511988888888",
    "content": "Lembrete: sua consulta está agendada para amanhã às 14:00.",
    "send_at": "2026-08-21T14:00:00-03:00"
}'</code></pre>

          <!-- Bulk Scheduling -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon
              color="success"
              :icon="mdiMicrosoftExcel"
              class="mr-1 mt-n1"
            />
            {{ $t('featureScheduling.bulk.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureScheduling.bulk.desc') }}
          </p>

          <VCard variant="outlined" class="pa-4 mb-6">
            <h4 class="text-title-medium font-weight-bold mb-2">
              {{ $t('featureScheduling.bulk.columnTitle') }}
            </h4>
            <p class="text-body-medium text-medium-emphasis mb-0">
              {{ $t('featureScheduling.bulk.columnDesc') }}
            </p>
          </VCard>

          <VBtn color="primary" variant="flat" size="large" to="/bulk-messages">
            {{ $t('featureScheduling.bulk.bulkBtn') }}
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
            />
            <VListItem
              to="/features/scheduling-sms-messages"
              :prepend-icon="mdiClockOutline"
              :title="$t('docs.sections.scheduling.title')"
              class="mb-1 rounded"
              active
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
