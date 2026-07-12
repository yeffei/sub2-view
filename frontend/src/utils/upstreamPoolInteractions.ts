type NamedEntity = {
  id: number
  name?: string
  code?: string
}

type PoolLike = NamedEntity & {
  platform?: string
  account_type_strategy?: string
}

type AccountSetLike = NamedEntity & {
  account_count?: number
}

type PoolMemberLike = {
  id: number
  pool_id?: number
  editable?: boolean
}

type MemberSetLike = {
  id: number
  pool_id: number
  set_id: number
  set_name?: string
  enabled?: boolean
}

type BindingLike = {
  id: number
  group_id?: number
  pool_id: number
  group_name?: string
  enabled?: boolean
}

type SyncableAccountLike = {
  id: number
  platform: string
  type: string
  group_ids?: number[]
}

type AccountSetMemberLike = {
  account_id: number
}

export type PoolRoutingPaginationState = {
  loading: boolean
  loadingMore: boolean
  total: number
  page: number
  pageSize: number
  logsLength: number
}

export type PoolHealthAlertLogLike = {
  id: number
  created_at: string
  message?: string
  extra?: Record<string, any>
}

const formatEntityLabel = (entity: NamedEntity): string => {
  if (entity.name) return entity.name
  if (entity.code) return entity.code
  return `#${entity.id}`
}

const formatCountLine = (label: string, count: number, suffix = '项') => `- ${label}：${count} ${suffix}`

export function buildDeletePoolConfirmMessage(params: {
  pool: PoolLike
  selectedPoolId?: number | null
  loadedMembers: PoolMemberLike[]
  memberSets: MemberSetLike[]
  bindings: BindingLike[]
}): string {
  const { pool, selectedPoolId, loadedMembers, memberSets, bindings } = params
  const directMembersLoaded = selectedPoolId === pool.id
  const directMemberCount = directMembersLoaded
    ? loadedMembers.filter(member => member.editable !== false).length
    : null
  const expandedMemberCount = directMembersLoaded
    ? loadedMembers.filter(member => member.editable === false).length
    : null
  const memberSetCount = directMembersLoaded ? memberSets.filter(item => item.pool_id === pool.id).length : null
  const bindingCount = bindings.filter(item => item.pool_id === pool.id).length

  const lines = [
    `确定删除上游池「${formatEntityLabel(pool)}」吗？`,
    '',
    '影响范围：',
    directMemberCount == null
      ? '- 直接成员：未加载当前池，删除前无法在前端精确预览'
      : formatCountLine('直接成员', directMemberCount),
    expandedMemberCount == null
      ? '- 集合展开成员：未加载当前池，删除前无法在前端精确预览'
      : formatCountLine('集合展开成员', expandedMemberCount),
    memberSetCount == null
      ? '- 集合绑定：未加载当前池，删除前无法在前端精确预览'
      : formatCountLine('集合绑定', memberSetCount),
    formatCountLine('分组绑定', bindingCount),
    '',
    '删除后这些关联将不可再用于路由；如需保留配置，请先取消或迁移绑定。'
  ]

  return lines.join('\n')
}

export function buildDeleteAccountSetConfirmMessage(params: {
  accountSet: AccountSetLike
  selectedAccountSetId?: number | null
  loadedMembers: AccountSetMemberLike[]
  memberSets: MemberSetLike[]
}): string {
  const { accountSet, selectedAccountSetId, loadedMembers, memberSets } = params
  const membersLoaded = selectedAccountSetId === accountSet.id
  const memberCount = membersLoaded ? loadedMembers.length : accountSet.account_count
  const attachedPoolCount = memberSets.filter(item => item.set_id === accountSet.id).length

  return [
    `确定删除账号集合「${formatEntityLabel(accountSet)}」吗？`,
    '',
    '影响范围：',
    memberCount == null
      ? '- 集合成员：未加载，删除前无法在前端精确预览'
      : formatCountLine('集合成员', memberCount),
    attachedPoolCount > 0
      ? formatCountLine('当前已加载池内绑定', attachedPoolCount, '个')
      : '- 当前已加载池内绑定：0 个；其它池是否引用需以后端删除校验为准',
    '',
    '删除后，使用该集合的池会失去这批集合成员。'
  ].join('\n')
}

export function buildDeleteMemberSetConfirmMessage(item: MemberSetLike): string {
  return [
    `确定删除集合绑定「${item.set_name || item.set_id}」吗？`,
    '',
    '影响范围：',
    '- 只移除这个集合与当前上游池的关系',
    '- 不会删除账号集合本身，也不会删除集合内账号',
    '',
    '删除后，当前池不再从该集合展开成员。'
  ].join('\n')
}

export function buildDeleteBindingConfirmMessage(binding: BindingLike): string {
  return [
    `确定删除分组绑定 #${binding.id} 吗？`,
    '',
    '影响范围：',
    `- 分组：${binding.group_name || `#${binding.id}`}`,
    '- 该分组不再通过这条绑定命中当前上游池',
    '',
    '如果这是唯一绑定，相关请求会回退到其它可用路由策略。'
  ].join('\n')
}

export function getPoolRoutingPaginationInfo(state: PoolRoutingPaginationState) {
  const loaded = state.logsLength
  const total = state.total || loaded
  const hasMore = total > loaded
  const nextPage = state.page + 1
  const disabled = state.loading || state.loadingMore || !hasMore
  return {
    loaded,
    total,
    hasMore,
    nextPage,
    disabled,
    label: hasMore ? `加载更多（${loaded}/${total}）` : `已加载全部 ${loaded} 条`
  }
}

// 系统日志按时间倒序返回。每个 alert_key 只保留最新状态，避免已恢复的
// 旧 firing 记录继续显示为当前异常。
export function getLatestPoolHealthAlertStates<T extends PoolHealthAlertLogLike>(logs: T[]): T[] {
  const latest = new Map<string, T>()
  for (const log of logs) {
    const key = String(log.extra?.alert_key || '').trim()
    if (!key || latest.has(key)) continue
    latest.set(key, log)
  }
  return Array.from(latest.values())
}

export function getEnabledPoolBindingGroupIDs(poolID: number, bindings: BindingLike[]): Set<number> {
  const groupIDs = new Set<number>()
  for (const binding of bindings) {
    if (binding.pool_id !== poolID || binding.enabled === false || !binding.group_id) continue
    groupIDs.add(binding.group_id)
  }
  return groupIDs
}

export function accountMatchesPoolAccountTypeStrategy(account: Pick<SyncableAccountLike, 'type'>, strategy?: string) {
  const normalized = String(strategy || 'all').trim()
  if (normalized === 'oauth_only') {
    return account.type === 'oauth' || account.type === 'setup-token'
  }
  return true
}

export function filterAccountsForPoolCompletion(params: {
  pool: PoolLike | null | undefined
  accounts: SyncableAccountLike[]
  bindings: BindingLike[]
}): SyncableAccountLike[] {
  const { pool, accounts, bindings } = params
  if (!pool?.platform) return []

  const bindingGroupIDs = getEnabledPoolBindingGroupIDs(pool.id, bindings)
  return accounts.filter(account => {
    if (account.platform !== pool.platform) return false
    if (!accountMatchesPoolAccountTypeStrategy(account, pool.account_type_strategy)) return false
    if (bindingGroupIDs.size === 0) return true
    return (account.group_ids || []).some(groupID => bindingGroupIDs.has(groupID))
  })
}
