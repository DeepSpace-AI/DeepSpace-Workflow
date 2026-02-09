<script setup lang="ts">
definePageMeta({
  layout: "default",
});
useHead({
  title: "工作流中心 - Deepspace Workflow",
})
</script>

<template>
  <div class="min-h-[calc(100vh-var(--ui-header-height))] bg-white dark:bg-neutral-950">
    <UContainer class="py-10">
      <UCard class="border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900">
        <template #header>
          <div class="space-y-1">
            <h1 class="text-2xl font-bold text-neutral-900 dark:text-neutral-100">工作流中心</h1>
            <p class="text-sm text-neutral-600 dark:text-neutral-400">
              当前提供最小入口，你可以先在项目内管理和调用工作流。
            </p>
          </div>
        </template>

        <div class="space-y-4">
          <div class="rounded-lg border border-dashed border-neutral-300 dark:border-neutral-700 p-5">
            <p class="text-sm text-neutral-700 dark:text-neutral-300">即将提供工作流模板库、运行记录与版本管理能力。</p>
          </div>

          <div class="flex gap-3">
            <UButton to="/projects" color="primary">前往项目工作区</UButton>
            <UButton to="/projects" color="neutral" variant="outline">查看我的项目</UButton>
          </div>
        </div>
      </UCard>
    </UContainer>
  </div>
</template>
