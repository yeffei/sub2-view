<template>
  <div class="mt-3 pt-2.5 border-t border-gray-100 dark:border-dark-700/60">
    <div
      class="mb-2 flex justify-between gap-3 text-[11px] font-medium text-gray-500 dark:text-gray-400"
    >
      <span>{{ t('monitorCommon.history60pts', { n: length }) }}</span>
      <span class="tabular-nums">{{ t('monitorCommon.nextUpdateIn', { n: countdownSeconds }) }}</span>
    </div>

    <div
      v-if="maintenance"
      class="flex h-5 w-full items-center justify-center rounded border border-dashed border-gray-300 dark:border-dark-600 text-[10px] text-gray-400"
    >
      {{ t('monitorCommon.maintenancePaused') }}
    </div>
    <div v-else class="monitor-timeline-bars" aria-hidden="true">
      <div
        v-for="(bar, idx) in displayBars"
        :key="idx"
        class="monitor-timeline-bar"
        :class="bar.colorClass"
        :style="{ height: bar.heightPct + '%' }"
        :title="bar.title"
      ></div>
    </div>

    <div
      class="mt-1 flex justify-between text-[10px] text-gray-400"
    >
      <span>{{ t('monitorCommon.past') }}</span>
      <span>{{ t('monitorCommon.now') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { MonitorTimelinePoint } from '@/api/channelMonitor'
import { useChannelMonitorFormat } from '@/composables/useChannelMonitorFormat'

const props = withDefaults(defineProps<{
  buckets?: MonitorTimelinePoint[]
  countdownSeconds: number
  length?: number
  maintenance?: boolean
}>(), {
  buckets: () => [],
  length: 60,
  maintenance: false,
})

const { t } = useI18n()
const { statusLabel, formatLatency, formatRelativeTime } = useChannelMonitorFormat()

interface Bar {
  colorClass: string
  heightPct: number
  title: string
}

// 缩略水尺：正常保持低矮稳定，异常用更高的竹节提醒，灰色代表未采样。
const STATUS_HEIGHT: Record<string, number> = {
  operational: 34,
  degraded: 62,
  failed: 86,
  error: 86,
  empty: 16,
}

const STATUS_COLOR: Record<string, string> = {
  operational: 'is-good',
  degraded: 'is-warn',
  failed: 'is-bad',
  error: 'is-bad',
  empty: 'is-empty',
}

const displayBars = computed<Bar[]>(() => {
  // Real points come newest-first; convert to oldest-first so the rightmost
  // bar represents "now". Pad the left with empty placeholders to keep the
  // bar count stable at `length`.
  const real = [...(props.buckets ?? [])]
    .slice(0, props.length)
    .reverse()

  const padCount = Math.max(0, props.length - real.length)
  const bars: Bar[] = []

  for (let i = 0; i < padCount; i += 1) {
    bars.push({
      colorClass: STATUS_COLOR.empty,
      heightPct: STATUS_HEIGHT.empty,
      title: '',
    })
  }

  for (const point of real) {
    const status = point.status as keyof typeof STATUS_HEIGHT
    const colorClass = STATUS_COLOR[status] ?? STATUS_COLOR.empty
    const heightPct = STATUS_HEIGHT[status] ?? STATUS_HEIGHT.empty
    const latency = formatLatency(point.latency_ms)
    const relative = formatRelativeTime(point.checked_at)
    const label = statusLabel(point.status)
    bars.push({
      colorClass,
      heightPct,
      title: `${relative} · ${label} · ${latency}ms`,
    })
  }

  return bars
})
</script>

<style scoped>
.monitor-timeline-bars {
  position: relative;
  display: grid;
  grid-template-columns: repeat(60, minmax(0, 1fr));
  align-items: center;
  gap: 3px;
  width: 100%;
  height: 26px;
  min-width: 0;
  overflow: hidden;
  border-radius: 999px;
  background:
    radial-gradient(ellipse at 50% 100%, rgba(81, 98, 79, 0.06), transparent 70%),
    linear-gradient(90deg, rgba(190, 176, 148, 0.13), transparent 1px) 0 0 / 1.18rem 100%;
  padding: 0 0.42rem;
}

.monitor-timeline-bars::before {
  content: "";
  position: absolute;
  left: 0.45rem;
  right: 0.45rem;
  top: 56%;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(122, 93, 58, 0.2), transparent);
}

.monitor-timeline-bars::after {
  content: "";
  position: absolute;
  left: 0.45rem;
  right: 0.45rem;
  top: calc(56% + 0.3rem);
  height: 0.22rem;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, rgba(81, 98, 79, 0.07), transparent);
}

.monitor-timeline-bar {
  position: relative;
  min-width: 0;
  align-self: center;
  border-radius: 999px;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.25),
    0 0.18rem 0.55rem -0.45rem rgba(48, 38, 25, 0.5);
}

.monitor-timeline-bar::before {
  content: "";
  position: absolute;
  left: 1px;
  top: 20%;
  bottom: 20%;
  width: 1px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.32);
}

.monitor-timeline-bar::after {
  content: "";
  position: absolute;
  inset: 1px 22% auto;
  height: 1px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.28);
}

.monitor-timeline-bar.is-good {
  background: #64ad8e;
  opacity: 0.9;
}

.monitor-timeline-bar.is-warn {
  background: #bc8d48;
  opacity: 0.96;
}

.monitor-timeline-bar.is-bad {
  background: #a73a2a;
  opacity: 0.98;
}

.monitor-timeline-bar.is-empty {
  background: rgba(190, 176, 148, 0.34);
  box-shadow: none;
}

@media (max-width: 640px) {
  .monitor-timeline-bars {
    gap: 1px;
  }
}
</style>
