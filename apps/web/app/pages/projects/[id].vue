<script setup lang="ts">
import type { ChatContextItem } from '~/app/composables/useChatContext'
import { useProjectWorkspaceStore } from '#imports'

definePageMeta({
    layout: 'editor',
})

const route = useRoute()
const projectId = String(route.params.id)
const requestHeaders = useRequestHeaders(['cookie'])
const store = useProjectWorkspaceStore()

// 初始化store
store.setProject(projectId)
store.loadLocal()

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
const showErrorMessage = ref('')
const editorHasSelection = ref(false)

// 使用组合式函数
const documents = useProjectDocuments(projectId)
const conversations = useProjectConversations(projectId)
const editorAi = useEditorAi()
const chatContext = useChatContext(projectId)
const skills = useProjectSkills(projectId)
const workflows = useProjectWorkflows(projectId)

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

const documentList = documents.getDocumentList(docSearch)

// 对话管理
async function handleCreateConversation() {
    try {
        const created = await conversations.createConversation()
        if (created?.id) {
            chatContext.ensureConversation(String(created.id))
        }
        return created
    } catch (error: any) {
        showErrorMessage.value = error.message || '创建对话失败'
    }
}

async function handleDeleteConversation(conversationId: string) {
    try {
        await conversations.deleteConversation(Number(conversationId))
        chatContext.removeConversation(String(conversationId))
    } catch (error: any) {
        showErrorMessage.value = error.message || '删除对话失败'
    }
}

async function handleSendMessage(text: string) {
    try {
        const activeId = conversations.activeConversationId.value
        const context = chatContext.getContext(activeId)
        const prefix = buildContextPrefix(context.items, context.skills, context.workflows)
        const finalText = prefix ? `${prefix}\n\n${text}` : text
        const contextView = context.items.map((item) => ({
            id: item.id,
            label: item.label,
        }))
        await conversations.sendMessage(finalText, { context: contextView })
        chatContext.clearTransient(activeId)
    } catch (error: any) {
        showErrorMessage.value = error.message || '发送消息失败'
    }
}

function openRenameDocument(docId: string) {
    const doc = store.documents[docId]
    if (!doc) return
    renameDocId.value = docId
    renameDocTitle.value = doc.title || ''
    showRenameDoc.value = true
}

async function submitRenameDocument() {
    if (!renameDocId.value) return
    const title = renameDocTitle.value.trim()
    if (!title) return
    await documents.updateDocument(renameDocId.value, { title })
    showRenameDoc.value = false
}

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
        .map((item) => item.trim())
        .filter(Boolean)
    await documents.updateDocument(editTagsDocId.value, { tags })
    showTagsModal.value = false
}

function startInlineRename() {
    editingDocTitle.value = true
    docTitleDraft.value = documents.documentTitle.value
}

async function submitInlineRename() {
    if (!editingDocTitle.value) return
    const title = docTitleDraft.value.trim()
    if (title && store.activeDocId) {
        await documents.updateDocument(store.activeDocId, { title })
    }
    editingDocTitle.value = false
    docTitleDraft.value = ''
}

watch(
    () => conversations.activeConversationId.value,
    (value) => {
        if (value) {
            chatContext.ensureConversation(value)
        }
    },
    { immediate: true }
)

const currentContext = computed(() => chatContext.getContext(conversations.activeConversationId.value))
const contextItems = computed(() => currentContext.value.items)
const selectedSkills = computed(() => currentContext.value.skills)
const selectedWorkflows = computed(() => currentContext.value.workflows)

function updateContextItems(items: ChatContextItem[]) {
    chatContext.setContextItems(conversations.activeConversationId.value, items)
}

function updateSelectedSkills(items: string[]) {
    chatContext.setSkills(conversations.activeConversationId.value, items)
}

function updateSelectedWorkflows(items: string[]) {
    chatContext.setWorkflows(conversations.activeConversationId.value, items)
}

async function handleRequestSelectionContext() {
    const selected = editorAreaRef.value?.editorRef?.getSelectedText?.() || ''
    if (!selected.trim()) {
        showErrorMessage.value = '请先在编辑器中选择内容'
        return
    }
    if (!conversations.activeConversationId.value) {
        await handleCreateConversation()
    }
    const label = selected.length > 24 ? `${selected.slice(0, 24)}…` : selected
    chatContext.addItem(conversations.activeConversationId.value, {
        type: 'selection',
        id: `sel-${Date.now()}`,
        label,
        content: selected,
        sourceDocId: store.activeDocId ?? undefined,
    })
}

const defaultKbId = ref<number | null>(null)

async function ensureChatKnowledgeBase() {
    if (defaultKbId.value) return defaultKbId.value
    const list = await $fetch<{ items: any[] }>('/api/knowledge-bases', {
        query: { scope: 'project', project_id: projectId },
        headers: requestHeaders,
    })
    const existing = list.items?.[0]
    if (existing?.id) {
        defaultKbId.value = Number(existing.id)
        return defaultKbId.value
    }
    const created = await $fetch<any>('/api/knowledge-bases', {
        method: 'POST',
        body: {
            scope: 'project',
            project_id: Number(projectId),
            name: '项目对话附件库',
            description: '用于对话附件上传的默认知识库',
        },
        headers: requestHeaders,
    })
    defaultKbId.value = Number(created.id)
    return defaultKbId.value
}

