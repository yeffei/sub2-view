<template>
  <AppLayout>
    <div class="sst-admin-page">
      <TablePageLayout>
        <template #filters>
        <div class="flex flex-wrap-reverse items-start justify-between gap-3">
          <AccountTableFilters
            v-model:searchQuery="params.search"
            :filters="params"
            :groups="groups"
            @update:filters="(newFilters) => Object.assign(params, newFilters)"
            @change="debouncedReload"
            @update:searchQuery="debouncedReload"
          />
          <AccountTableActions
            :loading="loading"
            @refresh="handleManualRefresh"
            @create="openCreateAccount"
          >
            <template #beforeCreate>
              <button
                type="button"
                class="btn btn-secondary gap-2"
                :disabled="scheduledCostRateSyncLoading || scheduledCostRateSyncSaving"
                :title="t('admin.accounts.scheduledCostRateSyncHint')"
                :aria-pressed="scheduledCostRateSyncEnabled"
                @click="toggleScheduledCostRateSync"
              >
                <span
                  class="relative inline-flex h-5 w-9 flex-shrink-0 rounded-full transition-colors"
                  :class="scheduledCostRateSyncEnabled ? 'bg-teal-500' : 'bg-gray-300 dark:bg-gray-600'"
                >
                  <span
                    class="mt-0.5 inline-block h-4 w-4 rounded-full bg-white shadow transition-transform"
                    :class="scheduledCostRateSyncEnabled ? 'translate-x-[1.125rem]' : 'translate-x-0.5'"
                  />
                </span>
                <span>{{ t('admin.accounts.scheduledCostRateSync') }}</span>
              </button>
              <button
                class="btn btn-secondary"
                :disabled="syncingAllUpstreamRates"
                :title="t('admin.accounts.syncAllUpstreamRatesHint')"
                @click="showSyncAllUpstreamRatesDialog = true"
              >
                <Icon name="sync" size="sm" :class="syncingAllUpstreamRates ? 'animate-spin' : ''" />
                <span>{{ syncingAllUpstreamRates ? t('admin.accounts.syncingAllUpstreamRates') : t('admin.accounts.syncAllUpstreamRates') }}</span>
              </button>
              <button
                v-if="failedUpstreamRateAccountIds.length > 0"
                class="btn btn-secondary"
                :disabled="syncingRateBatch || syncingAllUpstreamRates"
                @click="handleRetryFailedUpstreamRates"
              >
                <Icon name="refresh" size="sm" :class="syncingRateBatch ? 'animate-spin' : ''" />
                <span>{{ t('admin.accounts.retryFailedUpstreamRates', { count: failedUpstreamRateAccountIds.length }) }}</span>
              </button>
			  <button
				v-if="lastUpstreamRateSyncResult"
				class="btn btn-secondary"
				@click="showUpstreamRateSyncResult = true"
			  >
				<Icon name="document" size="sm" />
				<span>{{ t('admin.accounts.viewUpstreamRateSyncResult') }}</span>
			  </button>
            </template>
            <template #after>
              <!-- Auto Refresh Dropdown -->
              <div class="relative" ref="autoRefreshDropdownRef">
                <button
                  @click="
                    showAutoRefreshDropdown = !showAutoRefreshDropdown;
                    showAccountToolsDropdown = false
                  "
                  class="btn btn-secondary px-2 md:px-3"
                  :title="t('admin.accounts.autoRefresh')"
                >
                  <Icon name="refresh" size="sm" :class="[autoRefreshEnabled ? 'animate-spin' : '']" />
                  <span class="hidden md:inline">
                    {{
                      autoRefreshEnabled
                        ? t('admin.accounts.autoRefreshCountdown', { seconds: autoRefreshCountdown })
                        : t('admin.accounts.autoRefresh')
                    }}
                  </span>
                </button>
                <div
                  v-if="showAutoRefreshDropdown"
                  class="absolute right-0 z-50 mt-2 w-56 origin-top-right rounded-lg border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
                >
                  <div class="p-2">
                    <button
                      @click="setAutoRefreshEnabled(!autoRefreshEnabled)"
                      class="flex w-full items-center justify-between rounded-md px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
                    >
                      <span>{{ t('admin.accounts.enableAutoRefresh') }}</span>
                      <Icon v-if="autoRefreshEnabled" name="check" size="sm" class="text-primary-500" />
                    </button>
                    <div class="my-1 border-t border-gray-100 dark:border-gray-700"></div>
                    <button
                      v-for="sec in autoRefreshIntervals"
                      :key="sec"
                      @click="setAutoRefreshInterval(sec)"
                      class="flex w-full items-center justify-between rounded-md px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
                    >
                      <span>{{ autoRefreshIntervalLabel(sec) }}</span>
                      <Icon v-if="autoRefreshIntervalSeconds === sec" name="check" size="sm" class="text-primary-500" />
                    </button>
                  </div>
                </div>
              </div>

              <!-- More Tools Dropdown -->
              <div class="relative" ref="accountToolsDropdownRef">
                <button
                  @click="
                    toggleAccountToolsDropdown($event);
                    showAutoRefreshDropdown = false
                  "
                  class="btn btn-secondary px-2 md:px-3"
                  :title="t('admin.accounts.moreActions')"
                >
                  <Icon name="more" size="sm" class="md:mr-1.5" />
                  <span class="hidden md:inline">{{ t('admin.accounts.moreActions') }}</span>
                  <Icon name="chevronDown" size="xs" class="ml-1 hidden md:inline" />
                </button>
                <div
                  v-if="showAccountToolsDropdown && accountToolsDropdownPosition"
                  class="account-tools-dropdown fixed z-[9999] w-[min(20rem,calc(100vw-2rem))] origin-top-right overflow-hidden rounded-lg border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-800"
                  :style="{
                    top: `${accountToolsDropdownPosition.top}px`,
                    left: `${accountToolsDropdownPosition.left}px`
                  }"
                >
                  <div class="max-h-[70vh] overflow-y-auto p-2">
                    <div class="px-2 py-2">
                      <div class="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                        {{ t('admin.accounts.dataActions') }}
                      </div>
                    </div>
                    <button class="account-tools-menu-item" @click="openSyncFromCrs">
                      <span class="account-tools-menu-icon bg-blue-50 text-blue-600 dark:bg-blue-900/30 dark:text-blue-300">
                        <Icon name="sync" size="sm" />
                      </span>
                      <span class="flex-1 text-left">{{ t('admin.accounts.syncFromCrs') }}</span>
                    </button>
                    <button class="account-tools-menu-item" @click="openImportData">
                      <span class="account-tools-menu-icon bg-emerald-50 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-300">
                        <Icon name="upload" size="sm" />
                      </span>
                      <span class="flex-1 text-left">{{ t('admin.accounts.dataImport') }}</span>
                    </button>
                    <button class="account-tools-menu-item" @click="openExportDataDialogFromMenu">
                      <span class="account-tools-menu-icon bg-violet-50 text-violet-600 dark:bg-violet-900/30 dark:text-violet-300">
                        <Icon name="download" size="sm" />
                      </span>
                      <span class="flex-1 text-left">
                        {{ selIds.length ? t('admin.accounts.dataExportSelected') : t('admin.accounts.dataExport') }}
                      </span>
                      <span
                        v-if="selIds.length"
                        class="rounded-full bg-primary-100 px-2 py-0.5 text-xs font-medium text-primary-700 dark:bg-primary-900/40 dark:text-primary-300"
                      >
                        {{ t('admin.accounts.selectedCount', { count: selIds.length }) }}
                      </span>
                    </button>

                    <div class="my-2 border-t border-gray-100 dark:border-gray-700"></div>
                    <div class="px-2 py-2">
                      <div class="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                        {{ t('admin.accounts.toolActions') }}
                      </div>
                    </div>
                    <button class="account-tools-menu-item" @click="openErrorPassthrough">
                      <span class="account-tools-menu-icon bg-amber-50 text-amber-600 dark:bg-amber-900/30 dark:text-amber-300">
                        <Icon name="shield" size="sm" />
                      </span>
                      <span class="flex-1 text-left">{{ t('admin.errorPassthrough.title') }}</span>
                    </button>
                    <button class="account-tools-menu-item" @click="openTLSFingerprintProfiles">
                      <span class="account-tools-menu-icon bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-200">
                        <Icon name="lock" size="sm" />
                      </span>
                      <span class="flex-1 text-left">{{ t('admin.tlsFingerprintProfiles.title') }}</span>
                    </button>

                    <div class="my-2 border-t border-gray-100 dark:border-gray-700"></div>
                    <div class="px-2 py-2">
                      <div class="flex items-center justify-between gap-3">
                        <span class="text-xs font-semibold uppercase tracking-wide text-gray-400 dark:text-gray-500">
                          {{ t('admin.accounts.viewColumns') }}
                        </span>
                        <Icon name="grid" size="sm" class="text-gray-400" />
                      </div>
                    </div>
                    <div class="grid grid-cols-1 gap-1">
                      <button
                        v-for="col in toggleableColumns"
                        :key="col.key"
                        @click="toggleColumn(col.key)"
                        class="flex w-full items-center justify-between rounded-md px-3 py-2 text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
                      >
                        <span class="truncate">{{ col.label }}</span>
                        <Icon v-if="isColumnVisible(col.key)" name="check" size="sm" class="text-primary-500" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </AccountTableActions>
        </div>
        <div
          v-if="hasPendingListSync"
          class="mt-2 flex items-center justify-between rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-800 dark:border-amber-700/40 dark:bg-amber-900/20 dark:text-amber-200"
        >
          <span>{{ t('admin.accounts.listPendingSyncHint') }}</span>
          <button
            class="btn btn-secondary px-2 py-1 text-xs"
            @click="syncPendingListChanges"
          >
            {{ t('admin.accounts.listPendingSyncAction') }}
          </button>
        </div>
      </template>
      <template #table>
        <AccountBulkActionsBar
          :selected-ids="selIds"
          :syncing-upstream-rate="syncingRateBatch"
          @delete="handleBulkDelete"
          @reset-status="handleBulkResetStatus"
          @refresh-token="handleBulkRefreshToken"
          @sync-upstream-rate="handleBulkSyncUpstreamRate"
          @edit-selected="openBulkEditSelected"
          @edit-filtered="openBulkEditFiltered"
          @clear="clearSelection"
          @select-page="selectPage"
          @toggle-schedulable="handleBulkToggleSchedulable"
        />
        <div
          data-test="account-anomaly-summary"
          class="mb-3 flex flex-wrap items-center gap-2 rounded-xl border border-[#d9d1c3] bg-[#faf7ef] px-3 py-2 text-sm text-[#5e5447] dark:border-[#4a4336] dark:bg-[#1f1d19] dark:text-[#d6cdbf]"
        >
          <span class="text-xs font-semibold uppercase tracking-[0.18em] text-[#8d7d66] dark:text-[#a99a82]">
            {{ t('admin.accounts.anomalySummaryTitle') }}
          </span>
          <button
            v-for="item in accountAnomalySummaryItems"
            :key="item.code"
            type="button"
            data-test="account-anomaly-chip"
            class="inline-flex items-center gap-1 rounded-full border px-2.5 py-1 text-xs transition-colors"
            :class="item.code === params.anomaly_reason
              ? 'border-[#a73a2a] bg-[#a73a2a] text-white dark:border-[#cf7966] dark:bg-[#cf7966] dark:text-[#1a120f]'
              : 'border-[#d3c7b5] bg-white/80 text-[#6a5f51] hover:border-[#b8a389] hover:text-[#3f352b] dark:border-[#5a5244] dark:bg-[#26231e] dark:text-[#d9cfbe] dark:hover:border-[#8d7d66]'"
            @click="applyAnomalyReasonFilter(item.code)"
          >
            <span>{{ anomalyReasonLabel(item.code) }}</span>
            <strong class="font-mono text-[11px]">{{ item.count }}</strong>
          </button>
          <span v-if="!accountAnomalySummaryItems.length" class="text-xs text-[#8d7d66] dark:text-[#a99a82]">
            {{ t('admin.accounts.anomalySummaryEmpty') }}
          </span>
          <button
            v-if="params.anomaly_reason"
            type="button"
            class="ml-auto text-xs text-[#8f3e30] transition-colors hover:text-[#742f24] dark:text-[#d08a79] dark:hover:text-[#efb2a5]"
            @click="applyAnomalyReasonFilter('')"
          >
            {{ t('admin.accounts.clearAnomalyFilter') }}
          </button>
        </div>
        <div ref="accountTableRef" class="account-table-shell flex min-h-0 min-w-0 flex-1 flex-col">
          <DataTable
            ref="dataTableRef"
            class="account-table-grid"
            :style="accountTableViewportStyle"
            :columns="cols"
            :data="accounts"
            :loading="loading"
            :sticky-actions-column="false"
            row-key="id"
            :server-side-sort="true"
            @sort="handleSort"
            default-sort-key="name"
            default-sort-order="asc"
            :sort-storage-key="ACCOUNT_SORT_STORAGE_KEY"
            :estimate-row-height="72"
            :overscan="5"
          >
          <template #header-select>
            <input
              type="checkbox"
              class="h-4 w-4 cursor-pointer rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              :checked="allVisibleSelected"
              @click.stop
              @change="toggleSelectAllVisible($event)"
            />
          </template>
          <template #cell-select="{ row }">
            <input type="checkbox" :checked="isSelected(row.id)" @change="toggleSel(row.id)" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
          </template>
          <template #cell-id="{ value }">
            <span class="font-mono text-xs text-gray-500 dark:text-gray-400">#{{ value }}</span>
          </template>
          <template #cell-name="{ row, value }">
            <div class="flex w-full min-w-0 max-w-full flex-col">
              <span class="block w-full truncate font-medium text-gray-900 dark:text-white" :title="String(value || '')">{{ value }}</span>
              <span
                v-if="row.extra?.email_address || row.extra?.email || row.credentials?.email"
                class="block w-full truncate text-xs text-gray-500 dark:text-gray-400"
                :title="String(row.extra?.email_address || row.extra?.email || row.credentials?.email)"
              >
                {{ row.extra?.email_address || row.extra?.email || row.credentials?.email }}
              </span>
            </div>
          </template>
          <template #cell-notes="{ value }">
            <span v-if="value" :title="value" class="block max-w-xs truncate text-sm text-gray-600 dark:text-gray-300">{{ value }}</span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
          </template>
          <template #cell-platform_type="{ row }">
            <div class="flex min-w-0 items-center">
              <div class="flex items-center gap-1 overflow-hidden">
                <PlatformTypeBadge compact :platform="row.platform" :type="row.type" :plan-type="row.credentials?.plan_type" :privacy-mode="row.extra?.privacy_mode" :subscription-expires-at="row.credentials?.subscription_expires_at" />
                <span
                  v-if="getAntigravityTierLabel(row)"
                  :class="['inline-block rounded px-1.5 py-0.5 text-[10px] font-medium', getAntigravityTierClass(row)]"
                >
                  {{ getAntigravityTierLabel(row) }}
                </span>
              </div>
            </div>
          </template>
          <template #cell-capacity="{ row }">
            <AccountCapacityCell :account="row" />
          </template>
          <template #cell-status="{ row }">
            <div class="flex items-center" :title="accountHealthTitle(row)">
              <AccountStatusIndicator compact :account="row" @show-temp-unsched="handleShowTempUnsched" />
            </div>
          </template>
          <template #cell-schedulable="{ row }">
            <button @click="handleToggleSchedulable(row)" :disabled="togglingSchedulable === row.id" class="relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 dark:focus:ring-offset-dark-800" :class="[row.schedulable ? 'bg-primary-500 hover:bg-primary-600' : 'bg-gray-200 hover:bg-gray-300 dark:bg-dark-600 dark:hover:bg-dark-500']" :title="row.schedulable ? t('admin.accounts.schedulableEnabled') : t('admin.accounts.schedulableDisabled')">
              <span class="pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out" :class="[row.schedulable ? 'translate-x-4' : 'translate-x-0']" />
            </button>
          </template>
          <template #cell-today_stats="{ row }">
            <AccountTodayStatsCell
              :stats="todayStatsByAccountId[String(row.id)] ?? null"
              :loading="todayStatsLoading"
              :error="todayStatsError"
            />
          </template>
          <template #cell-groups="{ row }">
            <AccountGroupsCell :groups="row.groups" :max-display="4" />
          </template>
          <template #header-usage="{ column }">
            <div class="flex items-center">
              <span>{{ column.label }}</span>
              <HelpTooltip :content="t('admin.accounts.usageWindowsHint')" width-class="w-72" />
            </div>
          </template>
          <template #cell-usage="{ row }">
            <AccountUsageCell
              :account="row"
              :today-stats="todayStatsByAccountId[String(row.id)] ?? null"
              :today-stats-loading="todayStatsLoading"
              :manual-refresh-token="usageManualRefreshToken"
              compact
            />
          </template>
          <template #cell-proxy="{ row }">
            <div class="flex flex-col gap-1">
              <div v-if="row.proxy" class="flex items-center gap-2">
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ row.proxy.name }}</span>
                <span v-if="row.proxy.country_code" class="text-xs text-gray-500 dark:text-gray-400">
                  ({{ row.proxy.country_code }})
                </span>
              </div>
              <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
              <div v-if="row.proxy && row.proxy.expires_at" class="flex items-center gap-2 text-xs">
                <span class="text-gray-600 dark:text-gray-300">{{ formatDateTime(row.proxy.expires_at) }}</span>
                <span :class="proxyExpiryBadge(row.proxy)">{{ proxyExpiryText(row.proxy) }}</span>
              </div>
              <div v-if="row.proxy_fallback_origin_id" class="flex items-center gap-1">
                <span class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200" :title="t('admin.accounts.fallbackActiveTip', { origin: row.proxy_fallback_origin_name })">
                  {{ t('admin.accounts.fallbackActive') }}
                </span>
                <button class="text-xs px-1.5 py-0.5 rounded border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700" @click="onRevertFallback(row)">{{ t('admin.accounts.revertProxy') }}</button>
              </div>
            </div>
          </template>
          <template #cell-rate_multiplier="{ row }">
            <button
              v-if="canSyncUpstreamRate(row)"
              type="button"
              class="rounded px-1 py-0.5 font-mono text-sm text-teal-700 transition-colors hover:bg-teal-50 hover:text-teal-800 dark:text-teal-300 dark:hover:bg-teal-900/20"
              :title="t('admin.accounts.viewCostRateHistory')"
              @click="openCostRateHistory(row)"
            >
              {{ (row.rate_multiplier ?? 1).toFixed(2) }}x
            </button>
            <span v-else class="text-sm font-mono text-gray-700 dark:text-gray-300">
              {{ (row.rate_multiplier ?? 1).toFixed(2) }}x
            </span>
          </template>
          <template #cell-priority="{ value }">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value }}</span>
          </template>
          <template #cell-last_used_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatRelativeTime(value) }}</span>
          </template>
          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateOnly(value) }}</span>
          </template>
          <template #cell-expires_at="{ row, value }">
            <div class="flex flex-col items-start gap-1">
              <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatExpiresAt(value) }}</span>
              <div v-if="isExpired(value) || (row.auto_pause_on_expired && value)" class="flex items-center gap-1">
                <span
                  v-if="isExpired(value)"
                  class="inline-flex items-center rounded-md bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
                >
                  {{ t('admin.accounts.expired') }}
                </span>
                <span
                  v-if="row.auto_pause_on_expired && value"
                  class="inline-flex items-center rounded-md bg-emerald-100 px-2 py-0.5 text-xs font-medium text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300"
                >
                  {{ t('admin.accounts.autoPauseOnExpired') }}
                </span>
              </div>
            </div>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex items-center gap-0.5 whitespace-nowrap">
              <button @click="handleEdit(row)" :title="t('common.edit')" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" /></svg>
              </button>
              <button @click="handleDelete(row)" :title="t('common.delete')" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" /></svg>
              </button>
              <button @click="handleTest(row)" :title="t('admin.accounts.testAccountConnection')" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-emerald-600 transition-colors hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-900/20">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M5 5v14l11-7-11-7z" /></svg>
              </button>
              <button @click="handleViewStats(row)" :title="t('admin.accounts.usageStatistics')" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-indigo-600 transition-colors hover:bg-indigo-50 dark:text-indigo-400 dark:hover:bg-indigo-900/20">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M7.5 16.5v-9m4.5 9v-6m4.5 6V4.5m-12 12h15" /></svg>
              </button>
              <button @click="handleSchedule(row)" :title="t('admin.scheduledTests.schedule')" class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-amber-600 transition-colors hover:bg-amber-50 dark:text-amber-300 dark:hover:bg-amber-900/20">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6l4 2m5-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              </button>
              <button
                v-if="canSyncUpstreamRate(row)"
                @click="handleSyncUpstreamRate(row)"
                :title="t('admin.accounts.syncUpstreamRateMultiplier')"
                :disabled="syncingRateAccountId === row.id"
                class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-teal-600 transition-colors hover:bg-teal-50 disabled:cursor-not-allowed disabled:opacity-50 dark:text-teal-300 dark:hover:bg-teal-900/20"
              >
                <svg v-if="syncingRateAccountId === row.id" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" /><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v3a5 5 0 00-5 5H4z" /></svg>
                <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v3m0 12v3m7.8-9h-3M7.2 12h-3m13.3-5.5-2.1 2.1M8.6 15.4l-2.1 2.1m0-11 2.1 2.1m6.8 6.8 2.1 2.1" /></svg>
              </button>
            </div>
          </template>
          </DataTable>
        </div>
      </template>
      <template #pagination><Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="handlePageChange" @update:pageSize="handlePageSizeChange" /></template>
    </TablePageLayout>
    </div>
    <CreateAccountModal :show="showCreate" :proxies="proxies" :groups="groups" :quick-flow="createQuickFlow" @close="closeCreateModal" @created="reload" />
    <EditAccountModal :show="showEdit" :account="edAcc" :proxies="proxies" :groups="groups" @close="showEdit = false" @updated="handleAccountUpdated" />
    <ReAuthAccountModal :show="showReAuth" :account="reAuthAcc" @close="closeReAuthModal" @reauthorized="handleAccountUpdated" />
    <AccountTestModal :show="showTest" :account="testingAcc" @close="closeTestModal" />
    <AccountStatsModal :show="showStats" :account="statsAcc" @close="closeStatsModal" />
    <ScheduledTestsPanel :show="showSchedulePanel" :account-id="scheduleAcc?.id ?? null" :model-options="scheduleModelOptions" @close="closeSchedulePanel" />
    <SyncFromCrsModal :show="showSync" @close="showSync = false" @synced="reload" />
    <BaseDialog
      :show="showImportChooser"
      :title="t('admin.accounts.importChooserTitle')"
      width="wide"
      @close="showImportChooser = false"
    >
      <div class="grid gap-3 md:grid-cols-2">
        <button
          type="button"
          class="group rounded-xl border border-gray-200 bg-white p-4 text-left transition hover:border-emerald-300 hover:bg-emerald-50/60 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-emerald-700 dark:hover:bg-emerald-900/20"
          @click="openImportDataFromChooser"
        >
          <span class="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-300">
            <Icon name="upload" size="sm" />
          </span>
          <span class="block text-sm font-semibold text-gray-900 dark:text-white">
            {{ t('admin.accounts.importChooserDataTitle') }}
          </span>
          <span class="mt-1 block text-sm leading-6 text-gray-500 dark:text-gray-400">
            {{ t('admin.accounts.importChooserDataDesc') }}
          </span>
        </button>
        <button
          type="button"
          class="group rounded-xl border border-gray-200 bg-white p-4 text-left transition hover:border-blue-300 hover:bg-blue-50/60 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-blue-700 dark:hover:bg-blue-900/20"
          @click="openOpenAICodexImportFromChooser"
        >
          <span class="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-600 dark:bg-blue-900/30 dark:text-blue-300">
            <Icon name="key" size="sm" />
          </span>
          <span class="block text-sm font-semibold text-gray-900 dark:text-white">
            {{ t('admin.accounts.importChooserOpenAICodexTitle') }}
          </span>
          <span class="mt-1 block text-sm leading-6 text-gray-500 dark:text-gray-400">
            {{ t('admin.accounts.importChooserOpenAICodexDesc') }}
          </span>
        </button>
      </div>
      <p class="mt-4 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs leading-5 text-amber-800 dark:border-amber-800/40 dark:bg-amber-900/20 dark:text-amber-200">
        {{ t('admin.accounts.importChooserSafetyHint') }}
      </p>
    </BaseDialog>
    <ImportDataModal :show="showImportData" @close="showImportData = false" @imported="handleDataImported" />
    <BulkEditAccountModal
      :show="showBulkEdit"
      :account-ids="selIds"
      :selected-platforms="selPlatforms"
      :selected-types="selTypes"
      :target="bulkEditTarget ?? undefined"
      :proxies="proxies"
      :groups="groups"
      @close="showBulkEdit = false"
      @updated="handleBulkUpdated"
    />
    <TempUnschedStatusModal :show="showTempUnsched" :account="tempUnschedAcc" @close="showTempUnsched = false" @reset="handleTempUnschedReset" />
    <ConfirmDialog :show="showDeleteDialog" :title="t('admin.accounts.deleteAccount')" :message="t('admin.accounts.deleteConfirm', { name: deletingAcc?.name })" :confirm-text="t('common.delete')" :cancel-text="t('common.cancel')" :danger="true" @confirm="confirmDelete" @cancel="showDeleteDialog = false" />
    <ConfirmDialog
      :show="showSyncAllUpstreamRatesDialog"
      :title="t('admin.accounts.syncAllUpstreamRatesConfirmTitle')"
      :message="t('admin.accounts.syncAllUpstreamRatesConfirmMessage')"
      :confirm-text="t('admin.accounts.syncAllUpstreamRatesConfirm')"
      :cancel-text="t('common.cancel')"
      @confirm="handleSyncAllUpstreamRates"
      @cancel="showSyncAllUpstreamRatesDialog = false"
    />
	<BaseDialog
	  :show="showUpstreamRateSyncResult"
	  :title="t('admin.accounts.upstreamRateSyncResultTitle')"
	  width="wide"
	  @close="showUpstreamRateSyncResult = false"
	>
	  <div v-if="lastUpstreamRateSyncResult" class="space-y-4">
		<div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
		  <div class="rounded-lg border border-gray-200 p-3 dark:border-gray-700">
			<p class="text-xs text-gray-500">{{ t('admin.accounts.upstreamRateSyncResultTotal') }}</p>
			<p class="mt-1 text-lg font-semibold">{{ lastUpstreamRateSyncResult.total }}</p>
		  </div>
		  <div class="rounded-lg border border-gray-200 p-3 dark:border-gray-700">
			<p class="text-xs text-gray-500">{{ t('admin.accounts.upstreamRateSyncResultChanged') }}</p>
			<p class="mt-1 text-lg font-semibold text-emerald-600">{{ upstreamRateSyncChangedCount }}</p>
		  </div>
		  <div class="rounded-lg border border-gray-200 p-3 dark:border-gray-700">
			<p class="text-xs text-gray-500">{{ t('admin.accounts.upstreamRateSyncResultSignificant') }}</p>
			<p class="mt-1 text-lg font-semibold text-amber-600">{{ upstreamRateSyncSignificantCount }}</p>
		  </div>
		  <div class="rounded-lg border border-gray-200 p-3 dark:border-gray-700">
			<p class="text-xs text-gray-500">{{ t('admin.accounts.upstreamRateSyncResultFailed') }}</p>
			<p class="mt-1 text-lg font-semibold" :class="lastUpstreamRateSyncResult.failed > 0 ? 'text-rose-600' : 'text-gray-900 dark:text-white'">{{ lastUpstreamRateSyncResult.failed }}</p>
		  </div>
		</div>
		<div class="max-h-[26rem] overflow-auto rounded-lg border border-gray-200 dark:border-gray-700">
		  <table class="w-full text-sm">
			<thead class="sticky top-0 bg-gray-50 text-left text-xs text-gray-500 dark:bg-gray-800">
			  <tr>
				<th class="px-3 py-2">{{ t('admin.accounts.upstreamRateSyncResultAccount') }}</th>
				<th class="px-3 py-2">{{ t('admin.accounts.upstreamRateSyncResultChange') }}</th>
				<th class="px-3 py-2">{{ t('admin.accounts.upstreamRateSyncResultStatus') }}</th>
			  </tr>
			</thead>
			<tbody>
			  <tr v-for="item in lastUpstreamRateSyncResult.results" :key="item.account_id" class="border-t border-gray-100 dark:border-gray-700">
				<td class="px-3 py-2">
				  <span class="font-medium text-gray-900 dark:text-white">{{ upstreamRateSyncAccountLabel(item) }}</span>
				  <span class="ml-1 text-xs text-gray-400">#{{ item.account_id }}</span>
				</td>
				<td class="px-3 py-2 font-mono text-xs">
				  <template v-if="item.success">
					{{ formatRateMultiplier(item.previous_rate_multiplier) }} → {{ formatRateMultiplier(item.rate_multiplier) }}
				  </template>
				  <span v-else class="text-gray-400">—</span>
				</td>
				<td class="px-3 py-2">
				  <span v-if="!item.success" class="text-rose-600">{{ item.error || t('admin.accounts.syncUpstreamRateMultiplierFailed') }}</span>
				  <span v-else-if="item.significant_change" class="text-amber-600">{{ t('admin.accounts.upstreamRateSyncResultLargeChange') }}</span>
				  <span v-else-if="item.changed" class="text-emerald-600">{{ t('admin.accounts.upstreamRateSyncResultUpdated') }}</span>
				  <span v-else class="text-gray-500">{{ t('admin.accounts.upstreamRateSyncResultUnchanged') }}</span>
				</td>
			  </tr>
			</tbody>
		  </table>
		</div>
	  </div>
	  <template #footer>
		<button type="button" class="btn btn-primary" @click="showUpstreamRateSyncResult = false">{{ t('common.close') }}</button>
	  </template>
	</BaseDialog>
    <BaseDialog
      :show="showCostRateHistory"
      :title="t('admin.accounts.costRateHistoryTitle', { name: costRateHistoryAccount?.name || '' })"
      width="wide"
      @close="showCostRateHistory = false"
    >
      <div class="space-y-4">
        <div class="flex flex-wrap items-center justify-between gap-2 rounded-lg bg-gray-50 px-3 py-2 dark:bg-gray-800">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ t('admin.accounts.currentCostRate') }}</span>
          <span class="font-mono text-sm font-semibold text-gray-900 dark:text-white">
            {{ formatRateMultiplier(costRateHistoryAccount?.rate_multiplier) }}
          </span>
        </div>
        <div v-if="costRateHistoryLoading" class="py-10 text-center text-sm text-gray-500">
          {{ t('common.loading') }}
        </div>
        <div v-else-if="costRateHistoryError" class="rounded-lg border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700 dark:border-rose-900/50 dark:bg-rose-900/20 dark:text-rose-300">
          <p>{{ costRateHistoryError }}</p>
          <button type="button" class="mt-3 font-medium underline" @click="loadCostRateHistory">{{ t('common.retry') }}</button>
        </div>
        <div v-else-if="costRateHistoryLogs.length === 0" class="py-10 text-center">
          <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('admin.accounts.costRateHistoryEmpty') }}</p>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.accounts.costRateHistoryEmptyHint') }}</p>
        </div>
        <div v-else class="max-h-[28rem] overflow-auto rounded-lg border border-gray-200 dark:border-gray-700">
          <table class="w-full text-sm">
            <thead class="sticky top-0 bg-gray-50 text-left text-xs text-gray-500 dark:bg-gray-800 dark:text-gray-400">
              <tr>
                <th class="px-3 py-2">{{ t('admin.accounts.costRateHistoryTime') }}</th>
                <th class="px-3 py-2">{{ t('admin.accounts.upstreamRateSyncResultChange') }}</th>
                <th class="px-3 py-2">{{ t('admin.accounts.costRateHistoryDelta') }}</th>
                <th class="px-3 py-2">{{ t('admin.accounts.costRateHistorySource') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="log in costRateHistoryLogs" :key="log.id" class="border-t border-gray-100 dark:border-gray-700">
                <td class="whitespace-nowrap px-3 py-2 text-gray-600 dark:text-gray-300">{{ formatDateTime(log.created_at) }}</td>
                <td class="whitespace-nowrap px-3 py-2 font-mono text-xs text-gray-900 dark:text-white">
                  {{ formatRateMultiplier(costRateLogNumber(log, 'previous_rate_multiplier')) }} →
                  {{ formatRateMultiplier(costRateLogNumber(log, 'rate_multiplier')) }}
                </td>
                <td class="whitespace-nowrap px-3 py-2" :class="costRateDeltaClass(log)">
                  {{ formatCostRateDelta(log) }}
                  <span v-if="log.extra?.significant_change" class="ml-1 rounded bg-amber-100 px-1.5 py-0.5 text-[10px] text-amber-700 dark:bg-amber-900/40 dark:text-amber-300">
                    {{ t('admin.accounts.costRateHistoryLargeChange') }}
                  </span>
                </td>
                <td class="px-3 py-2 text-xs text-gray-500 dark:text-gray-400">{{ String(log.extra?.rate_source || '—') }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.accounts.costRateHistoryHint') }}</p>
      </div>
      <template #footer>
        <button type="button" class="btn btn-secondary" :disabled="costRateHistoryLoading" @click="loadCostRateHistory">
          {{ t('common.refresh') }}
        </button>
        <button type="button" class="btn btn-primary" @click="showCostRateHistory = false">{{ t('common.close') }}</button>
      </template>
    </BaseDialog>
    <ConfirmDialog :show="showExportDataDialog" :title="t('admin.accounts.dataExport')" :message="t('admin.accounts.dataExportConfirmMessage')" :confirm-text="t('admin.accounts.dataExportConfirm')" :cancel-text="t('common.cancel')" @confirm="handleExportData" @cancel="showExportDataDialog = false">
      <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
        <input type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" v-model="includeProxyOnExport" />
        <span>{{ t('admin.accounts.dataExportIncludeProxies') }}</span>
      </label>
    </ConfirmDialog>
    <ErrorPassthroughRulesModal :show="showErrorPassthrough" @close="showErrorPassthrough = false" />
    <TLSFingerprintProfilesModal :show="showTLSFingerprintProfiles" @close="showTLSFingerprintProfiles = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, toRaw, watch } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { adminAPI } from '@/api/admin'
import { useTableLoader } from '@/composables/useTableLoader'
import { useSwipeSelect, type SwipeSelectVirtualContext } from '@/composables/useSwipeSelect'
import { useTableSelection } from '@/composables/useTableSelection'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import type { Column } from '@/components/common/types'
import BaseDialog from '@/components/common/BaseDialog.vue'
import HelpTooltip from '@/components/common/HelpTooltip.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { CreateAccountModal, EditAccountModal, BulkEditAccountModal, SyncFromCrsModal, TempUnschedStatusModal } from '@/components/account'
import AccountTableActions from '@/components/admin/account/AccountTableActions.vue'
import AccountTableFilters from '@/components/admin/account/AccountTableFilters.vue'
import AccountBulkActionsBar from '@/components/admin/account/AccountBulkActionsBar.vue'
import ImportDataModal from '@/components/admin/account/ImportDataModal.vue'
import ReAuthAccountModal from '@/components/admin/account/ReAuthAccountModal.vue'
import AccountTestModal from '@/components/admin/account/AccountTestModal.vue'
import AccountStatsModal from '@/components/admin/account/AccountStatsModal.vue'
import ScheduledTestsPanel from '@/components/admin/account/ScheduledTestsPanel.vue'
import type { SelectOption } from '@/components/common/Select.vue'
import AccountStatusIndicator from '@/components/account/AccountStatusIndicator.vue'
import AccountUsageCell from '@/components/account/AccountUsageCell.vue'
import AccountTodayStatsCell from '@/components/account/AccountTodayStatsCell.vue'
import AccountGroupsCell from '@/components/account/AccountGroupsCell.vue'
import AccountCapacityCell from '@/components/account/AccountCapacityCell.vue'
import PlatformTypeBadge from '@/components/common/PlatformTypeBadge.vue'
import Icon from '@/components/icons/Icon.vue'
import ErrorPassthroughRulesModal from '@/components/admin/ErrorPassthroughRulesModal.vue'
import TLSFingerprintProfilesModal from '@/components/admin/TLSFingerprintProfilesModal.vue'
import { accountAnomalyReasonFilterOrder, deriveAccountAnomalyReason, type AccountAnomalyReasonCode } from '@/utils/accountAnomaly'
import { buildOpenAIUsageRefreshKey } from '@/utils/accountUsageRefresh'
import { formatDateOnly, formatDateTime, formatRelativeTime } from '@/utils/format'
import { proxyExpiryBadgeClass, proxyExpiryLabelKey } from '@/utils/proxyExpiry'
import type { Account, AccountPlatform, AccountType, Proxy as AccountProxy, AdminGroup, WindowStats, ClaudeModel } from '@/types'
import type { BatchSyncUpstreamRateMultiplierResult } from '@/api/admin/accounts'
import type { OpsSystemLog } from '@/api/admin/ops'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const proxies = ref<AccountProxy[]>([])
const groups = ref<AdminGroup[]>([])
const accountTableRef = ref<HTMLElement | null>(null)
const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null)
type AccountBulkEditTarget =
  | {
      mode: 'selected'
      accountIds: number[]
      selectedPlatforms: AccountPlatform[]
      selectedTypes: AccountType[]
    }
  | {
      mode: 'filtered'
      filters: {
        platform?: string
        type?: string
        status?: string
        anomaly_reason?: string
        group?: string
        search?: string
        privacy_mode?: string
        sort_by?: string
        sort_order?: AccountSortOrder
      }
      previewCount: number
      selectedPlatforms: AccountPlatform[]
      selectedTypes: AccountType[]
    }
const selPlatforms = computed<AccountPlatform[]>(() => {
  const platforms = new Set(
    accounts.value
      .filter(a => isSelected(a.id))
      .map(a => a.platform)
  )
  return [...platforms]
})
const selTypes = computed<AccountType[]>(() => {
  const types = new Set(
    accounts.value
      .filter(a => isSelected(a.id))
      .map(a => a.type)
  )
  return [...types]
})
const showCreate = ref(false)
const createQuickFlow = ref<'openai-codex-import' | null>(null)
const showEdit = ref(false)
const showSync = ref(false)
const showImportChooser = ref(false)
const showImportData = ref(false)
const showExportDataDialog = ref(false)
const includeProxyOnExport = ref(true)
const showBulkEdit = ref(false)
const bulkEditTarget = ref<AccountBulkEditTarget | null>(null)
const showTempUnsched = ref(false)
const showDeleteDialog = ref(false)
const showSyncAllUpstreamRatesDialog = ref(false)
const showUpstreamRateSyncResult = ref(false)
const showCostRateHistory = ref(false)
const showReAuth = ref(false)
const showTest = ref(false)
const showStats = ref(false)
const showErrorPassthrough = ref(false)
const showTLSFingerprintProfiles = ref(false)
const edAcc = ref<Account | null>(null)
const tempUnschedAcc = ref<Account | null>(null)
const deletingAcc = ref<Account | null>(null)
const reAuthAcc = ref<Account | null>(null)
const testingAcc = ref<Account | null>(null)
const statsAcc = ref<Account | null>(null)
const showSchedulePanel = ref(false)
const scheduleAcc = ref<Account | null>(null)
const scheduleModelOptions = ref<SelectOption[]>([])
const togglingSchedulable = ref<number | null>(null)
const syncingRateAccountId = ref<number | null>(null)
const syncingRateBatch = ref(false)
const syncingAllUpstreamRates = ref(false)
const failedUpstreamRateAccountIds = ref<number[]>([])
const lastUpstreamRateSyncResult = ref<BatchSyncUpstreamRateMultiplierResult | null>(null)
const costRateHistoryAccount = ref<Account | null>(null)
const costRateHistoryLogs = ref<OpsSystemLog[]>([])
const costRateHistoryLoading = ref(false)
const costRateHistoryError = ref('')
const scheduledCostRateSyncEnabled = ref(false)
const scheduledCostRateSyncLoading = ref(true)
const scheduledCostRateSyncSaving = ref(false)

const upstreamRateSyncChangedCount = computed(() => (
  lastUpstreamRateSyncResult.value?.results.filter(item => item.success && item.changed).length || 0
))
const upstreamRateSyncSignificantCount = computed(() => (
  lastUpstreamRateSyncResult.value?.results.filter(item => item.success && item.significant_change).length || 0
))
const formatRateMultiplier = (value?: number) => `${(value ?? 1).toFixed(4)}x`
const upstreamRateSyncAccountLabel = (item: BatchSyncUpstreamRateMultiplierResult['results'][number]) => (
  item.account?.name || item.account_name || t('admin.accounts.upstreamRateSyncResultUnknownAccount')
)
const recordUpstreamRateSyncResult = (result: BatchSyncUpstreamRateMultiplierResult, open = true) => {
  lastUpstreamRateSyncResult.value = result
  showUpstreamRateSyncResult.value = open
  if (upstreamRateSyncSignificantCount.value > 0) {
    appStore.showWarning(t('admin.accounts.upstreamRateSyncSignificantWarning', {
      count: upstreamRateSyncSignificantCount.value
    }))
  }
}
const costRateLogNumber = (log: OpsSystemLog, key: string) => {
  const value = Number(log.extra?.[key])
  return Number.isFinite(value) ? value : undefined
}
const formatCostRateDelta = (log: OpsSystemLog) => {
  const ratio = costRateLogNumber(log, 'change_ratio')
  if (ratio === undefined) return '—'
  const sign = ratio > 0 ? '+' : ''
  return `${sign}${(ratio * 100).toFixed(1)}%`
}
const costRateDeltaClass = (log: OpsSystemLog) => {
  const ratio = costRateLogNumber(log, 'change_ratio')
  if (ratio === undefined || Math.abs(ratio) < 1e-9) return 'text-gray-500 dark:text-gray-400'
  return ratio > 0 ? 'text-rose-600 dark:text-rose-300' : 'text-emerald-600 dark:text-emerald-300'
}
const loadCostRateHistory = async () => {
  if (!costRateHistoryAccount.value) return
  costRateHistoryLoading.value = true
  costRateHistoryError.value = ''
  try {
    const result = await adminAPI.ops.listSystemLogs({
      page: 1,
      page_size: 50,
      time_range: '30d',
      component: 'upstream.cost_rate_sync',
      account_id: costRateHistoryAccount.value.id
    })
    costRateHistoryLogs.value = result.items
  } catch (error: any) {
    console.error('Failed to load cost rate history:', error)
    costRateHistoryError.value = error?.message || error?.response?.data?.message || t('admin.accounts.costRateHistoryLoadFailed')
  } finally {
    costRateHistoryLoading.value = false
  }
}
const openCostRateHistory = (account: Account) => {
  costRateHistoryAccount.value = account
  costRateHistoryLogs.value = []
  showCostRateHistory.value = true
  void loadCostRateHistory()
}
const loadScheduledCostRateSyncSetting = async () => {
  scheduledCostRateSyncLoading.value = true
  if (!adminAPI.settings?.getSettings) {
    scheduledCostRateSyncLoading.value = false
    return
  }
  try {
    const settings = await adminAPI.settings.getSettings()
    scheduledCostRateSyncEnabled.value = Boolean(settings.upstream_rate_sync_enabled)
  } catch (error) {
    console.error('Failed to load scheduled cost rate sync setting:', error)
  } finally {
    scheduledCostRateSyncLoading.value = false
  }
}
const toggleScheduledCostRateSync = async () => {
  if (scheduledCostRateSyncLoading.value || scheduledCostRateSyncSaving.value) return
  const next = !scheduledCostRateSyncEnabled.value
  scheduledCostRateSyncSaving.value = true
  try {
    const settings = await adminAPI.settings.updateSettings({ upstream_rate_sync_enabled: next })
    scheduledCostRateSyncEnabled.value = Boolean(settings.upstream_rate_sync_enabled)
    appStore.showSuccess(t(next
      ? 'admin.accounts.scheduledCostRateSyncEnabled'
      : 'admin.accounts.scheduledCostRateSyncDisabled'))
  } catch (error: any) {
    console.error('Failed to update scheduled cost rate sync setting:', error)
    appStore.showError(error?.message || t('admin.accounts.scheduledCostRateSyncUpdateFailed'))
  } finally {
    scheduledCostRateSyncSaving.value = false
  }
}
const exportingData = ref(false)

// Account tools dropdown
const showAccountToolsDropdown = ref(false)
const accountToolsDropdownRef = ref<HTMLElement | null>(null)
const accountToolsDropdownPosition = ref<{ top: number; left: number } | null>(null)
const hiddenColumns = reactive<Set<string>>(new Set())
const DEFAULT_HIDDEN_COLUMNS = ['today_stats', 'proxy', 'notes', 'priority', 'rate_multiplier']
const HIDDEN_COLUMNS_KEY = 'account-hidden-columns'

// Sorting settings
const ACCOUNT_SORT_STORAGE_KEY = 'account-table-sort'
type AccountSortOrder = 'asc' | 'desc'
type AccountSortState = {
  sort_by: string
  sort_order: AccountSortOrder
}
const ACCOUNT_SORTABLE_KEYS = new Set([
  'id',
  'name',
  'status',
  'schedulable',
  'priority',
  'rate_multiplier',
  'last_used_at',
  'created_at',
  'expires_at'
])
const loadInitialAccountSortState = (): AccountSortState => {
  const fallback: AccountSortState = { sort_by: 'name', sort_order: 'asc' }
  try {
    const raw = localStorage.getItem(ACCOUNT_SORT_STORAGE_KEY)
    if (!raw) return fallback
    const parsed = JSON.parse(raw) as { key?: string; order?: string }
    const key = typeof parsed.key === 'string' ? parsed.key : ''
    if (!ACCOUNT_SORTABLE_KEYS.has(key)) return fallback
    return {
      sort_by: key,
      sort_order: parsed.order === 'desc' ? 'desc' : 'asc'
    }
  } catch {
    return fallback
  }
}
const sortState = reactive<AccountSortState>(loadInitialAccountSortState())

// Auto refresh settings
const showAutoRefreshDropdown = ref(false)
const autoRefreshDropdownRef = ref<HTMLElement | null>(null)
const AUTO_REFRESH_STORAGE_KEY = 'account-auto-refresh'
const autoRefreshIntervals = [5, 10, 15, 30] as const
const autoRefreshEnabled = ref(false)
const autoRefreshIntervalSeconds = ref<(typeof autoRefreshIntervals)[number]>(30)
const autoRefreshCountdown = ref(0)
const autoRefreshETag = ref<string | null>(null)
const autoRefreshFetching = ref(false)
const AUTO_REFRESH_SILENT_WINDOW_MS = 15000
const autoRefreshSilentUntil = ref(0)
const hasPendingListSync = ref(false)
const todayStatsByAccountId = ref<Record<string, WindowStats>>({})
const todayStatsLoading = ref(false)
const todayStatsError = ref<string | null>(null)
const todayStatsReqSeq = ref(0)
const pendingTodayStatsRefresh = ref(false)
const usageManualRefreshToken = ref(0)

const buildDefaultTodayStats = (): WindowStats => ({
  requests: 0,
  tokens: 0,
  cost: 0,
  standard_cost: 0,
  user_cost: 0
})

const refreshTodayStatsBatch = async () => {
  // Why this checks both columns:
  // - today_stats column shows dedicated today's metrics.
  // - usage column also embeds today's stats for Key/Bedrock rows.
  // So we only skip fetching when BOTH columns are hidden.
  if (hiddenColumns.has('today_stats') && hiddenColumns.has('usage')) {
    todayStatsLoading.value = false
    todayStatsError.value = null
    return
  }

  const accountIDs = accounts.value.map(account => account.id)
  const reqSeq = ++todayStatsReqSeq.value
  if (accountIDs.length === 0) {
    todayStatsByAccountId.value = {}
    todayStatsError.value = null
    todayStatsLoading.value = false
    return
  }

  todayStatsLoading.value = true
  todayStatsError.value = null

  try {
    const result = await adminAPI.accounts.getBatchTodayStats(accountIDs)
    if (reqSeq !== todayStatsReqSeq.value) return
    const serverStats = result.stats ?? {}
    const nextStats: Record<string, WindowStats> = {}
    for (const accountID of accountIDs) {
      const key = String(accountID)
      nextStats[key] = serverStats[key] ?? buildDefaultTodayStats()
    }
    todayStatsByAccountId.value = nextStats
  } catch (error) {
    if (reqSeq !== todayStatsReqSeq.value) return
    todayStatsError.value = 'Failed'
    console.error('Failed to load account today stats:', error)
  } finally {
    if (reqSeq === todayStatsReqSeq.value) {
      todayStatsLoading.value = false
    }
  }
}

const autoRefreshIntervalLabel = (sec: number) => {
  if (sec === 5) return t('admin.accounts.refreshInterval5s')
  if (sec === 10) return t('admin.accounts.refreshInterval10s')
  if (sec === 15) return t('admin.accounts.refreshInterval15s')
  if (sec === 30) return t('admin.accounts.refreshInterval30s')
  return `${sec}s`
}

const loadSavedColumns = () => {
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      const parsed = JSON.parse(saved) as string[]
      parsed.forEach(key => {
        hiddenColumns.add(key)
      })
    } else {
      DEFAULT_HIDDEN_COLUMNS.forEach(key => {
        hiddenColumns.add(key)
      })
    }
  } catch (e) {
    console.error('Failed to load saved columns:', e)
    DEFAULT_HIDDEN_COLUMNS.forEach(key => {
      hiddenColumns.add(key)
    })
  }
}

