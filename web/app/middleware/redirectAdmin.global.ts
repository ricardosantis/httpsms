export default defineNuxtRouteMiddleware(async (to) => {
  if (to.name === 'admin-users' || to.name === 'login') {
    return
  }

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

  if (!authStore.user && authStore.authUser) {
    await authStore.loadUser()
  }

  if (authStore.user?.is_admin) {
    return navigateTo('/admin/users', { replace: true })
  }
})
