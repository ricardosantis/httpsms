<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: ['auth', 'admin'],
})

await navigateTo('/admin/users', { replace: true })
</script>

<template>
  <div />
</template>
