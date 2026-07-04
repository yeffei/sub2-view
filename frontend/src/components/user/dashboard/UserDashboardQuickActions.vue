<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <div class="flex items-center justify-between gap-3">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">快捷入口</h2>
        <span class="action-caption">常用入口</span>
      </div>
    </div>
    <div class="space-y-3 p-4">
      <button @click="router.push('/keys')" class="group action-row flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-primary-100 transition-transform group-hover:scale-105 dark:bg-primary-900/30">
          <Icon name="key" size="lg" class="text-primary-600 dark:text-primary-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">密钥与调用</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">查看、创建与整理调用凭证</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-primary-500 dark:text-dark-500"
        />
      </button>

      <button @click="router.push('/usage')" class="group action-row flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-emerald-100 transition-transform group-hover:scale-105 dark:bg-emerald-900/30">
          <Icon name="chart" size="lg" class="text-emerald-600 dark:text-emerald-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">用量与账册</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">按时段回看请求、Token 与费用</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-emerald-500 dark:text-dark-500"
        />
      </button>

      <button @click="router.push('/profile')" class="group action-row flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-amber-100 transition-transform group-hover:scale-105 dark:bg-amber-900/30">
          <Icon :name="isSimple ? 'user' : 'shield'" size="lg" class="text-amber-600 dark:text-amber-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">{{ isSimple ? '账户与规则' : '账户与安全' }}</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">{{ isSimple ? '查看个人资料与站点规则' : '查看资料、绑定方式、提醒与安全设置' }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-amber-500 dark:text-dark-500"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'

defineProps<{
  isSimple: boolean
}>()

const router = useRouter()
</script>

<style scoped>
.action-caption {
  font-size: 0.72rem;
  letter-spacing: 0.12em;
  color: #7b6a53;
}

.action-row {
  border: 1px solid rgba(216, 205, 185, 0.72);
}

.action-row:hover {
  border-color: rgba(167, 58, 42, 0.36);
}

.dark .action-caption {
  color: #879186;
}

.dark .action-row {
  border-color: rgba(48, 52, 43, 0.9);
}

.dark .action-row:hover {
  border-color: rgba(167, 58, 42, 0.52);
}
</style>

