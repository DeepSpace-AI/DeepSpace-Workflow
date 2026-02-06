<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>

<template>
  <div class="flex flex-wrap items-start justify-between gap-4">
    <div class="space-y-1">
      <h1 class="text-2xl font-semibold text-white">{{ title }}</h1>
      <p v-if="description" class="text-sm text-white/60">{{ description }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
