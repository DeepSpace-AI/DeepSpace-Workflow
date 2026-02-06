import { defineStore } from 'pinia'

export type StandaloneConversationItem = {
  id: number
  project_id?: number | null
  title?: string | null
  created_at?: string
  updated_at?: string
}

export type StandaloneMessageItem = {
  id: number
  conversation_id: number
  role: string
  content: string
  model?: string | null
  trace_id?: string | null
  created_at?: string
}

export type ConversationParamState = {
  enableChatCache: boolean
  thinkingDepth: boolean
  reasoningEffort: 'low' | 'medium' | 'high'
  imageModel: string
  imageSize: '1024x1024' | '1024x1792' | '1792x1024'
  imageQuality: 'standard' | 'hd'
  imageStyle: 'vivid' | 'natural'
  maxTokenBudget: number
  temperature: number
  topP: number
}

type StandaloneChatSnapshot = {
  conversations: StandaloneConversationItem[]
  activeConversationId: number | null
  messagesByConversation: Record<string, StandaloneMessageItem[]>
  selectedModel: string
  paramsByConversation: Record<string, ConversationParamState>
  defaultParams: ConversationParamState
}

const STORAGE_KEY = 'standalone-chat-workspace-v1'

function createDefaultParams(): ConversationParamState {
  return {
    enableChatCache: true,
    thinkingDepth: false,
    reasoningEffort: 'medium',
    imageModel: '',
    imageSize: '1024x1024',
    imageQuality: 'standard',
    imageStyle: 'vivid',
    maxTokenBudget: 1024,
    temperature: 0.7,
    topP: 1,
  }
}

function normalizeParams(input?: Partial<ConversationParamState> | null): ConversationParamState {
  const fallback = createDefaultParams()
  const source = input || {}
  const reasoningEffort =
    source.reasoningEffort === 'low' || source.reasoningEffort === 'medium' || source.reasoningEffort === 'high'
      ? source.reasoningEffort
      : fallback.reasoningEffort
  const imageSize =
    source.imageSize === '1024x1024' || source.imageSize === '1024x1792' || source.imageSize === '1792x1024'
      ? source.imageSize
      : fallback.imageSize
  const imageQuality =
    source.imageQuality === 'standard' || source.imageQuality === 'hd'
      ? source.imageQuality
      : fallback.imageQuality
  const imageStyle =
    source.imageStyle === 'vivid' || source.imageStyle === 'natural'
      ? source.imageStyle
      : fallback.imageStyle

  return {
    enableChatCache:
      typeof source.enableChatCache === 'boolean'
        ? source.enableChatCache
        : fallback.enableChatCache,
    thinkingDepth:
      typeof source.thinkingDepth === 'boolean'
        ? source.thinkingDepth
        : fallback.thinkingDepth,
    reasoningEffort,
    imageModel: typeof source.imageModel === 'string' ? source.imageModel : fallback.imageModel,
    imageSize,
    imageQuality,
    imageStyle,
    maxTokenBudget:
      typeof source.maxTokenBudget === 'number' && Number.isFinite(source.maxTokenBudget)
        ? Math.max(256, Math.min(64000, Math.round(source.maxTokenBudget)))
        : fallback.maxTokenBudget,
    temperature:
      typeof source.temperature === 'number' && Number.isFinite(source.temperature)
        ? Math.max(0, Math.min(2, source.temperature))
        : fallback.temperature,
    topP:
      typeof source.topP === 'number' && Number.isFinite(source.topP)
        ? Math.max(0, Math.min(1, source.topP))
        : fallback.topP,
  }
}

