<script setup lang="ts">
useHead({
  title: '审计 · DeepSpace Admin'
})
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader
      title="审计日志"
      description="追踪 trace_id 与调用链路，支持合规检查。"
    >
      <template #actions>
        <UButton size="sm" color="primary" icon="i-lucide-search">查询</UButton>
      </template>
    </AdminPageHeader>

    <UCard class="bg-slate-950/40">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm font-semibold">审计记录</p>
        <div class="flex items-center gap-2">
          <UInput size="sm" placeholder="搜索 trace_id / 用户 / 项目" icon="i-lucide-search" />
          <UButton size="sm" variant="ghost" color="neutral">筛选</UButton>
        </div>
      </div>
      <div class="mt-4 grid grid-cols-4 gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-2 text-xs uppercase text-white/40">
        <span>trace_id</span>
        <span>用户</span>
        <span>项目</span>
        <span>时间</span>
      </div>
      <div class="mt-2 space-y-2 text-sm text-white/60">
        <div class="grid grid-cols-4 gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-3">
          <span>--</span>
          <span>--</span>
          <span>--</span>
          <span>--</span>
        </div>
        <div class="grid grid-cols-4 gap-2 rounded-xl border border-white/10 bg-white/5 px-3 py-3">
          <span>--</span>
          <span>--</span>
          <span>--</span>
          <span>--</span>
        </div>
      </div>
    </UCard>
  </div>
</template>
