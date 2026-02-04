<script setup lang="ts">
import { Chat } from '@ai-sdk/vue';

definePageMeta({
    layout: 'editor',
});

type ConversationItem = {
    id: number
    title?: string | null
    created_at?: string
    updated_at?: string
}

type MessageItem = {
    id: number
    conversation_id: number
    role: string
    content: string
    created_at?: string
}

const route = useRoute();
const projectId = String(route.params.id);
const currentTab = ref('folder');
const content = ref('');
const input = ref('');
const activeConversationId = ref<number | null>(null);
const requestHeaders = useRequestHeaders(['cookie']);

const { data: projectData } = await useAsyncData('project', () =>
    $fetch(`/api/projects/${projectId}`, { headers: requestHeaders })
);

const { data: conversationsData, refresh: refreshConversations } = await useAsyncData('conversations', () =>
    $fetch<{ items: ConversationItem[] }>(`/api/projects/${projectId}/conversations`, { headers: requestHeaders })
);

const conversations = computed(() => conversationsData.value?.items ?? []);

const { data: messagesData, refresh: refreshMessages } = await useAsyncData('messages', () => {
    if (!activeConversationId.value) return Promise.resolve({ items: [] });
    return $fetch<{ items: MessageItem[] }>(`/api/conversations/${activeConversationId.value}/messages`, { headers: requestHeaders });
}, { watch: [activeConversationId] });

const chat = new Chat({
    api: '/api/chat',
    onError(error) {
        console.error('Chat error:', error);
    },
    onFinish: async ({ message, isError, isDisconnect, isAbort }) => {
        if (isError || isDisconnect || isAbort) return;
        await persistAssistantMessage(message);
        await refreshMessages();
    },
});

function ensureActiveConversation() {
    if (!activeConversationId.value && conversations.value.length > 0) {
        activeConversationId.value = conversations.value[0].id;
    }
}

watch(conversations, () => {
    ensureActiveConversation();
});

onMounted(() => {
    ensureActiveConversation();
});

const hydratedMessages = computed(() => {
    const items = messagesData.value?.items ?? [];
    return items.map((item) => ({
        id: String(item.id),
        role: item.role,
        parts: [{ type: 'text', text: item.content }],
    }));
});

watch(hydratedMessages, (items) => {
    // Sync server messages into Chat UI
    // @ts-ignore - ai-sdk Chat accepts direct assign in runtime
    chat.messages = items;
}, { immediate: true });

const conversationTitle = computed(() => {
    if (!activeConversationId.value) return '新对话';
    const current = conversations.value.find((item) => item.id === activeConversationId.value);
    return current?.title || '新对话';
});

async function onSubmit(e: any) {
    e.preventDefault();
    const text = input.value.trim();
    if (!text) return;

    if (!activeConversationId.value) {
        const created = await $fetch<ConversationItem>(`/api/projects/${projectId}/conversations`, {
            method: 'POST',
            body: { title: text.slice(0, 20) }
        });
        activeConversationId.value = created.id;
        await refreshConversations();
    }

    const conversationId = activeConversationId.value as number;
    await $fetch(`/api/conversations/${conversationId}/messages`, {
        method: 'POST',
        body: { role: 'user', content: text }
    });

    chat.sendMessage({ text });
    input.value = '';

    await refreshMessages();
}

async function createNewChat() {
    const created = await $fetch<ConversationItem>(`/api/projects/${projectId}/conversations`, {
        method: 'POST',
        body: { title: '新对话' }
    });
    activeConversationId.value = created.id;
    // Clear UI immediately to avoid showing previous conversation
    // @ts-ignore - ai-sdk Chat accepts direct assign in runtime
    chat.messages = [];
    await refreshConversations();
    await refreshMessages();
}

async function persistAssistantMessage(message: any) {
    if (!activeConversationId.value) return;
    if (!message || message.role !== 'assistant') return;

    const parts = Array.isArray(message.parts) ? message.parts : [];
    const text = parts
        .filter((part: any) => part?.type === 'text' && typeof part.text === 'string')
        .map((part: any) => part.text)
        .join('');

    if (!text.trim()) return;

    await $fetch(`/api/conversations/${activeConversationId.value}/messages`, {
        method: 'POST',
        body: { role: 'assistant', content: text }
    });
}
</script>
<template>
    <div class="flex flex-row w-full h-full">
        <!-- 左侧导航栏 -->
        <div class="min-h-[calc(100vh-var(--ui-header-height))] flex">
            <ProjectSider @select="currentTab = $event" />
        </div>
        <!-- 中间文本编辑区域 -->
        <div class="flex-1 flex h-[calc(100vh-var(--ui-header-height))]">
            <div class="p-2">
                <div class="my-2 text-xs pb-2">资源管理</div>
                <ProjectFolder v-if="currentTab === 'folder'" :project-id="projectId" />
                <ProjectKnowledge v-if="currentTab === 'knowledges'" :project-id="projectId" />
            </div>
            <USeparator orientation="vertical" />
            <EditorDocumentEditor v-model="content" placeholder="开始写作或输入/唤醒AI助手" :editable="true" />
            <USeparator orientation="vertical" />
        </div>
        <!-- 右侧AI对话区域 -->
        <div class="w-100 h-[calc(100vh-var(--ui-header-height))] flex flex-col bg-white dark:bg-gray-900">
            <!-- AI Agent Header -->
            <div class="flex items-center justify-between px-3 py-2 border-b border-gray-200 dark:border-gray-800 shrink-0">
                <div class="flex items-center gap-2">
                    <UAvatar icon="i-lucide-bot" size="sm" alt="AI Agent" class="bg-primary-50 dark:bg-primary-900/20 text-primary-500" />
                <div class="flex flex-col">
                    <span class="font-medium text-sm leading-none">AI Agent</span>
                    <span class="text-[10px] text-gray-500 dark:text-gray-400 leading-none mt-1">{{ conversationTitle }}</span>
                </div>
                </div>
                <div class="flex items-center gap-1">
                    <UTooltip text="新对话">
                        <UButton icon="i-lucide-plus" color="neutral" variant="ghost" size="xs" @click="createNewChat" />
                    </UTooltip>
                    <UTooltip text="历史记录">
                        <UButton icon="i-lucide-history" color="neutral" variant="ghost" size="xs" />
                    </UTooltip>
                    <UTooltip text="设置">
                        <UButton icon="i-lucide-settings" color="neutral" variant="ghost" size="xs" />
                    </UTooltip>
                </div>
            </div>

            <div class="flex-1 overflow-y-auto p-4">
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
                }" :auto-scroll="{
                    color: 'neutral',
                    variant: 'outline'
                }" :messages="chat.messages" :status="chat.status">
                    <template #content="{ message }">
                        <template v-for="(part, index) in message.parts" :key="`${message.id}-${part.type}-${index}`">
                            <MDC v-if="part.type === 'text' && message.role === 'assistant'" :value="part.text"
                                :cache-key="`${message.id}-${index}`" class="*:first:mt-0 *:last:mb-0" />
                            <p v-else-if="part.type === 'text' && message.role === 'user'" class="whitespace-pre-wrap">
                                {{ part.text }}</p>
                        </template>
                    </template>
                </UChatMessages>
            </div>
            <div class="p-4 pt-2">
                <UChatPrompt v-model="input" :error="chat.error" @submit="onSubmit">
                    <UChatPromptSubmit :status="chat.status" @stop="chat.stop()" @reload="chat.regenerate()" />
                </UChatPrompt>
            </div>
        </div>
    </div>
</template>