const saveColumnsToStorage = () => {
  try {
    localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
  } catch (e) {
    console.error('Failed to save columns:', e)
  }
}

const loadSavedAutoRefresh = () => {
  try {
    const saved = localStorage.getItem(AUTO_REFRESH_STORAGE_KEY)
    if (!saved) return
    const parsed = JSON.parse(saved) as { enabled?: boolean; interval_seconds?: number }
    autoRefreshEnabled.value = parsed.enabled === true
    const interval = Number(parsed.interval_seconds)
    if (autoRefreshIntervals.includes(interval as any)) {
      autoRefreshIntervalSeconds.value = interval as any
    }
  } catch (e) {
    console.error('Failed to load saved auto refresh settings:', e)
  }
}

const saveAutoRefreshToStorage = () => {
  try {
    localStorage.setItem(
      AUTO_REFRESH_STORAGE_KEY,
      JSON.stringify({
        enabled: autoRefreshEnabled.value,
        interval_seconds: autoRefreshIntervalSeconds.value
      })
    )
  } catch (e) {
    console.error('Failed to save auto refresh settings:', e)
  }
}

if (typeof window !== 'undefined') {
  loadSavedColumns()
  loadSavedAutoRefresh()
}

const setAutoRefreshEnabled = (enabled: boolean) => {
  autoRefreshEnabled.value = enabled
  saveAutoRefreshToStorage()
  if (enabled) {
    autoRefreshCountdown.value = autoRefreshIntervalSeconds.value
    resumeAutoRefresh()
  } else {
    pauseAutoRefresh()
    autoRefreshCountdown.value = 0
  }
}

