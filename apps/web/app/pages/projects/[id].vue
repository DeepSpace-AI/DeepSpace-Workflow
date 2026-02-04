<script setup lang="ts">
import { useProjectWorkspaceStore } from '#imports'

definePageMeta({
    layout: 'editor',
})

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

const route = useRoute()
const projectId = String(route.params.id)
const requestHeaders = useRequestHeaders(['cookie'])
const store = useProjectWorkspaceStore()

// 初始化store
store.setProject(projectId)
store.loadLocal()

// 使用组合式函数
const documents = useProjectDocuments(projectId)
const conversations = useProjectConversations(projectId)
const editorAi = useEditorAi()

// 编辑器区域引用
const editorAreaRef = ref<any>(null)

// 将 editorRef 连接到 editorAi（更有效的双向绑定）
watch(
    () => editorAreaRef.value?.editorRef,
    (ref) => {
        if (ref) {
            editorAi.editorRef.value = ref
        }
    },
    { immediate: true, flush: 'post' }
)

// 直接暴露 editorRef 以便调试
const getEditorRef = () => editorAreaRef.value?.editorRef

// UI状态
const docSearch = ref('')
const showRenameDoc = ref(false)
const renameDocId = ref<string | null>(null)
const renameDocTitle = ref('')
const editingDocTitle = ref(false)
const docTitleDraft = ref('')
const showTagsModal = ref(false)
const editTagsDocId = ref<string | null>(null)
const editTagsValue = ref('')
const showRenameConversation = ref(false)
const renameConversationTitle = ref('')

// 加载项目数据
const { data: projectData } = await useAsyncData('project', () =>
    $fetch(`/api/projects/${projectId}`, { headers: requestHeaders })
)

// 文档列表
const documentList = documents.getDocumentList(docSearch)

// 监听活动文档变化
watch(documents.activeDocument, () => {
    editingDocTitle.value = false
    docTitleDraft.value = ''
    editorAi.closeAiPreview()
})

// 文档标题编辑
function openRenameDocument(docId: string) {
    const doc = store.documents[docId]
    if (!doc) return
    renameDocId.value = docId
    renameDocTitle.value = doc.title
    showRenameDoc.value = true
}

function startInlineRename() {
    if (!documents.activeDocument.value) return
    docTitleDraft.value = documents.activeDocument.value.title || ''
    editingDocTitle.value = true
    nextTick(() => {
        const el = document.getElementById('doc-title-input') as HTMLInputElement | null
        el?.focus()
        el?.select()
    })
}

async function submitInlineRename() {
    if (!documents.activeDocument.value) return
    const title = docTitleDraft.value.trim()
    if (!title) {
        editingDocTitle.value = false
        return
    }
    await documents.updateDocument(String(documents.activeDocument.value.id), { title })
    editingDocTitle.value = false
}

async function submitRenameDocument() {
    if (!renameDocId.value) return
    const title = renameDocTitle.value.trim()
    if (!title) return
    await documents.updateDocument(renameDocId.value, { title })
    showRenameDoc.value = false
}

// 标签编辑
function openEditTags(docId: string) {
    const doc = store.documents[docId]
    if (!doc) return
    editTagsDocId.value = docId
    editTagsValue.value = (doc.tags || []).join(', ')
    showTagsModal.value = true
}

async function submitTags() {
    if (!editTagsDocId.value) return
    const tags = editTagsValue.value
        .split(',')
        .map((tag) => tag.trim())
        .filter(Boolean)
    await documents.updateDocument(editTagsDocId.value, { tags })
    showTagsModal.value = false
}

// 对话管理
function openRenameConversation() {
    if (!store.activeConversationId) return
    const current = store.conversations.find((item) => item.id === store.activeConversationId)
    renameConversationTitle.value = current?.title || '新对话'
    showRenameConversation.value = true
}

