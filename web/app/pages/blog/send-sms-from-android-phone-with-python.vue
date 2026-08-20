<script setup lang="ts">
import { useDisplay } from 'vuetify'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t, locale } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.python.seoTitle')),
  description: computed(() => t('blog.articles.python.seoDescription')),
  ogTitle: computed(() => t('blog.articles.python.title')),
  ogDescription: computed(() => t('blog.articles.python.description')),
  ogImage: '/img/blog/send-sms-from-android-phone-with-python/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer v-highlight class="pt-8">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <VImg
          style="border-radius: 4px"
          alt="blog post header image"
          src="/img/blog/send-sms-from-android-phone-with-python/header.png"
        />

        <h1
          :class="
            mdAndUp ? 'text-display-medium mt-1' : 'text-display-small mt-n2'
          "
        >
          {{ $t('blog.articles.python.title') }}
        </h1>
        <BlogInfo
          :date="locale === 'pt-BR' ? '03 de Junho de 2023' : 'June 03, 2023'"
          :read-time="$t('blog.index.readTime', { time: '6 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.python.intro1') }}
        </p>
        <p>
          {{ $t('blog.articles.python.intro2') }}
        </p>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.common.prerequisites') }}
        </h3>
        <ul>
          <li>{{ $t('blog.articles.python.prereqPython') }}</li>
          <li>{{ $t('blog.articles.excel.prereqAndroid') }}</li>
          <li>
            <a
              class="text-decoration-none"
              href="https://www.python.org/"
              target="_blank"
            >
              Python
            </a>
            {{ $t('blog.articles.python.prereqPythonInstalled') }}
          </li>
        </ul>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.common.step1ApiKeyTitle') }}
        </h3>
        <p>
          {{ $t('blog.common.step1ApiKeyDesc') }}
          <NuxtLink class="text-decoration-none" to="/settings">
            {{ appStore.appData.url }}/settings
          </NuxtLink>
        </p>
        <VImg
          style="border-radius: 4px"
          alt="settings page"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/settings.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
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
          {{ $t('blog.common.step2AppDesc') }}
        </p>
        <VAlert type="info" variant="outlined">
          {{ $t('blog.common.step2PhoneFormatAlert') }}
        </VAlert>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/android-app.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.articles.python.step3Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.python.step3Desc') }}
        </p>
        <VAlert type="info" variant="outlined" class="mt-2 mb-4">
          {{ $t('blog.articles.python.step3Alert') }}
        </VAlert>
        <pre
          class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
        ><code class="language-python text-body-medium">import requests
import json

api_key = "" # Get API Key from /settings

url = '{{ appStore.appData.apiBaseUrl }}/v1/messages/send'

headers = {
    'x-api-key': api_key,
    'Accept': 'application/json',
    'Content-Type': 'application/json'
}

payload = {
    "content": "This is a sample text message sent via python",
    "from": "+18005550199", # This is the phone number of your android phone
    "to": "+18005550100" # This is the recipient phone number
}

response = requests.post(url, headers=headers, data=json.dumps(payload))

print(json.dumps(response.json(), indent=4))</code></pre>
        <p>
          {{ $t('blog.articles.python.step3RunDesc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="sms sent"
          height="800"
          src="/img/blog/send-sms-from-android-phone-with-python/sms-sent.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p>
          {{ $t('blog.articles.python.conclusionDesc') }}
        </p>
        <p>
          {{ $t('blog.articles.python.forwardPrompt') }}
          <NuxtLink
            class="text-decoration-none"
            to="/blog/forward-incoming-sms-from-phone-to-webhook"
          >
            {{ $t('blog.articles.python.forwardLink') }}
          </NuxtLink>
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
