<template>
  <AppLayout>
    <div class="admin-usage-shell space-y-6">

      <section class="card admin-usage-briefing-card p-4">
        <div class="admin-usage-briefing-head">
          <div>
            <span class="admin-page-kicker">{{ t('usage.adminLedger.briefing.kicker') }}</span>
            <h3>{{ t('usage.adminLedger.briefing.title') }}</h3>
            <p>{{ usageBriefingLead }}</p>
          </div>
          <div class="admin-usage-briefing-actions">
            <button class="btn btn-secondary" @click="switchToErrorsTab">{{ t('usage.adminLedger.briefing.primaryAction') }}</button>
            <button class="btn btn-secondary" @click="scrollToCharts">{{ t('usage.adminLedger.briefing.secondaryAction') }}</button>
          </div>
        </div>

        <div class="admin-usage-briefing-grid mt-3">
          <article class="admin-usage-briefing-panel">
            <span>{{ t('usage.adminLedger.briefing.requestsTitle') }}</span>
            <strong>{{ usageBriefingRequestsValue }}</strong>
            <p>{{ usageBriefingRequestsNote }}</p>
          </article>
          <article class="admin-usage-briefing-panel">
            <span>{{ t('usage.adminLedger.briefing.costTitle') }}</span>
            <strong>{{ usageBriefingCostValue }}</strong>
            <p>{{ usageBriefingCostNote }}</p>
          </article>
          <article class="admin-usage-briefing-panel">
            <span>{{ t('usage.adminLedger.briefing.anomalyTitle') }}</span>
            <strong>{{ usageBriefingAnomalyValue }}</strong>
            <p>{{ usageBriefingAnomalyNote }}</p>
          </article>
          <article class="admin-usage-briefing-panel">
            <span>{{ t('usage.adminLedger.briefing.actionTitle') }}</span>
            <strong>{{ usageBriefingActionValue }}</strong>
            <p>{{ usageBriefingActionNote }}</p>
          </article>
        </div>

        <div class="admin-usage-briefing-foot mt-3">
          <div class="admin-usage-briefing-tags">
            <span v-for="tag in usageBriefingTags" :key="tag" class="admin-usage-briefing-tag">{{ tag }}</span>
            <span v-if="!usageBriefingTags.length" class="admin-usage-briefing-tag admin-usage-briefing-tag-muted">{{ t('usage.adminLedger.briefing.tagObserve') }}</span>
          </div>
          <p class="admin-usage-briefing-summary">{{ usageBriefingSummary }}</p>
        </div>
      </section>

      <section class="admin-usage-summary-stack">
        <UsageStatsCards :stats="usageStats" />
      </section>
      <!-- Charts Section -->
      <div ref="modelChartSectionRef" class="space-y-4">
        <div class="card p-4">
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.timeRange') }}:</span>
              <DateRangePicker
                v-model:start-date="startDate"
                v-model:end-date="endDate"
                @change="onDateRangeChange"
              />
            </div>
            <div class="ml-auto flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.granularity') }}:</span>
              <div class="w-28">
                <Select v-model="granularity" :options="granularityOptions" @change="loadChartData" />
              </div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ModelDistributionChart
            v-model:source="modelDistributionSource"
            v-model:metric="modelDistributionMetric"
            :model-stats="requestedModelStats"
            :upstream-model-stats="upstreamModelStats"
            :mapping-model-stats="mappingModelStats"
            :loading="modelStatsLoading"
            :show-source-toggle="true"
            :show-metric-toggle="true"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
          <GroupDistributionChart
            v-model:metric="groupDistributionMetric"
            :group-stats="groupStats"
            :loading="chartsLoading"
            :show-metric-toggle="true"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
        </div>
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <EndpointDistributionChart
            v-model:source="endpointDistributionSource"
            v-model:metric="endpointDistributionMetric"
            :endpoint-stats="inboundEndpointStats"
            :upstream-endpoint-stats="upstreamEndpointStats"
            :endpoint-path-stats="endpointPathStats"
            :loading="endpointStatsLoading"
            :show-source-toggle="true"
            :show-metric-toggle="true"
            :title="t('usage.endpointDistribution')"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
          <TokenUsageTrend :trend-data="trendData" :loading="chartsLoading" />
        </div>
      </div>
      <UsageFilters v-model="filters" :start-date="startDate" :end-date="endDate" :exporting="exporting" :model-options="modelNameOptions" @change="applyFilters" @refresh="refreshData" @reset="resetFilters" @cleanup="openCleanupDialog" @export="exportToExcel">
        <template #after-reset>
          <div class="relative" ref="columnDropdownRef">
            <button
              @click="showColumnDropdown = !showColumnDropdown"
              class="btn btn-secondary px-2 md:px-3"
              :title="t('admin.users.columnSettings')"
            >
              <svg class="h-4 w-4 md:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 4.5v15m6-15v15m-10.875 0h15.75c.621 0 1.125-.504 1.125-1.125V5.625c0-.621-.504-1.125-1.125-1.125H4.125C3.504 4.5 3 5.004 3 5.625v12.75c0 .621.504 1.125 1.125 1.125z" />
              </svg>
              <span class="hidden md:inline">{{ t('admin.users.columnSettings') }}</span>
            </button>
            <div
              v-if="showColumnDropdown"
              class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
            >
              <button
                v-for="col in toggleableColumns"
                :key="col.key"
                @click="toggleColumn(col.key)"
                class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
              >
                <span>{{ col.label }}</span>
                <Icon
                  v-if="isColumnVisible(col.key)"
                  name="check"
                  size="sm"
                  class="text-primary-500"
                  :stroke-width="2"
                />
              </button>
            </div>
          </div>
        </template>
      </UsageFilters>
      <div class="mb-4 flex gap-2 border-b border-gray-200 dark:border-dark-700">
        <button class="tab" :class="{ 'tab-active': activeTab === 'usage' }" @click="activeTab = 'usage'">
          {{ t('usage.tabs.usage') }}
        </button>
        <button class="tab" :class="{ 'tab-active': activeTab === 'errors' }" @click="switchToErrorsTab">
          {{ t('usage.tabs.errors') }}
        </button>
      </div>
      <div v-show="activeTab === 'usage'" class="usage-ledger-panel">
        <UsageTable
          :data="usageLogs"
          :loading="loading"
          :columns="visibleColumns"
          :server-side-sort="true"
          :default-sort-key="'created_at'"
          :default-sort-order="'desc'"
          @sort="handleSort"
          @userClick="handleUserClick"
        />
        <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="handlePageChange" @update:pageSize="handlePageSizeChange" />
      </div>
      <div v-show="activeTab === 'errors'">
        <OpsErrorLogTable
          :rows="errRows" :total="errTotal" :loading="errLoading"
          :page="errPage" :page-size="errPageSize"
          @openErrorDetail="openError"
          @update:page="onErrPage"
          @update:pageSize="onErrPageSize" />
        <OpsErrorDetailModal v-model:show="showErrorModal" :error-id="selectedErrorId" :error-type="'request'" />
      </div>
    </div>
  </AppLayout>
  <UsageExportProgress :show="exportProgress.show" :progress="exportProgress.progress" :current="exportProgress.current" :total="exportProgress.total" :estimated-time="exportProgress.estimatedTime" @cancel="cancelExport" />
  <UsageCleanupDialog
    :show="cleanupDialogVisible"
    :filters="filters"
    :start-date="startDate"
    :end-date="endDate"
    @close="cleanupDialogVisible = false"
  />
  <!-- Balance history modal triggered from usage table user click -->
  <UserBalanceHistoryModal
    :show="showBalanceHistoryModal"
    :user="balanceHistoryUser"
    :hide-actions="true"
    @close="showBalanceHistoryModal = false; balanceHistoryUser = null"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { saveAs } from 'file-saver'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'; import { adminAPI } from '@/api/admin'; import { adminUsageAPI } from '@/api/admin/usage'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'
