<template>
  <AppLayout>
    <div class="keys-page space-y-6">
      <section class="keys-hero overflow-hidden rounded-zen border border-zen-paperLine bg-white/45 p-6 shadow-paper dark:border-zen-nightLine dark:bg-zen-nightPanel/70 lg:p-7">
        <div class="grid gap-6 lg:grid-cols-[1fr_auto] lg:items-end">
          <div>
            <div class="mb-4 flex items-center gap-4">
              <span class="h-px w-14 bg-zen-paperLine dark:bg-zen-nightLine"></span>
              <span class="font-mono text-xs uppercase tracking-[0.34em] text-zen-mist dark:text-zen-stone">{{ keysCopy.kicker }}</span>
            </div>
            <h1 class="font-serif text-3xl font-semibold text-zen-ink dark:text-zen-paper sm:text-4xl">{{ keysCopy.title }}</h1>
          </div>

          <div class="keys-ledger grid gap-3 sm:grid-cols-3 lg:min-w-[28rem]">
            <div class="keys-ledger-item">
              <span>{{ keysCopy.currentList }}</span>
              <strong>{{ pagination.total.toLocaleString() }}</strong>
            </div>
            <div class="keys-ledger-item">
              <span>{{ keysCopy.activeOnPage }}</span>
              <strong>{{ activeKeyCount.toLocaleString() }}</strong>
            </div>
            <div class="keys-ledger-item">
              <span>{{ keysCopy.groupsAvailable }}</span>
              <strong>{{ groups.length.toLocaleString() }}</strong>
            </div>
          </div>
        </div>
      </section>

      <section class="keys-workspace">
        <div class="keys-toolbar-shell">
          <div class="keys-toolbar-grid">
            <div class="keys-filter-field keys-filter-field-key">
              <Select
                :model-value="selectedFilterKeyId"
                class="w-full"
                :options="keyFilterOptions"
                :searchable="true"
                :aria-label="keysCopy.selectKeyAria"
                :search-placeholder="keysCopy.searchKeyPlaceholder"
                dropdown-class="keys-filter-dropdown"
                @update:model-value="onKeyFilterChange"
              />
            </div>

            <div class="keys-filter-field">
              <Select
                :model-value="filterGroupId"
                class="w-full"
                :options="groupFilterOptions"
                :aria-label="keysCopy.selectGroupAria"
                dropdown-class="keys-filter-dropdown"
                @update:model-value="onGroupFilterChange"
              />
            </div>

            <div class="keys-filter-field">
              <Select
                :model-value="filterStatus"
                class="w-full"
                :options="statusFilterOptions"
                :aria-label="keysCopy.selectStatusAria"
                dropdown-class="keys-filter-dropdown"
                @update:model-value="onStatusFilterChange"
              />
            </div>

            <div class="keys-toolbar-actions">
              <button @click="openCreateModal" class="btn btn-primary" data-tour="keys-create-btn">
                <Icon name="plus" size="md" class="mr-2" />
                {{ t('keys.createKey') }}
              </button>
              <button
                @click="loadApiKeys"
                :disabled="loading"
                class="btn btn-secondary"
                :title="t('common.refresh')"
              >
                <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
                {{ t('common.refresh') }}
              </button>
              <div class="relative" ref="columnDropdownRef">
                <button
                  @click="showColumnDropdown = !showColumnDropdown"
                  class="btn btn-secondary px-2 md:px-3"
                  :title="t('keys.columnSettings')"
                >
                  <Icon name="cog" size="md" />
                  <span class="hidden md:inline">{{ t('keys.columnSettings') }}</span>
                </button>
                <div
                  v-if="showColumnDropdown"
                  class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
                >
                  <button
                    v-for="col in toggleableColumns"
                    :key="col.key"
                    @click="toggleColumn(col.key)"
                    class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
                  >
                    <span>{{ col.label }}</span>
                    <Icon v-if="isColumnVisible(col.key)" name="check" size="sm" class="text-primary-500" :stroke-width="2" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          v-if="publicSettings?.api_base_url || (publicSettings?.custom_endpoints?.length ?? 0) > 0"
          class="keys-endpoint-popover-shell"
        >
          <EndpointPopover
            :api-base-url="publicSettings?.api_base_url || ''"
            :custom-endpoints="publicSettings?.custom_endpoints || []"
          />
        </div>

        <div class="keys-data-shell">
          <div class="keys-access-strip">
            <div class="keys-access-value">
              <span>{{ keysCopy.endpoint }}</span>
              <strong>{{ resolvedApiBaseUrl }}</strong>
            </div>
            <div class="keys-access-strip-actions">
              <button type="button" @click="copyEndpoint">
                <Icon name="clipboard" size="sm" />
                {{ keysCopy.copyEndpoint }}
              </button>
              <button type="button" :disabled="apiKeys.length === 0" @click="openConnectionTestDialog">
                <Icon name="link" size="sm" />
                {{ keysCopy.connectionCheck }}
              </button>
            </div>
          </div>
        <DataTable
          :columns="columns"
          :data="apiKeys"
          :loading="loading"
          :server-side-sort="true"
          :sticky-first-column="false"
          :sticky-actions-column="false"
          default-sort-key="created_at"
          default-sort-order="desc"
          @sort="handleSort"
        >
          <template #cell-id="{ value }">
            <span class="font-mono text-xs text-gray-500 dark:text-gray-400">#{{ value }}</span>
          </template>

          <template #cell-key="{ value, row }">
            <div class="flex items-center gap-2">
              <code class="code text-xs">
                {{ maskApiKey(value) }}
              </code>
              <button
                @click="copyToClipboard(value, row.id)"
                class="rounded-lg p-1 transition-colors hover:bg-gray-100 dark:hover:bg-dark-700"
                :class="
                  copiedKeyId === row.id
                    ? 'text-green-500'
                    : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                "
                :title="copiedKeyId === row.id ? t('keys.copied') : t('keys.copyToClipboard')"
              >
                <Icon
                  v-if="copiedKeyId === row.id"
                  name="check"
                  size="sm"
                  :stroke-width="2"
                />
                <Icon v-else name="clipboard" size="sm" />
              </button>
            </div>
          </template>

          <template #cell-name="{ value, row }">
            <div class="key-name-cell">
              <span class="key-name-text font-medium text-gray-900 dark:text-white">{{ value }}</span>
              <Icon
                v-if="row.ip_whitelist?.length > 0 || row.ip_blacklist?.length > 0"
                name="shield"
                size="sm"
                class="key-name-flag"
                :title="t('keys.ipRestrictionEnabled')"
              />
            </div>
          </template>

          <template #cell-group="{ row }">
            <div class="group/dropdown relative">
              <button
                type="button"
                :ref="(el) => setGroupButtonRef(row.id, el)"
                @click.stop="openGroupSelector(row)"
                class="key-group-trigger"
                :title="t('keys.clickToChangeGroup')"
                aria-haspopup="listbox"
                :aria-expanded="groupSelectorKeyId === row.id"
              >
                <GroupBadge
                  v-if="row.group"
                  :name="row.group.name"
                  :platform="row.group.platform"
                  :subscription-type="row.group.subscription_type"
                  :rate-multiplier="row.group.rate_multiplier"
                  :user-rate-multiplier="userGroupRates[row.group.id]"
                />
                <span v-else class="key-group-empty">{{
                  t('keys.noGroup')
                }}</span>
                <span class="key-group-helper">{{ t('keys.selectGroup') }}</span>
                <Icon
                  name="arrowsUpDown"
                  size="sm"
                  class="key-group-chevron"
                  :stroke-width="2"
                />
              </button>
            </div>
          </template>

          <template #cell-current_concurrency="{ value }">
            <span
              :class="[
                'inline-flex min-w-8 items-center justify-center rounded px-2 py-1 text-sm font-semibold tabular-nums',
                (value ?? 0) > 0
                  ? 'bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200 dark:bg-emerald-900/25 dark:text-emerald-300 dark:ring-emerald-800'
                  : 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-dark-400'
              ]"
            >
              {{ value ?? 0 }}
            </span>
          </template>

          <template #cell-usage="{ row }">
            <div class="key-usage-cell">
              <div class="key-usage-ledger">
                <div class="key-usage-stat">
                  <span>{{ t('keys.today') }}</span>
                  <strong>{{ formatUsd(workbenchStats[row.id]?.today_actual_cost, 4) }}</strong>
                </div>
                <div class="key-usage-stat">
                  <span>{{ t('keys.total') }}</span>
                  <strong>{{ formatUsd(workbenchStats[row.id]?.total_actual_cost, 4) }}</strong>
                </div>
                <div class="key-usage-stat">
                  <span>{{ keysCopy.success24h }}</span>
                  <strong>{{ (workbenchStats[row.id]?.success_requests_24h ?? 0).toLocaleString() }}</strong>
                </div>
                <div class="key-usage-stat">
                  <span>{{ keysCopy.winRate24h }}</span>
                  <strong>{{ formatSuccessRate(workbenchStats[row.id]?.success_rate_24h) }}</strong>
                </div>
              </div>

              <div v-if="getRateLimitWindows(row).length" class="key-usage-meta">
                <div v-if="getRateLimitWindows(row).length" class="key-rate-inline-list">
                  <span
                    v-for="window in getRateLimitWindows(row)"
                    :key="window.key"
                    class="key-rate-inline-chip"
                    :class="`tone-${getMeterTone(window.usage, window.limit)}`"
                  >
                    <strong>{{ window.label }}</strong>
                    <span>{{ formatUsd(window.usage) }}/{{ formatUsd(window.limit) }}</span>
                  </span>
                </div>
              </div>
            </div>
          </template>

          <template #cell-status="{ value, row }">
            <div class="key-health-cell">
              <span :class="[
                'badge',
                value === 'active' ? 'badge-success' :
                value === 'quota_exhausted' ? 'badge-warning' :
                value === 'expired' ? 'badge-danger' :
                'badge-gray'
              ]">
                {{ t('keys.status.' + value) }}
              </span>
              <div class="key-health-lines">
                <span>{{ getKeyHealth(row).lastCall }}</span>
                <span>{{ getKeyHealth(row).requestPulse }}</span>
                <span :class="{ 'key-health-lines-alert': getKeyHealth(row).hasAttention }">{{ getKeyHealth(row).summary }}</span>
                <span
                  v-for="hint in getKeyHealth(row).modelHints"
                  :key="hint"
                  class="key-health-model-hint"
                >
                  {{ hint }}
                </span>
                <button
                  v-if="getKeyHealth(row).canReview"
                  type="button"
                  class="key-health-link"
                  @click="openKeyErrorLedger(row.id)"
                >
                  {{ keysCopy.viewRecentFailures }}
                </button>
              </div>
            </div>
          </template>

          <template #cell-last_used_at="{ value }">
            <span v-if="value" class="text-sm text-gray-500 dark:text-dark-400">
              {{ formatDateTime(value) }}
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
          </template>

          <template #cell-last_used_ip="{ value }">
            <span v-if="value" class="text-sm text-gray-500 dark:text-dark-400">
              {{ value }}
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
          </template>

          <template #cell-created_at="{ value }">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateOnly(value) }}</span>
          </template>

          <template #cell-actions="{ row }">
            <div class="key-row-actions">
              <!-- Use Key Button -->
              <button
                @click="openUseKeyModal(row)"
                class="key-row-action"
                :title="t('keys.useKey')"
              >
                <Icon name="terminal" size="sm" />
              </button>
              <!-- Import to CC Switch Button -->
              <button
                v-if="!publicSettings?.hide_ccs_import_button"
                @click="importToCcswitch(row)"
                class="key-row-action"
                :title="t('keys.importToCcSwitch')"
              >
                <Icon name="upload" size="sm" />
              </button>
              <!-- Toggle Status Button -->
              <button
                @click="toggleKeyStatus(row)"
                class="key-row-action"
                :title="row.status === 'active' ? t('keys.disable') : t('keys.enable')"
              >
                <Icon v-if="row.status === 'active'" name="ban" size="sm" />
                <Icon v-else name="checkCircle" size="sm" />
              </button>
              <!-- Edit Button -->
              <button
                @click="editKey(row)"
                class="key-row-action"
                :title="t('common.edit')"
              >
                <Icon name="edit" size="sm" />
              </button>
              <!-- Delete Button -->
              <button
                @click="confirmDelete(row)"
                class="key-row-action is-danger"
                :title="t('common.delete')"
              >
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState
              :title="t('keys.noKeysYet')"
              :description="t('keys.createFirstKey')"
              :action-text="t('keys.createKey')"
              @action="openCreateModal"
            />
          </template>
        </DataTable>

        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          :show-controls-when-single-page="false"
          page-size-dropdown-class="keys-page-size-dropdown"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
        </div>
      </section>
    </div>


    <BaseDialog
      :show="showConnectionTestDialog"
      :title="keysCopy.connectionConfig"
      width="normal"
      class="connection-test-modal"
      @close="closeConnectionTestDialog"
    >
      <div class="connection-test-dialog">
        <div class="connection-key-list">
          <button
            v-for="key in apiKeys"
            :key="key.id"
            type="button"
            :class="{ 'is-selected': connectionTestKeyId === key.id }"
            @click="selectConnectionTestKey(key.id)"
          >
            <span>{{ key.name }}</span>
            <small>{{ maskApiKey(key.key) }}</small>
          </button>
        </div>

        <div class="integration-kit" v-if="selectedConnectionTestKey">
          <div class="integration-kit-actions">
            <button type="button" @click="copyIntegrationSnippet('curl')">{{ keysCopy.copyCurl }}</button>
            <button type="button" @click="copyIntegrationSnippet('openai')">{{ keysCopy.copyOpenAiSdk }}</button>
            <button type="button" @click="copyIntegrationSnippet('env')">{{ keysCopy.copyEnv }}</button>
            <button type="button" @click="openUseKeyModalFromConnectionTest">{{ keysCopy.viewFullConfig }}</button>
          </div>
        </div>

        <div
          v-if="selectedConnectionTestKey && (selectedConnectionModelHints.length || selectedConnectionWorkbenchSummary?.latest_error)"
          class="connection-model-brief"
        >
          <div class="connection-model-brief-head">
            <div>
              <span>{{ keysCopy.recentModelHints }}</span>
              <strong>{{ selectedConnectionLatestErrorLabel }}</strong>
            </div>
            <button type="button" @click="openKeyErrorLedger(selectedConnectionTestKey.id)">{{ keysCopy.viewRecentFailures }}</button>
          </div>
          <ul v-if="selectedConnectionModelHints.length" class="connection-model-brief-list">
            <li v-for="hint in selectedConnectionModelHints" :key="hint">{{ hint }}</li>
          </ul>
        </div>

        <div v-if="connectionTestResult" class="connection-test-report" :class="`is-${connectionTestResult.tone}`">
          <div class="connection-report-head">
            <div>
              <span>{{ keysCopy.checkRecord }}</span>
              <strong>{{ connectionTestResult.title }}</strong>
            </div>
            <button type="button" @click="copyConnectionDiagnosticReport">{{ keysCopy.copyReport }}</button>
          </div>
          <p>{{ connectionTestResult.detail }}</p>
          <div class="connection-report-grid">
            <div><span>{{ keysCopy.endpoint }}</span><strong>{{ resolvedApiBaseUrl }}/v1</strong></div>
            <div><span>{{ keysCopy.keyStatus }}</span><strong>{{ selectedConnectionTestKey?.status || '-' }}</strong></div>
            <div><span>{{ keysCopy.modelsEndpoint }}</span><strong>{{ connectionTestResult.statusCode ? 'HTTP ' + connectionTestResult.statusCode : connectionTestResult.tone === 'success' ? keysCopy.reachable : keysCopy.notFinished }}</strong></div>
            <div><span>{{ keysCopy.responseLatency }}</span><strong>{{ connectionTestResult.latencyMs !== null ? connectionTestResult.latencyMs + 'ms' : '-' }}</strong></div>
          </div>
          <div v-if="connectionTestResult.sampleModels.length" class="connection-model-samples">
            <div class="connection-model-samples-head">
              <span>{{ keysCopy.modelVisibility }}</span>
              <strong>{{ connectionTestResult.availableModelCount !== null ? keysCopy.visibleModels(connectionTestResult.availableModelCount) : '-' }}</strong>
            </div>
            <div class="connection-model-sample-list">
              <code v-for="model in visibleConnectionModels" :key="model">{{ model }}</code>
              <button
                v-if="hiddenConnectionModelCount > 0 || showAllConnectionModels"
                type="button"
                class="connection-model-sample-more"
                @click="toggleConnectionModelsExpanded"
              >
                {{ showAllConnectionModels ? t('common.collapse') : '+' + hiddenConnectionModelCount }}
              </button>
            </div>
          </div>
          <small v-if="connectionTestResult.action">{{ connectionTestResult.action }}</small>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-secondary" @click="closeConnectionTestDialog">{{ t('common.cancel') }}</button>
          <button
            type="button"
            class="btn btn-primary"
            :disabled="!selectedConnectionTestKey || testingKeyId !== null"
            @click="testKeyConnection(selectedConnectionTestKey)"
          >
            <Icon name="refresh" size="sm" :class="testingKeyId !== null ? 'animate-spin' : ''" />
            {{ testingKeyId !== null ? keysCopy.checking : keysCopy.startCheck }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Create/Edit Modal -->
    <BaseDialog
      :show="showCreateModal || showEditModal"
      :title="showEditModal ? t('keys.editKey') : t('keys.createKey')"
      width="normal"
      class="key-editor-modal"
      @close="closeModals"
    >
      <form id="key-form" @submit.prevent="handleSubmit" class="key-editor-form space-y-5">
        <div>
          <label class="input-label">{{ t('keys.nameLabel') }}</label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="input"
            :placeholder="t('keys.namePlaceholder')"
            data-tour="key-form-name"
          />
        </div>

        <div>
          <label class="input-label">{{ t('keys.groupLabel') }}</label>
          <Select
            v-model="formData.group_id"
            :options="groupOptions"
            :placeholder="t('keys.selectGroup')"
            :searchable="true"
            :search-placeholder="t('keys.searchGroup')"
            dropdown-class="key-editor-dropdown"
            data-tour="key-form-group"
          >
            <template #selected="{ option }">
              <GroupBadge
                v-if="option"
                :name="(option as unknown as GroupOption).label"
                :platform="(option as unknown as GroupOption).platform"
                :subscription-type="(option as unknown as GroupOption).subscriptionType"
                :rate-multiplier="(option as unknown as GroupOption).rate"
                :user-rate-multiplier="(option as unknown as GroupOption).userRate"
              />
              <span v-else class="text-gray-400">{{ t('keys.selectGroup') }}</span>
            </template>
            <template #option="{ option, selected }">
              <GroupOptionItem
                :name="(option as unknown as GroupOption).label"
                :platform="(option as unknown as GroupOption).platform"
                :subscription-type="(option as unknown as GroupOption).subscriptionType"
                :rate-multiplier="(option as unknown as GroupOption).rate"
                :user-rate-multiplier="(option as unknown as GroupOption).userRate"
                :description="(option as unknown as GroupOption).description"
                :selected="selected"
              />
            </template>
          </Select>
        </div>

        <!-- Custom Key Section (only for create) -->
        <div v-if="!showEditModal" class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.customKeyLabel') }}</label>
            <button
              type="button"
              @click="formData.use_custom_key = !formData.use_custom_key"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.use_custom_key ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.use_custom_key ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>
          <div v-if="formData.use_custom_key">
            <input
              v-model="formData.custom_key"
              type="text"
              class="input font-mono"
              :placeholder="t('keys.customKeyPlaceholder')"
              :class="{ 'border-red-500 dark:border-red-500': customKeyError }"
            />
            <p v-if="customKeyError" class="mt-1 text-sm text-red-500">{{ customKeyError }}</p>
            <p v-else class="input-hint">{{ t('keys.customKeyHint') }}</p>
          </div>
        </div>

        <div v-if="showEditModal">
          <label class="input-label">{{ t('keys.statusLabel') }}</label>
          <Select
            v-model="formData.status"
            :options="statusOptions"
            :placeholder="t('keys.selectStatus')"
            dropdown-class="key-editor-dropdown"
          />
        </div>

        <div class="key-advanced-shell">
          <button
            type="button"
            class="key-advanced-toggle"
            :aria-expanded="advancedSettingsExpanded"
            @click="advancedSettingsExpanded = !advancedSettingsExpanded"
          >
            <div class="key-advanced-copy">
              <span>{{ keysCopy.advancedControl }}</span>
              <strong>{{ advancedSettingsSummary }}</strong>
            </div>
            <div class="key-advanced-meta">
              <div v-if="advancedSummaryItems.length" class="key-advanced-chips">
                <span v-for="item in advancedSummaryItems" :key="item">{{ item }}</span>
              </div>
              <small>{{ advancedSettingsExpanded ? t('common.collapse') : t('common.expand') }}</small>
              <Icon name="chevronDown" size="sm" :class="['transition-transform duration-200', advancedSettingsExpanded && 'rotate-180']" />
            </div>
          </button>

          <div v-if="advancedSettingsExpanded" class="key-advanced-panel space-y-5">

        <!-- IP Restriction Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.ipRestriction') }}</label>
            <button
              type="button"
              @click="formData.enable_ip_restriction = !formData.enable_ip_restriction"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_ip_restriction ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_ip_restriction ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_ip_restriction" class="space-y-4 pt-2">
            <div>
              <label class="input-label">{{ t('keys.ipWhitelist') }}</label>
              <textarea
                v-model="formData.ip_whitelist"
                rows="3"
                class="input font-mono text-sm"
                :placeholder="t('keys.ipWhitelistPlaceholder')"
              />
              <p class="input-hint">{{ t('keys.ipWhitelistHint') }}</p>
            </div>

            <div>
              <label class="input-label">{{ t('keys.ipBlacklist') }}</label>
              <textarea
                v-model="formData.ip_blacklist"
                rows="3"
                class="input font-mono text-sm"
                :placeholder="t('keys.ipBlacklistPlaceholder')"
              />
              <p class="input-hint">{{ t('keys.ipBlacklistHint') }}</p>
            </div>
          </div>
        </div>

        <!-- Quota Limit Section -->
        <div class="space-y-3">
          <label class="input-label">{{ t('keys.quotaLimit') }}</label>
          <!-- Switch commented out - always show input, 0 = unlimited
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.quotaLimit') }}</label>
            <button
              type="button"
              @click="formData.enable_quota = !formData.enable_quota"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_quota ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_quota ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>
          -->

          <div class="space-y-4">
            <div>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.quota"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="t('keys.quotaAmountPlaceholder')"
                />
              </div>
              <p class="input-hint">{{ t('keys.quotaAmountHint') }}</p>
            </div>

            <!-- Quota used display (only in edit mode) -->
            <div v-if="showEditModal && selectedKey && selectedKey.quota > 0">
              <label class="input-label">{{ t('keys.quotaUsed') }}</label>
              <div class="flex items-center gap-2">
                <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-700">
                  <span class="font-medium text-gray-900 dark:text-white">
                    ${{ selectedKey.quota_used?.toFixed(4) || '0.0000' }}
                  </span>
                  <span class="mx-2 text-gray-400">/</span>
                  <span class="text-gray-500 dark:text-gray-400">
                    ${{ selectedKey.quota?.toFixed(2) || '0.00' }}
                  </span>
                </div>
                <button
                  type="button"
                  @click="confirmResetQuota"
                  class="btn btn-secondary text-sm"
                  :title="t('keys.resetQuotaUsed')"
                >
                  {{ t('keys.reset') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Rate Limit Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.rateLimitSection') }}</label>
            <button
              type="button"
              @click="formData.enable_rate_limit = !formData.enable_rate_limit"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_rate_limit ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_rate_limit ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_rate_limit" class="space-y-4 pt-2">
            <p class="input-hint -mt-2">{{ t('keys.rateLimitHint') }}</p>
            <!-- 5-Hour Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit5h') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_5h"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_5h > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-700 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h ? 'text-red-500' :
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_5h?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_5h?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h ? 'bg-red-500' :
                      selectedKey.usage_5h >= selectedKey.rate_limit_5h * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_5h / selectedKey.rate_limit_5h) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- Daily Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit1d') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_1d"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_1d > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-700 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d ? 'text-red-500' :
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_1d?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_1d?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d ? 'bg-red-500' :
                      selectedKey.usage_1d >= selectedKey.rate_limit_1d * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_1d / selectedKey.rate_limit_1d) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- 7-Day Limit -->
            <div>
              <label class="input-label">{{ t('keys.rateLimit7d') }}</label>
              <div class="relative">
                <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
                <input
                  v-model.number="formData.rate_limit_7d"
                  type="number"
                  step="0.01"
                  min="0"
                  class="input pl-7"
                  :placeholder="'0'"
                />
              </div>
              <!-- Usage info (edit mode only) -->
              <div v-if="showEditModal && selectedKey && selectedKey.rate_limit_7d > 0" class="mt-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 rounded-lg bg-gray-100 px-3 py-2 dark:bg-dark-700 text-sm">
                    <span :class="[
                      'font-medium',
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d ? 'text-red-500' :
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d * 0.8 ? 'text-yellow-500' :
                      'text-gray-900 dark:text-white'
                    ]">
                      ${{ selectedKey.usage_7d?.toFixed(4) || '0.0000' }}
                    </span>
                    <span class="mx-2 text-gray-400">/</span>
                    <span class="text-gray-500 dark:text-gray-400">
                      ${{ selectedKey.rate_limit_7d?.toFixed(2) || '0.00' }}
                    </span>
                  </div>
                </div>
                <div class="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-dark-600">
                  <div
                    :class="[
                      'h-full rounded-full transition-all',
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d ? 'bg-red-500' :
                      selectedKey.usage_7d >= selectedKey.rate_limit_7d * 0.8 ? 'bg-yellow-500' :
                      'bg-green-500'
                    ]"
                    :style="{ width: Math.min((selectedKey.usage_7d / selectedKey.rate_limit_7d) * 100, 100) + '%' }"
                  />
                </div>
              </div>
            </div>

            <!-- Reset Rate Limit button (edit mode only) -->
            <div v-if="showEditModal && selectedKey && (selectedKey.rate_limit_5h > 0 || selectedKey.rate_limit_1d > 0 || selectedKey.rate_limit_7d > 0)">
              <button
                type="button"
                @click="confirmResetRateLimit"
                class="btn btn-secondary text-sm"
              >
                {{ t('keys.resetRateLimitUsage') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Expiration Section -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="input-label mb-0">{{ t('keys.expiration') }}</label>
            <button
              type="button"
              @click="formData.enable_expiration = !formData.enable_expiration"
              :class="[
                'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none',
                formData.enable_expiration ? 'bg-primary-600' : 'bg-gray-200 dark:bg-dark-600'
              ]"
            >
              <span
                :class="[
                  'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
                  formData.enable_expiration ? 'translate-x-4' : 'translate-x-0'
                ]"
              />
            </button>
          </div>

          <div v-if="formData.enable_expiration" class="space-y-4 pt-2">
            <!-- Quick select buttons (for both create and edit mode) -->
            <div class="flex flex-wrap gap-2">
              <button
                v-for="days in ['7', '30', '90']"
                :key="days"
                type="button"
                @click="setExpirationDays(parseInt(days))"
                :class="[
                  'rounded-lg px-3 py-1.5 text-sm transition-colors',
                  formData.expiration_preset === days
                    ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-700 dark:text-gray-400 dark:hover:bg-dark-600'
                ]"
              >
                {{ showEditModal ? t('keys.extendDays', { days }) : t('keys.expiresInDays', { days }) }}
              </button>
              <button
                type="button"
                @click="formData.expiration_preset = 'custom'"
                :class="[
                  'rounded-lg px-3 py-1.5 text-sm transition-colors',
                  formData.expiration_preset === 'custom'
                    ? 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200 dark:bg-dark-700 dark:text-gray-400 dark:hover:bg-dark-600'
                ]"
              >
                {{ t('keys.customDate') }}
              </button>
            </div>

            <!-- Date picker (always show for precise adjustment) -->
            <div>
              <label class="input-label">{{ t('keys.expirationDate') }}</label>
              <input
                v-model="formData.expiration_date"
                type="datetime-local"
                class="input"
              />
              <p class="input-hint">{{ t('keys.expirationDateHint') }}</p>
            </div>

            <!-- Current expiration display (only in edit mode) -->
            <div v-if="showEditModal && selectedKey?.expires_at" class="text-sm">
              <span class="text-gray-500 dark:text-gray-400">{{ t('keys.currentExpiration') }}: </span>
              <span class="font-medium text-gray-900 dark:text-white">
                {{ formatDateTime(selectedKey.expires_at) }}
              </span>
            </div>
          </div>
        </div>
          </div>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button @click="closeModals" type="button" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
          <button
            form="key-form"
            type="submit"
            :disabled="submitting"
            class="btn btn-primary"
            data-tour="key-form-submit"
          >
            <Icon v-if="submitting" name="refresh" size="sm" class="-ml-1 mr-2 animate-spin" />
            {{
              submitting
                ? t('keys.saving')
                : showEditModal
                  ? t('common.update')
                  : t('common.create')
            }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('keys.deleteKey')"
      :message="t('keys.deleteConfirmMessage', { name: selectedKey?.name })"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="handleDelete"
      @cancel="showDeleteDialog = false"
    />

    <!-- Reset Quota Confirmation Dialog -->
    <ConfirmDialog
      :show="showResetQuotaDialog"
      :title="t('keys.resetQuotaTitle')"
      :message="t('keys.resetQuotaConfirmMessage', { name: selectedKey?.name, used: selectedKey?.quota_used?.toFixed(4) })"
      :confirm-text="t('keys.reset')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="resetQuotaUsed"
      @cancel="showResetQuotaDialog = false"
    />

    <!-- Reset Rate Limit Confirmation Dialog -->
    <ConfirmDialog
      :show="showResetRateLimitDialog"
      :title="t('keys.resetRateLimitTitle')"
      :message="t('keys.resetRateLimitConfirmMessage', { name: selectedKey?.name })"
      :confirm-text="t('keys.reset')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="resetRateLimitUsage"
      @cancel="showResetRateLimitDialog = false"
    />

    <!-- Use Key Modal -->
    <UseKeyModal
      :show="showUseKeyModal"
      :api-key="selectedKey?.key || ''"
      :base-url="publicSettings?.api_base_url || ''"
      :platform="selectedKey?.group?.platform || null"
      :allow-messages-dispatch="selectedKey?.group?.allow_messages_dispatch || false"
      @close="closeUseKeyModal"
    />

    <!-- CCS Client Selection Dialog for Antigravity -->
    <BaseDialog
      :show="showCcsClientSelect"
      :title="t('keys.ccsClientSelect.title')"
      width="narrow"
      @close="closeCcsClientSelect"
    >
      <div class="space-y-4">
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {{ t('keys.ccsClientSelect.description') }}
	        </p>
	        <div class="grid grid-cols-2 gap-3">
	          <button
	            @click="handleCcsClientSelect('claude')"
	            class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 border-gray-200 dark:border-dark-600 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 transition-all"
	          >
	            <Icon name="terminal" size="xl" class="text-gray-600 dark:text-gray-400" />
	            <span class="font-medium text-gray-900 dark:text-white">{{
	              t('keys.ccsClientSelect.claudeCode')
	            }}</span>
	            <span class="text-xs text-gray-500 dark:text-gray-400">{{
	              t('keys.ccsClientSelect.claudeCodeDesc')
	            }}</span>
	          </button>
	          <button
	            @click="handleCcsClientSelect('gemini')"
	            class="flex flex-col items-center gap-2 p-4 rounded-xl border-2 border-gray-200 dark:border-dark-600 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 transition-all"
	          >
	            <Icon name="sparkles" size="xl" class="text-gray-600 dark:text-gray-400" />
	            <span class="font-medium text-gray-900 dark:text-white">{{
	              t('keys.ccsClientSelect.geminiCli')
	            }}</span>
	            <span class="text-xs text-gray-500 dark:text-gray-400">{{
	              t('keys.ccsClientSelect.geminiCliDesc')
	            }}</span>
	          </button>
	        </div>
	      </div>
      <template #footer>
        <div class="flex justify-end">
          <button @click="closeCcsClientSelect" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Group Selector Dropdown (Teleported to body to avoid overflow clipping) -->
    <Teleport to="body">
      <div
        v-if="groupSelectorKeyId !== null && dropdownPosition"
        ref="dropdownRef"
        class="animate-in fade-in slide-in-from-top-2 fixed z-[100000020] w-max min-w-[380px] overflow-hidden rounded-xl bg-white shadow-lg ring-1 ring-black/5 duration-200 dark:bg-dark-800 dark:ring-white/10"
        style="pointer-events: auto !important;"
        :style="{
          top: dropdownPosition.top !== undefined ? dropdownPosition.top + 'px' : undefined,
          bottom: dropdownPosition.bottom !== undefined ? dropdownPosition.bottom + 'px' : undefined,
          left: dropdownPosition.left + 'px'
        }"
        @click.stop
      >
        <!-- Search box -->
        <div class="border-b border-gray-100 p-2 dark:border-dark-700">
          <div class="relative">
            <Icon name="search" size="sm" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" :stroke-width="2" />
            <input
              v-model="groupSearchQuery"
              type="text"
              class="w-full rounded-lg border border-gray-200 bg-gray-50 py-1.5 pl-8 pr-3 text-sm text-gray-900 placeholder-gray-400 outline-none focus:border-primary-300 focus:ring-1 focus:ring-primary-300 dark:border-dark-600 dark:bg-dark-700 dark:text-white dark:placeholder-gray-500 dark:focus:border-primary-600 dark:focus:ring-primary-600"
              :placeholder="t('keys.searchGroup')"
              @click.stop
            />
          </div>
        </div>
        <!-- Group list -->
        <div class="max-h-80 overflow-y-auto p-1.5">
          <button
            v-for="option in filteredGroupOptions"
            :key="option.value ?? 'null'"
            @click="changeGroup(selectedKeyForGroup!, option.value)"
            :class="[
              'flex w-full items-center justify-between rounded-lg px-3 py-2.5 text-sm transition-colors',
              'border-b border-gray-100 last:border-0 dark:border-dark-700',
              selectedKeyForGroup?.group_id === option.value ||
              (!selectedKeyForGroup?.group_id && option.value === null)
                ? 'bg-primary-50 dark:bg-primary-900/20'
                : 'hover:bg-gray-100 dark:hover:bg-dark-700'
            ]"
            :title="option.description || undefined"
          >
            <GroupOptionItem
              :name="option.label"
              :platform="option.platform"
              :subscription-type="option.subscriptionType"
              :rate-multiplier="option.rate"
              :user-rate-multiplier="option.userRate"
              :description="option.description"
              :selected="
                selectedKeyForGroup?.group_id === option.value ||
                (!selectedKeyForGroup?.group_id && option.value === null)
              "
            />
          </button>
          <!-- Empty state when search has no results -->
          <div v-if="filteredGroupOptions.length === 0" class="py-4 text-center text-sm text-gray-400 dark:text-gray-500">
            {{ t('keys.noGroupFound') }}
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
	import { ref, computed, reactive, onMounted, onUnmounted, nextTick, type ComponentPublicInstance } from 'vue'
	import { useI18n } from 'vue-i18n'
	import { useAppStore } from '@/stores/app'
	import { useOnboardingStore } from '@/stores/onboarding'
	import { useClipboard } from '@/composables/useClipboard'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'