const setAutoRefreshInterval = (seconds: (typeof autoRefreshIntervals)[number]) => {
  autoRefreshIntervalSeconds.value = seconds
  saveAutoRefreshToStorage()
  if (autoRefreshEnabled.value) {
    autoRefreshCountdown.value = seconds
  }
}

const toggleColumn = (key: string) => {
  const wasHidden = hiddenColumns.has(key)
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  } else {
    hiddenColumns.add(key)
  }
  saveColumnsToStorage()
  if ((key === 'today_stats' || key === 'usage') && wasHidden) {
    refreshTodayStatsBatch().catch((error) => {
      console.error('Failed to load account today stats after showing column:', error)
    })
  }
}

const isColumnVisible = (key: string) => !hiddenColumns.has(key)

const {
  items: accounts,
  loading,
  params,
  pagination,
  load: baseLoad,
  reload: baseReload,
  debouncedReload: baseDebouncedReload,
  handlePageChange: baseHandlePageChange,
  handlePageSizeChange: baseHandlePageSizeChange
} = useTableLoader<Account, any>({
  fetchFn: adminAPI.accounts.list,
  pageSize: 10,
  initialParams: {
    platform: '',
    type: '',
    status: '',
    anomaly_reason: '',
    privacy_mode: '',
    group: '',
    search: '',
    sort_by: sortState.sort_by,
    sort_order: sortState.sort_order
  }
})

