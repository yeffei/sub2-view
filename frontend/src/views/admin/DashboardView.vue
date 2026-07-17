<template>
  <AppLayout>
    <div class="admin-dashboard-shell space-y-6">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>

      <template v-else-if="stats">

        <section class="card admin-briefing-card p-5">
          <div class="admin-briefing-head">
            <div>
              <span class="admin-page-kicker">{{ briefingKicker }}</span>
              <h3>{{ briefingTitle }}</h3>
              <p>{{ briefingLead }}</p>
            </div>
            <div class="admin-attribution-actions">
              <router-link class="btn btn-secondary" :to="briefingPrimaryLink">{{ t('admin.dashboard.briefing.primaryAction') }}</router-link>
              <router-link class="btn btn-secondary" :to="briefingSecondaryLink">{{ t('admin.dashboard.briefing.secondaryAction') }}</router-link>
            </div>
          </div>

          <div class="admin-briefing-grid mt-4">
            <article class="admin-briefing-panel">
              <span>{{ t('admin.dashboard.briefing.trafficTitle') }}</span>
              <strong>{{ briefingTrafficValue }}</strong>
              <p>{{ briefingTrafficNote }}</p>
            </article>
            <article class="admin-briefing-panel">
              <span>{{ t('admin.dashboard.briefing.costTitle') }}</span>
              <strong>{{ briefingCostValue }}</strong>
              <p>{{ briefingCostNote }}</p>
            </article>
            <article class="admin-briefing-panel">
              <span>{{ t('admin.dashboard.briefing.anomalyTitle') }}</span>
              <strong>{{ briefingAnomalyValue }}</strong>
              <p>{{ briefingAnomalyNote }}</p>
            </article>
            <article class="admin-briefing-panel">
              <span>{{ t('admin.dashboard.briefing.actionTitle') }}</span>
              <strong>{{ briefingActionValue }}</strong>
              <p>{{ briefingActionNote }}</p>
            </article>
          </div>

          <div class="admin-briefing-foot mt-4">
            <div class="admin-briefing-tags">
              <span v-for="tag in briefingTags" :key="tag" class="admin-briefing-tag">{{ tag }}</span>
              <span v-if="!briefingTags.length" class="admin-briefing-tag admin-briefing-tag-muted">{{ t('admin.dashboard.briefing.tagObserve') }}</span>
            </div>
            <p class="admin-briefing-summary">{{ briefingSummary }}</p>
          </div>
        </section>

        <!-- Row 1: Core Stats -->
        <div class="admin-stat-grid grid grid-cols-2 gap-4 lg:grid-cols-4">
          <!-- Total API Keys -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30">
                <Icon name="key" size="md" class="text-blue-600 dark:text-blue-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.apiKeys') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ stats.total_api_keys }}
                </p>
                <p class="text-xs text-green-600 dark:text-green-400">
                  {{ stats.active_api_keys }} {{ t('common.active') }}
                </p>
              </div>
            </div>
          </div>

          <!-- Service Accounts -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-purple-100 p-2 dark:bg-purple-900/30">
                <Icon name="server" size="md" class="text-purple-600 dark:text-purple-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.accounts') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ stats.total_accounts }}
                </p>
                <p class="text-xs">
                  <span class="text-green-600 dark:text-green-400"
                    >{{ stats.normal_accounts }} {{ t('common.active') }}</span
                  >
                  <span v-if="stats.error_accounts > 0" class="ml-1 text-red-500"
                    >{{ stats.error_accounts }} {{ t('common.error') }}</span
                  >
                </p>
              </div>
            </div>
          </div>

          <!-- Today Requests -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-green-100 p-2 dark:bg-green-900/30">
                <Icon name="chart" size="md" class="text-green-600 dark:text-green-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.todayRequests') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ stats.today_requests }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('common.total') }}: {{ formatNumber(stats.total_requests) }}
                </p>
              </div>
            </div>
          </div>

          <!-- New Users Today -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-emerald-100 p-2 dark:bg-emerald-900/30">
                <Icon name="userPlus" size="md" class="text-emerald-600 dark:text-emerald-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.users') }}
                </p>
                <p class="text-xl font-bold text-emerald-600 dark:text-emerald-400">
                  +{{ stats.today_new_users }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('common.total') }}: {{ formatNumber(stats.total_users) }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Row 2: Token Stats -->
        <div class="admin-stat-grid grid grid-cols-2 gap-4 lg:grid-cols-4">
          <!-- Today Tokens -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
                <Icon name="cube" size="md" class="text-amber-600 dark:text-amber-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.todayTokens') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatTokens(stats.today_tokens) }}
                </p>
                <p class="text-xs">
                  <span
                    class="text-green-600 dark:text-green-400"
                    :title="t('admin.dashboard.actual')"
                    >${{ formatCost(stats.today_actual_cost) }}</span
                  >
                  <span class="text-gray-400 dark:text-gray-500"> / </span>
                  <span
                    class="text-orange-500 dark:text-orange-400"
                    :title="t('admin.dashboard.accountCost')"
                    >${{ formatCost(stats.today_account_cost) }}</span
                  >
                  <span class="text-gray-400 dark:text-gray-500"> / </span>
                  <span
                    class="text-gray-400 dark:text-gray-500"
                    :title="t('admin.dashboard.standard')"
                    >${{ formatCost(stats.today_cost) }}</span
                  >
                </p>
              </div>
            </div>
          </div>

          <!-- Total Tokens -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-indigo-100 p-2 dark:bg-indigo-900/30">
                <Icon name="database" size="md" class="text-indigo-600 dark:text-indigo-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.totalTokens') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatTokens(stats.total_tokens) }}
                </p>
                <p class="text-xs">
                  <span
                    class="text-green-600 dark:text-green-400"
                    :title="t('admin.dashboard.actual')"
                    >${{ formatCost(stats.total_actual_cost) }}</span
                  >
                  <span class="text-gray-400 dark:text-gray-500"> / </span>
                  <span
                    class="text-orange-500 dark:text-orange-400"
                    :title="t('admin.dashboard.accountCost')"
                    >${{ formatCost(stats.total_account_cost) }}</span
                  >
                  <span class="text-gray-400 dark:text-gray-500"> / </span>
                  <span
                    class="text-gray-400 dark:text-gray-500"
                    :title="t('admin.dashboard.standard')"
                    >${{ formatCost(stats.total_cost) }}</span
                  >
                </p>
              </div>
            </div>
          </div>

          <!-- Performance (RPM/TPM) -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-violet-100 p-2 dark:bg-violet-900/30">
                <Icon name="bolt" size="md" class="text-violet-600 dark:text-violet-400" :stroke-width="2" />
              </div>
              <div class="flex-1">
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.performance') }}
                </p>
                <div class="flex items-baseline gap-2">
                  <p class="text-xl font-bold text-gray-900 dark:text-white">
                    {{ formatTokens(stats.rpm) }}
                  </p>
                  <span class="text-xs text-gray-500 dark:text-gray-400">RPM</span>
                </div>
                <div class="flex items-baseline gap-2">
                  <p class="text-sm font-semibold text-violet-600 dark:text-violet-400">
                    {{ formatTokens(stats.tpm) }}
                  </p>
                  <span class="text-xs text-gray-500 dark:text-gray-400">TPM</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Avg Response Time -->
          <div class="card admin-stat-card p-4">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-rose-100 p-2 dark:bg-rose-900/30">
                <Icon name="clock" size="md" class="text-rose-600 dark:text-rose-400" :stroke-width="2" />
              </div>
              <div>
                <p class="text-xs font-medium text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.avgResponse') }}
                </p>
                <p class="text-xl font-bold text-gray-900 dark:text-white">
                  {{ formatDuration(stats.average_duration_ms) }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ stats.active_users }} {{ t('admin.dashboard.activeUsers') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <section class="card admin-morning-sheet-card p-5">
          <div class="admin-morning-sheet-head">
            <div>
              <span class="admin-page-kicker">{{ morningSheetKicker }}</span>
              <h3>{{ morningSheetTitle }}</h3>
              <p>{{ morningSheetLead }}</p>
            </div>
            <div class="admin-attribution-actions">
              <router-link class="btn btn-secondary" :to="morningSheetPrimaryLink">{{ t('admin.dashboard.morningSheet.primaryAction') }}</router-link>
              <router-link class="btn btn-secondary" :to="morningSheetSecondaryLink">{{ t('admin.dashboard.morningSheet.secondaryAction') }}</router-link>
              <button type="button" class="btn btn-secondary" @click="copyMorningSheet">{{ t('admin.dashboard.morningSheet.copyAction') }}</button>
              <button type="button" class="btn btn-secondary" @click="downloadMorningSheet">{{ t('admin.dashboard.morningSheet.downloadAction') }}</button>
            </div>
          </div>

          <div class="admin-morning-sheet-layout mt-4">
            <article class="admin-morning-sheet-summary">
              <p>{{ morningSheetSummary }}</p>
            </article>
            <div class="admin-morning-sheet-columns">
              <article class="admin-morning-sheet-panel">
                <span>{{ t('admin.dashboard.morningSheet.windowsTitle') }}</span>
                <strong>{{ morningSheetWindowValue }}</strong>
                <p>{{ morningSheetWindowNote }}</p>
              </article>
              <article class="admin-morning-sheet-panel">
                <span>{{ t('admin.dashboard.morningSheet.anomalyTitle') }}</span>
                <strong>{{ morningSheetAnomalyValue }}</strong>
                <p>{{ morningSheetAnomalyNote }}</p>
              </article>
              <article class="admin-morning-sheet-panel">
                <span>{{ t('admin.dashboard.morningSheet.watchTitle') }}</span>
                <strong>{{ morningSheetWatchValue }}</strong>
                <p>{{ morningSheetWatchNote }}</p>
              </article>
            </div>
          </div>
        </section>

        <section class="card admin-attribution-card p-5">
          <div class="admin-attribution-head">
            <div>
              <span class="admin-page-kicker">{{ t('admin.dashboard.attribution.kicker') }}</span>
              <h3>{{ t('admin.dashboard.attribution.title') }}</h3>
              <p>{{ attributionHeadline }}</p>
            </div>
            <div class="admin-attribution-actions">
              <router-link class="btn btn-secondary" :to="attributionPrimaryLink">{{ attributionPrimaryAction }}</router-link>
              <router-link class="btn btn-secondary" :to="attributionSecondaryLink">{{ attributionSecondaryAction }}</router-link>
            </div>
          </div>

          <div v-if="attributionLoading" class="mt-4 flex items-center justify-center py-6">
            <LoadingSpinner size="md" />
          </div>
          <div v-else class="admin-attribution-grid mt-4">
            <article class="admin-attribution-panel">
              <span>{{ t('admin.dashboard.attribution.ownerTitle') }}</span>
              <strong>{{ topOwnerLabel }}</strong>
              <p>{{ ownerInsight }}</p>
            </article>
            <article class="admin-attribution-panel">
              <span>{{ t('admin.dashboard.attribution.sourceTitle') }}</span>
              <strong>{{ topSourceLabel }}</strong>
              <p>{{ sourceInsight }}</p>
            </article>
            <article class="admin-attribution-panel">
              <span>{{ t('admin.dashboard.attribution.platformTitle') }}</span>
              <strong>{{ topPlatformLabel }}</strong>
              <p>{{ platformInsight }}</p>
            </article>
            <article class="admin-attribution-panel">
              <span>{{ t('admin.dashboard.attribution.statusTitle') }}</span>
              <strong>{{ topStatusLabel }}</strong>
              <p>{{ statusInsight }}</p>
            </article>
          </div>
        </section>

        <section class="card admin-risk-card p-5">
          <div class="admin-risk-head">
            <div>
              <span class="admin-page-kicker">{{ t('admin.dashboard.risk.kicker') }}</span>
              <h3>{{ t('admin.dashboard.risk.title') }}</h3>
              <p>{{ riskHeadline }}</p>
            </div>
            <div class="admin-attribution-actions">
              <router-link class="btn btn-secondary" :to="riskPrimaryLink">{{ t('admin.dashboard.risk.primaryAction') }}</router-link>
              <router-link class="btn btn-secondary" :to="riskSecondaryLink">{{ t('admin.dashboard.risk.secondaryAction') }}</router-link>
            </div>
          </div>

          <div v-if="riskListLoading" class="mt-4 flex items-center justify-center py-6">
            <LoadingSpinner size="md" />
          </div>
          <div v-else-if="highRiskUsers.length" class="admin-risk-list mt-4">
            <article v-for="user in highRiskUsers" :key="user.userId" class="admin-risk-item">
              <div class="admin-risk-item-main">
                <div class="admin-risk-item-copy">
                  <div class="admin-risk-item-title">
                    <strong>{{ user.displayName }}</strong>
                    <span>{{ user.email }}</span>
                  </div>
                  <p>{{ user.summary }}</p>
                </div>
                <div class="admin-risk-item-metrics">
                  <div>
                    <span>{{ t('admin.dashboard.risk.metricFailedCount') }}</span>
                    <strong>{{ formatNumber(user.failedCount) }}</strong>
                  </div>
                  <div>
                    <span>{{ t('admin.dashboard.risk.metricFailureRate') }}</span>
                    <strong>{{ formatPercent(user.failureRate) }}</strong>
                  </div>
                  <div>
                    <span>{{ t('admin.dashboard.risk.metricRequestDensity') }}</span>
                    <strong>{{ user.requestDensityLabel }}</strong>
                  </div>
                </div>
              </div>

              <div class="admin-risk-item-foot">
                <div class="admin-risk-tags">
                  <span v-for="tag in user.tags" :key="tag.label" class="admin-risk-tag" :data-tone="tag.tone">{{ tag.label }}</span>
                </div>
                <div class="admin-risk-links">
                  <router-link class="btn btn-secondary" :to="getRiskUsageLink(user.userId)">{{ t('admin.dashboard.risk.viewUsage') }}</router-link>
                  <router-link class="btn btn-secondary" :to="getRiskOpsLink(user.userId)">{{ t('admin.dashboard.risk.viewErrors') }}</router-link>
                </div>
              </div>
            </article>
          </div>
          <div v-else class="admin-risk-empty mt-4">
            {{ riskEmptyMessage }}
          </div>
        </section>
        <!-- Quick Actions -->
        <div class="card p-4">
          <div class="mb-3 flex items-center justify-between">
            <h2 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.dashboard.quickActions') }}
            </h2>
          </div>
          <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
            <button
              v-if="canUseBatchImage"
              type="button"
              class="group flex items-center gap-3 rounded-lg bg-gray-50 p-3 text-left transition-colors hover:bg-sky-50 dark:bg-dark-800/50 dark:hover:bg-sky-900/20"
              @click="router.push('/batch-image')"
            >
              <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-sky-100 text-sky-600 dark:bg-sky-900/30 dark:text-sky-400">
                <Icon name="sparkles" size="md" :stroke-width="2" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block text-sm font-medium text-gray-900 dark:text-white">
                  {{ t('admin.dashboard.batchImage') }}
                </span>
                <span class="block text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.batchImageDesc') }}
                </span>
              </span>
              <Icon name="chevronRight" size="sm" class="text-gray-400 group-hover:text-sky-500" />
            </button>
            <button
              type="button"
              class="group flex items-center gap-3 rounded-lg bg-gray-50 p-3 text-left transition-colors hover:bg-emerald-50 dark:bg-dark-800/50 dark:hover:bg-emerald-900/20"
              @click="router.push('/admin/groups')"
            >
              <span class="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400">
                <Icon name="grid" size="md" :stroke-width="2" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="block text-sm font-medium text-gray-900 dark:text-white">
                  {{ t('admin.dashboard.groupPricing') }}
                </span>
                <span class="block text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.dashboard.groupPricingDesc') }}
                </span>
              </span>
              <Icon name="chevronRight" size="sm" class="text-gray-400 group-hover:text-emerald-500" />
            </button>
          </div>
        </div>

        <!-- Charts Section -->
        <div class="space-y-6">
          <!-- Date Range Filter -->
          <div class="card admin-filter-bar p-4">
            <div class="admin-filter-bar-inner flex flex-wrap items-center gap-4">
              <div class="admin-filter-bar-group flex items-center gap-2">
                <span class="admin-filter-label text-sm font-medium"
                  >{{ t('admin.dashboard.timeRange') }}:</span
                >
                <DateRangePicker
                  v-model:start-date="startDate"
                  v-model:end-date="endDate"
                  @change="onDateRangeChange"
                />
              </div>
              <button @click="loadDashboardStats" :disabled="chartsLoading" class="btn btn-secondary admin-filter-refresh">
                {{ t('common.refresh') }}
              </button>
              <div class="admin-filter-bar-group ml-auto flex items-center gap-2">
                <span class="admin-filter-label text-sm font-medium"
                  >{{ t('admin.dashboard.granularity') }}:</span
                >
                <div class="w-28">
                  <Select
                    v-model="granularity"
                    :options="granularityOptions"
                    @change="loadChartData"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Charts Grid -->
          <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <ModelDistributionChart
              :model-stats="modelStats"
              :enable-ranking-view="true"
              :ranking-items="rankingItems"
              :ranking-total-actual-cost="rankingTotalActualCost"
              :ranking-total-requests="rankingTotalRequests"
              :ranking-total-tokens="rankingTotalTokens"
              :loading="chartsLoading"
              :ranking-loading="rankingLoading"
              :ranking-error="rankingError"
              :start-date="startDate"
              :end-date="endDate"
              @ranking-click="goToUserUsage"
            />
            <TokenUsageTrend :trend-data="trendData" :loading="chartsLoading" />
          </div>

          <!-- User Usage Trend (Full Width) -->
          <div class="card admin-trend-card p-4">
            <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.dashboard.recentUsage') }} (Top 12)
            </h3>
            <div class="admin-trend-stage h-64">
              <div v-if="userTrendLoading" class="flex h-full items-center justify-center">
                <LoadingSpinner size="md" />
              </div>
              <Line v-else-if="userTrendChartData" :data="userTrendChartData" :options="lineOptions" />
              <div
                v-else
                class="admin-trend-empty flex h-full items-center justify-center text-sm"
              >
                {{ t('admin.dashboard.noDataAvailable') }}
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { adminAPI } from '@/api/admin'
import type { OpsErrorLog } from '@/api/admin/ops'
import type {
  DashboardStats,
  TrendDataPoint,
  ModelStat,
  UserUsageTrendPoint,
  UserSpendingRankingItem
} from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Select from '@/components/common/Select.vue'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'
import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import { useAppStore } from '@/stores/app'
import { useBatchImageAccess } from '@/composables/useBatchImageAccess'

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { Line } from 'vue-chartjs'

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  Filler
)

const { t } = useI18n()
const appStore = useAppStore()
const router = useRouter()
const { canUseBatchImage, refreshBatchImageAccess } = useBatchImageAccess()
const stats = ref<DashboardStats | null>(null)
const loading = ref(false)
const chartsLoading = ref(false)
const userTrendLoading = ref(false)
const rankingLoading = ref(false)
const rankingError = ref(false)
const attributionLoading = ref(false)
const attributionError = ref(false)
const topOwner = ref('')
const topSource = ref('')
const topPlatform = ref('')
const topStatusCode = ref<number | null>(null)
const attributionErrorRate = ref(0)
const recentErrorItems = ref<OpsErrorLog[]>([])

// Chart data
const trendData = ref<TrendDataPoint[]>([])
const modelStats = ref<ModelStat[]>([])
const userTrend = ref<UserUsageTrendPoint[]>([])
const rankingItems = ref<UserSpendingRankingItem[]>([])
const rankingTotalActualCost = ref(0)
const rankingTotalRequests = ref(0)
const rankingTotalTokens = ref(0)
let chartLoadSeq = 0
let usersTrendLoadSeq = 0
let rankingLoadSeq = 0
const rankingLimit = 12
const riskListLimit = 5

type RiskTagTone = 'critical' | 'warning' | 'watch'

interface RiskUserItem {
  userId: number
  displayName: string
  email: string
  failedCount: number
  requests: number
  failureRate: number
  requestsPerHour: number
  requestDensityLabel: string
  summary: string
  tags: Array<{
    label: string
    tone: RiskTagTone
  }>
  score: number
}

// Helper function to format date in local timezone
const formatLocalDate = (date: Date): string => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const getLast24HoursRangeDates = (): { start: string; end: string } => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  return {
    start: formatLocalDate(start),
    end: formatLocalDate(end)
  }
}

