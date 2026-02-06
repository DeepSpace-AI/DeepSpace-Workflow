<script setup lang="ts">
useHead({ title: '模型 - DeepSpace 管理台' })

const models = [
  { name: 'gpt-4.1', provider: 'OpenAI', status: '可用', latency: '850ms' },
  { name: 'deepseek-v3', provider: 'DeepSeek', status: '可用', latency: '720ms' },
  { name: 'claude-3.5', provider: 'Anthropic', status: '限流', latency: '1.1s' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="模型" description="管理上游模型、调用速率与健康状态。">
      <UButton color="neutral" variant="outline" icon="i-lucide-cloud">同步模型</UButton>
      <UButton color="primary" icon="i-lucide-plus">新增模型</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="mt-1 text-sm text-slate-500">建议为关键模型配置备用渠道与风险阈值。</div>
      <div class="mt-4 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">模型</th>
              <th class="px-3 py-2">供应商</th>
              <th class="px-3 py-2">状态</th>
              <th class="px-3 py-2">延迟</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="model in models" :key="model.name" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ model.name }}</td>
              <td class="px-3 py-2 text-slate-500">{{ model.provider }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ model.status }}</UBadge>
              </td>
              <td class="px-3 py-2">{{ model.latency }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
