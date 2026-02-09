<template>
  <div class="space-y-4">
    <UCard>
      <template #header>
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-3">
            <h3 class="text-lg font-semibold">套餐管理</h3>
            <UBadge color="neutral" variant="soft">共 {{ planCount }} 个</UBadge>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <UInput v-model="searchTerm" icon="i-heroicons-magnifying-glass" placeholder="搜索套餐名称"
              class="w-full sm:w-52" />
            <USelect v-model="statusFilter" :items="statusFilterOptions" class="w-full sm:w-40" />
            <UButton icon="i-heroicons-arrow-path" variant="soft" color="neutral" :loading="isLoading"
              @click="() => refreshList()" />
            <UButton icon="i-heroicons-plus" label="新建套餐" color="primary" @click="openCreateModal" />
          </div>
        </div>
      </template>

      <div v-if="isLoading && planCount === 0" class="py-12">
        <UEmpty icon="i-heroicons-archive-box" title="正在加载套餐" description="正在从数据库获取套餐列表。" />
      </div>
      <div v-else-if="filteredPlans.length === 0" class="py-12">
        <UEmpty icon="i-heroicons-archive-box" title="暂无套餐" description="可以先创建一个套餐，配置配额与状态。">
          <template #actions>
            <div class="flex items-center gap-2">
              <UButton color="primary" label="新建套餐" icon="i-heroicons-plus" @click="openCreateModal" />
              <UButton color="neutral" variant="soft" label="查看模型管理" icon="i-heroicons-cpu-chip" to="/models" />
            </div>
          </template>
        </UEmpty>
      </div>
      <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        <UCard
          v-for="plan in filteredPlans"
          :key="plan.ID"
          variant="subtle"
          class="group transition duration-200 hover:-translate-y-0.5 hover:shadow-lg dark:hover:shadow-black/40"
        >
          <template #header>
            <div class="flex items-start justify-between gap-3">
              <div class="space-y-1">
                <div class="text-base font-semibold text-gray-900 dark:text-gray-100">{{ plan.Name }}</div>
                <div class="flex flex-wrap items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                  <span class="rounded-full bg-gray-100 px-2 py-0.5 text-gray-600 dark:bg-gray-800 dark:text-gray-300">
                    ID {{ plan.ID }}
                  </span>
                  <span class="text-gray-300 dark:text-gray-600">|</span>
                  <span>更新于 {{ formatDate(plan.UpdatedAt || plan.CreatedAt) }}</span>
                </div>
              </div>
              <UBadge :color="statusColor(plan.Status)" variant="soft" class="shrink-0">
                {{ statusLabel(plan.Status) }}
              </UBadge>
            </div>
          </template>
          <div class="grid gap-3 text-sm text-gray-600 dark:text-gray-300 sm:grid-cols-2">
            <div class="rounded-lg border border-gray-100 bg-white/70 p-3 dark:border-gray-800 dark:bg-gray-900/60">
              <div class="text-xs text-gray-500 dark:text-gray-400">配额</div>
              <div class="mt-1 font-medium text-gray-900 dark:text-gray-100">{{ quotaLabel(plan) }}</div>
            </div>
            <div class="rounded-lg border border-gray-100 bg-white/70 p-3 dark:border-gray-800 dark:bg-gray-900/60">
              <div class="text-xs text-gray-500 dark:text-gray-400">重置周期</div>
              <div class="mt-1 font-medium text-gray-900 dark:text-gray-100">{{ cycleLabel(plan) }}</div>
            </div>
            <div class="rounded-lg border border-gray-100 bg-white/70 p-3 dark:border-gray-800 dark:bg-gray-900/60 sm:col-span-2">
              <div class="text-xs text-gray-500 dark:text-gray-400">价格</div>
              <div class="mt-1 text-base font-semibold text-gray-900 dark:text-gray-100">{{ priceLabel(plan) }}</div>
            </div>
          </div>
          <template #footer>
            <div class="flex items-center justify-end">
              <UButton variant="soft" color="primary" label="编辑套餐" @click="openEditModal(plan)" />
            </div>
          </template>
        </UCard>
      </div>
    </UCard>

    <UModal v-model:open="isModalOpen" :title="modalTitle">
      <template #body>
        <UForm :key="formKey" :state="formState" :validate="validateForm" :validate-on="[]" class="space-y-4"
          @submit="submitForm">
          <UFormField label="套餐名称" name="name" required>
            <UInput v-model="formState.name" class="w-full" placeholder="如：专业版" />
          </UFormField>
          <UFormField label="状态" name="status" required>
            <USelect v-model="formState.status" class="w-full" :items="statusOptions" />
          </UFormField>
          <UFormField label="配额类型" required>
            <URadioGroup v-model="quotaType" :items="quotaOptions" orientation="horizontal" variant="list" />
          </UFormField>
          <UFormField v-if="quotaType === 'tokens'" label="包含 Token" name="included_tokens" required>
            <UInput v-model.number="formState.included_tokens" type="number" min="0" />
          </UFormField>
          <UFormField v-else label="包含请求次数" name="included_requests" required>
            <UInput v-model.number="formState.included_requests" type="number" min="0" />
          </UFormField>
          <div class="grid gap-6 sm:grid-cols-2">
            <UFormField label="重置周期（天）" name="reset_interval_days" required>
              <UInput v-model.number="formState.reset_interval_days" class="w-full" type="number" min="1" max="365" />
            </UFormField>
            <UFormField label="套餐价格" name="price" required>
              <UInput v-model.number="formState.price" class="w-full" type="number" min="0" step="0.01" />
            </UFormField>
          </div>
          <UFormField label="价格单位" name="currency" required>
            <USelect v-model="formState.currency" class="w-full" :items="currencyOptions" />
          </UFormField>
          <div v-if="generalError" class="text-sm text-red-500">{{ generalError }}</div>
          <div class="flex items-center justify-end gap-2 pt-2">
            <UButton variant="soft" color="neutral" label="取消" type="button" @click="closeModal" />
            <UButton color="primary" :loading="isSubmitting" :label="submitLabel" type="submit" />
          </div>
        </UForm>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { FormError, FormSubmitEvent } from '@nuxt/ui'

