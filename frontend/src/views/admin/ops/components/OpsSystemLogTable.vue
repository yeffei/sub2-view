<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { opsAPI, type OpsRuntimeLogConfig, type OpsSystemLog, type OpsSystemLogSinkHealth } from '@/api/admin/ops'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import { useAppStore } from '@/stores'
import {
  buildCacheInstrumentationSummary,
  formatCacheInstrumentationDetail,
  isCacheInstrumentationLog
} from '../utils/cacheInstrumentation'

const appStore = useAppStore()

const props = withDefaults(defineProps<{
  platformFilter?: string
  refreshToken?: number
}>(), {
  platformFilter: '',
  refreshToken: 0
})

const loading = ref(false)
const logs = ref<OpsSystemLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const health = ref<OpsSystemLogSinkHealth>({
  queue_depth: 0,
  queue_capacity: 0,
  dropped_count: 0,
  write_failed_count: 0,
  written_count: 0,
  avg_write_delay_ms: 0
})

const runtimeLoading = ref(false)
const runtimeSaving = ref(false)
const runtimeConfig = reactive<OpsRuntimeLogConfig>({
  level: 'info',
  enable_sampling: false,
  sampling_initial: 100,
  sampling_thereafter: 100,
  caller: true,
  stacktrace_level: 'error',
  retention_days: 30
})

const filters = reactive({
  time_range: '1h' as '5m' | '30m' | '1h' | '6h' | '24h' | '7d' | '30d',
  start_time: '',
  end_time: '',
  level: '',
  component: '',
  request_id: '',
  client_request_id: '',
  user_id: '',
  account_id: '',
  platform: '',
  model: '',
  q: ''
})

const runtimeLevelOptions = [
  { value: 'debug', label: 'debug' },
  { value: 'info', label: 'info' },
  { value: 'warn', label: 'warn' },
  { value: 'error', label: 'error' }
]

const stacktraceLevelOptions = [
  { value: 'none', label: 'none' },
  { value: 'error', label: 'error' },
  { value: 'fatal', label: 'fatal' }
]

const timeRangeOptions = [
  { value: '5m', label: '5m' },
  { value: '30m', label: '30m' },
  { value: '1h', label: '1h' },
  { value: '6h', label: '6h' },
  { value: '24h', label: '24h' },
  { value: '7d', label: '7d' },
  { value: '30d', label: '30d' }
]

const filterLevelOptions = [
  { value: '', label: '全部' },
  { value: 'debug', label: 'debug' },
  { value: 'info', label: 'info' },
  { value: 'warn', label: 'warn' },
  { value: 'error', label: 'error' }
]

const levelBadgeClass = (level: string) => {
  const v = String(level || '').toLowerCase()
  if (v === 'error' || v === 'fatal') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
  if (v === 'warn' || v === 'warning') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  if (v === 'debug') return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
  return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
}

const formatTime = (value: string) => {
  if (!value) return '-'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString()
}

const getExtraString = (extra: Record<string, any> | undefined, key: string) => {
  if (!extra) return ''
  const v = extra[key]
  if (v == null) return ''
  if (typeof v === 'string') return v.trim()
  if (typeof v === 'number' || typeof v === 'boolean') return String(v)
  return ''
}

const isRoutingExplanationLog = (row: OpsSystemLog) => {
  const component = String(row.component || row.extra?.component || '').trim()
  return component === 'routing.explanation' || String(row.message || '').trim() === 'routing_explanation'
}

