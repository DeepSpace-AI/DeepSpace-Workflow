import { Chat } from '@ai-sdk/vue'
import { useProjectWorkspaceStore } from '#imports'

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

export function useProjectConversations(projectId: string) {
    const store = useProjectWorkspaceStore()
    const requestHeaders = useRequestHeaders(['cookie'])

    // 加载对话列表
    const { data: conversationsData, refresh: refreshConversations } = useAsyncData('conversations', () =>
        $fetch<{ items: ConversationItem[] }>(`/api/projects/${projectId}/conversations`, { headers: requestHeaders })
    )

    // 监听对话数据变化
    watch(
        () => conversationsData.value?.items,
        (items) => {
            if (!items) return
            store.setConversations(items)
            if (!store.activeConversationId && items.length > 0) {
                store.setActiveConversation(items[0].id)
            }
        },
        { immediate: true }
    )

    // 当前对话ID
    const activeConversationId = computed({
        get: () => (store.activeConversationId ? String(store.activeConversationId) : null),
        set: (value) => store.setActiveConversation(value ? Number(value) : null),
    })

    // 对话选项列表
    const conversationOptions = computed(() =>
        store.conversations.map((item) => ({
            label: item.title || '新对话',
            value: String(item.id),
        }))
    )

    // 当前对话标题
    const conversationTitle = computed(() => {
        if (!store.activeConversationId) return '新对话'
        const current = store.conversations.find((item) => item.id === store.activeConversationId)
        return current?.title || '新对话'
    })

    // 加载消息列表
    const { data: messagesData, refresh: refreshMessages } = useAsyncData(
        'messages',
        () => {
            if (!store.activeConversationId) return Promise.resolve({ items: [] })
            return $fetch<{ items: MessageItem[] }>(
                `/api/conversations/${store.activeConversationId}/messages`,
                { headers: requestHeaders }
            )
        },
        { watch: [() => store.activeConversationId] }
    )

    // 监听消息数据变化
    watch(
        () => messagesData.value?.items,
        (items) => {
            if (!items || !store.activeConversationId) return
            store.setMessages(store.activeConversationId, items)
        },
        { immediate: true }
    )

    // 转换为Chat组件所需格式
    const hydratedMessages = computed(() => {
        if (!store.activeConversationId) return []
        const items = store.messagesByConversation[String(store.activeConversationId)] ?? []
        return items.map((item) => ({
            id: String(item.id),
            role: item.role,
            parts: [{ type: 'text', text: item.content }],
        }))
    })

    // 创建Chat实例
    const chat = new Chat({
        onError(error) {
            console.error('Chat error:', error)
        },
        onFinish: async ({ message, isError, isDisconnect, isAbort }) => {
            if (isError || isDisconnect || isAbort) return
            await persistAssistantMessage(message)
            await refreshMessages()
        },
    })

    // 同步消息到Chat实例
    watch(
        hydratedMessages,
        (items) => {
            // @ts-ignore - ai-sdk Chat accepts direct assign in runtime
            chat.messages = items
        },
        { immediate: true }
    )

    // 持久化助手消息
    async function persistAssistantMessage(message: any) {
        if (!store.activeConversationId) return
        if (!message || message.role !== 'assistant') return

        const parts = Array.isArray(message.parts) ? message.parts : []
        const text = parts
            .filter((part: any) => part?.type === 'text' && typeof part.text === 'string')
            .map((part: any) => part.text)
            .join('')

        if (!text.trim()) return

        await $fetch(`/api/conversations/${store.activeConversationId}/messages`, {
            method: 'POST',
            body: { role: 'assistant', content: text },
        })
    }

    // CRUD操作
    async function createConversation(title?: string) {
        const created = await $fetch<ConversationItem>(`/api/projects/${projectId}/conversations`, {
            method: 'POST',
            body: { title: title || '新对话' },
        })
        store.setConversations([created, ...store.conversations])
        store.setActiveConversation(created.id)
        // @ts-ignore
        chat.messages = []
        await refreshConversations()
        await refreshMessages()
    }

    async function updateConversation(conversationId: number, title: string) {
        const updated = await $fetch<ConversationItem>(`/api/conversations/${conversationId}`, {
            method: 'PATCH',
            body: { title },
        })

        store.setConversations(
            store.conversations.map((item) =>
                item.id === updated.id
                    ? { ...item, title: updated.title, updated_at: updated.updated_at }
                    : item
            )
        )
    }

    async function deleteConversation(conversationId: number) {
        await $fetch(`/api/conversations/${conversationId}`, { method: 'DELETE' })
        const next = store.conversations.filter((item) => item.id !== conversationId)
        store.setConversations(next)
        store.setActiveConversation(next[0]?.id ?? null)
        await refreshConversations()
        await refreshMessages()
    }

    async function sendMessage(text: string) {
        if (!text.trim()) return

        // 如果没有活动对话，先创建一个
        if (!store.activeConversationId) {
            await createConversation(text.slice(0, 20))
        }

        const conversationId = store.activeConversationId as number

        // 保存用户消息
        await $fetch(`/api/conversations/${conversationId}/messages`, {
            method: 'POST',
            body: { role: 'user', content: text },
        })

        // 发送到AI
        chat.sendMessage({ text })

        // 刷新消息列表
        await refreshMessages()
    }

    return {
        // 数据
        chat,
        activeConversationId,
        conversationOptions,
        conversationTitle,
        hydratedMessages,

        // 方法
        createConversation,
        updateConversation,
        deleteConversation,
        sendMessage,
        refreshConversations,
        refreshMessages,
    }
}
