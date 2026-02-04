<template>
  <UContainer>
    <div class="my-12 space-y-8">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-primary">设置</p>
          <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100 md:text-4xl">
            账户安全与偏好
          </h1>
          <p class="max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">
            管理密码安全、系统偏好与组织入口，确保团队协作稳定可控。
          </p>
        </div>
      </div>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="space-y-4">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">账户安全</h2>
          <div class="grid gap-4 md:grid-cols-3">
            <UFormField label="当前密码">
              <UInput v-model="passwordForm.old" type="password" placeholder="请输入当前密码" class="w-full" />
            </UFormField>
            <UFormField label="新密码">
              <UInput v-model="passwordForm.next" type="password" placeholder="至少 8 位" class="w-full" />
            </UFormField>
            <UFormField label="确认新密码">
              <UInput v-model="passwordForm.confirm" type="password" placeholder="再次输入新密码" class="w-full" />
            </UFormField>
          </div>
          <div class="flex items-center gap-3">
            <UButton color="primary" :loading="changingPassword" @click="changePassword">
              更新密码
            </UButton>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="space-y-4">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">体验偏好</h2>
          <div class="grid gap-4 md:grid-cols-3">
            <UFormField label="主题">
              <USelectMenu v-model="settingsForm.theme" :options="themeOptions" class="w-full" />
            </UFormField>
            <UFormField label="语言">
              <USelectMenu v-model="settingsForm.locale" :options="localeOptions" class="w-full" />
            </UFormField>
            <UFormField label="时区">
              <UInput v-model="settingsForm.timezone" placeholder="Asia/Shanghai" class="w-full" />
            </UFormField>
          </div>
          <div class="flex items-center gap-3">
            <UButton color="primary" :loading="savingSettings" @click="saveSettings">
              保存偏好
            </UButton>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="space-y-3">
          <h2 class="text-lg font-semibold text-slate-900 dark:text-white">组织与 API Key</h2>
          <p class="text-sm text-slate-600 dark:text-slate-300">
            管理组织信息、成员角色与 API Key（需要管理员权限）。
          </p>
          <div class="flex flex-wrap gap-3">
            <UButton color="neutral" variant="soft" to="/admin">
              前往组织控制台
            </UButton>
          </div>
        </div>
      </UCard>
    </div>
  </UContainer>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'

type UserMeResponse = {
  user: { id: number; email: string; status: string }
  profile: {
    display_name?: string | null
    full_name?: string | null
    title?: string | null
    avatar_url?: string | null
    bio?: string | null
    phone?: string | null
  }
  settings: {
    theme?: string | null
    locale?: string | null
    timezone?: string | null
  }
}

const requestHeaders = useRequestHeaders(['cookie'])
const { data, refresh } = await useAsyncData<UserMeResponse>('user-me-settings', () =>
  $fetch('/api/users/me', { headers: requestHeaders })
)

const themeOptions = [
  { label: '系统', value: 'system' },
  { label: '明亮', value: 'light' },
  { label: '暗色', value: 'dark' }
]
const localeOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' }
]

const settingsForm = reactive({
  theme: 'system',
  locale: 'zh-CN',
  timezone: 'Asia/Shanghai'
})

watch(
  () => data.value,
  (value) => {
    if (!value?.settings) return
    settingsForm.theme = value.settings.theme || 'system'
    settingsForm.locale = value.settings.locale || 'zh-CN'
    settingsForm.timezone = value.settings.timezone || 'Asia/Shanghai'
  },
  { immediate: true }
)

const savingSettings = ref(false)
const toast = useToast()

const saveSettings = async () => {
  savingSettings.value = true
  try {
    await $fetch('/api/users/me', {
      method: 'PATCH',
      body: {
        settings: {
          theme: settingsForm.theme,
          locale: settingsForm.locale,
          timezone: settingsForm.timezone
        }
      }
    })
    toast.add({ title: '偏好已保存', color: 'green' })
    await refresh()
  } catch (err: any) {
    toast.add({ title: '保存失败', description: err?.data?.message || err?.message, color: 'red' })
  } finally {
    savingSettings.value = false
  }
}

const passwordForm = reactive({
  old: '',
  next: '',
  confirm: ''
})
const changingPassword = ref(false)

const changePassword = async () => {
  if (!passwordForm.old || !passwordForm.next) {
    toast.add({ title: '请输入完整密码', color: 'red' })
    return
  }
  if (passwordForm.next.length < 8) {
    toast.add({ title: '新密码至少 8 位', color: 'red' })
    return
  }
  if (passwordForm.next !== passwordForm.confirm) {
    toast.add({ title: '两次输入的新密码不一致', color: 'red' })
    return
  }

  changingPassword.value = true
  try {
    await $fetch('/api/users/me/password', {
      method: 'POST',
      body: {
        old_password: passwordForm.old,
        new_password: passwordForm.next
      }
    })
    toast.add({ title: '密码已更新', color: 'green' })
    passwordForm.old = ''
    passwordForm.next = ''
    passwordForm.confirm = ''
  } catch (err: any) {
    toast.add({ title: '更新失败', description: err?.data?.message || err?.message, color: 'red' })
  } finally {
    changingPassword.value = false
  }
}
</script>
