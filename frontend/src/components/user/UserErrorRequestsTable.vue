<template>
  <div class="usage-error-table flex min-h-0 flex-1 flex-col">
    <div class="usage-error-toolbar flex-shrink-0">
      <div class="usage-error-control-panel">
        <div class="usage-error-chips">
          <button
            v-for="chip in quickFilterChips"
            :key="chip.value"
            type="button"
            class="usage-error-chip rounded-full border px-3 py-1 text-xs font-medium transition-colors"
            :class="localCategory === chip.value
              ? 'usage-error-chip-active border-amber-300 bg-amber-50 text-amber-900 dark:border-amber-800 dark:bg-amber-950/20 dark:text-amber-100'
              : 'border-gray-200 bg-white text-gray-600 hover:border-amber-200 hover:text-amber-900 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-300 dark:hover:border-amber-900/50 dark:hover:text-amber-100'"
            @click="applyQuickFilter(chip.value)"
          >
            {{ chip.label }}
          </button>
        </div>

        <div class="usage-error-filters">
        <div class="usage-error-filter usage-error-filter-model">
          <label class="input-label">{{ t('usage.errors.model') }}</label>
          <Select
            v-model="localModel"
            :options="modelOptions"
            searchable
            creatable
            clearable
            :placeholder="t('usage.errors.modelPlaceholder')"
            dropdown-class="usage-error-dropdown"
            @change="apply"
          />
        </div>
        <div class="usage-error-filter usage-error-filter-key">
          <label class="input-label">{{ t('usage.errors.keyName') }}</label>
          <Select
            v-model="localApiKeyId"
            :options="keyOptions"
            :placeholder="t('usage.errors.allKeys')"
            dropdown-class="usage-error-dropdown"
            @change="apply"
          />
        </div>
        <div class="usage-error-filter usage-error-filter-category">
          <label class="input-label">{{ t('usage.errors.category') }}</label>
          <Select
            v-model="localCategory"
            :options="categoryOptions"
            :placeholder="t('usage.errors.allCategories')"
            dropdown-class="usage-error-dropdown"
            @change="apply"
          />
        </div>
        <button class="btn btn-primary usage-error-search" @click="apply">
          <Icon name="search" size="sm" />
          {{ t('common.search') }}
        </button>
        </div>
      </div>
    </div>

    <div class="usage-error-scroll min-h-0 flex-1 overflow-auto">
      <table class="usage-error-grid min-w-full text-sm">
        <colgroup>
          <col class="usage-error-col-model" />
          <col class="usage-error-col-key" />
          <col class="usage-error-col-endpoint" />
          <col class="usage-error-col-status" />
          <col class="usage-error-col-category" />
          <col class="usage-error-col-hint" />
          <col class="usage-error-col-message" />
          <col class="usage-error-col-platform" />
          <col class="usage-error-col-time" />
        </colgroup>
        <thead>
          <tr>
            <th class="text-left">{{ t('usage.errors.model') }}</th>
            <th class="text-left">{{ t('usage.errors.keyName') }}</th>
            <th class="text-left">{{ t('usage.errors.endpoint') }}</th>
            <th class="text-left">{{ t('usage.errors.status') }}</th>
            <th class="text-left">{{ t('usage.errors.category') }}</th>
            <th class="text-left">{{ t('usage.errors.hint') }}</th>
            <th class="text-left">{{ t('usage.errors.message') }}</th>
            <th class="text-left">{{ t('usage.errors.platform') }}</th>
            <th class="text-left">{{ t('usage.errors.time') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(row, i) in rows"
            :key="i"
            class="usage-error-row border-t border-gray-100 dark:border-dark-700 cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-800"
            @click="openDetail(row.id)"
          >
            <td>{{ row.model || '-' }}</td>
            <td>
              <span>{{ row.key_name || '-' }}</span>
              <span
                v-if="row.key_deleted"
                class="usage-error-deleted ml-1 inline-flex items-center rounded px-1 py-px text-[10px] font-medium leading-tight bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400"
              >{{ t('usage.errors.keyDeleted') }}</span>
            </td>
            <td>{{ row.inbound_endpoint || '-' }}</td>
            <td><span class="badge" :class="statusClass(row.status_code)">{{ row.status_code || '-' }}</span></td>
            <td>{{ t('usage.errors.categories.' + row.category) }}</td>
            <td>
              <span class="usage-error-hint inline-flex max-w-[220px] rounded-full border border-amber-200/80 bg-amber-50/70 px-2.5 py-1 text-xs leading-5 text-amber-900 dark:border-amber-900/50 dark:bg-amber-950/10 dark:text-amber-100">
                {{ getHintLabel(row) }}
              </span>
            </td>
            <td class="truncate" :title="getDisplayMessage(row)">{{ getDisplayMessage(row) || '-' }}</td>
            <td>{{ row.platform || '-' }}</td>
            <td>{{ formatDateTime(row.created_at) }}</td>
          </tr>
          <tr v-if="!loading && rows.length === 0">
            <td colspan="9" class="usage-error-empty px-4 py-8 text-center text-gray-400">{{ t('usage.errors.empty') }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="usage-error-pagination flex-shrink-0">
      <Pagination :page="page" :page-size="pageSize" :total="total" page-size-dropdown-class="usage-error-dropdown"
        @update:page="$emit('update:page', $event)"
        @update:pageSize="$emit('update:pageSize', $event)" />
    </div>

    <UserErrorDetailModal v-model:show="showDetail" :error-id="selectedId" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import UserErrorDetailModal from '@/components/user/UserErrorDetailModal.vue'
import { formatDateTime } from '@/utils/format'
import type { UserErrorRequest, ApiKey } from '@/types'

const props = defineProps<{
  rows: UserErrorRequest[]
  total: number
  loading: boolean
  page: number
  pageSize: number
  apiKeys?: ApiKey[]
  modelFilter?: string
  categoryFilter?: string
  apiKeyIdFilter?: number | null
}>()

const emit = defineEmits<{
  (e: 'update:page', v: number): void
  (e: 'update:pageSize', v: number): void
  (e: 'filter', v: { model: string; category: string; api_key_id: number | null }): void
}>()

const { t } = useI18n()
// string | null:clearable 清空时 Select 回传 null,apply 中归一为空串
const localModel = ref<string | null>('')
const localCategory = ref<string>('')
const localApiKeyId = ref<number | null>(null)

const categoryCodes = ['auth', 'rate_limit', 'quota', 'invalid_request', 'service_unavailable', 'upstream', 'internal', 'cyber']
const quickFilterCategories = ['rate_limit', 'quota', 'service_unavailable', 'upstream'] as const

const categoryOptions = computed(() => [
  { value: '', label: t('usage.errors.allCategories') },
  ...categoryCodes.map((c) => ({ value: c, label: t('usage.errors.categories.' + c) })),
])

const quickFilterChips = computed(() => [
  { value: '', label: t('usage.errors.quickFilters.all') },
  ...quickFilterCategories.map((c) => ({ value: c, label: t('usage.errors.quickFilters.' + c) })),
])

// 首项 value: null 表示不按 key 过滤；其余项取自父组件传入的 apiKeys 候选列表。
const keyOptions = computed(() => [
  { value: null, label: t('usage.errors.allKeys') },
  ...(props.apiKeys ?? []).map((k) => ({ value: k.id, label: k.name })),
])

// 模型候选取自当前已加载错误中出现过的模型；creatable 允许输入任意片段做后端模糊。
const modelOptions = computed(() => {
  const seen = new Set<string>()
  const opts: { value: string; label: string }[] = []
  for (const r of props.rows) {
    if (r.model && !seen.has(r.model)) {
      seen.add(r.model)
      opts.push({ value: r.model, label: r.model })
    }
  }
  return opts
})

const showDetail = ref(false)
const selectedId = ref<number | null>(null)

watch(
  () => [props.modelFilter, props.categoryFilter, props.apiKeyIdFilter] as const,
  ([model, category, apiKeyId]) => {
    localModel.value = model || ''
    localCategory.value = category || ''
    localApiKeyId.value = apiKeyId ?? null
  },
  { immediate: true }
)

function openDetail(id: number) {
  selectedId.value = id
  showDetail.value = true
}

function apply() {
  emit('filter', {
    model: (localModel.value ?? '').trim(),
    category: localCategory.value || '',
    api_key_id: localApiKeyId.value,
  })
}

function applyQuickFilter(category: string) {
  localCategory.value = category
  apply()
}

function statusClass(code: number) {
  if (code >= 500) return 'badge-danger'
  if (code === 429) return 'badge-warning'
  return 'badge-gray'
}

function getHintLabel(row: UserErrorRequest) {
  const key = 'usage.errors.hints.' + row.category
  return t(key)
}

function getDisplayMessage(row: UserErrorRequest) {
  return summarizeErrorMessage(row.message || '', row.status_code)
}

function summarizeErrorMessage(message: string, statusCode?: number) {
  const trimmed = message.trim()
  if (!trimmed) return ''
  if (!looksLikeHTML(trimmed)) return truncateDisplayText(normalizeDisplayWhitespace(trimmed))

  const parts: string[] = []
  const title = extractHTMLTitle(trimmed)
  if (title) parts.push(title)
  if (statusCode && statusCode > 0) parts.push(`HTTP ${statusCode}`)
  const text = normalizeDisplayWhitespace(decodeBasicHTMLEntities(trimmed.replace(/<[^>]*>/gs, ' ')))
  if (text) parts.push(text)
  return truncateDisplayText(uniqueSummaryParts(parts).join(' · ') || t('usage.errors.detail.htmlSummaryFallback'))
}

function looksLikeHTML(value: string) {
  const lower = value.trim().toLowerCase()
  return lower.startsWith('<!doctype html') ||
    lower.startsWith('<html') ||
    lower.includes('<head') ||
    lower.includes('<body') ||
    lower.includes('</html>')
}

function extractHTMLTitle(value: string) {
  const match = value.match(/<title[^>]*>([\s\S]*?)<\/title>/i)
  return match ? normalizeDisplayWhitespace(decodeBasicHTMLEntities(match[1])) : ''
}

function decodeBasicHTMLEntities(value: string) {
  return value
    .replace(/&nbsp;/gi, ' ')
    .replace(/&amp;/gi, '&')
    .replace(/&lt;/gi, '<')
    .replace(/&gt;/gi, '>')
    .replace(/&quot;/gi, '"')
    .replace(/&#39;/gi, "'")
}

function normalizeDisplayWhitespace(value: string) {
  return value.trim().replace(/\s+/g, ' ')
}

function uniqueSummaryParts(parts: string[]) {
  const seen = new Set<string>()
  return parts.filter((part) => {
    const normalized = part.toLowerCase()
    if (!part || seen.has(normalized)) return false
    seen.add(normalized)
    return true
  })
}

function truncateDisplayText(value: string) {
  const chars = Array.from(value.trim())
  return chars.length > 240 ? `${chars.slice(0, 240).join('')}...` : chars.join('')
}
</script>

<style scoped>
.usage-error-table {
  color: #38413a;
}

.usage-error-toolbar {
  padding: 0.85rem 1.35rem 0.95rem;
  border-bottom: 1px solid rgba(216, 205, 185, 0.52);
  background: linear-gradient(180deg, rgba(250, 247, 239, 0.66), rgba(250, 247, 239, 0.36));
}

.usage-error-control-panel {
  display: grid;
  gap: 0.78rem;
  max-width: 55rem;
}

.usage-error-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.usage-error-chip {
  border-color: rgba(216, 205, 185, 0.78);
  background: rgba(255, 252, 246, 0.68);
  color: #59645a;
}

.usage-error-chip:hover,
.usage-error-chip-active {
  border-color: rgba(167, 58, 42, 0.34) !important;
  background: rgba(167, 58, 42, 0.08) !important;
  color: #8f3024 !important;
}

.usage-error-filters {
  display: grid;
  grid-template-columns: minmax(10rem, 13rem) minmax(10rem, 11.5rem) minmax(9rem, 10rem) auto;
  align-items: end;
  gap: 0.72rem;
}

.usage-error-filter {
  min-width: 0;
}

.usage-error-filters .input-label {
  display: block;
  margin-bottom: 0.35rem;
  color: #59645a;
  font-size: 0.78rem;
  font-weight: 650;
}

.usage-error-filters :deep(.select-trigger) {
  min-height: 2.62rem;
  border-color: rgba(216, 205, 185, 0.76);
  border-radius: 10px;
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.42), rgba(250, 247, 239, 0.7));
  color: #38413a;
  box-shadow: inset 0 1px 0 rgba(255, 252, 246, 0.72);
}

.usage-error-filters :deep(.select-trigger:hover),
.usage-error-filters :deep(.select-trigger:focus-visible),
.usage-error-filters :deep(.select-trigger-open) {
  border-color: rgba(167, 58, 42, 0.34);
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(250, 247, 239, 0.84));
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.08), inset 0 1px 0 rgba(255, 252, 246, 0.78);
}

