export type ChatContextItem =
  | {
      type: "selection";
      id: string;
      label: string;
      content: string;
      sourceDocId?: string | number;
    }
  | {
      type: "file";
      id: string;
      label: string;
      kbId: number;
      docId: number;
      contentType?: string | null;
    }
  | {
      type: "skill";
      id: string;
      label: string;
      prompt?: string | null;
    }
  | {
      type: "workflow";
      id: string;
      label: string;
      steps?: unknown;
    };

export type ConversationContextState = {
  [conversationId: string]: {
    items: ChatContextItem[];
    skills: string[];
    workflows: string[];
  };
};

export function useChatContext(projectId: string) {
  const contexts = ref<ConversationContextState>({});
  const storageKey = `project:${projectId}:conversation-contexts`;
  let saveTimer: ReturnType<typeof setTimeout> | null = null;

  function ensureConversation(conversationId: string | null) {
    if (!conversationId) return;
    if (!contexts.value[conversationId]) {
      contexts.value[conversationId] = { items: [], skills: [], workflows: [] };
    }
  }

  function getContext(conversationId: string | null) {
    if (!conversationId) return { items: [], skills: [], workflows: [] };
    ensureConversation(conversationId);
    return contexts.value[conversationId];
  }

  function setContextItems(conversationId: string | null, items: ChatContextItem[]) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].items = items;
  }

  function addItem(conversationId: string | null, item: ChatContextItem) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].items = [
      item,
      ...contexts.value[conversationId].items,
    ];
  }

  function removeItem(conversationId: string | null, itemId: string) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].items = contexts.value[
      conversationId
    ].items.filter((item) => item.id !== itemId);
  }

  function setSkills(conversationId: string | null, skills: string[]) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].skills = skills;
  }

  function setWorkflows(conversationId: string | null, workflows: string[]) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].workflows = workflows;
  }

  function clearTransient(conversationId: string | null) {
    if (!conversationId) return;
    ensureConversation(conversationId);
    contexts.value[conversationId].items = contexts.value[
      conversationId
    ].items.filter((item) => item.type !== "selection" && item.type !== "file");
  }

  function removeConversation(conversationId: string) {
    const next = { ...contexts.value };
    delete next[conversationId];
    contexts.value = next;
  }

  function loadFromStorage() {
    if (import.meta.server) return;
    try {
      const raw = window.localStorage.getItem(storageKey);
      if (!raw) return;
      const parsed = JSON.parse(raw);
      if (parsed && typeof parsed === "object") {
        contexts.value = parsed;
      }
    } catch (error) {
      console.warn("Failed to load chat context:", error);
    }
  }

  function saveToStorage() {
    if (import.meta.server) return;
    if (saveTimer) clearTimeout(saveTimer);
    saveTimer = setTimeout(() => {
      try {
        window.localStorage.setItem(storageKey, JSON.stringify(contexts.value));
      } catch (error) {
        console.warn("Failed to save chat context:", error);
      }
    }, 400);
  }

  onMounted(loadFromStorage);

  watch(
    contexts,
    () => {
      saveToStorage();
    },
    { deep: true }
  );

  return {
    contexts,
    ensureConversation,
    getContext,
    setContextItems,
    addItem,
    removeItem,
    setSkills,
    setWorkflows,
    clearTransient,
    removeConversation,
  };
}
