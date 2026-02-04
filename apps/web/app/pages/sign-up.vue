<template>
  <div class="min-h-screen flex items-center justify-center">
    <UCard>
      <div class="min-w-md space-y-2 text-center">
        <h2 class="text-2xl font-black">注册 Deepspace Workflow</h2>
        <div class="text-sm text-muted">开始你的科研写作协作之旅</div>
        <div>已有账号？<NuxtLink class="text-primary" to="/sign-in">登录</NuxtLink></div>
        <UButton class="w-full flex items-center justify-center" color="neutral" variant="outline"
          icon="i-lucide-github">使用 Github 注册</UButton>
        <UButton class="w-full flex items-center justify-center" color="neutral" variant="outline"
          icon="i-mage-google">使用 Google 注册</UButton>
      </div>

      <USeparator class="my-8" label="or" />

      <UForm class="space-y-2" @submit.prevent="submit">
        <UFormField label="组织名称" name="org">
          <UInput v-model="orgName" class="w-full" type="text" placeholder="可选，填写组织名称" />
        </UFormField>

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

        <UButton class="w-full flex items-center justify-center" color="primary" :loading="loading" type="submit">
          注册
        </UButton>
      </UForm>

      <p v-if="error" class="mt-4 text-sm text-red-500">{{ error }}</p>

      <div class="mt-6 text-sm text-muted">
        当您注册时，即表示您同意我们的
        <NuxtLink class="text-primary" to="/terms-of-service">服务条款</NuxtLink>
        及
        <NuxtLink class="text-primary" to="/privacy-policy">隐私政策</NuxtLink>。
      </div>
    </UCard>
  </div>
</template>
<script setup lang="ts">
definePageMeta({ layout: false })

const orgName = ref('')
const email = ref('')
const password = ref('')
const show = ref(false)
const loading = ref(false)
const error = ref('')

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
    await $fetch('/api/auth/register', {
      method: 'POST',
      body: { email: email.value, password: password.value, org_name: orgName.value }
    })
    await navigateTo('/projects')
  } catch (err: any) {
    error.value = err?.data?.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>