.usage-error-filters :deep(.select-value),
.usage-error-filters :deep(.select-search-input) {
  color: #1f2320;
}

.usage-error-filters :deep(.select-icon),
.usage-error-filters :deep(.select-clear) {
  color: #8d978a;
}

.usage-error-search {
  min-height: 2.62rem;
  align-self: end;
  min-width: 6rem;
}

.usage-error-scroll {
  border-top: 1px solid rgba(255, 252, 246, 0.42);
}

.usage-error-grid {
  min-width: 72rem;
  background: rgba(250, 247, 239, 0.52);
  table-layout: fixed;
}

.usage-error-col-model { width: 9.5rem; }
.usage-error-col-key { width: 17.5rem; }
.usage-error-col-endpoint { width: 9.5rem; }
.usage-error-col-status { width: 6.5rem; }
.usage-error-col-category { width: 8rem; }
.usage-error-col-hint { width: 14rem; }
.usage-error-col-message { width: 22rem; }
.usage-error-col-platform { width: 8rem; }
.usage-error-col-time { width: 10.5rem; }

.usage-error-grid thead tr {
  background: rgba(237, 229, 212, 0.72);
}

.usage-error-grid th {
  border-bottom: 1px solid rgba(198, 184, 157, 0.54);
  color: #59645a;
  font-size: 0.76rem;
  font-weight: 650;
  letter-spacing: 0.08em;
  padding-left: 1rem;
  padding-right: 1rem;
  padding-top: 0.95rem;
  padding-bottom: 0.95rem;
  white-space: nowrap;
}

