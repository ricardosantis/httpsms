import { defineStore } from 'pinia'
import type { User as FirebaseUser } from 'firebase/auth'
import { setAuthHeader, setApiKey } from '~/composables/useApi'
import type { EntitiesUser } from '~~/shared/types/api'
import { getApiErrorMessage } from '~/utils/api-error'

export interface AuthUser {
  email: string | null
  displayName: string | null
  id: string
}

export const useAuthStore = defineStore('auth', () => {
  const authStateChanged = ref(false)
  const authUser = ref<AuthUser | null>(null)
  const user = ref<EntitiesUser | null>(null)
  const { apiFetch } = useApi()
  const notificationsStore = useNotificationsStore()

  function t(key: string, fallback: string): string {
    try {
      return useNuxtApp().$i18n.t(key)
    } catch {
      return fallback
    }
  }

  async function setAuthUserAction(newUser: AuthUser | null | undefined) {
    const userChanged = newUser?.id !== authUser.value?.id
    authUser.value = newUser ?? null
    authStateChanged.value = true

    if (userChanged && newUser !== null) {
      await Promise.all([loadUser(), loadPhones()])
    }
  }

  async function onAuthStateChanged(firebaseUser: FirebaseUser | null) {
    if (firebaseUser == null) {
      authUser.value = null
      user.value = null
      authStateChanged.value = true
      setApiKey('')
      return
    }
    setAuthHeader(await firebaseUser.getIdToken(true))
    const { uid, email, displayName } = firebaseUser
    authUser.value = { id: uid, email, displayName }
    authStateChanged.value = true
  }

  async function onIdTokenChanged(firebaseUser: FirebaseUser | null) {
    if (firebaseUser == null) {
      setApiKey('')
      return
    }
    setAuthHeader(await firebaseUser.getIdToken(true))
  }

  async function loadUser() {
    let setLocale: ((locale: string) => void) | undefined
    let locales: globalThis.Ref<Array<string | { code: string }>> | undefined
    try {
      const i18n = useI18n()
      setLocale = i18n.setLocale
      locales = i18n.locales as unknown as globalThis.Ref<
        Array<string | { code: string }>
      >
    } catch {
      // Ignore if called outside Nuxt context
    }

    try {
      const response = await apiFetch<{ data: EntitiesUser }>('/v1/users/me')
      user.value = response.data
      if (user.value?.locale && setLocale && locales) {
        try {
          const availableLocales = locales.value.map((l) =>
            typeof l === 'string' ? l : l.code,
          )
          const requested = user.value.locale
          const matchedLocale =
            availableLocales.find((code: string) => code === requested) ??
            availableLocales.find((code: string) =>
              requested.toLowerCase().startsWith(`${code.toLowerCase()}-`),
            )
          if (matchedLocale) {
            setLocale(matchedLocale)
          }
        } catch (e) {
          console.error('Failed to set locale in loadUser:', e)
        }
      }
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t('authStore.errorLoadingUserProfile', 'Error loading user profile'),
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function updateUser(payload: {
    owner?: string
    timezone?: string
    locale?: string
  }) {
    const phonesStore = usePhonesStore()
    if (payload.owner) {
      phonesStore.setOwner(payload.owner)
    }

    const activePhone = phonesStore.activePhone
    const activePhoneId =
      activePhone?.id ||
      (payload.owner
        ? phonesStore.phones.find((p) => p.phone_number === payload.owner)?.id
        : user.value?.active_phone_id)

    const body: {
      active_phone_id?: string
      timezone?: string
      locale?: string
    } = {}
    if (activePhoneId) {
      body.active_phone_id = activePhoneId
    }
    if (payload.timezone) {
      body.timezone = payload.timezone
    } else if (user.value?.timezone) {
      body.timezone = user.value.timezone
    }
    if (payload.locale) {
      body.locale = payload.locale
    } else if (user.value?.locale) {
      body.locale = user.value.locale
    }

    const response = await apiFetch<{ data: EntitiesUser }>('/v1/users/me', {
      method: 'PUT',
      body,
    })

    setApiKey(response.data.api_key)
    user.value = response.data
  }

  async function deleteUserAccount(): Promise<string> {
    try {
      await apiFetch<{ message: string }>('/v1/users/me', {
        method: 'DELETE',
      })
      return 'Your account has been deleted successfully'
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t('authStore.errorDeletingAccount', 'Error deleting account'),
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function rotateApiKey(userId: string): Promise<EntitiesUser> {
    try {
      const response = await apiFetch<{ data: EntitiesUser }>(
        `/v1/users/${userId}/api-keys`,
        {
          method: 'DELETE',
        },
      )
      user.value = response.data
      setApiKey(response.data.api_key)
      return response.data
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t('authStore.errorRotatingApiKey', 'Error rotating API key'),
        ),
        type: 'error',
      })
      throw error
    }
  }

  function resetState() {
    user.value = null
    authUser.value = null
    authStateChanged.value = true
    setApiKey('')
  }

  function loadPhones() {
    const phonesStore = usePhonesStore()
    return phonesStore.loadPhones(false)
  }

  return {
    authStateChanged,
    authUser,
    user,
    setAuthUserAction,
    onAuthStateChanged,
    onIdTokenChanged,
    loadUser,
    updateUser,
    deleteUserAccount,
    rotateApiKey,
    resetState,
  }
})
