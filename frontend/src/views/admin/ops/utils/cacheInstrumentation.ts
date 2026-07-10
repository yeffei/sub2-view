import type { OpsSystemLog } from '@/api/admin/ops'

export interface CacheInstrumentationSummaryItem {
  label: string
  count: number
}

export interface CacheInstrumentationSummary {
  total: number
  cacheHitCount: number
  accountSwitchCount: number
  stickyMissCount: number
  previousResponseCount: number
  promptCacheKeyMissingCount: number
  autoInjectedCount: number
  families: CacheInstrumentationSummaryItem[]
  signalSources: CacheInstrumentationSummaryItem[]
}

function getExtra(row: OpsSystemLog): Record<string, any> {
  return row.extra ?? {}
}

function getString(extra: Record<string, any>, key: string): string {
  const value = extra[key]
  if (value == null) return ''
  if (typeof value === 'string') return value.trim()
  if (typeof value === 'number' || typeof value === 'boolean') return String(value)
  return ''
}

function getBool(extra: Record<string, any>, key: string): boolean {
  const value = extra[key]
  if (typeof value === 'boolean') return value
  if (typeof value === 'number') return value !== 0
  if (typeof value === 'string') {
    const normalized = value.trim().toLowerCase()
    return normalized === 'true' || normalized === '1' || normalized === 'yes'
  }
  return false
}

function getNumber(extra: Record<string, any>, key: string): number | null {
  const value = extra[key]
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string' && value.trim()) {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }
  return null
}

function formatStickyStatus(extra: Record<string, any>): string {
  const stickyAccountID = getNumber(extra, 'sticky_account_id')
  const stickyHit = getBool(extra, 'sticky_account_hit')
  if (!stickyAccountID || stickyAccountID <= 0) return 'sticky=none'
  return `sticky=${stickyHit ? 'hit' : 'miss'}(${stickyAccountID})`
}

function formatFocus(extra: Record<string, any>): string {
  const focus: string[] = []
  if (getBool(extra, 'account_switch_happened')) focus.push('账号切换')
  const stickyAccountID = getNumber(extra, 'sticky_account_id')
  if (stickyAccountID && stickyAccountID > 0 && !getBool(extra, 'sticky_account_hit')) {
    focus.push('粘性未命中')
  }
  if (!getBool(extra, 'prompt_cache_key_present')) focus.push('缺少 prompt_cache_key')
  if (getBool(extra, 'prompt_cache_key_auto_injected')) focus.push('使用自动派生 key')
  if (focus.length === 0) return ''
  return `关注=${focus.join(' / ')}`
}

export function isCacheInstrumentationLog(row: OpsSystemLog): boolean {
  const component = String(row.component || row.extra?.component || '').trim()
  return component === 'cache.instrumentation' || String(row.message || '').trim() === 'cache_instrumentation'
}

export function formatCacheInstrumentationDetail(row: OpsSystemLog): string {
  const extra = getExtra(row)
  const family = getString(extra, 'cache_family') || '-'
  const session = getString(extra, 'session_hash_short') || '-'
  const signalSource = getString(extra, 'session_signal_source') || '-'
  const promptKeyPresent = getBool(extra, 'prompt_cache_key_present')
  const promptKeyAutoInjected = getBool(extra, 'prompt_cache_key_auto_injected')
  const previousResponse = getBool(extra, 'previous_response_id_present')
  const accountSwitch = getBool(extra, 'account_switch_happened')
  const selectedAccountID = getString(extra, 'selected_account_id') || String(row.account_id || '')
  const cacheReadTokens = getString(extra, 'cache_read_tokens') || '0'
  const cacheCreationTokens = getString(extra, 'cache_creation_tokens') || '0'
  const requestType = getString(extra, 'request_type') || '-'
  const requestedModel = getString(extra, 'requested_model') || row.model || '-'
  const upstreamModel = getString(extra, 'upstream_model') || '-'
  const inboundEndpoint = getString(extra, 'inbound_endpoint') || '-'
  const upstreamEndpoint = getString(extra, 'upstream_endpoint') || '-'
  const sampleReason = getString(extra, 'sample_reason')
  const sampleRate = getString(extra, 'sample_rate')

  const parts = [
    `缓存诊断 family=${family}`,
    `session=${session}`,
    `source=${signalSource}`,
    `prompt_key=${promptKeyPresent ? 'yes' : 'no'}${promptKeyAutoInjected ? '(auto)' : ''}`,
    `previous_response_id=${previousResponse ? 'yes' : 'no'}`,
    formatStickyStatus(extra),
    `switch=${accountSwitch ? 'yes' : 'no'}`,
    `selected=${selectedAccountID || '-'}`,
    `cache_read=${cacheReadTokens}`,
    `cache_create=${cacheCreationTokens}`,
    `request_type=${requestType}`,
    `requested=${requestedModel}`,
    `upstream=${upstreamModel}`,
    `inbound=${inboundEndpoint}`,
    `upstream_endpoint=${upstreamEndpoint}`
  ]

  const focus = formatFocus(extra)
  if (focus) parts.push(focus)
  if (sampleReason) parts.push(`sample=${sampleReason}${sampleRate ? `(${sampleRate})` : ''}`)

  return parts.join('  ')
}

export function buildCacheInstrumentationSummary(rows: OpsSystemLog[]): CacheInstrumentationSummary {
  const cacheRows = rows.filter(isCacheInstrumentationLog)
  const familyCount = new Map<string, number>()
  const signalSourceCount = new Map<string, number>()

  let cacheHitCount = 0
  let accountSwitchCount = 0
  let stickyMissCount = 0
  let previousResponseCount = 0
  let promptCacheKeyMissingCount = 0
  let autoInjectedCount = 0

  for (const row of cacheRows) {
    const extra = getExtra(row)
    const family = getString(extra, 'cache_family') || 'unknown'
    const signalSource = getString(extra, 'session_signal_source') || 'unknown'
    familyCount.set(family, (familyCount.get(family) ?? 0) + 1)
    signalSourceCount.set(signalSource, (signalSourceCount.get(signalSource) ?? 0) + 1)

    if ((getNumber(extra, 'cache_read_tokens') ?? 0) > 0) cacheHitCount++
    if (getBool(extra, 'account_switch_happened')) accountSwitchCount++
    if ((getNumber(extra, 'sticky_account_id') ?? 0) > 0 && !getBool(extra, 'sticky_account_hit')) stickyMissCount++
    if (getBool(extra, 'previous_response_id_present')) previousResponseCount++
    if (!getBool(extra, 'prompt_cache_key_present')) promptCacheKeyMissingCount++
    if (getBool(extra, 'prompt_cache_key_auto_injected')) autoInjectedCount++
  }

  const toItems = (input: Map<string, number>) =>
    Array.from(input.entries())
      .sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))
      .slice(0, 4)
      .map(([label, count]) => ({ label, count }))

  return {
    total: cacheRows.length,
    cacheHitCount,
    accountSwitchCount,
    stickyMissCount,
    previousResponseCount,
    promptCacheKeyMissingCount,
    autoInjectedCount,
    families: toItems(familyCount),
    signalSources: toItems(signalSourceCount)
  }
}