const toOpsDateTime = (dateStr: string, endOfDay = false): string =>
  new Date(dateStr + (endOfDay ? 'T23:59:59.999' : 'T00:00:00')).toISOString()

// Date range
const granularity = ref<'day' | 'hour'>('hour')
const defaultRange = getLast24HoursRangeDates()
const startDate = ref(defaultRange.start)
const endDate = ref(defaultRange.end)
const opsStartTime = computed(() => toOpsDateTime(startDate.value))
const opsEndTime = computed(() => toOpsDateTime(endDate.value, true))

// Granularity options for Select component
const granularityOptions = computed(() => [
  { value: 'day', label: t('admin.dashboard.day') },
  { value: 'hour', label: t('admin.dashboard.hour') }
])

// Dark mode detection
const isDarkMode = computed(() => {
  return document.documentElement.classList.contains('dark')
})

// Chart colors
const chartColors = computed(() => ({
  text: isDarkMode.value ? '#e5e7eb' : '#374151',
  grid: isDarkMode.value ? '#374151' : '#e5e7eb'
}))

// Line chart options (for user trend chart)
const lineOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index' as const
  },
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: chartColors.value.text,
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 15,
        font: {
          size: 11
        }
      }
    },
    tooltip: {
      itemSort: (a: any, b: any) => {
        const aValue = typeof a?.raw === 'number' ? a.raw : Number(a?.parsed?.y ?? 0)
        const bValue = typeof b?.raw === 'number' ? b.raw : Number(b?.parsed?.y ?? 0)
        return bValue - aValue
      },
      callbacks: {
        label: (context: any) => {
          return `${context.dataset.label}: ${formatTokens(context.raw)}`
        }
      }
    }
  },
  scales: {
    x: {
      grid: {
        color: chartColors.value.grid
      },
      ticks: {
        color: chartColors.value.text,
        font: {
          size: 10
        }
      }
    },
    y: {
      grid: {
        color: chartColors.value.grid
      },
      ticks: {
        color: chartColors.value.text,
        font: {
          size: 10
        },
        callback: (value: string | number) => formatTokens(Number(value))
      }
    }
  }
}))

