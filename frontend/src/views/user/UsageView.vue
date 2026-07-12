<template>
  <AppLayout>
    <div class="usage-page space-y-6" :class="{ 'usage-page-en': isEnglishUsage }">
      <section class="usage-hero overflow-hidden rounded-zen border border-zen-paperLine bg-white/45 p-6 shadow-paper dark:border-zen-nightLine dark:bg-zen-nightPanel/70 lg:p-7">
        <div class="grid gap-6 lg:grid-cols-[1fr_auto] lg:items-end">
          <div>
            <div class="mb-4 flex items-center gap-4">
              <span class="h-px w-14 bg-zen-paperLine dark:bg-zen-nightLine"></span>
              <span class="font-mono text-xs uppercase tracking-[0.34em] text-zen-mist dark:text-zen-stone">{{ usageCopy.kicker }}</span>
            </div>
            <h1 class="font-serif text-3xl font-semibold text-zen-ink dark:text-zen-paper sm:text-4xl">{{ usageCopy.title }}</h1>
          </div>

          <div class="usage-ledger grid gap-3 sm:grid-cols-3 lg:min-w-0 lg:max-w-[30rem]">
            <div class="usage-ledger-item">
              <span>{{ usageCopy.records }}</span>
              <strong>{{ pagination.total.toLocaleString() }}</strong>
            </div>
            <div class="usage-ledger-item">
              <span>{{ usageCopy.keyScope }}</span>
              <strong>{{ apiKeys.length.toLocaleString() }}</strong>
            </div>
            <div class="usage-ledger-item">
              <span>{{ usageCopy.cacheHit }}</span>
              <strong>{{ cacheStats.ratePercent }}</strong>
            </div>
          </div>
        </div>
      </section>

      <section class="usage-workspace">
        <div class="usage-summary-shell">
          <div class="usage-summary-grid grid grid-cols-2 gap-4 lg:grid-cols-4">
            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-zen-seal/10 p-2 dark:bg-zen-seal/15">
                  <Icon name="document" size="md" class="text-zen-seal dark:text-zen-paper" />
                </div>
                <div class="usage-summary-copy min-w-0 flex-1">
                  <p class="text-xs font-medium text-zen-mist dark:text-zen-stone">
                    {{ t('usage.totalRequests') }}
                  </p>
                  <p class="text-xl font-bold text-zen-ink dark:text-zen-paper">
                    {{ usageStats?.total_requests?.toLocaleString() || '0' }}
                  </p>
                  <p class="text-xs text-zen-mist dark:text-zen-stone">
                    {{ t('usage.inSelectedRange') }}
                  </p>
                </div>
              </div>
            </div>

            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-zen-seal/10 p-2 dark:bg-zen-seal/15">
                  <Icon name="cube" size="md" class="text-zen-seal dark:text-zen-paper" />
                </div>
                <div class="usage-summary-copy min-w-0 flex-1">
                  <p class="text-xs font-medium text-zen-mist dark:text-zen-stone">
                    {{ t('usage.totalTokens') }}
                  </p>
                  <p class="usage-summary-total-line text-xl font-bold text-zen-ink dark:text-zen-paper">
                    <span>{{ formatTokens(usageStats?.total_tokens || 0) }}</span>
                    <span class="usage-cache-rate-badge" :title="`${usageCopy.cacheRate} ${cacheStats.ratePercent}`">
                      {{ usageCopy.cacheRate }} {{ cacheStats.ratePercent }}
                    </span>
                  </p>
                  <div class="usage-summary-detail usage-summary-detail-token text-xs text-zen-mist dark:text-zen-stone">
                    <p class="usage-summary-token-line">
                      <span>{{ t('usage.in') }} {{ formatTokens(usageStats?.total_input_tokens || 0) }}</span>
                      <span> · </span>
                      <span>{{ t('usage.out') }} {{ formatTokens(usageStats?.total_output_tokens || 0) }}</span>
                    </p>
                    <p class="usage-summary-token-line usage-summary-token-cache-line">
                      <span class="usage-summary-token-cache-values">
                        <span>{{ t('usage.cacheHit') }} {{ formatTokens(usageStats?.total_cache_read_tokens || 0) }}</span>
                        <span> · </span>
                        <span>{{ t('usage.cacheCreate') }} {{ formatTokens(usageStats?.total_cache_creation_tokens || 0) }}</span>
                      </span>
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div class="card usage-summary-card usage-summary-card-primary p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-zen-seal/10 p-2 dark:bg-zen-seal/15">
                  <Icon name="dollar" size="md" class="text-zen-seal dark:text-zen-paper" />
                </div>
                <div class="usage-summary-copy min-w-0 flex-1">
                  <p class="text-xs font-medium text-zen-mist dark:text-zen-stone">
                    {{ t('usage.totalCost') }}
                  </p>
                  <p class="text-xl font-bold text-zen-seal dark:text-zen-paper">
                    ${{ (usageStats?.total_actual_cost || 0).toFixed(4) }}
                  </p>
                  <p class="text-xs text-zen-mist dark:text-zen-stone">
                    {{ t('usage.actualCost') }} / ${{ (usageStats?.total_cost || 0).toFixed(4) }} {{ t('usage.standardCost') }}
                  </p>
                </div>
              </div>
            </div>

            <div class="card p-4">
              <div class="flex items-center gap-3">
                <div class="rounded-lg bg-zen-seal/10 p-2 dark:bg-zen-seal/15">
                  <Icon name="clock" size="md" class="text-zen-seal dark:text-zen-paper" />
                </div>
                <div class="usage-summary-copy min-w-0 flex-1">
                  <p class="text-xs font-medium text-zen-mist dark:text-zen-stone">
                    {{ t('usage.avgDuration') }}
                  </p>
                  <p class="text-xl font-bold text-zen-ink dark:text-zen-paper">
                    {{ formatDuration(usageStats?.average_duration_ms || 0) }}
                  </p>
                  <p class="text-xs text-zen-mist dark:text-zen-stone">{{ t('usage.perRequest') }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="usage-toolbar-shell">
          <div class="usage-filters">
            <div class="usage-toolbar-grid">
              <div class="usage-filter-field usage-filter-field-key">
                <Select
                  v-model="filters.api_key_id"
                  :options="apiKeyOptions"
                  :placeholder="t('usage.allApiKeys')"
                  dropdown-class="keys-filter-dropdown"
                  @change="applyFilters"
                />
              </div>

              <div class="usage-filter-field usage-filter-field-range">
                <DateRangePicker
                  v-model:start-date="startDate"
                  v-model:end-date="endDate"
                  @change="onDateRangeChange"
                />
              </div>

              <!-- Tab 切换栏 -->
              <div v-if="errorViewEnabled" class="usage-tab-strip">
                <button type="button" class="tab usage-tab-button" :class="{ 'tab-active': activeTab === 'usage' }" @click="activeTab = 'usage'">
                  {{ t('usage.tabs.usage') }}
                </button>
                <button type="button" class="tab usage-tab-button" :class="{ 'tab-active': activeTab === 'errors' }" @click="switchToErrors">
                  {{ t('usage.tabs.errors') }}
                </button>
              </div>

              <div class="usage-toolbar-actions">
                <button @click="exportToCSV" :disabled="exporting || pagination.total === 0" class="btn btn-primary">
                  <Icon v-if="exporting" name="refresh" size="sm" class="-ml-1 mr-2 animate-spin" />
                  {{ exporting ? t('usage.exporting') : pagination.total === 0 ? usageCopy.noExportRecords : t('usage.exportCsv') }}
                </button>
                <button @click="applyFilters" :disabled="loading" class="btn btn-secondary">
                  {{ t('common.refresh') }}
                </button>
                <button @click="resetFilters" class="btn btn-secondary">
                  {{ t('common.reset') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="usage-data-shell">
        <!-- 用量明细表 -->
        <!-- flex 链让 DataTable 根 .table-wrapper(flex:1)拿到有界高度以启用内部滚动。
             虚拟化器测高 race 导致的概率空白,已在 DataTable 内用「就绪门控 + initialRect 兜底」根治。 -->
        <div v-show="activeTab === 'usage'" class="flex min-h-0 flex-1 flex-col">
          <DataTable
          :key="`usage-${pagination.page}-${pagination.page_size}`"
          :columns="columns"
          :data="usageLogs"
          :loading="loading"
          :server-side-sort="true"
          :estimate-row-height="88"
          :overscan="12"
          default-sort-key="created_at"
          default-sort-order="desc"
          @sort="handleSort"
        >
          <template #cell-api_key="{ row }">
            <span class="text-sm text-gray-900 dark:text-white">{{
              row.api_key?.name || '-'
            }}</span>
          </template>

          <template #cell-model="{ row, value }">
            <div class="usage-cell-stack usage-cell-break">
              <strong>{{ value }}</strong>
              <span>{{ getUsageServiceTierLabel(row.service_tier, t) }}</span>
            </div>
          </template>

          <template #cell-reasoning_effort="{ row }">
            <span class="usage-inline-chip usage-inline-chip-muted">{{ formatReasoningEffort(row.reasoning_effort) }}</span>
          </template>

          <template #cell-endpoint="{ row }">
            <div class="usage-cell-stack usage-cell-stack-compact usage-cell-break">
              <strong>{{ formatUsageEndpoints(row) }}</strong>
              <span>{{ row.upstream_endpoint?.trim() || t('usage.endpoint') }}</span>
            </div>
          </template>

          <template #cell-stream="{ row }">
            <span class="usage-inline-chip usage-inline-chip-muted">
              {{ getRequestTypeLabel(row) }}
            </span>
          </template>

          <template #cell-billing_mode="{ row }">
            <span class="usage-inline-chip usage-inline-chip-muted">
              {{ getBillingModeLabel(getDisplayBillingMode(row), t) }}
            </span>
          </template>

          <template #cell-tokens="{ row }">
            <!-- 图片生成请求 -->
            <div v-if="isImageUsage(row)" class="usage-token-cell usage-token-cell-image">
              <Icon name="image" size="sm" class="text-indigo-500" :stroke-width="2" />
              <div class="usage-cell-stack usage-cell-stack-compact">
                <strong>{{ row.image_count }}{{ t('usage.imageUnit') }}</strong>
                <span>{{ formatImageBillingSize(row, t) }}</span>
              </div>
            </div>
            <!-- Token 请求 -->
            <div v-else class="usage-token-cell">
              <div class="usage-token-ledger usage-token-ledger-compact">
                <div class="usage-token-stat">
                  <span>IN</span>
                  <strong>{{ (row.input_tokens ?? 0).toLocaleString() }}</strong>
                </div>
                <div class="usage-token-stat">
                  <span>OUT</span>
                  <strong>{{ (row.output_tokens ?? 0).toLocaleString() }}</strong>
                </div>
              </div>
              <div
                v-if="row.cache_read_tokens > 0 || row.cache_creation_tokens > 0 || hasImageOutputTokens(row)"
                class="usage-token-meta"
              >
                <span v-if="row.cache_read_tokens > 0" class="usage-meta-chip">
                  Cache read {{ formatCacheTokens(row.cache_read_tokens) }}
                </span>
                <span v-if="row.cache_creation_tokens > 0" class="usage-meta-chip">
                  Cache write {{ formatCacheTokens(row.cache_creation_tokens) }}
                </span>
                <span v-if="row.cache_creation_1h_tokens > 0" class="usage-meta-chip">TTL 1h</span>
                <span v-if="row.cache_ttl_overridden" :title="t('usage.cacheTtlOverriddenHint')" class="usage-meta-chip">TTL reset</span>
                <span v-if="hasImageOutputTokens(row)" class="usage-meta-chip">
                  Image {{ row.image_output_tokens.toLocaleString() }}
                </span>
              </div>
            </div>
          </template>

          <template #cell-cost="{ row }">
            <div class="usage-cost-cell">
              <div class="usage-cell-stack usage-cell-stack-compact">
                <strong>${{ (row.actual_cost ?? 0).toFixed(6) }}</strong>
                <span>{{ usageCopy.standardCostPrefix }} ${{ (row.total_cost ?? 0).toFixed(6) }}</span>
              </div>
              <!-- Cost Detail Tooltip -->
              <div
                class="group relative"
                @mouseenter="showTooltip($event, row)"
                @mouseleave="hideTooltip"
              >
                <div
                  class="flex h-4 w-4 cursor-help items-center justify-center rounded-full bg-gray-100 transition-colors group-hover:bg-blue-100 dark:bg-gray-700 dark:group-hover:bg-blue-900/50"
                >
                  <Icon
                    name="infoCircle"
                    size="xs"
                    class="text-gray-400 group-hover:text-blue-500 dark:text-gray-500 dark:group-hover:text-blue-400"
                  />
                </div>
              </div>
            </div>
          </template>

          <template #cell-first_token="{ row }">
            <div class="usage-cell-stack usage-cell-stack-compact">
              <strong v-if="row.first_token_ms != null">{{ formatDuration(row.first_token_ms) }}</strong>
              <strong v-else>-</strong>
              <span>{{ row.stream ? usageCopy.firstTokenLatency : usageCopy.nonStreamRequest }}</span>
            </div>
          </template>

          <template #cell-duration="{ row }">
            <div class="usage-cell-stack usage-cell-stack-compact">
              <strong>{{ formatDuration(row.duration_ms) }}</strong>
              <span>{{ row.stream ? usageCopy.fullResponse : usageCopy.singleCompletion }}</span>
            </div>
          </template>

          <template #cell-created_at="{ value }">
            <div class="usage-cell-stack usage-cell-stack-compact">
              <strong>{{ formatUsageDate(value) }}</strong>
              <span>{{ formatUsageTime(value) }}</span>
            </div>
          </template>

          <template #cell-user_agent="{ row }">
            <div v-if="row.user_agent" class="usage-cell-stack usage-cell-stack-compact usage-user-agent" :title="row.user_agent">
              <strong>{{ formatUserAgent(row.user_agent).primary }}</strong>
              <span>{{ formatUserAgent(row.user_agent).secondary || row.api_key?.name || usageCopy.clientSource }}</span>
            </div>
            <span v-else class="text-sm text-gray-400 dark:text-gray-500">-</span>
          </template>

          <template #empty>
            <div class="usage-empty-state">
              <span>{{ usageCopy.emptyKicker }}</span>
              <strong>{{ t('usage.noRecords') }}</strong>
              <p>{{ usageCopy.emptyCopy }}</p>
              <button type="button" @click="resetFilters">{{ t('common.reset') }}</button>
            </div>
          </template>
        </DataTable>
        </div>

        <!-- 错误请求表 -->
        <div v-if="errorViewEnabled" v-show="activeTab === 'errors'" class="flex min-h-0 flex-1 flex-col">
          <UserErrorRequestsTable
            :rows="errorRows"
            :total="errorTotal"
            :loading="errorLoading"
            :page="errorPage"
            :page-size="errorPageSize"
            :api-keys="apiKeys"
            :model-filter="errorFilter.model"
            :category-filter="errorFilter.category"
            :api-key-id-filter="errorFilter.api_key_id"
            @filter="onErrorFilter"
            @update:page="onErrorPage"
            @update:pageSize="onErrorPageSize"
          />
        </div>

          <div v-if="pagination.total > 0 && activeTab === 'usage'" class="usage-pagination-shell">
            <Pagination
              :page="pagination.page"
              :total="pagination.total"
              :page-size="pagination.page_size"
              :show-controls-when-single-page="false"
              page-size-dropdown-class="keys-page-size-dropdown"
              @update:page="handlePageChange"
              @update:pageSize="handlePageSizeChange"
            />
          </div>
        </div>
      </section>
    </div>
  </AppLayout>

  <!-- Tooltip Portal -->
  <Teleport to="body">
    <div
      v-if="tooltipVisible"
      class="fixed z-[9999] pointer-events-none -translate-y-1/2"
      :style="{
        left: tooltipPosition.x + 'px',
        top: tooltipPosition.y + 'px'
      }"
    >
      <div
        class="whitespace-nowrap rounded-lg border border-gray-700 bg-gray-900 px-3 py-2.5 text-xs text-white shadow-xl dark:border-gray-600 dark:bg-gray-800"
      >
        <div class="space-y-1.5">
          <!-- Cost Breakdown -->
          <div class="mb-2 border-b border-gray-700 pb-1.5">
            <div class="text-xs font-semibold text-gray-300 mb-1">{{ t('usage.costDetails') }}</div>
            <div v-if="tooltipData && tooltipData.input_cost > 0" class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('admin.usage.inputCost') }}</span>
              <span class="font-medium text-white">${{ tooltipData.input_cost.toFixed(6) }}</span>
            </div>
            <div v-if="tooltipData && tooltipData.output_cost > 0" class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('admin.usage.outputCost') }}</span>
              <span class="font-medium text-white">${{ tooltipData.output_cost.toFixed(6) }}</span>
            </div>
            <div v-if="tooltipData && hasImageOutputCost(tooltipData)" class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('usage.imageOutputCost') }}</span>
              <span class="font-medium text-pink-300">${{ tooltipData.image_output_cost.toFixed(6) }}</span>
            </div>
            <!-- Token billing: show unit prices per 1M tokens -->
            <template v-if="tooltipData && !isImageUsage(tooltipData) && (!tooltipData.billing_mode || tooltipData.billing_mode === BILLING_MODE_TOKEN)">
              <div v-if="tooltipData && tooltipData.input_tokens > 0" class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.inputTokenPrice') }}</span>
                <span class="font-medium text-sky-300">{{ formatTokenPricePerMillion(tooltipData.input_cost, tooltipData.input_tokens) }} {{ t('usage.perMillionTokens') }}</span>
              </div>
              <div v-if="tooltipData && tooltipData.output_cost > 0 && textOutputTokens(tooltipData) > 0" class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.outputTokenPrice') }}</span>
                <span class="font-medium text-violet-300">{{ formatTokenPricePerMillion(tooltipData.output_cost, textOutputTokens(tooltipData)) }} {{ t('usage.perMillionTokens') }}</span>
              </div>
              <div v-if="tooltipData && hasImageOutputTokens(tooltipData)" class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageOutputTokenPrice') }}</span>
                <span class="font-medium text-pink-300">{{ formatTokenPricePerMillion(tooltipData.image_output_cost ?? 0, tooltipData.image_output_tokens) }} {{ t('usage.perMillionTokens') }}</span>
              </div>
            </template>
            <!-- Per-image billing: show image metadata and unit price -->
            <template v-else-if="tooltipData && isImageUsage(tooltipData)">
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageCount') }}</span>
                <span class="font-medium text-white">{{ tooltipData.image_count }}{{ t('usage.imageUnit') }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageBillingSize') }}</span>
                <span class="font-medium text-white">{{ formatImageBillingSize(tooltipData, t) }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageSizeSource') }}</span>
                <span class="font-medium text-white">{{ formatImageSizeSource(tooltipData, t) }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageInputSize') }}</span>
                <span class="font-medium text-white">{{ formatImageInputSize(tooltipData, t) }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageOutputSize') }}</span>
                <span class="font-medium text-white">{{ formatImageOutputSize(tooltipData, t) }}</span>
              </div>
              <div v-if="formatImageSizeBreakdown(tooltipData)" class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageSizeBreakdown') }}</span>
                <span class="font-medium text-white">{{ formatImageSizeBreakdown(tooltipData) }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageUnitPrice') }}</span>
                <span class="font-medium text-sky-300">${{ imageUnitPrice(tooltipData).toFixed(6) }}</span>
              </div>
              <div class="flex items-center justify-between gap-4">
                <span class="text-gray-400">{{ t('usage.imageTotalPrice') }}</span>
                <span class="font-medium text-white">${{ tooltipData.total_cost?.toFixed(6) || '0.000000' }}</span>
              </div>
            </template>
            <div v-else class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('usage.unitPrice') }}</span>
              <span class="font-medium text-sky-300">${{ tooltipData?.total_cost?.toFixed(6) || '0.000000' }}</span>
            </div>
            <div v-if="tooltipData && tooltipData.cache_creation_cost > 0" class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('admin.usage.cacheCreationCost') }}</span>
              <span class="font-medium text-white">${{ tooltipData.cache_creation_cost.toFixed(6) }}</span>
            </div>
            <div v-if="tooltipData && tooltipData.cache_read_cost > 0" class="flex items-center justify-between gap-4">
              <span class="text-gray-400">{{ t('admin.usage.cacheReadCost') }}</span>
              <span class="font-medium text-white">${{ tooltipData.cache_read_cost.toFixed(6) }}</span>
            </div>
          </div>
          <!-- Rate and Summary -->
          <div class="flex items-center justify-between gap-6">
            <span class="text-gray-400">{{ t('usage.serviceTier') }}</span>
            <span class="font-semibold text-cyan-300">{{ getUsageServiceTierLabel(tooltipData?.service_tier, t) }}</span>
          </div>
          <div class="flex items-center justify-between gap-6">
            <span class="text-gray-400">{{ t('usage.rate') }}</span>
            <span class="font-semibold text-blue-400"
              >{{ formatMultiplier(tooltipData?.rate_multiplier || 1) }}x</span
            >
          </div>
          <div class="flex items-center justify-between gap-6">
            <span class="text-gray-400">{{ t('usage.original') }}</span>
            <span class="font-medium text-white">${{ tooltipData?.total_cost.toFixed(6) }}</span>
          </div>
          <div class="flex items-center justify-between gap-6 border-t border-gray-700 pt-1.5">
            <span class="text-gray-400">{{ t('usage.billed') }}</span>
            <span class="font-semibold text-green-400"
              >${{ tooltipData?.actual_cost.toFixed(6) }}</span
            >
          </div>
        </div>
        <!-- Tooltip Arrow (left side) -->
        <div
          class="absolute right-full top-1/2 h-0 w-0 -translate-y-1/2 border-b-[6px] border-r-[6px] border-t-[6px] border-b-transparent border-r-gray-900 border-t-transparent dark:border-r-gray-800"
        ></div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { usageAPI, keysAPI } from '@/api'
import AppLayout from '@/components/layout/AppLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import DateRangePicker from '@/components/common/DateRangePicker.vue'
import Icon from '@/components/icons/Icon.vue'
import UserErrorRequestsTable from '@/components/user/UserErrorRequestsTable.vue'
import type { UsageLog, ApiKey, UsageQueryParams, UsageStatsResponse, UserErrorRequest } from '@/types'
import type { Column } from '@/components/common/types'
import { formatDateOnly, formatReasoningEffort, formatTime } from '@/utils/format'
import { getPersistedPageSize } from '@/composables/usePersistedPageSize'
import { formatCacheTokens, formatMultiplier } from '@/utils/formatters'
import { formatTokenPricePerMillion } from '@/utils/usagePricing'
import { getUsageServiceTierLabel } from '@/utils/usageServiceTier'
import { resolveUsageRequestType } from '@/utils/usageRequestType'
import {
  BILLING_MODE_TOKEN,
  getBillingModeLabel,
  isImageUsage,
  getDisplayBillingMode,
  imageUnitPrice,
} from '@/utils/billingMode'
import {
  formatImageBillingSize,
  formatImageInputSize,
  formatImageOutputSize,
  formatImageSizeBreakdown,
  formatImageSizeSource,
  hasImageOutputTokens,
  textOutputTokens,
  hasImageOutputCost,
} from '@/utils/imageUsage'

const { t, locale } = useI18n()
const route = useRoute()
const appStore = useAppStore()

let abortController: AbortController | null = null

const zhUsageCopy = {
  kicker: '用量账册',
  title: '流转明细',
  records: '当前记录',
  keyScope: '密钥范围',
  cacheHit: '缓存命中',
  cacheRate: '缓存率',
  noExportRecords: '当前范围暂无记录',
  standardCostPrefix: '标准',
  firstTokenLatency: '首包耗时',
  nonStreamRequest: '非流式请求',
  fullResponse: '完整响应',
  singleCompletion: '单次完成',
  clientSource: '客户端来源',
  emptyKicker: '账册暂无记录',
  emptyCopy: '当前筛选范围内还没有调用明细。可以切换时间范围、改看全部密钥，或重置筛选后再查看。',
  clientLabels: {
    desktop: '桌面客户端',
    anthropic: 'Anthropic 客户端',
    openai: 'OpenAI 客户端',
    editor: '编辑器客户端',
    browser: '浏览器请求',
    source: '客户端来源'
  }
}

const enUsageCopy = {
  kicker: 'Usage ledger',
  title: 'Flow details',
  records: 'Records',
  keyScope: 'Key scope',
  cacheHit: 'Cache hit',
  cacheRate: 'Cache rate',
  noExportRecords: 'No records in this range',
  standardCostPrefix: 'Standard',
  firstTokenLatency: 'First-token latency',
  nonStreamRequest: 'Non-stream request',
  fullResponse: 'Full response',
  singleCompletion: 'Single completion',
  clientSource: 'Client source',
  emptyKicker: 'No ledger records',
  emptyCopy: 'No usage details match the current filters. Try another date range, show all keys, or reset filters.',
  clientLabels: {
    desktop: 'Desktop client',
    anthropic: 'Anthropic client',
    openai: 'OpenAI client',
    editor: 'Editor client',
    browser: 'Browser request',
    source: 'Client source'
  }
}

const usageCopy = computed(() => locale.value === 'zh' ? zhUsageCopy : enUsageCopy)
const isEnglishUsage = computed(() => !locale.value.startsWith('zh'))

// Tooltip state
const tooltipVisible = ref(false)
const tooltipPosition = ref({ x: 0, y: 0 })
const tooltipData = ref<UsageLog | null>(null)

// Usage stats from API
const usageStats = ref<UsageStatsResponse | null>(null)

// 缓存命中率 = cache_read / (input + cache_creation + cache_read)
// 分母为 0（无任何输入）时显示 0.0%
const cacheStats = computed(() => {
  // 总输入 token = 普通输入 + 缓存写入 + 缓存读取（命中）
  // 缓存命中率 = 缓存读取 / 总输入。
  const cacheRead = usageStats.value?.total_cache_read_tokens || 0
  const cacheCreate = usageStats.value?.total_cache_creation_tokens || 0
  const input = usageStats.value?.total_input_tokens || 0
  const totalInput = input + cacheCreate + cacheRead
  const ratePercent = totalInput > 0 ? `${((cacheRead / totalInput) * 100).toFixed(1)}%` : '0.0%'
  return { cacheRead, totalInput, ratePercent }
})

const columns = computed<Column[]>(() => [
  { key: 'api_key', label: t('usage.apiKeyFilter'), sortable: false, class: 'usage-col-api-key' },
  { key: 'model', label: t('usage.model'), sortable: true, class: 'usage-col-model' },
  { key: 'reasoning_effort', label: t('usage.reasoningEffort'), sortable: false, class: 'usage-col-reasoning' },
  { key: 'endpoint', label: t('usage.endpoint'), sortable: false, class: 'usage-col-endpoint' },
  { key: 'stream', label: t('usage.type'), sortable: false, class: 'usage-col-type' },
  { key: 'billing_mode', label: t('admin.usage.billingMode'), sortable: false, class: 'usage-col-billing-mode' },
  { key: 'tokens', label: t('usage.tokens'), sortable: false, class: 'usage-col-tokens' },
  { key: 'cost', label: t('usage.cost'), sortable: false, class: 'usage-col-cost' },
  { key: 'first_token', label: t('usage.firstToken'), sortable: false, class: 'usage-col-first-token' },
  { key: 'duration', label: t('usage.duration'), sortable: false, class: 'usage-col-duration' },
  { key: 'created_at', label: t('usage.time'), sortable: true, class: 'usage-col-created-at' },
  { key: 'user_agent', label: usageCopy.value.clientSource, sortable: false, class: 'usage-col-user-agent' }
])

const usageLogs = ref<UsageLog[]>([])
const apiKeys = ref<ApiKey[]>([])
const loading = ref(false)
const exporting = ref(false)

const apiKeyOptions = computed(() => {
  return [
    { value: null, label: t('usage.allApiKeys') },
    ...apiKeys.value.map((key) => ({
      value: key.id,
      label: key.name
    }))
  ]
})

// Helper function to format date in local timezone
const formatLocalDate = (date: Date): string => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// Initialize date range immediately
const now = new Date()
const weekAgo = new Date(now)
weekAgo.setDate(weekAgo.getDate() - 6)

// Date range state
const startDate = ref(formatLocalDate(weekAgo))
const endDate = ref(formatLocalDate(now))

const filters = ref<UsageQueryParams>({
  api_key_id: undefined,
  start_date: undefined,
  end_date: undefined
})

// Initialize filters with date range
filters.value.start_date = startDate.value
filters.value.end_date = endDate.value

// Handle date range change from DateRangePicker
const onDateRangeChange = (range: {
  startDate: string
  endDate: string
  preset: string | null
}) => {
  filters.value.start_date = range.startDate
  filters.value.end_date = range.endDate
  applyFilters()
  errorPage.value = 1
  if (activeTab.value === 'errors') {
    loadErrors()
  } else {
    errorRows.value = []  // 失效，下次切到 errors tab 时按新日期重新加载
  }
}

const pagination = reactive({
  page: 1,
  page_size: getPersistedPageSize(),
  total: 0,
  pages: 0
})
const sortState = reactive({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
})

const formatDuration = (ms: number | null | undefined): string => {
  if (ms == null) return '-'
  if (ms < 1000) return `${ms.toFixed(0)}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

type FormattedUserAgent = {
  primary: string
  secondary: string
}

const extractUaVersion = (source: string): string => {
  const match = source.match(/\/([0-9]+(?:\.[0-9A-Za-z]+){0,3})/)
  return match?.[1] || ''
}

const formatUserAgent = (ua: string): FormattedUserAgent => {
  const raw = ua.trim()
  if (!raw) {
    return { primary: '-', secondary: '' }
  }

  const normalized = raw.replace(/\s+/g, ' ')
  const lower = normalized.toLowerCase()

  const platformParts: string[] = []
  if (lower.includes('windows')) platformParts.push('Windows')
  else if (lower.includes('mac os') || lower.includes('macintosh')) platformParts.push('macOS')
  else if (lower.includes('android')) platformParts.push('Android')
  else if (lower.includes('iphone') || lower.includes('ios')) platformParts.push('iOS')
  else if (lower.includes('linux')) platformParts.push('Linux')

  if (lower.includes('x86_64') || lower.includes('win64') || lower.includes('x64')) {
    platformParts.push('64-bit')
  } else if (lower.includes('arm64') || lower.includes('aarch64')) {
    platformParts.push('ARM64')
  }

  let primary = normalized
  let secondary = platformParts.join(' · ')

  const codexMatch = normalized.match(/(Codex Desktop)\/([^\s(]+)/i)
  if (codexMatch) {
    primary = `${codexMatch[1]} ${codexMatch[2]}`
    secondary = secondary || usageCopy.value.clientLabels.desktop
    return { primary, secondary }
  }

  const claudeMatch = normalized.match(/(Claude(?: Code)?)/i)
  if (claudeMatch) {
    primary = claudeMatch[1]
    const version = extractUaVersion(normalized)
    if (version) primary = `${primary} ${version}`
    secondary = secondary || usageCopy.value.clientLabels.anthropic
    return { primary, secondary }
  }

  const openAiMatch = normalized.match(/(ChatGPT|OpenAI|openai-node|openai-python)/i)
  if (openAiMatch) {
    primary = openAiMatch[1]
    const version = extractUaVersion(normalized)
    if (version) primary = `${primary} ${version}`
    secondary = secondary || usageCopy.value.clientLabels.openai
    return { primary, secondary }
  }

  const cursorMatch = normalized.match(/(Cursor)/i)
  if (cursorMatch) {
    primary = cursorMatch[1]
    const version = extractUaVersion(normalized)
    if (version) primary = `${primary} ${version}`
    secondary = secondary || usageCopy.value.clientLabels.editor
    return { primary, secondary }
  }

  const browserMatch = normalized.match(/(Chrome|Edg|Edge|Firefox|Safari)\/([^\s]+)/i)
  if (browserMatch) {
    const browser = browserMatch[1] === 'Edg' ? 'Edge' : browserMatch[1]
    primary = `${browser} ${browserMatch[2]}`
    secondary = secondary || usageCopy.value.clientLabels.browser
    return { primary, secondary }
  }

  if (normalized.includes('/')) {
    const firstToken = normalized.split(' ')[0]
    const [name, version] = firstToken.split('/', 2)
    if (name) {
      primary = version ? `${name} ${version}` : name
    }
  } else if (normalized.length > 36) {
    primary = `${normalized.slice(0, 33)}...`
  }

  return {
    primary,
    secondary: secondary || usageCopy.value.clientLabels.source
  }
}

const formatUsageDate = (date: string | null | undefined): string => {
  return formatDateOnly(date) || '-'
}

const formatUsageTime = (date: string | null | undefined): string => {
  return formatTime(date) || '-'
}

const getRequestTypeLabel = (log: UsageLog): string => {
  const requestType = resolveUsageRequestType(log)
  if (requestType === 'cyber') return t('usage.cyber')
  if (requestType === 'ws_v2') return t('usage.ws')
  if (requestType === 'stream') return t('usage.stream')
  if (requestType === 'sync') return t('usage.sync')
  return t('usage.unknown')
}

const getRequestTypeExportText = (log: UsageLog): string => {
  const requestType = resolveUsageRequestType(log)
  if (requestType === 'cyber') return 'Cyber'
  if (requestType === 'ws_v2') return 'WS'
  if (requestType === 'stream') return 'Stream'
  if (requestType === 'sync') return 'Sync'
  return 'Unknown'
}

const formatUsageEndpoints = (log: UsageLog): string => {
  const inbound = log.inbound_endpoint?.trim()
  return inbound || '-'
}

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B`
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M`
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K`
  }
  return value.toLocaleString()
}

