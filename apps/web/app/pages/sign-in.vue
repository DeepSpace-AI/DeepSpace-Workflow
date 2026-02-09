<template>
  <div class="min-h-screen flex items-center justify-center bg-muted">
    <UCard>
      <div class="min-w-md space-y-2 text-center">
        <h2 class="text-2xl font-black">登录 Deepspace Workflow</h2>
        <div class="text-sm text-muted">欢迎回来👏</div>
        <div>还没有账号？<NuxtLink class="text-primary" to="/sign-up">注册</NuxtLink>
        </div>
        <UButton class="w-full flex items-center justify-center" color="neutral" variant="outline"
          icon="i-lucide-github">使用 Github 登录</UButton>
        <UButton class="w-full flex items-center justify-center" color="neutral" variant="outline" icon="i-mage-google">
          使用 Google 登录</UButton>
      </div>

      <USeparator class="my-8" label="or" />

      <UForm class="space-y-2" @submit.prevent="submit">
        <UFormField required label="邮箱" name="email">
          <UInput v-model="email" class="w-full" type="email" placeholder="请输入邮箱地址" />
        </UFormField>

        <UFormField required label="密码" name="password">
          <UInput class="w-full" v-model="password" placeholder="请输入密码" :type="show ? 'text' : 'password'"
            :ui="{ trailing: 'pe-1' }">
            <template #trailing>
              <UButton color="neutral" variant="link" size="sm" :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password"
                @click="show = !show" />
            </template>
          </UInput>
        </UFormField>

        <UFormField class="text-right">
          <UButton color="neutral" variant="link" size="sm" class="px-0" type="button" @click="openResetModal">
            忘记密码？
          </UButton>
        </UFormField>

        <UButton class="w-full flex items-center justify-center" color="primary" :loading="loading" type="submit">
          登录
        </UButton>
      </UForm>

      <p v-if="error" class="mt-4 text-sm text-red-500">{{ error }}</p>

      <div class="mt-6 text-sm text-muted">
        当您登录时，即表示您同意我们的
        <NuxtLink class="text-primary" to="/terms-of-service">服务条款</NuxtLink>
        及
        <NuxtLink class="text-primary" to="/privacy-policy">隐私政策</NuxtLink>。
      </div>
    </UCard>
  </div>

  <UModal v-model:open="resetModalOpen">
    <template #content>
      <div class="p-6 space-y-4">
        <h3 class="text-lg font-semibold">找回密码</h3>
        <p class="text-sm text-muted">输入你的注册邮箱，我们会发送重置密码链接。</p>
        <UForm class="space-y-3" @submit.prevent="submitResetRequest">
          <UFormField required label="邮箱" name="resetEmail">
            <UInput v-model="resetEmail" type="email" placeholder="请输入注册邮箱" />
          </UFormField>

          <p v-if="resetError" class="text-sm text-red-500">{{ resetError }}</p>
          <p v-else-if="resetSuccess" class="text-sm text-green-600 dark:text-green-400">{{ resetSuccess }}</p>

          <div class="flex justify-end gap-2">
            <UButton color="neutral" variant="outline" type="button" @click="resetModalOpen = false">取消</UButton>
            <UButton type="submit" :loading="resetLoading">发送重置邮件</UButton>
          </div>
        </UForm>
      </div>
    </template>
  </UModal>
</template>
<script setup lang="ts">
definePageMeta({ layout: false })
useHead({
  title: "登录 - Deepspace Workflow",
})

const email = ref('')
const password = ref('')
const show = ref(false)
const loading = ref(false)
const error = ref('')
const resetModalOpen = ref(false)
const resetEmail = ref('')
const resetLoading = ref(false)
const resetError = ref('')
const resetSuccess = ref('')

const validateEmail = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return '请输入邮箱地址'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(trimmed)) return '邮箱格式不正确'
  return ''
}

const validatePassword = (value: string) => {
  if (!value) return '请输入密码'
  if (value.length < 8) return '密码至少 8 位'
  if (!/[A-Z]/.test(value)) return '需包含至少 1 个大写字母'
  if (!/[a-z]/.test(value)) return '需包含至少 1 个小写字母'
  if (!/[0-9]/.test(value)) return '需包含至少 1 个数字'
  return ''
}

const submit = async () => {
  error.value = ''
  const emailError = validateEmail(email.value)
  if (emailError) {
    error.value = emailError
    return
  }
  const passwordError = validatePassword(password.value)
  if (passwordError) {
    error.value = passwordError
    return
  }
  loading.value = true
  try {
    await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })
    await navigateTo('/projects')
  } catch (err: any) {
    error.value = err?.data?.message || '登录失败'
  } finally {
    loading.value = false
  }
}

const openResetModal = () => {
  resetModalOpen.value = true
  resetEmail.value = email.value.trim()
  resetError.value = ''
  resetSuccess.value = ''
}

const submitResetRequest = async () => {
  resetError.value = ''
  resetSuccess.value = ''
  const emailError = validateEmail(resetEmail.value)
  if (emailError) {
    resetError.value = emailError
    return
  }

  resetLoading.value = true
  try {
    await $fetch('/api/auth/password-reset/request', {
      method: 'POST',
      body: { email: resetEmail.value.trim() }
    })
    resetSuccess.value = '重置邮件已发送，请检查你的邮箱。'
  } catch (err: any) {
    resetError.value = err?.data?.message || err?.data?.error || '发送失败，请稍后重试'
  } finally {
    resetLoading.value = false
  }
}
</script>
