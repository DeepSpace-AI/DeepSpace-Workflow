<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">用量记录</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索模型或 trace_id" icon="i-heroicons-magnifying-glass" class="w-64" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="pagedItems" :columns="columns" />
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
  model: string
  total_tokens: number
  cost: number
  trace_id: string
  created_at: string
}

const rawItems = ref<UsageRow[]>([
  {
    model: 'gpt-4.1',
    total_tokens: 15800,
    cost: 3.82,
    trace_id: 'trace_20250201_001',
    created_at: '2025-02-01T10:22:30Z'
  },
  {
    model: 'gpt-4o-mini',
    total_tokens: 6200,
    cost: 0.94,
    trace_id: 'trace_20250120_009',
    created_at: '2025-01-20T08:12:08Z'
  },
  {
    model: 'deepseek-r1',
    total_tokens: 4200,
    cost: 0.51,
    trace_id: 'trace_20250115_002',
    created_at: '2025-01-15T09:31:50Z'
  },
  {
    model: 'qwen2.5-72b',
    total_tokens: 9600,
    cost: 2.12,
    trace_id: 'trace_20250105_777',
    created_at: '2025-01-05T17:05:00Z'
  }
])

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')

const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

const filteredItems = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return rawItems.value.filter((item) => {
    const matchesKeyword = !keyword || item.model.toLowerCase().includes(keyword) || item.trace_id.toLowerCase().includes(keyword)
    return matchesKeyword
  })
})

const totalCount = computed(() => filteredItems.value.length)

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch([searchTerm, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')

const formatAmount = (value: number) => value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const columns = computed<TableColumn<UsageRow>[]>(() => [
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
