<template>
  <div class="min-h-screen flex items-center justify-center">
    <UCard>
      <div class="min-w-md space-y-2 text-center">
        <h2 class="text-2xl font-black">退出登录</h2>
        <div class="text-sm text-muted">正在安全退出，请稍候...</div>
      </div>
    </UCard>
  </div>
</template>
<script setup lang="ts">
definePageMeta({ layout: false })

onMounted(async () => {
  try {
    await $fetch('/api/auth/logout', { method: 'POST' })
  } finally {
    await navigateTo('/')
  }
})
</script>
