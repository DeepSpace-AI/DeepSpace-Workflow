<template>
  <div class="space-y-4">
    <UCard>
      <template #header>
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-3">
            <h3 class="text-lg font-semibold">模型管理</h3>
            <UBadge color="neutral" variant="soft">共 {{ totalCount }} 个</UBadge>
            <UBadge v-if="selectedCount" color="primary" variant="soft">已选 {{ selectedCount }} 个</UBadge>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <UInput v-model="searchTerm" placeholder="搜索模型名称或ID" icon="i-heroicons-magnifying-glass" class="w-64" />
            <USelect v-model="providerFilter" :items="providerOptions" class="w-28" />
            <USelect v-model="statusFilter" :items="statusOptions" class="w-28" />
            <UDropdownMenu :items="columnToggleItems" :content="{ align: 'end' }">
              <UButton label="列展示" color="neutral" variant="outline" trailing-icon="i-lucide-chevron-down" />
            </UDropdownMenu>
            <UButton icon="i-heroicons-rectangle-stack" label="批量管理" color="primary" variant="outline"
              :disabled="!selectedCount" @click="openBulkModal" />
            <UButton label="清空选择" color="neutral" variant="ghost" :disabled="!selectedCount" @click="clearSelection" />
            <UButton icon="i-heroicons-arrow-path" label="同步上游" color="neutral" variant="outline" :loading="isSyncing"
              @click="handleSync" />
            <UButton icon="i-heroicons-plus" label="添加模型" color="primary" @click="openCreateModal" />
          </div>
        </div>
      </template>

      <div class="flex h-[calc(100vh-12rem)] flex-col gap-4">
        <div v-if="listError" class="text-sm text-red-600">{{ listError }}</div>
        <div v-if="actionError" class="text-sm text-red-600">{{ actionError }}</div>
        <div v-if="bulkSuccess" class="text-sm text-green-600">{{ bulkSuccess }}</div>
        <div v-if="syncError" class="text-sm text-red-600">{{ syncError }}</div>
        <UTable 
        ref="table" 
        v-model:sorting="sorting" 
        v-model:column-visibility="columnVisibility"
          v-model:column-pinning="columnPinning" :sorting-options="{ manualSorting: true }" :data="pagedItems"
          :columns="columns" :loading="isLoading" sticky />
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-500">共 {{ totalCount }} 条</div>
          <div class="flex items-center gap-2">
            <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
            <UPagination 
            v-model:page="page" :total="totalCount" :items-per-page="pageSize" :sibling-count="1"
              show-edges />
          </div>
        </div>
      </div>
    </UCard>

    <UModal v-model:open="isModalOpen" title="添加模型">
      <template #body>
        <div class="grid gap-4 md:grid-cols-3">
          <UFormField label="模型名称" required class="md:col-span-3">
            <UInput v-model="formState.name" class="w-full" placeholder="请输入模型名称" />
          </UFormField>
          <UFormField label="提供方" required>
            <UInputMenu v-model="formState.provider" create-item class="w-full" :items="providerInputItems"
              placeholder="请选择或输入提供方" value-key="value" label-key="label" @create="onCreateProviderForCreate" />
          </UFormField>
          <UFormField label="状态">
            <USelect v-model="formState.status" class="w-full" :items="modelStatusOptions" placeholder="请选择状态" />
          </UFormField>
          <UFormField label="计价币种">
            <USelect v-model="formState.currency" class="w-full" :items="currencyOptions" placeholder="请选择币种" />
          </UFormField>
          <UFormField label="输入价">
            <UInput v-model="formState.priceInput" class="w-full" placeholder="例如 0.002">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M {{ formState.currency || '—' }}</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="输出价">
            <UInput v-model="formState.priceOutput" class="w-full" placeholder="例如 0.006">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M {{ formState.currency || '—' }}</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="能力枚举" class="md:col-span-3">
            <UCheckboxGroup v-model="formState.capabilities" :items="capabilityOptions"
              class="grid grid-cols-1 gap-2 sm:grid-cols-3" />
          </UFormField>
          <div v-if="formError" class="text-sm text-red-600 md:col-span-3">{{ formError }}</div>
        </div>
      </template>
      <template #footer="{ close }">
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" :disabled="isSubmitting" @click="close">取消</UButton>
          <UButton color="primary" :loading="isSubmitting" @click="submitForm">创建</UButton>
        </div>
      </template>
    </UModal>
    <UModal v-model:open="isEditOpen" title="编辑模型">
      <template #body>
        <div class="grid gap-4 md:grid-cols-3">
          <UFormField label="提供方" required>
            <UInputMenu v-model="editFormState.provider" create-item class="w-full" :items="providerInputItems"
              placeholder="请选择或输入提供方" value-key="value" label-key="label" @create="onCreateProviderForEdit" />
          </UFormField>
          <UFormField label="计价币种">
            <USelect v-model="editFormState.currency" class="w-full" :items="currencyOptions" placeholder="请选择币种" />
          </UFormField>
          <UFormField label="输入价">
            <UInput v-model="editFormState.priceInput" class="w-full" placeholder="例如 0.002">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M {{ editFormState.currency || '—' }}</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="输出价">
            <UInput v-model="editFormState.priceOutput" class="w-full" placeholder="例如 0.006">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M {{ editFormState.currency || '—' }}</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="能力枚举" class="md:col-span-3">
            <UCheckboxGroup v-model="editFormState.capabilities" :items="capabilityOptions"
              class="grid grid-cols-3 gap-2" />
          </UFormField>
          <UFormField label="提供方图标" class="md:col-span-3">
            <UInput v-model="editFormState.providerIcon" class="w-full" placeholder="请输入 Base64 字符串" />
          </UFormField>
          <div v-if="editError" class="text-sm text-red-600 md:col-span-3">{{ editError }}</div>
        </div>
      </template>
      <template #footer="{ close }">
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" :disabled="isEditSubmitting" @click="close">取消</UButton>
          <UButton color="primary" :loading="isEditSubmitting" @click="submitEdit">保存</UButton>
        </div>
      </template>
    </UModal>
    <UModal v-model:open="isBulkOpen" title="批量管理模型">
      <template #body>
        <div class="grid gap-4 md:grid-cols-2">
          <div class="text-sm text-gray-500 md:col-span-2">已选 {{ selectedCount }} 个模型</div>
          <UFormField label="更新状态">
            <USelect v-model="bulkForm.status" class="w-full" :items="bulkStatusOptions" />
          </UFormField>
          <UFormField label="更新币种">
            <USelect v-model="bulkForm.currency" class="w-full" :items="bulkCurrencyOptions" />
          </UFormField>
          <UFormField label="输入价">
            <UInput v-model="bulkForm.priceInput" class="w-full" placeholder="例如 0.002">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="输出价">
            <UInput v-model="bulkForm.priceOutput" class="w-full" placeholder="例如 0.006">
              <template #trailing>
                <span class="text-xs text-gray-500">/1M</span>
              </template>
            </UInput>
          </UFormField>
          <UFormField label="能力更新">
            <USelect v-model="bulkForm.capabilityMode" class="w-full" :items="capabilityModeOptions" />
          </UFormField>
          <UFormField label="能力标签" class="md:col-span-2">
            <UCheckboxGroup v-model="bulkForm.capabilities" :items="capabilityOptions"
              :disabled="bulkForm.capabilityMode !== 'override'" class="grid grid-cols-1 gap-2 sm:grid-cols-3" />
          </UFormField>
          <div v-if="bulkError" class="text-sm text-red-600 md:col-span-2">{{ bulkError }}</div>
        </div>
      </template>
      <template #footer="{ close }">
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" :disabled="isBulkSubmitting" @click="close">取消</UButton>
          <UButton color="primary" :loading="isBulkSubmitting" @click="submitBulk">保存</UButton>
        </div>
      </template>
    </UModal>
    <UModal v-model:open="isSyncPreviewOpen" title="同步上游模型预览" :ui="{ content: 'sm:max-w-6xl' }">
      <template #body>
        <div class="space-y-4">
          <div v-if="syncPreviewError" class="text-sm text-red-600">{{ syncPreviewError }}</div>
          <div class="flex flex-wrap items-center gap-2 text-sm">
            <UBadge color="neutral" variant="soft">同步模型 {{ upstreamItems.length }}</UBadge>
            <UBadge color="neutral" variant="soft">本地模型 {{ listItems.length }}</UBadge>
            <UBadge color="primary" variant="soft">新增 {{ addedItems.length }}</UBadge>
          </div>
          <UTable v-model:expanded="previewExpanded" :data="syncRows" :columns="previewColumns"
            :ui="{ tr: 'data-[expanded=true]:bg-elevated/50' }">
            <template #expanded="{ row }">
              <div class="grid gap-3 p-3 text-sm">
                <div class="grid gap-4 sm:grid-cols-2">
                  <div class="rounded-md border border-default p-3">
                    <div class="mb-2 text-gray-500">上游</div>
                    <div class="flex items-center gap-2">
                      <span class="text-gray-500">模型名称</span>
                      <span class="font-medium">{{ row.original.upstream?.name || '—' }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-gray-500">提供方</span>
                      <span class="font-medium">{{ formatProvider(row.original.upstream?.provider) }}</span>
                    </div>
                  </div>
                  <div class="rounded-md border border-default p-3">
                    <div class="mb-2 text-gray-500">本地</div>
                    <div class="flex items-center gap-2">
                      <span class="text-gray-500">模型名称</span>
                      <span class="font-medium">{{ row.original.local?.name || '—' }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-gray-500">提供方</span>
                      <span class="font-medium">{{ formatProvider(row.original.local?.provider) }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="text-gray-500">模型ID</span>
                      <span class="font-medium">{{ row.original.local?.id || '—' }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </UTable>
        </div>
      </template>
      <template #footer="{ close }">
        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" :disabled="isSyncConfirming" @click="close">取消</UButton>
          <UButton color="primary" :loading="isSyncConfirming" :disabled="!addedItems.length" @click="confirmSync">确认同步
          </UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import { computed, h, ref, resolveComponent, watch } from 'vue'
import type { TableColumn } from '@nuxt/ui'

type ModelRow = {
  id: string
  name: string
  provider: string
  status: string
  price_input?: number | null
  price_output?: number | null
  currency?: string | null
  capabilities?: string[] | null
  provider_icon?: string | null
  metadata?: Record<string, unknown> | null
  updated_at?: string | null
}

type ModelListResponse = {
  items: ModelRow[]
}
type ModelProviderListResponse = {
  items: string[]
}

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const providerFilter = ref('all')
const statusFilter = ref('all')
const isModalOpen = ref(false)
const isEditOpen = ref(false)
const isSubmitting = ref(false)
const isEditSubmitting = ref(false)
const isSyncing = ref(false)
const isSyncConfirming = ref(false)
const isSyncPreviewOpen = ref(false)
const formError = ref('')
const editError = ref('')
const syncError = ref('')
const syncPreviewError = ref('')
const actionError = ref('')
const isBulkOpen = ref(false)
const isBulkSubmitting = ref(false)
const bulkError = ref('')
const bulkSuccess = ref('')
const selectedIds = ref<string[]>([])

const providerOptions = computed(() => [
  { label: '全部提供方', value: 'all' },
  ...providerBaseItems.value
])
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
const modelStatusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'disabled' }
]
const bulkStatusOptions = [
  { label: '不修改', value: 'keep' },
  ...modelStatusOptions
]
const currencyOptions = [
  { label: 'USD', value: 'USD' },
  { label: 'CNY', value: 'CNY' }
]
const bulkCurrencyOptions = [
  { label: '不修改', value: 'keep' },
  ...currencyOptions
]
const customProviders = ref<string[]>([])
const capabilityOptions = [
  { label: '对话', value: 'chat' },
  { label: '补全', value: 'completion' },
  { label: '向量', value: 'embedding' },
  { label: '视觉', value: 'vision' },
  { label: '图像', value: 'image' },
  { label: '音频', value: 'audio' },
  { label: '工具', value: 'tool' },
  { label: '技能', value: 'skill' },
  { label: '流式', value: 'stream' },
  { label: 'JSON 模式', value: 'json_mode' },
  { label: '函数调用', value: 'function_call' }
]
const capabilityModeOptions = [
  { label: '不修改', value: 'keep' },
  { label: '覆盖更新', value: 'override' }
]

const formState = ref({
  name: '',
  provider: '',
  status: 'active',
  currency: 'USD',
  priceInput: '',
  priceOutput: '',
  capabilities: [] as string[]
})
const editFormState = ref({
  provider: '',
  currency: 'USD',
  priceInput: '',
  priceOutput: '',
  capabilities: [] as string[],
  providerIcon: ''
})
const bulkForm = ref({
  status: 'keep',
  currency: 'keep',
  priceInput: '',
  priceOutput: '',
  capabilityMode: 'keep',
  capabilities: [] as string[]
})
const editingModel = ref<ModelRow | null>(null)
const updatingIds = ref<Record<string, boolean>>({})
const sorting = ref<{ id: string; desc: boolean }[]>([])

const queryParams = computed(() => ({
  search: searchTerm.value || undefined,
  provider: providerFilter.value === 'all' ? undefined : providerFilter.value,
  status: statusFilter.value === 'all' ? undefined : statusFilter.value
}))

const { data: listData, pending: isLoading, refresh: refreshList, error: listFetchError } = await useFetch<ModelListResponse>('/api/admin/models', {
  query: queryParams,
  default: () => ({
    items: []
  })
})
const { data: providerData } = await useFetch<ModelProviderListResponse>('/api/admin/models/providers', {
  default: () => ({
    items: []
  })
})

const listItems = computed(() => listData.value?.items ?? [])
const providerItems = computed(() => providerData.value?.items ?? [])

const filteredItems = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return listItems.value.filter((item) => {
    const matchesKeyword = !keyword || item.id.toLowerCase().includes(keyword) || item.name.toLowerCase().includes(keyword)
    const matchesProvider = providerFilter.value === 'all' || item.provider === providerFilter.value
    const matchesStatus = statusFilter.value === 'all' || item.status === statusFilter.value
    return matchesKeyword && matchesProvider && matchesStatus
  })
})

