<script setup lang="ts">
import type { ProjectDocumentState } from '~/app/stores/projectWorkspace'

const props = defineProps<{
  projectId: string
  documents?: ProjectDocumentState[]
  activeDocId: string | null
  search: string
  syncStatus?: {
    saving: boolean
    lastSavedAt?: string
    lastError?: string
  }
}>()

const emit = defineEmits<{
  (e: 'update:search', value: string): void
  (e: 'select-doc', id: string): void
  (e: 'create-doc'): void
  (e: 'rename-doc', id: string): void
  (e: 'edit-tags', id: string): void
  (e: 'delete-doc', id: string): void
}>()

const section = ref<'documents' | 'knowledge'>('documents')
const safeDocuments = computed(() => props.documents ?? [])

const searchModel = computed({
  get: () => props.search,
  set: (value) => emit('update:search', value),
})

const savingLabel = computed(() => {
  if (!props.syncStatus) return ''
  if (props.syncStatus.lastError) return '保存失败'
  if (props.syncStatus.saving) return '正在保存...'
  if (props.syncStatus.lastSavedAt) return '已保存'
  return ''
})

function formatTime(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toISOString().slice(0, 10)
}
</script>

<template>
  <div class="w-80 border-r border-slate-200/70 dark:border-slate-800 flex flex-col">
    <div class="p-4 border-b border-slate-200/70 dark:border-slate-800 space-y-3">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">资源管理</p>
          <p class="text-sm font-semibold text-slate-900 dark:text-white">项目资源</p>
        </div>
        <UBadge v-if="savingLabel" color="neutral" variant="soft">{{ savingLabel }}</UBadge>
      </div>

      <div class="flex items-center gap-2">
        <UButton
          size="xs"
          :color="section === 'documents' ? 'primary' : 'neutral'"
          variant="soft"
          @click="section = 'documents'"
        >
          文档
        </UButton>
        <UButton
          size="xs"
          :color="section === 'knowledge' ? 'primary' : 'neutral'"
          variant="soft"
          @click="section = 'knowledge'"
        >
          知识库
        </UButton>
      </div>

      <template v-if="section === 'documents'">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">项目文档</p>
            <p class="text-[11px] text-slate-400">共 {{ safeDocuments.length }} 篇</p>
          </div>
          <UButton icon="i-lucide-plus" size="xs" color="primary" @click="emit('create-doc')" />
        </div>
        <UInput v-model="searchModel" size="sm" placeholder="搜索文档或标签" />
      </template>
    </div>

    <div class="flex-1 overflow-y-auto p-2 space-y-1" v-if="section === 'documents'">
      <button
        v-for="doc in safeDocuments"
        :key="doc.id"
        @click="emit('select-doc', String(doc.id))"
        class="w-full text-left rounded-lg px-3 py-2 transition"
        :class="[
          activeDocId === String(doc.id)
            ? 'bg-primary/10 text-primary'
            : 'hover:bg-slate-100 dark:hover:bg-slate-800'
        ]"
      >
        <div class="flex items-center justify-between gap-2">
          <p class="text-sm font-medium truncate">{{ doc.title || '未命名文档' }}</p>
          <div class="flex items-center gap-2">
            <span v-if="doc.dirty" class="h-2 w-2 rounded-full bg-amber-400" />
            <UDropdownMenu
              :items="[[ 
                { label: '重命名', icon: 'i-lucide-pencil', onSelect: () => emit('rename-doc', String(doc.id)) },
                { label: '编辑标签', icon: 'i-lucide-tag', onSelect: () => emit('edit-tags', String(doc.id)) },
                { label: '删除', icon: 'i-lucide-trash-2', onSelect: () => emit('delete-doc', String(doc.id)) }
              ]]"
            >
              <UButton icon="i-lucide-more-vertical" size="xs" color="neutral" variant="ghost" />
            </UDropdownMenu>
          </div>
        </div>
        <div v-if="doc.tags?.length" class="mt-2 flex flex-wrap gap-1">
          <UBadge v-for="tag in doc.tags" :key="tag" size="xs" color="neutral" variant="soft">
            {{ tag }}
          </UBadge>
        </div>
        <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">更新于 {{ formatTime(doc.updatedAt) }}</p>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto" v-else>
      <ProjectKnowledge :project-id="projectId" />
    </div>
  </div>
</template>
