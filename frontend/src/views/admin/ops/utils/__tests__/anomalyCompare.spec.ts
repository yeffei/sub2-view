import { describe, expect, it } from 'vitest'
import type {
  OpsDashboardOverview,
  OpsErrorDistributionResponse,
  OpsErrorTrendResponse,
  OpsThroughputTrendResponse
} from '@/api/admin/ops'
import { buildAnomalyCompareSummary, buildPreviousWindowRange } from '../anomalyCompare'

function makeOverview(overrides: Partial<OpsDashboardOverview> = {}): OpsDashboardOverview {
  return {
    start_time: '2026-06-23T00:00:00.000Z',
    end_time: '2026-06-23T01:00:00.000Z',
    platform: '',
    success_count: 900,
    error_count_total: 100,
    business_limited_count: 30,
    error_count_sla: 70,
    request_count_total: 1000,
    request_count_sla: 970,
    token_consumed: 10000,
    sla: 0.97,
    error_rate: 0.1,
    upstream_error_rate: 0.03,
    upstream_error_count_excl_429_529: 20,
    upstream_429_count: 5,
    upstream_529_count: 5,
    qps: { current: 1, peak: 2, avg: 1.5 },
    tps: { current: 10, peak: 20, avg: 15 },
    duration: {},
    ttft: {},
    ...overrides
  }
}

function makeTrend(points: OpsErrorTrendResponse['points']): OpsErrorTrendResponse {
  return { bucket: '5m', points }
}

function makeDistribution(items: OpsErrorDistributionResponse['items']): OpsErrorDistributionResponse {
  return {
    total: items.reduce((sum, item) => sum + item.total, 0),
    items
  }
}

function makeThroughput(byPlatform: OpsThroughputTrendResponse['by_platform']): OpsThroughputTrendResponse {
  return {
    bucket: '5m',
    points: [],
    by_platform: byPlatform ?? []
  }
}

describe('anomalyCompare', () => {
  it('builds previous equal-length window for preset ranges', () => {
    const now = new Date('2026-06-23T12:00:00.000Z')
    const range = buildPreviousWindowRange('1h', null, null, now)

    expect(range.startTime).toBe('2026-06-23T10:00:00.000Z')
    expect(range.endTime).toBe('2026-06-23T11:00:00.000Z')
  })

  it('builds previous equal-length window for custom ranges', () => {
    const range = buildPreviousWindowRange('custom', '2026-06-23T09:00:00.000Z', '2026-06-23T11:30:00.000Z')

    expect(range.startTime).toBe('2026-06-23T06:30:00.000Z')
    expect(range.endTime).toBe('2026-06-23T09:00:00.000Z')
  })

  it('summarizes rate deltas and structure shifts', () => {
    const current = {
      overview: makeOverview({ request_count_total: 1000, error_rate: 0.12 }),
      errorTrend: makeTrend([
        {
          bucket_start: '2026-06-23T00:00:00.000Z',
          error_count_total: 120,
          business_limited_count: 30,
          error_count_sla: 60,
          upstream_error_count_excl_429_529: 20,
          upstream_429_count: 10,
          upstream_529_count: 5
        }
      ]),
      errorDistribution: makeDistribution([
        { status_code: 500, total: 40, sla: 40, business_limited: 0 },
        { status_code: 429, total: 20, sla: 20, business_limited: 0 },
        { status_code: 400, total: 10, sla: 10, business_limited: 0 }
      ]),
      throughputTrend: makeThroughput([
        { platform: 'openai', request_count: 600, token_consumed: 0 },
        { platform: 'claude', request_count: 400, token_consumed: 0 }
      ])
    }

    const previous = {
      overview: makeOverview({ request_count_total: 800, error_rate: 0.08 }),
      errorTrend: makeTrend([
        {
          bucket_start: '2026-06-23T00:00:00.000Z',
          error_count_total: 64,
          business_limited_count: 8,
          error_count_sla: 24,
          upstream_error_count_excl_429_529: 12,
          upstream_429_count: 4,
          upstream_529_count: 0
        }
      ]),
      errorDistribution: makeDistribution([
        { status_code: 400, total: 18, sla: 18, business_limited: 0 },
        { status_code: 404, total: 10, sla: 10, business_limited: 0 },
        { status_code: 502, total: 6, sla: 6, business_limited: 0 }
      ]),
      throughputTrend: makeThroughput([
        { platform: 'claude', request_count: 500, token_consumed: 0 },
        { platform: 'openai', request_count: 300, token_consumed: 0 }
      ])
    }

    const summary = buildAnomalyCompareSummary(current, previous)

    expect(summary.errorRate.currentValue).toBe(12)
    expect(summary.errorRate.previousValue).toBe(8)
    expect(summary.requestErrorShare.currentValue).toBe(6)
    expect(summary.requestErrorShare.previousValue).toBe(3)
    expect(summary.upstreamErrorShare.currentValue).toBeCloseTo(3.5)
    expect(summary.businessLimitedShare.currentValue).toBe(3)
    expect(summary.statusShift.currentLabel).toBe('5xx')
    expect(summary.statusShift.previousLabel).toBe('4xx')
    expect(summary.platformShift.currentLabel).toBe('openai')
    expect(summary.platformShift.previousLabel).toBe('claude')
  })
})
