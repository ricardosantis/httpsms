<script setup lang="ts">
import {
  mdiArrowLeft,
  mdiAccountCircle,
  mdiShieldCheck,
  mdiEye,
  mdiEyeOff,
  mdiQrcode,
  mdiRefresh,
  mdiLinkVariant,
  mdiSquareEditOutline,
  mdiDelete,
  mdiContentSave,
  // mdiConnection, // Temporarily unused: Discord integration hidden for future use
  mdiCalendarClock,
  mdiPlus,
} from '@mdi/js'
import {
  getAuth,
  sendEmailVerification,
  signOut,
  type User as FirebaseUser,
} from 'firebase/auth'
import QRCode from 'qrcode'
import { ErrorMessages } from '~/utils/errors'
import { toApiError } from '~/utils/api-error'
import type {
  EntitiesPhone,
  EntitiesWebhook,
  // EntitiesDiscord, // Temporarily unused: Discord integration hidden for future use
  EntitiesMessageSendSchedule,
} from '~~/shared/types/api'

definePageMeta({
  middleware: ['auth'],
})

const { t, setLocale } = useI18n()

useHead({
  title: computed(() => `${t('settings.title')} - httpSMS`),
})

const config = useRuntimeConfig()
const route = useRoute()
const router = useRouter()
const { mdAndDown, mdAndUp, lgAndUp, xlAndUp, smAndUp } = useVDisplay()
const authStore = useAuthStore()
const appStore = useAppStore()
const phonesStore = usePhonesStore()
const billingStore = useBillingStore()
const notificationsStore = useNotificationsStore()
const redirectPreferenceStore = useRedirectPreferenceStore()

const firebaseUser = ref<FirebaseUser | null>(null)
const gravatarUrl = ref<string | null>(null)
const sendingVerificationEmail = ref(false)
const verificationEmailSent = ref(false)

async function sendVerificationEmail() {
  if (!firebaseUser.value) return
  sendingVerificationEmail.value = true
  try {
    await sendEmailVerification(firebaseUser.value)
    verificationEmailSent.value = true
    notificationsStore.addNotification({
      message: t('settings.verificationEmailSent'),
      type: 'success',
    })
  } catch (error) {
    console.error('sendEmailVerification failed:', error)
    notificationsStore.addNotification({
      message: t('settings.verificationEmailFailed'),
      type: 'error',
    })
  } finally {
    sendingVerificationEmail.value = false
  }
}

const computeGravatarUrl = async (email: string): Promise<string> => {
  try {
    const normalized = email.trim().toLowerCase()
    if (typeof crypto !== 'undefined' && crypto?.subtle) {
      const data = new TextEncoder().encode(normalized)
      const digest = await crypto.subtle.digest('SHA-256', data)
      const hash = Array.from(new Uint8Array(digest))
        .map((byte) => byte.toString(16).padStart(2, '0'))
        .join('')
      return `https://www.gravatar.com/avatar/${hash}?d=identicon&s=200`
    }
  } catch {
    // fallback
  }
  return 'https://www.gravatar.com/avatar/?d=identicon&s=200'
}

const avatarUrl = computed(
  () => firebaseUser.value?.photoURL ?? gravatarUrl.value,
)

const apiKeyShow = ref(false)
const showQrCodeDialog = ref(false)
const showRotateApiKey = ref(false)
const rotatingApiKey = ref(false)
const qrCodeCanvas = ref<HTMLCanvasElement | null>(null)

const errorMessages = ref(new ErrorMessages())

// Timezones
const standardTimezones = [
  'Africa/Abidjan',
  'Africa/Accra',
  'Africa/Cairo',
  'Africa/Johannesburg',
  'Africa/Lagos',
  'Africa/Nairobi',
  'America/Anchorage',
  'America/Argentina/Buenos_Aires',
  'America/Bogota',
  'America/Caracas',
  'America/Chicago',
  'America/Denver',
  'America/Halifax',
  'America/Lima',
  'America/Los_Angeles',
  'America/Manaus',
  'America/Mexico_City',
  'America/New_York',
  'America/Phoenix',
  'America/Santiago',
  'America/Sao_Paulo',
  'America/Toronto',
  'America/Vancouver',
  'Asia/Bangkok',
  'Asia/Dubai',
  'Asia/Hong_Kong',
  'Asia/Jakarta',
  'Asia/Jerusalem',
  'Asia/Kolkata',
  'Asia/Manila',
  'Asia/Seoul',
  'Asia/Shanghai',
  'Asia/Singapore',
  'Asia/Tokyo',
  'Atlantic/Azores',
  'Australia/Brisbane',
  'Australia/Melbourne',
  'Australia/Perth',
  'Australia/Sydney',
  'Europe/Amsterdam',
  'Europe/Berlin',
  'Europe/Dublin',
  'Europe/Lisbon',
  'Europe/London',
  'Europe/Madrid',
  'Europe/Paris',
  'Europe/Rome',
  'Pacific/Auckland',
  'Pacific/Honolulu',
  'UTC',
]

const timezones = (() => {
  try {
    const intlWithTimeZones = Intl as typeof Intl & {
      supportedValuesOf?: (key: string) => string[]
    }
    if (typeof intlWithTimeZones.supportedValuesOf === 'function') {
      const list = intlWithTimeZones.supportedValuesOf('timeZone')
      if (list && list.length > 0) return list
    }
    return standardTimezones
  } catch {
    return standardTimezones
  }
})()

const apiKey = computed(() => authStore.user?.api_key ?? '')

const hasActiveSubscription = computed(() => {
  if (authStore.user === null) return true
  return authStore.user.subscription_renews_at != null
})

const phoneNumbers = computed(() =>
  phonesStore.phones.map((phone) => phone.phone_number),
)

const webhookEventOptions = [
  'message.phone.received',
  'message.phone.sent',
  'message.phone.delivered',
  'message.send.failed',
  'message.send.expired',
  'message.call.missed',
  'phone.heartbeat.offline',
  'phone.heartbeat.online',
]

function resetErrors() {
  errorMessages.value = new ErrorMessages()
}

function parseErrors(error: unknown): ErrorMessages {
  const bag = new ErrorMessages()
  const data = toApiError(error).data?.data
  if (data && typeof data === 'object') {
    Object.keys(data).forEach((key) => bag.addMany(key, data[key]!))
  }
  return bag
}

// ---------------------------------------------------------------------------
// API Key
// ---------------------------------------------------------------------------
async function rotateApiKey() {
  if (!authStore.user) return
  rotatingApiKey.value = true
  try {
    await authStore.rotateApiKey(authStore.user.id)
    showRotateApiKey.value = false
    notificationsStore.addNotification({
      message: t('settings.rotateApiKeySuccess'),
      type: 'success',
    })
  } catch {
    notificationsStore.addNotification({
      message: t('settings.rotateApiKeyFailed'),
      type: 'error',
    })
  } finally {
    rotatingApiKey.value = false
  }
}

function generateQrCode() {
  showQrCodeDialog.value = true
  nextTick(() => {
    if (qrCodeCanvas.value) {
      QRCode.toCanvas(
        qrCodeCanvas.value,
        apiKey.value,
        { errorCorrectionLevel: 'H', width: 300, margin: 2 },
        (err) => {
          if (err) {
            notificationsStore.addNotification({
              message: t('settings.rotateApiKeyFailed'),
              type: 'error',
            })
          }
        },
      )
    }
  })
}