import { useRoute, useRouter } from 'vue-router'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
import { keysAPI, authAPI, userGroupsAPI, usageAPI } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
	import DataTable from '@/components/common/DataTable.vue'
	import Pagination from '@/components/common/Pagination.vue'
	import BaseDialog from '@/components/common/BaseDialog.vue'
	import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
	import EmptyState from '@/components/common/EmptyState.vue'
	import Select from '@/components/common/Select.vue'
	import Icon from '@/components/icons/Icon.vue'
	import UseKeyModal from '@/components/keys/UseKeyModal.vue'
	import EndpointPopover from '@/components/keys/EndpointPopover.vue'
	import GroupBadge from '@/components/common/GroupBadge.vue'
	import GroupOptionItem from '@/components/common/GroupOptionItem.vue'
	import type { ApiKey, Group, PublicSettings, SubscriptionType, GroupPlatform } from '@/types'
import type { Column } from '@/components/common/types'
import type { ApiKeyWorkbenchSummary } from '@/api/usage'
import { formatDateTime } from '@/utils/format'
import { maskApiKey } from '@/utils/maskApiKey'
import { buildWorkbenchLatestErrorLabel, buildWorkbenchModelHints } from '@/utils/apiKeyWorkbench'
import { extractErrorTextFromPayload, extractModelIdsFromModelsPayload, safeParseJson } from '@/utils/modelsEndpoint'
import {
  buildCcSwitchImportDeeplink,
  type CcSwitchClientType
} from '@/utils/ccswitchImport'