// User trend chart data
const userTrendChartData = computed(() => {
  if (!userTrend.value?.length) return null

  const getDisplayName = (point: UserUsageTrendPoint): string => {
    const username = point.username?.trim()
    if (username) {
      return username
    }

    const email = point.email?.trim()
    if (email) {
      return email
    }

    return t('admin.redeem.userPrefix', { id: point.user_id })
  }

  // Group by user_id to avoid merging different users with the same display name
  const userGroups = new Map<number, { name: string; data: Map<string, number> }>()
  const allDates = new Set<string>()

  userTrend.value.forEach((point) => {
    allDates.add(point.date)
    const key = point.user_id
    if (!userGroups.has(key)) {
      userGroups.set(key, { name: getDisplayName(point), data: new Map() })
    }
    userGroups.get(key)!.data.set(point.date, point.tokens)
  })

  const sortedDates = Array.from(allDates).sort()
  const colors = [
    '#3b82f6',
    '#10b981',
    '#f59e0b',
    '#ef4444',
    '#8b5cf6',
    '#ec4899',
    '#14b8a6',
    '#f97316',
    '#6366f1',
    '#84cc16',
    '#06b6d4',
    '#a855f7'
  ]

  const datasets = Array.from(userGroups.values()).map((group, idx) => ({
    label: group.name,
    data: sortedDates.map((date) => group.data.get(date) || 0),
    borderColor: colors[idx % colors.length],
    backgroundColor: `${colors[idx % colors.length]}20`,
    fill: false,
    tension: 0.3
  }))

  return {
    labels: sortedDates,
    datasets
  }
})

