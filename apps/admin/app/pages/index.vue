<script setup lang="ts">
useHead({
  title: '管理概览 - DeepSpace 管理台'
})

const stats = [
  { label: '活跃组织', value: '42', trend: '+8%' },
  { label: '今日调用', value: '18,640', trend: '+12%' },
  { label: '冻结余额', value: '¥ 312,890', trend: '-4%' },
  { label: '风险拦截', value: '14', trend: '+3' }
]

const alerts = [
  { title: '异常调用峰值', detail: '模型 gpt-4.1 在 09:10 出现 3 倍峰值', level: '高' },
  { title: '企业钱包余额不足', detail: '组织「星曜实验室」余额低于阈值', level: '中' },
  { title: '策略变更待审核', detail: '新增跨区访问策略', level: '低' }
]

const audits = [
  { time: '09:12', actor: '风控机器人', action: '拦截异常调用', scope: 'Project #2094', result: '已阻断' },
  { time: '08:57', actor: '财务管理员', action: '调整定价', scope: '模型 gpt-4.1', result: '已发布' },
  { time: '08:30', actor: '系统', action: '钱包冻结', scope: 'Org #1202', result: '冻结 ¥4,800' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader
      title="管理概览"
      description="监控计费、审计与风控状态，确保组织与项目调用安全可控。"
    >
      <UButton color="neutral" variant="outline" icon="i-lucide-download">导出日报</UButton>
      <UButton color="primary" icon="i-lucide-refresh-ccw">刷新</UButton>
    </AdminPageHeader>

    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <UCard v-for="stat in stats" :key="stat.label" class="border-black/5">
        <div class="space-y-2">
          <p class="text-sm text-slate-500">{{ stat.label }}</p>
          <div class="flex items-end justify-between">
            <p class="text-2xl font-semibold text-slate-900">{{ stat.value }}</p>
            <UBadge color="neutral" variant="subtle">{{ stat.trend }}</UBadge>
          </div>
        </div>
      </UCard>
    </section>

    <section class="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
      <UCard class="border-black/5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-lg font-semibold text-slate-900">最新审计</p>
            <p class="text-sm text-slate-500">关键操作与系统动作的实时记录。</p>
          </div>
          <UButton color="neutral" variant="ghost" icon="i-lucide-arrow-right">查看全部</UButton>
        </div>
        <div class="mt-4 overflow-hidden rounded-xl border border-black/5">
          <table class="w-full text-left text-sm">
            <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
              <tr>
                <th class="px-3 py-2">时间</th>
                <th class="px-3 py-2">操作人</th>
                <th class="px-3 py-2">动作</th>
                <th class="px-3 py-2">范围</th>
                <th class="px-3 py-2">结果</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="audit in audits" :key="audit.time" class="border-t border-black/5">
                <td class="px-3 py-2 text-slate-500">{{ audit.time }}</td>
                <td class="px-3 py-2 font-medium text-slate-900">{{ audit.actor }}</td>
                <td class="px-3 py-2">{{ audit.action }}</td>
                <td class="px-3 py-2 text-slate-500">{{ audit.scope }}</td>
                <td class="px-3 py-2">
                  <UBadge color="neutral" variant="subtle">{{ audit.result }}</UBadge>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </UCard>

      <div class="space-y-6">
        <UCard class="border-black/5">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-lg font-semibold text-slate-900">待处理</p>
              <p class="text-sm text-slate-500">优先处理风险与余额预警。</p>
            </div>
            <UButton color="neutral" variant="ghost" icon="i-lucide-filter">筛选</UButton>
          </div>
          <div class="mt-4 space-y-3">
            <div
              v-for="alert in alerts"
              :key="alert.title"
              class="rounded-xl border border-black/5 bg-slate-50 px-3 py-3"
            >
              <div class="flex items-center justify-between">
                <p class="text-sm font-semibold text-slate-900">{{ alert.title }}</p>
                <UBadge color="neutral" variant="subtle">{{ alert.level }}</UBadge>
              </div>
              <p class="mt-1 text-xs text-slate-500">{{ alert.detail }}</p>
            </div>
          </div>
        </UCard>

        <UCard class="border-black/5">
          <p class="text-lg font-semibold text-slate-900">模型健康</p>
          <p class="text-sm text-slate-500">最近 24 小时模型调用与稳定性。</p>
          <div class="mt-4 space-y-3">
            <div class="flex items-center justify-between rounded-xl border border-black/5 px-3 py-2">
              <span class="text-sm font-medium">gpt-4.1</span>
              <span class="text-xs text-slate-500">可用率 99.8%</span>
            </div>
            <div class="flex items-center justify-between rounded-xl border border-black/5 px-3 py-2">
              <span class="text-sm font-medium">deepseek-v3</span>
              <span class="text-xs text-slate-500">可用率 99.5%</span>
            </div>
            <div class="flex items-center justify-between rounded-xl border border-black/5 px-3 py-2">
              <span class="text-sm font-medium">claude-3.5</span>
              <span class="text-xs text-slate-500">可用率 99.7%</span>
            </div>
          </div>
        </UCard>
      </div>
    </section>
  </div>
</template>
