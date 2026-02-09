<template>
  <div class="min-h-screen bg-white text-slate-900">
    <UContainer class="py-16 max-w-3xl">
      <h1 class="text-3xl font-bold tracking-tight">服务条款</h1>
      <p class="mt-4 text-sm text-slate-500">生效日期：2026-02-04</p>

      <section class="mt-10 space-y-6 text-sm leading-7 text-slate-700">
        <p>
          欢迎使用 DeepSpace Workflows（以下简称“本服务”）。请在使用前仔细阅读本服务条款。
          当您访问或使用本服务，即表示您同意受本条款约束。
        </p>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">1. 服务范围</h2>
          <p>本服务提供 AI 辅助写作、项目管理、知识库与工作流等功能。我们有权根据业务需要进行功能调整。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">2. 账户与安全</h2>
          <p>您需提供真实有效的信息注册账号并妥善保管凭证。因您保管不当造成的损失由您自行承担。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">3. 合规使用</h2>
          <p>您不得利用本服务从事任何违法或侵权行为，包括但不限于恶意攻击、侵犯隐私或传播违法内容。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">4. 知识产权</h2>
          <p>本服务相关的产品、界面、商标与文档等权益归我们所有。您在服务中上传的内容仍归您所有。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">5. 费用与计费</h2>
          <p>如涉及付费服务，将以页面或合同约定为准。我们将提供使用记录与计费明细。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">6. 免责声明</h2>
          <p>本服务基于 AI 生成，输出仅供参考。您应对最终使用的内容承担责任。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">7. 条款变更</h2>
          <p>我们可能适时更新条款并在页面公示。继续使用即视为接受最新条款。</p>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-slate-900">8. 联系我们</h2>
          <p>如有问题请通过官方渠道联系我们。</p>
        </div>
      </section>
    </UContainer>
  </div>
</template>
<script setup lang="ts">
definePageMeta({ layout: false })
useHead({
  title: "服务条款 - Deepspace Workflow",
})
</script>
