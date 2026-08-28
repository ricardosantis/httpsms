<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiAlert,
  mdiArrowLeft,
  mdiCallMade,
  mdiCallMissed,
  mdiCallReceived,
  mdiCheck,
  mdiCheckAll,
  mdiDelete,
  mdiExport,
  mdiMagnify,
  mdiProgressCheck,
  mdiRefresh,
} from '@mdi/js'
import type { EntitiesMessage, EntitiesPhone } from '~~/shared/types/api'
import type { SearchMessagesRequest } from '~~/shared/types/message'
import { ErrorMessages } from '~/utils/errors'
import { toApiError } from '~/utils/api-error'

interface Turnstile {
  ready(callback: () => void): void
  render(
    container: string | HTMLElement,
    params?: {
      sitekey: string
      action: string
      callback?: (token: string) => void
      'error-callback'?: (error: string) => void
    },
  ): string | null | undefined
  remove(widgetId?: string): void
}

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()

useHead({
  title: computed(() => `${t('searchMessages.title')} - httpSMS`),
})

const route = useRoute()
const config = useRuntimeConfig()
const { mdAndUp, lgAndUp } = useDisplay()
const messagesStore = useMessagesStore()
const phonesStore = usePhonesStore()
const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()
const { useApi } = useApiComposable()
const { formatPhoneNumber, formatTimestamp, capitalize } = useFilters()

const loading = ref(true)
const initialLoadComplete = ref(false)
const errorTitle = ref('')
const showDeleteDialog = ref(false)
const showResendDialog = ref(false)
const errorMessages = ref(new ErrorMessages())

const formOwners = ref<string[]>([])
const formTypes = ref<string[]>([])
const formStatuses = ref<string[]>([])
const formQuery = ref('')

const messages = ref<EntitiesMessage[]>([])
const totalMessages = ref(0)
const selectedIds = ref<string[]>([])
let turnstileWidgetId: string | null = null

const page = ref(1)
const itemsPerPage = ref(100)
const sortBy = ref<{ key: string; order: 'asc' | 'desc' }[]>([
  { key: 'created_at', order: 'desc' },
])

const itemsPerPageOptions = [
  { value: 10, title: '10' },
  { value: 50, title: '50' },
  { value: 100, title: '100' },
  { value: 200, title: '200' },
]

const pageText = computed(() => {
  if (totalMessages.value === 0) return `0-0 ${t('common.of')} 0`
  const start = (page.value - 1) * itemsPerPage.value + 1
  const end = Math.min(page.value * itemsPerPage.value, totalMessages.value)
  return `${start}-${end} ${t('common.of')} ${totalMessages.value}`
})

const headers = computed(() => [
  { title: t('searchMessages.createdAt'), key: 'created_at' },
  { title: t('searchMessages.owner'), key: 'owner' },
  { title: t('searchMessages.contact'), key: 'contact' },
  { title: t('searchMessages.type'), key: 'type' },
  { title: t('searchMessages.status'), key: 'status' },
  { title: t('searchMessages.content'), key: 'content', sortable: false },
])

const selectedMessages = computed<EntitiesMessage[]>(() =>
  messages.value.filter((message) => selectedIds.value.includes(message.id)),
)

const canResendSelected = computed<boolean>(
  () =>
    selectedMessages.value.length > 0 &&
    selectedMessages.value.every(
      (message) =>
        message.type === 'mobile-terminated' &&
        (message.status === 'expired' || message.status === 'failed'),
    ),
)

const phoneNumberSelectItems = computed(() =>
  phonesStore.phones.map((phone: EntitiesPhone) => ({
    title: formatPhoneNumber(phone.phone_number),
    value: phone.phone_number,
  })),
)

const messageTypeSelectItems = computed(() => [
  { title: t('searchMessages.outbound'), value: 'mobile-terminated' },
  { title: t('searchMessages.inbound'), value: 'mobile-originated' },
  { title: t('searchMessages.missedCalls'), value: 'call/missed' },
])

