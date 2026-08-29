<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="d-flex align-center justify-space-between">
        <h1 class="text-h4">Painel Admin - Usuários</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-text-field
              v-model="search"
              append-inner-icon="mdi-magnify"
              label="Buscar por e-mail"
              single-line
              hide-details
              @update:model-value="onSearch"
            />
          </v-card-title>

          <v-data-table-server
            v-model:items-per-page="itemsPerPage"
            :headers="headers"
            :items="users"
            :items-length="totalItems"
            :loading="loading"
            @update:options="fetchUsers"
          >
            <template #[`item.active`]="{ item }">
              <v-chip :color="item.active ? 'success' : 'error'" size="small">
                {{ item.active ? 'Ativo' : 'Bloqueado' }}
              </v-chip>
            </template>
            <template #[`item.created_at`]="{ item }">
              {{ new Date(item.created_at).toLocaleDateString() }}
            </template>
            <template #[`item.actions`]="{ item }">
              <div class="d-flex ga-2">
                <v-btn
                  v-if="item.active"
                  variant="tonal"
                  color="warning"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="toggleBlock(item)"
                >
                  Bloquear
                </v-btn>
                <v-btn
                  v-else
                  variant="tonal"
                  color="success"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="toggleBlock(item)"
                >
                  Restaurar
                </v-btn>
                <v-btn
                  variant="tonal"
                  color="error"
                  size="small"
                  :loading="actionLoading === item.id"
                  @click="confirmDelete(item)"
                >
                  Remover
                </v-btn>
              </div>
            </template>
          </v-data-table-server>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type {
  EntitiesUser,
  ResponsesUserListResponse,
} from '~~/shared/types/api'

definePageMeta({
  layout: 'admin',
  middleware: ['auth', 'admin'],
})

const authStore = useAuthStore()
const { apiFetch } = useApi()

const users = ref<EntitiesUser[]>([])
const totalItems = ref(0)
const loading = ref(false)
const search = ref('')
const itemsPerPage = ref(20)
const actionLoading = ref<string | null>(null)

const headers = [
  { title: 'Email', key: 'email', sortable: false },
  { title: 'Telefone', key: 'phone', sortable: false },
  { title: 'Assinatura', key: 'subscription_name', sortable: false },
  { title: 'Status', key: 'active', sortable: false },
  { title: 'Assin. Status', key: 'subscription_status', sortable: false },
  { title: 'Data Cadastro', key: 'created_at', sortable: false },
  { title: 'Ações', key: 'actions', sortable: false },
]

const fetchUsers = async (options?: { page: number; itemsPerPage: number }) => {
  if (!authStore.user?.is_admin) return

  loading.value = true

  const page = options?.page || 1
  const limit = options?.itemsPerPage || itemsPerPage.value
  const skip = (page - 1) * limit

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
    console.error('Failed to fetch users', error)
  } finally {
    loading.value = false
  }
}

const onSearch = () => {
  fetchUsers()
}

const toggleBlock = async (item: EntitiesUser) => {
  actionLoading.value = item.id
  try {
    if (item.active) {
      await apiFetch(`/v1/admin/users/${item.id}/block`, { method: 'POST' })
      item.active = false
    } else {
      await apiFetch(`/v1/admin/users/${item.id}/unblock`, { method: 'POST' })
      item.active = true
    }
  } catch (error) {
    console.error('Failed to toggle block', error)
  } finally {
    actionLoading.value = null
  }
}

const confirmDelete = async (item: EntitiesUser) => {
  if (
    !confirm(
      `Tem certeza que deseja remover ${item.email}? Esta ação é irreversível.`,
    )
  ) {
    return
  }

  actionLoading.value = item.id
  try {
    await apiFetch(`/v1/admin/users/${item.id}`, { method: 'DELETE' })
    users.value = users.value.filter((u) => u.id !== item.id)
    totalItems.value = Math.max(0, totalItems.value - 1)
  } catch (error) {
    console.error('Failed to delete user', error)
  } finally {
    actionLoading.value = null
  }
}
</script>
