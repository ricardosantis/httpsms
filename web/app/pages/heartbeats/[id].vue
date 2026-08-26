<script setup lang="ts">
import { mdiArrowLeft, mdiCircle } from '@mdi/js'
import {
  formatDuration as formatDurationFns,
  intervalToDuration,
} from 'date-fns'
import { Bar } from 'vue-chartjs'
import type { ChartData, ChartOptions } from 'chart.js'
import type { EntitiesHeartbeat } from '~~/shared/types/api'

definePageMeta({
  middleware: ['auth'],
})

const { t } = useI18n()

useHead({
  title: computed(() => `${t('heartbeats.title')} - httpSMS`),
})

const route = useRoute()
const { mdAndDown, mdAndUp, lgAndUp } = useVDisplay()
const authStore = useAuthStore()
const phonesStore = usePhonesStore()
const { formatPhoneNumber, formatTimestamp } = useFilters()

const loading = ref(true)
const heartbeats = ref<EntitiesHeartbeat[]>([])
const phoneId = computed(() => route.params.id as string)

interface HeartbeatTableItem {
  id: string
  owner: string
  timestamp: string
  interval: number
}

const dataTableHeaders = computed(() => [
  { title: t('heartbeats.id'), key: 'id', sortable: false },
  { title: t('heartbeats.phoneNumber'), key: 'owner', sortable: false },
  { title: t('heartbeats.receivedAt'), key: 'timestamp' },
  { title: t('heartbeats.timeInterval'), key: 'interval' },
])

function getDiff(a: string, b: string): number {
  try {
    const da = new Date(a).getTime()
    const db = new Date(b).getTime()
    if (isNaN(da) || isNaN(db)) return 0
    return da - db
  } catch {
    return 0
  }
}

function formatInterval(duration: number): string {
  if (!duration || duration <= 0) {
    return '-'
  }
  try {
    const start = new Date()
    start.setMilliseconds(start.getMilliseconds() + duration)
    return (
      formatDurationFns(
        intervalToDuration({ start: new Date(), end: start }),
      ) || '0 seconds'
    )
  } catch {
    return '-'
  }
}

const dataTableItems = computed<HeartbeatTableItem[]>(() => {
  return heartbeats.value.map((heartbeat, index) => {
    let interval = 0
    if (index < heartbeats.value.length - 1 && heartbeats.value[index + 1]) {
      interval = getDiff(
        heartbeat.timestamp,
        heartbeats.value[index + 1]!.timestamp,
      )
    }
    return {
      id: heartbeat.id,
      timestamp: heartbeat.timestamp,
      owner: heartbeat.owner,
      interval,
    }
  })
})

const chartData = computed<ChartData<'bar'>>(() => {
  const data: Array<{ x: string; y: number }> = []
  for (const heartbeat of heartbeats.value) {
    try {
      const d = new Date(heartbeat.timestamp)
      if (!isNaN(d.getTime())) {
        data.push({ x: d.toISOString(), y: 1 })
      }
    } catch {
      // ignore
    }
  }

  if (!data.length) {
    return {
      datasets: [{ data: [], backgroundColor: '#2196f3' }],
    } as unknown as ChartData<'bar'>
  }

  let prev = new Date(data[0]!.x)
  const newData = [] as Array<{ x: string; y: number }>
  for (let i = 1; i < data.length; i++) {
    const current = new Date(data[i]!.x)
    const diff = prev.getTime() - current.getTime()
    if (diff > 600000) {
      // 10 minutes
      newData.push(data[i]!)
      prev = current
    }
  }

  return {
    datasets: [{ data: newData, backgroundColor: '#2196f3' }],
  } as unknown as ChartData<'bar'>
})

const chartOptions = computed<ChartOptions<'bar'>>(() => {
  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label(context) {
            const dataset = context.dataset.data as unknown as Array<{
              x: string
            }>
            if (!dataset || context.dataIndex >= dataset.length - 1) {
              return '-'
            }
            try {
              const start = new Date(dataset[context.dataIndex + 1]!.x)
              const end = new Date(dataset[context.dataIndex]!.x)
              if (isNaN(start.getTime()) || isNaN(end.getTime())) return '-'
              const duration = intervalToDuration({ start, end })
              return formatDurationFns(duration)
            } catch {
              return '-'
            }
          },
        },
      },
    },
    scales: {
      x: { type: 'time' },
      y: { display: false },
    },
  } as ChartOptions<'bar'>
})

async function loadHeartbeats() {
  loading.value = true
  try {
    heartbeats.value = await phonesStore.getHeartbeat(100)
  } catch {
    heartbeats.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await Promise.allSettled([authStore.loadUser(), phonesStore.loadPhones()])
    if (!phonesStore.owner) {
      phonesStore.setOwner(phoneId.value)
    }
    await loadHeartbeats()
  } catch (err) {
    console.error('Error mounting heartbeats:', err)
  }
})
</script>

<template>
  <VContainer fluid class="px-0 pt-0" :class="{ 'fill-height': lgAndUp }">
    <div class="w-100 h-100">
      <VAppBar height="60" :density="mdAndDown ? 'compact' : 'default'">
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>
          {{ $t('heartbeats.title') }}
          <VIcon size="12" class="mx-2" color="primary" :icon="mdiCircle" />
          <span v-if="phonesStore.owner">{{
            formatPhoneNumber(phonesStore.owner)
          }}</span>
        </VToolbarTitle>
        <VSpacer />
        <LanguageSwitcher class="mr-2" />
      </VAppBar>
      <VContainer>
        <VRow>
          <VCol cols="12" class="mt-n4">
            <p>
              {{ $t('heartbeats.description1') }}
              <a
                href="https://dontkillmyapp.com"
                class="text-decoration-none hover:text-decoration-underline"
                target="_blank"
                >https://dontkillmyapp.com</a
              >.
            </p>
            <p>
              {{ $t('heartbeats.description2') }}
            </p>
          </VCol>
          <VCol v-if="mdAndUp" cols="12" class="px-0">
            <div class="heartbeat--chart">
              <ClientOnly>
                <Bar :data="chartData" :options="chartOptions" />
              </ClientOnly>
            </div>
          </VCol>
          <VCol cols="12">
            <p>
              {{ $t('heartbeats.tableDescription') }}
            </p>
            <VProgressLinear
              v-if="loading"
              color="primary"
              indeterminate
              class="mb-4"
            />
            <VDataTable
              v-else
              hover
              :headers="dataTableHeaders"
              :items="dataTableItems"
              :items-per-page="100"
              :sort-by="[{ key: 'timestamp', order: 'desc' }]"
              hide-default-footer
              class="heartbeat--table"
              :row-props="
                ({ item }) =>
                  item.interval > 3600000 ? { class: 'heartbeat--missed' } : {}
              "
            >
              <template #[`item.interval`]="{ item }">
                {{ formatInterval(item.interval) }}
              </template>
              <template #[`item.owner`]="{ item }">
                {{ formatPhoneNumber(item.owner) }}
              </template>
              <template #[`item.timestamp`]="{ item }">
                {{ formatTimestamp(item.timestamp) }}
              </template>
            </VDataTable>
          </VCol>
        </VRow>
      </VContainer>
    </div>
  </VContainer>
</template>

<style lang="scss">
.heartbeat--chart {
  height: 200px;
}

.v-application {
  .heartbeat--table tbody tr.heartbeat--missed {
    background: #b71c1c;
  }
}
</style>
