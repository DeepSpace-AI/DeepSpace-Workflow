<script setup lang="ts">
type DocumentItem = {
    id: string | number
    title: string
    content: string
    tags: string[]
    status: string
    updatedAt?: string
    dirty?: boolean
}

type SyncStatus = {
    saving: boolean
    lastSavedAt?: string
    lastError?: string
}

interface Props {
    projectId: string
    documents: DocumentItem[]
    activeDocId: string | null
    documentTitle: string
    documentCharCount: number
    documentLastSavedAt?: string
    content: string
    search: string
    syncStatus: SyncStatus
}

interface Emits {
    (e: 'update:content', value: string): void
    (e: 'update:search', value: string): void
    (e: 'select-doc', docId: string): void
    (e: 'create-doc'): void
    (e: 'rename-doc', docId: string): void
    (e: 'edit-tags', docId: string): void
    (e: 'delete-doc', docId: string): void
    (e: 'rename-title'): void
    (e: 'ai-action', payload: { type: 'polish' | 'expand' | 'summary'; selection: string; fullText: string }): void
    (e: 'ai-apply'): void
    (e: 'ai-dismiss'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const editorRef = ref<any>(null)

// 暴露 editorRef 供父组件使用
defineExpose({ editorRef })

const localContent = computed({
    get: () => props.content,
    set: (value) => emit('update:content', value),
})

const localSearch = computed({
    get: () => props.search,
    set: (value) => emit('update:search', value),
})
</script>

<template>
    <div class="flex-1 flex h-full">
        <!-- 文档列表 -->
        <ProjectResources
            :project-id="projectId"
            :documents="documents"
            :active-doc-id="activeDocId"
            :search="localSearch"
            :sync-status="syncStatus"
            @update:search="localSearch = $event"
            @select-doc="emit('select-doc', $event)"
            @create-doc="emit('create-doc')"
            @rename-doc="emit('rename-doc', $event)"
            @edit-tags="emit('edit-tags', $event)"
            @delete-doc="emit('delete-doc', $event)"
        />

        <!-- 编辑器 -->
        <div class="flex-1 flex flex-col">
            <!-- 标题栏 -->
            <div class="flex items-center justify-between border-b border-slate-200/70 px-6 py-3 text-xs text-slate-500 dark:border-slate-800 dark:text-slate-400">
                <div class="flex items-center gap-3">
                    <div
                        class="text-sm font-semibold text-slate-900 dark:text-white cursor-text select-none"
                        role="button"
                        @dblclick="emit('rename-title')"
                    >
                        {{ documentTitle }}
                    </div>
                    <UBadge color="neutral" variant="soft">{{ documentCharCount }} 字符</UBadge>
                </div>
                <div class="flex items-center gap-2">
                    <span>最后保存：</span>
                    <span>{{ documentLastSavedAt ? new Date(documentLastSavedAt).toLocaleString() : '未保存' }}</span>
                </div>
            </div>

            <!-- 文档编辑器 -->
            <EditorDocumentEditor
                ref="editorRef"
                v-model="localContent"
                placeholder="开始写作或输入/唤醒AI助手"
                :editable="true"
                @ai="emit('ai-action', $event)"
                @ai-apply="emit('ai-apply')"
                @ai-dismiss="emit('ai-dismiss')"
                class="flex-1"
            />
        </div>
    </div>
</template>
