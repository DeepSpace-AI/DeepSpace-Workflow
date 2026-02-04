import type { Ref } from 'vue'
import { useProjectWorkspaceStore } from '#imports'

type DocumentItem = {
    id: number
    project_id: number
    title: string
    content: string
    tags: string[]
    status: string
    created_at?: string
    updated_at?: string
}

export function useProjectDocuments(projectId: string) {
    const store = useProjectWorkspaceStore()
    const requestHeaders = useRequestHeaders(['cookie'])

    // 加载文档列表
    const { data: documentsData, refresh: refreshDocuments } = useAsyncData('documents', () =>
        $fetch<{ items: DocumentItem[] }>(`/api/projects/${projectId}/documents`, { headers: requestHeaders })
    )

    // 监听文档数据变化并更新store
    watch(
        () => documentsData.value?.items,
        async (items) => {
            if (!items) return
            store.mergeRemoteDocuments(items)
            if (!store.activeDocId && items.length > 0) {
                store.setActiveDoc(String(items[0].id))
            }
            if (items.length === 0) {
                await createDocument('新建文档')
            }
        },
        { immediate: true }
    )

    // 当前活动文档
    const activeDocument = computed(() => store.activeDocument)
    const documentTitle = computed(() => activeDocument.value?.title || '未命名文档')
    const documentCharCount = computed(() => (activeDocument.value?.content || '').length)
    const documentLastSavedAt = computed(() => activeDocument.value?.updatedAt || store.sync.lastSavedAt)

    // 文档内容的双向绑定
    const content = computed({
        get: () => activeDocument.value?.content ?? '',
        set: (value: string) => {
            if (!activeDocument.value) return
            store.updateDocumentContent(String(activeDocument.value.id), value)
            scheduleLocalPersist()
            scheduleRemotePersist()
        },
    })

    // 自动保存定时器
    let localPersistTimer: ReturnType<typeof setTimeout> | null = null
    let remotePersistTimer: ReturnType<typeof setTimeout> | null = null

    function scheduleLocalPersist() {
        if (localPersistTimer) clearTimeout(localPersistTimer)
        localPersistTimer = setTimeout(() => {
            store.persistLocal()
        }, 1500)
    }

    function scheduleRemotePersist() {
        if (remotePersistTimer) clearTimeout(remotePersistTimer)
        remotePersistTimer = setTimeout(() => {
            flushDirtyDocuments()
        }, 2000)
    }

    // 持久化脏文档
    async function flushDirtyDocuments() {
        const dirtyIds = store.dirtyDocumentIds
        if (!dirtyIds.length) return
        if (store.sync.saving) return

        store.sync.saving = true
        store.sync.lastError = undefined
        try {
            for (const docId of dirtyIds) {
                const doc = store.documents[docId]
                if (!doc) continue
                const updated = await $fetch<DocumentItem>(
                    `/api/projects/${projectId}/documents/${docId}`,
                    {
                        method: 'PATCH',
                        body: { title: doc.title, content: doc.content, tags: doc.tags },
                    }
                )
                store.markDocumentSaved(docId, updated.updated_at)
            }
            store.sync.lastSavedAt = new Date().toISOString()
        } catch (error: any) {
            store.sync.lastError = error?.message || '保存失败'
        } finally {
            store.sync.saving = false
        }
    }

    // 监听store变化，触发本地持久化
    watch(
        () => store.$state,
        () => scheduleLocalPersist(),
        { deep: true }
    )

    // CRUD 操作
    async function createDocument(title: string) {
        const created = await $fetch<DocumentItem>(`/api/projects/${projectId}/documents`, {
            method: 'POST',
            body: { title, content: '', tags: [] },
        })
        store.upsertDocument({
            id: created.id,
            title: created.title,
            content: created.content || '',
            tags: created.tags || [],
            status: created.status || 'draft',
            updatedAt: created.updated_at,
            localUpdatedAt: created.updated_at,
            dirty: false,
        })
        store.setActiveDoc(String(created.id))
        await refreshDocuments()
    }

    async function updateDocument(docId: string, updates: Partial<Pick<DocumentItem, 'title' | 'tags'>>) {
        const updated = await $fetch<DocumentItem>(
            `/api/projects/${projectId}/documents/${docId}`,
            { method: 'PATCH', body: updates }
        )
        store.upsertDocument({
            id: updated.id,
            title: updated.title,
            content: updated.content || '',
            tags: updated.tags || [],
            status: updated.status || 'draft',
            updatedAt: updated.updated_at,
            localUpdatedAt: updated.updated_at,
            dirty: false,
        })
    }

    async function deleteDocument(docId: string) {
        await $fetch(`/api/projects/${projectId}/documents/${docId}`, { method: 'DELETE' })
        store.removeDocument(docId)
        await refreshDocuments()
    }

    // 文档列表（带搜索和排序）
    function getDocumentList(searchKeyword: Ref<string>) {
        return computed(() =>
            Object.values(store.documents)
                .filter((doc) => {
                    const keyword = searchKeyword.value.trim().toLowerCase()
                    if (!keyword) return true
                    const haystack = [doc.title, ...(doc.tags || [])].join(' ').toLowerCase()
                    return haystack.includes(keyword)
                })
                .sort((a, b) => {
                    const aTime = new Date(a.updatedAt || 0).getTime()
                    const bTime = new Date(b.updatedAt || 0).getTime()
                    return bTime - aTime
                })
        )
    }

    return {
        // 数据
        activeDocument,
        documentTitle,
        documentCharCount,
        documentLastSavedAt,
        content,
        documentsData,
        syncStatus: computed(() => store.sync),

        // 方法
        createDocument,
        updateDocument,
        deleteDocument,
        refreshDocuments,
        flushDirtyDocuments,
        getDocumentList,
        setActiveDoc: (docId: string) => store.setActiveDoc(docId),
    }
}