// Helper to format date for datetime-local input
const formatDateTimeLocal = (isoDate: string): string => {
  const date = new Date(isoDate)
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

interface GroupOption {
  value: number
  label: string
  description: string | null
  rate: number
  userRate: number | null
  subscriptionType: SubscriptionType
  platform: GroupPlatform
}

const appStore = useAppStore()
const onboardingStore = useOnboardingStore()
const { copyToClipboard: clipboardCopy } = useClipboard()
const showColumnDropdown = ref(false)
const columnDropdownRef = ref<HTMLElement | null>(null)

const zhKeysCopy = {
  kicker: '密钥庭册',
  title: '调用凭证',
  currentList: '当前清册',
  activeOnPage: '本页启用',
  groupsAvailable: '分组可用',
  selectKeyAria: '选择密钥',
  searchKeyPlaceholder: '搜索密钥名称或前缀',
  selectGroupAria: '选择分组',
  selectStatusAria: '选择状态',
  endpoint: '接入地址',
  copyEndpoint: '复制地址',
  connectionCheck: '接入体检',
  success24h: '24h 成功',
  winRate24h: '24h 胜率',
  viewRecentFailures: '查看最近失败',
  connectionConfig: '接入配置',
  copyCurl: '复制 curl',
  copyOpenAiSdk: '复制 OpenAI SDK',
  copyEnv: '复制环境变量',
  viewFullConfig: '查看全量配置',
  recentModelHints: '近期模型提示',
  checkRecord: '体检记录',
  copyReport: '复制报告',
  keyStatus: '密钥状态',
  modelsEndpoint: '模型接口',
  responseLatency: '响应耗时',
  reachable: '可访问',
  notFinished: '未完成',
  modelVisibility: '模型可见性',
  checking: '体检中…',
  startCheck: '开始体检',
  advancedControl: '高级控制',
  ledger: '账册',
  createdAt: '创建时间',
  healthStatus: '健康状态',
  allKeys: '全部密钥',
  ungrouped: '未分组',
  defaultLimits: '当前使用默认限制',
  advancedEnabled: (count: number) => `已启用 ${count} 项高级设置`,
  lastCall: (date: string) => `最近调用 ${date}`,
  noCalls: '暂无调用',
  requestPulse: (attempts: number, successes: number) => `24h 共 ${attempts} 次，请求成功 ${successes} 次`,
  noRequests24h: '24h 暂无请求',
  loadingRecent: '正在整理最近请求',
  groupOk: '分组正常',
  noGroup: '未绑定分组',
  inactive: '当前已停用',
  quotaExhausted: '当前额度已耗尽',
  expired: '当前已过期',
  copiedEndpoint: '接入地址已复制',
  errorLedgerUnavailable: '错误账册暂未开放，请联系管理员开启后查看。',
  copiedCurl: 'curl 已复制',
  copiedSdk: 'SDK 配置已复制',
  copiedEnv: '环境变量已复制',
  reportCopied: '诊断报告已复制',
  reportTitle: '山枢庭 SST 接入体检报告',
  checkedAt: '检查时间',
  conclusion: '结论',
  detail: '说明',
  visibleModelCount: '可见模型',
  sampleModels: '示例模型',
  suggestion: '建议',
  visibleModels: (count: number) => `${count} 个可见模型`,
  modelCount: (count: number) => `${count} 个`,
  joinList: '，',
  joinHints: '；',
  health: {
    disabledTitle: '密钥未启用',
    disabledDetail: '当前密钥未处于 active 状态，建议先启用，或换一枚可用密钥后再体检。',
    disabledAction: '先启用这枚密钥，或换一枚 active 状态的可用密钥后再检测。',
    noModelsTitle: '模型接口已连通，但当前未返回可用模型',
    noModelsDetail: '这把 Key 可以访问 /v1/models，但当前分组下没有返回可见模型，通常意味着模型列表未开放、线路未覆盖，或上游还没准备好。',
    noModelsAction: '优先检查当前 Key 绑定分组的模型列表配置，以及这条线路是否真的向用户公开了目标模型。',
    passedTitle: '接入体检通过',
    passedDetail: '模型接口可访问，当前接入地址与密钥可用于基础调用。',
    failedTitle: '接入体检未通过',
    unauthorized: (status: number) => `模型接口返回 HTTP ${status}，更像密钥未通过校验，或当前分组没有查看模型列表的权限。`,
    rateLimited: '模型接口返回 HTTP 429，当前更像线路或上游在限流窗口内。',
    serverError: (status: number) => `模型接口返回 HTTP ${status}，当前更像线路或上游暂时不可用。`,
    genericHttp: (status: number) => `模型接口返回 HTTP ${status}，请检查密钥权限、分组或接入地址。`,
    returned: '接口返回：',
    unauthorizedAction: '优先检查密钥是否启用、是否复制完整，以及分组权限是否允许访问模型列表。',
    rateLimitAction: '先拉开重试间隔；如果最近失败里已经提示某个模型限流，优先按那个模型处理。',
    genericAction: '建议稍后重试；若持续失败，可复制体检记录并结合最近模型提示一起排查。',
    failedToast: (status: number) => `接入体检未通过：HTTP ${status}`,
    unreachableTitle: '无法完成接入体检',
    unreachableDetail: '请求未能抵达模型接口，请检查接入地址、浏览器跨域限制或当前网络。',
    unreachableAction: '优先确认当前站点接入地址、浏览器网络和跨域限制；必要时复制体检记录给客服排查。',
    unreachableToast: '接入体检失败，请检查接入地址、浏览器跨域限制或当前网络'
  }
}

const enKeysCopy = {
  kicker: 'Key ledger',
  title: 'API credentials',
  currentList: 'Current list',
  activeOnPage: 'Active keys',
  groupsAvailable: 'Groups',
  selectKeyAria: 'Select key',
  searchKeyPlaceholder: 'Search key name or prefix',
  selectGroupAria: 'Select group',
  selectStatusAria: 'Select status',
  endpoint: 'Endpoint',
  copyEndpoint: 'Copy endpoint',
  connectionCheck: 'Connection check',
  success24h: '24h success',
  winRate24h: '24h win rate',
  viewRecentFailures: 'View recent failures',
  connectionConfig: 'Integration config',
  copyCurl: 'Copy curl',
  copyOpenAiSdk: 'Copy OpenAI SDK',
  copyEnv: 'Copy env vars',
  viewFullConfig: 'View full config',
  recentModelHints: 'Recent model hints',
  checkRecord: 'Check record',
  copyReport: 'Copy report',
  keyStatus: 'Key status',
  modelsEndpoint: 'Models endpoint',
  responseLatency: 'Response latency',
  reachable: 'Reachable',
  notFinished: 'Not finished',
  modelVisibility: 'Model visibility',
  checking: 'Checking...',
  startCheck: 'Start check',
  advancedControl: 'Advanced controls',
  ledger: 'Ledger',
  createdAt: 'Created',
  healthStatus: 'Health',
  allKeys: 'All keys',
  ungrouped: 'Ungrouped',
  defaultLimits: 'Using default limits',
  advancedEnabled: (count: number) => `${count} advanced settings enabled`,
  lastCall: (date: string) => `Last call ${date}`,
  noCalls: 'No calls yet',
  requestPulse: (attempts: number, successes: number) => `24h total ${attempts}, successful ${successes}`,
  noRequests24h: 'No requests in 24h',
  loadingRecent: 'Collecting recent requests',
  groupOk: 'Group OK',
  noGroup: 'No group bound',
  inactive: 'Currently disabled',
  quotaExhausted: 'Quota exhausted',
  expired: 'Expired',
  copiedEndpoint: 'Endpoint copied',
  errorLedgerUnavailable: 'Error ledger is not enabled. Ask an administrator to turn it on.',
  copiedCurl: 'curl copied',
  copiedSdk: 'SDK config copied',
  copiedEnv: 'Environment variables copied',
  reportCopied: 'Diagnostic report copied',
  reportTitle: 'SST connection check report',
  checkedAt: 'Checked at',
  conclusion: 'Conclusion',
  detail: 'Details',
  visibleModelCount: 'Visible models',
  sampleModels: 'Sample models',
  suggestion: 'Suggestion',
  visibleModels: (count: number) => `${count} visible models`,
  modelCount: (count: number) => `${count}`,
  joinList: ', ',
  joinHints: '; ',
  health: {
    disabledTitle: 'Key is not active',
    disabledDetail: 'This key is not in active status. Enable it first, or choose another active key for the check.',
    disabledAction: 'Enable this key, or switch to an active key and test again.',
    noModelsTitle: 'Models endpoint is reachable, but no models were returned',
    noModelsDetail: 'This key can access /v1/models, but the current group returned no visible models. The model list may be closed, the route may not cover it, or upstream may not be ready.',
    noModelsAction: 'Check the model list configuration for this key group, and confirm the route is actually exposed to users.',
    passedTitle: 'Connection check passed',
    passedDetail: 'The models endpoint is reachable. This endpoint and key are ready for basic calls.',
    failedTitle: 'Connection check failed',
    unauthorized: (status: number) => `The models endpoint returned HTTP ${status}. The key likely failed validation, or this group lacks permission to view the model list.`,
    rateLimited: 'The models endpoint returned HTTP 429. The route or upstream is likely in a rate-limit window.',
    serverError: (status: number) => `The models endpoint returned HTTP ${status}. The route or upstream may be temporarily unavailable.`,
    genericHttp: (status: number) => `The models endpoint returned HTTP ${status}. Check key permission, group, or endpoint configuration.`,
    returned: 'Response: ',
    unauthorizedAction: 'Check that the key is active, copied completely, and allowed to access the model list.',
    rateLimitAction: 'Increase retry intervals first. If recent failures mention a specific model rate limit, start there.',
    genericAction: 'Try again later. If it keeps failing, copy this check record and compare it with recent model hints.',
    failedToast: (status: number) => `Connection check failed: HTTP ${status}`,
    unreachableTitle: 'Connection check could not complete',
    unreachableDetail: 'The request did not reach the models endpoint. Check the endpoint address, browser CORS limits, or current network.',
    unreachableAction: 'Confirm the site endpoint, browser network, and CORS settings first. Copy the record for support if needed.',
    unreachableToast: 'Connection check failed. Check endpoint address, CORS limits, or current network.'
  }
}

const keysCopy = computed(() => locale.value === 'zh' ? zhKeysCopy : enKeysCopy)

const allColumns = computed<Column[]>(() => [
  { key: 'name', label: t('common.name'), sortable: true, class: 'min-w-[8rem]' },
  { key: 'id', label: t('keys.id'), sortable: true, class: 'min-w-[5rem]' },
  { key: 'key', label: t('keys.apiKey'), sortable: false, class: 'min-w-[8.75rem]' },
  { key: 'group', label: t('keys.group'), sortable: false, class: 'min-w-[8rem]' },
  { key: 'current_concurrency', label: t('keys.currentConcurrency'), sortable: true, class: 'min-w-[6.5rem]' },
  { key: 'usage', label: keysCopy.value.ledger, sortable: false, class: 'min-w-[12rem]' },
  { key: 'rate_limit', label: t('keys.rateLimitColumn'), sortable: false, class: 'min-w-[8rem]' },
  { key: 'expires_at', label: t('keys.expiresAt'), sortable: true, class: 'min-w-[8rem]' },
  { key: 'last_used_at', label: t('keys.lastUsedAt'), sortable: true, class: 'min-w-[8rem]' },
  { key: 'last_used_ip', label: t('keys.lastUsedIP'), sortable: false, class: 'min-w-[8rem]' },
  { key: 'created_at', label: keysCopy.value.createdAt, sortable: true, class: 'min-w-[6.25rem]' },
  { key: 'status', label: keysCopy.value.healthStatus, sortable: true, class: 'min-w-[7.5rem]' },
  { key: 'actions', label: t('common.actions'), sortable: false, class: 'min-w-[6rem]' }
])

const ALWAYS_VISIBLE_COLUMNS = new Set(['name', 'actions'])
const DEFAULT_HIDDEN_COLUMNS = ['id', 'rate_limit', 'last_used_at', 'last_used_ip']
const HIDDEN_COLUMNS_KEY = 'api-key-hidden-columns'
const COLUMN_SETTINGS_VERSION_KEY = 'api-key-column-settings-version'
const COLUMN_SETTINGS_VERSION = 3
const VERSION_NEW_HIDDEN_COLUMNS: Record<number, string[]> = {
  2: ['last_used_ip'],
  3: ['id']
}

const toggleableColumns = computed(() =>
  allColumns.value.filter((col) => !ALWAYS_VISIBLE_COLUMNS.has(col.key))
)

const hiddenColumns = reactive<Set<string>>(new Set())

const saveColumnsToStorage = () => {
  try {
    localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
    localStorage.setItem(COLUMN_SETTINGS_VERSION_KEY, String(COLUMN_SETTINGS_VERSION))
  } catch (error) {
    console.error('Failed to save API key table columns:', error)
  }
}

const loadSavedColumns = () => {
  hiddenColumns.clear()
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      const parsed = JSON.parse(saved) as string[]
      const validColumnKeys = new Set(allColumns.value.map((col) => col.key))
      parsed
        .filter((key) =>
          typeof key === 'string' &&
          validColumnKeys.has(key) &&
          !ALWAYS_VISIBLE_COLUMNS.has(key)
        )
        .forEach((key) => hiddenColumns.add(key))
      const storedVersion = Number(localStorage.getItem(COLUMN_SETTINGS_VERSION_KEY) ?? '1')
      if (storedVersion < COLUMN_SETTINGS_VERSION) {
        for (let v = storedVersion + 1; v <= COLUMN_SETTINGS_VERSION; v++) {
          for (const key of VERSION_NEW_HIDDEN_COLUMNS[v] ?? []) {
            if (validColumnKeys.has(key) && !ALWAYS_VISIBLE_COLUMNS.has(key)) {
              hiddenColumns.add(key)
            }
          }
        }
        saveColumnsToStorage()
      } else {
        localStorage.setItem(COLUMN_SETTINGS_VERSION_KEY, String(COLUMN_SETTINGS_VERSION))
      }
    } else {
      DEFAULT_HIDDEN_COLUMNS.forEach((key) => hiddenColumns.add(key))
      localStorage.setItem(COLUMN_SETTINGS_VERSION_KEY, String(COLUMN_SETTINGS_VERSION))
    }
  } catch (error) {
    console.error('Failed to load API key table columns:', error)
    DEFAULT_HIDDEN_COLUMNS.forEach((key) => hiddenColumns.add(key))
  }
}

const toggleColumn = (key: string) => {
  if (ALWAYS_VISIBLE_COLUMNS.has(key)) return
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  } else {
    hiddenColumns.add(key)
  }
  saveColumnsToStorage()
}

const isColumnVisible = (key: string) => !hiddenColumns.has(key)

