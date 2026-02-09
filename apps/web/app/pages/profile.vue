<template>
  <UContainer>
    <div class="my-12 space-y-8">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-primary">个人资料</p>
          <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100 md:text-4xl">
            个人资料与研究档案
          </h1>
          <p class="max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">
            管理你的个人信息、研究方向与团队归属，让协作更清晰。
          </p>
        </div>
        <div class="flex w-full flex-col gap-3 sm:flex-row md:w-auto">
          <UButton color="primary" size="lg" class="w-full sm:w-auto" :loading="saving" @click="saveProfile">
            保存资料
          </UButton>
        </div>
      </div>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex items-center gap-4">
            <UAvatar
              size="xl"
              :src="form.avatar_url || undefined"
              icon="i-lucide-user"
            />
            <div>
              <p class="text-lg font-semibold text-slate-900 dark:text-white">
                {{ form.display_name || form.full_name || '未设置姓名' }}
              </p>
              <p class="text-sm text-slate-500 dark:text-slate-400">
                {{ form.title || '研究成员' }}
              </p>
            </div>
          </div>
          <div class="flex w-full flex-col gap-2 sm:w-auto">
            <UInput v-model="form.avatar_url" placeholder="头像 URL" size="lg" class="w-full" />
            <span class="text-xs text-slate-500 dark:text-slate-400">
              建议使用 256px 以上清晰头像。
            </span>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="space-y-4">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">基本信息</h2>
          <div class="grid gap-4 md:grid-cols-2">
            <UFormField label="显示名">
              <UInput v-model="form.display_name" placeholder="用于系统内展示" class="w-full" />
            </UFormField>
            <UFormField label="真实姓名">
              <UInput v-model="form.full_name" placeholder="身份证件姓名" class="w-full" />
            </UFormField>
            <UFormField label="职位 / 头衔">
              <UInput v-model="form.title" placeholder="Research Engineer" class="w-full" />
            </UFormField>
            <UFormField label="联系电话">
              <UInput v-model="form.phone" placeholder="可选" class="w-full" />
            </UFormField>
          </div>
          <UFormField label="个人简介">
            <UTextarea v-model="form.bio" :rows="4" placeholder="简要描述你的研究方向与经验" class="w-full" />
          </UFormField>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="space-y-3">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">研究偏好</h2>
          <p class="text-sm text-slate-600 dark:text-slate-300">
            即将开放：研究领域偏好、常用模型、默认项目类型等设置。
          </p>
        </div>
      </UCard>

    </div>
  </UContainer>

  <UModal v-model:open="forceModalOpen" :dismissible="false">
    <template #content>
      <div class="p-6 space-y-5">
        <div class="space-y-1">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-slate-100">请先完善个人信息</h2>
          <p class="text-sm text-slate-600 dark:text-slate-300">
            为保障系统功能可用，你需要先填写显示名和真实姓名，完成后才能进入系统。
          </p>
        </div>

        <div class="space-y-3">
          <UFormField label="显示名" required>
            <UInput v-model="form.display_name" placeholder="用于系统内展示" class="w-full" />
          </UFormField>
          <UFormField label="真实姓名" required>
            <UInput v-model="form.full_name" placeholder="身份证件姓名" class="w-full" />
          </UFormField>
        </div>

        <div class="flex justify-end">
          <UButton color="primary" :loading="saving" @click="saveProfile">
            保存并继续
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
useHead({
  title: "个人资料 - Deepspace Workflow",
})

type UserMeResponse = {
  user: { id: number; email: string; status: string }
  profile: {
    display_name?: string | null
    full_name?: string | null
    title?: string | null
    avatar_url?: string | null
    bio?: string | null
    phone?: string | null
    DisplayName?: string | null
    FullName?: string | null
    Title?: string | null
    AvatarURL?: string | null
    Bio?: string | null
    Phone?: string | null
  }
  settings: {
    theme?: string | null
    locale?: string | null
    timezone?: string | null
  }
}

const requestHeaders = useRequestHeaders(['cookie'])
const route = useRoute()
const { data, refresh } = await useAsyncData<UserMeResponse>('user-me', () =>
  $fetch('/api/users/me', { headers: requestHeaders })
)

const forceInit = computed(() => String(route.query.force_init || '') === '1')
const profileIncomplete = computed(() => {
  const profile = data.value?.profile
  if (!profile) return false
  const displayName = String(profile.display_name || profile.DisplayName || '').trim()
  const fullName = String(profile.full_name || profile.FullName || '').trim()
  return !(displayName && fullName)
})
const mustForceInit = computed(() => forceInit.value || profileIncomplete.value)
const redirectPath = computed(() => {
  const raw = String(route.query.redirect || '').trim()
  return raw.startsWith('/') ? raw : '/projects'
})
const forceModalOpen = ref(false)

const form = reactive({
  display_name: '',
  full_name: '',
  title: '',
  avatar_url: '',
  bio: '',
  phone: ''
})

watch(
  () => data.value,
  (value) => {
    if (!value?.profile) return
    form.display_name = value.profile.display_name || value.profile.DisplayName || ''
    form.full_name = value.profile.full_name || value.profile.FullName || ''
    form.title = value.profile.title || value.profile.Title || ''
    form.avatar_url = value.profile.avatar_url || value.profile.AvatarURL || ''
    form.bio = value.profile.bio || value.profile.Bio || ''
    form.phone = value.profile.phone || value.profile.Phone || ''
  },
  { immediate: true }
)

watch(
  () => mustForceInit.value,
  (value) => {
    forceModalOpen.value = value
  },
  { immediate: true }
)

watch(forceModalOpen, (value) => {
  if (mustForceInit.value && !value) {
    forceModalOpen.value = true
  }
})

const saving = ref(false)
const toast = useToast()

function isProfileInitialized() {
  return Boolean(form.display_name.trim() && form.full_name.trim())
}

const saveProfile = async () => {
  const shouldRedirectAfterSave = mustForceInit.value

  if (!isProfileInitialized()) {
    toast.add({ title: '请先填写显示名和真实姓名', color: 'warning' })
    forceModalOpen.value = forceInit.value
    return
  }

  saving.value = true
  try {
    await $fetch('/api/users/me', {
      method: 'PATCH',
      body: {
        profile: {
          display_name: form.display_name.trim(),
          full_name: form.full_name.trim(),
          title: form.title,
          avatar_url: form.avatar_url,
          bio: form.bio,
          phone: form.phone
        }
      }
    })
    toast.add({ title: '保存成功', color: 'green' })
    await refresh()
    if (shouldRedirectAfterSave) {
      forceModalOpen.value = false
      await navigateTo(redirectPath.value)
    }
  } catch (err: any) {
    toast.add({ title: '保存失败', description: err?.data?.message || err?.message, color: 'red' })
  } finally {
    saving.value = false
  }
}
</script>