const formatRoutingExplanationDetail = (row: OpsSystemLog) => {
  const extra = row.extra || {}
  const reason = getExtraString(extra, 'reason') || 'selection'
  const layer = getExtraString(extra, 'layer') || '-'
  const accountID = getExtraString(extra, 'account_id') || String(row.account_id || '')
  const candidateCount = getExtraString(extra, 'candidate_count')
  const topK = getExtraString(extra, 'top_k')
  const waitPlanned = getExtraString(extra, 'wait_planned')
  const latency = getExtraString(extra, 'routing_latency_ms') || getExtraString(extra, 'latency_ms')
  const transport = getExtraString(extra, 'required_transport')
  const capability = getExtraString(extra, 'required_capability') || getExtraString(extra, 'required_image_capability')
  const fallback = getExtraString(extra, 'fallback')
  const poolName = getExtraString(extra, 'pool_name')
  const poolCode = getExtraString(extra, 'pool_code')
  const stickyEscapeSource = getExtraString(extra, 'sticky_escape_source')
  const stickyEscapeEnabled = getExtraString(extra, 'sticky_escape_enabled')
  const stickyEscapeTTFT = getExtraString(extra, 'sticky_escape_ttft_ms')
  const stickyEscapeErrorRate = getExtraString(extra, 'sticky_escape_error_rate')
  const stickyEscapeTriggered = getExtraString(extra, 'sticky_escape_triggered')
  const stickyEscapeReason = getExtraString(extra, 'sticky_escape_reason')
  const error = getExtraString(extra, 'error')

  const parts = [`调度解释 reason=${reason}`, `layer=${layer}`]
  if (accountID) parts.push(`account=${accountID}`)
  if (candidateCount) parts.push(`candidates=${candidateCount}`)
  if (topK) parts.push(`top_k=${topK}`)
  if (waitPlanned) parts.push(`wait=${waitPlanned}`)
  if (latency) parts.push(`latency_ms=${latency}`)
  if (transport) parts.push(`transport=${transport}`)
  if (capability) parts.push(`capability=${capability}`)
  if (fallback) parts.push(`fallback=${fallback}`)
  if (poolName || poolCode) parts.push(`pool=${poolName || '-'}${poolCode ? `(${poolCode})` : ''}`)
  if (stickyEscapeSource) parts.push(`escape_source=${stickyEscapeSource}`)
  if (stickyEscapeEnabled) parts.push(`escape_enabled=${stickyEscapeEnabled}`)
  if (stickyEscapeTTFT) parts.push(`escape_ttft_ms=${stickyEscapeTTFT}`)
  if (stickyEscapeErrorRate) parts.push(`escape_error_rate=${stickyEscapeErrorRate}`)
  if (stickyEscapeTriggered) parts.push(`escape_triggered=${stickyEscapeTriggered}`)
  if (stickyEscapeReason) parts.push(`escape_reason=${stickyEscapeReason}`)
  if (row.model) parts.push(`model=${row.model}`)
  if (error) parts.push(`error=${error}`)
  return parts.join('  ')
}

const formatSystemLogDetail = (row: OpsSystemLog) => {
  if (isCacheInstrumentationLog(row)) return formatCacheInstrumentationDetail(row)
  if (isRoutingExplanationLog(row)) return formatRoutingExplanationDetail(row)

  const parts: string[] = []
  const msg = String(row.message || '').trim()
  if (msg) parts.push(msg)

  const extra = row.extra || {}
  const statusCode = getExtraString(extra, 'status_code')
  const latencyMs = getExtraString(extra, 'latency_ms')
  const method = getExtraString(extra, 'method')
  const path = getExtraString(extra, 'path')
  const clientIP = getExtraString(extra, 'client_ip')
  const protocol = getExtraString(extra, 'protocol')

  const accessParts: string[] = []
  if (statusCode) accessParts.push(`status=${statusCode}`)
  if (latencyMs) accessParts.push(`latency_ms=${latencyMs}`)
  if (method) accessParts.push(`method=${method}`)
  if (path) accessParts.push(`path=${path}`)
  if (clientIP) accessParts.push(`ip=${clientIP}`)
  if (protocol) accessParts.push(`proto=${protocol}`)
  if (accessParts.length > 0) parts.push(accessParts.join(' '))

  const corrParts: string[] = []
  if (row.request_id) corrParts.push(`req=${row.request_id}`)
  if (row.client_request_id) corrParts.push(`client_req=${row.client_request_id}`)
  if (row.user_id != null) corrParts.push(`user=${row.user_id}`)
  if (row.account_id != null) corrParts.push(`acc=${row.account_id}`)
  if (row.platform) corrParts.push(`platform=${row.platform}`)
  if (row.model) corrParts.push(`model=${row.model}`)
  if (corrParts.length > 0) parts.push(corrParts.join(' '))

  const errors = getExtraString(extra, 'errors')
  if (errors) parts.push(`errors=${errors}`)
  const err = getExtraString(extra, 'err') || getExtraString(extra, 'error')
  if (err) parts.push(`error=${err}`)

  // 用空格拼接，交给 CSS 自动换行，尽量“填满再换行”。
  return parts.join('  ')
}

