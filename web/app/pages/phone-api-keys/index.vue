<script setup lang="ts">
import { mdiArrowLeft, mdiPlus, mdiDelete, mdiEye } from '@mdi/js'
import QRCode from 'qrcode'
import Pusher from 'pusher-js'
import type { Channel } from 'pusher-js'
import { ErrorMessages } from '~/utils/errors'
import { toApiError } from '~/utils/api-error'
import type { EntitiesPhoneAPIKey } from '~~/shared/types/api'

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()

useHead({
  title: computed(() => `${t('phoneApiKeys.title')} - httpSMS`),
})

const config = useRuntimeConfig()
const { lgAndUp } = useDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const notificationsStore = useNotificationsStore()
const { formatTimestamp, formatPhoneNumber } = useFilters()
const { useApi } = useApiComposable()

const loading = ref(true)
const phoneApiKeys = ref<EntitiesPhoneAPIKey[]>([])
const errorMessages = ref(new ErrorMessages())

const showCreateApiKeyDialog = ref(false)
const formPhoneApiKeyName = ref('')

const showPhoneApiKeyQrCode = ref(false)
const deleteApiKeyDialog = ref(false)
const removePhoneFromApiKeyDialog = ref(false)
const activePhoneApiKey = ref<EntitiesPhoneAPIKey | null>(null)
const activePhoneNumber = ref('')

const qrCodeCanvas = ref<HTMLCanvasElement | null>(null)
let webhookChannel: Channel | null = null

function parseErrors(error: unknown): ErrorMessages {
  const bag = new ErrorMessages()
  const data = toApiError(error).data?.data
  if (data && typeof data === 'object') {
    Object.keys(data).forEach((key) => bag.addMany(key, data[key]))
  }
  return bag
}

async function loadPhoneApiKeys() {
  loading.value = true
  try {
    const api = useApi()
    const response = await api<{ data: EntitiesPhoneAPIKey[] }>(
      '/v1/phone-api-keys',
      { query: { limit: 100 } },
    )
    phoneApiKeys.value = response.data ?? []
  } catch {
    notificationsStore.addNotification({
      message: t('phoneApiKeys.loadFailed'),
      type: 'error',
    })
  } finally {
    loading.value = false
  }
}

async function createPhoneApiKey() {
  errorMessages.value = new ErrorMessages()
  loading.value = true
  try {
    const api = useApi()
    await api('/v1/phone-api-keys', {
      method: 'POST',
      body: { name: formPhoneApiKeyName.value },
    })
    notificationsStore.addNotification({
      message: t('phoneApiKeys.createSuccess'),
      type: 'success',
    })
    formPhoneApiKeyName.value = ''
    showCreateApiKeyDialog.value = false
    await loadPhoneApiKeys()
  } catch (error: unknown) {
    errorMessages.value = parseErrors(error)
    if (errorMessages.value.size() === 0) {
      notificationsStore.addNotification({
        message: t('phoneApiKeys.createFailed'),
        type: 'error',
      })
    }
  } finally {
    loading.value = false
  }
}

function generateQrCode(text: string) {
  const canvas = qrCodeCanvas.value
  if (!canvas) {
    return
  }
  QRCode.toCanvas(
    canvas,
    text,
    { errorCorrectionLevel: 'H' },
    (err: Error | null | undefined) => {
      if (err) {
        notificationsStore.addNotification({
          message: t('phoneApiKeys.qrFailed'),
          type: 'error',
        })
      }
    },
  )
}

function showPhoneApiKey(apiKey: EntitiesPhoneAPIKey) {
  activePhoneApiKey.value = apiKey
  showPhoneApiKeyQrCode.value = true
  nextTick(() => {
    generateQrCode(apiKey.api_key)
  })
}

function showDeletePhoneApiKeyDialog(apiKey: EntitiesPhoneAPIKey) {
  activePhoneApiKey.value = apiKey
  deleteApiKeyDialog.value = true
}

function showRemovePhoneFromApiKeyDialog(
  apiKey: EntitiesPhoneAPIKey,
  phoneNumber: string,
) {
  activePhoneApiKey.value = apiKey
  activePhoneNumber.value = phoneNumber
  removePhoneFromApiKeyDialog.value = true
}

async function deleteApiKey() {
  if (!activePhoneApiKey.value) {
    return
  }
  loading.value = true
  try {
    const api = useApi()
    await api(`/v1/phone-api-keys/${activePhoneApiKey.value.id}`, {
      method: 'DELETE',
    })
    notificationsStore.addNotification({
      message: t('phoneApiKeys.deleteSuccess'),
      type: 'success',
    })
    deleteApiKeyDialog.value = false
    await loadPhoneApiKeys()
  } catch {
    notificationsStore.addNotification({
      message: t('phoneApiKeys.deleteFailed'),
      type: 'error',
    })
    loading.value = false
  }
}