// ---------------------------------------------------------------------------
// Timezone & Language
// ---------------------------------------------------------------------------
async function updateTimezone(timezone: string) {
  try {
    await authStore.updateUser({ timezone })
    notificationsStore.addNotification({
      message: t('settings.timezoneSuccess'),
      type: 'success',
    })
  } catch {
    notificationsStore.addNotification({
      message: t('settings.timezoneFailed'),
      type: 'error',
    })
  }
}

async function updateLocale(localeCode: string) {
  try {
    setLocale(localeCode as 'pt-BR' | 'en')
    await authStore.updateUser({ locale: localeCode })
    notificationsStore.addNotification({
      message: t('common.savedSuccessfully') || 'Preferences saved',
      type: 'success',
    })
  } catch {
    // silently catch or notify
  }
}

// ---------------------------------------------------------------------------
// Webhooks
// ---------------------------------------------------------------------------
const loadingWebhooks = ref(true)
const webhooks = ref<EntitiesWebhook[]>([])
const updatingWebhook = ref(false)
const showWebhookEdit = ref(false)
const activeWebhook = ref<{
  id: string | null
  url: string
  signing_key: string
  phone_numbers: string[]
  events: string[]
}>({
  id: null,
  url: '',
  signing_key: '',
  phone_numbers: [],
  events: ['message.phone.received'],
})

async function loadWebhooks() {
  loadingWebhooks.value = true
  try {
    webhooks.value = await billingStore.getWebhooks()
  } finally {
    loadingWebhooks.value = false
  }
}

function onWebhookCreate() {
  resetErrors()
  activeWebhook.value = {
    id: null,
    url: '',
    signing_key: '',
    phone_numbers: phoneNumbers.value.slice(),
    events: [
      'message.phone.received',
      'message.phone.sent',
      'message.phone.delivered',
      'message.send.failed',
      'message.send.expired',
    ],
  }
  showWebhookEdit.value = true
}

function onWebhookEdit(webhookId: string) {
  const webhook = webhooks.value.find((x) => x.id === webhookId)
  if (!webhook) return
  resetErrors()
  activeWebhook.value = {
    id: webhook.id,
    url: webhook.url,
    signing_key: webhook.signing_key,
    phone_numbers: (webhook.phone_numbers ?? []).filter((x) =>
      phoneNumbers.value.includes(x),
    ),
    events: webhook.events ?? [],
  }
  showWebhookEdit.value = true
}

async function saveWebhook() {
  resetErrors()
  updatingWebhook.value = true
  try {
    const payload = {
      url: activeWebhook.value.url,
      signing_key: activeWebhook.value.signing_key,
      phone_numbers: activeWebhook.value.phone_numbers,
      events: activeWebhook.value.events,
    }
    if (activeWebhook.value.id) {
      await billingStore.updateWebhook({
        id: activeWebhook.value.id,
        ...payload,
      })
    } else {
      await billingStore.createWebhook(payload)
    }
    notificationsStore.addNotification({
      message: t('settings.webhookSaveSuccess', {
        action: activeWebhook.value.id ? t('common.save') : t('common.create'),
      }),
      type: 'success',
    })
    showWebhookEdit.value = false
    await loadWebhooks()
  } catch (error: unknown) {
    errorMessages.value = parseErrors(error)
    if (errorMessages.value.size() === 0) {
      notificationsStore.addNotification({
        message: t('settings.webhookSaveFailed'),
        type: 'error',
      })
    }
  } finally {
    updatingWebhook.value = false
  }
}

async function deleteWebhook(id: string) {
  updatingWebhook.value = true
  try {
    await billingStore.deleteWebhook(id)
    notificationsStore.addNotification({
      message: t('settings.webhookDeleteSuccess'),
      type: 'success',
    })
    showWebhookEdit.value = false
    await loadWebhooks()
  } catch {
    notificationsStore.addNotification({
      message: t('settings.webhookDeleteFailed'),
      type: 'error',
    })
  } finally {
    updatingWebhook.value = false
  }
}

// ---------------------------------------------------------------------------
// ---------------------------------------------------------------------------
// [OCULTAÇÃO TEMPORÁRIA / FUTURE USE]
// Funções e variáveis de integração com Discord temporariamente comentadas.
// Descomente esta seção ao reativar a funcionalidade.
// ---------------------------------------------------------------------------
/*
const loadingDiscordIntegrations = ref(true)
const discords = ref<EntitiesDiscord[]>([])
const updatingDiscord = ref(false)
const showDiscordEdit = ref(false)
const activeDiscord = ref<{
  id: string | null
  name: string
  server_id: string
  incoming_channel_id: string
}>({
  id: null,
  name: '',
  server_id: '',
  incoming_channel_id: '',
})

async function loadDiscordIntegrations() {
  loadingDiscordIntegrations.value = true
  try {
    discords.value = await billingStore.getDiscordIntegrations()
  } finally {
    loadingDiscordIntegrations.value = false
  }
}

function onDiscordCreate() {
  resetErrors()
  activeDiscord.value = {
    id: null,
    name: '',
    server_id: '',
    incoming_channel_id: '',
  }
  showDiscordEdit.value = true
}

function onDiscordEdit(discordId: string) {
  const discord = discords.value.find((x) => x.id === discordId)
  if (!discord) return
  resetErrors()
  activeDiscord.value = {
    id: discord.id,
    name: discord.name,
    server_id: discord.server_id,
    incoming_channel_id: discord.incoming_channel_id,
  }
  showDiscordEdit.value = true
}

async function saveDiscord() {
  resetErrors()
  updatingDiscord.value = true
  try {
    const payload = {
      name: activeDiscord.value.name,
      server_id: activeDiscord.value.server_id,
      incoming_channel_id: activeDiscord.value.incoming_channel_id,
    }
    if (activeDiscord.value.id) {
      await billingStore.updateDiscordIntegration({
        id: activeDiscord.value.id,
        ...payload,
      })
    } else {
      await billingStore.createDiscord(payload)
    }
    notificationsStore.addNotification({
      message: t('settings.discordSaveSuccess', {
        action: activeDiscord.value.id ? t('common.save') : t('common.create'),
      }),
      type: 'success',
    })
    showDiscordEdit.value = false
    await loadDiscordIntegrations()
  } catch (error: unknown) {
    errorMessages.value = parseErrors(error)
    if (errorMessages.value.size() === 0) {
      notificationsStore.addNotification({
        message: t('settings.discordSaveFailed'),
        type: 'error',
      })
    }
  } finally {
    updatingDiscord.value = false
  }
}

async function deleteDiscord(id: string) {
  updatingDiscord.value = true
  try {
    await billingStore.deleteDiscordIntegration(id)
    notificationsStore.addNotification({
      message: t('settings.discordDeleteSuccess'),
      type: 'success',
    })
    showDiscordEdit.value = false
    await loadDiscordIntegrations()
  } catch {
    notificationsStore.addNotification({
      message: t('settings.discordDeleteFailed'),
      type: 'error',
    })
  } finally {
    updatingDiscord.value = false
  }
}
*/

// ---------------------------------------------------------------------------
// Phones
// ---------------------------------------------------------------------------
const updatingPhone = ref(false)
const showPhoneEdit = ref(false)
const activePhone = ref<EntitiesPhone | null>(null)

function showEditPhone(phoneId: string) {
  const phone = phonesStore.phones.find((x) => x.id === phoneId)
  if (!phone) return
  resetErrors()
  activePhone.value = { ...phone }
  showPhoneEdit.value = true
}

async function updatePhone() {
  if (!activePhone.value) return
  updatingPhone.value = true
  try {
    await phonesStore.updatePhone(activePhone.value)
    showPhoneEdit.value = false
    activePhone.value = null
  } finally {
    updatingPhone.value = false
  }
}

