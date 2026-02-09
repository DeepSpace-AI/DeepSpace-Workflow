<template>
  <div class="space-y-4">
    <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">风控策略</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="userIdTerm" placeholder="用户ID" icon="i-heroicons-magnifying-glass" class="w-36" />
          <UInput v-model="projectIdTerm" placeholder="项目ID" class="w-36" />
          <USelect v-model="scopeFilter" :items="scopeOptions" class="w-28" />
          <USelect v-model="statusFilter" :items="statusOptions" class="w-28" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
          <UButton icon="i-heroicons-plus" label="添加策略" color="primary" @click="openCreateModal" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="policyItems" :columns="columns" :loading="isLoading" />
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-500">共 {{ totalCount }} 条</div>
        <UPagination v-model:page="page" :total="totalCount" :items-per-page="pageSize" :sibling-count="1" show-edges />
      </div>
    </div>
    </UCard>

    <UModal v-model:open="isModalOpen" :title="modalTitle">
    <template #body>
      <div class="grid gap-4 md:grid-cols-2">
        <UFormField label="策略名称" required class="md:col-span-2">
          <UInput v-model="formState.name" class="w-full" placeholder="请输入策略名称" :disabled="isPolicyLocked" />
        </UFormField>
        <UFormField label="范围" required>
          <USelect v-model="formState.scope" class="w-full" :items="formScopeOptions" :disabled="isPolicyLocked" />
        </UFormField>
        <UFormField label="状态" required>
          <USelect v-model="formState.status" class="w-full" :items="formStatusOptions" :disabled="isPolicyLocked" />
        </UFormField>
        <UFormField label="优先级" required>
          <UInput v-model="formState.priority" class="w-full" type="number" min="0" placeholder="请输入优先级" :disabled="isPolicyLocked" />
        </UFormField>
        <UFormField label="用户ID">
          <UInput v-model="formState.userId" class="w-full" type="number" min="1" placeholder="可选" :disabled="isPolicyLocked" />
        </UFormField>
        <UFormField label="项目ID">
          <UInput v-model="formState.projectId" class="w-full" type="number" min="1" placeholder="可选" :disabled="isPolicyLocked" />
        </UFormField>
        <template v-if="isCreateMode">
          <UFormField label="策略类型" required class="md:col-span-2">
            <USelect v-model="formState.policyType" class="w-full" :items="policyTypeOptions" :disabled="createdPolicyId !== null" />
          </UFormField>
          <template v-if="formState.policyType === 'ip_rule'">
            <UFormField label="IP 类型" required>
              <USelect v-model="formState.ipType" class="w-full" :items="ipTypeOptions" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="IP 地址">
              <UInput v-model="formState.ipValue" class="w-full" placeholder="例如 192.168.1.1" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="CIDR">
              <UInput v-model="formState.cidrValue" class="w-full" placeholder="例如 192.168.1.0/24" :disabled="createdPolicyId !== null" />
            </UFormField>
          </template>
          <template v-else-if="formState.policyType === 'rate_limit'">
            <UFormField label="时间窗口（秒）" required>
              <UInput v-model="formState.windowSeconds" class="w-full" type="number" min="1" placeholder="例如 60" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="最大请求数">
              <UInput v-model="formState.maxRequests" class="w-full" type="number" min="0" placeholder="例如 100" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="最大 Tokens">
              <UInput v-model="formState.maxTokens" class="w-full" type="number" min="0" placeholder="例如 10000" :disabled="createdPolicyId !== null" />
            </UFormField>
          </template>
          <template v-else-if="formState.policyType === 'budget_cap'">
            <UFormField label="统计周期" required>
              <USelect v-model="formState.budgetCycle" class="w-full" :items="budgetCycleOptions" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="预算上限" required>
              <UInput v-model="formState.budgetMaxCost" class="w-full" type="number" min="0.01" placeholder="例如 1000" :disabled="createdPolicyId !== null" />
            </UFormField>
            <UFormField label="币种">
              <UInput v-model="formState.budgetCurrency" class="w-full" placeholder="默认 CNY" :disabled="createdPolicyId !== null" />
            </UFormField>
          </template>
        </template>
        <div v-if="createdPolicyId !== null" class="text-sm text-amber-600 md:col-span-2">策略已创建，请补充子规则后再次保存</div>
        <div v-if="formError" class="text-sm text-red-600 md:col-span-2">{{ formError }}</div>
      </div>
    </template>
    <template #footer="{ close }">
      <div class="flex justify-end gap-2">
        <UButton color="neutral" variant="outline" :disabled="isSubmitting" @click="close">取消</UButton>
        <UButton color="primary" :loading="isSubmitting" @click="submitForm">保存</UButton>
      </div>
    </template>
    </UModal>

    <UModal v-model:open="isDeleteOpen" title="确认删除">
    <template #body>
      <div class="text-sm text-gray-600">
        将删除策略 <span class="font-medium text-gray-900">{{ deleteTarget?.name || '' }}</span>，此操作不可撤销。
      </div>
      <div v-if="deleteError" class="mt-2 text-sm text-red-600">{{ deleteError }}</div>
    </template>
    <template #footer="{ close }">
      <div class="flex justify-end gap-2">
        <UButton color="neutral" variant="outline" :disabled="isDeleting" @click="close">取消</UButton>
        <UButton color="error" :loading="isDeleting" @click="confirmDelete">删除</UButton>
      </div>
    </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'