async function removePhoneFromPhoneKey() {
  if (!activePhoneApiKey.value) {
    return
  }
  const phoneId = phonesStore.phones.find(
    (phone) => phone.phone_number === activePhoneNumber.value,
  )?.id
  if (!phoneId) {
    notificationsStore.addNotification({
      message: t('phoneApiKeys.phoneNotFound'),
      type: 'error',
    })
    return
  }
  loading.value = true
  try {
    const api = useApi()
    await api(
      `/v1/phone-api-keys/${activePhoneApiKey.value.id}/phones/${phoneId}`,
      { method: 'DELETE' },
    )
    notificationsStore.addNotification({
      message: t('phoneApiKeys.removePhoneSuccess'),
      type: 'success',
    })
    removePhoneFromApiKeyDialog.value = false
    await loadPhoneApiKeys()
  } catch {
    notificationsStore.addNotification({
      message: t('phoneApiKeys.removePhoneFailed'),
      type: 'error',
    })
    loading.value = false
  }
}

onMounted(async () => {
  await authStore.loadUser()
  await phonesStore.loadPhones()
  await loadPhoneApiKeys()

  const pusherKey = config.public.pusherKey as string
  const pusherCluster = config.public.pusherCluster as string
  if (pusherKey && pusherCluster && authStore.user?.id) {
    try {
      const pusher = new Pusher(pusherKey, { cluster: pusherCluster })
      webhookChannel = pusher.subscribe(authStore.user.id)
      webhookChannel.bind('phone.updated', () => {
        if (!loading.value) {
          loadPhoneApiKeys()
        }
      })
    } catch {
      // Pusher failed to initialize
    }
  }
})

