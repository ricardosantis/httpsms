<script setup lang="ts">
import { mdiBookOpenVariant } from '@mdi/js'
</script>

<template>
  <div>
    <NuxtLink to="/" class="text-decoration-none d-flex align-center">
      <VAvatar color="#121212" class="pa-1" size="36" rounded="0">
        <VImg height="33" width="33" :src="'/img/logo.svg'" />
      </VAvatar>
      <h3 class="text-headline-large ml-1 mb-3 mt-3 text-white">httpSMS</h3>
    </NuxtLink>
    <p class="mt-0">
      {{ $t('blog.sidebar.desc') }}
    </p>
    <VBtn variant="tonal" to="/docs">
      <VIcon start :icon="mdiBookOpenVariant" />
      {{ $t('blog.sidebar.documentation') }}
    </VBtn>
  </div>
</template>
