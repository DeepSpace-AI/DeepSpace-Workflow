<script setup lang="ts">
useHead({
  title: "钱包概览 - Deepspace Workflow",
})
const requestHeaders = useRequestHeaders(['cookie'])

type WalletRecord = {
  balance?: number
  frozen_balance?: number
  updated_at?: string
  Balance?: number
  FrozenBalance?: number
  UpdatedAt?: string
}

type WalletResponse = {
  wallet?: WalletRecord
  usage_24h?: number
  usage24h?: number
}

const { data, pending, error, refresh } = await useAsyncData('billing-wallet', () =>
  $fetch<WalletResponse>('/api/billing/wallet', { headers: requestHeaders })
)

const wallet = computed(() => {
  const raw = data.value?.wallet
  return {
    balance: raw?.balance ?? raw?.Balance ?? 0,
    frozen: raw?.frozen_balance ?? raw?.FrozenBalance ?? 0,
    updatedAt: raw?.updated_at ?? raw?.UpdatedAt ?? '',
  }
})

const usage24h = computed(() => data.value?.usage_24h ?? data.value?.usage24h ?? 0)

const tabs = [
  { label: '钱包概览', to: '/billing/wallet' },
  { label: 'Usage 明细', to: '/billing/usage' },
]

function formatAmount(value: number) {
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
}

function formatTime(value?: string) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toISOString().slice(0, 19).replace('T', ' ')
}
</script>

<template>
  <UContainer class="py-10 space-y-6">
    <div class="space-y-2">
      <p class="text-xs font-semibold uppercase tracking-widest text-primary">Billing</p>
      <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100">
        钱包概览
      </h1>
      <p class="text-sm text-slate-600 dark:text-slate-300">
        查看当前余额、冻结金额与近 24 小时 usage。
      </p>
    </div>

    <div class="flex flex-wrap gap-2">
      <UButton
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        size="sm"
        :color="tab.to === '/billing/wallet' ? 'primary' : 'neutral'"
        :variant="tab.to === '/billing/wallet' ? 'solid' : 'soft'"
      >
        {{ tab.label }}
      </UButton>
      <UButton size="sm" color="neutral" variant="soft" :loading="pending" @click="refresh">
        刷新
      </UButton>
    </div>

    <UCard class="border-slate-200/70 dark:border-slate-800">
      <div v-if="error" class="text-sm text-red-500">
        加载失败，请稍后重试。
      </div>
      <div v-else class="grid gap-6 md:grid-cols-3">
        <div>
          <p class="text-xs text-slate-500 dark:text-slate-400">可用余额</p>
          <p class="mt-2 text-2xl font-semibold text-slate-900 dark:text-white">
            {{ formatAmount(wallet.balance) }}
          </p>
        </div>
        <div>
          <p class="text-xs text-slate-500 dark:text-slate-400">冻结金额</p>
          <p class="mt-2 text-2xl font-semibold text-slate-900 dark:text-white">
            {{ formatAmount(wallet.frozen) }}
          </p>
        </div>
        <div>
          <p class="text-xs text-slate-500 dark:text-slate-400">近 24 小时 Usage</p>
          <p class="mt-2 text-2xl font-semibold text-slate-900 dark:text-white">
            {{ formatAmount(usage24h) }}
          </p>
        </div>
      </div>
      <p class="mt-4 text-xs text-slate-500 dark:text-slate-400">
        更新时间：{{ formatTime(wallet.updatedAt) }}
      </p>
    </UCard>
  </UContainer>
</template>
