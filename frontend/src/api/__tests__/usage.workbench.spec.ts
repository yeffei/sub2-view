import { beforeEach, describe, expect, it, vi } from 'vitest'

const postMock = vi.fn()

vi.mock('@/api/client', () => ({
  apiClient: {
    post: postMock,
  },
}))

describe('usage api key workbench api', () => {
  beforeEach(() => {
    postMock.mockReset()
  })

  it('posts api key ids to the workbench summary endpoint', async () => {
    postMock.mockResolvedValue({
      data: {
        stats: {
          '12': {
            api_key_id: 12,
            success_requests_24h: 8,
            error_count_24h: 2,
          },
        },
      },
    })

    const { getDashboardApiKeysWorkbench } = await import('@/api/usage')
    const result = await getDashboardApiKeysWorkbench([12, 15])

    expect(postMock).toHaveBeenCalledWith(
      '/usage/dashboard/api-keys-workbench',
      { api_key_ids: [12, 15] },
      { signal: undefined },
    )
    expect(result.stats['12']?.api_key_id).toBe(12)
  })

  it('falls back to the legacy usage endpoint when workbench route is unavailable', async () => {
    postMock
      .mockRejectedValueOnce({ status: 404, message: 'page not found' })
      .mockResolvedValueOnce({
        data: {
          stats: {
            '12': {
              api_key_id: 12,
              today_actual_cost: 1.25,
              total_actual_cost: 9.5,
              success_requests_24h: 7,
            },
          },
        },
      })

    const { getDashboardApiKeysWorkbench } = await import('@/api/usage')
    const result = await getDashboardApiKeysWorkbench([12])

    expect(postMock).toHaveBeenNthCalledWith(
      1,
      '/usage/dashboard/api-keys-workbench',
      { api_key_ids: [12] },
      { signal: undefined },
    )
    expect(postMock).toHaveBeenNthCalledWith(
      2,
      '/usage/dashboard/api-keys-usage',
      { api_key_ids: [12] },
      { signal: undefined },
    )
    expect(result.stats['12']).toMatchObject({
      api_key_id: 12,
      today_actual_cost: 1.25,
      total_actual_cost: 9.5,
      success_requests_24h: 7,
      error_count_24h: 0,
      attempt_count_24h: 7,
    })
  })
})
