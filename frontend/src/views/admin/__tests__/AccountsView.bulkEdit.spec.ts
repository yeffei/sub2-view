import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import AccountsView from '../AccountsView.vue'

const {
  listAccounts,
  listWithEtag,
  getBatchTodayStats,
  getAllProxies,
  getAllGroups,
  syncAllUpstreamRateMultipliers,
  batchSyncUpstreamRateMultiplier,
  syncUpstreamRateMultiplier,
  listSystemLogs,
  getSettings,
  updateSettings,
  showError,
  showSuccess,
  showWarning
} = vi.hoisted(() => ({
  listAccounts: vi.fn(),
  listWithEtag: vi.fn(),
  getBatchTodayStats: vi.fn(),
  getAllProxies: vi.fn(),
  getAllGroups: vi.fn(),
  syncAllUpstreamRateMultipliers: vi.fn(),
  batchSyncUpstreamRateMultiplier: vi.fn(),
  syncUpstreamRateMultiplier: vi.fn(),
  listSystemLogs: vi.fn(),
  getSettings: vi.fn(),
  updateSettings: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn(),
  showWarning: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      list: listAccounts,
      listWithEtag,
      getBatchTodayStats,
      delete: vi.fn(),
      batchClearError: vi.fn(),
      batchRefresh: vi.fn(),
      toggleSchedulable: vi.fn(),
      syncAllUpstreamRateMultipliers,
      batchSyncUpstreamRateMultiplier,
      syncUpstreamRateMultiplier
    },
    proxies: {
      getAll: getAllProxies
    },
    groups: {
      getAll: getAllGroups
    },
    ops: {
      listSystemLogs
    },
    settings: {
      getSettings,
      updateSettings
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
	showError,
	showSuccess,
	showWarning,
    showInfo: vi.fn()
  })
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    token: 'test-token'
  })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

const DataTableStub = {
  props: ['columns', 'data'],
  template: `
    <div data-test="data-table">
      <span v-for="column in columns" :key="column.key" data-test="column-key">{{ column.key }}</span>
      <div v-for="row in data" :key="row.id">
        <slot name="cell-created_at" :value="row.created_at" :row="row" />
        <slot name="cell-rate_multiplier" :value="row.rate_multiplier" :row="row" />
        <slot name="cell-actions" :row="row" />
      </div>
    </div>
  `
}

const AccountBulkActionsBarStub = {
  props: ['selectedIds'],
  emits: ['edit-filtered'],
  template: '<button data-test="edit-filtered" @click="$emit(\'edit-filtered\')">edit filtered</button>'
}

const BulkEditAccountModalStub = {
  props: ['show', 'target'],
  template: '<div data-test="bulk-edit-modal" :data-show="String(show)" :data-target-mode="target?.mode ?? \'\'"></div>'
}

const ConfirmDialogStub = {
  props: ['show'],
  emits: ['confirm', 'cancel'],
  template: '<button v-if="show" data-test="confirm-dialog-confirm" @click="$emit(\'confirm\')">confirm</button>'
}

const mountAccountsView = () => mount(AccountsView, {
  global: {
    stubs: {
      AppLayout: { template: '<div><slot /></div>' },
      TablePageLayout: {
        template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
      },
      DataTable: DataTableStub,
      Pagination: true,
      ConfirmDialog: ConfirmDialogStub,
      AccountTableActions: { template: '<div><slot name="beforeCreate" /><slot name="after" /></div>' },
      AccountTableFilters: { template: '<div></div>' },
      AccountBulkActionsBar: AccountBulkActionsBarStub,
      AccountActionMenu: true,
      ImportDataModal: true,
      ReAuthAccountModal: true,
      AccountTestModal: true,
      AccountStatsModal: true,
      ScheduledTestsPanel: true,
      SyncFromCrsModal: true,
      TempUnschedStatusModal: true,
      ErrorPassthroughRulesModal: true,
      TLSFingerprintProfilesModal: true,
      CreateAccountModal: true,
      EditAccountModal: true,
      BulkEditAccountModal: BulkEditAccountModalStub,
      PlatformTypeBadge: true,
      AccountCapacityCell: true,
      AccountStatusIndicator: true,
      AccountTodayStatsCell: true,
      AccountGroupsCell: true,
      AccountUsageCell: true,
      Icon: true
    }
  }
})

