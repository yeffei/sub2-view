<template>
  <AppLayout>
    <div class="sst-admin-page upstream-pools-night">

      <TablePageLayout>
        <template #filters>
          <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
            <div class="flex flex-1 flex-wrap items-center gap-3">
              <div class="relative w-full sm:w-72">
                <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <input
                  v-model="searchQuery"
                  type="text"
                  class="input pl-10"
                  placeholder="搜索池名、编码、描述"
                  @input="handleSearch"
                />
              </div>
              <Select
                v-model="filters.platform"
                class="w-40"
                :options="platformOptions"
                placeholder="全部平台"
                @change="handleFilterChange"
              />
              <Select
                v-model="filters.enabled"
                class="w-32"
                :options="enabledOptions"
                placeholder="全部状态"
                @change="handleFilterChange"
              />
            </div>
            <div class="flex flex-wrap items-center justify-end gap-2">
              <button class="btn btn-secondary" :disabled="loading" @click="loadAll">
                <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
              </button>
              <button class="btn btn-primary" @click="openPoolModal()">
                <Icon name="plus" size="md" class="mr-2" />
                新建上游池
              </button>
            </div>
          </div>
        </template>

        <template #table>
          <DataTable
            :columns="poolColumns"
            :data="pagedPools"
            :loading="loading"
            row-key="id"
          >
            <template #cell-name="{ row, value }">
              <div class="flex flex-col gap-1">
                <div class="font-medium text-gray-900 dark:text-white">{{ value }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">编码：{{ row.code }}</div>
              </div>
            </template>

            <template #cell-platform="{ value }">
              <span class="badge badge-gray">{{ platformLabel(value) }}</span>
            </template>

            <template #cell-enabled="{ row }">
              <button
                type="button"
                :class="['pool-status-toggle', row.enabled ? 'is-enabled' : 'is-disabled']"
                :disabled="isPoolStatusToggling(row.id)"
                :title="row.enabled ? `点击关闭上游池 ${row.name}` : `点击启用上游池 ${row.name}`"
                :aria-label="row.enabled ? `关闭上游池 ${row.name}` : `启用上游池 ${row.name}`"
                @click.stop="togglePoolEnabled(row)"
              >
                <span class="pool-status-toggle__dot" aria-hidden="true"></span>
                <span>{{ isPoolStatusToggling(row.id) ? '处理中' : (row.enabled ? '启用' : '关闭') }}</span>
              </button>
            </template>

            <template #cell-sticky_enabled="{ value }">
              <span :class="['badge', value ? 'badge-primary' : 'badge-gray']">
                {{ value ? '开启' : '关闭' }}
              </span>
            </template>

            <template #cell-routing="{ row }">
              <div class="flex flex-col gap-1 text-xs text-gray-600 dark:text-gray-300">
                <div class="flex flex-wrap gap-1">
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">负载均衡 {{ row.load_balance_enabled ? '开' : '关' }}</span>
                  <span v-if="row.auto_weight_mode && row.auto_weight_mode !== 'off'" class="rounded bg-emerald-50 px-2 py-0.5 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300">调权 {{ row.auto_weight_mode === 'observe' ? '观察' : '应用' }}</span>
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">故障转移 {{ row.failover_enabled ? '开' : '关' }}</span>
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">TopK {{ row.top_k }}</span>
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">{{ accountTypeStrategyLabel(row.account_type_strategy) }}</span>
                </div>
                <div class="flex flex-wrap gap-1">
                  <span class="rounded bg-primary-50 px-2 py-0.5 text-primary-700 dark:bg-primary-950/30 dark:text-primary-200">
                    逃逸 {{ row.sticky_escape_enabled ? '开' : '关' }}
                  </span>
                  <span class="rounded bg-primary-50 px-2 py-0.5 text-primary-700 dark:bg-primary-950/30 dark:text-primary-200">
                    TTFT {{ formatThresholdMs(row.sticky_escape_ttft_ms_threshold) }}
                  </span>
                  <span class="rounded bg-primary-50 px-2 py-0.5 text-primary-700 dark:bg-primary-950/30 dark:text-primary-200">
                    错率 {{ formatRateThreshold(row.sticky_escape_error_rate_threshold) }}
                  </span>
                </div>
              </div>
            </template>

            <template #cell-members="{ row }">
              <button
                class="text-sm text-primary-700 hover:underline dark:text-primary-300"
                :title="`启用 ${row.member_enabled_count ?? 0}，共 ${row.member_total_count ?? 0} 个成员`"
                @click="selectPool(row)"
              >
                {{ row.member_enabled_count ?? 0 }} / {{ row.member_total_count ?? 0 }}
              </button>
            </template>

            <template #cell-bindings="{ row }">
              <button
                class="text-sm text-primary-700 hover:underline dark:text-primary-300"
                :title="`启用 ${row.binding_enabled_count ?? 0}，共 ${row.binding_total_count ?? 0} 个绑定`"
                @click="selectPool(row)"
              >
                {{ row.binding_enabled_count ?? 0 }} / {{ row.binding_total_count ?? 0 }}
              </button>
            </template>

            <template #cell-actions="{ row }">
              <div class="flex flex-wrap items-center gap-1">
                <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="selectPool(row)">
                  查看
                </button>
                <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="openPoolModal(row)">
                  编辑
                </button>
                <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeletePool(row)">
                  删除
                </button>
              </div>
            </template>

            <template #empty>
              <EmptyState
                title="暂无上游池"
                description="先创建一个上游池，再添加成员和分组绑定。"
                action-text="新建上游池"
                @action="openPoolModal()"
              />
            </template>
          </DataTable>
        </template>

        <template #pagination>
          <Pagination
            v-if="pagination.total > 0"
            :page="pagination.page"
            :total="pagination.total"
            :page-size="pagination.page_size"
            @update:page="handlePageChange"
            @update:pageSize="handlePageSizeChange"
          />
        </template>
      </TablePageLayout>

      <section class="sst-admin-panel account-sets-panel mb-4 p-4">
        <div class="account-sets-header mb-3 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">账号集合</h3>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              先把同平台账号放进集合，再把集合绑定到上游池成员。
            </p>
          </div>
          <button class="btn btn-secondary btn-sm" @click="openAccountSetModal()">
            新建账号集合
          </button>
        </div>
        <div class="account-sets-guide mb-4">
          <span>新建同平台集合</span>
          <span>加入账号</span>
          <span>挂到上游池</span>
          <span>分组命中池</span>
        </div>
        <DataTable :columns="accountSetColumns" :data="accountSets" :loading="accountSetsLoading" row-key="id">
          <template #cell-name="{ row, value }">
            <div class="flex flex-col gap-1">
              <button class="text-left font-medium text-primary-700 hover:underline dark:text-primary-300" @click="selectAccountSet(row)">
                {{ value }}
              </button>
              <div class="text-xs text-gray-500 dark:text-gray-400">编码：{{ row.code }}</div>
            </div>
          </template>

          <template #cell-platform="{ value }">
            <span class="badge badge-gray">{{ platformLabel(value) }}</span>
          </template>

          <template #cell-enabled="{ value }">
            <span :class="['badge', value ? 'badge-success' : 'badge-danger']">{{ value ? '启用' : '停用' }}</span>
          </template>

          <template #cell-account_count="{ value }">
            <span class="font-mono text-sm">{{ value || 0 }}</span>
          </template>

          <template #cell-shared_concurrency_limit="{ value }">
            <span class="font-mono text-sm">{{ value || '未启用' }}</span>
          </template>

          <template #cell-updated_at="{ value }">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex flex-wrap gap-1">
              <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="selectAccountSet(row)">
                成员
              </button>
              <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="openAccountSetModal(row)">
                编辑
              </button>
              <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeleteAccountSet(row)">
                删除
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState title="暂无账号集合" description="先新建集合，再批量加入同平台账号。" action-text="新建账号集合" @action="openAccountSetModal()" />
          </template>
        </DataTable>

        <div v-if="selectedAccountSet" class="account-set-members-panel mt-4">
          <div class="mb-3 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <div>
              <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ selectedAccountSet.name }} · 集合成员</h4>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                当前集合平台：{{ platformLabel(selectedAccountSet.platform) }}，只会加入同平台账号。
              </p>
            </div>
            <div class="flex flex-wrap gap-2">
              <button class="btn btn-primary btn-sm" :disabled="accountSetSelectableAccounts.length === 0" @click="openAccountSetMemberPicker">
                <Icon name="plus" size="sm" class="mr-1" />
                选择账号
              </button>
              <button class="btn btn-secondary btn-sm" :disabled="accountSetSelectableAccounts.length === 0" @click="addAccountsToSelectedSet()">
                补齐平台账号
              </button>
            </div>
          </div>
          <DataTable :columns="accountSetMemberColumns" :data="pagedAccountSetMembers" :loading="accountSetMembersLoading" :row-key="accountSetMemberRowKey">
            <template #cell-account_name="{ row, value }">
              <div class="flex flex-col gap-1">
                <div class="font-medium text-gray-900 dark:text-white">{{ value || `账号 #${row.account_id}` }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">{{ row.account_type || '-' }} · {{ row.account_status || '-' }}</div>
              </div>
            </template>

            <template #cell-account_platform="{ value }">
              <span class="badge badge-gray">{{ platformLabel(value) }}</span>
            </template>

            <template #cell-runtime_status="{ row, value }">
              <div class="flex flex-col gap-1">
                <span :class="['badge', runtimeStatusClass(value)]">{{ runtimeStatusLabel(value) }}</span>
                <div
                  v-if="row.runtime_ttft_ms != null || row.runtime_error_rate != null"
                  class="flex flex-wrap gap-1 text-[11px] text-gray-500 dark:text-gray-400"
                >
                  <span
                    v-if="row.runtime_ttft_ms != null"
                    class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
                  >
                    TTFT {{ formatThresholdMs(row.runtime_ttft_ms) }}
                  </span>
                  <span
                    v-if="row.runtime_error_rate != null"
                    class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
                  >
                    错率 {{ formatRateThreshold(row.runtime_error_rate) }}
                  </span>
                </div>
                <span v-if="row.runtime_reason" class="max-w-[320px] whitespace-normal text-xs text-gray-500 dark:text-gray-400">
                  {{ row.runtime_reason }}
                </span>
                <span
                  v-if="row.runtime_rate_limit_reset_at || row.runtime_overload_until || row.runtime_temp_unschedulable_until"
                  class="text-xs text-gray-500 dark:text-gray-400"
                >
                  {{
                    formatCompactDateTime(
                      row.runtime_rate_limit_reset_at ||
                      row.runtime_overload_until ||
                      row.runtime_temp_unschedulable_until
                    )
                  }}
                </span>
              </div>
            </template>

            <template #cell-usage="{ row }">
              <div class="account-set-usage-cell">
                <div class="account-set-usage-cell__main">
                  {{ formatRuntimeLastUsed(row.runtime_last_used_at) }}
                </div>
                <div class="account-set-usage-cell__meta">
                  <span>TTFT {{ formatThresholdMs(row.runtime_ttft_ms) }}</span>
                  <span>错率 {{ formatRateThreshold(row.runtime_error_rate) }}</span>
                </div>
              </div>
            </template>

            <template #cell-capacity="{ row }">
              <div v-if="capacityEditingAccountId === row.account_id" class="flex min-w-[12rem] flex-wrap gap-1">
                <input v-model.number="capacityDraft.hard" class="input w-24 px-2 py-1 text-xs" type="number" min="1" placeholder="硬上限" />
                <input v-model.number="capacityDraft.soft" class="input w-24 px-2 py-1 text-xs" type="number" min="1" placeholder="软份额" />
              </div>
              <div v-else class="text-xs text-gray-600 dark:text-gray-300">
                <span>硬 {{ row.capacity_hard_limit || '共享' }}</span>
                <span class="ml-2">软 {{ row.capacity_soft_share || '—' }}</span>
              </div>
            </template>

            <template #cell-added_at="{ value }">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(value) }}</span>
            </template>

            <template #cell-actions="{ row }">
              <div class="flex flex-wrap gap-1">
                <template v-if="capacityEditingAccountId === row.account_id">
                  <button class="rounded px-2 py-1 text-xs text-primary-700 hover:bg-primary-50 dark:text-primary-300" @click="saveAccountSetMemberCapacity(row)">保存</button>
                  <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300" @click="cancelAccountSetMemberCapacity">取消</button>
                </template>
                <button v-else class="rounded px-2 py-1 text-xs text-primary-700 hover:bg-primary-50 dark:text-primary-300" @click="editAccountSetMemberCapacity(row)">容量</button>
                <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeleteAccountSetMember(row)">移除</button>
              </div>
            </template>

            <template #empty>
              <EmptyState title="集合暂无成员" description="点击右上角把当前平台账号批量加入集合。" />
            </template>
          </DataTable>
          <Pagination
            v-if="accountSetMembersPagination.total > 0"
            class="mt-3"
            :page="accountSetMembersPagination.page"
            :total="accountSetMembersPagination.total"
            :page-size="accountSetMembersPagination.page_size"
            :page-size-options="[10, 20, 50, 100]"
            @update:page="handleAccountSetMembersPageChange"
            @update:pageSize="handleAccountSetMembersPageSizeChange"
          />
        </div>
      </section>

      <section class="sst-admin-panel mb-4 p-4">
        <div class="mb-3 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">共享容量压力</h3>
            <p class="text-xs text-gray-500 dark:text-gray-400">同一上游账号集合的实时并发、余量、排队和近 5 分钟满载观测。</p>
          </div>
          <button
            type="button"
            class="btn btn-secondary btn-sm self-start md:self-auto"
            :disabled="capacityPressuresLoading"
            @click="loadCapacityPressures"
          >
            <Icon name="refresh" size="sm" :class="capacityPressuresLoading ? 'mr-1 animate-spin' : 'mr-1'" />
            刷新压力
          </button>
        </div>
        <div v-if="capacityPressuresLoading" class="rounded-xl border border-dashed border-gray-300 px-3 py-6 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-gray-400">
          正在读取共享容量…
        </div>
        <div v-else-if="capacityPressuresError" class="rounded-xl border border-rose-200 bg-rose-50/70 px-3 py-3 text-sm text-rose-800 dark:border-rose-900/40 dark:bg-rose-950/20 dark:text-rose-200">
          {{ capacityPressuresError }}
        </div>
        <div v-else-if="capacityPressures.length === 0" class="rounded-xl border border-dashed border-gray-300 px-3 py-6 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-gray-400">
          当前没有启用共享容量组。
        </div>
        <div v-else class="grid gap-3 lg:grid-cols-2">
          <article
            v-for="pressure in capacityPressures"
            :key="pressure.set_id"
            class="rounded-2xl border border-gray-200 bg-white/80 p-4 dark:border-dark-700 dark:bg-dark-900/60"
          >
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <h4 class="truncate text-sm font-semibold text-gray-900 dark:text-white">{{ pressure.set_name }}</h4>
                  <span :class="['badge', capacityPressureClass(pressure)]">{{ capacityPressureLabel(pressure) }}</span>
                </div>
                <div class="mt-1 text-[11px] text-gray-500 dark:text-gray-400">{{ pressure.set_code }} · {{ platformLabel(pressure.platform) }}</div>
              </div>
              <div class="text-right">
                <div class="font-mono text-2xl font-semibold text-gray-900 dark:text-white">{{ pressure.current_concurrency }} / {{ pressure.capacity_limit }}</div>
                <div class="text-[11px] text-gray-500 dark:text-gray-400">可用 {{ pressure.available_capacity }} · 排队 {{ pressure.waiting_count }}</div>
              </div>
            </div>
            <div class="mt-3 grid grid-cols-3 gap-2 text-xs">
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">组满载</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.group_full_count }}</strong></div>
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">成员满载</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.member_full_count }}</strong></div>
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">借用槽位</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.borrowed_slot_count }}</strong></div>
            </div>
            <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">5 分钟峰值</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.peak_concurrency_5m }}</strong></div>
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">P95 负载</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.p95_load_rate_5m }}%</strong></div>
              <div class="rounded-xl bg-gray-50 px-3 py-2 dark:bg-dark-800/70"><div class="text-gray-500 dark:text-gray-400">调度集中度</div><strong class="font-mono text-gray-900 dark:text-white">{{ pressure.scheduling_concentration }}%</strong></div>
            </div>
            <div v-if="pressure.members.length > 0" class="mt-3 divide-y divide-gray-200/80 rounded-xl border border-gray-200 dark:divide-dark-700/70 dark:border-dark-700">
              <div v-for="member in pressure.members" :key="member.account_id" class="flex items-center justify-between gap-3 px-3 py-2 text-xs">
                <div class="min-w-0"><div class="truncate font-medium text-gray-800 dark:text-gray-200">{{ member.account_name || `账号 #${member.account_id}` }}</div><div class="text-[11px] text-gray-500 dark:text-gray-400">{{ member.current_concurrency }} 并发 · 排队 {{ member.waiting_count }}</div></div>
                <div class="text-right"><div class="font-mono text-gray-800 dark:text-gray-200">{{ member.load_rate }}%</div><div class="text-[11px] text-gray-500 dark:text-gray-400">{{ member.hard_concurrency_limit ? `硬 ${member.hard_concurrency_limit}` : '共享' }}</div></div>
              </div>
            </div>
          </article>
        </div>
      </section>

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1.28fr)_minmax(460px,0.72fr)]">
        <section class="sst-admin-panel overview-panel p-4">
          <div class="overview-panel__header mb-3">
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">池运行概览</h3>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ selectedPool ? `${selectedPool.name} 的运行状态、命中与诊断；直接成员列表已移到下方整行区域。` : '选择一个池后查看运行概览；直接成员列表会显示在下方整行区域。' }}
            </p>
          </div>
          <div class="overview-sheet mb-3">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ selectedPool?.name || '未选择上游池' }}</h4>
                  <span :class="['badge', poolHealthStatusClass]">{{ poolHealthStatusLabel }}</span>
                  <span class="text-xs text-gray-500 dark:text-gray-400">成员、绑定与 24 小时选路</span>
                </div>
                <div
                  v-if="poolHealthSignals.length > 0"
                  class="mt-2 flex flex-wrap gap-2 text-[11px] overview-signal-list"
                >
                  <span
                    v-for="signal in poolHealthSignals"
                    :key="signal"
                    class="rounded-full bg-amber-50 px-2 py-0.5 text-amber-800 dark:bg-amber-900/20 dark:text-amber-200 overview-signal-chip"
                  >
                    {{ signal }}
                  </span>
                </div>
              </div>
              <button
                class="btn btn-secondary btn-sm overview-refresh-button"
                :disabled="!selectedPool || membersLoading || poolRoutingObservability.loading || poolHealthAlerts.loading"
                @click="refreshSelectedPoolHealth"
              >
                <Icon name="refresh" size="sm" :class="(membersLoading || poolRoutingObservability.loading || poolHealthAlerts.loading) ? 'animate-spin' : ''" class="mr-1" />
                刷新概览
              </button>
            </div>

            <div class="overview-metric-grid mt-3 grid gap-2 md:grid-cols-2 xl:grid-cols-4">
              <div class="overview-metric-card">
                <div class="text-[11px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500 overview-metric-label">调度就绪</div>
                <div class="mt-2 text-3xl font-semibold leading-none text-gray-900 dark:text-white overview-metric-value">{{ poolHealth.readyMembers }}/{{ poolHealth.totalMembers }}</div>
                <div class="mt-2 text-[11px] text-gray-500 dark:text-gray-400 overview-metric-meta">健康 {{ poolHealth.healthyMembers }} · 排空 {{ poolHealth.drainedMembers }}</div>
              </div>
              <div class="overview-metric-card">
                <div class="text-[11px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500 overview-metric-label">运行风险</div>
                <div class="mt-2 text-3xl font-semibold leading-none text-gray-900 dark:text-white overview-metric-value">{{ poolHealth.blockedMembers + poolHealth.errorMembers }}</div>
                <div class="mt-2 text-[11px] text-gray-500 dark:text-gray-400 overview-metric-meta">熔断 {{ poolHealth.blockedMembers }} · 错误 {{ poolHealth.errorMembers }}</div>
              </div>
              <div class="overview-metric-card">
                <div class="text-[11px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500 overview-metric-label">最近命中</div>
                <div class="mt-2 text-3xl font-semibold leading-none text-gray-900 dark:text-white overview-metric-value">{{ poolHealth.observedRoutes }}</div>
                <div class="mt-2 text-[11px] text-gray-500 dark:text-gray-400 overview-metric-meta">账号 {{ poolHealth.uniqueRoutedAccounts }} · 模型 {{ topObservedModelSummary }}</div>
              </div>
              <div class="overview-metric-card">
                <div class="text-[11px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500 overview-metric-label">调度集中度</div>
                <div class="mt-2 truncate text-2xl font-semibold leading-none text-gray-900 dark:text-white overview-metric-value">{{ topRoutedAccountName }}</div>
                <div class="mt-2 text-[11px] text-gray-500 dark:text-gray-400 overview-metric-meta">{{ topRoutedAccountShare }} · 粘性逃逸 {{ poolHealth.stickyEscapes }}</div>
              </div>
            </div>

            <div class="overview-detail-grid mt-3 grid gap-4 border-t border-gray-200/80 pt-3 dark:border-dark-700/80 xl:grid-cols-[minmax(0,1.06fr)_minmax(0,0.94fr)]">
              <div class="min-w-0 xl:pr-4">
                <div class="mb-3 flex items-center justify-between gap-2">
                  <div class="text-xs font-semibold text-gray-700 dark:text-gray-200">绑定与命中</div>
                  <span class="text-[11px] text-gray-500 dark:text-gray-400">覆盖 {{ poolModelCoverageRows.length }} 项</span>
                </div>
                <div class="space-y-3">
                  <div class="min-w-0 overview-subpanel">
                    <div class="mb-2 text-[11px] font-semibold uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500 overview-subpanel-label">绑定覆盖</div>
                    <div v-if="poolModelCoverageRows.length > 0" class="divide-y divide-gray-200/80 dark:divide-dark-700/70">
                      <div
                        v-for="row in poolModelCoverageRows.slice(0, 4)"
                        :key="row.key"
                        class="flex items-center justify-between gap-3 py-2 text-xs overview-coverage-row"
                      >
                        <div class="min-w-0">
                          <div class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ row.label }}</div>
                          <div class="mt-0.5 flex items-center gap-1.5 text-[11px] text-gray-500 dark:text-gray-400">
                            <span class="h-1.5 w-1.5 rounded-full" :class="row.ready ? 'bg-emerald-500' : 'bg-rose-500'"></span>
                            <span>{{ row.ready ? `${poolHealth.readyMembers} 个就绪账号` : '当前无就绪账号' }}</span>
                          </div>
                        </div>
                        <span :class="['badge', row.ready ? 'badge-success' : 'badge-danger']">
                          {{ row.ready ? '可调度' : '待恢复' }}
                        </span>
                      </div>
                    </div>
                    <div v-else class="rounded-lg border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400">
                      当前池还没有启用中的分组绑定。
                    </div>
                  </div>

                  <div class="min-w-0 overview-subpanel">
                    <div class="mb-2 text-xs font-semibold text-gray-700 dark:text-gray-200">最近账号命中</div>
                    <div v-if="poolRoutedAccountRows.length > 0" class="space-y-2">
                      <div
                        v-for="row in poolRoutedAccountRows.slice(0, 2)"
                        :key="row.accountId"
                        class="text-xs overview-hit-row"
                      >
                        <div class="flex items-center justify-between gap-2">
                          <span class="min-w-0 truncate text-gray-700 dark:text-gray-200">{{ row.label }}</span>
                          <span class="font-mono text-gray-500 dark:text-gray-400">{{ row.count }}</span>
                        </div>
                        <div class="mt-1 h-1.5 rounded-full bg-gray-200 dark:bg-dark-700">
                          <div class="h-1.5 rounded-full bg-primary-500" :style="{ width: `${row.percent}%` }"></div>
                        </div>
                      </div>
                    </div>
                    <div v-else class="rounded-lg border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400">
                      最近 24 小时还没有账号命中记录。
                    </div>
                  </div>
                </div>
              </div>

              <div class="min-w-0 xl:border-l xl:border-gray-200/80 xl:pl-4 dark:xl:border-dark-700/80 overview-diagnostic-panel">
                <div class="mb-4 border-b border-gray-200/80 pb-4 dark:border-dark-700/80">
                  <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
                    <div class="flex items-center gap-2">
                      <div class="text-xs font-semibold text-gray-700 dark:text-gray-200">异常预警</div>
                      <span :class="['badge', poolActiveHealthAlerts.length > 0 ? 'badge-warning' : 'badge-success']">
                        {{ poolActiveHealthAlerts.length > 0 ? `${poolActiveHealthAlerts.length} 项持续中` : '当前正常' }}
                      </span>
                    </div>
                    <span class="text-[11px] text-gray-500 dark:text-gray-400">近 7 天</span>
                  </div>

                  <div v-if="poolHealthAlerts.loading" class="rounded-lg bg-gray-50 px-3 py-3 text-xs text-gray-500 dark:bg-dark-800/70 dark:text-gray-400">
                    正在加载最近预警…
                  </div>
                  <div
                    v-else-if="poolHealthAlerts.error"
                    class="flex flex-col gap-2 rounded-lg border border-rose-200 bg-rose-50/70 px-3 py-3 text-xs text-rose-800 dark:border-rose-900/40 dark:bg-rose-950/20 dark:text-rose-200 sm:flex-row sm:items-center sm:justify-between"
                  >
                    <span>{{ poolHealthAlerts.error }}</span>
                    <button type="button" class="font-medium underline underline-offset-2" @click="selectedPool && loadPoolHealthAlerts(selectedPool)">重试</button>
                  </div>
                  <div v-else-if="poolActiveHealthAlerts.length > 0" class="space-y-2">
                    <div
                      v-for="log in poolActiveHealthAlerts.slice(0, 3)"
                      :key="log.id"
                      :class="['rounded-lg border px-3 py-2.5', poolHealthAlertCardClass(log)]"
                    >
                      <div class="flex flex-wrap items-center justify-between gap-2">
                        <div class="flex items-center gap-2">
                          <div class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ formatPoolHealthAlertTitle(log) }}</div>
                          <span :class="['badge', getLogExtraString(log, 'severity') === 'critical' ? 'badge-danger' : 'badge-warning']">
                            {{ getLogExtraString(log, 'severity') === 'critical' ? '严重' : '注意' }}
                          </span>
                        </div>
                        <span class="text-[11px] text-gray-600 dark:text-gray-300">{{ formatCompactDateTime(log.created_at) }}</span>
                      </div>
                      <div class="mt-1 text-xs leading-5 text-gray-700 dark:text-gray-200">{{ formatPoolHealthAlertDetail(log) }}</div>
                    </div>
                    <div v-if="poolActiveHealthAlerts.length > 3" class="text-[11px] text-gray-500 dark:text-gray-400">
                      另有 {{ poolActiveHealthAlerts.length - 3 }} 项持续异常，可在运维日志中查看。
                    </div>
                  </div>
                  <div v-else class="rounded-lg border border-dashed border-emerald-200 bg-emerald-50/50 px-3 py-3 text-xs text-emerald-800 dark:border-emerald-900/40 dark:bg-emerald-950/10 dark:text-emerald-200">
                    当前没有持续异常。
                    <span v-if="latestResolvedPoolHealthAlert" class="ml-1 text-emerald-700/80 dark:text-emerald-300/80">
                      最近恢复：{{ formatPoolHealthAlertTitle(latestResolvedPoolHealthAlert) }} · {{ formatCompactDateTime(latestResolvedPoolHealthAlert.created_at) }}
                    </span>
                  </div>
                </div>

                <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
                  <div class="text-xs font-semibold text-gray-700 dark:text-gray-200">路由诊断</div>
                  <button
                    v-if="poolRoutingObservabilitySupported"
                    class="rounded-lg px-3 py-1.5 text-xs font-medium text-gray-600 transition hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-800 overview-detail-button"
                    :disabled="!selectedPool || poolRoutingObservability.loading || poolRoutingObservability.loadingMore"
                    @click="openPoolObservabilityModal"
                  >
                    查看明细
                  </button>
                </div>
                <div
                  v-if="!poolRoutingObservabilitySupported"
                  class="rounded-lg border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400"
                >
                  当前仅对 OpenAI / Anthropic 池提供路由观测。
                </div>
                <template v-else>
                  <div class="grid gap-4 sm:grid-cols-3">
                    <div class="min-w-0">
                      <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">粘性逃逸</div>
                      <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">{{ stickyEscapeCount }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">{{ stickyEscapeReasonSummary }}</div>
                    </div>
                    <div class="min-w-0">
                      <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">最近账号</div>
                      <div class="mt-1 truncate text-sm font-medium leading-5 text-gray-900 dark:text-white" :title="recentRoutedAccountSummary">{{ recentRoutedAccountSummary }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">按最新记录去重</div>
                    </div>
                    <div class="min-w-0 sm:pl-4">
                      <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">最近时间</div>
                      <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ latestRoutingLogAt }}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">{{ (poolRoutingObservability.loading || poolRoutingObservability.loadingMore) ? '正在刷新…' : '24h 观测窗口' }}</div>
                    </div>
                  </div>
                  <div
                    v-if="poolRoutingObservability.logs.length === 0 && !poolRoutingObservability.loading"
                    class="mt-2 rounded-lg border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400"
                  >
                    最近 24 小时还没有这个池的路由解释日志。
                  </div>
                </template>
              </div>
            </div>
          </div>
        </section>

        <section class="sst-admin-panel current-sets-panel p-4">
          <div class="current-sets-flow mb-4">
            <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <div>
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">当前集合</h3>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ selectedPool ? `${selectedPool.name} 当前挂载的账号集合；运行时会自动展开为成员。` : '先选择一个池，再查看它当前挂载的集合。' }}
                </p>
              </div>
              <button class="btn btn-secondary btn-sm current-sets-manage-button self-start md:self-auto" :disabled="!selectedPool" @click="openMemberSetManagerModal">
                管理集合绑定
              </button>
            </div>
            <div
              v-if="selectedPool"
              class="current-sets-title-strip mt-3 flex flex-wrap items-center gap-2 text-sm"
            >
              <span class="text-base font-semibold text-gray-900 dark:text-white">{{ selectedPool.name }}</span>
              <span class="current-sets-pill">
                {{ platformLabel(selectedPool.platform) }}
              </span>
              <span class="current-sets-pill">
                编码 {{ selectedPool.code }}
              </span>
            </div>
            <div class="current-sets-stats mt-3">
              <div class="current-sets-stat">
                <span class="text-[11px] text-gray-500 dark:text-gray-400">已绑定集合</span>
                <span class="font-mono text-sm font-semibold text-gray-900 dark:text-white">{{ memberSets.length }}</span>
              </div>
              <div class="current-sets-stat">
                <span class="text-[11px] text-gray-500 dark:text-gray-400">启用绑定</span>
                <span class="font-mono text-sm font-semibold text-gray-900 dark:text-white">{{ enabledMemberSetCount }}</span>
              </div>
              <div class="current-sets-stat">
                <span class="text-[11px] text-gray-500 dark:text-gray-400">空集合</span>
                <span :class="['font-mono text-sm font-semibold', emptyBoundMemberSets.length > 0 ? 'text-amber-700 dark:text-amber-200' : 'text-gray-900 dark:text-white']">
                  {{ emptyBoundMemberSets.length }}
                </span>
              </div>
            </div>
            <div v-if="selectedPool && memberSets.length > 0" class="current-sets-list mt-4">
              <div
                v-for="item in memberSets"
                :key="item.id"
                class="current-sets-item"
              >
                <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
                  <div class="min-w-0">
                    <div class="flex flex-wrap items-center gap-2">
                      <div class="truncate text-base font-semibold text-gray-900 dark:text-white">
                        {{ item.set_name || `集合 #${item.set_id}` }}
                      </div>
                      <span class="rounded-full bg-white/85 px-2.5 py-1 text-[11px] text-[#7b6041] ring-1 ring-[#e7d9c3] dark:bg-dark-900/55 dark:text-[#d7c1a6] dark:ring-white/10">
                        {{ platformLabel(item.set_platform) }}
                      </span>
                      <span :class="['badge', item.enabled ? 'badge-success' : 'badge-danger']">{{ item.enabled ? '启用' : '停用' }}</span>
                      <span
                        v-if="getAccountSetMemberCount(item.set_id) === 0"
                        class="rounded-full bg-amber-50 px-2.5 py-1 text-[11px] text-amber-800 ring-1 ring-amber-200 dark:bg-amber-900/20 dark:text-amber-200 dark:ring-amber-900/40"
                      >
                        空集合
                      </span>
                    </div>
                    <div class="mt-1 truncate text-xs text-gray-500 dark:text-gray-400">编码：{{ item.set_code }}</div>
                    <div v-if="item.notes" class="mt-2 max-w-2xl text-sm leading-6 text-[#6e5a45] dark:text-[#cfbea8]">
                      {{ item.notes }}
                    </div>
                    <div
                      v-else-if="getAccountSetMemberCount(item.set_id) === 0"
                      class="mt-2 text-sm leading-6 text-amber-800 dark:text-amber-200"
                    >
                      当前集合为空，不会为这个池展开出可用成员。
                    </div>
                  </div>
                  <div class="current-sets-item-meta">
                    <span class="font-mono text-lg font-semibold text-gray-900 dark:text-white">{{ getAccountSetMemberCount(item.set_id) }}</span>
                    <span class="text-xs text-gray-500 dark:text-gray-400">账号</span>
                    <span class="hidden h-3 w-px bg-[#ddd1bf] lg:inline-block"></span>
                    <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(item.updated_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div
              v-else-if="selectedPool"
              class="mt-4 rounded-2xl border border-dashed border-[#d6c6ae] bg-[rgba(252,248,242,0.85)] px-4 py-5 text-sm text-[#7a654b] dark:border-[#4e4031] dark:bg-[rgba(24,20,17,0.72)] dark:text-[#ccb79a]"
            >
              当前池还没有挂载任何集合。左侧可维护直接成员，若希望按集合管理，请点击“管理集合绑定”添加当前集合。
            </div>
            <div
              v-if="emptyBoundMemberSets.length > 0"
              class="mt-3 rounded-xl border border-amber-200 bg-amber-50/80 px-3 py-2 text-xs text-amber-900 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-100"
            >
              已绑定但为空的集合：
              {{ emptyBoundMemberSets.map(item => item.set_name || `集合 #${item.set_id}`).join(' / ') }}。
              这些集合当前不会为池提供可用成员。
            </div>
          </div>

          <div class="mb-3 flex items-center justify-between gap-2">
            <div>
              <h3 class="text-base font-semibold text-gray-900 dark:text-white">绑定</h3>
              <p class="text-xs text-gray-500 dark:text-gray-400">分组绑定到池，路由层据此接管请求。</p>
            </div>
            <button class="btn btn-secondary btn-sm" @click="openBindingModal()">
              添加绑定
            </button>
          </div>
          <p class="mb-3 text-xs text-gray-500 dark:text-gray-400">
            分组一旦绑定到池，请求只会从这个池的成员里选账号。
          </p>
          <DataTable :columns="bindingColumns" :data="bindings" :loading="bindingsLoading" row-key="id">
            <template #cell-group_name="{ row, value }">
              <div class="flex flex-col gap-1">
                <div class="font-medium text-gray-900 dark:text-white">{{ value || `分组 #${row.group_id}` }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">{{ row.group_platform || '未知平台' }}</div>
              </div>
            </template>

            <template #cell-enabled="{ value }">
              <span :class="['badge', value ? 'badge-success' : 'badge-danger']">{{ value ? '启用' : '停用' }}</span>
            </template>

            <template #cell-platform="{ value }">
              <span class="badge badge-gray">{{ platformLabel(value) }}</span>
            </template>

            <template #cell-actions="{ row }">
              <div class="flex gap-1">
                <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="openBindingModal(row)">
                  编辑
                </button>
                <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeleteBinding(row)">
                  删除
                </button>
              </div>
            </template>

            <template #empty>
              <EmptyState title="暂无绑定" description="先将分组绑定到上游池。" action-text="添加绑定" @action="openBindingModal()" />
            </template>
          </DataTable>
        </section>
      </div>

      <section class="sst-admin-panel mt-4 p-4">
        <div class="mb-3 flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
          <div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">直接成员</h3>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ selectedPool ? `${selectedPool.name} 的直接挂载账号；集合请看上方右侧“当前集合”。` : '选择一个池后查看直接成员。' }}
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button
              class="btn btn-secondary btn-sm"
              :disabled="!selectedPool || !canSyncSelectedPool || syncingMembers"
              @click="addMissingSelectedPoolMembers()"
            >
              <Icon name="plus" size="sm" class="mr-1" />
              补齐缺失账号
            </button>
            <button
              class="btn btn-secondary btn-sm"
              :disabled="!selectedPool || !canSyncSelectedPool || syncingMembers"
              @click="syncSelectedPoolMembers()"
            >
              <Icon name="refresh" size="sm" :class="syncingMembers ? 'animate-spin' : ''" class="mr-1" />
              全量同步账号
            </button>
            <button class="btn btn-secondary btn-sm" :disabled="!selectedPool" @click="openMemberModal()">
              添加成员
            </button>
          </div>
        </div>
        <div
          class="mb-3 flex flex-wrap items-center gap-2 rounded-xl border border-primary-200 bg-primary-50/70 px-3 py-2 text-sm text-primary-900 dark:border-primary-900/40 dark:bg-primary-950/20 dark:text-primary-100"
        >
          <span class="font-medium">当前操作池：</span>
          <span v-if="selectedPool">{{ selectedPool.name }}</span>
          <span v-else>未选择</span>
          <span v-if="selectedPool" class="rounded bg-white/80 px-2 py-0.5 text-xs dark:bg-dark-900/60">{{ platformLabel(selectedPool.platform) }}</span>
          <span v-if="selectedPool" class="rounded bg-white/80 px-2 py-0.5 text-xs dark:bg-dark-900/60">编码 {{ selectedPool.code }}</span>
          <span v-if="selectedPool" class="rounded bg-white/80 px-2 py-0.5 text-xs dark:bg-dark-900/60">
            逃逸 {{ selectedPool.sticky_escape_enabled ? '开' : '关' }}
          </span>
          <span v-if="selectedPool" class="rounded bg-white/80 px-2 py-0.5 text-xs dark:bg-dark-900/60">
            TTFT {{ formatThresholdMs(selectedPool.sticky_escape_ttft_ms_threshold) }}
          </span>
          <span v-if="selectedPool" class="rounded bg-white/80 px-2 py-0.5 text-xs dark:bg-dark-900/60">
            错率 {{ formatRateThreshold(selectedPool.sticky_escape_error_rate_threshold) }}
          </span>
        </div>
        <div class="mb-3 flex flex-wrap items-center gap-2 rounded-xl border border-amber-200 bg-amber-50/80 px-3 py-2 text-xs text-amber-900 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-100">
          <span class="font-medium">恢复探针</span>
          <span>自动探测异常/可恢复成员；每分钟扫描，近 2 分钟内刚用过的成员先跳过；API Key 优先使用轻量 Hi 探针，运行态同时参考 Claude Code 与 OpenAI 最近探针。</span>
        </div>
        <DataTable :columns="memberColumns" :data="directMembers" :loading="membersLoading" :row-key="poolMemberRowKey">
          <template #cell-account_name="{ row, value }">
            <div class="flex flex-col gap-1">
              <div class="pool-member-account-name whitespace-normal break-words font-medium leading-6 text-gray-900 dark:text-white">
                {{ formatPoolMemberAccountName(value, row.account_id) }}
              </div>
              <div class="flex flex-wrap gap-1 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ row.account_platform || '未知平台' }}</span>
                <span>·</span>
                <span>{{ row.account_type || '-' }}</span>
              </div>
            </div>
          </template>

          <template #cell-enabled="{ value }">
            <span :class="['badge', value ? 'badge-success' : 'badge-danger']">{{ value ? '启用' : '停用' }}</span>
          </template>

          <template #cell-runtime_status="{ row, value }">
            <div class="flex flex-col gap-1">
              <span :class="['badge', runtimeStatusClass(value)]">{{ runtimeStatusLabel(value) }}</span>
              <div
                v-if="row.runtime_ttft_ms != null || row.runtime_error_rate != null"
                class="flex flex-wrap gap-1 text-[11px] text-gray-500 dark:text-gray-400"
              >
                <span
                  v-if="row.runtime_ttft_ms != null"
                  class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
                >
                  TTFT {{ formatThresholdMs(row.runtime_ttft_ms) }}
                </span>
                <span
                  v-if="row.runtime_error_rate != null"
                  class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
                >
                  错率 {{ formatRateThreshold(row.runtime_error_rate) }}
                </span>
              </div>
              <span v-if="row.runtime_reason" class="max-w-[320px] whitespace-normal text-xs text-gray-500 dark:text-gray-400">
                {{ row.runtime_reason }}
              </span>
              <span
                v-if="row.runtime_rate_limit_reset_at || row.runtime_overload_until || row.runtime_temp_unschedulable_until"
                class="text-xs text-gray-500 dark:text-gray-400"
              >
                {{
                  formatCompactDateTime(
                    row.runtime_rate_limit_reset_at ||
                    row.runtime_overload_until ||
                    row.runtime_temp_unschedulable_until
                  )
                }}
              </span>
            </div>
          </template>

          <template #cell-manual_drained="{ value }">
            <span :class="['badge', value ? 'badge-warning' : 'badge-gray']">{{ value ? '已排空' : '正常' }}</span>
          </template>

          <template #cell-weight="{ value, row }">
            <div class="flex flex-col">
              <span class="font-mono text-sm">
                {{ value }}<template v-if="row.effective_weight && row.effective_weight !== value"> → {{ row.effective_weight }}</template>
              </span>
              <span v-if="row.runtime_weight_factor && row.runtime_weight_factor !== 1" class="text-xs text-gray-500 dark:text-gray-400">
                ×{{ Number(row.runtime_weight_factor).toFixed(2) }} · {{ formatAutoWeightReason(row.runtime_weight_reason) }}
              </span>
            </div>
          </template>

          <template #cell-updated_at="{ value }">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex gap-1">
              <button
                class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-50 dark:text-gray-300 dark:hover:bg-dark-700"
                :disabled="row.editable === false"
                @click="openMemberModal(row)"
              >
                编辑
              </button>
              <button
                class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50 dark:text-red-300 dark:hover:bg-red-900/20"
                :disabled="row.editable === false"
                @click="confirmDeleteMember(row)"
              >
                删除
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState title="暂无直接成员" description="当前池没有直接挂载账号；如果你使用集合，请看上方右侧“当前集合”。" action-text="添加成员" @action="openMemberModal()" />
          </template>
        </DataTable>
      </section>
    </div>

    <BaseDialog
      :show="showPoolObservabilityModal"
      :title="selectedPool ? `${selectedPool.name} · 最近路由观测` : '最近路由观测'"
      width="extra-wide"
      :close-on-click-outside="true"
      @close="closePoolObservabilityModal"
    >
      <div class="space-y-4">
        <div class="flex flex-col gap-3 rounded-2xl border border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/50 md:flex-row md:items-center md:justify-between">
          <div class="min-w-0">
            <div class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ selectedPool?.name || '-' }}
              <span class="ml-2 rounded bg-white/80 px-2 py-0.5 text-xs font-normal text-gray-500 dark:bg-dark-800 dark:text-gray-400">
                最近 24 小时
              </span>
            </div>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              已载入 {{ poolRoutingPaginationInfo.loaded }} / {{ poolRoutingPaginationInfo.total }} 条记录，按日志时间倒序展示。
            </p>
          </div>
          <button
            class="btn btn-secondary btn-sm self-start md:self-auto"
            :disabled="!selectedPool || !poolRoutingObservabilitySupported || poolRoutingObservability.loading || poolRoutingObservability.loadingMore"
            @click="selectedPool && loadPoolObservability(selectedPool)"
          >
            <Icon name="refresh" size="sm" :class="(poolRoutingObservability.loading || poolRoutingObservability.loadingMore) ? 'animate-spin' : ''" class="mr-1" />
            刷新
          </button>
        </div>

        <div
          v-if="poolRoutingObservability.loading"
          class="rounded-xl border border-dashed border-gray-300 px-3 py-6 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-gray-400"
        >
          正在加载最近路由观测…
        </div>

        <div
          v-else-if="poolRoutingObservability.logs.length > 0"
          class="max-h-[65vh] space-y-2 overflow-y-auto pr-1"
        >
          <div
            v-for="log in poolRoutingObservability.logs"
            :key="log.id"
            class="rounded-xl border border-gray-200 bg-white/85 px-3 py-3 dark:border-dark-700 dark:bg-dark-900/70"
          >
            <div class="flex flex-wrap items-center gap-2 text-xs">
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ formatCompactDateTime(log.created_at) }}</span>
              <span class="rounded bg-gray-100 px-2 py-0.5 text-gray-700 dark:bg-dark-700 dark:text-gray-200">
                {{ formatRoutingReason(getLogExtraString(log, 'reason')) }}
              </span>
              <span class="rounded bg-gray-100 px-2 py-0.5 text-gray-700 dark:bg-dark-700 dark:text-gray-200">
                命中 {{ formatObservabilityAccountLabel(getObservabilitySelectedAccountId(log)) }}
              </span>
              <span
                v-if="getLogExtraBool(log, 'sticky_escape_triggered')"
                class="rounded bg-amber-100 px-2 py-0.5 text-amber-800 dark:bg-amber-900/30 dark:text-amber-200"
              >
                触发逃逸 · {{ formatStickyEscapeReason(getLogExtraString(log, 'sticky_escape_reason')) }}
              </span>
            </div>
            <div
              v-if="getObservabilityCandidateCount(log) != null || getObservabilityTopKLabels(log).length > 0 || getObservabilitySkippedEntries(log).length > 0"
              class="mt-2 rounded-xl border border-[#e6dccd] bg-[#faf7f2]/90 px-3 py-3 dark:border-[#4e4031] dark:bg-[rgba(41,33,26,0.72)]"
            >
              <div class="flex flex-wrap items-center gap-2 text-[11px] text-gray-500 dark:text-gray-400">
                <span
                  v-if="getObservabilityCandidateCount(log) != null"
                  class="rounded bg-white/90 px-2 py-0.5 text-gray-700 ring-1 ring-gray-200 dark:bg-dark-800/70 dark:text-gray-200 dark:ring-dark-600"
                >
                  候选 {{ getObservabilityCandidateCount(log) }}
                </span>
                <span
                  v-if="getObservabilityTopKLabels(log).length > 0"
                  class="rounded bg-white/90 px-2 py-0.5 text-gray-700 ring-1 ring-gray-200 dark:bg-dark-800/70 dark:text-gray-200 dark:ring-dark-600"
                >
                  TopK {{ getObservabilityTopKLabels(log).length }}
                </span>
              </div>
              <div v-if="getObservabilityTopKLabels(log).length > 0" class="mt-2">
                <div class="text-[11px] font-medium uppercase tracking-[0.14em] text-[#8a6d4a] dark:text-[#c5ac8f]">TopK 候选池</div>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="label in getObservabilityTopKLabels(log)"
                    :key="label"
                    class="rounded-full border border-[#e1d4c3] bg-white/90 px-2.5 py-1 text-xs text-[#6f573d] dark:border-[#5a4838] dark:bg-dark-800/70 dark:text-[#dbc7b0]"
                  >
                    {{ label }}
                  </span>
                </div>
              </div>
              <div v-if="getObservabilitySkippedEntries(log).length > 0" class="mt-3">
                <div class="text-[11px] font-medium uppercase tracking-[0.14em] text-[#8a6d4a] dark:text-[#c5ac8f]">跳过原因</div>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="entry in getObservabilitySkippedEntries(log)"
                    :key="entry.key"
                    class="rounded-full border border-amber-200 bg-amber-50/90 px-2.5 py-1 text-xs text-amber-800 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-200"
                  >
                    {{ entry.label }} {{ entry.count }}
                  </span>
                </div>
              </div>
            </div>
            <div class="mt-2 flex flex-wrap gap-2 text-[11px] text-gray-500 dark:text-gray-400">
              <span
                v-if="getLogExtraNumber(log, 'sticky_escape_observed_ttft_ms') != null"
                class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
              >
                观测 TTFT {{ formatThresholdMs(getLogExtraNumber(log, 'sticky_escape_observed_ttft_ms')) }}
              </span>
              <span
                v-if="getLogExtraNumber(log, 'sticky_escape_observed_error_rate') != null"
                class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
              >
                观测错率 {{ formatRateThreshold(getLogExtraNumber(log, 'sticky_escape_observed_error_rate')) }}
              </span>
              <span
                v-if="getLogExtraNumber(log, 'routing_latency_ms') != null"
                class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
              >
                选路耗时 {{ getLogExtraNumber(log, 'routing_latency_ms') }}ms
              </span>
              <span
                v-if="getLogExtraString(log, 'pool_name') || getLogExtraString(log, 'pool_code')"
                class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
              >
                池 {{ getLogExtraString(log, 'pool_name') || selectedPool?.name || '-' }}
                <template v-if="getLogExtraString(log, 'pool_code')">· {{ getLogExtraString(log, 'pool_code') }}</template>
              </span>
              <span
                v-if="getLogExtraString(log, 'required_capability')"
                class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700"
              >
                能力 {{ getLogExtraString(log, 'required_capability') }}
              </span>
            </div>
          </div>

          <div class="sticky bottom-0 flex flex-col gap-2 rounded-xl border border-gray-200 bg-white/95 px-3 py-3 shadow-sm backdrop-blur dark:border-dark-700 dark:bg-dark-900/95 sm:flex-row sm:items-center sm:justify-between">
            <div class="text-xs text-gray-500 dark:text-gray-400">
              {{ poolRoutingPaginationInfo.hasMore ? `还有 ${poolRoutingPaginationInfo.total - poolRoutingPaginationInfo.loaded} 条未载入` : '当前筛选范围已全部载入' }}
            </div>
            <button
              type="button"
              class="btn btn-secondary btn-sm"
              :disabled="poolRoutingPaginationInfo.disabled"
              @click="loadMorePoolObservability"
            >
              <Icon
                v-if="poolRoutingObservability.loadingMore"
                name="refresh"
                size="sm"
                class="mr-1 animate-spin"
              />
              {{ poolRoutingObservability.loadingMore ? '加载中…' : poolRoutingPaginationInfo.label }}
            </button>
          </div>
        </div>

        <div
          v-else
          class="rounded-xl border border-dashed border-gray-300 px-3 py-6 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-gray-400"
        >
          最近 24 小时还没有这个池的路由解释日志。
        </div>
      </div>
    </BaseDialog>

    <BaseDialog :show="showPoolModal" :title="editingPool ? '编辑上游池' : '新建上游池'" width="wide" @close="closePoolModal">
      <form class="space-y-4" @submit.prevent="submitPool">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">名称</label>
            <p class="field-hint">用于后台识别这个池。</p>
            <input v-model="poolForm.name" class="input" />
          </div>
          <div>
            <label class="input-label">编码</label>
            <p class="field-hint">唯一标识，保存后尽量不改。</p>
            <input v-model="poolForm.code" class="input" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">平台</label>
            <p class="field-hint">决定池内账号的平台类型。</p>
            <Select v-model="poolForm.platform" :options="platformOptions" :disabled="editingPoolPlatformLocked" />
          </div>
          <div>
            <label class="input-label">描述</label>
            <p class="field-hint">给管理员看的备注。</p>
            <input v-model="poolForm.description" class="input" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-4">
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">启用</span>
              <span class="field-hint mb-0 block">关闭后整个池不参与调度。</span>
            </span>
            <Toggle :model-value="poolForm.enabled" @update:modelValue="poolForm.enabled = $event" />
          </label>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">粘性会话</span>
              <span class="field-hint mb-0 block">同一会话尽量走同一账号。</span>
            </span>
            <Toggle :model-value="poolForm.sticky_enabled" @update:modelValue="poolForm.sticky_enabled = $event" />
          </label>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">负载均衡</span>
              <span class="field-hint mb-0 block">在池内账号间分摊请求。</span>
            </span>
            <Toggle :model-value="poolForm.load_balance_enabled" @update:modelValue="poolForm.load_balance_enabled = $event" />
          </label>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">粘性逃逸</span>
              <span class="field-hint mb-0 block">sticky 账号太慢或错误率过高时跳出。</span>
            </span>
            <Toggle :model-value="poolForm.sticky_escape_enabled" @update:modelValue="poolForm.sticky_escape_enabled = $event" />
          </label>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">故障转移</span>
              <span class="field-hint mb-0 block">请求失败后允许切换到候选账号。</span>
            </span>
            <Toggle :model-value="poolForm.failover_enabled" @update:modelValue="poolForm.failover_enabled = $event" />
          </label>
          <div>
            <label class="input-label">自动健康调权</label>
            <p class="field-hint">观察模式只保存建议，不改变实际权重。</p>
            <Select v-model="poolForm.auto_weight_mode" :options="autoWeightModeOptions" :disabled="poolForm.platform !== 'openai'" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-3">
          <div>
            <label class="input-label">逃逸 TTFT 阈值</label>
            <p class="field-hint">首 token 均值超过这里，就不再强粘这个账号。</p>
            <input v-model.number="poolForm.sticky_escape_ttft_ms_threshold" type="number" min="1" class="input" />
          </div>
          <div>
            <label class="input-label">逃逸错误率阈值</label>
            <p class="field-hint">错误率超过这里，就不再强粘这个账号。</p>
            <input v-model.number="poolForm.sticky_escape_error_rate_threshold" type="number" min="0" max="1" step="0.01" class="input" />
          </div>
          <div>
            <label class="input-label">TopK</label>
            <p class="field-hint">每次从前 N 个候选里选。</p>
            <input v-model.number="poolForm.top_k" type="number" min="1" class="input" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-3">
          <div>
            <label class="input-label">故障跳数</label>
            <p class="field-hint">失败后最多切换几次。</p>
            <input v-model.number="poolForm.max_failover_hops" type="number" min="0" class="input" />
          </div>
          <div>
            <label class="input-label">等待毫秒</label>
            <p class="field-hint">等待可用账号的最长时间。</p>
            <input v-model.number="poolForm.wait_timeout_ms" type="number" min="0" class="input" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-3">
          <div>
            <label class="input-label">成员上限</label>
            <p class="field-hint">同时排队等待的请求上限。</p>
            <input v-model.number="poolForm.max_waiting" type="number" min="0" class="input" />
          </div>
          <div>
            <label class="input-label">调度模式</label>
            <p class="field-hint">默认 advanced 即可。</p>
            <input v-model="poolForm.scheduler_mode" class="input" />
          </div>
          <div>
            <label class="input-label">账号类型策略</label>
            <p class="field-hint">OAuth 优先时，缓存亲和仍在 OAuth 候选内生效。</p>
            <Select v-model="poolForm.account_type_strategy" :options="accountTypeStrategyOptions" />
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="closePoolModal">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">保存</button>
        </div>
      </form>
    </BaseDialog>

    <BaseDialog :show="showMemberModal" :title="editingMember ? '编辑成员' : '添加成员'" width="wide" @close="closeMemberModal">
      <form class="space-y-4" @submit.prevent="submitMember">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">账号</label>
            <p class="field-hint">加入这个池的上游账号。</p>
            <Select v-model="memberForm.account_id" :options="accountOptions" />
          </div>
          <div>
            <label class="input-label">权重</label>
            <p class="field-hint">影响分流比例；想固定优先走请用优先级覆盖。</p>
            <input v-model.number="memberForm.weight" type="number" min="1" class="input" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">启用</span>
              <span class="field-hint mb-0 block">关闭后该成员不参与调度。</span>
            </span>
            <Toggle :model-value="memberForm.enabled" @update:modelValue="memberForm.enabled = $event" />
          </label>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">手动排空</span>
              <span class="field-hint mb-0 block">临时停用，保留配置。</span>
            </span>
            <Toggle :model-value="memberForm.manual_drained" @update:modelValue="memberForm.manual_drained = $event" />
          </label>
        </div>
        <div class="grid gap-4 md:grid-cols-3">
          <div>
            <label class="input-label">调度覆盖</label>
            <p class="field-hint">覆盖账号自身是否可调度。</p>
            <Select v-model="memberForm.schedulable_override" :options="boolSelectOptions" />
          </div>
          <div>
            <label class="input-label">优先级覆盖</label>
            <p class="field-hint">覆盖账号原有优先级。</p>
            <input v-model.number="memberForm.priority_override" type="number" class="input" />
          </div>
          <div>
            <label class="input-label">并发覆盖</label>
            <p class="field-hint">覆盖账号并发上限。</p>
            <input v-model.number="memberForm.max_concurrency_override" type="number" class="input" />
          </div>
        </div>
        <div>
          <label class="input-label">备注</label>
          <p class="field-hint">仅给管理员查看。</p>
          <textarea v-model="memberForm.notes" rows="2" class="input" />
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="closeMemberModal">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">保存</button>
        </div>
      </form>
    </BaseDialog>

    <BaseDialog :show="showAccountSetModal" :title="editingAccountSet ? '编辑账号集合' : '新建账号集合'" width="wide" @close="closeAccountSetModal">
      <form class="space-y-4" @submit.prevent="submitAccountSet">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">集合名称</label>
            <p class="field-hint">建议按用途或账号来源命名。</p>
            <input v-model="accountSetForm.name" class="input" />
          </div>
          <div v-if="editingAccountSet">
            <label class="input-label">集合编码</label>
            <p class="field-hint">系统生成的唯一标识，创建后不再手动修改。</p>
            <input :value="accountSetForm.code" class="input bg-gray-50 text-gray-500 dark:bg-dark-800 dark:text-gray-400" readonly />
          </div>
          <div v-else class="rounded-xl border border-dashed border-gray-200 px-3 py-3 dark:border-dark-700">
            <label class="input-label">集合编码</label>
            <p class="field-hint mb-0">保存时由系统自动生成，无需手动填写。</p>
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">平台</label>
            <p class="field-hint">集合成员必须和这里的平台一致。</p>
            <Select v-model="accountSetForm.platform" :options="platformOptions" :disabled="editingAccountSetPlatformLocked" />
          </div>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">启用</span>
              <span class="field-hint mb-0 block">停用后，绑定它的池不会再展开这个集合。</span>
            </span>
            <Toggle :model-value="accountSetForm.enabled" @update:modelValue="accountSetForm.enabled = $event" />
          </label>
        </div>
        <div>
          <label class="input-label">共享并发上限</label>
          <p class="field-hint">留空表示不启用共享容量组；启用后集合成员共享此总上限。</p>
          <input v-model.number="accountSetForm.shared_concurrency_limit" type="number" min="1" step="1" class="input" placeholder="例如 3000" />
        </div>
        <div>
          <label class="input-label">备注</label>
          <p class="field-hint">说明这批账号的来源或用途。</p>
          <textarea v-model="accountSetForm.description" rows="2" class="input" />
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="closeAccountSetModal">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">保存</button>
        </div>
      </form>
    </BaseDialog>

    <BaseDialog
      :show="showAccountSetMemberPicker"
      :title="selectedAccountSet ? `${selectedAccountSet.name} · 选择账号` : '选择账号'"
      width="extra-wide"
      @close="closeAccountSetMemberPicker"
    >
      <div class="space-y-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div class="relative min-w-0 flex-1 md:max-w-md">
            <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input v-model="accountSetMemberSearch" class="input pl-9" placeholder="搜索账号名称或 ID" />
          </div>
          <div class="inline-flex w-full rounded border border-gray-200 p-1 dark:border-dark-600 md:w-auto" role="group" aria-label="账号类型">
            <button
              v-for="option in accountSetMemberTypeOptions"
              :key="option.value"
              type="button"
              :class="['min-h-9 flex-1 px-3 text-sm md:flex-none', accountSetMemberType === option.value ? 'bg-primary-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700']"
              @click="accountSetMemberType = option.value"
            >
              {{ option.label }}
            </button>
          </div>
        </div>

        <div class="flex flex-col gap-2 border-b border-gray-200 pb-3 text-sm dark:border-dark-700 sm:flex-row sm:items-center sm:justify-between">
          <label class="inline-flex items-center gap-2 text-gray-700 dark:text-gray-200">
            <input
              type="checkbox"
              :checked="allFilteredAccountSetCandidatesSelected"
              :disabled="filteredAccountSetCandidates.length === 0"
              @change="toggleAllFilteredAccountSetCandidates"
            />
            选择当前结果
          </label>
          <span class="text-gray-500 dark:text-gray-400">已选 {{ selectedAccountSetMemberIDs.length }} 个</span>
        </div>

        <div class="max-h-[26rem] overflow-y-auto rounded border border-gray-200 dark:border-dark-700">
          <label
            v-for="account in filteredAccountSetCandidates"
            :key="account.id"
            class="grid min-h-14 cursor-pointer grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-3 border-b border-gray-100 px-3 py-2 last:border-b-0 hover:bg-gray-50 dark:border-dark-700 dark:hover:bg-dark-800"
          >
            <input v-model="selectedAccountSetMemberIDs" type="checkbox" :value="account.id" />
            <span class="min-w-0">
              <span class="block truncate text-sm font-medium text-gray-900 dark:text-white">{{ account.name }}</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400">#{{ account.id }} · {{ accountSetAccountTypeLabel(account.type) }}</span>
            </span>
            <span :class="['badge', account.status === 'active' ? 'badge-success' : 'badge-gray']">{{ account.status }}</span>
          </label>
          <div v-if="filteredAccountSetCandidates.length === 0" class="px-4 py-10 text-center text-sm text-gray-500 dark:text-gray-400">
            没有符合条件的可选账号
          </div>
        </div>

        <div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button type="button" class="btn btn-secondary" @click="closeAccountSetMemberPicker">取消</button>
          <button type="button" class="btn btn-primary" :disabled="submitting || selectedAccountSetMemberIDs.length === 0" @click="addSelectedAccountsToSet">
            加入 {{ selectedAccountSetMemberIDs.length }} 个账号
          </button>
        </div>
      </div>
    </BaseDialog>

    <BaseDialog
      :show="showMemberSetManagerModal"
      :title="selectedPool ? `${selectedPool.name} · 集合绑定管理` : '集合绑定管理'"
      width="extra-wide"
      @close="closeMemberSetManagerModal"
    >
      <div class="space-y-4">
        <div class="flex flex-col gap-3 rounded-2xl border border-gray-200 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/50 md:flex-row md:items-center md:justify-between">
          <div class="min-w-0">
            <div class="text-sm font-semibold text-gray-900 dark:text-white">集合成员绑定</div>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">把账号集合挂到当前池，运行时会自动展开为池成员。</p>
          </div>
          <button class="btn btn-primary btn-sm self-start md:self-auto" :disabled="!selectedPool" @click="openMemberSetEditor()">
            绑定集合
          </button>
        </div>
        <div
          v-if="emptyBoundMemberSets.length > 0"
          class="rounded-xl border border-amber-200 bg-amber-50/80 px-3 py-2 text-xs text-amber-900 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-100"
        >
          已绑定但为空的集合：
          {{ emptyBoundMemberSets.map(item => item.set_name || `集合 #${item.set_id}`).join(' / ') }}。
          这些集合当前不会为池提供可用成员。
        </div>
        <DataTable :columns="memberSetColumns" :data="memberSets" :loading="memberSetsLoading" row-key="id">
          <template #cell-set_name="{ row, value }">
            <div class="flex flex-col gap-1">
              <div class="font-medium text-gray-900 dark:text-white">{{ value || `集合 #${row.set_id}` }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">编码：{{ row.set_code }}</div>
              <div v-if="getAccountSetMemberCount(row.set_id) === 0" class="text-xs text-amber-700 dark:text-amber-200">
                当前集合为空，不会展开出成员
              </div>
            </div>
          </template>

          <template #cell-set_platform="{ value }">
            <span class="badge badge-gray">{{ platformLabel(value) }}</span>
          </template>

          <template #cell-account_count="{ row }">
            <div class="flex flex-col gap-1">
              <span class="font-mono text-sm text-gray-900 dark:text-white">{{ getAccountSetMemberCount(row.set_id) }}</span>
              <span
                v-if="getAccountSetMemberCount(row.set_id) === 0"
                class="inline-flex w-fit rounded bg-amber-50 px-2 py-0.5 text-[11px] text-amber-800 dark:bg-amber-900/20 dark:text-amber-200"
              >
                空集合
              </span>
            </div>
          </template>

          <template #cell-enabled="{ value }">
            <span :class="['badge', value ? 'badge-success' : 'badge-danger']">{{ value ? '启用' : '停用' }}</span>
          </template>

          <template #cell-updated_at="{ value }">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="flex gap-1">
              <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="openMemberSetEditor(row)">
                编辑
              </button>
              <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeleteMemberSet(row)">
                删除
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState title="暂无集合绑定" description="先选择集合，再把它绑定到当前池。" action-text="绑定集合" @action="openMemberSetEditor()" />
          </template>
        </DataTable>
      </div>
    </BaseDialog>

    <BaseDialog :show="showMemberSetModal" :title="editingMemberSet ? '编辑集合绑定' : '绑定集合到池'" width="wide" @close="closeMemberSetModal">
      <form class="space-y-4" @submit.prevent="submitMemberSet">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">账号集合</label>
            <p class="field-hint">绑定后，池运行时会自动展开其中账号。</p>
            <Select v-model="memberSetForm.set_id" :options="accountSetOptionsForPool" />
          </div>
          <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
            <span>
              <span class="block text-sm">启用</span>
              <span class="field-hint mb-0 block">关闭后保留绑定配置，但不参与展开。</span>
            </span>
            <Toggle :model-value="memberSetForm.enabled" @update:modelValue="memberSetForm.enabled = $event" />
          </label>
        </div>
        <div>
          <label class="input-label">备注</label>
          <p class="field-hint">可记录这个绑定的用途。</p>
          <textarea v-model="memberSetForm.notes" rows="2" class="input" />
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="closeMemberSetModal">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">保存</button>
        </div>
      </form>
    </BaseDialog>

    <BaseDialog :show="showBindingModal" :title="editingBinding ? '编辑绑定' : '添加绑定'" width="wide" @close="closeBindingModal">
      <form class="space-y-4" @submit.prevent="submitBinding">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">分组</label>
            <p class="field-hint">哪些用户分组走这个池。</p>
            <Select v-model="bindingForm.group_id" :options="groupOptionsForBinding" />
          </div>
          <div>
            <label class="input-label">上游池</label>
            <p class="field-hint">请求要分配到的账号池。</p>
            <Select v-model="bindingForm.pool_id" :options="poolOptions" @change="handleBindingPoolChange" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">平台</label>
            <p class="field-hint">必须和池的平台一致。</p>
            <Select v-model="bindingForm.platform" :options="platformOptions" disabled />
          </div>
          <div>
            <label class="input-label">优先级</label>
            <p class="field-hint">数字越小越先匹配。</p>
            <input v-model.number="bindingForm.priority" type="number" class="input" />
          </div>
        </div>
        <div>
          <label class="input-label">模型</label>
          <p class="field-hint">留空表示全部模型。</p>
          <input v-model="bindingForm.modelsText" class="input" placeholder="用英文逗号分隔" />
        </div>
        <div>
          <label class="input-label">路径范围</label>
          <p class="field-hint">留空表示全部接口路径。</p>
          <input v-model="bindingForm.request_path_scopeText" class="input" placeholder="用英文逗号分隔" />
        </div>
        <label class="flex items-center justify-between gap-3 rounded-xl border border-gray-200 px-3 py-2 dark:border-dark-600">
          <span>
            <span class="block text-sm">启用</span>
            <span class="field-hint mb-0 block">关闭后这条绑定不参与调度。</span>
          </span>
          <Toggle :model-value="bindingForm.enabled" @update:modelValue="bindingForm.enabled = $event" />
        </label>
        <div class="flex justify-end gap-2 pt-2">
          <button type="button" class="btn btn-secondary" @click="closeBindingModal">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="submitting">保存</button>
        </div>
      </form>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'
import { adminAPI } from '@/api/admin'
import type { OpsSystemLog } from '@/api/admin/ops'
import type {
  Account,
  AdminGroup,
  UpstreamAccountSet,
  UpstreamAccountSetMember,
  UpstreamPool,
  UpstreamPoolBinding,
  UpstreamPoolMember,
  UpstreamPoolMemberSet,
  UpstreamPoolMemberSyncMode,
  UpstreamPoolMemberSyncResult,
  UpstreamCapacityPressure
} from '@/types'
import { extractApiErrorMessage } from '@/utils/apiError'
import {
  buildDeleteAccountSetConfirmMessage,
  buildDeleteBindingConfirmMessage,
  buildDeleteMemberSetConfirmMessage,
  buildDeletePoolConfirmMessage,
  filterAccountsForPoolCompletion,
  getLatestPoolHealthAlertStates,
  getPoolRoutingPaginationInfo
} from '@/utils/upstreamPoolInteractions'
import { useAppStore } from '@/stores'

type PoolForm = {
  name: string
  code: string
  platform: string
  description: string
  enabled: boolean
  scheduler_mode: string
  account_type_strategy: string
  sticky_enabled: boolean
  sticky_escape_enabled: boolean
  sticky_escape_error_rate_threshold: number
  sticky_escape_ttft_ms_threshold: number
  load_balance_enabled: boolean
  auto_weight_enabled: boolean
  auto_weight_mode: 'off' | 'observe' | 'active'
  failover_enabled: boolean
  top_k: number
  max_failover_hops: number
  wait_timeout_ms: number
  max_waiting: number
}

type MemberForm = {
  account_id: number | null
  enabled: boolean
  schedulable_override: boolean | null
  manual_drained: boolean
  weight: number
  priority_override: number | null
  max_concurrency_override: number | null
  notes: string
}

type SyncableAccount = {
  id: number
  name: string
  platform: string
  type: string
  status: Account['status']
  group_ids?: number[]
}

type BindingForm = {
  group_id: number | null
  pool_id: number | null
  platform: string
  modelsText: string
  request_path_scopeText: string
  priority: number
  enabled: boolean
}

type AccountSetForm = {
  name: string
  code: string
  platform: string
  description: string
  enabled: boolean
  shared_concurrency_limit: number | null
}

type CapacityDraft = { hard: number | null; soft: number | null }
type AccountSetMemberTypeFilter = 'all' | 'oauth' | 'apikey'

type MemberSetForm = {
  set_id: number | null
  enabled: boolean
  notes: string
}

type PoolRoutingObservabilityState = {
  loading: boolean
  loadingMore: boolean
  total: number
  page: number
  pageSize: number
  logs: OpsSystemLog[]
}

type PoolHealthAlertState = {
  loading: boolean
  error: string
  logs: OpsSystemLog[]
}

const appStore = useAppStore()
const loading = ref(false)
const submitting = ref(false)
const pools = ref<UpstreamPool[]>([])
const members = ref<UpstreamPoolMember[]>([])
const accountSets = ref<UpstreamAccountSet[]>([])
const accountSetMembers = ref<UpstreamAccountSetMember[]>([])
const memberSets = ref<UpstreamPoolMemberSet[]>([])
const bindings = ref<UpstreamPoolBinding[]>([])
const allGroups = ref<AdminGroup[]>([])
const syncableAccounts = ref<SyncableAccount[]>([])
const selectedPool = ref<UpstreamPool | null>(null)
const selectedAccountSet = ref<UpstreamAccountSet | null>(null)
const membersLoading = ref(false)
const accountSetsLoading = ref(false)
const accountSetMembersLoading = ref(false)
const memberSetsLoading = ref(false)
const bindingsLoading = ref(false)
const syncingMembers = ref(false)
const capacityPressures = ref<UpstreamCapacityPressure[]>([])
const capacityPressuresLoading = ref(false)
const capacityPressuresError = ref('')
const togglingPoolStatusIds = ref<Set<number>>(new Set())
const accountSetMembersPagination = ref({ page: 1, page_size: 10, total: 0 })
const poolRoutingObservability = ref<PoolRoutingObservabilityState>({
  loading: false,
  loadingMore: false,
  total: 0,
  page: 1,
  pageSize: 60,
  logs: [],
})
const poolHealthAlerts = ref<PoolHealthAlertState>({
  loading: false,
  error: '',
  logs: [],
})
const searchQuery = ref('')
const filters = ref({ platform: '', enabled: '' })
const pagination = ref({ page: 1, page_size: 20, total: 0 })
let poolObservabilityRequestToken = 0
let poolHealthAlertRequestToken = 0

const showPoolModal = ref(false)
const showMemberModal = ref(false)
const showAccountSetModal = ref(false)
const showAccountSetMemberPicker = ref(false)
const showMemberSetManagerModal = ref(false)
const showMemberSetModal = ref(false)
const showBindingModal = ref(false)
const showPoolObservabilityModal = ref(false)
const editingPool = ref<UpstreamPool | null>(null)
const editingMember = ref<UpstreamPoolMember | null>(null)
const editingAccountSet = ref<UpstreamAccountSet | null>(null)
const editingMemberSet = ref<UpstreamPoolMemberSet | null>(null)
const editingBinding = ref<UpstreamPoolBinding | null>(null)
const accountSetMemberSearch = ref('')
const accountSetMemberType = ref<AccountSetMemberTypeFilter>('all')
const selectedAccountSetMemberIDs = ref<number[]>([])

function createEmptyPoolRoutingObservabilityState(): PoolRoutingObservabilityState {
  return {
    loading: false,
    loadingMore: false,
    total: 0,
    page: 1,
    pageSize: 60,
    logs: [],
  }
}

function createEmptyPoolHealthAlertState(): PoolHealthAlertState {
  return { loading: false, error: '', logs: [] }
}

const poolRoutingPaginationInfo = computed(() => getPoolRoutingPaginationInfo({
  loading: poolRoutingObservability.value.loading,
  loadingMore: poolRoutingObservability.value.loadingMore,
  total: poolRoutingObservability.value.total,
  page: poolRoutingObservability.value.page,
  pageSize: poolRoutingObservability.value.pageSize,
  logsLength: poolRoutingObservability.value.logs.length,
}))

function openPoolObservabilityModal() {
  showPoolObservabilityModal.value = true
}

function closePoolObservabilityModal() {
  showPoolObservabilityModal.value = false
}

const poolForm = ref<PoolForm>({
  name: '',
  code: '',
  platform: 'openai',
  description: '',
  enabled: true,
  scheduler_mode: 'advanced',
  account_type_strategy: 'all',
  sticky_enabled: true,
  sticky_escape_enabled: true,
  sticky_escape_error_rate_threshold: 0.3,
  sticky_escape_ttft_ms_threshold: 6000,
  load_balance_enabled: true,
  auto_weight_enabled: false,
  auto_weight_mode: 'off',
  failover_enabled: true,
  top_k: 2,
  max_failover_hops: 3,
  wait_timeout_ms: 30000,
  max_waiting: 100,
})

const memberForm = ref<MemberForm>({
  account_id: null,
  enabled: true,
  schedulable_override: null,
  manual_drained: false,
  weight: 100,
  priority_override: null,
  max_concurrency_override: null,
  notes: '',
})

const accountSetForm = ref<AccountSetForm>({
  name: '',
  code: '',
  platform: 'openai',
  description: '',
  enabled: true,
  shared_concurrency_limit: null,
})
const capacityEditingAccountId = ref<number | null>(null)
const capacityDraft = ref<CapacityDraft>({ hard: null, soft: null })

const memberSetForm = ref<MemberSetForm>({
  set_id: null,
  enabled: true,
  notes: '',
})

const bindingForm = ref<BindingForm>({
  group_id: null,
  pool_id: null,
  platform: 'openai',
  modelsText: '',
  request_path_scopeText: '',
  priority: 0,
  enabled: true,
})

const platformOptions = [
  { value: 'openai', label: 'OpenAI' },
  { value: 'anthropic', label: 'Anthropic' },
  { value: 'gemini', label: 'Gemini' },
  { value: 'antigravity', label: 'Antigravity' },
]

const enabledOptions = [
  { value: '', label: '全部状态' },
  { value: 'true', label: '启用' },
  { value: 'false', label: '停用' },
]

const accountTypeStrategyOptions = [
  { value: 'all', label: '全部参与' },
  { value: 'oauth_preferred', label: 'OAuth 优先，API Key 兜底' },
  { value: 'oauth_only', label: '仅 OAuth' },
  { value: 'apikey_preferred', label: 'API Key 优先，OAuth 兜底' },
]

const autoWeightModeOptions = [
  { value: 'off', label: '关闭' },
  { value: 'observe', label: '观察模式' },
  { value: 'active', label: '正式应用' },
]

const boolSelectOptions = [
  { value: null, label: '不覆盖' },
  { value: true, label: '是' },
  { value: false, label: '否' },
]

const poolAccountSyncPlatforms = new Set(['openai', 'anthropic'])
const poolRoutingObservabilityPlatforms = new Set(['openai', 'anthropic'])
const canSyncSelectedPool = computed(() => poolAccountSyncPlatforms.has(selectedPool.value?.platform || ''))
const poolRoutingObservabilitySupported = computed(() => poolRoutingObservabilityPlatforms.has(selectedPool.value?.platform || ''))
const latestPoolHealthAlertStates = computed(() => getLatestPoolHealthAlertStates(poolHealthAlerts.value.logs))
const poolActiveHealthAlerts = computed(() => latestPoolHealthAlertStates.value.filter(log => {
  const status = getLogExtraString(log, 'alert_status')
  return status === 'firing' || status === 'reminder'
}))
const latestResolvedPoolHealthAlert = computed(() =>
  latestPoolHealthAlertStates.value.find(log => getLogExtraString(log, 'alert_status') === 'resolved') || null
)
const memberAccountNameMap = computed(() =>
  new Map(
    members.value.map(member => [
      member.account_id,
      formatPoolMemberAccountName(member.account_name, member.account_id),
    ])
  )
)
const stickyEscapeCount = computed(() =>
  poolRoutingObservability.value.logs.filter(log => getLogExtraBool(log, 'sticky_escape_triggered')).length
)
const stickyEscapeReasonSummary = computed(() => {
  const counts = new Map<string, number>()
  for (const log of poolRoutingObservability.value.logs) {
    if (!getLogExtraBool(log, 'sticky_escape_triggered')) continue
    const reason = formatStickyEscapeReason(getLogExtraString(log, 'sticky_escape_reason'))
    counts.set(reason, (counts.get(reason) || 0) + 1)
  }
  if (counts.size === 0) return '无'
  return Array.from(counts.entries())
    .sort((left, right) => right[1] - left[1])
    .map(([reason, count]) => `${reason} ${count}`)
    .join(' / ')
})
const recentRoutedAccountSummary = computed(() => {
  const labels: string[] = []
  const seen = new Set<number>()
  for (const log of poolRoutingObservability.value.logs) {
    const accountID = typeof log.account_id === 'number' ? log.account_id : null
    if (!accountID || seen.has(accountID)) continue
    seen.add(accountID)
    labels.push(formatObservabilityAccountLabel(accountID))
    if (labels.length >= 3) break
  }
  return labels.length > 0 ? labels.join(' / ') : '-'
})
const latestRoutingLogAt = computed(() => {
  const latest = poolRoutingObservability.value.logs[0]
  return latest ? formatCompactDateTime(latest.created_at) : '-'
})
const selectedPoolBindings = computed(() =>
  bindings.value.filter(binding => binding.pool_id === selectedPool.value?.id)
)
const selectedPoolEnabledBindings = computed(() =>
  selectedPoolBindings.value.filter(binding => binding.enabled)
)
const accountSetCountMap = computed(() =>
  new Map(accountSets.value.map(item => [item.id, item.account_count || 0]))
)
const emptyBoundMemberSets = computed(() =>
  memberSets.value.filter(item => item.enabled && (accountSetCountMap.value.get(item.set_id) || 0) === 0)
)
const enabledMemberSetCount = computed(() => memberSets.value.filter(item => item.enabled).length)
const poolHealth = computed(() => {
  const totalMembers = members.value.length
  const healthyMembers = members.value.filter(member => member.runtime_status === 'healthy').length
  const readyMembers = members.value.filter(isPoolMemberReady).length
  const drainedMembers = members.value.filter(member => member.manual_drained).length
  const disabledMembers = members.value.filter(member => !member.enabled || member.schedulable_override === false || member.account_schedulable === false).length
  const blockedMembers = members.value.filter(member =>
    ['rate_limited', 'overloaded', 'temp_unschedulable'].includes(String(member.runtime_status || ''))
  ).length
  const errorMembers = members.value.filter(member => member.runtime_status === 'error_recovering').length
  const logs = poolRoutingObservability.value.logs
  const routedAccountIds = new Set<number>()
  let stickyEscapes = 0
  for (const log of logs) {
    if (typeof log.account_id === 'number' && log.account_id > 0) {
      routedAccountIds.add(log.account_id)
    }
    if (getLogExtraBool(log, 'sticky_escape_triggered')) {
      stickyEscapes += 1
    }
  }
  return {
    totalMembers,
    healthyMembers,
    readyMembers,
    drainedMembers,
    disabledMembers,
    blockedMembers,
    errorMembers,
    observedRoutes: poolRoutingObservability.value.total || logs.length,
    uniqueRoutedAccounts: routedAccountIds.size,
    stickyEscapes,
  }
})
const poolHealthStatusLabel = computed(() => {
  if (!selectedPool.value) return '未选择'
  if (!selectedPool.value.enabled) return '已停用'
  if (poolHealth.value.totalMembers === 0) return '空池'
  if (poolHealth.value.readyMembers === 0) return '不可调度'
  if (poolHealth.value.blockedMembers > 0 || poolHealth.value.errorMembers > 0) return '有风险'
  return '健康'
})
const poolHealthStatusClass = computed(() => {
  switch (poolHealthStatusLabel.value) {
    case '健康':
      return 'badge-success'
    case '有风险':
      return 'badge-warning'
    case '不可调度':
    case '空池':
      return 'badge-danger'
    default:
      return 'badge-gray'
  }
})
const poolRoutedAccountRows = computed(() => {
  const counts = new Map<number, { accountId: number; label: string; count: number; latestAt: string }>()
  for (const log of poolRoutingObservability.value.logs) {
    const accountId = typeof log.account_id === 'number' ? log.account_id : 0
    if (!accountId) continue
    const current = counts.get(accountId)
    if (current) {
      current.count += 1
      continue
    }
    counts.set(accountId, {
      accountId,
      label: formatObservabilityAccountLabel(accountId),
      count: 1,
      latestAt: log.created_at,
    })
  }
  const max = Math.max(...Array.from(counts.values()).map(row => row.count), 1)
  return Array.from(counts.values())
    .sort((left, right) => right.count - left.count)
    .map(row => ({
      ...row,
      percent: Math.max(6, Math.round((row.count / max) * 100)),
    }))
})
const topRoutedAccountName = computed(() => {
  const top = poolRoutedAccountRows.value[0]
  return top ? top.label : '-'
})
const topRoutedAccountShare = computed(() => {
  const top = poolRoutedAccountRows.value[0]
  if (!top) return '-'
  const total = poolRoutingObservability.value.logs.length || top.count
  return formatPercent(top.count / total)
})
const topObservedModelSummary = computed(() => {
  const counts = new Map<string, number>()
  for (const log of poolRoutingObservability.value.logs) {
    const model = String(log.model || getLogExtraString(log, 'model') || '').trim()
    if (!model) continue
    counts.set(model, (counts.get(model) || 0) + 1)
  }
  const top = Array.from(counts.entries()).sort((left, right) => right[1] - left[1])[0]
  return top ? top[0] : '-'
})
const poolModelCoverageRows = computed(() => {
  if (!selectedPool.value) return []
  if (selectedPoolEnabledBindings.value.length === 0) return []
  const rows: Array<{ key: string; label: string; ready: boolean }> = []
  const seen = new Set<string>()
  for (const binding of selectedPoolEnabledBindings.value) {
    const models = (binding.models || []).map(model => model.trim()).filter(Boolean)
    if (models.length === 0) {
      if (!seen.has('*')) {
        rows.push({ key: '*', label: '全部模型', ready: poolHealth.value.readyMembers > 0 })
        seen.add('*')
      }
      continue
    }
    for (const model of models) {
      const key = model.toLowerCase()
      if (seen.has(key)) continue
      rows.push({ key, label: model, ready: poolHealth.value.readyMembers > 0 })
      seen.add(key)
    }
  }
  return rows
})
const poolHealthSignals = computed(() => {
  const signals: string[] = []
  if (!selectedPool.value) return signals
  if (!selectedPool.value.enabled) signals.push('池已停用')
  if (poolHealth.value.totalMembers === 0) signals.push('还没有成员账号')
  if (poolHealth.value.readyMembers === 0 && poolHealth.value.totalMembers > 0) signals.push('当前没有可调度成员')
  if (poolHealth.value.blockedMembers > 0) signals.push(`${poolHealth.value.blockedMembers} 个成员处于限流/过载/临时熔断`)
  if (poolHealth.value.errorMembers > 0) signals.push(`${poolHealth.value.errorMembers} 个成员错误待恢复`)
  if (selectedPoolEnabledBindings.value.length === 0) signals.push('没有启用中的分组绑定')
  const top = poolRoutedAccountRows.value[0]
  if (top && poolRoutingObservability.value.logs.length >= 10 && top.count / poolRoutingObservability.value.logs.length >= 0.8) {
    signals.push('最近调度集中在单个账号')
  }
  return signals
})

const poolColumns = [
  { key: 'name', label: '池' },
  { key: 'platform', label: '平台' },
  { key: 'enabled', label: '状态' },
  { key: 'routing', label: '路由' },
  { key: 'sticky_enabled', label: '粘性' },
  { key: 'members', label: '成员' },
  { key: 'bindings', label: '绑定' },
  { key: 'actions', label: '操作' },
]

const memberColumns = [
  { key: 'account_name', label: '账号', class: 'pool-member-account-col' },
  { key: 'enabled', label: '状态' },
  { key: 'runtime_status', label: '运行态' },
  { key: 'manual_drained', label: '排空' },
  { key: 'weight', label: '权重' },
  { key: 'updated_at', label: '更新时间' },
  { key: 'actions', label: '操作' },
]

const accountSetColumns = [
  { key: 'name', label: '集合' },
  { key: 'platform', label: '平台' },
  { key: 'enabled', label: '状态' },
  { key: 'account_count', label: '账号数' },
  { key: 'shared_concurrency_limit', label: '共享上限' },
  { key: 'updated_at', label: '更新时间' },
  { key: 'actions', label: '操作' },
]

const accountSetMemberColumns = [
  { key: 'account_name', label: '账号' },
  { key: 'account_platform', label: '平台' },
  { key: 'runtime_status', label: '运行态' },
  { key: 'usage', label: '使用情况' },
  { key: 'capacity', label: '容量策略' },
  { key: 'added_at', label: '加入时间' },
  { key: 'actions', label: '操作' },
]

const memberSetColumns = [
  { key: 'set_name', label: '集合' },
  { key: 'set_platform', label: '平台' },
  { key: 'account_count', label: '账号数' },
  { key: 'enabled', label: '状态' },
  { key: 'updated_at', label: '更新时间' },
  { key: 'actions', label: '操作' },
]

const bindingColumns = [
  { key: 'group_name', label: '分组' },
  { key: 'platform', label: '平台' },
  { key: 'enabled', label: '状态' },
  { key: 'priority', label: '优先级' },
  { key: 'actions', label: '操作' },
]

const platformLabel = (value: string) => platformOptions.find(item => item.value === value)?.label || value || '-'
const accountTypeStrategyLabel = (value?: string) => accountTypeStrategyOptions.find(item => item.value === value)?.label || '全部参与'

const runtimeStatusLabel = (value?: string) => {
  switch (value) {
    case 'healthy':
      return '正常'
    case 'error_recovering':
      return '错误待恢复'
    case 'rate_limited':
      return '限流中'
    case 'overloaded':
      return '过载保护'
    case 'temp_unschedulable':
      return '临时熔断'
    case 'disabled':
      return '账号停调度'
    case 'degraded':
      return '未调度'
    default:
      return '未知'
  }
}

const runtimeStatusClass = (value?: string) => {
  switch (value) {
    case 'healthy':
      return 'badge-success'
    case 'error_recovering':
      return 'badge-danger'
    case 'rate_limited':
    case 'overloaded':
    case 'temp_unschedulable':
      return 'badge-warning'
    default:
      return 'badge-gray'
  }
}

const formatCompactDateTime = (value?: string | null) => {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${month}/${day} ${hours}:${minutes}`
}

const formatRuntimeLastUsed = (value?: string | null) => {
  const formatted = formatCompactDateTime(value)
  return formatted === '-' ? '最近未使用' : `最近 ${formatted}`
}

const formatPoolMemberAccountName = (value?: string | null, accountID?: number) => {
  const text = String(value || `账号 #${accountID ?? '-'}`)
  const atIndex = text.indexOf('@')
  if (atIndex <= 0 || atIndex >= text.length - 1) {
    return text
  }
  return `${text.slice(0, atIndex + 1)}\u200b${text.slice(atIndex + 1)}`
}

const poolMemberRowKey = (row: UpstreamPoolMember) =>
  `${row.source_type || 'direct'}:${row.source_set_id || 0}:${row.account_id}:${row.id || 0}`

const accountSetMemberRowKey = (row: UpstreamAccountSetMember) => `${row.set_id}:${row.account_id}`

const getAccountSetMemberCount = (setID: number) => accountSetCountMap.value.get(setID) || 0

const formatThresholdMs = (value?: number | null) => {
  if (!value || value <= 0) return '-'
  return `${Math.round(value)}ms`
}

const formatRateThreshold = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value) || value < 0) return '-'
  return value.toFixed(2)
}

