<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { mdiCommentTextMultipleOutline } from '@mdi/js'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t, locale } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.excel.seoTitle')),
  description: computed(() => t('blog.articles.excel.seoDescription')),
  ogTitle: computed(() => t('blog.articles.excel.title')),
  ogDescription: computed(() => t('blog.articles.excel.description')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer class="pt-8">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <h1
          :class="
            mdAndUp ? 'text-display-medium mt-1' : 'text-display-small mt-n2'
          "
        >
          {{ $t('blog.articles.excel.title') }}
        </h1>
        <BlogInfo
          :date="
            locale === 'pt-BR' ? '29 de Outubro de 2023' : 'October 29, 2023'
          "
          :read-time="$t('blog.index.readTime', { time: '5 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.excel.intro') }}
        </p>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.common.prerequisites') }}
        </h3>
        <ul>
          <li>{{ $t('blog.articles.excel.prereqExcel') }}</li>
          <li>{{ $t('blog.articles.excel.prereqAndroid') }}</li>
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
        <VAlert type="info" variant="tonal">
          {{ $t('blog.common.step2PhoneFormatAlert') }}
        </VAlert>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/forward-incoming-sms-from-phone-to-webhook/android-app.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.articles.excel.step3Title') }}
        </h3>
        <p>
          <a
            class="text-decoration-none"
            href="/templates/httpsms-bulk.xlsx"
            download
          >
            {{ $t('blog.articles.excel.excelTemplateLink') }}
          </a>
          {{ $t('blog.articles.excel.step3Desc') }}
        </p>
        <VAlert type="info" variant="tonal" class="mt-2 mb-4">
          {{ $t('blog.articles.excel.step3Alert') }}
        </VAlert>
        <VImg
          style="border-radius: 4px"
          alt="httpsms spreadsheet"
          height="800"
          src="/img/blog/send-bulk-sms-from-csv-file-with-no-code/httpms-spreedsheet.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.articles.excel.step4Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.excel.step4Desc') }}
          <NuxtLink class="text-decoration-none" to="/bulk-messages">
            <VIcon
              color="primary"
              size="small"
              class="mr-1"
              :icon="mdiCommentTextMultipleOutline"
            />
            {{ $t('blog.articles.excel.bulkMessagesLink') }} </NuxtLink
          >.
        </p>
        <VImg
          style="border-radius: 4px"
          alt="bulk csv upload"
          height="800"
          src="/img/blog/send-bulk-sms-from-csv-file-with-no-code/bulk-csv-upload.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p class="mt-4">
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