const {
  selectedIds: selIds,
  allVisibleSelected,
  isSelected,
  setSelectedIds,
  select,
  deselect,
  toggle: toggleSel,
  clear: clearSelection,
  removeMany: removeSelectedAccounts,
  toggleVisible,
  selectVisible: selectPage,
  batchUpdate
} = useTableSelection<Account>({
  rows: accounts,
  getId: (account) => account.id
})

const swipeVirtualContext: SwipeSelectVirtualContext = {
  getVirtualizer: () => dataTableRef.value?.virtualizer ?? null,
  getSortedData: () => dataTableRef.value?.sortedData ?? accounts.value,
  getRowId: (row: any) => row.id,
}

const ACCOUNT_TABLE_ROW_HEIGHT_PX = 84
const ACCOUNT_TABLE_HEADER_HEIGHT_PX = 60
const ACCOUNT_TABLE_SCROLLBAR_HEIGHT_PX = 16

const accountTableViewportStyle = computed(() => {
  const rowCount = Math.max(1, Math.min(accounts.value.length || 0, pagination.page_size || 10))
  const tableHeight = ACCOUNT_TABLE_HEADER_HEIGHT_PX + ACCOUNT_TABLE_SCROLLBAR_HEIGHT_PX + (rowCount * ACCOUNT_TABLE_ROW_HEIGHT_PX)
  return {
    '--account-table-visible-height': `${tableHeight}px`
  } as Record<string, string>
})

