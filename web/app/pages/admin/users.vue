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
              @update:model-value="fetchUsers"
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
            <template #[`item.is_admin`]="{ item }">
              <v-chip
                :color="item.is_admin ? 'success' : 'default'"
                size="small"
              >
                {{ item.is_admin ? 'Admin' : 'Usuário' }}
              </v-chip>
            </template>
            <template #[`item.created_at`]="{ item }">
              {{ new Date(item.created_at).toLocaleDateString() }}
            </template>
          </v-data-table-server>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'
import type { EntitiesUser } from '~/shared/types/api'
import { useNuxtApp } from '#app'

const authStore = useAuthStore()
const router = useRouter()
const { $api } = useNuxtApp()

const users = ref<EntitiesUser[]>([])
const totalItems = ref(0)
const loading = ref(false)
const search = ref('')
const itemsPerPage = ref(20)

const headers = [
  { title: 'Email', key: 'email', sortable: false },
  { title: 'Telefone', key: 'phone', sortable: false },
  { title: 'Assinatura', key: 'subscription_name', sortable: false },
  { title: 'Status', key: 'subscription_status', sortable: false },
  { title: 'Papel', key: 'is_admin', sortable: false },
  { title: 'Data Cadastro', key: 'created_at', sortable: false },
]

onMounted(() => {
  if (!authStore.user?.is_admin) {
    router.push('/')
  }
})

const fetchUsers = async (options?: { page: number; itemsPerPage: number }) => {
  if (!authStore.user?.is_admin) return

  loading.value = true

  const page = options?.page || 1
  const limit = options?.itemsPerPage || itemsPerPage.value
  const skip = (page - 1) * limit

  try {
    const response = await $api.v1.adminUsersList({
      skip,
      limit,
      query: search.value,
    })

    users.value = response.data.data?.items || []
    totalItems.value = response.data.data?.total_count || 0
  } catch (error) {
    console.error('Failed to fetch users', error)
  } finally {
    loading.value = false
  }
}
</script>