function capacityPressureClass(pressure: UpstreamCapacityPressure) {
  if (pressure.waiting_count > 0 || pressure.available_capacity <= 0) return 'badge-danger'
  const ratio = pressure.capacity_limit > 0 ? pressure.current_concurrency / pressure.capacity_limit : 0
  if (ratio >= 0.8) return 'badge-warning'
  return 'badge-success'
}

function capacityPressureLabel(pressure: UpstreamCapacityPressure) {
  if (pressure.waiting_count > 0) return '排队中'
  if (pressure.available_capacity <= 0) return '余量紧张'
  const ratio = pressure.capacity_limit > 0 ? pressure.current_concurrency / pressure.capacity_limit : 0
  if (ratio >= 0.8) return '需要观察'
  return '余量充足'
}

const isPoolStatusToggling = (poolID: number) => togglingPoolStatusIds.value.has(poolID)

function setPoolStatusToggling(poolID: number, toggling: boolean) {
  const next = new Set(togglingPoolStatusIds.value)
  if (toggling) {
    next.add(poolID)
  } else {
    next.delete(poolID)
  }
  togglingPoolStatusIds.value = next
}

const formatPercent = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value) || value < 0) return '-'
  return `${Math.round(value * 100)}%`
}

const isPoolMemberReady = (member: UpstreamPoolMember) => (
  member.enabled &&
  !member.manual_drained &&
  member.schedulable_override !== false &&
  member.account_schedulable !== false &&
  member.runtime_status === 'healthy'
)