async function deletePhone(phoneId: string) {
  updatingPhone.value = true
  try {
    await phonesStore.deletePhone(phoneId)
    notificationsStore.addNotification({
      message: t('settings.phoneDeleteSuccess'),
      type: 'success',
    })
    showPhoneEdit.value = false
    activePhone.value = null
  } catch {
    notificationsStore.addNotification({
      message: t('settings.phoneDeleteFailed'),
      type: 'error',
    })
  } finally {
    updatingPhone.value = false
  }
}

// ---------------------------------------------------------------------------
// Send Schedules
// ---------------------------------------------------------------------------
const loadingSendSchedules = ref(true)
const sendSchedules = ref<EntitiesMessageSendSchedule[]>([])
const showScheduleEdit = ref(false)
const showScheduleDelete = ref(false)
const savingSchedule = ref(false)
const activeSchedule = ref<{
  id: string | null
  name: string
  timezone: string
  windows: Array<{
    day_of_week: number
    start_time: string
    end_time: string
  }>
}>({
  id: null,
  name: '',
  timezone: '',
  windows: [],
})

const weekDays = computed(() => [
  { value: 1, label: t('settings.monday') },
  { value: 2, label: t('settings.tuesday') },
  { value: 3, label: t('settings.wednesday') },
  { value: 4, label: t('settings.thursday') },
  { value: 5, label: t('settings.friday') },
  { value: 6, label: t('settings.saturday') },
  { value: 0, label: t('settings.sunday') },
])

async function loadSendSchedules() {
  loadingSendSchedules.value = true
  try {
    sendSchedules.value = await billingStore.getSendSchedules()
  } finally {
    loadingSendSchedules.value = false
  }
}

function minuteToClock(value: number): string {
  const hours = String(Math.floor(value / 60)).padStart(2, '0')
  const minutes = String(value % 60).padStart(2, '0')
  return `${hours}:${minutes}`
}

function clockToMinute(value: string): number {
  if (!value || !value.includes(':')) return 0
  const [hours = 0, minutes = 0] = value.split(':').map((x) => parseInt(x, 10))
  return hours * 60 + minutes
}

function scheduleSummary(schedule: EntitiesMessageSendSchedule): string[][] {
  if (!schedule) return []
  return weekDays.value
    .map((day) => {
      const windows = (schedule.windows || []).filter(
        (x) => x && x.day_of_week === day.value,
      )
      if (windows.length === 0) return []
      return [
        day.label,
        windows
          .map(
            (w) =>
              `${minuteToClock(w.start_minute ?? 0)} - ${minuteToClock(w.end_minute ?? 0)}`,
          )
          .join(', '),
      ]
    })
    .filter((x) => x.length > 0)
}

function defaultTimezone(): string {
  try {
    return (
      authStore.user?.timezone ||
      Intl.DateTimeFormat().resolvedOptions().timeZone ||
      'America/Sao_Paulo'
    )
  } catch {
    return 'America/Sao_Paulo'
  }
}

function openCreateSchedule() {
  resetErrors()
  activeSchedule.value = {
    id: null,
    name: '',
    timezone: defaultTimezone(),
    windows: [1, 2, 3, 4, 5].map((day) => ({
      day_of_week: day,
      start_time: '09:00',
      end_time: '17:00',
    })),
  }
  showScheduleEdit.value = true
}

function openEditSchedule(schedule: EntitiesMessageSendSchedule) {
  resetErrors()
  activeSchedule.value = {
    id: schedule.id,
    name: schedule.name,
    timezone: schedule.timezone,
    windows: (schedule.windows || []).map((x) => ({
      day_of_week: x.day_of_week,
      start_time: minuteToClock(x.start_minute),
      end_time: minuteToClock(x.end_minute),
    })),
  }
  showScheduleEdit.value = true
}

function scheduleWindowsForDay(dayOfWeek: number) {
  return activeSchedule.value.windows.filter((w) => w.day_of_week === dayOfWeek)
}

function scheduleDayEnabled(dayOfWeek: number): boolean {
  return activeSchedule.value.windows.some((w) => w.day_of_week === dayOfWeek)
}

function scheduleToggleDay(dayOfWeek: number, enabled: boolean | null) {
  if (enabled) {
    if (!scheduleDayEnabled(dayOfWeek)) {
      activeSchedule.value.windows.push({
        day_of_week: dayOfWeek,
        start_time: '09:00',
        end_time: '17:00',
      })
    }
  } else {
    activeSchedule.value.windows = activeSchedule.value.windows.filter(
      (w) => w.day_of_week !== dayOfWeek,
    )
  }
}

function scheduleAddWindow(dayOfWeek: number) {
  activeSchedule.value.windows.push({
    day_of_week: dayOfWeek,
    start_time: '09:00',
    end_time: '17:00',
  })
}

function scheduleRemoveWindow(dayOfWeek: number, indexInDay: number) {
  const matches = activeSchedule.value.windows
    .map((w, idx) => ({ w, idx }))
    .filter(({ w }) => w.day_of_week === dayOfWeek)
  const target = matches[indexInDay]
  if (target) {
    activeSchedule.value.windows.splice(target.idx, 1)
  }
}

function scheduleWindowError(dayOfWeek: number): string | null {
  const windows = scheduleWindowsForDay(dayOfWeek)
  for (const w of windows) {
    if (!w.start_time || !w.end_time) return t('settings.scheduleSaveFailed')
    if (clockToMinute(w.start_time) >= clockToMinute(w.end_time)) {
      return `${t('settings.scheduleName')}: ${w.start_time} >= ${w.end_time}`
    }
  }
  return null
}

async function saveSchedule() {
  resetErrors()
  savingSchedule.value = true
  try {
    const rawWindows = activeSchedule.value.windows.map((w) => ({
      day_of_week: w.day_of_week,
      start_minute: clockToMinute(w.start_time),
      end_minute: clockToMinute(w.end_time),
    }))

    const payload = {
      name: activeSchedule.value.name,
      timezone: activeSchedule.value.timezone,
      windows: rawWindows,
    }

    if (activeSchedule.value.id) {
      await billingStore.updateSendSchedule({
        id: activeSchedule.value.id,
        ...payload,
      })
    } else {
      await billingStore.createSendSchedule(payload)
    }

    notificationsStore.addNotification({
      type: 'success',
      message: t('settings.scheduleSaveSuccess'),
    })
    showScheduleEdit.value = false
    await loadSendSchedules()
  } catch (error: unknown) {
    errorMessages.value = parseErrors(error)
    if (errorMessages.value.size() != 0) {
      notificationsStore.addNotification({
        type: 'error',
        message: t('settings.scheduleSaveFailed'),
      })
    }
  } finally {
    savingSchedule.value = false
  }
}

function confirmDeleteSchedule() {
  showScheduleDelete.value = true
}

async function deleteSchedule() {
  if (!activeSchedule.value.id) return
  savingSchedule.value = true
  try {
    await billingStore.deleteSendSchedule(activeSchedule.value.id)
    notificationsStore.addNotification({
      type: 'success',
      message: t('settings.scheduleDeleteSuccess'),
    })
    showScheduleDelete.value = false
    showScheduleEdit.value = false
    await loadSendSchedules()
  } catch {
    notificationsStore.addNotification({
      type: 'error',
      message: t('settings.scheduleDeleteFailed'),
    })
  } finally {
    savingSchedule.value = false
  }
}

// ---------------------------------------------------------------------------
// Email Notifications
// ---------------------------------------------------------------------------
const updatingEmailNotifications = ref(false)
const notificationSettings = ref({
  webhook_enabled: true,
  message_status_enabled: true,
  newsletter_enabled: true,
  heartbeat_enabled: true,
})

