<template>
  <AppLayout>
    <div class="sst-admin-page">

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

            <template #cell-enabled="{ value }">
              <span :class="['badge', value ? 'badge-success' : 'badge-danger']">
                {{ value ? '启用' : '停用' }}
              </span>
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
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">故障转移 {{ row.failover_enabled ? '开' : '关' }}</span>
                  <span class="rounded bg-gray-100 px-2 py-0.5 dark:bg-dark-700">TopK {{ row.top_k }}</span>
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
                <button class="text-sm text-primary-700 hover:underline dark:text-primary-300" @click="selectPool(row)">
                {{ selectedPool?.id === row.id ? (members.length || 0) : '—' }} 个成员
                </button>
              </template>

            <template #cell-bindings="{ row }">
              <button class="text-sm text-primary-700 hover:underline dark:text-primary-300" @click="selectPool(row)">
                {{ selectedPool?.id === row.id ? (bindings.filter(item => item.pool_id === row.id).length || 0) : '—' }} 个绑定
              </button>
            </template>

            <template #cell-created_at="{ value }">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDateTime(value) }}</span>
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

      <div class="grid gap-4 xl:grid-cols-[minmax(0,1.45fr)_minmax(320px,0.85fr)]">
        <section class="sst-admin-panel p-4">
          <div class="mb-3 flex items-center justify-between gap-2">
            <div>
              <h3 class="text-base font-semibold text-gray-900 dark:text-white">成员</h3>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ selectedPool ? `${selectedPool.name} 的成员列表` : '选择一个池后查看成员' }}
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
          <div class="mb-3 rounded-xl border border-amber-200 bg-amber-50/80 px-3 py-3 text-xs leading-6 text-amber-900 dark:border-amber-900/40 dark:bg-amber-950/20 dark:text-amber-100">
            <div class="font-medium">恢复探针</div>
            <div>
              后台自动运行，这里没有手动按钮。只探“异常/可恢复”成员，正常成员不定时探测。探针走 compact 最小测试单元，尽量少耗 token。每 1 分钟扫描一次；成员 2 分钟内刚被使用过会先跳过。探测成功后，会自动清理临时故障/限流状态，并重新加入调度。
            </div>
          </div>
          <div class="mb-3 rounded-2xl border border-gray-200 bg-gray-50/80 p-3 dark:border-dark-700 dark:bg-dark-900/40">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <div class="text-sm font-semibold text-gray-900 dark:text-white">最近路由观测</div>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  看这个池最近 24 小时的 `routing.explanation`，确认阈值是否真的在影响选路。
                </p>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <button
                  class="btn btn-secondary btn-sm"
                  :disabled="!selectedPool || !poolRoutingObservabilitySupported || poolRoutingObservability.loading"
                  @click="selectedPool && loadPoolObservability(selectedPool)"
                >
                  <Icon name="refresh" size="sm" :class="poolRoutingObservability.loading ? 'animate-spin' : ''" class="mr-1" />
                  刷新观测
                </button>
                <button
                  class="btn btn-secondary btn-sm"
                  :disabled="!selectedPool || !poolRoutingObservabilitySupported || poolRoutingObservability.loading"
                  @click="openPoolObservabilityModal"
                >
                  查看明细
                </button>
              </div>
            </div>

            <div
              v-if="!poolRoutingObservabilitySupported"
              class="mt-3 rounded-xl border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400"
            >
              当前仅对 OpenAI 池提供这组路由观测。
            </div>

            <template v-else>
              <div class="mt-3 grid gap-2 md:grid-cols-4">
                <div class="rounded-xl border border-white/70 bg-white/90 px-3 py-2 dark:border-dark-700 dark:bg-dark-900/70">
                  <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">观测数</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">
                    {{ poolRoutingObservability.total || poolRoutingObservability.logs.length }}
                  </div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">最近 24 小时命中记录</div>
                </div>
                <div class="rounded-xl border border-white/70 bg-white/90 px-3 py-2 dark:border-dark-700 dark:bg-dark-900/70">
                  <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">粘性逃逸</div>
                  <div class="mt-1 text-lg font-semibold text-gray-900 dark:text-white">{{ stickyEscapeCount }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">{{ stickyEscapeReasonSummary }}</div>
                </div>
                <div class="rounded-xl border border-white/70 bg-white/90 px-3 py-2 dark:border-dark-700 dark:bg-dark-900/70">
                  <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">最近账号</div>
                  <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ recentRoutedAccountSummary }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">按最新路由记录去重</div>
                </div>
                <div class="rounded-xl border border-white/70 bg-white/90 px-3 py-2 dark:border-dark-700 dark:bg-dark-900/70">
                  <div class="text-[11px] uppercase tracking-[0.16em] text-gray-400 dark:text-gray-500">最近时间</div>
                  <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ latestRoutingLogAt }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">
                    {{ poolRoutingObservability.loading ? '正在刷新…' : '按日志倒序展示' }}
                  </div>
                </div>
              </div>

              <div
                v-if="poolRoutingObservability.logs.length === 0 && !poolRoutingObservability.loading"
                class="mt-3 rounded-xl border border-dashed border-gray-300 px-3 py-3 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400"
              >
                最近 24 小时还没有这个池的路由解释日志。
              </div>
            </template>
          </div>
          <DataTable :columns="memberColumns" :data="members" :loading="membersLoading" row-key="id">
            <template #cell-account_name="{ row, value }">
              <div class="flex flex-col gap-1">
                <div class="max-w-[320px] whitespace-normal break-words font-medium leading-6 text-gray-900 dark:text-white">
                  {{ formatPoolMemberAccountName(value, row.account_id) }}
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">{{ row.account_platform || '未知平台' }}</div>
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

            <template #cell-weight="{ value }">
              <span class="font-mono text-sm">{{ value }}</span>
            </template>

            <template #cell-updated_at="{ value }">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatCompactDateTime(value) }}</span>
            </template>

            <template #cell-actions="{ row }">
              <div class="flex gap-1">
                <button class="rounded px-2 py-1 text-xs text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700" @click="openMemberModal(row)">
                  编辑
                </button>
                <button class="rounded px-2 py-1 text-xs text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-900/20" @click="confirmDeleteMember(row)">
                  删除
                </button>
              </div>
            </template>

            <template #empty>
              <EmptyState title="暂无成员" description="先把账号加入当前上游池。" action-text="添加成员" @action="openMemberModal()" />
            </template>
          </DataTable>
        </section>

        <section class="sst-admin-panel p-4">
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

            <template #cell-created_at="{ value }">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDateTime(value) }}</span>
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
              共 {{ poolRoutingObservability.total || poolRoutingObservability.logs.length }} 条记录，按日志时间倒序展示。
            </p>
          </div>
          <button
            class="btn btn-secondary btn-sm self-start md:self-auto"
            :disabled="!selectedPool || !poolRoutingObservabilitySupported || poolRoutingObservability.loading"
            @click="selectedPool && loadPoolObservability(selectedPool)"
          >
            <Icon name="refresh" size="sm" :class="poolRoutingObservability.loading ? 'animate-spin' : ''" class="mr-1" />
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
                {{ formatObservabilityAccountLabel(log.account_id) }}
              </span>
              <span
                v-if="getLogExtraBool(log, 'sticky_escape_triggered')"
                class="rounded bg-amber-100 px-2 py-0.5 text-amber-800 dark:bg-amber-900/30 dark:text-amber-200"
              >
                触发逃逸 · {{ formatStickyEscapeReason(getLogExtraString(log, 'sticky_escape_reason')) }}
              </span>
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
            <Select v-model="poolForm.platform" :options="platformOptions" />
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
        <div class="grid gap-4 md:grid-cols-2">
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

    <BaseDialog :show="showBindingModal" :title="editingBinding ? '编辑绑定' : '添加绑定'" width="wide" @close="closeBindingModal">
      <form class="space-y-4" @submit.prevent="submitBinding">
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">分组</label>
            <p class="field-hint">哪些用户分组走这个池。</p>
            <Select v-model="bindingForm.group_id" :options="groupOptions" />
          </div>
          <div>
            <label class="input-label">上游池</label>
            <p class="field-hint">请求要分配到的账号池。</p>
            <Select v-model="bindingForm.pool_id" :options="poolOptions" />
          </div>
        </div>
        <div class="grid gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">平台</label>
            <p class="field-hint">必须和池的平台一致。</p>
            <Select v-model="bindingForm.platform" :options="platformOptions" />
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
import type { Account, AdminGroup, UpstreamPool, UpstreamPoolBinding, UpstreamPoolMember } from '@/types'
import { extractApiErrorMessage } from '@/utils/apiError'
import { formatDateTime } from '@/utils/format'
import { useAppStore } from '@/stores'

