import { Chat } from "@ai-sdk/vue";
import { useProjectWorkspaceStore } from "#imports";

type ConversationItem = {
  id: number;
  title?: string | null;
  created_at?: string;
  updated_at?: string;
};

type MessageItem = {
  id: number;
  conversation_id: number;
  role: string;
  content: string;
  created_at?: string;
};

type ComposableState = {
  loading: boolean;
  error: string | null;
};

export function useProjectConversations(projectId: string) {
  const store = useProjectWorkspaceStore();
  const requestHeaders = useRequestHeaders(["cookie"]);

  // Toast
  const toast = useToast();

  // 状态管理
  const state = ref<ComposableState>({
    loading: false,
    error: null,
  });

  // 加载对话列表
  const { data: conversationsData, refresh: refreshConversations } =
    useAsyncData("conversations", () =>
      $fetch<{ items: ConversationItem[] }>(
        `/api/projects/${projectId}/conversations`,
        { headers: requestHeaders },
      ).catch((error) => {
        state.value.error = `加载对话列表失败: ${error.message}`;
        console.error("Failed to load conversations:", error);
        return { items: [] };
      }),
    );

  // 监听对话数据变化
  watch(
    () => conversationsData.value?.items,
    (items) => {
      if (!items) return;
      store.setConversations(items);
      if (!store.activeConversationId && items.length > 0) {
        store.setActiveConversation(items[0].id);
      }
      state.value.error = null;
    },
    { immediate: true },
  );

  // 当前对话ID
  const activeConversationId = computed({
    get: () =>
      store.activeConversationId ? String(store.activeConversationId) : null,
    set: (value) => store.setActiveConversation(value ? Number(value) : null),
  });

  // 对话选项列表
  const conversationOptions = computed(() =>
    store.conversations.map((item) => ({
      label: item.title || "新对话",
      value: String(item.id),
    })),
  );

  // 当前对话标题
  const conversationTitle = computed(() => {
    if (!store.activeConversationId) return "新对话";
    const current = store.conversations.find(
      (item) => item.id === store.activeConversationId,
    );
    return current?.title || "新对话";
  });

  // 加载消息列表
  const { data: messagesData, refresh: refreshMessages } = useAsyncData(
    "messages",
    () => {
      if (!store.activeConversationId) return Promise.resolve({ items: [] });
      return $fetch<{ items: MessageItem[] }>(
        `/api/conversations/${store.activeConversationId}/messages`,
        { headers: requestHeaders },
      ).catch((error) => {
        state.value.error = `加载消息失败: ${error.message}`;
        console.error("Failed to load messages:", error);
        return { items: [] };
      });
    },
    { watch: [() => store.activeConversationId] },
  );

  // 监听消息数据变化
  watch(
    () => messagesData.value?.items,
    (items) => {
      if (!items || !store.activeConversationId) return;
      store.setMessages(store.activeConversationId, items);
      state.value.error = null;
    },
    { immediate: true },
  );

  // 转换为Chat组件所需格式
  const hydratedMessages = computed(() => {
    if (!store.activeConversationId) return [];
    const items =
      store.messagesByConversation[String(store.activeConversationId)] ?? [];
    return items.map((item) => ({
      id: String(item.id),
      role: item.role,
      parts: [{ type: "text", text: item.content }],
      context: (item as any).context ?? [],
    }));
  });

  // 创建Chat实例
  const chat = new Chat({
    onError(error: any) {
      const status =
        error?.statusCode ||
        error?.status ||
        error?.cause?.status ||
        error?.response?.status;
      const rawMessage = String(error?.statusMessage || error?.message || "");
      const lower = rawMessage.toLowerCase();
      let message = rawMessage ? `聊天错误: ${rawMessage}` : "聊天错误";
      if (
        status === 402 ||
        lower.includes("payment required") ||
        lower.includes("insufficient balance")
      ) {
        message = "余额不足，请先充值后再试";
      } else if (status === 409 || lower.includes("ref_id conflict")) {
        message = "计费信息冲突，请稍后重试";
      }
      console.log(message);
      state.value.error = message;
      // 尽可能写回Chat实例，便于UI展示
      // try {
      //     // @ts-ignore - runtime允许赋值
      //     chat.error = message
      // } catch {
      //     // ignore
      // }
      toast.add({
        title: "错误",
        description: message,
        color: "error",
        icon: "i-lucide-circle-x",
      });
      console.error("Chat error:", error);
    },
    onFinish: async ({ message, isError, isDisconnect, isAbort }) => {
      if (isError || isDisconnect || isAbort) {
        const parts = Array.isArray(message?.parts) ? message.parts : [];
        const errorText = parts
          .map((part: any) => (typeof part?.text === "string" ? part.text : ""))
          .join(" ")
          .trim();
        const lower = errorText.toLowerCase();
        let display = errorText || "AI 响应失败";
        if (
          lower.includes("payment required") ||
          lower.includes("insufficient balance")
        ) {
          display = "余额不足，请先充值后再试";
        }
        state.value.error = display;
        // try {
        //     // @ts-ignore
        //     chat.error = display
        // } catch {
        //     // ignore
        // }
        return;
      }
      await persistAssistantMessage(message);
      await refreshMessages();
    },
  });

  // 同步消息到Chat实例
  watch(
    hydratedMessages,
    (items) => {
      // @ts-ignore - ai-sdk Chat accepts direct assign in runtime
      chat.messages = items;
    },
    { immediate: true },
  );

  // 持久化助手消息
  async function persistAssistantMessage(message: any) {
    if (!store.activeConversationId) return;
    if (!message || message.role !== "assistant") return;

    const parts = Array.isArray(message.parts) ? message.parts : [];
    const text = parts
      .filter(
        (part: any) => part?.type === "text" && typeof part.text === "string",
      )
      .map((part: any) => part.text)
      .join("");

    if (!text.trim()) return;

    try {
      await $fetch(
        `/api/conversations/${store.activeConversationId}/messages`,
        {
          method: "POST",
          body: { role: "assistant", content: text },
        },
      );
      state.value.error = null;
    } catch (error: any) {
      state.value.error = `保存消息失败: ${error.message}`;
      console.error("Failed to persist message:", error);
    }
  }

  // CRUD操作
  async function createConversation(title?: string) {
    state.value.loading = true;
    try {
      const created = await $fetch<ConversationItem>(
        `/api/projects/${projectId}/conversations`,
        {
          method: "POST",
          body: { title: title || "新对话" },
        },
      );
      store.setConversations([created, ...store.conversations]);
      store.setActiveConversation(created.id);
      // @ts-ignore
      chat.messages = [];
      await refreshConversations();
      await refreshMessages();
      state.value.error = null;
      return created;
    } catch (error: any) {
      state.value.error = `创建对话失败: ${error.message}`;
      console.error("Failed to create conversation:", error);
      throw error;
    } finally {
      state.value.loading = false;
    }
  }

  async function updateConversation(conversationId: number, title: string) {
    state.value.loading = true;
    try {
      const updated = await $fetch<ConversationItem>(
        `/api/conversations/${conversationId}`,
        {
          method: "PATCH",
          body: { title },
        },
      );

      store.setConversations(
        store.conversations.map((item) =>
          item.id === updated.id
            ? { ...item, title: updated.title, updated_at: updated.updated_at }
            : item,
        ),
      );
      state.value.error = null;
    } catch (error: any) {
      state.value.error = `重命名对话失败: ${error.message}`;
      console.error("Failed to update conversation:", error);
      throw error;
    } finally {
      state.value.loading = false;
    }
  }

  async function deleteConversation(conversationId: number) {
    state.value.loading = true;
    try {
      await $fetch(`/api/conversations/${conversationId}`, {
        method: "DELETE",
      });
      const next = store.conversations.filter(
        (item) => item.id !== conversationId,
      );
      store.setConversations(next);
      store.setActiveConversation(next[0]?.id ?? null);
      await refreshConversations();
      await refreshMessages();
      state.value.error = null;
      return conversationId;
    } catch (error: any) {
      state.value.error = `删除对话失败: ${error.message}`;
      console.error("Failed to delete conversation:", error);
      throw error;
    } finally {
      state.value.loading = false;
    }
  }

  async function sendMessage(text: string, meta?: { context?: any[] }) {
    if (!text.trim()) return;

    state.value.loading = true;
    try {
      // 如果没有活动对话，先创建一个
      if (!store.activeConversationId) {
        await createConversation(text.slice(0, 20));
      }

      const conversationId = store.activeConversationId as number;

      // 保存用户消息
      await $fetch(`/api/conversations/${conversationId}/messages`, {
        method: "POST",
        body: { role: "user", content: text },
      });

      // 将上下文展示用数据附加到本地消息
      const existing =
        store.messagesByConversation[String(conversationId)] ?? [];
      store.setMessages(conversationId, [
        ...existing,
        {
          id: Date.now(),
          conversation_id: conversationId,
          role: "user",
          content: text,
          created_at: new Date().toISOString(),
          context: meta?.context ?? [],
        } as any,
      ]);

      // 发送到AI
      await chat.sendMessage({ text });

      // 刷新消息列表
      await refreshMessages();
      state.value.error = null;
    } catch (error: any) {
      const status =
        error?.statusCode ||
        error?.status ||
        error?.cause?.status ||
        error?.response?.status;
      const rawMessage = String(error?.statusMessage || error?.message || "");
      const lower = rawMessage.toLowerCase();
      let message = rawMessage ? `发送消息失败: ${rawMessage}` : "发送消息失败";
      if (
        status === 402 ||
        lower.includes("payment required") ||
        lower.includes("insufficient balance")
      ) {
        message = "余额不足，请先充值后再试";
      }
      state.value.error = message;
      try {
        // @ts-ignore
        chat.error = message;
      } catch {
        // ignore
      }
      console.error("Failed to send message:", error);
      throw error;
    } finally {
      state.value.loading = false;
    }
  }

  return {
    // 数据
    chat,
    activeConversationId,
    conversationOptions,
    conversationTitle,
    hydratedMessages,
    state,

    // 方法
    createConversation,
    updateConversation,
    deleteConversation,
    sendMessage,
    refreshConversations,
    refreshMessages,
  };
}
