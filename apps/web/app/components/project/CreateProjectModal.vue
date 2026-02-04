<script setup lang="ts">
const open = defineModel<boolean>('open', { required: true })

type ProjectTypeOption = {
  id: string
  label: string
  description: string
  icon: string
  disabled?: boolean
}

const typeOptions: ProjectTypeOption[] = [
  {
    id: 'research',
    label: '科研写作',
    description: '适合论文撰写与资料整理',
    icon: 'i-heroicons-document-text'
  },
  {
    id: 'agent',
    label: '多 Agent 协作',
    description: '即将支持',
    icon: 'i-heroicons-rectangle-group',
    disabled: true
  }
]

const selectedType = ref('research')
const createName = ref('')
const createDescription = ref('')
const isSubmitting = ref(false)
const showErrors = ref(false)

const emit = defineEmits<{
  (e: 'created'): void
}>()

const nameError = computed(() => {
  const value = createName.value.trim()
  if (!showErrors.value) return ''
  if (!value) return '请输入项目名称'
  if (value.length < 2) return '名称至少 2 个字'
  if (value.length > 40) return '名称最多 40 个字'
  return ''
})

const descriptionError = computed(() => {
  if (!showErrors.value) return ''
  if (createDescription.value && createDescription.value.length > 200) return '描述最多 200 字'
  return ''
})

const isFormValid = computed(() => !nameError.value && !descriptionError.value)

function resetForm() {
  selectedType.value = 'research'
  createName.value = ''
  createDescription.value = ''
  showErrors.value = false
}

watch(open, (value) => {
  if (!value) resetForm()
})

async function createProject() {
  showErrors.value = true
  if (!isFormValid.value) return
  if (isSubmitting.value) return

  isSubmitting.value = true
  try {
    await $fetch('/api/projects', {
      method: 'POST',
      body: {
        name: createName.value,
        description: createDescription.value || undefined,
        type: selectedType.value
      }
    })
    resetForm()
    open.value = false
    emit('created')
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
                新建项目
              </p>
              <h3 class="mt-2 text-lg font-semibold text-slate-900 dark:text-white">
                创建写作空间
              </h3>
            </div>
            <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="open = false" />
          </div>
        </template>

        <div class="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div class="space-y-5">
            <div class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">项目类型</p>
              <div class="grid gap-3 sm:grid-cols-2">
                <button
                  v-for="option in typeOptions"
                  :key="option.id"
                  type="button"
                  :disabled="option.disabled"
                  :aria-pressed="selectedType === option.id"
                  @click="selectedType = option.id"
                  class="text-left"
                >
                  <div
                    :class="[
                      'rounded-xl border p-4 transition',
                      option.disabled
                        ? 'border-slate-200/70 bg-white/50 opacity-70 dark:border-slate-800 dark:bg-slate-900/40'
                        : selectedType === option.id
                          ? 'border-primary/70 bg-white shadow-sm dark:border-primary/50 dark:bg-slate-900/70'
                          : 'border-slate-200/70 bg-white/70 hover:border-primary/60 dark:border-slate-800 dark:bg-slate-900/60'
                    ]"
                  >
                    <div class="flex items-start gap-3">
                      <div
                        :class="[
                          'flex h-10 w-10 items-center justify-center rounded-lg',
                          selectedType === option.id
                            ? 'bg-primary/10 text-primary'
                            : 'bg-slate-200/60 text-slate-500 dark:bg-slate-800 dark:text-slate-400'
                        ]"
                      >
                        <UIcon :name="option.icon" class="h-5 w-5" />
                      </div>
                      <div>
                        <p
                          :class="[
                            'text-sm font-semibold',
                            option.disabled
                              ? 'text-slate-600 dark:text-slate-300'
                              : 'text-slate-900 dark:text-white'
                          ]"
                        >
                          {{ option.label }}
                        </p>
                        <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                          {{ option.description }}
                        </p>
                      </div>
                    </div>
                  </div>
                </button>
              </div>
            </div>

            <div class="space-y-2">
              <div class="flex items-center justify-between text-xs font-semibold text-slate-500 dark:text-slate-400">
                <span>应用名称 & 图标</span>
                <span class="text-[10px] font-medium text-slate-400">必填</span>
              </div>
              <div
                class="flex items-center gap-3 rounded-xl border border-slate-200/70 bg-white/70 p-3 dark:border-slate-800 dark:bg-slate-900/60"
              >
                <div
                  class="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 text-amber-600 dark:bg-amber-500/10 dark:text-amber-300"
                >
                  <UIcon name="i-heroicons-sparkles" class="h-6 w-6" />
                </div>
                <UInput
                  v-model="createName"
                  placeholder="给你的项目起个名字"
                  size="lg"
                  class="flex-1"
                  maxlength="40"
                  @keyup.enter="createProject"
                />
              </div>
              <p v-if="nameError" class="text-xs text-rose-500">{{ nameError }}</p>
            </div>

            <div class="space-y-2">
              <p class="text-xs font-semibold text-slate-500 dark:text-slate-400">描述（可选）</p>
              <UTextarea
                v-model="createDescription"
                placeholder="简要描述目标、研究方向或阶段"
                :rows="4"
                size="lg"
                class="w-full"
                maxlength="200"
              />
              <p v-if="descriptionError" class="text-xs text-rose-500">{{ descriptionError }}</p>
            </div>
          </div>

          <div
            class="rounded-2xl border border-slate-200/70 bg-linear-to-b from-slate-50 to-white p-5 text-sm text-slate-600 dark:border-slate-800 dark:from-slate-900/80 dark:to-slate-900/30 dark:text-slate-300"
          >
            <div class="space-y-4">
              <div>
                <p class="text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">
                  建议
                </p>
                <h4 class="mt-2 text-base font-semibold text-slate-900 dark:text-white">
                  创建后你可以
                </h4>
                <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
                  组织资料、启动 AI 对话并持续跟踪进度。
                </p>
              </div>
              <div class="space-y-3 text-xs">
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-folder-open" class="h-4 w-4 text-primary" />
                  <span>集中管理文献、数据与稿件</span>
                </div>
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-chat-bubble-left-right" class="h-4 w-4 text-emerald-500" />
                  <span>与 AI 协作完成摘要和段落</span>
                </div>
                <div class="flex items-center gap-2">
                  <UIcon name="i-heroicons-clock" class="h-4 w-4 text-amber-500" />
                  <span>追踪最近修改与里程碑</span>
                </div>
              </div>
              <div
                class="rounded-xl border border-dashed border-slate-200/80 bg-white/80 p-3 text-xs text-slate-500 dark:border-slate-700 dark:bg-slate-900/60 dark:text-slate-400"
              >
                没有灵感？稍后可从模板库创建项目。
              </div>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p class="text-xs text-slate-500 dark:text-slate-400">
              创建后可在项目详情中进一步配置。
            </p>
            <div class="flex items-center gap-2">
              <UButton color="neutral" variant="ghost" @click="open = false">取消</UButton>
              <UButton
                color="primary"
                :disabled="!isFormValid || isSubmitting"
                :loading="isSubmitting"
                @click="createProject"
              >
                创建项目
              </UButton>
            </div>
          </div>
        </template>
      </UCard>
    </template>
  </UModal>
</template>
