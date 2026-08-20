<script setup lang="ts">
import { useDisplay } from 'vuetify'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t, locale } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.zapierSheets.seoTitle')),
  description: computed(() => t('blog.articles.zapierSheets.seoDescription')),
  ogTitle: computed(() => t('blog.articles.zapierSheets.title')),
  ogDescription: computed(() => t('blog.articles.zapierSheets.description')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer v-highlight class="pt-8">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <h1
          :class="
            mdAndUp ? 'text-display-medium mt-1' : 'text-display-small mt-n2'
          "
        >
          {{ $t('blog.articles.zapierSheets.title') }}
        </h1>
        <BlogInfo
          :date="
            locale === 'pt-BR' ? '29 de Novembro de 2023' : 'November 29, 2023'
          "
          :read-time="$t('blog.index.readTime', { time: '5 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.zapierSheets.intro') }}
        </p>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.common.prerequisites') }}
        </h3>
        <ul>
          <li>{{ $t('blog.articles.zapierSheets.prereqSheets') }}</li>
          <li>{{ $t('blog.articles.zapierSheets.prereqZapier') }}</li>
          <li>{{ $t('blog.articles.zapierSheets.prereqAccount') }}</li>
        </ul>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.articles.zapierSheets.step1Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.zapierSheets.step1Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="zapier trigger"
          src="/img/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier/zapier-trigger.png"
        />
        <p class="mt-8">
          {{ $t('blog.articles.zapierSheets.step1SpreadsheetDesc') }}
        </p>
        <VAlert type="info" variant="tonal">
          {{ $t('blog.articles.zapierSheets.step1Alert') }}
        </VAlert>
        <VImg
          style="border-radius: 4px"
          alt="google sheets"
          src="/img/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier/google-sheets.png"
        />

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.articles.zapierSheets.step2Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.zapierSheets.step2Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="zapier action event"
          src="/img/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier/zapier-action-event.png"
        />
        <p class="mt-8">
          {{ $t('blog.articles.zapierSheets.step2ActionDesc') }}
        </p>
        <pre
          class="pa-4 mb-6 rounded bg-surface-variant overflow-x-auto"
        ><code class="language-json text-body-medium">{
  "content": "Hello [Name]\nThanks for ordering [Product] via our store. Your order will be shipped today!",
  "from": "+18005550199",
  "to": "[ToPhoneNumber]"
}</code></pre>
        <VAlert type="info" variant="tonal" class="mt-4">
          {{ $t('blog.articles.zapierSheets.step2Alert') }}
        </VAlert>
        <p class="mt-8">
          {{ $t('blog.articles.zapierSheets.step2HeadersDesc') }}
          <NuxtLink class="text-decoration-none" to="/settings">
            {{ appStore.appData.url }}/settings </NuxtLink
          >.
        </p>
        <VImg
          style="border-radius: 4px"
          alt="zapier action action"
          src="/img/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier/zapier-action-action.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p>
          {{ $t('blog.articles.zapierSheets.conclusionDesc') }}
          {{ $t('blog.common.contactPrompt') }}
        </p>
        <p>{{ $t('blog.common.untilNextTime') }}</p>

        <BlogAuthorBio />
        <VDivider class="mx-16" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
      <VCol v-if="$vuetify.display.mdAndUp" md="3">
        <BlogSidebar />
      </VCol>
    </VRow>
  </VContainer>
</template>