const totalCount = computed(() => filteredItems.value.length)

const normalizeNumber = (value?: number | null) => {
  if (value === null || value === undefined || Number.isNaN(Number(value))) return null
  return Number(value)
}
const normalizeTimeValue = (value?: string | null) => {
  const time = value ? new Date(value).getTime() : Number.NaN
  if (Number.isNaN(time)) return null
  return time
}
const sortedItems = computed(() => {
  const sortRule = sorting.value?.[0]
  if (!sortRule) return filteredItems.value
  const { id, desc } = sortRule
  const list = [...filteredItems.value]
  const resolveValue = (item: ModelRow) => {
    if (id === 'price_input') return normalizeNumber(item.price_input ?? null)
    if (id === 'price_output') return normalizeNumber(item.price_output ?? null)
    if (id === 'updated_at') return normalizeTimeValue(item.updated_at ?? null)
    return null
  }
  list.sort((a, b) => {
    const aValue = resolveValue(a)
    const bValue = resolveValue(b)
    if (aValue === null && bValue === null) return 0
    if (aValue === null) return 1
    if (bValue === null) return -1
    if (aValue === bValue) return 0
    const result = aValue > bValue ? 1 : -1
    return desc ? -result : result
  })
  return list
})

const pagedItems = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return sortedItems.value.slice(start, start + pageSize.value)
})
const selectedCount = computed(() => selectedIds.value.length)
const selectedIdSet = computed(() => new Set(selectedIds.value))
const pageIds = computed(() => pagedItems.value.map((item) => item.id).filter(Boolean))
const isPageAllSelected = computed(() => pageIds.value.length > 0 && pageIds.value.every((id) => selectedIdSet.value.has(id)))
const isPageSomeSelected = computed(() => pageIds.value.some((id) => selectedIdSet.value.has(id)) && !isPageAllSelected.value)