async function submitRenameConversation() {
    if (!store.activeConversationId) return
    const title = renameConversationTitle.value.trim()
    if (!title) return
    await conversations.updateConversation(store.activeConversationId, title)
    showRenameConversation.value = false
}
</script>

<template>
    <div class="flex flex-row w-full h-[calc(100vh-var(--ui-header-height))]">
        <!-- 左侧导航栏 -->
        <div class="flex">
            <ProjectSider />
        </div>

        <!-- 中间编辑器区域 -->
        <ProjectEditorArea
            ref="editorAreaRef"
            :project-id="projectId"
            :documents="documentList"
            :active-doc-id="store.activeDocId"
            :document-title="documents.documentTitle.value"
            :document-char-count="documents.documentCharCount.value"
            :document-last-saved-at="documents.documentLastSavedAt.value"
            v-model:content="documents.content.value"
            v-model:search="docSearch"
            :sync-status="documents.syncStatus.value"
            @select-doc="documents.setActiveDoc"
            @create-doc="documents.createDocument('新建文档')"
            @rename-doc="openRenameDocument"
            @edit-tags="openEditTags"
            @delete-doc="documents.deleteDocument"
            @rename-title="editingDocTitle ? submitInlineRename() : startInlineRename()"
            @ai-action="editorAi.handleAiAction"
            @ai-apply="editorAi.applyAiResult"
            @ai-dismiss="editorAi.closeAiPreview"
        />

        <!-- 右侧AI对话区域 -->
        <ProjectChat
            :project-id="projectId"
            :chat="conversations.chat"
            :conversation-options="conversations.conversationOptions.value"
            v-model:active-conversation-id="conversations.activeConversationId.value"
            :conversation-title="conversations.conversationTitle.value"
            @submit="conversations.sendMessage"
            @create="conversations.createConversation()"
            @rename="openRenameConversation"
            @delete="conversations.deleteConversation(store.activeConversationId!)"
        />

        <!-- 重命名文档弹窗 -->
        <UModal v-model:open="showRenameDoc">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">重命名文档</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="showRenameDoc = false" />
                        </div>
                    </template>
                    <UInput v-model="renameDocTitle" size="lg" placeholder="输入文档标题" />
                    <template #footer>
                        <div class="flex items-center justify-end gap-2">
                            <UButton color="neutral" variant="ghost" @click="showRenameDoc = false">取消</UButton>
                            <UButton color="primary" @click="submitRenameDocument">保存</UButton>
                        </div>
                    </template>
                </UCard>
            </template>
        </UModal>

        <!-- 重命名对话弹窗 -->
        <UModal v-model:open="showRenameConversation">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">重命名对话</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="showRenameConversation = false" />
                        </div>
                    </template>
                    <UInput v-model="renameConversationTitle" size="lg" placeholder="输入对话标题" />
                    <template #footer>
                        <div class="flex items-center justify-end gap-2">
                            <UButton color="neutral" variant="ghost" @click="showRenameConversation = false">取消</UButton>
                            <UButton color="primary" @click="submitRenameConversation">保存</UButton>
                        </div>
                    </template>
                </UCard>
            </template>
        </UModal>

        <!-- 编辑标签弹窗 -->
        <UModal v-model:open="showTagsModal">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">编辑标签</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="showTagsModal = false" />
                        </div>
                    </template>
                    <div class="space-y-2">
                        <p class="text-xs text-slate-500 dark:text-slate-400">使用英文逗号分隔标签</p>
                        <UInput v-model="editTagsValue" size="lg" placeholder="如：方法, 数据集, 摘要" />
                    </div>
                    <template #footer>
                        <div class="flex items-center justify-end gap-2">
                            <UButton color="neutral" variant="ghost" @click="showTagsModal = false">取消</UButton>
                            <UButton color="primary" @click="submitTags">保存</UButton>
                        </div>
                    </template>
                </UCard>
            </template>
        </UModal>
    </div>
</template>