const columns = computed<Column[]>(() =>
  allColumns.value.filter((col) => ALWAYS_VISIBLE_COLUMNS.has(col.key) || !hiddenColumns.has(col.key))
)
const apiKeys = ref<ApiKey[]>([])
const activeKeyCount = computed(() => apiKeys.value.filter((key) => key.status === 'active').length)
const firstTestableKey = computed(function () { return apiKeys.value.find(function (key) { return key.status === 'active' }) || apiKeys.value[0] || null })
const selectedConnectionTestKey = computed(function () { return apiKeys.value.find(function (key) { return key.id === connectionTestKeyId.value }) || firstTestableKey.value })
const selectedConnectionWorkbenchSummary = computed(() => {
  if (!selectedConnectionTestKey.value) return null
  return getWorkbenchStat(selectedConnectionTestKey.value.id)
})
const selectedConnectionModelHints = computed(() => buildWorkbenchModelHints(selectedConnectionWorkbenchSummary.value?.latest_error, locale.value))
const selectedConnectionLatestErrorLabel = computed(() => describeLatestError(selectedConnectionWorkbenchSummary.value))
const formatDateOnly = (value: string | null | undefined) => value ? new Intl.DateTimeFormat('zh-CN', {
  timeZone: 'Asia/Shanghai',
  year: 'numeric',
  month: '2-digit',
  day: '2-digit'
}).format(new Date(value)).replace(/\//g, '-') : '-'

const keyReasonLabels: Record<string, string> = {
  auth_key_deleted: '旧 Key 已删除',
  auth_invalid_credentials: '凭证未通过校验',
  quota_balance_exhausted: '余额或额度耗尽',
  quota_subscription_exhausted: '套餐或订阅失效',
  rate_limit_window_exhausted: '速率窗口已打满',
  request_model_not_supported: '模型不可用',
  request_payload_too_large: '请求体过大',
  request_invalid: '请求参数不匹配',
  service_model_not_available: '当前线路未开放该模型',
  service_model_rate_limited: '该模型当前线路已限流',
  service_no_route_available: '当前无可用线路',
  upstream_temporarily_unavailable: '上游临时不可用',
  upstream_transport_error: '上游链路异常',
  internal_gateway_error: '平台内部异常',
  cyber_policy_blocked: '触发安全策略',
}

const keyCategoryLabels: Record<string, string> = {
  auth: '认证失败',
  rate_limit: '限流',
  quota: '额度不足',
  invalid_request: '参数错误',
  service_unavailable: '线路暂不可用',
  upstream: '上游错误',
  internal: '平台错误',
  cyber: '安全策略',
  other: '其他异常',
}

const enKeyReasonLabels: Record<string, string> = {
  auth_key_deleted: 'Old key deleted',
  auth_invalid_credentials: 'Credentials failed validation',
  quota_balance_exhausted: 'Balance or quota exhausted',
  quota_subscription_exhausted: 'Plan or subscription inactive',
  rate_limit_window_exhausted: 'Rate-limit window exhausted',
  request_model_not_supported: 'Model unavailable',
  request_payload_too_large: 'Payload too large',
  request_invalid: 'Request parameters mismatch',
  service_model_not_available: 'Model not enabled on this route',
  service_model_rate_limited: 'Model route is rate-limited',
  service_no_route_available: 'No available route',
  upstream_temporarily_unavailable: 'Upstream temporarily unavailable',
  upstream_transport_error: 'Upstream transport error',
  internal_gateway_error: 'Gateway internal error',
  cyber_policy_blocked: 'Blocked by security policy',
}

const enKeyCategoryLabels: Record<string, string> = {
  auth: 'Auth failed',
  rate_limit: 'Rate limited',
  quota: 'Insufficient quota',
  invalid_request: 'Invalid request',
  service_unavailable: 'Route unavailable',
  upstream: 'Upstream error',
  internal: 'Gateway error',
  cyber: 'Security policy',
  other: 'Other error',
}

const describeLatestError = (summary?: ApiKeyWorkbenchSummary | null) => {
  return buildWorkbenchLatestErrorLabel(
    summary?.latest_error,
    locale.value.startsWith('zh') ? keyReasonLabels : enKeyReasonLabels,
    locale.value.startsWith('zh') ? keyCategoryLabels : enKeyCategoryLabels,
    locale.value,
  )
}

const getWorkbenchStat = (keyId: number) => workbenchStats.value[keyId] || workbenchStats.value[String(keyId)] || null

const getKeyHealth = (key: ApiKey) => {
  const summary = getWorkbenchStat(key.id)
  const lastCall = key.last_used_at ? keysCopy.value.lastCall(formatDateOnly(key.last_used_at)) : keysCopy.value.noCalls
  const requestPulse = summary
    ? summary.attempt_count_24h > 0
      ? keysCopy.value.requestPulse(summary.attempt_count_24h, summary.success_requests_24h)
      : keysCopy.value.noRequests24h
    : keysCopy.value.loadingRecent
  const group = key.group ? keysCopy.value.groupOk : keysCopy.value.noGroup
  const statusFallback = key.status === 'inactive'
    ? keysCopy.value.inactive
    : key.status === 'quota_exhausted'
      ? keysCopy.value.quotaExhausted
      : key.status === 'expired'
        ? keysCopy.value.expired
        : describeLatestError(summary)
  const summaryLine = `${group} / ${statusFallback}`

  return {
    lastCall,
    requestPulse,
    summary: summaryLine,
    modelHints: buildWorkbenchModelHints(summary?.latest_error, locale.value),
    hasAttention: !!summary?.latest_error || key.status !== 'active' || !key.group,
    canReview: !!summary?.latest_error && canViewErrorRequests.value,
  }
}

const formatUsd = (value: number | null | undefined, digits = 2) => `$${(value ?? 0).toFixed(digits)}`

const formatSuccessRate = (value: number | null | undefined) => {
  if (value === null || value === undefined || Number.isNaN(value)) return '-'
  return `${(value * 100).toFixed(value >= 0.995 ? 0 : 1)}%`
}

const getMeterTone = (used: number | null | undefined, limit: number | null | undefined) => {
  const current = used ?? 0
  const total = limit ?? 0
  if (total <= 0) return 'muted'
  if (current >= total) return 'danger'
  if (current >= total * 0.8) return 'warning'
  return 'safe'
}

const getRateLimitWindows = (row: ApiKey) => ([
  {
    key: '5h',
    label: '5h',
    usage: row.usage_5h ?? 0,
    limit: row.rate_limit_5h ?? 0,
    resetAt: row.reset_5h_at ?? null,
  },
  {
    key: '1d',
    label: '1d',
    usage: row.usage_1d ?? 0,
    limit: row.rate_limit_1d ?? 0,
    resetAt: row.reset_1d_at ?? null,
  },
  {
    key: '7d',
    label: '7d',
    usage: row.usage_7d ?? 0,
    limit: row.rate_limit_7d ?? 0,
    resetAt: row.reset_7d_at ?? null,
  }
]).filter((window) => window.limit > 0)

type ConnectionTestResult = {
  tone: 'success' | 'warning' | 'danger'
  title: string
  detail: string
  latencyMs: number | null
  statusCode: number | null
  checkedAt: string
  action: string
  availableModelCount: number | null
  allModels: string[]
  sampleModels: string[]
}
const groups = ref<Group[]>([])
const loading = ref(false)
const submitting = ref(false)
const workbenchStats = ref<Record<string, ApiKeyWorkbenchSummary>>({})
const userGroupRates = ref<Record<number, number>>({})
const filterKeyCatalog = ref<ApiKey[]>([])

const pagination = ref({
  page: 1,
  page_size: getPersistedPageSize(),
  total: 0,
  pages: 0
})
const sortState = ref({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
})

// Filter state
const filterSearch = ref('')
const filterStatus = ref('')
const filterGroupId = ref<string | number>('')
const selectedFilterKeyId = ref<string | number | null>('')

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)
const showResetQuotaDialog = ref(false)
const showResetRateLimitDialog = ref(false)
const showUseKeyModal = ref(false)
const showConnectionTestDialog = ref(false)
const advancedSettingsExpanded = ref(false)
const connectionTestKeyId = ref<number | null>(null)
const showCcsClientSelect = ref(false)
const pendingCcsRow = ref<ApiKey | null>(null)
const selectedKey = ref<ApiKey | null>(null)
const copiedKeyId = ref<number | null>(null)
const groupSelectorKeyId = ref<number | null>(null)
const publicSettings = ref<PublicSettings | null>(null)
const canViewErrorRequests = computed(() => publicSettings.value?.allow_user_view_error_requests ?? false)
const testingKeyId = ref<number | null>(null)
const connectionTestResult = ref<ConnectionTestResult | null>(null)
const showAllConnectionModels = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const dropdownPosition = ref<{ top?: number; bottom?: number; left: number } | null>(null)
const groupButtonRefs = ref<Map<number, HTMLElement>>(new Map())
let abortController: AbortController | null = null
let workbenchAbortController: AbortController | null = null

const visibleConnectionModels = computed(() => {
  if (!connectionTestResult.value) return [] as string[]
  return showAllConnectionModels.value
    ? connectionTestResult.value.allModels
    : connectionTestResult.value.sampleModels
})

const hiddenConnectionModelCount = computed(() => {
  if (!connectionTestResult.value) return 0
  return Math.max(0, connectionTestResult.value.allModels.length - connectionTestResult.value.sampleModels.length)
})

// Get the currently selected key for group change
const selectedKeyForGroup = computed(() => {
  if (groupSelectorKeyId.value === null) return null
  return apiKeys.value.find((k) => k.id === groupSelectorKeyId.value) || null
})

const GROUP_SELECTOR_MIN_WIDTH = 380
const GROUP_SELECTOR_ESTIMATED_HEIGHT = 360
const GROUP_SELECTOR_VIEWPORT_GAP = 12
const GROUP_SELECTOR_OFFSET = 8

const setGroupButtonRef = (keyId: number, el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    groupButtonRefs.value.set(keyId, el)
  } else {
    groupButtonRefs.value.delete(keyId)
  }
}

const formData = ref({
  name: '',
  group_id: null as number | null,
  status: 'active' as 'active' | 'inactive',
  use_custom_key: false,
  custom_key: '',
  enable_ip_restriction: false,
  ip_whitelist: '',
  ip_blacklist: '',
  // Quota settings (empty = unlimited)
  enable_quota: false,
  quota: null as number | null,
  // Rate limit settings
  enable_rate_limit: false,
  rate_limit_5h: null as number | null,
  rate_limit_1d: null as number | null,
  rate_limit_7d: null as number | null,
  enable_expiration: false,
  expiration_preset: '30' as '7' | '30' | '90' | 'custom',
  expiration_date: ''
})

// 自定义Key验证
const customKeyError = computed(() => {
  if (!formData.value.use_custom_key || !formData.value.custom_key) {
    return ''
  }
  const key = formData.value.custom_key
  if (key.length < 16) {
    return t('keys.customKeyTooShort')
  }
  // 检查字符：只允许字母、数字、下划线、连字符
  if (!/^[a-zA-Z0-9_-]+$/.test(key)) {
    return t('keys.customKeyInvalidChars')
  }
  return ''
})

const advancedSummaryItems = computed(() => {
  const items: string[] = []
  if (!showEditModal.value && formData.value.use_custom_key) items.push(t('keys.customKeyLabel'))
  if (formData.value.enable_ip_restriction) items.push(t('keys.ipRestriction'))
  if (formData.value.quota && formData.value.quota > 0) items.push(t('keys.quotaLimit'))
  if (formData.value.enable_rate_limit) items.push(t('keys.rateLimitSection'))
  if (formData.value.enable_expiration) items.push(t('keys.expiration'))
  return items
})

const advancedSettingsSummary = computed(() => {
  if (advancedSummaryItems.value.length === 0) {
    return keysCopy.value.defaultLimits
  }
  return keysCopy.value.advancedEnabled(advancedSummaryItems.value.length)
})

const statusOptions = computed(() => [
  { value: 'active', label: t('common.active') },
  { value: 'inactive', label: t('common.inactive') }
])

// Filter dropdown options
const groupFilterOptions = computed(() => [
  { value: '', label: t('keys.allGroups') },
  { value: 0, label: t('keys.noGroup') },
  ...groups.value.map((g) => ({ value: g.id, label: g.name }))
])

const statusFilterOptions = computed(() => [
  { value: '', label: t('keys.allStatus') },
  { value: 'active', label: t('keys.status.active') },
  { value: 'inactive', label: t('keys.status.inactive') },
  { value: 'quota_exhausted', label: t('keys.status.quota_exhausted') },
  { value: 'expired', label: t('keys.status.expired') }
])

const keyFilterOptions = computed(() => {
  const allOption = { value: '', label: keysCopy.value.allKeys }
  const source = filterKeyCatalog.value.length ? filterKeyCatalog.value : apiKeys.value
  const options = source.map((key) => ({
    value: key.id,
    label: `${key.name} · ${maskApiKey(key.key)}`,
    description: key.group?.name || keysCopy.value.ungrouped,
  }))
  return [allOption, ...options]
})

const onFilterChange = () => {
  pagination.value.page = 1
  loadApiKeys()
}

const onKeyFilterChange = (value: string | number | boolean | null) => {
  selectedFilterKeyId.value = value as string | number | null
  const selectedKey = typeof value === 'number'
    ? apiKeys.value.find((key) => key.id === value)
    : null
  filterSearch.value = selectedKey ? `${selectedKey.name} ${selectedKey.key}` : ''
  onFilterChange()
}

const onGroupFilterChange = (value: string | number | boolean | null) => {
  filterGroupId.value = value as string | number
  onFilterChange()
}

const onStatusFilterChange = (value: string | number | boolean | null) => {
  filterStatus.value = value as string
  onFilterChange()
}

// Convert groups to Select options format with rate multiplier and subscription type
const groupOptions = computed(() =>
  groups.value.map((group) => ({
    value: group.id,
    label: group.name,
    description: group.description,
    rate: group.rate_multiplier,
    userRate: userGroupRates.value[group.id] ?? null,
    subscriptionType: group.subscription_type,
    platform: group.platform
  }))
)

// Group dropdown search
const groupSearchQuery = ref('')
const resolvedApiBaseUrl = computed(() => (publicSettings.value?.api_base_url || window.location.origin).replace(/\/$/, ''))

const filteredGroupOptions = computed(() => {
  const query = groupSearchQuery.value.trim().toLowerCase()
  if (!query) return groupOptions.value
  return groupOptions.value.filter((opt) => {
    return opt.label.toLowerCase().includes(query) ||
      (opt.description && opt.description.toLowerCase().includes(query))
  })
})

const copyToClipboard = async (text: string, keyId: number) => {
  const success = await clipboardCopy(text, t('keys.copied'))
  if (success) {
    copiedKeyId.value = keyId
    setTimeout(() => {
      copiedKeyId.value = null
    }, 800)
  }
}

const copyEndpoint = async () => {
  await clipboardCopy(resolvedApiBaseUrl.value, keysCopy.value.copiedEndpoint)
}

const loadKeyWorkbenchStats = async (keys: ApiKey[]) => {
  if (workbenchAbortController) workbenchAbortController.abort()
  workbenchAbortController = new AbortController()
  if (!keys.length) {
    workbenchStats.value = {}
    return
  }
  try {
    const response = await usageAPI.getDashboardApiKeysWorkbench(
      keys.map((key) => key.id),
      { signal: workbenchAbortController.signal }
    )
    workbenchStats.value = response.stats || {}
  } catch (error: any) {
    if (error?.name === 'AbortError' || error?.code === 'ERR_CANCELED') return
    console.error('Failed to load key workbench stats:', error)
    workbenchStats.value = {}
  }
}

const loadApiKeys = async () => {
  if (abortController) abortController.abort()
  abortController = new AbortController()
  loading.value = true
  try {
    const groupId = filterGroupId.value === '' ? undefined : filterGroupId.value
    const response = await keysAPI.list(pagination.value.page, pagination.value.page_size, {
      search: filterSearch.value || undefined,
      status: filterStatus.value || undefined,
      group_id: groupId,
      sort_by: sortState.value.sort_by,
      sort_order: sortState.value.sort_order,
    }, { signal: abortController.signal })
    apiKeys.value = response.items || []
    if (typeof selectedFilterKeyId.value === 'number') {
      const stillExists = apiKeys.value.some((key) => key.id === selectedFilterKeyId.value)
      if (!stillExists) {
        selectedFilterKeyId.value = ''
      }
    }
    pagination.value.total = response.total || 0
    pagination.value.pages = response.pages || 0
    await loadKeyWorkbenchStats(apiKeys.value)
  } catch (error: any) {
    if (error?.name !== 'AbortError' && error?.code !== 'ERR_CANCELED') appStore.showError(t('keys.failedToLoad'))
  } finally {
    loading.value = false
  }
}

const loadGroups = async () => {
  try { groups.value = await userGroupsAPI.getAvailable() } catch (error) { console.error('Failed to load groups:', error) }
}

const loadUserGroupRates = async () => {
  try { userGroupRates.value = await userGroupsAPI.getUserGroupRates() } catch (error) { console.error('Failed to load user group rates:', error) }
}

const loadPublicSettings = async () => {
  try { publicSettings.value = await authAPI.getPublicSettings() } catch (error) { console.error('Failed to load public settings:', error) }
}

const loadFilterKeyCatalog = async () => {
  try {
    const response = await keysAPI.list(1, 100, {
      sort_by: 'created_at',
      sort_order: 'desc'
    })
    filterKeyCatalog.value = response.items || []
  } catch (error) {
    console.error('Failed to load key filter catalog:', error)
  }
}

const handleSort = (key: string, order: 'asc' | 'desc') => {
  sortState.value = { sort_by: key, sort_order: order }
  pagination.value.page = 1
  loadApiKeys()
}

const handlePageChange = (page: number) => { pagination.value.page = page; loadApiKeys() }
const handlePageSizeChange = (pageSize: number) => { pagination.value.page_size = pageSize; pagination.value.page = 1; loadApiKeys() }
const openUseKeyModal = (key: ApiKey) => { selectedKey.value = key; showUseKeyModal.value = true }
const closeUseKeyModal = () => { showUseKeyModal.value = false; selectedKey.value = null }
const openKeyErrorLedger = (apiKeyId: number) => {
  if (!canViewErrorRequests.value) {
    appStore.showError(keysCopy.value.errorLedgerUnavailable)
    return
  }
  router.push({ path: '/usage', query: { tab: 'errors', api_key_id: String(apiKeyId) } })
}

const toggleKeyStatus = async (key: ApiKey) => {
  try { await keysAPI.toggleStatus(key.id, key.status === 'active' ? 'inactive' : 'active'); appStore.showSuccess(t('common.success')); await loadApiKeys() }
  catch (error: any) { appStore.showError(error?.message || t('common.error')) }
}