import { formatReasoningEffort } from '@/utils/format'
import { resolveUsageRequestType, requestTypeToLegacyStream } from '@/utils/usageRequestType'
import AppLayout from '@/components/layout/AppLayout.vue'; import Pagination from '@/components/common/Pagination.vue'; import Select from '@/components/common/Select.vue'; import DateRangePicker from '@/components/common/DateRangePicker.vue'
import UsageStatsCards from '@/components/admin/usage/UsageStatsCards.vue'; import UsageFilters from '@/components/admin/usage/UsageFilters.vue'
import UsageTable from '@/components/admin/usage/UsageTable.vue'; import UsageExportProgress from '@/components/admin/usage/UsageExportProgress.vue'
import UsageCleanupDialog from '@/components/admin/usage/UsageCleanupDialog.vue'
import UserBalanceHistoryModal from '@/components/admin/user/UserBalanceHistoryModal.vue'
import OpsErrorLogTable from '@/views/admin/ops/components/OpsErrorLogTable.vue'
import OpsErrorDetailModal from '@/views/admin/ops/components/OpsErrorDetailModal.vue'
import { listErrorLogs } from '@/api/admin/ops'
import type { OpsErrorLog } from '@/api/admin/ops'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'; import GroupDistributionChart from '@/components/charts/GroupDistributionChart.vue'; import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import EndpointDistributionChart from '@/components/charts/EndpointDistributionChart.vue'
import Icon from '@/components/icons/Icon.vue'
import type { AdminUsageLog, TrendDataPoint, ModelStat, GroupStat, EndpointStat, AdminUser } from '@/types'; import type { AdminUsageStatsResponse, AdminUsageQueryParams } from '@/api/admin/usage'

const { t } = useI18n()
const appStore = useAppStore()
type DistributionMetric = 'tokens' | 'actual_cost'
type CacheDistributionMetric = DistributionMetric | 'cache_hit_ratio' | 'cache_read_per_hit'
type EndpointSource = 'inbound' | 'upstream' | 'path'
type ModelDistributionSource = 'requested' | 'upstream' | 'mapping'
const route = useRoute()
const usageStats = ref<AdminUsageStatsResponse | null>(null); const usageLogs = ref<AdminUsageLog[]>([]); const loading = ref(false); const exporting = ref(false)
const trendData = ref<TrendDataPoint[]>([]); const requestedModelStats = ref<ModelStat[]>([]); const upstreamModelStats = ref<ModelStat[]>([]); const mappingModelStats = ref<ModelStat[]>([]); const groupStats = ref<GroupStat[]>([]); const chartsLoading = ref(false); const modelStatsLoading = ref(false); const granularity = ref<'day' | 'hour'>('hour')
const modelChartSectionRef = ref<HTMLElement | null>(null)

