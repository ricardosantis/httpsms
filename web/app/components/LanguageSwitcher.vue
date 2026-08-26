<script setup lang="ts">
import { mdiTranslate } from '@mdi/js'

const { locale, locales, setLocale } = useI18n()
const authStore = useAuthStore()

const currentLocale = computed(() => locale.value)
const availableLocales = computed(() => locales.value)

async function switchLocale(code: string) {
  setLocale(code)
  if (authStore.user) {
    try {
      await authStore.updateUser({ locale: code })
    } catch {
      // ignore sync errors if offline or unauthorized
    }
  }
}
</script>

<template>
  <v-menu location="bottom end">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        variant="text"
        size="small"
        class="text-caption text-capitalize px-2"
      >
        <v-icon :icon="mdiTranslate" size="18" class="mr-1" />
        {{ currentLocale === 'pt-BR' ? 'PT-BR' : 'EN' }}
      </v-btn>
    </template>
    <v-list density="compact" bg-color="surface">
      <v-list-item
        v-for="loc in availableLocales"
        :key="typeof loc === 'string' ? loc : loc.code"
        :active="currentLocale === (typeof loc === 'string' ? loc : loc.code)"
        @click="switchLocale(typeof loc === 'string' ? loc : loc.code)"
      >
        <v-list-item-title class="text-caption">
          {{ typeof loc === 'string' ? loc : loc.name }}
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
