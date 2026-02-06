<script setup lang="ts">
useHead({ title: '组织 - DeepSpace 管理台' })

const orgs = [
  { name: '星曜实验室', projects: 12, wallet: '¥ 78,200', status: '正常' },
  { name: '深空研究院', projects: 7, wallet: '¥ 24,900', status: '关注' },
  { name: '云渡团队', projects: 4, wallet: '¥ 6,300', status: '冻结' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="组织" description="查看组织状态、项目规模与钱包余额。">
      <UButton color="neutral" variant="outline" icon="i-lucide-file-text">导出名单</UButton>
      <UButton color="primary" icon="i-lucide-building">新增组织</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="grid gap-4 md:grid-cols-3">
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">活跃组织</p>
          <p class="mt-2 text-2xl font-semibold">26</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">本月新增</p>
          <p class="mt-2 text-2xl font-semibold">5</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">平均项目数</p>
          <p class="mt-2 text-2xl font-semibold">8.4</p>
        </div>
      </div>

      <div class="mt-6 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">组织</th>
              <th class="px-3 py-2">项目</th>
              <th class="px-3 py-2">钱包</th>
              <th class="px-3 py-2">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="org in orgs" :key="org.name" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ org.name }}</td>
              <td class="px-3 py-2">{{ org.projects }}</td>
              <td class="px-3 py-2 text-slate-500">{{ org.wallet }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ org.status }}</UBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