const modelDistributionMetric = ref<CacheDistributionMetric>('tokens')
const modelDistributionSource = ref<ModelDistributionSource>('requested')
const loadedModelSources = reactive<Record<ModelDistributionSource, boolean>>({
  requested: false,
  upstream: false,
  mapping: false,
})
const groupDistributionMetric = ref<DistributionMetric>('tokens')
const endpointDistributionMetric = ref<CacheDistributionMetric>('tokens')
const endpointDistributionSource = ref<EndpointSource>('inbound')
const inboundEndpointStats = ref<EndpointStat[]>([])
const upstreamEndpointStats = ref<EndpointStat[]>([])
const endpointPathStats = ref<EndpointStat[]>([])
const endpointStatsLoading = ref(false)
let abortController: AbortController | null = null; let exportAbortController: AbortController | null = null
let chartReqSeq = 0
let statsReqSeq = 0
let modelStatsReqSeq = 0
const exportProgress = reactive({ show: false, progress: 0, current: 0, total: 0, estimatedTime: '' })
const cleanupDialogVisible = ref(false)
// Balance history modal state
const showBalanceHistoryModal = ref(false)
const balanceHistoryUser = ref<AdminUser | null>(null)

const breakdownFilters = computed(() => {
  const f: Record<string, any> = {}
  if (filters.value.user_id) f.user_id = filters.value.user_id
  if (filters.value.api_key_id) f.api_key_id = filters.value.api_key_id
  if (filters.value.account_id) f.account_id = filters.value.account_id
  if (filters.value.group_id) f.group_id = filters.value.group_id
  if (filters.value.request_type != null) f.request_type = filters.value.request_type
  if (filters.value.billing_type != null) f.billing_type = filters.value.billing_type
  return f
})

const modelNameOptions = computed(() =>
  Array.from(new Set(requestedModelStats.value.map((m) => m.model).filter(Boolean))).sort()
)

const handleUserClick = async (userId: number) => {
  try {
    const user = await adminAPI.users.getById(userId, true)
    balanceHistoryUser.value = user
    showBalanceHistoryModal.value = true
  } catch {
    appStore.showError(t('admin.usage.failedToLoadUser'))
  }
}

const granularityOptions = computed(() => [{ value: 'day', label: t('admin.dashboard.day') }, { value: 'hour', label: t('admin.dashboard.hour') }])
const formatUsageCost = (value: number) => value.toFixed(4)
const formatUsageDuration = (ms: number) => ms < 1000 ? ms.toFixed(0) + 'ms' : (ms / 1000).toFixed(2) + 's'
const formatUsageTokens = (value: number) => value.toLocaleString()

// Use local timezone to avoid UTC timezone issues
const formatLD = (d: Date) => {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
const getLast24HoursRangeDates = (): { start: string; end: string } => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  return {
    start: formatLD(start),
    end: formatLD(end)
  }
}
const getGranularityForRange = (start: string, end: string): 'day' | 'hour' => {
  const startTime = new Date(`${start}T00:00:00`).getTime()
  const endTime = new Date(`${end}T00:00:00`).getTime()
  const daysDiff = Math.ceil((endTime - startTime) / (1000 * 60 * 60 * 24))
  return daysDiff <= 1 ? 'hour' : 'day'
}
const defaultRange = getLast24HoursRangeDates()
const startDate = ref(defaultRange.start); const endDate = ref(defaultRange.end)
const filters = ref<AdminUsageQueryParams>({ user_id: undefined, model: undefined, group_id: undefined, request_type: undefined, billing_type: null, start_date: startDate.value, end_date: endDate.value })
const pagination = reactive({ page: 1, page_size: getPersistedPageSize(), total: 0 })
const sortState = reactive({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
})

const getSingleQueryValue = (value: string | null | Array<string | null> | undefined): string | undefined => {
  if (Array.isArray(value)) return value.find((item): item is string => typeof item === 'string' && item.length > 0)
  return typeof value === 'string' && value.length > 0 ? value : undefined
}

const getNumericQueryValue = (value: string | null | Array<string | null> | undefined): number | undefined => {
  const raw = getSingleQueryValue(value)
  if (!raw) return undefined
  const parsed = Number(raw)
  return Number.isFinite(parsed) ? parsed : undefined
}

const scrollToCharts = () => {
  modelChartSectionRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

const usageBriefingLead = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.leadLoading')
  return t('usage.adminLedger.briefing.leadPrefix')
})

const usageBriefingRequestsValue = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.requestsPending')
  return t('usage.adminLedger.briefing.requestsValue', { value: formatUsageTokens(usageStats.value.total_requests || 0) })
})

const usageBriefingRequestsNote = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.requestsNotePending')
  return t('usage.adminLedger.briefing.requestsNote', {
    value: formatUsageTokens(usageStats.value.total_tokens || 0),
    duration: formatUsageDuration(usageStats.value.average_duration_ms || 0)
  })
})

const usageBriefingCostValue = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.costPending')
  return t('usage.adminLedger.briefing.costValue', { value: formatUsageCost(usageStats.value.total_actual_cost || 0) })
})

const usageBriefingCostNote = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.costNotePending')
  return t('usage.adminLedger.briefing.costNote', {
    accountCost: formatUsageCost(usageStats.value.total_account_cost || 0),
    standardCost: formatUsageCost(usageStats.value.total_cost || 0)
  })
})

