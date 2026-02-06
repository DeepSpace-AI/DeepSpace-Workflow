<script setup lang="ts">
const route = useRoute()

const navItems = [
  { label: '管理概览', to: '/', icon: 'i-lucide-layout-dashboard' },
  { label: '用户', to: '/users', icon: 'i-lucide-users' },
  { label: '组织', to: '/orgs', icon: 'i-lucide-building-2' },
  { label: '项目', to: '/projects', icon: 'i-lucide-folder-kanban' },
  { label: '模型', to: '/models', icon: 'i-lucide-brain' },
  { label: '定价', to: '/pricing', icon: 'i-lucide-tag' },
  { label: '审计', to: '/audit', icon: 'i-lucide-shield-check' },
  { label: '策略', to: '/policy', icon: 'i-lucide-sliders' }
]

const billingItems = [
  { label: '钱包', to: '/billing/wallets', icon: 'i-lucide-wallet' },
  { label: '交易', to: '/billing/transactions', icon: 'i-lucide-receipt' },
  { label: '用量', to: '/billing/usage', icon: 'i-lucide-bar-chart-3' }
]

const isActive = (path: string) => {
  if (path === '/') {
    return route.path === '/'
  }

  return route.path === path || route.path.startsWith(`${path}/`)
}
</script>

<template>
  <UApp>
    <div class="min-h-screen bg-[var(--admin-bg)] text-[var(--admin-ink)]">
      <header class="sticky top-0 z-30 border-b border-black/5 bg-white/80 backdrop-blur">
        <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-6 py-4">
          <div class="flex items-center gap-3">
            <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-[var(--admin-accent)] text-white">
              <AppLogo class="h-5 w-auto" />
            </div>
            <div>
              <p class="text-sm font-semibold">DeepSpace 管理台</p>
              <p class="text-xs text-slate-500">Billing · Audit · Policy</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <UBadge color="neutral" variant="subtle">Production</UBadge>
            <UButton color="neutral" variant="ghost" icon="i-lucide-bell" aria-label="通知" />
            <UAvatar text="LA" size="sm" />
          </div>
        </div>
      </header>

      <div class="mx-auto grid max-w-7xl grid-cols-1 gap-6 px-6 pb-10 pt-6 lg:grid-cols-[240px_1fr]">
        <aside class="space-y-6">
          <div class="rounded-2xl border border-black/5 bg-white/90 p-4 shadow-sm">
            <p class="text-xs font-semibold uppercase tracking-wide text-slate-500">导航</p>
            <nav class="mt-4 space-y-1">
              <NuxtLink
                v-for="item in navItems"
                :key="item.to"
                :to="item.to"
                class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition"
                :class="isActive(item.to) ? 'bg-[var(--admin-accent)]/10 text-slate-900 shadow-sm' : 'text-slate-600 hover:bg-black/5 hover:text-slate-900'"
              >
                <UIcon :name="item.icon" class="text-base" />
                <span>{{ item.label }}</span>
              </NuxtLink>
            </nav>
          </div>

          <div class="rounded-2xl border border-black/5 bg-white/90 p-4 shadow-sm">
            <p class="text-xs font-semibold uppercase tracking-wide text-slate-500">计费</p>
            <nav class="mt-4 space-y-1">
              <NuxtLink
                v-for="item in billingItems"
                :key="item.to"
                :to="item.to"
                class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition"
                :class="isActive(item.to) ? 'bg-[var(--admin-accent)]/10 text-slate-900 shadow-sm' : 'text-slate-600 hover:bg-black/5 hover:text-slate-900'"
              >
                <UIcon :name="item.icon" class="text-base" />
                <span>{{ item.label }}</span>
              </NuxtLink>
            </nav>
          </div>

          <div class="rounded-2xl border border-black/5 bg-[var(--admin-accent)] text-white p-4 shadow-sm">
            <p class="text-xs font-semibold uppercase tracking-wide text-white/70">今日状态</p>
            <div class="mt-3 space-y-2">
              <div class="flex items-center justify-between text-sm">
                <span>冻结额度</span>
                <span class="font-semibold">¥ 128,400</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span>风险拦截</span>
                <span class="font-semibold">14</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span>活跃项目</span>
                <span class="font-semibold">86</span>
              </div>
            </div>
          </div>
        </aside>

        <main class="space-y-6">
          <slot />
        </main>
      </div>
    </div>
  </UApp>
</template>
