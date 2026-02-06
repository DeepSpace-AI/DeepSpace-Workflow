<script setup lang="ts">
useHead({ title: '定价 - DeepSpace 管理台' })

const plans = [
  { name: '研究版', tokens: '500k', price: '¥ 1,280/月', status: '启用' },
  { name: '企业版', tokens: '2M', price: '¥ 4,800/月', status: '启用' },
  { name: '自定义', tokens: '不限', price: '议价', status: '草稿' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="定价" description="统一维护套餐、计费规则与折扣策略。">
      <UButton color="neutral" variant="outline" icon="i-lucide-file-text">导出配置</UButton>
      <UButton color="primary" icon="i-lucide-tag">新增套餐</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="grid gap-4 md:grid-cols-3">
        <div v-for="plan in plans" :key="plan.name" class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <div class="flex items-center justify-between">
            <p class="text-base font-semibold text-slate-900">{{ plan.name }}</p>
            <UBadge color="neutral" variant="subtle">{{ plan.status }}</UBadge>
          </div>
          <p class="mt-2 text-sm text-slate-500">额度：{{ plan.tokens }}</p>
          <p class="mt-1 text-lg font-semibold">{{ plan.price }}</p>
        </div>
      </div>
    </UCard>
  </div>
</template>
