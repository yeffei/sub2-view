import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import type { DashboardStats } from '@/types'
import zh from '@/i18n/locales/zh'
import DashboardView from '../DashboardView.vue'

const { getSnapshotV2, getUserUsageTrend, getUserSpendingRanking, getDashboardOverview, listErrorLogs, showErrorMock, showSuccessMock, writeTextMock } = vi.hoisted(() => ({
  getSnapshotV2: vi.fn(),
  getUserUsageTrend: vi.fn(),
  getUserSpendingRanking: vi.fn(),
  getDashboardOverview: vi.fn(),
  listErrorLogs: vi.fn(),
  showErrorMock: vi.fn(),
  showSuccessMock: vi.fn(),
  writeTextMock: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    dashboard: {
      getSnapshotV2,
      getUserUsageTrend,
      getUserSpendingRanking
    },
    ops: {
      getDashboardOverview,
      listErrorLogs
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: showErrorMock,
    showSuccess: showSuccessMock
  })
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  })
}))

const translate = (key: string, params?: Record<string, unknown>) => {
  const value = key.split('.').reduce<unknown>((acc, part) => {
    if (acc && typeof acc === 'object' && part in (acc as Record<string, unknown>)) {
      return (acc as Record<string, unknown>)[part]
    }
    return undefined
  }, zh as unknown as Record<string, unknown>)

  if (typeof value !== 'string') return key
  return value.replace(/\{(\w+)\}/g, (_, token) => String(params?.[token] ?? {}))
}

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, unknown>) => translate(key, params)
    })
  }
})

const formatLocalDate = (date: Date): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const createDashboardStats = (): DashboardStats => ({
  total_users: 0,
  today_new_users: 0,
  active_users: 0,
  hourly_active_users: 0,
  stats_updated_at: '',
  stats_stale: false,
  total_api_keys: 0,
  active_api_keys: 0,
  total_accounts: 0,
  normal_accounts: 0,
  error_accounts: 0,
  ratelimit_accounts: 0,
  overload_accounts: 0,
  total_requests: 0,
  total_input_tokens: 0,
  total_output_tokens: 0,
  total_cache_creation_tokens: 0,
  total_cache_read_tokens: 0,
  total_tokens: 0,
  total_cost: 0,
  total_actual_cost: 0,
  total_account_cost: 0,
  today_requests: 0,
  today_input_tokens: 0,
  today_output_tokens: 0,
  today_cache_creation_tokens: 0,
  today_cache_read_tokens: 0,
  today_tokens: 0,
  today_cost: 0,
  today_actual_cost: 0,
  today_account_cost: 0,
  average_duration_ms: 0,
  uptime: 0,
  rpm: 0,
  tpm: 0
})

const createWrapper = () =>
  mount(DashboardView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' },
        LoadingSpinner: true,
        Icon: true,
        DateRangePicker: true,
        Select: true,
        ModelDistributionChart: true,
        TokenUsageTrend: true,
        Line: true,
        RouterLink: {
          props: ['to'],
          template: '<a><slot /></a>'
        }
      }
    }
  })

