<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">审计日志</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索 trace_id 或资源" icon="i-heroicons-magnifying-glass" class="w-64" />
          <USelect v-model="actionFilter" :items="actionOptions" class="w-28" />
          <USelect v-model="resultFilter" :items="resultOptions" class="w-28" />
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

type AuditRow = {
  trace_id: string
  action: 'login' | 'update' | 'delete' | 'create' | 'export'
  operator: string
  resource: string
  result: 'success' | 'failed'
  created_at: string
}

const rawItems = ref<AuditRow[]>([
  {
    trace_id: 'trace_20250201_001',
    action: 'login',
    operator: 'admin@deepspace.ai',
    resource: '登录控制台',
    result: 'success',
    created_at: '2025-02-01T10:22:30Z'
  },
  {
    trace_id: 'trace_20250120_009',
    action: 'update',
    operator: 'ops@deepspace.ai',
    resource: '定价规则 pricing-001',
    result: 'success',
    created_at: '2025-01-20T08:12:08Z'
  },
  {
    trace_id: 'trace_20250115_002',
    action: 'delete',
    operator: 'admin@deepspace.ai',
    resource: '用户 1002',
    result: 'failed',
    created_at: '2025-01-15T09:31:50Z'
  },
  {
    trace_id: 'trace_20250105_777',
    action: 'create',
    operator: 'ops@deepspace.ai',
    resource: '风控策略 policy-03',
    result: 'success',
    created_at: '2025-01-05T17:05:00Z'
  }
])

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const actionFilter = ref<'all' | AuditRow['action']>('all')
const resultFilter = ref<'all' | AuditRow['result']>('all')

const actionOptions = [
  { label: '全部动作', value: 'all' },
  { label: '登录', value: 'login' },
  { label: '新增', value: 'create' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '导出', value: 'export' }
]
const resultOptions = [
  { label: '全部结果', value: 'all' },
  { label: '成功', value: 'success' },
  { label: '失败', value: 'failed' }
]
const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

const filteredItems = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return rawItems.value.filter((item) => {
    const matchesKeyword = !keyword || item.trace_id.toLowerCase().includes(keyword) || item.resource.toLowerCase().includes(keyword)
    const matchesAction = actionFilter.value === 'all' || item.action === actionFilter.value
    const matchesResult = resultFilter.value === 'all' || item.result === resultFilter.value
    return matchesKeyword && matchesAction && matchesResult
  })
})

const totalCount = computed(() => filteredItems.value.length)

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredItems.value.slice(start, start + pageSize.value)
})

watch([searchTerm, actionFilter, resultFilter, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')

const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const columns = computed<TableColumn<AuditRow>[]>(() => [
  {
    accessorKey: 'trace_id',
    header: 'Trace ID',
    meta: { class: { th: 'w-44' } }
  },
  {
    accessorKey: 'action',
    header: '动作',
    cell: ({ row }) => {
      const actionValue = String(row.getValue('action') || '')
      const label =
        actionValue === 'login'
          ? '登录'
          : actionValue === 'create'
            ? '新增'
            : actionValue === 'update'
              ? '更新'
              : actionValue === 'delete'
                ? '删除'
                : '导出'
      return h(UBadge, { color: 'neutral', variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'operator',
    header: '操作者'
  },
  {
    accessorKey: 'resource',
    header: '资源'
  },
  {
    accessorKey: 'result',
    header: '结果',
    cell: ({ row }) => {
      const resultValue = String(row.getValue('result') || '')
      const color = resultValue === 'success' ? 'success' : 'error'
      const label = resultValue === 'success' ? '成功' : '失败'
      return h(UBadge, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'created_at',
    header: '时间',
    cell: ({ row }) => formatTime(String(row.getValue('created_at')))
  }
])
</script>
