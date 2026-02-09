<template>
  <div class="space-y-4">
    <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">钱包管理</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 条</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索用户ID" icon="i-heroicons-magnifying-glass" class="w-60" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
          <UButton icon="i-heroicons-plus" label="系统充值" color="primary" @click="openTopupModal" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="walletItems" :columns="columns" :loading="isLoading" />
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-500">共 {{ totalCount }} 条</div>
        <UPagination v-model:page="page" :total="totalCount" :items-per-page="pageSize" :sibling-count="1" show-edges />
      </div>
    </div>
    </UCard>

    <UModal v-model:open="isTopupOpen" title="系统充值">
      <template #body>
        <UForm :state="topupFormState" :validate="validateTopup" class="grid gap-4" @submit="submitTopup">
          <UFormField label="用户搜索" name="user" required>
            <div class="grid gap-2">
              <div class="flex flex-wrap items-center gap-2">
                <UInput v-model="userSearchKeyword" class="flex-1 min-w-[220px]" placeholder="输入用户邮箱或ID" />
                <UButton color="neutral" variant="outline" :loading="isUserSearching" @click="searchUsers">搜索</UButton>
                <UButton v-if="selectedUser" color="neutral" variant="ghost" @click="clearSelectedUser">重新选择</UButton>
              </div>
              <div v-if="selectedUser" class="rounded-md border border-gray-200 dark:border-gray-800 px-3 py-2 text-sm">
                <div class="text-gray-500">已选择用户</div>
                <div class="font-medium text-gray-900 dark:text-gray-100">ID {{ selectedUser.id }} · {{ selectedUser.email || '—' }}</div>
                <div class="text-gray-500">状态：{{ formatStatus(selectedUser.status) }}</div>
              </div>
              <div v-if="userSearchError" class="text-sm text-red-600">{{ userSearchError }}</div>
              <div v-else-if="isUserSearching" class="text-sm text-gray-500">正在搜索用户...</div>
              <div v-else-if="searchPerformed && !userSearchResults.length" class="text-sm text-gray-500">未找到匹配用户</div>
              <div v-if="userSearchResults.length" class="rounded-md border border-gray-200 dark:border-gray-800 divide-y divide-gray-200 dark:divide-gray-800">
                <div v-for="user in userSearchResults" :key="user.id" class="flex items-center justify-between px-3 py-2 text-sm">
                  <div class="flex flex-col">
                    <span class="font-medium text-gray-900 dark:text-gray-100">ID {{ user.id }}</span>
                    <span class="text-gray-500">{{ user.email || '—' }}</span>
                  </div>
                  <UButton size="xs" color="primary" variant="soft" @click="confirmUser(user)">确认</UButton>
                </div>
              </div>
            </div>
          </UFormField>
          <div class="grid gap-4 sm:grid-cols-2">
            <UFormField label="充值金额" name="amount" required :help="amountHelp">
              <UInput v-model.number="topupFormState.amount" type="number" step="0.01" class="w-full" placeholder="例如 100.00" />
            </UFormField>
            <UFormField label="币种" name="currency" required>
              <USelect v-model="topupFormState.currency" class="w-full" :items="currencyOptions" />
            </UFormField>
          </div>
          <UFormField label="引用ID" name="ref_id" :help="refIdHelp">
            <UInput v-model="topupFormState.ref_id" class="w-full" placeholder="用于幂等处理，可选" />
          </UFormField>
          <UFormField label="备注" name="remark">
            <UInput v-model="topupFormState.remark" class="w-full" placeholder="可填写原因或工单号" />
          </UFormField>
          <div v-if="topupError" class="text-sm text-red-600">{{ topupError }}</div>
          <div class="flex items-center justify-end gap-2">
            <UButton color="neutral" variant="outline" type="button" :disabled="isTopupSubmitting" @click="closeTopupModal">取消</UButton>
            <UButton color="primary" type="submit" :loading="isTopupSubmitting">确认充值</UButton>
          </div>
        </UForm>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { FormError, FormSubmitEvent, TableColumn } from '@nuxt/ui'

type WalletUser = {
  id: number
  email: string
  status: string
  role: string
  created_at: string
}

type WalletRow = {
  user: WalletUser
  balance: number
  frozen_balance: number
  updated_at: string
}

type WalletListResponse = {
  items: WalletRow[]
  total: number
  page: number
  page_size: number
}

type TopupFormState = {
  amount: number | null
  currency: string
  ref_id: string
  remark: string
}

type UserSearchItem = {
  id: number
  email: string
  status: string
  role: string
}

