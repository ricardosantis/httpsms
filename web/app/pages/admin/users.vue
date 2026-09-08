<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { mdiMagnify } from '@mdi/js'
import type {
  EntitiesUser,
  ResponsesUserListResponse,
} from '~~/shared/types/api'
import { getApiErrorMessage } from '~/utils/api-error'

definePageMeta({
  layout: 'admin',
  middleware: ['auth', 'admin'],
})

const { t, locale } = useI18n()
useHead({
  title: computed(() => `${t('admin.usersPageTitle')} - httpSMS`),
})

const authStore = useAuthStore()
const notificationsStore = useNotificationsStore()
const { apiFetch } = useApi()

const users = ref<EntitiesUser[]>([])
const totalItems = ref(0)
const loading = ref(false)
const search = ref('')
const page = ref(1)
const itemsPerPage = ref(20)
const actionLoading = ref<string | null>(null)

let searchTimer: ReturnType<typeof setTimeout> | undefined

const headers = computed(() => [
  { title: t('admin.headers.email'), key: 'email', sortable: false },
  {
    title: t('admin.headers.subscription'),
    key: 'subscription_name',
    sortable: false,
  },
  { title: t('admin.headers.status'), key: 'active', sortable: false },
  {
    title: t('admin.headers.subStatus'),
    key: 'subscription_status',
    sortable: false,
  },
  { title: t('admin.headers.createdAt'), key: 'created_at', sortable: false },
  {
    title: t('admin.headers.actions'),
    key: 'actions',
    sortable: false,
    align: 'end' as const,
  },
])

const fetchUsers = async (options?: { page: number; itemsPerPage: number }) => {
  if (!authStore.user?.is_admin) return

  loading.value = true

  const currentPage = options?.page || page.value
  const limit = options?.itemsPerPage || itemsPerPage.value
  const skip = (currentPage - 1) * limit

  try {
    const response = await apiFetch<ResponsesUserListResponse>(
      '/v1/admin/users',
      {
        params: { skip, limit, query: search.value },
      },
    )

    users.value = response.data.items || []
    totalItems.value = response.data.total_count || 0
  } catch (error) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(error, t('admin.errorGeneric')),
    })
  } finally {
    loading.value = false
  }
}

watch(search, () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  searchTimer = setTimeout(() => {
    if (page.value !== 1) {
      page.value = 1
    } else {
      fetchUsers()
    }
  }, 350)
})

onBeforeUnmount(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})

const toggleBlock = async (item: EntitiesUser) => {
  actionLoading.value = item.id
  try {
    if (item.active !== false) {
      await apiFetch(`/v1/admin/users/${item.id}/block`, { method: 'POST' })
      item.active = false
      notificationsStore.addNotification({
        type: 'success',
        message: t('admin.blockSuccess'),
      })
    } else {
      await apiFetch(`/v1/admin/users/${item.id}/unblock`, { method: 'POST' })
      item.active = true
      notificationsStore.addNotification({
        type: 'success',
        message: t('admin.unblockSuccess'),
      })
    }
  } catch (error) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(error, t('admin.errorGeneric')),
    })
  } finally {
    actionLoading.value = null
  }
}

const confirmDelete = async (item: EntitiesUser) => {
  if (!confirm(t('admin.confirmDelete', { email: item.email }))) {
    return
  }

  actionLoading.value = item.id
  try {
    await apiFetch(`/v1/admin/users/${item.id}`, { method: 'DELETE' })
    users.value = users.value.filter((u) => u.id !== item.id)
    totalItems.value = Math.max(0, totalItems.value - 1)
    notificationsStore.addNotification({
      type: 'success',
      message: t('admin.deleteSuccess'),
    })
  } catch (error) {
    notificationsStore.addNotification({
      type: 'error',
      message: getApiErrorMessage(error, t('admin.errorGeneric')),
    })
  } finally {
    actionLoading.value = null
  }
}
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="d-flex align-center justify-space-between">
        <h1 class="text-h4">{{ t('admin.usersPageTitle') }}</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-text-field
              v-model="search"
              :append-inner-icon="mdiMagnify"
              :label="t('admin.searchPlaceholder')"
              single-line
              hide-details
              clearable
            />
          </v-card-title>

          <v-data-table-server
            v-model:page="page"
            v-model:items-per-page="itemsPerPage"
            :headers="headers"
            :items="users"
            :items-length="totalItems"
            :loading="loading"
            @update:options="fetchUsers"
          >
            <template #[`item.active`]="{ item }">
              <v-chip
                :color="item.active !== false ? 'success' : 'error'"
                size="small"
              >
                {{
                  item.active !== false
                    ? t('admin.statusActive')
                    : t('admin.statusBlocked')
                }}
              </v-chip>
            </template>
            <template #[`item.created_at`]="{ item }">
              {{ new Date(item.created_at).toLocaleDateString(locale) }}
            </template>
            <template #[`item.actions`]="{ item }">
              <div class="d-flex ga-2 justify-end">
                <v-btn
                  v-if="item.active !== false"
                  variant="tonal"
                  color="warning"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="toggleBlock(item)"
                >
                  {{ t('admin.block') }}
                </v-btn>
                <v-btn
                  v-else
                  variant="tonal"
                  color="success"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="toggleBlock(item)"
                >
                  {{ t('admin.unblock') }}
                </v-btn>
                <v-btn
                  variant="tonal"
                  color="error"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="confirmDelete(item)"
                >
                  {{ t('admin.delete') }}
                </v-btn>
              </div>
            </template>
          </v-data-table-server>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
