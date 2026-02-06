<script setup lang="ts">
import type { ChatContextItem } from '#imports'

type ConversationOption = {
    label: string
    value: string
}

interface Props {
    projectId: string
    chat: any
    conversationOptions: ConversationOption[]
    activeConversationId: string | null
    conversationTitle: string
    contextItems: ChatContextItem[]
    selectedSkills: string[]
    selectedWorkflows: string[]
    skillOptions: ConversationOption[]
    workflowOptions: ConversationOption[]
    hasSelection: boolean
}

interface Emits {
    (e: 'update:activeConversationId', value: string | null): void
    (e: 'submit', text: string): void
    (e: 'create'): void
    (e: 'delete', conversationId: string): void
    (e: 'settings'): void
    (e: 'request-selection-context'): void
    (e: 'upload-file', file: File): void
    (e: 'update-context-items', items: ChatContextItem[]): void
    (e: 'update-selected-skills', items: string[]): void
    (e: 'update-selected-workflows', items: string[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const input = ref('')
const searchQuery = ref('')
const showHistory = ref(false)
const showAgentDialog = ref(false)
const agentName = ref('')
const isLoading = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)

const localActiveId = computed({
    get: () => props.activeConversationId,
    set: (value) => emit('update:activeConversationId', value),
})

const skillModel = computed({
    get: () => props.selectedSkills,
    set: (value: string[]) => emit('update-selected-skills', value),
})

const workflowModel = computed({
    get: () => props.selectedWorkflows,
    set: (value: string[]) => emit('update-selected-workflows', value),
})

// 过滤对话列表
const filteredConversations = computed(() => {
    if (!searchQuery.value.trim()) return props.conversationOptions

    const query = searchQuery.value.toLowerCase()
    return props.conversationOptions.filter(option =>
        option.label.toLowerCase().includes(query)
    )
})

function onSubmit(e: Event) {
    e.preventDefault()
    const text = input.value.trim()
    if (!text) return

    emit('submit', text)
    input.value = ''
}

async function handleCreateNew() {
    showHistory.value = false
    input.value = ''
    await emit('create')
}

function openHistory() {
    showHistory.value = true
}

function handleSelectHistory(value: string) {
    localActiveId.value = value
    showHistory.value = false
}

function handleRemoveContext(itemId: string) {
    emit('update-context-items', props.contextItems.filter((item) => item.id !== itemId))
}

function triggerFileSelect() {
    fileInputRef.value?.click()
}

function handleFileChange(event: Event) {
    const target = event.target as HTMLInputElement | null
    const file = target?.files?.[0]
    if (file) {
        emit('upload-file', file)
    }
    if (target) {
        target.value = ''
    }
}

function openAgentDialog() {
    showAgentDialog.value = true
}

function closeAgentDialog() {
    showAgentDialog.value = false
    agentName.value = ''
}

function submitAgentDialog() {
    // 先占位：后续接入真正的 Agent 创建流程
    closeAgentDialog()
    emit('settings')
}
</script>

<template>
    <div class="w-96 h-full flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-800">
        <!-- 头部 -->
        <div class="flex items-center justify-between px-3 py-2 border-b border-gray-200 dark:border-gray-800 shrink-0">
            <div class="flex items-center gap-2">
                <UAvatar icon="i-lucide-bot" size="sm" alt="AI Agent"
                    class="bg-primary-50 dark:bg-primary-900/20 text-primary-500" />
                <div class="flex flex-col">
                    <span class="font-medium text-sm leading-none">AI Agent</span>
                    <span class="text-[10px] text-gray-500 dark:text-gray-400 leading-none mt-1">
                        {{ conversationTitle }}
                    </span>
                </div>
            </div>
            <div class="flex items-center gap-2">
                <UTooltip text="新建对话">
                    <UButton icon="i-lucide-plus" color="primary" variant="soft" size="xs" @click="handleCreateNew"
                        :disabled="isLoading" />
                </UTooltip>
                <UTooltip text="历史对话">
                    <UButton icon="i-lucide-history" color="neutral" variant="ghost" size="xs" @click="openHistory"
                        :disabled="isLoading" />
                </UTooltip>
                <UDropdownMenu :items="[[{ label: '新建 Agent', icon: 'i-lucide-bot', onSelect: openAgentDialog }]]"
                    :disabled="isLoading">
                    <template #default="{ open }">
                        <UTooltip text="更多">
                            <UButton icon="i-lucide-more-horizontal" color="neutral" variant="ghost" size="xs"
                                :active="open" :disabled="isLoading" />
                        </UTooltip>
                    </template>
                </UDropdownMenu>
            </div>
        </div>

        <!-- 搜索对话框 -->
        <div v-if="showHistory" class="px-3 py-2 border-b border-gray-200 dark:border-gray-800 shrink-0">
            <UInput v-model="searchQuery" placeholder="搜索对话..." icon="i-lucide-search" size="sm"
                :disabled="isLoading" />
        </div>

        <!-- 历史对话列表 -->
        <div v-if="showHistory" class="flex-1 overflow-y-auto p-3">
            <div v-if="filteredConversations.length === 0" class="text-center py-6 text-sm text-gray-500">
                {{ searchQuery ? '没有找到对话' : '暂无对话' }}
            </div>
            <div v-else class="space-y-2">
                <button v-for="option in filteredConversations" :key="option.value"
                    class="w-full rounded-lg border border-transparent px-3 py-2 text-left transition hover:border-primary/30 hover:bg-primary/5"
                    :class="{ 'border-primary/40 bg-primary/10': localActiveId === option.value }"
                    @click="handleSelectHistory(option.value)" :disabled="isLoading">
                    <div class="flex items-center justify-between gap-2">
                        <p class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ option.label }}</p>
                        <UButton icon="i-lucide-trash-2" color="neutral" variant="ghost" size="xs"
                            @click.stop="emit('delete', option.value)" />
                    </div>
                </button>
            </div>
        </div>

        <!-- 消息列表 -->
        <div v-else class="flex-1 overflow-y-auto p-4">
            <UChatMessages :assistant="{
                side: 'left',
                variant: 'outline',
                avatar: {
                    icon: 'i-lucide-bot'
                },
                actions: [
                    {
                        label: 'Copy to clipboard',
                        icon: 'i-lucide-copy',
                    }
                ]
            }" class="text-sm" :messages="chat.messages" :status="chat.status">
                <template #content="{ message }">
                    <template v-for="(part, index) in message.parts" :key="`${message.id}-${part.type}-${index}`">
                        <MDC v-if="part.type === 'text' && message.role === 'assistant'" :value="part.text"
                            :cache-key="`${message.id}-${index}`" class="*:first:mt-0 *:last:mb-0" />
                        <div v-else-if="part.type === 'text' && message.role === 'user'" class="space-y-2">
                            <!-- <div v-if="message.context?.length" class="flex flex-wrap gap-2">
                                <span v-for="item in message.context" :key="item.id"
                                    class="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-1 text-[11px] text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                                    <span class="truncate max-w-[180px]">{{ item.label }}</span>
                                </span>
                            </div> -->
                            <p class="whitespace-pre-wrap">{{ part.text }}</p>
                        </div>
                    </template>
                </template>
            </UChatMessages>
        </div>

        <!-- 输入框 -->
        <div v-if="!showHistory"
            class="p-1 rounded-2xl border-slate-200/80 bg-white shadow-sm dark:border-slate-800 dark:bg-slate-900">
            <UChatPrompt v-model="input" placeholder="使用 @ 唤起更多功能" @submit="onSubmit" :error="chat?.error"
                :disabled="isLoading" size="sm" :ui="{ base: 'text-sm', body: 'text-sm' }" class="bg-transparent">
                <template #header>
                    <div class="space-y-2">
                        <div v-if="contextItems.length" class="flex flex-wrap gap-2">
                            <span v-for="item in contextItems" :key="item.id"
                                class="inline-flex items-center gap-1 rounded-full bg-slate-100 px-2 py-1 text-[11px] text-slate-700 dark:bg-slate-800 dark:text-slate-200">
                                <span class="truncate max-w-[180px]">{{ item.label }}</span>
                                <button type="button" class="text-slate-400 hover:text-slate-600"
                                    @click="handleRemoveContext(item.id)">
                                    ×
                                </button>
                            </span>
                        </div>

                    </div>
                </template>

                <template #trailing>
                    <div class="flex items-center gap-2">
                        <UChatPromptSubmit square size="sm" :status="chat.status" @stop="chat.stop()"
                            @reload="chat.regenerate()" />
                    </div>
                </template>

                <template #footer>
                    <div class="flex items-center justify-between text-xs text-slate-500">
                        <div class="flex items-center gap-1">
                            <UTooltip text="上传附件">
                                <UButton icon="i-lucide-paperclip" size="xs" color="neutral" variant="ghost"
                                    @click="triggerFileSelect" :disabled="isLoading" />
                            </UTooltip>
                            <UTooltip text="引用选中内容">
                                <UButton icon="i-lucide-quote" size="xs" color="neutral" variant="ghost"
                                    @click="emit('request-selection-context')" :disabled="isLoading || !hasSelection" />
                            </UTooltip>
                            <input ref="fileInputRef" type="file" accept="image/*,.pdf,.doc,.docx,.txt,.md"
                                class="hidden" @change="handleFileChange" />
                        </div>
                    </div>
                </template>
            </UChatPrompt>
        </div>

        <!-- 新建 Agent -->
        <UModal v-model:open="showAgentDialog">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">新建 Agent</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs"
                                @click="closeAgentDialog" />
                        </div>
                    </template>
                    <div class="space-y-3">
                        <UInput v-model="agentName" placeholder="输入 Agent 名称" />
                        <p class="text-xs text-slate-500 dark:text-slate-400">创建后可在设置中进一步配置。</p>
                    </div>
                    <template #footer>
                        <div class="flex gap-2 justify-end">
                            <UButton label="取消" color="neutral" variant="ghost" @click="closeAgentDialog" />
                            <UButton label="创建" color="primary" @click="submitAgentDialog" />
                        </div>
                    </template>
                </UCard>
            </template>
        </UModal>

    </div>
</template>