const messageStatusSelectItems = computed(() => [
  { value: 'pending', title: t('searchMessages.pending') },
  { value: 'sent', title: t('searchMessages.sent') },
  { value: 'delivered', title: t('searchMessages.delivered') },
  { value: 'failed', title: t('searchMessages.failed') },
  { value: 'expired', title: t('searchMessages.expired') },
  { value: 'received', title: t('searchMessages.received') },
])

function getCaptcha(): Promise<string> {
  const siteKey = (config.public as Record<string, string>)
    .cloudflareTurnstileSiteKey
  if (!siteKey) {
    return Promise.resolve('')
  }
  return new Promise<string>((resolve) => {
    const turnstile = (window as unknown as { turnstile?: Turnstile })
      ?.turnstile
    if (!turnstile) {
      return resolve('')
    }
    const timer = setTimeout(() => resolve(''), 3000)
    try {
      turnstile.ready(() => {
        if (turnstileWidgetId) {
          turnstile.remove(turnstileWidgetId)
          turnstileWidgetId = null
        }

        turnstileWidgetId =
          turnstile.render('#cloudflare-turnstile', {
            sitekey: siteKey,
            action: 'search_messages',
            callback: (token) => {
              clearTimeout(timer)
              resolve(token)
            },
            'error-callback': () => {
              clearTimeout(timer)
              resolve('')
            },
          }) ?? null
      })
    } catch {
      clearTimeout(timer)
      resolve('')
    }
  })
}

function parseErrors(error: unknown): ErrorMessages {
  const bag = new ErrorMessages()
  const data = toApiError(error).data?.data
  if (data && typeof data === 'object') {
    Object.keys(data).forEach((key) => bag.addMany(key, data[key] ?? []))
  }
  return bag
}

async function fetchMessages(reset = false) {
  loading.value = true
  errorMessages.value = new ErrorMessages()
  errorTitle.value = ''

  if (reset) {
    page.value = 1
  }

  try {
    const token = await getCaptcha()
    const sort = sortBy.value[0]
    const results = await messagesStore.searchMessages({
      token,
      owners: formOwners.value,
      types: formTypes.value,
      statuses: formStatuses.value,
      query: formQuery.value,
      sort_by: sort?.key ?? 'created_at',
      sort_descending: sort ? sort.order === 'desc' : true,
      skip: (page.value - 1) * itemsPerPage.value,
      limit: itemsPerPage.value,
    } as SearchMessagesRequest)

    messages.value = results
    totalMessages.value = (page.value - 1) * itemsPerPage.value + results.length
    if (results.length === itemsPerPage.value) {
      totalMessages.value += 1
    }
  } catch (error: unknown) {
    errorTitle.value = capitalize(
      toApiError(error).data?.message ?? t('searchMessages.fetchError'),
    )
    errorMessages.value = parseErrors(error)
  } finally {
    loading.value = false
  }
}

function onUpdateOptions(options: {
  page: number
  itemsPerPage: number
  sortBy: { key: string; order: 'asc' | 'desc' }[]
}) {
  if (!initialLoadComplete.value) return
  page.value = options.page
  itemsPerPage.value = options.itemsPerPage
  if (options.sortBy.length > 0) {
    sortBy.value = options.sortBy
  }
  fetchMessages()
}