type PlanItem = {
  ID: number
  Name: string
  Status: string
  IncludedTokens: number
  IncludedRequests: number
  ResetIntervalDays: number
  Price: number
  Currency: string
  CreatedAt?: string
  UpdatedAt?: string
}

type PlanListResponse = {
  items: PlanItem[]
}

type PlanFormState = {
  name: string
  status: string
  included_tokens: number
  included_requests: number
  reset_interval_days: number
  price: number
  currency: string
}

const statusOptions = [
  { label: '生效', value: 'active' },
  { label: '停用', value: 'disabled' },
  { label: '过期', value: 'expired' },
  { label: '已取消', value: 'canceled' }
]

const statusFilterOptions = [
  { label: '全部状态', value: 'all' },
  ...statusOptions
]

const searchTerm = ref('')
const statusFilter = ref('all')
const isModalOpen = ref(false)
const isSubmitting = ref(false)
const editingPlan = ref<PlanItem | null>(null)
const generalError = ref('')
const formKey = ref(0)

const formState = reactive<PlanFormState>({
  name: '',
  status: 'active',
  included_tokens: 0,
  included_requests: 0,
  reset_interval_days: 30,
  price: 0,
  currency: 'CNY'
})

const quotaOptions = [
  { label: '按 Token 计', value: 'tokens', description: '填写包含 Token 数量' },
  { label: '按请求次数计', value: 'requests', description: '填写包含请求次数' }
]

const quotaType = ref<'tokens' | 'requests'>('tokens')

const currencyOptions = [
  { label: '人民币 (CNY)', value: 'CNY' },
  { label: '美元 (USD)', value: 'USD' }
]

const { data: listData, pending: isLoading, refresh: refreshList } = await useFetch<PlanListResponse>('/api/admin/plans', {
  default: () => ({ items: [] })
})

const planItems = computed(() => listData.value?.items ?? [])

const filteredPlans = computed(() => {
  const keyword = searchTerm.value.trim().toLowerCase()
  return planItems.value.filter((plan) => {
    const matchName = keyword ? plan.Name.toLowerCase().includes(keyword) : true
    const matchStatus = statusFilter.value === 'all' ? true : plan.Status === statusFilter.value
    return matchName && matchStatus
  })
})

const planCount = computed(() => filteredPlans.value.length)

const modalTitle = computed(() => (editingPlan.value ? '编辑套餐' : '新建套餐'))
const submitLabel = computed(() => (editingPlan.value ? '保存修改' : '创建套餐'))

const statusLabel = (status: string) => {
  const item = statusOptions.find((option) => option.value === status)
  return item?.label || status || '未知'
}

const statusColor = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'disabled':
      return 'neutral'
    case 'expired':
      return 'warning'
    case 'canceled':
      return 'error'
    default:
      return 'neutral'
  }
}

const quotaLabel = (plan: PlanItem) => {
  if (plan.IncludedTokens > 0) {
    return `包含 Token ${plan.IncludedTokens.toLocaleString('zh-CN')}`
  }
  if (plan.IncludedRequests > 0) {
    return `包含请求次数 ${plan.IncludedRequests.toLocaleString('zh-CN')}`
  }
  return '未配置配额'
}

