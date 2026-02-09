<script setup lang="ts">
useHead({
  title: "Usage 明细 - Deepspace Workflow",
})
const requestHeaders = useRequestHeaders(['cookie'])

type UsageRecord = {
  id?: number
  model?: string
  prompt_tokens?: number
  completion_tokens?: number
  total_tokens?: number
  cost?: number
  trace_id?: string
  created_at?: string
  Model?: string
  PromptTokens?: number
  CompletionTokens?: number
  TotalTokens?: number
  Cost?: number
  TraceID?: string
  CreatedAt?: string
}

type UsageResponse = {
  items?: UsageRecord[]
  page?: number
  page_size?: number
  total?: number
}

const range = ref<'today' | '7d' | '30d'>('7d')
const page = ref(1)
const pageSize = 20

const tabs = [
  { label: '钱包概览', to: '/billing/wallet' },
  { label: 'Usage 明细', to: '/billing/usage' },
]

const rangeOptions = [
  { label: '今日', value: 'today' },
  { label: '近 7 天', value: '7d' },
  { label: '近 30 天', value: '30d' },
]

const queryParams = computed(() => {
  const now = new Date()
  let start: Date
  if (range.value === 'today') {
    start = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  } else if (range.value === '7d') {
    start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
  } else {
    start = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000)
  }
  return {
    page: page.value,
    page_size: pageSize,
    start: start.toISOString(),
    end: now.toISOString(),
  }
})

const { data, pending, error, refresh } = await useAsyncData('billing-usage', () =>
  $fetch<UsageResponse>('/api/billing/usage', {
    headers: requestHeaders,
    query: queryParams.value,
  }),
  { watch: [range, page] }
)

const items = computed(() => data.value?.items ?? [])
const total = computed(() => data.value?.total ?? 0)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

watch(range, () => {
  page.value = 1
})

function formatAmount(value?: number) {
  if (value == null) return '0.00'
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 6 })
}

function formatTime(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toISOString().slice(0, 19).replace('T', ' ')
}

function normalize(record: UsageRecord) {
  return {
    model: record.model ?? record.Model ?? 'unknown',
    totalTokens: record.total_tokens ?? record.TotalTokens ?? 0,
    cost: record.cost ?? record.Cost ?? 0,
    traceId: record.trace_id ?? record.TraceID ?? '',
    createdAt: record.created_at ?? record.CreatedAt ?? '',
  }
}
</script>

<template>
  <UContainer class="py-10 space-y-6">
    <div class="space-y-2">
      <p class="text-xs font-semibold uppercase tracking-widest text-primary">Billing</p>
      <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100">
        Usage 明细
      </h1>
      <p class="text-sm text-slate-600 dark:text-slate-300">
        查看模型调用产生的 Usage 记录。
      </p>
    </div>

    <div class="flex flex-wrap gap-2 items-center">
      <UButton
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        size="sm"
        :color="tab.to === '/billing/usage' ? 'primary' : 'neutral'"
        :variant="tab.to === '/billing/usage' ? 'solid' : 'soft'"
      >
        {{ tab.label }}
      </UButton>

      <div class="ml-auto flex items-center gap-2">
        <USelectMenu
          v-model="range"
          :options="rangeOptions"
          size="sm"
        />
        <UButton size="sm" color="neutral" variant="soft" :loading="pending" @click="refresh">
          刷新
        </UButton>
      </div>
    </div>

    <UCard class="border-slate-200/70 dark:border-slate-800">
      <div v-if="error" class="text-sm text-red-500">
        加载失败，请稍后重试。
      </div>
      <div v-else-if="items.length === 0 && !pending" class="text-sm text-slate-500">
        暂无 Usage 记录。
      </div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="text-left text-xs uppercase tracking-wider text-slate-500">
            <tr>
              <th class="py-2">时间</th>
              <th class="py-2">模型</th>
              <th class="py-2">请求数</th>
              <th class="py-2">Tokens</th>
              <th class="py-2">费用</th>
              <th class="py-2">Trace ID</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200/70 dark:divide-slate-800">
            <tr v-for="record in items" :key="record.id">
              <td class="py-3">{{ formatTime(normalize(record).createdAt) }}</td>
              <td class="py-3">{{ normalize(record).model }}</td>
              <td class="py-3">1</td>
              <td class="py-3">{{ normalize(record).totalTokens }}</td>
              <td class="py-3">{{ formatAmount(normalize(record).cost) }}</td>
              <td class="py-3 text-xs text-slate-500">
                <span v-if="normalize(record).traceId" class="font-mono">
                  {{ normalize(record).traceId }}
                </span>
                <span v-else>-</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <div class="flex items-center justify-between text-sm">
      <div class="text-slate-500">
        第 {{ page }} / {{ totalPages }} 页，共 {{ total }} 条
      </div>
      <div class="flex items-center gap-2">
        <UButton
          size="sm"
          color="neutral"
          variant="soft"
          :disabled="page <= 1 || pending"
          @click="page = Math.max(1, page - 1)"
        >
          上一页
        </UButton>
        <UButton
          size="sm"
          color="neutral"
          variant="soft"
          :disabled="page >= totalPages || pending"
          @click="page = Math.min(totalPages, page + 1)"
        >
          下一页
        </UButton>
      </div>
    </div>
  </UContainer>
</template>