.usage-error-grid td {
  border-bottom: 1px solid rgba(198, 184, 157, 0.32);
  color: #38413a;
  overflow: hidden;
  padding-left: 1rem;
  padding-right: 1rem;
  padding-top: 0.95rem;
  padding-bottom: 0.95rem;
  text-overflow: ellipsis;
  vertical-align: top;
  white-space: normal;
  word-break: break-word;
}

.usage-error-grid td:nth-child(3),
.usage-error-grid td:nth-child(4),
.usage-error-grid td:nth-child(5),
.usage-error-grid td:nth-child(8),
.usage-error-grid td:nth-child(9) {
  white-space: nowrap;
}

.usage-error-row {
  transition: background-color 160ms ease;
}

.usage-error-row:hover {
  background: rgba(167, 58, 42, 0.055) !important;
}

.usage-error-deleted {
  background: rgba(216, 205, 185, 0.46) !important;
  color: #7b6a53 !important;
}

.usage-error-hint {
  border-color: rgba(184, 156, 116, 0.28) !important;
  background: rgba(184, 156, 116, 0.12) !important;
  color: #8a7456 !important;
}

.usage-error-empty {
  color: #7b6a53 !important;
}

.usage-error-pagination {
  border-top: 1px solid rgba(216, 205, 185, 0.3);
}