watch([searchTerm, providerFilter, statusFilter, pageSize, sorting], () => {
  page.value = 1
})
watch(listItems, (items) => {
  const existingIds = new Set(items.map((item) => item.id))
  const next = selectedIds.value.filter((id) => existingIds.has(id))
  if (next.length !== selectedIds.value.length) {
    selectedIds.value = next
  }
})

const listError = computed(() => {
  if (!listFetchError.value) return ''
  const data = listFetchError.value.data as { message?: string; error?: string } | undefined
  return data?.message || data?.error || listFetchError.value.statusMessage || '加载模型列表失败'
})

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')
const UBadgeNeutral = resolveComponent('UBadge')
const UCheckbox = resolveComponent('UCheckbox')

const formatPrice = (value?: number | null, currency?: string | null) => {
  if (value === null || value === undefined || Number.isNaN(Number(value))) return '—'
  const currencyLabel = currency || '—'
  return `${Number(value).toFixed(4)} / 1M ${currencyLabel}`
}
const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

type UpstreamModel = { name: string; provider: string }
type ModelSyncResponse = { items: UpstreamModel[] }
type ModelConfirmItem = { name: string; provider: string }
type ModelConfirmRequest = { items: ModelConfirmItem[] }
type ModelConfirmResponse = { items?: ModelRow[]; created?: number; updated?: number }
type ModelPricingItem = {
  id: string
  price_input?: number | null
  price_output?: number | null
  currency?: string | null
  capabilities?: string[] | null
  status?: string | null
}
type ModelPricingRequest = { items: ModelPricingItem[] }
type ModelPricingResponse = { items?: ModelRow[]; missed_ids?: string[] }