const toRFC3339 = (value: string) => {
  if (!value) return undefined
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return undefined
  return d.toISOString()
}

const buildQuery = () => {
  const query: Record<string, any> = {
    page: page.value,
    page_size: pageSize.value,
    time_range: filters.time_range
  }

  if (filters.time_range === '30d') {
    query.time_range = '30d'
  }
  if (filters.start_time) query.start_time = toRFC3339(filters.start_time)
  if (filters.end_time) query.end_time = toRFC3339(filters.end_time)
  if (filters.level.trim()) query.level = filters.level.trim()
  if (filters.component.trim()) query.component = filters.component.trim()
  if (filters.request_id.trim()) query.request_id = filters.request_id.trim()
  if (filters.client_request_id.trim()) query.client_request_id = filters.client_request_id.trim()
  if (filters.user_id.trim()) {
    const v = Number.parseInt(filters.user_id.trim(), 10)
    if (Number.isFinite(v) && v > 0) query.user_id = v
  }
  if (filters.account_id.trim()) {
    const v = Number.parseInt(filters.account_id.trim(), 10)
    if (Number.isFinite(v) && v > 0) query.account_id = v
  }
  if (filters.platform.trim()) query.platform = filters.platform.trim()
  if (filters.model.trim()) query.model = filters.model.trim()
  if (filters.q.trim()) query.q = filters.q.trim()
  return query
}

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await opsAPI.listSystemLogs(buildQuery())
    logs.value = res.items || []
    total.value = res.total || 0
  } catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to fetch logs', err)
    appStore.showError(err?.response?.data?.detail || '系统日志加载失败')
  } finally {
    loading.value = false
  }
}

const fetchHealth = async () => {
  try {
    health.value = await opsAPI.getSystemLogSinkHealth()
  } catch {
    // 忽略健康数据读取失败，不影响主流程。
  }
}

const loadRuntimeConfig = async () => {
  runtimeLoading.value = true
  try {
    const cfg = await opsAPI.getRuntimeLogConfig()
    runtimeConfig.level = cfg.level
    runtimeConfig.enable_sampling = cfg.enable_sampling
    runtimeConfig.sampling_initial = cfg.sampling_initial
    runtimeConfig.sampling_thereafter = cfg.sampling_thereafter
    runtimeConfig.caller = cfg.caller
    runtimeConfig.stacktrace_level = cfg.stacktrace_level
    runtimeConfig.retention_days = cfg.retention_days
  } catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to load runtime log config', err)
  } finally {
    runtimeLoading.value = false
  }
}

const saveRuntimeConfig = async () => {
  runtimeSaving.value = true
  try {
    const saved = await opsAPI.updateRuntimeLogConfig({ ...runtimeConfig })
    runtimeConfig.level = saved.level
    runtimeConfig.enable_sampling = saved.enable_sampling
    runtimeConfig.sampling_initial = saved.sampling_initial
    runtimeConfig.sampling_thereafter = saved.sampling_thereafter
    runtimeConfig.caller = saved.caller
    runtimeConfig.stacktrace_level = saved.stacktrace_level
    runtimeConfig.retention_days = saved.retention_days
    appStore.showSuccess('日志运行时配置已生效')
  } catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to save runtime log config', err)
    appStore.showError(err?.response?.data?.detail || '保存日志配置失败')
  } finally {
    runtimeSaving.value = false
  }
}

const resetRuntimeConfig = async () => {
  const ok = window.confirm('确认回滚为启动配置（env/yaml）并立即生效？')
  if (!ok) return

  runtimeSaving.value = true
  try {
    const saved = await opsAPI.resetRuntimeLogConfig()
    runtimeConfig.level = saved.level
    runtimeConfig.enable_sampling = saved.enable_sampling
    runtimeConfig.sampling_initial = saved.sampling_initial
    runtimeConfig.sampling_thereafter = saved.sampling_thereafter
    runtimeConfig.caller = saved.caller
    runtimeConfig.stacktrace_level = saved.stacktrace_level
    runtimeConfig.retention_days = saved.retention_days
    appStore.showSuccess('已回滚到启动日志配置')
    await fetchHealth()
  } catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to reset runtime log config', err)
    appStore.showError(err?.response?.data?.detail || '回滚日志配置失败')
  } finally {
    runtimeSaving.value = false
  }
}

