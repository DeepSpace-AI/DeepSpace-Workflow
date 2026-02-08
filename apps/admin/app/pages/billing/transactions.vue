<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">交易流水</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索引用ID或用户ID" icon="i-heroicons-magnifying-glass" class="w-64" />
          <USelect v-model="typeFilter" :items="typeOptions" class="w-28" />
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
import { computed, h, ref, resolveComponent, watch } from 'vue'
import type { TableColumn } from '@nuxt/ui'

type TransactionRow = {
  id: number
  user_id: number
  type: 'hold' | 'capture' | 'release'
  amount: number
  ref_id: string
  created_at: string
}

const rawItems = ref<TransactionRow[]>([
  {
    id: 90231,
    user_id: 1001,
    type: 'hold',
    amount: 18.5,
    ref_id: 'req_20250201_001',
    created_at: '2025-02-01T10:22:10Z'
  },
  {
    id: 90232,
    user_id: 1001,
    type: 'capture',
    amount: 16.8,
    ref_id: 'req_20250201_001',
    created_at: '2025-02-01T10:22:30Z'
  },
  {
    id: 90233,
    user_id: 1003,
    type: 'hold',
    amount: 6.2,
    ref_id: 'req_20250115_002',
    created_at: '2025-01-15T09:31:10Z'
  },
  {
    id: 90234,
    user_id: 1003,
    type: 'release',
    amount: 1.1,
    ref_id: 'req_20250115_002',
    created_at: '2025-01-15T09:31:50Z'
  },
  {
    id: 90235,
    user_id: 1002,
    type: 'capture',
    amount: 12.6,
    ref_id: 'req_20250120_009',
    created_at: '2025-01-20T08:12:08Z'
  }
])

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const typeFilter = ref<'all' | TransactionRow['type']>('all')

const typeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '预扣', value: 'hold' },
  { label: '扣款', value: 'capture' },
  { label: '释放', value: 'release' }
]
const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

const filteredItems = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return rawItems.value.filter((item) => {
    const matchesKeyword =
      !keyword ||
      item.ref_id.toLowerCase().includes(keyword) ||
      String(item.user_id).includes(keyword)
    const matchesType = typeFilter.value === 'all' || item.type === typeFilter.value
    return matchesKeyword && matchesType
  })
})

const totalCount = computed(() => filteredItems.value.length)

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch([searchTerm, typeFilter, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')

const formatAmount = (value: number) => value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const columns = computed<TableColumn<TransactionRow>[]>(() => [
  {
    accessorKey: 'id',
    header: '流水ID',
    meta: { class: { th: 'w-28' } }
  },
  {
    accessorKey: 'user_id',
    header: '用户ID'
  },
  {
    accessorKey: 'type',
    header: '类型',
    cell: ({ row }) => {
      const typeValue = String(row.getValue('type') || '')
      const color = typeValue === 'hold' ? 'warning' : typeValue === 'capture' ? 'success' : 'neutral'
      const label = typeValue === 'hold' ? '预扣' : typeValue === 'capture' ? '扣款' : '释放'
      return h(UBadge, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'amount',
    header: '金额',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => formatAmount(Number(row.getValue('amount')))
  },
  {
    accessorKey: 'ref_id',
    header: '引用ID'
  },
  {
    accessorKey: 'created_at',
    header: '创建时间',
    cell: ({ row }) => formatTime(String(row.getValue('created_at')))
  }
])
</script>
