<script setup lang="ts">
useHead({ title: '钱包 - DeepSpace 管理台' })

const wallets = [
  { org: '星曜实验室', balance: '¥ 42,300', frozen: '¥ 4,800', status: '正常' },
  { org: '深空研究院', balance: '¥ 12,600', frozen: '¥ 1,200', status: '预警' },
  { org: '云渡团队', balance: '¥ 1,900', frozen: '¥ 0', status: '冻结' }
]
</script>

<template>
  <div class="space-y-6">
    <AdminPageHeader title="钱包" description="查看组织余额、冻结额度与充值记录。">
      <UButton color="neutral" variant="outline" icon="i-lucide-download">导出账单</UButton>
      <UButton color="primary" icon="i-lucide-wallet">充值</UButton>
    </AdminPageHeader>

    <UCard class="border-black/5">
      <div class="grid gap-4 md:grid-cols-3">
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">平台余额</p>
          <p class="mt-2 text-2xl font-semibold">¥ 1,203,900</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">冻结额度</p>
          <p class="mt-2 text-2xl font-semibold">¥ 128,400</p>
        </div>
        <div class="rounded-xl border border-black/5 bg-slate-50 p-4">
          <p class="text-sm text-slate-500">当日充值</p>
          <p class="mt-2 text-2xl font-semibold">¥ 24,000</p>
        </div>
      </div>

      <div class="mt-6 overflow-hidden rounded-xl border border-black/5">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
            <tr>
              <th class="px-3 py-2">组织</th>
              <th class="px-3 py-2">余额</th>
              <th class="px-3 py-2">冻结</th>
              <th class="px-3 py-2">状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wallet in wallets" :key="wallet.org" class="border-t border-black/5">
              <td class="px-3 py-2 font-medium text-slate-900">{{ wallet.org }}</td>
              <td class="px-3 py-2">{{ wallet.balance }}</td>
              <td class="px-3 py-2 text-slate-500">{{ wallet.frozen }}</td>
              <td class="px-3 py-2">
                <UBadge color="neutral" variant="subtle">{{ wallet.status }}</UBadge>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>