const cleanupCurrentFilter = async () => {
  const ok = window.confirm('确认按当前筛选条件清理系统日志？该操作不可撤销。')
  if (!ok) return
  try {
    const payload = {
      start_time: toRFC3339(filters.start_time),
      end_time: toRFC3339(filters.end_time),
      level: filters.level.trim() || undefined,
      component: filters.component.trim() || undefined,
      request_id: filters.request_id.trim() || undefined,
      client_request_id: filters.client_request_id.trim() || undefined,
      user_id: filters.user_id.trim() ? Number.parseInt(filters.user_id.trim(), 10) : undefined,
      account_id: filters.account_id.trim() ? Number.parseInt(filters.account_id.trim(), 10) : undefined,
      platform: filters.platform.trim() || undefined,
      model: filters.model.trim() || undefined,
      q: filters.q.trim() || undefined
    }
    const res = await opsAPI.cleanupSystemLogs(payload)
    appStore.showSuccess(`清理完成，删除 ${res.deleted || 0} 条日志`)
    page.value = 1
    await Promise.all([fetchLogs(), fetchHealth()])
  } catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to cleanup logs', err)
    appStore.showError(err?.response?.data?.detail || '清理系统日志失败')
  }
}

const resetFilters = () => {
  filters.time_range = '1h'
  filters.start_time = ''
  filters.end_time = ''
  filters.level = ''
  filters.component = ''
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.account_id = ''
  filters.platform = props.platformFilter || ''
  filters.model = ''
  filters.q = ''
  page.value = 1
  fetchLogs()
}

watch(() => props.platformFilter, (v) => {
  if (v && !filters.platform) {
    filters.platform = v
    page.value = 1
    fetchLogs()
  }
})

watch(() => props.refreshToken, () => {
  fetchLogs()
  fetchHealth()
})

const onPageChange = (next: number) => {
  page.value = next
  fetchLogs()
}

const onPageSizeChange = (next: number) => {
  pageSize.value = next
  page.value = 1
  fetchLogs()
}

const applyFilters = () => {
  page.value = 1
  fetchLogs()
}

const showCacheInstrumentation = () => {
  filters.time_range = '24h'
  filters.level = ''
  filters.component = 'cache.instrumentation'
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.account_id = ''
  filters.model = ''
  filters.q = ''
  page.value = 1
  fetchLogs()
}

const showCacheStickyDrift = () => {
  filters.time_range = '24h'
  filters.level = ''
  filters.component = 'cache.instrumentation'
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.account_id = ''
  filters.model = ''
  filters.q = 'account_switch_happened'
  page.value = 1
  fetchLogs()
}

const showCachePreviousResponse = () => {
  filters.time_range = '24h'
  filters.level = ''
  filters.component = 'cache.instrumentation'
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.account_id = ''
  filters.model = ''
  filters.q = 'previous_response_id_present'
  page.value = 1
  fetchLogs()
}

const showCachePromptKeyGaps = () => {
  filters.time_range = '24h'
  filters.level = ''
  filters.component = 'cache.instrumentation'
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.account_id = ''
  filters.model = ''
  filters.q = 'prompt_cache_key_present\":false'
  page.value = 1
  fetchLogs()
}

const showRoutingExplanations = () => {
  filters.time_range = '24h'
  filters.level = ''
  filters.component = 'routing.explanation'
  filters.q = ''
  page.value = 1
  fetchLogs()
}

const hasData = computed(() => logs.value.length > 0)
const cacheMode = computed(() => filters.component.trim() === 'cache.instrumentation')
const cacheSummary = computed(() => buildCacheInstrumentationSummary(logs.value))
const cacheHitRateLabel = computed(() => {
  const totalSamples = cacheSummary.value.total
  if (totalSamples <= 0) return '-'
  return `${Math.round((cacheSummary.value.cacheHitCount / totalSamples) * 100)}%`
})
const formatSummaryList = (items: Array<{ label: string; count: number }>) =>
  items.length > 0 ? items.map((item) => `${item.label} (${item.count})`).join(' / ') : '-'