useSwipeSelect(accountTableRef, {
  isSelected,
  select,
  deselect,
  batchUpdate
}, swipeVirtualContext)

const resetAutoRefreshCache = () => {
  autoRefreshETag.value = null
}

const isFirstLoad = ref(true)

const load = async () => {
  const requestParams = params as any
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = false
  if (isFirstLoad.value) {
    requestParams.lite = '1'
  }
  await baseLoad()
  if (isFirstLoad.value) {
    isFirstLoad.value = false
    delete requestParams.lite
  }
  await refreshTodayStatsBatch()
}

const reload = async () => {
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = false
  await baseReload()
  await refreshTodayStatsBatch()
}

const debouncedReload = () => {
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = true
  baseDebouncedReload()
}

const handlePageChange = (page: number) => {
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = true
  baseHandlePageChange(page)
}

const handlePageSizeChange = (size: number) => {
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = true
  baseHandlePageSizeChange(size)
}

const handleSort = (key: string, order: AccountSortOrder) => {
  sortState.sort_by = key
  sortState.sort_order = order
  const requestParams = params as any
  requestParams.sort_by = key
  requestParams.sort_order = order
  pagination.page = 1
  hasPendingListSync.value = false
  resetAutoRefreshCache()
  pendingTodayStatsRefresh.value = true
  load()
}

watch(loading, (isLoading, wasLoading) => {
  if (wasLoading && !isLoading && pendingTodayStatsRefresh.value) {
    pendingTodayStatsRefresh.value = false
    refreshTodayStatsBatch().catch((error) => {
      console.error('Failed to refresh account today stats after table load:', error)
    })
  }
})

const isAnyModalOpen = computed(() => {
  return (
    showCreate.value ||
    showEdit.value ||
    showSync.value ||
    showImportChooser.value ||
    showImportData.value ||
    showExportDataDialog.value ||
    showBulkEdit.value ||
    showTempUnsched.value ||
    showDeleteDialog.value ||
    showReAuth.value ||
    showTest.value ||
    showStats.value ||
    showSchedulePanel.value ||
    showErrorPassthrough.value ||
    showTLSFingerprintProfiles.value
  )
})

const enterAutoRefreshSilentWindow = () => {
  autoRefreshSilentUntil.value = Date.now() + AUTO_REFRESH_SILENT_WINDOW_MS
  autoRefreshCountdown.value = autoRefreshIntervalSeconds.value
}

const inAutoRefreshSilentWindow = () => {
  return Date.now() < autoRefreshSilentUntil.value
}

const shouldReplaceAutoRefreshRow = (current: Account, next: Account) => {
  return (
    current.updated_at !== next.updated_at ||
    current.current_concurrency !== next.current_concurrency ||
    current.current_window_cost !== next.current_window_cost ||
    current.active_sessions !== next.active_sessions ||
    current.health?.score !== next.health?.score ||
    current.health?.level !== next.health?.level ||
    current.schedulable !== next.schedulable ||
    current.status !== next.status ||
    current.rate_limit_reset_at !== next.rate_limit_reset_at ||
    current.overload_until !== next.overload_until ||
    current.temp_unschedulable_until !== next.temp_unschedulable_until ||
    buildOpenAIUsageRefreshKey(current) !== buildOpenAIUsageRefreshKey(next)
  )
}

const isFutureTime = (value?: string | null) => {
  if (!value) return false
  return new Date(value).getTime() > Date.now()
}

const resolveAccountHealth = (account: Account) => {
  if (account.health) return account.health
  let score = 100
  const reasons: string[] = []
  let nextAction = ''
  const addReason = (points: number, reason: string, action = '') => {
    score -= points
    reasons.push(reason)
    if (!nextAction && action) nextAction = action
  }

  if (account.status === 'error') {
    addReason(55, '账号处于错误状态', '查看错误信息并执行恢复状态或重新授权')
  } else if (account.status !== 'active') {
    addReason(40, '账号未启用', '启用账号后再参与调度')
  }
  if (!account.schedulable) addReason(35, '已暂停调度', '确认后可重新开启调度')
  if (account.auto_pause_on_expired && account.expires_at && account.expires_at * 1000 <= Date.now()) {
    addReason(45, '账号已过期并自动暂停', '刷新凭证或更新过期时间')
  }
  if (isFutureTime(account.temp_unschedulable_until)) {
    addReason(35, account.temp_unschedulable_reason ? `临时冷却中：${account.temp_unschedulable_reason}` : '临时冷却中', '等待冷却结束，或确认风险后恢复状态')
  }
  if (isFutureTime(account.rate_limit_reset_at)) addReason(30, '上游限流中', '等待限流窗口结束，或切换/补充账号池')
  if (isFutureTime(account.overload_until)) addReason(30, '上游过载冷却中', '等待过载冷却结束，或降低该账号调度权重')

  const quotaExceeded = (used?: number | null, limit?: number | null) =>
    typeof limit === 'number' && limit > 0 && typeof used === 'number' && used >= limit
  if (
    quotaExceeded(account.quota_used, account.quota_limit) ||
    quotaExceeded(account.quota_daily_used, account.quota_daily_limit) ||
    quotaExceeded(account.quota_weekly_used, account.quota_weekly_limit)
  ) {
    addReason(45, '账号配额已耗尽', '重置配额或调整预算上限')
  }

  const currentConcurrency = account.current_concurrency ?? 0
  if (account.concurrency > 0 && currentConcurrency >= account.concurrency) {
    addReason(18, '并发已满', '等待请求完成，或提高并发上限')
  } else if (account.concurrency > 0 && currentConcurrency >= account.concurrency * 0.8) {
    addReason(8, '并发接近上限', '观察请求堆积情况')
  }

  if (typeof account.current_window_cost === 'number' && typeof account.window_cost_limit === 'number' && account.window_cost_limit > 0) {
    const ratio = account.current_window_cost / account.window_cost_limit
    if (ratio >= 1) addReason(28, '5h 费用窗口已达阈值', '等待窗口重置或调整费用阈值')
    else if (ratio >= 0.8) addReason(12, '5h 费用窗口接近阈值', '关注成本水位')
  }

  if (typeof account.active_sessions === 'number' && typeof account.max_sessions === 'number' && account.max_sessions > 0) {
    const ratio = account.active_sessions / account.max_sessions
    if (ratio >= 1) addReason(20, '活跃会话已满', '等待会话空闲释放或提高会话上限')
    else if (ratio >= 0.8) addReason(8, '活跃会话接近上限', '观察会话占用')
  }

  if (typeof account.current_rpm === 'number' && typeof account.base_rpm === 'number' && account.base_rpm > 0) {
    const ratio = account.current_rpm / account.base_rpm
    if (ratio >= 1) addReason(20, 'RPM 已达基础上限', '等待分钟窗口恢复或提高 RPM 上限')
    else if (ratio >= 0.8) addReason(8, 'RPM 接近基础上限', '观察短时请求峰值')
  }

  score = Math.max(0, score)
  const level = score < 50 ? 'critical' : score < 75 ? 'warning' : 'good'
  const label = score < 50 ? '需处理' : score < 75 ? '需留意' : '安稳'
  return {
    score,
    level,
    label,
    reasons: reasons.length ? reasons : ['运行状态正常'],
    next_action: nextAction
  }
}

const accountHealthTitle = (account: Account) => {
  const anomaly = deriveAccountAnomalyReason(account)
  const health = resolveAccountHealth(account)
  if (!health) return ''
  const parts = [`账号脉象：${health.label} · ${health.score}`]
  if (anomaly.code !== 'healthy') {
    parts.push(`异常：${anomalyReasonLabel(anomaly.code)}${anomaly.detail ? `（${anomaly.detail}）` : ''}`)
    parts.push(`建议：${anomalyReasonAction(anomaly.code)}`)
  }
  if (health.reasons?.length) {
    parts.push(`原因：${health.reasons.join('；')}`)
  }
  if (health.next_action) {
    parts.push(`建议：${health.next_action}`)
  }
  return parts.join('\n')
}

const syncAccountRefs = (nextAccount: Account) => {
  if (edAcc.value?.id === nextAccount.id) edAcc.value = nextAccount
  if (reAuthAcc.value?.id === nextAccount.id) reAuthAcc.value = nextAccount
  if (tempUnschedAcc.value?.id === nextAccount.id) tempUnschedAcc.value = nextAccount
  if (deletingAcc.value?.id === nextAccount.id) deletingAcc.value = nextAccount
}

const mergeAccountsIncrementally = (nextRows: Account[]) => {
  const currentRows = accounts.value
  const currentByID = new Map(currentRows.map(row => [row.id, row]))
  let changed = nextRows.length !== currentRows.length
  const mergedRows = nextRows.map((nextRow) => {
    const currentRow = currentByID.get(nextRow.id)
    if (!currentRow) {
      changed = true
      return nextRow
    }
    if (shouldReplaceAutoRefreshRow(currentRow, nextRow)) {
      changed = true
      syncAccountRefs(nextRow)
      return nextRow
    }
    return currentRow
  })
  if (!changed) {
    for (let i = 0; i < mergedRows.length; i += 1) {
      if (mergedRows[i].id !== currentRows[i]?.id) {
        changed = true
        break
      }
    }
  }
  if (changed) {
    accounts.value = mergedRows
  }
}