.usage-error-pagination :deep(> div) {
  border-top: 0;
  background: transparent;
}

.usage-error-pagination :deep(p),
.usage-error-pagination :deep(span) {
  color: #667066;
}

.usage-error-pagination :deep(.page-size-select .select-trigger),
.usage-error-pagination :deep(.input) {
  border-color: rgba(216, 205, 185, 0.72);
  background: rgba(250, 247, 239, 0.52);
  color: #38413a;
  box-shadow: inset 0 1px 0 rgba(255, 252, 246, 0.66);
}

.usage-error-pagination :deep(.page-size-select .select-trigger:hover),
.usage-error-pagination :deep(.page-size-select .select-trigger:focus-visible),
.usage-error-pagination :deep(.input:focus),
.usage-error-pagination :deep(.input:focus-visible) {
  border-color: rgba(167, 58, 42, 0.34);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.08);
}

.usage-error-pagination :deep(.btn-ghost),
.usage-error-pagination :deep(nav button) {
  border-color: rgba(216, 205, 185, 0.72);
  background: rgba(255, 252, 246, 0.5);
  color: #59645a;
}

.usage-error-pagination :deep(.btn-ghost:hover),
.usage-error-pagination :deep(nav button:hover) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  color: #8f3024;
}

.usage-error-pagination :deep(nav button[aria-current='page']) {
  border-color: rgba(167, 58, 42, 0.3);
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
}