type UsageTableQueryParams = UsageQueryParams & {
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

const buildUsageQueryParams = (page: number, pageSize: number): UsageTableQueryParams => ({
  page,
  page_size: pageSize,
  ...filters.value,
  sort_by: sortState.sort_by,
  sort_order: sortState.sort_order
})

const loadUsageLogs = async () => {
  if (abortController) {
    abortController.abort()
  }
  const currentAbortController = new AbortController()
  abortController = currentAbortController
  const { signal } = currentAbortController
  loading.value = true
  try {
    const response = await usageAPI.query(
      buildUsageQueryParams(pagination.page, pagination.page_size),
      { signal }
    )
    if (signal.aborted) {
      return
    }
    usageLogs.value = response.items
    pagination.page = response.page || pagination.page
    pagination.page_size = response.page_size || pagination.page_size
    pagination.total = response.total
    pagination.pages = response.pages
  } catch (error) {
    if (signal.aborted) {
      return
    }
    const abortError = error as { name?: string; code?: string }
    if (abortError?.name === 'AbortError' || abortError?.code === 'ERR_CANCELED') {
      return
    }
    appStore.showError(t('usage.failedToLoad'))
  } finally {
    if (abortController === currentAbortController) {
      loading.value = false
    }
  }
}

const loadApiKeys = async () => {
  try {
    const response = await keysAPI.list(1, 100)
    apiKeys.value = response.items
  } catch (error) {
    console.error('Failed to load API keys:', error)
  }
}

const loadUsageStats = async () => {
  try {
    const apiKeyId = filters.value.api_key_id ? Number(filters.value.api_key_id) : undefined
    const stats = await usageAPI.getStatsByDateRange(
      filters.value.start_date || startDate.value,
      filters.value.end_date || endDate.value,
      apiKeyId
    )
    usageStats.value = stats
  } catch (error) {
    console.error('Failed to load usage stats:', error)
  }
}

const applyFilters = () => {
  pagination.page = 1
  loadUsageLogs()
  loadUsageStats()
}

const resetFilters = () => {
  filters.value = {
    api_key_id: undefined,
    start_date: undefined,
    end_date: undefined
  }
  // Reset date range to default (last 7 days)
  const now = new Date()
  const weekAgo = new Date(now)
  weekAgo.setDate(weekAgo.getDate() - 6)
  startDate.value = formatLocalDate(weekAgo)
  endDate.value = formatLocalDate(now)
  filters.value.start_date = startDate.value
  filters.value.end_date = endDate.value
  sortState.sort_by = 'created_at'
  sortState.sort_order = 'desc'
  pagination.page = 1
  loadUsageLogs()
  loadUsageStats()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadUsageLogs()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.page_size = pageSize
  pagination.page = 1
  loadUsageLogs()
}

const handleSort = (key: string, order: 'asc' | 'desc') => {
  sortState.sort_by = key
  sortState.sort_order = order
  pagination.page = 1
  loadUsageLogs()
}

/**
 * Escape CSV value to prevent injection and handle special characters
 */
const escapeCSVValue = (value: unknown): string => {
  if (value == null) return ''

  const str = String(value)
  const escaped = str.replace(/"/g, '""')

  // Prevent formula injection by prefixing dangerous characters with single quote
  if (/^[=+\-@\t\r]/.test(str)) {
    return `"\'${escaped}"`
  }

  // Escape values containing comma, quote, or newline
  if (/[,"\n\r]/.test(str)) {
    return `"${escaped}"`
  }

  return str
}

const exportToCSV = async () => {
  if (pagination.total === 0) {
    appStore.showWarning(t('usage.noDataToExport'))
    return
  }

  exporting.value = true
  appStore.showInfo(t('usage.preparingExport'))

  try {
    const allLogs: UsageLog[] = []
    const pageSize = 100 // Use a larger page size for export to reduce requests
    const totalRequests = Math.ceil(pagination.total / pageSize)

    for (let page = 1; page <= totalRequests; page++) {
      const response = await usageAPI.query(buildUsageQueryParams(page, pageSize))
      allLogs.push(...response.items)
    }

    if (allLogs.length === 0) {
      appStore.showWarning(t('usage.noDataToExport'))
      return
    }

    const headers = [
      'Time',
      'API Key Name',
      'Model',
      'Reasoning Effort',
      'Inbound Endpoint',
      'Type',
      'Billing Mode',
      'Input Tokens',
      'Output Tokens',
      'Cache Read Tokens',
      'Cache Creation Tokens',
      'Rate Multiplier',
      'Billed Cost',
      'Original Cost',
      'First Token (ms)',
      'Duration (ms)'
    ]
    const rows = allLogs.map((log) =>
      [
        log.created_at,
        log.api_key?.name || '',
        log.model,
        formatReasoningEffort(log.reasoning_effort),
        log.inbound_endpoint || '',
        getRequestTypeExportText(log),
        getBillingModeLabel(getDisplayBillingMode(log), t),
        log.input_tokens,
        log.output_tokens,
        log.cache_read_tokens,
        log.cache_creation_tokens,
        log.rate_multiplier,
        (log.actual_cost ?? 0).toFixed(8),
        (log.total_cost ?? 0).toFixed(8),
        log.first_token_ms ?? '',
        log.duration_ms
      ].map(escapeCSVValue)
    )

    const csvContent = [
      headers.map(escapeCSVValue).join(','),
      ...rows.map((row) => row.join(','))
    ].join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `usage_${filters.value.start_date}_to_${filters.value.end_date}.csv`
    link.click()
    window.URL.revokeObjectURL(url)

    appStore.showSuccess(t('usage.exportSuccess'))
  } catch (error) {
    appStore.showError(t('usage.exportFailed'))
    console.error('CSV Export failed:', error)
  } finally {
    exporting.value = false
  }
}

// Tooltip functions
const showTooltip = (event: MouseEvent, row: UsageLog) => {
  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()

  tooltipData.value = row
  // Position to the right of the icon, vertically centered
  tooltipPosition.value.x = rect.right + 8
  tooltipPosition.value.y = rect.top + rect.height / 2
  tooltipVisible.value = true
}

const hideTooltip = () => {
  tooltipVisible.value = false
  tooltipData.value = null
}

// ── Error Requests Tab ──────────────────────────────────────────────────────
const activeTab = ref<'usage' | 'errors'>('usage')
const errorViewEnabled = computed(() => appStore.cachedPublicSettings?.allow_user_view_error_requests ?? false)

const errorRows = ref<UserErrorRequest[]>([])
const errorLoading = ref(false)
const errorPage = ref(1)
const errorPageSize = ref(20)
const errorTotal = ref(0)
const errorFilter = ref<{ model: string; category: string; api_key_id: number | null }>({ model: '', category: '', api_key_id: null })

const parsePositiveQueryInt = (value: unknown): number | null => {
  const raw = Array.isArray(value) ? value[0] : value
  if (typeof raw !== 'string' || raw.trim() === '') return null
  const parsed = Number.parseInt(raw, 10)
  return Number.isFinite(parsed) && parsed > 0 ? parsed : null
}

const loadErrors = async () => {
  errorLoading.value = true
  try {
    const resp = await usageAPI.listMyErrorRequests({
      page: errorPage.value,
      page_size: errorPageSize.value,
      start_date: startDate.value,
      end_date: endDate.value,
      model: errorFilter.value.model || undefined,
      category: errorFilter.value.category || undefined,
      api_key_id: errorFilter.value.api_key_id ?? undefined,
    })
    errorRows.value = resp.items
    errorTotal.value = resp.total
  } catch (error) {
    console.error('[UsageView] loadErrors failed:', error)
    appStore.showError(t('usage.errors.failedToLoad'))
  } finally {
    errorLoading.value = false
  }
}

const onErrorFilter = (f: { model: string; category: string; api_key_id: number | null }) => {
  errorFilter.value = f
  errorPage.value = 1
  loadErrors()
}
const onErrorPage = (p: number) => { errorPage.value = p; loadErrors() }
const onErrorPageSize = (s: number) => { errorPageSize.value = s; errorPage.value = 1; loadErrors() }

const switchToErrors = () => {
  activeTab.value = 'errors'
  if (errorRows.value.length === 0) loadErrors()
}

watch(
  () => route?.query ?? {},
  (query) => {
    const tab = typeof query.tab === 'string' ? query.tab : ''
    const category = typeof query.category === 'string' ? query.category : ''
    const apiKeyId = parsePositiveQueryInt(query.api_key_id)
    if (tab === 'errors') {
      activeTab.value = 'errors'
      if (
        errorFilter.value.category !== category ||
        errorFilter.value.api_key_id !== apiKeyId
      ) {
        errorFilter.value = { ...errorFilter.value, category, api_key_id: apiKeyId }
        filters.value.api_key_id = apiKeyId ?? undefined
        errorPage.value = 1
      }
      loadErrors()
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadApiKeys()
  loadUsageLogs()
  loadUsageStats()
})
</script>

<style scoped>
.usage-page {
  color: #1f2320;
}

.usage-hero {
  position: relative;
  background-image:
    linear-gradient(90deg, rgba(250, 247, 239, 0.78), rgba(250, 247, 239, 0.42)),
    linear-gradient(rgba(31, 35, 32, 0.024) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.018) 1px, transparent 1px);
  background-size: auto, 86px 86px, 86px 86px;
}

.usage-hero::after {
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

.usage-ledger-item {
  border-left: 1px solid rgba(216, 205, 185, 0.95);
  padding-left: 1rem;
}

.usage-ledger-item span {
  display: block;
  font-size: 0.74rem;
  color: #59645a;
  overflow-wrap: anywhere;
}

.usage-ledger-item strong {
  display: block;
  margin-top: 0.35rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 1.35rem;
  color: #1f2320;
}

.usage-workspace {
  display: grid;
  gap: 1rem;
}

.usage-summary-shell,
.usage-toolbar-shell,
.usage-data-shell {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  min-width: 0;
}

.usage-summary-shell,
.usage-toolbar-shell {
  padding: 0;
  min-width: 0;
}

.usage-toolbar-shell {
  display: grid;
  gap: 0.75rem;
}

.usage-filters {
  display: grid;
  gap: 0.78rem;
  min-width: 0;
}

.usage-toolbar-grid {
  display: grid;
  grid-template-columns: minmax(12rem, 18rem) minmax(12rem, max-content) auto minmax(14rem, 1fr);
  align-items: center;
  gap: 0.65rem;
  min-width: 0;
}

.usage-filter-field {
  display: grid;
  gap: 0;
  min-width: 0;
}

.usage-filter-field-key {
  width: 100%;
  min-width: 0;
}

.usage-filter-field-range {
  width: 100%;
  min-width: 0;
}

.usage-filter-field-range :deep(.relative) {
  display: flex;
  width: 100%;
  min-width: 0;
}

.usage-filter-field-range :deep(.date-picker-trigger) {
  display: inline-flex;
  width: 100%;
  min-width: 0;
  gap: 0.5rem;
  min-height: 2.8rem;
  padding-inline: 0.95rem;
  border-color: rgba(216, 205, 185, 0.58);
  border-radius: 12px;
  background: rgba(255, 252, 246, 0.34);
  color: #38413a;
  box-shadow: none;
  white-space: nowrap;
}

.usage-filter-field-range :deep(.date-picker-value) {
  color: #1f2320;
  font-weight: 650;
  white-space: nowrap;
}

.usage-filter-field-range :deep(.date-picker-icon),
.usage-filter-field-range :deep(.date-picker-chevron),
.usage-filter-field-range :deep(svg) {
  color: #8d978a;
}

.usage-toolbar-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.55rem;
  margin-left: 0;
  min-width: 0;
  justify-self: end;
}

.usage-toolbar-actions .btn {
  min-width: 0;
  max-width: 100%;
  line-height: 1.2;
  white-space: normal;
}

.usage-data-shell {
  display: grid;
  gap: 0.72rem;
  padding: 0;
}

.usage-pagination-shell {
  margin-top: 0;
  min-width: 0;
}

.usage-pagination-shell :deep(.pagination-root) {
  border-top: 1px solid rgba(216, 205, 185, 0.58);
  border-radius: 0;
  background: transparent;
  padding-inline: 0;
  min-width: 0;
}

.usage-pagination-shell :deep(.pagination-desktop) {
  align-items: center;
  justify-content: flex-start;
  gap: 1rem;
  min-width: 0;
}

.usage-pagination-shell :deep(.pagination-meta) {
  flex: 1 1 auto;
  min-width: 0;
}

.usage-pagination-shell :deep(nav),
.usage-pagination-shell :deep(.pagination-controls) {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
  border: 0;
  border-radius: 0;
  background: transparent;
  padding: 0;
  box-shadow: none;
  min-width: 0;
}

.usage-pagination-shell :deep(.page-size-select) {
  width: 5.9rem;
  min-width: 5.9rem;
  flex: 0 0 5.9rem;
}

.usage-pagination-shell :deep(.page-size-select .select-trigger) {
  min-height: 2.04rem;
  padding-inline: 0.7rem;
  border-color: rgba(216, 205, 185, 0.5);
  border-radius: 999px;
  background: rgba(250, 247, 239, 0.24);
  color: #38413a;
  box-shadow: none;
}

.usage-pagination-shell :deep(.page-size-select .select-value) {
  min-width: 1.8rem;
}

.usage-pagination-shell :deep(.page-size-select .select-trigger:hover),
.usage-pagination-shell :deep(.page-size-select .select-trigger:focus-visible) {
  border-color: rgba(167, 58, 42, 0.24);
  background: rgba(250, 247, 239, 0.36);
  box-shadow: 0 0 0 2px rgba(167, 58, 42, 0.06);
}

.usage-pagination-shell :deep(.btn-ghost),
.usage-pagination-shell :deep(nav button) {
  border-color: rgba(216, 205, 185, 0.72);
  background: rgba(255, 252, 246, 0.5);
  color: #59645a;
}

.usage-pagination-shell :deep(.btn-ghost:hover),
.usage-pagination-shell :deep(nav button:hover) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  color: #8f3024;
}

.usage-pagination-shell :deep(nav button[aria-current='page']) {
  border-color: rgba(167, 58, 42, 0.3);
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
}

.usage-filters :deep(.input),
.usage-filters :deep(.select-trigger) {
  min-height: 2.7rem;
  border-color: rgba(216, 205, 185, 0.58);
  border-radius: 12px;
  background: rgba(255, 252, 246, 0.34);
  color: #38413a;
  box-shadow: none;
}

.usage-filters :deep(.input-label) {
  color: #59645a;
  font-size: 0.78rem;
  font-weight: 650;
}

.usage-filter-field-key :deep(.select-trigger) {
  min-height: 2.8rem;
  border-color: rgba(216, 205, 185, 0.58);
  border-radius: 12px;
  background: rgba(255, 252, 246, 0.34);
  color: #38413a;
  box-shadow: none;
}

.usage-filter-field-key :deep(.select-trigger:hover),
.usage-filter-field-key :deep(.select-trigger:focus-visible),
.usage-filter-field-key :deep(.select-trigger-open) {
  border-color: rgba(167, 58, 42, 0.24);
  background: rgba(255, 252, 246, 0.5);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.05);
}

.usage-filter-field-key :deep(.select-value) {
  color: #1f2320;
  font-weight: 650;
}

.usage-filter-field-key :deep(.select-icon),
.usage-filter-field-key :deep(svg) {
  color: #8d978a;
}

.usage-filters :deep(.select-value) {
  color: #1f2320;
  font-weight: 650;
}

.usage-filters :deep(.select-icon),
.usage-filters :deep(svg) {
  color: #8d978a;
}

.usage-filters :deep(.input:hover),
.usage-filters :deep(.select-trigger:hover),
.usage-filters :deep(.input:focus),
.usage-filters :deep(.select-trigger:focus-visible),
.usage-filters :deep(.select-trigger-open),
.usage-filter-field-range :deep(.date-picker-trigger:hover),
.usage-filter-field-range :deep(.date-picker-trigger:focus-visible),
.usage-filter-field-range :deep(.date-picker-trigger-open) {
  border-color: rgba(167, 58, 42, 0.55);
  background: rgba(255, 252, 246, 0.58);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.14);
  outline: none;
}

