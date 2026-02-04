<template>
    <UContainer>
        <div class="my-12 box-border space-y-8">
            <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
                <div class="space-y-2">
                    <p class="text-sm font-semibold text-primary">AI 科研写作</p>
                    <h1 class="text-3xl font-black tracking-tight text-slate-900 dark:text-slate-100 md:text-4xl">
                        写作项目
                    </h1>
                    <p class="max-w-2xl text-sm leading-6 text-slate-600 dark:text-slate-300">
                        汇总论文写作所需的文献、工作流与 Agent 协作状态，快速掌握每个课题的进度与最近编辑情况。
                    </p>
                </div>
                <div class="flex w-full flex-col gap-3 sm:flex-row md:w-auto">
                    <UButton color="primary" size="lg" class="w-full sm:w-auto" @click="showCreate = true">
                        新建项目
                    </UButton>
                    <UButton color="neutral" variant="soft" size="lg" class="w-full sm:w-auto">
                        导出报告
                    </UButton>
                </div>
            </div>

            <div class="grid gap-4 md:grid-cols-4">
                <UCard class="border-slate-200/70 dark:border-slate-800">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">写作项目</p>
                            <p class="mt-1 text-2xl font-semibold text-slate-900 dark:text-white">
                                {{ stats.projects }}
                            </p>
                        </div>
                        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
                            <UIcon name="i-heroicons-document-text" class="h-5 w-5" />
                        </div>
                    </div>
                </UCard>
                <UCard class="border-slate-200/70 dark:border-slate-800">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">可引用文献</p>
                            <p class="mt-1 text-2xl font-semibold text-slate-900 dark:text-white">
                                {{ stats.literature }}
                            </p>
                        </div>
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-500">
                            <UIcon name="i-heroicons-book-open" class="h-5 w-5" />
                        </div>
                    </div>
                </UCard>
                <UCard class="border-slate-200/70 dark:border-slate-800">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">自动化工作流</p>
                            <p class="mt-1 text-2xl font-semibold text-slate-900 dark:text-white">
                                {{ stats.workflows }}
                            </p>
                        </div>
                        <div
                            class="flex h-10 w-10 items-center justify-center rounded-full bg-amber-500/10 text-amber-500">
                            <UIcon name="i-heroicons-rectangle-group" class="h-5 w-5" />
                        </div>
                    </div>
                </UCard>
                <UCard class="border-slate-200/70 dark:border-slate-800">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-xs font-medium text-slate-500 dark:text-slate-400">协作 Agent</p>
                            <p class="mt-1 text-2xl font-semibold text-slate-900 dark:text-white">
                                {{ stats.agents }}
                            </p>
                        </div>
                        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-sky-500/10 text-sky-500">
                            <UIcon name="i-heroicons-cpu-chip" class="h-5 w-5" />
                        </div>
                    </div>
                </UCard>
            </div>

            <UCard class="border-slate-200/70 dark:border-slate-800">
                <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                    <div class="flex w-full flex-col gap-4 sm:flex-row">
                        <UInput v-model="query" icon="i-heroicons-magnifying-glass" placeholder="搜索项目、主题、Agent 或文献"
                            size="lg" class="flex-1" />
                        <USelectMenu v-model="status" :options="statusOptions" size="lg" class="min-w-[160px]" />
                        <USelectMenu v-model="sort" :options="sortOptions" size="lg" class="min-w-[160px]" />
                    </div>
                    <div class="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400">
                        <span>筛选结果</span>
                        <UBadge color="primary" variant="soft">{{ filteredProjects.length }}</UBadge>
                    </div>
                </div>
            </UCard>

            <div class="grid gap-6 lg:grid-cols-3">
                <UCard v-for="project in filteredProjects" :key="project.id"
                    class="flex h-full flex-col border-slate-200/70 transition-shadow duration-200 hover:shadow-lg dark:border-slate-800">
                    <template #header>
                        <div class="flex items-start justify-between gap-4">
                            <div>
                                <p class="text-xs font-medium text-slate-500 dark:text-slate-400">{{ project.owner }}
                                </p>
                                <h2 class="mt-1 text-lg font-semibold text-slate-900 dark:text-white">
                                    {{ project.name }}
                                </h2>
                            </div>
                            <UBadge :color="project.badgeColor" variant="soft">{{ project.status }}</UBadge>
                        </div>
                    </template>

                    <div class="flex flex-1 flex-col gap-4">
                        <p class="text-sm leading-6 text-slate-600 dark:text-slate-300">
                            {{ project.description }}
                        </p>
                        <div
                            class="grid gap-2 rounded-xl bg-slate-50 p-3 text-xs text-slate-600 dark:bg-slate-900/60 dark:text-slate-300">
                            <div class="flex items-center justify-between">
                                <span>文献</span>
                                <span class="font-semibold text-slate-900 dark:text-white">{{ project.literature
                                    }}</span>
                            </div>
                            <div class="flex items-center justify-between">
                                <span>工作流</span>
                                <span class="font-semibold text-slate-900 dark:text-white">{{ project.workflow }}</span>
                            </div>
                            <div class="flex items-center justify-between">
                                <span>Agent</span>
                                <span class="font-semibold text-slate-900 dark:text-white">{{ project.agents }}</span>
                            </div>
                            <div class="flex items-center justify-between">
                                <span>最后编辑</span>
                                <span class="font-semibold text-slate-900 dark:text-white">{{ project.lastEdited
                                    }}</span>
                            </div>
                        </div>
                        <div class="flex flex-wrap gap-2">
                            <UBadge v-for="tag in project.tags" :key="tag" color="neutral" variant="soft">
                                {{ tag }}
                            </UBadge>
                        </div>
                    </div>

                    <template #footer>
                        <div class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400">
                            <div class="flex items-center gap-2">
                                <UIcon name="i-heroicons-clock" class="h-4 w-4" />
                                <span>更新于 {{ project.updatedAt }}</span>
                            </div>
                            <NuxtLink :to="`/projects/${project.id}`"
                                class="font-medium text-primary transition-colors hover:text-primary/80">
                                继续写作
                            </NuxtLink>
                        </div>
                    </template>
                </UCard>
            </div>

            <UCard v-if="filteredProjects.length === 0" class="border-slate-200/70 dark:border-slate-800">
                <div class="flex flex-col items-center gap-3 py-10 text-center">
                    <div
                        class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-500 dark:bg-slate-800 dark:text-slate-300">
                        <UIcon name="i-heroicons-magnifying-glass" class="h-6 w-6" />
                    </div>
                    <div>
                        <h3 class="text-base font-semibold text-slate-900 dark:text-white">未找到匹配项目</h3>
                        <p class="mt-1 text-sm text-slate-500 dark:text-slate-400">尝试调整关键词或筛选条件。</p>
                    </div>
                    <UButton color="neutral" variant="soft" size="sm" @click="resetFilters">
                        重置筛选
                    </UButton>
                </div>
            </UCard>
        </div>

        <ProjectCreateProjectModal v-model:open="showCreate" @created="onProjectCreated" />
    </UContainer>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type ProjectItem = {
    id: number
    name: string
    description?: string | null
    created_at?: string
    updated_at?: string
}