export const useStandaloneChatStore = defineStore('standalone-chat', () => {
  const conversations = ref<StandaloneConversationItem[]>([])
  const activeConversationId = ref<number | null>(null)
  const messagesByConversation = ref<Record<string, StandaloneMessageItem[]>>({})
  const selectedModel = ref('deepseek-chat')

  const paramsByConversation = ref<Record<string, ConversationParamState>>({})
  const defaultParams = ref<ConversationParamState>(createDefaultParams())

  const loading = ref(false)
  const error = ref<string | null>(null)
  const localLoaded = ref(false)
  const titleGenStatusByConversation = ref<Record<string, 'idle' | 'pending' | 'done' | 'error'>>({})

  function setConversations(items: StandaloneConversationItem[]) {
    conversations.value = items
  }

  function setActiveConversation(id: number | null) {
    activeConversationId.value = id
  }

  function setMessages(conversationId: number, items: StandaloneMessageItem[]) {
    messagesByConversation.value = {
      ...messagesByConversation.value,
      [String(conversationId)]: items,
    }
  }

  function ensureParams(conversationId: number | null) {
    if (!conversationId) return
    const key = String(conversationId)
    if (!paramsByConversation.value[key]) {
      paramsByConversation.value = {
        ...paramsByConversation.value,
        [key]: { ...defaultParams.value },
      }
    }
  }

  function getParams(conversationId: number | null): ConversationParamState {
    if (!conversationId) return { ...defaultParams.value }
    ensureParams(conversationId)
    return { ...paramsByConversation.value[String(conversationId)] }
  }

  function setParams(conversationId: number | null, patch: Partial<ConversationParamState>) {
    if (!patch || typeof patch !== 'object') return

    if (!conversationId) {
      defaultParams.value = normalizeParams({
        ...defaultParams.value,
        ...patch,
      })
      return
    }

    ensureParams(conversationId)
    const key = String(conversationId)
    paramsByConversation.value = {
      ...paramsByConversation.value,
      [key]: normalizeParams({
        ...paramsByConversation.value[key],
        ...patch,
      }),
    }
  }

  function removeConversation(conversationId: number) {
    const key = String(conversationId)
    conversations.value = conversations.value.filter((item) => item.id !== conversationId)

    const nextMessages = { ...messagesByConversation.value }
    delete nextMessages[key]
    messagesByConversation.value = nextMessages

    const nextTitleStatus = { ...titleGenStatusByConversation.value }
    delete nextTitleStatus[key]
    titleGenStatusByConversation.value = nextTitleStatus

    const nextParams = { ...paramsByConversation.value }
    delete nextParams[key]
    paramsByConversation.value = nextParams

    if (activeConversationId.value === conversationId) {
      activeConversationId.value = conversations.value[0]?.id ?? null
    }
  }

  function upsertConversation(item: StandaloneConversationItem) {
    const exists = conversations.value.some((value) => value.id === item.id)
    if (exists) {
      conversations.value = conversations.value.map((value) =>
        value.id === item.id ? { ...value, ...item } : value,
      )
    } else {
      conversations.value = [item, ...conversations.value]
    }

    ensureParams(item.id)
  }

  function setTitleGenStatus(conversationId: number, status: 'idle' | 'pending' | 'done' | 'error') {
    titleGenStatusByConversation.value = {
      ...titleGenStatusByConversation.value,
      [String(conversationId)]: status,
    }
  }

  function loadLocal() {
    if (!process.client || localLoaded.value) return
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      localLoaded.value = true
      return
    }

    try {
      const parsed = JSON.parse(raw) as StandaloneChatSnapshot
      conversations.value = Array.isArray(parsed?.conversations) ? parsed.conversations : []
      activeConversationId.value = parsed?.activeConversationId ?? null
      messagesByConversation.value = parsed?.messagesByConversation ?? {}
      selectedModel.value = parsed?.selectedModel || 'deepseek-chat'
      paramsByConversation.value = parsed?.paramsByConversation ?? {}
      defaultParams.value = normalizeParams(parsed?.defaultParams)
    } catch {
      // ignore bad cache
    }

    localLoaded.value = true
  }

  function persistLocal() {
    if (!process.client) return
    const payload: StandaloneChatSnapshot = {
      conversations: conversations.value,
      activeConversationId: activeConversationId.value,
      messagesByConversation: messagesByConversation.value,
      selectedModel: selectedModel.value,
      paramsByConversation: paramsByConversation.value,
      defaultParams: defaultParams.value,
    }
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(payload))
  }

  let saveTimer: ReturnType<typeof setTimeout> | null = null

  watch(
    [
      conversations,
      activeConversationId,
      messagesByConversation,
      selectedModel,
      paramsByConversation,
      defaultParams,
    ],
    () => {
      if (!process.client) return
      if (saveTimer) clearTimeout(saveTimer)
      saveTimer = setTimeout(() => persistLocal(), 300)
    },
    { deep: true },
  )

  return {
    conversations,
    activeConversationId,
    messagesByConversation,
    selectedModel,
    paramsByConversation,
    defaultParams,
    loading,
    error,
    localLoaded,
    titleGenStatusByConversation,
    setConversations,
    setActiveConversation,
    setMessages,
    ensureParams,
    getParams,
    setParams,
    removeConversation,
    upsertConversation,
    setTitleGenStatus,
    loadLocal,
    persistLocal,
  }
})
