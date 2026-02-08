<template>
  <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
    <UCard>
      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-gray-500 dark:text-gray-400">总用户数</div>
          <div class="text-3xl font-bold mt-1">1,234</div>
        </div>
        <UIcon name="i-heroicons-users" class="w-8 h-8 text-primary-500" />
      </div>
      <div class="mt-4 text-sm text-green-500 flex items-center gap-1">
        <UIcon name="i-heroicons-arrow-trending-up" />
        <span>较上月增长 12%</span>
      </div>
    </UCard>

    <UCard>
      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-gray-500 dark:text-gray-400">活跃组织</div>
          <div class="text-3xl font-bold mt-1">56</div>
        </div>
        <UIcon name="i-heroicons-building-office" class="w-8 h-8 text-blue-500" />
      </div>
      <div class="mt-4 text-sm text-green-500 flex items-center gap-1">
        <UIcon name="i-heroicons-arrow-trending-up" />
        <span>本周新增 3 个</span>
      </div>
    </UCard>

    <UCard>
      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-gray-500 dark:text-gray-400">API 调用 (24h)</div>
          <div class="text-3xl font-bold mt-1">89.2k</div>
        </div>
        <UIcon name="i-heroicons-server-stack" class="w-8 h-8 text-purple-500" />
      </div>
      <div class="mt-4 text-sm text-red-500 flex items-center gap-1">
        <UIcon name="i-heroicons-arrow-trending-down" />
        <span>较昨日下降 5%</span>
      </div>
    </UCard>

    <UCard>
      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-gray-500 dark:text-gray-400">本月营收</div>
          <div class="text-3xl font-bold mt-1">$4,520</div>
        </div>
        <UIcon name="i-heroicons-currency-dollar" class="w-8 h-8 text-green-500" />
      </div>
      <div class="mt-4 text-sm text-green-500 flex items-center gap-1">
        <UIcon name="i-heroicons-arrow-trending-up" />
        <span>较上月增长 18%</span>
      </div>
    </UCard>
  </div>
</template>
