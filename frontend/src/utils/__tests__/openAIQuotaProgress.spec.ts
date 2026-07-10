import { describe, expect, it } from 'vitest'
import {
  buildOpenAIQuotaProgressBars,
  formatOpenAIQuotaWindowLabel
} from '../openAIQuotaProgress'

describe('openAIQuotaProgress', () => {
  it('按真实窗口时长生成标签', () => {
    expect(formatOpenAIQuotaWindowLabel(5 * 60 * 60)).toBe('5h')
    expect(formatOpenAIQuotaWindowLabel(7 * 24 * 60 * 60)).toBe('7d')
    expect(formatOpenAIQuotaWindowLabel(2_628_000)).toBe('30d')
  })

  it('只有 30 天 primary window 时恢复真实 60% 进度', () => {
    const fiveHourStats = { requests: 11, tokens: 1500, cost: 1.5 }
    const bars = buildOpenAIQuotaProgressBars({
      fetched_at: 1_783_652_760,
      rate_limit: {
        allowed: true,
        limit_reached: false,
        primary_window: {
          used_percent: 60,
          limit_window_seconds: 2_628_000,
          reset_after_seconds: 2_612_728,
          reset_at: 0
        },
        secondary_window: null
      }
    }, { fiveHour: fiveHourStats })

    expect(bars).toHaveLength(1)
    expect(bars[0]).toMatchObject({
      key: 'primary',
      label: '30d',
      utilization: 60,
      windowSeconds: 2_628_000
    })
    expect(bars[0]?.resetsAt).not.toBeNull()
    expect(bars[0]?.windowStats).toEqual(fiveHourStats)
  })

  it('双窗口按短窗口优先，并只给同周期窗口挂载本地统计', () => {
    const fiveHourStats = { requests: 5, tokens: 500, cost: 0.5 }
    const sevenDayStats = { requests: 70, tokens: 7000, cost: 7 }
    const bars = buildOpenAIQuotaProgressBars({
      fetched_at: 1_783_652_760,
      rate_limit: {
        allowed: true,
        limit_reached: false,
        primary_window: {
          used_percent: 75,
          limit_window_seconds: 7 * 24 * 60 * 60,
          reset_after_seconds: 10,
          reset_at: 0
        },
        secondary_window: {
          used_percent: 25,
          limit_window_seconds: 5 * 60 * 60,
          reset_after_seconds: 20,
          reset_at: 0
        }
      }
    }, {
      fiveHour: fiveHourStats,
      sevenDay: sevenDayStats
    })

    expect(bars.map((bar) => `${bar.label}:${bar.utilization}`)).toEqual(['5h:25', '7d:75'])
    expect(bars[0]?.windowStats).toEqual(fiveHourStats)
    expect(bars[1]?.windowStats).toEqual(sevenDayStats)
  })
})