const refreshAccountsIncrementally = async () => {
  if (autoRefreshFetching.value) return
  autoRefreshFetching.value = true
  try {
    const result = await adminAPI.accounts.listWithEtag(
      pagination.page,
      pagination.page_size,
      toRaw(params) as {
        platform?: string
        type?: string
        status?: string
        anomaly_reason?: string
        privacy_mode?: string
        group?: string
        search?: string
        sort_by?: string
        sort_order?: AccountSortOrder

      },
      { etag: autoRefreshETag.value }
    )

    if (result.etag) {
      autoRefreshETag.value = result.etag
    }
    if (!result.notModified && result.data) {
      pagination.total = result.data.total || 0
      pagination.pages = result.data.pages || 0
      mergeAccountsIncrementally(result.data.items || [])
      hasPendingListSync.value = false
    }

    await refreshTodayStatsBatch()
  } catch (error) {
    console.error('Auto refresh failed:', error)
  } finally {
    autoRefreshFetching.value = false
  }
}

const handleManualRefresh = async () => {
  await load()
  // Force usage cells to refetch /usage on explicit user refresh.
  usageManualRefreshToken.value += 1
}

const openCreateAccount = () => {
  createQuickFlow.value = null
  showCreate.value = true
}

const closeCreateModal = () => {
  showCreate.value = false
  createQuickFlow.value = null
}

const closeAccountToolsDropdown = () => {
  showAccountToolsDropdown.value = false
  accountToolsDropdownPosition.value = null
}

const openAccountToolsDropdown = (event: MouseEvent) => {
  const target = event.currentTarget as HTMLElement | null
  if (!target) {
    showAccountToolsDropdown.value = true
    accountToolsDropdownPosition.value = { top: event.clientY, left: Math.max(16, event.clientX - 240) }
    return
  }

  const rect = target.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  const padding = 16
  const menuWidth = Math.min(320, viewportWidth - padding * 2)
  const estimatedMenuHeight = 440

  const left = Math.max(
    padding,
    Math.min(rect.right - menuWidth, viewportWidth - menuWidth - padding)
  )

  let top = rect.bottom + 8
  if (top + estimatedMenuHeight > viewportHeight - padding) {
    top = Math.max(padding, rect.top - estimatedMenuHeight - 8)
  }

  accountToolsDropdownPosition.value = { top, left }
  showAccountToolsDropdown.value = true
}

const toggleAccountToolsDropdown = (event: MouseEvent) => {
  if (showAccountToolsDropdown.value) {
    closeAccountToolsDropdown()
    return
  }
  openAccountToolsDropdown(event)
}

const openSyncFromCrs = () => {
  closeAccountToolsDropdown()
  showSync.value = true
}

const openImportData = () => {
  closeAccountToolsDropdown()
  showImportChooser.value = true
}

const openImportDataFromChooser = () => {
  showImportChooser.value = false
  showImportData.value = true
}

const openOpenAICodexImportFromChooser = () => {
  showImportChooser.value = false
  createQuickFlow.value = 'openai-codex-import'
  showCreate.value = true
}

const openExportDataDialogFromMenu = () => {
  closeAccountToolsDropdown()
  openExportDataDialog()
}

const openErrorPassthrough = () => {
  closeAccountToolsDropdown()
  showErrorPassthrough.value = true
}

const openTLSFingerprintProfiles = () => {
  closeAccountToolsDropdown()
  showTLSFingerprintProfiles.value = true
}

const syncPendingListChanges = async () => {
  hasPendingListSync.value = false
  await load()
  // Keep behavior consistent with manual refresh.
  usageManualRefreshToken.value += 1
}

const { pause: pauseAutoRefresh, resume: resumeAutoRefresh } = useIntervalFn(
  async () => {
    if (!autoRefreshEnabled.value) return
    if (document.hidden) return
    if (loading.value || autoRefreshFetching.value) return
    if (isAnyModalOpen.value) return
    if (showAccountToolsDropdown.value || showAutoRefreshDropdown.value) return
    if (inAutoRefreshSilentWindow()) {
      autoRefreshCountdown.value = Math.max(
        0,
        Math.ceil((autoRefreshSilentUntil.value - Date.now()) / 1000)
      )
      return
    }

    if (autoRefreshCountdown.value <= 0) {
      autoRefreshCountdown.value = autoRefreshIntervalSeconds.value
      await refreshAccountsIncrementally()
      return
    }

    autoRefreshCountdown.value -= 1
  },
  1000,
  { immediate: false }
)

// Antigravity 订阅等级辅助函数
function getAntigravityTierFromRow(row: any): string | null {
  if (row.platform !== 'antigravity') return null
  const extra = row.extra as Record<string, unknown> | undefined
  if (!extra) return null
  const lca = extra.load_code_assist as Record<string, unknown> | undefined
  if (!lca) return null
  const paid = lca.paidTier as Record<string, unknown> | undefined
  if (paid && typeof paid.id === 'string') return paid.id
  const current = lca.currentTier as Record<string, unknown> | undefined
  if (current && typeof current.id === 'string') return current.id
  return null
}

function getAntigravityTierLabel(row: any): string | null {
  const tier = getAntigravityTierFromRow(row)
  switch (tier) {
    case 'free-tier': return t('admin.accounts.tier.free')
    case 'g1-pro-tier': return t('admin.accounts.tier.pro')
    case 'g1-ultra-tier': return t('admin.accounts.tier.ultra')
    default: return null
  }
}

function getAntigravityTierClass(row: any): string {
  const tier = getAntigravityTierFromRow(row)
  switch (tier) {
    case 'free-tier': return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
    case 'g1-pro-tier': return 'bg-blue-100 text-blue-600 dark:bg-blue-900/40 dark:text-blue-300'
    case 'g1-ultra-tier': return 'bg-purple-100 text-purple-600 dark:bg-purple-900/40 dark:text-purple-300'
    default: return ''
  }
}

const ACCOUNT_COLUMN_CLASS_MAP: Record<string, string> = {
  select: 'w-12 min-w-[3rem] max-w-[3rem]',
  name: 'w-[24ch] min-w-[24ch] max-w-[24ch]',
  id: 'w-[5.5rem] min-w-[5.5rem]',
  platform_type: 'w-[8rem] min-w-[8rem]',
  capacity: 'w-[6rem] min-w-[6rem]',
  status: 'w-[5.5rem] min-w-[5.5rem]',
  schedulable: 'w-[4.75rem] min-w-[4.75rem]',
  today_stats: 'w-[8.5rem] min-w-[8.5rem]',
  groups: 'w-[10rem] min-w-[10rem]',
  usage: 'w-[8.25rem] min-w-[8.25rem]',
  proxy: 'w-[7.75rem] min-w-[7.75rem]',
  priority: 'w-[4rem] min-w-[4rem]',
  rate_multiplier: 'w-[5.25rem] min-w-[5.25rem]',
  last_used_at: 'w-[8rem] min-w-[8rem]',
  created_at: 'w-[6.5rem] min-w-[6.5rem]',
  expires_at: 'w-[8.25rem] min-w-[8.25rem]',
  notes: 'w-[8rem] min-w-[8rem]',
  actions: 'w-[12rem] min-w-[12rem]'
}

const getAccountColumnClass = (key: string) => ACCOUNT_COLUMN_CLASS_MAP[key] ?? ''

// All available columns
const allColumns = computed<Column[]>(() => {
  const c: Column[] = [
    { key: 'select', label: '', sortable: false, class: getAccountColumnClass('select') },
    { key: 'name', label: t('admin.accounts.columns.name'), sortable: true, class: getAccountColumnClass('name') },
    { key: 'id', label: t('admin.accounts.columns.id'), sortable: true, class: getAccountColumnClass('id') },
    { key: 'platform_type', label: t('admin.accounts.columns.platformType'), sortable: false, class: getAccountColumnClass('platform_type') },
    { key: 'capacity', label: t('admin.accounts.columns.capacity'), sortable: false, class: getAccountColumnClass('capacity') },
    { key: 'status', label: t('admin.accounts.columns.status'), sortable: true, class: getAccountColumnClass('status') },
    { key: 'schedulable', label: t('admin.accounts.columns.schedulable'), sortable: true, class: getAccountColumnClass('schedulable') },
    { key: 'today_stats', label: t('admin.accounts.columns.todayStats'), sortable: false, class: getAccountColumnClass('today_stats') }
  ]
  if (!authStore.isSimpleMode) {
    c.push({ key: 'groups', label: t('admin.accounts.columns.groups'), sortable: false, class: getAccountColumnClass('groups') })
  }
  c.push(
    { key: 'usage', label: t('admin.accounts.columns.usageWindows'), sortable: false, class: getAccountColumnClass('usage') },
    { key: 'proxy', label: t('admin.accounts.columns.proxy'), sortable: false, class: getAccountColumnClass('proxy') },
    { key: 'priority', label: t('admin.accounts.columns.priority'), sortable: true, class: getAccountColumnClass('priority') },
    { key: 'rate_multiplier', label: t('admin.accounts.columns.billingRateMultiplier'), sortable: true, class: getAccountColumnClass('rate_multiplier') },
    { key: 'last_used_at', label: t('admin.accounts.columns.lastUsed'), sortable: true, class: getAccountColumnClass('last_used_at') },
    { key: 'created_at', label: t('admin.accounts.columns.createdAt'), sortable: true, class: getAccountColumnClass('created_at') },
    { key: 'expires_at', label: t('admin.accounts.columns.expiresAt'), sortable: true, class: getAccountColumnClass('expires_at') },
    { key: 'notes', label: t('admin.accounts.columns.notes'), sortable: false, class: getAccountColumnClass('notes') },
    { key: 'actions', label: t('admin.accounts.columns.actions'), sortable: false, class: getAccountColumnClass('actions') }
  )
  return c
})

// Columns that can be toggled (exclude select, name, and actions)
const toggleableColumns = computed(() =>
  allColumns.value.filter(col => col.key !== 'select' && col.key !== 'name' && col.key !== 'actions')
)

// Filtered columns based on visibility
const cols = computed(() =>
  allColumns.value.filter(col =>
    col.key === 'select' || col.key === 'name' || col.key === 'actions' || !hiddenColumns.has(col.key)
  )
)