const formatStickyEscapeReason = (value?: string | null) => {
  switch (String(value || '').trim()) {
    case 'ttft':
      return 'TTFT 过高'
    case 'error_rate':
      return '错率过高'
    case 'concurrency_full':
      return '并发已满'
    default:
      return String(value || '').trim() || '-'
  }
}

const formatAutoWeightReason = (value?: string | null) => {
  const labels: Record<string, string> = {
    rate_limited: '限流',
    overloaded: '过载',
    temporarily_unschedulable: '暂不可调度',
    probe_failed: '探测失败',
    probe_degraded: '探测降级',
    faster_than_pool: '快于同池',
    much_slower_than_pool: '明显慢于同池',
    slower_than_pool: '慢于同池',
    healthy: '健康',
  }
  const reason = String(value || '').trim()
  return labels[reason] || reason || '-'
}

const formatRoutingReason = (value?: string | null) => {
  switch (String(value || '').trim()) {
    case 'previous_response_sticky':
      return '上一响应粘性'
    case 'session_sticky':
      return '会话粘性'
    case 'session_sticky_wait':
      return '会话粘性等待'
    case 'load_balance':
      return '负载均衡'
    case 'load_balance_wait':
      return '负载均衡等待'
    case 'fallback_load_balance':
      return '回退负载均衡'
    case 'fallback_load_balance_wait':
      return '回退负载均衡等待'
    default:
      return String(value || '').trim() || '-'
  }
}