const usageErrorCount = computed(() => errTotal.value || errRows.value.length)

const usageBriefingAnomalyValue = computed(() => {
  if (errLoading.value) return t('usage.adminLedger.briefing.anomalyPending')
  if (usageErrorCount.value > 0) return t('usage.adminLedger.briefing.anomalyValue', { value: formatUsageTokens(usageErrorCount.value) })
  return t('usage.adminLedger.briefing.anomalyEmpty')
})

const usageBriefingAnomalyNote = computed(() => {
  if (errLoading.value) return t('usage.adminLedger.briefing.anomalyNotePending')
  if (usageErrorCount.value > 0) {
    return t('usage.adminLedger.briefing.anomalyNote', { value: formatUsageTokens(usageErrorCount.value) })
  }
  return t('usage.adminLedger.briefing.anomalyNoteEmpty')
})

const usageBriefingActionValue = computed(() => {
  if (errLoading.value || loading.value) return t('usage.adminLedger.briefing.actionPending')
  if (usageErrorCount.value > 0) return t('usage.adminLedger.briefing.actionErrors')
  if ((groupStats.value?.length || 0) > 0 || (inboundEndpointStats.value?.length || 0) > 0) return t('usage.adminLedger.briefing.actionGroups')
  return t('usage.adminLedger.briefing.actionDefault')
})

const usageBriefingActionNote = computed(() => {
  if (errLoading.value || loading.value) return t('usage.adminLedger.briefing.actionNotePending')
  if (usageErrorCount.value > 0) return t('usage.adminLedger.briefing.actionNoteErrors')
  if ((groupStats.value?.length || 0) > 0 || (inboundEndpointStats.value?.length || 0) > 0) return t('usage.adminLedger.briefing.actionNoteGroups')
  return t('usage.adminLedger.briefing.actionNoteDefault')
})

const usageBriefingTags = computed(() => {
  const tags: string[] = []
  if (usageStats.value?.total_tokens) tags.push(t('usage.adminLedger.briefing.tagTokens', { value: formatUsageTokens(usageStats.value.total_tokens) }))
  if (usageErrorCount.value > 0) tags.push(t('usage.adminLedger.briefing.tagErrors', { value: formatUsageTokens(usageErrorCount.value) }))
  if ((groupStats.value?.length || 0) > 0) tags.push(t('usage.adminLedger.briefing.tagGroups', { value: formatUsageTokens(groupStats.value.length) }))
  if ((inboundEndpointStats.value?.length || 0) > 0) tags.push(t('usage.adminLedger.briefing.tagEndpoints', { value: formatUsageTokens(inboundEndpointStats.value.length) }))
  return tags.slice(0, 5)
})

const usageBriefingSummary = computed(() => {
  if (!usageStats.value) return t('usage.adminLedger.briefing.summaryPending')
  const parts = [
    t('usage.adminLedger.briefing.summaryRequests', {
      start: startDate.value,
      end: endDate.value,
      requests: formatUsageTokens(usageStats.value.total_requests || 0)
    }),
    t('usage.adminLedger.briefing.summaryCost', { cost: formatUsageCost(usageStats.value.total_actual_cost || 0) })
  ]
  if (usageErrorCount.value > 0) {
    parts.push(t('usage.adminLedger.briefing.summaryErrors', { errors: formatUsageTokens(usageErrorCount.value) }))
    parts.push(t('usage.adminLedger.briefing.summaryActionErrors'))
  } else {
    parts.push(t('usage.adminLedger.briefing.summaryNoErrors'))
    parts.push(t('usage.adminLedger.briefing.summaryActionGroups'))
  }
  return parts.join('，') + '。'
})

const applyRouteQueryFilters = () => {
  const queryStartDate = getSingleQueryValue(route.query.start_date)
  const queryEndDate = getSingleQueryValue(route.query.end_date)
  const queryUserId = getNumericQueryValue(route.query.user_id)

  if (queryStartDate) {
    startDate.value = queryStartDate
  }
  if (queryEndDate) {
    endDate.value = queryEndDate
  }

  filters.value = {
    ...filters.value,
    user_id: queryUserId,
    start_date: startDate.value,
    end_date: endDate.value
  }
  granularity.value = getGranularityForRange(startDate.value, endDate.value)
}

const onDateRangeChange = (range: { startDate: string; endDate: string; preset: string | null }) => {
  startDate.value = range.startDate
  endDate.value = range.endDate
  filters.value = {
    ...filters.value,
    start_date: range.startDate,
    end_date: range.endDate
  }
  granularity.value = getGranularityForRange(range.startDate, range.endDate)
  applyFilters()
}

const buildUsageListParams = (
  page: number,
  pageSize: number,
  exactTotal: boolean
): AdminUsageQueryParams => {
  const requestType = filters.value.request_type
  const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
  return {
    page,
    page_size: pageSize,
    exact_total: exactTotal,
    ...filters.value,
    stream: legacyStream === null ? undefined : legacyStream,
    sort_by: sortState.sort_by,
    sort_order: sortState.sort_order
  }
}