const openCreateModal = () => {
  closeModals()
  showCreateModal.value = true
  advancedSettingsExpanded.value = false
}

const editKey = (key: ApiKey) => {
  selectedKey.value = key
  formData.value = {
    name: key.name,
    group_id: key.group_id ?? null,
    status: key.status as 'active' | 'inactive',
    use_custom_key: false,
    custom_key: '',
    enable_ip_restriction: !!((key.ip_whitelist?.length || 0) > 0 || (key.ip_blacklist?.length || 0) > 0),
    ip_whitelist: (key.ip_whitelist || []).join('\n'),
    ip_blacklist: (key.ip_blacklist || []).join('\n'),
    enable_quota: (key.quota || 0) > 0,
    quota: (key.quota || 0) > 0 ? key.quota : null,
    enable_rate_limit: (key.rate_limit_5h || 0) > 0 || (key.rate_limit_1d || 0) > 0 || (key.rate_limit_7d || 0) > 0,
    rate_limit_5h: (key.rate_limit_5h || 0) > 0 ? key.rate_limit_5h : null,
    rate_limit_1d: (key.rate_limit_1d || 0) > 0 ? key.rate_limit_1d : null,
    rate_limit_7d: (key.rate_limit_7d || 0) > 0 ? key.rate_limit_7d : null,
    enable_expiration: !!key.expires_at,
    expiration_preset: 'custom',
    expiration_date: key.expires_at ? formatDateTimeLocal(key.expires_at) : ''
  }
  advancedSettingsExpanded.value = false
  showEditModal.value = true
}

const confirmDelete = (key: ApiKey) => { selectedKey.value = key; showDeleteDialog.value = true }

const updateGroupSelectorPosition = () => {
  if (groupSelectorKeyId.value === null) {
    dropdownPosition.value = null
    return
  }

  const buttonEl = groupButtonRefs.value.get(groupSelectorKeyId.value)
  if (!buttonEl) {
    dropdownPosition.value = null
    return
  }

  const rect = buttonEl.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  const menuWidth = Math.max(GROUP_SELECTOR_MIN_WIDTH, rect.width)
  const maxLeft = Math.max(GROUP_SELECTOR_VIEWPORT_GAP, viewportWidth - menuWidth - GROUP_SELECTOR_VIEWPORT_GAP)
  const left = Math.min(Math.max(GROUP_SELECTOR_VIEWPORT_GAP, rect.left), maxLeft)
  const shouldOpenAbove =
    rect.bottom + GROUP_SELECTOR_OFFSET + GROUP_SELECTOR_ESTIMATED_HEIGHT > viewportHeight - GROUP_SELECTOR_VIEWPORT_GAP &&
    rect.top > GROUP_SELECTOR_ESTIMATED_HEIGHT / 2

  dropdownPosition.value = shouldOpenAbove
    ? {
        bottom: Math.max(GROUP_SELECTOR_VIEWPORT_GAP, viewportHeight - rect.top + GROUP_SELECTOR_OFFSET),
        left
      }
    : {
        top: rect.bottom + GROUP_SELECTOR_OFFSET,
        left
      }
}

const openGroupSelector = async (key: ApiKey) => {
  if (groupSelectorKeyId.value === key.id) {
    closeGroupSelector()
    return
  }

  groupSelectorKeyId.value = key.id
  groupSearchQuery.value = ''
  await nextTick()
  updateGroupSelectorPosition()
}

const closeGroupSelector = () => {
  groupSelectorKeyId.value = null
  groupSearchQuery.value = ''
  dropdownPosition.value = null
}

const handleGroupSelectorViewportChange = () => {
  if (groupSelectorKeyId.value !== null) {
    updateGroupSelectorPosition()
  }
}

const handleDocumentClick = (event: MouseEvent) => {
  const target = event.target
  if (!(target instanceof Node)) {
    closeGroupSelector()
    showColumnDropdown.value = false
    return
  }

  if (showColumnDropdown.value && !columnDropdownRef.value?.contains(target)) {
    showColumnDropdown.value = false
  }

  if (groupSelectorKeyId.value !== null) {
    if (dropdownRef.value?.contains(target)) return

    const activeButton = groupButtonRefs.value.get(groupSelectorKeyId.value)
    if (activeButton?.contains(target)) return

    closeGroupSelector()
  }
}

const changeGroup = async (key: ApiKey, groupId: number | string | boolean | null) => {
  if (typeof groupId !== 'number') return
  if (key.group_id === groupId) {
    closeGroupSelector()
    return
  }
  try {
    await keysAPI.update(key.id, { group_id: groupId === 0 ? null : groupId })
    appStore.showSuccess(t('keys.groupChangedSuccess'))
    closeGroupSelector()
    await loadApiKeys()
  }
  catch (error: any) {
    appStore.showError(error?.message || t('keys.failedToChangeGroup'))
  }
}

const selectConnectionTestKey = (keyId: number) => {
  connectionTestKeyId.value = keyId
  connectionTestResult.value = null
  showAllConnectionModels.value = false
}

const openConnectionTestDialog = () => {
  connectionTestKeyId.value = firstTestableKey.value?.id ?? null
  connectionTestResult.value = null
  showAllConnectionModels.value = false
  showConnectionTestDialog.value = true
}

const closeConnectionTestDialog = () => {
  showConnectionTestDialog.value = false
  connectionTestResult.value = null
  showAllConnectionModels.value = false
}

const toggleConnectionModelsExpanded = () => {
  showAllConnectionModels.value = !showAllConnectionModels.value
}

const buildIntegrationSnippet = (type: 'curl' | 'openai' | 'env', key: ApiKey) => {
  const modelsUrl = `${resolvedApiBaseUrl.value}/v1/models`
  if (type === 'curl') {
    return `curl ${modelsUrl} \\
  -H "Authorization: Bearer ${key.key}"`
  }

  if (type === 'env') {
    return `OPENAI_API_KEY=${key.key}
OPENAI_BASE_URL=${resolvedApiBaseUrl.value}/v1`
  }

  return `import OpenAI from 'openai'

const client = new OpenAI({
  apiKey: '${key.key}',
  baseURL: '${resolvedApiBaseUrl.value}/v1'
})

const models = await client.models.list()`
}

const copyIntegrationSnippet = async (type: 'curl' | 'openai' | 'env') => {
  if (!selectedConnectionTestKey.value) return
  const message = type === 'curl' ? keysCopy.value.copiedCurl : type === 'openai' ? keysCopy.value.copiedSdk : keysCopy.value.copiedEnv
  await clipboardCopy(buildIntegrationSnippet(type, selectedConnectionTestKey.value), message)
}

const buildConnectionDiagnosticReport = () => {
  const key = selectedConnectionTestKey.value
  const result = connectionTestResult.value
  if (!key || !result) return ''
  const lines = [
    keysCopy.value.reportTitle,
    keysCopy.value.checkedAt + ': ' + new Date(result.checkedAt).toLocaleString(),
    'API Key：' + key.name + ' (' + maskApiKey(key.key) + ')',
    keysCopy.value.endpoint + ': ' + resolvedApiBaseUrl.value + '/v1',
    keysCopy.value.keyStatus + ': ' + key.status,
    keysCopy.value.modelsEndpoint + ': ' + (result.statusCode ? 'HTTP ' + result.statusCode : result.tone === 'success' ? keysCopy.value.reachable : keysCopy.value.notFinished),
    keysCopy.value.responseLatency + ': ' + (result.latencyMs !== null ? result.latencyMs + 'ms' : '-'),
    keysCopy.value.conclusion + ': ' + result.title,
    keysCopy.value.detail + ': ' + result.detail,
  ]
  if (result.availableModelCount !== null) {
    lines.push(keysCopy.value.visibleModelCount + ': ' + keysCopy.value.modelCount(result.availableModelCount))
  }
  if (result.sampleModels.length > 0) {
    lines.push(keysCopy.value.sampleModels + ': ' + result.sampleModels.join(keysCopy.value.joinList))
  }
  if (selectedConnectionModelHints.value.length > 0) {
    lines.push(keysCopy.value.recentModelHints + ': ' + selectedConnectionModelHints.value.join(keysCopy.value.joinHints))
  }
  if (result.action) {
    lines.push(keysCopy.value.suggestion + ': ' + result.action)
  }
  return lines.join('\n')
}

const copyConnectionDiagnosticReport = async () => {
  const report = buildConnectionDiagnosticReport()
  if (!report) return
  await clipboardCopy(report, keysCopy.value.reportCopied)
}

const openUseKeyModalFromConnectionTest = () => {
  if (!selectedConnectionTestKey.value) return
  selectedKey.value = selectedConnectionTestKey.value
  showConnectionTestDialog.value = false
  connectionTestResult.value = null
  showUseKeyModal.value = true
}

const testKeyConnection = async (key: ApiKey | null) => {
  if (!key || testingKeyId.value) return
  testingKeyId.value = key.id
  connectionTestResult.value = null
  showAllConnectionModels.value = false

  if (key.status !== 'active') {
    connectionTestResult.value = {
      tone: 'warning',
      title: keysCopy.value.health.disabledTitle,
      detail: keysCopy.value.health.disabledDetail,
      latencyMs: null,
      statusCode: null,
      checkedAt: new Date().toISOString(),
      action: keysCopy.value.health.disabledAction,
      availableModelCount: null,
      allModels: [],
      sampleModels: [],
    }
    testingKeyId.value = null
    return
  }

  const startedAt = performance.now()
  try {
    const response = await fetch(resolvedApiBaseUrl.value + '/v1/models', {
      headers: { Authorization: 'Bearer ' + key.key }
    })
    const latencyMs = Math.round(performance.now() - startedAt)
    const rawBody = await response.text()
    const payload = safeParseJson(rawBody)
    if (response.ok) {
      const modelIds = extractModelIdsFromModelsPayload(payload)
      if (modelIds.length === 0) {
        connectionTestResult.value = {
          tone: 'warning',
          title: keysCopy.value.health.noModelsTitle,
          detail: keysCopy.value.health.noModelsDetail,
          latencyMs,
          statusCode: response.status,
          checkedAt: new Date().toISOString(),
          action: keysCopy.value.health.noModelsAction,
          availableModelCount: 0,
          allModels: [],
          sampleModels: [],
        }
        appStore.showError(keysCopy.value.health.noModelsTitle)
        return
      }
      connectionTestResult.value = {
        tone: 'success',
        title: keysCopy.value.health.passedTitle,
        detail: keysCopy.value.health.passedDetail,
        latencyMs,
        statusCode: response.status,
        checkedAt: new Date().toISOString(),
        action: '',
        availableModelCount: modelIds.length,
        allModels: modelIds,
        sampleModels: modelIds.slice(0, 5),
      }
      appStore.showSuccess(keysCopy.value.health.passedTitle)
      return
    }
    const errorText = extractErrorTextFromPayload(payload) || rawBody.trim()
    const statusDetail = response.status === 401 || response.status === 403
      ? keysCopy.value.health.unauthorized(response.status)
      : response.status === 429
        ? keysCopy.value.health.rateLimited
        : response.status >= 500
          ? keysCopy.value.health.serverError(response.status)
          : keysCopy.value.health.genericHttp(response.status)
    connectionTestResult.value = {
      tone: response.status === 401 || response.status === 403 ? 'danger' : 'warning',
      title: keysCopy.value.health.failedTitle,
      detail: errorText ? `${statusDetail} ${keysCopy.value.health.returned}${errorText}` : statusDetail,
      latencyMs,
      statusCode: response.status,
      checkedAt: new Date().toISOString(),
      action: response.status === 401 || response.status === 403
        ? keysCopy.value.health.unauthorizedAction
        : response.status === 429
          ? keysCopy.value.health.rateLimitAction
          : keysCopy.value.health.genericAction,
      availableModelCount: null,
      allModels: [],
      sampleModels: [],
    }
    appStore.showError(keysCopy.value.health.failedToast(response.status))
  } catch (error) {
    connectionTestResult.value = {
      tone: 'danger',
      title: keysCopy.value.health.unreachableTitle,
      detail: keysCopy.value.health.unreachableDetail,
      latencyMs: Math.round(performance.now() - startedAt),
      statusCode: null,
      checkedAt: new Date().toISOString(),
      action: keysCopy.value.health.unreachableAction,
      availableModelCount: null,
      allModels: [],
      sampleModels: [],
    }
    appStore.showError(keysCopy.value.health.unreachableToast)
  } finally {
    testingKeyId.value = null
  }
}