const getLogExtraString = (log: OpsSystemLog, key: string) => {
  const value = log.extra?.[key]
  if (value == null) return ''
  if (typeof value === 'string') return value.trim()
  if (typeof value === 'number' || typeof value === 'boolean') return String(value)
  return ''
}

const getLogExtraBool = (log: OpsSystemLog, key: string) => {
  const value = log.extra?.[key]
  if (typeof value === 'boolean') return value
  if (typeof value === 'string') return value.trim().toLowerCase() === 'true'
  return false
}

const getLogExtraNumber = (log: OpsSystemLog, key: string) => {
  const value = log.extra?.[key]
  if (typeof value === 'number' && Number.isFinite(value)) return value
  if (typeof value === 'string') {
    const parsed = Number(value)
    if (Number.isFinite(parsed)) return parsed
  }
  return null
}

const poolHealthAlertTypeLabels: Record<string, string> = {
  pool_capacity_low: '池可用账号不足',
  pool_unavailable: '池无可用账号',
  account_rate_limited: '账号持续限流',
  account_error_rate_high: '账号错误率偏高',
  account_probe_failed: '账号连续探测失败',
  account_latency_degraded: '账号延迟恶化',
  account_runtime_weight_low: '自动权重持续偏低',
}

const formatPoolHealthAlertTitle = (log: OpsSystemLog) => {
  const type = getLogExtraString(log, 'alert_type')
  return poolHealthAlertTypeLabels[type] || log.message || '上游异常'
}