:global(.usage-error-dropdown) {
  border-color: rgba(216, 205, 185, 0.72) !important;
  border-radius: 12px !important;
  background: rgba(250, 247, 239, 0.98) !important;
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.34) !important;
}

:global(.usage-error-dropdown .select-search) {
  border-bottom-color: rgba(216, 205, 185, 0.56) !important;
}

:global(.usage-error-dropdown .select-search-input),
:global(.usage-error-dropdown .select-option),
:global(.usage-error-dropdown .select-option-label) {
  color: #38413a !important;
}

:global(.usage-error-dropdown .select-option:hover),
:global(.usage-error-dropdown .select-option-focused) {
  background: rgba(167, 58, 42, 0.075) !important;
  color: #a73a2a !important;
}

:global(.usage-error-dropdown .select-option-selected) {
  background: rgba(167, 58, 42, 0.1) !important;
  color: #a73a2a !important;
}

:global(.usage-error-dropdown .select-empty) {
  color: #7b6a53 !important;
}

.dark .usage-error-table {
  color: #d8cdb9;
}

.dark .usage-error-toolbar {
  border-bottom-color: rgba(48, 52, 43, 0.72);
  background: linear-gradient(180deg, rgba(24, 26, 21, 0.78), rgba(17, 19, 15, 0.42));
}

.dark .usage-error-chip {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.34);
  color: #879186;
}

.dark .usage-error-chip:hover,
.dark .usage-error-chip-active {
  border-color: rgba(216, 205, 185, 0.34) !important;
  background: rgba(216, 205, 185, 0.06) !important;
  color: #d8cdb9 !important;
}

.dark .usage-error-filters .input-label {
  color: #879186;
}

.dark .usage-error-filters :deep(.select-trigger) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
  box-shadow: none;
}

.dark .usage-error-filters :deep(.select-trigger:hover),
.dark .usage-error-filters :deep(.select-trigger:focus-visible),
.dark .usage-error-filters :deep(.select-trigger-open) {
  border-color: rgba(216, 205, 185, 0.28);
  background: rgba(24, 26, 21, 0.9);
  box-shadow: 0 0 0 3px rgba(184, 156, 116, 0.08);
}