type ProjectStats = {
    projects: number
    literature: number
    workflows: number
    agents: number
}

const query = ref('')
const statusOptions = ['全部', '进行中', '评审中', '已上线', '已归档']
const status = ref('全部')
const showCreate = ref(false)
const sortOptions = [
    { label: '最近更新', value: 'updated' },
    { label: '进度优先', value: 'progress' },
    { label: '名称排序', value: 'name' }
]
const sort = ref('updated')

const requestHeaders = useRequestHeaders(['cookie'])

const { data, refresh: refreshProjects } = await useAsyncData('projects', () =>
    $fetch<{ items: ProjectItem[] }>('/api/projects', { headers: requestHeaders })
)

const { data: statsData, refresh: refreshStats } = await useAsyncData('project-stats', () =>
    $fetch<ProjectStats>('/api/projects/stats', { headers: requestHeaders })
)

const stats = computed<ProjectStats>(() => ({
    projects: statsData.value?.projects ?? 0,
    literature: statsData.value?.literature ?? 0,
    workflows: statsData.value?.workflows ?? 0,
    agents: statsData.value?.agents ?? 0
}))

const projects = computed(() => {
    const items = data.value?.items ?? []
    return items.map((item) => {
        const updatedAt = formatRelative(item.updated_at || item.created_at)
        return {
            id: item.id,
            name: item.name,
            description: item.description || '暂无描述',
            status: '进行中',
            badgeColor: 'emerald',
            progress: 0,
            updatedAt,
            owner: '当前组织',
            tags: [],
            literature: '-',
            workflow: '-',
            agents: '-',
            lastEdited: item.updated_at || item.created_at || ''
        }
    })
})

const filteredProjects = computed(() => {
    const keyword = query.value.trim().toLowerCase()
    const normalized = projects.value.filter((project) => {
        const matchesStatus = status.value === '全部' || project.status === status.value
        if (!matchesStatus) return false

        if (!keyword) return true

        const haystack = [
            project.name,
            project.description,
            project.owner,
            project.tags.join(' ')
        ]
            .join(' ')
            .toLowerCase()

        return haystack.includes(keyword)
    })

    if (sort.value === 'progress') {
        return [...normalized].sort((a, b) => b.progress - a.progress)
    }

    if (sort.value === 'name') {
        return [...normalized].sort((a, b) => a.name.localeCompare(b.name, 'zh-CN'))
    }

    return normalized
})

const resetFilters = () => {
    query.value = ''
    status.value = '全部'
    sort.value = 'updated'
}

async function onProjectCreated() {
    await refreshProjects()
    await refreshStats()
}

function formatRelative(value?: string) {
    if (!value) return '刚刚'
    const date = new Date(value)
    const diffMs = Date.now() - date.getTime()
    if (Number.isNaN(diffMs)) return value
    const diffMin = Math.floor(diffMs / 60000)
    if (diffMin < 1) return '刚刚'
    if (diffMin < 60) return `${diffMin} 分钟前`
    const diffHr = Math.floor(diffMin / 60)
    if (diffHr < 24) return `${diffHr} 小时前`
    const diffDay = Math.floor(diffHr / 24)
    if (diffDay < 7) return `${diffDay} 天前`
    return date.toISOString().slice(0, 10)
}
</script>