const upstreamItems = ref<UpstreamModel[]>([])

const keyForUpstream = (u: UpstreamModel) => `${u.provider.toLowerCase()}::${u.name.toLowerCase()}`
const keyForLocal = (m: ModelRow) => `${String(m.provider || '').toLowerCase()}::${String(m.name || '').toLowerCase()}`

const addedItems = computed(() => {
  const localKeys = new Set(listItems.value.map(keyForLocal))
  return upstreamItems.value.filter((u) => !localKeys.has(keyForUpstream(u)))
})
type SyncRow = {
  key: string
  id?: string | null
  name: string
  provider: string
  upstream?: UpstreamModel | null
  local?: ModelRow | null
  status: 'added' | 'existing' | 'missing'
}

const syncRows = computed<SyncRow[]>(() => {
  const uk = new Set(upstreamItems.value.map(keyForUpstream))
  const lk = new Set(listItems.value.map(keyForLocal))
  const keys = new Set<string>([...uk, ...lk])
  const rows: SyncRow[] = []
  for (const key of keys) {
    const upstream = upstreamItems.value.find((u) => keyForUpstream(u) === key) || null
    const local = listItems.value.find((m) => keyForLocal(m) === key) || null
    const status = upstream && local ? 'existing' : upstream && !local ? 'added' : 'missing'
    const name = upstream?.name || local?.name || ''
    const provider = upstream?.provider || local?.provider || ''
    rows.push({
      key,
      id: local?.id || null,
      name,
      provider,
      upstream,
      local,
      status
    })
  }
  return rows
})

