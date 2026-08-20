<script setup lang="ts">
import { useDisplay } from 'vuetify'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t, locale } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.forwardWebhook.seoTitle')),
  description: computed(() => t('blog.articles.forwardWebhook.seoDescription')),
  ogTitle: computed(() => t('blog.articles.forwardWebhook.title')),
  ogDescription: computed(() => t('blog.articles.forwardWebhook.description')),
  ogImage: '/img/blog/forward-incoming-sms-from-phone-to-webhook/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer class="pt-8">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <VImg
          style="border-radius: 4px"
          alt="blog post header image"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/header.png"
        />

        <h1
          :class="
            mdAndUp ? 'text-display-medium mt-1' : 'text-display-small mt-n2'
          "
        >
          {{ $t('blog.articles.forwardWebhook.title') }}
        </h1>
        <BlogInfo
          :date="locale === 'pt-BR' ? '08 de Abril de 2023' : 'April 08, 2023'"
          :read-time="$t('blog.index.readTime', { time: '5 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.forwardWebhook.intro') }}
        </p>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.common.step1ApiKeyTitle') }}
        </h3>
        <p>
          {{ $t('blog.articles.forwardWebhook.step1Desc') }}
          <NuxtLink class="text-decoration-none" to="/settings">
            {{ appStore.appData.url }}/settings
          </NuxtLink>
        </p>
        <VImg
          style="border-radius: 4px"
          alt="settings page"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/settings.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.step2AppTitle') }}
        </h3>
        <p>
          <a
            class="text-decoration-none"
            :href="appStore.appData.appDownloadUrl"
            download
          >
            {{ $t('blog.common.downloadAndInstall') }}
          </a>
          {{ $t('blog.articles.forwardWebhook.step2Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/android-app.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.articles.forwardWebhook.step3Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.forwardWebhook.step3Desc') }}
          <NuxtLink class="text-decoration-none" to="/settings">
            {{ appStore.appData.url }}/settings
          </NuxtLink>
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpSMS webhook configuration"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/webhook.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p>
          {{ $t('blog.articles.forwardWebhook.conclusionDesc') }}
        </p>
        <p>
          {{ $t('blog.articles.forwardWebhook.docPrompt') }}
          <a
            class="text-decoration-none"
            :href="appStore.appData.apiBaseUrl + '/docs/'"
            target="_blank"
          >
            {{ appStore.appData.apiBaseUrl }}/docs/
          </a>
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
