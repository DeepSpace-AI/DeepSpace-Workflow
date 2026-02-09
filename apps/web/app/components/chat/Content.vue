<script setup lang="ts">
import type { ButtonProps } from '@nuxt/ui'
import type { ChatStatus, UIMessage } from 'ai'
import type { ConversationParamState } from '~/stores/standaloneChat'

type Props = {
  chatId: string | null
  messages: UIMessage[]
  status?: ChatStatus
  error?: string | null
  model: string
  models: string[]
  imageModelOptions: string[]
  drawingLoading?: boolean
  sessionConfig: ConversationParamState
}

type Emits = {
  (e: 'submit', text: string): void
  (e: 'draw', prompt: string): void
  (e: 'update:model', value: string): void
  (e: 'update:sessionConfig', patch: Partial<ConversationParamState>): void
}

const props = withDefaults(defineProps<Props>(), {
  status: 'ready',
  error: null,
  drawingLoading: false,
})
const emit = defineEmits<Emits>()

const input = ref('')
const showDrawModal = ref(false)
const drawPrompt = ref('')
const reasoningEffortItems = [
  { label: '低 (low)', value: 'low' },
  { label: '中 (medium)', value: 'medium' },
  { label: '高 (high)', value: 'high' },
]
const imageSizeItems = [
  { label: '1024 x 1024', value: '1024x1024' },
  { label: '1024 x 1792', value: '1024x1792' },
  { label: '1792 x 1024', value: '1792x1024' },
]
const imageQualityItems = [
  { label: 'standard', value: 'standard' },
  { label: 'hd', value: 'hd' },
]
const imageStyleItems = [
  { label: 'vivid', value: 'vivid' },
  { label: 'natural', value: 'natural' },
]
const canSubmitDraw = computed(() => {
  return drawPrompt.value.trim().length > 0 && !!props.sessionConfig.imageModel && !props.drawingLoading
})

const toast = useToast()

type MessageAction = Omit<ButtonProps, 'onClick'> & {
  onClick?: (event: MouseEvent, message: UIMessage) => void
}
const messageActions = ref<MessageAction[]>([
  {
    icon: 'i-lucide-copy',
    label: '复制',
    onClick: (_event: MouseEvent, message: UIMessage) => {
      const textParts = message.parts.filter((part) => part.type === 'text') as { type: 'text'; text: string }[]
      const fullText = textParts.map((part) => part.text).join('\n')
      navigator.clipboard.writeText(fullText)
      toast.add({
        title: '复制成功',
        color: 'success'
      })
    },
  },
])

function onSubmit(e: Event) {
  e.preventDefault()
  const text = input.value.trim()
  if (!text) return
  emit('submit', text)
  input.value = ''
}

function patchConfig(patch: Partial<ConversationParamState>) {
  emit('update:sessionConfig', patch)
}

function normalizeSliderValue(value: number | number[] | undefined) {
  return Array.isArray(value) ? Number(value[0] ?? 0) : Number(value)
}

function updateTokenBudget(value: number | number[] | undefined) {
  patchConfig({ maxTokenBudget: normalizeSliderValue(value) })
}

function updateTemperature(value: number | number[] | undefined) {
  patchConfig({ temperature: Number(normalizeSliderValue(value).toFixed(2)) })
}

function updateTopP(value: number | number[] | undefined) {
  patchConfig({ topP: Number(normalizeSliderValue(value).toFixed(2)) })
}

function updateReasoningEffort(value: string | undefined) {
  if (value !== 'low' && value !== 'medium' && value !== 'high') return
  patchConfig({ reasoningEffort: value })
}

function updateImageModel(value: string | undefined) {
  patchConfig({ imageModel: String(value || '') })
}

function updateImageSize(value: string | undefined) {
  if (value !== '1024x1024' && value !== '1024x1792' && value !== '1792x1024') return
  patchConfig({ imageSize: value })
}

function updateImageQuality(value: string | undefined) {
  if (value !== 'standard' && value !== 'hd') return
  patchConfig({ imageQuality: value })
}

function updateImageStyle(value: string | undefined) {
  if (value !== 'vivid' && value !== 'natural') return
  patchConfig({ imageStyle: value })
}

function openDrawModal() {
  showDrawModal.value = true
}

function closeDrawModal() {
  showDrawModal.value = false
}

