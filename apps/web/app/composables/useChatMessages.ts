/**
 * 聊天消息管理 Composable
 * 职责：仅处理指定对话的消息加载、持久化等
 */

type MessageItem = {
  id: number;
  conversation_id: number;
  role: string;
  content: string;
  model?: string | null;
  created_at?: string;
};

type MessageState = {
  loading: boolean;
  error: string | null;
};

export function useChatMessages(conversationId: Ref<number | null>) {
  const requestHeaders = useRequestHeaders(["cookie"]);

  const state = ref<MessageState>({
    loading: false,
    error: null,
  });

  const messages = ref<MessageItem[]>([]);

  const { data: messagesData, refresh: refreshMessages } = useAsyncData(
    "chat-messages",
    () => {
      if (!conversationId.value) return Promise.resolve({ items: [] });
      return $fetch<{ items: MessageItem[] }>(
        `/api/conversations/${conversationId.value}/messages`,
        { headers: requestHeaders },
      ).catch((error) => {
        state.value.error = `加载消息失败: ${error.message}`;
        console.error("Failed to load messages:", error);
        return { items: [] };
      });
    },
    { watch: [conversationId] },
  );

  watch(
    () => messagesData.value?.items,
    (items) => {
      if (!items) return;
      messages.value = items;
      state.value.error = null;
    },
    { immediate: true },
  );

  // 为 Chat 组件转换消息格式
  const hydratedMessages = computed(() => {
    return messages.value.map((item) => ({
      id: String(item.id),
      role: item.role,
      parts: [{ type: "text", text: item.content }],
    }));
  });

  async function persistMessage(
    content: string,
    role: string,
    model?: string | null,
  ) {
    if (!conversationId.value) return;
    if (!content.trim()) return;

    try {
      state.value.loading = true;
      await $fetch(
        `/api/conversations/${conversationId.value}/messages`,
        {
          method: "POST",
          body: { role, content, model: model || undefined },
          headers: requestHeaders,
        },
      );
      state.value.error = null;
      await refreshMessages();
    } catch (error: any) {
      state.value.error = `保存消息失败: ${error.message}`;
      console.error("Failed to persist message:", error);
      throw error;
    } finally {
      state.value.loading = false;
    }
  }

  return {
    messages,
    hydratedMessages,
    state,
    persistMessage,
    refresh: refreshMessages,
  };
}
