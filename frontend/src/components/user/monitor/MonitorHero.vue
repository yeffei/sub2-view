<template>
  <section class="monitor-hero py-2 md:py-3">
    <div class="monitor-hero-shell">
      <div class="monitor-hero-topline">
        <div class="monitor-hero-copy">
          <div class="monitor-hero-heading">
            <h1>{{ t('channelStatus.title') }}</h1>
            <span class="monitor-hero-status" :class="overallChipClass">
              <span class="w-1.5 h-1.5 rounded-full" :class="overallDotClass"></span>
              {{ overallLabel }}
            </span>
          </div>
        </div>

        <div class="monitor-summary-strip" role="list" :aria-label="t('channelStatus.summary.eyebrow')">
          <div class="monitor-summary-item" role="listitem">
            <span>{{ t('channelStatus.summary.monitored') }}</span>
            <strong>{{ totalCount }}</strong>
          </div>
          <div class="monitor-summary-item monitor-summary-item-good" role="listitem">
            <span>{{ t('channelStatus.summary.healthy') }}</span>
            <strong>{{ healthyCount }}</strong>
          </div>
          <div class="monitor-summary-item monitor-summary-item-warn" role="listitem">
            <span>{{ t('channelStatus.summary.attention') }}</span>
            <strong>{{ attentionCount }}</strong>
          </div>
        </div>
      </div>

      <div class="monitor-toolbar">
        <div role="tablist" class="monitor-window-tabs">
          <button
            v-for="opt in windowOptions"
            :key="opt.value"
            type="button"
            role="tab"
            :aria-selected="window === opt.value"
            class="monitor-window-tab"
            :class="window === opt.value ? 'is-active' : ''"
            @click="emit('update:window', opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>

        <div class="monitor-toolbar-actions">
          <button
            type="button"
            class="monitor-refresh-button"
            :disabled="loading"
            :title="t('common.refresh')"
            @click="emit('refresh')"
          >
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            <span>{{ t('common.refresh') }}</span>
          </button>

          <AutoRefreshButton
            v-if="autoRefresh"
            class="monitor-auto-refresh"
            :enabled="autoRefresh.enabled.value"
            :interval-seconds="autoRefresh.intervalSeconds.value"
            :countdown="autoRefresh.countdown.value"
            :intervals="autoRefresh.intervals"
            @update:enabled="autoRefresh.setEnabled"
            @update:interval="autoRefresh.setInterval"
          />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AutoRefreshButton from '@/components/common/AutoRefreshButton.vue'
export type MonitorWindow = '7d' | '15d' | '30d'
export type OverallStatus = 'operational' | 'degraded'

const props = defineProps<{
  overallStatus: OverallStatus
  intervalSeconds: number
  window: MonitorWindow
  loading: boolean
  totalCount: number
  healthyCount: number
  attentionCount: number
  latestCheckedAt: string | null
  autoRefresh?: {
    enabled: { value: boolean }
    intervalSeconds: { value: number }
    countdown: { value: number }
    intervals: readonly number[]
    setEnabled: (v: boolean) => void
    setInterval: (v: number) => void
  }
}>()

const emit = defineEmits<{
  (e: 'update:window', value: MonitorWindow): void
  (e: 'refresh'): void
}>()

const { t } = useI18n()

const windowOptions = computed<{ value: MonitorWindow; label: string }[]>(() => [
  { value: '7d', label: t('channelStatus.windowTab.7d') },
  { value: '15d', label: t('channelStatus.windowTab.15d') },
  { value: '30d', label: t('channelStatus.windowTab.30d') },
])

const overallLabel = computed(() => t(`channelStatus.overall.${props.overallStatus}`))

const overallChipClass = computed(() => {
  switch (props.overallStatus) {
    case 'operational':
      return 'border-[#51624f]/25 bg-[#51624f]/10 text-[#51624f] dark:bg-[#51624f]/20 dark:text-emerald-300'
    case 'degraded':
    default:
      return 'border-[#9b8155]/25 bg-[#9b8155]/12 text-[#7b6a53] dark:bg-[#9b8155]/20 dark:text-amber-300'
  }
})

