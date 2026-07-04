import { beforeEach, describe, expect, it, vi } from 'vitest'

const { get, post } = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
}))

vi.mock('@/api/client', () => ({
  apiClient: {
    get,
    post,
  },
}))

import { systemAPI } from '@/api/admin/system'

describe('admin system api', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    get.mockResolvedValue({ data: {} })
    post.mockResolvedValue({ data: {} })
  })

  it('performs release update with idempotency key and long timeout', async () => {
    await systemAPI.performUpdate('system-update-test-key')

    expect(post).toHaveBeenCalledWith('/admin/system/update', undefined, {
      headers: { 'Idempotency-Key': 'system-update-test-key' },
      timeout: 10 * 60 * 1000,
    })
  })

  it('performs release update without idempotency header when key is omitted', async () => {
    await systemAPI.performUpdate()

    expect(post).toHaveBeenCalledWith('/admin/system/update', undefined, {
      headers: undefined,
      timeout: 10 * 60 * 1000,
    })
  })

  it('checks update preflight with forced refresh', async () => {
    await systemAPI.checkUpdatePreflight(true)

    expect(get).toHaveBeenCalledWith('/admin/system/update/preflight', {
      params: { force: 'true' },
    })
  })

  it('rolls back with idempotency key and guarded timeout', async () => {
    await systemAPI.rollback('system-rollback-test-key')

    expect(post).toHaveBeenCalledWith('/admin/system/rollback', undefined, {
      headers: { 'Idempotency-Key': 'system-rollback-test-key' },
      timeout: 5 * 60 * 1000,
    })
  })

  it('restarts with idempotency key and short timeout', async () => {
    await systemAPI.restartService('system-restart-test-key')

    expect(post).toHaveBeenCalledWith('/admin/system/restart', undefined, {
      headers: { 'Idempotency-Key': 'system-restart-test-key' },
      timeout: 60 * 1000,
    })
  })
})