const activeQuickView = computed(() => {
  const component = filters.component.trim()
  const q = filters.q.trim()
  if (component === 'cache.instrumentation' && q === '') return 'cache_all'
  if (component === 'cache.instrumentation' && q === 'account_switch_happened') return 'cache_switch'
  if (component === 'cache.instrumentation' && q === 'previous_response_id_present') return 'cache_previous_response'
  if (component === 'cache.instrumentation' && q === 'prompt_cache_key_present\":false') return 'cache_prompt_gap'
  if (component === 'routing.explanation') return 'routing'
  return 'custom'
})
const quickViews = computed(() => [
  {
    id: 'cache_all',
    title: '缓存诊断总览',
    description: '先看最近 24h 的 cache 样本，再决定是查会话漂移、key 缺失，还是 continuation 没接上。',
    tone: 'emerald',
    action: showCacheInstrumentation
  },
  {
    id: 'cache_switch',
    title: '会话漂移 / 切账号',
    description: '优先看 `account_switch_happened` 和 sticky miss，判断缓存失败是不是调度层导致。',
    tone: 'amber',
    action: showCacheStickyDrift
  },
  {
    id: 'cache_previous_response',
    title: 'Continuation 接续',
    description: '聚焦带 `previous_response_id` 的样本，确认 continuation 是否真的接上。',
    tone: 'sky',
    action: showCachePreviousResponse
  },
  {
    id: 'cache_prompt_gap',
    title: 'Prompt Cache Key 缺口',
    description: '直接筛缺少 `prompt_cache_key` 的样本，判断问题在客户端透传还是网关注入。',
    tone: 'rose',
    action: showCachePromptKeyGaps
  },
  {
    id: 'routing',
    title: '调度解释回看',
    description: '如果缓存问题像是账号选择导致，再切到路由解释看 account 为什么被选中。',
    tone: 'slate',
    action: showRoutingExplanations
  }
])
const cacheFocusNotes = computed(() => {
  if (!cacheMode.value || cacheSummary.value.total <= 0) return []
  const notes: string[] = []
  if (cacheSummary.value.accountSwitchCount > 0) notes.push(`有 ${cacheSummary.value.accountSwitchCount} 条样本发生账号切换，优先排查 sticky/session 漂移。`)
  if (cacheSummary.value.promptCacheKeyMissingCount > 0) notes.push(`有 ${cacheSummary.value.promptCacheKeyMissingCount} 条样本缺少 prompt_cache_key，建议先看客户端透传或自动注入链路。`)
  if (cacheSummary.value.previousResponseCount > 0) notes.push(`有 ${cacheSummary.value.previousResponseCount} 条样本带 previous_response_id，可继续核对 continuation 命中质量。`)
  if (cacheSummary.value.cacheHitCount === 0) notes.push('当前结果页没有 cache_read 命中样本，建议优先看 signal source 和切账号情况。')
  return notes.slice(0, 3)
})
const quickViewCardClass = (id: string, tone: string) => {
  const active = activeQuickView.value === id
  const base = 'rounded-xl border px-3 py-3 text-left transition-colors'
  if (tone === 'emerald') {
    return `${base} ${active ? 'border-emerald-400 bg-emerald-50 dark:border-emerald-700 dark:bg-emerald-950/20' : 'border-emerald-100 bg-white hover:border-emerald-300 dark:border-emerald-900/30 dark:bg-dark-900/40'}`
  }
  if (tone === 'amber') {
    return `${base} ${active ? 'border-amber-400 bg-amber-50 dark:border-amber-700 dark:bg-amber-950/20' : 'border-amber-100 bg-white hover:border-amber-300 dark:border-amber-900/30 dark:bg-dark-900/40'}`
  }
  if (tone === 'sky') {
    return `${base} ${active ? 'border-sky-400 bg-sky-50 dark:border-sky-700 dark:bg-sky-950/20' : 'border-sky-100 bg-white hover:border-sky-300 dark:border-sky-900/30 dark:bg-dark-900/40'}`
  }
  if (tone === 'rose') {
    return `${base} ${active ? 'border-rose-400 bg-rose-50 dark:border-rose-700 dark:bg-rose-950/20' : 'border-rose-100 bg-white hover:border-rose-300 dark:border-rose-900/30 dark:bg-dark-900/40'}`
  }
  return `${base} ${active ? 'border-slate-400 bg-slate-50 dark:border-slate-700 dark:bg-slate-900/30' : 'border-slate-200 bg-white hover:border-slate-300 dark:border-dark-700 dark:bg-dark-900/40'}`
}

