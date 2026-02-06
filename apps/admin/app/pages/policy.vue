<script setup lang="ts">
useHead({
  title: '策略 · DeepSpace Admin'
})
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader
      title="风控策略"
      description="定义预算冻结、阈值告警与访问控制规则。"
    >
      <template #actions>
        <UButton size="sm" color="primary" icon="i-lucide-plus">新增策略</UButton>
      </template>
    </AdminPageHeader>

    <UCard class="bg-slate-950/40">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm font-semibold">策略列表</p>
        <div class="flex items-center gap-2">
          <UButton size="sm" variant="ghost" color="neutral">导入模板</UButton>
          <UButton size="sm" variant="ghost" color="neutral">历史版本</UButton>
        </div>
      </div>
      <div class="mt-4 space-y-2 text-sm text-white/60">
        <div class="rounded-xl border border-white/10 bg-white/5 px-3 py-3">
          尚未创建策略。你可以添加预算冻结、API 限制或异常检测规则。
        </div>
      </div>
    </UCard>
  </div>
</template>