const handleSubmit = async () => {
  // Validate group_id is required
  if (formData.value.group_id === null) {
    appStore.showError(t('keys.groupRequired'))
    return
  }

  // Validate custom key if enabled
  if (!showEditModal.value && formData.value.use_custom_key) {
    if (!formData.value.custom_key) {
      appStore.showError(t('keys.customKeyRequired'))
      return
    }
    if (customKeyError.value) {
      appStore.showError(customKeyError.value)
      return
    }
  }

  // Parse IP lists only if IP restriction is enabled
  const parseIPList = (text: string): string[] =>
    text.split('\n').map(ip => ip.trim()).filter(ip => ip.length > 0)
  const ipWhitelist = formData.value.enable_ip_restriction ? parseIPList(formData.value.ip_whitelist) : []
  const ipBlacklist = formData.value.enable_ip_restriction ? parseIPList(formData.value.ip_blacklist) : []

  // Calculate quota value (null/empty/0 = unlimited, stored as 0)
  const quota = formData.value.quota && formData.value.quota > 0 ? formData.value.quota : 0

  // Calculate expiration
  let expiresInDays: number | undefined
  let expiresAt: string | null | undefined
  if (formData.value.enable_expiration && formData.value.expiration_date) {
    if (!showEditModal.value) {
      // Create mode: calculate days from date
      const expDate = new Date(formData.value.expiration_date)
      const now = new Date()
      const diffDays = Math.ceil((expDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
      expiresInDays = diffDays > 0 ? diffDays : 1
    } else {
      // Edit mode: use custom date directly
      expiresAt = new Date(formData.value.expiration_date).toISOString()
    }
  } else if (showEditModal.value) {
    // Edit mode: if expiration disabled or date cleared, send empty string to clear
    expiresAt = ''
  }

  // Calculate rate limit values (send 0 when toggle is off)
  const rateLimitData = formData.value.enable_rate_limit ? {
    rate_limit_5h: formData.value.rate_limit_5h && formData.value.rate_limit_5h > 0 ? formData.value.rate_limit_5h : 0,
    rate_limit_1d: formData.value.rate_limit_1d && formData.value.rate_limit_1d > 0 ? formData.value.rate_limit_1d : 0,
    rate_limit_7d: formData.value.rate_limit_7d && formData.value.rate_limit_7d > 0 ? formData.value.rate_limit_7d : 0,
  } : { rate_limit_5h: 0, rate_limit_1d: 0, rate_limit_7d: 0 }

  submitting.value = true
  try {
    if (showEditModal.value && selectedKey.value) {
      await keysAPI.update(selectedKey.value.id, {
        name: formData.value.name,
        group_id: formData.value.group_id,
        status: formData.value.status,
        ip_whitelist: ipWhitelist,
        ip_blacklist: ipBlacklist,
        quota: quota,
        expires_at: expiresAt,
        rate_limit_5h: rateLimitData.rate_limit_5h,
        rate_limit_1d: rateLimitData.rate_limit_1d,
        rate_limit_7d: rateLimitData.rate_limit_7d,
      })
      appStore.showSuccess(t('keys.keyUpdatedSuccess'))
    } else {
      const customKey = formData.value.use_custom_key ? formData.value.custom_key : undefined
      await keysAPI.create(
        formData.value.name,
        formData.value.group_id,
        customKey,
        ipWhitelist,
        ipBlacklist,
        quota,
        expiresInDays,
        rateLimitData
      )
      appStore.showSuccess(t('keys.keyCreatedSuccess'))
      // Only advance tour if active, on submit step, and creation succeeded
      if (onboardingStore.isCurrentStep('[data-tour="key-form-submit"]')) {
        onboardingStore.nextStep(500)
      }
    }
    closeModals()
    loadApiKeys()
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToSave')
    appStore.showError(errorMsg)
    // Don't advance tour on error
  } finally {
    submitting.value = false
  }
}

/**
 * 删除 API Key 时优先展示后端返回的具体错误信息。
 */
const handleDelete = async () => {
  if (!selectedKey.value) return

  try {
    await keysAPI.delete(selectedKey.value.id)
    appStore.showSuccess(t('keys.keyDeletedSuccess'))
    showDeleteDialog.value = false
    loadApiKeys()
  } catch (error: any) {
    // Prefer backend error messages when available.
    const errorMsg = error?.message || t('keys.failedToDelete')
    appStore.showError(errorMsg)
  }
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  advancedSettingsExpanded.value = false
  selectedKey.value = null
  formData.value = {
    name: '',
    group_id: null,
    status: 'active',
    use_custom_key: false,
    custom_key: '',
    enable_ip_restriction: false,
    ip_whitelist: '',
    ip_blacklist: '',
    enable_quota: false,
    quota: null,
    enable_rate_limit: false,
    rate_limit_5h: null,
    rate_limit_1d: null,
    rate_limit_7d: null,
    enable_expiration: false,
    expiration_preset: '30',
    expiration_date: ''
  }
}

// Show reset quota confirmation dialog
const confirmResetQuota = () => {
  showResetQuotaDialog.value = true
}

// Set expiration date based on quick select days
const setExpirationDays = (days: number) => {
  formData.value.expiration_preset = days.toString() as '7' | '30' | '90'
  const expDate = new Date()
  expDate.setDate(expDate.getDate() + days)
  formData.value.expiration_date = formatDateTimeLocal(expDate.toISOString())
}

// Reset quota used for an API key
const resetQuotaUsed = async () => {
  if (!selectedKey.value) return
  showResetQuotaDialog.value = false
  try {
    await keysAPI.update(selectedKey.value.id, { reset_quota: true })
    appStore.showSuccess(t('keys.quotaResetSuccess'))
    // Update local state
    if (selectedKey.value) {
      selectedKey.value.quota_used = 0
    }
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToResetQuota')
    appStore.showError(errorMsg)
  }
}

// Show reset rate limit confirmation dialog (from edit modal)
const confirmResetRateLimit = () => {
  showResetRateLimitDialog.value = true
}

// Reset rate limit usage for an API key
const resetRateLimitUsage = async () => {
  if (!selectedKey.value) return
  showResetRateLimitDialog.value = false
  try {
    await keysAPI.update(selectedKey.value.id, { reset_rate_limit_usage: true })
    appStore.showSuccess(t('keys.rateLimitResetSuccess'))
    // Refresh key data
    await loadApiKeys()
    // Update the editing key with fresh data
    const refreshedKey = apiKeys.value.find(k => k.id === selectedKey.value!.id)
    if (refreshedKey) {
      selectedKey.value = refreshedKey
    }
  } catch (error: any) {
    const errorMsg = error.response?.data?.detail || t('keys.failedToResetRateLimit')
    appStore.showError(errorMsg)
  }
}

const importToCcswitch = (row: ApiKey) => {
  const platform = row.group?.platform || 'anthropic'

  // For antigravity platform, show client selection dialog
  if (platform === 'antigravity') {
    pendingCcsRow.value = row
    showCcsClientSelect.value = true
    return
  }

  // For other platforms, execute directly
  executeCcsImport(row, platform === 'gemini' ? 'gemini' : 'claude')
}

const executeCcsImport = (row: ApiKey, clientType: CcSwitchClientType) => {
  const baseUrl = publicSettings.value?.api_base_url || window.location.origin
  const platform = row.group?.platform || 'anthropic'

  const usageScript = `({
    request: {
      url: "{{baseUrl}}/v1/usage",
      method: "GET",
      headers: { "Authorization": "Bearer {{apiKey}}" }
    },
    extractor: function(response) {
      const remaining = response?.remaining ?? response?.quota?.remaining ?? response?.balance;
      const unit = response?.unit ?? response?.quota?.unit ?? "USD";
      return {
        isValid: response?.is_active ?? response?.isValid ?? true,
        remaining,
        unit
      };
    }
  })`
  const providerName = (publicSettings.value?.site_name || 'sub2api').trim() || 'sub2api'
  const deeplink = buildCcSwitchImportDeeplink({
    baseUrl,
    platform,
    clientType,
    providerName,
    apiKey: row.key,
    usageScript
  })

  try {
    window.open(deeplink, '_self')

    // Check if the protocol handler worked by detecting if we're still focused
    setTimeout(() => {
      if (document.hasFocus()) {
        // Still focused means the protocol handler likely failed
        appStore.showError(t('keys.ccSwitchNotInstalled'))
      }
    }, 100)
  } catch (error) {
    appStore.showError(t('keys.ccSwitchNotInstalled'))
  }
}

const handleCcsClientSelect = (clientType: CcSwitchClientType) => {
  if (pendingCcsRow.value) {
    executeCcsImport(pendingCcsRow.value, clientType)
  }
  showCcsClientSelect.value = false
  pendingCcsRow.value = null
}

const closeCcsClientSelect = () => {
  showCcsClientSelect.value = false
  pendingCcsRow.value = null
}

onMounted(() => {
  loadSavedColumns()
  loadApiKeys().then(() => {
    if (route.query.panel === 'connection-test') {
      nextTick(openConnectionTestDialog)
    }
  })
  loadGroups()
  loadUserGroupRates()
  loadPublicSettings()
  loadFilterKeyCatalog()
  document.addEventListener('click', handleDocumentClick)
  window.addEventListener('resize', handleGroupSelectorViewportChange)
  window.addEventListener('scroll', handleGroupSelectorViewportChange, true)
})

onUnmounted(() => {
  if (abortController) abortController.abort()
  if (workbenchAbortController) workbenchAbortController.abort()
  document.removeEventListener('click', handleDocumentClick)
  window.removeEventListener('resize', handleGroupSelectorViewportChange)
  window.removeEventListener('scroll', handleGroupSelectorViewportChange, true)
})
</script>

<style scoped>
.keys-page {
  color: #1f2320;
}

.keys-hero {
  position: relative;
  background-image:
    linear-gradient(90deg, rgba(250, 247, 239, 0.78), rgba(250, 247, 239, 0.42)),
    linear-gradient(rgba(31, 35, 32, 0.024) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.018) 1px, transparent 1px);
  background-size: auto, 86px 86px, 86px 86px;
}

.keys-hero::after {
  content: '';
  position: absolute;
  right: 1.75rem;
  top: 1.75rem;
  width: 0.72rem;
  height: 0.72rem;
  border-radius: 999px;
  background: #a73a2a;
  opacity: 0.82;
}

.keys-ledger-item {
  border-left: 1px solid rgba(216, 205, 185, 0.95);
  padding-left: 1rem;
}

.keys-ledger-item span {
  display: block;
  font-size: 0.74rem;
  color: #59645a;
}

.keys-ledger-item strong {
  display: block;
  margin-top: 0.35rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 1.35rem;
  color: #1f2320;
}

.workbench-panel-title {
  border-bottom: 1px solid rgba(216, 205, 185, 0.76);
  padding-bottom: 0.72rem;
}

.workbench-panel-title span {
  display: block;
  margin-bottom: 0.25rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  color: #59645a;
}

.workbench-panel-title strong {
  display: block;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.05rem;
  color: #1f2320;
}

.keys-rail-card {
  display: grid;
  gap: 0.9rem;
  border: 1px solid rgba(216, 205, 185, 0.62);
  border-radius: 16px;
  background: linear-gradient(180deg, rgba(250, 247, 239, 0.62), rgba(255, 252, 246, 0.38));
  padding: 0.9rem;
  box-shadow: 0 18px 42px -36px rgba(31, 35, 32, 0.24);
}

.keys-rail-card + .keys-rail-card {
  margin-top: 0.8rem;
}

.keys-rail-card-filters {
  gap: 0.95rem;
}

.keys-rail-intro {
  display: grid;
  gap: 0.3rem;
}

.keys-rail-intro span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.16em;
}

.keys-rail-intro p {
  color: #59645a;
  font-size: 0.82rem;
  line-height: 1.72;
}

.keys-rail-card-endpoint {
  padding: 0.75rem 0.9rem;
}

.keys-rail-card-actions {
  gap: 0;
}

.keys-filters :deep(.input),
.keys-filters :deep(.select-trigger) {
  min-height: 2.9rem;
  border-color: rgba(216, 205, 185, 0.76);
  border-radius: 10px;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.42), rgba(250, 247, 239, 0.7));
  color: #38413a;
  box-shadow: inset 0 1px 0 rgba(255, 252, 246, 0.72);
}

.keys-filters :deep(.input::placeholder) {
  color: #9aa391;
}

.keys-filters :deep(.select-trigger) {
  padding-inline: 0.95rem;
}

.keys-filters :deep(.select-value) {
  color: #1f2320;
  font-weight: 650;
}

.keys-filters :deep(.select-icon),
.keys-filters :deep(svg) {
  color: #8d978a;
}

.keys-filters :deep(.input:hover),
.keys-filters :deep(.select-trigger:hover),
.keys-filters :deep(.input:focus),
.keys-filters :deep(.select-trigger:focus-visible),
.keys-filters :deep(.select-trigger-open) {
  border-color: rgba(167, 58, 42, 0.34);
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(250, 247, 239, 0.84));
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.08), inset 0 1px 0 rgba(255, 252, 246, 0.78);
  outline: none;
}

.keys-toolbar-grid :deep(.select-trigger) {
  min-height: 2.8rem;
  border-color: rgba(216, 205, 185, 0.58);
  border-radius: 12px;
  background: rgba(255, 252, 246, 0.34);
  color: #38413a;
  box-shadow: none;
}

.keys-toolbar-grid :deep(.select-trigger:hover),
.keys-toolbar-grid :deep(.select-trigger:focus-visible),
.keys-toolbar-grid :deep(.select-trigger-open) {
  border-color: rgba(167, 58, 42, 0.24);
  background: rgba(255, 252, 246, 0.5);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.05);
}

.keys-toolbar-grid :deep(.select-value) {
  color: #1f2320;
  font-weight: 650;
}

.keys-toolbar-grid :deep(.select-icon),
.keys-toolbar-grid :deep(svg) {
  color: #8d978a;
}

:global(.keys-filter-dropdown) {
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 12px !important;
  background: rgba(250, 247, 239, 0.98) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34) !important;
}

:global(.keys-filter-dropdown .select-option) {
  color: #38413a !important;
}

:global(.keys-filter-dropdown .select-option:hover),
:global(.keys-filter-dropdown .select-option-focused) {
  background: rgba(167, 58, 42, 0.075) !important;
  color: #a73a2a !important;
}

:global(.keys-filter-dropdown .select-option-selected) {
  background: rgba(167, 58, 42, 0.1) !important;
  color: #a73a2a !important;
}

.keys-workspace {
  display: grid;
  gap: 1rem;
}

.keys-toolbar-shell,
.keys-data-shell,
.keys-endpoint-popover-shell {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.keys-toolbar-shell {
  display: grid;
  gap: 0.75rem;
  padding: 0;
}

.keys-toolbar-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
  align-items: end;
}

.keys-filter-field {
  display: grid;
  gap: 0;
  min-width: 10rem;
  flex: 1 1 11rem;
}

.keys-filter-field-key {
  min-width: 14rem;
  flex-basis: 16rem;
}

.keys-toolbar-actions {
  display: inline-flex;
  gap: 0.55rem;
  justify-content: flex-end;
  flex-wrap: wrap;
  margin-left: auto;
  min-width: fit-content;
}

.keys-toolbar-actions .btn {
  min-width: 8.4rem;
}

.keys-endpoint-popover-shell {
  padding: 0;
}

.keys-data-shell {
  padding: 0;
}

.keys-data-shell .keys-access-strip {
  margin-bottom: 0.7rem;
}

.keys-data-shell :deep(.page-size-select) {
  width: 5.5rem;
  min-width: 5.5rem;
}

.keys-data-shell :deep(.page-size-select .select-trigger) {
  min-height: 2.04rem;
  border-color: rgba(216, 205, 185, 0.5);
  border-radius: 999px;
  background: rgba(250, 247, 239, 0.24);
  color: #38413a;
  box-shadow: none;
}

.keys-data-shell :deep(.page-size-select .select-trigger:hover),
.keys-data-shell :deep(.page-size-select .select-trigger:focus-visible) {
  border-color: rgba(167, 58, 42, 0.24);
  background: rgba(250, 247, 239, 0.36);
  box-shadow: 0 0 0 2px rgba(167, 58, 42, 0.06);
}

.keys-data-shell :deep(nav) {
  border: 0;
  border-radius: 0;
  background: transparent;
  padding: 0;
  box-shadow: none;
}

.keys-access-strip {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: flex-start;
  gap: 0.75rem;
  margin: 0 0 1rem;
  padding: 0.18rem 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.keys-access-value {
  display: flex;
  align-items: baseline;
  flex: 0 1 auto;
  min-width: 0;
  gap: 0.56rem;
  padding: 0;
}

.keys-access-value > span {
  flex: 0 0 auto;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
  white-space: nowrap;
}

.keys-access-value > strong {
  display: block;
  flex: 0 1 auto;
  min-width: 0;
  overflow: hidden;
  max-width: min(42vw, 28rem);
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #38413a;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.88rem;
  font-weight: 600;
}

.keys-access-strip-actions {
  display: inline-flex;
  flex: 0 0 auto;
  flex-wrap: wrap;
  gap: 0.38rem;
  margin-left: 0.1rem;
}

.keys-access-strip-actions button {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  border: 1px solid rgba(216, 205, 185, 0.5);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.32);
  padding: 0.3rem 0.52rem;
  color: #59645a;
  font-size: 0.74rem;
  font-weight: 650;
  white-space: nowrap;
}

.keys-access-strip-actions button:disabled {
  opacity: 0.56;
}

.keys-access-strip-actions button:hover:not(:disabled),
.keys-access-strip-actions button:focus-visible {
  border-color: rgba(167, 58, 42, 0.28);
  background: rgba(167, 58, 42, 0.06);
  color: #a73a2a;
  outline: none;
}

.keys-page :deep(.table-wrapper) {
  overflow: auto;
  min-height: 0;
  max-height: min(34rem, calc(100vh - 19rem));
  flex: 0 0 auto !important;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  isolation: auto;
}

.keys-page :deep(table) {
  width: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  background: transparent;
  table-layout: auto;
}

.keys-page :deep(thead) {
  background: transparent;
}

.keys-page :deep(th) {
  border-bottom: 1px solid rgba(167, 58, 42, 0.2);
  background: rgba(237, 229, 212, 0.58) !important;
  color: #7b6a53;
  font-size: 0.65rem;
  font-weight: 650;
  letter-spacing: 0.12em;
  padding-top: 0.62rem;
  padding-bottom: 0.62rem;
}

.keys-page :deep(td) {
  border-bottom: 1px solid rgba(216, 205, 185, 0.46);
  background: rgba(255, 252, 246, 0.24);
  color: #38413a;
  padding-top: 0.68rem;
  padding-bottom: 0.68rem;
  vertical-align: middle;
  white-space: normal;
}

.keys-page :deep(tbody) {
  background: transparent;
}

.keys-page :deep(tbody .sticky-col) {
  background: rgba(250, 247, 239, 0.92) !important;
}

.keys-page :deep(tbody tr:hover) {
  background: transparent;
}

.keys-page :deep(tbody tr:hover td),
.keys-page :deep(tbody tr:hover .sticky-col) {
  background: rgba(167, 58, 42, 0.04) !important;
}

.keys-page :deep(.code) {
  border: 0;
  border-bottom: 1px solid rgba(216, 205, 185, 0.88);
  border-radius: 0;
  background: transparent;
  padding-inline: 0;
  color: #38413a;
  max-width: 8.2rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.endpoint-cell {
  display: inline-flex;
  max-width: 18rem;
  align-items: center;
  gap: 0.42rem;
}

.endpoint-code {
  overflow: hidden;
  max-width: 14rem;
  text-overflow: ellipsis;
  white-space: nowrap;
  border-bottom: 1px solid rgba(167, 58, 42, 0.28);
  color: #59645a;
  font-size: 0.78rem;
}

.endpoint-copy {
  display: inline-grid;
  width: 1.75rem;
  height: 1.75rem;
  place-items: center;
  border-radius: 999px;
  color: #8b6f5b;
  transition: background-color 160ms ease, color 160ms ease;
}

.endpoint-copy:hover,
.endpoint-copy:focus-visible {
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
  outline: none;
}

.key-name-cell {
  display: flex;
  align-items: center;
  gap: 0.34rem;
  min-width: 0;
}

.key-name-text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.key-name-flag {
  color: #8b7d67;
}

.key-group-trigger {
  display: inline-flex;
  align-items: center;
  gap: 0.38rem;
  min-height: 1.9rem;
  min-width: 0;
  max-width: min(100%, 13.5rem);
  overflow: hidden;
  padding: 0.24rem 0.48rem;
  border: 1px solid rgba(216, 205, 185, 0.56);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.46);
  color: #59645a;
  white-space: nowrap;
  transition: border-color 160ms ease, background 160ms ease, color 160ms ease, transform 160ms ease;
}

.key-group-trigger:hover,
.key-group-trigger:focus-visible {
  border-color: rgba(167, 58, 42, 0.28);
  background: rgba(167, 58, 42, 0.06);
  color: #a73a2a;
  outline: none;
  transform: translateY(-1px);
}

.key-group-empty {
  color: #7b6a53;
  font-size: 0.76rem;
}

.key-group-helper {
  color: #8b7d67;
  display: none;
  font-size: 0.66rem;
  white-space: nowrap;
}

.key-group-chevron {
  width: 0.72rem;
  height: 0.72rem;
  color: #8b7d67;
  opacity: 0.76;
  transition: opacity 160ms ease, color 160ms ease;
}

.key-group-trigger:hover .key-group-chevron,
.key-group-trigger:focus-visible .key-group-chevron {
  color: #a73a2a;
  opacity: 1;
}

.key-usage-cell,
.key-expiry-cell {
  display: grid;
  gap: 0.38rem;
  justify-items: start;
  min-width: 0;
}

.key-usage-ledger {
  display: inline-grid;
  grid-template-columns: repeat(2, minmax(5rem, max-content));
  justify-content: start;
  gap: 0.42rem 1.25rem;
}

