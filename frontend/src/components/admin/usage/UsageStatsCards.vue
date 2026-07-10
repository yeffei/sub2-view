<template>
  <div class="usage-stats-grid">
    <div class="card usage-stats-card">
      <div class="usage-stats-icon rounded-lg bg-blue-100 p-2 text-blue-600 dark:bg-blue-900/30">
        <Icon name="document" size="md" />
      </div>
      <div class="min-w-0">
        <p class="text-[11px] font-medium uppercase tracking-[0.12em] text-gray-500">{{ t('usage.totalRequests') }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">{{ stats?.total_requests?.toLocaleString() || '0' }}</p>
        <p class="text-[11px] text-gray-400">{{ t('usage.inSelectedRange') }}</p>
      </div>
    </div>
    <div class="card usage-stats-card">
      <div class="usage-stats-icon rounded-lg bg-amber-100 p-2 text-amber-600 dark:bg-amber-900/30"><svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m21 7.5-9-5.25L3 7.5m18 0-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9" /></svg></div>
      <div class="min-w-0">
        <p class="text-[11px] font-medium uppercase tracking-[0.12em] text-gray-500">{{ t('usage.totalTokens') }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">{{ formatTokens(stats?.total_tokens || 0) }}</p>
        <p class="flex flex-wrap items-center gap-x-1 text-[11px] text-gray-500">
          <span>{{ t('usage.in') }}: {{ formatTokens(stats?.total_input_tokens || 0) }}</span>
          <span>/</span>
          <span>{{ t('usage.out') }}: {{ formatTokens(stats?.total_output_tokens || 0) }}</span>
          <span>/</span>
          <span class="group relative inline-flex cursor-help items-center gap-0.5" tabindex="0">
            <span>{{ cacheLabel() }}: {{ formatTokens(stats?.total_cache_tokens || 0) }}</span>
            <svg
              class="h-3.5 w-3.5 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span
              class="pointer-events-none absolute left-1/2 top-full z-30 mt-2 w-56 -translate-x-1/2 rounded-lg border border-gray-200 bg-white p-3 text-left text-xs text-gray-700 opacity-0 shadow-lg transition-opacity duration-150 group-hover:opacity-100 group-focus:opacity-100 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-200"
            >
              <span class="mb-2 block font-medium text-gray-900 dark:text-white">
                {{ cacheDetailLabel() }}
              </span>
              <span class="flex items-center justify-between gap-3">
                <span>{{ t('usage.cacheCreationTokensLabel') }}</span>
                <span class="tabular-nums">
                  {{ formatTokens(stats?.total_cache_creation_tokens || 0) }}
                </span>
              </span>
              <span class="mt-1 flex items-center justify-between gap-3">
                <span>{{ t('usage.cacheReadTokensLabel') }}</span>
                <span class="tabular-nums">
                  {{ formatTokens(stats?.total_cache_read_tokens || 0) }}
                </span>
              </span>
            </span>
          </span>
        </p>
      </div>
    </div>
    <div class="card usage-stats-card">
      <div class="usage-stats-icon rounded-lg bg-green-100 p-2 text-green-600 dark:bg-green-900/30">
        <Icon name="dollar" size="md" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-[11px] font-medium uppercase tracking-[0.12em] text-gray-500">{{ t('usage.totalCost') }}</p>
        <p class="text-lg font-bold text-green-600">
          ${{ (stats?.total_actual_cost || 0).toFixed(4) }}
        </p>
        <p class="text-[11px] text-gray-400">
          <span class="text-orange-500">{{ t('usage.accountCost') }} ${{ (stats?.total_account_cost || 0).toFixed(4) }}</span>
          <span> · </span>
          <span>{{ t('usage.standardCost') }} ${{ (stats?.total_cost || 0).toFixed(4) }}</span>
        </p>
      </div>
    </div>
    <div class="card usage-stats-card">
      <div class="usage-stats-icon rounded-lg bg-rose-100 p-2 text-rose-600 dark:bg-rose-900/30">
        <Icon name="sparkles" size="md" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-[11px] font-medium uppercase tracking-[0.12em] text-gray-500">{{ t('usage.cacheHitTitle') }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">{{ formatPercent(stats?.cache_read_hit_ratio || 0) }}</p>
        <p class="text-[11px] text-gray-400">
          {{ t('usage.cacheHitRequestsLabel', { value: (stats?.cache_read_hit_requests || 0).toLocaleString() }) }}
          <span> · </span>
          {{ t('usage.cacheReadPerHitLabel', { value: formatTokens(stats?.average_cache_read_tokens_per_hit || 0) }) }}
        </p>
        <p class="text-[11px] text-gray-400">
          {{ t('usage.avgActualInputTokensLabel', { value: formatTokens(stats?.average_actual_input_tokens || 0) }) }}
        </p>
      </div>
    </div>
    <div class="card usage-stats-card">
      <div class="usage-stats-icon rounded-lg bg-purple-100 p-2 text-purple-600 dark:bg-purple-900/30">
        <Icon name="clock" size="md" />
      </div>
      <div class="min-w-0">
        <p class="text-[11px] font-medium uppercase tracking-[0.12em] text-gray-500">{{ t('usage.avgDuration') }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">{{ formatDuration(stats?.average_duration_ms || 0) }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { AdminUsageStatsResponse } from '@/api/admin/usage'
import Icon from '@/components/icons/Icon.vue'

defineProps<{ stats: AdminUsageStatsResponse | null }>()

const { t } = useI18n()

const formatDuration = (ms: number) =>
  ms < 1000 ? `${ms.toFixed(0)}ms` : `${(ms / 1000).toFixed(2)}s`

const formatTokens = (value: number) => {
  if (value >= 1e9) return (value / 1e9).toFixed(2) + 'B'
  if (value >= 1e6) return (value / 1e6).toFixed(2) + 'M'
  if (value >= 1e3) return (value / 1e3).toFixed(2) + 'K'
  return value.toLocaleString()
}

const formatPercent = (value: number) => `${(value * 100).toFixed(1)}%`

const cacheLabel = () => t('usage.cacheTotal')
const cacheDetailLabel = () => t('usage.cacheBreakdown')
</script>

<style scoped>
.usage-stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.85rem;
}

.usage-stats-card {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  min-height: 5.5rem;
  padding: 0.95rem 1rem;
  border: 1px solid rgba(198, 184, 157, 0.38);
  background: rgba(250, 247, 239, 0.72);
  box-shadow: 0 16px 36px -34px rgba(31, 35, 32, 0.24);
}

.usage-stats-icon {
  flex-shrink: 0;
  border: 1px solid rgba(198, 184, 157, 0.24);
}

@media (min-width: 1024px) {
  .usage-stats-grid {
    grid-template-columns: repeat(5, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .usage-stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>

<style>
.dark .usage-stats-card {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.74);
  box-shadow: 0 18px 34px -34px rgba(0, 0, 0, 0.42);
}

.dark .usage-stats-icon {
  border-color: rgba(48, 52, 43, 0.72);
}
</style>
