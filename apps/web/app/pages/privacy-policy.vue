<template>
  <div class="min-h-screen bg-white text-slate-900">
    <UContainer class="py-16 max-w-3xl">
      <h1 class="text-3xl font-bold tracking-tight">隐私政策</h1>
      <p class="mt-4 text-sm text-slate-500">生效日期：2026-02-04</p>

      <section class="mt-10 space-y-6 text-sm leading-7 text-slate-700">
        <p>
          本隐私政策说明我们如何收集、使用、存储和保护您的个人信息。使用本服务即表示您同意本政策。
        </p>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">1. 信息收集</h2>
          <p>我们可能收集账号信息、使用日志、设备与浏览器信息，用于提供与改进服务。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">2. 信息使用</h2>
          <p>收集的信息用于身份验证、计费、审计、安全风控及产品优化，不会用于未经授权的用途。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">3. 信息共享</h2>
          <p>除法律法规要求或获得您授权外，我们不会向第三方共享您的个人信息。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">4. 数据安全</h2>
          <p>我们采取合理的技术与管理措施保护数据安全，防止泄露、篡改或丢失。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">5. 数据保留与删除</h2>
          <p>在满足法律法规或业务需要的前提下，我们会保留必要数据。您可申请删除相关数据。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">6. Cookie 使用</h2>
          <p>我们使用 Cookie 用于登录态与安全目的。您可以通过浏览器设置管理 Cookie。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">7. 政策变更</h2>
          <p>本政策可能更新，我们会在页面公示。继续使用即视为同意变更。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">8. 联系我们</h2>
          <p>如对隐私政策有疑问，请通过官方渠道联系我们。</p>
        </div>
      </section>
    </UContainer>
  </div>
</template>
<script setup lang="ts">
definePageMeta({ layout: false })
useHead({
  title: "隐私政策 - Deepspace Workflow",
})
</script>