// Format helpers
const formatTokens = (value: number | undefined): string => {
  if (value === undefined || value === null) return '0'
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B`
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M`
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K`
  }
  return value.toLocaleString()
}

const formatNumber = (value: number): string => {
  return value.toLocaleString()
}

const formatCost = (value: number | undefined | null): string => {
  if (value === undefined || value === null) return '0.0000'
  if (value >= 1000) {
    return (value / 1000).toFixed(2) + 'K'
  } else if (value >= 1) {
    return value.toFixed(2)
  } else if (value >= 0.01) {
    return value.toFixed(3)
  }
  return value.toFixed(4)
}

const formatDuration = (ms: number): string => {
  if (ms >= 1000) {
    return `${(ms / 1000).toFixed(2)}s`
  }
  return `${Math.round(ms)}ms`
}

const formatPercent = (value: number): string => `${(value * 100).toFixed(1)}%`

const ownerLabelMap: Record<string, string> = {
  client: t('admin.dashboard.attribution.ownerClient'),
  provider: t('admin.dashboard.attribution.ownerProvider'),
  platform: t('admin.dashboard.attribution.ownerPlatform'),
}

const sourceLabelMap: Record<string, string> = {
  client_request: t('admin.dashboard.attribution.sourceClientRequest'),
  upstream_http: t('admin.dashboard.attribution.sourceUpstreamHttp'),
  gateway: t('admin.dashboard.attribution.sourceGateway'),
}

const topOwnerLabel = computed(() => ownerLabelMap[topOwner.value] || (topOwner.value || t('admin.dashboard.attribution.ownerNone')))
const topSourceLabel = computed(() => sourceLabelMap[topSource.value] || (topSource.value || t('admin.dashboard.attribution.sourceNone')))
const topPlatformLabel = computed(() => topPlatform.value || t('admin.dashboard.attribution.platformNone'))
const topStatusLabel = computed(() => topStatusCode.value != null ? String(topStatusCode.value) : t('admin.dashboard.attribution.statusNone'))

const ownerInsight = computed(() => {
  if (topOwner.value === 'client') return t('admin.dashboard.attribution.ownerInsightClient')
  if (topOwner.value === 'provider') return t('admin.dashboard.attribution.ownerInsightProvider')
  if (topOwner.value === 'platform') return t('admin.dashboard.attribution.ownerInsightPlatform')
  return t('admin.dashboard.attribution.ownerInsightNone')
})

const sourceInsight = computed(() => {
  if (topSource.value === 'client_request') return t('admin.dashboard.attribution.sourceInsightClientRequest')
  if (topSource.value === 'upstream_http') return t('admin.dashboard.attribution.sourceInsightUpstreamHttp')
  if (topSource.value === 'gateway') return t('admin.dashboard.attribution.sourceInsightGateway')
  return t('admin.dashboard.attribution.sourceInsightNone')
})

const platformInsight = computed(() => {
  if (!topPlatform.value) return t('admin.dashboard.attribution.platformInsightNone')
  return t('admin.dashboard.attribution.platformInsight', { platform: topPlatform.value })
})

const statusInsight = computed(() => {
  if (topStatusCode.value == null) return t('admin.dashboard.attribution.statusInsightNone')
  if (topStatusCode.value === 429) return t('admin.dashboard.attribution.statusInsight429')
  if (topStatusCode.value >= 500) return t('admin.dashboard.attribution.statusInsight5xx')
  if (topStatusCode.value >= 400) return t('admin.dashboard.attribution.statusInsight4xx')
  return t('admin.dashboard.attribution.statusInsightOther')
})

const attributionHeadline = computed(() => {
  if (attributionLoading.value) return t('admin.dashboard.attribution.loading')
  if (attributionError.value) return t('admin.dashboard.attribution.error')
  const rate = (attributionErrorRate.value * 100).toFixed(1)
  if (!topOwner.value && !topSource.value) return t('admin.dashboard.attribution.empty')
  return t('admin.dashboard.attribution.headline', {
    rate,
    owner: topOwnerLabel.value,
    source: topSourceLabel.value
  })
})

const attributionPrimaryLink = computed(() => ({
  path: '/admin/ops',
  query: {
    start_time: opsStartTime.value,
    end_time: opsEndTime.value,
    ...(topOwner.value ? { error_owner: topOwner.value } : {}),
    ...(topSource.value ? { error_source: topSource.value } : {}),
  }
}))

const attributionSecondaryLink = computed(() => ({
  path: '/admin/usage',
  query: {
    start_date: startDate.value,
    end_date: endDate.value,
    ...(topPlatform.value ? { platform: topPlatform.value } : {}),
  }
}))

const attributionPrimaryAction = computed(() => topOwner.value === 'client'
  ? t('admin.dashboard.attribution.primaryActionClient')
  : t('admin.dashboard.attribution.primaryActionOps'))
const attributionSecondaryAction = computed(() => topPlatform.value
  ? t('admin.dashboard.attribution.secondaryActionPlatform')
  : t('admin.dashboard.attribution.secondaryActionDefault'))

const riskListLoading = computed(() => rankingLoading.value || attributionLoading.value)

const rangeHours = computed(() => {
  const start = new Date(startDate.value)
  const end = new Date(endDate.value)
  const diff = end.getTime() - start.getTime()
  if (!Number.isFinite(diff) || diff <= 0) {
    return 24
  }
  return Math.max(1, diff / (1000 * 60 * 60))
})

const riskPrimaryLink = computed(() => ({
  path: '/admin/usage',
  query: {
    start_date: startDate.value,
    end_date: endDate.value,
  }
}))

const riskSecondaryLink = computed(() => ({
  path: '/admin/ops',
  query: {
    start_time: opsStartTime.value,
    end_time: opsEndTime.value,
  }
}))

const getRiskUsageLink = (userId: number) => ({
  path: '/admin/usage',
  query: {
    user_id: String(userId),
    start_date: startDate.value,
    end_date: endDate.value,
  }
})

const getRiskOpsLink = (userId: number) => ({
  path: '/admin/ops',
  query: {
    user_id: String(userId),
    start_time: opsStartTime.value,
    end_time: opsEndTime.value,
  }
})

const userNameMap = computed(() => {
  const map = new Map<number, string>()
  for (const point of userTrend.value) {
    const username = point.username?.trim()
    if (username && !map.has(point.user_id)) {
      map.set(point.user_id, username)
    }
  }
  return map
})

const highRiskUsers = computed<RiskUserItem[]>(() => {
  const userMap = new Map<number, { email: string; requests: number; failedCount: number }>()

  for (const item of rankingItems.value) {
    userMap.set(item.user_id, {
      email: item.email?.trim() || '',
      requests: item.requests || 0,
      failedCount: 0,
    })
  }

  for (const item of recentErrorItems.value) {
    if (!item.user_id) continue
    const current = userMap.get(item.user_id) || {
      email: item.user_email?.trim() || '',
      requests: 0,
      failedCount: 0,
    }
    if (!current.email && item.user_email?.trim()) {
      current.email = item.user_email.trim()
    }
    current.failedCount += 1
    userMap.set(item.user_id, current)
  }

  return Array.from(userMap.entries())
    .map(([userId, item]) => {
      const requests = item.requests || 0
      const failedCount = item.failedCount || 0
      const failureRate = requests > 0 ? failedCount / requests : (failedCount > 0 ? 1 : 0)
      const requestsPerHour = requests / rangeHours.value
      const displayName = userNameMap.value.get(userId) || item.email || t('admin.dashboard.risk.displayNameFallback', { id: userId })
      const tags: RiskUserItem['tags'] = []

      if (failedCount >= 8) {
        tags.push({ label: t('admin.dashboard.risk.tagFailedHigh'), tone: 'critical' })
      } else if (failedCount >= 4) {
        tags.push({ label: t('admin.dashboard.risk.tagFailedElevated'), tone: 'warning' })
      }

      if (requests >= 10 && failureRate >= 0.35) {
        tags.push({ label: t('admin.dashboard.risk.tagFailureRateHigh'), tone: 'critical' })
      } else if (requests >= 6 && failureRate >= 0.2) {
        tags.push({ label: t('admin.dashboard.risk.tagFailureRateRaised'), tone: 'warning' })
      }

      if (requestsPerHour >= 30) {
        tags.push({ label: t('admin.dashboard.risk.tagRequestDensityHigh'), tone: 'watch' })
      }

      if (!tags.length && failedCount > 0) {
        tags.push({ label: t('admin.dashboard.risk.tagWatch'), tone: 'watch' })
      }

      const score = failedCount * 4 + failureRate * 100 + requestsPerHour * 0.6
      const summaryParts = [
        t('admin.dashboard.risk.summaryFailedCount', { value: formatNumber(failedCount) }),
        t('admin.dashboard.risk.summaryFailureRate', { value: formatPercent(failureRate) }),
        t('admin.dashboard.risk.summaryRequests', { value: formatNumber(requests) })
      ]
      if (requestsPerHour >= 20) {
        summaryParts.push(t('admin.dashboard.risk.summaryDense'))
      } else if (failedCount >= 6) {
        summaryParts.push(t('admin.dashboard.risk.summarySupportHeavy'))
      }

      return {
        userId,
        displayName,
        email: item.email || t('admin.dashboard.risk.emailFallback', { id: userId }),
        failedCount,
        requests,
        failureRate,
        requestsPerHour,
        requestDensityLabel: `${requestsPerHour.toFixed(1)}/h`,
        summary: summaryParts.join('，'),
        tags,
        score,
      }
    })
    .filter((item) => item.failedCount > 0 || item.requestsPerHour >= 18)
    .sort((a, b) => b.score - a.score)
    .slice(0, riskListLimit)
})

const riskHeadline = computed(() => {
  if (riskListLoading.value) return t('admin.dashboard.risk.loading')
  if (rankingError.value && attributionError.value) return t('admin.dashboard.risk.errorBoth')
  if (rankingError.value) return t('admin.dashboard.risk.errorRanking')
  if (attributionError.value) return t('admin.dashboard.risk.errorAttribution')
  if (!highRiskUsers.value.length) return t('admin.dashboard.risk.empty')

  const topUser = highRiskUsers.value[0]
  const topTag = topUser.tags[0]?.label || t('admin.dashboard.risk.tagWatch')
  return t('admin.dashboard.risk.headline', {
    user: topUser.displayName,
    tag: topTag
  })
})

const morningSheetKicker = computed(() => rangeHours.value > 48
  ? t('admin.dashboard.morningSheet.weeklyKicker')
  : t('admin.dashboard.morningSheet.dailyKicker'))
const morningSheetTitle = computed(() => rangeHours.value > 48
  ? t('admin.dashboard.morningSheet.weeklyTitle')
  : t('admin.dashboard.morningSheet.dailyTitle'))
const morningSheetPrimaryLink = computed(() => briefingPrimaryLink.value)
const morningSheetSecondaryLink = computed(() => briefingSecondaryLink.value)
const morningSheetLead = computed(() => {
  if (!stats.value) return t('admin.dashboard.morningSheet.leadLoading')
  return rangeHours.value > 48
    ? t('admin.dashboard.morningSheet.leadWeekly')
    : t('admin.dashboard.morningSheet.leadDaily')
})
const morningSheetSummary = computed(() => {
  if (!stats.value) return t('admin.dashboard.morningSheet.summaryPending')
  const topRisk = highRiskUsers.value[0]?.displayName || t('admin.dashboard.morningSheet.summaryRiskNone')
  return t('admin.dashboard.morningSheet.summary', {
    start: startDate.value,
    end: endDate.value,
    requests: formatNumber(stats.value.total_requests),
    spend: formatCost(stats.value.total_actual_cost),
    owner: topOwnerLabel.value,
    platform: topPlatformLabel.value,
    risk: topRisk,
  })
})
const morningSheetWindowValue = computed(() => {
  if (!stats.value) return t('admin.dashboard.morningSheet.pending')
  return rangeHours.value > 48
    ? t('admin.dashboard.morningSheet.windowWeeklyValue', { value: formatNumber(stats.value.total_requests) })
    : t('admin.dashboard.morningSheet.windowDailyValue', { value: formatNumber(stats.value.today_requests) })
})
const morningSheetWindowNote = computed(() => {
  if (!stats.value) return t('admin.dashboard.morningSheet.windowNotePending')
  return t('admin.dashboard.morningSheet.windowNote', {
    active: formatNumber(stats.value.active_users),
    accounts: formatNumber(stats.value.error_accounts),
  })
})
const morningSheetAnomalyValue = computed(() => {
  if (attributionLoading.value) return t('admin.dashboard.morningSheet.pending')
  if (attributionError.value) return t('admin.dashboard.morningSheet.anomalyFallback')
  if (!topOwner.value) return t('admin.dashboard.morningSheet.anomalyNone')
  return t('admin.dashboard.morningSheet.anomalyValue', {
    owner: topOwnerLabel.value,
    source: topSourceLabel.value,
  })
})
const morningSheetAnomalyNote = computed(() => {
  if (attributionLoading.value) return t('admin.dashboard.morningSheet.anomalyNotePending')
  if (attributionError.value) return t('admin.dashboard.morningSheet.anomalyNoteError')
  return topOwner.value ? ownerInsight.value : t('admin.dashboard.morningSheet.anomalyNoteNone')
})
const morningSheetWatchValue = computed(() => {
  if (riskListLoading.value) return t('admin.dashboard.morningSheet.pending')
  if (highRiskUsers.value[0]) return highRiskUsers.value[0].displayName
  return t('admin.dashboard.morningSheet.watchNone')
})
const morningSheetWatchNote = computed(() => {
  if (riskListLoading.value) return t('admin.dashboard.morningSheet.watchNotePending')
  if (highRiskUsers.value[0]) return highRiskUsers.value[0].summary
  return t('admin.dashboard.morningSheet.watchNoteNone')
})

const morningSheetTopRiskUsers = computed(() => highRiskUsers.value.slice(0, 3))

const morningSheetAttributionLines = computed(() => {
  if (attributionLoading.value) {
    return [
      `${t('admin.dashboard.morningSheet.attributionPending')}`,
      `${t('admin.dashboard.morningSheet.lineRecommendation', { value: t('admin.dashboard.morningSheet.anomalyNotePending') })}`
    ]
  }

  if (attributionError.value) {
    return [
      `${t('admin.dashboard.morningSheet.attributionPending')}`,
      `${t('admin.dashboard.morningSheet.lineRecommendation', { value: t('admin.dashboard.morningSheet.anomalyNoteError') })}`
    ]
  }

  if (!topOwner.value && !topSource.value && !topPlatform.value && topStatusCode.value == null) {
    return [
      `${t('admin.dashboard.morningSheet.attributionEmpty')}`,
      `${t('admin.dashboard.morningSheet.lineRecommendation', { value: t('admin.dashboard.morningSheet.anomalyNoteNone') })}`
    ]
  }

  return [
    `${t('admin.dashboard.morningSheet.lineOwner', { value: topOwnerLabel.value })}`,
    `${t('admin.dashboard.morningSheet.lineSource', { value: topSourceLabel.value })}`,
    `${t('admin.dashboard.morningSheet.linePlatform', { value: topPlatformLabel.value })}`,
    `${t('admin.dashboard.morningSheet.lineStatus', { value: topStatusLabel.value })}`,
    `${t('admin.dashboard.morningSheet.lineRecommendation', { value: attributionHeadline.value })}`
  ]
})

const morningSheetRiskLines = computed(() => {
  if (riskListLoading.value) {
    return [
      `${t('admin.dashboard.morningSheet.riskPending')}`,
      `${t('admin.dashboard.morningSheet.lineRecommendation', { value: t('admin.dashboard.morningSheet.watchNotePending') })}`
    ]
  }

  if (!morningSheetTopRiskUsers.value.length) {
    return [
      `${t('admin.dashboard.morningSheet.riskEmpty')}`,
      `${t('admin.dashboard.morningSheet.lineRecommendation', { value: riskEmptyMessage.value })}`
    ]
  }

  return morningSheetTopRiskUsers.value.flatMap((user, index) => {
    const tags = user.tags.map((tag) => tag.label).join(' / ') || t('admin.dashboard.risk.tagWatch')
    return [
      `### ${index + 1}. ${user.displayName}`,
      `${t('auth.emailLabel')}: ${user.email}`,
      `${t('admin.dashboard.morningSheet.lineRiskTags', { value: tags })}`,
      `${t('admin.dashboard.morningSheet.lineSummary', { value: user.summary })}`,
    ]
  })
})