const cycleLabel = (plan: PlanItem) => {
  const days = plan.ResetIntervalDays || 30
  return `重置周期 ${days} 天`
}

const priceLabel = (plan: PlanItem) => {
  const price = Number.isFinite(plan.Price) ? plan.Price : 0
  const currency = plan.Currency || 'CNY'
  return `套餐价格 ${price.toLocaleString('zh-CN')} ${currency}`
}

const formatDate = (value?: string) => {
  if (!value) {
    return '-'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN')
}

const resetForm = () => {
  formState.name = ''
  formState.status = 'active'
  formState.included_tokens = 0
  formState.included_requests = 0
  formState.reset_interval_days = 30
  formState.price = 0
  formState.currency = 'CNY'
  quotaType.value = 'tokens'
  generalError.value = ''
}

const openCreateModal = () => {
  editingPlan.value = null
  resetForm()
  formKey.value += 1
  isModalOpen.value = true
}

const openEditModal = (plan: PlanItem) => {
  editingPlan.value = plan
  formState.name = plan.Name
  formState.status = plan.Status || 'active'
  formState.included_tokens = plan.IncludedTokens || 0
  formState.included_requests = plan.IncludedRequests || 0
  formState.reset_interval_days = plan.ResetIntervalDays || 30
  formState.price = Number.isFinite(plan.Price) ? plan.Price : 0
  formState.currency = plan.Currency || 'CNY'
  quotaType.value = plan.IncludedTokens > 0 ? 'tokens' : plan.IncludedRequests > 0 ? 'requests' : 'tokens'
  generalError.value = ''
  formKey.value += 1
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const validateForm = (state: PlanFormState): FormError[] => {
  const errors: FormError[] = []
  generalError.value = ''
  if (!state.name.trim()) {
    errors.push({ name: 'name', message: '请输入套餐名称' })
  }
  if (!state.status) {
    errors.push({ name: 'status', message: '请选择状态' })
  }
  const tokens = Number(state.included_tokens || 0)
  const requests = Number(state.included_requests || 0)
  const quotaField = quotaType.value === 'tokens' ? 'included_tokens' : 'included_requests'
  if (tokens < 0 || requests < 0) {
    errors.push({ name: quotaField, message: '配额不能为负数' })
  } else if (quotaType.value === 'tokens' && tokens <= 0) {
    errors.push({ name: 'included_tokens', message: '请输入包含 Token 数量' })
  } else if (quotaType.value === 'requests' && requests <= 0) {
    errors.push({ name: 'included_requests', message: '请输入包含请求次数' })
  }
  const resetInterval = Number(state.reset_interval_days || 0)
  if (!Number.isFinite(resetInterval) || resetInterval < 1 || resetInterval > 365) {
    errors.push({ name: 'reset_interval_days', message: '重置周期需在 1-365 天之间' })
  }
  const price = Number(state.price || 0)
  if (!Number.isFinite(price) || price < 0) {
    errors.push({ name: 'price', message: '套餐价格不能为负数' })
  }
  if (!state.currency) {
    errors.push({ name: 'currency', message: '请选择价格单位' })
  }
  return errors
}

const submitForm = async (event: FormSubmitEvent<PlanFormState>) => {
  generalError.value = ''
  isSubmitting.value = true
  const state = event.data
  const payload = {
    name: state.name.trim(),
    status: state.status,
    included_tokens: quotaType.value === 'tokens' ? Number(state.included_tokens || 0) : 0,
    included_requests: quotaType.value === 'requests' ? Number(state.included_requests || 0) : 0,
    reset_interval_days: Number(state.reset_interval_days || 0),
    price: Number(state.price || 0),
    currency: state.currency
  }
  try {
    if (editingPlan.value) {
      await $fetch(`/api/admin/plans/${editingPlan.value.ID}`, {
        method: 'PATCH',
        body: payload
      })
    } else {
      await $fetch('/api/admin/plans', {
        method: 'POST',
        body: payload
      })
    }
    await refreshList()
    isModalOpen.value = false
  } catch (error: unknown) {
    if (typeof error === 'object' && error && 'statusMessage' in error && typeof (error as { statusMessage?: string }).statusMessage === 'string') {
      generalError.value = (error as { statusMessage?: string }).statusMessage || '操作失败'
    } else if (typeof error === 'object' && error && 'message' in error && typeof (error as { message?: string }).message === 'string') {
      generalError.value = (error as { message?: string }).message || '操作失败'
    } else {
      generalError.value = '操作失败'
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>