const loadLogs = async () => {
  abortController?.abort(); const c = new AbortController(); abortController = c; loading.value = true
  try {
    const res = await adminAPI.usage.list(
      buildUsageListParams(pagination.page, pagination.page_size, false),
      { signal: c.signal }
    )
    if(!c.signal.aborted) { usageLogs.value = res.items; pagination.total = res.total }
  } catch (error: any) { if(error?.name !== 'AbortError') console.error('Failed to load usage logs:', error) } finally { if(abortController === c) loading.value = false }
}
const loadStats = async (force = false) => {
  const seq = ++statsReqSeq
  endpointStatsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const s = await adminAPI.usage.getStats({
      ...filters.value,
      stream: legacyStream === null ? undefined : legacyStream,
      ...(force ? { nocache: 1 } : {}),
    })
    if (seq !== statsReqSeq) return
    usageStats.value = s
    inboundEndpointStats.value = s.endpoints || []
    upstreamEndpointStats.value = s.upstream_endpoints || []
    endpointPathStats.value = s.endpoint_paths || []
  } catch (error) {
    if (seq !== statsReqSeq) return
    console.error('Failed to load usage stats:', error)
    inboundEndpointStats.value = []
    upstreamEndpointStats.value = []
    endpointPathStats.value = []
  } finally {
    if (seq === statsReqSeq) endpointStatsLoading.value = false
  }
}

// 失效模型统计缓存:仅标记需要重取,保留旧数据直到新数据到达(避免刷新时图表闪空)。
const invalidateModelStatsCache = () => {
  loadedModelSources.requested = false
  loadedModelSources.upstream = false
  loadedModelSources.mapping = false
}

const loadModelStats = async (source: ModelDistributionSource, force = false) => {
  if (!force && loadedModelSources[source]) {
    return
  }

  const seq = ++modelStatsReqSeq
  modelStatsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const baseParams = {
      start_date: filters.value.start_date || startDate.value,
      end_date: filters.value.end_date || endDate.value,
      user_id: filters.value.user_id,
      model: filters.value.model,
      api_key_id: filters.value.api_key_id,
      account_id: filters.value.account_id,
      group_id: filters.value.group_id,
      request_type: requestType,
      stream: legacyStream === null ? undefined : legacyStream,
      billing_type: filters.value.billing_type,
    }

    const response = await adminAPI.dashboard.getModelStats({ ...baseParams, model_source: source })

    if (seq !== modelStatsReqSeq) return

    const models = response.models || []
    if (source === 'requested') {
      requestedModelStats.value = models
    } else if (source === 'upstream') {
      upstreamModelStats.value = models
    } else {
      mappingModelStats.value = models
    }
    loadedModelSources[source] = true
  } catch (error) {
    if (seq !== modelStatsReqSeq) return
    console.error('Failed to load model stats:', error)
    if (source === 'requested') {
      requestedModelStats.value = []
    } else if (source === 'upstream') {
      upstreamModelStats.value = []
    } else {
      mappingModelStats.value = []
    }
    loadedModelSources[source] = false
  } finally {
    if (seq === modelStatsReqSeq) modelStatsLoading.value = false
  }
}

const loadChartData = async () => {
  const seq = ++chartReqSeq
  chartsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const snapshot = await adminAPI.dashboard.getSnapshotV2({
      start_date: filters.value.start_date || startDate.value,
      end_date: filters.value.end_date || endDate.value,
      granularity: granularity.value,
      user_id: filters.value.user_id,
      model: filters.value.model,
      api_key_id: filters.value.api_key_id,
      account_id: filters.value.account_id,
      group_id: filters.value.group_id,
      request_type: requestType,
      stream: legacyStream === null ? undefined : legacyStream,
      billing_type: filters.value.billing_type,
      include_stats: false,
      include_trend: true,
      include_model_stats: false,
      include_group_stats: true,
      include_users_trend: false
    })
    if (seq !== chartReqSeq) return
    trendData.value = snapshot.trend || []
    groupStats.value = snapshot.groups || []
  } catch (error) { console.error('Failed to load chart data:', error) } finally { if (seq === chartReqSeq) chartsLoading.value = false }
}
const applyFilters = () => {
  pagination.page = 1
  invalidateModelStatsCache()
  loadLogs()
  loadStats()
  loadModelStats(modelDistributionSource.value, true)
  loadChartData()
  errPage.value = 1
  if (activeTab.value === 'errors') {
    loadAdminErrors()
  } else {
    errRows.value = []
  }
}
const refreshData = () => {
  invalidateModelStatsCache()
  loadLogs()
  loadStats(true)
  loadModelStats(modelDistributionSource.value, true)
  loadChartData()
  if (activeTab.value === 'errors') loadAdminErrors()
}
const resetFilters = () => {
  const range = getLast24HoursRangeDates()
  startDate.value = range.start
  endDate.value = range.end
  filters.value = { start_date: startDate.value, end_date: endDate.value, request_type: undefined, billing_type: null, billing_mode: undefined }
  granularity.value = getGranularityForRange(startDate.value, endDate.value)
  applyFilters()
}
const handlePageChange = (p: number) => { pagination.page = p; loadLogs() }
const handlePageSizeChange = (s: number) => { pagination.page_size = s; pagination.page = 1; loadLogs() }
const handleSort = (key: string, order: 'asc' | 'desc') => {
  sortState.sort_by = key
  sortState.sort_order = order
  pagination.page = 1
  loadLogs()
}
const cancelExport = () => exportAbortController?.abort()
const openCleanupDialog = () => { cleanupDialogVisible.value = true }
const getRequestTypeLabel = (log: AdminUsageLog): string => {
  const requestType = resolveUsageRequestType(log)
  if (requestType === 'cyber') return t('usage.cyber')
  if (requestType === 'ws_v2') return t('usage.ws')
  if (requestType === 'stream') return t('usage.stream')
  if (requestType === 'sync') return t('usage.sync')
  return t('usage.unknown')
}