.usage-filters :deep(.date-picker-dropdown) {
  width: min(17.2rem, calc(100vw - 2rem));
  min-width: 0;
  max-width: calc(100vw - 2rem);
  right: auto;
  z-index: 100000030;
  overflow: hidden;
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 8px;
  background: rgba(250, 247, 239, 0.98) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34);
}

.usage-filters :deep(.date-picker-presets) {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 0.2rem 0.26rem;
  padding: 0.28rem 0.3rem 0.32rem;
}

.usage-filters :deep(.date-picker-preset) {
  min-height: 1.34rem;
  border-radius: 6px;
  background: transparent;
  color: #59645a;
  padding: 0.12rem 0.24rem;
  font-size: 0.62rem;
  font-weight: 650;
  line-height: 1.15;
}

.usage-filters :deep(.date-picker-preset:hover),
.usage-filters :deep(.date-picker-preset-active) {
  background: rgba(167, 58, 42, 0.12);
  color: #a73a2a;
}

.usage-filters :deep(.date-picker-custom) {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 0.34rem;
  padding: 0.3rem 0.32rem 0.28rem;
}

.usage-filters :deep(.date-picker-label) {
  margin-bottom: 0.12rem;
  font-size: 0.6rem;
}

.usage-filters :deep(.date-picker-separator) {
  display: none;
}

