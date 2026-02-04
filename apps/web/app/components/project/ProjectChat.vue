<script setup lang="ts">
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
}

interface Emits {
    (e: 'update:activeConversationId', value: string | null): void
    (e: 'submit', text: string): void
    (e: 'create'): void
    (e: 'rename'): void
    (e: 'delete'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const input = ref('')

const localActiveId = computed({
    get: () => props.activeConversationId,
    set: (value) => emit('update:activeConversationId', value),
})

function onSubmit(e: Event) {
    e.preventDefault()
    const text = input.value.trim()
    if (!text) return
    
    emit('submit', text)
    input.value = ''
}
</script>

<template>
    <div class="w-96 h-full flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-800">
        <!-- 头部 -->
        <div class="flex items-center justify-between px-3 py-2 border-b border-gray-200 dark:border-gray-800 shrink-0">
            <div class="flex items-center gap-2">
                <UAvatar 
                    icon="i-lucide-bot" 
                    size="sm" 
                    alt="AI Agent" 
                    class="bg-primary-50 dark:bg-primary-900/20 text-primary-500" 
                />
                <div class="flex flex-col">
                    <span class="font-medium text-sm leading-none">AI Agent</span>
                    <span class="text-[10px] text-gray-500 dark:text-gray-400 leading-none mt-1">
                        {{ conversationTitle }}
                    </span>
                </div>
            </div>
            <div class="flex items-center gap-2">
                <USelectMenu
                    v-model="localActiveId"
                    :options="conversationOptions"
                    size="xs"
                    placeholder="选择对话"
                    class="min-w-[140px]"
                />
                <UDropdownMenu
                    :items="[[
                        { label: '重命名', icon: 'i-lucide-pencil', onSelect: () => emit('rename') },
                        { label: '删除', icon: 'i-lucide-trash-2', onSelect: () => emit('delete') }
                    ]]"
                >
                    <UButton icon="i-lucide-more-vertical" color="neutral" variant="ghost" size="xs" />
                </UDropdownMenu>
                <UTooltip text="新对话">
                    <UButton icon="i-lucide-plus" color="neutral" variant="ghost" size="xs" @click="emit('create')" />
                </UTooltip>
            </div>
        </div>

        <!-- 消息列表 -->
        <div class="flex-1 overflow-y-auto p-4">
            <UChatMessages 
                :assistant="{
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
                }" 
                :auto-scroll="{
                    color: 'neutral',
                    variant: 'outline'
                }" 
                :messages="chat.messages" 
                :status="chat.status"
            >
                <template #content="{ message }">
                    <template v-for="(part, index) in message.parts" :key="`${message.id}-${part.type}-${index}`">
                        <MDC 
                            v-if="part.type === 'text' && message.role === 'assistant'" 
                            :value="part.text"
                            :cache-key="`${message.id}-${index}`" 
                            class="*:first:mt-0 *:last:mb-0" 
                        />
                        <p v-else-if="part.type === 'text' && message.role === 'user'" class="whitespace-pre-wrap">
                            {{ part.text }}
                        </p>
                    </template>
                </template>
            </UChatMessages>
        </div>

        <!-- 输入框 -->
        <div class="p-4 pt-2">
            <UChatPrompt v-model="input" :error="chat.error" @submit="onSubmit">
                <UChatPromptSubmit :status="chat.status" @stop="chat.stop()" @reload="chat.regenerate()" />
            </UChatPrompt>
        </div>
    </div>
</template>