const exportToExcel = async () => {
  if (exporting.value) return; exporting.value = true; exportProgress.show = true
  const c = new AbortController(); exportAbortController = c
  try {
    let p = 1; let total = pagination.total; let exportedCount = 0
    const XLSX = await import('xlsx')
    const headers = [
      t('usage.time'), t('admin.usage.user'), t('usage.apiKeyFilter'),
      t('admin.usage.account'), t('usage.model'), t('usage.upstreamModel'), t('usage.reasoningEffort'), t('admin.usage.group'),
      t('usage.inboundEndpoint'), t('usage.upstreamEndpoint'),
      t('usage.type'),
      t('admin.usage.inputTokens'), t('admin.usage.outputTokens'),
      t('admin.usage.cacheReadTokens'), t('admin.usage.cacheCreationTokens'),
      t('admin.usage.inputCost'), t('admin.usage.outputCost'),
      t('admin.usage.cacheReadCost'), t('admin.usage.cacheCreationCost'),
      t('usage.rate'), t('usage.accountMultiplier'), t('usage.original'), t('usage.userBilled'), t('usage.accountBilled'),
      t('usage.firstToken'), t('usage.duration'),
      t('admin.usage.requestId'), t('usage.userAgent'), t('admin.usage.ipAddress')
    ]
    const ws = XLSX.utils.aoa_to_sheet([headers])
    while (true) {
      const res = await adminUsageAPI.list(
        buildUsageListParams(p, 100, true),
        { signal: c.signal }
      )
      if (c.signal.aborted) break; if (p === 1) { total = res.total; exportProgress.total = total }
      const rows = (res.items || []).map((log: AdminUsageLog) => [
        log.created_at, log.user?.email || '', log.api_key?.name || '', log.account?.name || '', log.model,
        log.upstream_model || '', formatReasoningEffort(log.reasoning_effort), log.group?.name || '',
        log.inbound_endpoint || '', log.upstream_endpoint || '', getRequestTypeLabel(log),
        log.input_tokens, log.output_tokens, log.cache_read_tokens, log.cache_creation_tokens,
        log.input_cost?.toFixed(6) || '0.000000', log.output_cost?.toFixed(6) || '0.000000',
        log.cache_read_cost?.toFixed(6) || '0.000000', log.cache_creation_cost?.toFixed(6) || '0.000000',
        log.rate_multiplier?.toPrecision(4) || '1.00', (log.account_rate_multiplier ?? 1).toPrecision(4),
        log.total_cost?.toFixed(6) || '0.000000', log.actual_cost?.toFixed(6) || '0.000000',
        ((log.account_stats_cost ?? log.total_cost) * (log.account_rate_multiplier ?? 1)).toFixed(6), log.first_token_ms ?? '', log.duration_ms,
        log.request_id || '', log.user_agent || '', log.ip_address || ''
      ])
      if (rows.length) {
        XLSX.utils.sheet_add_aoa(ws, rows, { origin: -1 })
      }
      exportedCount += rows.length
      exportProgress.current = exportedCount
      exportProgress.progress = total > 0 ? Math.min(100, Math.round(exportedCount / total * 100)) : 0
      if (exportedCount >= total || res.items.length < 100) break; p++
    }
    if(!c.signal.aborted) {
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, 'Usage')
      saveAs(new Blob([XLSX.write(wb, { bookType: 'xlsx', type: 'array' })], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }), `usage_${filters.value.start_date}_to_${filters.value.end_date}.xlsx`)
      appStore.showSuccess(t('usage.exportSuccess'))
    }
  } catch (error) { console.error('Failed to export:', error); appStore.showError('Export Failed') }
  finally { if(exportAbortController === c) { exportAbortController = null; exporting.value = false; exportProgress.show = false } }
}

// Column visibility
const ALWAYS_VISIBLE = ['user', 'created_at']
const DEFAULT_HIDDEN_COLUMNS = ['reasoning_effort', 'user_agent']
const HIDDEN_COLUMNS_KEY = 'usage-hidden-columns'

const allColumns = computed(() => [
  { key: 'user', label: t('admin.usage.user'), sortable: false },
  { key: 'api_key', label: t('usage.apiKeyFilter'), sortable: false },
  { key: 'account', label: t('admin.usage.account'), sortable: false },
  { key: 'model', label: t('usage.model'), sortable: true },
  { key: 'reasoning_effort', label: t('usage.reasoningEffort'), sortable: false },
  { key: 'endpoint', label: t('usage.endpoint'), sortable: false },
  { key: 'group', label: t('admin.usage.group'), sortable: false },
  { key: 'stream', label: t('usage.type'), sortable: false },
  { key: 'billing_mode', label: t('admin.usage.billingMode'), sortable: false },
  { key: 'tokens', label: t('usage.tokens'), sortable: false },
  { key: 'cost', label: t('usage.cost'), sortable: false },
  { key: 'first_token', label: t('usage.firstToken'), sortable: false },
  { key: 'duration', label: t('usage.duration'), sortable: false },
  { key: 'created_at', label: t('usage.time'), sortable: true },
  { key: 'user_agent', label: t('usage.userAgent'), sortable: false },
  { key: 'ip_address', label: t('admin.usage.ipAddress'), sortable: false }
])

