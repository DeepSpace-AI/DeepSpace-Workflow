<script setup lang="ts">
useHead({ title: '交易 - DeepSpace 管理台' })

const transactions = [
  { id: 'TX-4832', org: '星曜实验室', type: '冻结', amount: '¥ 1,200', time: '09:18' },
  { id: 'TX-4821', org: '深空研究院', type: '扣费', amount: '¥ 680', time: '08:44' },
  { id: 'TX-4819', org: '云渡团队', type: '充值', amount: '¥ 5,000', time: '08:02' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="交易" description="追踪充值、冻结、扣费与解冻记录。">
      <UButton color="neutral" variant="outline" icon="i-lucide-filter">筛选</UButton>
      <UButton color="primary" icon="i-lucide-download">导出</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <UInput placeholder="搜索交易号或组织" icon="i-lucide-search" />
        <UButton color="neutral" variant="ghost" icon="i-lucide-refresh-ccw">刷新</UButton>
      </div>

      <div class="mt-4 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">交易号</th>
              <th class="px-3 py-2">组织</th>
              <th class="px-3 py-2">类型</th>
              <th class="px-3 py-2">金额</th>
              <th class="px-3 py-2">时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tx in transactions" :key="tx.id" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ tx.id }}</td>
              <td class="px-3 py-2 text-slate-500">{{ tx.org }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ tx.type }}</UBadge>
              </td>
              <td class="px-3 py-2">{{ tx.amount }}</td>
              <td class="px-3 py-2 text-slate-500">{{ tx.time }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
