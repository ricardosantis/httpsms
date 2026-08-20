<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiFileDocumentOutline,
  mdiShieldAccount,
  mdiAlertCircle,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const { t } = useI18n()

definePageMeta({
  layout: 'website',
})

useSeoMeta({
  title: computed(() => `${t('terms.title')} - httpSMS`),
  description: computed(() => t('terms.subtitle')),
  ogTitle: computed(() => `${t('terms.title')} | httpSMS`),
  ogDescription: computed(() => t('terms.subtitle')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer class="pt-8 pb-16">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="10" offset-md="1">
        <VCard class="pa-6 pa-md-10" elevation="2">
          <div class="d-flex align-center mb-3">
            <VChip color="primary" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiFileDocumentOutline" size="small" />
              {{ $t('website.legal') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">
              {{ $t('terms.lastUpdated') }}
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('terms.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Section 1 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('terms.section1.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section1.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.section1.p2') }}
          </p>

          <!-- Section 2 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('terms.section2.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section2.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section2.p2') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.section2.p3') }}
          </p>

          <!-- Section 3 -->
          <h2 class="text-headline-small font-weight-bold mb-3 text-error">
            <VIcon color="error" :icon="mdiAlertCircle" class="mr-1 mt-n1" />
            {{ $t('terms.section3.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section3.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.section3.p2') }}
          </p>

          <!-- Section 4 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('terms.section4.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section4.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section4.p2') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.section4.p3') }}
          </p>

          <!-- Section 5 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('terms.section5.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section5.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('terms.section5.p2') }}
          </p>

          <!-- Section 6 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon
              color="primary"
              :icon="mdiShieldAccount"
              class="mr-1 mt-n1"
            />
            {{ $t('terms.section6.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('terms.section6.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-0">
            {{ $t('terms.section6.p2') }}
          </p>
        </VCard>
      </VCol>
    </VRow>
  </VContainer>
</template>
