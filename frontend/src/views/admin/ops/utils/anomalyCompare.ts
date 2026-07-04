import type {
  OpsDashboardOverview,
  OpsErrorDistributionResponse,
  OpsErrorTrendResponse,
  OpsThroughputTrendResponse
} from '@/api/admin/ops'
import { parseTimeRangeMinutes, sumNumbers } from './opsFormatters'

export interface CompareWindowRange {
  startTime: string
  endTime: string
}

export interface CompareMetricSummary {
  currentValue: number
  previousValue: number
  changeValue: number
  changeRatio: number | null
}

export interface CompareShiftSummary {
  currentLabel: string
  previousLabel: string
  currentValue: number
  previousValue: number
  changed: boolean
}

export interface AnomalyCompareSummary {
  errorRate: CompareMetricSummary
  requestErrorShare: CompareMetricSummary
  upstreamErrorShare: CompareMetricSummary
  businessLimitedShare: CompareMetricSummary
  statusShift: CompareShiftSummary
  platformShift: CompareShiftSummary
}

interface CompareSnapshotInput {
  overview: OpsDashboardOverview | null
  errorTrend: OpsErrorTrendResponse | null
  errorDistribution: OpsErrorDistributionResponse | null
  throughputTrend: OpsThroughputTrendResponse | null
}

interface DistributionBucket {
  label: string
  value: number
}

function safePercent(value: number, base: number): number {
  if (!Number.isFinite(value) || !Number.isFinite(base) || base <= 0) return 0
  return (value / base) * 100
}

function buildMetricSummary(currentValue: number, previousValue: number): CompareMetricSummary {
  const delta = currentValue - previousValue
  const ratio = Math.abs(previousValue) > 0 ? delta / previousValue : null
  return {
    currentValue,
    previousValue,
    changeValue: delta,
    changeRatio: ratio
  }
}

function sumUpstreamErrors(trend: OpsErrorTrendResponse | null): number {
  if (!trend) return 0
  return sumNumbers(
    trend.points.map((point) =>
      sumNumbers([
        point.upstream_error_count_excl_429_529,
        point.upstream_429_count,
        point.upstream_529_count
      ])
    )
  )
}

function sumRequestErrors(trend: OpsErrorTrendResponse | null): number {
  if (!trend) return 0
  return sumNumbers(trend.points.map((point) => point.error_count_sla ?? 0))
}

function sumBusinessLimited(trend: OpsErrorTrendResponse | null): number {
  if (!trend) return 0
  return sumNumbers(trend.points.map((point) => point.business_limited_count ?? 0))
}

function pickTopBucket(buckets: DistributionBucket[]): DistributionBucket {
  if (buckets.length === 0) return { label: '-', value: 0 }
  return buckets.reduce((max, item) => (item.value > max.value ? item : max))
}

function buildStatusBuckets(distribution: OpsErrorDistributionResponse | null): DistributionBucket[] {
  if (!distribution) return []

  let group4xx = 0
  let group5xx = 0
  let groupOther = 0

  const topStatusItems = [...(distribution.items ?? [])]
    .map((item) => ({ label: String(item.status_code || '-'), value: Number(item.sla || 0) }))
    .filter((item) => item.value > 0)
    .sort((a, b) => b.value - a.value)
    .slice(0, 2)

  for (const item of distribution.items ?? []) {
    const code = Number(item.status_code || 0)
    const count = Number(item.sla || 0)
    if (!Number.isFinite(code) || !Number.isFinite(count) || count <= 0) continue
    if (code >= 400 && code < 500) {
      group4xx += count
    } else if (code >= 500 && code < 600) {
      group5xx += count
    } else {
      groupOther += count
    }
  }

  const buckets: DistributionBucket[] = []
  if (group5xx > 0) buckets.push({ label: '5xx', value: group5xx })
  if (group4xx > 0) buckets.push({ label: '4xx', value: group4xx })
  if (groupOther > 0) buckets.push({ label: 'other', value: groupOther })

  for (const item of topStatusItems) {
    buckets.push(item)
  }

  return buckets
}

function buildPlatformBuckets(throughputTrend: OpsThroughputTrendResponse | null): DistributionBucket[] {
  return [...(throughputTrend?.by_platform ?? [])]
    .map((item) => ({
      label: item.platform || '-',
      value: Number(item.request_count || 0)
    }))
    .filter((item) => item.value > 0)
    .sort((a, b) => b.value - a.value)
}

function buildShiftSummary(currentBuckets: DistributionBucket[], previousBuckets: DistributionBucket[]): CompareShiftSummary {
  const current = pickTopBucket(currentBuckets)
  const previous = pickTopBucket(previousBuckets)
  return {
    currentLabel: current.label,
    previousLabel: previous.label,
    currentValue: current.value,
    previousValue: previous.value,
    changed: current.label !== previous.label
  }
}

export function buildPreviousWindowRange(
  timeRange: string,
  customStartTime?: string | null,
  customEndTime?: string | null,
  now = new Date()
): CompareWindowRange {
  if (timeRange === 'custom' && customStartTime && customEndTime) {
    const start = new Date(customStartTime)
    const end = new Date(customEndTime)
    const durationMs = Math.max(60_000, end.getTime() - start.getTime())
    return {
      startTime: new Date(start.getTime() - durationMs).toISOString(),
      endTime: start.toISOString()
    }
  }

  const durationMs = parseTimeRangeMinutes(timeRange) * 60 * 1000
  const endTime = now.getTime()
  return {
    startTime: new Date(endTime - durationMs * 2).toISOString(),
    endTime: new Date(endTime - durationMs).toISOString()
  }
}

export function buildAnomalyCompareSummary(
  current: CompareSnapshotInput,
  previous: CompareSnapshotInput
): AnomalyCompareSummary {
  const currentRequestTotal = Number(current.overview?.request_count_total || 0)
  const previousRequestTotal = Number(previous.overview?.request_count_total || 0)

  const currentRequestErrors = sumRequestErrors(current.errorTrend)
  const previousRequestErrors = sumRequestErrors(previous.errorTrend)
  const currentUpstreamErrors = sumUpstreamErrors(current.errorTrend)
  const previousUpstreamErrors = sumUpstreamErrors(previous.errorTrend)
  const currentBusinessLimited = sumBusinessLimited(current.errorTrend)
  const previousBusinessLimited = sumBusinessLimited(previous.errorTrend)

  return {
    errorRate: buildMetricSummary(
      Number(current.overview?.error_rate || 0) * 100,
      Number(previous.overview?.error_rate || 0) * 100
    ),
    requestErrorShare: buildMetricSummary(
      safePercent(currentRequestErrors, currentRequestTotal),
      safePercent(previousRequestErrors, previousRequestTotal)
    ),
    upstreamErrorShare: buildMetricSummary(
      safePercent(currentUpstreamErrors, currentRequestTotal),
      safePercent(previousUpstreamErrors, previousRequestTotal)
    ),
    businessLimitedShare: buildMetricSummary(
      safePercent(currentBusinessLimited, currentRequestTotal),
      safePercent(previousBusinessLimited, previousRequestTotal)
    ),
    statusShift: buildShiftSummary(
      buildStatusBuckets(current.errorDistribution),
      buildStatusBuckets(previous.errorDistribution)
    ),
    platformShift: buildShiftSummary(
      buildPlatformBuckets(current.throughputTrend),
      buildPlatformBuckets(previous.throughputTrend)
    )
  }
}
