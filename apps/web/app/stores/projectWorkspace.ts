import { defineStore } from "pinia";

export type ProjectDocumentState = {
  id: number;
  title: string;
  content: string;
  tags: string[];
  status: string;
  updatedAt?: string;
  localUpdatedAt?: string;
  dirty: boolean;
};

export type ConversationState = {
  id: number;
  title?: string | null;
  updatedAt?: string;
};

export type MessageState = {
  id: number;
  conversation_id: number;
  role: string;
  content: string;
  created_at?: string;
};

type WorkspaceSnapshot = {
  projectId: string;
  documents: Record<string, ProjectDocumentState>;
  activeDocId: string | null;
  conversations: ConversationState[];
  activeConversationId: number | null;
  messagesByConversation: Record<string, MessageState[]>;
};

const STORAGE_VERSION = 1;

export const useProjectWorkspaceStore = defineStore("project-workspace", () => {
  const projectId = ref("");
  const documents = ref<Record<string, ProjectDocumentState>>({});
  const activeDocId = ref<string | null>(null);
  const conversations = ref<ConversationState[]>([]);
  const activeConversationId = ref<number | null>(null);
  const messagesByConversation = ref<Record<string, MessageState[]>>({});
  const sync = ref({
    saving: false,
    lastSavedAt: undefined as string | undefined,
    lastError: undefined as string | undefined,
  });

  const dirtyDocumentIds = computed(() =>
    Object.values(documents.value)
      .filter((doc) => doc.dirty)
      .map((doc) => String(doc.id))
  );

  const activeDocument = computed(() => {
    if (!activeDocId.value) return null;
    return documents.value[activeDocId.value] ?? null;
  });

  function setProject(value: string) {
    projectId.value = value;
  }

  function loadLocal() {
    if (!process.client || !projectId.value) return;
    const key = storageKey(projectId.value);
    const raw = localStorage.getItem(key);
    if (!raw) return;
    try {
      const parsed = JSON.parse(raw) as { version?: number; data?: WorkspaceSnapshot };
      if (!parsed?.data) return;
      const snapshot = parsed.data;
      if (snapshot.projectId !== projectId.value) return;
      documents.value = snapshot.documents || {};
      activeDocId.value = snapshot.activeDocId ?? null;
      conversations.value = snapshot.conversations || [];
      activeConversationId.value = snapshot.activeConversationId ?? null;
      messagesByConversation.value = snapshot.messagesByConversation || {};
    } catch {
      // ignore invalid local cache
    }
  }

  function persistLocal() {
    if (!process.client || !projectId.value) return;
    const snapshot: WorkspaceSnapshot = {
      projectId: projectId.value,
      documents: documents.value,
      activeDocId: activeDocId.value,
      conversations: conversations.value,
      activeConversationId: activeConversationId.value,
      messagesByConversation: messagesByConversation.value,
    };
    const payload = JSON.stringify({ version: STORAGE_VERSION, data: snapshot });
    localStorage.setItem(storageKey(projectId.value), payload);
  }

  function mergeRemoteDocuments(remote: Array<any>) {
    const next: Record<string, ProjectDocumentState> = { ...documents.value };
    remote.forEach((doc) => {
      const key = String(doc.id);
      const local = next[key];
      const remoteUpdatedAt = doc.updated_at || doc.updatedAt;
      if (local?.localUpdatedAt && remoteUpdatedAt) {
        const localTime = new Date(local.localUpdatedAt).getTime();
        const remoteTime = new Date(remoteUpdatedAt).getTime();
        if (localTime > remoteTime) {
          return;
        }
      }
      next[key] = {
        id: doc.id,
        title: doc.title,
        content: doc.content || "",
        tags: Array.isArray(doc.tags) ? doc.tags : [],
        status: doc.status || "draft",
        updatedAt: remoteUpdatedAt,
        localUpdatedAt: remoteUpdatedAt,
        dirty: false,
      };
    });
    documents.value = next;
  }

  function setActiveDoc(docId: string | null) {
    activeDocId.value = docId;
  }

  function upsertDocument(doc: ProjectDocumentState) {
    documents.value = { ...documents.value, [String(doc.id)]: doc };
  }

  function updateDocumentTags(docId: string, tags: string[]) {
    const doc = documents.value[docId];
    if (!doc) return;
    documents.value = {
      ...documents.value,
      [docId]: {
        ...doc,
        tags,
        dirty: true,
      },
    };
  }

  function removeDocument(docId: string) {
    const next = { ...documents.value };
    delete next[docId];
    documents.value = next;
    if (activeDocId.value === docId) {
      activeDocId.value = null;
    }
  }

  function updateDocumentContent(docId: string, content: string) {
    const doc = documents.value[docId];
    if (!doc) return;
    const now = new Date().toISOString();
    documents.value = {
      ...documents.value,
      [docId]: {
        ...doc,
        content,
        localUpdatedAt: now,
        dirty: true,
      },
    };
  }

  function markDocumentSaved(docId: string, updatedAt?: string) {
    const doc = documents.value[docId];
    if (!doc) return;
    documents.value = {
      ...documents.value,
      [docId]: {
        ...doc,
        updatedAt,
        localUpdatedAt: updatedAt || doc.localUpdatedAt,
        dirty: false,
      },
    };
  }

  function setConversations(items: ConversationState[]) {
    conversations.value = items;
  }

  function setActiveConversation(id: number | null) {
    activeConversationId.value = id;
  }

  function setMessages(conversationId: number, items: MessageState[]) {
    messagesByConversation.value = {
      ...messagesByConversation.value,
      [String(conversationId)]: items,
    };
  }

  function storageKey(id: string) {
    return `project:${id}:workspace`;
  }

  return {
    projectId,
    documents,
    activeDocId,
    conversations,
    activeConversationId,
    messagesByConversation,
    sync,
    dirtyDocumentIds,
    activeDocument,
    setProject,
    loadLocal,
    persistLocal,
    mergeRemoteDocuments,
    setActiveDoc,
    upsertDocument,
    removeDocument,
    updateDocumentContent,
    updateDocumentTags,
    markDocumentSaved,
    setConversations,
    setActiveConversation,
    setMessages,
  };
});
