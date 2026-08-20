<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiKey,
  mdiAlertCircle,
  mdiCellphone,
  mdiQrcodeScan,
  mdiTrashCanOutline,
  mdiLinkOff,
  mdiDownload,
  mdiArrowRight,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('featurePhoneApiKeys.title')} - httpSMS`),
  description: computed(() => t('featurePhoneApiKeys.subtitle')),
  ogTitle: computed(() => `${t('featurePhoneApiKeys.title')} | httpSMS`),
  ogDescription: computed(() => t('featurePhoneApiKeys.subtitle')),
  ogImage: '/img/phone-api-key.png',
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
            <VChip color="primary" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiKey" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              /features/phone-api-keys
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('featurePhoneApiKeys.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('featurePhoneApiKeys.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Intro Paragraph -->
          <p class="text-body-1 leading-relaxed mb-6">
            {{ $t('featurePhoneApiKeys.intro') }}
          </p>

          <!-- Important Note Alert -->
          <VAlert
            color="info"
            variant="tonal"
            class="mb-8"
            border="start"
            :icon="mdiAlertCircle"
          >
            <div class="text-subtitle-2 font-weight-bold mb-1">
              {{ $t('featurePhoneApiKeys.warningTitle') }}
            </div>
            <div class="text-body-2">
              {{ $t('featurePhoneApiKeys.warningText') }}
              <NuxtLink
                to="/settings"
                class="font-weight-bold text-decoration-none text-info"
              >
                {{ $t('featurePhoneApiKeys.settingsLink') }}
              </NuxtLink>
              (<code>{{ appStore.appData.url }}/settings</code>).
            </div>
          </VAlert>

          <!-- Section 1: Create Phone API Key -->
          <h2 class="text-headline-medium font-weight-bold mt-8 mb-3">
            <VIcon :icon="mdiCellphone" color="primary" class="mr-1" />
            {{ $t('featurePhoneApiKeys.createTitle') }}
          </h2>

          <p class="text-body-1 mb-4">
            {{ $t('featurePhoneApiKeys.createStep1') }}
            <NuxtLink
              to="/phone-api-keys"
              class="text-primary font-weight-bold text-decoration-none"
            >
              {{ appStore.appData.url }}/phone-api-keys
            </NuxtLink>
            {{ $t('featurePhoneApiKeys.createStep2') }}
            <VChip size="small" color="primary" class="font-weight-bold mx-1">
              {{ $t('featurePhoneApiKeys.createBtnText') }}
            </VChip>
            .
          </p>

          <VCard
            variant="outlined"
            class="my-6 pa-2 bg-grey-darken-4 overflow-hidden"
          >
            <VImg
              src="/img/phone-api-key.png"
              alt="Phone API Keys Interface"
              cover
              class="rounded"
            />
            <div class="text-caption text-center text-medium-emphasis mt-2">
              {{ $t('featurePhoneApiKeys.title') }} - Painel de Controle
            </div>
          </VCard>

          <p class="text-body-1 mb-4">
            {{ $t('featurePhoneApiKeys.createStep3') }}
            <VChip size="small" color="info" class="font-weight-bold mx-1">
              {{ $t('featurePhoneApiKeys.viewBtnText') }}
            </VChip>
            {{ $t('featurePhoneApiKeys.createStep4') }}
          </p>

          <!-- APK Download Alert -->
          <VAlert
            color="primary"
            variant="tonal"
            class="my-6"
            :icon="mdiQrcodeScan"
          >
            <div
              class="d-flex flex-column flex-sm-row align-sm-center justify-space-between gap-4"
            >
              <div class="text-body-2 mr-sm-4">
                {{ $t('featurePhoneApiKeys.appDownloadNotice') }}
              </div>
              <VBtn
                color="primary"
                variant="flat"
                :href="appStore.appData.appDownloadUrl"
                target="_blank"
                class="text-none mt-3 mt-sm-0"
              >
                <VIcon start :icon="mdiDownload" />
                {{ $t('featurePhoneApiKeys.downloadAppBtn') }}
              </VBtn>
            </div>
          </VAlert>

          <!-- Section 2: Remove Phone from Key -->
          <h2 class="text-headline-medium font-weight-bold mt-10 mb-3">
            <VIcon :icon="mdiLinkOff" color="warning" class="mr-1" />
            {{ $t('featurePhoneApiKeys.removeTitle') }}
          </h2>

          <p class="text-body-1 mb-3">
            {{ $t('featurePhoneApiKeys.removeDesc1') }}
          </p>

          <p class="text-body-1 mb-6">
            {{ $t('featurePhoneApiKeys.removeDesc2') }}
          </p>

          <!-- Section 3: Delete Phone API Key -->
          <h2 class="text-headline-medium font-weight-bold mt-10 mb-3">
            <VIcon :icon="mdiTrashCanOutline" color="error" class="mr-1" />
            {{ $t('featurePhoneApiKeys.deleteTitle') }}
          </h2>

          <p class="text-body-1 mb-6">
            {{ $t('featurePhoneApiKeys.deleteDesc') }}
          </p>

          <!-- Quick Action Buttons -->
          <VDivider class="my-8" />

          <div
            class="d-flex flex-wrap gap-4 align-center justify-space-between"
          >
            <VBtn
              color="primary"
              variant="flat"
              to="/phone-api-keys"
              class="text-none"
            >
              {{ $t('featurePhoneApiKeys.goToApiKeys') }}
              <VIcon end :icon="mdiArrowRight" />
            </VBtn>

            <VBtn variant="tonal" to="/settings" class="text-none">
              {{ $t('nav.settings') }}
            </VBtn>
          </div>
        </VCard>
      </VCol>

      <!-- Sidebar Column -->
      <VCol cols="12" md="3" class="mt-6 mt-md-0">
        <BlogSidebar />
      </VCol>
    </VRow>
  </VContainer>
</template>

<style scoped>
.leading-relaxed {
  line-height: 1.75;
}

.gap-4 {
  gap: 16px;
}
</style>
