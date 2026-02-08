<template>
    <UDashboardGroup>
        <UDashboardSidebar collapsible resizable
            class="bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800">
            <template #header="{ collapsed }">
                <div class="flex items-center gap-2 px-2" :class="{ 'justify-center': collapsed }">
                    <UIcon name="i-heroicons-rocket-launch" class="w-8 h-8 text-primary-500" />
                    <span v-if="!collapsed" class="text-xl font-bold truncate">DeepSpace</span>
                </div>
            </template>

            <UNavigationMenu :items="links" orientation="vertical" />

            <div class="flex-1" />

            <template #footer="{ collapsed }">
                <UDropdownMenu :items="userItems" :content="{ side: 'top', align: 'center' }" :ui="{
                    content: 'min-w-fit'
                }">
                    <UButton color="neutral" variant="ghost" class="w-full"
                        :class="{ 'justify-center': collapsed, 'justify-start': !collapsed }">
                        <UAvatar :src="avatarSrc || undefined" :alt="displayName" size="xs" />

                        <div v-if="!collapsed" class="text-left truncate">
                            <div class="font-medium text-gray-900 dark:text-white truncate">{{ displayName }}</div>
                            <div class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ displayEmail }}</div>
                        </div>

                        <UIcon v-if="!collapsed" name="i-heroicons-ellipsis-vertical" class="ml-auto text-gray-400" />
                    </UButton>
                </UDropdownMenu>
            </template>
        </UDashboardSidebar>

        <UDashboardPanel class="bg-gray-50 dark:bg-gray-950">

            <UDashboardNavbar :title="routeName"
                class="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800">
                <template #right>
                    <UColorModeButton />
                    <UButton icon="i-heroicons-bell" color="neutral" variant="ghost" />
                </template>
            </UDashboardNavbar>

            <div class="p-4 overflow-y-auto h-full">
                <slot />
            </div>
        </UDashboardPanel>
    </UDashboardGroup>
</template>

<script setup lang="ts">
const route = useRoute()
const sessionCookie = useCookie('dsp_session')

type UserMeResponse = {
    user: {
        id: number
        email: string
        status: string
        role: string
        last_login_at?: string | null
        created_at?: string | null
    }
    profile?: {
        display_name?: string | null
        full_name?: string | null
        title?: string | null
        avatar_url?: string | null
        bio?: string | null
        phone?: string | null
    }
    settings?: {
        theme?: string | null
        locale?: string | null
        timezone?: string | null
    }
}

const { data: meData } = await useFetch<UserMeResponse | null>('/api/users/me', {
    default: () => null
})

const displayName = computed(() => {
    const profile = meData.value?.profile
    const user = meData.value?.user
    return profile?.display_name || profile?.full_name || user?.email || '管理员'
})

const displayEmail = computed(() => meData.value?.user?.email || '—')
const avatarSrc = computed(() => meData.value?.profile?.avatar_url || '')

const handleLogout = async () => {
    await $fetch('/api/auth/logout', { method: 'POST' }).catch(() => null)
    sessionCookie.value = null
    await navigateTo('/sign-in')
}

const userItems = computed(() => [
    [
        {
            label: '个人设置',
            icon: 'i-heroicons-cog-6-tooth',
            to: '/profile'
        }
    ],
    [{ label: '退出登录', icon: 'i-heroicons-arrow-left-on-rectangle', onSelect: handleLogout }]
])

const routeName = computed(() => {
    const name = String(route.name || '').toLowerCase()
    const map: Record<string, string> = {
        'index': '仪表盘',
        'users': '用户管理',
        'models': '模型管理',
        'pricing': '套餐管理',
        'billing-wallets': '钱包管理',
        'billing-transactions': '交易流水',
        'billing-usage': '用量记录',
        'audit': '审计日志',
        'policy': '风控策略'
    }
    return map[name] || 'DeepSpace'
})

const links = [
    { label: '仪表盘', icon: 'i-heroicons-home', to: '/' },
    {
        label: '用户管理',
        icon: 'i-heroicons-users',
        children: [
            { label: '用户管理', to: '/users' }
        ]
    },
    { label: '模型管理', icon: 'i-heroicons-cpu-chip', to: '/models' },
    { label: '套餐管理', icon: 'i-heroicons-currency-dollar', to: '/pricing' },
    {
        label: '财务管理', icon: 'i-heroicons-credit-card', children: [
            { label: '钱包管理', to: '/billing/wallets' },
            { label: '交易流水', to: '/billing/transactions' },
            { label: '用量记录', to: '/billing/usage' }
        ]
    },
    { label: '审计日志', icon: 'i-heroicons-clipboard-document-list', to: '/audit' },
    { label: '风控策略', icon: 'i-heroicons-shield-check', to: '/policy' }
]
</script>
