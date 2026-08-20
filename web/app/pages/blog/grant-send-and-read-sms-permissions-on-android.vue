<script setup lang="ts">
import { mdiDotsVertical } from '@mdi/js'
import { useDisplay } from 'vuetify'

const { mdAndUp } = useDisplay()
const { t, locale } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.grantPermissions.seoTitle')),
  description: computed(() =>
    t('blog.articles.grantPermissions.seoDescription'),
  ),
  ogTitle: computed(() => t('blog.articles.grantPermissions.title')),
  ogDescription: computed(() =>
    t('blog.articles.grantPermissions.description'),
  ),
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
          {{ $t('blog.articles.grantPermissions.title') }}
        </h1>
        <BlogInfo
          :date="
            locale === 'pt-BR' ? '18 de Fevereiro de 2025' : 'February 18, 2025'
          "
          :read-time="$t('blog.index.readTime', { time: '5 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.grantPermissions.intro') }}
          <a
            class="text-decoration-none"
            href="https://developer.android.com/training/permissions/requesting"
            target="_blank"
          >
            {{ $t('blog.articles.grantPermissions.runtimeLink') }} </a
          >.
        </p>

        <VAlert type="warning" variant="tonal" class="my-4">
          {{ $t('blog.articles.grantPermissions.alert') }}
        </VAlert>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.articles.grantPermissions.step1Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.grantPermissions.step1Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/grant-send-and-read-sms-permissions-on-android/app-info.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
          {{ $t('blog.articles.grantPermissions.step2Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.grantPermissions.step2Desc') }}
          <VIcon :icon="mdiDotsVertical" />
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/grant-send-and-read-sms-permissions-on-android/allow-restricted-settings.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
          {{ $t('blog.articles.grantPermissions.step3Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.grantPermissions.step3Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/grant-send-and-read-sms-permissions-on-android/allow.png"
        />

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p>
          {{ $t('blog.articles.grantPermissions.conclusionDesc') }}
          {{ $t('blog.common.contactPrompt') }}
        </p>

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
