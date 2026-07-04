import type { Account } from '@/types'

export type AccountAnomalyReasonCode =
  | 'healthy'
  | 'expired_pause'
  | 'temp_unschedulable'
  | 'rate_limited'
  | 'overloaded'
  | 'auth_failed'
  | 'error_state'
  | 'quota_exhausted'
  | 'manual_pause'
  | 'inactive'

export interface AccountAnomalyReason {
  code: AccountAnomalyReasonCode
  detail?: string
}

export const accountAnomalyReasonFilterOrder: AccountAnomalyReasonCode[] = [
  'expired_pause',
  'temp_unschedulable',
  'rate_limited',
  'overloaded',
  'auth_failed',
  'error_state',
  'quota_exhausted',
  'manual_pause',
  'inactive'
]

const authFailureKeywords = [
  '401',
  '403',
  'unauthorized',
  'forbidden',
  'invalid api key',
  'api key not found',
  'authentication',
  'token expired',
  'refresh token',
  'invalid_grant',
  'session expired',
  'login'
]

const isFutureTime = (value?: string | null) => {
  if (!value) return false
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) && timestamp > Date.now()
}

const isExpiredTime = (value?: number | null) => {
  if (!value) return false
  return value * 1000 <= Date.now()
}

const isQuotaDimensionExceeded = (used?: number | null, limit?: number | null) => (
  typeof limit === 'number' && limit > 0 && typeof used === 'number' && used >= limit
)

export const isAccountAuthFailure = (message?: string | null) => {
  const normalized = String(message || '').trim().toLowerCase()
  if (!normalized) return false
  return authFailureKeywords.some(keyword => normalized.includes(keyword))
}

export const isAccountQuotaExceeded = (account: Pick<Account, 'quota_used' | 'quota_limit' | 'quota_daily_used' | 'quota_daily_limit' | 'quota_weekly_used' | 'quota_weekly_limit'>) => (
  isQuotaDimensionExceeded(account.quota_used, account.quota_limit) ||
  isQuotaDimensionExceeded(account.quota_daily_used, account.quota_daily_limit) ||
  isQuotaDimensionExceeded(account.quota_weekly_used, account.quota_weekly_limit)
)

export const deriveAccountAnomalyReason = (
  account: Pick<Account, 'status' | 'error_message' | 'schedulable' | 'auto_pause_on_expired' | 'expires_at' | 'temp_unschedulable_until' | 'temp_unschedulable_reason' | 'rate_limit_reset_at' | 'overload_until' | 'quota_used' | 'quota_limit' | 'quota_daily_used' | 'quota_daily_limit' | 'quota_weekly_used' | 'quota_weekly_limit'>
): AccountAnomalyReason => {
  if (account.auto_pause_on_expired && isExpiredTime(account.expires_at)) {
    return { code: 'expired_pause' }
  }
  if (isFutureTime(account.temp_unschedulable_until)) {
    return {
      code: 'temp_unschedulable',
      detail: account.temp_unschedulable_reason || undefined
    }
  }
  if (isFutureTime(account.rate_limit_reset_at)) {
    return { code: 'rate_limited' }
  }
  if (isFutureTime(account.overload_until)) {
    return { code: 'overloaded' }
  }
  if (account.status === 'error' && isAccountAuthFailure(account.error_message)) {
    return { code: 'auth_failed' }
  }
  if (account.status === 'error') {
    return { code: 'error_state' }
  }
  if (isAccountQuotaExceeded(account)) {
    return { code: 'quota_exhausted' }
  }
  if (account.status === 'active' && !account.schedulable) {
    return { code: 'manual_pause' }
  }
  if (account.status === 'inactive') {
    return { code: 'inactive' }
  }
  return { code: 'healthy' }
}
