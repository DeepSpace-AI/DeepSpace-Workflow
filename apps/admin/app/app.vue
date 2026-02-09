<template>
  <UApp :locale="zh_cn" class="min-h-screen">
    <NuxtLayout v-if="!isAuthRoute" name="dashboard">
      <NuxtPage />
    </NuxtLayout>
    <NuxtPage v-else />
  </UApp>
</template>

<script setup lang="ts">
import { zh_cn } from '@nuxt/ui/locale'
const route = useRoute()
const isAuthRoute = computed(() => route.path === '/sign-in')
</script>
