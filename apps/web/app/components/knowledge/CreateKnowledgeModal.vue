<script setup lang="ts">
type ProjectOption = {
  label: string
  value: string
}

const open = defineModel<boolean>("open", { required: true })

const props = withDefaults(
  defineProps<{
    mode?: "global" | "project"
    projectOptions?: ProjectOption[]
    projectId?: string | null
  }>(),
  {
    mode: "global",
    projectOptions: () => [],
    projectId: null,
  }
)

const createName = ref("")
const createDescription = ref("")
const createScope = ref<"org" | "project">("org")
const createProjectId = ref<string | null>(null)
const isSubmitting = ref(false)

const emit = defineEmits<{
  (e: "created"): void
}>()

const isProjectMode = computed(() => props.mode === "project")

watch(
  () => open.value,
  (value) => {
    if (!value) return
    createName.value = ""
    createDescription.value = ""
    if (isProjectMode.value) {
      createScope.value = "project"
      createProjectId.value = props.projectId ?? null
    } else {
      createScope.value = "org"
      createProjectId.value = null
    }
  }
)

const canSubmit = computed(() => {
  if (!createName.value.trim()) return false
  if (createScope.value === "project" && !createProjectId.value) return false
  return true
})

async function createKnowledgeBase() {
  if (!canSubmit.value || isSubmitting.value) return

  isSubmitting.value = true
  try {
    const payload: Record<string, any> = {
      name: createName.value,
      description: createDescription.value || undefined,
      scope: createScope.value,
    }
    if (createScope.value === "project") {
      payload.project_id = Number(createProjectId.value)
    }

    await $fetch("/api/knowledge-bases", {
      method: "POST",
      body: payload,
    })

    open.value = false
    emit("created")
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <UModal v-model:open="open" class="min-w-6xl">
    <template #content>
      <UCard class="w-full overflow-hidden">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">
                新建知识库
              </p>
              <h3 class="mt-2 text-lg font-semibold text-slate-900 dark:text-white">
                创建知识库空间
              </h3>
            </div>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="open = false" />
          </div>
        </template>

        <div class="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div class="space-y-5">
            <div class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">知识库范围</p>
              <div class="grid gap-3 sm:grid-cols-2">
                <button
                  type="button"
                  class="rounded-xl border border-slate-200/70 bg-white/70 p-4 text-left transition hover:border-primary/60 dark:border-slate-800 dark:bg-slate-900/60"
                  :class="[
                    createScope === 'org' ? 'border-primary/60 ring-1 ring-primary/30' : 'opacity-80',
                    isProjectMode ? 'cursor-not-allowed opacity-60' : ''
                  ]"
                  :disabled="isProjectMode"
                  @click="createScope = 'org'"
                >
                  <div class="flex items-start gap-3">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                      <UIcon name="i-heroicons-building-library" class="h-5 w-5" />
                    </div>
                    <div>
                      <p class="text-sm font-semibold text-slate-900 dark:text-white">组织级</p>
                      <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                        适合团队共享的公共资料库
                      </p>
                    </div>
                  </div>
                </button>
                <button
                  type="button"
                  class="rounded-xl border border-slate-200/70 bg-white/70 p-4 text-left transition hover:border-primary/60 dark:border-slate-800 dark:bg-slate-900/60"
                  :class="createScope === 'project' ? 'border-primary/60 ring-1 ring-primary/30' : 'opacity-80'"
                  @click="createScope = 'project'"
                >
                  <div class="flex items-start gap-3">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-500">
                      <UIcon name="i-heroicons-folder-open" class="h-5 w-5" />
                    </div>
                    <div>
                      <p class="text-sm font-semibold text-slate-900 dark:text-white">项目级</p>
                      <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                        与项目进度强绑定的资料库
                      </p>
                    </div>
                  </div>
                </button>
              </div>
            </div>

            <div class="space-y-2">
              <div class="flex items-center justify-between text-xs font-semibold text-slate-500 dark:text-slate-400">
                <span>知识库名称</span>
                <span class="text-[10px] font-medium text-slate-400">必填</span>
              </div>
              <div
                class="flex items-center gap-3 rounded-xl border border-slate-200/70 bg-white/70 p-3 dark:border-slate-800 dark:bg-slate-900/60"
              >
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-xl bg-sky-100 text-sky-600 dark:bg-sky-500/10 dark:text-sky-300"
                >
                  <UIcon name="i-heroicons-book-open" class="h-6 w-6" />
                </div>
                <UInput v-model="createName" placeholder="给知识库起个名字" size="lg" class="flex-1" />
              </div>
            </div>

            <div class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">描述（可选）</p>
              <UTextarea
                v-model="createDescription"
                placeholder="说明资料范围、主题或使用场景"
                :rows="4"
                size="lg"
                class="w-full"
              />
            </div>

            <div v-if="createScope === 'project' && !isProjectMode" class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">绑定项目</p>
              <USelectMenu
                v-model="createProjectId"
                :options="projectOptions"
                placeholder="选择项目"
                size="lg"
              />
            </div>
            <div v-else-if="createScope === 'project' && isProjectMode" class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">绑定项目</p>
              <div
                class="rounded-xl border border-dashed border-slate-200/80 bg-white/80 p-3 text-xs text-slate-500 dark:border-slate-700 dark:bg-slate-900/60 dark:text-slate-400"
              >
                当前项目已自动绑定到该知识库。
              </div>
            </div>
          </div>

          <div
            class="rounded-2xl border border-slate-200/70 bg-linear-to-b from-slate-50 to-white p-5 text-sm text-slate-600 dark:border-slate-800 dark:from-slate-900/80 dark:to-slate-900/30 dark:text-slate-300"
          >
            <div class="space-y-4">
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">
                  指南
                </p>
                <h4 class="mt-2 text-base font-semibold text-slate-900 dark:text-white">
                  知识库可以帮你
                </h4>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                  汇总文献、数据、代码与策略文件，方便后续检索与引用。
                </p>
              </div>
              <div class="space-y-3 text-xs">
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-cloud-arrow-up" class="h-4 w-4 text-primary" />
                  <span>快速上传论文、数据与实验记录</span>
                </div>
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-squares-plus" class="h-4 w-4 text-emerald-500" />
                  <span>统一来源，便于 AI 引用和追踪</span>
                </div>
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-clock" class="h-4 w-4 text-amber-500" />
                  <span>记录上传时间，方便归档整理</span>
                </div>
              </div>
              <div
                class="rounded-xl border border-dashed border-slate-200/80 bg-white/80 p-3 text-xs text-slate-500 dark:border-slate-700 dark:bg-slate-900/60 dark:text-slate-400"
              >
                后续可在知识库内开启 RAG 检索能力。
              </div>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p class="text-xs text-slate-500 dark:text-slate-400">
              创建后可继续上传文档并设置访问范围。
            </p>
            <div class="flex items-center gap-2">
              <UButton color="neutral" variant="ghost" @click="open = false">取消</UButton>
              <UButton color="primary" :disabled="!canSubmit || isSubmitting" :loading="isSubmitting" @click="createKnowledgeBase">
                创建知识库
              </UButton>
            </div>
          </div>
        </template>
      </UCard>
    </template>
  </UModal>
</template>
