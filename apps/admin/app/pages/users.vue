<script setup lang="ts">
useHead({ title: '用户 - DeepSpace 管理台' })

const users = [
  { name: '林澜', role: '组织管理员', org: '星曜实验室', status: '活跃' },
  { name: '杨策', role: '研究员', org: '深空研究院', status: '活跃' },
  { name: '乔伊', role: '访客', org: '云渡团队', status: '冻结' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="用户" description="管理组织成员、角色与访问状态。">
      <UButton color="neutral" variant="outline" icon="i-lucide-upload">导入</UButton>
      <UButton color="primary" icon="i-lucide-user-plus">新增用户</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <UInput placeholder="搜索用户或组织" icon="i-lucide-search" />
        <div class="flex items-center gap-2">
          <UButton color="neutral" variant="ghost" icon="i-lucide-filter">筛选</UButton>
          <UButton color="neutral" variant="ghost" icon="i-lucide-sliders">排序</UButton>
        </div>
      </div>

      <div class="mt-4 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">用户</th>
              <th class="px-3 py-2">角色</th>
              <th class="px-3 py-2">组织</th>
              <th class="px-3 py-2">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.name" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ user.name }}</td>
              <td class="px-3 py-2">{{ user.role }}</td>
              <td class="px-3 py-2 text-slate-500">{{ user.org }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ user.status }}</UBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