const morningSheetExportTitle = computed(() => rangeHours.value > 48
  ? t('admin.dashboard.morningSheet.exportWeeklyTitle')
  : t('admin.dashboard.morningSheet.exportDailyTitle'))

const morningSheetExportText = computed(() => {
  const sections = [
    `# ${morningSheetExportTitle.value}`,
    '',
    `- ${t('admin.dashboard.hero.statsRange')}: ${startDate.value} 至 ${endDate.value}`,
    '',
    morningSheetSummary.value,
    '',
    `## ${t('admin.dashboard.morningSheet.windowsTitle')}`,
    `- ${morningSheetWindowValue.value}`,
    `- ${morningSheetWindowNote.value}`,
    '',
    `## ${t('admin.dashboard.morningSheet.anomalyTitle')}`,
    `- ${morningSheetAnomalyValue.value}`,
    `- ${morningSheetAnomalyNote.value}`,
    '',
    `## ${t('admin.dashboard.morningSheet.watchTitle')}`,
    `- ${morningSheetWatchValue.value}`,
    `- ${morningSheetWatchNote.value}`,
    '',
    `## ${t('admin.dashboard.morningSheet.attributionSectionTitle')}`,
    ...morningSheetAttributionLines.value,
    '',
    `## ${t('admin.dashboard.morningSheet.riskSectionTitle')}`,
    ...morningSheetRiskLines.value,
  ]
  return sections.join('\n')
})

async function copyMorningSheet() {
  try {
    await navigator.clipboard.writeText(morningSheetExportText.value)
    appStore.showSuccess(t('admin.dashboard.morningSheet.copySuccess'))
  } catch (error) {
    console.error('Failed to copy morning sheet:', error)
    appStore.showError(t('admin.dashboard.morningSheet.copyFailed'))
  }
}