function submitDraw() {
  const prompt = drawPrompt.value.trim()
  if (!prompt || !props.sessionConfig.imageModel || props.drawingLoading) return
  emit('draw', prompt)
  drawPrompt.value = ''
  showDrawModal.value = false
}
</script>

<template>
  <div class="flex-1 min-h-0 flex flex-col overflow-hidden bg-default">
    <!-- 初始 -->
    <div v-if="!props.chatId" class="mx-auto min-w-3xl px-6 py-6">
      <div class="my-12 text-center">
        <h1 class="text-2xl font-black">与你一起创建无限可能</h1>
      </div>

      <UChatPrompt v-model="input" placeholder="从任何想法开始... 按 Enter 发送" variant="soft" :rows="4" size="xl"
        @submit="onSubmit">
        <UChatPromptSubmit icon="i-lucide:send" size="lg" :status="props.status" />
        <template #footer>
          <div class="flex gap-x-1 items-center">
            <USelectMenu placeholder="请选择模型" :model-value="props.model" :items="props.models"
              :ui="{ content: 'min-w-fit' }" @update:model-value="emit('update:model', String($event || ''))" />
            <UTooltip text="绘画">
              <UButton icon="i-lucide-image" variant="ghost" color="neutral" size="sm" @click="openDrawModal" />
            </UTooltip>
            <UPopover>
              <UTooltip text="对话参数">
                <UButton icon="i-lucide:sliders-vertical" variant="ghost" color="neutral" size="sm" />
              </UTooltip>
              <template #content>
                <div class="p-3 min-w-72 max-w-lg space-y-6">
                  <div class="flex justify-between items-center gap-x-3">
                    <div>
                      <h1 class="text-sm font-bold">开启上下文缓存</h1>
                      <p class="text-xs text-muted max-w-lg">最多可减少90%的单次对话生成成本，并将速度提升至最多4倍。</p>
                    </div>
                    <USwitch :model-value="props.sessionConfig.enableChatCache"
                      @update:model-value="patchConfig({ enableChatCache: Boolean($event) })" />
                  </div>

                  <div class="flex justify-between items-center gap-x-3">
                    <div>
                      <h1 class="text-sm font-bold">开启深度思考</h1>
                      <p class="text-xs text-muted max-w-lg">开启后展示高级采样参数，便于微调推理稳定性与发散程度。</p>
                    </div>
                    <USwitch :model-value="props.sessionConfig.thinkingDepth"
                      @update:model-value="patchConfig({ thinkingDepth: Boolean($event) })" />
                  </div>

                  <div v-if="props.sessionConfig.thinkingDepth" class="flex flex-col gap-y-2">
                    <div>
                      <h1 class="text-sm font-bold">推理强度</h1>
                      <p class="text-xs text-muted">控制推理深度档位（low / medium / high）。</p>
                    </div>
                    <USelectMenu :items="reasoningEffortItems" value-key="value" label-key="label"
                      :model-value="props.sessionConfig.reasoningEffort"
                      @update:model-value="updateReasoningEffort(String($event || ''))" />
                  </div>

                  <div class="flex flex-col gap-y-3">
                    <div>
                      <h1 class="text-sm font-bold">最大输出 Token</h1>
                      <p class="text-xs text-muted max-w-lg">控制单次回复的最大长度预算。</p>
                    </div>
                    <USlider :step="1024" :model-value="props.sessionConfig.maxTokenBudget" :min="1024" :max="64000"
                      tooltip @update:model-value="updateTokenBudget" />
                    <div class="text-xs text-muted">当前: {{ props.sessionConfig.maxTokenBudget }}</div>
                  </div>
                  <div class="flex flex-col gap-y-3">
                    <div>
                      <h1 class="text-sm font-bold">Temperature</h1>
                      <p class="text-xs text-muted">值越低越稳定，值越高越发散（0 - 2）。</p>
                    </div>
                    <USlider :step="0.05" :model-value="props.sessionConfig.temperature" :min="0" :max="2" tooltip
                      @update:model-value="updateTemperature" />
                    <div class="text-xs text-muted">当前: {{ props.sessionConfig.temperature.toFixed(2) }}</div>
                  </div>

                  <div class="flex flex-col gap-y-3">
                    <div>
                      <h1 class="text-sm font-bold">Top P</h1>
                      <p class="text-xs text-muted">限制候选词累计概率阈值（0 - 1）。</p>
                    </div>
                    <USlider :step="0.01" :model-value="props.sessionConfig.topP" :min="0" :max="1" tooltip
                      @update:model-value="updateTopP" />
                    <div class="text-xs text-muted">当前: {{ props.sessionConfig.topP.toFixed(2) }}</div>
                  </div>
                </div>
              </template>
            </UPopover>
          </div>
        </template>
      </UChatPrompt>

      <div class="flex justify-center items-center gap-x-4 mt-6">
        <UButton color="neutral" variant="soft" icon="i-lucide:bot" class="rounded-full">创建助理</UButton>
        <UButton color="neutral" variant="soft" icon="i-lucide:messages-square" class="rounded-full">创建群聊</UButton>
        <UButton color="neutral" variant="soft" icon="i-lucide:pen-line" class="rounded-full">开始写作</UButton>
      </div>
    </div>

    <div v-else class="flex-1 min-h-0 flex flex-col w-full px-4 pt-4">
      <div
        class="chat-scrollbar flex-1 min-h-0 overflow-y-auto rounded-t-2xl border border-neutral-200 dark:border-neutral-800 bg-neutral-50 dark:bg-neutral-950/60">
        <UChatMessages class="p-6" :assistant="{
          side: 'left',
          variant: 'naked',
          actions: messageActions,
        }" :messages="props.messages" :status="props.status">
          <template #content="{ message }">
            <template v-for="(part, index) in message.parts" :key="`${message.id}-${part.type}-${index}`">
              <MDC v-if="part.type === 'text' && message.role === 'assistant'" :value="part.text"
                :cache-key="`${message.id}-${index}`" class="chat-message-typography *:first:mt-0 *:last:mb-0" />
              <p v-else-if="part.type === 'text'" class="chat-message-typography whitespace-pre-wrap">{{ part.text }}
              </p>
            </template>
          </template>
        </UChatMessages>
      </div>

      <div
        class="shrink-0 px-6 pt-3 pb-2 rounded-b-2xl border-x border-b border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900">
        <UChatPrompt v-model="input" placeholder="有问题尽管问" variant="soft" @submit="onSubmit"
          :error="props.error ? new Error(props.error) : undefined">
          <UChatPromptSubmit icon="i-lucide:send" size="sm" :status="props.status" />
          <template #footer>
            <div class="flex gap-x-1 items-center">
              <USelectMenu placeholder="请选择模型" :model-value="props.model" :items="props.models"
                :ui="{ content: 'min-w-fit' }" @update:model-value="emit('update:model', String($event || ''))" />
              <UTooltip text="绘画">
                <UButton icon="i-lucide-image" variant="ghost" color="neutral" size="sm" @click="openDrawModal" />
              </UTooltip>
              <UPopover>
                <UTooltip text="对话参数">
                  <UButton icon="i-lucide:sliders-vertical" variant="ghost" color="neutral" size="sm" />
                </UTooltip>
                <template #content>
                  <div class="p-3 min-w-72 max-w-lg space-y-6">
                    <div class="flex justify-between items-center gap-x-3">
                      <div>
                        <h1 class="text-sm font-bold">开启上下文缓存</h1>
                        <p class="text-xs text-muted max-w-lg">最多可减少90%的单次对话生成成本，并将速度提升至最多4倍。</p>
                      </div>
                      <USwitch :model-value="props.sessionConfig.enableChatCache"
                        @update:model-value="patchConfig({ enableChatCache: Boolean($event) })" />
                    </div>

                    <div class="flex justify-between items-center gap-x-3">
                      <div>
                        <h1 class="text-sm font-bold">开启深度思考</h1>
                        <p class="text-xs text-muted max-w-lg">开启后展示高级采样参数，便于微调推理稳定性与发散程度。</p>
                      </div>
                      <USwitch :model-value="props.sessionConfig.thinkingDepth"
                        @update:model-value="patchConfig({ thinkingDepth: Boolean($event) })" />
                    </div>

                    <div v-if="props.sessionConfig.thinkingDepth" class="flex flex-col gap-y-2">
                      <div>
                        <h1 class="text-sm font-bold">推理强度</h1>
                        <p class="text-xs text-muted">控制推理深度档位（low / medium / high）。</p>
                      </div>
                      <USelectMenu :items="reasoningEffortItems" value-key="value" label-key="label"
                        :model-value="props.sessionConfig.reasoningEffort"
                        @update:model-value="updateReasoningEffort(String($event || ''))" />
                    </div>

                    <div class="flex flex-col gap-y-3">
                      <div>
                        <h1 class="text-sm font-bold">最大输出 Token</h1>
                        <p class="text-xs text-muted max-w-lg">控制单次回复的最大长度预算。</p>
                      </div>
                      <USlider :step="1024" :model-value="props.sessionConfig.maxTokenBudget" :min="1024" :max="64000"
                        tooltip @update:model-value="updateTokenBudget" />
                      <div class="text-xs text-muted">当前: {{ props.sessionConfig.maxTokenBudget }}</div>
                    </div>

                    <div class="flex flex-col gap-y-3">
                      <div>
                        <h1 class="text-sm font-bold">Temperature</h1>
                        <p class="text-xs text-muted">值越低越稳定，值越高越发散（0 - 2）。</p>
                      </div>
                      <USlider :step="0.05" :model-value="props.sessionConfig.temperature" :min="0" :max="2" tooltip
                        @update:model-value="updateTemperature" />
                      <div class="text-xs text-muted">当前: {{ props.sessionConfig.temperature.toFixed(2) }}</div>
                    </div>

                    <div class="flex flex-col gap-y-3">
                      <div>
                        <h1 class="text-sm font-bold">Top P</h1>
                        <p class="text-xs text-muted">限制候选词累计概率阈值（0 - 1）。</p>
                      </div>
                      <USlider :step="0.01" :model-value="props.sessionConfig.topP" :min="0" :max="1" tooltip
                        @update:model-value="updateTopP" />
                      <div class="text-xs text-muted">当前: {{ props.sessionConfig.topP.toFixed(2) }}</div>
                    </div>
                  </div>
                </template>
              </UPopover>
            </div>
          </template>
        </UChatPrompt>
      </div>

      <div class="shrink-0 mt-2 px-6 pb-4 text-xs text-muted text-center">AI 不是万能 AI 也可能会犯错。请核查重要信息。</div>
    </div>

    <UModal v-model:open="showDrawModal">
      <template #content>
        <UCard class="w-full max-w-xl">
          <template #header>
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-base font-semibold text-slate-900 dark:text-slate-100">绘画</h3>
                <p class="text-xs text-slate-500 dark:text-slate-400">填写提示词并选择图片参数</p>
              </div>
              <UButton icon="i-lucide-x" color="neutral" variant="ghost" size="xs" @click="closeDrawModal" />
            </div>
          </template>

          <div class="space-y-4">
            <UFormField label="提示词" required>
              <UTextarea v-model="drawPrompt" :rows="4" placeholder="描述你想生成的画面" class="w-full" />
            </UFormField>

            <div class="grid gap-4 sm:grid-cols-2">
              <UFormField label="图片模型" required>
                <USelectMenu placeholder="请选择图片模型" :items="props.imageModelOptions"
                  :model-value="props.sessionConfig.imageModel"
                  @update:model-value="updateImageModel(String($event || ''))" />
              </UFormField>
              <UFormField label="尺寸">
                <USelectMenu :items="imageSizeItems" value-key="value" label-key="label"
                  :model-value="props.sessionConfig.imageSize"
                  @update:model-value="updateImageSize(String($event || ''))" />
              </UFormField>
              <UFormField label="质量">
                <USelectMenu :items="imageQualityItems" value-key="value" label-key="label"
                  :model-value="props.sessionConfig.imageQuality"
                  @update:model-value="updateImageQuality(String($event || ''))" />
              </UFormField>
              <UFormField label="风格">
                <USelectMenu :items="imageStyleItems" value-key="value" label-key="label"
                  :model-value="props.sessionConfig.imageStyle"
                  @update:model-value="updateImageStyle(String($event || ''))" />
              </UFormField>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-end gap-2">
              <UButton color="neutral" variant="ghost" @click="closeDrawModal">取消</UButton>
              <UButton color="primary" :loading="props.drawingLoading" :disabled="!canSubmitDraw" @click="submitDraw">
                生成
              </UButton>
            </div>
          </template>
        </UCard>
      </template>
    </UModal>
  </div>
</template>
