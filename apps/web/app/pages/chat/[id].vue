<script setup lang="ts">
definePageMeta({
  layout: 'editor',
})

const route = useRoute()
const router = useRouter()
const conversations = useStandaloneChatConversations()
const models = useModels()

const syncingRoute = ref(false)
const consumedInitialQuery = ref<string | null>(null)

const chatId = computed(() => String(route.params.id || ''))

watch(
  () => models.defaultModel.value,
  (value) => {
    if (!conversations.store.selectedModel && value) {
      conversations.store.selectedModel = value
    }
  },
  { immediate: true },
)

watch(
  () => models.defaultImageModel.value,
  (value) => {
    if (!value) return
    const current = conversations.activeParams.value.imageModel
    if (!current) {
      conversations.updateActiveParams({ imageModel: value })
    }
  },
  { immediate: true },
)

function parseConversationId(value: string) {
  const id = Number(value)
  if (!Number.isInteger(id) || id <= 0) {
    return null
  }
  return id
}

async function syncRouteConversation() {
  if (syncingRoute.value) return
  syncingRoute.value = true

  try {
    const id = parseConversationId(chatId.value)
    if (!id) {
      await router.replace('/chat')
      return
    }

    await conversations.loadConversations()
    const exists = conversations.store.conversations.some((item) => item.id === id)
    if (!exists) {
      await router.replace('/chat')
      return
    }

    conversations.activeConversationId.value = String(id)
    await conversations.loadMessages(id)

    const seed = String(route.query.q || '').trim()
    if (seed && consumedInitialQuery.value !== `${id}:${seed}`) {
      consumedInitialQuery.value = `${id}:${seed}`
      await conversations.sendMessage(seed)
      await router.replace(`/chat/${id}`)
    }
  } finally {
    syncingRoute.value = false
  }
}

watch(
  () => route.params.id,
  async () => {
    await syncRouteConversation()
  },
  { immediate: true },
)

async function handleCreateConversation() {
  const created = await conversations.createConversation('新对话')
  await router.push(`/chat/${created.id}`)
}

async function handleSelectConversation(conversationId: string) {
  conversations.activeConversationId.value = conversationId
  await router.push(`/chat/${conversationId}`)
}

async function handleDeleteConversation(conversationId: string) {
  const id = Number(conversationId)
  await conversations.deleteConversation(id)

  const nextId = conversations.store.activeConversationId
  if (nextId) {
    await router.push(`/chat/${nextId}`)
    return
  }

  await router.push('/chat')
}

async function handleSubmit(text: string) {
  await conversations.sendMessage(text)
}

async function handleDraw(prompt: string) {
  await conversations.sendDrawing(prompt)
}
</script>

<template>
  <div class="flex h-[calc(100vh-var(--ui-header-height))] min-h-0 overflow-hidden">
    <ChatSider
      :conversations="conversations.conversationOptions.value"
      :active-conversation-id="conversations.activeConversationId.value"
      :loading="conversations.store.loading"
      @create="handleCreateConversation"
      @select="handleSelectConversation"
      @delete="handleDeleteConversation"
    />

    <ChatContent
      :chat-id="chatId"
      :messages="conversations.chat.messages"
      :status="conversations.chat.status"
      :error="conversations.store.error"
      :model="conversations.store.selectedModel"
      :models="models.menuItems"
      :image-model-options="models.imageMenuItems"
      :drawing-loading="conversations.drawingLoading.value"
      :session-config="conversations.activeParams.value"
      @submit="handleSubmit"
      @draw="handleDraw"
      @update:model="conversations.store.selectedModel = $event"
      @update:sessionConfig="conversations.updateActiveParams($event)"
    />
  </div>
</template>
