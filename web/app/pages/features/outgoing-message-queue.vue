<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiTrayFull,
  mdiClockOutline,
  mdiSpeedometer,
  mdiCalendarRange,
  mdiInformation,
  mdiKey,
  mdiWebhook,
  mdiCodeBraces,
  mdiArrowRight,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('featureQueue.title')} - httpSMS`),
  description: computed(() => t('featureQueue.subtitle')),
  ogTitle: computed(() => `${t('featureQueue.title')} | httpSMS`),
  ogDescription: computed(() => t('featureQueue.subtitle')),
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
            <VChip color="cyan" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiTrayFull" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              /features/outgoing-message-queue
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('featureQueue.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('featureQueue.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Introduction -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('featureQueue.intro.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureQueue.intro.desc') }}
          </p>

          <VAlert
            color="info"
            variant="tonal"
            class="mb-6"
            :icon="mdiInformation"
          >
            {{ $t('featureQueue.intro.alert') }}
          </VAlert>

          <!-- 3 Priority Dispatch Rules -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('featureQueue.rules.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('featureQueue.rules.desc') }}
          </p>

          <VRow class="mb-6">
            <VCol cols="12" md="4">
              <VCard variant="outlined" class="pa-4 h-100">
                <div class="d-flex align-center mb-2">
                  <VAvatar
                    color="purple"
                    size="32"
                    variant="tonal"
                    class="mr-2"
                  >
                    <VIcon :icon="mdiClockOutline" size="18" />
                  </VAvatar>
                  <h3 class="text-title-medium font-weight-bold">
                    {{ $t('featureQueue.rules.rule1Title') }}
                  </h3>
                </div>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('featureQueue.rules.rule1Desc') }}
                </p>
              </VCard>
            </VCol>

            <VCol cols="12" md="4">
              <VCard variant="outlined" class="pa-4 h-100">
                <div class="d-flex align-center mb-2">
                  <VAvatar
                    color="success"
                    size="32"
                    variant="tonal"
                    class="mr-2"
                  >
                    <VIcon :icon="mdiSpeedometer" size="18" />
                  </VAvatar>
                  <h3 class="text-title-medium font-weight-bold">
                    {{ $t('featureQueue.rules.rule2Title') }}
                  </h3>
                </div>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('featureQueue.rules.rule2Desc') }}
                </p>
              </VCard>
            </VCol>

            <VCol cols="12" md="4">
              <VCard variant="outlined" class="pa-4 h-100">
                <div class="d-flex align-center mb-2">
                  <VAvatar color="info" size="32" variant="tonal" class="mr-2">
                    <VIcon :icon="mdiCalendarRange" size="18" />
                  </VAvatar>
                  <h3 class="text-title-medium font-weight-bold">
                    {{ $t('featureQueue.rules.rule3Title') }}
                  </h3>
                </div>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t('featureQueue.rules.rule3Desc') }}
                </p>
              </VCard>
            </VCol>
          </VRow>

          <!-- Bulk Send API -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('featureQueue.bulkApi.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('featureQueue.bulkApi.desc') }}
          </p>

          <pre
            class="pa-4 bg-surface rounded mb-6"
          ><code class="language-bash">curl --location --request POST '{{ appStore.appData.apiBaseUrl }}/v1/messages/bulk-send' \
--header 'x-api-key: SUA_CHAVE_DE_API' \
--header 'Content-Type: application/json' \
--data-raw '[
  {
    "from": "+5511999999999",
    "to": "+5511988888881",
    "content": "Olá João, seu pedido foi enviado!"
  },
  {
    "from": "+5511999999999",
    "to": "+5511988888882",
    "content": "Olá Maria, seu pedido foi enviado!"
  }
]'</code></pre>

          <VBtn color="primary" variant="flat" size="large" to="/settings">
            {{ $t('featureQueue.settingsBtn') }}
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
              active
            />
          </VList>
        </VCard>
      </VCol>
    </VRow>
  </VContainer>
</template>