type PoolForm = {
  name: string
  code: string
  platform: string
  description: string
  enabled: boolean
  scheduler_mode: string
  sticky_enabled: boolean
  sticky_escape_enabled: boolean
  sticky_escape_error_rate_threshold: number
  sticky_escape_ttft_ms_threshold: number
  load_balance_enabled: boolean
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

type PoolRoutingObservabilityState = {
  loading: boolean
  total: number
  logs: OpsSystemLog[]
}

const appStore = useAppStore()
const loading = ref(false)
const submitting = ref(false)
const pools = ref<UpstreamPool[]>([])
const members = ref<UpstreamPoolMember[]>([])
const bindings = ref<UpstreamPoolBinding[]>([])
const allGroups = ref<AdminGroup[]>([])
const syncableAccounts = ref<SyncableAccount[]>([])
const selectedPool = ref<UpstreamPool | null>(null)
const membersLoading = ref(false)
const bindingsLoading = ref(false)
const syncingMembers = ref(false)
const poolRoutingObservability = ref<PoolRoutingObservabilityState>({
  loading: false,
  total: 0,
  logs: [],
})
const searchQuery = ref('')
const filters = ref({ platform: '', enabled: '' })
const pagination = ref({ page: 1, page_size: 20, total: 0 })
let poolObservabilityRequestToken = 0

const showPoolModal = ref(false)
const showMemberModal = ref(false)
const showBindingModal = ref(false)
const showPoolObservabilityModal = ref(false)
const editingPool = ref<UpstreamPool | null>(null)
const editingMember = ref<UpstreamPoolMember | null>(null)
const editingBinding = ref<UpstreamPoolBinding | null>(null)

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
  sticky_enabled: true,
  sticky_escape_enabled: true,
  sticky_escape_error_rate_threshold: 0.3,
  sticky_escape_ttft_ms_threshold: 6000,
  load_balance_enabled: true,
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

const boolSelectOptions = [
  { value: null, label: '不覆盖' },
  { value: true, label: '是' },
  { value: false, label: '否' },
]

const canSyncSelectedPool = computed(() => selectedPool.value?.platform === 'openai')
const poolRoutingObservabilitySupported = computed(() => selectedPool.value?.platform === 'openai')
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

const poolColumns = [
  { key: 'name', label: '池' },
  { key: 'platform', label: '平台' },
  { key: 'enabled', label: '状态' },
  { key: 'routing', label: '路由' },
  { key: 'sticky_enabled', label: '粘性' },
  { key: 'members', label: '成员' },
  { key: 'bindings', label: '绑定' },
  { key: 'created_at', label: '创建时间' },
  { key: 'actions', label: '操作' },
]

const memberColumns = [
  { key: 'account_name', label: '账号' },
  { key: 'enabled', label: '状态' },
  { key: 'runtime_status', label: '运行态' },
  { key: 'manual_drained', label: '排空' },
  { key: 'weight', label: '权重' },
  { key: 'updated_at', label: '更新时间' },
  { key: 'actions', label: '操作' },
]

const bindingColumns = [
  { key: 'group_name', label: '分组' },
  { key: 'platform', label: '平台' },
  { key: 'enabled', label: '状态' },
  { key: 'priority', label: '优先级' },
  { key: 'created_at', label: '创建时间' },
  { key: 'actions', label: '操作' },
]

const platformLabel = (value: string) => platformOptions.find(item => item.value === value)?.label || value || '-'

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

const formatPoolMemberAccountName = (value?: string | null, accountID?: number) => {
  const text = String(value || `账号 #${accountID ?? '-'}`)
  const atIndex = text.indexOf('@')
  if (atIndex <= 0 || atIndex >= text.length - 1) {
    return text
  }
  return `${text.slice(0, atIndex + 1)}\u200b${text.slice(atIndex + 1)}`
}

const formatThresholdMs = (value?: number | null) => {
  if (!value || value <= 0) return '-'
  return `${Math.round(value)}ms`
}

const formatRateThreshold = (value?: number | null) => {
  if (typeof value !== 'number' || Number.isNaN(value) || value < 0) return '-'
  return value.toFixed(2)
}

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

const availablePoolAccounts = computed(() => {
  const platform = selectedPool.value?.platform
  if (!platform) return []
  return syncableAccounts.value.filter(account => account.platform === platform)
})

const accountOptions = computed(() =>
  availablePoolAccounts.value.map((account) => ({
    value: account.id,
    label: `${account.name} · ${platformLabel(account.platform)} · ${account.type === 'apikey' ? 'API Key' : account.type}`,
  }))
)

const groupOptions = computed(() =>
  allGroups.value.map((group) => ({
    value: group.id,
    label: `${group.name} · ${platformLabel(group.platform)}`,
  }))
)

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
    sticky_enabled: true,
    sticky_escape_enabled: true,
    sticky_escape_error_rate_threshold: 0.3,
    sticky_escape_ttft_ms_threshold: 6000,
    load_balance_enabled: true,
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

async function loadAll() {
  loading.value = true
  try {
    const [poolList, groupList, accountList, bindingList] = await Promise.all([
      adminAPI.upstreamPools.list(),
      adminAPI.groups.getAllIncludingInactive(),
      loadSyncableAccounts(),
      adminAPI.upstreamPools.getBindings(),
    ])
    pools.value = poolList
    allGroups.value = groupList
    syncableAccounts.value = accountList
    bindings.value = bindingList
    pagination.value.total = filteredPools.value.length
    if (poolList.length === 0) {
      selectedPool.value = null
      members.value = []
      poolRoutingObservability.value = { loading: false, total: 0, logs: [] }
    } else if (selectedPool.value) {
      const nextSelected = poolList.find(pool => pool.id === selectedPool.value?.id)
      selectPool(nextSelected || poolList[0])
    } else {
      selectPool(poolList[0])
    }
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载上游池失败'))
  } finally {
    loading.value = false
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
  loadPoolObservability(pool)
  pagination.value.total = filteredPools.value.length
}

async function loadPoolObservability(pool: UpstreamPool) {
  if (pool.platform !== 'openai') {
    poolRoutingObservability.value = { loading: false, total: 0, logs: [] }
    return
  }

  const requestToken = ++poolObservabilityRequestToken
  poolRoutingObservability.value = {
    loading: true,
    total: poolRoutingObservability.value.total,
    logs: poolRoutingObservability.value.logs,
  }

  try {
    const result = await adminAPI.ops.listSystemLogs({
      page: 1,
      page_size: 60,
      time_range: '24h',
      component: 'routing.explanation',
      platform: pool.platform,
      pool_id: pool.id,
    })
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = {
      loading: false,
      total: result.total || 0,
      logs: result.items || [],
    }
  } catch (error) {
    if (requestToken !== poolObservabilityRequestToken) return
    poolRoutingObservability.value = { loading: false, total: 0, logs: [] }
    appStore.showError(extractApiErrorMessage(error, '加载最近路由观测失败'))
  }
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
      sticky_enabled: pool.sticky_enabled,
      sticky_escape_enabled: pool.sticky_escape_enabled,
      sticky_escape_error_rate_threshold: pool.sticky_escape_error_rate_threshold,
      sticky_escape_ttft_ms_threshold: pool.sticky_escape_ttft_ms_threshold,
      load_balance_enabled: pool.load_balance_enabled,
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

async function syncSelectedPoolMembers() {
  if (!selectedPool.value) return
  if (!canSyncSelectedPool.value) {
    appStore.showError('当前只支持同步 OpenAI 上游池')
    return
  }
  if (availablePoolAccounts.value.length === 0) {
    appStore.showError('当前平台没有可同步的账号')
    return
  }
  if (!window.confirm('全量同步会把当前平台账号同步到这个池，并更新或移除不再匹配的成员。若只想给当前池单独加账号，请用「添加成员」。继续吗？')) {
    return
  }

  syncingMembers.value = true
  try {
    const poolMembers = await adminAPI.upstreamPools.getMembers(selectedPool.value.id)
    const memberByAccountId = new Map(poolMembers.map(member => [member.account_id, member]))
    const sourceAccountIds = new Set(availablePoolAccounts.value.map(account => account.id))

    let created = 0
    let updated = 0
    let removed = 0

    for (const account of availablePoolAccounts.value) {
      const member = memberByAccountId.get(account.id)
      const enabled = account.status === 'active'
      const payload = {
        account_id: account.id,
        enabled,
        schedulable_override: null,
        manual_drained: !enabled,
        weight: 100,
        priority_override: null,
        max_concurrency_override: null,
        notes: '从账号管理同步',
      }

      if (member) {
        await adminAPI.upstreamPools.updateMember(member.id, payload)
        updated += 1
      } else {
        await adminAPI.upstreamPools.createMember(selectedPool.value.id, payload)
        created += 1
      }
    }

    for (const member of poolMembers) {
      if (sourceAccountIds.has(member.account_id)) continue
      await adminAPI.upstreamPools.removeMember(member.id)
      removed += 1
    }

    appStore.showSuccess(`同步完成：新增 ${created}，更新 ${updated}，移除 ${removed}`)
    await loadMembers(selectedPool.value.id)
    await loadBindings()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '同步账号失败'))
  } finally {
    syncingMembers.value = false
  }
}

async function addMissingSelectedPoolMembers() {
  if (!selectedPool.value) return
  if (!canSyncSelectedPool.value) {
    appStore.showError('当前只支持补齐 OpenAI 上游池')
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
  if (!window.confirm(`确定删除上游池「${pool.name}」吗？`)) return
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

async function confirmDeleteMember(member: UpstreamPoolMember) {
  if (!window.confirm(`确定删除成员 #${member.id} 吗？`)) return
  try {
    await adminAPI.upstreamPools.removeMember(member.id)
    appStore.showSuccess('成员已删除')
    if (selectedPool.value) await loadMembers(selectedPool.value.id)
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '删除成员失败'))
  }
}

async function confirmDeleteBinding(binding: UpstreamPoolBinding) {
  if (!window.confirm(`确定删除绑定 #${binding.id} 吗？`)) return
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

onMounted(loadAll)
</script>

<style scoped>
.field-hint {
  margin-bottom: 0.4rem;
  font-size: 0.75rem;
  line-height: 1.35;
  color: rgb(107 114 128);
}

:global(.dark) .field-hint {
  color: rgb(156 163 175);
}
</style>