const poolHealthAlertCardClass = (log: OpsSystemLog) => (
  getLogExtraString(log, 'severity') === 'critical'
    ? 'border-rose-200 bg-rose-50/70 dark:border-rose-900/40 dark:bg-rose-950/20'
    : 'border-amber-200 bg-amber-50/70 dark:border-amber-900/40 dark:bg-amber-950/20'
)

const formatPoolHealthAlertDetail = (log: OpsSystemLog) => {
  const parts: string[] = []
  const accountName = getLogExtraString(log, 'account_name')
  const reason = getLogExtraString(log, 'reason')
  const available = getLogExtraNumber(log, 'available_members')
  const total = getLogExtraNumber(log, 'total_members')
  const errorRate = getLogExtraNumber(log, 'error_rate')
  const latency = getLogExtraNumber(log, 'latency_ms')
  const medianLatency = getLogExtraNumber(log, 'pool_median_latency_ms')
  const factor = getLogExtraNumber(log, 'runtime_weight_factor')

  if (accountName) parts.push(`账号 ${accountName}`)
  if (available != null && total != null) parts.push(`可用 ${available}/${total}`)
  if (errorRate != null) parts.push(`错误率 ${(errorRate * 100).toFixed(1)}%`)
  if (latency != null) parts.push(`延迟 ${latency}ms${medianLatency != null ? `，池中位数 ${medianLatency}ms` : ''}`)
  if (factor != null) parts.push(`运行因子 ${factor.toFixed(2)}`)
  if (parts.length === 0 && reason) parts.push(`原因 ${formatAutoWeightReason(reason)}`)
  return parts.join(' · ') || log.message || '请检查当前池成员状态'
}