function syncEmailNotifications() {
  if (!authStore.user) return
  notificationSettings.value = {
    webhook_enabled: authStore.user.notification_webhook_enabled,
    message_status_enabled: authStore.user.notification_message_status_enabled,
    heartbeat_enabled: authStore.user.notification_heartbeat_enabled,
    newsletter_enabled: authStore.user.notification_newsletter_enabled,
  }
}

async function saveEmailNotifications() {
  if (!authStore.user) return
  updatingEmailNotifications.value = true
  try {
    await billingStore.saveEmailNotifications(
      authStore.user.id,
      notificationSettings.value,
    )
    notificationsStore.addNotification({
      message: t('settings.notificationSettingsSuccess'),
      type: 'success',
    })
    syncEmailNotifications()
  } catch {
    notificationsStore.addNotification({
      message: t('settings.notificationSettingsFailed'),
      type: 'error',
    })
  } finally {
    updatingEmailNotifications.value = false
  }
}

// ---------------------------------------------------------------------------
// Delete account
// ---------------------------------------------------------------------------
const deletingAccount = ref(false)
const showDeleteAccountDialog = ref(false)

async function deleteUserAccount() {
  deletingAccount.value = true
  try {
    const message = await authStore.deleteUserAccount()
    notificationsStore.addNotification({
      message: message ?? t('settings.deleteAccountConfirmTitle'),
      type: 'success',
    })
    const auth = getAuth()
    await signOut(auth)
    authStore.resetState()
    phonesStore.resetState()
    redirectPreferenceStore.resetState()
    notificationsStore.addNotification({
      type: 'info',
      message: t('auth.logoutSuccess'),
    })
    await router.push({ name: 'index' })
  } catch {
    notificationsStore.addNotification({
      message: t('auth.unexpectedError'),
      type: 'error',
    })
  } finally {
    deletingAccount.value = false
    showDeleteAccountDialog.value = false
  }
}

watch(showQrCodeDialog, (open) => {
  if (open && apiKey.value) {
    nextTick(() => generateQrCode())
  }
})

onMounted(() => {
  try {
    firebaseUser.value = getAuth().currentUser
    if (firebaseUser.value?.email) {
      computeGravatarUrl(firebaseUser.value.email).then((url) => {
        gravatarUrl.value = url
      })
    }
    void Promise.allSettled([
      authStore.loadUser().then(() => syncEmailNotifications()),
      phonesStore.loadPhones(),
      loadWebhooks(),
      // loadDiscordIntegrations(), // Temporarily hidden: Discord integration is hidden for future use
      loadSendSchedules(),
    ])
    if (route.hash) {
      nextTick(() => {
        try {
          const el = document.querySelector(route.hash)
          if (el) el.scrollIntoView({ behavior: 'smooth' })
        } catch {
          // ignore invalid query selector
        }
      })
    }
  } catch (err) {
    console.error('Error mounting settings:', err)
  }
})
</script>

