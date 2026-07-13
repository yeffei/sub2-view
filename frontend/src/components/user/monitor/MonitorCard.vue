<template>
  <button
    type="button"
    class="monitor-card group text-left p-4 rounded-[1rem] min-h-[236px] w-full overflow-hidden bg-white/70 backdrop-blur-xl border border-gray-200/80 shadow-card dark:border-[#44473a]/80 hover:-translate-y-0.5 hover:shadow-card-hover dark:hover:border-[#a73a2a]/40 hover:border-gray-300 transition-all duration-300 ease-out flex flex-col"
    @click="emit('click')"
  >
    <!-- Header: pool name + status chip -->
    <div class="flex items-start gap-2.5">
      <span class="w-8 h-8 rounded-lg ring-1 ring-black/5 dark:ring-white/10 grid place-items-center flex-shrink-0 bg-stone-100 text-stone-500 dark:bg-dark-700 dark:text-stone-300">
        <span class="h-2.5 w-2.5 rounded-full" :class="poolSignalClass"></span>
      </span>
      <div class="flex-1 min-w-0">
        <div class="text-[15px] font-semibold truncate text-gray-900 dark:text-gray-100">
          {{ item.name }}
        </div>
        <div class="mt-0.5 text-[11px] text-stone-500 dark:text-stone-400">
          {{ t('channelStatus.card.poolRuntimeOnly') }}
        </div>
      </div>
      <span
        class="px-2.5 py-1 rounded-full text-[11px] font-semibold flex-shrink-0"
        :class="statusBadgeClass(item.status)"
      >
        {{ statusLabel(item.status) }}
      </span>
    </div>
    <div class="mt-2 text-[11px] text-stone-500 dark:text-stone-400">
      容量：{{ capacityStatusLabel(item.capacity_status) }}
    </div>

    <!-- Metrics -->
    <MonitorMetricPair
      primary-icon="bolt"
      :primary-label="t('monitorCommon.dialogLatency')"
      :primary-value="formatLatency(item.best_latency_ms)"
      primary-unit="ms"
      secondary-icon="globe"
      :secondary-label="t('monitorCommon.endpointPing')"
      :secondary-value="formatLatency(item.best_ping_latency_ms)"
      secondary-unit="ms"
    />

    <!-- Divider -->
  <div class="mt-3 border-t border-gray-100 dark:border-dark-700/60"></div>

    <!-- Availability row -->
    <MonitorAvailabilityRow
      :window-label="availabilityLabel"
      :value="availabilityValue"
      :samples-label="t('channelStatus.detailSummary.recentSamplesHint')"
    />

    <!-- Timeline -->
    <MonitorTimeline
      :buckets="item.timeline"
      :countdown-seconds="countdownSeconds"
    />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { PoolHealthView } from '@/api/poolHealth'
import {
  useChannelMonitorFormat,
} from '@/composables/useChannelMonitorFormat'
import MonitorMetricPair from './MonitorMetricPair.vue'
import MonitorAvailabilityRow from './MonitorAvailabilityRow.vue'
import MonitorTimeline from './MonitorTimeline.vue'

const props = defineProps<{
  item: PoolHealthView
  window: '7d' | '15d' | '30d'
  availabilityValue: number | null
  countdownSeconds: number
}>()

const emit = defineEmits<{
  (e: 'click'): void
}>()

const { t } = useI18n()
const {
  statusLabel,
  statusBadgeClass,
  formatLatency,
} = useChannelMonitorFormat()

const poolSignalClass = computed(() => {
  if (props.item.status === 'operational') return 'bg-emerald-500'
  if (props.item.status === 'degraded') return 'bg-amber-500'
  return 'bg-rose-500'
})

function capacityStatusLabel(status: PoolHealthView['capacity_status']) {
  return ({ ample: '余量充足', observe: '需要观察', tight: '余量紧张', queueing: '排队中' } as const)[status]
}

const availabilityLabel = computed(() => {
  const win = t(`channelStatus.windowTab.${props.window}`)
  return `${t('monitorCommon.availabilityPrefix')} · ${win}`
})

</script>

<style scoped>
.monitor-card {
  min-width: 0;
}

.monitor-card :deep(*) {
  min-width: 0;
}

.monitor-card :deep(.monitor-timeline) {
  min-width: 0;
}

.dark .monitor-card {
  background:
    linear-gradient(180deg, rgba(255, 247, 235, 0.035), transparent 24%),
    rgba(17, 19, 15, 0.78);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.025), 0 24px 54px -46px rgba(0, 0, 0, 0.72);
}
</style>
