<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import EmptyState from '@/components/common/EmptyState.vue'
import type { AnomalyCompareSummary } from '../utils/anomalyCompare'

interface Props {
  summary: AnomalyCompareSummary | null
  loading: boolean
  errorMessage?: string
  currentWindowLabel?: string
  previousWindowLabel?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'openCurrentDetails'): void
}>()
const { t } = useI18n()

function formatPercent(value: number): string {
  const digits = Math.abs(value) >= 10 ? 1 : 2
  return `${value.toFixed(digits)}%`
}

function formatDelta(value: number): string {
  if (Math.abs(value) < 0.05) return t('admin.ops.compare.delta.flat')
  const digits = Math.abs(value) >= 10 ? 1 : 2
  const amount = Math.abs(value).toFixed(digits)
  return value > 0
    ? t('admin.ops.compare.delta.higher', { value: amount })
    : t('admin.ops.compare.delta.lower', { value: amount })
}

function getDeltaTone(value: number): 'up' | 'down' | 'flat' {
  if (Math.abs(value) < 0.05) return 'flat'
  return value > 0 ? 'up' : 'down'
}

const metricCards = computed(() => {
  if (!props.summary) return []

  return [
    {
      key: 'errorRate',
      label: t('admin.ops.compare.metrics.errorRate'),
      current: formatPercent(props.summary.errorRate.currentValue),
      previous: formatPercent(props.summary.errorRate.previousValue),
      delta: formatDelta(props.summary.errorRate.changeValue),
      tone: getDeltaTone(props.summary.errorRate.changeValue)
    },
    {
      key: 'requestErrorShare',
      label: t('admin.ops.compare.metrics.requestErrorShare'),
      current: formatPercent(props.summary.requestErrorShare.currentValue),
      previous: formatPercent(props.summary.requestErrorShare.previousValue),
      delta: formatDelta(props.summary.requestErrorShare.changeValue),
      tone: getDeltaTone(props.summary.requestErrorShare.changeValue)
    },
    {
      key: 'upstreamErrorShare',
      label: t('admin.ops.compare.metrics.upstreamErrorShare'),
      current: formatPercent(props.summary.upstreamErrorShare.currentValue),
      previous: formatPercent(props.summary.upstreamErrorShare.previousValue),
      delta: formatDelta(props.summary.upstreamErrorShare.changeValue),
      tone: getDeltaTone(props.summary.upstreamErrorShare.changeValue)
    },
    {
      key: 'businessLimitedShare',
      label: t('admin.ops.compare.metrics.businessLimitedShare'),
      current: formatPercent(props.summary.businessLimitedShare.currentValue),
      previous: formatPercent(props.summary.businessLimitedShare.previousValue),
      delta: formatDelta(props.summary.businessLimitedShare.changeValue),
      tone: getDeltaTone(props.summary.businessLimitedShare.changeValue)
    }
  ]
})

const shiftCards = computed(() => {
  if (!props.summary) return []

  return [
    {
      key: 'statusShift',
      label: t('admin.ops.compare.shifts.status'),
      from: props.summary.statusShift.previousLabel,
      to: props.summary.statusShift.currentLabel,
      changed: props.summary.statusShift.changed
    },
    {
      key: 'platformShift',
      label: t('admin.ops.compare.shifts.platform'),
      from: props.summary.platformShift.previousLabel,
      to: props.summary.platformShift.currentLabel,
      changed: props.summary.platformShift.changed
    }
  ]
})

const compareWindowLabel = computed(() => {
  if (!props.currentWindowLabel || !props.previousWindowLabel) return ''
  return t('admin.ops.compare.compareWindow', {
    current: props.currentWindowLabel,
    previous: props.previousWindowLabel
  })
})
</script>