.dark .usage-error-filters :deep(.select-value),
.dark .usage-error-filters :deep(.select-search-input) {
  color: #f4efe4;
}

.dark .usage-error-filters :deep(.select-icon),
.dark .usage-error-filters :deep(.select-clear) {
  color: #879186;
}

.dark .usage-error-grid {
  background: rgba(24, 26, 21, 0.72);
}

.dark .usage-error-grid thead tr {
  background: rgba(17, 19, 15, 0.62);
}

.dark .usage-error-grid th,
.dark .usage-error-grid td {
  border-color: rgba(48, 52, 43, 0.82);
}

.dark .usage-error-grid th {
  color: #879186;
}

.dark .usage-error-grid td {
  color: #d8cdb9;
}

.dark .usage-error-row:hover {
  background: rgba(216, 205, 185, 0.045) !important;
}

.dark .usage-error-deleted {
  background: rgba(48, 52, 43, 0.92) !important;
  color: #879186 !important;
}

.dark .usage-error-hint {
  border-color: rgba(184, 156, 116, 0.24) !important;
  background: rgba(184, 156, 116, 0.08) !important;
  color: #d8cdb9 !important;
}

.dark .usage-error-empty {
  color: #879186 !important;
}

.dark .usage-error-pagination {
  border-top-color: rgba(48, 52, 43, 0.58);
}

.dark .usage-error-pagination :deep(p),
.dark .usage-error-pagination :deep(span) {
  color: #879186;
}

.dark .usage-error-pagination :deep(.page-size-select .select-trigger),
.dark .usage-error-pagination :deep(.input),
.dark .usage-error-pagination :deep(.btn-ghost),
.dark .usage-error-pagination :deep(nav button) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.42);
  color: #d8cdb9;
}

.dark .usage-error-pagination :deep(.page-size-select .select-trigger:hover),
.dark .usage-error-pagination :deep(.page-size-select .select-trigger:focus-visible),
.dark .usage-error-pagination :deep(.input:focus),
.dark .usage-error-pagination :deep(.input:focus-visible),
.dark .usage-error-pagination :deep(.btn-ghost:hover),
.dark .usage-error-pagination :deep(nav button:hover) {
  border-color: rgba(216, 205, 185, 0.34);
  background: rgba(216, 205, 185, 0.06);
  color: #f4efe4;
}

.dark .usage-error-pagination :deep(nav button[aria-current='page']) {
  border-color: rgba(184, 156, 116, 0.34);
  background: rgba(184, 156, 116, 0.1);
  color: #f4efe4;
}

:global(.dark .usage-error-dropdown) {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.98) !important;
}

:global(.dark .usage-error-dropdown .select-search) {
  border-bottom-color: rgba(48, 52, 43, 0.72) !important;
}

:global(.dark .usage-error-dropdown .select-search-input),
:global(.dark .usage-error-dropdown .select-option),
:global(.dark .usage-error-dropdown .select-option-label) {
  color: #d8cdb9 !important;
}

:global(.dark .usage-error-dropdown .select-option:hover),
:global(.dark .usage-error-dropdown .select-option-focused) {
  background: rgba(216, 205, 185, 0.06) !important;
  color: #f4efe4 !important;
}

:global(.dark .usage-error-dropdown .select-option-selected) {
  background: rgba(184, 156, 116, 0.1) !important;
  color: #f4efe4 !important;
}

:global(.dark .usage-error-dropdown .select-empty) {
  color: #879186 !important;
}

@media (max-width: 768px) {
  .usage-error-toolbar {
    padding: 0.85rem 1rem 1rem;
  }

  .usage-error-control-panel {
    max-width: none;
  }

  .usage-error-filters {
    grid-template-columns: 1fr;
  }

  .usage-error-scroll {
    overflow-x: auto;
  }

  .usage-error-grid th,
  .usage-error-grid td {
    white-space: nowrap;
  }

  .usage-error-search {
    width: 100%;
    justify-content: center;
  }
}
</style>
