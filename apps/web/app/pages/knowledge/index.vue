<script setup lang="ts">
type KnowledgeBaseItem = {
  id: number
  name: string
  scope: string
  project_id?: number | null
  description?: string | null
  doc_count?: number
  updated_at?: string
}

type ProjectItem = {
  id: number
  name: string
}

const scope = ref("all")
const projectId = ref<string | null>(null)
const showCreate = ref(false)
const requestHeaders = useRequestHeaders(["cookie"])

const { data: projectsData } = await useAsyncData("knowledge-projects", () =>
  $fetch<{ items: ProjectItem[] }>("/api/projects", { headers: requestHeaders })
)

const projectOptions = computed(() => {
  const items = projectsData.value?.items ?? []
  return items.map((item) => ({ label: item.name, value: String(item.id) }))
})

const { data: basesData, refresh: refreshBases } = await useAsyncData(
  "knowledge-bases",
  () => {
    const query: Record<string, string> = { scope: scope.value }
    if (projectId.value) {
      query.project_id = projectId.value
    }
    return $fetch<{ items: KnowledgeBaseItem[] }>("/api/knowledge-bases", { query, headers: requestHeaders })
  },
  { watch: [scope, projectId] }
)

const bases = computed(() => basesData.value?.items ?? [])

const scopeLabel = (value: string) => {
  if (value === "org") return "组织级"
  if (value === "project") return "项目级"
  return "未知"
}

async function handleCreated() {
  await refreshBases()
}
</script>

<template>
  <UContainer>
    <div class="my-10 space-y-8">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-primary">Knowledge Base</p>
          <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100 md:text-4xl">
            知识库管理
          </h1>
          <p class="max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">
            统一管理组织级与项目级知识库文档，支持快速上传与追踪。
          </p>
        </div>
        <div class="flex w-full flex-col gap-3 sm:flex-row md:w-auto">
          <UButton color="primary" size="lg" class="w-full sm:w-auto" @click="showCreate = true">
            新建知识库
          </UButton>
        </div>
      </div>

      <UCard class="border-slate-200/70 dark:border-slate-800">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div class="flex w-full flex-col gap-4 sm:flex-row">
            <USelectMenu
              v-model="scope"
              :options="[
                { label: '全部范围', value: 'all' },
                { label: '组织级', value: 'org' },
                { label: '项目级', value: 'project' }
              ]"
              size="lg"
              class="min-w-[160px]"
            />
            <USelectMenu
              v-model="projectId"
              :options="projectOptions"
              placeholder="选择项目(可选)"
              size="lg"
              class="min-w-[200px]"
              clearable
            />
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
            <span>知识库数量</span>
            <UBadge color="primary" variant="soft">{{ bases.length }}</UBadge>
          </div>
        </div>
      </UCard>

      <div class="grid gap-6 lg:grid-cols-3">
        <UCard
          v-for="base in bases"
          :key="base.id"
          class="flex h-full flex-col border-slate-200/70 transition-shadow duration-200 hover:shadow-lg dark:border-slate-800"
        >
          <template #header>
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-xs font-medium text-slate-500 dark:text-slate-400">
                  {{ scopeLabel(base.scope) }}
                </p>
                <h2 class="mt-1 text-lg font-semibold text-slate-900 dark:text-white">
                  {{ base.name }}
                </h2>
              </div>
              <UBadge color="neutral" variant="soft">
                {{ base.doc_count ?? 0 }} 文档
              </UBadge>
            </div>
          </template>

          <div class="flex flex-1 flex-col gap-4">
            <p class="text-sm leading-6 text-slate-600 dark:text-slate-300">
              {{ base.description || "暂无描述" }}
            </p>
            <div
              class="rounded-xl bg-slate-50 p-3 text-xs text-slate-600 dark:bg-slate-900/60 dark:text-slate-300"
            >
              <div class="flex items-center justify-between">
                <span>最近更新</span>
                <span class="font-semibold text-slate-900 dark:text-white">
                  {{ base.updated_at || "-" }}
                </span>
              </div>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
              <NuxtLink
                :to="`/knowledge/${base.id}`"
                class="font-medium text-primary transition-colors hover:text-primary/80"
              >
                查看详情
              </NuxtLink>
            </div>
          </template>
        </UCard>
      </div>

      <UCard v-if="bases.length === 0" class="border-slate-200/70 dark:border-slate-800">
        <div class="flex flex-col items-center gap-3 py-10 text-center">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300"
          >
            <UIcon name="i-heroicons-magnifying-glass" class="h-6 w-6" />
          </div>
          <div>
            <h3 class="text-base font-semibold text-slate-900 dark:text-white">暂无知识库</h3>
            <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">创建一个知识库开始管理文档。</p>
          </div>
          <UButton color="neutral" variant="soft" size="sm" @click="showCreate = true">
            新建知识库
          </UButton>
        </div>
      </UCard>
    </div>

    <KnowledgeCreateKnowledgeModal
      v-model:open="showCreate"
      :project-options="projectOptions"
      @created="handleCreated"
    />
  </UContainer>
</template>