const hiddenColumns = reactive<Set<string>>(new Set())

const toggleableColumns = computed(() =>
  allColumns.value.filter(col => !ALWAYS_VISIBLE.includes(col.key))
)

const visibleColumns = computed(() =>
  allColumns.value.filter(col =>
    ALWAYS_VISIBLE.includes(col.key) || !hiddenColumns.has(col.key)
  )
)

const isColumnVisible = (key: string) => !hiddenColumns.has(key)

const toggleColumn = (key: string) => {
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  } else {
    hiddenColumns.add(key)
  }
  try {
    localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
  } catch (e) {
    console.error('Failed to save columns:', e)
  }
}

const loadSavedColumns = () => {
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      (JSON.parse(saved) as string[]).forEach((key) => {
        hiddenColumns.add(key)
      })
    } else {
      DEFAULT_HIDDEN_COLUMNS.forEach((key) => {
        hiddenColumns.add(key)
      })
    }
  } catch {
    DEFAULT_HIDDEN_COLUMNS.forEach((key) => {
      hiddenColumns.add(key)
    })
  }
}

// Error tab state
const activeTab = ref<'usage' | 'errors'>('usage')
const errRows = ref<OpsErrorLog[]>([])
const errLoading = ref(false)
const errPage = ref(1)
const errPageSize = ref(20)
const errTotal = ref(0)
const showErrorModal = ref(false)
const selectedErrorId = ref<number | null>(null)

// 注意：'YYYY-MM-DDT00:00:00' 无时区后缀，按本地时区解析后再转 UTC——与页面其它日期处理语义一致，刻意如此，勿改成 'T00:00:00Z'
const toRFC3339 = (d: string | undefined, endOfDay = false): string | undefined =>
  d ? new Date(d + (endOfDay ? 'T23:59:59.999' : 'T00:00:00')).toISOString() : undefined

const loadAdminErrors = async () => {
  errLoading.value = true
  try {
    const resp = await listErrorLogs({
      page: errPage.value,
      page_size: errPageSize.value,
      view: 'all',
      start_time: toRFC3339(filters.value.start_date),
      end_time: toRFC3339(filters.value.end_date, true),
      user_id: filters.value.user_id ?? undefined,
      api_key_id: filters.value.api_key_id ?? undefined,
      account_id: filters.value.account_id ?? undefined,
      group_id: filters.value.group_id ?? undefined,
      model: filters.value.model || undefined,
    })
    errRows.value = resp.items
    errTotal.value = resp.total
  } catch (error) {
    console.error('Failed to load admin errors:', error)
    appStore.showError(t('usage.errors.failedToLoad'))
  } finally {
    errLoading.value = false
  }
}

const onErrPage = (p: number) => { errPage.value = p; loadAdminErrors() }
const onErrPageSize = (s: number) => { errPageSize.value = s; errPage.value = 1; loadAdminErrors() }
const openError = (id: number) => { selectedErrorId.value = id; showErrorModal.value = true }
const switchToErrorsTab = () => { activeTab.value = 'errors'; if (errRows.value.length === 0) loadAdminErrors() }

const showColumnDropdown = ref(false)
const columnDropdownRef = ref<HTMLElement | null>(null)

const handleColumnClickOutside = (event: MouseEvent) => {
  if (columnDropdownRef.value && !columnDropdownRef.value.contains(event.target as HTMLElement)) {
    showColumnDropdown.value = false
  }
}

onMounted(() => {
  applyRouteQueryFilters()
  loadLogs()
  loadStats()
  loadModelStats(modelDistributionSource.value, true)
  window.setTimeout(() => {
    void loadChartData()
  }, 120)
  loadSavedColumns()
  document.addEventListener('click', handleColumnClickOutside)
})
onUnmounted(() => { abortController?.abort(); exportAbortController?.abort(); document.removeEventListener('click', handleColumnClickOutside) })

watch(modelDistributionSource, (source) => {
  void loadModelStats(source)
})

defineExpose({ requestedModelStats, refreshData })
</script>

<style scoped>
.admin-usage-shell {
  max-width: 88rem;
}

.admin-usage-briefing-card {
  display: none;
}

.admin-usage-topline {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
}

@media (max-width: 1024px) {
  .admin-usage-topline {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .admin-usage-topline {
    grid-template-columns: 1fr;
  }
}
</style>

<style scoped>
.admin-usage-shell {
  padding: clamp(0.2rem, 0.6vw, 0.4rem);
}

.admin-usage-summary-stack {
  display: grid;
  gap: 0.85rem;
}

.admin-page-hero {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.15rem 1.25rem 1.05rem;
  border: 1px solid rgba(198, 184, 157, 0.52);
  border-radius: 12px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.055), transparent 30%),
    linear-gradient(180deg, rgba(255, 252, 245, 0.9), rgba(246, 241, 231, 0.78));
}

.admin-page-kicker {
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
  font-size: clamp(1.4rem, 1.8vw, 1.82rem);
  font-weight: 600;
}

.admin-page-hero p {
  max-width: 38rem;
  margin-top: 0.5rem;
  color: #5f675d;
  font-size: 0.94rem;
  line-height: 1.7;
}

.admin-usage-briefing-card {
  background:
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.1), transparent 30%),
    linear-gradient(135deg, rgba(255, 252, 246, 0.94), rgba(244, 238, 227, 0.88));
}