onBeforeUnmount(() => {
  if (webhookChannel) {
    webhookChannel.unsubscribe()
  }
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0" :class="{ 'fill-height': lgAndUp }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>{{ $t('phoneApiKeys.title') }}</VToolbarTitle>
        <VSpacer />
        <LanguageSwitcher class="mr-2" />
        <VProgressLinear
          color="primary"
          :active="loading"
          :indeterminate="loading"
          absolute
          location="bottom"
        />
      </VAppBar>
      <VContainer class="pt-0">
        <VRow>
          <VCol cols="12" md="9" offset-md="1" xl="8" offset-xl="2">
            <div class="d-flex align-center flex-wrap mt-3 mb-4">
              <VProgressCircular
                v-if="loading"
                :size="24"
                :width="2"
                color="primary"
                class="mt-1 mr-2"
                indeterminate
              />
              <h5 class="text-md-display-small text-title-large my-0">
                {{ $t('phoneApiKeys.title') }}
              </h5>
              <VBtn
                color="primary"
                class="ml-4 mt-1"
                @click="showCreateApiKeyDialog = true"
              >
                <VIcon start :icon="mdiPlus" />
                {{ $t('phoneApiKeys.createApiKey') }}
              </VBtn>
              <VSpacer />
              <VBtn
                v-if="lgAndUp"
                :href="`${appStore.appData.documentationUrl}/features/phone-api-keys`"
                target="_blank"
                variant="tonal"
                class="mt-1"
              >
                {{ $t('common.documentation') }}
              </VBtn>
            </div>
            <p class="text-medium-emphasis">
              {{
                $t('phoneApiKeys.description', {
                  url: `${appStore.appData.url}/settings`,
                })
              }}
            </p>
            <VTable class="mb-4 api-key-table" density="comfortable">
              <thead>
                <tr class="text-uppercase text-medium-emphasis">
                  <th class="text-left">{{ $t('common.name') }}</th>
                  <th class="text-left">{{ $t('common.createdAt') }}</th>
                  <th class="text-left">{{ $t('nav.phoneNumbers') }}</th>
                  <th class="text-left">{{ $t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="phoneApiKey in phoneApiKeys" :key="phoneApiKey.id">
                  <td class="text-left">{{ phoneApiKey.name }}</td>
                  <td>{{ formatTimestamp(phoneApiKey.created_at) }}</td>
                  <td>
                    <ul v-if="phoneApiKey.phone_numbers.length" class="ml-n3">
                      <li
                        v-for="phoneNumber in phoneApiKey.phone_numbers"
                        :key="phoneNumber"
                        class="my-3"
                      >
                        <b>{{ formatPhoneNumber(phoneNumber) }}</b>
                        <VBtn
                          class="ml-2"
                          size="small"
                          color="error"
                          @click="
                            showRemovePhoneFromApiKeyDialog(
                              phoneApiKey,
                              phoneNumber,
                            )
                          "
                        >
                          {{ $t('common.remove') }}
                        </VBtn>
                      </li>
                    </ul>
                    <span v-else class="text-medium-emphasis">-</span>
                  </td>
                  <td>
                    <VBtn
                      size="small"
                      color="primary"
                      :disabled="loading"
                      @click="showPhoneApiKey(phoneApiKey)"
                    >
                      <VIcon start :icon="mdiEye" /> {{ $t('common.view') }}
                    </VBtn>
                    <VBtn
                      class="ml-2"
                      size="small"
                      color="error"
                      :disabled="loading"
                      @click="showDeletePhoneApiKeyDialog(phoneApiKey)"
                    >
                      <VIcon start :icon="mdiDelete" />
                      {{ $t('common.delete') }}
                    </VBtn>
                  </td>
                </tr>
              </tbody>
            </VTable>
          </VCol>
        </VRow>
      </VContainer>
    </div>

    <VDialog v-model="showCreateApiKeyDialog" max-width="600" opacity="0.9">
      <VCard>
        <VCardTitle>{{ $t('phoneApiKeys.createTitle') }}</VCardTitle>
        <VCardSubtitle class="mt-2" style="white-space: normal">
          {{ $t('phoneApiKeys.createSubtitle') }}
        </VCardSubtitle>
        <VCardText>
          <VForm @submit.prevent="createPhoneApiKey">
            <VTextField
              v-model="formPhoneApiKeyName"
              variant="outlined"
              :label="$t('common.name')"
              class="mt-4"
              persistent-placeholder
              :placeholder="$t('phoneApiKeys.namePlaceholder')"
              name="api-key"
              :disabled="loading"
              :error="errorMessages.has('name')"
              :error-messages="errorMessages.get('name')"
            />
          </VForm>
        </VCardText>
        <VCardActions class="mt-n6 mb-1">
          <loading-button
            color="primary"
            :loading="loading"
            @click="createPhoneApiKey"
          >
            {{ $t('common.create') }}
          </loading-button>
          <VSpacer />
          <VBtn
            variant="text"
            color="warning"
            @click="showCreateApiKeyDialog = false"
          >
            {{ $t('common.close') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

    <VDialog v-model="showPhoneApiKeyQrCode" max-width="600" opacity="0.9">
      <VCard>
        <VCardTitle>{{ $t('phoneApiKeys.qrTitle') }}</VCardTitle>
        <VCardSubtitle class="mt-2" style="white-space: normal">
          {{ $t('phoneApiKeys.qrSubtitle') }}
        </VCardSubtitle>
        <VCardText class="text-center">
          <VTextField
            :model-value="activePhoneApiKey?.api_key"
            readonly
            name="api-key"
            variant="outlined"
          />
          <canvas ref="qrCodeCanvas"></canvas>
        </VCardText>
        <VCardActions>
          <CopyButton
            :value="activePhoneApiKey?.api_key ?? ''"
            color="primary"
            :copy-text="$t('phoneApiKeys.copyApiKey')"
            :notification-text="$t('phoneApiKeys.copySuccess')"
          />
          <VSpacer />
          <VBtn
            color="warning"
            variant="text"
            @click="showPhoneApiKeyQrCode = false"
          >
            {{ $t('common.close') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

    <VDialog v-model="deleteApiKeyDialog" max-width="600" opacity="0.9">
      <VCard>
        <VCardTitle class="text-h5 text-break">
          {{ $t('phoneApiKeys.deleteConfirmTitle') }}
        </VCardTitle>
        <VCardText class="text-medium-emphasis">
          {{ $t('phoneApiKeys.deleteConfirmDesc') }}
        </VCardText>
        <VCardActions class="pb-2 mt-n2">
          <VBtn
            color="error"
            variant="flat"
            :loading="loading"
            @click="deleteApiKey"
          >
            <VIcon start :icon="mdiDelete" />
            {{ $t('common.delete') }}
          </VBtn>
          <VSpacer />
          <VBtn
            variant="text"
            color="warning"
            @click="deleteApiKeyDialog = false"
          >
            {{ $t('common.close') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

    <VDialog
      v-model="removePhoneFromApiKeyDialog"
      max-width="600"
      opacity="0.9"
    >
      <VCard>
        <VCardTitle class="text-h5 text-break">
          {{ $t('phoneApiKeys.removePhoneConfirmTitle') }}
        </VCardTitle>
        <VCardText>
          {{
            $t('phoneApiKeys.removePhoneConfirmDesc', {
              phone: formatPhoneNumber(activePhoneNumber),
            })
          }}
        </VCardText>
        <VCardActions class="pb-4">
          <VBtn
            color="error"
            :loading="loading"
            @click="removePhoneFromPhoneKey"
          >
            <VIcon start :icon="mdiDelete" />
            {{ $t('phoneApiKeys.removePhoneFromKey') }}
          </VBtn>
          <VSpacer />
          <VBtn
            variant="text"
            color="warning"
            @click="removePhoneFromApiKeyDialog = false"
          >
            {{ $t('common.close') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VContainer>
</template>

<style scoped lang="scss">
.api-key-table {
  tbody {
    tr:hover {
      background-color: transparent !important;
    }
  }
}
</style>
