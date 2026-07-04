import { describe, expect, it, vi } from 'vitest'

import { deriveAccountAnomalyReason, isAccountAuthFailure } from '../accountAnomaly'

describe('accountAnomaly', () => {
  it('detects auth failures from common upstream error text', () => {
    expect(isAccountAuthFailure('401 Unauthorized: token expired')).toBe(true)
    expect(isAccountAuthFailure('Invalid API key provided')).toBe(true)
    expect(isAccountAuthFailure('rate limit exceeded')).toBe(false)
  })

  it('prioritizes current runtime blockers over generic error state', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-07-03T12:00:00Z'))

    expect(deriveAccountAnomalyReason({
      status: 'error',
      error_message: '401 unauthorized',
      schedulable: false,
      auto_pause_on_expired: false,
      expires_at: null,
      temp_unschedulable_until: '2026-07-03T12:30:00Z',
      temp_unschedulable_reason: '429 overload rule',
      rate_limit_reset_at: null,
      overload_until: null,
      quota_used: null,
      quota_limit: null,
      quota_daily_used: null,
      quota_daily_limit: null,
      quota_weekly_used: null,
      quota_weekly_limit: null
    }).code).toBe('temp_unschedulable')

    vi.useRealTimers()
  })

  it('returns quota_exhausted for active accounts that hit quota limits', () => {
    expect(deriveAccountAnomalyReason({
      status: 'active',
      error_message: null,
      schedulable: true,
      auto_pause_on_expired: false,
      expires_at: null,
      temp_unschedulable_until: null,
      temp_unschedulable_reason: null,
      rate_limit_reset_at: null,
      overload_until: null,
      quota_used: 100,
      quota_limit: 100,
      quota_daily_used: null,
      quota_daily_limit: null,
      quota_weekly_used: null,
      quota_weekly_limit: null
    }).code).toBe('quota_exhausted')
  })
})