type PolicyRow = {
  id: number
  name: string
  scope: 'global' | 'user' | 'project'
  user_id?: number | null
  project_id?: number | null
  status: 'active' | 'disabled'
  priority: number
  created_at: string
  updated_at: string
}

type PolicyListResponse = {
  items: PolicyRow[]
  total: number
  page: number
  page_size: number
}

type PolicyFormState = {
  name: string
  scope: PolicyRow['scope']
  status: PolicyRow['status']
  priority: string
  userId: string
  projectId: string
  policyType: 'ip_rule' | 'rate_limit' | 'budget_cap'
  ipType: 'allow' | 'deny'
  ipValue: string
  cidrValue: string
  windowSeconds: string
  maxRequests: string
  maxTokens: string
  budgetCycle: 'daily' | 'weekly' | 'monthly'
  budgetMaxCost: string
  budgetCurrency: string
}

const page = ref(1)
const pageSize = ref(10)
const userIdTerm = ref('')
const projectIdTerm = ref('')
const scopeFilter = ref<'all' | PolicyRow['scope']>('all')
const statusFilter = ref<'all' | PolicyRow['status']>('all')
const isModalOpen = ref(false)
const isDeleteOpen = ref(false)
const isSubmitting = ref(false)
const isDeleting = ref(false)
const formError = ref('')
const deleteError = ref('')
const editingPolicy = ref<PolicyRow | null>(null)
const deleteTarget = ref<PolicyRow | null>(null)
const createdPolicyId = ref<number | null>(null)
const formState = ref<PolicyFormState>({
  name: '',
  scope: 'global',
  status: 'active',
  priority: '1',
  userId: '',
  projectId: '',
  policyType: 'ip_rule',
  ipType: 'allow',
  ipValue: '',
  cidrValue: '',
  windowSeconds: '60',
  maxRequests: '100',
  maxTokens: '',
  budgetCycle: 'monthly',
  budgetMaxCost: '1000',
  budgetCurrency: 'CNY'
})

const scopeOptions = [
  { label: '全部范围', value: 'all' },
  { label: '全局', value: 'global' },
  { label: '用户', value: 'user' },
  { label: '项目', value: 'project' }
]
const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'disabled' }
]
const formScopeOptions = [
  { label: '全局', value: 'global' },
  { label: '用户', value: 'user' },
  { label: '项目', value: 'project' }
]
const formStatusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'disabled' }
]
const policyTypeOptions = [
  { label: 'IP 规则', value: 'ip_rule' },
  { label: '速率限制', value: 'rate_limit' },
  { label: '预算上限', value: 'budget_cap' }
]
const ipTypeOptions = [
  { label: '允许', value: 'allow' },
  { label: '拒绝', value: 'deny' }
]
const budgetCycleOptions = [
  { label: '每日', value: 'daily' },
  { label: '每周', value: 'weekly' },
  { label: '每月', value: 'monthly' }
]
const pageSizeOptions = [
  { label: '10 / 页', value: 10 },
  { label: '20 / 页', value: 20 },
  { label: '50 / 页', value: 50 }
]