const handleEdit = (a: Account) => { edAcc.value = a; showEdit.value = true }
const toggleSelectAllVisible = (event: Event) => {
  const target = event.target as HTMLInputElement
  toggleVisible(target.checked)
}
const handleBulkDelete = async () => { if(!confirm(t('common.confirm'))) return; try { await Promise.all(selIds.value.map(id => adminAPI.accounts.delete(id))); clearSelection(); reload() } catch (error) { console.error('Failed to bulk delete accounts:', error) } }
const handleBulkResetStatus = async () => {
  if (!confirm(t('common.confirm'))) return
  try {
    const result = await adminAPI.accounts.batchClearError(selIds.value)
    if (result.failed > 0) {
      appStore.showError(t('admin.accounts.bulkActions.partialSuccess', { success: result.success, failed: result.failed }))
    } else {
      appStore.showSuccess(t('admin.accounts.bulkActions.resetStatusSuccess', { count: result.success }))
      clearSelection()
    }
    reload()
  } catch (error) {
    console.error('Failed to bulk reset status:', error)
    appStore.showError(String(error))
  }
}
const handleBulkRefreshToken = async () => {
  if (!confirm(t('common.confirm'))) return
  try {
    const result = await adminAPI.accounts.batchRefresh(selIds.value)
    if (result.failed > 0) {
      appStore.showError(t('admin.accounts.bulkActions.partialSuccess', { success: result.success, failed: result.failed }))
    } else {
      appStore.showSuccess(t('admin.accounts.bulkActions.refreshTokenSuccess', { count: result.success }))
      clearSelection()
    }
    reload()
  } catch (error) {
    console.error('Failed to bulk refresh token:', error)
    appStore.showError(String(error))
  }
}
const handleBulkSyncUpstreamRate = async () => {
  if (syncingRateBatch.value || selIds.value.length === 0) return
  syncingRateBatch.value = true
  try {
    const result = await adminAPI.accounts.batchSyncUpstreamRateMultiplier(selIds.value)
    const updated = result.results
      .map(item => item.account)
      .filter((account): account is Account => Boolean(account))
    if (updated.length > 0) {
      updated.forEach(syncAccountRefs)
      mergeAccountsIncrementally(updated)
    }
		recordUpstreamRateSyncResult(result)
    if (result.failed > 0) {
      appStore.showError(t('admin.accounts.bulkActions.syncUpstreamRatePartial', {
        success: result.success,
        failed: result.failed
      }))
    } else {
      appStore.showSuccess(t('admin.accounts.bulkActions.syncUpstreamRateSuccess', { count: result.success }))
      clearSelection()
    }
  } catch (error: any) {
    console.error('Failed to batch sync upstream rate multiplier:', error)
    appStore.showError(error?.message || t('admin.accounts.syncUpstreamRateMultiplierFailed'))
  } finally {
    syncingRateBatch.value = false
  }
}
const handleSyncAllUpstreamRates = async () => {
  if (syncingAllUpstreamRates.value) return
  showSyncAllUpstreamRatesDialog.value = false
  syncingAllUpstreamRates.value = true
  try {
    const result = await adminAPI.accounts.syncAllUpstreamRateMultipliers()
    const updated = result.results
      .map(item => item.account)
      .filter((account): account is Account => Boolean(account))
    updated.forEach(syncAccountRefs)
    mergeAccountsIncrementally(updated)
		recordUpstreamRateSyncResult(result)
    failedUpstreamRateAccountIds.value = result.results
      .filter(item => !item.success && item.account_id > 0)
      .map(item => item.account_id)
    if (result.failed > 0) {
      appStore.showError(t('admin.accounts.syncAllUpstreamRatesPartial', {
        success: result.success,
        failed: result.failed
      }))
    } else {
      failedUpstreamRateAccountIds.value = []
      appStore.showSuccess(t('admin.accounts.syncAllUpstreamRatesSuccess', { count: result.success }))
    }
  } catch (error: any) {
    console.error('Failed to sync all upstream rate multipliers:', error)
    appStore.showError(error?.message || t('admin.accounts.syncUpstreamRateMultiplierFailed'))
  } finally {
    syncingAllUpstreamRates.value = false
  }
}
const handleRetryFailedUpstreamRates = async () => {
  const accountIds = [...new Set(failedUpstreamRateAccountIds.value)]
  if (syncingRateBatch.value || accountIds.length === 0) return
  syncingRateBatch.value = true
  try {
    const result = await adminAPI.accounts.batchSyncUpstreamRateMultiplier(accountIds)
    const updated = result.results
      .map(item => item.account)
      .filter((account): account is Account => Boolean(account))
    updated.forEach(syncAccountRefs)
    mergeAccountsIncrementally(updated)
		recordUpstreamRateSyncResult(result)
    failedUpstreamRateAccountIds.value = result.results
      .filter(item => !item.success && item.account_id > 0)
      .map(item => item.account_id)
    if (result.failed > 0) {
      appStore.showError(t('admin.accounts.retryFailedUpstreamRatesPartial', {
        success: result.success,
        failed: result.failed
      }))
    } else {
      appStore.showSuccess(t('admin.accounts.retryFailedUpstreamRatesSuccess', { count: result.success }))
    }
  } catch (error: any) {
    console.error('Failed to retry upstream rate multiplier sync:', error)
    appStore.showError(error?.message || t('admin.accounts.syncUpstreamRateMultiplierFailed'))
  } finally {
    syncingRateBatch.value = false
  }
}
const updateSchedulableInList = (accountIds: number[], schedulable: boolean) => {
  if (accountIds.length === 0) return
  const idSet = new Set(accountIds)
  accounts.value = accounts.value.map((account) => (idSet.has(account.id) ? { ...account, schedulable } : account))
}
const normalizeBulkSchedulableResult = (
  result: {
    success?: number
    failed?: number
    success_ids?: number[]
    failed_ids?: number[]
    results?: Array<{ account_id: number; success: boolean }>
  },
  accountIds: number[]
) => {
  const responseSuccessIds = Array.isArray(result.success_ids) ? result.success_ids : []
  const responseFailedIds = Array.isArray(result.failed_ids) ? result.failed_ids : []
  if (responseSuccessIds.length > 0 || responseFailedIds.length > 0) {
    return {
      successIds: responseSuccessIds,
      failedIds: responseFailedIds,
      successCount: typeof result.success === 'number' ? result.success : responseSuccessIds.length,
      failedCount: typeof result.failed === 'number' ? result.failed : responseFailedIds.length,
      hasIds: true,
      hasCounts: true
    }
  }

  const results = Array.isArray(result.results) ? result.results : []
  if (results.length > 0) {
    const successIds = results.filter(item => item.success).map(item => item.account_id)
    const failedIds = results.filter(item => !item.success).map(item => item.account_id)
    return {
      successIds,
      failedIds,
      successCount: typeof result.success === 'number' ? result.success : successIds.length,
      failedCount: typeof result.failed === 'number' ? result.failed : failedIds.length,
      hasIds: true,
      hasCounts: true
    }
  }

  const hasExplicitCounts = typeof result.success === 'number' || typeof result.failed === 'number'
  const successCount = typeof result.success === 'number' ? result.success : 0
  const failedCount = typeof result.failed === 'number' ? result.failed : 0
  if (hasExplicitCounts && failedCount === 0 && successCount === accountIds.length && accountIds.length > 0) {
    return {
      successIds: accountIds,
      failedIds: [],
      successCount,
      failedCount,
      hasIds: true,
      hasCounts: true
    }
  }

  return {
    successIds: [],
    failedIds: [],
    successCount,
    failedCount,
    hasIds: false,
    hasCounts: hasExplicitCounts
  }
}
const handleBulkToggleSchedulable = async (schedulable: boolean) => {
  const accountIds = [...selIds.value]
  try {
    const result = await adminAPI.accounts.bulkUpdate(accountIds, { schedulable })
    const { successIds, failedIds, successCount, failedCount, hasIds, hasCounts } = normalizeBulkSchedulableResult(result, accountIds)
    if (!hasIds && !hasCounts) {
      appStore.showError(t('admin.accounts.bulkSchedulableResultUnknown'))
      setSelectedIds(accountIds)
      load().catch((error) => {
        console.error('Failed to refresh accounts:', error)
      })
      return
    }
    if (successIds.length > 0) {
      updateSchedulableInList(successIds, schedulable)
    }
    if (successCount > 0 && failedCount === 0) {
      const message = schedulable
        ? t('admin.accounts.bulkSchedulableEnabled', { count: successCount })
        : t('admin.accounts.bulkSchedulableDisabled', { count: successCount })
      appStore.showSuccess(message)
    }
    if (failedCount > 0) {
      const message = hasCounts || hasIds
        ? t('admin.accounts.bulkSchedulablePartial', { success: successCount, failed: failedCount })
        : t('admin.accounts.bulkSchedulableResultUnknown')
      appStore.showError(message)
      setSelectedIds(failedIds.length > 0 ? failedIds : accountIds)
    } else {
      if (hasIds) clearSelection()
      else setSelectedIds(accountIds)
    }
  } catch (error) {
    console.error('Failed to bulk toggle schedulable:', error)
    appStore.showError(t('common.error'))
  }
}
const buildBulkEditFilterSnapshot = () => {
  const rawParams = toRaw(params) as Record<string, unknown>
  const sortOrder: AccountSortOrder = rawParams.sort_order === 'desc' ? 'desc' : 'asc'
  return {
    platform: typeof rawParams.platform === 'string' ? rawParams.platform : '',
    type: typeof rawParams.type === 'string' ? rawParams.type : '',
    status: typeof rawParams.status === 'string' ? rawParams.status : '',
    anomaly_reason: typeof rawParams.anomaly_reason === 'string' ? rawParams.anomaly_reason : '',
    group: typeof rawParams.group === 'string' ? rawParams.group : '',
    search: typeof rawParams.search === 'string' ? rawParams.search : '',
    privacy_mode: typeof rawParams.privacy_mode === 'string' ? rawParams.privacy_mode : '',
    sort_by: typeof rawParams.sort_by === 'string' ? rawParams.sort_by : '',
    sort_order: sortOrder
  }
}

const collectSelectionMetadata = (rows: Account[]) => {
  const selectedPlatforms = Array.from(new Set(rows.map(account => account.platform)))
  const selectedTypes = Array.from(new Set(rows.map(account => account.type)))
  return { selectedPlatforms, selectedTypes }
}

const openBulkEditSelected = () => {
  bulkEditTarget.value = {
    mode: 'selected',
    accountIds: [...selIds.value],
    selectedPlatforms: [...selPlatforms.value],
    selectedTypes: [...selTypes.value]
  }
  showBulkEdit.value = true
}

const openBulkEditFiltered = async () => {
  const filters = buildBulkEditFilterSnapshot()
  const preview = await adminAPI.accounts.list(1, 100, filters)
  const { selectedPlatforms, selectedTypes } = collectSelectionMetadata(preview.items)
  bulkEditTarget.value = {
    mode: 'filtered',
    filters,
    previewCount: preview.total,
    selectedPlatforms,
    selectedTypes
  }
  showBulkEdit.value = true
}

const handleBulkUpdated = () => {
  showBulkEdit.value = false
  bulkEditTarget.value = null
  clearSelection()
  reload()
}
const handleDataImported = () => { showImportData.value = false; reload() }
const canSyncUpstreamRate = (account: Account) => account.type === 'apikey'
const handleSyncUpstreamRate = async (account: Account) => {
  if (syncingRateAccountId.value !== null) return
  syncingRateAccountId.value = account.id
  try {
    const result = await adminAPI.accounts.syncUpstreamRateMultiplier(account.id)
    syncAccountRefs(result.account)
    mergeAccountsIncrementally([result.account])
		recordUpstreamRateSyncResult({
			total: 1,
			success: 1,
			failed: 0,
			results: [{
				account_id: account.id,
				account_name: account.name,
				success: true,
				previous_rate_multiplier: result.previous_rate_multiplier,
				rate_multiplier: result.rate_multiplier,
				changed: result.changed,
				significant_change: result.significant_change,
				account: result.account,
				source: result.source,
			}],
		}, result.significant_change)
    appStore.showSuccess(t('admin.accounts.syncUpstreamRateMultiplierSuccessWithChange', {
      previous: result.previous_rate_multiplier.toFixed(4),
	  rate: result.rate_multiplier.toFixed(4)
    }))
  } catch (error: any) {
    console.error('Failed to sync upstream rate multiplier:', error)
    const message = error?.status === 404
      ? t('admin.accounts.syncUpstreamRateMultiplierRouteMissing')
      : (error?.message || error?.response?.data?.message || t('admin.accounts.syncUpstreamRateMultiplierFailed'))
    appStore.showError(message)
  } finally {
    syncingRateAccountId.value = null
  }
}
const ACCOUNT_UNGROUPED_GROUP_QUERY_VALUE = 'ungrouped'
const ACCOUNT_PRIVACY_MODE_UNSET_QUERY_VALUE = '__unset__'
const anomalyReasonLabel = (code: AccountAnomalyReasonCode) => t(`admin.accounts.anomalyReasons.${code}.label`)
const anomalyReasonAction = (code: AccountAnomalyReasonCode) => t(`admin.accounts.anomalyReasons.${code}.action`)
const buildAccountQueryFilters = () => ({
  platform: params.platform || '',
  type: params.type || '',
  status: params.status || '',
  anomaly_reason: params.anomaly_reason || '',
  group: params.group || '',
  privacy_mode: params.privacy_mode || '',
  search: params.search || '',
  sort_by: sortState.sort_by,
  sort_order: sortState.sort_order
})
const accountAnomalySummaryItems = computed(() => {
  const counts = new Map<AccountAnomalyReasonCode, number>()
  for (const account of accounts.value) {
    const code = deriveAccountAnomalyReason(account).code
    if (code === 'healthy') continue
    counts.set(code, (counts.get(code) || 0) + 1)
  }
  return accountAnomalyReasonFilterOrder
    .map(code => ({ code, count: counts.get(code) || 0 }))
    .filter(item => item.count > 0)
})
const applyAnomalyReasonFilter = (code: AccountAnomalyReasonCode | '') => {
  params.anomaly_reason = code
  debouncedReload()
}
const accountMatchesCurrentFilters = (account: Account) => {
  const filters = buildAccountQueryFilters()
  if (filters.platform && account.platform !== filters.platform) return false
  if (filters.type && account.type !== filters.type) return false
  if (filters.status) {
    const now = Date.now()
    const rateLimitResetAt = account.rate_limit_reset_at ? new Date(account.rate_limit_reset_at).getTime() : Number.NaN
    const isRateLimited = Number.isFinite(rateLimitResetAt) && rateLimitResetAt > now
    const tempUnschedUntil = account.temp_unschedulable_until ? new Date(account.temp_unschedulable_until).getTime() : Number.NaN
    const isTempUnschedulable = Number.isFinite(tempUnschedUntil) && tempUnschedUntil > now

    if (filters.status === 'active') {
      if (account.status !== 'active' || isRateLimited || isTempUnschedulable || !account.schedulable) return false
    } else if (filters.status === 'rate_limited') {
      if (account.status !== 'active' || !isRateLimited || isTempUnschedulable) return false
    } else if (filters.status === 'temp_unschedulable') {
      if (account.status !== 'active' || !isTempUnschedulable) return false
    } else if (filters.status === 'unschedulable') {
      if (account.status !== 'active' || account.schedulable || isRateLimited || isTempUnschedulable) return false
    } else if (account.status !== filters.status) {
      return false
    }
  }
  if (filters.anomaly_reason) {
    if (deriveAccountAnomalyReason(account).code !== filters.anomaly_reason) return false
  }
  if (filters.group) {
    const groupIds = account.group_ids ?? account.groups?.map((group) => group.id) ?? []
    if (filters.group === ACCOUNT_UNGROUPED_GROUP_QUERY_VALUE) {
      if (groupIds.length > 0) return false
    } else if (!groupIds.includes(Number(filters.group))) {
      return false
    }
  }
  const privacyMode = typeof account.extra?.privacy_mode === 'string' ? account.extra.privacy_mode : ''
  if (filters.privacy_mode) {
    if (filters.privacy_mode === ACCOUNT_PRIVACY_MODE_UNSET_QUERY_VALUE) {
      if (privacyMode.trim() !== '') return false
    } else if (privacyMode !== filters.privacy_mode) {
      return false
    }
  }
  const search = String(filters.search || '').trim().toLowerCase()
  if (search && !account.name.toLowerCase().includes(search)) return false
  return true
}
const mergeRuntimeFields = (oldAccount: Account, updatedAccount: Account): Account => ({
  ...updatedAccount,
  current_concurrency: updatedAccount.current_concurrency ?? oldAccount.current_concurrency,
  current_window_cost: updatedAccount.current_window_cost ?? oldAccount.current_window_cost,
  active_sessions: updatedAccount.active_sessions ?? oldAccount.active_sessions
})