describe('admin AccountsView bulk edit scope', () => {
  beforeEach(() => {
    localStorage.clear()

    listAccounts.mockReset()
    listWithEtag.mockReset()
    getBatchTodayStats.mockReset()
    getAllProxies.mockReset()
    getAllGroups.mockReset()
		syncAllUpstreamRateMultipliers.mockReset()
		batchSyncUpstreamRateMultiplier.mockReset()
		syncUpstreamRateMultiplier.mockReset()
		listSystemLogs.mockReset()
		getSettings.mockReset()
		updateSettings.mockReset()
		showError.mockReset()
		showSuccess.mockReset()
		showWarning.mockReset()

    listAccounts.mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 20,
      pages: 0
    })
    listWithEtag.mockResolvedValue({
      notModified: true,
      etag: null,
      data: null
    })
    getBatchTodayStats.mockResolvedValue({ stats: {} })
    getAllProxies.mockResolvedValue([])
    getAllGroups.mockResolvedValue([])
		syncAllUpstreamRateMultipliers.mockResolvedValue({ total: 0, success: 0, failed: 0, results: [] })
		batchSyncUpstreamRateMultiplier.mockResolvedValue({ total: 0, success: 0, failed: 0, results: [] })
		syncUpstreamRateMultiplier.mockResolvedValue({})
		listSystemLogs.mockResolvedValue({ items: [], total: 0, page: 1, page_size: 50, pages: 0 })
		getSettings.mockResolvedValue({ upstream_rate_sync_enabled: false })
		updateSettings.mockImplementation(async payload => payload)
  })

  it('opens bulk edit in filtered-results mode from the bulk actions dropdown', async () => {
    const wrapper = mount(AccountsView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          TablePageLayout: {
            template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
          },
          DataTable: DataTableStub,
          Pagination: true,
          ConfirmDialog: true,
          AccountTableActions: { template: '<div><slot name="beforeCreate" /><slot name="after" /></div>' },
          AccountTableFilters: { template: '<div></div>' },
          AccountBulkActionsBar: AccountBulkActionsBarStub,
          AccountActionMenu: true,
          ImportDataModal: true,
          ReAuthAccountModal: true,
          AccountTestModal: true,
          AccountStatsModal: true,
          ScheduledTestsPanel: true,
          SyncFromCrsModal: true,
          TempUnschedStatusModal: true,
          ErrorPassthroughRulesModal: true,
          TLSFingerprintProfilesModal: true,
          CreateAccountModal: true,
          EditAccountModal: true,
          BulkEditAccountModal: BulkEditAccountModalStub,
          PlatformTypeBadge: true,
          AccountCapacityCell: true,
          AccountStatusIndicator: true,
          AccountTodayStatsCell: true,
          AccountGroupsCell: true,
          AccountUsageCell: true,
          Icon: true
        }
      }
    })

    await flushPromises()
    await wrapper.get('[data-test="edit-filtered"]').trigger('click')
    await flushPromises()

    expect(wrapper.get('[data-test="bulk-edit-modal"]').attributes('data-show')).toBe('true')
    expect(wrapper.get('[data-test="bulk-edit-modal"]').attributes('data-target-mode')).toBe('filtered')
  })

  it('renders the created_at column by default', async () => {
    listAccounts.mockResolvedValue({
      items: [
        {
          id: 1,
          name: 'test-account',
          platform: 'anthropic',
          type: 'oauth',
          status: 'active',
          schedulable: true,
          created_at: '2026-03-07T10:00:00Z',
          updated_at: '2026-03-07T10:00:00Z'
        }
      ],
      total: 1,
      page: 1,
      page_size: 20,
      pages: 1
    })

    const wrapper = mount(AccountsView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' },
          TablePageLayout: {
            template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
          },
          DataTable: DataTableStub,
          Pagination: true,
          ConfirmDialog: true,
          AccountTableActions: { template: '<div><slot name="beforeCreate" /><slot name="after" /></div>' },
          AccountTableFilters: { template: '<div></div>' },
          AccountBulkActionsBar: AccountBulkActionsBarStub,
          AccountActionMenu: true,
          ImportDataModal: true,
          ReAuthAccountModal: true,
          AccountTestModal: true,
          AccountStatsModal: true,
          ScheduledTestsPanel: true,
          SyncFromCrsModal: true,
          TempUnschedStatusModal: true,
          ErrorPassthroughRulesModal: true,
          TLSFingerprintProfilesModal: true,
          CreateAccountModal: true,
          EditAccountModal: true,
          BulkEditAccountModal: BulkEditAccountModalStub,
          PlatformTypeBadge: true,
          AccountCapacityCell: true,
          AccountStatusIndicator: true,
          AccountTodayStatsCell: true,
          AccountGroupsCell: true,
          AccountUsageCell: true,
          Icon: true
        }
      }
    })

    await flushPromises()

    expect(listAccounts.mock.calls[0]?.[1]).toBe(20)

    const columnKeys = wrapper.findAll('[data-test="column-key"]').map(node => node.text())
    expect(columnKeys).toContain('created_at')
    const columns = wrapper.getComponent(DataTableStub).props('columns') as Array<{ key: string; label: string; sortable: boolean }>
    expect(columns.find(column => column.key === 'created_at')).toMatchObject({
      label: 'admin.accounts.columns.createdAt',
      sortable: true
    })
  })

	it('confirms sync-all and retries only failed account IDs', async () => {
		syncAllUpstreamRateMultipliers.mockResolvedValue({
			total: 2,
			success: 1,
			failed: 1,
			results: [
				{ account_id: 7, success: false, error: 'upstream unavailable' },
				{
					account_id: 8,
					account_name: 'relay-8',
					success: true,
					previous_rate_multiplier: 0.04,
					rate_multiplier: 0.08,
					changed: true,
					significant_change: true
				}
			]
		})
		batchSyncUpstreamRateMultiplier.mockResolvedValue({
			total: 1,
			success: 1,
			failed: 0,
			results: [{ account_id: 7, success: true, rate_multiplier: 0.08 }]
		})
		const wrapper = mountAccountsView()
		await flushPromises()

		const syncAllButton = wrapper.findAll('button').find(button => button.text().includes('admin.accounts.syncAllUpstreamRates'))
		expect(syncAllButton).toBeTruthy()
		await syncAllButton!.trigger('click')
		expect(syncAllUpstreamRateMultipliers).not.toHaveBeenCalled()
		await wrapper.get('[data-test="confirm-dialog-confirm"]').trigger('click')
		await flushPromises()
		expect(syncAllUpstreamRateMultipliers).toHaveBeenCalledTimes(1)
		expect(showWarning).toHaveBeenCalledWith('admin.accounts.upstreamRateSyncSignificantWarning')
		expect(document.body.textContent).toContain('admin.accounts.upstreamRateSyncResultTitle')
		expect(document.body.textContent).toContain('0.0400x → 0.0800x')

		const retryButton = wrapper.findAll('button').find(button => button.text().includes('admin.accounts.retryFailedUpstreamRates'))
		expect(retryButton).toBeTruthy()
		await retryButton!.trigger('click')
		await flushPromises()
		expect(batchSyncUpstreamRateMultiplier).toHaveBeenCalledWith([7])
		expect(wrapper.findAll('button').some(button => button.text().includes('admin.accounts.retryFailedUpstreamRates'))).toBe(false)
	})

  it('syncs multiple account rates concurrently and keeps the current account list visible', async () => {
    const timeoutSpy = vi.spyOn(globalThis, 'setTimeout')
    localStorage.setItem('account-hidden-columns', JSON.stringify([]))
    const accounts = [
      {
        id: 21,
        name: 'relay-21',
        platform: 'openai',
        type: 'apikey',
        status: 'active',
        schedulable: true,
        rate_multiplier: 0.04,
        created_at: '2026-07-12T00:00:00Z',
        updated_at: '2026-07-12T00:00:00Z'
      },
      {
        id: 22,
        name: 'relay-22',
        platform: 'openai',
        type: 'apikey',
        status: 'active',
        schedulable: true,
        rate_multiplier: 0.05,
        created_at: '2026-07-12T00:00:00Z',
        updated_at: '2026-07-12T00:00:00Z'
      }
    ]
    listAccounts.mockResolvedValue({ items: accounts, total: 2, page: 1, page_size: 20, pages: 1 })

    let resolveFirst!: (value: any) => void
    let resolveSecond!: (value: any) => void
    syncUpstreamRateMultiplier.mockImplementation((accountId: number) => new Promise(resolve => {
      if (accountId === 21) resolveFirst = resolve
      else resolveSecond = resolve
    }))

    const wrapper = mountAccountsView()
    await flushPromises()

    const syncButtons = wrapper.findAll('button[title="admin.accounts.syncUpstreamRateMultiplier"]')
    expect(syncButtons).toHaveLength(2)
    await syncButtons[0].trigger('click')
    await syncButtons[1].trigger('click')

    expect(syncUpstreamRateMultiplier).toHaveBeenCalledTimes(2)
    expect(syncUpstreamRateMultiplier).toHaveBeenNthCalledWith(1, 21)
    expect(syncUpstreamRateMultiplier).toHaveBeenNthCalledWith(2, 22)
    expect(syncButtons[0].attributes('disabled')).toBeDefined()
    expect(syncButtons[1].attributes('disabled')).toBeDefined()

    resolveSecond({
      previous_rate_multiplier: 0.05,
      rate_multiplier: 0.08,
      changed: true,
      significant_change: true,
      source: 'user_group_override',
      account: { ...accounts[1], rate_multiplier: 0.08 }
    })
    await flushPromises()

    const rowsAfterSecondSync = wrapper.getComponent(DataTableStub).props('data') as Array<{ id: number; rate_multiplier: number }>
    expect(rowsAfterSecondSync).toHaveLength(2)
    expect(rowsAfterSecondSync.find(account => account.id === 22)?.rate_multiplier).toBe(0.08)
    expect(wrapper.text()).toContain('+0.03x')
    const resultDialog = wrapper.findAllComponents({ name: 'BaseDialog' })
      .find(component => component.props('title') === 'admin.accounts.upstreamRateSyncResultTitle')
    expect(resultDialog?.props('show')).toBe(false)

    resolveFirst({
      previous_rate_multiplier: 0.04,
      rate_multiplier: 0.06,
      changed: true,
      significant_change: false,
      source: 'user_group_override',
      account: { ...accounts[0], rate_multiplier: 0.06 }
    })
    await flushPromises()

    const finalRows = wrapper.getComponent(DataTableStub).props('data') as Array<{ id: number; rate_multiplier: number }>
    expect(finalRows).toHaveLength(2)
    expect(finalRows.find(account => account.id === 21)?.rate_multiplier).toBe(0.06)
    expect(finalRows.find(account => account.id === 22)?.rate_multiplier).toBe(0.08)
    expect(wrapper.text()).toContain('+0.02x')
    expect(timeoutSpy).toHaveBeenCalledWith(expect.any(Function), 10_000)
    wrapper.unmount()
    timeoutSpy.mockRestore()
  })

  it('places upstream rate sync before usage stats and hides the row schedule action', async () => {
    localStorage.setItem('account-hidden-columns', JSON.stringify([]))
    listAccounts.mockResolvedValue({
      items: [{
        id: 31,
        name: 'relay-31',
        platform: 'openai',
        type: 'apikey',
        status: 'active',
        schedulable: true,
        rate_multiplier: 0.03,
        created_at: '2026-07-12T00:00:00Z',
        updated_at: '2026-07-12T00:00:00Z'
      }],
      total: 1,
      page: 1,
      page_size: 20,
      pages: 1
    })

    const wrapper = mountAccountsView()
    await flushPromises()

    const actionTitles = wrapper.findAll('[data-test="data-table"] button')
      .map(button => button.attributes('title'))
      .filter((title): title is string => Boolean(title))
    expect(actionTitles).not.toContain('admin.scheduledTests.schedule')
    expect(actionTitles.indexOf('admin.accounts.syncUpstreamRateMultiplier'))
      .toBeLessThan(actionTitles.indexOf('admin.accounts.usageStatistics'))
    const columns = wrapper.getComponent(DataTableStub).props('columns') as Array<{ key: string; class?: string }>
    expect(columns.find(column => column.key === 'rate_multiplier')?.class).toContain('w-[7.75rem]')
    expect(columns.find(column => column.key === 'actions')?.class).toContain('w-[10.5rem]')
  })

  it('loads the selected API-key account cost rate history', async () => {
    localStorage.setItem('account-hidden-columns', JSON.stringify([]))
    listAccounts.mockResolvedValue({
      items: [
        {
          id: 12,
          name: 'relay-12',
          platform: 'openai',
          type: 'apikey',
          status: 'active',
          schedulable: true,
          rate_multiplier: 0.08,
          created_at: '2026-07-12T00:00:00Z',
          updated_at: '2026-07-12T00:00:00Z'
        }
      ],
      total: 1,
      page: 1,
      page_size: 20,
      pages: 1
    })
    listSystemLogs.mockResolvedValue({
      items: [
        {
          id: 99,
          created_at: '2026-07-12T01:00:00Z',
          level: 'warn',
          component: 'upstream.cost_rate_sync',
          account_id: 12,
          message: '账号上游成本倍率已同步',
          extra: {
            previous_rate_multiplier: 0.04,
            rate_multiplier: 0.08,
            change_ratio: 1,
            significant_change: true,
            rate_source: 'user_group_override'
          }
        }
      ],
      total: 1,
      page: 1,
      page_size: 50,
      pages: 1
    })

    const wrapper = mountAccountsView()
    await flushPromises()
    await wrapper.get('button[title="admin.accounts.viewCostRateHistory"]').trigger('click')
    await flushPromises()

    expect(listSystemLogs).toHaveBeenCalledWith({
      page: 1,
      page_size: 50,
      time_range: '30d',
      component: 'upstream.cost_rate_sync',
      account_id: 12
    })
    expect(document.body.textContent).toContain('0.0400x → 0.0800x')
    expect(document.body.textContent).toContain('+100.0%')
  })

  it('enables the fixed three-times-daily cost rate sync switch', async () => {
    const wrapper = mountAccountsView()
    await flushPromises()

    const toggle = wrapper.findAll('button').find(button => button.text().includes('admin.accounts.scheduledCostRateSync'))
    expect(toggle).toBeTruthy()
    expect(toggle!.attributes('aria-pressed')).toBe('false')
    await toggle!.trigger('click')
    await flushPromises()

    expect(updateSettings).toHaveBeenCalledWith({ upstream_rate_sync_enabled: true })
    expect(toggle!.attributes('aria-pressed')).toBe('true')
    expect(showSuccess).toHaveBeenCalledWith('admin.accounts.scheduledCostRateSyncEnabled')
  })
})