watch([userIdTerm, projectIdTerm, scopeFilter, statusFilter, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')

const isCreateMode = computed(() => !editingPolicy.value)
const isPolicyLocked = computed(() => isCreateMode.value && createdPolicyId.value !== null)

watch(isModalOpen, (open) => {
  if (!open) {
    createdPolicyId.value = null
  }
})

const userIdQuery = computed(() => {
  const keyword = userIdTerm.value.trim()
  if (!keyword) return undefined
  const parsed = Number(keyword)
  if (!Number.isInteger(parsed) || parsed <= 0) return undefined
  return parsed
})

const projectIdQuery = computed(() => {
  const keyword = projectIdTerm.value.trim()
  if (!keyword) return undefined
  const parsed = Number(keyword)
  if (!Number.isInteger(parsed) || parsed <= 0) return undefined
  return parsed
})

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  scope: scopeFilter.value === 'all' ? undefined : scopeFilter.value,
  status: statusFilter.value === 'all' ? undefined : statusFilter.value,
  user_id: userIdQuery.value,
  project_id: projectIdQuery.value
}))

const { data: listData, pending: isLoading, refresh: refreshList } = await useFetch<PolicyListResponse | { data?: PolicyListResponse }>(
  '/api/admin/risk/policies',
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

const rawPolicyItems = computed(() => {
  const dataValue = listData.value as PolicyListResponse | { data?: PolicyListResponse } | undefined
  if (!dataValue) return []
  if (Array.isArray((dataValue as PolicyListResponse).items)) return (dataValue as PolicyListResponse).items
  if (Array.isArray(dataValue.data?.items)) return dataValue.data.items
  return []
})

const normalizePolicyRow = (item: Record<string, unknown>): PolicyRow => {
  return {
    id: Number(item.id ?? item.ID ?? 0),
    name: String(item.name ?? item.Name ?? ''),
    scope: (item.scope ?? item.Scope ?? 'global') as PolicyRow['scope'],
    user_id: (item.user_id ?? item.UserID ?? null) as number | null,
    project_id: (item.project_id ?? item.ProjectID ?? null) as number | null,
    status: (item.status ?? item.Status ?? 'active') as PolicyRow['status'],
    priority: Number(item.priority ?? item.Priority ?? 0),
    created_at: String(item.created_at ?? item.CreatedAt ?? ''),
    updated_at: String(item.updated_at ?? item.UpdatedAt ?? '')
  }
}

const policyItems = computed(() => rawPolicyItems.value.map((item) => normalizePolicyRow(item as Record<string, unknown>)))
const totalCount = computed(() => {
  const dataValue = listData.value as PolicyListResponse | { data?: PolicyListResponse } | undefined
  return dataValue?.total ?? dataValue?.data?.total ?? 0
})

const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}
const formatScope = (value: string) => {
  if (value === 'global') return '全局'
  if (value === 'user') return '用户'
  if (value === 'project') return '项目'
  return value || '—'
}

const modalTitle = computed(() => (editingPolicy.value ? '编辑策略' : '添加策略'))

const resetForm = () => {
  formState.value = {
    name: '',
    scope: 'global',
    status: 'active',
    priority: '1',
    userId: '',
    projectId: '',
    policyType: 'ip_rule',
    ipType: 'allow',
    ipValue: '',
    cidrValue: '',
    windowSeconds: '60',
    maxRequests: '100',
    maxTokens: '',
    budgetCycle: 'monthly',
    budgetMaxCost: '1000',
    budgetCurrency: 'CNY'
  }
  formError.value = ''
}

const openCreateModal = () => {
  editingPolicy.value = null
  resetForm()
  isModalOpen.value = true
}

