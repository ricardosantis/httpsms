<script setup lang="ts">
import { useDisplay } from 'vuetify'

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()

useHead({
  title: computed(() => `${t('threads.title')} - httpSMS`),
})

const { lgAndUp } = useDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const threadsStore = useThreadsStore()

onMounted(async () => {
  try {
    await authStore.loadUser().catch(() => {})
    await phonesStore.loadPhones().catch(() => {})
    await threadsStore.loadThreads()
  } catch {
    // ignore
  }
})
</script>

<template>
  <VContainer class="pa-0" fluid :class="{ 'fill-height': lgAndUp }">
    <VRow
      v-if="lgAndUp"
      align="center"
      :class="{ 'fill-height': lgAndUp }"
      justify="center"
    >
      <div :class="{ 'mt-n16': lgAndUp }">
        <VImg
          class="mx-auto mb-4"
          :class="{ 'mt-n16': lgAndUp }"
          max-height="400"
          max-width="90%"
          :src="'/img/person-texting.svg'"
        />
        <div class="text-center">
          <h3 class="text-headline-medium mt-4 mb-0">
            {{ $t('threads.selectMessage') }}
          </h3>
        </div>
      </div>
    </VRow>
    <VRow v-else justify="end">
      <VCol class="px-0 py-0">
        <MessageThreadHeader />
        <MessageThread />
      </VCol>
    </VRow>
  </VContainer>
</template>