function downloadMorningSheet() {
  try {
    const blob = new Blob([morningSheetExportText.value], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const safeDate = `${startDate.value}_${endDate.value}`
    link.href = url
    link.download = `sst-admin-brief-${safeDate}.md`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    appStore.showSuccess(t('admin.dashboard.morningSheet.downloadSuccess'))
  } catch (error) {
    console.error('Failed to download morning sheet:', error)
    appStore.showError(t('admin.dashboard.morningSheet.downloadFailed'))
  }
}

const briefingKicker = computed(() => rangeHours.value > 48
  ? t('admin.dashboard.briefing.weeklyKicker')
  : t('admin.dashboard.briefing.dailyKicker'))
const briefingTitle = computed(() => rangeHours.value > 48
  ? t('admin.dashboard.briefing.weeklyTitle')
  : t('admin.dashboard.briefing.dailyTitle'))

const briefingPrimaryLink = computed(() => ({
  path: '/admin/ops',
  query: {
    start_time: opsStartTime.value,
    end_time: opsEndTime.value,
    ...(topPlatform.value ? { platform: topPlatform.value } : {})
  }
}))

const briefingSecondaryLink = computed(() => ({
  path: '/admin/usage',
  query: {
    start_date: startDate.value,
    end_date: endDate.value,
    ...(highRiskUsers.value[0] ? { user_id: String(highRiskUsers.value[0].userId) } : {})
  }
}))

const briefingLead = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.leadLoading')
  const modeLabel = rangeHours.value > 48
    ? t('admin.dashboard.briefing.leadWeekly')
    : t('admin.dashboard.briefing.leadDaily')
  return t('admin.dashboard.briefing.leadPrefix') + modeLabel
})

const briefingTrafficValue = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.trafficPending')
  return formatNumber(stats.value.today_requests) + ' 次请求'
})

const briefingTrafficNote = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.trafficNotePending')
  return t('admin.dashboard.briefing.trafficNote', {
    active: formatNumber(stats.value.active_users),
    newUsers: formatNumber(stats.value.today_new_users)
  })
})

const briefingCostValue = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.costPending')
  return '$' + formatCost(stats.value.today_actual_cost)
})

const briefingCostNote = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.costNotePending')
  return t('admin.dashboard.briefing.costNote', {
    tokens: formatTokens(stats.value.today_tokens),
    accountCost: formatCost(stats.value.today_account_cost)
  })
})

const briefingAnomalyValue = computed(() => {
  if (attributionLoading.value) return t('admin.dashboard.briefing.anomalyPending')
  if (attributionError.value) return t('admin.dashboard.briefing.anomalyError')
  if (!topOwner.value) return t('admin.dashboard.briefing.anomalyNone')
  return topOwnerLabel.value + ' / ' + topSourceLabel.value
})

const briefingAnomalyNote = computed(() => {
  if (attributionLoading.value) return t('admin.dashboard.briefing.anomalyNotePending')
  if (attributionError.value) return t('admin.dashboard.briefing.anomalyErrorNote')
  const rate = (attributionErrorRate.value * 100).toFixed(1)
  if (!topOwner.value) return t('admin.dashboard.briefing.anomalyNoteNone')
  return t('admin.dashboard.briefing.anomalyNote', {
    rate,
    platform: topPlatformLabel.value,
    status: topStatusLabel.value
  })
})

const briefingActionValue = computed(() => {
  if (riskListLoading.value) return t('admin.dashboard.briefing.actionPending')
  if (rankingError.value) return t('admin.dashboard.briefing.actionRetry')
  if (!highRiskUsers.value.length) return t('admin.dashboard.briefing.actionObserve')
  return '先看 ' + highRiskUsers.value[0].displayName
})

const briefingActionNote = computed(() => {
  if (riskListLoading.value) return t('admin.dashboard.briefing.actionNotePending')
  if (rankingError.value) return t('admin.dashboard.briefing.actionNoteError')
  if (!highRiskUsers.value.length) return t('admin.dashboard.briefing.actionNoteNone')
  return highRiskUsers.value[0].summary
})

const briefingTags = computed(() => {
  const tags: string[] = []
  if (stats.value) {
    if (stats.value.today_new_users > 0) tags.push(t('admin.dashboard.briefing.tagNewUsers', { value: formatNumber(stats.value.today_new_users) }))
    if (stats.value.error_accounts > 0) tags.push(t('admin.dashboard.briefing.tagErrorAccounts', { value: formatNumber(stats.value.error_accounts) }))
    if (stats.value.rpm > 0) tags.push(t('admin.dashboard.briefing.tagRpm', { value: formatNumber(Math.round(stats.value.rpm)) }))
  }
  if (topPlatform.value) tags.push(t('admin.dashboard.briefing.tagPlatform', { value: topPlatform.value }))
  if (topStatusCode.value != null) tags.push(t('admin.dashboard.briefing.tagStatus', { value: topStatusLabel.value }))
  if (highRiskUsers.value[0]?.tags[0]?.label) tags.push(highRiskUsers.value[0].tags[0].label)
  return tags.slice(0, 5)
})

const briefingSummary = computed(() => {
  if (!stats.value) return t('admin.dashboard.briefing.summaryPending')
  const parts = [
    t('admin.dashboard.briefing.summaryRequests', {
      start: startDate.value,
      end: endDate.value,
      value: formatNumber(stats.value.total_requests)
    }),
    t('admin.dashboard.briefing.summarySpend', { value: formatCost(stats.value.total_actual_cost) }),
  ]
  if (topOwner.value) parts.push(t('admin.dashboard.briefing.summaryOwner', { owner: topOwnerLabel.value }))
  else if (attributionError.value) parts.push(t('admin.dashboard.briefing.summaryAttributionPending'))
  if (highRiskUsers.value[0]) parts.push(t('admin.dashboard.briefing.summaryRisk', { user: highRiskUsers.value[0].displayName }))
  else if (rankingError.value) parts.push(t('admin.dashboard.briefing.summaryRankingPending'))
  return parts.join('，') + '。'
})

const riskEmptyMessage = computed(() => {
  if (rankingError.value && attributionError.value) return t('admin.dashboard.risk.errorBoth')
  if (rankingError.value) return t('admin.dashboard.risk.errorRanking')
  if (attributionError.value) return t('admin.dashboard.risk.errorAttribution')
  return t('admin.dashboard.risk.empty')
})

const goToUserUsage = (item: UserSpendingRankingItem) => {
  void router.push({
    path: '/admin/usage',
    query: {
      user_id: String(item.user_id),
      start_date: startDate.value,
      end_date: endDate.value
    }
  })
}

// Date range change handler
const onDateRangeChange = (range: {
  startDate: string
  endDate: string
  preset: string | null
}) => {
  // Auto-select granularity based on date range
  const start = new Date(range.startDate)
  const end = new Date(range.endDate)
  const daysDiff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))

  // If range is 1 day, use hourly granularity
  if (daysDiff <= 1) {
    granularity.value = 'hour'
  } else {
    granularity.value = 'day'
  }

  loadChartData()
}

// Load data
const loadDashboardSnapshot = async (includeStats: boolean) => {
  const currentSeq = ++chartLoadSeq
  if (includeStats && !stats.value) {
    loading.value = true
  }
  chartsLoading.value = true
  try {
    const response = await adminAPI.dashboard.getSnapshotV2({
      start_date: startDate.value,
      end_date: endDate.value,
      granularity: granularity.value,
      include_stats: includeStats,
      include_trend: true,
      include_model_stats: true,
      include_group_stats: false,
      include_users_trend: false
    })
    if (currentSeq !== chartLoadSeq) return
    if (includeStats && response.stats) {
      stats.value = response.stats
    }
    trendData.value = response.trend || []
    modelStats.value = response.models || []
  } catch (error) {
    if (currentSeq !== chartLoadSeq) return
    appStore.showError(t('admin.dashboard.failedToLoad'))
    console.error('Error loading dashboard snapshot:', error)
  } finally {
    if (currentSeq === chartLoadSeq) {
      loading.value = false
      chartsLoading.value = false
    }
  }
}

const loadUsersTrend = async () => {
  const currentSeq = ++usersTrendLoadSeq
  userTrendLoading.value = true
  try {
    const response = await adminAPI.dashboard.getUserUsageTrend({
      start_date: startDate.value,
      end_date: endDate.value,
      granularity: granularity.value,
      limit: 12
    })
    if (currentSeq !== usersTrendLoadSeq) return
    userTrend.value = response.trend || []
  } catch (error) {
    if (currentSeq !== usersTrendLoadSeq) return
    console.error('Error loading users trend:', error)
    userTrend.value = []
  } finally {
    if (currentSeq === usersTrendLoadSeq) {
      userTrendLoading.value = false
    }
  }
}

