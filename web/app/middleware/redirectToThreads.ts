import { STORAGE_KEY } from '~/stores/redirectPreference'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  await new Promise<void>((resolve) => {
    if (authStore.authStateChanged) {
      resolve()
      return
    }
    const unwatch = watch(
      () => authStore.authStateChanged,
      (value) => {
        if (value) {
          unwatch()
          resolve()
        }
      },
    )
  })

  try {
    if (!authStore.user && authStore.authUser) {
      await authStore.loadUser()
    }

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