.admin-usage-briefing-head {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
}

.admin-usage-briefing-head h3 {
  margin-top: 0.45rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.18rem;
  font-weight: 600;
}

.admin-usage-briefing-head p,
.admin-usage-briefing-panel p,
.admin-usage-briefing-summary {
  color: #5f675d;
  line-height: 1.7;
  overflow-wrap: anywhere;
}

.admin-usage-briefing-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
}

.admin-usage-briefing-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.9rem;
}

.admin-usage-briefing-panel {
  display: grid;
  gap: 0.4rem;
  border: 1px solid rgba(198, 184, 157, 0.38);
  border-radius: 12px;
  background: rgba(255, 252, 245, 0.82);
  padding: 1rem;
}

.admin-usage-briefing-panel span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-usage-briefing-panel strong {
  color: #1f2320;
  font-size: 1rem;
  font-weight: 700;
}

.admin-usage-briefing-foot {
  display: flex;
  flex-wrap: wrap;
  align-items: start;
  justify-content: space-between;
  gap: 0.9rem;
}

.admin-usage-briefing-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.55rem;
}

.admin-usage-briefing-tag {
  display: inline-flex;
  align-items: center;
  max-width: 100%;
  padding: 0.34rem 0.7rem;
  border-radius: 999px;
  border: 1px solid rgba(167, 58, 42, 0.16);
  color: #7d4d3d;
  background: rgba(167, 58, 42, 0.06);
  font-size: 0.76rem;
  line-height: 1;
}

.admin-usage-briefing-tag-muted {
  border-color: rgba(198, 184, 157, 0.24);
  color: #8b7e6a;
  background: rgba(246, 240, 229, 0.9);
}

.usage-ledger-panel {
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 22px;
  background: rgba(250, 247, 239, 0.52);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
  overflow: hidden;
}

.usage-ledger-panel :deep(.table-wrapper) {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  max-height: none;
}

.usage-ledger-panel :deep(table) {
  background: transparent;
}

.usage-ledger-panel :deep(thead) {
  background: rgba(237, 229, 212, 0.72);
  backdrop-filter: blur(8px);
}

.usage-ledger-panel :deep(th) {
  border-bottom: 1px solid rgba(198, 184, 157, 0.54);
  color: #59645a;
  font-size: 0.76rem;
  font-weight: 650;
  letter-spacing: 0.12em;
  padding-top: 1.1rem;
  padding-bottom: 1.1rem;
}

.usage-ledger-panel :deep(td) {
  border-bottom: 1px solid rgba(198, 184, 157, 0.32);
  color: #38413a;
  padding-top: 1.15rem;
  padding-bottom: 1.15rem;
}

.usage-ledger-panel :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.055);
}

.usage-ledger-panel :deep(tbody .sticky-col) {
  background: rgba(250, 247, 239, 0.52) !important;
}

.usage-ledger-panel :deep(tbody tr:hover .sticky-col) {
  background: rgba(167, 58, 42, 0.055) !important;
}

.usage-ledger-panel :deep(.console-workbench-footer) {
  border: 0;
  border-top: 1px solid rgba(198, 184, 157, 0.32);
  background: rgba(250, 247, 239, 0.52);
  padding: 0.7rem 1rem 0.9rem;
}

@media (max-width: 980px) {
  .admin-usage-briefing-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 780px) {
  .admin-usage-briefing-grid {
    grid-template-columns: 1fr;
  }

  .admin-usage-briefing-actions,
  .admin-usage-briefing-foot {
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

.dark .admin-page-hero h2 {
  color: #f4efe4;
}

.dark .admin-page-hero p {
  color: #bdb5a8;
}

.dark .admin-usage-briefing-card {
  background:
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.14), transparent 30%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.9), rgba(18, 20, 16, 0.86));
}

.dark .admin-usage-briefing-head h3,
.dark .admin-usage-briefing-panel strong {
  color: #f4efe4;
}

.dark .admin-usage-briefing-head p,
.dark .admin-usage-briefing-panel p,
.dark .admin-usage-briefing-summary {
  color: #bdb5a8;
}

.dark .admin-usage-briefing-panel {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.72);
}

.dark .admin-usage-briefing-tag {
  border-color: rgba(167, 58, 42, 0.28);
  color: #f0c1b4;
  background: rgba(167, 58, 42, 0.14);
}

.dark .admin-usage-briefing-tag-muted {
  border-color: rgba(80, 74, 58, 0.5);
  color: #d8ccb8;
  background: rgba(40, 33, 25, 0.72);
}

.dark .usage-ledger-panel {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .usage-ledger-panel thead {
  background: rgba(17, 19, 15, 0.62);
}

.dark .usage-ledger-panel th,
.dark .usage-ledger-panel td {
  border-color: rgba(48, 52, 43, 0.82);
}

.dark .usage-ledger-panel th {
  color: #879186;
}

.dark .usage-ledger-panel td {
  color: #d9d0be;
}

.dark .usage-ledger-panel tbody .sticky-col {
  background: rgba(24, 26, 21, 0.72) !important;
}

.dark .usage-ledger-panel tbody tr:hover .sticky-col {
  background: rgba(216, 205, 185, 0.045) !important;
}

.dark .usage-ledger-panel .console-workbench-footer {
  background: rgba(24, 26, 21, 0.72);
}
</style>