<template>
  <VContainer fluid :class="{ 'fill-height': lgAndUp }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>{{ $t('settings.title') }}</VToolbarTitle>
        <VSpacer />
        <LanguageSwitcher class="mr-2" />
      </VAppBar>
      <VContainer class="pa-0">
        <VRow>
          <VCol cols="12" md="9" offset-md="1" xl="8" offset-xl="2">
            <!-- Profile -->
            <div v-if="firebaseUser" class="text-center">
              <VAvatar v-if="avatarUrl" size="100" :image="avatarUrl" />
              <v-avatar v-else size="100">
                <VIcon size="80" :icon="mdiAccountCircle" />
              </v-avatar>

              <h3
                v-if="firebaseUser.displayName"
                class="text-title-large mt-2 mb-0"
              >
                {{ firebaseUser.displayName }}
              </h3>
              <h4 class="text-medium-emphasis mb-2 mt-0">
                {{ firebaseUser.email }}
                <VIcon
                  v-if="firebaseUser.emailVerified"
                  size="small"
                  color="primary"
                  :icon="mdiShieldCheck"
                />
                <VBtn
                  v-else
                  size="x-small"
                  variant="tonal"
                  color="warning"
                  :loading="sendingVerificationEmail"
                  :disabled="verificationEmailSent"
                  @click="sendVerificationEmail"
                >
                  {{ $t('settings.verifyEmail') }}
                </VBtn>
              </h4>
              <VAutocomplete
                v-if="authStore.user"
                density="compact"
                variant="outlined"
                :model-value="authStore.user.timezone"
                class="mx-auto mt-2"
                style="max-width: 250px"
                :label="$t('settings.timezone')"
                :items="timezones"
                @update:model-value="updateTimezone"
              />
              <VSelect
                v-if="authStore.user"
                density="compact"
                variant="outlined"
                :model-value="authStore.user.locale || 'pt-BR'"
                class="mx-auto mt-2"
                style="max-width: 250px"
                label="Idioma / Language"
                :items="[
                  { title: 'Português (Brasil)', value: 'pt-BR' },
                  { title: 'English', value: 'en' },
                ]"
                @update:model-value="updateLocale"
              />
            </div>

            <!-- API Key -->
            <h5 class="text-headline-large mb-3 mt-0">
              {{ $t('settings.apiKey') }}
            </h5>
            <p class="text-medium-emphasis">
              {{
                $t('settings.apiKeyHeaderDesc', {
                  code: 'x-api-key',
                  endpoint:
                    appStore.appData.apiBaseUrl ||
                    config.public.apiBaseUrl ||
                    'https://api.smsandroid.com.br',
                })
              }}
            </p>
            <div v-if="apiKey === ''" class="mb-n9 pl-3 pt-5">
              <VProgressCircular
                :size="20"
                :width="2"
                color="primary"
                indeterminate
              />
            </div>
            <form v-else autocomplete="off" @submit.prevent>
              <VTextField
                :append-inner-icon="apiKeyShow ? mdiEye : mdiEyeOff"
                :type="apiKeyShow ? 'text' : 'password'"
                :model-value="apiKey"
                readonly
                name="api-key"
                autocomplete="new-password"
                variant="outlined"
                class="mb-n2"
                @click:append-inner="apiKeyShow = !apiKeyShow"
              />
            </form>
            <div class="d-flex flex-wrap">
              <CopyButton
                :value="apiKey"
                color="primary"
                :copy-text="$t('settings.copyApiKey')"
                :notification-text="$t('settings.copyApiKeySuccess')"
              />
              <VBtn
                v-if="mdAndUp"
                color="primary"
                class="ml-4"
                @click="generateQrCode"
              >
                <VIcon start :icon="mdiQrcode" />
                {{ $t('settings.showQrCode') }}
              </VBtn>
              <VDialog
                v-model="showQrCodeDialog"
                max-width="400px"
                opacity="0.9"
              >
                <VCard>
                  <VCardTitle class="text-center">{{
                    $t('settings.qrCodeTitle')
                  }}</VCardTitle>
                  <VCardText class="text-center">
                    <p class="text-body-large mt-0">
                      {{
                        $t('settings.qrCodeDesc', {
                          link: $t('settings.httpSmsApp'),
                        })
                      }}
                    </p>
                    <canvas ref="qrCodeCanvas" />
                  </VCardText>
                  <VCardActions>
                    <VBtn
                      color="primary"
                      block
                      variant="flat"
                      class="mt-n4"
                      @click="showQrCodeDialog = false"
                      >{{ $t('common.close') }}</VBtn
                    >
                  </VCardActions>
                </VCard>
              </VDialog>
              <VBtn
                v-if="lgAndUp"
                class="ml-4"
                :href="config.public.appDocumentationUrl"
                >{{ $t('common.documentation') }}</VBtn
              >
              <VSpacer />
              <VDialog v-model="showRotateApiKey" max-width="550">
                <template #activator="{ props }">
                  <VBtn
                    :size="mdAndDown ? 'small' : 'default'"
                    :variant="lgAndUp ? 'text' : 'elevated'"
                    color="warning"
                    v-bind="props"
                  >
                    <VIcon start :icon="mdiRefresh" />
                    {{ $t('settings.rotateApiKey') }}
                  </VBtn>
                </template>
                <VCard>
                  <VCardTitle class="text-headline-small">{{
                    $t('settings.rotateApiKeyConfirm')
                  }}</VCardTitle>
                  <VCardText class="text-medium-emphasis">
                    {{ $t('settings.rotateApiKeyDesc', { app: 'httpSMS' }) }}
                  </VCardText>
                  <VCardActions class="pb-4">
                    <VBtn
                      color="primary"
                      variant="flat"
                      :loading="rotatingApiKey"
                      @click="rotateApiKey"
                    >
                      <VIcon start :icon="mdiRefresh" />
                      {{ $t('settings.yesRotateKey') }}
                    </VBtn>
                    <VSpacer />
                    <VBtn
                      variant="text"
                      color="warning"
                      @click="showRotateApiKey = false"
                      >{{ $t('common.close') }}</VBtn
                    >
                  </VCardActions>
                </VCard>
              </VDialog>
            </div>

            <!-- Webhooks -->
            <h5 id="webhook-settings" class="text-headline-large mb-3 mt-12">
              {{ $t('settings.webhooks') }}
            </h5>
            <p class="text-medium-emphasis">
              {{ $t('settings.webhooksDesc') }}
            </p>
            <div v-if="loadingWebhooks">
              <VProgressCircular
                :size="60"
                :width="2"
                color="primary"
                class="mb-4"
                indeterminate
              />
            </div>
            <VTable v-else-if="webhooks.length" class="mb-4">
              <thead>
                <tr class="text-uppercase text-title-medium">
                  <th v-if="xlAndUp" class="text-left">ID</th>
                  <th class="text-left text-break">
                    {{ $t('settings.callbackUrl') }}
                  </th>
                  <th v-if="lgAndUp" class="text-center">
                    {{ $t('settings.events') }}
                  </th>
                  <th class="text-center">{{ $t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="webhook in webhooks" :key="webhook.id">
                  <td v-if="xlAndUp" class="text-left">{{ webhook.id }}</td>
                  <td class="text-break">{{ webhook.url }}</td>
                  <td v-if="lgAndUp" class="text-center">
                    <VChip
                      v-for="event in webhook.events ?? []"
                      :key="event"
                      class="ma-1"
                      size="small"
                      >{{ event }}</VChip
                    >
                  </td>
                  <td class="text-center">
                    <VBtn
                      :icon="mdAndDown"
                      size="small"
                      color="info"
                      :disabled="updatingWebhook"
                      @click.prevent="onWebhookEdit(webhook.id)"
                    >
                      <VIcon size="small" :icon="mdiSquareEditOutline" />
                      <span v-if="!mdAndDown" class="ml-1">{{
                        $t('common.edit')
                      }}</span>
                    </VBtn>
                  </td>
                </tr>
              </tbody>
            </VTable>
            <div class="d-flex">
              <VBtn color="primary" @click="onWebhookCreate">
                <VIcon start :icon="mdiLinkVariant" />
                {{ $t('settings.addWebhook') }}
              </VBtn>
              <VBtn
                v-if="lgAndUp"
                class="ml-4"
                :href="`${appStore.appData.documentationUrl}/webhooks/introduction`"
                >{{ $t('common.documentation') }}</VBtn
              >
            </div>

            <!--
              [OCULTAÇÃO TEMPORÁRIA / FUTURE USE]
              O bloco de integração com o Discord foi temporariamente ocultado da interface.
              Toda a lógica e componentes foram preservados para reativação em uso futuro.
            -->
            <!--
            <h5 id="discord-settings" class="text-headline-large mb-3 mt-12">
              {{ $t('settings.discordTitle') }}
            </h5>
            <p class="text-medium-emphasis">
              {{ $t('settings.discordDesc', { code: '/httpsms' }) }}
            </p>
            <div v-if="loadingDiscordIntegrations">
              <VProgressCircular
                :size="60"
                :width="2"
                color="primary"
                class="mb-4"
                indeterminate
              />
            </div>
            <VTable v-else-if="discords.length" class="mb-4">
              <thead>
                <tr class="text-uppercase text-title-medium">
                  <th class="text-left">{{ $t('common.name') }}</th>
                  <th class="text-left">
                    {{ $t('settings.discordServerId') }}
                  </th>
                  <th class="text-left">
                    {{ $t('settings.discordChannelId') }}
                  </th>
                  <th class="text-center">{{ $t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="discord in discords" :key="discord.id">
                  <td class="text-left">{{ discord.name }}</td>
                  <td class="text-left">{{ discord.server_id }}</td>
                  <td class="text-left">{{ discord.incoming_channel_id }}</td>
                  <td class="text-center">
                    <VBtn
                      :icon="mdAndDown"
                      size="small"
                      color="info"
                      :disabled="updatingDiscord"
                      @click.prevent="onDiscordEdit(discord.id)"
                    >
                      <VIcon size="small" :icon="mdiSquareEditOutline" />
                      <span v-if="!mdAndDown" class="ml-1">{{
                        $t('common.edit')
                      }}</span>
                    </VBtn>
                  </td>
                </tr>
              </tbody>
            </VTable>
            <VBtn color="primary" @click="onDiscordCreate">
              <VIcon start :icon="mdiConnection" />
              {{ $t('settings.addDiscord') }}
            </VBtn>
            -->

            <!-- Phones -->
            <h5 id="phones" class="text-headline-large mb-3 mt-12">
              {{ $t('settings.phonesTitle') }}
            </h5>
            <p class="text-medium-emphasis">
              {{ $t('settings.phonesDesc') }}
            </p>
            <VTable class="mb-4" density="comfortable">
              <thead>
                <tr class="text-uppercase text-medium-emphasis">
                  <th v-if="xlAndUp" class="text-left">ID</th>
                  <th class="text-left">{{ $t('newMessage.phoneNumber') }}</th>
                  <th v-if="lgAndUp" class="text-center">
                    {{ $t('settings.phoneRetries') }}
                  </th>
                  <th class="text-center">{{ $t('settings.phoneRate') }}</th>
                  <th class="text-center">
                    {{ $t('settings.phoneUpdatedAt') }}
                  </th>
                  <th class="text-center">{{ $t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="phone in phonesStore.phones" :key="phone.id">
                  <td v-if="xlAndUp" class="text-left">{{ phone.id }}</td>
                  <td>
                    {{ useFilters().formatPhoneNumber(phone.phone_number) }}
                  </td>
                  <td v-if="lgAndUp" class="text-center">
                    {{ phone.max_send_attempts ? phone.max_send_attempts : 1 }}
                  </td>
                  <td class="text-center">
                    <span v-if="phone.messages_per_minute"
                      >{{ phone.messages_per_minute }}/min</span
                    >
                    <span v-else>{{ $t('common.all') }}</span>
                  </td>
                  <td class="text-center">
                    {{ useFilters().formatTimestamp(phone.updated_at) }}
                  </td>
                  <td class="text-center">
                    <VBtn
                      :icon="mdAndDown"
                      size="small"
                      color="info"
                      :disabled="updatingPhone"
                      @click.prevent="showEditPhone(phone.id)"
                    >
                      <VIcon size="small" :icon="mdiSquareEditOutline" />
                      <span v-if="!mdAndDown" class="ml-1">{{
                        $t('common.edit')
                      }}</span>
                    </VBtn>
                  </td>
                </tr>
              </tbody>
            </VTable>

            <!-- Send Schedules -->
            <h5 id="send-schedules" class="text-headline-large mb-3 mt-12">
              {{ $t('settings.sendSchedules') }}
            </h5>
            <p class="text-medium-emphasis">
              {{
                $t('settings.sendSchedulesDesc', {
                  link: $t('settings.configuredSendRate'),
                })
              }}
            </p>
            <div v-if="loadingSendSchedules">
              <VProgressCircular
                :size="60"
                :width="2"
                color="primary"
                class="mb-4"
                indeterminate
              />
            </div>
            <VTable class="mb-4" density="comfortable">
              <thead>
                <tr class="text-uppercase text-medium-emphasis">
                  <th class="text-left">{{ $t('common.name') }}</th>
                  <th class="text-left">{{ $t('settings.timezone') }}</th>
                  <th class="text-left">{{ $t('settings.sendSchedules') }}</th>
                  <th class="text-center">{{ $t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="schedule in sendSchedules" :key="schedule.id">
                  <td class="text-left pt-2" style="vertical-align: top">
                    {{ schedule.name }}
                  </td>
                  <td class="pt-2" style="vertical-align: top">
                    {{ schedule.timezone }}
                  </td>
                  <td class="py-2">
                    <div
                      v-for="line in scheduleSummary(schedule)"
                      :key="`${schedule.id}-${line[0]}`"
                      class="mb-1"
                    >
                      {{ line[0] }}:
                      <span class="text-medium-emphasis">{{ line[1] }}</span>
                    </div>
                  </td>
                  <td class="text-center pt-2" style="vertical-align: top">
                    <VBtn
                      :icon="mdAndDown"
                      size="small"
                      color="info"
                      :disabled="loadingSendSchedules"
                      @click.prevent="openEditSchedule(schedule)"
                    >
                      <VIcon size="small" :icon="mdiSquareEditOutline" />
                      <span v-if="!mdAndDown" class="ml-1">{{
                        $t('common.edit')
                      }}</span>
                    </VBtn>
                  </td>
                </tr>
              </tbody>
            </VTable>
            <div class="d-flex mt-4">
              <VBtn color="primary" @click="openCreateSchedule">
                <VIcon start :icon="mdiCalendarClock" />
                {{ $t('settings.createSendSchedule') }}
              </VBtn>
              <VBtn
                v-if="lgAndUp"
                class="ml-4"
                :href="`${appStore.appData.documentationUrl}/features/outgoing-message-queue`"
                >{{ $t('common.documentation') }}</VBtn
              >
            </div>

            <!-- Email Notifications -->
            <h5 id="email-notifications" class="text-headline-large mb-3 mt-12">
              {{ $t('settings.emailNotifications') }}
            </h5>
            <p class="text-medium-emphasis">
              {{ $t('settings.emailNotificationsDesc') }}
            </p>
            <VSwitch
              v-model="notificationSettings.heartbeat_enabled"
              color="primary"
              :label="$t('settings.heartbeatEmails')"
              :disabled="updatingEmailNotifications"
              :hint="$t('settings.heartbeatEmailsHint')"
              persistent-hint
            />
            <VSwitch
              v-model="notificationSettings.webhook_enabled"
              color="primary"
              :label="$t('settings.webhookEmails')"
              :disabled="updatingEmailNotifications"
              :hint="$t('settings.webhookEmailsHint')"
              persistent-hint
            />
            <VSwitch
              v-model="notificationSettings.message_status_enabled"
              color="primary"
              :label="$t('settings.messageStatusEmails')"
              :disabled="updatingEmailNotifications"
              :hint="$t('settings.messageStatusEmailsHint')"
              persistent-hint
            />
            <VSwitch
              v-model="notificationSettings.newsletter_enabled"
              color="primary"
              :label="$t('settings.newsletterEmails')"
              :disabled="updatingEmailNotifications"
              :hint="$t('settings.newsletterEmailsHint')"
              persistent-hint
            />
            <VBtn
              color="primary"
              :loading="updatingEmailNotifications"
              class="mt-4"
              @click="saveEmailNotifications"
            >
              <VIcon start :icon="mdiContentSave" />
              {{ $t('settings.saveNotificationSettings') }}
            </VBtn>

            <!-- Message Data Retention -->
            <h5 class="text-headline-large mb-3 mt-12">
              {{ $t('settings.dataRetention') }}
            </h5>
            <p class="text-medium-emphasis">
              {{
                $t('settings.dataRetentionDesc', {
                  link: $t('settings.messageSearchPage'),
                })
              }}
            </p>
            <VSelect
              :items="[$t('settings.year')]"
              :model-value="$t('settings.year')"
              :label="$t('settings.retentionPeriod')"
              variant="outlined"
              density="compact"
              class="mt-4"
              style="max-width: 300px"
            />

            <!-- Delete Account -->
            <h5 class="text-headline-large text-error mb-3 mt-10">
              {{ $t('settings.deleteAccount') }}
            </h5>
            <p v-if="hasActiveSubscription" class="text-medium-emphasis">
              {{
                $t('settings.deleteAccountActiveSub', {
                  link: $t('settings.cancelSubscription'),
                })
              }}
            </p>
            <p v-else class="text-medium-emphasis">
              {{
                $t('settings.deleteAccountDesc', {
                  irreversible: $t('settings.irreversible'),
                })
              }}
            </p>
            <VBtn
              color="error"
              :loading="deletingAccount"
              class="mt-2"
              :disabled="hasActiveSubscription"
              @click="showDeleteAccountDialog = true"
            >
              <VIcon start :icon="mdiDelete" />
              {{ $t('settings.deleteYourAccount') }}
            </VBtn>
            <VDialog v-model="showDeleteAccountDialog" max-width="600px">
              <VCard>
                <VCardTitle class="text-center">{{
                  $t('settings.deleteAccountConfirmTitle')
                }}</VCardTitle>
                <VCardText class="mt-2 text-center text-medium-emphasis">
                  {{
                    $t('settings.deleteAccountConfirmDesc', {
                      irreversible: $t('settings.irreversible'),
                    })
                  }}
                </VCardText>
                <VCardActions>
                  <VBtn
                    color="error"
                    variant="text"
                    :loading="deletingAccount"
                    @click="deleteUserAccount"
                  >
                    <VIcon v-if="lgAndUp" start :icon="mdiDelete" />
                    {{ $t('settings.deleteMyAccount') }}
                  </VBtn>
                  <VSpacer />
                  <VBtn
                    color="primary"
                    variant="flat"
                    @click="showDeleteAccountDialog = false"
                  >
                    <span v-if="lgAndUp">{{
                      $t('settings.keepMyAccount')
                    }}</span>
                    <span v-else>{{ $t('common.close') }}</span>
                  </VBtn>
                </VCardActions>
              </VCard>
            </VDialog>
          </VCol>
        </VRow>
      </VContainer>
    </div>

    <!-- Webhook Edit Dialog -->
    <VDialog v-model="showWebhookEdit" max-width="600px" opacity="0.9">
      <VCard>
        <VCardTitle>
          {{
            activeWebhook.id
              ? $t('settings.editWebhookTitle')
              : $t('settings.addWebhookTitle')
          }}
        </VCardTitle>
        <VCardText>
          <VRow>
            <VCol>
              <VTextField
                v-if="activeWebhook.id"
                variant="outlined"
                density="compact"
                disabled
                label="ID"
                :model-value="activeWebhook.id"
              />
              <VTextField
                v-model="activeWebhook.url"
                variant="outlined"
                density="compact"
                :label="$t('settings.callbackUrl')"
                persistent-placeholder
                persistent-hint
                :error="errorMessages.has('url')"
                :error-messages="errorMessages.get('url')"
                :hint="$t('settings.webhookCallbackHint')"
                placeholder="https://example.com/webhook"
              />
              <VTextField
                v-model="activeWebhook.signing_key"
                variant="outlined"
                density="compact"
                class="mt-6"
                persistent-placeholder
                persistent-hint
                :label="$t('settings.signingKeyOptional')"
                placeholder="******************"
                :error="errorMessages.has('signing_key')"
                :error-messages="errorMessages.get('signing_key')"
                :hint="$t('settings.webhookSigningKeyHint')"
              />
              <VSelect
                v-model="activeWebhook.events"
                :items="webhookEventOptions"
                :label="$t('settings.events')"
                multiple
                chips
                variant="outlined"
                persistent-placeholder
                class="mt-6"
                density="compact"
                :error="errorMessages.has('events')"
                :error-messages="errorMessages.get('events')"
                :hint="$t('settings.webhookEventsHint')"
                persistent-hint
              />
              <VSelect
                v-model="activeWebhook.phone_numbers"
                :items="phoneNumbers"
                :label="$t('settings.phoneNumbers')"
                multiple
                chips
                variant="outlined"
                persistent-placeholder
                class="mt-6"
                density="compact"
                :error="errorMessages.has('phone_numbers')"
                :error-messages="errorMessages.get('phone_numbers')"
                :hint="$t('settings.webhookPhonesHint')"
                persistent-hint
              />
            </VCol>
          </VRow>
        </VCardText>
        <VCardActions class="pb-4 px-4">
          <LoadingButton
            :icon="mdiContentSave"
            :loading="updatingWebhook"
            @click="saveWebhook"
          >
            {{
              activeWebhook.id
                ? $t('settings.updateWebhook')
                : $t('settings.saveWebhook')
            }}
          </LoadingButton>
          <VSpacer />
          <VBtn
            v-if="activeWebhook.id"
            :disabled="updatingWebhook"
            size="small"
            color="error"
            variant="text"
            @click="deleteWebhook(activeWebhook.id)"
          >
            <VIcon v-if="lgAndUp" start :icon="mdiDelete" />
            {{ $t('common.delete') }}
          </VBtn>
          <VBtn
            v-else
            variant="text"
            color="warning"
            @click="showWebhookEdit = false"
            >{{ $t('common.close') }}</VBtn
          >
        </VCardActions>
      </VCard>
    </VDialog>

    <!--
      [OCULTAÇÃO TEMPORÁRIA / FUTURE USE]
      Diálogo de edição do Discord temporariamente ocultado.
    -->
    <!--
    <VDialog v-model="showDiscordEdit" max-width="700px">
      <VCard>
        <VCardTitle>
          {{
            activeDiscord.id
              ? $t('settings.editDiscordTitle')
              : $t('settings.addDiscordTitle')
          }}
        </VCardTitle>
        <VCardText>
          <VRow>
            <VCol class="pt-8">
              <p class="mt-n4 text-body-1">
                {{ $t('settings.discordBotDesc') }}
              </p>
              <VBtn
                color="#5865f2"
                class="mb-6"
                target="_blank"
                href="https://discord.com/api/oauth2/authorize?client_id=1095780203256627291&permissions=2147485760&scope=bot%20applications.commands"
              >
                <VIcon start :icon="mdiConnection" />
                {{ $t('settings.addDiscordBot') }}
              </VBtn>
              <VTextField
                v-if="activeDiscord.id"
                variant="outlined"
                density="compact"
                disabled
                label="ID"
                :model-value="activeDiscord.id"
              />
              <VTextField
                v-model="activeDiscord.name"
                variant="outlined"
                density="compact"
                :label="$t('common.name')"
                persistent-placeholder
                persistent-hint
                :error="errorMessages.has('name')"
                :error-messages="errorMessages.get('name')"
                hint=""
                placeholder="e.g Game Server"
              />
              <VTextField
                v-model="activeDiscord.server_id"
                variant="outlined"
                density="compact"
                class="mt-6"
                persistent-placeholder
                persistent-hint
                :label="$t('settings.discordServerId')"
                placeholder="e.g 1095778291488653372"
                :error="errorMessages.has('server_id')"
                :error-messages="errorMessages.get('server_id')"
                :hint="$t('settings.discordServerIdHint')"
              />
              <VTextField
                v-model="activeDiscord.incoming_channel_id"
                variant="outlined"
                density="compact"
                class="mt-6"
                persistent-placeholder
                persistent-hint
                :label="$t('settings.discordChannelId')"
                placeholder="e.g 1095778291488653372"
                :error="errorMessages.has('incoming_channel_id')"
                :error-messages="errorMessages.get('incoming_channel_id')"
                :hint="$t('settings.discordChannelIdHint')"
              />
            </VCol>
          </VRow>
        </VCardText>
        <VCardActions class="pb-4 pl-6">
          <LoadingButton
            :icon="mdiContentSave"
            :loading="updatingDiscord"
            @click="saveDiscord"
          >
            {{
              activeDiscord.id
                ? $t('settings.updateDiscord')
                : $t('settings.saveDiscord')
            }}
          </LoadingButton>
          <VSpacer />
          <VBtn
            v-if="activeDiscord.id"
            :disabled="updatingDiscord"
            color="error"
            variant="text"
            @click="deleteDiscord(activeDiscord.id)"
          >
            <VIcon v-if="lgAndUp" start :icon="mdiDelete" />
            {{ $t('common.delete') }}
          </VBtn>
          <VBtn
            v-else
            variant="text"
            color="warning"
            @click="showDiscordEdit = false"
            >{{ $t('common.close') }}</VBtn
          >
        </VCardActions>
      </VCard>
    </VDialog>
    -->

    <!-- Phone Edit Dialog -->
    <VDialog v-model="showPhoneEdit" max-width="700px" opacity="0.9">
      <VCard>
        <VCardTitle>{{ $t('settings.phoneEditTitle') }}</VCardTitle>
        <VCardText v-if="activePhone">
          <VContainer>
            <VRow>
              <VCol>
                <VTextField
                  variant="outlined"
                  density="compact"
                  disabled
                  label="ID"
                  :model-value="activePhone.id"
                />
                <VTextField
                  variant="outlined"
                  disabled
                  density="compact"
                  :label="$t('newMessage.phoneNumber')"
                  :model-value="activePhone.phone_number"
                />
                <VTextField
                  variant="outlined"
                  disabled
                  density="compact"
                  label="SIM"
                  :model-value="activePhone.sim"
                />
                <VTextarea
                  variant="outlined"
                  disabled
                  density="compact"
                  label="FCM Token"
                  :model-value="activePhone.fcm_token"
                />
                <VTextField
                  v-model="activePhone.message_expiration_seconds"
                  variant="outlined"
                  type="number"
                  density="compact"
                  :label="$t('settings.messageExpirationSeconds')"
                />
                <VTextField
                  v-model="activePhone.messages_per_minute"
                  variant="outlined"
                  type="number"
                  density="compact"
                  :label="$t('settings.phoneMessagesPerMinute')"
                />
                <VTextField
                  v-model="activePhone.max_send_attempts"
                  variant="outlined"
                  type="number"
                  density="compact"
                  :placeholder="$t('settings.retriesPlaceholder')"
                  :label="$t('settings.phoneMaxSendAttempts')"
                  min="1"
                  max="5"
                  :rules="[
                    (v: number) =>
                      (v >= 1 && v <= 5) || $t('settings.retriesRule'),
                  ]"
                />
                <VAutocomplete
                  v-model="activePhone.message_send_schedule_id"
                  variant="outlined"
                  :readonly="sendSchedules.length === 0"
                  density="compact"
                  clearable
                  :label="$t('settings.sendSchedules')"
                  :items="sendSchedules"
                  item-title="name"
                  item-value="id"
                  hint=""
                />
                <VTextarea
                  v-model="activePhone.missed_call_auto_reply"
                  variant="outlined"
                  density="compact"
                  class="mt-6"
                  :label="$t('settings.missedCallAutoReply')"
                  persistent-placeholder
                  persistent-hint
                  :placeholder="$t('settings.missedCallAutoReplyPlaceholder')"
                  hint=""
                />
                <VSwitch
                  v-model="activePhone.unarchive_thread"
                  class="mt-4"
                  color="primary"
                  density="compact"
                  :label="$t('settings.unarchiveThreads')"
                  hint=""
                />
              </VCol>
            </VRow>
          </VContainer>
        </VCardText>
        <VCardActions class="pb-4 px-4 mt-n4">
          <loading-button :loading="updatingPhone" @click="updatePhone">
            <VIcon v-if="lgAndUp" start :icon="mdiContentSave" />
            {{ $t('common.save') }}
          </loading-button>
          <VSpacer />
          <VBtn
            color="error"
            variant="text"
            :disabled="updatingPhone"
            @click="deletePhone(activePhone?.id ?? '')"
          >
            <VIcon v-if="lgAndUp" start :icon="mdiDelete" />
            {{ $t('common.delete') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

    <!-- Send Schedule Edit Dialog -->
    <VDialog v-model="showScheduleEdit" max-width="800px" opacity="0.9">
      <VCard>
        <VCardTitle>
          <span v-if="!activeSchedule.id">{{
            $t('settings.createSendSchedule')
          }}</span>
          <span v-else>{{ $t('settings.createSendSchedule') }}</span>
        </VCardTitle>
        <VCardText class="mt-4" :class="{ 'px-2': mdAndDown }">
          <VRow>
            <VCol cols="12" md="6">
              <VTextField
                v-model="activeSchedule.name"
                variant="outlined"
                density="compact"
                persistent-placeholder
                :label="$t('settings.scheduleName')"
                placeholder="e.g Business Hours"
                :error="errorMessages.has('name')"
                :error-messages="errorMessages.get('name')"
              />
            </VCol>
            <VCol cols="12" md="6">
              <VAutocomplete
                v-model="activeSchedule.timezone"
                density="compact"
                variant="outlined"
                :items="timezones"
                :label="$t('settings.timezone')"
                :error="errorMessages.has('timezone')"
                :error-messages="errorMessages.get('timezone')"
              />
            </VCol>
          </VRow>
          <VCard variant="flat" :border="lgAndUp" class="px-0">
            <VCardText :class="mdAndDown ? 'px-2 mt-n4' : 'px-4'">
              <div
                v-for="day in weekDays"
                :key="day.value"
                :class="[smAndUp ? 'd-flex align-start' : '', 'mb-4']"
              >
                <div
                  :class="[smAndUp ? 'pr-4' : '', 'pt-2']"
                  :style="smAndUp ? 'min-width: 160px' : ''"
                >
                  <VSwitch
                    :model-value="scheduleDayEnabled(day.value)"
                    inset
                    density="compact"
                    color="primary"
                    :label="day.label"
                    hide-details
                    class="mt-0 pt-0"
                    @update:model-value="scheduleToggleDay(day.value, $event)"
                  />
                </div>
                <div class="pt-2 flex-grow-1">
                  <div
                    v-for="(window, index) in scheduleWindowsForDay(day.value)"
                    :key="`${day.value}-${index}`"
                    class="d-flex align-center flex-wrap mb-2"
                  >
                    <div
                      class="mr-2 mb-2"
                      style="width: 130px; max-width: 100%"
                    >
                      <VTextField
                        v-model="window.start_time"
                        density="compact"
                        variant="outlined"
                        :error="!!scheduleWindowError(day.value)"
                        type="time"
                        :label="$t('settings.scheduleStart')"
                        hide-details="auto"
                      />
                    </div>
                    <div class="mb-2 mr-2">–</div>
                    <div
                      class="mr-2 mb-2"
                      style="width: 130px; max-width: 100%"
                    >
                      <VTextField
                        v-model="window.end_time"
                        density="compact"
                        variant="outlined"
                        :error="!!scheduleWindowError(day.value)"
                        type="time"
                        :label="$t('settings.scheduleEnd')"
                        hide-details="auto"
                      />
                    </div>
                    <div class="mb-2">
                      <VBtn
                        v-if="index == 0"
                        icon
                        variant="text"
                        density="comfortable"
                        color="primary"
                        @click="scheduleAddWindow(day.value)"
                      >
                        <VIcon :icon="mdiPlus" />
                      </VBtn>
                      <VBtn
                        icon
                        density="comfortable"
                        variant="text"
                        class="ml-1"
                        color="error"
                        @click="scheduleRemoveWindow(day.value, index)"
                      >
                        <VIcon :icon="mdiDelete" />
                      </VBtn>
                    </div>
                  </div>
                  <div
                    v-if="scheduleWindowError(day.value)"
                    class="w-100 text-error mt-n2 mb-4"
                  >
                    {{ scheduleWindowError(day.value) }}
                  </div>
                </div>
              </div>
            </VCardText>
          </VCard>
        </VCardText>
        <VCardActions class="pb-4 mt-n2">
          <LoadingButton
            :icon="mdiContentSave"
            :loading="savingSchedule"
            @click="saveSchedule"
          >
            {{
              activeSchedule.id
                ? $t('common.save')
                : $t('settings.createSendSchedule')
            }}
          </LoadingButton>
          <VSpacer />
          <VBtn
            v-if="activeSchedule.id"
            :disabled="savingSchedule"
            color="error"
            variant="text"
            @click="confirmDeleteSchedule"
          >
            <VIcon v-if="lgAndUp" start :icon="mdiDelete" />
            {{ $t('common.delete') }}
          </VBtn>
          <VBtn
            v-else
            variant="text"
            color="warning"
            @click="showScheduleEdit = false"
          >
            {{ $t('common.close') }}
          </VBtn>
        </VCardActions>
      </VCard>
    </VDialog>

    <!-- Send Schedule Delete Confirmation -->
    <VDialog v-model="showScheduleDelete" max-width="500" opacity="0.9">
      <VCard>
        <VCardTitle>{{ $t('common.delete') }}</VCardTitle>
        <VCardText class="text-medium-emphasis">
          {{ activeSchedule.name }}
        </VCardText>
        <VCardActions>
          <VBtn
            variant="flat"
            color="error"
            :loading="savingSchedule"
            @click="deleteSchedule"
          >
            {{ $t('common.delete') }}
          </VBtn>
          <VSpacer />
          <VBtn variant="text" @click="showScheduleDelete = false">{{
            $t('common.cancel')
          }}</VBtn>
        </VCardActions>
      </VCard>
    </VDialog>
  </VContainer>
</template>