function exportMessages() {
  const headers = [
    'id',
    'created_at',
    'owner',
    'contact',
    'type',
    'status',
    'content',
  ]
  const rows = selectedMessages.value.map((m) => [
    m.id,
    m.created_at,
    m.owner,
    m.contact,
    m.type,
    m.status,
    `"${(m.content || '').replace(/"/g, '""')}"`,
  ])
  const csvContent =
    'data:text/csv;charset=utf-8,' +
    [headers.join(','), ...rows.map((e) => e.join(','))].join('\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', `httpsms_messages_${Date.now()}.csv`)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function deleteMessages() {
  loading.value = true
  try {
    await Promise.all(
      selectedMessages.value.map((message) =>
        messagesStore.deleteMessage(message.id),
      ),
    )
    notificationsStore.addNotification({
      message: t('searchMessages.deleteSuccess'),
      type: 'success',
    })
    selectedIds.value = []
  } catch {
    notificationsStore.addNotification({
      message: t('searchMessages.deleteError'),
      type: 'error',
    })
  } finally {
    loading.value = false
    showDeleteDialog.value = false
    fetchMessages()
  }
}

async function resendMessages() {
  loading.value = true
  const api = useApi()
  try {
    const results = await Promise.allSettled(
      selectedMessages.value.map((message) =>
        api('/v1/messages/send', {
          method: 'POST',
          body: {
            from: message.owner,
            to: message.contact,
            content: message.content,
            sim: message.sim,
            request_id: message.request_id,
          },
        }),
      ),
    )

    const failed = results.filter((r) => r.status === 'rejected')
    if (failed.length === 0) {
      notificationsStore.addNotification({
        message: t('searchMessages.resendSuccess'),
        type: 'success',
      })
      selectedIds.value = []
    } else if (failed.length === results.length) {
      notificationsStore.addNotification({
        message: t('searchMessages.resendError'),
        type: 'error',
      })
    } else {
      notificationsStore.addNotification({
        message: t('searchMessages.resendPartial', {
          success: results.length - failed.length,
          failed: failed.length,
        }),
        type: 'info',
      })
      selectedIds.value = []
    }
  } finally {
    loading.value = false
    showResendDialog.value = false
    fetchMessages()
  }
}

onMounted(async () => {
  try {
    await authStore.loadUser()
  } catch {
    // user load failed — continue with cached state
  }
  try {
    await phonesStore.loadPhones()
  } catch {
    // phones load failed — continue
  }

  const queryParam = route.query.query
  if (queryParam && typeof queryParam === 'string') {
    formQuery.value = queryParam
  }

  loading.value = false
  initialLoadComplete.value = true

  if (formQuery.value) {
    await fetchMessages(true)
  }
})

onBeforeUnmount(() => {
  const turnstile = (window as unknown as { turnstile?: Turnstile }).turnstile
  if (turnstile && turnstileWidgetId) {
    turnstile.remove(turnstileWidgetId)
    turnstileWidgetId = null
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
        <VToolbarTitle>
          <div class="py-16">{{ $t('searchMessages.title') }}</div>
        </VToolbarTitle>
        <VSpacer />
        <LanguageSwitcher class="mr-2" />
        <VProgressLinear
          :active="loading"
          color="primary"
          :indeterminate="loading"
          location="bottom"
          absolute
        />
      </VAppBar>
      <VContainer>
        <VRow>
          <VCol cols="12">
            <h5 class="text-headline-large mb-3 mt-0">
              {{ $t('searchMessages.title') }}
            </h5>
            <p>
              {{ $t('searchMessages.desc') }}
            </p>
            <VAlert v-if="errorTitle" variant="tonal" prominent type="warning">
              <h6 class="text-title-large font-weight-bold">
                {{ errorTitle }}
              </h6>
            </VAlert>
          </VCol>
        </VRow>
        <VCard>
          <VCardText class="pt-4" :class="{ 'pb-0': mdAndUp }">
            <VRow>
              <VCol cols="12" md="4">
                <VSelect
                  v-model="formOwners"
                  color="primary"
                  :error="errorMessages.has('owners')"
                  :error-messages="errorMessages.get('owners')"
                  :items="phoneNumberSelectItems"
                  multiple
                  density="compact"
                  :label="$t('searchMessages.phoneNumbers')"
                  variant="outlined"
                />
              </VCol>
              <VCol cols="12" md="4">
                <VSelect
                  v-model="formTypes"
                  color="primary"
                  :error="errorMessages.has('types')"
                  :error-messages="errorMessages.get('types')"
                  :items="messageTypeSelectItems"
                  density="compact"
                  multiple
                  :label="$t('searchMessages.messageTypes')"
                  variant="outlined"
                />
              </VCol>
              <VCol cols="12" md="4">
                <VSelect
                  v-model="formStatuses"
                  color="primary"
                  :error="errorMessages.has('statuses')"
                  :error-messages="errorMessages.get('statuses')"
                  :items="messageStatusSelectItems"
                  density="compact"
                  multiple
                  :label="$t('searchMessages.messageStatus')"
                  variant="outlined"
                />
              </VCol>
            </VRow>
            <VRow class="mt-n3">
              <VCol cols="12" md="8">
                <VTextField
                  v-model="formQuery"
                  color="primary"
                  :error="errorMessages.has('query')"
                  :error-messages="errorMessages.get('query')"
                  :label="$t('searchMessages.searchQuery')"
                  variant="outlined"
                  density="compact"
                  clearable
                  @keyup.enter="fetchMessages(true)"
                />
              </VCol>
              <VCol cols="12" md="4">
                <div id="cloudflare-turnstile" class="d-none"></div>
                <VBtn
                  :loading="loading"
                  :disabled="loading"
                  color="primary"
                  class="py-5"
                  :block="!mdAndUp"
                  @click="fetchMessages(true)"
                >
                  <VIcon v-if="mdAndUp" start :icon="mdiMagnify" />
                  <span>{{ $t('searchMessages.search') }}</span>
                </VBtn>
              </VCol>
            </VRow>
          </VCardText>
        </VCard>
        <VRow>
          <VCol cols="12" class="mt-16 mb-n2 d-flex align-center">
            <h2 class="text-md-headline-large text-headline-medium mb-0 mt-0">
              {{ $t('searchMessages.searchResults') }}
            </h2>
            <VDialog v-model="showDeleteDialog" opacity="0.9" max-width="550">
              <template #activator="{ props }">
                <VBtn
                  :loading="loading"
                  :disabled="loading || selectedMessages.length < 1"
                  size="small"
                  class="ml-2"
                  color="error"
                  v-bind="props"
                >
                  <VIcon v-if="mdAndUp" start :icon="mdiDelete" />
                  <span>{{ $t('searchMessages.deleteMessages') }}</span>
                </VBtn>
              </template>
              <VCard>
                <VCardTitle>
                  {{
                    $t('searchMessages.deleteConfirmTitle', {
                      count: selectedMessages.length,
                    })
                  }}
                </VCardTitle>
                <VCardText class="text-medium-emphasis">
                  {{ $t('searchMessages.deleteConfirmDesc') }}
                </VCardText>
                <VCardActions class="pb-4">
                  <VBtn
                    color="error"
                    :loading="loading"
                    variant="flat"
                    @click="deleteMessages"
                  >
                    {{ $t('searchMessages.deleteMessages') }}
                  </VBtn>
                  <VSpacer />
                  <VBtn color="warning" @click="showDeleteDialog = false">
                    {{ $t('common.close') }}
                  </VBtn>
                </VCardActions>
              </VCard>
            </VDialog>
            <VDialog v-model="showResendDialog" opacity="0.9" max-width="550">
              <template #activator="{ props }">
                <VBtn
                  :loading="loading"
                  :disabled="loading || !canResendSelected"
                  size="small"
                  class="ml-2 d-none d-md-inline-flex"
                  v-bind="props"
                >
                  <VIcon start :icon="mdiRefresh" />
                  {{ $t('searchMessages.resendMessages') }}
                </VBtn>
              </template>
              <VCard>
                <VCardTitle class="text-headline-medium text-break">
                  {{
                    $t('searchMessages.resendConfirmTitle', {
                      count: selectedMessages.length,
                    })
                  }}
                </VCardTitle>
                <VCardText class="text-medium-emphasis">
                  {{ $t('searchMessages.resendConfirmDesc') }}
                </VCardText>
                <VCardActions class="pb-4">
                  <VBtn
                    color="primary"
                    variant="flat"
                    :loading="loading"
                    @click="resendMessages"
                  >
                    {{ $t('searchMessages.resendMessages') }}
                  </VBtn>
                  <VSpacer />
                  <VBtn color="warning" @click="showResendDialog = false">
                    {{ $t('common.close') }}
                  </VBtn>
                </VCardActions>
              </VCard>
            </VDialog>
            <VSpacer />
            <VBtn
              :loading="loading"
              :disabled="loading || selectedMessages.length < 1"
              size="small"
              color="primary"
              @click="exportMessages"
            >
              <VIcon v-if="mdAndUp" start :icon="mdiExport" />
              <span>{{ $t('searchMessages.exportCsv') }}</span>
            </VBtn>
          </VCol>
          <VCol cols="12">
            <VDataTableServer
              v-model="selectedIds"
              v-model:items-per-page="itemsPerPage"
              v-model:page="page"
              v-model:sort-by="sortBy"
              color="primary"
              item-value="id"
              :headers="headers"
              :items="messages"
              :items-length="totalMessages"
              :items-per-page-options="itemsPerPageOptions"
              :page-text="pageText"
              :loading="loading"
              show-select
              :loading-text="$t('searchMessages.loading')"
              :no-data-text="$t('searchMessages.noData')"
              class="elevation-1"
              @update:options="onUpdateOptions"
            >
              <template #[`item.created_at`]="{ item }">
                {{ formatTimestamp(item.created_at) }}
              </template>
              <template #[`item.type`]="{ item }">
                <span v-if="item.type === 'call/missed'">
                  <VIcon size="small" color="error" :icon="mdiCallMissed" />
                  {{ $t('searchMessages.missedCalls') }}
                </span>
                <span v-else-if="item.type === 'mobile-originated'">
                  <VIcon size="small" :icon="mdiCallReceived" />
                  {{ $t('searchMessages.inbound') }}
                </span>
                <span v-else-if="item.type === 'mobile-terminated'">
                  <VIcon size="small" color="secondary" :icon="mdiCallMade" />
                  {{ $t('searchMessages.outbound') }}
                </span>
              </template>
              <template #[`item.status`]="{ item }">
                <VChip
                  v-if="item.status === 'expired'"
                  color="warning"
                  size="small"
                  variant="outlined"
                >
                  <VIcon size="small" start :icon="mdiAlert" />
                  {{ $t('searchMessages.expired') }}
                </VChip>
                <VChip
                  v-else-if="item.status === 'delivered'"
                  color="primary"
                  size="small"
                  variant="outlined"
                >
                  <VIcon size="small" start :icon="mdiCheckAll" />
                  {{ $t('searchMessages.delivered') }}
                </VChip>
                <VChip
                  v-else-if="item.status === 'received'"
                  color="success"
                  size="small"
                  variant="outlined"
                >
                  <VIcon size="small" start :icon="mdiCheckAll" />
                  {{ $t('searchMessages.received') }}
                </VChip>
                <VChip
                  v-else-if="item.status === 'sent'"
                  color="success"
                  size="small"
                  variant="outlined"
                >
                  <VIcon size="small" start :icon="mdiCheck" />
                  {{ $t('searchMessages.sent') }}
                </VChip>
                <VChip
                  v-else-if="item.status === 'failed'"
                  color="error"
                  size="small"
                  variant="outlined"
                >
                  <VIcon size="small" start :icon="mdiAlert" />
                  {{ $t('searchMessages.failed') }}
                </VChip>
                <VChip v-else size="small" color="cyan" variant="outlined">
                  <VIcon size="small" start :icon="mdiProgressCheck" />
                  {{ capitalize(item.status) }}
                </VChip>
              </template>
              <template #[`item.content`]="{ item }">
                <pre
                  style="
                    white-space: pre-wrap;
                    max-width: 300px;
                    word-break: break-all;
                  "
                  >{{ item.content }}</pre>
              </template>
            </VDataTableServer>
          </VCol>
        </VRow>
      </VContainer>
    </div>
  </VContainer>
</template>