const getLogExtraNumberArray = (log: OpsSystemLog, key: string) => {
  const value = log.extra?.[key]
  if (!Array.isArray(value)) return []
  return value
    .map(item => {
      if (typeof item === 'number' && Number.isFinite(item)) return item
      if (typeof item === 'string') {
        const parsed = Number(item)
        return Number.isFinite(parsed) ? parsed : null
      }
      return null
    })
    .filter((item): item is number => item != null)
}

const getLogExtraRecord = (log: OpsSystemLog, key: string) => {
  const value = log.extra?.[key]
  if (!value || Array.isArray(value) || typeof value !== 'object') return null
  return value as Record<string, unknown>
}

const getObservabilitySelectedAccountId = (log: OpsSystemLog) => {
  const extraSelected = getLogExtraNumber(log, 'selected_account_id')
  if (extraSelected != null && extraSelected > 0) return extraSelected
  if (typeof log.account_id === 'number' && log.account_id > 0) return log.account_id
  return null
}

const getObservabilityCandidateCount = (log: OpsSystemLog) => {
  const value = getLogExtraNumber(log, 'candidate_count')
  return value != null && value >= 0 ? value : null
}

const getObservabilityTopKLabels = (log: OpsSystemLog) =>
  getLogExtraNumberArray(log, 'cache_affinity_top_k_account_ids').map(accountId =>
    formatObservabilityAccountLabel(accountId)
  )

const formatRoutingSkippedReason = (value: string) => {
  switch (value) {
    case 'excluded_by_failover':
      return '故障回退排除'
    case 'platform_mismatch':
      return '平台不匹配'
    case 'not_schedulable':
      return '账号不可调度'
    case 'runtime_blocked':
      return '运行态阻塞'
    case 'privacy_not_set':
      return '隐私集未设置'
    case 'request_incompatible':
      return '请求不兼容'
    case 'transport_incompatible':
      return '传输不兼容'
    default:
      return value || '未知原因'
  }
}

const getObservabilitySkippedEntries = (log: OpsSystemLog) => {
  const record = getLogExtraRecord(log, 'skipped')
  if (!record) return []
  return Object.entries(record)
    .map(([key, value]) => {
      const count = typeof value === 'number'
        ? value
        : typeof value === 'string'
          ? Number(value)
          : NaN
      if (!Number.isFinite(count) || count <= 0) return null
      return {
        key,
        count,
        label: formatRoutingSkippedReason(key),
      }
    })
    .filter((entry): entry is { key: string; count: number; label: string } => entry != null)
    .sort((left, right) => right.count - left.count)
}

const formatObservabilityAccountLabel = (accountID?: number | null) => {
  if (!accountID) return '-'
  return memberAccountNameMap.value.get(accountID) || `账号 #${accountID}`
}

const filteredPools = computed(() => {
  const keyword = searchQuery.value.trim().toLowerCase()
  return pools.value.filter((pool) => {
    if (filters.value.platform && pool.platform !== filters.value.platform) return false
    if (filters.value.enabled && String(pool.enabled) !== filters.value.enabled) return false
    if (!keyword) return true
    return [pool.name, pool.code, pool.description, pool.platform].some((value) =>
      String(value || '').toLowerCase().includes(keyword)
    )
  })
})

const pagedPools = computed(() => {
  const start = (pagination.value.page - 1) * pagination.value.page_size
  return filteredPools.value.slice(start, start + pagination.value.page_size)
})

const pagedAccountSetMembers = computed(() => {
  const start = (accountSetMembersPagination.value.page - 1) * accountSetMembersPagination.value.page_size
  return accountSetMembers.value.slice(start, start + accountSetMembersPagination.value.page_size)
})

const directMembers = computed(() =>
  members.value.filter(member => String(member.source_type || 'direct') !== 'account_set')
)

const availablePoolAccounts = computed(() => {
  return filterAccountsForPoolCompletion({
    pool: selectedPool.value,
    accounts: syncableAccounts.value,
    bindings: bindings.value,
  }) as SyncableAccount[]
})

const accountSetAvailableAccounts = computed(() => {
  const platform = selectedAccountSet.value?.platform
  if (!platform) return []
  return syncableAccounts.value.filter(account => account.platform === platform)
})

const accountSetSelectableAccounts = computed(() => {
  const existing = new Set(accountSetMembers.value.map(item => item.account_id))
  return accountSetAvailableAccounts.value.filter(account => !existing.has(account.id))
})

const accountSetMemberTypeOptions: Array<{ value: AccountSetMemberTypeFilter; label: string }> = [
  { value: 'all', label: '全部' },
  { value: 'oauth', label: 'OAuth' },
  { value: 'apikey', label: 'API Key' },
]

const filteredAccountSetCandidates = computed(() => {
  const query = accountSetMemberSearch.value.trim().toLowerCase()
  return accountSetSelectableAccounts.value.filter((account) => {
    const typeMatches = accountSetMemberType.value === 'all'
      || (accountSetMemberType.value === 'oauth' && (account.type === 'oauth' || account.type === 'setup-token'))
      || (accountSetMemberType.value === 'apikey' && account.type === 'apikey')
    if (!typeMatches) return false
    if (!query) return true
    return account.name.toLowerCase().includes(query) || String(account.id).includes(query)
  })
})

const allFilteredAccountSetCandidatesSelected = computed(() => {
  if (filteredAccountSetCandidates.value.length === 0) return false
  const selected = new Set(selectedAccountSetMemberIDs.value)
  return filteredAccountSetCandidates.value.every(account => selected.has(account.id))
})

const accountOptions = computed(() =>
  availablePoolAccounts.value.map((account) => ({
    value: account.id,
    label: `${account.name} · ${platformLabel(account.platform)} · ${account.type === 'apikey' ? 'API Key' : account.type}`,
  }))
)

const accountSetOptionsForPool = computed(() => {
  const platform = selectedPool.value?.platform
  return accountSets.value
    .filter((item) => !platform || item.platform === platform)
    .map((item) => ({
      value: item.id,
      label: `${item.name} · ${platformLabel(item.platform)} · ${item.account_count || 0} 个账号`,
    }))
})

const editingPoolPlatformLocked = computed(() => {
  if (!editingPool.value) return false
  return directMembers.value.length > 0 || memberSets.value.length > 0 || selectedPoolBindings.value.length > 0
})

const editingAccountSetPlatformLocked = computed(() => {
  if (!editingAccountSet.value) return false
  return getAccountSetMemberCount(editingAccountSet.value.id) > 0 || memberSets.value.some(item => item.set_id === editingAccountSet.value?.id)
})

const groupOptionsForBinding = computed(() => {
  const pool = pools.value.find(item => item.id === bindingForm.value.pool_id)
  const platform = pool?.platform || bindingForm.value.platform
  return allGroups.value
    .filter(group => !platform || group.platform === platform)
    .map((group) => ({
      value: group.id,
      label: `${group.name} · ${platformLabel(group.platform)}`,
    }))
})

const poolOptions = computed(() =>
  pools.value.map((pool) => ({
    value: pool.id,
    label: `${pool.name} · ${platformLabel(pool.platform)}`,
  }))
)

function resetPoolForm() {
  poolForm.value = {
    name: '',
    code: '',
    platform: 'openai',
    description: '',
    enabled: true,
    scheduler_mode: 'advanced',
    account_type_strategy: 'all',
    sticky_enabled: true,
    sticky_escape_enabled: true,
    sticky_escape_error_rate_threshold: 0.3,
    sticky_escape_ttft_ms_threshold: 6000,
    load_balance_enabled: true,
    auto_weight_enabled: false,
    auto_weight_mode: 'off',
    failover_enabled: true,
    top_k: 2,
    max_failover_hops: 3,
    wait_timeout_ms: 30000,
    max_waiting: 100,
  }
}

function resetMemberForm() {
  memberForm.value = {
    account_id: null,
    enabled: true,
    schedulable_override: null,
    manual_drained: false,
    weight: 100,
    priority_override: null,
    max_concurrency_override: null,
    notes: '',
  }
}

function resetBindingForm() {
  bindingForm.value = {
    group_id: null,
    pool_id: selectedPool.value?.id ?? null,
    platform: selectedPool.value?.platform ?? 'openai',
    modelsText: '',
    request_path_scopeText: '',
    priority: 0,
    enabled: true,
  }
}

function resetAccountSetForm() {
  accountSetForm.value = {
    name: '',
    code: '',
    platform: selectedPool.value?.platform ?? 'openai',
    description: '',
    enabled: true,
    shared_concurrency_limit: null,
  }
}

function resetMemberSetForm() {
  memberSetForm.value = {
    set_id: null,
    enabled: true,
    notes: '',
  }
}

async function loadAll() {
  loading.value = true
  accountSetsLoading.value = true
  try {
    const [poolList, groupList, accountList, bindingList, accountSetList] = await Promise.all([
      adminAPI.upstreamPools.list(),
      adminAPI.groups.getAllIncludingInactive(),
      loadSyncableAccounts(),
      adminAPI.upstreamPools.getBindings(),
      adminAPI.upstreamPools.getAccountSets(),
    ])
    pools.value = poolList
    allGroups.value = groupList
    syncableAccounts.value = accountList
    bindings.value = bindingList
    accountSets.value = accountSetList
    pagination.value.total = filteredPools.value.length
    if (poolList.length === 0) {
      selectedPool.value = null
      members.value = []
      memberSets.value = []
      poolRoutingObservability.value = createEmptyPoolRoutingObservabilityState()
      poolHealthAlerts.value = createEmptyPoolHealthAlertState()
    } else if (selectedPool.value) {
      const nextSelected = poolList.find(pool => pool.id === selectedPool.value?.id)
      selectPool(nextSelected || poolList[0])
    } else {
      selectPool(poolList[0])
    }
    if (accountSetList.length === 0) {
      selectedAccountSet.value = null
      accountSetMembers.value = []
    } else if (selectedAccountSet.value) {
      const nextSelectedSet = accountSetList.find(item => item.id === selectedAccountSet.value?.id)
      selectAccountSet(nextSelectedSet || accountSetList[0])
    } else {
      selectAccountSet(accountSetList[0])
    }
    await loadCapacityPressures()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载上游池失败'))
  } finally {
    loading.value = false
    accountSetsLoading.value = false
  }
}

async function loadCapacityPressures() {
  capacityPressuresLoading.value = true
  capacityPressuresError.value = ''
  try {
    capacityPressures.value = await adminAPI.upstreamPools.getCapacityPressures()
  } catch (error) {
    capacityPressures.value = []
    capacityPressuresError.value = extractApiErrorMessage(error, '加载共享容量压力失败')
  } finally {
    capacityPressuresLoading.value = false
  }
}

async function loadMembers(poolId: number) {
  membersLoading.value = true
  try {
    members.value = await adminAPI.upstreamPools.getMembers(poolId)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载成员失败'))
  } finally {
    membersLoading.value = false
  }
}

async function loadMemberSets(poolId: number) {
  memberSetsLoading.value = true
  try {
    memberSets.value = await adminAPI.upstreamPools.getMemberSets(poolId)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载集合绑定失败'))
  } finally {
    memberSetsLoading.value = false
  }
}

