<script setup lang="ts">
definePageMeta({
  layout: 'editor',
})
useHead({
  title: "智能助手 - Deepspace Workflow",
})



const conversations = useStandaloneChatConversations()
const models = useModels()

onMounted(async () => {
  await conversations.loadConversations()
  if (!conversations.store.selectedModel && models.defaultModel.value) {
    conversations.store.selectedModel = models.defaultModel.value
  }
})

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

async function handleCreateConversation() {
  const created = await conversations.createConversation('新对话')
  await navigateTo(`/chat/${created.id}`)
}

async function handleSelectConversation(conversationId: string) {
  conversations.activeConversationId.value = conversationId
  await navigateTo(`/chat/${conversationId}`)
}

async function handleDeleteConversation(conversationId: string) {
  await conversations.deleteConversation(Number(conversationId))
}

async function handleSubmit(text: string) {
  const created = await conversations.createConversation('新对话')
  await navigateTo({
    path: `/chat/${created.id}`,
    query: { q: text },
  })
}

async function handleDraw(prompt: string) {
  await conversations.sendDrawing(prompt)
  const id = conversations.store.activeConversationId
  if (id) {
    await navigateTo(`/chat/${id}`)
  }
}
</script>

<template>
  <div class="flex min-h-[calc(100vh-4rem)] bg-neutral-100 dark:bg-neutral-950">
    <ChatSider
      :conversations="conversations.conversationOptions.value"
      :active-conversation-id="null"
      :loading="conversations.store.loading"
      @create="handleCreateConversation"
      @select="handleSelectConversation"
      @delete="handleDeleteConversation"
    />
    
    <ChatContent
      :chat-id="null"
      :messages="[]"
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