onMounted(async () => {
  if (props.platformFilter) {
    filters.platform = props.platformFilter
  }
  await Promise.all([fetchLogs(), fetchHealth(), loadRuntimeConfig()])
})
</script>

<template>
  <section class="ops-system-log-shell rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900/60">
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h3 class="text-[13px] font-semibold text-gray-900 dark:text-white">系统日志</h3>
        <p class="mt-1 text-[11px] text-gray-500 dark:text-gray-400">默认按最新时间倒序，支持筛选搜索与按条件清理；调度解释用于查看账号为什么被选中。</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 text-xs">
        <span class="rounded-md bg-gray-100 px-2 py-1 text-gray-700 dark:bg-dark-700 dark:text-gray-200">队列 {{ health.queue_depth }}/{{ health.queue_capacity }}</span>
        <span class="rounded-md bg-gray-100 px-2 py-1 text-gray-700 dark:bg-dark-700 dark:text-gray-200">写入 {{ health.written_count }}</span>
        <span class="rounded-md bg-amber-100 px-2 py-1 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300">丢弃 {{ health.dropped_count }}</span>
        <span class="rounded-md bg-red-100 px-2 py-1 text-red-700 dark:bg-red-900/30 dark:text-red-300">失败 {{ health.write_failed_count }}</span>
      </div>
    </div>

    <div class="ops-system-log-runtime-panel mb-4 rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800/70">
      <div class="mb-2 flex items-center justify-between">
        <div class="text-[11px] font-semibold text-gray-700 dark:text-gray-200">运行时日志配置（实时生效）</div>
        <span v-if="runtimeLoading" class="text-xs text-gray-500">加载中...</span>
      </div>
      <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-6">
        <label class="text-xs text-gray-600 dark:text-gray-300">
          级别
          <Select v-model="runtimeConfig.level" class="mt-1" :options="runtimeLevelOptions" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          堆栈阈值
          <Select v-model="runtimeConfig.stacktrace_level" class="mt-1" :options="stacktraceLevelOptions" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          采样初始
          <input v-model.number="runtimeConfig.sampling_initial" type="number" min="1" class="input mt-1" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          采样后续
          <input v-model.number="runtimeConfig.sampling_thereafter" type="number" min="1" class="input mt-1" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          保留天数
          <input v-model.number="runtimeConfig.retention_days" type="number" min="1" max="3650" class="input mt-1" />
        </label>
        <div class="md:col-span-2 xl:col-span-6">
          <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-end">
            <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
              <label class="inline-flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                <input v-model="runtimeConfig.caller" type="checkbox" />
                caller
              </label>
              <label class="inline-flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                <input v-model="runtimeConfig.enable_sampling" type="checkbox" />
                sampling
              </label>
            </div>
            <div class="flex flex-wrap items-center gap-2 lg:justify-end">
              <button type="button" class="btn btn-primary btn-sm" :disabled="runtimeSaving" @click="saveRuntimeConfig">
                {{ runtimeSaving ? '保存中...' : '保存并生效' }}
              </button>
              <button type="button" class="btn btn-secondary btn-sm" :disabled="runtimeSaving" @click="resetRuntimeConfig">
                回滚默认值
              </button>
            </div>
          </div>
        </div>
      </div>
      <p v-if="health.last_error" class="mt-2 text-xs text-red-600 dark:text-red-400">最近写入错误：{{ health.last_error }}</p>
    </div>

    <div class="mb-4 rounded-xl border border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/50">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <div class="text-[12px] font-semibold text-gray-900 dark:text-white">预设诊断入口</div>
          <p class="mt-1 text-[11px] leading-6 text-gray-500 dark:text-gray-400">
            先用预置视角快速定位问题，再按下方筛选补细节。当前视角：{{ activeQuickView === 'custom' ? '自定义筛选' : quickViews.find((item) => item.id === activeQuickView)?.title }}
          </p>
        </div>
        <div class="rounded-lg bg-white px-3 py-2 text-[11px] text-gray-600 shadow-sm dark:bg-dark-800 dark:text-gray-300">
          建议顺序：总览 → 切账号/漂移 → previous_response_id → prompt_cache_key 缺口
        </div>
      </div>

      <div class="mt-4 grid gap-3 lg:grid-cols-5">
        <button
          v-for="view in quickViews"
          :key="view.id"
          type="button"
          :class="quickViewCardClass(view.id, view.tone)"
          @click="view.action"
        >
          <div class="flex items-center justify-between gap-2">
            <div class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ view.title }}</div>
            <span
              v-if="activeQuickView === view.id"
              class="rounded-full bg-gray-900/90 px-2 py-0.5 text-[10px] text-white dark:bg-white/90 dark:text-gray-900"
            >
              当前
            </span>
          </div>
          <p class="mt-2 text-[11px] leading-6 text-gray-600 dark:text-gray-400">
            {{ view.description }}
          </p>
        </button>
      </div>
    </div>

    <div class="mb-4 grid grid-cols-1 gap-3 md:grid-cols-5">
      <label class="text-xs text-gray-600 dark:text-gray-300">
        时间范围
        <Select v-model="filters.time_range" class="mt-1" :options="timeRangeOptions" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        开始时间（可选）
        <input v-model="filters.start_time" type="datetime-local" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        结束时间（可选）
        <input v-model="filters.end_time" type="datetime-local" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        级别
        <Select v-model="filters.level" class="mt-1" :options="filterLevelOptions" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        组件
        <input v-model="filters.component" type="text" class="input mt-1" placeholder="如 http.access" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        request_id
        <input v-model="filters.request_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        client_request_id
        <input v-model="filters.client_request_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        user_id
        <input v-model="filters.user_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        account_id
        <input v-model="filters.account_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        平台
        <input v-model="filters.platform" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        模型
        <input v-model="filters.model" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        关键词
        <input v-model="filters.q" type="text" class="input mt-1" placeholder="消息/request_id" />
      </label>
    </div>

    <div class="mb-3 flex flex-wrap gap-2">
      <button type="button" class="btn btn-primary btn-sm" @click="applyFilters">查询</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="showCacheInstrumentation">只看缓存诊断</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="showCacheStickyDrift">只看会话漂移/切账号</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="showCachePreviousResponse">只看 previous_response_id</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="showCachePromptKeyGaps">只看缺少 prompt_cache_key</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="showRoutingExplanations">只看调度解释</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="resetFilters">重置</button>
      <button type="button" class="btn btn-danger btn-sm" @click="cleanupCurrentFilter">按当前筛选清理</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="fetchHealth">刷新健康指标</button>
    </div>

    <div
      v-if="cacheMode"
      class="mb-4 rounded-xl border border-amber-200 bg-amber-50/70 p-4 dark:border-amber-900/40 dark:bg-amber-950/10"
    >
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div>
          <div class="text-sm font-semibold text-amber-900 dark:text-amber-200">缓存诊断视角</div>
          <p class="mt-1 max-w-3xl text-xs leading-6 text-amber-800/90 dark:text-amber-100/80">
            这批样本来自 `cache.instrumentation` 采样事件，适合先看会不会话漂移、`prompt_cache_key` 是否缺失、以及 `previous_response_id` 是否真正接上。
            当前表格仍按时间倒序，下面的摘要只基于当前筛选结果页。
          </p>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 text-xs text-amber-900 shadow-sm dark:bg-amber-950/20 dark:text-amber-100">
          当前样本 {{ cacheSummary.total }} 条
        </div>
      </div>

      <div class="mt-4 grid grid-cols-2 gap-3 xl:grid-cols-6">
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">缓存读取命中</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.cacheHitCount }} / {{ cacheHitRateLabel }}</div>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">账号切换</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.accountSwitchCount }}</div>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">粘性未命中</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.stickyMissCount }}</div>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">带 previous_response_id</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.previousResponseCount }}</div>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">缺少 prompt_cache_key</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.promptCacheKeyMissingCount }}</div>
        </div>
        <div class="rounded-lg bg-white/80 px-3 py-2 shadow-sm dark:bg-amber-950/20">
          <div class="text-[11px] text-gray-500 dark:text-amber-100/60">自动派生 key</div>
          <div class="mt-1 text-sm font-semibold text-gray-900 dark:text-amber-50">{{ cacheSummary.autoInjectedCount }}</div>
        </div>
      </div>

      <div class="mt-3 grid gap-3 xl:grid-cols-2">
        <div class="rounded-lg bg-white/70 px-3 py-2 text-xs text-gray-700 shadow-sm dark:bg-amber-950/15 dark:text-amber-100/80">
          <span class="font-medium text-gray-900 dark:text-amber-50">cache_family</span>
          <span class="ml-2">{{ formatSummaryList(cacheSummary.families) }}</span>
        </div>
        <div class="rounded-lg bg-white/70 px-3 py-2 text-xs text-gray-700 shadow-sm dark:bg-amber-950/15 dark:text-amber-100/80">
          <span class="font-medium text-gray-900 dark:text-amber-50">session_signal_source</span>
          <span class="ml-2">{{ formatSummaryList(cacheSummary.signalSources) }}</span>
        </div>
      </div>

      <div v-if="cacheFocusNotes.length > 0" class="mt-3 grid gap-2">
        <div
          v-for="note in cacheFocusNotes"
          :key="note"
          class="rounded-lg border border-white/70 bg-white/70 px-3 py-2 text-xs text-amber-900 shadow-sm dark:border-amber-900/30 dark:bg-amber-950/15 dark:text-amber-100"
        >
          {{ note }}
        </div>
      </div>
    </div>

    <div class="ops-system-log-table-wrap overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
      <div v-if="loading" class="px-4 py-8 text-center text-sm text-gray-500">加载中...</div>
      <div v-else-if="!hasData" class="px-4 py-8 text-center text-sm text-gray-500">暂无系统日志</div>
      <div v-else class="overflow-auto">
        <table class="min-w-full table-fixed divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="bg-gray-50 dark:bg-dark-900">
            <tr>
              <th class="w-[170px] px-3 py-2 text-left text-[11px] font-semibold text-gray-500">时间</th>
              <th class="w-[80px] px-3 py-2 text-left text-[11px] font-semibold text-gray-500">级别</th>
              <th class="px-3 py-2 text-left text-[11px] font-semibold text-gray-500">日志详细信息</th>
            </tr>
          </thead>
          <tbody class="ops-system-log-table-body divide-y divide-gray-100 dark:divide-dark-800">
            <tr
              v-for="row in logs"
              :key="row.id"
              class="align-top"
              :class="isCacheInstrumentationLog(row)
                ? 'bg-emerald-50/60 dark:bg-emerald-950/10'
                : (isRoutingExplanationLog(row) ? 'bg-amber-50/60 dark:bg-amber-950/10' : '')"
            >
              <td class="px-3 py-2 text-xs text-gray-700 dark:text-gray-300">{{ formatTime(row.created_at) }}</td>
              <td class="px-3 py-2 text-xs">
                <span class="inline-flex rounded-full px-2 py-0.5 font-semibold" :class="levelBadgeClass(row.level)">
                  {{ row.level }}
                </span>
              </td>
              <td class="px-3 py-2 text-xs text-gray-700 dark:text-gray-300 whitespace-normal break-all">
                {{ formatSystemLogDetail(row) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        :total="total"
        :page="page"
        :page-size="pageSize"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </div>
  </section>
</template>

<style>
.dark .ops-system-log-shell {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.88) !important;
}

.dark .ops-system-log-shell :is(.bg-white, .bg-gray-50, .bg-gray-100) {
  background-color: rgba(24, 26, 21, 0.76) !important;
}

.dark .ops-system-log-runtime-panel {
  background: rgba(22, 24, 19, 0.96) !important;
  border-color: rgba(48, 52, 43, 0.82) !important;
}

.dark .ops-system-log-table-wrap {
  background: rgba(22, 24, 19, 0.98) !important;
  border-color: rgba(48, 52, 43, 0.82) !important;
}

.dark .ops-system-log-table-body {
  background: rgba(24, 26, 21, 0.98) !important;
}

.dark .ops-system-log-shell :is(.text-gray-900, .text-gray-800, .text-gray-700) {
  color: #f4efe4 !important;
}

.dark .ops-system-log-shell :is(.text-gray-600, .text-gray-500, .text-gray-400) {
  color: #bdb5a8 !important;
}

.dark .ops-system-log-shell thead {
  background: rgba(17, 19, 15, 0.68) !important;
}

.dark .ops-system-log-shell th {
  border-color: rgba(48, 52, 43, 0.82) !important;
  color: #b5aa94 !important;
}

.dark .ops-system-log-shell td {
  border-color: rgba(48, 52, 43, 0.78) !important;
  color: #d7d0c2 !important;
}

.dark .ops-system-log-shell tbody tr:hover {
  background: rgba(167, 58, 42, 0.08) !important;
}
</style>
