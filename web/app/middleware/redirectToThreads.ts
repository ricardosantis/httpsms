import { STORAGE_KEY } from '~/stores/redirectPreference'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  try {
    if (authStore.user?.is_admin) {
      return navigateTo('/admin/users', { replace: true })
    }
    if (localStorage.getItem(STORAGE_KEY) === 'true') {
      return navigateTo('/threads', { replace: true })
    }
  } catch (error) {
    console.error(error)
  }
})