const table = ref()
const columnVisibility = ref<Record<string, boolean>>({})
const columnLabelMap: Record<string, string> = {
  select: '选择',
  id: '模型ID',
  name: '模型名称',
  provider: '提供方',
  status: '状态',
  price_input: '输入价/1M',
  price_output: '输出价/1M',
  updated_at: '更新时间',
  actions: '操作'
}
const columnPinning = ref<Record<string, string[]>>({
  right: ['actions']
})
const columnToggleItems = computed(() => {
  const columns = table.value?.tableApi?.getAllColumns?.() || []
  return columns
    .filter((column: { getCanHide?: () => boolean }) => column.getCanHide?.())
    .map((column: { id: string; getIsVisible?: () => boolean; toggleVisibility?: (value: boolean) => void }) => ({
      label: columnLabelMap[column.id] || column.id,
      type: 'checkbox' as const,
      checked: column.getIsVisible?.() ?? true,
      onUpdateChecked(checked: boolean) {
        column.toggleVisibility?.(!!checked)
      },
      onSelect(event: Event) {
        event.preventDefault()
      }
    }))
})

const previewExpanded = ref<Record<string, boolean>>({})

const formatProvider = (value?: string | null) => {
  if (!value) return '—'
  return value === 'openai' ? 'OpenAI' : value === 'deepseek' ? 'DeepSeek' : value === 'qwen' ? '通义千问' : value
}
const providerBaseItems = computed(() => {
  const seen = new Map<string, string>()
  for (const provider of providerItems.value) {
    const value = provider?.trim()
    if (!value) continue
    if (!seen.has(value)) {
      seen.set(value, formatProvider(value))
    }
  }
  return Array.from(seen.entries()).map(([value, label]) => ({ label, value }))
})
const providerInputItems = computed(() => {
  const baseValues = new Set(providerBaseItems.value.map((item) => item.value))
  const customItems = customProviders.value
    .filter((value) => value && !baseValues.has(value))
    .map((value) => ({ label: formatProvider(value), value }))
  return [...providerBaseItems.value, ...customItems]
})

