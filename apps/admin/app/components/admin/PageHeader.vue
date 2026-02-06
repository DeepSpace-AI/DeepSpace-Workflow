<script setup lang="ts">
defineProps<{ title: string; description?: string; badge?: string }>()
</script>

<template>
  <div class="flex flex-wrap items-start justify-between gap-4">
    <div>
      <div class="flex items-center gap-2">
        <h1 class="text-2xl font-semibold text-slate-900">{{ title }}</h1>
        <UBadge v-if="badge" color="neutral" variant="subtle">{{ badge }}</UBadge>
      </div>
      <p v-if="description" class="mt-1 text-sm text-slate-500">
        {{ description }}
      </p>
    </div>
    <div class="flex items-center gap-2">
      <slot />
    </div>
  </div>
</template>
