<template>
  <div class="space-y-4">
    <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">用户管理</h3>
          <UBadge color="neutral" variant="soft">共 {{ totalCount }} 人</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UInput v-model="searchTerm" placeholder="搜索邮箱" icon="i-heroicons-magnifying-glass" class="w-60" />
          <USelect v-model="roleFilter" :items="roleOptions" class="w-28" />
          <USelect v-model="statusFilter" :items="statusOptions" class="w-28" />
          <USelect v-model="pageSize" :items="pageSizeOptions" class="w-24" />
          <UButton
            icon="i-heroicons-arrow-path"
            label="刷新"
            color="neutral"
            variant="outline"
            :loading="isLoading"
            :disabled="isLoading"
            @click="refreshList"
          />
          <UButton icon="i-heroicons-plus" label="添加用户" color="primary" @click="openCreateModal" />
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <UTable :data="userItems" :columns="columns" :loading="isLoading" />
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-500">共 {{ totalCount }} 条</div>
        <UPagination v-model:page="page" :total="totalCount" :items-per-page="pageSize" :sibling-count="1" show-edges />
      </div>
    </div>
    </UCard>

    <UModal v-model:open="isModalOpen" :title="modalTitle">
    <template #body>
      <div class="grid gap-4">
        <UFormField label="邮箱" required>
          <UInput  v-model="formState.email" class="w-full" placeholder="user@deepspace.ai" />
        </UFormField>
        <UFormField label="密码" :help="passwordHelp">
          <UInput v-model="formState.password" class="w-full" type="password" placeholder="请输入密码" />
        </UFormField>
        <div class="grid gap-4 sm:grid-cols-2">
          <UFormField label="角色">
            <USelect v-model="formState.role" class="w-full" :items="roleOptions" />
          </UFormField>
          <UFormField label="状态">
            <USelect v-model="formState.status" class="w-full" :items="statusOptions" />
          </UFormField>
        </div>
        <div class="grid gap-4 sm:grid-cols-2">
          <UFormField label="显示名称">
            <UInput v-model="formState.displayName" class="w-full" placeholder="如：张三" />
          </UFormField>
          <UFormField label="姓名">
            <UInput v-model="formState.fullName" class="w-full" placeholder="如：张三" />
          </UFormField>
        </div>
        <div v-if="formError" class="text-sm text-red-600">{{ formError }}</div>
        <div v-if="isDetailLoading" class="text-sm text-gray-500">正在加载用户信息...</div>
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
        将删除用户 <span class="font-medium text-gray-900">{{ deleteTarget?.email || '' }}</span>，此操作不可撤销。
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

type UserRow = {
  id: number
  email: string
  username: string
  role: string
  status: string
  last_login_at?: string | null
  created_at?: string | null
}

type UserListResponse = {
  items: UserRow[]
  total: number
  page: number
  page_size: number
}

type UserProfile = {
  UserID?: number | null
  DisplayName?: string | null
  FullName?: string | null
  Title?: string | null
  AvatarURL?: string | null
  Bio?: string | null
  Phone?: string | null
  CreatedAt?: string | null
  UpdatedAt?: string | null
  display_name?: string | null
  full_name?: string | null
  title?: string | null
  avatar_url?: string | null
  bio?: string | null
  phone?: string | null
}

type UserDetailResponse = {
  user: UserRow
  profile?: UserProfile
}

const page = ref(1)
const pageSize = ref(10)
const searchTerm = ref('')
const roleFilter = ref('all')
const statusFilter = ref('all')
const isModalOpen = ref(false)
const isDeleteOpen = ref(false)
const isSubmitting = ref(false)
const isDeleting = ref(false)
const isDetailLoading = ref(false)
const formError = ref('')
const deleteError = ref('')
const editingUser = ref<UserRow | null>(null)
const deleteTarget = ref<UserRow | null>(null)
const formState = ref({
  email: '',
  password: '',
  role: 'user',
  status: 'active',
  displayName: '',
  fullName: ''
})

const roleOptions = [
  { label: '全部角色', value: 'all' },
  { label: '管理员', value: 'admin' },
  { label: '开发者', value: 'developer' },
  { label: '普通用户', value: 'user' }
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

const modalTitle = computed(() => (editingUser.value ? '编辑用户' : '添加用户'))
const passwordHelp = computed(() => (editingUser.value ? '如需重置密码请填写' : '至少 8 位字符'))

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  search: searchTerm.value || undefined,
  role: roleFilter.value === 'all' ? undefined : roleFilter.value,
  status: statusFilter.value === 'all' ? undefined : statusFilter.value
}))

const { data: listData, pending: isLoading, refresh: refreshList } = await useFetch<UserListResponse>('/api/admin/users', {
  query: queryParams,
  default: () => ({
    items: [],
    total: 0,
    page: 1,
    page_size: pageSize.value
  })
})

