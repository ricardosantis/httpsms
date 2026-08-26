<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiShieldLockOutline,
  mdiShieldCheck,
  mdiDatabaseLock,
  mdiEmailOutline,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const { t } = useI18n()

definePageMeta({
  layout: 'website',
})

useSeoMeta({
  title: computed(() => `${t('privacy.title')} - httpSMS`),
  description: computed(() => t('privacy.subtitle')),
  ogTitle: computed(() => `${t('privacy.title')} | httpSMS`),
  ogDescription: computed(() => t('privacy.subtitle')),
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
            <VChip color="success" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiShieldLockOutline" size="small" />
              {{ $t('website.legal') }} (LGPD)
            </VChip>
            <span class="text-caption text-medium-emphasis">
              {{ $t('privacy.lastUpdated') }}
            </span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('privacy.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('privacy.subtitle') }}
          </p>

          <VDivider class="mb-6" />

          <!-- Section 1 -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('privacy.section1.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('privacy.section1.p1') }}
          </p>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('privacy.section1.p2') }}
          </p>

          <!-- Section 2: Data Collection -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('privacy.section2.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('privacy.section2.p1') }}
          </p>
          <ul class="text-body-large text-medium-emphasis pl-6 mb-6">
            <li class="mb-2">{{ $t('privacy.section2.item1') }}</li>
            <li class="mb-2">{{ $t('privacy.section2.item2') }}</li>
            <li class="mb-2">{{ $t('privacy.section2.item3') }}</li>
            <li class="mb-2">{{ $t('privacy.section2.item4') }}</li>
            <li class="mb-2">{{ $t('privacy.section2.item5') }}</li>
          </ul>

          <!-- Section 3: Legal Bases -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('privacy.section3.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('privacy.section3.p1') }}
          </p>
          <ul class="text-body-large text-medium-emphasis pl-6 mb-6">
            <li class="mb-2">{{ $t('privacy.section3.item1') }}</li>
            <li class="mb-2">{{ $t('privacy.section3.item2') }}</li>
            <li class="mb-2">{{ $t('privacy.section3.item3') }}</li>
          </ul>

          <!-- Section 4: End to End Encryption -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="success" :icon="mdiDatabaseLock" class="mr-1 mt-n1" />
            {{ $t('privacy.section4.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('privacy.section4.p1') }}
          </p>

          <!-- Section 5: Third Party Sharing -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            {{ $t('privacy.section5.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('privacy.section5.p1') }}
          </p>
          <ul class="text-body-large text-medium-emphasis pl-6 mb-6">
            <li class="mb-2">{{ $t('privacy.section5.item1') }}</li>
            <li class="mb-2">{{ $t('privacy.section5.item2') }}</li>
            <li class="mb-2">{{ $t('privacy.section5.item3') }}</li>
          </ul>

          <!-- Section 6: Data Subject Rights -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="primary" :icon="mdiShieldCheck" class="mr-1 mt-n1" />
            {{ $t('privacy.section6.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-3">
            {{ $t('privacy.section6.p1') }}
          </p>
          <ul class="text-body-large text-medium-emphasis pl-6 mb-6">
            <li class="mb-2">{{ $t('privacy.section6.item1') }}</li>
            <li class="mb-2">{{ $t('privacy.section6.item2') }}</li>
            <li class="mb-2">{{ $t('privacy.section6.item3') }}</li>
            <li class="mb-2">{{ $t('privacy.section6.item4') }}</li>
            <li class="mb-2">{{ $t('privacy.section6.item5') }}</li>
          </ul>

          <!-- Section 7: DPO Contact -->
          <h2 class="text-headline-small font-weight-bold mb-3">
            <VIcon color="primary" :icon="mdiEmailOutline" class="mr-1 mt-n1" />
            {{ $t('privacy.section7.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-2">
            {{ $t('privacy.section7.p1') }}
          </p>
          <VBtn
            variant="tonal"
            color="primary"
            href="mailto:contato@smsandroid.com.br"
          >
            <VIcon start :icon="mdiEmailOutline" />
            contato@smsandroid.com.br
          </VBtn>
        </VCard>
      </VCol>
    </VRow>
  </VContainer>
</template>
