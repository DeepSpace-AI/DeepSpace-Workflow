<template>
  <div class="min-h-screen bg-white text-gray-900 font-sans">
    <!-- 顶部导航 -->
    <header class="fixed top-0 left-0 right-0 z-50 bg-white/90 backdrop-blur-sm border-b border-gray-200">
      <nav class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        <div class="text-xl font-bold tracking-tight flex items-center gap-2">
          <div class="w-8 h-8 bg-gray-900 flex items-center justify-center">
            <span class="text-white text-sm">DS</span>
          </div>
          <span>DeepSpace AI</span>
        </div>
        <div class="hidden md:flex items-center gap-8">
          <a href="#features" class="text-sm text-gray-600 hover:text-gray-900 transition-colors">核心功能</a>
          <a href="#technology" class="text-sm text-gray-600 hover:text-gray-900 transition-colors">技术优势</a>
          <a href="#scenarios" class="text-sm text-gray-600 hover:text-gray-900 transition-colors">应用场景</a>
          <a href="#workflow" class="text-sm text-gray-600 hover:text-gray-900 transition-colors">协作流程</a>
        </div>
        <NuxtLink :to="user ? '/projects' : '/sign-in'"
          class="px-5 py-2 bg-gray-900 text-white text-sm font-medium hover:bg-gray-800 transition-colors">
          {{ user ? '控制台' : '开始探索' }}
        </NuxtLink>
      </nav>
    </header>

    <!-- Hero 区域 -->
    <section class="pt-32 pb-20 border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6">
        <div class="max-w-4xl">
          <div class="inline-block px-3 py-1 mb-6 border border-gray-300 text-xs font-medium tracking-wide">
            AI 科研写作协作平台
          </div>
          <h1 class="text-5xl md:text-6xl font-bold leading-tight tracking-tight mb-6">
            让科研写作<br>
            <span class="text-gray-500">更高效、更智能</span>
          </h1>
          <p class="text-xl text-gray-600 max-w-2xl mb-10 leading-relaxed">
            DeepSpace AI 专注于为科研工作者提供智能写作辅助，通过 AI 技术提升论文撰写、文献分析、协作编辑的效率。
          </p>
          <div class="flex flex-wrap gap-4">
            <button class="px-8 py-3 bg-gray-900 text-white font-medium hover:bg-gray-800 transition-colors">
              免费试用
            </button>
            <button class="px-8 py-3 border border-gray-300 font-medium hover:bg-gray-50 transition-colors">
              观看演示
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- 核心功能 -->
    <section id="features" class="py-20 border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6">
        <div class="mb-12">
          <h2 class="text-3xl font-bold tracking-tight mb-4">核心功能</h2>
          <p class="text-gray-600 max-w-xl">全方位智能写作工具，助力科研工作</p>
        </div>
        <div class="grid md:grid-cols-3 gap-0 md:gap-px bg-gray-200">
          <div v-for="(item, index) in features" :key="index"
            class="bg-white p-8 hover:bg-gray-50 transition-colors group">
            <div
              class="w-12 h-12 border border-gray-300 flex items-center justify-center mb-6 group-hover:border-gray-400 transition-colors">
              <span class="text-xl font-bold">{{ index + 1 }}</span>
            </div>
            <h3 class="text-xl font-semibold mb-3">{{ item.title }}</h3>
            <p class="text-gray-600 text-sm leading-relaxed">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- 技术优势 -->
    <section id="technology" class="py-20 border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6">
        <div class="grid md:grid-cols-2 gap-16 items-center">
          <div>
            <h2 class="text-3xl font-bold tracking-tight mb-4">技术优势</h2>
            <p class="text-gray-600 mb-8 leading-relaxed">
              基于前沿大语言模型技术，专为科研场景优化，提供精准、可靠的写作辅助。
            </p>
            <div class="space-y-6">
              <div v-for="(tech, index) in techStack" :key="index" class="border-l-2 border-gray-300 pl-6">
                <h4 class="font-semibold mb-1">{{ tech.name }}</h4>
                <p class="text-sm text-gray-600">{{ tech.desc }}</p>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 border border-gray-200 p-8">
            <div class="space-y-4 font-mono text-sm">
              <div v-for="(line, index) in codeExample" :key="index" class="flex gap-4">
                <span class="text-gray-400 w-6 text-right">{{ index + 1 }}</span>
                <span :class="line.class">{{ line.text }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 应用场景 -->
    <section id="scenarios" class="py-20 border-b border-gray-200 bg-gray-50">
      <div class="max-w-7xl mx-auto px-6">
        <div class="mb-12">
          <h2 class="text-3xl font-bold tracking-tight mb-4">应用场景</h2>
        </div>
        <div class="grid md:grid-cols-4 gap-0 md:gap-px bg-gray-300">
          <div v-for="(stat, index) in stats" :key="index" class="bg-white p-8 text-center">
            <div class="text-4xl font-bold mb-2 tracking-tight">{{ stat.value }}</div>
            <div class="text-sm text-gray-600">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 协作流程 -->
    <section id="workflow" class="py-20 border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-6">
        <div class="mb-12">
          <h2 class="text-3xl font-bold tracking-tight mb-4">协作流程</h2>
          <p class="text-gray-600 max-w-xl">从构思到发表，全程智能辅助</p>
        </div>
        <div class="relative">
          <div class="grid md:grid-cols-6 gap-0 md:gap-px bg-gray-200">
            <div v-for="(step, index) in pipelineSteps" :key="index"
              class="bg-white p-6 hover:bg-gray-50 transition-colors relative">
              <div
                class="absolute -top-3 left-6 w-6 h-6 bg-gray-900 text-white text-xs flex items-center justify-center font-medium">
                {{ index + 1 }}
              </div>
              <div class="w-10 h-10 border border-gray-300 flex items-center justify-center mb-4 mt-2">
                <span class="text-lg">{{ step.icon }}</span>
              </div>
              <h4 class="font-semibold mb-2">{{ step.title }}</h4>
              <p class="text-xs text-gray-600 leading-relaxed">{{ step.description }}</p>
            </div>
          </div>
          <div class="hidden md:block absolute top-1/2 -right-4 transform -translate-y-1/2 z-10">
            <div class="w-8 h-8 bg-white border border-gray-300 flex items-center justify-center">
              <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="square" stroke-linejoin="miter" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 页脚 -->
    <footer class="py-8 bg-gray-900 text-white">
      <div class="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-4">
        <div class="text-sm text-gray-400">
          © 2025 DeepSpace AI. All rights reserved.
        </div>
        <div class="flex gap-6 text-sm text-gray-400">
          <a href="#" class="hover:text-white transition-colors">隐私政策</a>
          <a href="#" class="hover:text-white transition-colors">使用条款</a>
          <a href="#" class="hover:text-white transition-colors">帮助中心</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const { data } = await useAsyncData('me', () => $fetch('/api/auth/me').catch(() => null))
const user = computed(() => !!data.value)

const features = [
  {
    title: '智能文献综述',
    description: '自动分析、归纳和总结大量文献，快速生成结构化的文献综述框架。'
  },
  {
    title: '论文写作辅助',
    description: '提供段落扩写、句子润色、学术表达建议，提升论文表达质量。'
  },
  {
    title: '参考文献管理',
    description: '自动格式化引用、生成参考文献列表，支持多种学术格式。'
  },
  {
    title: '数据可视化建议',
    description: '根据数据特点推荐合适的图表类型，辅助结果展示。'
  },
  {
    title: '多语言翻译',
    description: '专业术语精准翻译，助您跨越语言障碍进行国际学术交流。'
  },
  {
    title: '协作审阅',
    description: '多人实时协作、批注反馈、版本管理，提升团队协作效率。'
  }
]

const techStack = [
  { name: '大语言模型', desc: '基于前沿 LLM 技术，专为学术写作优化' },
  { name: '知识图谱', desc: '构建学科知识网络，提供精准知识推荐' },
  { name: '自然语言处理', desc: '深度语义理解，确保表达准确专业' },
  { name: '实时协作引擎', desc: '多人同步编辑，毫秒级响应体验' }
]

const codeExample = [
  { text: 'import DeepSpace from "@deepspace/ai";', class: 'text-purple-600' },
  { text: '', class: '' },
  { text: 'const paper = new Paper({', class: 'text-gray-700' },
  { text: '  type: "research_article",', class: 'text-blue-600' },
  { text: '  language: "zh-CN",', class: 'text-blue-600' },
  { text: '  collaborators: 5', class: 'text-blue-600' },
  { text: '});', class: 'text-gray-700' },
  { text: '', class: '' },
  { text: 'await paper.generateOutline();', class: 'text-gray-700' }
]

const stats = [
  { value: '100K+', label: '活跃科研用户' },
  { value: '500K+', label: '辅助论文数' },
  { value: '50+', label: '学科领域' },
  { value: '98%', label: '用户满意度' }
]

const pipelineSteps = [
  {
    title: '选题构思',
    description: 'AI 辅助分析研究热点，推荐创新选题方向',
    icon: '�'
  },
  {
    title: '文献调研',
    description: '智能文献检索、分类管理与摘要生成',
    icon: '📚'
  },
  {
    title: '大纲规划',
    description: '自动生成论文结构框架，明确各章节内容',
    icon: '📋'
  },
  {
    title: '内容撰写',
    description: '段落生成、润色改写、表达优化',
    icon: '✍️'
  },
  {
    title: '协作审校',
    description: '团队在线批注、修改追踪、质量把控',
    icon: '👥'
  },
  {
    title: '投稿发表',
    description: '格式排版、查重检测、一键投稿',
    icon: '�'
  }
]
</script>

<style>
html {
  scroll-behavior: smooth;
}
</style>
