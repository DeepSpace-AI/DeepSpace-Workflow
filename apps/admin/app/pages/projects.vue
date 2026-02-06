<script setup lang="ts">
useHead({ title: '项目 - DeepSpace 管理台' })

const projects = [
  { name: '星图对话', org: '星曜实验室', budget: '¥ 12,000', status: '进行中' },
  { name: '材料检索', org: '深空研究院', budget: '¥ 8,500', status: '预警' },
  { name: '天文知识库', org: '云渡团队', budget: '¥ 2,300', status: '暂停' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="项目" description="掌握项目预算、运行状态与策略绑定情况。">
      <UButton color="neutral" variant="outline" icon="i-lucide-file-plus">批量创建</UButton>
      <UButton color="primary" icon="i-lucide-folder-plus">新增项目</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="grid gap-4 md:grid-cols-2">
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">预算预警项目</p>
          <p class="mt-2 text-2xl font-semibold">6</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">策略未绑定</p>
          <p class="mt-2 text-2xl font-semibold">3</p>
        </div>
      </div>

      <div class="mt-6 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">项目</th>
              <th class="px-3 py-2">组织</th>
              <th class="px-3 py-2">预算</th>
              <th class="px-3 py-2">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="project in projects" :key="project.name" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ project.name }}</td>
              <td class="px-3 py-2 text-slate-500">{{ project.org }}</td>
              <td class="px-3 py-2">{{ project.budget }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ project.status }}</UBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
