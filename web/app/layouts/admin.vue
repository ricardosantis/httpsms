<script setup lang="ts">
import { getAuth, signOut } from 'firebase/auth'
import { mdiAccountGroup, mdiLogout, mdiWeb } from '@mdi/js'

const authStore = useAuthStore()
const router = useRouter()
const notificationsStore = useNotificationsStore()
const { t } = useI18n()

async function logout() {
  const auth = getAuth()
  await signOut(auth)
  authStore.resetState()
  notificationsStore.addNotification({
    type: 'info',
    message: t('auth.logoutSuccess'),
  })
  router.push({ name: 'index' })
}
</script>

<template>
  <v-app>
    <v-app-bar color="#121212" elevation="0">
      <v-container>
        <v-row>
          <v-col class="w-full d-flex align-center">
            <NuxtLink
              :to="{ name: 'admin-users' }"
              class="text-decoration-none d-flex align-baseline"
            >
              <v-avatar
                color="#121212"
                class="pa-1"
                rounded="0"
                :image="'/img/logo.svg'"
                :size="38"
              />
              <h3 class="text-headline-large mb-0 ml-1 mt-6 text-white">
                {{ t('admin.title') }}
              </h3>
            </NuxtLink>
            <v-spacer />
            <v-btn
              to="/"
              color="grey-lighten-1"
              variant="text"
              size="large"
              class="mr-2"
            >
              <v-icon start :icon="mdiWeb" />
              {{ t('admin.viewSite') }}
            </v-btn>
            <v-btn
              :to="{ name: 'admin-users' }"
              color="primary"
              variant="flat"
              size="large"
            >
              <v-icon start :icon="mdiAccountGroup" />
              {{ t('admin.users') }}
            </v-btn>
            <v-btn color="primary" variant="text" size="large" @click="logout">
              <v-icon start :icon="mdiLogout" />
              {{ t('admin.logout') }}
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-app-bar>
    <v-main>
      <AppToast />
      <slot v-if="authStore.authStateChanged" />
      <LoadingDashboard v-else />
    </v-main>
  </v-app>
</template>
