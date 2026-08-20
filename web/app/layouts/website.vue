<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiGithub,
  mdiCircle,
  mdiTwitter,
  mdiHeart,
  mdiLightbulbOn50,
  mdiCreation,
  mdiEyeOffOutline,
  mdiPost,
  mdiCreditCardOutline,
  mdiScaleBalance,
  mdiEmailOutline,
  mdiBookOpenVariant,
} from '@mdi/js'

const router = useRouter()
const route = useRoute()
const { lgAndUp, mdAndUp } = useDisplay()
const authStore = useAuthStore()

function goToPricing() {
  if (route.name === 'index') {
    document.getElementById('pricing')?.scrollIntoView({ behavior: 'smooth' })
  } else {
    router.push('/').then(() => {
      setTimeout(() => {
        document
          .getElementById('pricing')
          ?.scrollIntoView({ behavior: 'smooth' })
      }, 300)
    })
  }
}
</script>

<template>
  <v-app>
    <v-app-bar color="#121212" elevation="0">
      <v-container>
        <v-row>
          <v-col class="w-full d-flex align-center">
            <NuxtLink to="/" class="text-decoration-none d-flex align-baseline">
              <v-avatar
                color="#121212"
                class="pa-1"
                rounded="0"
                :image="'/img/logo.svg'"
                :size="38"
              />
              <h3
                v-if="lgAndUp"
                class="text-headline-large mb-0 ml-1 mt-6 text-white"
              >
                httpSMS
              </h3>
            </NuxtLink>
            <v-spacer />
            <v-btn
              v-show="lgAndUp"
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              @click="goToPricing"
            >
              {{ $t('website.pricing') }}
            </v-btn>
            <v-btn
              v-show="lgAndUp"
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              to="/docs"
            >
              {{ $t('website.documentation') }}
            </v-btn>
            <v-btn
              v-show="lgAndUp"
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              :to="{ name: 'blog' }"
            >
              {{ $t('website.blog') }}
            </v-btn>
            <v-btn
              v-show="
                lgAndUp &&
                authStore.authStateChanged &&
                authStore.authUser === null
              "
              size="large"
              variant="text"
              color="primary"
              class="my-5 mr-2"
              :to="{ name: 'login' }"
            >
              {{ $t('website.login') }}
            </v-btn>
            <v-btn
              v-show="authStore.authStateChanged && authStore.authUser === null"
              color="primary"
              variant="flat"
              :class="{ 'mt-5': mdAndUp, 'mt-1': !mdAndUp }"
              :size="lgAndUp ? 'large' : 'default'"
              :to="{ name: 'login' }"
            >
              {{ $t('website.getStarted') }}
              <span v-show="lgAndUp">&nbsp;{{ $t('website.forFree') }}</span>
            </v-btn>
            <div
              v-show="authStore.authStateChanged && authStore.authUser !== null"
              class="position-relative d-inline-block"
            >
              <v-btn
                color="primary"
                variant="flat"
                :class="{ 'mt-5': mdAndUp, 'mt-1': !mdAndUp }"
                :size="lgAndUp ? 'large' : 'default'"
                :to="{ name: 'threads' }"
              >
                {{ $t('website.dashboard') }}
              </v-btn>
              <RedirectPromptPopover />
            </div>
            <LanguageSwitcher class="ml-4" />
          </v-col>
        </v-row>
      </v-container>
    </v-app-bar>
    <v-main>
      <AppToast />
      <slot />
    </v-main>
    <v-footer>
      <v-container>
        <v-row>
          <v-col cols="12" md="3">
            <NuxtLink to="/" class="text-decoration-none d-flex mt-n6">
              <v-avatar
                color="#212121"
                class="mt-8 pa-1"
                rounded="0"
                :image="'/img/logo.svg'"
                :size="38"
              />
              <h3 class="text-headline-large ml-1 mb-0 text-white">httpSMS</h3>
            </NuxtLink>
            <div class="text-title-medium mb-4 text-medium-emphasis">
              {{ $t('website.madeWith') }}
              <v-icon color="#cf1112" :icon="mdiHeart" />
              {{ $t('website.inSaoPaulo') }}
              <v-img
                class="d-inline-block"
                width="20"
                src="https://upload.wikimedia.org/wikipedia/commons/0/05/Flag_of_Brazil.svg"
              />
            </div>
            <p class="mt-n3">
              <v-btn
                href="#"
                color="#1DA1F2"
                class="ml-n3"
                variant="text"
                :icon="mdiTwitter"
              />
              <v-btn
                href="#"
                color="#ffffff"
                variant="text"
                :icon="mdiGithub"
              />
              <v-btn href="#" icon variant="text" color="#5865f2">
                <v-img
                  contain
                  height="24"
                  width="24"
                  src="/img/discord-logo-blue.svg"
                />
              </v-btn>
            </p>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">
              {{ $t('website.resources') }}
            </h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <a
                  class="text-white text-decoration-none footer-link"
                  style="cursor: pointer"
                  @click.stop="goToPricing"
                >
                  {{ $t('website.pricing') }}
                  <v-icon size="small" :icon="mdiCreditCardOutline" />
                </a>
              </li>
              <li class="mb-2">
                <a
                  href="https://status.httpsms.com"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.siteStatus') }}
                  <v-icon color="success" size="x-small" :icon="mdiCircle" />
                </a>
              </li>
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/blog"
                >
                  {{ $t('website.blog') }}
                  <v-icon size="small" :icon="mdiPost" />
                </NuxtLink>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">
              {{ $t('website.developers') }}
            </h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <NuxtLink
                  to="/docs"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.documentation') }}
                  <v-icon size="small" :icon="mdiBookOpenVariant" />
                </NuxtLink>
              </li>
              <li class="mb-2">
                <a href="#" class="text-white text-decoration-none footer-link">
                  Github <v-icon size="small" :icon="mdiGithub" />
                </a>
              </li>
              <li class="mb-2">
                <a
                  href="https://sandbox.httpsms.com"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.sandbox') }}
                  <v-icon size="small" color="pink" :icon="mdiCreation" />
                </a>
              </li>
              <li class="mb-2">
                <a
                  href="https://httpsms.featurebase.app"
                  class="text-white text-decoration-none footer-link"
                >
                  {{ $t('website.requestFeature') }}
                  <v-icon
                    size="small"
                    color="yellow"
                    :icon="mdiLightbulbOn50"
                  />
                </a>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="3">
            <h2 class="text-headline-small mb-2">{{ $t('website.legal') }}</h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/terms-and-conditions"
                >
                  {{ $t('website.termsAndConditions') }}
                  <v-icon size="small" :icon="mdiScaleBalance" />
                </NuxtLink>
              </li>
              <li class="mb-2">
                <NuxtLink
                  class="text-white text-decoration-none footer-link"
                  to="/privacy-policy"
                >
                  {{ $t('website.privacyPolicy') }}
                  <v-icon size="small" :icon="mdiEyeOffOutline" />
                </NuxtLink>
              </li>
              <li class="mt-2">
                <a
                  class="text-white text-decoration-none footer-link"
                  href="mailto:contato@mesaquevende.com.br"
                >
                  {{ $t('website.contactSupport') }}
                  <v-icon size="small" :icon="mdiEmailOutline" />
                </a>
              </li>
            </ul>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>
  </v-app>
</template>

<style scoped>
.footer-link:hover {
  text-decoration: underline !important;
}
</style>