const previewColumns = computed<TableColumn<SyncRow>[]>(() => [
  {
    id: 'expand',
    cell: ({ row }) =>
      h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        icon: 'i-lucide-chevron-down',
        square: true,
        size: 'xs',
        'aria-label': '展开',
        ui: {
          leadingIcon: ['transition-transform', row.getIsExpanded() ? 'duration-200 rotate-180' : '']
        },
        onClick: () => row.toggleExpanded()
      }),
    meta: { class: { th: 'w-10' } }
  },
  {
    accessorKey: 'name',
    header: '模型名称',
    cell: ({ row }) => row.original.name || '—'
  },
  {
    id: 'existing',
    header: '已存在',
    cell: ({ row }) => {
      const exists = !!row.original.local && !!row.original.upstream
      const color = exists ? 'success' : 'neutral'
      const label = exists ? '是' : '否'
      return h(UBadgeNeutral, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    id: 'local',
    header: '本地',
    cell: ({ row }) => {
      const exists = !!row.original.local
      const color = exists ? 'info' : 'neutral'
      const label = exists ? '是' : '否'
      return h(UBadgeNeutral, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    id: 'added',
    header: '新增',
    cell: ({ row }) => {
      const added = row.original.status === 'added'
      const color = added ? 'primary' : 'neutral'
      const label = added ? '是' : '否'
      return h(UBadgeNeutral, { color, variant: 'subtle' }, () => label)
    }
  }
])

const handleSync = async () => {
  syncError.value = ''
  syncPreviewError.value = ''
  isSyncing.value = true
  try {
    const res = await $fetch<ModelSyncResponse>('/api/admin/models/sync', { method: 'POST' })
    upstreamItems.value = res?.items || []
    isSyncPreviewOpen.value = true
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    syncPreviewError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '拉取上游模型失败'
  } finally {
    isSyncing.value = false
  }
}

const confirmSync = async () => {
  syncPreviewError.value = ''
  isSyncConfirming.value = true
  try {
    const payloadItems = addedItems.value.map((item) => ({
      name: item.name,
      provider: item.provider
    }))
    if (!payloadItems.length) {
      syncPreviewError.value = '暂无需要入库的模型'
      return
    }
    await $fetch<ModelConfirmResponse>('/api/admin/models/confirm', {
      method: 'POST',
      body: { items: payloadItems } satisfies ModelConfirmRequest
    })
    isSyncPreviewOpen.value = false
    upstreamItems.value = []
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    syncPreviewError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '同步上游模型失败'
  } finally {
    isSyncConfirming.value = false
  }
}

const resetForm = () => {
  formState.value = {
    name: '',
    provider: '',
    status: 'active',
    currency: 'USD',
    priceInput: '',
    priceOutput: '',
    capabilities: []
  }
  formError.value = ''
}

const openCreateModal = () => {
  resetForm()
  isModalOpen.value = true
}
const resetBulkForm = () => {
  bulkForm.value = {
    status: 'keep',
    currency: 'keep',
    priceInput: '',
    priceOutput: '',
    capabilityMode: 'keep',
    capabilities: []
  }
  bulkError.value = ''
}
const openBulkModal = () => {
  bulkSuccess.value = ''
  resetBulkForm()
  isBulkOpen.value = true
}
const clearSelection = () => {
  selectedIds.value = []
}

const addProviderItem = (value: string) => {
  const item = value.trim()
  if (!item) return ''
  if (!providerBaseItems.value.some((provider) => provider.value === item) && !customProviders.value.includes(item)) {
    customProviders.value = [...customProviders.value, item]
  }
  return item
}

const onCreateProviderForCreate = (value: string) => {
  const item = addProviderItem(value)
  if (item) {
    formState.value.provider = item
  }
}

const onCreateProviderForEdit = (value: string) => {
  const item = addProviderItem(value)
  if (item) {
    editFormState.value.provider = item
  }
}

const updateSelectedIds = (updater: (current: Set<string>) => Set<string>) => {
  const next = updater(new Set(selectedIds.value))
  selectedIds.value = Array.from(next)
}
const toggleRowSelection = (id: string, checked: boolean) => {
  updateSelectedIds((set) => {
    if (checked) {
      set.add(id)
    } else {
      set.delete(id)
    }
    return set
  })
}
const togglePageSelection = (value: boolean | 'indeterminate') => {
  updateSelectedIds((set) => {
    if (value === true) {
      pageIds.value.forEach((id) => set.add(id))
    } else {
      pageIds.value.forEach((id) => set.delete(id))
    }
    return set
  })
}

const setUpdating = (id: string, value: boolean) => {
  updatingIds.value = { ...updatingIds.value, [id]: value }
}

const openEditModal = (row: ModelRow) => {
  editingModel.value = row
  addProviderItem(row.provider || '')
  editFormState.value = {
    provider: row.provider || '',
    currency: row.currency || 'USD',
    priceInput: row.price_input !== null && row.price_input !== undefined ? String(row.price_input) : '',
    priceOutput: row.price_output !== null && row.price_output !== undefined ? String(row.price_output) : '',
    capabilities: Array.isArray(row.capabilities) ? [...row.capabilities] : [],
    providerIcon: row.provider_icon || ''
  }
  editError.value = ''
  isEditOpen.value = true
}

const submitForm = async () => {
  formError.value = ''
  if (!formState.value.name.trim()) {
    formError.value = '请输入模型名称'
    return
  }
  if (!formState.value.provider) {
    formError.value = '请选择提供方'
    return
  }
  const priceInput = formState.value.priceInput ? Number(formState.value.priceInput) : undefined
  const priceOutput = formState.value.priceOutput ? Number(formState.value.priceOutput) : undefined
  if (formState.value.priceInput && Number.isNaN(priceInput)) {
    formError.value = '请输入有效的输入价'
    return
  }
  if (formState.value.priceOutput && Number.isNaN(priceOutput)) {
    formError.value = '请输入有效的输出价'
    return
  }
  const capabilities = formState.value.capabilities
  isSubmitting.value = true
  try {
    await $fetch('/api/admin/models', {
      method: 'POST',
      body: {
        name: formState.value.name.trim(),
        provider: formState.value.provider,
        status: formState.value.status || undefined,
        currency: formState.value.currency || undefined,
        price_input: priceInput,
        price_output: priceOutput,
        capabilities: capabilities.length ? capabilities : undefined
      }
    })
    isModalOpen.value = false
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    formError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '创建模型失败'
  } finally {
    isSubmitting.value = false
  }
}

const submitEdit = async () => {
  editError.value = ''
  actionError.value = ''
  const current = editingModel.value
  if (!current?.id) {
    editError.value = '缺少模型ID'
    return
  }
  if (!editFormState.value.provider) {
    editError.value = '请选择提供方'
    return
  }
  const priceInputValue = editFormState.value.priceInput.trim()
  const priceOutputValue = editFormState.value.priceOutput.trim()
  const priceInput = priceInputValue ? Number(priceInputValue) : null
  const priceOutput = priceOutputValue ? Number(priceOutputValue) : null
  if (priceInputValue && Number.isNaN(priceInput)) {
    editError.value = '请输入有效的输入价'
    return
  }
  if (priceOutputValue && Number.isNaN(priceOutput)) {
    editError.value = '请输入有效的输出价'
    return
  }
  isEditSubmitting.value = true
  try {
    await $fetch(`/api/admin/models/${current.id}`, {
      method: 'PATCH',
      body: {
        provider: editFormState.value.provider || undefined,
        currency: editFormState.value.currency || undefined,
        price_input: priceInput,
        price_output: priceOutput,
        capabilities: editFormState.value.capabilities.length ? editFormState.value.capabilities : [],
        provider_icon: editFormState.value.providerIcon ? editFormState.value.providerIcon : null
      }
    })
    isEditOpen.value = false
    editingModel.value = null
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    editError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '更新模型失败'
  } finally {
    isEditSubmitting.value = false
  }
}

const submitBulk = async () => {
  bulkError.value = ''
  bulkSuccess.value = ''
  if (!selectedIds.value.length) {
    bulkError.value = '请选择需要批量管理的模型'
    return
  }
  const priceInputValue = bulkForm.value.priceInput.trim()
  const priceOutputValue = bulkForm.value.priceOutput.trim()
  const priceInput = priceInputValue ? Number(priceInputValue) : undefined
  const priceOutput = priceOutputValue ? Number(priceOutputValue) : undefined
  if (priceInputValue && Number.isNaN(priceInput)) {
    bulkError.value = '请输入有效的输入价'
    return
  }
  if (priceOutputValue && Number.isNaN(priceOutput)) {
    bulkError.value = '请输入有效的输出价'
    return
  }
  const shouldUpdateCapabilities = bulkForm.value.capabilityMode === 'override'
  const shouldUpdateStatus = bulkForm.value.status !== 'keep'
  const shouldUpdateCurrency = bulkForm.value.currency !== 'keep'
  const hasUpdates =
    shouldUpdateStatus ||
    shouldUpdateCurrency ||
    !!priceInputValue ||
    !!priceOutputValue ||
    shouldUpdateCapabilities
  if (!hasUpdates) {
    bulkError.value = '请选择需要批量更新的字段'
    return
  }
  isBulkSubmitting.value = true
  try {
    const payloadItems: ModelPricingItem[] = selectedIds.value.map((id) => ({
      id,
      ...(shouldUpdateStatus ? { status: bulkForm.value.status } : {}),
      ...(shouldUpdateCurrency ? { currency: bulkForm.value.currency } : {}),
      ...(priceInputValue ? { price_input: priceInput } : {}),
      ...(priceOutputValue ? { price_output: priceOutput } : {}),
      ...(shouldUpdateCapabilities ? { capabilities: bulkForm.value.capabilities } : {})
    }))
    const res = await $fetch<ModelPricingResponse>('/api/admin/models/pricing', {
      method: 'POST',
      body: { items: payloadItems } satisfies ModelPricingRequest
    })
    const missedCount = res?.missed_ids?.length ?? 0
    const updatedCount = res?.items?.length ?? Math.max(selectedIds.value.length - missedCount, 0)
    bulkSuccess.value = missedCount
      ? `已更新 ${updatedCount} 个，未找到 ${missedCount} 个模型`
      : `已更新 ${updatedCount} 个模型`
    isBulkOpen.value = false
    clearSelection()
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    bulkError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '批量更新模型失败'
  } finally {
    isBulkSubmitting.value = false
  }
}

const toggleModelStatus = async (row: ModelRow) => {
  actionError.value = ''
  const nextStatus = row.status === 'active' ? 'disabled' : 'active'
  if (!row.id) {
    actionError.value = '缺少模型ID'
    return
  }
  setUpdating(row.id, true)
  try {
    await $fetch(`/api/admin/models/${row.id}`, {
      method: 'PATCH',
      body: { status: nextStatus }
    })
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    actionError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '更新模型状态失败'
  } finally {
    setUpdating(row.id, false)
  }
}

const columns = computed<TableColumn<ModelRow>[]>(() => [
  {
    id: 'select',
    header: () =>
      h(UCheckbox, {
        modelValue: isPageAllSelected.value ? true : isPageSomeSelected.value ? 'indeterminate' : false,
        disabled: pageIds.value.length === 0,
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => togglePageSelection(value)
      }),
    cell: ({ row }) =>
      h(UCheckbox, {
        modelValue: selectedIdSet.value.has(row.original.id),
        disabled: !row.original.id,
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => toggleRowSelection(row.original.id, value === true)
      }),
    meta: { class: { th: 'w-10' } }
  },
  {
    accessorKey: 'id',
    header: '模型ID',
    meta: { class: { th: 'w-40' } }
  },
  {
    accessorKey: 'name',
    header: '模型名称'
  },
  {
    accessorKey: 'provider',
    header: '提供方',
    cell: ({ row }) => {
      const providerValue = String(row.getValue('provider') || '')
      const display = providerValue === 'openai' ? 'OpenAI' : providerValue === 'deepseek' ? 'DeepSeek' : providerValue === 'qwen' ? '通义千问' : providerValue
      return h(UBadge, { color: 'neutral', variant: 'subtle' }, () => display)
    }
  },
  {
    accessorKey: 'status',
    header: '状态',
    cell: ({ row }) => {
      const statusValue = String(row.getValue('status') || '')
      const color = statusValue === 'active' ? 'success' : statusValue === 'disabled' ? 'error' : 'neutral'
      const label = statusValue === 'active' ? '启用' : statusValue === 'disabled' ? '禁用' : '未知'
      return h(UBadge, { color, variant: 'subtle' }, () => label)
    }
  },
  {
    accessorKey: 'price_input',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: '输入价/1M',
        icon:
          isSorted === 'asc'
            ? 'i-lucide-arrow-up-narrow-wide'
            : isSorted === 'desc'
              ? 'i-lucide-arrow-down-wide-narrow'
              : 'i-lucide-arrow-up-down',
        class: '-mx-2.5',
        onClick: () => column.toggleSorting(isSorted === 'asc')
      })
    },
    cell: ({ row }) => formatPrice(row.getValue('price_input') as number | null | undefined, row.original.currency)
  },
  {
    accessorKey: 'price_output',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: '输出价/1M',
        icon:
          isSorted === 'asc'
            ? 'i-lucide-arrow-up-narrow-wide'
            : isSorted === 'desc'
              ? 'i-lucide-arrow-down-wide-narrow'
              : 'i-lucide-arrow-up-down',
        class: '-mx-2.5',
        onClick: () => column.toggleSorting(isSorted === 'asc')
      })
    },
    cell: ({ row }) => formatPrice(row.getValue('price_output') as number | null | undefined, row.original.currency)
  },
  {
    accessorKey: 'updated_at',
    header: ({ column }) => {
      const isSorted = column.getIsSorted()
      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: '更新时间',
        icon:
          isSorted === 'asc'
            ? 'i-lucide-arrow-up-narrow-wide'
            : isSorted === 'desc'
              ? 'i-lucide-arrow-down-wide-narrow'
              : 'i-lucide-arrow-up-down',
        class: '-mx-2.5',
        onClick: () => column.toggleSorting(isSorted === 'asc')
      })
    },
    cell: ({ row }) => formatTime(String(row.getValue('updated_at')))
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) => {
      const model = row.original
      const isUpdating = !!updatingIds.value[model.id]
      const statusLabel = model.status === 'active' ? '停用' : '启用'
      const statusColor = model.status === 'active' ? 'error' : 'success'
      return h('div', { class: 'flex items-center gap-2' }, [
        h(UButton, {
          label: statusLabel,
          color: statusColor,
          variant: 'ghost',
          size: 'xs',
          loading: isUpdating,
          disabled: isUpdating,
          onClick: () => toggleModelStatus(model)
        }),
        h(UButton, {
          label: '编辑',
          color: 'primary',
          variant: 'ghost',
          size: 'xs',
          onClick: () => openEditModal(model)
        })
      ])
    }
  }
])
</script>
