<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">钱包管理</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索组织ID" icon="i-heroicons-magnifying-glass" class="w-60" />
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

type WalletRow = {
  org_id: number
  balance: number
  frozen_balance: number
  updated_at: string
}

const rawItems = ref<WalletRow[]>([
  {
    org_id: 1001,
    balance: 1280.5,
    frozen_balance: 120,
    updated_at: '2025-02-01T10:20:00Z'
  },
  {
    org_id: 1002,
    balance: 320.12,
    frozen_balance: 0,
    updated_at: '2025-01-20T08:10:00Z'
  },
  {
    org_id: 1003,
    balance: 58.9,
    frozen_balance: 12.5,
    updated_at: '2025-01-15T09:30:00Z'
  },
  {
    org_id: 1004,
    balance: 860.75,
    frozen_balance: 42,
    updated_at: '2025-01-05T17:05:00Z'
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
  const keyword = searchTerm.value.trim()
  return rawItems.value.filter((item) => {
    return !keyword || String(item.org_id).includes(keyword)
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

const columns = computed<TableColumn<WalletRow>[]>(() => [
  {
    accessorKey: 'org_id',
    header: '组织ID',
    meta: { class: { th: 'w-28' } }
  },
  {
    accessorKey: 'balance',
    header: '余额',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => formatAmount(Number(row.getValue('balance')))
  },
  {
    accessorKey: 'frozen_balance',
    header: '冻结余额',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => formatAmount(Number(row.getValue('frozen_balance')))
  },
  {
    accessorKey: 'updated_at',
    header: '更新时间',
    cell: ({ row }) => formatTime(String(row.getValue('updated_at')))
  }
])
</script>
