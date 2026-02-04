<script setup lang="ts">
const emit = defineEmits<{
    (e: 'select', value: string): void;
}>();

const tabs = [
    {
        label: '文件管理',
        icon: 'i-lucide-folder',
        value: 'folder'
    },
    {
        label: '工作流库',
        icon: 'i-lucide-activity',
        value: 'workflows'
    },
    {
        label: '知识库',
        icon: 'i-lucide-database',
        value: 'knowledges'
    },
    {
        label: '成员管理',
        icon: 'i-lucide-users',
        value: 'members'
    },
    {
        label: '项目设置',
        icon: 'i-lucide-settings',
        value: 'settings'
    }
];

const selectedTab = ref('folder');

const handleTabClick = (value: string) => {
    selectedTab.value = value;
    emit('select', value);
};
</script>

<template>
    <div class="flex items-start bg-default min-h-[calc(100vh-var(--ui-header-height))]">
        <nav class="flex flex-col gap-1 p-2">
            <UTooltip v-for="tab in tabs" :key="tab.value" :text="tab.label" :content="{
                align: 'center',
                side: 'right',
            }">
                <button @click="handleTabClick(tab.value)" :class="[
                    'relative flex items-center justify-center w-12 h-12 rounded-lg transition-all duration-200',
                    'hover:bg-gray-100 dark:hover:bg-gray-800',
                    'before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:h-6 before:w-1 before:rounded-r before:transition-all before:duration-200',
                    selectedTab === tab.value
                        ? 'bg-primary-50 dark:bg-primary-950 text-primary-600 dark:text-primary-400 before:bg-primary-600 dark:before:bg-primary-400'
                        : 'text-gray-600 dark:text-gray-400 before:bg-transparent'
                ]">
                    <UIcon :name="tab.icon" class="w-6 h-6" />
                </button>
            </UTooltip>
        </nav>

        <USeparator orientation="vertical" />
    </div>
</template>