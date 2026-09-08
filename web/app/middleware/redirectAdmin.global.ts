export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    return
  }

  // Allow admin routes and login
  if (to.path.startsWith('/admin') || to.name === 'login') {
    return
  }

  // Only redirect from customer dashboard app routes
  const customerAppPrefixes = [
    '/threads',
    '/messages',
    '/contacts',
    '/settings',
    '/billing',
    '/bulk-messages',
    '/heartbeats',
    '/phone-api-keys',
    '/search-messages',
  ]

  const isCustomerRoute = customerAppPrefixes.some((prefix) =>
    to.path.startsWith(prefix),
  )

  if (!isCustomerRoute) {
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
