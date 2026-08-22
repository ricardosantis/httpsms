import { defineStore } from 'pinia'
import type {
  EntitiesBillingUsage,
  EntitiesUser,
  EntitiesWebhook,
  EntitiesDiscord,
  EntitiesMessageSendSchedule,
  EntitiesPhoneAPIKey,
  RequestsWebhookStore,
  RequestsWebhookUpdate,
  RequestsDiscordStore,
  RequestsDiscordUpdate,
  RequestsMessageSendScheduleStore,
  RequestsUserNotificationUpdate,
  RequestsUserPaymentInvoice,
  ResponsesUserSubscriptionPaymentsResponse,
} from '~~/shared/types/api'
import { getApiErrorMessage } from '~/utils/api-error'

export const useBillingStore = defineStore('billing', () => {
  const billingUsage = ref<EntitiesBillingUsage | null>(null)
  const billingUsageHistory = ref<EntitiesBillingUsage[]>([])
  const { apiFetch } = useApi()
  const notificationsStore = useNotificationsStore()

  async function loadBillingUsage() {
    try {
      const response = await apiFetch<{ data: EntitiesBillingUsage }>(
        '/v1/billing/usage',
      )
      billingUsage.value = response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to load billing usage'),
        type: 'error',
      })
      throw error
    }
  }

  async function loadBillingUsageHistory() {
    try {
      const response = await apiFetch<{ data: EntitiesBillingUsage[] }>(
        '/v1/billing/usage-history',
      )
      billingUsageHistory.value = response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to load billing usage history',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function getSubscriptionUpdateLink(): Promise<string> {
    try {
      const response = await apiFetch<{ data: { url: string } }>(
        '/v1/stripe/customer-portal',
        { method: 'POST' },
      )
      return response.data.url
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to get subscription management link',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function createStripeCheckoutSession(
    planId: string = 'pro-monthly',
    priceId?: string,
  ): Promise<string> {
    try {
      const response = await apiFetch<{ data: { url: string } }>(
        '/v1/stripe/checkout-session',
        {
          method: 'POST',
          body: {
            plan_id: planId,
            ...(priceId ? { price_id: priceId } : {}),
          },
        },
      )
      return response.data.url
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to create checkout session'),
        type: 'error',
      })
      throw error
    }
  }

  async function cancelSubscription(): Promise<string> {
    try {
      const response = await apiFetch<{ message: string }>(
        '/v1/users/subscription',
        {
          method: 'DELETE',
        },
      )
      return response.message
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to cancel subscription'),
        type: 'error',
      })
      throw error
    }
  }

  async function indexSubscriptionPayments(): Promise<ResponsesUserSubscriptionPaymentsResponse> {
    try {
      const response =
        await apiFetch<ResponsesUserSubscriptionPaymentsResponse>(
          '/v1/users/subscription/payments',
          { params: { limit: 100 } },
        )
      return response
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to load subscription payments',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function generateSubscriptionPaymentInvoice(
    subscriptionInvoiceId: string,
    request: RequestsUserPaymentInvoice,
  ): Promise<void> {
    try {
      const response = await apiFetch(
        `/v1/users/subscription/invoices/${subscriptionInvoiceId}`,
        {
          method: 'POST',
          body: request,
          responseType: 'blob',
        },
      )

      const pdfBlob = new Blob([response as Blob], {
        type: 'application/pdf',
      })
      const url = window.URL.createObjectURL(pdfBlob)
      const tempLink = document.createElement('a')
      tempLink.href = url
      tempLink.setAttribute('download', 'Invoice.pdf')
      document.body.appendChild(tempLink)
      tempLink.click()
      document.body.removeChild(tempLink)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to generate invoice'),
        type: 'error',
      })
      throw error
    }
  }

  // Webhooks
  async function createWebhook(
    payload: RequestsWebhookStore,
  ): Promise<EntitiesWebhook> {
    try {
      const response = await apiFetch<{ data: EntitiesWebhook }>(
        '/v1/webhooks',
        {
          method: 'POST',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to create webhook'),
        type: 'error',
      })
      throw error
    }
  }

  async function getWebhooks(): Promise<EntitiesWebhook[]> {
    try {
      const response = await apiFetch<{ data: EntitiesWebhook[] }>(
        '/v1/webhooks',
        {
          params: { limit: 100 },
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to load webhooks'),
        type: 'error',
      })
      throw error
    }
  }

  async function updateWebhook(
    payload: RequestsWebhookUpdate & { id: string },
  ): Promise<EntitiesWebhook> {
    try {
      const response = await apiFetch<{ data: EntitiesWebhook }>(
        `/v1/webhooks/${payload.id}`,
        {
          method: 'PUT',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to update webhook'),
        type: 'error',
      })
      throw error
    }
  }

  async function deleteWebhook(id: string): Promise<void> {
    try {
      await apiFetch(`/v1/webhooks/${id}`, { method: 'DELETE' })
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to delete webhook'),
        type: 'error',
      })
      throw error
    }
  }

  // Discord
  async function createDiscord(
    payload: RequestsDiscordStore,
  ): Promise<EntitiesDiscord> {
    try {
      const response = await apiFetch<{ data: EntitiesDiscord }>(
        '/v1/discord-integrations',
        {
          method: 'POST',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to create Discord integration',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function getDiscordIntegrations(): Promise<EntitiesDiscord[]> {
    try {
      const response = await apiFetch<{ data: EntitiesDiscord[] }>(
        '/v1/discord-integrations',
        {
          params: { limit: 100 },
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to load Discord integrations',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function updateDiscordIntegration(
    payload: RequestsDiscordUpdate & { id: string },
  ): Promise<EntitiesDiscord> {
    try {
      const response = await apiFetch<{ data: EntitiesDiscord }>(
        `/v1/discord-integrations/${payload.id}`,
        {
          method: 'PUT',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to update Discord integration',
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function deleteDiscordIntegration(id: string): Promise<void> {
    try {
      await apiFetch(`/v1/discord-integrations/${id}`, { method: 'DELETE' })
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to delete Discord integration',
        ),
        type: 'error',
      })
      throw error
    }
  }

  // Send Schedules
  async function getSendSchedules(): Promise<EntitiesMessageSendSchedule[]> {
    try {
      const response = await apiFetch<{ data: EntitiesMessageSendSchedule[] }>(
        '/v1/send-schedules',
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to load send schedules'),
        type: 'error',
      })
      throw error
    }
  }

  async function createSendSchedule(
    payload: RequestsMessageSendScheduleStore,
  ): Promise<EntitiesMessageSendSchedule> {
    try {
      const response = await apiFetch<{ data: EntitiesMessageSendSchedule }>(
        '/v1/send-schedules',
        {
          method: 'POST',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to create send schedule'),
        type: 'error',
      })
      throw error
    }
  }

  async function updateSendSchedule(
    payload: RequestsMessageSendScheduleStore & { id: string },
  ): Promise<EntitiesMessageSendSchedule> {
    try {
      const response = await apiFetch<{ data: EntitiesMessageSendSchedule }>(
        `/v1/send-schedules/${payload.id}`,
        {
          method: 'PUT',
          body: payload,
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to update send schedule'),
        type: 'error',
      })
      throw error
    }
  }

  async function deleteSendSchedule(id: string): Promise<void> {
    try {
      await apiFetch(`/v1/send-schedules/${id}`, { method: 'DELETE' })
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to delete send schedule'),
        type: 'error',
      })
      throw error
    }
  }

  // Phone API Keys
  async function storePhoneApiKey(name: string): Promise<EntitiesPhoneAPIKey> {
    try {
      const response = await apiFetch<{
        data: EntitiesPhoneAPIKey
        message: string
      }>('/v1/phone-api-keys', {
        method: 'POST',
        body: { name },
      })
      notificationsStore.addNotification({
        message: response.message,
        type: 'success',
      })
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to create phone API key'),
        type: 'error',
      })
      throw error
    }
  }

  async function indexPhoneApiKeys(): Promise<EntitiesPhoneAPIKey[]> {
    try {
      const response = await apiFetch<{ data: EntitiesPhoneAPIKey[] }>(
        '/v1/phone-api-keys',
        {
          params: { limit: 100 },
        },
      )
      return response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to load phone API keys'),
        type: 'error',
      })
      throw error
    }
  }

  async function deletePhoneApiKey(id: string): Promise<void> {
    try {
      const response = await apiFetch<{ message: string }>(
        `/v1/phone-api-keys/${id}`,
        { method: 'DELETE' },
      )
      notificationsStore.addNotification({
        message: response.message,
        type: 'success',
      })
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to delete phone API key'),
        type: 'error',
      })
      throw error
    }
  }

  async function deletePhoneFromPhoneApiKey(
    phoneApiKeyId: string,
    phoneId: string,
  ): Promise<void> {
    try {
      const response = await apiFetch<{ message: string }>(
        `/v1/phone-api-keys/${phoneApiKeyId}/phones/${phoneId}`,
        { method: 'DELETE' },
      )
      notificationsStore.addNotification({
        message: response.message,
        type: 'success',
      })
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to remove phone from API key',
        ),
        type: 'error',
      })
      throw error
    }
  }

  // Email notifications
  async function saveEmailNotifications(
    userId: string,
    payload: RequestsUserNotificationUpdate,
  ): Promise<void> {
    try {
      const authStore = useAuthStore()
      const response = await apiFetch<{ data: EntitiesUser }>(
        `/v1/users/${userId}/notifications`,
        {
          method: 'PUT',
          body: payload,
        },
      )
      authStore.user = response.data
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          'Failed to save notification settings',
        ),
        type: 'error',
      })
      throw error
    }
  }

  return {
    billingUsage,
    billingUsageHistory,
    loadBillingUsage,
    loadBillingUsageHistory,
    getSubscriptionUpdateLink,
    createStripeCheckoutSession,
    cancelSubscription,
    indexSubscriptionPayments,
    generateSubscriptionPaymentInvoice,
    createWebhook,
    getWebhooks,
    updateWebhook,
    deleteWebhook,
    createDiscord,
    getDiscordIntegrations,
    updateDiscordIntegration,
    deleteDiscordIntegration,
    getSendSchedules,
    createSendSchedule,
    updateSendSchedule,
    deleteSendSchedule,
    storePhoneApiKey,
    indexPhoneApiKeys,
    deletePhoneApiKey,
    deletePhoneFromPhoneApiKey,
    saveEmailNotifications,
  }
})
