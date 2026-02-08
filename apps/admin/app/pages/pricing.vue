<template>
  <UCard>
    <template #header>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3">
          <h3 class="text-lg font-semibold">套餐管理</h3>
          <UBadge color="neutral" variant="soft">规划中</UBadge>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <UButton icon="i-heroicons-plus" label="新建套餐" color="primary" :disabled="true" />
        </div>
      </div>
    </template>

    <div class="py-12">
      <UEmpty icon="i-heroicons-archive-box" title="套餐管理正在建设" description="后续将支持套餐配置、价格分层与上架管理。">
        <template #actions>
          <div class="flex items-center gap-2">
            <UButton color="primary" label="新建套餐" icon="i-heroicons-plus" :disabled="true" />
            <UButton color="neutral" variant="soft" label="查看模型管理" icon="i-heroicons-cpu-chip" to="/models" />
          </div>
        </template>
      </UEmpty>
    </div>
  </UCard>
</template>
