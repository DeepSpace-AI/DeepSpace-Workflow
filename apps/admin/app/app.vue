<template>
  <UApp class="min-h-screen">
    <NuxtLayout v-if="!isAuthRoute" name="dashboard">
      <NuxtPage />
    </NuxtLayout>
    <NuxtPage v-else />
  </UApp>
</template>

<script setup lang="ts">
const route = useRoute()

const isAuthRoute = computed(() => route.path === '/sign-in')
</script>