describe('admin DashboardView', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => {})
    getSnapshotV2.mockReset()
    showErrorMock.mockReset()
    showSuccessMock.mockReset()
    writeTextMock.mockReset()
    writeTextMock.mockResolvedValue(undefined)
    Object.defineProperty(navigator, 'clipboard', {
      configurable: true,
      value: {
        writeText: writeTextMock,
      },
    })
    getUserUsageTrend.mockReset()
    getUserSpendingRanking.mockReset()
    getDashboardOverview.mockReset()
    listErrorLogs.mockReset()

    getSnapshotV2.mockResolvedValue({
      stats: createDashboardStats(),
      trend: [],
      models: []
    })
    getUserUsageTrend.mockResolvedValue({
      trend: [],
      start_date: '',
      end_date: '',
      granularity: 'hour'
    })
    getUserSpendingRanking.mockResolvedValue({
      ranking: [],
      total_actual_cost: 0,
      total_requests: 0,
      total_tokens: 0,
      start_date: '',
      end_date: ''
    })
    getDashboardOverview.mockResolvedValue({
      error_rate: 0.18,
      start_time: '',
      end_time: '',
      platform: '',
      success_count: 0,
      error_count_total: 0,
      business_limited_count: 0,
      error_count_sla: 0,
      request_count_total: 0,
      request_count_sla: 0,
      token_consumed: 0,
      sla: 0,
      upstream_error_rate: 0,
      upstream_error_count_excl_429_529: 0,
      upstream_429_count: 0,
      upstream_529_count: 0,
      qps: { current: 0, peak: 0, avg: 0 },
      tps: { current: 0, peak: 0, avg: 0 },
      duration: {},
      ttft: {}
    })
    listErrorLogs.mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 50,
      total_pages: 1
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('uses last 24 hours as default dashboard range', async () => {
    createWrapper()

    await flushPromises()

    const now = new Date()
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)

    expect(getSnapshotV2).toHaveBeenCalledTimes(1)
    expect(getSnapshotV2).toHaveBeenCalledWith(expect.objectContaining({
      start_date: formatLocalDate(yesterday),
      end_date: formatLocalDate(now),
      granularity: 'hour'
    }))
  })

  it('renders branded briefing summary panels', async () => {
    const wrapper = createWrapper()

    await flushPromises()

    expect(wrapper.text()).toContain('品牌化晨报')
    expect(wrapper.text()).toContain('今晨案头快报')
    expect(wrapper.text()).toContain('流量概览')
    expect(wrapper.text()).toContain('异常焦点')
    expect(wrapper.text()).toContain('案头简报')
    expect(wrapper.text()).toContain('一页式值守简报')
    expect(wrapper.text()).toContain('重点观察')
  })
  it('copies the one-page morning sheet summary', async () => {
    getUserSpendingRanking.mockResolvedValue({
      ranking: [
        { user_id: 7, email: 'risk@example.com', actual_cost: 12, requests: 12, tokens: 3200 }
      ],
      total_actual_cost: 12,
      total_requests: 12,
      total_tokens: 3200,
      start_date: '',
      end_date: ''
    })
    listErrorLogs.mockResolvedValue({
      items: [
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'client', error_source: 'client_request', platform: 'openai', status_code: 429 },
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'client', error_source: 'client_request', platform: 'openai', status_code: 429 },
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'provider', error_source: 'upstream_http', platform: 'openai', status_code: 500 }
      ],
      total: 3,
      page: 1,
      page_size: 50,
      total_pages: 1
    })

    const wrapper = createWrapper()

    await flushPromises()

    const copyButton = wrapper.findAll('button').find((button) => button.text().includes('复制简报'))
    expect(copyButton).toBeTruthy()

    await copyButton!.trigger('click')

    expect(writeTextMock).toHaveBeenCalledTimes(1)
    const copied = writeTextMock.mock.calls[0][0] as string
    expect(copied).toContain('山枢庭 / 管理后台一页式值守简报')
    expect(copied).toContain('窗口摘要')
    expect(copied).toContain('异常落点')
    expect(copied).toContain('重点观察')
    expect(copied).toContain('异常归因小结')
    expect(copied).toContain('主要归属：用户侧请求')
    expect(copied).toContain('集中平台：openai')
    expect(copied).toContain('高风险用户明细')
    expect(copied).toContain('### 1. risk@example.com')
    expect(copied).toContain('邮箱: risk@example.com')
    expect(copied).toContain('风险标签：失败率走高')
    expect(showSuccessMock).toHaveBeenCalledWith('简报已复制')
  })
  it('downloads the morning sheet as markdown', async () => {
    const originalCreateElement = document.createElement.bind(document)
    const clickMock = vi.fn()
    const anchor = originalCreateElement('a')
    anchor.click = clickMock as unknown as typeof anchor.click
    const createElementMock = vi.spyOn(document, 'createElement').mockImplementation(((tagName: string) => {
      if (tagName === 'a') return anchor
      return originalCreateElement(tagName)
    }) as typeof document.createElement)
    const originalCreateObjectURL = URL.createObjectURL
    const originalRevokeObjectURL = URL.revokeObjectURL
    URL.createObjectURL = vi.fn(() => 'blob:mock') as typeof URL.createObjectURL
    URL.revokeObjectURL = vi.fn() as typeof URL.revokeObjectURL

    const wrapper = createWrapper()
    await flushPromises()

    const downloadButton = wrapper.findAll('button').find((button) => button.text().includes('导出 Markdown'))
    expect(downloadButton).toBeTruthy()

    await downloadButton!.trigger('click')

    expect(URL.createObjectURL).toHaveBeenCalledTimes(1)
    const blob = (URL.createObjectURL as any).mock.calls[0][0] as Blob
    expect(blob.type).toBe('text/markdown;charset=utf-8')
    expect(anchor.download.endsWith('.md')).toBe(true)
    expect(clickMock).toHaveBeenCalledTimes(1)
    expect(showSuccessMock).toHaveBeenCalledWith('Markdown 简报已导出')
    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:mock')

    createElementMock.mockRestore()
    URL.createObjectURL = originalCreateObjectURL
    URL.revokeObjectURL = originalRevokeObjectURL
  })
  it('uses RFC3339 time range for ops dashboard requests', async () => {
    createWrapper()

    await flushPromises()

    const now = new Date()
    const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
    const expectedStart = new Date(formatLocalDate(yesterday) + 'T00:00:00').toISOString()
    const expectedEnd = new Date(formatLocalDate(now) + 'T23:59:59.999').toISOString()

    expect(getDashboardOverview).toHaveBeenCalledWith(expect.objectContaining({
      start_time: expectedStart,
      end_time: expectedEnd
    }))
    expect(listErrorLogs).toHaveBeenCalledWith(expect.objectContaining({
      start_time: expectedStart,
      end_time: expectedEnd
    }))
  })
  it('renders fallback risk message when ranking data fails', async () => {
    getUserSpendingRanking.mockRejectedValue(new Error('ranking failed'))

    const wrapper = createWrapper()

    await flushPromises()

    expect(wrapper.text()).toContain('消费榜数据暂未取回')
    expect(wrapper.text()).toContain('回看榜单数据')
  })
  it('renders high risk users from ranking and error logs', async () => {
    getUserSpendingRanking.mockResolvedValue({
      ranking: [
        { user_id: 7, email: 'risk@example.com', actual_cost: 12, requests: 12, tokens: 3200 },
        { user_id: 8, email: 'steady@example.com', actual_cost: 6, requests: 20, tokens: 2800 }
      ],
      total_actual_cost: 18,
      total_requests: 32,
      total_tokens: 6000,
      start_date: '',
      end_date: ''
    })
    listErrorLogs.mockResolvedValue({
      items: [
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'client', error_source: 'client_request', platform: 'openai', status_code: 429 },
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'client', error_source: 'client_request', platform: 'openai', status_code: 429 },
        { user_id: 7, user_email: 'risk@example.com', error_owner: 'provider', error_source: 'upstream_http', platform: 'openai', status_code: 500 }
      ],
      total: 3,
      page: 1,
      page_size: 50,
      total_pages: 1
    })

    const wrapper = createWrapper()

    await flushPromises()

    expect(wrapper.text()).toContain('高风险用户榜单')
    expect(wrapper.text()).toContain('risk@example.com')
    expect(wrapper.text()).toContain('失败率走高')
  })
})