const syncPaginationAfterLocalRemoval = () => {
  const nextTotal = Math.max(0, pagination.total - 1)
  pagination.total = nextTotal
  pagination.pages = nextTotal > 0 ? Math.ceil(nextTotal / pagination.page_size) : 0

  const maxPage = Math.max(1, pagination.pages || 1)

  if (pagination.page > maxPage) {
    pagination.page = maxPage
  }
  // 行被本地移除后不立刻全量补页，改为提示用户手动同步。
  hasPendingListSync.value = nextTotal > 0
}

const patchAccountInList = (updatedAccount: Account) => {
  const index = accounts.value.findIndex(account => account.id === updatedAccount.id)
  if (index === -1) return
  const mergedAccount = mergeRuntimeFields(accounts.value[index], updatedAccount)
  if (!accountMatchesCurrentFilters(mergedAccount)) {
    accounts.value = accounts.value.filter(account => account.id !== mergedAccount.id)
    syncPaginationAfterLocalRemoval()
    removeSelectedAccounts([mergedAccount.id])
    return
  }
  const nextAccounts = [...accounts.value]
  nextAccounts[index] = mergedAccount
  accounts.value = nextAccounts
  syncAccountRefs(mergedAccount)
}
const handleAccountUpdated = (updatedAccount: Account) => {
  patchAccountInList(updatedAccount)
  enterAutoRefreshSilentWindow()
}
const formatExportTimestamp = () => {
  const now = new Date()
  const pad2 = (value: number) => String(value).padStart(2, '0')
  return `${now.getFullYear()}${pad2(now.getMonth() + 1)}${pad2(now.getDate())}${pad2(now.getHours())}${pad2(now.getMinutes())}${pad2(now.getSeconds())}`
}
const openExportDataDialog = () => {
  includeProxyOnExport.value = true
  showExportDataDialog.value = true
}
const handleExportData = async () => {
  if (exportingData.value) return
  exportingData.value = true
  try {
    const dataPayload = await adminAPI.accounts.exportData(
      selIds.value.length > 0
        ? { ids: selIds.value, includeProxies: includeProxyOnExport.value }
        : {
            includeProxies: includeProxyOnExport.value,
            filters: buildAccountQueryFilters()
          }
    )
    const timestamp = formatExportTimestamp()
    const filename = `sub2api-account-${timestamp}.json`
    const blob = new Blob([JSON.stringify(dataPayload, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    link.click()
    URL.revokeObjectURL(url)
    appStore.showSuccess(t('admin.accounts.dataExported'))
  } catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.dataExportFailed'))
  } finally {
    exportingData.value = false
    showExportDataDialog.value = false
  }
}
const closeTestModal = () => { showTest.value = false; testingAcc.value = null }
const closeStatsModal = () => { showStats.value = false; statsAcc.value = null }
const closeReAuthModal = () => { showReAuth.value = false; reAuthAcc.value = null }
const handleTest = (a: Account) => { testingAcc.value = a; showTest.value = true }
const handleViewStats = (a: Account) => { statsAcc.value = a; showStats.value = true }
const handleSchedule = async (a: Account) => {
  scheduleAcc.value = a
  scheduleModelOptions.value = []
  showSchedulePanel.value = true
  try {
    const models = await adminAPI.accounts.getAvailableModels(a.id)
    scheduleModelOptions.value = models.map((m: ClaudeModel) => ({ value: m.id, label: m.display_name || m.id }))
  } catch {
    scheduleModelOptions.value = []
  }
}
const closeSchedulePanel = () => { showSchedulePanel.value = false; scheduleAcc.value = null; scheduleModelOptions.value = [] }
const onRevertFallback = async (a: Account) => {
  try {
    await adminAPI.accounts.revertProxyFallback(a.id)
    appStore.showSuccess(t('admin.accounts.revertProxySuccess'))
    reload()
  } catch (error: any) {
    console.error('Failed to revert proxy fallback:', error)
    appStore.showError(error?.response?.data?.message || t('admin.accounts.revertProxyFailed'))
  }
}
const handleDelete = (a: Account) => { deletingAcc.value = a; showDeleteDialog.value = true }
const confirmDelete = async () => { if(!deletingAcc.value) return; try { await adminAPI.accounts.delete(deletingAcc.value.id); showDeleteDialog.value = false; deletingAcc.value = null; reload() } catch (error) { console.error('Failed to delete account:', error) } }
const handleToggleSchedulable = async (a: Account) => {
  const nextSchedulable = !a.schedulable
  togglingSchedulable.value = a.id
  try {
    const updated = await adminAPI.accounts.setSchedulable(a.id, nextSchedulable)
    updateSchedulableInList([a.id], updated?.schedulable ?? nextSchedulable)
    enterAutoRefreshSilentWindow()
  } catch (error) {
    console.error('Failed to toggle schedulable:', error)
    appStore.showError(t('admin.accounts.failedToToggleSchedulable'))
  } finally {
    togglingSchedulable.value = null
  }
}
const handleShowTempUnsched = (a: Account) => { tempUnschedAcc.value = a; showTempUnsched.value = true }
const handleTempUnschedReset = async (updated: Account) => {
  showTempUnsched.value = false
  tempUnschedAcc.value = null
  patchAccountInList(updated)
  enterAutoRefreshSilentWindow()
}
const formatExpiresAt = (value: number | null) => {
  if (!value) return '-'
  return formatDateTime(
    new Date(value * 1000),
    {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    },
    'sv-SE'
  )
}
const isExpired = (value: number | null) => {
  if (!value) return false
  return value * 1000 <= Date.now()
}
// 所绑定代理的有效期(逻辑同 /admin/proxies,见 utils/proxyExpiry)
const proxyExpiryBadge = (p: AccountProxy): string => proxyExpiryBadgeClass(p.expires_at, p.status)
const proxyExpiryText = (p: AccountProxy): string => {
  const { key, params } = proxyExpiryLabelKey(p.expires_at, p.status)
  return params ? t(key, params) : t(key)
}

// 点击外部关闭顶部下拉菜单
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (accountToolsDropdownRef.value && !accountToolsDropdownRef.value.contains(target)) {
    closeAccountToolsDropdown()
  }
  if (autoRefreshDropdownRef.value && !autoRefreshDropdownRef.value.contains(target)) {
    showAutoRefreshDropdown.value = false
  }
}

onMounted(async () => {
  load()
  void loadScheduledCostRateSyncSetting()
  try {
    const [p, g] = await Promise.all([adminAPI.proxies.getAll(), adminAPI.groups.getAll()])
    proxies.value = p
    groups.value = g
  } catch (error) {
    console.error('Failed to load proxies/groups:', error)
  }
  document.addEventListener('click', handleClickOutside)

  if (autoRefreshEnabled.value) {
    autoRefreshCountdown.value = autoRefreshIntervalSeconds.value
    resumeAutoRefresh()
  } else {
    pauseAutoRefresh()
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.sst-admin-page {
  width: 100%;
  min-width: 0;
  overflow-x: hidden;
}

.sst-admin-page :deep(.table-page-layout),
.sst-admin-page :deep(.layout-section-fixed),
.sst-admin-page :deep(.layout-section-scrollable),
.sst-admin-page :deep(.table-scroll-container) {
  width: 100%;
  max-width: 100%;
  min-width: 0;
}

.sst-admin-page :deep(.table-page-layout) {
  height: auto;
  min-height: calc(100vh - 64px - 4rem);
}

.sst-admin-page :deep(.layout-section-scrollable) {
  flex: 0 0 auto;
}

.sst-admin-page :deep(.table-scroll-container) {
  height: auto;
}

.sst-admin-page :deep(.layout-section-fixed) {
  overflow-x: hidden;
}

.account-table-shell {
  width: 100%;
  min-width: 0;
}

.account-table-grid {
  min-width: 0;
  flex: 1 1 auto;
}

.account-table-grid :deep(.table-wrapper) {
  width: 100%;
  max-width: 100%;
  min-width: 0;
  overflow-x: auto;
  overflow-y: auto;
  height: var(--account-table-visible-height, min(72vh, 52rem));
  min-height: var(--account-table-visible-height, min(72vh, 52rem));
  scrollbar-gutter: stable;
}

.account-table-grid :deep(.sticky-header-cell) {
  padding-top: 0.65rem;
  padding-bottom: 0.65rem;
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}

.account-table-grid :deep(tbody td) {
  padding-top: 0.55rem;
  padding-bottom: 0.55rem;
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}

.account-table-grid :deep(table) {
  min-width: max-content;
}

.account-table-grid :deep(.table-wrapper::-webkit-scrollbar) {
  height: 12px;
}

.account-table-grid :deep(.table-wrapper::-webkit-scrollbar-track) {
  background: rgba(0, 0, 0, 0.03);
  border-radius: 6px;
}

.account-table-grid :deep(.table-wrapper::-webkit-scrollbar-thumb) {
  background: rgba(107, 114, 128, 0.72);
  border-radius: 6px;
}

.account-table-grid :deep(.table-wrapper::-webkit-scrollbar-thumb:hover) {
  background: rgba(75, 85, 99, 0.88);
}

.account-tools-menu-item {
  @apply flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors;
  color: #38413a;
  background: transparent;
}

.account-tools-menu-item:hover {
  background: rgba(167, 58, 42, 0.055);
}

.account-tools-menu-icon {
  @apply inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-md;
  border: 1px solid rgba(198, 184, 157, 0.34);
  background: rgba(250, 247, 239, 0.78);
}

.account-tools-dropdown {
  max-width: min(20rem, calc(100vw - 2rem));
}

.account-health-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  max-width: 8.5rem;
  border-radius: 999px;
  border: 1px solid rgba(198, 184, 157, 0.42);
  padding: 0.125rem 0.45rem;
  font-size: 0.6875rem;
  line-height: 1rem;
  font-weight: 600;
  white-space: nowrap;
}

.account-health-pill strong {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.625rem;
  opacity: 0.82;
}

.account-health-dot {
  width: 0.375rem;
  height: 0.375rem;
  flex: 0 0 auto;
  border-radius: 999px;
  background: currentColor;
  opacity: 0.72;
}

.account-health-pill-good {
  color: #51624f;
  background: rgba(233, 239, 225, 0.72);
}

.account-health-pill-warning {
  color: #8a6424;
  background: rgba(250, 238, 204, 0.82);
}

.account-health-pill-critical {
  color: #9b3f31;
  background: rgba(248, 225, 217, 0.82);
}

</style>
<style>
.dark .account-tools-menu-item {
  color: #d7d0c2;
}

.dark .account-tools-menu-item:hover {
  background: rgba(167, 58, 42, 0.12);
}

.dark .account-tools-menu-icon {
  border-color: rgba(48, 52, 43, 0.84);
  background: rgba(24, 26, 21, 0.74);
}

.dark .account-health-pill {
  border-color: rgba(91, 82, 64, 0.78);
}

.dark .account-health-pill-good {
  color: #a8b89e;
  background: rgba(64, 80, 56, 0.34);
}

.dark .account-health-pill-warning {
  color: #d4b16b;
  background: rgba(105, 76, 30, 0.34);
}

.dark .account-health-pill-critical {
  color: #df9a8d;
  background: rgba(111, 46, 38, 0.34);
}
</style>