const overallDotClass = computed(() => {
  switch (props.overallStatus) {
    case 'operational':
      return 'bg-[#51624f] animate-pulse'
    case 'degraded':
    default:
      return 'bg-[#9b8155] animate-pulse'
  }
})

</script>

<style scoped>
.monitor-hero-shell {
  padding: 0.2rem 0;
}

.monitor-hero-topline {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.55rem 0.82rem;
  flex-wrap: wrap;
}

.monitor-hero-copy {
  display: flex;
  flex: 0 1 auto;
  min-width: 0;
  max-width: 17.25rem;
  flex-direction: column;
  gap: 0.14rem;
}

.monitor-hero-heading {
  display: flex;
  align-items: center;
  gap: 0.34rem;
  flex-wrap: wrap;
}

.monitor-hero-heading h1 {
  margin: 0;
  font-size: 1.16rem;
  line-height: 1.2;
  font-weight: 600;
  color: #2f241d;
}


.monitor-hero-status {
  display: inline-flex;
  align-items: center;
  gap: 0.34rem;
  border-radius: 999px;
  border-width: 1px;
  padding: 0.18rem 0.52rem;
  font-size: 0.66rem;
  font-weight: 600;
}

.monitor-summary-strip {
  display: inline-grid;
  grid-template-columns: repeat(3, minmax(5.75rem, auto));
  gap: 0.14rem;
  flex: 0 0 auto;
  align-self: flex-start;
  margin-top: 0.01rem;
  width: fit-content;
  max-width: 100%;
  border-top: 1px solid rgba(198, 184, 157, 0.22);
  border-bottom: 1px solid rgba(198, 184, 157, 0.18);
}

.monitor-summary-item {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.34rem;
  min-width: 5.75rem;
  padding: 0.34rem 0.02rem;
}

.monitor-summary-item + .monitor-summary-item {
  border-left: 1px solid rgba(198, 184, 157, 0.2);
  padding-left: 0.42rem;
}

.monitor-summary-item span {
  font-size: 0.58rem;
  color: rgba(123, 106, 83, 0.84);
  white-space: nowrap;
}

.monitor-summary-item strong {
  flex-shrink: 0;
  font-size: 0.84rem;
  line-height: 1.05;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  color: #2f241d;
}

.monitor-summary-item-good strong {
  color: #51624f;
}

.monitor-summary-item-warn strong {
  color: #8e5f3a;
}

.monitor-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
  margin-top: 0.7rem;
  padding-top: 0.75rem;
  border-top: 1px solid rgba(198, 184, 157, 0.3);
}

.monitor-window-tabs {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 0.14rem;
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 0.85rem;
  background: rgba(250, 247, 239, 0.82);
  padding: 0.24rem;
}

.monitor-window-tab {
  border: 0;
  border-radius: 0.68rem;
  background: transparent;
  padding: 0.46rem 0.82rem;
  font-size: 0.79rem;
  font-weight: 500;
  color: rgba(106, 86, 67, 0.86);
  transition: color 0.18s ease, background-color 0.18s ease, box-shadow 0.18s ease;
}

.monitor-window-tab:hover,
.monitor-window-tab:focus-visible {
  color: #8f3426;
  outline: none;
}

.monitor-window-tab.is-active {
  background: rgba(255, 252, 247, 0.98);
  color: #8f3426;
  box-shadow: 0 0 0 1px rgba(167, 58, 42, 0.16) inset;
}

.monitor-toolbar-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.55rem;
}

.monitor-refresh-button {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  height: 2.1rem;
  border: 1px solid rgba(198, 184, 157, 0.46);
  border-radius: 0.8rem;
  background: rgba(250, 247, 239, 0.88);
  padding: 0 0.9rem;
  font-size: 0.78rem;
  font-weight: 500;
  color: rgba(86, 68, 52, 0.9);
  transition: color 0.18s ease, border-color 0.18s ease, background-color 0.18s ease;
}

.monitor-refresh-button:hover,
.monitor-refresh-button:focus-visible {
  border-color: rgba(167, 58, 42, 0.3);
  background: rgba(247, 240, 228, 0.96);
  color: #8f3426;
  outline: none;
}