.usage-filters :deep(.date-picker-input) {
  border-color: rgba(216, 205, 185, 0.78);
  background: rgba(255, 252, 246, 0.62);
  color: #1f2320;
  min-height: 1.54rem;
  border-radius: 6px;
  padding: 0.2rem 0.32rem;
  font-size: 0.66rem;
}

.usage-filters :deep(.date-picker-actions) {
  padding: 0 0.32rem 0.32rem;
}

.usage-filters :deep(.date-picker-apply) {
  border-radius: 999px;
  background: #a73a2a;
  color: #f4efe4;
  min-height: 1.42rem;
  padding: 0.2rem 0.58rem;
  font-size: 0.64rem;
  font-weight: 800;
}

:global(.usage-filter-dropdown),
:global(.usage-page-size-dropdown),
:global(.keys-filter-dropdown),
:global(.keys-page-size-dropdown) {
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 12px !important;
  background: rgba(250, 247, 239, 0.96) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34) !important;
}

:global(.usage-filter-dropdown .select-option),
:global(.usage-page-size-dropdown .select-option),
:global(.keys-filter-dropdown .select-option),
:global(.keys-page-size-dropdown .select-option) {
  color: #38413a !important;
}

:global(.usage-filter-dropdown .select-option:hover),
:global(.usage-filter-dropdown .select-option-focused),
:global(.usage-page-size-dropdown .select-option:hover),
:global(.usage-page-size-dropdown .select-option-focused),
:global(.keys-filter-dropdown .select-option:hover),
:global(.keys-filter-dropdown .select-option-focused),
:global(.keys-page-size-dropdown .select-option:hover),
:global(.keys-page-size-dropdown .select-option-focused) {
  background: rgba(167, 58, 42, 0.075) !important;
  color: #a73a2a !important;
}

