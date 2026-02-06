<script setup lang="ts">
import type { NavigationMenuItem, DropdownMenuItem } from '@nuxt/ui';
import { zh_cn } from '@nuxt/ui/locale'

const route = useRoute();

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

</script>
<template>
    <UApp :locale="zh_cn">
        <UHeader>
            <template #title>
                <NuxtLink to="/">DeepSpace Workflows</NuxtLink>
            </template>

            <UNavigationMenu :items="items" />

            <template #right>
                <div class="flex justify-around items-center gap-4">
                    <UColorModeSelect />

                    <UDropdownMenu :items="dropdownItems">
                        <UUser name="Jhon" description="Software Engineer" :avatar="{
                            src: 'https://i.pravatar.cc/150?u=john-doe',
                            icon: 'i-lucide-image'
                        }" />
                    </UDropdownMenu>
                </div>
            </template>
        </UHeader>

        <UMain>
            <slot />
        </UMain>

        <UFooter>
            <p class="text-muted text-sm">Copyright © DeepSpace {{ new Date().getFullYear() }}</p>
        </UFooter>
    </UApp>
</template>
