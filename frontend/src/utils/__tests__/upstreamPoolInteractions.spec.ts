import { describe, expect, it } from 'vitest'
import {
  buildDeleteAccountSetConfirmMessage,
  buildDeleteBindingConfirmMessage,
  buildDeleteMemberSetConfirmMessage,
  buildDeletePoolConfirmMessage,
  filterAccountsForPoolCompletion,
  getLatestPoolHealthAlertStates,
  getPoolRoutingPaginationInfo,
} from '../upstreamPoolInteractions'

describe('upstream pool interaction helpers', () => {
  it('builds a delete-pool confirmation with loaded impact counts', () => {
    const message = buildDeletePoolConfirmMessage({
      pool: { id: 10, name: 'OpenAI 主池' },
      selectedPoolId: 10,
      loadedMembers: [
        { id: 1, editable: true },
        { id: 2, editable: false },
      ],
      memberSets: [
        { id: 1, pool_id: 10, set_id: 30, set_name: 'OAuth 集合' },
      ],
      bindings: [
        { id: 1, pool_id: 10, group_name: '默认分组' },
        { id: 2, pool_id: 11, group_name: '其它分组' },
      ],
    })

    expect(message).toContain('确定删除上游池「OpenAI 主池」吗？')
    expect(message).toContain('直接成员：1 项')
    expect(message).toContain('集合展开成员：1 项')
    expect(message).toContain('集合绑定：1 项')
    expect(message).toContain('分组绑定：1 项')
  })

  it('does not invent unloaded pool member-set counts', () => {
    const message = buildDeletePoolConfirmMessage({
      pool: { id: 10, name: '未选中池' },
      selectedPoolId: 20,
      loadedMembers: [],
      memberSets: [],
      bindings: [],
    })

    expect(message).toContain('直接成员：未加载当前池')
    expect(message).toContain('集合绑定：未加载当前池')
  })

  it('builds account-set delete impact without pretending global pool bindings are loaded', () => {
    const message = buildDeleteAccountSetConfirmMessage({
      accountSet: { id: 30, name: 'OAuth 集合', account_count: 8 },
      selectedAccountSetId: 31,
      loadedMembers: [],
      memberSets: [],
    })

    expect(message).toContain('集合成员：8 项')
    expect(message).toContain('其它池是否引用需以后端删除校验为准')
  })

  it('explains member-set and binding delete scope', () => {
    expect(buildDeleteMemberSetConfirmMessage({
      id: 1,
      pool_id: 10,
      set_id: 30,
      set_name: 'OAuth 集合',
    })).toContain('不会删除账号集合本身')

    expect(buildDeleteBindingConfirmMessage({
      id: 9,
      pool_id: 10,
      group_name: '默认分组',
    })).toContain('该分组不再通过这条绑定命中当前上游池')
  })

  it('computes routing observability pagination state', () => {
    expect(getPoolRoutingPaginationInfo({
      loading: false,
      loadingMore: false,
      total: 120,
      page: 1,
      pageSize: 60,
      logsLength: 60,
    })).toMatchObject({
      loaded: 60,
      total: 120,
      hasMore: true,
      nextPage: 2,
      disabled: false,
    })

    expect(getPoolRoutingPaginationInfo({
      loading: false,
      loadingMore: false,
      total: 60,
      page: 1,
      pageSize: 60,
      logsLength: 60,
    })).toMatchObject({
      hasMore: false,
      disabled: true,
    })
  })

  it('keeps only the latest state for each pool health alert', () => {
    const logs = [
      { id: 3, created_at: '2026-07-12T10:10:00Z', extra: { alert_key: 'probe:1', alert_status: 'resolved' } },
      { id: 2, created_at: '2026-07-12T10:05:00Z', extra: { alert_key: 'capacity:1', alert_status: 'firing' } },
      { id: 1, created_at: '2026-07-12T10:00:00Z', extra: { alert_key: 'probe:1', alert_status: 'firing' } },
    ]

    expect(getLatestPoolHealthAlertStates(logs).map(log => log.id)).toEqual([3, 2])
  })

  it('filters completion accounts by enabled pool binding groups', () => {
    const accounts = [
      { id: 1, platform: 'openai', type: 'oauth', group_ids: [101] },
      { id: 2, platform: 'openai', type: 'oauth', group_ids: [102] },
      { id: 3, platform: 'openai', type: 'apikey', group_ids: [101] },
      { id: 4, platform: 'anthropic', type: 'oauth', group_ids: [101] },
    ]

    expect(filterAccountsForPoolCompletion({
      pool: { id: 10, name: 'Codex Pro', platform: 'openai', account_type_strategy: 'all' },
      accounts,
      bindings: [
        { id: 1, pool_id: 10, group_id: 101, enabled: true, group_name: 'Codex Pro' },
        { id: 2, pool_id: 10, group_id: 102, enabled: false, group_name: 'Codex Plus' },
      ],
    }).map(account => account.id)).toEqual([1, 3])
  })

  it('keeps oauth-only completion from adding API key accounts', () => {
    expect(filterAccountsForPoolCompletion({
      pool: { id: 10, name: 'OAuth only', platform: 'openai', account_type_strategy: 'oauth_only' },
      accounts: [
        { id: 1, platform: 'openai', type: 'oauth', group_ids: [] },
        { id: 2, platform: 'openai', type: 'setup-token', group_ids: [] },
        { id: 3, platform: 'openai', type: 'apikey', group_ids: [] },
      ],
      bindings: [],
    }).map(account => account.id)).toEqual([1, 2])
  })
})
