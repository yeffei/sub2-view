import type {
  OpenAIQuotaUsage,
  OpenAIRateLimitWindow
} from '@/api/admin/accounts'
import type { WindowStats } from '@/types'

const HOUR_SECONDS = 60 * 60
const DAY_SECONDS = 24 * HOUR_SECONDS

export interface OpenAIQuotaProgressBar {
  key: 'primary' | 'secondary'
  label: string
  utilization: number
  resetsAt: string | null
  windowSeconds: number
  windowStats?: WindowStats | null
  color: 'indigo' | 'emerald' | 'purple' | 'amber'
}

export const formatOpenAIQuotaWindowLabel = (seconds: number): string => {
  if (!Number.isFinite(seconds) || seconds <= 0) return '-'

  if (seconds >= DAY_SECONDS) {
    return `${Math.max(1, Math.round(seconds / DAY_SECONDS))}d`
  }
  if (seconds >= HOUR_SECONDS) {
    return `${Math.max(1, Math.round(seconds / HOUR_SECONDS))}h`
  }
  return `${Math.max(1, Math.round(seconds / 60))}m`
}

const resolveResetAt = (
  quota: OpenAIQuotaUsage,
  window: OpenAIRateLimitWindow
): string | null => {
  const resetAt = Number(window.reset_at)
  if (Number.isFinite(resetAt) && resetAt > 0) {
    return new Date(resetAt * 1000).toISOString()
  }

  const fetchedAt = Number(quota.fetched_at)
  const resetAfter = Number(window.reset_after_seconds)
  if (Number.isFinite(fetchedAt) && fetchedAt > 0 && Number.isFinite(resetAfter) && resetAfter >= 0) {
    return new Date((fetchedAt + resetAfter) * 1000).toISOString()
  }

  return null
}

const resolveWindowStats = (
  seconds: number,
  fiveHourStats?: WindowStats | null,
  sevenDayStats?: WindowStats | null
): WindowStats | null | undefined => {
  // The local /usage endpoint currently exposes exact 5h and 7d aggregates.
  // Only attach them when the upstream quota window has the same semantics.
  if (seconds >= 4 * HOUR_SECONDS && seconds <= 6 * HOUR_SECONDS) {
    return fiveHourStats
  }
  if (seconds >= 6 * DAY_SECONDS && seconds <= 8 * DAY_SECONDS) {
    return sevenDayStats
  }
  return undefined
}

export const buildOpenAIQuotaProgressBars = (
  quota: OpenAIQuotaUsage | null | undefined,
  stats?: {
    fiveHour?: WindowStats | null
    sevenDay?: WindowStats | null
  }
): OpenAIQuotaProgressBar[] => {
  const rateLimit = quota?.rate_limit
  if (!rateLimit) return []

  const windows: Array<{
    key: 'primary' | 'secondary'
    window: OpenAIRateLimitWindow
  }> = []

  if (rateLimit.primary_window && rateLimit.primary_window.limit_window_seconds > 0) {
    windows.push({ key: 'primary', window: rateLimit.primary_window })
  }
  if (rateLimit.secondary_window && rateLimit.secondary_window.limit_window_seconds > 0) {
    windows.push({ key: 'secondary', window: rateLimit.secondary_window })
  }

  const colors: OpenAIQuotaProgressBar['color'][] = ['indigo', 'emerald', 'purple', 'amber']

  const bars = windows
    .sort((a, b) => a.window.limit_window_seconds - b.window.limit_window_seconds)
    .map(({ key, window }, index) => ({
      key,
      label: formatOpenAIQuotaWindowLabel(window.limit_window_seconds),
      utilization: Number.isFinite(window.used_percent) ? window.used_percent : 0,
      resetsAt: resolveResetAt(quota, window),
      windowSeconds: window.limit_window_seconds,
      windowStats: resolveWindowStats(
        window.limit_window_seconds,
        stats?.fiveHour,
        stats?.sevenDay
      ),
      color: colors[index] ?? 'amber'
    }))

  // Keep the existing local activity badges visible in compact account rows.
  // They are local 5h traffic stats; the progress percentage and reset time
  // still come exclusively from the real upstream quota window.
  if (bars.length > 0 && !bars[0]?.windowStats && stats?.fiveHour) {
    bars[0]!.windowStats = stats.fiveHour
  }

  return bars
}
