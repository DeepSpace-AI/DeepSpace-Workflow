<script setup lang="ts">
useHead({
  title: "知识库详情 - Deepspace Workflow",
})

type KnowledgeBaseItem = {
  id: number
  name: string
  scope: string
  project_id?: number | null
  description?: string | null
  created_at?: string
  updated_at?: string
}

type KnowledgeDocumentItem = {
  id: number
  file_name: string
  content_type?: string | null
  size_bytes?: number | null
  status: string
  created_at?: string
}

const route = useRoute()
const kbId = String(route.params.id)
const file = ref<File | null>(null)
const uploading = ref(false)
const errorMessage = ref("")

const { data: baseData } = await useAsyncData("knowledge-base", () =>
  $fetch<KnowledgeBaseItem>(`/api/knowledge-bases/${kbId}`)
)

const { data: docsData, refresh: refreshDocs } = await useAsyncData("knowledge-docs", () =>
  $fetch<{ items: KnowledgeDocumentItem[] }>(`/api/knowledge-bases/${kbId}/documents`)
)

const base = computed(() => baseData.value)
const docs = computed(() => docsData.value?.items ?? [])

const scopeLabel = (value?: string) => {
  if (value === "org") return "组织级"
  if (value === "project") return "项目级"
  return "未知"
}

const formatSize = (value?: number | null) => {
  if (!value) return "-"
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  if (value < 1024 * 1024 * 1024) return `${(value / (1024 * 1024)).toFixed(1)} MB`
  return `${(value / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

watch(file, async (newFile) => {
  if (!newFile) return
  errorMessage.value = ""
  uploading.value = true

  try {
    const form = new FormData()
    form.append("file", newFile)
    await $fetch(`/api/knowledge-bases/${kbId}/documents`, {
      method: "POST",
      body: form,
    })
    await refreshDocs()
  } catch (error: any) {
    errorMessage.value = error?.message || "上传失败"
  } finally {
    uploading.value = false
    file.value = null
  }
})

async function deleteDocument(docId: number) {
  await $fetch(`/api/knowledge-bases/${kbId}/documents/${docId}`, { method: "DELETE" })
  await refreshDocs()
}
</script>

<template>
  <UContainer>
    <div class="my-10 space-y-8">
      <div class="flex flex-col gap-3">
        <NuxtLink to="/knowledge" class="text-sm text-primary hover:text-primary/80">← 返回知识库</NuxtLink>
        <div class="space-y-1">
          <p class="text-xs font-semibold uppercase tracking-widest text-primary">Knowledge Base</p>
          <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100">
            {{ base?.name || "知识库详情" }}
          </h1>
          <p class="text-sm text-slate-600 dark:text-slate-300">
            {{ base?.description || "暂无描述" }}
          </p>
        </div>
      </div>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="grid gap-4 md:grid-cols-3">
          <div>
            <p class="text-xs text-slate-500 dark:text-slate-400">范围</p>
            <p class="mt-1 text-sm font-semibold text-slate-900 dark:text-white">{{ scopeLabel(base?.scope) }}</p>
          </div>
          <div>
            <p class="text-xs text-slate-500 dark:text-slate-400">最近更新</p>
            <p class="mt-1 text-sm font-semibold text-slate-900 dark:text-white">{{ base?.updated_at || "-" }}</p>
          </div>
          <div>
            <p class="text-xs text-slate-500 dark:text-slate-400">文档数量</p>
            <p class="mt-1 text-sm font-semibold text-slate-900 dark:text-white">{{ docs.length }}</p>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <h2 class="text-base font-semibold text-slate-900 dark:text-white">上传文档</h2>
            <p class="text-xs text-slate-500 dark:text-slate-400">支持任意格式文件上传，仅做管理。</p>
          </div>
          <div class="w-full md:w-80">
            <UFileUpload
              v-model="file"
              :preview="false"
              label="选择文件"
              description="最大 25MB"
              :disabled="uploading"
            >
              <template #leading>
                <UAvatar :icon="uploading ? 'i-lucide-loader-circle' : 'i-lucide-upload'" size="lg" :ui="{ icon: [uploading && 'animate-spin'] }" />
              </template>
            </UFileUpload>
            <p v-if="errorMessage" class="mt-2 text-xs text-red-500">{{ errorMessage }}</p>
          </div>
        </div>
      </UCard>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-base font-semibold text-slate-900 dark:text-white">文档列表</h2>
          </div>
        </template>

        <div v-if="docs.length === 0" class="py-10 text-center text-sm text-slate-500 dark:text-slate-400">
          暂无文档，请上传第一份资料。
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="doc in docs"
            :key="doc.id"
            class="flex flex-col gap-3 rounded-xl border border-slate-200/60 p-4 text-sm text-slate-700 dark:border-slate-800 dark:text-slate-200 md:flex-row md:items-center md:justify-between"
          >
            <div class="space-y-1">
              <p class="font-semibold text-slate-900 dark:text-white">{{ doc.file_name }}</p>
              <p class="text-xs text-slate-500 dark:text-slate-400">
                {{ doc.content_type || "未知格式" }} · {{ formatSize(doc.size_bytes) }} · {{ doc.created_at || "-" }}
              </p>
            </div>
            <div class="flex items-center gap-2">
              <UBadge color="neutral" variant="soft">{{ doc.status }}</UBadge>
              <NuxtLink
                :to="`/api/knowledge-bases/${kbId}/documents/${doc.id}/download?disposition=inline`"
                target="_blank"
                class="text-xs text-primary hover:text-primary/80"
              >
                预览
              </NuxtLink>
              <NuxtLink
                :to="`/api/knowledge-bases/${kbId}/documents/${doc.id}/download`"
                class="text-xs text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white"
              >
                下载
              </NuxtLink>
              <UButton color="red" variant="ghost" size="xs" @click="deleteDocument(doc.id)">删除</UButton>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </UContainer>
</template>
