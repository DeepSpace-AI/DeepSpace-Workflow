<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">交易流水</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索用户ID" icon="i-heroicons-magnifying-glass" class="w-64" />
          <USelect v-model="typeFilter" :items="typeOptions" class="w-28" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="transactionItems" :columns="columns" :loading="isLoading" />
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
  type: 'hold' | 'capture' | 'release' | 'topup'
  amount: number
  ref_id: string
  created_at: string
}

type TransactionListResponse = {
  items: TransactionRow[]
  total: number
  page: number
  page_size: number
}

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const typeFilter = ref<'all' | TransactionRow['type']>('all')

const typeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '预扣', value: 'hold' },
  { label: '扣款', value: 'capture' },
  { label: '释放', value: 'release' },
  { label: '充值', value: 'topup' }
]
const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

watch([searchTerm, typeFilter, pageSize], () => {
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
  user_id: userIdQuery.value,
  type: typeFilter.value === 'all' ? undefined : typeFilter.value
}))

const { data: listData, pending: isLoading } = await useFetch<TransactionListResponse | { data?: TransactionListResponse }>(
  '/api/admin/billing/transactions',
  {
  query: queryParams,
  default: () => ({
    items: [],
    total: 0,
    page: 1,
    page_size: pageSize.value
  })
  }
)

const rawTransactionItems = computed(() => {
  const dataValue = listData.value
  if (!dataValue) return []
  if ('data' in dataValue && Array.isArray(dataValue.data?.items)) return dataValue.data.items
  if ('items' in dataValue && Array.isArray(dataValue.items)) return dataValue.items
  return []
})

const normalizeTransactionRow = (item: Record<string, unknown>): TransactionRow => {
  return {
    id: Number(item.id ?? item.ID ?? 0),
    user_id: Number(item.user_id ?? item.UserID ?? 0),
    type: String(item.type ?? item.Type ?? 'release') as TransactionRow['type'],
    amount: Number(item.amount ?? item.Amount ?? 0),
    ref_id: String(item.ref_id ?? item.RefID ?? ''),
    created_at: String(item.created_at ?? item.CreatedAt ?? '')
  }
}

const transactionItems = computed(() => rawTransactionItems.value.map(normalizeTransactionRow))
const totalCount = computed(() => {
  const dataValue = listData.value
  if (!dataValue) return 0
  if ('total' in dataValue && typeof dataValue.total === 'number') return dataValue.total
  if ('data' in dataValue && typeof dataValue.data?.total === 'number') return dataValue.data.total
  return 0
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
      const color = typeValue === 'hold'
        ? 'warning'
        : typeValue === 'capture'
          ? 'success'
          : typeValue === 'topup'
            ? 'primary'
            : 'neutral'
      const label = typeValue === 'hold'
        ? '预扣'
        : typeValue === 'capture'
          ? '扣款'
          : typeValue === 'release'
            ? '释放'
            : typeValue === 'topup'
              ? '充值'
              : typeValue || '未知'
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
    header: '幂等ID'
  },
  {
    accessorKey: 'created_at',
    header: '创建时间',
    cell: ({ row }) => formatTime(String(row.getValue('created_at')))
  }
])
</script>