const loadAttributionWorkbench = async () => {
  attributionLoading.value = true
  try {
    const [overview, errors] = await Promise.all([
      adminAPI.ops.getDashboardOverview({
        start_time: opsStartTime.value,
        end_time: opsEndTime.value,
      }),
      adminAPI.ops.listErrorLogs({
        start_time: opsStartTime.value,
        end_time: opsEndTime.value,
        page: 1,
        page_size: 50,
      }),
    ])

    attributionErrorRate.value = overview.error_rate || 0

    const ownerCounts = new Map<string, number>()
    const sourceCounts = new Map<string, number>()
    const platformCounts = new Map<string, number>()
    const statusCounts = new Map<number, number>()

    for (const item of errors.items || []) {
      if (item.error_owner) ownerCounts.set(item.error_owner, (ownerCounts.get(item.error_owner) || 0) + 1)
      if (item.error_source) sourceCounts.set(item.error_source, (sourceCounts.get(item.error_source) || 0) + 1)
      if (item.platform) platformCounts.set(item.platform, (platformCounts.get(item.platform) || 0) + 1)
      if (typeof item.status_code === 'number') statusCounts.set(item.status_code, (statusCounts.get(item.status_code) || 0) + 1)
    }

    recentErrorItems.value = errors.items || []
    topOwner.value = pickTopKey(ownerCounts)
    topSource.value = pickTopKey(sourceCounts)
    topPlatform.value = pickTopKey(platformCounts)
    topStatusCode.value = pickTopNumericKey(statusCounts)
  } catch (error) {
    console.error('Error loading attribution workbench:', error)
    topOwner.value = ''
    topSource.value = ''
    topPlatform.value = ''
    topStatusCode.value = null
    attributionErrorRate.value = 0
    recentErrorItems.value = []
  } finally {
    attributionLoading.value = false
  }
}

const loadUserSpendingRanking = async () => {
  const currentSeq = ++rankingLoadSeq
  rankingLoading.value = true
  rankingError.value = false
  try {
    const response = await adminAPI.dashboard.getUserSpendingRanking({
      start_date: startDate.value,
      end_date: endDate.value,
      limit: rankingLimit
    })
    if (currentSeq !== rankingLoadSeq) return
    rankingItems.value = response.ranking || []
    rankingTotalActualCost.value = response.total_actual_cost || 0
    rankingTotalRequests.value = response.total_requests || 0
    rankingTotalTokens.value = response.total_tokens || 0
  } catch (error) {
    if (currentSeq !== rankingLoadSeq) return
    console.error('Error loading user spending ranking:', error)
    rankingItems.value = []
    rankingTotalActualCost.value = 0
    rankingTotalRequests.value = 0
    rankingTotalTokens.value = 0
    rankingError.value = true
  } finally {
    if (currentSeq === rankingLoadSeq) {
      rankingLoading.value = false
    }
  }
}

function pickTopKey(map: Map<string, number>) {
  let bestKey = ''
  let bestValue = -1
  for (const [key, value] of map.entries()) {
    if (value > bestValue) {
      bestKey = key
      bestValue = value
    }
  }
  return bestKey
}

function pickTopNumericKey(map: Map<number, number>) {
  let bestKey: number | null = null
  let bestValue = -1
  for (const [key, value] of map.entries()) {
    if (value > bestValue) {
      bestKey = key
      bestValue = value
    }
  }
  return bestKey
}

const loadDashboardStats = async () => {
  await Promise.all([
    loadDashboardSnapshot(true),
    loadUsersTrend(),
    loadUserSpendingRanking(),
    loadAttributionWorkbench()
  ])
}

const loadChartData = async () => {
  await Promise.all([
    loadDashboardSnapshot(false),
    loadUsersTrend(),
    loadUserSpendingRanking(),
    loadAttributionWorkbench()
  ])
}

onMounted(() => {
  void refreshBatchImageAccess()
  loadDashboardStats()
})
</script>

<style scoped>
.admin-dashboard-shell {
  padding: clamp(0.2rem, 0.6vw, 0.4rem);
}

.admin-page-hero {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.25rem 1.3rem 1.1rem;
  border: 1px solid rgba(198, 184, 157, 0.52);
  border-radius: 12px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.055), transparent 30%),
    linear-gradient(180deg, rgba(255, 252, 245, 0.9), rgba(246, 241, 231, 0.78));
}

.admin-page-kicker,
.admin-page-meta dt {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-page-hero h2 {
  margin-top: 0.45rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.45rem, 2vw, 1.9rem);
  font-weight: 600;
  line-height: 1.15;
}

.admin-page-hero p {
  max-width: 38rem;
  margin-top: 0.55rem;
  color: #5f675d;
  font-size: 0.95rem;
  line-height: 1.7;
}

.admin-page-meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(10rem, 1fr));
  gap: 0.75rem;
  min-width: min(100%, 24rem);
}

.admin-page-meta div {
  padding: 0.85rem 0.95rem;
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.7);
}

.admin-page-meta dd {
  margin-top: 0.32rem;
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 600;
}

.admin-stat-card {
  position: relative;
  overflow: hidden;
}

.admin-stat-card::after {
  content: '';
  position: absolute;
  inset: 0 auto 0 0;
  width: 3px;
  background: linear-gradient(180deg, rgba(167, 58, 42, 0.58), rgba(155, 129, 85, 0.12));
  opacity: 0.7;
}

.admin-attribution-card {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 30%),
    rgba(250, 247, 239, 0.76);
}

.admin-briefing-card {
  background:
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.1), transparent 30%),
    linear-gradient(135deg, rgba(255, 252, 246, 0.94), rgba(244, 238, 227, 0.88));
}

.admin-risk-card {
  background:
    linear-gradient(90deg, rgba(130, 95, 44, 0.05), transparent 34%),
    rgba(252, 248, 240, 0.8);
}

.admin-morning-sheet-card {
  background:
    linear-gradient(140deg, rgba(255, 251, 243, 0.96), rgba(245, 238, 226, 0.92)),
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.08), transparent 30%);
}

.admin-attribution-head {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
}

.admin-briefing-head {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
}

.admin-morning-sheet-head {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
}

.admin-attribution-head h3,
.admin-briefing-head h3,
.admin-morning-sheet-head h3 {
  margin-top: 0.45rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.35;
}

.admin-attribution-head p,
.admin-briefing-head p,
.admin-morning-sheet-head p {
  max-width: 44rem;
  margin-top: 0.45rem;
  color: #5f675d;
  font-size: 0.88rem;
  line-height: 1.7;
}

.admin-stat-card p.text-xl {
  font-size: clamp(0.96rem, 1.08vw, 1.06rem);
  line-height: 1.15;
  letter-spacing: -0.01em;
}

.admin-stat-card p.text-sm.font-semibold {
  font-size: 0.84rem;
  line-height: 1.35;
}

.admin-stat-card p.text-xs,
.admin-stat-card span.text-xs {
  line-height: 1.45;
}

.admin-attribution-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
}

.admin-briefing-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.9rem;
}

.admin-briefing-panel {
  display: grid;
  gap: 0.4rem;
  border: 1px solid rgba(198, 184, 157, 0.38);
  border-radius: 12px;
  background: rgba(255, 252, 245, 0.82);
  padding: 1rem;
}

.admin-briefing-panel span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-briefing-panel strong {
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.35;
}

.admin-briefing-panel p,
.admin-briefing-summary {
  color: #5f675d;
  font-size: 0.84rem;
  line-height: 1.7;
}

.admin-briefing-foot {
  display: flex;
  flex-wrap: wrap;
  align-items: start;
  justify-content: space-between;
  gap: 0.9rem;
}

.admin-briefing-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
}

.admin-briefing-tag {
  display: inline-flex;
  align-items: center;
  padding: 0.34rem 0.7rem;
  border-radius: 999px;
  border: 1px solid rgba(167, 58, 42, 0.16);
  color: #7d4d3d;
  background: rgba(167, 58, 42, 0.06);
  font-size: 0.76rem;
  line-height: 1;
}

.admin-morning-sheet-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 1.8fr);
  gap: 1rem;
}

.admin-morning-sheet-summary {
  padding: 1.1rem 1.15rem;
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 16px;
  background: rgba(255, 252, 245, 0.86);
}

.admin-morning-sheet-summary p {
  color: #433d33;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.9rem;
  line-height: 1.9;
}

.admin-morning-sheet-columns {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.85rem;
}

.admin-morning-sheet-panel {
  display: grid;
  gap: 0.45rem;
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 14px;
  background: rgba(252, 248, 239, 0.82);
  padding: 1rem;
}

.admin-morning-sheet-panel span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-morning-sheet-panel strong {
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.35;
}

.admin-morning-sheet-panel p {
  color: #5f675d;
  font-size: 0.84rem;
  line-height: 1.7;
}

.admin-risk-head {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
}

.admin-risk-head h3 {
  margin-top: 0.45rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.35;
}

.admin-risk-head p {
  max-width: 44rem;
  margin-top: 0.45rem;
  color: #5f675d;
  font-size: 0.88rem;
  line-height: 1.7;
}