.key-usage-stat,
.key-meter-card,
.key-expiry-cell {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.key-usage-stat {
  padding: 0;
}

.key-usage-stat span,
.key-meter-label,
.key-expiry-label {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.key-usage-stat strong,
.key-expiry-cell strong {
  display: block;
  margin-top: 0.12rem;
  color: #1f2320;
  font-size: 0.8rem;
  font-weight: 650;
  line-height: 1.35;
  font-variant-numeric: tabular-nums;
}

.key-meter-card {
  padding: 0;
}

.key-meter-card-compact {
  padding-block: 0;
}

.key-meter-head,
.key-meter-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.key-meter-value {
  color: #38413a;
  font-size: 0.74rem;
  font-weight: 650;
  line-height: 1.35;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.key-meter-track {
  height: 0.22rem;
  margin-top: 0.18rem;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(216, 205, 185, 0.62);
}

.key-meter-fill {
  height: 100%;
  border-radius: inherit;
  background: #51624f;
  transition: width 180ms ease, background-color 180ms ease;
}

.key-meter-foot {
  margin-top: 0.2rem;
  color: #7b6a53;
  font-size: 0.64rem;
  line-height: 1.4;
}

.key-usage-meta {
  display: grid;
  gap: 0.28rem;
}

.key-rate-inline-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.22rem 0.52rem;
}

.key-rate-inline-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.34rem;
  border: 0;
  border-radius: 0;
  background: transparent;
  padding: 0;
  color: #667066;
  font-size: 0.64rem;
  line-height: 1.2;
  white-space: nowrap;
}

.key-rate-inline-chip + .key-rate-inline-chip::before {
  content: '·';
  color: rgba(123, 106, 83, 0.62);
  margin-right: 0.18rem;
}

.key-rate-inline-chip strong {
  color: #38413a;
  font-size: 0.62rem;
  letter-spacing: 0.08em;
}

.key-rate-inline-chip.tone-warning {
  color: #8f7853;
}

.key-rate-inline-chip.tone-danger {
  color: #a73a2a;
}

.key-meter-card.tone-safe .key-meter-fill,
.key-expiry-cell.tone-safe::before {
  background: #51624f;
}

.key-meter-card.tone-warning {
  color: #8f7853;
}

.key-meter-card.tone-warning .key-meter-fill,
.key-expiry-cell.tone-warning::before {
  background: #9b8155;
}

.key-meter-card.tone-danger {
  color: #a73a2a;
}

.key-meter-card.tone-danger .key-meter-fill,
.key-expiry-cell.tone-danger::before {
  background: #a73a2a;
}

.key-inline-note {
  color: #8b7d67;
  font-size: 0.72rem;
  line-height: 1.55;
}

.key-meter-reset {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: 0.3rem;
  border: 1px solid rgba(216, 205, 185, 0.68);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.72);
  padding: 0.34rem 0.58rem;
  color: #59645a;
  font-size: 0.7rem;
  font-weight: 650;
  transition: border-color 160ms ease, background-color 160ms ease, color 160ms ease;
}

.key-meter-reset:hover,
.key-meter-reset:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.07);
  color: #a73a2a;
  outline: none;
}

.key-expiry-cell {
  position: relative;
  padding: 0.18rem 0 0 0.66rem;
}

.key-expiry-cell::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0.12rem;
  bottom: 0.12rem;
  width: 2px;
  border-radius: 999px;
  background: rgba(81, 98, 79, 0.72);
}

.key-expiry-cell small {
  display: block;
  margin-top: 0.28rem;
  color: #667066;
  font-size: 0.68rem;
  line-height: 1.45;
}

.connection-test-modal :deep(.modal-content) {
  overflow: hidden;
  border-color: rgba(198, 184, 157, 0.58);
  border-radius: 18px;
  background:
    radial-gradient(circle at 16% 0%, rgba(167, 58, 42, 0.055), transparent 18rem),
    linear-gradient(180deg, rgba(250, 247, 239, 0.98), rgba(246, 241, 230, 0.96));
  box-shadow: 0 34px 90px -58px rgba(31, 35, 32, 0.62);
}

.connection-test-modal :deep(.modal-header) {
  border-bottom-color: rgba(216, 205, 185, 0.64);
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.84), rgba(250, 247, 239, 0.68));
}

.connection-test-modal :deep(.modal-title) {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.95rem;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.connection-test-modal :deep(.modal-body) {
  padding-top: 1.25rem;
  padding-bottom: 1.15rem;
}

.connection-test-modal :deep(.modal-footer) {
  border-top-color: rgba(216, 205, 185, 0.64);
  background: linear-gradient(180deg, rgba(250, 247, 239, 0.58), rgba(246, 241, 230, 0.8));
}

.connection-test-dialog {
  display: grid;
  gap: 1.2rem;
}

.connection-test-intro {
  display: grid;
  gap: 0.36rem;
  padding-bottom: 0.2rem;
}

.connection-test-intro span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.connection-test-intro p {
  color: #59645a;
  font-size: 0.96rem;
  line-height: 1.75;
}

.connection-test-dialog code {
  color: #38413a;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.connection-key-list {
  display: grid;
  gap: 0.58rem;
  max-height: 14rem;
  overflow: auto;
  padding-right: 0.16rem;
}

.connection-key-list button {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  align-items: center;
  gap: 1rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 16px;
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.04), transparent 36%), rgba(255, 252, 246, 0.72);
  padding: 0.9rem 1rem;
  color: #38413a;
  text-align: left;
  transition: border-color 160ms ease, background 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.connection-key-list button span {
  min-width: 0;
  color: #1f2320;
  font-size: 1rem;
  font-weight: 650;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.connection-key-list button.is-selected,
.connection-key-list button:hover,
.connection-key-list button:focus-visible {
  border-color: rgba(167, 58, 42, 0.4);
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 42%), rgba(250, 247, 239, 0.86);
  box-shadow: inset 3px 0 0 rgba(167, 58, 42, 0.72), 0 18px 42px -36px rgba(31, 35, 32, 0.32);
  outline: none;
  transform: translateY(-1px);
}

.connection-key-list small {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.82rem;
}

.integration-kit {
  display: grid;
  gap: 0.82rem;
  border-top: 1px solid rgba(216, 205, 185, 0.72);
  padding-top: 0.95rem;
}

.connection-model-brief {
  display: grid;
  gap: 0.7rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 16px;
  background: rgba(250, 247, 239, 0.68);
  padding: 0.92rem 1rem;
}

.connection-model-brief-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.connection-model-brief-head span,
.connection-model-samples > span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.connection-model-brief-head strong {
  display: block;
  margin-top: 0.24rem;
  color: #38413a;
  font-size: 0.86rem;
  line-height: 1.6;
}

.connection-model-brief-head button {
  flex: 0 0 auto;
  border: 1px solid rgba(216, 205, 185, 0.82);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.76);
  padding: 0.38rem 0.68rem;
  color: #59645a;
  font-size: 0.74rem;
  font-weight: 700;
}

.connection-model-brief-list {
  display: grid;
  gap: 0.36rem;
  padding-left: 1rem;
  color: #59645a;
  font-size: 0.8rem;
  line-height: 1.6;
}

.integration-kit-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 1rem;
  min-width: 0;
}

.integration-kit-head span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.16em;
}

.integration-kit-head strong {
  overflow: hidden;
  color: #38413a;
  font-size: 0.92rem;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.integration-kit-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 0.58rem;
}

.integration-kit-actions button {
  border: 1px solid rgba(216, 205, 185, 0.82);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.76);
  padding: 0.46rem 0.76rem;
  color: #59645a;
  font-size: 0.78rem;
  font-weight: 700;
  transition: border-color 160ms ease, background 160ms ease, color 160ms ease, transform 160ms ease;
}

.integration-kit-actions button:hover,
.integration-kit-actions button:focus-visible {
  border-color: rgba(167, 58, 42, 0.38);
  background: rgba(167, 58, 42, 0.075);
  color: #a73a2a;
  outline: none;
  transform: translateY(-1px);
}

.connection-test-report {
  display: grid;
  gap: 0.78rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-left: 3px solid rgba(155, 129, 85, 0.68);
  border-radius: 16px;
  background: rgba(255, 252, 246, 0.72);
  padding: 1rem 1rem 1.05rem;
}

.connection-report-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.connection-report-head span,
.connection-report-grid span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.connection-report-head strong {
  display: block;
  margin-top: 0.24rem;
  color: #38413a;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
}

.connection-report-head button {
  flex: 0 0 auto;
  border: 1px solid rgba(216, 205, 185, 0.82);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.76);
  padding: 0.38rem 0.68rem;
  color: #59645a;
  font-size: 0.74rem;
  font-weight: 700;
}

.connection-test-report p,
.connection-test-report small {
  color: #59645a;
  font-size: 0.8rem;
  line-height: 1.65;
}

.connection-report-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.62rem;
}

.connection-report-grid div {
  min-width: 0;
  border-top: 1px solid rgba(216, 205, 185, 0.52);
  padding-top: 0.52rem;
}

.connection-report-grid strong {
  display: block;
  overflow: hidden;
  margin-top: 0.18rem;
  color: #38413a;
  font-size: 0.82rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.connection-model-samples {
  display: grid;
  gap: 0.45rem;
  border-top: 1px solid rgba(216, 205, 185, 0.52);
  padding-top: 0.52rem;
}

.connection-model-samples-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.8rem;
}

.connection-model-samples-head span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.connection-model-samples-head strong {
  color: #38413a;
  font-size: 0.8rem;
  font-weight: 600;
}

.connection-model-sample-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.42rem;
}

.connection-model-sample-list code {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.8);
  padding: 0.22rem 0.56rem;
  font-size: 0.74rem;
}

.connection-model-sample-more {
  display: inline-flex;
  align-items: center;
  border: 1px dashed rgba(216, 205, 185, 0.82);
  border-radius: 999px;
  background: rgba(250, 247, 239, 0.72);
  padding: 0.22rem 0.56rem;
  color: #7b6a53;
  font-size: 0.74rem;
  cursor: pointer;
  transition: border-color 160ms ease, background-color 160ms ease, color 160ms ease;
}

.connection-model-sample-more:hover,
.connection-model-sample-more:focus-visible {
  border-color: rgba(167, 58, 42, 0.38);
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
  outline: none;
}

.connection-test-report.is-success {
  border-left-color: #51624f;
  background: rgba(81, 98, 79, 0.055);
}

.connection-test-report.is-warning {
  border-left-color: #9b8155;
  background: rgba(155, 129, 85, 0.065);
}

.connection-test-report.is-danger {
  border-left-color: #a73a2a;
  background: rgba(167, 58, 42, 0.055);
}

.key-editor-modal :deep(.modal-content) {
  overflow: hidden;
  border-color: rgba(198, 184, 157, 0.58);
  border-radius: 18px;
  background:
    radial-gradient(circle at 16% 0%, rgba(167, 58, 42, 0.05), transparent 18rem),
    linear-gradient(180deg, rgba(250, 247, 239, 0.98), rgba(246, 241, 230, 0.96));
  box-shadow: 0 34px 90px -58px rgba(31, 35, 32, 0.62);
}

.key-editor-modal :deep(.modal-header) {
  border-bottom-color: rgba(216, 205, 185, 0.64);
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.84), rgba(250, 247, 239, 0.68));
}

.key-editor-modal :deep(.modal-title) {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.85rem;
  font-weight: 600;
}

.key-editor-modal :deep(.modal-body) {
  padding-top: 1.2rem;
  padding-bottom: 1rem;
}

.key-editor-modal :deep(.modal-footer) {
  border-top-color: rgba(216, 205, 185, 0.64);
  background: linear-gradient(180deg, rgba(250, 247, 239, 0.58), rgba(246, 241, 230, 0.8));
}

.key-editor-form {
  display: grid;
  gap: 1.1rem;
}

.key-editor-form .input-label {
  color: #59645a;
  font-size: 0.8rem;
  font-weight: 650;
}

.key-editor-form :deep(.input),
.key-editor-form :deep(.select-trigger),
.key-editor-form textarea,
.key-editor-form input[type='datetime-local'] {
  border-color: rgba(216, 205, 185, 0.76);
  border-radius: 10px;
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.42), rgba(250, 247, 239, 0.7));
  color: #38413a;
  box-shadow: inset 0 1px 0 rgba(255, 252, 246, 0.72);
}

.key-editor-form :deep(.input:focus),
.key-editor-form :deep(.select-trigger:focus-visible),
.key-editor-form :deep(.select-trigger-open),
.key-editor-form textarea:focus,
.key-editor-form input[type='datetime-local']:focus {
  border-color: rgba(167, 58, 42, 0.34);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.08), inset 0 1px 0 rgba(255, 252, 246, 0.78);
  outline: none;
}

.key-editor-form :deep(.select-value),
.key-editor-form :deep(.select-icon),
.key-editor-form :deep(.select-clear),
.key-editor-form input::placeholder,
.key-editor-form textarea::placeholder {
  color: #8d978a;
}

.key-editor-form .input-hint {
  color: #7b6a53;
  font-size: 0.78rem;
  line-height: 1.65;
}

.key-advanced-shell {
  display: grid;
  gap: 0.82rem;
}

.key-advanced-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  width: 100%;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 16px;
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 42%), rgba(255, 252, 246, 0.68);
  padding: 0.95rem 1rem;
  text-align: left;
  transition: border-color 160ms ease, background 160ms ease, box-shadow 160ms ease, transform 160ms ease;
}

.key-advanced-toggle:hover,
.key-advanced-toggle:focus-visible {
  border-color: rgba(167, 58, 42, 0.36);
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 46%), rgba(250, 247, 239, 0.82);
  box-shadow: 0 18px 42px -36px rgba(31, 35, 32, 0.26);
  outline: none;
  transform: translateY(-1px);
}

.key-advanced-copy {
  display: grid;
  gap: 0.18rem;
}

.key-advanced-copy span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.18em;
}

.key-advanced-copy strong {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 600;
}

.key-advanced-meta {
  display: flex;
  align-items: center;
  gap: 0.58rem;
  color: #59645a;
}

.key-advanced-meta small {
  color: #8b7d67;
  font-size: 0.74rem;
  font-weight: 650;
}

.key-advanced-chips {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.35rem;
}

.key-advanced-chips span {
  border: 1px solid rgba(216, 205, 185, 0.68);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.7);
  padding: 0.22rem 0.56rem;
  color: #59645a;
  font-size: 0.72rem;
  font-weight: 650;
}

.key-advanced-panel {
  border: 1px solid rgba(216, 205, 185, 0.62);
  border-radius: 16px;
  background: rgba(255, 252, 246, 0.42);
  padding: 1rem;
}

:global(.key-editor-dropdown) {
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 12px !important;
  background: rgba(250, 247, 239, 0.98) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34) !important;
}

:global(.key-editor-dropdown .select-option) {
  color: #38413a !important;
}

:global(.key-editor-dropdown .select-option:hover),
:global(.key-editor-dropdown .select-option-focused) {
  background: rgba(167, 58, 42, 0.075) !important;
  color: #a73a2a !important;
}

:global(.key-editor-dropdown .select-option-selected) {
  background: rgba(167, 58, 42, 0.1) !important;
  color: #a73a2a !important;
}

.dark .key-editor-modal :deep(.modal-content) {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    radial-gradient(circle at 16% 0%, rgba(167, 58, 42, 0.09), transparent 18rem),
    linear-gradient(180deg, rgba(24, 26, 21, 0.96), rgba(17, 19, 15, 0.94));
  box-shadow: 0 42px 100px -60px rgba(0, 0, 0, 0.72);
}

.dark .key-editor-modal :deep(.modal-header) {
  border-bottom-color: rgba(48, 52, 43, 0.78);
  background: linear-gradient(180deg, rgba(24, 26, 21, 0.88), rgba(17, 19, 15, 0.72));
}

.dark .key-editor-modal :deep(.modal-title) {
  color: #f4efe4;
}

.dark .key-editor-modal :deep(.modal-footer) {
  border-top-color: rgba(48, 52, 43, 0.78);
  background: linear-gradient(180deg, rgba(17, 19, 15, 0.62), rgba(24, 26, 21, 0.86));
}

.dark .key-editor-form .input-label,
.dark .key-editor-form .input-hint,
.dark .key-advanced-copy span,
.dark .key-advanced-meta,
.dark .key-advanced-meta small,
.dark .key-advanced-chips span {
  color: #879186;
}

.dark .key-editor-form :deep(.input),
.dark .key-editor-form :deep(.select-trigger),
.dark .key-editor-form textarea,
.dark .key-editor-form input[type='datetime-local'] {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
  box-shadow: none;
}

.dark .key-editor-form :deep(.input:focus),
.dark .key-editor-form :deep(.select-trigger:focus-visible),
.dark .key-editor-form :deep(.select-trigger-open),
.dark .key-editor-form textarea:focus,
.dark .key-editor-form input[type='datetime-local']:focus {
  border-color: rgba(216, 205, 185, 0.28);
  box-shadow: 0 0 0 3px rgba(184, 156, 116, 0.08);
}

.dark .key-editor-form :deep(.select-value),
.dark .key-editor-form :deep(.select-icon),
.dark .key-editor-form :deep(.select-clear),
.dark .key-editor-form input::placeholder,
.dark .key-editor-form textarea::placeholder {
  color: #879186;
}

.dark .key-advanced-toggle {
  border-color: rgba(48, 52, 43, 0.9);
  background: linear-gradient(90deg, rgba(184, 156, 116, 0.08), transparent 44%), rgba(24, 26, 21, 0.72);
}

.dark .key-advanced-toggle:hover,
.dark .key-advanced-toggle:focus-visible {
  border-color: rgba(184, 156, 116, 0.42);
  background: linear-gradient(90deg, rgba(184, 156, 116, 0.12), transparent 46%), rgba(24, 26, 21, 0.9);
  box-shadow: 0 18px 42px -36px rgba(0, 0, 0, 0.48);
}

.dark .key-advanced-copy strong {
  color: #f4efe4;
}

.dark .key-advanced-chips span {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.52);
}

.dark .key-advanced-panel {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.42);
}

:global(.dark .key-editor-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.98) !important;
}

:global(.dark .key-editor-dropdown .select-option) {
  color: #d8cdb9 !important;
}

:global(.dark .key-editor-dropdown .select-option:hover),
:global(.dark .key-editor-dropdown .select-option-focused) {
  background: rgba(216, 205, 185, 0.06) !important;
  color: #f4efe4 !important;
}

