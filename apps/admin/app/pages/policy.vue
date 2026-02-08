<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">风控策略</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索策略名称" icon="i-heroicons-magnifying-glass" class="w-60" />
          <USelect v-model="typeFilter" :items="typeOptions" class="w-28" />
          <USelect v-model="statusFilter" :items="statusOptions" class="w-28" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
          <UButton icon="i-heroicons-plus" label="添加策略" color="primary" />
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
import type { TableColumn } from '@nuxt/ui'

type PolicyRow = {
  id: string
  name: string
  type: 'rate_limit' | 'ip_allowlist' | 'ip_blocklist' | 'budget'
  status: 'active' | 'disabled'
  threshold: string
  updated_at: string
}

const rawItems = ref<PolicyRow[]>([
  {
    id: 'policy-01',
    name: '默认速率限制',
    type: 'rate_limit',
    status: 'active',
    threshold: '60 次/分钟',
    updated_at: '2025-02-01T10:20:00Z'
  },
  {
    id: 'policy-02',
    name: '高风险 IP 黑名单',
    type: 'ip_blocklist',
    status: 'active',
    threshold: '12 个 IP',
    updated_at: '2025-01-20T08:10:00Z'
  },
  {
    id: 'policy-03',
    name: '企业白名单',
    type: 'ip_allowlist',
    status: 'disabled',
    threshold: '6 个 IP',
    updated_at: '2025-01-15T09:30:00Z'
  },
  {
    id: 'policy-04',
    name: '组织预算上限',
    type: 'budget',
    status: 'active',
    threshold: '¥5,000 / 月',
    updated_at: '2025-01-05T17:05:00Z'
  }
])

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const typeFilter = ref<'all' | PolicyRow['type']>('all')
const statusFilter = ref<'all' | PolicyRow['status']>('all')

const typeOptions = [
  { label: '全部类型', value: 'all' },
  { label: '速率限制', value: 'rate_limit' },
  { label: 'IP 白名单', value: 'ip_allowlist' },
  { label: 'IP 黑名单', value: 'ip_blocklist' },
  { label: '预算上限', value: 'budget' }
]
const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'disabled' }
]
const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

const filteredItems = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return rawItems.value.filter((item) => {
    const matchesKeyword = !keyword || item.name.toLowerCase().includes(keyword)
    const matchesType = typeFilter.value === 'all' || item.type === typeFilter.value
    const matchesStatus = statusFilter.value === 'all' || item.status === statusFilter.value
    return matchesKeyword && matchesType && matchesStatus
  })
})

const totalCount = computed(() => filteredItems.value.length)

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch([searchTerm, typeFilter, statusFilter, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')

const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const columns = computed<TableColumn<PolicyRow>[]>(() => [
  {
    accessorKey: 'id',
    header: '策略ID',
    meta: { class: { th: 'w-28' } }
  },
  {
    accessorKey: 'name',
    header: '策略名称'
  },
  {
    accessorKey: 'type',
    header: '类型',
    cell: ({ row }) => {
      const typeValue = String(row.getValue('type') || '')
      const label =
        typeValue === 'rate_limit'
          ? '速率限制'
          : typeValue === 'ip_allowlist'
            ? 'IP 白名单'
            : typeValue === 'ip_blocklist'
              ? 'IP 黑名单'
              : '预算上限'
      return h(UBadge, { color: 'neutral', variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => {
      const statusValue = String(row.getValue('status') || '')
      const color = statusValue === 'active' ? 'success' : 'error'
      const label = statusValue === 'active' ? '启用' : '禁用'
      return h(UBadge, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'threshold',
    header: '阈值/参数'
  },
  {
    accessorKey: 'updated_at',
    header: '更新时间',
    cell: ({ row }) => formatTime(String(row.getValue('updated_at')))
  }
])
</script>
