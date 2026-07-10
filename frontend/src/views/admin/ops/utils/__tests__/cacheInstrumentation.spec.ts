import { describe, expect, it } from 'vitest'
import type { OpsSystemLog } from '@/api/admin/ops'
import {
  buildCacheInstrumentationSummary,
  formatCacheInstrumentationDetail,
  isCacheInstrumentationLog
} from '../cacheInstrumentation'

function buildRow(extra: Record<string, any>, overrides: Partial<OpsSystemLog> = {}): OpsSystemLog {
  return {
    id: overrides.id ?? 1,
    created_at: overrides.created_at ?? '2026-07-09T10:00:00Z',
    level: overrides.level ?? 'info',
    component: overrides.component ?? 'cache.instrumentation',
    message: overrides.message ?? 'cache_instrumentation',
    request_id: overrides.request_id ?? 'local:req-1',
    platform: overrides.platform ?? 'openai',
    model: overrides.model ?? 'gpt-5',
    extra,
    ...overrides
  }
}

describe('cacheInstrumentation utils', () => {
  it('detects cache instrumentation rows', () => {
    expect(isCacheInstrumentationLog(buildRow({}))).toBe(true)
    expect(isCacheInstrumentationLog(buildRow({}, { component: 'routing.explanation', message: 'routing_explanation' }))).toBe(false)
  })

  it('formats cache instrumentation detail with key diagnostics', () => {
    const detail = formatCacheInstrumentationDetail(buildRow({
      cache_family: 'anthropic_to_openai_responses',
      session_hash_short: 'abcd1234',
      session_signal_source: 'anthropic_digest',
      prompt_cache_key_present: false,
      prompt_cache_key_auto_injected: true,
      previous_response_id_present: true,
      sticky_account_hit: false,
      sticky_account_id: 11,
      selected_account_id: 22,
      account_switch_happened: true,
      cache_read_tokens: 0,
      cache_creation_tokens: 1280,
      request_type: 'stream',
      requested_model: 'claude-sonnet-4-5',
      upstream_model: 'gpt-5',
      inbound_endpoint: '/v1/messages',
      upstream_endpoint: '/backend-api/codex/responses',
      sample_reason: 'sticky_miss_or_switch',
      sample_rate: 1
    }))

    expect(detail).toContain('family=anthropic_to_openai_responses')
    expect(detail).toContain('source=anthropic_digest')
    expect(detail).toContain('prompt_key=no(auto)')
    expect(detail).toContain('sticky=miss(11)')
    expect(detail).toContain('switch=yes')
    expect(detail).toContain('关注=账号切换 / 粘性未命中 / 缺少 prompt_cache_key / 使用自动派生 key')
  })

  it('builds page-level summary for cache diagnostics', () => {
    const rows = [
      buildRow({
        cache_family: 'openai_responses',
        session_signal_source: 'header_session_id',
        prompt_cache_key_present: true,
        sticky_account_hit: true,
        cache_read_tokens: 40,
        cache_creation_tokens: 0
      }),
      buildRow({
        cache_family: 'anthropic_to_openai_responses',
        session_signal_source: 'anthropic_digest',
        prompt_cache_key_present: false,
        prompt_cache_key_auto_injected: true,
        previous_response_id_present: true,
        sticky_account_hit: false,
        sticky_account_id: 11,
        account_switch_happened: true,
        cache_read_tokens: 0,
        cache_creation_tokens: 200
      }, { id: 2 }),
      buildRow({}, { id: 3, component: 'routing.explanation', message: 'routing_explanation' })
    ]

    const summary = buildCacheInstrumentationSummary(rows)

    expect(summary.total).toBe(2)
    expect(summary.cacheHitCount).toBe(1)
    expect(summary.accountSwitchCount).toBe(1)
    expect(summary.stickyMissCount).toBe(1)
    expect(summary.previousResponseCount).toBe(1)
    expect(summary.promptCacheKeyMissingCount).toBe(1)
    expect(summary.autoInjectedCount).toBe(1)
    expect(summary.families[0]).toEqual({ label: 'anthropic_to_openai_responses', count: 1 })
    expect(summary.signalSources.map((item) => item.label)).toContain('header_session_id')
  })
})
