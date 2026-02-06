<script setup lang="ts">
useHead({ title: '用量 - DeepSpace 管理台' })

const usage = [
  { model: 'gpt-4.1', tokens: '4.2M', cost: '¥ 12,300', trend: '+6%' },
  { model: 'deepseek-v3', tokens: '3.1M', cost: '¥ 6,480', trend: '+2%' },
  { model: 'claude-3.5', tokens: '1.8M', cost: '¥ 5,900', trend: '-4%' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="用量" description="跟踪模型调用量与成本走势。">
      <UButton color="neutral" variant="outline" icon="i-lucide-calendar">选择周期</UButton>
      <UButton color="primary" icon="i-lucide-download">导出报表</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="grid gap-4 md:grid-cols-3">
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">本月调用</p>
          <p class="mt-2 text-2xl font-semibold">9.1M</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">平均成本</p>
          <p class="mt-2 text-2xl font-semibold">¥ 0.23 / 1k</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">预算利用率</p>
          <p class="mt-2 text-2xl font-semibold">63%</p>
        </div>
      </div>

      <div class="mt-6 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">模型</th>
              <th class="px-3 py-2">Token</th>
              <th class="px-3 py-2">成本</th>
              <th class="px-3 py-2">趋势</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in usage" :key="row.model" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ row.model }}</td>
              <td class="px-3 py-2">{{ row.tokens }}</td>
              <td class="px-3 py-2 text-slate-500">{{ row.cost }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ row.trend }}</UBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