async function handleUploadFile(file: File) {
    try {
        if (!conversations.activeConversationId.value) {
            await handleCreateConversation()
        }
        const kbId = await ensureChatKnowledgeBase()
        const form = new FormData()
        form.append('file', file)
        const uploaded = await $fetch<any>(`/api/knowledge-bases/${kbId}/documents`, {
            method: 'POST',
            body: form,
        })
        chatContext.addItem(conversations.activeConversationId.value, {
            type: 'file',
            id: `file-${Date.now()}`,
            label: uploaded?.file_name || file.name,
            kbId,
            docId: Number(uploaded.id),
            contentType: uploaded?.content_type ?? file.type ?? null,
        })
    } catch (error: any) {
        showErrorMessage.value = error.message || '上传附件失败'
    }
}

function buildContextPrefix(
    items: ChatContextItem[],
    skillIds: string[],
    workflowIds: string[]
) {
    const lines: string[] = []
    items.forEach((item) => {
        if (item.type === 'selection') {
            lines.push(`- 选中内容：${item.content}`)
        } else if (item.type === 'file') {
            lines.push(`- 附件：${item.label} (KB:${item.kbId}/Doc:${item.docId})`)
        }
    })
    skillIds.forEach((id) => {
        const skill = skills.itemMap.value.get(id)
        if (!skill) return
        const prompt = skill.prompt ? `\n  ${skill.prompt}` : ''
        lines.push(`- 技能：${skill.name}${prompt}`)
    })
    workflowIds.forEach((id) => {
        const workflow = workflows.itemMap.value.get(id)
        if (!workflow) return
        let stepsSummary = ''
        if (Array.isArray(workflow.steps)) {
            stepsSummary = workflow.steps.join(' > ')
        } else if (workflow.steps) {
            stepsSummary = JSON.stringify(workflow.steps)
        }
        lines.push(`- 工作流：${workflow.name}${stepsSummary ? `\n  ${stepsSummary}` : ''}`)
    })
    if (!lines.length) return ''
    return `【上下文】\n${lines.join('\n')}`
}
</script>

<template>
    <div class="flex flex-row w-full h-[calc(100vh-var(--ui-header-height))]">
        <!-- 左侧导航栏 -->
        <div class="flex">
            <ProjectSider />
        </div>

        <!-- 中间编辑器区域 -->
        <ProjectEditorArea ref="editorAreaRef" :project-id="projectId" :documents="documentList"
            :active-doc-id="store.activeDocId" :document-title="documents.documentTitle.value"
            :document-char-count="documents.documentCharCount.value"
            :document-last-saved-at="documents.documentLastSavedAt.value" v-model:content="documents.content.value"
            v-model:search="docSearch" :sync-status="documents.syncStatus.value" @select-doc="documents.setActiveDoc"
            @create-doc="documents.createDocument('新建文档')" @rename-doc="openRenameDocument" @edit-tags="openEditTags"
            @delete-doc="documents.deleteDocument"
            @rename-title="editingDocTitle ? submitInlineRename() : startInlineRename()"
            @ai-action="editorAi.handleAiAction" @ai-apply="editorAi.applyAiResult"
            @ai-dismiss="editorAi.closeAiPreview" @selection-change="editorHasSelection = $event" />

        <!-- 右侧AI对话区域 -->
        <ProjectChat :project-id="projectId" :chat="conversations.chat" :error-message="conversations.state"
            :conversation-options="conversations.conversationOptions.value"
            v-model:active-conversation-id="conversations.activeConversationId.value"
            :conversation-title="conversations.conversationTitle.value" :context-items="contextItems"
            :selected-skills="selectedSkills" :selected-workflows="selectedWorkflows" :skill-options="skills.options"
            :workflow-options="workflows.options" :has-selection="editorHasSelection" @submit="handleSendMessage"
            @create="handleCreateConversation" @delete="handleDeleteConversation"
            @request-selection-context="handleRequestSelectionContext" @upload-file="handleUploadFile"
            @update-context-items="updateContextItems" @update-selected-skills="updateSelectedSkills"
            @update-selected-workflows="updateSelectedWorkflows" @settings="showErrorMessage = 'Agent 创建暂未接入'" />

        <!-- 重命名文档弹窗 -->
        <UModal v-model:open="showRenameDoc">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">重命名文档</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs"
                                @click="showRenameDoc = false" />
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

        <!-- 编辑标签弹窗 -->
        <UModal v-model:open="showTagsModal">
            <template #content>
                <UCard>
                    <template #header>
                        <div class="flex items-center justify-between">
                            <h3 class="text-base font-semibold text-slate-900 dark:text-white">编辑标签</h3>
                            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs"
                                @click="showTagsModal = false" />
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

        <!-- 错误提示 -->
        <UAlert v-if="showErrorMessage" :title="'错误'" :description="showErrorMessage" :timeout="5000" close
            color="error" icon="i-lucide-alert-circle" @update:open="showErrorMessage = ''"
            class="absolute top-5 z-100 box-border" />
    </div>
</template>