:global(.usage-filter-dropdown .select-option-selected),
:global(.usage-page-size-dropdown .select-option-selected),
:global(.keys-filter-dropdown .select-option-selected),
:global(.keys-page-size-dropdown .select-option-selected) {
  background: rgba(167, 58, 42, 0.1) !important;
  color: #a73a2a !important;
}

.usage-page :deep(.card) {
  border-color: rgba(216, 205, 185, 0.76);
  border-radius: 6px;
  background: rgba(250, 247, 239, 0.62);
  box-shadow: 0 14px 38px -32px rgba(31, 35, 32, 0.26);
}

.usage-summary-grid :deep(.card) {
  min-height: 6.5rem;
  min-width: 0;
}

.usage-summary-copy {
  display: grid;
  gap: 0.18rem;
  min-width: 0;
  overflow-wrap: anywhere;
}

.usage-summary-detail {
  display: flex;
  flex-wrap: wrap;
  gap: 0.12rem 0.35rem;
  line-height: 1.55;
}

.usage-summary-detail-token {
  display: grid;
  gap: 0.08rem;
}

.usage-summary-total-line {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 0.35rem 0.75rem;
  margin: 0;
}

.usage-summary-token-line {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 0 0.35rem;
  margin: 0;
}

.usage-summary-token-cache-line {
  justify-content: space-between;
}