type UserSearchResponse = {
  items: UserSearchItem[]
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

const { data: listData, pending: isLoading, refresh: refreshList } = await useFetch<WalletListResponse>('/api/admin/billing/wallets', {
  query: queryParams,
  default: () => ({
    items: [],
    total: 0,
    page: 1,
    page_size: pageSize.value
  })
})

const walletItems = computed(() => listData.value?.items ?? [])
const totalCount = computed(() => listData.value?.total ?? 0)

const UBadge = resolveComponent('UBadge')

const currencyOptions = [
  { label: 'USD', value: 'USD' },
  { label: 'CNY', value: 'CNY' }
]
const amountHelp = '可输入负数表示冲正'
const refIdHelp = '建议填写唯一引用ID避免重复充值'
const isTopupOpen = ref(false)
const isTopupSubmitting = ref(false)
const topupError = ref('')
const userSearchKeyword = ref('')
const userSearchResults = ref<UserSearchItem[]>([])
const userSearchError = ref('')
const isUserSearching = ref(false)
const searchPerformed = ref(false)
const selectedUser = ref<UserSearchItem | null>(null)
const topupFormState = ref<TopupFormState>({
  amount: null,
  currency: 'USD',
  ref_id: '',
  remark: ''
})

const resetTopupForm = () => {
  topupFormState.value = {
    amount: null,
    currency: 'USD',
    ref_id: '',
    remark: ''
  }
  topupError.value = ''
  userSearchKeyword.value = searchTerm.value.trim()
  userSearchResults.value = []
  userSearchError.value = ''
  isUserSearching.value = false
  searchPerformed.value = false
  selectedUser.value = null
}

const openTopupModal = () => {
  resetTopupForm()
  isTopupOpen.value = true
}

const closeTopupModal = () => {
  isTopupOpen.value = false
}

const validateTopup = (state: TopupFormState): FormError[] => {
  const errors: FormError[] = []
  if (!selectedUser.value) {
    errors.push({ name: 'user', message: '请选择用户后再提交' })
  }
  const amountValue = Number(state.amount ?? 0)
  if (!Number.isFinite(amountValue) || amountValue === 0) {
    errors.push({ name: 'amount', message: '请输入有效的充值金额' })
  }
  if (!state.currency) {
    errors.push({ name: 'currency', message: '请选择币种' })
  }
  return errors
}

const submitTopup = async (event: FormSubmitEvent<TopupFormState>) => {
  topupError.value = ''
  isTopupSubmitting.value = true
  const state = event.data
  if (!selectedUser.value) {
    topupError.value = '请先确认用户'
    isTopupSubmitting.value = false
    return
  }
  if (!state.ref_id.trim()) {
    const now = new Date()
    const timePart = now.toISOString().replace(/[-:.TZ]/g, '')
    const randomPart = Math.random().toString(36).slice(2, 10).toUpperCase()
    state.ref_id = `TOPUP-${timePart}-${randomPart}`
    topupFormState.value.ref_id = state.ref_id
  }
  const payload = {
    user_id: selectedUser.value.id,
    amount: Number(state.amount || 0),
    currency: state.currency,
    ...(state.ref_id.trim() ? { ref_id: state.ref_id.trim() } : {}),
    ...(state.remark.trim() ? { metadata: { remark: state.remark.trim() } } : {})
  }
  try {
    await $fetch('/api/admin/billing/topups', {
      method: 'POST',
      body: payload
    })
    isTopupOpen.value = false
    await refreshList()
  } catch (error: unknown) {
    if (typeof error === 'object' && error && 'statusMessage' in error && typeof (error as { statusMessage?: string }).statusMessage === 'string') {
      topupError.value = (error as { statusMessage?: string }).statusMessage || '系统充值失败'
    } else if (typeof error === 'object' && error && 'message' in error && typeof (error as { message?: string }).message === 'string') {
      topupError.value = (error as { message?: string }).message || '系统充值失败'
    } else {
      topupError.value = '系统充值失败'
    }
  } finally {
    isTopupSubmitting.value = false
  }
}

const searchUsers = async () => {
  userSearchError.value = ''
  userSearchResults.value = []
  searchPerformed.value = false
  const keyword = userSearchKeyword.value.trim()
  if (!keyword) {
    userSearchError.value = '请输入用户邮箱或ID'
    return
  }
  isUserSearching.value = true
  try {
    const data = await $fetch<UserSearchResponse>('/api/admin/users', {
      query: {
        page: 1,
        page_size: 5,
        keyword: keyword
      }
    })
    userSearchResults.value = data?.items ?? []
    searchPerformed.value = true
  } catch (error: unknown) {
    if (typeof error === 'object' && error && 'statusMessage' in error && typeof (error as { statusMessage?: string }).statusMessage === 'string') {
      userSearchError.value = (error as { statusMessage?: string }).statusMessage || '搜索用户失败'
    } else if (typeof error === 'object' && error && 'message' in error && typeof (error as { message?: string }).message === 'string') {
      userSearchError.value = (error as { message?: string }).message || '搜索用户失败'
    } else {
      userSearchError.value = '搜索用户失败'
    }
  } finally {
    isUserSearching.value = false
  }
}

const confirmUser = (user: UserSearchItem) => {
  selectedUser.value = user
  userSearchResults.value = []
  searchPerformed.value = false
  userSearchError.value = ''
}

const clearSelectedUser = () => {
  selectedUser.value = null
}

const formatAmount = (value: number) => value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
const formatTime = (value: string) => {
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}
const formatStatus = (value: string) => {
  if (value === 'active') return '启用'
  if (value === 'disabled') return '禁用'
  return value || '—'
}
const formatRole = (value: string) => {
  if (value === 'admin') return '管理员'
  if (value === 'user') return '用户'
  return value || '—'
}

const columns = computed<TableColumn<WalletRow>[]>(() => [
  {
    accessorKey: 'user',
    header: '用户ID',
    meta: { class: { th: 'w-28' } },
    cell: ({ row }) => row.original.user.id
  },
  {
    accessorKey: 'user',
    header: '邮箱',
    cell: ({ row }) => row.original.user.email || '—'
  },
  {
    accessorKey: 'user',
    header: '状态',
    cell: ({ row }) => formatStatus(row.original.user.status)
  },
  {
    accessorKey: 'user',
    header: '角色',
    cell: ({ row }) => formatRole(row.original.user.role)
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
