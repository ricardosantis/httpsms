import { useAuthStore } from '../stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()

  if (!authStore.authStateChanged) {
    await new Promise<void>((resolve) => {
      const stop = watch(
        () => authStore.authStateChanged,
        (changed) => {
          if (changed) {
            stop()
            resolve()
          }
        },
        { immediate: true },
      )
    })
  }

  if (authStore.authUser === null) {
    return navigateTo({ path: '/login', query: { to: '/admin/users' } })
  }

  if (authStore.user === null) {
    await authStore.loadUser().catch(() => {})
  }

  if (!authStore.user?.is_admin) {
    return navigateTo('/')
  }
})