.monitor-refresh-button:disabled {
  opacity: 0.56;
}

.monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger) {
  height: 2.1rem;
  border-color: rgba(198, 184, 157, 0.46);
  background: rgba(250, 247, 239, 0.88);
  color: rgba(86, 68, 52, 0.9);
  padding-inline: 0.88rem;
  font-size: 0.78rem;
}

.monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger:hover),
.monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger:focus-visible),
.monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger-open) {
  border-color: rgba(167, 58, 42, 0.3);
  background: rgba(247, 240, 228, 0.96);
  color: #8f3426;
}

.monitor-hero :deep(.monitor-auto-refresh .auto-refresh-panel) {
  min-width: 11.5rem;
}

.dark .monitor-hero-shell {
  background: transparent;
}

.dark .monitor-hero-heading h1 {
  color: rgba(244, 239, 228, 0.94);
}

.dark .monitor-summary-item span {
  color: rgba(214, 205, 185, 0.76);
}

.dark .monitor-summary-strip {
  border-top-color: rgba(88, 86, 70, 0.44);
  border-bottom-color: rgba(88, 86, 70, 0.38);
}

.dark .monitor-summary-item + .monitor-summary-item {
  border-left-color: rgba(88, 86, 70, 0.34);
}

.dark .monitor-summary-item strong {
  color: rgba(244, 239, 228, 0.94);
}

.dark .monitor-summary-item-good strong {
  color: #9bc59a;
}

.dark .monitor-summary-item-warn strong {
  color: #d0ad7c;
}

.dark .monitor-toolbar {
  border-top-color: rgba(88, 86, 70, 0.54);
}

.dark .monitor-window-tabs {
  border-color: rgba(88, 86, 70, 0.64);
  background: rgba(27, 29, 24, 0.88);
}

.dark .monitor-window-tab {
  color: rgba(214, 205, 185, 0.76);
}

.dark .monitor-window-tab:hover,
.dark .monitor-window-tab:focus-visible,
.dark .monitor-window-tab.is-active,
.dark .monitor-refresh-button:hover,
.dark .monitor-refresh-button:focus-visible {
  color: #f0b4a8;
}

.dark .monitor-window-tab.is-active {
  background: rgba(44, 46, 39, 0.92);
  box-shadow: 0 0 0 1px rgba(167, 58, 42, 0.28) inset;
}

.dark .monitor-refresh-button {
  border-color: rgba(88, 86, 70, 0.64);
  background: rgba(27, 29, 24, 0.88);
  color: rgba(214, 205, 185, 0.84);
}

.dark .monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger) {
  border-color: rgba(88, 86, 70, 0.64);
  background: rgba(27, 29, 24, 0.88);
  color: rgba(214, 205, 185, 0.84);
}

.dark .monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger:hover),
.dark .monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger:focus-visible),
.dark .monitor-hero :deep(.monitor-auto-refresh .auto-refresh-trigger-open) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(36, 39, 31, 0.94);
  color: #f0b4a8;
}

@media (max-width: 900px) {
  .monitor-hero-topline {
    gap: 0.34rem;
  }

  .monitor-hero-copy {
    flex-basis: 100%;
    max-width: none;
  }

  .monitor-summary-strip {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    width: 100%;
  }

  .monitor-summary-item:nth-child(3) {
    border-left: 0;
    padding-left: 0.2rem;
  }
}

@media (max-width: 640px) {
  .monitor-summary-strip {
    grid-template-columns: 1fr;
  }

  .monitor-summary-item {
    min-width: 0;
    padding: 0.6rem 0;
  }

  .monitor-summary-item + .monitor-summary-item {
    border-left: 0;
    border-top: 1px solid rgba(198, 184, 157, 0.18);
    padding-left: 0;
  }

  .dark .monitor-summary-item + .monitor-summary-item {
    border-top-color: rgba(88, 86, 70, 0.28);
  }

  .monitor-toolbar {
    align-items: stretch;
  }

  .monitor-window-tabs,
  .monitor-toolbar-actions {
    width: 100%;
  }

  .monitor-toolbar-actions {
    justify-content: flex-end;
  }
}
</style>









