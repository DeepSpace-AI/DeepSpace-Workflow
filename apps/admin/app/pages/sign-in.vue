<template>
  <UMain class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-950 p-4">
    <UCard class="w-full max-w-md">
      <template #header>
        <div class="flex items-center gap-3">
          <UIcon name="i-heroicons-rocket-launch" class="w-8 h-8 text-primary-500" />
          <div>
            <div class="text-lg font-semibold">DS Admin 登录</div>
            <div class="text-sm text-gray-500">仅管理员可访问控制台</div>
          </div>
        </div>
      </template>

      <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
        <div class="flex flex-col gap-2">
          <label class="text-sm text-gray-600">邮箱</label>
          <UInput
            v-model="form.email"
            type="email"
            placeholder="admin@deepspace.ai"
            icon="i-heroicons-envelope"
            autocomplete="email"
          />
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm text-gray-600">密码</label>
          <UInput
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            icon="i-heroicons-lock-closed"
            autocomplete="current-password"
          />
        </div>
        <div v-if="errorMessage" class="text-sm text-red-600">{{ errorMessage }}</div>
        <UButton type="submit" color="primary" :loading="isSubmitting" block>登录</UButton>
      </form>

      <template #footer>
        <div class="text-xs text-gray-500">如需开通管理员权限，请联系系统管理员</div>
      </template>
    </UCard>
  </UMain>
</template>

<script setup lang="ts">
const form = reactive({
  email: '',
  password: ''
})
const isSubmitting = ref(false)
const errorMessage = ref('')
const sessionCookie = useCookie('dsp_session')

const redirectIfAuthed = async () => {
  if (!sessionCookie.value) return
  const me = await $fetch('/api/auth/me').catch(() => null)
  if (me) {
    return navigateTo('/')
  }
}

await redirectIfAuthed()

const handleSubmit = async () => {
  errorMessage.value = ''
  if (!form.email || !form.password) {
    errorMessage.value = '请输入邮箱和密码'
    return
  }
  isSubmitting.value = true
  try {
    await $fetch('/api/auth/login', {
      method: 'POST',
      body: {
        email: form.email,
        password: form.password
      }
    })
    await navigateTo('/')
  } catch (error) {
    const fetchError = error as { data?: { message?: string; error?: string }; statusMessage?: string }
    const message =
      fetchError?.data?.message ||
      fetchError?.data?.error ||
      fetchError?.statusMessage ||
      '登录失败'
    errorMessage.value = message
  } finally {
    isSubmitting.value = false
  }
}
</script>