const openEditModal = (row: PolicyRow) => {
  editingPolicy.value = row
  formState.value = {
    name: row.name || '',
    scope: row.scope || 'global',
    status: row.status || 'active',
    priority: Number.isFinite(row.priority) ? String(row.priority) : '0',
    userId: row.user_id ? String(row.user_id) : '',
    projectId: row.project_id ? String(row.project_id) : '',
    policyType: 'ip_rule',
    ipType: 'allow',
    ipValue: '',
    cidrValue: '',
    windowSeconds: '60',
    maxRequests: '100',
    maxTokens: '',
    budgetCycle: 'monthly',
    budgetMaxCost: '1000',
    budgetCurrency: 'CNY'
  }
  formError.value = ''
  isModalOpen.value = true
}

const openDeleteModal = (row: PolicyRow) => {
  deleteTarget.value = row
  deleteError.value = ''
  isDeleteOpen.value = true
}

const parseOptionalPositiveInt = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  const parsed = Number(trimmed)
  if (!Number.isInteger(parsed) || parsed <= 0) return undefined
  return parsed
}

const parseOptionalPositiveNumber = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return undefined
  const parsed = Number(trimmed)
  if (!Number.isFinite(parsed) || parsed <= 0) return undefined
  return parsed
}

const parseRequiredPositiveInt = (value: string) => {
  const parsed = Number(value)
  if (!Number.isInteger(parsed) || parsed <= 0) return undefined
  return parsed
}

const parseNonNegativeInt = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return 0
  const parsed = Number(trimmed)
  if (!Number.isInteger(parsed) || parsed < 0) return undefined
  return parsed
}

const createSubRule = async (policyId: number) => {
  if (formState.value.policyType === 'ip_rule') {
    const ipValue = formState.value.ipValue.trim()
    const cidrValue = formState.value.cidrValue.trim()
    await $fetch('/api/admin/risk/ip-rules', {
      method: 'POST',
      body: {
        policy_id: policyId,
        type: formState.value.ipType,
        status: formState.value.status,
        ...(ipValue ? { ip: ipValue } : {}),
        ...(cidrValue ? { cidr: cidrValue } : {})
      }
    })
    return
  }
  if (formState.value.policyType === 'rate_limit') {
    const windowSeconds = parseRequiredPositiveInt(formState.value.windowSeconds)
    const maxRequests = parseNonNegativeInt(formState.value.maxRequests)
    const maxTokens = parseNonNegativeInt(formState.value.maxTokens)
    await $fetch('/api/admin/risk/rate-limits', {
      method: 'POST',
      body: {
        policy_id: policyId,
        window_seconds: windowSeconds,
        max_requests: maxRequests ?? 0,
        max_tokens: maxTokens ?? 0,
        status: formState.value.status
      }
    })
    return
  }
  const maxCost = parseOptionalPositiveNumber(formState.value.budgetMaxCost)
  const currencyValue = formState.value.budgetCurrency.trim() || 'CNY'
  await $fetch('/api/admin/risk/budget-caps', {
    method: 'POST',
    body: {
      policy_id: policyId,
      cycle: formState.value.budgetCycle,
      max_cost: maxCost ?? 0,
      currency: currencyValue,
      status: formState.value.status
    }
  })
}

