import { defineStore } from 'pinia'
import type { EntitiesMessage, EntitiesBulkMessage } from '~~/shared/types/api'
import type { SearchMessagesRequest } from '~~/shared/types/message'
import { getApiErrorMessage } from '~/utils/api-error'

export type SIM = 'SIM1' | 'SIM2' | 'DEFAULT'

export interface SendMessageRequest {
  from: string
  to: string
  content: string
  sim: SIM
  request_id?: string
}

export const useMessagesStore = defineStore('messages', () => {
  const { apiFetch } = useApi()
  const notificationsStore = useNotificationsStore()

  function t(key: string, fallback: string): string {
    try {
      return useNuxtApp().$i18n.t(key)
    } catch {
      return fallback
    }
  }

  async function sendMessage(request: SendMessageRequest) {
    try {
      const response = await apiFetch<{ message: string }>(
        '/v1/messages/send',
        {
          method: 'POST',
          body: request,
        },
      )
      notificationsStore.addNotification({
        message: response.message,
        type: 'success',
      })
      const threadsStore = useThreadsStore()
      await threadsStore.loadThreads()
    } catch (e: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          e,
          t(
            'messagesStore.errorWhileSendingMessage',
            'Error while sending message',
          ),
        ),
        type: 'error',
      })
      throw e
    }
  }

  async function deleteMessage(messageId: string) {
    try {
      await apiFetch(`/v1/messages/${messageId}`, { method: 'DELETE' })
      notificationsStore.addNotification({
        message: t(
          'messagesStore.theMessageHasBeenDeletedSuccessfully',
          'The message has been deleted successfully',
        ),
        type: 'success',
      })
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t('messagesStore.errorDeletingMessage', 'Error deleting message'),
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function searchMessages(
    payload: SearchMessagesRequest,
  ): Promise<EntitiesMessage[]> {
    try {
      const token = payload.token
      const params = { ...payload }
      delete params.token

      const response = await apiFetch<{ data: EntitiesMessage[] }>(
        '/v1/messages/search',
        {
          params,
          headers: token ? { token } : undefined,
        },
      )
      return response.data
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t('messagesStore.errorSearchingMessages', 'Error searching messages'),
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function sendBulkMessages(document: File): Promise<void> {
    try {
      const formData = new FormData()
      formData.append('document', document)
      const response = await apiFetch<{ message?: string }>(
        '/v1/bulk-messages',
        {
          method: 'POST',
          body: formData,
        },
      )
      notificationsStore.addNotification({
        message: response.message ?? 'Bulk messages sent successfully',
        type: 'success',
      })
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t(
            'messagesStore.errorSendingBulkMessages',
            'Error sending bulk messages',
          ),
        ),
        type: 'error',
      })
      throw error
    }
  }

  async function fetchBulkMessageOrders(): Promise<EntitiesBulkMessage[]> {
    try {
      const response = await apiFetch<{ data: EntitiesBulkMessage[] }>(
        '/v1/bulk-messages',
      )
      return response.data ?? []
    } catch (error: unknown) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(
          error,
          t(
            'messagesStore.errorLoadingBulkMessageOrders',
            'Error loading bulk message orders',
          ),
        ),
        type: 'error',
      })
      throw error
    }
  }

  return {
    sendMessage,
    deleteMessage,
    searchMessages,
    sendBulkMessages,
    fetchBulkMessageOrders,
  }
})