.usage-summary-token-cache-values {
  display: inline-flex;
  min-width: 0;
  flex-wrap: wrap;
  align-items: center;
  gap: 0 0.35rem;
}

.usage-cache-rate-badge {
  margin-left: auto;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
  font-size: 0.875rem;
  font-weight: 600;
  color: #8b3d2f;
}

.usage-page-en .usage-summary-detail {
  flex-wrap: nowrap;
  align-items: baseline;
  gap: 0.24rem;
  overflow: hidden;
  white-space: nowrap;
}

.usage-page-en .usage-summary-detail span {
  flex: 0 0 auto;
}

.usage-page-en .usage-summary-detail span:nth-child(n + 5) {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

.usage-page-en .usage-summary-detail-token {
  flex-wrap: initial;
  align-items: initial;
  gap: 0.08rem;
  overflow: visible;
  white-space: normal;
}

.usage-page-en .usage-summary-detail-token span {
  flex: initial;
}

.usage-page-en .usage-summary-detail-token span:nth-child(n + 5) {
  overflow: visible;
  text-overflow: clip;
}

.usage-summary-card-primary {
  border-color: rgba(167, 58, 42, 0.24) !important;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 38%),
    rgba(250, 247, 239, 0.72) !important;
}

.usage-summary-card-primary :deep(.text-green-600),
.usage-summary-card-primary :deep(p.text-green-600) {
  color: #2f6f5e !important;
}

.usage-cell-stack {
  display: grid;
  gap: 0.14rem;
  min-width: 0;
}

.usage-cell-stack-compact {
  gap: 0.08rem;
}