<template>
  <section class="ops-anomaly-compare-shell rounded-[28px] p-6">
    <div class="ops-anomaly-compare-header flex flex-col gap-4 pb-4 md:flex-row md:items-start md:justify-between">
      <div class="space-y-1">
        <div class="ops-anomaly-compare-title flex items-center gap-2 text-[13px] font-semibold">
          <span class="ops-anomaly-compare-mark inline-flex h-6 w-6 items-center justify-center rounded-full">异</span>
          {{ t('admin.ops.compare.title') }}
        </div>
        <p class="ops-anomaly-compare-subtitle text-[13px]">{{ t('admin.ops.compare.subtitle') }}</p>
        <p v-if="compareWindowLabel" class="ops-anomaly-compare-window text-xs">{{ compareWindowLabel }}</p>
      </div>
      <button
        type="button"
        class="ops-anomaly-compare-action inline-flex items-center justify-center rounded-full px-4 py-2 text-[13px] font-semibold transition disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="loading || !summary"
        @click="emit('openCurrentDetails')"
      >
        {{ t('admin.ops.compare.currentWindowDetails') }}
      </button>
    </div>

    <div v-if="errorMessage" class="ops-anomaly-compare-error mt-4 rounded-2xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div v-else-if="loading" class="mt-5 grid grid-cols-1 gap-4 xl:grid-cols-6">
      <div v-for="item in 6" :key="item" class="ops-anomaly-compare-skeleton h-28 animate-pulse rounded-2xl"></div>
    </div>

    <div v-else-if="summary" class="mt-5 space-y-5">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <article
          v-for="item in metricCards"
          :key="item.key"
          class="ops-anomaly-compare-metric-card rounded-2xl px-4 py-4"
        >
          <p class="ops-anomaly-compare-metric-label text-xs font-medium tracking-[0.14em]">{{ item.label }}</p>
          <div class="mt-3 flex items-end justify-between gap-3">
            <div>
              <p class="ops-anomaly-compare-metric-value text-xl font-semibold">{{ item.current }}</p>
              <p class="ops-anomaly-compare-metric-prev mt-1 text-xs">
                {{ t('admin.ops.compare.previousWindow') }} {{ item.previous }}
              </p>
            </div>
            <span
              class="ops-anomaly-compare-chip rounded-full px-2.5 py-1 text-xs font-semibold"
              :class="{
                'ops-anomaly-compare-chip-up': item.tone === 'up',
                'ops-anomaly-compare-chip-down': item.tone === 'down',
                'ops-anomaly-compare-chip-flat': item.tone === 'flat'
              }"
            >
              {{ item.delta }}
            </span>
          </div>
        </article>
      </div>

      <div class="grid grid-cols-1 gap-4 xl:grid-cols-2">
        <article
          v-for="item in shiftCards"
          :key="item.key"
          class="ops-anomaly-compare-shift-card rounded-2xl px-4 py-4"
        >
          <div class="flex items-center justify-between gap-3">
            <p class="ops-anomaly-compare-shift-label text-[13px] font-semibold">{{ item.label }}</p>
            <span
              class="ops-anomaly-compare-chip rounded-full px-2.5 py-1 text-xs font-semibold"
              :class="item.changed
                ? 'ops-anomaly-compare-chip-shift'
                : 'ops-anomaly-compare-chip-flat'"
            >
              {{ item.changed ? t('admin.ops.compare.shifts.changed') : t('admin.ops.compare.shifts.stable') }}
            </span>
          </div>
          <p class="ops-anomaly-compare-shift-route mt-4 text-base font-semibold">
            {{ t('admin.ops.compare.shifts.fromTo', { from: item.from, to: item.to }) }}
          </p>
        </article>
      </div>
    </div>

    <div v-else class="mt-4">
      <EmptyState :title="t('common.noData')" :description="t('admin.ops.compare.empty')" />
    </div>
  </section>
</template>

<style>
.ops-anomaly-compare-shell {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  border: 1px solid rgba(201, 188, 165, 0.72);
  background: linear-gradient(180deg, rgba(251, 248, 242, 0.98), rgba(244, 238, 228, 0.96));
  box-shadow:
    0 22px 44px -36px rgba(82, 63, 41, 0.28),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
}

.ops-anomaly-compare-shell::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at top right, rgba(198, 184, 157, 0.24), transparent 36%),
    linear-gradient(135deg, rgba(167, 58, 42, 0.04), transparent 26%);
  opacity: 0.95;
}

.ops-anomaly-compare-shell > * {
  position: relative;
  z-index: 1;
}

.ops-anomaly-compare-header {
  border-bottom: 1px solid rgba(198, 184, 157, 0.42);
}

.ops-anomaly-compare-title {
  color: #3c3327;
}

.ops-anomaly-compare-mark {
  border: 1px solid rgba(167, 58, 42, 0.24);
  background: rgba(167, 58, 42, 0.08);
  color: #8f4032;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.42);
}

.ops-anomaly-compare-subtitle {
  color: #5e564c;
}

.ops-anomaly-compare-window {
  color: #8c8375;
}

.ops-anomaly-compare-action {
  border: 1px solid rgba(167, 58, 42, 0.24);
  background: rgba(255, 252, 247, 0.82);
  color: #6f3329;
  box-shadow: 0 10px 20px -18px rgba(104, 55, 47, 0.34);
}

.ops-anomaly-compare-action:hover {
  background: rgba(255, 248, 240, 0.96);
  border-color: rgba(167, 58, 42, 0.34);
}

.ops-anomaly-compare-error {
  border: 1px solid rgba(202, 95, 81, 0.22);
  background: rgba(253, 243, 241, 0.88);
  color: #9a4538;
}

.ops-anomaly-compare-skeleton {
  border: 1px solid rgba(214, 203, 183, 0.7);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.88), rgba(245, 238, 227, 0.88));
}

.ops-anomaly-compare-metric-card,
.ops-anomaly-compare-shift-card {
  border: 1px solid rgba(214, 203, 183, 0.82);
  box-shadow:
    0 16px 34px -30px rgba(90, 68, 45, 0.26),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
}

