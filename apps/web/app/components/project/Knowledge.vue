<script setup lang="ts">
type KnowledgeBaseItem = {
  id: number
  name: string
  scope: string
  description?: string | null
  doc_count?: number
  updated_at?: string
}

const props = defineProps<{
  projectId: string
}>()

const showCreate = ref(false)
const selectedBaseId = ref<string | null>(null)
const file = ref<File | null>(null)
const uploading = ref(false)

const { data: basesData, refresh: refreshBases } = await useAsyncData(
  () => `project-knowledge-${props.projectId}`,
  () =>
    $fetch<{ items: KnowledgeBaseItem[] }>("/api/knowledge-bases", {
      query: { scope: "project", project_id: props.projectId },
    })
)

const bases = computed(() => basesData.value?.items ?? [])

const baseOptions = computed(() =>
  bases.value.map((item) => ({ label: item.name, value: String(item.id) }))
)

watch(bases, (items) => {
  if (!selectedBaseId.value && items.length > 0) {
    selectedBaseId.value = String(items[0].id)
  }
})

async function handleCreated() {
  await refreshBases()
}

watch(file, async (newFile) => {
  if (!newFile || !selectedBaseId.value) return
  uploading.value = true
  try {
    const form = new FormData()
    form.append("file", newFile)
    await $fetch(`/api/knowledge-bases/${selectedBaseId.value}/documents`, {
      method: "POST",
      body: form,
    })
    await refreshBases()
  } finally {
    uploading.value = false
    file.value = null
  }
})
</script>

<template>
  <div class="w-80 p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-semibold text-slate-900 dark:text-white">项目知识库</h3>
      <UButton color="primary" size="xs" @click="showCreate = true">新建</UButton>
    </div>

    <UCard class="border-slate-200/70 dark:border-slate-800">
      <div v-if="bases.length === 0" class="text-xs text-slate-500 dark:text-slate-400">
        暂无知识库，请先创建。
      </div>
      <div v-else class="space-y-3">
        <div
          v-for="base in bases"
          :key="base.id"
          class="rounded-lg border border-slate-200/60 p-3 text-xs text-slate-600 dark:border-slate-800 dark:text-slate-300"
        >
          <div class="flex items-center justify-between">
            <p class="font-semibold text-slate-900 dark:text-white">{{ base.name }}</p>
            <UBadge color="neutral" variant="soft">{{ base.doc_count ?? 0 }} 文档</UBadge>
          </div>
          <p class="mt-2 line-clamp-2">{{ base.description || "暂无描述" }}</p>
          <NuxtLink :to="`/knowledge/${base.id}`" class="mt-2 inline-flex text-primary">
            进入详情
          </NuxtLink>
        </div>
      </div>
    </UCard>

    <UCard class="border-slate-200/70 dark:border-slate-800">
      <div class="space-y-3">
        <p class="text-xs font-semibold text-slate-900 dark:text-white">上传文档</p>
        <USelectMenu
          v-model="selectedBaseId"
          :options="baseOptions"
          placeholder="选择知识库"
          size="sm"
        />
        <UFileUpload
          v-model="file"
          :preview="false"
          label="选择文件"
          description="上传到选中的知识库"
          :disabled="uploading || !selectedBaseId"
        >
          <template #leading>
            <UAvatar :icon="uploading ? 'i-lucide-loader-circle' : 'i-lucide-upload'" size="md" :ui="{ icon: [uploading && 'animate-spin'] }" />
          </template>
        </UFileUpload>
      </div>
    </UCard>

    <KnowledgeCreateKnowledgeModal
      v-model:open="showCreate"
      mode="project"
      :project-id="props.projectId"
      @created="handleCreated"
    />
  </div>
</template>