async function loadAccountSetMembers(setId: number) {
  accountSetMembersLoading.value = true
  try {
    accountSetMembers.value = await adminAPI.upstreamPools.getAccountSetMembers(setId)
    accountSetMembersPagination.value.total = accountSetMembers.value.length
    const totalPages = Math.max(1, Math.ceil(accountSetMembers.value.length / accountSetMembersPagination.value.page_size))
    if (accountSetMembersPagination.value.page > totalPages) {
      accountSetMembersPagination.value.page = totalPages
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载集合成员失败'))
  } finally {
    accountSetMembersLoading.value = false
  }
}

async function loadBindings() {
  bindingsLoading.value = true
  try {
    bindings.value = await adminAPI.upstreamPools.getBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载绑定失败'))
  } finally {
    bindingsLoading.value = false
  }
}

function selectPool(pool: UpstreamPool) {
  selectedPool.value = pool
  bindingForm.value.pool_id = pool.id
  bindingForm.value.platform = pool.platform
  loadMembers(pool.id)
  loadMemberSets(pool.id)
  loadPoolObservability(pool)
  loadPoolHealthAlerts(pool)
  pagination.value.total = filteredPools.value.length
}

function selectAccountSet(item: UpstreamAccountSet) {
  selectedAccountSet.value = item
  accountSetMembersPagination.value.page = 1
  loadAccountSetMembers(item.id)
}

async function loadPoolObservability(pool: UpstreamPool) {
  if (!poolRoutingObservabilityPlatforms.has(pool.platform)) {
    poolRoutingObservability.value = createEmptyPoolRoutingObservabilityState()
    return
  }

  const requestToken = ++poolObservabilityRequestToken
  poolRoutingObservability.value = {
    loading: true,
    loadingMore: false,
    total: 0,
    page: 1,
    pageSize: poolRoutingObservability.value.pageSize,
    logs: [],
  }

  try {
    const result = await adminAPI.ops.listSystemLogs({
      page: 1,
      page_size: poolRoutingObservability.value.pageSize,
      time_range: '24h',
      component: 'routing.explanation',
      platform: pool.platform,
      pool_id: pool.id,
    })
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = {
      loading: false,
      loadingMore: false,
      total: result.total || 0,
      page: result.page || 1,
      pageSize: result.page_size || poolRoutingObservability.value.pageSize,
      logs: result.items || [],
    }
  } catch (error) {
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = createEmptyPoolRoutingObservabilityState()
    appStore.showError(extractApiErrorMessage(error, '加载最近路由观测失败'))
  }
}

async function loadPoolHealthAlerts(pool: UpstreamPool) {
  const requestToken = ++poolHealthAlertRequestToken
  poolHealthAlerts.value = { loading: true, error: '', logs: [] }
  try {
    const result = await adminAPI.ops.listSystemLogs({
      page: 1,
      page_size: 50,
      time_range: '7d',
      component: 'upstream.health_alert',
      platform: pool.platform,
      pool_id: pool.id,
    })
    if (requestToken !== poolHealthAlertRequestToken) return
    poolHealthAlerts.value = { loading: false, error: '', logs: result.items || [] }
  } catch (error) {
    if (requestToken !== poolHealthAlertRequestToken) return
    poolHealthAlerts.value = {
      loading: false,
      error: extractApiErrorMessage(error, '加载最近预警失败'),
      logs: [],
    }
  }
}

async function loadMorePoolObservability() {
  const pool = selectedPool.value
  if (!pool || !poolRoutingObservabilitySupported.value || !poolRoutingPaginationInfo.value.hasMore) return

  const requestToken = ++poolObservabilityRequestToken
  poolRoutingObservability.value = {
    ...poolRoutingObservability.value,
    loadingMore: true,
  }

  try {
    const result = await adminAPI.ops.listSystemLogs({
      page: poolRoutingPaginationInfo.value.nextPage,
      page_size: poolRoutingObservability.value.pageSize,
      time_range: '24h',
      component: 'routing.explanation',
      platform: pool.platform,
      pool_id: pool.id,
    })
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = {
      loading: false,
      loadingMore: false,
      total: result.total || poolRoutingObservability.value.total,
      page: result.page || poolRoutingPaginationInfo.value.nextPage,
      pageSize: result.page_size || poolRoutingObservability.value.pageSize,
      logs: [
        ...poolRoutingObservability.value.logs,
        ...(result.items || []),
      ],
    }
  } catch (error) {
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = {
      ...poolRoutingObservability.value,
      loadingMore: false,
    }
    appStore.showError(extractApiErrorMessage(error, '加载更多路由观测失败'))
  }
}

async function refreshSelectedPoolHealth() {
  if (!selectedPool.value) return
  await Promise.all([
    loadMembers(selectedPool.value.id),
    loadMemberSets(selectedPool.value.id),
    loadPoolObservability(selectedPool.value),
    loadPoolHealthAlerts(selectedPool.value),
  ])
}

function openPoolModal(pool?: UpstreamPool | null) {
  editingPool.value = pool ?? null
  if (pool) {
    poolForm.value = {
      name: pool.name,
      code: pool.code,
      platform: pool.platform,
      description: pool.description || '',
      enabled: pool.enabled,
      scheduler_mode: pool.scheduler_mode || 'advanced',
      account_type_strategy: pool.account_type_strategy || 'all',
      sticky_enabled: pool.sticky_enabled,
      sticky_escape_enabled: pool.sticky_escape_enabled,
      sticky_escape_error_rate_threshold: pool.sticky_escape_error_rate_threshold,
      sticky_escape_ttft_ms_threshold: pool.sticky_escape_ttft_ms_threshold,
      load_balance_enabled: pool.load_balance_enabled,
      auto_weight_enabled: Boolean(pool.auto_weight_enabled),
      auto_weight_mode: pool.auto_weight_mode || (pool.auto_weight_enabled ? 'active' : 'off'),
      failover_enabled: pool.failover_enabled,
      top_k: pool.top_k,
      max_failover_hops: pool.max_failover_hops,
      wait_timeout_ms: pool.wait_timeout_ms,
      max_waiting: pool.max_waiting,
    }
  } else {
    resetPoolForm()
  }
  showPoolModal.value = true
}

function closePoolModal() {
  showPoolModal.value = false
  editingPool.value = null
}

function openMemberModal(member?: UpstreamPoolMember | null) {
  if (!selectedPool.value) return
  if (member && member.editable === false) {
    appStore.showError('集合展开成员不能直接编辑，请编辑集合绑定或直接成员配置')
    return
  }
  editingMember.value = member ?? null
  if (member) {
    memberForm.value = {
      account_id: member.account_id,
      enabled: member.enabled,
      schedulable_override: member.schedulable_override,
      manual_drained: member.manual_drained,
      weight: member.weight,
      priority_override: member.priority_override,
      max_concurrency_override: member.max_concurrency_override,
      notes: member.notes || '',
    }
  } else {
    resetMemberForm()
  }
  showMemberModal.value = true
}

function closeMemberModal() {
  showMemberModal.value = false
  editingMember.value = null
}

function openAccountSetModal(item?: UpstreamAccountSet | null) {
  editingAccountSet.value = item ?? null
  if (item) {
    accountSetForm.value = {
      name: item.name,
      code: item.code,
      platform: item.platform,
      description: item.description || '',
      enabled: item.enabled,
      shared_concurrency_limit: item.shared_concurrency_limit,
    }
  } else {
    resetAccountSetForm()
  }
  showAccountSetModal.value = true
}

function closeAccountSetModal() {
  showAccountSetModal.value = false
  editingAccountSet.value = null
}

function openMemberSetManagerModal() {
  if (!selectedPool.value) return
  showMemberSetManagerModal.value = true
}

function closeMemberSetManagerModal() {
  showMemberSetManagerModal.value = false
}

function openMemberSetEditor(item?: UpstreamPoolMemberSet | null) {
  showMemberSetManagerModal.value = false
  openMemberSetModal(item)
}

function openMemberSetModal(item?: UpstreamPoolMemberSet | null) {
  if (!selectedPool.value) return
  editingMemberSet.value = item ?? null
  if (item) {
    memberSetForm.value = {
      set_id: item.set_id,
      enabled: item.enabled,
      notes: item.notes || '',
    }
  } else {
    resetMemberSetForm()
    memberSetForm.value.set_id = selectedAccountSet.value?.id ?? null
  }
  showMemberSetModal.value = true
}

function closeMemberSetModal() {
  showMemberSetModal.value = false
  editingMemberSet.value = null
  if (selectedPool.value) {
    showMemberSetManagerModal.value = true
  }
}

function openBindingModal(binding?: UpstreamPoolBinding | null) {
  editingBinding.value = binding ?? null
  if (binding) {
    bindingForm.value = {
      group_id: binding.group_id,
      pool_id: binding.pool_id,
      platform: binding.platform,
      modelsText: (binding.models || []).join(', '),
      request_path_scopeText: (binding.request_path_scope || []).join(', '),
      priority: binding.priority,
      enabled: binding.enabled,
    }
  } else {
    resetBindingForm()
  }
  showBindingModal.value = true
}

function handleBindingPoolChange(value: string | number | boolean | null) {
  const poolID = typeof value === 'number' ? value : Number(value)
  const pool = pools.value.find(item => item.id === poolID)
  if (!pool) return
  bindingForm.value.platform = pool.platform
  const selectedGroup = allGroups.value.find(group => group.id === bindingForm.value.group_id)
  if (selectedGroup && selectedGroup.platform !== pool.platform) {
    bindingForm.value.group_id = null
  }
}

function closeBindingModal() {
  showBindingModal.value = false
  editingBinding.value = null
}

async function submitPool() {
  submitting.value = true
  try {
    const payload = {
      ...poolForm.value,
      description: poolForm.value.description || undefined,
      auto_weight_enabled: poolForm.value.platform === 'openai' && poolForm.value.auto_weight_enabled,
      auto_weight_mode: poolForm.value.platform === 'openai' ? poolForm.value.auto_weight_mode : 'off',
    }
    if (editingPool.value) {
      await adminAPI.upstreamPools.update(editingPool.value.id, payload)
      appStore.showSuccess('上游池已更新')
    } else {
      await adminAPI.upstreamPools.create(payload)
      appStore.showSuccess('上游池已创建')
    }
    closePoolModal()
    await loadAll()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存上游池失败'))
  } finally {
    submitting.value = false
  }
}

async function togglePoolEnabled(pool: UpstreamPool) {
  if (isPoolStatusToggling(pool.id)) return

  const previousEnabled = pool.enabled
  const nextEnabled = !previousEnabled
  const replacePool = (nextPool: UpstreamPool) => {
    pools.value = pools.value.map(item => (item.id === nextPool.id ? nextPool : item))
    if (selectedPool.value?.id === nextPool.id) {
      selectedPool.value = nextPool
    }
    if (editingPool.value?.id === nextPool.id) {
      editingPool.value = nextPool
      poolForm.value.enabled = nextPool.enabled
    }
  }

  setPoolStatusToggling(pool.id, true)
  replacePool({ ...pool, enabled: nextEnabled })
  try {
    const updatedPool = await adminAPI.upstreamPools.update(pool.id, { enabled: nextEnabled })
    replacePool(updatedPool)
    appStore.showSuccess(nextEnabled ? '上游池已启用' : '上游池已关闭')
  } catch (error) {
    replacePool({ ...pool, enabled: previousEnabled })
    appStore.showError(extractApiErrorMessage(error, nextEnabled ? '启用上游池失败' : '关闭上游池失败'))
  } finally {
    setPoolStatusToggling(pool.id, false)
  }
}

async function submitMember() {
  if (!selectedPool.value || !memberForm.value.account_id) {
    appStore.showError('请选择账号')
    return
  }
  submitting.value = true
  try {
    const payload = {
      account_id: memberForm.value.account_id,
      enabled: memberForm.value.enabled,
      schedulable_override: memberForm.value.schedulable_override,
      manual_drained: memberForm.value.manual_drained,
      weight: memberForm.value.weight,
      priority_override: memberForm.value.priority_override,
      max_concurrency_override: memberForm.value.max_concurrency_override,
      notes: memberForm.value.notes || null,
    }
    if (editingMember.value) {
      await adminAPI.upstreamPools.updateMember(editingMember.value.id, payload)
      appStore.showSuccess('成员已更新')
    } else {
      await adminAPI.upstreamPools.createMember(selectedPool.value.id, payload)
      appStore.showSuccess('成员已添加')
    }
    closeMemberModal()
    await loadMembers(selectedPool.value.id)
    await loadBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存成员失败'))
  } finally {
    submitting.value = false
  }
}

async function submitAccountSet() {
  submitting.value = true
  try {
    const payload = {
      name: accountSetForm.value.name,
      platform: accountSetForm.value.platform,
      enabled: accountSetForm.value.enabled,
      description: accountSetForm.value.description || undefined,
      shared_concurrency_limit: accountSetForm.value.shared_concurrency_limit || null,
      ...(editingAccountSet.value ? { code: accountSetForm.value.code } : {}),
    }
    if (editingAccountSet.value) {
      await adminAPI.upstreamPools.updateAccountSet(editingAccountSet.value.id, payload)
      appStore.showSuccess('账号集合已更新')
    } else {
      await adminAPI.upstreamPools.createAccountSet(payload)
      appStore.showSuccess('账号集合已创建')
    }
    closeAccountSetModal()
    await loadAll()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存账号集合失败'))
  } finally {
    submitting.value = false
  }
}

async function submitMemberSet() {
  if (!selectedPool.value || !memberSetForm.value.set_id) {
    appStore.showError('请选择账号集合')
    return
  }
  submitting.value = true
  try {
    const payload = {
      set_id: memberSetForm.value.set_id,
      enabled: memberSetForm.value.enabled,
      notes: memberSetForm.value.notes || null,
    }
    if (editingMemberSet.value) {
      await adminAPI.upstreamPools.updateMemberSet(editingMemberSet.value.id, {
        enabled: payload.enabled,
        notes: payload.notes,
      })
      appStore.showSuccess('集合绑定已更新')
    } else {
      await adminAPI.upstreamPools.createMemberSet(selectedPool.value.id, payload)
      appStore.showSuccess('集合已绑定到当前池')
    }
    closeMemberSetModal()
    await loadMemberSets(selectedPool.value.id)
    await loadMembers(selectedPool.value.id)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存集合绑定失败'))
  } finally {
    submitting.value = false
  }
}

async function syncSelectedPoolMembers() {
  if (!selectedPool.value) return
  if (!canSyncSelectedPool.value) {
    appStore.showError('当前只支持同步 OpenAI / Anthropic 上游池')
    return
  }

  syncingMembers.value = true
  try {
    const mode: UpstreamPoolMemberSyncMode = 'membership_only'
    const preview = await adminAPI.upstreamPools.previewMemberSync(selectedPool.value.id, { mode })
    if (!confirmMemberSyncPreview(preview)) return
    const result = await adminAPI.upstreamPools.applyMemberSync(selectedPool.value.id, { mode })
    appStore.showSuccess(`同步完成：新增 ${result.create_count}，保留 ${result.skip_count}，移除 ${result.delete_count}`)
    await loadMembers(selectedPool.value.id)
    await loadBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '同步账号失败'))
  } finally {
    syncingMembers.value = false
  }
}

function confirmMemberSyncPreview(preview: UpstreamPoolMemberSyncResult) {
  if (preview.create_count === 0 && preview.update_count === 0 && preview.delete_count === 0) {
    appStore.showSuccess(`无需同步：已保留 ${preview.skip_count} 个直接成员`)
    return false
  }
  const risk = preview.overwrite_risk_count > 0 ? `\n会影响 ${preview.overwrite_risk_count} 个带人工配置的成员。` : ''
  return window.confirm(
    `将按后端预演结果事务同步直接成员：\n新增 ${preview.create_count}\n更新 ${preview.update_count}\n移除 ${preview.delete_count}\n保留 ${preview.skip_count}${risk}\n\n继续执行吗？`
  )
}

async function addMissingSelectedPoolMembers() {
  if (!selectedPool.value) return
  if (!canSyncSelectedPool.value) {
    appStore.showError('当前只支持补齐 OpenAI / Anthropic 上游池')
    return
  }
  if (availablePoolAccounts.value.length === 0) {
    appStore.showError('当前平台没有可添加的账号')
    return
  }

  syncingMembers.value = true
  try {
    const poolMembers = await adminAPI.upstreamPools.getMembers(selectedPool.value.id)
    const memberAccountIds = new Set(poolMembers.map(member => member.account_id))
    let created = 0

    for (const account of availablePoolAccounts.value) {
      if (memberAccountIds.has(account.id)) continue
      const enabled = account.status === 'active'
      await adminAPI.upstreamPools.createMember(selectedPool.value.id, {
        account_id: account.id,
        enabled,
        schedulable_override: null,
        manual_drained: !enabled,
        weight: 100,
        priority_override: null,
        max_concurrency_override: null,
        notes: '补齐缺失账号',
      })
      created += 1
    }

    appStore.showSuccess(`补齐完成：新增 ${created}`)
    await loadMembers(selectedPool.value.id)
    await loadBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '补齐账号失败'))
  } finally {
    syncingMembers.value = false
  }
}

function accountSetAccountTypeLabel(type: string) {
  if (type === 'apikey') return 'API Key'
  if (type === 'setup-token') return 'Setup Token'
  if (type === 'oauth') return 'OAuth'
  return type || '-'
}

function openAccountSetMemberPicker() {
  accountSetMemberSearch.value = ''
  accountSetMemberType.value = 'all'
  selectedAccountSetMemberIDs.value = []
  showAccountSetMemberPicker.value = true
}

function closeAccountSetMemberPicker() {
  showAccountSetMemberPicker.value = false
  selectedAccountSetMemberIDs.value = []
}

function toggleAllFilteredAccountSetCandidates() {
  const visibleIDs = filteredAccountSetCandidates.value.map(account => account.id)
  const selected = new Set(selectedAccountSetMemberIDs.value)
  if (allFilteredAccountSetCandidatesSelected.value) {
    visibleIDs.forEach(id => selected.delete(id))
  } else {
    visibleIDs.forEach(id => selected.add(id))
  }
  selectedAccountSetMemberIDs.value = Array.from(selected)
}

async function addAccountIDsToSelectedSet(targetIDs: number[], successMessage: string): Promise<boolean> {
  if (!selectedAccountSet.value) return false
  if (targetIDs.length === 0) {
    appStore.showError('当前没有可加入的新账号')
    return false
  }
  submitting.value = true
  try {
    await adminAPI.upstreamPools.addAccountSetMembers(selectedAccountSet.value.id, {
      account_ids: targetIDs,
    })
    appStore.showSuccess(successMessage)
    await Promise.all([
      loadAccountSetMembers(selectedAccountSet.value.id),
      loadAll(),
    ])
    return true
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加入账号集合失败'))
    return false
  } finally {
    submitting.value = false
  }
}

async function addSelectedAccountsToSet() {
  const targetIDs = [...selectedAccountSetMemberIDs.value]
  const added = await addAccountIDsToSelectedSet(targetIDs, `已加入 ${targetIDs.length} 个账号`)
  if (added) closeAccountSetMemberPicker()
}

async function addAccountsToSelectedSet() {
  const targetIDs = accountSetSelectableAccounts.value.map(account => account.id)
  await addAccountIDsToSelectedSet(targetIDs, `已补齐 ${targetIDs.length} 个账号`)
}

async function submitBinding() {
  if (!bindingForm.value.group_id || !bindingForm.value.pool_id) {
    appStore.showError('请选择分组和上游池')
    return
  }
  submitting.value = true
  try {
    const payload = {
      group_id: bindingForm.value.group_id ?? 0,
      pool_id: bindingForm.value.pool_id ?? 0,
      platform: bindingForm.value.platform,
      models: bindingForm.value.modelsText.split(',').map(v => v.trim()).filter(Boolean),
      request_path_scope: bindingForm.value.request_path_scopeText.split(',').map(v => v.trim()).filter(Boolean),
      priority: bindingForm.value.priority,
      enabled: bindingForm.value.enabled,
    }
    if (editingBinding.value) {
      await adminAPI.upstreamPools.updateBinding(editingBinding.value.id, payload)
      appStore.showSuccess('绑定已更新')
    } else {
      await adminAPI.upstreamPools.createBinding(payload)
      appStore.showSuccess('绑定已添加')
    }
    closeBindingModal()
    await loadBindings()
    if (selectedPool.value) await loadMembers(selectedPool.value.id)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存绑定失败'))
  } finally {
    submitting.value = false
  }
}

async function confirmDeletePool(pool: UpstreamPool) {
  if (!window.confirm(buildDeletePoolConfirmMessage({
    pool,
    selectedPoolId: selectedPool.value?.id,
    loadedMembers: members.value,
    memberSets: memberSets.value,
    bindings: bindings.value,
  }))) return
  try {
    await adminAPI.upstreamPools.remove(pool.id)
    appStore.showSuccess('上游池已删除')
    if (selectedPool.value?.id === pool.id) {
      selectedPool.value = null
      members.value = []
    }
    await loadAll()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除上游池失败'))
  }
}

async function confirmDeleteAccountSet(item: UpstreamAccountSet) {
  if (!window.confirm(buildDeleteAccountSetConfirmMessage({
    accountSet: item,
    selectedAccountSetId: selectedAccountSet.value?.id,
    loadedMembers: accountSetMembers.value,
    memberSets: memberSets.value,
  }))) return
  try {
    await adminAPI.upstreamPools.removeAccountSet(item.id)
    appStore.showSuccess('账号集合已删除')
    if (selectedAccountSet.value?.id === item.id) {
      selectedAccountSet.value = null
      accountSetMembers.value = []
    }
    await loadAll()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除账号集合失败'))
  }
}

async function confirmDeleteMember(member: UpstreamPoolMember) {
  if (member.editable === false) {
    appStore.showError('集合展开成员不能直接删除，请去删除集合绑定或从集合里移除账号')
    return
  }
  if (!window.confirm(`确定删除成员 #${member.id} 吗？`)) return
  try {
    await adminAPI.upstreamPools.removeMember(member.id)
    appStore.showSuccess('成员已删除')
    if (selectedPool.value) await loadMembers(selectedPool.value.id)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除成员失败'))
  }
}

async function confirmDeleteAccountSetMember(item: UpstreamAccountSetMember) {
  if (!selectedAccountSet.value) return
  if (!window.confirm(`确定将账号「${item.account_name || item.account_id}」从集合中移除吗？`)) return
  try {
    await adminAPI.upstreamPools.removeAccountSetMember(selectedAccountSet.value.id, item.account_id)
    appStore.showSuccess('集合成员已移除')
    await Promise.all([
      loadAccountSetMembers(selectedAccountSet.value.id),
      loadAll(),
      selectedPool.value ? loadMembers(selectedPool.value.id) : Promise.resolve(),
    ])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '移除集合成员失败'))
  }
}

function editAccountSetMemberCapacity(item: UpstreamAccountSetMember) {
  capacityEditingAccountId.value = item.account_id
  capacityDraft.value = {
    hard: item.capacity_hard_limit || null,
    soft: item.capacity_soft_share || null,
  }
}

function cancelAccountSetMemberCapacity() {
  capacityEditingAccountId.value = null
  capacityDraft.value = { hard: null, soft: null }
}

async function saveAccountSetMemberCapacity(item: UpstreamAccountSetMember) {
  if (!selectedAccountSet.value) return
  const normalize = (value: number | null) => typeof value === 'number' && value > 0 ? Math.round(value) : null
  const hard = normalize(capacityDraft.value.hard)
  const soft = normalize(capacityDraft.value.soft)
  if (hard && selectedAccountSet.value.shared_concurrency_limit && hard > selectedAccountSet.value.shared_concurrency_limit) {
    appStore.showError('成员硬上限不能超过集合共享上限')
    return
  }
  try {
    await adminAPI.upstreamPools.updateAccountSetMemberCapacity(selectedAccountSet.value.id, item.account_id, {
      hard_concurrency_limit: hard,
      soft_concurrency_share: soft,
    })
    appStore.showSuccess(hard || soft ? '成员容量策略已更新' : '成员容量策略已清空')
    cancelAccountSetMemberCapacity()
    await Promise.all([loadAccountSetMembers(selectedAccountSet.value.id), loadCapacityPressures()])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '保存成员容量策略失败'))
  }
}

async function confirmDeleteMemberSet(item: UpstreamPoolMemberSet) {
  if (!window.confirm(buildDeleteMemberSetConfirmMessage(item))) return
  try {
    await adminAPI.upstreamPools.removeMemberSet(item.id)
    appStore.showSuccess('集合绑定已删除')
    if (selectedPool.value) {
      await Promise.all([
        loadMemberSets(selectedPool.value.id),
        loadMembers(selectedPool.value.id),
      ])
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除集合绑定失败'))
  }
}

async function confirmDeleteBinding(binding: UpstreamPoolBinding) {
  if (!window.confirm(buildDeleteBindingConfirmMessage(binding))) return
  try {
    await adminAPI.upstreamPools.removeBinding(binding.id)
    appStore.showSuccess('绑定已删除')
    await loadBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除绑定失败'))
  }
}

let searchTimer: ReturnType<typeof setTimeout> | undefined
function handleSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.value.page = 1
  }, 200)
}

function handleFilterChange() {
  pagination.value.page = 1
  pagination.value.total = filteredPools.value.length
}

function handlePageChange(page: number) {
  pagination.value.page = page
}

function handlePageSizeChange(pageSize: number) {
  pagination.value.page_size = pageSize
  pagination.value.page = 1
}

function handleAccountSetMembersPageChange(page: number) {
  accountSetMembersPagination.value.page = page
}

