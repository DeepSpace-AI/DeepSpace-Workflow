<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';
import { zh_cn as uiLocale } from '@nuxt/ui/locale';
const route = useRoute();
const { data: currentUser } = await useAsyncData('layout-current-user-editor', () =>
    $fetch<Record<string, any> | null>('/api/users/me').catch(() => null),
);

const items = computed<NavigationMenuItem[]>(() => [
    {
        label: '对话',
        to: '/chat',
        active: route.path.startsWith('/chat'),
    },
    {
        label: '项目',
        to: '/projects',
        active: route.path.startsWith('/projects'),
    },
    {
        label: '工作流',
        to: '/workflows',
        active: route.path.startsWith('/workflows'),
    },
    {
        label: '知识库',
        to: '/knowledge',
        active: route.path.startsWith('/knowledge'),
    },
    {
        label: '计费',
        to: '/billing/wallet',
        active: route.path.startsWith('/billing'),
    },
]);

const dropdownItems = computed<DropdownMenuItem[]>(() => [
    {
        label: '个人资料',
        icon: 'i-lucide-user',
        to: '/profile',
    },
    {
        label: '设置',
        icon: 'i-lucide-settings',
        to: '/settings',
    },
    {
        label: '登出',
        icon: 'i-lucide-log-out',
        to: '/sign-out'
    },
]);

const userName = computed(() => {
    const value = currentUser.value;
    const profile = value?.profile || {};
    return profile?.DisplayName || profile?.FullName || value?.nickname || value?.name || value?.username || '已登录用户';
});

const userDescription = computed(() => {
    const value = currentUser.value;
    const profile = value?.profile || {};
    return value?.email || profile?.title || profile?.phone || value?.phone || '个人信息';
});

const userAvatar = computed(() => {
    const value = currentUser.value;
    const profile = value?.profile || {};
    return {
        src: profile?.avatar_url || value?.avatar || value?.avatar_url || undefined,
        icon: 'i-lucide-user',
    };
});

</script>
<template>
    <UApp :locale="uiLocale">
        <UHeader>
            <template #title>
                <NuxtLink to="/">DeepSpace Workflows</NuxtLink>
            </template>

            <UNavigationMenu :items="items" />

            <template #right>
                <div class="flex justify-around items-center gap-4">
                    <UColorModeSelect />

                    <UDropdownMenu :items="dropdownItems">
                        <UUser :name="userName" :description="userDescription" :avatar="userAvatar" />
                    </UDropdownMenu>
                </div>
            </template>
        </UHeader>

        <UMain>
            <slot />
        </UMain>
    </UApp>
</template>