.ops-anomaly-compare-metric-card {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(249, 244, 237, 0.96));
}

.ops-anomaly-compare-shift-card {
  background: linear-gradient(180deg, rgba(245, 239, 228, 0.9), rgba(239, 232, 219, 0.94));
}

.ops-anomaly-compare-metric-label {
  color: #928879;
}

.ops-anomaly-compare-metric-value,
.ops-anomaly-compare-shift-label,
.ops-anomaly-compare-shift-route {
  color: #3b3227;
}

.ops-anomaly-compare-metric-prev {
  color: #847a6c;
}

.ops-anomaly-compare-chip {
  border: 1px solid transparent;
}

.ops-anomaly-compare-chip-up {
  border-color: rgba(182, 77, 63, 0.2);
  background: rgba(167, 58, 42, 0.1);
  color: #974739;
}

.ops-anomaly-compare-chip-down {
  border-color: rgba(88, 126, 101, 0.24);
  background: rgba(87, 124, 98, 0.12);
  color: #53745d;
}

.ops-anomaly-compare-chip-flat {
  border-color: rgba(137, 127, 109, 0.2);
  background: rgba(109, 99, 83, 0.08);
  color: #6f6657;
}

.ops-anomaly-compare-chip-shift {
  border-color: rgba(189, 141, 67, 0.22);
  background: rgba(201, 157, 78, 0.14);
  color: #966f2f;
}

.dark .ops-anomaly-compare-shell {
  border-color: rgba(74, 67, 55, 0.9);
  background: linear-gradient(180deg, rgba(27, 29, 23, 0.96), rgba(20, 22, 18, 0.98)) !important;
  box-shadow:
    0 22px 44px -36px rgba(0, 0, 0, 0.72),
    inset 0 1px 0 rgba(255, 255, 255, 0.04) !important;
}

.dark .ops-anomaly-compare-shell::before {
  background:
    radial-gradient(circle at top right, rgba(92, 78, 59, 0.16), transparent 40%),
    linear-gradient(135deg, rgba(167, 58, 42, 0.08), transparent 24%);
}

.dark .ops-anomaly-compare-header {
  border-bottom-color: rgba(61, 65, 54, 0.82);
}

.dark .ops-anomaly-compare-title,
.dark .ops-anomaly-compare-metric-value,
.dark .ops-anomaly-compare-shift-label,
.dark .ops-anomaly-compare-shift-route {
  color: #f4efe4;
}

.dark .ops-anomaly-compare-subtitle,
.dark .ops-anomaly-compare-metric-prev {
  color: #bdb5a8;
}

.dark .ops-anomaly-compare-window,
.dark .ops-anomaly-compare-metric-label {
  color: #968f84;
}

.dark .ops-anomaly-compare-mark {
  border-color: rgba(145, 53, 40, 0.46);
  background: rgba(136, 40, 27, 0.22);
  color: #f5d5cf;
  box-shadow: none;
}

.dark .ops-anomaly-compare-action {
  border-color: rgba(145, 53, 40, 0.42);
  background: rgba(112, 33, 23, 0.18);
  color: #f2e8df;
  box-shadow: none;
}

.dark .ops-anomaly-compare-action:hover {
  background: rgba(132, 37, 26, 0.26);
}

.dark .ops-anomaly-compare-error {
  border-color: rgba(133, 47, 35, 0.35);
  background: rgba(90, 22, 14, 0.22);
  color: #f2c1b9;
}

.dark .ops-anomaly-compare-skeleton {
  border-color: rgba(52, 56, 47, 0.84);
  background: linear-gradient(180deg, rgba(24, 26, 21, 0.94), rgba(18, 20, 16, 0.98));
}

.dark .ops-anomaly-compare-metric-card,
.dark .ops-anomaly-compare-shift-card {
  border-color: rgba(48, 52, 43, 0.82);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.03) !important;
}

.dark .ops-anomaly-compare-metric-card {
  background: linear-gradient(180deg, rgba(26, 28, 22, 0.98), rgba(19, 21, 16, 0.98)) !important;
}

.dark .ops-anomaly-compare-shift-card {
  background: rgba(22, 24, 19, 0.96) !important;
}

.dark .ops-anomaly-compare-chip-up {
  border-color: rgba(145, 53, 40, 0.34);
  background: rgba(120, 34, 24, 0.22);
  color: #f2c7bf;
}

.dark .ops-anomaly-compare-chip-down {
  border-color: rgba(60, 101, 73, 0.36);
  background: rgba(47, 85, 61, 0.22);
  color: #c6ead0;
}

.dark .ops-anomaly-compare-chip-flat {
  border-color: rgba(60, 64, 54, 0.82);
  background: rgba(42, 45, 36, 0.96);
  color: #bdb5a8;
}

.dark .ops-anomaly-compare-chip-shift {
  border-color: rgba(145, 112, 54, 0.34);
  background: rgba(133, 95, 32, 0.22);
  color: #f1ddb5;
}
</style>
