<script setup lang="ts">
import { useDisplay } from 'vuetify'

const { lgAndUp, mdAndDown } = useDisplay()
</script>

<template>
  <v-app-bar elevation="0" color="#121212" height="70">
    <v-container>
      <v-row>
        <v-col class="w-full d-flex">
          <NuxtLink
            :to="{ name: 'index' }"
            class="text-on-surface text-headline-large text-decoration-none"
            :class="{ 'mt-5': mdAndDown, 'mt-4': !mdAndDown }"
          >
            <v-avatar v-if="lgAndUp" :image="'/img/logo.svg'" :size="30" />
            HTTP SMS
          </NuxtLink>
          <v-spacer />
          <div class="d-flex align-center mt-3 mb-3">
            <LanguageSwitcher class="mr-2" />
            <v-btn
              color="primary"
              :size="lgAndUp ? 'large' : 'default'"
              :to="{ name: 'login' }"
            >
              {{ $t('nav.getStarted') }}
              <span v-if="lgAndUp">&nbsp;{{ $t('nav.forFree') }}</span>
            </v-btn>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </v-app-bar>
</template>