:global(.dark .key-editor-dropdown .select-option-selected) {
  background: rgba(184, 156, 116, 0.1) !important;
  color: #f4efe4 !important;
}

.dark .connection-test-modal :deep(.modal-content) {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    radial-gradient(circle at 16% 0%, rgba(167, 58, 42, 0.09), transparent 18rem),
    linear-gradient(180deg, rgba(24, 26, 21, 0.96), rgba(17, 19, 15, 0.94));
  box-shadow: 0 42px 100px -60px rgba(0, 0, 0, 0.72);
}

.dark .connection-test-modal :deep(.modal-header) {
  border-bottom-color: rgba(48, 52, 43, 0.78);
  background: linear-gradient(180deg, rgba(24, 26, 21, 0.88), rgba(17, 19, 15, 0.72));
}

.dark .connection-test-modal :deep(.modal-title) {
  color: #f4efe4;
}

.dark .connection-test-modal :deep(.modal-footer) {
  border-top-color: rgba(48, 52, 43, 0.78);
  background: linear-gradient(180deg, rgba(17, 19, 15, 0.62), rgba(24, 26, 21, 0.86));
}

.dark .connection-test-intro p {
  color: #879186;
}

.dark .connection-key-list button {
  border-color: rgba(48, 52, 43, 0.9);
  background: linear-gradient(90deg, rgba(184, 156, 116, 0.06), transparent 38%), rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
}

.dark .connection-key-list button span {
  color: #f4efe4;
}

.dark .connection-key-list button.is-selected,
.dark .connection-key-list button:hover,
.dark .connection-key-list button:focus-visible {
  border-color: rgba(184, 156, 116, 0.42);
  background: linear-gradient(90deg, rgba(184, 156, 116, 0.12), transparent 44%), rgba(24, 26, 21, 0.9);
  box-shadow: inset 3px 0 0 rgba(167, 58, 42, 0.82), 0 18px 42px -36px rgba(0, 0, 0, 0.5);
}

.dark .connection-key-list small,
.dark .integration-kit-head span,
.dark .connection-model-brief-head span,
.dark .connection-report-head span,
.dark .connection-report-grid span,
.dark .connection-model-samples > span {
  color: #b9aa91;
}

.dark .integration-kit {
  border-top-color: rgba(48, 52, 43, 0.72);
}

.dark .connection-model-brief {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.68);
}

.dark .integration-kit-head strong,
.dark .connection-model-brief-head strong,
.dark .connection-report-head strong,
.dark .connection-report-grid strong {
  color: #f4efe4;
}

.dark .integration-kit-actions button,
.dark .connection-model-brief-head button,
.dark .connection-report-head button {
  border-color: rgba(48, 52, 43, 0.92);
  background: rgba(17, 19, 15, 0.52);
  color: #d8cdb9;
}

.dark .integration-kit-actions button:hover,
.dark .integration-kit-actions button:focus-visible,
.dark .connection-model-brief-head button:hover,
.dark .connection-model-brief-head button:focus-visible,
.dark .connection-report-head button:hover,
.dark .connection-report-head button:focus-visible {
  border-color: rgba(184, 156, 116, 0.38);
  background: rgba(167, 58, 42, 0.14);
  color: #f4efe4;
}

.dark .connection-test-report {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.72);
}

.dark .connection-test-report p,
.dark .connection-test-report small,
.dark .connection-model-brief-list {
  color: #bdb5a8;
}

.dark .connection-report-grid div {
  border-top-color: rgba(48, 52, 43, 0.64);
}

.dark .connection-model-samples {
  border-top-color: rgba(48, 52, 43, 0.64);
}

.dark .connection-model-samples-head span {
  color: #b9aa91;
}

.dark .connection-model-samples-head strong {
  color: #f4efe4;
}

.dark .connection-model-sample-list code {
  border-color: rgba(48, 52, 43, 0.92);
  background: rgba(17, 19, 15, 0.56);
  color: #f4efe4;
}

.dark .connection-model-sample-more {
  border-color: rgba(48, 52, 43, 0.92);
  background: rgba(24, 26, 21, 0.58);
  color: #b9aa91;
}

.dark .connection-model-sample-more:hover,
.dark .connection-model-sample-more:focus-visible {
  border-color: rgba(184, 156, 116, 0.38);
  background: rgba(167, 58, 42, 0.14);
  color: #f4efe4;
}

.dark .connection-test-report.is-success {
  border-left-color: #6e8a67;
  background: rgba(81, 98, 79, 0.12);
}

.dark .connection-test-report.is-warning {
  border-left-color: #b89c74;
  background: rgba(155, 129, 85, 0.14);
}

.dark .connection-test-report.is-danger {
  border-left-color: #c06c5d;
  background: rgba(167, 58, 42, 0.14);
}

@media (max-width: 640px) {
  .connection-test-modal :deep(.modal-title) {
    font-size: 1.55rem;
  }

  .integration-kit-head {
    align-items: flex-start;
    flex-direction: column;
    gap: 0.35rem;
  }

  .connection-model-brief-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .connection-model-samples-head {
    align-items: flex-start;
    flex-direction: column;
    gap: 0.28rem;
  }

  .integration-kit-head strong {
    max-width: 100%;
    text-align: left;
  }

  .connection-key-list button {
    grid-template-columns: 1fr;
    gap: 0.38rem;
  }

  .connection-report-head {
    align-items: flex-start;
    flex-direction: column;
  }

  .connection-report-grid {
    grid-template-columns: 1fr;
  }
}

.key-row-actions {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  white-space: nowrap;
}

.key-health-cell {
  display: flex;
  align-items: flex-start;
  gap: 0.45rem;
  min-width: 0;
  padding: 0.12rem 0.12rem 0.12rem 0;
}

.key-health-cell .badge {
  width: fit-content;
  flex: 0 0 auto;
}

.key-health-lines {
  display: grid;
  gap: 0.14rem;
  padding-top: 0.04rem;
  color: #667066;
  font-size: 0.64rem;
  line-height: 1.42;
}

.key-health-lines span {
  overflow: hidden;
  max-width: 9rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.key-health-lines-alert {
  color: #8f3b2f;
}

.key-health-model-hint {
  max-width: 11rem;
  color: #8f3b2f;
  overflow: visible;
  text-overflow: initial;
  white-space: normal;
}

.key-health-link {
  width: fit-content;
  border: 0;
  padding: 0;
  background: transparent;
  color: #8f3b2f;
  font-size: 0.64rem;
  line-height: 1.42;
  text-decoration: underline;
  text-underline-offset: 0.18rem;
}

.key-health-link:hover,
.key-health-link:focus-visible {
  color: #a73a2a;
}

.key-row-action {
  display: inline-grid;
  width: 1.9rem;
  height: 1.9rem;
  place-items: center;
  border: 1px solid rgba(216, 205, 185, 0.48);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.34);
  color: #667066;
  transition: background-color 160ms ease, border-color 160ms ease, color 160ms ease, transform 160ms ease;
}

.key-row-action:hover,
.key-row-action:focus-visible {
  border-color: rgba(167, 58, 42, 0.22);
  background: rgba(167, 58, 42, 0.07);
  color: #a73a2a;
  outline: none;
  transform: translateY(-1px);
}

.key-row-action.is-danger:hover,
.key-row-action.is-danger:focus-visible {
  border-color: rgba(167, 58, 42, 0.26);
  background: rgba(167, 58, 42, 0.1);
}

.keys-page :deep(.rounded-lg),
.keys-page :deep(.rounded-xl),
.keys-page :deep(.rounded-2xl) {
  border-radius: 6px;
}

.keys-page :deep(.bg-blue-50),
.keys-page :deep(.bg-blue-100),
.keys-page :deep(.bg-green-50),
.keys-page :deep(.bg-green-100),
.keys-page :deep(.bg-yellow-50),
.keys-page :deep(.bg-amber-100),
.keys-page :deep(.bg-primary-50),
.keys-page :deep(.bg-primary-100),
.keys-page :deep(.bg-gray-50) {
  background-color: rgba(216, 205, 185, 0.46);
}

.keys-page :deep(.text-blue-500),
.keys-page :deep(.text-blue-600),
.keys-page :deep(.text-green-500),
.keys-page :deep(.text-green-600),
.keys-page :deep(.text-yellow-500),
.keys-page :deep(.text-amber-600),
.keys-page :deep(.text-primary-500),
.keys-page :deep(.text-primary-600) {
  color: #59645a;
}

.keys-page :deep(button.text-primary-600),
.keys-page :deep(a.text-primary-600),
.keys-page :deep(.btn-primary) {
  color: inherit;
}

.keys-page :deep(.btn-primary) {
  background-color: #a73a2a;
  border-color: #a73a2a;
  color: #f4efe4;
}

.keys-page :deep(.btn-primary:hover) {
  background-color: #8f3024;
  border-color: #8f3024;
}

.keys-page :deep(.btn-secondary) {
  border-color: rgba(216, 205, 185, 0.9);
  background: rgba(255, 255, 255, 0.38);
  color: #38413a;
}

.keys-page :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(255, 255, 255, 0.66);
}

.dark .keys-page {
  color: #f4efe4;
}

.dark .keys-ledger-item {
  border-left-color: rgba(48, 52, 43, 0.95);
}

.dark .keys-ledger-item span {
  color: #879186;
}

.dark .keys-ledger-item strong {
  color: #f4efe4;
}

.dark .keys-page :deep(thead) {
  background: transparent;
}

.dark .workbench-panel-title {
  border-bottom-color: rgba(48, 52, 43, 0.95);
}

.dark .workbench-panel-title span {
  color: #879186;
}

.dark .workbench-panel-title strong {
  color: #f4efe4;
}

.dark .keys-rail-card {
  border-color: rgba(48, 52, 43, 0.9);
  background: linear-gradient(180deg, rgba(24, 26, 21, 0.72), rgba(17, 19, 15, 0.46));
  box-shadow: 0 18px 42px -36px rgba(0, 0, 0, 0.42);
}

.dark .keys-rail-intro span,
.dark .keys-rail-intro p {
  color: #879186;
}

.dark .keys-filters :deep(.input),
.dark .keys-filters :deep(.select-trigger) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
}

.dark .keys-filters :deep(.select-value) {
  color: #f4efe4;
}

.dark .keys-filters :deep(.input::placeholder),
.dark .keys-filters :deep(.select-icon),
.dark .keys-filters :deep(svg) {
  color: #879186;
}

:global(.dark .keys-filter-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.98) !important;
}

:global(.dark .keys-filter-dropdown .select-option) {
  color: #d8cdb9 !important;
}

.dark .keys-access-strip {
  border-color: transparent;
  background: transparent;
}

.dark .keys-access-value > span,
.dark .keys-access-value > strong {
  color: #d8cdb9;
}

.dark .keys-access-strip-actions button {
  border-color: rgba(48, 52, 43, 0.7);
  background: rgba(17, 19, 15, 0.24);
  color: #d8cdb9;
}

.dark .keys-access-strip-actions button:hover:not(:disabled),
.dark .keys-access-strip-actions button:focus-visible {
  border-color: rgba(167, 58, 42, 0.4);
  background: rgba(167, 58, 42, 0.14);
  color: #f0b4a8;
}

.dark .keys-page :deep(.table-wrapper) {
  border-color: transparent;
  background: transparent;
}

.dark .keys-page :deep(th) {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.9) !important;
  color: #b9aa91;
}

.dark .keys-page :deep(td) {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.34);
  color: #d8cdb9;
}

.dark .keys-page :deep(tbody .sticky-col) {
  background: rgba(17, 19, 15, 0.92) !important;
}

.dark .keys-page :deep(tbody tr:hover td),
.dark .keys-page :deep(tbody tr:hover .sticky-col) {
  background: rgba(167, 58, 42, 0.11) !important;
}

.dark .keys-toolbar-shell,
.dark .keys-data-shell,
.dark .keys-endpoint-popover-shell {
  border-color: transparent;
  background: transparent;
  box-shadow: none;
}

.dark .keys-toolbar-grid :deep(.select-trigger),
.dark .keys-data-shell :deep(.page-size-select .select-trigger) {
  border-color: rgba(48, 52, 43, 0.72);
  background: rgba(17, 19, 15, 0.26);
  color: #d8cdb9;
}

.dark .keys-toolbar-grid :deep(.select-trigger:hover),
.dark .keys-toolbar-grid :deep(.select-trigger:focus-visible),
.dark .keys-toolbar-grid :deep(.select-trigger-open),
.dark .keys-data-shell :deep(.page-size-select .select-trigger:hover),
.dark .keys-data-shell :deep(.page-size-select .select-trigger:focus-visible) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  box-shadow: 0 0 0 2px rgba(167, 58, 42, 0.08);
}

.dark .keys-toolbar-grid :deep(.select-value) {
  color: #f4efe4;
}

.dark .keys-toolbar-grid :deep(.select-icon),
.dark .keys-toolbar-grid :deep(svg) {
  color: #879186;
}

.dark .keys-data-shell :deep(nav) {
  border-color: transparent;
  background: transparent;
}

:global(.keys-page-size-dropdown) {
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 12px !important;
  background: rgba(250, 247, 239, 0.96) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34) !important;
}

:global(.keys-page-size-dropdown .select-option) {
  color: #38413a !important;
}

:global(.keys-page-size-dropdown .select-option:hover),
:global(.keys-page-size-dropdown .select-option-focused) {
  background: rgba(167, 58, 42, 0.075) !important;
  color: #a73a2a !important;
}

:global(.keys-page-size-dropdown .select-option-selected) {
  background: rgba(167, 58, 42, 0.1) !important;
  color: #a73a2a !important;
}

:global(.dark .keys-page-size-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.98) !important;
}

.dark .keys-page :deep(.code) {
  border-color: rgba(48, 52, 43, 0.95);
  background: transparent;
  color: #d8cdb9;
}

.dark .key-row-action {
  border-color: rgba(48, 52, 43, 0.76);
  background: rgba(17, 19, 15, 0.38);
  color: #879186;
}

.dark .key-usage-stat,
.dark .key-meter-card,
.dark .key-expiry-cell {
  border-color: transparent;
  background: transparent;
  box-shadow: none;
}

.dark .key-usage-stat span,
.dark .key-meter-label,
.dark .key-expiry-label,
.dark .key-meter-foot,
.dark .key-inline-note {
  color: #879186;
}

.dark .key-rate-inline-chip {
  color: #a4ada2;
}

.dark .key-rate-inline-chip + .key-rate-inline-chip::before {
  color: rgba(135, 145, 134, 0.7);
}

.dark .key-rate-inline-chip strong {
  color: #f4efe4;
}

.dark .key-usage-stat strong,
.dark .key-meter-value,
.dark .key-expiry-cell strong {
  color: #f4efe4;
}

.dark .key-meter-track {
  background: rgba(48, 52, 43, 0.9);
}

.dark .key-meter-card.tone-warning {
  color: #c7ab73;
}

.dark .key-meter-card.tone-danger {
  color: #f0b4a8;
}

.dark .key-expiry-cell small {
  color: #a4ada2;
}

.dark .key-meter-reset {
  border-color: rgba(48, 52, 43, 0.9);
  background: rgba(17, 19, 15, 0.44);
  color: #d8cdb9;
}

.dark .key-meter-reset:hover,
.dark .key-meter-reset:focus-visible {
  border-color: rgba(167, 58, 42, 0.4);
  background: rgba(167, 58, 42, 0.14);
  color: #f0b4a8;
}

.dark .key-health-lines {
  color: #879186;
}

.dark .key-health-lines-alert,
.dark .key-health-model-hint,
.dark .key-health-link {
  color: #d7988c;
}

.dark .key-health-link:hover,
.dark .key-health-link:focus-visible {
  color: #efb3a8;
}

.dark .key-name-flag,
.dark .key-group-helper,
.dark .key-group-chevron,
.dark .key-group-empty {
  color: #879186;
}

.dark .key-group-trigger {
  border-color: rgba(48, 52, 43, 0.76);
  background: rgba(17, 19, 15, 0.36);
  color: #d8cdb9;
}

.dark .key-group-trigger:hover,
.dark .key-group-trigger:focus-visible {
  border-color: rgba(184, 156, 116, 0.34);
  background: rgba(184, 156, 116, 0.08);
  color: #f4efe4;
}

.dark .key-row-action:hover,
.dark .key-row-action:focus-visible {
  border-color: rgba(167, 58, 42, 0.38);
  background: rgba(167, 58, 42, 0.14);
  color: #f0b4a8;
}

.dark .keys-page :deep(.btn-secondary) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.42);
  color: #f4efe4;
}

@media (max-width: 768px) {
  .keys-toolbar-actions {
    justify-content: stretch;
    margin-left: 0;
    width: 100%;
  }

  .keys-toolbar-actions .btn {
    width: 100%;
  }

  .keys-ledger {
    grid-template-columns: 1fr;
  }

  .keys-ledger-item {
    border-left: 0;
    border-top: 1px solid rgba(216, 205, 185, 0.76);
    padding-left: 0;
    padding-top: 0.85rem;
  }

  .keys-actions .btn {
    flex: 1 1 0;
  }

  .keys-access-strip {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.7rem;
    border-radius: 12px;
    padding: 0.7rem;
  }

  .keys-access-value {
    width: 100%;
  }

  .keys-access-value > strong {
    max-width: none;
  }

  .keys-access-strip-actions {
    width: 100%;
  }

  .keys-access-strip-actions button {
    flex: 1 1 0;
    justify-content: center;
  }

  .key-usage-ledger {
    grid-template-columns: 1fr;
  }

  .key-meter-head,
  .key-meter-foot {
    flex-direction: column;
    align-items: flex-start;
  }

  .key-meter-value {
    text-align: left;
  }
}
</style>