const userItems = computed(() => listData.value?.items ?? [])
const totalCount = computed(() => listData.value?.total ?? 0)

watch([searchTerm, roleFilter, statusFilter, pageSize], () => {
  page.value = 1
})

const UBadge = resolveComponent('UBadge')
const UButton = resolveComponent('UButton')

const formatTime = (value?: string | null) => {
  if (!value) return '—'
  const time = new Date(value)
  if (Number.isNaN(time.getTime())) return '—'
  return time.toLocaleString('zh-CN', { hour12: false })
}

const resetForm = () => {
  formState.value = {
    email: '',
    password: '',
    role: 'user',
    status: 'active',
    displayName: '',
    fullName: ''
  }
  formError.value = ''
}

const initFormFromRow = (row: UserRow) => {
  formState.value = {
    email: row.email || '',
    password: '',
    role: row.role || 'user',
    status: row.status || 'active',
    displayName: '',
    fullName: ''
  }
  formError.value = ''
}

const normalizeProfile = (profile?: UserProfile | null) => ({
  displayName: profile?.display_name ?? profile?.DisplayName ?? '',
  fullName: profile?.full_name ?? profile?.FullName ?? ''
})

const loadUserDetail = async (userId: number) => {
  isDetailLoading.value = true
  formError.value = ''
  try {
    const detail = await $fetch<UserDetailResponse>(`/api/admin/users/${userId}`)
    const profile = normalizeProfile(detail.profile)
    if (detail?.user?.email) {
      formState.value.email = detail.user.email || ''
    }
    formState.value.password = ''
    if (detail?.user?.role) {
      formState.value.role = detail.user.role || 'user'
    }
    if (detail?.user?.status) {
      formState.value.status = detail.user.status || 'active'
    }
    formState.value.displayName = profile.displayName
    formState.value.fullName = profile.fullName
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    formError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '获取用户详情失败'
  } finally {
    isDetailLoading.value = false
  }
}

const openCreateModal = () => {
  editingUser.value = null
  resetForm()
  isModalOpen.value = true
}

const openEditModal = async (row: UserRow) => {
  editingUser.value = row
  initFormFromRow(row)
  isModalOpen.value = true
  await loadUserDetail(row.id)
}

const submitForm = async () => {
  formError.value = ''
  if (!formState.value.email) {
    formError.value = '请输入邮箱'
    return
  }
  if (!editingUser.value && !formState.value.password) {
    formError.value = '请输入密码'
    return
  }
  isSubmitting.value = true
  try {
    const profilePayload =
      formState.value.displayName || formState.value.fullName
        ? {
            display_name: formState.value.displayName || undefined,
            full_name: formState.value.fullName || undefined
          }
        : undefined
    const payload = {
      email: formState.value.email,
      role: formState.value.role,
      status: formState.value.status,
      ...(formState.value.password ? { password: formState.value.password } : {}),
      ...(profilePayload ? { profile: profilePayload } : {})
    }
    if (editingUser.value) {
      await $fetch(`/api/admin/users/${editingUser.value.id}`, {
        method: 'PATCH',
        body: payload
      })
    } else {
      await $fetch('/api/admin/users', {
        method: 'POST',
        body: payload
      })
    }
    isModalOpen.value = false
    await refreshList()
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    formError.value =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      (editingUser.value ? '更新用户失败' : '创建用户失败')
  } finally {
    isSubmitting.value = false
  }
}

const openDeleteModal = (row: UserRow) => {
  deleteTarget.value = row
  deleteError.value = ''
  isDeleteOpen.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  isDeleting.value = true
  deleteError.value = ''
  try {
    await $fetch(`/api/admin/users/${deleteTarget.value.id}`, {
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
      '删除用户失败'
  } finally {
    isDeleting.value = false
  }
}

const columns = computed<TableColumn<UserRow>[]>(() => [
  {
    accessorKey: 'id',
    header: 'ID',
    meta: { class: { th: 'w-20' } }
  },
  {
    accessorKey: 'email',
    header: '邮箱'
  },
  {
    accessorKey: 'username',
    header: '用户名'
  },
  {
    accessorKey: 'role',
    header: '角色',
    cell: ({ row }) => {
      const roleValue = String(row.getValue('role') || '')
      const color = roleValue === 'admin' ? 'primary' : roleValue === 'developer' ? 'info' : 'neutral'
      const label = roleValue === 'admin' ? '管理员' : roleValue === 'developer' ? '开发者' : '普通用户'
      return h(UBadge, { color, variant: 'subtle' }, () => label)
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
    accessorKey: 'last_login_at',
    header: '最近登录',
    cell: ({ row }) => formatTime(row.getValue('last_login_at') as string | null | undefined)
  },
  {
    accessorKey: 'created_at',
    header: '创建时间',
    cell: ({ row }) => formatTime(row.getValue('created_at') as string | null | undefined)
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
