import { describe, expect, it } from 'vitest'
import { buildWorkbenchLatestErrorLabel, buildWorkbenchModelHints, buildWorkbenchModelInsight } from '@/utils/apiKeyWorkbench'

describe('apiKeyWorkbench helpers', () => {
  it('prefers reason label when latest error has a specific reason', () => {
    const label = buildWorkbenchLatestErrorLabel(
      {
        error_id: 1,
        created_at: '2026-07-03T01:00:00Z',
        category: 'service_unavailable',
        message: 'no available accounts',
        reason_code: 'service_model_not_available',
      },
      { service_model_not_available: '当前线路未开放该模型' },
      { service_unavailable: '线路暂不可用' },
    )

    expect(label).toBe('最近失败 当前线路未开放该模型')
  })

  it('builds route-level model availability hint from requested model', () => {
    expect(
      buildWorkbenchModelHints({
        error_id: 2,
        created_at: '2026-07-03T01:00:00Z',
        category: 'service_unavailable',
        message: 'no available accounts supporting model',
        reason_code: 'service_model_not_available',
        requested_model: 'gpt-5.5',
      }),
    ).toEqual(['模型 gpt-5.5 当前线路未开放'])
  })

  it('includes upstream mapping hint when request and upstream models differ', () => {
    expect(
      buildWorkbenchModelHints({
        error_id: 3,
        created_at: '2026-07-03T01:00:00Z',
        category: 'invalid_request',
        message: 'model not found',
        reason_code: 'request_model_not_supported',
        requested_model: 'gpt-5.5',
        upstream_model: 'gpt-5.5-mini',
      }),
    ).toEqual([
      '模型 gpt-5.5 当前分组未开放',
      '请求 gpt-5.5 -> 实际上游 gpt-5.5-mini',
    ])
  })

  it('returns no hints when the latest error has no model evidence', () => {
    expect(
      buildWorkbenchModelHints({
        error_id: 4,
        created_at: '2026-07-03T01:00:00Z',
        category: 'upstream',
        message: 'temporary upstream issue',
      }),
    ).toEqual([])
  })

  it('builds a dashboard-style insight for route model unavailability', () => {
    expect(
      buildWorkbenchModelInsight({
        error_id: 5,
        created_at: '2026-07-03T01:00:00Z',
        category: 'service_unavailable',
        message: 'no available accounts supporting model',
        reason_code: 'service_model_not_available',
        requested_model: 'gpt-5.5',
      }, 'prod-key'),
    ).toEqual({
      status: '待处理',
      detail: '密钥 prod-key 最近请求的模型 gpt-5.5 当前线路未开放，更像路由覆盖不到该模型而不是瞬时波动。',
      action: '去做接入体检',
      to: '/keys?panel=connection-test',
      tone: 'alert',
    })
  })

  it('falls back to mapping insight when requested and upstream models differ', () => {
    expect(
      buildWorkbenchModelInsight({
        error_id: 6,
        created_at: '2026-07-03T01:00:00Z',
        category: 'invalid_request',
        message: 'model not found',
        requested_model: 'gpt-5.5',
        upstream_model: 'gpt-5.5-mini',
      }),
    ).toEqual({
      status: '留意',
      detail: '这把 Key 最近一次请求以 gpt-5.5 进入，但实际上游模型是 gpt-5.5-mini，如果这不是你的预期，优先检查分组映射。',
      action: '查看最近失败',
      to: '/usage?tab=errors',
      tone: 'notice',
    })
  })
})