.admin-attribution-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.9rem;
}

.admin-attribution-panel {
  display: grid;
  gap: 0.4rem;
  border: 1px solid rgba(198, 184, 157, 0.38);
  border-radius: 12px;
  background: rgba(255, 252, 245, 0.75);
  padding: 1rem;
}

.admin-attribution-panel span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-attribution-panel strong {
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.35;
}

.admin-attribution-panel p {
  color: #5f675d;
  font-size: 0.84rem;
  line-height: 1.65;
}

.admin-risk-list {
  display: grid;
  gap: 0.85rem;
}

.admin-risk-item {
  display: grid;
  gap: 0.9rem;
  padding: 1rem;
  border: 1px solid rgba(198, 184, 157, 0.38);
  border-radius: 12px;
  background: rgba(255, 252, 245, 0.78);
}

.admin-risk-item-main {
  display: flex;
  flex-wrap: wrap;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
}

.admin-risk-item-copy {
  display: grid;
  gap: 0.45rem;
  min-width: min(100%, 18rem);
}

.admin-risk-item-title {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.6rem;
}

.admin-risk-item-title strong {
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.35;
}

.admin-risk-item-title span,
.admin-risk-item-copy p,
.admin-risk-empty {
  color: #5f675d;
  font-size: 0.88rem;
  line-height: 1.7;
}

.admin-risk-item-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(5.4rem, 1fr));
  gap: 0.75rem;
  min-width: min(100%, 20rem);
}

.admin-risk-item-metrics div {
  padding: 0.8rem 0.9rem;
  border: 1px solid rgba(198, 184, 157, 0.32);
  border-radius: 10px;
  background: rgba(249, 245, 236, 0.92);
}

.admin-risk-item-metrics span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.admin-risk-item-metrics strong {
  display: block;
  margin-top: 0.38rem;
  color: #1f2320;
  font-size: 0.9rem;
  font-weight: 700;
  line-height: 1.3;
}

.admin-risk-item-foot {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
}

.admin-risk-tags,
.admin-risk-links {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
}

.admin-risk-tag {
  display: inline-flex;
  align-items: center;
  padding: 0.34rem 0.7rem;
  border-radius: 999px;
  border: 1px solid rgba(198, 184, 157, 0.34);
  color: #6a5a45;
  background: rgba(248, 242, 230, 0.94);
  font-size: 0.76rem;
  line-height: 1;
}

.admin-risk-tag[data-tone='critical'] {
  border-color: rgba(167, 58, 42, 0.32);
  color: #8d3529;
  background: rgba(167, 58, 42, 0.08);
}

.admin-risk-tag[data-tone='warning'] {
  border-color: rgba(174, 120, 41, 0.28);
  color: #8b6427;
  background: rgba(191, 138, 52, 0.08);
}

.admin-risk-tag[data-tone='watch'] {
  border-color: rgba(117, 107, 78, 0.22);
  color: #6c624e;
  background: rgba(123, 106, 83, 0.08);
}

.admin-filter-bar,
.admin-trend-card {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.04), transparent 28%),
    rgba(250, 247, 239, 0.76);
}

.admin-filter-bar-inner {
  min-width: 0;
}

.admin-filter-bar-group {
  min-width: 0;
}

.admin-filter-label {
  color: #59645a;
}

.admin-filter-refresh {
  box-shadow: none;
}

.admin-trend-stage {
  display: grid;
  min-width: 0;
}

.admin-trend-empty {
  color: #6f7a70;
}

:deep(.admin-stat-card .rounded-lg) {
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 0.7rem;
}

@media (max-width: 980px) {
  .admin-attribution-grid,
  .admin-briefing-grid,
  .admin-morning-sheet-columns {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .admin-morning-sheet-layout {
    grid-template-columns: 1fr;
  }

  .admin-risk-item-metrics {
    grid-template-columns: repeat(3, minmax(0, 1fr));
    width: 100%;
  }
}

@media (max-width: 780px) {
  .admin-page-meta,
  .admin-attribution-grid,
  .admin-briefing-grid,
  .admin-morning-sheet-columns {
    grid-template-columns: 1fr;
    width: 100%;
  }

  .admin-attribution-actions,
  .admin-briefing-foot {
    width: 100%;
  }

  .admin-risk-item-metrics {
    grid-template-columns: 1fr;
  }

  .admin-risk-links {
    width: 100%;
  }
}
</style>

<style>
.dark .admin-page-hero {
  border-color: rgba(48, 52, 43, 0.92);
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 32%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.92), rgba(17, 19, 15, 0.84));
}

.dark .admin-page-hero h2,
.dark .admin-page-meta dd {
  color: #f4efe4;
}

.dark .admin-page-kicker,
.dark .admin-page-meta dt,
.dark .admin-briefing-panel span,
.dark .admin-attribution-panel span,
.dark .admin-morning-sheet-panel span,
.dark .admin-risk-item-metrics span {
  color: #a89b84;
}

.dark .admin-page-hero p {
  color: #bdb5a8;
}

.dark .admin-page-meta div {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.68);
}

.dark .admin-stat-card {
  border-color: rgba(48, 52, 43, 0.88);
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.88), rgba(18, 20, 16, 0.8)),
    rgba(24, 26, 21, 0.76);
}

.dark .admin-stat-card .text-gray-900,
.dark .admin-stat-card .dark\:text-white,
.dark .admin-stat-card .text-xl,
.dark .admin-stat-card .font-bold {
  color: #f4efe4 !important;
}

.dark .admin-stat-card .text-gray-500,
.dark .admin-stat-card .text-gray-400,
.dark .admin-stat-card .dark\:text-gray-400,
.dark .admin-stat-card .dark\:text-gray-500 {
  color: #bdb5a8 !important;
}

.dark .admin-attribution-card,
.dark .admin-risk-card,
.dark .admin-filter-bar,
.dark .admin-trend-card {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 30%),
    rgba(24, 26, 21, 0.78);
}

.dark .admin-filter-label {
  color: #d9d0be;
}

.dark .admin-trend-empty {
  color: #9f9787;
}

.dark .admin-morning-sheet-card {
  background:
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.12), transparent 30%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.88), rgba(18, 20, 16, 0.84));
}

.dark .admin-briefing-card {
  background:
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.14), transparent 30%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.9), rgba(18, 20, 16, 0.86));
}

.dark .admin-attribution-head h3,
.dark .admin-risk-head h3,
.dark .admin-briefing-head h3,
.dark .admin-morning-sheet-head h3,
.dark .admin-attribution-panel strong,
.dark .admin-briefing-panel strong,
.dark .admin-morning-sheet-panel strong,
.dark .admin-morning-sheet-summary p {
  color: #f4efe4;
}

.dark .admin-attribution-head p,
.dark .admin-risk-head p,
.dark .admin-briefing-head p,
.dark .admin-morning-sheet-head p,
.dark .admin-attribution-panel p,
.dark .admin-briefing-panel p,
.dark .admin-briefing-summary,
.dark .admin-morning-sheet-panel p {
  color: #bdb5a8;
}

.dark .admin-attribution-panel,
.dark .admin-briefing-panel,
.dark .admin-morning-sheet-panel,
.dark .admin-morning-sheet-summary {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.72);
}

.dark .admin-risk-item {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.72);
}

.dark .admin-risk-item-title strong,
.dark .admin-risk-item-metrics strong {
  color: #f4efe4;
}

.dark .admin-risk-item-title span,
.dark .admin-risk-item-copy p,
.dark .admin-risk-empty {
  color: #bdb5a8;
}

.dark .admin-risk-item-metrics div {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(19, 22, 17, 0.84);
}

.dark .admin-risk-tag {
  border-color: rgba(80, 74, 58, 0.62);
  color: #d8ccb8;
  background: rgba(40, 33, 25, 0.72);
}

.dark .admin-risk-tag[data-tone='critical'] {
  border-color: rgba(167, 58, 42, 0.45);
  color: #f0b2a9;
  background: rgba(167, 58, 42, 0.14);
}

.dark .admin-risk-tag[data-tone='warning'] {
  border-color: rgba(184, 141, 72, 0.36);
  color: #e4c78f;
  background: rgba(160, 115, 44, 0.14);
}

.dark .admin-risk-tag[data-tone='watch'] {
  border-color: rgba(112, 103, 76, 0.42);
  color: #d2c4aa;
  background: rgba(92, 83, 61, 0.16);
}

.dark .admin-briefing-tag {
  border-color: rgba(167, 58, 42, 0.28);
  color: #f0c1b4;
  background: rgba(167, 58, 42, 0.14);
}
</style>


