import { beforeEach, describe, expect, it, vi } from 'vitest'

const { get, post, put, del } = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  del: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    get,
    post,
    put,
    delete: del,
  },
}))

import upstreamPoolsAPI from '@/api/admin/upstreamPools'

describe('admin upstream pools api', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    put.mockReset()
    del.mockReset()
    get.mockResolvedValue({ data: { items: [] } })
    post.mockResolvedValue({ data: {} })
    put.mockResolvedValue({ data: {} })
    del.mockResolvedValue({ data: { message: 'ok' } })
  })

  it('requests member sync preview before destructive apply', async () => {
    await upstreamPoolsAPI.previewMemberSync(12, { mode: 'membership_only' })

    expect(post).toHaveBeenCalledWith(
      '/admin/upstream-pools/12/member-sync/preview',
      { mode: 'membership_only' },
    )
  })

  it('applies member sync with the selected mode', async () => {
    await upstreamPoolsAPI.applyMemberSync(12, { mode: 'overwrite_scheduler_fields' })

    expect(post).toHaveBeenCalledWith(
      '/admin/upstream-pools/12/member-sync/apply',
      { mode: 'overwrite_scheduler_fields' },
    )
  })

  it('keeps binding priority 0 in update payloads', async () => {
    await upstreamPoolsAPI.updateBinding(5, {
      group_id: 3,
      pool_id: 12,
      platform: 'openai',
      models: [],
      request_path_scope: [],
      priority: 0,
      enabled: true,
    })

    expect(put).toHaveBeenCalledWith('/admin/upstream-pools/bindings/5', {
      group_id: 3,
      pool_id: 12,
      platform: 'openai',
      models: [],
      request_path_scope: [],
      priority: 0,
      enabled: true,
    })
  })
})
