<script setup lang="ts">
const { mdAndUp, smAndDown } = useDisplay()
const { t, locale } = useI18n()

definePageMeta({
  layout: 'website',
})

useSeoMeta({
  title: computed(() => t('blog.index.seo.title')),
  description: computed(() => t('blog.index.seo.description')),
  ogTitle: computed(() => t('blog.index.seo.ogTitle')),
  ogDescription: computed(() => t('blog.index.seo.ogDescription')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})

const articleDefs = [
  {
    slug: 'grant-send-and-read-sms-permissions-on-android',
    key: 'grantPermissions',
    datePt: '18 de Fevereiro de 2025',
    dateEn: 'February 18, 2025',
    readTimeMin: '5',
    sortDate: '2025-02-18',
  },
  {
    slug: 'end-to-end-encryption-to-sms-messages',
    key: 'e2eEncryption',
    datePt: '21 de Janeiro de 2024',
    dateEn: 'January 21, 2024',
    readTimeMin: '10',
    sortDate: '2024-01-21',
  },
  {
    slug: 'send-sms-when-new-row-is-added-to-google-sheets-using-zapier',
    key: 'zapierSheets',
    datePt: '29 de Novembro de 2023',
    dateEn: 'November 29, 2023',
    readTimeMin: '7',
    sortDate: '2023-11-29',
  },
  {
    slug: 'how-to-send-sms-messages-from-excel',
    key: 'excel',
    datePt: '29 de Outubro de 2023',
    dateEn: 'October 29, 2023',
    readTimeMin: '5',
    sortDate: '2023-10-29',
  },
  {
    slug: 'send-bulk-sms-from-csv-file-with-no-code',
    key: 'bulkCsv',
    datePt: '29 de Outubro de 2023',
    dateEn: 'October 29, 2023',
    readTimeMin: '7',
    sortDate: '2023-10-29',
  },
  {
    slug: 'send-sms-from-android-phone-with-python',
    key: 'python',
    datePt: '03 de Junho de 2023',
    dateEn: 'June 03, 2023',
    readTimeMin: '6',
    sortDate: '2023-06-03',
  },
  {
    slug: 'forward-incoming-sms-from-phone-to-webhook',
    key: 'forwardWebhook',
    datePt: '08 de Abril de 2023',
    dateEn: 'April 08, 2023',
    readTimeMin: '5',
    sortDate: '2023-04-08',
  },
]

const sortedArticles = computed(() =>
  [...articleDefs].sort((a, b) => b.sortDate.localeCompare(a.sortDate)),
)
</script>

<template>
  <VContainer>
    <VRow>
      <VCol cols="12" md="9">
        <VRow>
          <VCol cols="12">
            <h1 class="text-display-large mb-2" :class="{ 'mt-0': smAndDown }">
              {{ $t('blog.index.title') }}
            </h1>
            <h2
              class="text-medium-emphasis mt-2 mb-n4 text-title-large font-weight-light"
            >
              {{ $t('blog.index.subtitle') }}
            </h2>
          </VCol>
        </VRow>
        <VRow>
          <VCol
            v-for="article in sortedArticles"
            :key="article.slug"
            cols="12"
            xl="6"
          >
            <NuxtLink
              :to="`/blog/${article.slug}`"
              class="text-decoration-none"
            >
              <VHover v-slot="{ isHovering, props: hoverProps }">
                <VCard
                  v-bind="hoverProps"
                  :elevation="isHovering ? 8 : 2"
                  :color="isHovering ? 'blue-darken-4' : undefined"
                  class="blog-card"
                >
                  <VCardTitle class="text-headline-large text-wrap title-clamp">
                    {{ $t(`blog.articles.${article.key}.title`) }}
                  </VCardTitle>
                  <VCardSubtitle>
                    <span class="text-uppercase text-blue">
                      {{ locale === 'pt-BR' ? article.datePt : article.dateEn }}
                    </span>
                    •
                    <span class="text-uppercase">
                      {{
                        $t('blog.index.readTime', {
                          time: `${article.readTimeMin} min`,
                        })
                      }}
                    </span>
                  </VCardSubtitle>
                  <VCardText class="mt-n2">
                    <p class="text-title-medium description-clamp">
                      {{ $t(`blog.articles.${article.key}.description`) }}
                    </p>
                  </VCardText>
                </VCard>
              </VHover>
            </NuxtLink>
          </VCol>
        </VRow>
      </VCol>
      <VCol v-if="mdAndUp" md="3" class="pt-6">
        <BlogSidebar />
      </VCol>
    </VRow>
  </VContainer>
</template>

<style scoped>
.blog-card {
  transition: all 0.3s ease;
}

.title-clamp {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (width >= 1280px) {
  .title-clamp {
    -webkit-line-clamp: 2;
    height: calc(2.7rem * 2);
  }
}

.description-clamp {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (width >= 1280px) {
  .description-clamp {
    -webkit-line-clamp: 3;
    height: calc(1.5rem * 3);
  }
}

@media (width >= 2560px) {
  .description-clamp {
    -webkit-line-clamp: 2;
    height: calc(1.5rem * 2);
  }
}
</style>
