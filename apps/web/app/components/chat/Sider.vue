<script setup lang="ts">
type ConversationOption = {
  label: string
  value: string
  updated_at?: string
}

type Props = {
  conversations: ConversationOption[]
  activeConversationId: string | null
  loading?: boolean
}

type Emits = {
  (e: 'create'): void
  (e: 'select', conversationId: string): void
  (e: 'delete', conversationId: string): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
})
const emit = defineEmits<Emits>()

const showSearchModal = ref(false)
const searchKeyword = ref('')
const searching = ref(false)

type SearchResult = {
  value: string
  label: string
  matchType: 'title' | 'content'
  preview?: string
}

const searchResults = ref<SearchResult[]>([])

function formatTime(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function openSearchModal() {
  showSearchModal.value = true
}

function closeSearchModal() {
  showSearchModal.value = false
  searchKeyword.value = ''
  searchResults.value = []
}

function buildPreview(content: string, query: string) {
  const lower = content.toLowerCase()
  const index = lower.indexOf(query)
  if (index < 0) return content.slice(0, 80)
  const start = Math.max(0, index - 24)
  const end = Math.min(content.length, index + query.length + 24)
  const excerpt = content.slice(start, end).replace(/\s+/g, ' ')
  return start > 0 ? `...${excerpt}` : excerpt
}

async function runSearch() {
  const query = searchKeyword.value.trim().toLowerCase()
  if (!query) {
    searchResults.value = []
    return
  }

  searching.value = true
  try {
    const titleMatches = props.conversations
      .filter((item) => item.label.toLowerCase().includes(query))
      .map<SearchResult>((item) => ({
        value: item.value,
        label: item.label,
        matchType: 'title',
      }))

    const titleMatchedIds = new Set(titleMatches.map((item) => item.value))
    const contentCandidates = props.conversations.filter((item) => !titleMatchedIds.has(item.value))

    const contentMatches = await Promise.all(
      contentCandidates.map(async (item): Promise<SearchResult | null> => {
        try {
          const response = await $fetch<{ items?: Array<{ content?: string }> }>(
            `/api/conversations/${item.value}/messages`,
          )
          const matched = (response.items || []).find((msg) =>
            String(msg.content || '').toLowerCase().includes(query),
          )
          if (!matched?.content) return null
          return {
            value: item.value,
            label: item.label,
            matchType: 'content',
            preview: buildPreview(matched.content, query),
          }
        } catch {
          return null
        }
      }),
    )

    searchResults.value = [...titleMatches, ...(contentMatches.filter(Boolean) as SearchResult[])]
  } finally {
    searching.value = false
  }
}

function handleSelectSearch(item: SearchResult) {
  emit('select', item.value)
  closeSearchModal()
}
</script>

<template>
  <div class="p-3 w-80 bg-white/70 dark:bg-neutral-900/70 border-r border-neutral-200 dark:border-neutral-800 flex flex-col gap-4">
    <div class="flex w-full gap-4 items-center">
      <UAvatar src="https://avatars.githubusercontent.com/u/259289102?s=48&v=4" size="lg">Lan</UAvatar>
      <div>
        <h2 class="font-black">DeepSpace Chat</h2>
        <span class="text-xs text-neutral-400">与你的 AGENT 一起发现灵感</span>
      </div>
    </div>

    <div class="mt-4 space-y-1">
      <UButton variant="ghost" color="neutral" icon="i-material-symbols:chat-add-on-outline" class="w-full"
        :disabled="props.loading" :to="'/chat'">
        新建聊天
      </UButton>
      <UButton variant="ghost" color="neutral" icon="i-lucide-search" class="w-full" :disabled="props.loading"
        @click="openSearchModal">
        搜索聊天
      </UButton>
    </div>

    <div>
      <h3 class="mt-6 mb-2 px-2 text-sm font-semibold text-neutral-500">所有聊天</h3>
      <div v-if="!props.conversations.length" class="px-2 py-4 text-xs text-neutral-400">暂无会话</div>
      <div v-else class="space-y-1">
        <button v-for="item in props.conversations" :key="item.value"
          class="w-full text-left flex justify-between py-2 px-2 rounded-lg items-center gap-3 border border-transparent transition-all duration-200"
          :class="props.activeConversationId === item.value
            ? 'bg-neutral-200 dark:bg-neutral-800 border-neutral-300 dark:border-neutral-700 shadow-sm'
            : 'hover:bg-neutral-100 dark:hover:bg-neutral-900 hover:border-neutral-200 dark:hover:border-neutral-800 hover:-translate-y-0.5'"
          :disabled="props.loading" @click="emit('select', item.value)">
          <div class="min-w-0">
            <h1 class="text-sm truncate">{{ item.label }}</h1>
            <span class="text-xs text-neutral-400">{{ formatTime(item.updated_at) }}</span>
          </div>
          <UButton variant="ghost" color="neutral" size="xs" icon="i-lucide:trash-2" :disabled="props.loading"
            @click.stop="emit('delete', item.value)" />
        </button>
      </div>
    </div>
  </div>

  <UModal v-model:open="showSearchModal">
    <template #content>
      <div class="w-full max-w-2xl">
        <UChatPalette class="w-full bg-default">

          <div v-if="searchResults.length" class="space-y-2 px-3">
            <button v-for="item in searchResults" :key="`${item.value}-${item.matchType}`"
              class="w-full text-left rounded-lg border border-neutral-200 dark:border-neutral-800 p-3 hover:bg-neutral-50 dark:hover:bg-neutral-900 transition"
              @click="handleSelectSearch(item)">
              <p class="text-sm font-medium truncate">{{ item.label }}</p>
              <p class="text-xs text-neutral-500 mt-1">
                {{ item.matchType === 'title' ? '标题匹配' : '内容匹配' }}
              </p>
              <p v-if="item.preview" class="text-xs text-neutral-400 mt-1 line-clamp-2">{{ item.preview }}</p>
            </button>
          </div>

          <div v-else-if="searchKeyword.trim() && !searching" class="px-4 py-6 text-xs text-neutral-400">
            没有匹配结果
          </div>
          <div v-else class="px-4 py-6 text-xs text-neutral-400">
            输入关键词后点击搜索
          </div>

          <template #prompt>
            <UChatPrompt v-model="searchKeyword" icon="i-lucide-search" variant="naked" @submit="runSearch" />
          </template>
        </UChatPalette>
      </div>
    </template>
  </UModal>
</template>