const submitForm = async () => {
  formError.value = ''
  const nameValue = formState.value.name.trim()
  if (!nameValue) {
    formError.value = '请输入策略名称'
    return
  }
  const priorityValue = Number(formState.value.priority)
  if (!Number.isInteger(priorityValue) || priorityValue < 0) {
    formError.value = '请输入有效的优先级'
    return
  }
  const userIdValue = parseOptionalPositiveInt(formState.value.userId)
  if (formState.value.userId.trim() && userIdValue === undefined) {
    formError.value = '请输入有效的用户ID'
    return
  }
  const projectIdValue = parseOptionalPositiveInt(formState.value.projectId)
  if (formState.value.projectId.trim() && projectIdValue === undefined) {
    formError.value = '请输入有效的项目ID'
    return
  }
  if (formState.value.scope === 'user' && !userIdValue) {
    formError.value = '请选择用户范围时需填写用户ID'
    return
  }
  if (formState.value.scope === 'project' && !projectIdValue) {
    formError.value = '请选择项目范围时需填写项目ID'
    return
  }
  if (isCreateMode.value) {
    if (formState.value.policyType === 'ip_rule') {
      if (!formState.value.ipType) {
        formError.value = '请选择 IP 类型'
        return
      }
      if (!formState.value.ipValue.trim() && !formState.value.cidrValue.trim()) {
        formError.value = '请输入 IP 地址或 CIDR'
        return
      }
    }
    if (formState.value.policyType === 'rate_limit') {
      const windowSeconds = parseRequiredPositiveInt(formState.value.windowSeconds)
      if (!windowSeconds) {
        formError.value = '请输入有效的时间窗口'
        return
      }
      const maxRequests = parseNonNegativeInt(formState.value.maxRequests)
      const maxTokens = parseNonNegativeInt(formState.value.maxTokens)
      if (maxRequests === undefined || maxTokens === undefined || (maxRequests <= 0 && maxTokens <= 0)) {
        formError.value = '最大请求数与最大 Tokens 至少填写一个'
        return
      }
    }
    if (formState.value.policyType === 'budget_cap') {
      const maxCost = parseOptionalPositiveNumber(formState.value.budgetMaxCost)
      if (!maxCost) {
        formError.value = '请输入有效的预算上限'
        return
      }
      if (!formState.value.budgetCycle) {
        formError.value = '请选择统计周期'
        return
      }
    }
  }
  isSubmitting.value = true
  let createdPolicy = false
  let policyId = createdPolicyId.value
  try {
    const payload = {
      name: nameValue,
      scope: formState.value.scope,
      status: formState.value.status,
      priority: priorityValue,
      ...(userIdValue ? { user_id: userIdValue } : {}),
      ...(projectIdValue ? { project_id: projectIdValue } : {})
    }
    if (editingPolicy.value) {
      await $fetch(`/api/admin/risk/policies/${editingPolicy.value.id}`, {
        method: 'PATCH',
        body: payload
      })
    } else {
      if (!policyId) {
        const created = await $fetch<PolicyRow>('/api/admin/risk/policies', {
          method: 'POST',
          body: payload
        })
        policyId = created?.id
        if (!policyId) {
          throw new Error('创建策略失败')
        }
        createdPolicyId.value = policyId
        createdPolicy = true
      }
      await createSubRule(policyId)
    }
    isModalOpen.value = false
    createdPolicyId.value = null
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    formError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      (editingPolicy.value ? '更新策略失败' : '创建策略失败')
    if (createdPolicy || createdPolicyId.value) {
      formError.value = '策略已创建，子规则创建失败，请检查后再次保存'
      await refreshList()
    }
  } finally {
    isSubmitting.value = false
  }
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  isDeleting.value = true
  deleteError.value = ''
  try {
    await $fetch(`/api/admin/risk/policies/${deleteTarget.value.id}`, {
      method: 'DELETE'
    })
    isDeleteOpen.value = false
    deleteTarget.value = null
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    deleteError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '删除策略失败'
  } finally {
    isDeleting.value = false
  }
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
    accessorKey: 'scope',
    header: '范围',
    cell: ({ row }) => {
      const scopeValue = String(row.getValue('scope') || '')
      return h(UBadge, { color: 'neutral', variant: 'subtle' }, () => formatScope(scopeValue))
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
    accessorKey: 'priority',
    header: '优先级',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => Number(row.getValue('priority')).toLocaleString('zh-CN')
  },
  {
    accessorKey: 'user_id',
    header: '用户ID',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => row.getValue('user_id') ?? '—'
  },
  {
    accessorKey: 'project_id',
    header: '项目ID',
    meta: { class: { td: 'text-right' } },
    cell: ({ row }) => row.getValue('project_id') ?? '—'
  },
  {
    accessorKey: 'updated_at',
    header: '更新时间',
    cell: ({ row }) => formatTime(String(row.getValue('updated_at')))
  },
  {
    id: 'actions',
    header: '操作',
    cell: ({ row }) =>
      h(
        'div',
        { class: 'flex items-center gap-2' },
        [
          h(UButton, {
            label: '编辑',
            color: 'primary',
            variant: 'ghost',
            size: 'xs',
            onClick: () => openEditModal(row.original)
          }),
          h(UButton, {
            label: '删除',
            color: 'error',
            variant: 'ghost',
            size: 'xs',
            onClick: () => openDeleteModal(row.original)
          })
        ]
      )
  }
])
</script>