.usage-cell-stack strong {
  overflow: hidden;
  color: #1f2320;
  font-size: 0.8rem;
  font-weight: 650;
  line-height: 1.4;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.usage-cell-stack span {
  overflow: hidden;
  color: #7b6a53;
  font-size: 0.68rem;
  line-height: 1.4;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.usage-cell-break strong,
.usage-cell-break span {
  overflow: visible;
  text-overflow: unset;
  white-space: normal;
  word-break: break-all;
}

.usage-page-en .usage-cell-break strong,
.usage-page-en .usage-cell-break span {
  word-break: normal;
  overflow-wrap: anywhere;
}

.usage-user-agent {
  max-width: 15rem;
}

.usage-user-agent strong,
.usage-user-agent span {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.usage-inline-chip {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 0.24rem 0.62rem;
  font-size: 0.66rem;
  font-weight: 650;
  line-height: 1;
}

.usage-inline-chip-muted {
  border: 1px solid rgba(216, 205, 185, 0.68);
  background: rgba(255, 252, 246, 0.56);
  color: #667066;
}

.usage-token-cell {
  display: grid;
  gap: 0.38rem;
  min-width: 0;
}

.usage-page-en .usage-token-cell {
  gap: 0.3rem;
}

.usage-token-cell-image {
  grid-template-columns: auto minmax(0, 1fr);
  align-items: center;
}

.usage-token-ledger {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.usage-token-ledger-compact {
  gap: 0.62rem;
}

.usage-token-stat {
  display: inline-flex;
  align-items: baseline;
  gap: 0.34rem;
  min-width: 0;
}

.usage-token-stat span {
  color: #7b6a53;
  font-size: 0.62rem;
  font-weight: 650;
  letter-spacing: 0.08em;
  line-height: 1;
}

.usage-token-stat strong {
  color: #1f2320;
  font-size: 0.8rem;
  font-weight: 650;
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.usage-token-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.28rem;
}

.usage-page-en .usage-token-meta {
  gap: 0.22rem;
}

.usage-meta-chip {
  display: inline-flex;
  align-items: center;
  border: 1px solid rgba(216, 205, 185, 0.68);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.6);
  padding: 0.16rem 0.42rem;
  color: #667066;
  font-size: 0.62rem;
  font-weight: 650;
  line-height: 1.2;
}

.usage-cost-cell {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
}

.usage-cost-cell strong {
  color: #2f6f5e;
  font-variant-numeric: tabular-nums;
}

.usage-cost-cell span {
  font-variant-numeric: tabular-nums;
}

.usage-tab-strip {
  display: flex;
  gap: 0.35rem;
  margin-bottom: 0;
  min-width: 0;
  padding: 0;
  align-self: center;
  justify-self: start;
}

.usage-tab-button {
  position: relative;
  min-height: 2.16rem;
  padding: 0.46rem 0.92rem 0.5rem;
  border: 1px solid transparent;
  border-radius: 999px;
  font-size: 0.82rem;
  font-weight: 650;
  line-height: 1.1;
  white-space: normal;
  transition:
    border-color 160ms ease,
    background-color 160ms ease,
    color 160ms ease,
    box-shadow 160ms ease;
}

.usage-tab-button::after {
  content: '';
  position: absolute;
  left: 0.92rem;
  right: 0.92rem;
  bottom: 0.2rem;
  height: 1px;
  border-radius: 999px;
  background: transparent;
  transition:
    background-color 160ms ease,
    opacity 160ms ease;
}

.usage-page :deep(.table-wrapper) {
  overflow: auto;
  width: 100%;
  min-height: 0;
  max-height: min(42rem, calc(100vh - 15.5rem));
  flex: 0 0 auto !important;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  isolation: auto;
  scrollbar-color: rgba(167, 58, 42, 0.48) rgba(237, 229, 212, 0.18) !important;
}

.usage-page :deep(.table-wrapper::-webkit-scrollbar-track) {
  background-color: rgba(237, 229, 212, 0.18) !important;
}

.usage-page :deep(.table-wrapper::-webkit-scrollbar-thumb) {
  background-color: rgba(167, 58, 42, 0.52) !important;
}

.usage-page :deep(.table-wrapper::-webkit-scrollbar-thumb:hover) {
  background-color: rgba(143, 48, 36, 0.78) !important;
}

.usage-page :deep(table) {
  width: 100%;
  min-width: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  background: transparent;
  table-layout: auto;
}

.usage-page-en :deep(table) {
  min-width: 104rem;
  table-layout: fixed;
}

.usage-page-en :deep(.usage-col-api-key) {
  width: 5.4rem;
  min-width: 5.4rem;
  max-width: 5.4rem;
}

.usage-page-en :deep(.usage-col-model) {
  width: 8.8rem;
  min-width: 8.8rem;
  max-width: 8.8rem;
}

.usage-page-en :deep(.usage-col-reasoning) {
  width: 8.2rem;
  min-width: 8.2rem;
  max-width: 8.2rem;
}

.usage-page-en :deep(.usage-col-endpoint) {
  width: 9.2rem;
  min-width: 9.2rem;
  max-width: 9.2rem;
}

.usage-page-en :deep(.usage-col-type),
.usage-page-en :deep(.usage-col-billing-mode) {
  width: 6.2rem;
  min-width: 6.2rem;
  max-width: 6.2rem;
}

.usage-page-en :deep(.usage-col-tokens) {
  width: 10.6rem;
  min-width: 10.6rem;
  max-width: 10.6rem;
}

.usage-page-en :deep(.usage-col-first-token),
.usage-page-en :deep(.usage-col-duration) {
  width: 7.4rem;
  min-width: 7.4rem;
  max-width: 7.4rem;
}

.usage-page :deep(.usage-col-user-agent) {
  width: 15rem;
  min-width: 15rem;
  max-width: 15rem;
}

.usage-page :deep(.usage-col-created-at) {
  min-width: 8.5rem;
}

.usage-page :deep(.usage-col-cost) {
  min-width: 9rem;
}

.usage-page :deep(thead) {
  background: transparent;
}

.usage-page :deep(th),
.usage-page :deep(.sticky-header-cell) {
  border-bottom: 1px solid rgba(167, 58, 42, 0.2);
  background: rgba(237, 229, 212, 0.58) !important;
  color: #7b6a53;
  font-size: 0.65rem;
  font-weight: 650;
  letter-spacing: 0.12em;
  padding-top: 0.62rem;
  padding-bottom: 0.62rem;
}

.usage-page :deep(td) {
  border-bottom: 1px solid rgba(216, 205, 185, 0.46);
  background: rgba(255, 252, 246, 0.24);
  color: #38413a;
  font-size: 0.78rem;
  line-height: 1.7;
  padding-top: 0.68rem;
  padding-bottom: 0.68rem;
  vertical-align: middle;
  white-space: normal;
}

.usage-page :deep(tbody) {
  background: transparent;
}

.usage-page :deep(tbody .sticky-col) {
  background: rgba(250, 247, 239, 0.92) !important;
}

.usage-page :deep(tbody tr:hover) {
  background: transparent;
}

.usage-page :deep(tbody tr:hover td),
.usage-page :deep(tbody tr:hover .sticky-col) {
  background: rgba(167, 58, 42, 0.04) !important;
}

.usage-empty-state {
  display: grid;
  place-items: center;
  gap: 0.36rem;
  min-height: 12.5rem;
  padding: 1.2rem 1.45rem;
  color: #59645a;
  text-align: center;
}

.usage-empty-state span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.usage-empty-state strong {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.08rem;
  font-weight: 600;
}

.usage-empty-state p {
  max-width: 25rem;
  font-size: 0.82rem;
  line-height: 1.8;
}

.usage-empty-state button {
  margin-top: 0.35rem;
  border: 1px solid rgba(216, 205, 185, 0.78);
  border-radius: 999px;
  background: rgba(250, 247, 239, 0.62);
  padding: 0.42rem 0.86rem;
  color: #59645a;
  font-size: 0.78rem;
  font-weight: 700;
}

.usage-empty-state button:hover,
.usage-empty-state button:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.075);
  color: #a73a2a;
  outline: none;
}

.usage-page :deep(.badge),
.usage-page :deep(.rounded-full) {
  box-shadow: none;
}

.usage-page :deep(.tab) {
  border-radius: 999px;
  color: #59645a;
  letter-spacing: 0.03em;
}

.usage-page :deep(.tab:hover) {
  border-color: rgba(216, 205, 185, 0.66);
  background: rgba(255, 252, 246, 0.48);
  color: #7b6a53;
}

.usage-page :deep(.tab-active) {
  border-color: rgba(167, 58, 42, 0.2);
  background: rgba(167, 58, 42, 0.075);
  color: #a73a2a;
  box-shadow: inset 0 1px 0 rgba(255, 252, 246, 0.72);
}

.usage-page :deep(.tab-active)::after {
  background: rgba(167, 58, 42, 0.72);
}

.usage-page :deep(.rounded),
.usage-page :deep(.rounded-lg),
.usage-page :deep(.rounded-xl),
.usage-page :deep(.rounded-2xl) {
  border-radius: 6px;
}

.usage-page :deep(.bg-blue-50),
.usage-page :deep(.bg-blue-100),
.usage-page :deep(.bg-sky-50),
.usage-page :deep(.bg-green-50),
.usage-page :deep(.bg-green-100),
.usage-page :deep(.bg-amber-100),
.usage-page :deep(.bg-purple-100),
.usage-page :deep(.bg-violet-100),
.usage-page :deep(.bg-indigo-100),
.usage-page :deep(.bg-gray-50),
.usage-page :deep(.bg-gray-100) {
  background-color: rgba(216, 205, 185, 0.46);
}

.usage-page :deep(.text-blue-500),
.usage-page :deep(.text-blue-600),
.usage-page :deep(.text-sky-500),
.usage-page :deep(.text-sky-600),
.usage-page :deep(.text-green-500),
.usage-page :deep(.text-green-600),
.usage-page :deep(.text-amber-500),
.usage-page :deep(.text-amber-600),
.usage-page :deep(.text-purple-600),
.usage-page :deep(.text-violet-500),
.usage-page :deep(.text-indigo-500),
.usage-page :deep(.text-primary-600) {
  color: #59645a;
}

.usage-page :deep(.text-green-600.font-medium),
.usage-page :deep(p.text-green-600) {
  color: #2f6f5e;
}

.usage-page :deep(.btn-primary) {
  border-color: #a73a2a;
  background-color: #a73a2a;
  color: #f4efe4;
}

.usage-page :deep(.btn-primary:hover) {
  border-color: #8f3024;
  background-color: #8f3024;
}

.usage-page :deep(.btn-primary:disabled) {
  cursor: not-allowed;
  border-color: rgba(216, 205, 185, 0.86);
  background: rgba(250, 247, 239, 0.78);
  color: #38413a;
  opacity: 1;
}

.usage-page :deep(.btn-secondary) {
  border-color: rgba(216, 205, 185, 0.9);
  background: rgba(255, 255, 255, 0.38);
  color: #38413a;
}

.usage-page :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(255, 255, 255, 0.66);
}

.usage-pagination-shell :deep(p),
.usage-pagination-shell :deep(span) {
  color: #667066;
  font-size: 0.82rem;
}

.dark .usage-page {
  color: #f4efe4;
}

.dark .usage-ledger-item {
  border-left-color: rgba(48, 52, 43, 0.95);
}

.dark .usage-ledger-item span {
  color: #879186;
}

.dark .usage-ledger-item strong {
  color: #f4efe4;
}

