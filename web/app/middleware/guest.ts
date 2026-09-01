export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server) {
    return
  }

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

  if (authStore.authUser !== null) {
    if (authStore.user === null) {
      await authStore.loadUser().catch(() => {})
    }
    if (authStore.user?.is_admin) {
      return navigateTo('/admin/users')
    }
    return navigateTo('/threads')
  }
})
