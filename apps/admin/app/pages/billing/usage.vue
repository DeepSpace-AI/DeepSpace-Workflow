<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">用量记录</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索用户ID" icon="i-heroicons-magnifying-glass" class="w-64" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="usageItems" :columns="columns" :loading="isLoading" />
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-500">共 {{ totalCount }} 条</div>
        <UPagination v-model:page="page" :total="totalCount" :items-per-page="pageSize" :sibling-count="1" show-edges />
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { computed, ref, resolveComponent, watch } from 'vue'
import type { TableColumn } from '@nuxt/ui'

type UsageRow = {
  user_id: number
  model: string
  total_tokens: number
  cost: number
  trace_id: string
  created_at: string
}

type UsageListResponse = {
  items: UsageRow[]
  total: number
  page: number
  page_size: number
}

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')

const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

watch([searchTerm, pageSize], () => {
  page.value = 1
})

const userIdQuery = computed(() => {
  const keyword = searchTerm.value.trim()
  if (!keyword) return undefined
  const parsed = Number(keyword)
  if (!Number.isInteger(parsed) || parsed <= 0) return undefined
  return parsed
})

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  user_id: userIdQuery.value
}))

const { data: listData, pending: isLoading } = await useFetch<UsageListResponse>('/api/admin/billing/usage', {
  query: queryParams,
  default: () => ({
    items: [],
    total: 0,
    page: 1,
    page_size: pageSize.value
  })
})

const usageItems = computed(() => listData.value?.items ?? [])
const totalCount = computed(() => listData.value?.total ?? 0)

const UBadge = resolveComponent('UBadge')

const formatAmount = (value: number) => value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const columns = computed<TableColumn<UsageRow>[]>(() => [
  {
    accessorKey: 'user_id',
    header: '用户ID',
    meta: { class: { th: 'w-28' } }
  },
  {
    accessorKey: 'model',
    header: '模型'
  },
  {
    accessorKey: 'total_tokens',
    header: '总 Tokens',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => Number(row.getValue('total_tokens')).toLocaleString('zh-CN')
  },
  {
    accessorKey: 'cost',
    header: '费用',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => `$${formatAmount(Number(row.getValue('cost')))}`
  },
  {
    accessorKey: 'trace_id',
    header: 'Trace ID'
  },
  {
    accessorKey: 'created_at',
    header: '创建时间',
    cell: ({ row }) => formatTime(String(row.getValue('created_at')))
  }
])
</script>