.dark .usage-page :deep(.card) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .usage-summary-card-primary {
  border-color: rgba(167, 58, 42, 0.36) !important;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 38%),
    rgba(24, 26, 21, 0.78) !important;
}

.dark .usage-filters :deep(.input),
.dark .usage-filters :deep(.select-trigger) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
}

.dark .usage-filter-field-range :deep(.date-picker-trigger) {
  border-color: rgba(48, 52, 43, 0.72);
  background: rgba(17, 19, 15, 0.26);
  color: #d8cdb9;
}

.dark .usage-filters :deep(.input:hover),
.dark .usage-filters :deep(.select-trigger:hover),
.dark .usage-filters :deep(.input:focus),
.dark .usage-filters :deep(.select-trigger:focus-visible),
.dark .usage-filters :deep(.select-trigger-open),
.dark .usage-filters :deep(.date-picker-trigger:hover),
.dark .usage-filters :deep(.date-picker-trigger:focus-visible),
.dark .usage-filters :deep(.date-picker-trigger-open) {
  border-color: rgba(167, 58, 42, 0.55);
  background: rgba(24, 26, 21, 0.92);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.16);
  outline: none;
}

.dark .usage-filters :deep(.select-value) {
  color: #f4efe4;
}

.dark .usage-filter-field-key :deep(.select-trigger) {
  border-color: rgba(48, 52, 43, 0.72);
  background: rgba(17, 19, 15, 0.26);
  color: #d8cdb9;
}

.dark .usage-filter-field-key :deep(.select-trigger:hover),
.dark .usage-filter-field-key :deep(.select-trigger:focus-visible),
.dark .usage-filter-field-key :deep(.select-trigger-open) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  box-shadow: 0 0 0 2px rgba(167, 58, 42, 0.08);
}

.dark .usage-filter-field-key :deep(.select-value) {
  color: #f4efe4;
}

.dark .usage-filter-field-range :deep(.date-picker-value) {
  color: #f4efe4;
}

.dark .usage-filter-field-key :deep(.select-icon),
.dark .usage-filter-field-key :deep(svg) {
  color: #879186;
}

.dark .usage-filter-field-range :deep(.date-picker-icon),
.dark .usage-filter-field-range :deep(.date-picker-chevron),
.dark .usage-filter-field-range :deep(svg) {
  color: #879186;
}

.dark .usage-filters :deep(.select-icon),
.dark .usage-filters :deep(svg) {
  color: #879186;
}

.dark .usage-filters :deep(.date-picker-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(17, 19, 15, 0.98) !important;
  box-shadow: 0 18px 42px -34px rgba(0, 0, 0, 0.62);
}

.dark .usage-filters :deep(.date-picker-preset) {
  color: #d8cdb9;
}

.dark .usage-filters :deep(.date-picker-preset:hover),
.dark .usage-filters :deep(.date-picker-preset-active) {
  background: rgba(167, 58, 42, 0.18);
  color: #f0b4a8;
}

.dark .usage-filters :deep(.date-picker-divider) {
  border-color: rgba(48, 52, 43, 0.78);
}

.dark .usage-filters :deep(.date-picker-label) {
  color: #b9aa91;
}

.dark .usage-filters :deep(.date-picker-input) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.58);
  color: #f4efe4;
}

.dark .usage-filters :deep(.date-picker-input:hover),
.dark .usage-filters :deep(.date-picker-input:focus) {
  border-color: rgba(167, 58, 42, 0.45);
  background: rgba(24, 26, 21, 0.78);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.12);
}

:global(.dark .usage-filter-dropdown),
:global(.dark .usage-page-size-dropdown),
:global(.dark .keys-filter-dropdown),
:global(.dark .keys-page-size-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.98) !important;
}

:global(.dark .usage-filter-dropdown .select-option),
:global(.dark .usage-page-size-dropdown .select-option),
:global(.dark .keys-filter-dropdown .select-option),
:global(.dark .keys-page-size-dropdown .select-option) {
  color: #d8cdb9 !important;
}

.dark .usage-page :deep(.table-wrapper),
.dark .usage-page :deep(table) {
  border-color: transparent;
  background: transparent;
}

.dark .usage-page :deep(.table-wrapper) {
  scrollbar-color: rgba(184, 156, 116, 0.4) rgba(24, 26, 21, 0.32) !important;
}

.dark .usage-page :deep(.table-wrapper::-webkit-scrollbar-track) {
  background-color: rgba(24, 26, 21, 0.32) !important;
}

.dark .usage-page :deep(.table-wrapper::-webkit-scrollbar-thumb) {
  background-color: rgba(184, 156, 116, 0.44) !important;
}

.dark .usage-page :deep(.table-wrapper::-webkit-scrollbar-thumb:hover) {
  background-color: rgba(216, 205, 185, 0.62) !important;
}

.dark .usage-page :deep(.table-header),
.dark .usage-page :deep(thead),
.dark .usage-page :deep(.sticky-header-cell) {
  background: transparent !important;
}

.dark .usage-page :deep(th),
.dark .usage-page :deep(.sticky-header-cell) {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(24, 26, 21, 0.9) !important;
  color: #b9aa91;
}

.dark .usage-page :deep(td) {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.34);
  color: #d8cdb9;
}

.dark .usage-page :deep(tbody .sticky-col) {
  background: rgba(17, 19, 15, 0.92) !important;
}

.dark .usage-page :deep(tbody tr:hover td),
.dark .usage-page :deep(tbody tr:hover .sticky-col) {
  background: rgba(167, 58, 42, 0.11) !important;
}

.dark .usage-page :deep(.tab:hover) {
  border-color: rgba(216, 205, 185, 0.16);
  background: rgba(216, 205, 185, 0.045);
  color: #d8cdb9;
}

.dark .usage-empty-state strong {
  color: #f4efe4;
}

.dark .usage-empty-state button {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.58);
  color: #d8cdb9;
}

.dark .usage-page :deep(.tab) {
  color: #879186;
}

.dark .usage-page :deep(.tab-active) {
  border-color: rgba(184, 156, 116, 0.22);
  background: rgba(184, 156, 116, 0.09);
  color: #d8cdb9;
  box-shadow: inset 0 1px 0 rgba(244, 239, 228, 0.06);
}

.dark .usage-page :deep(.tab-active)::after {
  background: rgba(184, 156, 116, 0.72);
}

.dark .usage-page :deep(.btn-secondary) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.42);
  color: #f4efe4;
}

.dark .usage-page :deep(.btn-primary:disabled) {
  border-color: rgba(216, 205, 185, 0.24);
  background: rgba(216, 205, 185, 0.16);
  color: #f4efe4;
}

.dark .usage-cell-stack strong,
.dark .usage-token-stat strong {
  color: #f4efe4;
}

.dark .usage-cell-stack span,
.dark .usage-token-stat span {
  color: #b9aa91;
}

.dark .usage-meta-chip {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.5);
  color: #b9aa91;
}

.dark .usage-inline-chip-muted {
  border-color: rgba(48, 52, 43, 0.82);
  background: rgba(17, 19, 15, 0.5);
  color: #b9aa91;
}

.dark .usage-pagination-shell :deep(.page-size-select .select-trigger) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.58);
  color: #d8cdb9;
}

.dark .usage-pagination-shell :deep(.btn-ghost),
.dark .usage-pagination-shell :deep(nav button) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.42);
  color: #d8cdb9;
}

.dark .usage-pagination-shell :deep(.page-size-select .select-trigger:hover),
.dark .usage-pagination-shell :deep(.page-size-select .select-trigger:focus-visible),
.dark .usage-pagination-shell :deep(.btn-ghost:hover),
.dark .usage-pagination-shell :deep(nav button:hover) {
  border-color: rgba(216, 205, 185, 0.34);
  background: rgba(216, 205, 185, 0.06);
  color: #f4efe4;
}

.dark .usage-pagination-shell :deep(nav button[aria-current='page']) {
  border-color: rgba(184, 156, 116, 0.34);
  background: rgba(184, 156, 116, 0.1);
  color: #f4efe4;
}

@media (max-width: 1180px) {
  .usage-toolbar-grid {
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
    align-items: stretch;
  }

  .usage-tab-strip,
  .usage-toolbar-actions {
    justify-self: stretch;
  }

  .usage-toolbar-actions {
    justify-content: flex-start;
  }
}

@media (max-width: 768px) {
  .usage-ledger {
    grid-template-columns: 1fr;
  }

  .usage-ledger-item {
    border-left: 0;
    border-top: 1px solid rgba(216, 205, 185, 0.76);
    padding-left: 0;
    padding-top: 0.85rem;
  }

  .usage-summary-grid {
    grid-template-columns: 1fr;
  }

  .usage-toolbar-grid {
    grid-template-columns: 1fr;
    justify-content: stretch;
    align-items: stretch;
  }

  .usage-tab-strip {
    overflow-x: auto;
    padding-bottom: 0.08rem;
    justify-self: stretch;
  }

  .usage-filter-field,
  .usage-filter-field-key,
  .usage-filter-field-range,
  .usage-toolbar-actions {
    width: 100%;
    min-width: 0;
  }

  .usage-toolbar-actions {
    margin-left: 0;
    justify-content: stretch;
    justify-self: stretch;
  }

  .usage-toolbar-actions .btn {
    flex: 1 1 8rem;
    min-width: 0;
  }

  .usage-filters :deep(.date-picker-dropdown) {
    width: min(17.2rem, calc(100vw - 2rem));
  }

  .usage-filters :deep(.date-picker-custom) {
    grid-template-columns: 1fr;
  }

  .usage-filters :deep(.date-picker-separator) {
    display: none;
  }

  .usage-pagination-shell {
    margin-top: 0.6rem;
  }

  .usage-user-agent {
    max-width: none;
  }

}
</style>