function handleAccountSetMembersPageSizeChange(pageSize: number) {
  accountSetMembersPagination.value.page_size = pageSize
  accountSetMembersPagination.value.page = 1
  accountSetMembersPagination.value.total = accountSetMembers.value.length
}

async function loadSyncableAccounts() {
  const result: SyncableAccount[] = []
  let page = 1
  let pages = 1
  do {
    const response = await adminAPI.accounts.list(page, 200, {
      lite: '1',
    })
    for (const account of response.items) {
      result.push({
        id: account.id,
        name: account.name,
        platform: account.platform,
        type: account.type,
        status: account.status,
        group_ids: account.group_ids || [],
      })
    }
    pages = response.pages || 1
    page += 1
  } while (page <= pages)
  return result
}

watch(filteredPools, () => {
  pagination.value.total = filteredPools.value.length
})

watch(accountSetMembers, () => {
  accountSetMembersPagination.value.total = accountSetMembers.value.length
  const totalPages = Math.max(1, Math.ceil(accountSetMembers.value.length / accountSetMembersPagination.value.page_size))
  if (accountSetMembersPagination.value.page > totalPages) {
    accountSetMembersPagination.value.page = totalPages
  }
})

onMounted(loadAll)
</script>

<style scoped>
.field-hint {
  margin-bottom: 0.4rem;
  font-size: 0.75rem;
  line-height: 1.35;
  color: rgb(107 114 128);
}

.account-sets-panel,
.overview-panel,
.current-sets-panel {
  border-color: rgba(218, 205, 181, 0.62);
  background:
    linear-gradient(180deg, rgba(255, 253, 249, 0.92), rgba(250, 247, 240, 0.78)),
    rgba(255, 253, 248, 0.82);
  box-shadow: 0 14px 36px rgba(96, 72, 43, 0.035);
}

.account-sets-header,
.overview-panel__header {
  padding-bottom: 0.65rem;
  border-bottom: 1px solid rgba(215, 200, 166, 0.42);
}

.account-sets-guide {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  align-items: center;
  color: #8a4a2b;
  font-size: 0.78rem;
  line-height: 1.5;
}

.account-sets-guide span {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
}

.account-sets-guide span:not(:last-child)::after {
  content: '';
  width: 1.8rem;
  height: 1px;
  background: linear-gradient(90deg, rgba(138, 74, 43, 0.42), transparent);
}

.account-set-members-panel {
  padding-top: 1rem;
  border-top: 1px solid rgba(215, 200, 166, 0.5);
}

.account-set-usage-cell {
  min-width: 9.5rem;
  color: rgb(55 65 81);
}

.account-set-usage-cell__main {
  font-size: 0.82rem;
  font-weight: 600;
  color: rgb(31 41 55);
}

.account-set-usage-cell__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
  margin-top: 0.25rem;
  font-size: 0.68rem;
  color: rgb(107 114 128);
}

.account-set-usage-cell__meta span {
  padding: 0.08rem 0.4rem;
  border-radius: 999px;
  background: rgba(232, 224, 208, 0.56);
}

.overview-sheet {
  padding: 0.25rem 0 0;
}

.overview-refresh-button {
  min-height: 2.5rem;
  padding-inline: 0.95rem;
}

.overview-metric-card {
  min-height: 5.75rem;
  padding: 0.9rem 1rem;
  border: 1px solid rgba(225, 214, 195, 0.68);
  border-radius: 0.7rem;
  background: rgba(255, 252, 246, 0.74);
}

.overview-metric-label {
  letter-spacing: 0.12em;
}

.overview-metric-value {
  letter-spacing: 0;
}

.overview-metric-meta {
  line-height: 1.45;
}

.overview-subpanel {
  padding: 0.2rem 0 0.7rem;
}

.overview-subpanel + .overview-subpanel {
  border-top: 1px solid rgba(215, 200, 166, 0.42);
  padding-top: 0.9rem;
}

.overview-hit-row {
  padding-bottom: 0.1rem;
}

.overview-diagnostic-panel {
  position: relative;
  padding-top: 0.2rem;
}

.overview-detail-button {
  border: 1px solid rgba(215, 200, 166, 0.22);
}

.current-sets-flow {
  padding-bottom: 0.35rem;
}

.current-sets-manage-button {
  min-width: 7.75rem;
  white-space: nowrap;
}

.current-sets-title-strip {
  padding: 0.85rem 0;
  border-top: 1px solid rgba(218, 205, 181, 0.58);
  border-bottom: 1px solid rgba(218, 205, 181, 0.4);
  color: #6f5436;
  background: transparent;
}

.current-sets-pill {
  border-radius: 999px;
  padding: 0.22rem 0.62rem;
  background: rgba(255, 255, 255, 0.72);
  color: #7b6041;
  font-size: 0.72rem;
  box-shadow: inset 0 0 0 1px rgba(215, 200, 166, 0.52);
}

.current-sets-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.5rem;
}

.current-sets-stat {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 0.55rem;
  border-radius: 0.55rem;
  padding: 0.62rem 0.7rem;
  background: rgba(247, 242, 232, 0.64);
}

.current-sets-list {
  overflow: hidden;
  border-top: 1px solid rgba(218, 205, 181, 0.54);
}

.current-sets-item {
  padding: 0.95rem 0;
  border-bottom: 1px solid rgba(218, 205, 181, 0.44);
  background: transparent;
}

.current-sets-item:last-child {
  border-bottom: 0;
}

.current-sets-item-meta {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.5rem;
  color: rgb(75 85 99);
}

.pool-member-account-name {
  width: 13.35rem;
  max-width: 13.35rem;
}

:deep(.pool-member-account-col) {
  width: 13.35rem;
  min-width: 13.35rem;
  max-width: 13.35rem;
}

:deep(.upstream-pools-night .table-wrapper) {
  border-color: rgba(218, 205, 181, 0.54);
  border-radius: 0.7rem;
  background: rgba(255, 253, 248, 0.76);
}

:deep(.upstream-pools-night .table-wrapper .table-header),
:deep(.upstream-pools-night .sticky-header-cell) {
  background: linear-gradient(180deg, rgba(239, 231, 214, 0.92), rgba(231, 221, 201, 0.86));
  color: #6f6758;
}

.pool-status-toggle {
  display: inline-flex;
  min-width: 4.5rem;
  align-items: center;
  justify-content: center;
  gap: 0.35rem;
  border: 1px solid transparent;
  border-radius: 999px;
  padding: 0.22rem 0.62rem;
  font-size: 0.75rem;
  font-weight: 600;
  line-height: 1.25;
  transition:
    border-color 0.16s ease,
    background-color 0.16s ease,
    box-shadow 0.16s ease,
    color 0.16s ease,
    transform 0.16s ease;
}

.pool-status-toggle:hover:not(:disabled) {
  transform: translateY(-1px);
}

.pool-status-toggle:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px rgba(14, 165, 233, 0.2);
}

.pool-status-toggle:disabled {
  cursor: wait;
  opacity: 0.72;
}

.pool-status-toggle.is-enabled {
  border-color: rgba(16, 185, 129, 0.28);
  background-color: rgba(16, 185, 129, 0.12);
  color: rgb(4, 120, 87);
}

.pool-status-toggle.is-disabled {
  border-color: rgba(239, 68, 68, 0.26);
  background-color: rgba(239, 68, 68, 0.1);
  color: rgb(185, 28, 28);
}

.pool-status-toggle__dot {
  width: 0.42rem;
  height: 0.42rem;
  border-radius: 999px;
  background-color: currentColor;
  box-shadow: 0 0 0 3px color-mix(in srgb, currentColor 16%, transparent);
}

</style>

<style>
.dark .upstream-pools-night {
  --upstream-night-bg: #0b100d;
  --upstream-night-panel: rgba(17, 22, 17, 0.92);
  --upstream-night-panel-soft: rgba(24, 29, 23, 0.86);
  --upstream-night-panel-raised: rgba(31, 34, 28, 0.9);
  --upstream-night-line: rgba(229, 220, 198, 0.16);
  --upstream-night-line-strong: rgba(229, 220, 198, 0.24);
  --upstream-night-text: #f1eadb;
  --upstream-night-muted: #b9b09d;
  --upstream-night-subtle: #8f8878;
  --upstream-night-cyan: #41d7c7;
  --upstream-night-cyan-soft: rgba(65, 215, 199, 0.16);
  --upstream-night-amber: #d8bd72;
  --upstream-night-control: rgba(35, 43, 35, 0.94);
  --upstream-night-control-hover: rgba(43, 52, 42, 0.96);
  color: var(--upstream-night-text);
  background:
    radial-gradient(circle at top left, rgba(216, 189, 114, 0.08), transparent 34%),
    radial-gradient(circle at top right, rgba(65, 215, 199, 0.05), transparent 28%),
    linear-gradient(180deg, rgba(10, 13, 10, 0.96), rgba(9, 11, 9, 0.98));
}

.dark .upstream-pools-night .field-hint {
  color: rgb(156 163 175);
}

.dark .upstream-pools-night .sst-admin-panel,
.dark .upstream-pools-night .sst-admin-panel {
  border-color: var(--upstream-night-line) !important;
  background:
    linear-gradient(180deg, rgba(229, 220, 198, 0.035), rgba(229, 220, 198, 0.012)),
    var(--upstream-night-panel);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.035);
}

.dark .upstream-pools-night .rounded-2xl,
.dark .upstream-pools-night .rounded-xl {
  border-color: var(--upstream-night-line) !important;
}

.dark .upstream-pools-night .bg-white\/80,
.dark .upstream-pools-night .bg-white\/85,
.dark .upstream-pools-night .bg-white\/90,
.dark .upstream-pools-night .bg-white\/70,
.dark .upstream-pools-night .bg-gray-50\/55,
.dark .upstream-pools-night .bg-gray-50\/65,
.dark .upstream-pools-night .bg-gray-50\/70,
.dark .upstream-pools-night .bg-gray-50\/80,
.dark .upstream-pools-night .dark\:bg-dark-900\/40,
.dark .upstream-pools-night .dark\:bg-dark-900\/50,
.dark .upstream-pools-night .dark\:bg-dark-900\/60,
.dark .upstream-pools-night .dark\:bg-dark-900\/70,
.dark .upstream-pools-night .dark\:bg-dark-800\/35,
.dark .upstream-pools-night .dark\:bg-dark-800\/55,
.dark .upstream-pools-night .dark\:bg-dark-800\/50,
.dark .upstream-pools-night .dark\:bg-dark-800\/60 {
  background-color: var(--upstream-night-panel-soft) !important;
}

.dark .upstream-pools-night .text-gray-900,
.dark .upstream-pools-night .dark\:text-white,
.dark .upstream-pools-night .dark\:text-gray-100 {
  color: var(--upstream-night-text);
}

.dark .upstream-pools-night .text-gray-700,
.dark .upstream-pools-night .dark\:text-gray-200,
.dark .upstream-pools-night .dark\:text-gray-300 {
  color: #d8cfbd;
}

.dark .upstream-pools-night .text-gray-500,
.dark .upstream-pools-night .dark\:text-gray-400,
.dark .upstream-pools-night .dark\:text-dark-400 {
  color: var(--upstream-night-muted);
}

.dark .upstream-pools-night .text-gray-400,
.dark .upstream-pools-night .dark\:text-gray-500 {
  color: var(--upstream-night-subtle);
}

.dark .upstream-pools-night .border-gray-200,
.dark .upstream-pools-night .border-gray-300,
.dark .upstream-pools-night .border-white\/70,
.dark .upstream-pools-night .dark\:border-dark-700,
.dark .upstream-pools-night .dark\:border-dark-600 {
  border-color: var(--upstream-night-line) !important;
}

.dark .upstream-pools-night .bg-gray-100,
.dark .upstream-pools-night .dark\:bg-dark-700,
.dark .upstream-pools-night .dark\:bg-dark-800 {
  background-color: rgba(229, 220, 198, 0.08) !important;
}

.dark .upstream-pools-night .bg-primary-50,
.dark .upstream-pools-night .dark\:bg-primary-950\/20,
.dark .upstream-pools-night .dark\:bg-primary-950\/30 {
  background-color: var(--upstream-night-cyan-soft) !important;
}

.dark .upstream-pools-night .text-primary-700,
.dark .upstream-pools-night .dark\:text-primary-100,
.dark .upstream-pools-night .dark\:text-primary-200,
.dark .upstream-pools-night .dark\:text-primary-300 {
  color: #8ee9dc;
}

.dark .upstream-pools-night .border-primary-200,
.dark .upstream-pools-night .dark\:border-primary-900\/40 {
  border-color: rgba(65, 215, 199, 0.28) !important;
}

.dark .upstream-pools-night .bg-amber-50\/80,
.dark .upstream-pools-night .bg-amber-100,
.dark .upstream-pools-night .dark\:bg-amber-950\/20,
.dark .upstream-pools-night .dark\:bg-amber-900\/20 {
  background-color: rgba(144, 105, 42, 0.18) !important;
}

.dark .upstream-pools-night .border-amber-200,
.dark .upstream-pools-night .dark\:border-amber-900\/40 {
  border-color: rgba(216, 189, 114, 0.28) !important;
}

.dark .upstream-pools-night .text-amber-900,
.dark .upstream-pools-night .text-amber-800,
.dark .upstream-pools-night .dark\:text-amber-100,
.dark .upstream-pools-night .dark\:text-amber-200 {
  color: #f1dfaa;
}

.dark .upstream-pools-night .badge-gray {
  background-color: rgba(229, 220, 198, 0.12);
  color: #cfc5b2;
}

.dark .upstream-pools-night .badge-success {
  background-color: rgba(20, 148, 104, 0.22);
  color: #69e0b4;
}

.dark .upstream-pools-night .badge-warning {
  background-color: rgba(180, 117, 28, 0.24);
  color: #f1c46f;
}

.dark .upstream-pools-night .badge-danger {
  background-color: rgba(174, 48, 45, 0.24);
  color: #ff9a92;
}

.dark .upstream-pools-night .pool-status-toggle.is-enabled {
  border-color: rgba(105, 224, 180, 0.32);
  background-color: rgba(20, 148, 104, 0.22);
  color: #69e0b4;
}

.dark .upstream-pools-night .pool-status-toggle.is-disabled {
  border-color: rgba(255, 154, 146, 0.3);
  background-color: rgba(174, 48, 45, 0.24);
  color: #ff9a92;
}

.dark .upstream-pools-night .pool-status-toggle:focus-visible {
  box-shadow: 0 0 0 3px rgba(65, 215, 199, 0.18);
}

.dark .upstream-pools-night .btn.btn-secondary,
.dark .upstream-pools-night button.rounded,
.dark .upstream-pools-night select,
.dark .upstream-pools-night input:not([type='checkbox']):not([type='radio']),
.dark .upstream-pools-night textarea {
  border-color: var(--upstream-night-line-strong) !important;
  background-color: var(--upstream-night-control) !important;
  color: var(--upstream-night-text) !important;
}

.dark .upstream-pools-night .btn.btn-secondary:hover,
.dark .upstream-pools-night button.rounded:hover,
.dark .upstream-pools-night select:hover,
.dark .upstream-pools-night input:not([type='checkbox']):not([type='radio']):hover,
.dark .upstream-pools-night textarea:hover {
  background-color: var(--upstream-night-control-hover) !important;
}

.dark .upstream-pools-night select option {
  background-color: #182017;
  color: var(--upstream-night-text);
}

.dark .upstream-pools-night .table-wrapper {
  border: 1px solid var(--upstream-night-line);
  border-radius: 0.75rem;
  background: var(--upstream-night-panel);
}

.dark .upstream-pools-night .table-wrapper table {
  color: var(--upstream-night-text);
}

.dark .upstream-pools-night .table-wrapper .table-header,
.dark .upstream-pools-night .sticky-header-cell {
  background-color: #23241d !important;
  color: #cfc5b2;
}

.dark .upstream-pools-night .table-body {
  background-color: var(--upstream-night-bg) !important;
}

.dark .upstream-pools-night .table-body tr {
  background-color: rgba(13, 18, 14, 0.92) !important;
}

.dark .upstream-pools-night .table-body tr:hover,
.dark .upstream-pools-night tbody tr:hover .sticky-col {
  background-color: rgba(33, 42, 32, 0.95) !important;
}

.dark .upstream-pools-night tbody .sticky-col {
  background-color: #10160f !important;
}

.dark .upstream-pools-night .sticky-col-right {
  background-color: #151a14 !important;
}

.dark .upstream-pools-night .divide-gray-200 > :not([hidden]) ~ :not([hidden]),
.dark .upstream-pools-night .dark\:divide-dark-700 > :not([hidden]) ~ :not([hidden]) {
  border-color: var(--upstream-night-line) !important;
}

.dark .upstream-pools-night .table-wrapper::-webkit-scrollbar-track {
  background-color: rgba(229, 220, 198, 0.08) !important;
}

.dark .upstream-pools-night .table-wrapper::-webkit-scrollbar-thumb {
  background-color: rgba(207, 197, 178, 0.72) !important;
}

.dark .upstream-pools-night .overview-metric-card {
  background:
    linear-gradient(180deg, rgba(229, 220, 198, 0.035), rgba(229, 220, 198, 0.012)),
    var(--upstream-night-panel-raised) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.03),
    0 0 0 1px rgba(229, 220, 198, 0.04);
}

.dark .upstream-pools-night .overview-metric-label {
  color: #9e9788;
}

.dark .upstream-pools-night .overview-metric-meta {
  color: #b5ac99;
}

.dark .upstream-pools-night .overview-subpanel {
  background: rgba(24, 30, 24, 0.72) !important;
  border: 1px solid rgba(229, 220, 198, 0.08);
}

.dark .upstream-pools-night .overview-subpanel-label {
  color: #8f8878;
}

.dark .upstream-pools-night .overview-hit-row .text-gray-700,
.dark .upstream-pools-night .overview-hit-row .dark\:text-gray-200 {
  color: #ddd4c2;
}

.dark .upstream-pools-night .overview-diagnostic-panel {
  background: linear-gradient(180deg, rgba(20, 26, 20, 0.32), rgba(20, 26, 20, 0.12));
  border-radius: 1rem;
  padding-top: 0.35rem;
  padding-bottom: 0.35rem;
}

.dark .upstream-pools-night .overview-detail-button {
  background: rgba(29, 35, 29, 0.88) !important;
  color: #d8cfbd !important;
  border-color: rgba(229, 220, 198, 0.16);
}

.dark .upstream-pools-night .overview-detail-button:hover {
  background: rgba(41, 49, 41, 0.96) !important;
  color: var(--upstream-night-text) !important;
}

.dark .upstream-pools-night .overview-signal-chip {
  background: rgba(112, 72, 18, 0.22) !important;
  color: #f0d99b !important;
}

.dark .upstream-pools-night .h-1\.5.bg-gray-200 {
  background-color: rgba(229, 220, 198, 0.13);
}

.dark .upstream-pools-night .h-1\.5.bg-primary-500 {
  background-color: var(--upstream-night-cyan);
}
</style>
