import { Chat } from '@ai-sdk/vue'
import { DefaultChatTransport } from 'ai'
import type {
  ConversationParamState,
  StandaloneConversationItem,
  StandaloneMessageItem,
} from '~/app/stores/standaloneChat'

function fallbackTitleFromText(text: string) {
  const trimmed = text.trim().replace(/\s+/g, ' ')
  if (!trimmed) return '新对话'
  return trimmed.length > 24 ? `${trimmed.slice(0, 24)}…` : trimmed
}

function extractAssistantText(message: any) {
  const parts = Array.isArray(message?.parts) ? message.parts : []
  return parts
    .filter((part: any) => part?.type === 'text' && typeof part?.text === 'string')
    .map((part: any) => part.text)
    .join('')
    .trim()
}

export function useStandaloneChatConversations() {
  const requestHeaders = useRequestHeaders(['cookie'])
  const store = useStandaloneChatStore()
  store.loadLocal()
  const drawingLoading = ref(false)

  const activeConversationIdNumber = computed(() => store.activeConversationId)

  const activeParams = computed(() => store.getParams(store.activeConversationId))

  function updateActiveParams(patch: Partial<ConversationParamState>) {
    store.setParams(store.activeConversationId, patch)
  }

  const { data: conversationsData, refresh: refreshConversationsData } = useAsyncData(
    'standalone:conversations:list',
    () =>
      $fetch<{ items: StandaloneConversationItem[] }>('/api/conversations', {
        headers: requestHeaders,
      }).catch(() => ({ items: [] })),
  )

  watch(
    () => conversationsData.value?.items,
    (items) => {
      if (!items) return
      store.setConversations(items)
      items.forEach((item) => store.ensureParams(item.id))

      if (store.activeConversationId && !items.some((item) => item.id === store.activeConversationId)) {
        store.setActiveConversation(items[0]?.id ?? null)
      }
      if (!store.activeConversationId && items.length > 0) {
        store.setActiveConversation(items[0].id)
      }
      store.error = null
    },
    { immediate: true },
  )

  const { data: messagesData, refresh: refreshMessagesData } = useAsyncData(
    () => `standalone:messages:${activeConversationIdNumber.value ?? 'none'}`,
    () => {
      if (!activeConversationIdNumber.value) return Promise.resolve({ items: [] as StandaloneMessageItem[] })
      return $fetch<{ items: StandaloneMessageItem[] }>(
        `/api/conversations/${activeConversationIdNumber.value}/messages`,
        {
          headers: requestHeaders,
        },
      ).catch(() => ({ items: [] }))
    },
    {
      watch: [activeConversationIdNumber],
    },
  )

  watch(
    () => messagesData.value?.items,
    (items) => {
      if (!items || !activeConversationIdNumber.value) return
      store.setMessages(activeConversationIdNumber.value, items)
      store.error = null
    },
    { immediate: true },
  )

  const activeConversationId = computed({
    get: () =>
      store.activeConversationId ? String(store.activeConversationId) : null,
    set: (value) => store.setActiveConversation(value ? Number(value) : null),
  })

  const conversationOptions = computed(() =>
    store.conversations.map((item) => ({
      label: item.title || '新对话',
      value: String(item.id),
      updated_at: item.updated_at,
    })),
  )

  const activeConversation = computed(
    () =>
      store.conversations.find(
        (item) => item.id === store.activeConversationId,
      ) ?? null,
  )

  const chat = new Chat({
    transport: new DefaultChatTransport({ api: '/api/chat/global' }),
    onError(error: any) {
      const status =
        error?.statusCode ||
        error?.status ||
        error?.cause?.status ||
        error?.response?.status
      const rawMessage = String(error?.statusMessage || error?.message || '')
      const lower = rawMessage.toLowerCase()
      let message = rawMessage ? `聊天错误: ${rawMessage}` : '聊天错误'
      if (
        status === 402 ||
        lower.includes('payment required') ||
        lower.includes('insufficient balance')
      ) {
        message = '余额不足，请先充值后再试'
      }
      store.error = message
    },
    async onFinish({ message, isError, isDisconnect, isAbort }) {
      if (isError || isDisconnect || isAbort) {
        return
      }

      const conversationId = store.activeConversationId
      if (!conversationId || message?.role !== 'assistant') return

      const text = extractAssistantText(message)
      if (!text) return

      try {
        await $fetch(`/api/conversations/${conversationId}/messages`, {
          method: 'POST',
          body: {
            role: 'assistant',
            content: text,
            model: store.selectedModel || undefined,
          },
        })

        await refreshMessages(conversationId)
      } catch (error: any) {
        store.error = `保存消息失败: ${error.message || '未知错误'}`
        return
      }

      const userSeed = getFirstUserMessage(conversationId)
      if (userSeed) {
        await generateTitleIfNeeded(conversationId, userSeed)
      }
    },
  })

  const hydratedMessages = computed(() => {
    if (!store.activeConversationId) return []
    const items =
      store.messagesByConversation[String(store.activeConversationId)] ?? []
    return items.map((item) => ({
      id: String(item.id),
      role: item.role,
      parts: [{ type: 'text', text: item.content }],
    }))
  })

  watch(
    hydratedMessages,
    (items) => {
      // @ts-expect-error runtime assignment supported by ai-sdk
      chat.messages = items
    },
    { immediate: true },
  )

  async function loadConversations() {
    store.loading = true
    try {
      await refreshConversationsData()
      return store.conversations
    } catch (error: any) {
      store.error = `加载对话列表失败: ${error.message || '未知错误'}`
      return []
    } finally {
      store.loading = false
    }
  }

  async function refreshConversations() {
    return loadConversations()
  }

  async function loadMessages(conversationId: number) {
    if (!conversationId) return []

    if (store.activeConversationId === conversationId) {
      store.loading = true
      try {
        await refreshMessagesData()
        return store.messagesByConversation[String(conversationId)] ?? []
      } finally {
        store.loading = false
      }
    }

    store.loading = true
    try {
      const response = await $fetch<{ items: StandaloneMessageItem[] }>(
        `/api/conversations/${conversationId}/messages`,
        {
          headers: requestHeaders,
        },
      )
      const items = response.items ?? []
      store.setMessages(conversationId, items)
      store.error = null
      return items
    } catch (error: any) {
      store.error = `加载消息失败: ${error.message || '未知错误'}`
      return []
    } finally {
      store.loading = false
    }
  }

  async function refreshMessages(conversationId: number) {
    return loadMessages(conversationId)
  }

  async function createConversation(title?: string) {
    store.loading = true
    try {
      const created = await $fetch<StandaloneConversationItem>(
        '/api/conversations',
        {
          method: 'POST',
          body: { title: title || '新对话' },
        },
      )
      store.upsertConversation(created)
      store.setActiveConversation(created.id)
      store.setMessages(created.id, [])
      store.ensureParams(created.id)
      store.error = null
      return created
    } catch (error: any) {
      store.error = `创建对话失败: ${error.message || '未知错误'}`
      throw error
    } finally {
      store.loading = false
    }
  }

  async function renameConversation(conversationId: number, title: string) {
    const normalized = title.trim()
    if (!normalized) return null

    try {
      const updated = await $fetch<StandaloneConversationItem>(
        `/api/conversations/${conversationId}`,
        {
          method: 'PATCH',
          body: { title: normalized },
        },
      )
      store.upsertConversation(updated)
      store.error = null
      return updated
    } catch (error: any) {
      store.error = `重命名失败: ${error.message || '未知错误'}`
      return null
    }
  }

  async function deleteConversation(conversationId: number) {
    store.loading = true
    try {
      await $fetch(`/api/conversations/${conversationId}`, {
        method: 'DELETE',
      })
      store.removeConversation(conversationId)
      store.error = null
      await refreshConversationsData()
    } catch (error: any) {
      store.error = `删除对话失败: ${error.message || '未知错误'}`
      throw error
    } finally {
      store.loading = false
    }
  }

  function getFirstUserMessage(conversationId: number) {
    const items = store.messagesByConversation[String(conversationId)] ?? []
    const firstUser = items.find(
      (item) => item.role === 'user' && item.content?.trim(),
    )
    return firstUser?.content?.trim() ?? ''
  }

  async function generateTitleIfNeeded(
    conversationId: number,
    userSeed: string,
  ) {
    const key = String(conversationId)
    const status = store.titleGenStatusByConversation[key]
    if (status === 'pending' || status === 'done') return

    const currentConversation = store.conversations.find(
      (item) => item.id === conversationId,
    )
    const currentTitle = (currentConversation?.title || '').trim()
    if (currentTitle && currentTitle !== '新对话') {
      store.setTitleGenStatus(conversationId, 'done')
      return
    }

    store.setTitleGenStatus(conversationId, 'pending')

    try {
      const response = await $fetch<{ title: string }>('/api/chat/title', {
        method: 'POST',
        body: {
          text: userSeed,
          model: store.selectedModel || undefined,
        },
      })

      const nextTitle = (response?.title || '').trim()
      if (!nextTitle) {
        throw new Error('标题为空')
      }

      const updated = await renameConversation(conversationId, nextTitle)
      if (!updated) {
        throw new Error('更新标题失败')
      }

      store.setTitleGenStatus(conversationId, 'done')
      await refreshConversationsData()
    } catch {
      const fallbackTitle = fallbackTitleFromText(userSeed)
      await renameConversation(conversationId, fallbackTitle)
      store.setTitleGenStatus(conversationId, 'error')
      await refreshConversationsData()
    }
  }

  async function sendMessage(text: string) {
    const trimmed = text.trim()
    if (!trimmed) return

    let conversationId = store.activeConversationId
    if (!conversationId) {
      const created = await createConversation('新对话')
      conversationId = created.id
    }

    const params = store.getParams(conversationId)

    store.loading = true
    try {
      await $fetch(`/api/conversations/${conversationId}/messages`, {
        method: 'POST',
        body: {
          role: 'user',
          content: trimmed,
          model: store.selectedModel || undefined,
        },
      })

      await refreshMessages(conversationId)

      await chat.sendMessage(
        { text: trimmed },
        {
          body: {
            model: store.selectedModel || undefined,
            max_tokens: params.maxTokenBudget,
            temperature: params.temperature,
            top_p: params.topP,
            reasoning_effort: params.thinkingDepth ? params.reasoningEffort : undefined,
          },
        },
      )

      await refreshMessages(conversationId)
      store.error = null
    } catch (error: any) {
      const status =
        error?.statusCode ||
        error?.status ||
        error?.cause?.status ||
        error?.response?.status
      const rawMessage = String(error?.statusMessage || error?.message || '')
      const lower = rawMessage.toLowerCase()
      let message = rawMessage ? `发送消息失败: ${rawMessage}` : '发送消息失败'
      if (
        status === 402 ||
        lower.includes('payment required') ||
        lower.includes('insufficient balance')
      ) {
        message = '余额不足，请先充值后再试'
      }
      store.error = message
      throw error
    } finally {
      store.loading = false
    }
  }

  async function sendDrawing(prompt: string) {
    const trimmed = prompt.trim()
    if (!trimmed) return

    let conversationId = store.activeConversationId
    if (!conversationId) {
      const created = await createConversation('新对话')
      conversationId = created.id
    }

    const params = store.getParams(conversationId)

    drawingLoading.value = true
    store.loading = true
    try {
      await $fetch(`/api/conversations/${conversationId}/messages`, {
        method: 'POST',
        body: {
          role: 'user',
          content: trimmed,
          model: store.selectedModel || undefined,
        },
      })

      await refreshMessages(conversationId)

      const imageResp = await $fetch<{
        imageUrl: string
        revisedPrompt?: string
        model: string
      }>('/api/chat/image', {
        method: 'POST',
        body: {
          prompt: trimmed,
          model: params.imageModel || undefined,
          size: params.imageSize,
          quality: params.imageQuality,
          style: params.imageStyle,
          n: 1,
        },
      })

      const assistantParts = [`![绘画结果](${imageResp.imageUrl})`]
      if (imageResp.revisedPrompt) {
        assistantParts.push(`\n\n> 提示词优化：${imageResp.revisedPrompt}`)
      }

      await $fetch(`/api/conversations/${conversationId}/messages`, {
        method: 'POST',
        body: {
          role: 'assistant',
          content: assistantParts.join(''),
          model: imageResp.model || params.imageModel || undefined,
        },
      })

      await refreshMessages(conversationId)
      store.error = null

      const userSeed = getFirstUserMessage(conversationId)
      if (userSeed) {
        await generateTitleIfNeeded(conversationId, userSeed)
      }
    } catch (error: any) {
      const status =
        error?.statusCode ||
        error?.status ||
        error?.cause?.status ||
        error?.response?.status
      const rawMessage = String(error?.statusMessage || error?.message || '')
      const lower = rawMessage.toLowerCase()
      let message = rawMessage ? `绘画生成失败: ${rawMessage}` : '绘画生成失败'
      if (
        status === 402 ||
        lower.includes('payment required') ||
        lower.includes('insufficient balance')
      ) {
        message = '余额不足，请先充值后再试'
      }
      store.error = message
      throw error
    } finally {
      drawingLoading.value = false
      store.loading = false
    }
  }

  return {
    store,
    chat,
    activeConversation,
    activeConversationId,
    activeParams,
    updateActiveParams,
    conversationOptions,
    hydratedMessages,
    loadConversations,
    refreshConversations,
    loadMessages,
    refreshMessages,
    createConversation,
    renameConversation,
    deleteConversation,
    sendMessage,
    sendDrawing,
    drawingLoading,
    generateTitleIfNeeded,
  }
}
