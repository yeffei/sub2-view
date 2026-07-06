<template>
  <BaseDialog
    :show="show"
    :title="title"
    width="content"
    panel-class="pool-health-dialog"
    header-class="pool-health-dialog-header"
    body-class="pool-health-dialog-body"
    @close="$emit('close')"
  >
    <div v-if="loading" class="pool-health-state">
      {{ t('common.loading') }}
    </div>
    <div v-else-if="!detail" class="pool-health-state">
      {{ t('channelStatus.detailLoadError') }}
    </div>
    <div v-else class="pool-health-scroll">
      <section class="pool-health-hero">
        <div class="pool-health-seal" :class="statusSealClass(detail.status)" aria-hidden="true">
          {{ sealGlyph }}
        </div>
        <div class="pool-health-hero-copy">
          <div class="pool-health-kicker">{{ t('channelStatus.detailSections.poolRuntime') }}</div>
          <div class="pool-health-status-line">
            <span class="pool-health-status-dot" :class="statusDotClass(detail.status)"></span>
            <strong>{{ statusLabel(detail.status) }}</strong>
          </div>
        </div>
      </section>

      <section class="pool-health-paper-grid">
        <div class="pool-health-stat pool-health-stat-primary">
          <span>{{ t('channelStatus.detailSummary.availability7d') }}</span>
          <strong>{{ formatPercent(detail.availability_7d) }}</strong>
        </div>
        <div class="pool-health-stat">
          <span>{{ t('channelStatus.detailSummary.availability15d') }}</span>
          <strong>{{ formatPercent(detail.availability_15d) }}</strong>
        </div>
        <div class="pool-health-stat">
          <span>{{ t('channelStatus.detailSummary.availability30d') }}</span>
          <strong>{{ formatPercent(detail.availability_30d) }}</strong>
        </div>
        <div class="pool-health-stat">
          <span>{{ t('channelStatus.detailSummary.bestLatency') }}</span>
          <strong>{{ formatLatencyWithUnit(detail.best_latency_ms) }}</strong>
        </div>
        <div class="pool-health-stat">
          <span>{{ t('monitorCommon.endpointPing') }}</span>
          <strong>{{ formatLatencyWithUnit(detail.best_ping_latency_ms) }}</strong>
        </div>
      </section>

      <section class="pool-health-samples">
        <div class="pool-health-samples-head">
          <div>
            <span>{{ t('channelStatus.detailSections.recentSamples') }}</span>
            <strong>{{ t('channelStatus.detailSummary.recentSamplesHint') }}</strong>
          </div>
        </div>
        <div class="pool-health-timeline">
          <div v-if="detail.timeline.length > 0" class="pool-health-timeline-track">
            <span
              v-for="point in detail.timeline"
              :key="point.checked_at"
              class="pool-health-timeline-mark"
              :class="timelineDotClass(point.status)"
              :title="`${statusLabel(point.status)} · ${formatRelativeTime(point.checked_at)}`"
            ></span>
          </div>
          <div v-else class="pool-health-empty">{{ t('channelStatus.detailSummary.noSamples') }}</div>
        </div>
      </section>
    </div>

  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import {
  status as fetchPoolHealthDetail,
  type PoolHealthDetail,
} from '@/api/poolHealth'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { useChannelMonitorFormat } from '@/composables/useChannelMonitorFormat'

const props = defineProps<{
  show: boolean
  monitorId: number | null
  title: string
}>()

defineEmits<{
  (e: 'close'): void
}>()

const { t } = useI18n()
const appStore = useAppStore()
const {
  statusLabel,
  formatLatency,
  formatPercent,
  formatRelativeTime,
} = useChannelMonitorFormat()

const sealGlyph = '池'

const detail = ref<PoolHealthDetail | null>(null)
const loading = ref(false)

async function load(id: number) {
  detail.value = null
  loading.value = true
  try {
    detail.value = await fetchPoolHealthDetail(id)
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.detailLoadError')))
  } finally {
    loading.value = false
  }
}

watch(
  () => [props.show, props.monitorId] as const,
  ([show, id]) => {
    if (!show) {
      detail.value = null
      return
    }
    if (id != null) void load(id)
  },
  { immediate: true },
)

function timelineDotClass(status: string) {
  if (status === 'operational') return 'is-good'
  if (status === 'degraded') return 'is-warn'
  return 'is-bad'
}

function statusDotClass(status: string) {
  if (status === 'operational') return 'is-good'
  if (status === 'degraded') return 'is-warn'
  return 'is-bad'
}

function statusSealClass(status: string) {
  if (status === 'operational') return 'is-good'
  if (status === 'degraded') return 'is-warn'
  return 'is-bad'
}

function formatLatencyWithUnit(ms: number | null | undefined) {
  if (ms == null) return formatLatency(ms)
  return `${formatLatency(ms)} ms`
}

</script>

<style>
.modal-content.pool-health-dialog {
  position: relative;
  flex: none !important;
  overflow: hidden;
  width: min(42rem, calc(100vw - 2rem)) !important;
  height: auto !important;
  min-height: 0 !important;
  max-height: calc(100vh - 3rem) !important;
  border: 1px solid rgba(179, 159, 119, 0.52);
  border-radius: 1.35rem;
  background:
    linear-gradient(90deg, rgba(122, 93, 58, 0.015) 1px, transparent 1px) 0 0 / 3.2rem 100%,
    linear-gradient(180deg, rgba(122, 93, 58, 0.014) 1px, transparent 1px) 0 0 / 100% 3.2rem,
    radial-gradient(circle at 96% 8%, rgba(167, 58, 42, 0.05), transparent 16rem),
    radial-gradient(circle at 8% 4%, rgba(170, 137, 83, 0.07), transparent 15rem),
    linear-gradient(180deg, #fffefa 0%, #faf5eb 100%);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.92),
    0 28px 84px -50px rgba(48, 38, 25, 0.62);
}

.modal-content.pool-health-dialog::before {
  content: "";
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 24%, rgba(72, 56, 35, 0.026) 0 1px, transparent 1px) 0 0 / 1rem 1rem,
    linear-gradient(115deg, transparent 0 43%, rgba(255, 255, 255, 0.18) 48%, transparent 53% 100%);
  opacity: 0.38;
}

.modal-header.pool-health-dialog-header {
  position: relative;
  z-index: 1;
  border-bottom: 1px solid rgba(190, 176, 148, 0.44);
  background: rgba(255, 253, 248, 0.82);
  padding: 0.82rem 1rem 0.72rem;
}

.modal-header.pool-health-dialog-header .modal-title {
  color: #25281f;
  font-size: clamp(1.12rem, 1.7vw, 1.35rem);
  font-weight: 650;
  letter-spacing: -0.03em;
}

.modal-body.pool-health-dialog-body {
  position: relative;
  z-index: 1;
  flex: 0 0 auto !important;
  overflow: visible !important;
  padding: 0.72rem 1rem 1rem !important;
}

.pool-health-state {
  padding: 2.6rem 1rem;
  text-align: center;
  color: rgba(86, 68, 52, 0.72);
  font-size: 0.92rem;
}

.pool-health-scroll {
  display: flex;
  flex-direction: column;
  gap: 0.58rem;
}

.pool-health-hero {
  position: relative;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-items: center;
  gap: 0.72rem;
  min-height: 4.85rem;
  border: 1px solid rgba(179, 159, 119, 0.34);
  border-radius: 1rem;
  background:
    radial-gradient(circle at 10% 36%, rgba(81, 98, 79, 0.06), transparent 12rem),
    linear-gradient(90deg, rgba(255, 255, 255, 0.52), rgba(248, 242, 231, 0.64)),
    linear-gradient(180deg, rgba(255, 255, 255, 0.4), rgba(239, 230, 211, 0.12));
  padding: 0.78rem 0.88rem;
  overflow: hidden;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.72),
    0 14px 34px -32px rgba(69, 51, 31, 0.55);
}

.pool-health-hero::after {
  content: "";
  position: absolute;
  right: clamp(1.6rem, 8vw, 6.2rem);
  top: 0.8rem;
  bottom: 0.8rem;
  width: 1px;
  background: rgba(167, 58, 42, 0.16);
  box-shadow:
    -3.4rem 0 0 rgba(179, 159, 119, 0.12),
    -6.8rem 0 0 rgba(179, 159, 119, 0.06);
}

.pool-health-seal {
  display: grid;
  place-items: center;
  width: 2.58rem;
  height: 2.58rem;
  border-radius: 0.7rem;
  color: #9f382b;
  font-family: "STKaiti", "KaiTi", serif;
  font-size: 1.28rem;
  font-weight: 700;
  background: rgba(167, 58, 42, 0.08);
  border: 1px solid rgba(167, 58, 42, 0.24);
  transform: rotate(-2deg);
  box-shadow: inset 0 0 0 0.22rem rgba(255, 250, 241, 0.6);
}

.pool-health-seal.is-good {
  color: #51624f;
  background: rgba(81, 98, 79, 0.09);
  border-color: rgba(81, 98, 79, 0.24);
}

.pool-health-seal.is-warn {
  color: #a45e1d;
  background: rgba(155, 129, 85, 0.12);
  border-color: rgba(155, 129, 85, 0.26);
}

.pool-health-hero-copy {
  min-width: 0;
}

.pool-health-kicker,
.pool-health-stat span,
.pool-health-hero-number span,
.pool-health-samples-head span {
  color: rgba(86, 68, 52, 0.72);
  font-size: 0.66rem;
  font-weight: 600;
  letter-spacing: 0.12em;
}

.pool-health-status-line {
  display: flex;
  align-items: center;
  gap: 0.48rem;
  margin-top: 0.38rem;
  color: #1f241b;
}

.pool-health-status-line strong {
  font-size: clamp(1.26rem, 2vw, 1.62rem);
  line-height: 1.04;
  letter-spacing: -0.045em;
}

.pool-health-status-dot {
  width: 0.55rem;
  height: 0.55rem;
  border-radius: 999px;
  box-shadow:
    inset 0 1px 1px rgba(255, 255, 255, 0.38),
    0 0 0 0.26rem rgba(167, 58, 42, 0.07);
}

.pool-health-status-dot.is-good {
  background:
    radial-gradient(circle at 34% 28%, rgba(232, 255, 237, 0.95), transparent 0.14rem),
    linear-gradient(135deg, #5cb897, #2e806b);
  box-shadow:
    inset 0 1px 1px rgba(255, 255, 255, 0.44),
    0 0 0 0.26rem rgba(63, 162, 123, 0.1),
    0 0.26rem 0.8rem -0.38rem rgba(63, 162, 123, 0.7);
}

.pool-health-status-dot.is-warn {
  background:
    radial-gradient(circle at 34% 28%, rgba(255, 244, 212, 0.95), transparent 0.14rem),
    linear-gradient(135deg, #d39a47, #98671f);
  box-shadow:
    inset 0 1px 1px rgba(255, 255, 255, 0.42),
    0 0 0 0.26rem rgba(192, 130, 45, 0.1);
}

.pool-health-status-dot.is-bad {
  background:
    radial-gradient(circle at 34% 28%, rgba(255, 228, 220, 0.9), transparent 0.14rem),
    linear-gradient(135deg, #c24a38, #8c2e24);
  box-shadow:
    inset 0 1px 1px rgba(255, 255, 255, 0.36),
    0 0 0 0.26rem rgba(167, 58, 42, 0.1);
}

.pool-health-paper-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  border: 1px solid rgba(179, 159, 119, 0.3);
  border-radius: 0.95rem;
  overflow: hidden;
  background:
    radial-gradient(circle at 12% 18%, rgba(72, 56, 35, 0.025) 0 1px, transparent 1px) 0 0 / 1rem 1rem,
    rgba(255, 253, 248, 0.54);
}

.pool-health-stat {
  position: relative;
  min-height: 4.35rem;
  padding: 0.62rem 0.72rem;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.26), rgba(255, 255, 255, 0.1));
}

.pool-health-stat + .pool-health-stat {
  border-left: 1px solid rgba(179, 159, 119, 0.24);
}

.pool-health-stat:nth-child(4) {
  border-left: 0;
}

.pool-health-stat:nth-child(n + 4) {
  border-top: 1px solid rgba(179, 159, 119, 0.24);
}

.pool-health-stat strong {
  display: block;
  margin-top: 0.68rem;
  color: #1f241b;
  font-size: clamp(0.98rem, 1.42vw, 1.2rem);
  line-height: 1;
  letter-spacing: -0.045em;
  font-variant-numeric: tabular-nums;
}

.pool-health-stat-primary {
  background:
    linear-gradient(180deg, rgba(81, 98, 79, 0.075), transparent),
    rgba(255, 255, 255, 0.36);
}

.pool-health-samples {
  border: 1px solid rgba(179, 159, 119, 0.3);
  border-radius: 0.95rem;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.24), rgba(255, 255, 255, 0.08)),
    rgba(255, 253, 248, 0.48);
  padding: 0.68rem;
}

.pool-health-samples-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 0.48rem;
}

.pool-health-samples-head div {
  display: flex;
  align-items: baseline;
  gap: 0.65rem;
}

.pool-health-samples-head strong {
  color: rgba(86, 68, 52, 0.56);
  font-size: 0.72rem;
  font-weight: 500;
}

.pool-health-timeline {
  position: relative;
  overflow: hidden;
  border-radius: 0.72rem;
  background:
    radial-gradient(ellipse at 50% 100%, rgba(81, 98, 79, 0.045), transparent 68%),
    linear-gradient(90deg, rgba(179, 159, 119, 0.11), transparent 1px) 0 0 / 1.25rem 100%,
    linear-gradient(180deg, transparent 49%, rgba(122, 93, 58, 0.09) 49% 51%, transparent 51%),
    rgba(255, 255, 255, 0.38);
  padding: 0.9rem 0.86rem;
}

.pool-health-timeline::before {
  content: "";
  position: absolute;
  left: 0.86rem;
  right: 0.86rem;
  top: 50%;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(122, 93, 58, 0.2), transparent);
}

.pool-health-timeline::after {
  content: "";
  position: absolute;
  left: 0.86rem;
  right: 0.86rem;
  top: calc(50% + 0.34rem);
  height: 0.18rem;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, rgba(81, 98, 79, 0.06), transparent);
  filter: blur(0.5px);
}

.pool-health-timeline-track {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 0.16rem;
  overflow: hidden;
}

.pool-health-timeline-mark {
  position: relative;
  flex: 1 1 0;
  min-width: 0.16rem;
  max-width: 0.58rem;
  height: 1.05rem;
  border-radius: 0.08rem;
  opacity: 0.92;
  transform: translateZ(0);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.24),
    0 0.22rem 0.5rem -0.45rem rgba(48, 38, 25, 0.5);
}

.pool-health-timeline-mark.is-good {
  background: linear-gradient(180deg, #72c1a1, #3d9879);
}

.pool-health-timeline-mark.is-warn {
  background: linear-gradient(180deg, #d7a353, #b27427);
}

.pool-health-timeline-mark.is-bad {
  background: linear-gradient(180deg, #c85b49, #a73a2a);
}

.pool-health-timeline-mark::before {
  content: "";
  position: absolute;
  left: 50%;
  top: 0.14rem;
  bottom: 0.14rem;
  width: 1px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.2);
  transform: translateX(-50%);
}

.pool-health-timeline-mark::after {
  content: "";
  display: none;
}

.pool-health-empty {
  color: rgba(86, 68, 52, 0.62);
  font-size: 0.88rem;
}

html.dark .modal-content.pool-health-dialog {
  border-color: rgba(118, 102, 70, 0.58);
  background:
    linear-gradient(90deg, rgba(226, 211, 178, 0.014) 1px, transparent 1px) 0 0 / 3.2rem 100%,
    linear-gradient(180deg, rgba(226, 211, 178, 0.012) 1px, transparent 1px) 0 0 / 100% 3.2rem,
    radial-gradient(circle at 96% 8%, rgba(167, 58, 42, 0.065), transparent 16rem),
    radial-gradient(circle at 7% 0%, rgba(195, 153, 85, 0.08), transparent 14rem),
    linear-gradient(180deg, rgba(18, 22, 18, 0.985), rgba(10, 13, 11, 0.99));
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 28px 90px -46px rgba(0, 0, 0, 0.9);
}

html.dark .modal-content.pool-health-dialog::before {
  opacity: 0.23;
}

html.dark .modal-header.pool-health-dialog-header {
  border-color: rgba(91, 83, 62, 0.62);
  background: rgba(20, 23, 19, 0.88);
}

html.dark .modal-header.pool-health-dialog-header .modal-title,
html.dark .pool-health-status-line,
html.dark .pool-health-stat strong {
  color: rgba(245, 238, 222, 0.94);
}

html.dark .pool-health-kicker,
html.dark .pool-health-stat span,
html.dark .pool-health-samples-head span {
  color: rgba(220, 209, 184, 0.7);
}

html.dark .pool-health-hero,
html.dark .pool-health-paper-grid,
html.dark .pool-health-samples {
  border-color: rgba(112, 97, 67, 0.48);
  background:
    radial-gradient(circle at 16% 20%, rgba(232, 213, 173, 0.026) 0 1px, transparent 1px) 0 0 / 1rem 1rem,
    rgba(24, 27, 22, 0.68);
}

html.dark .pool-health-stat {
  background: linear-gradient(180deg, rgba(255, 247, 235, 0.026), rgba(255, 247, 235, 0.008));
}

html.dark .pool-health-stat + .pool-health-stat {
  border-left-color: rgba(112, 97, 67, 0.42);
}

html.dark .pool-health-stat:nth-child(n + 4) {
  border-top-color: rgba(112, 97, 67, 0.42);
}

html.dark .pool-health-seal {
  box-shadow: inset 0 0 0 0.35rem rgba(20, 23, 19, 0.42);
}

html.dark .pool-health-timeline {
  background:
    radial-gradient(ellipse at 50% 100%, rgba(155, 197, 154, 0.045), transparent 68%),
    linear-gradient(90deg, rgba(226, 211, 178, 0.038), transparent 1px) 0 0 / 1.25rem 100%,
    linear-gradient(180deg, transparent 49%, rgba(226, 211, 178, 0.055) 49% 51%, transparent 51%),
    rgba(11, 13, 12, 0.3);
}

html.dark .pool-health-samples-head strong,
html.dark .pool-health-empty,
html.dark .pool-health-state {
  color: rgba(220, 209, 184, 0.62);
}

@media (max-width: 760px) {
  .pool-health-dialog-header,
  .pool-health-dialog-body {
    padding-inline: 1rem;
  }

  .pool-health-hero {
    grid-template-columns: auto minmax(0, 1fr);
  }

  .pool-health-paper-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .pool-health-stat + .pool-health-stat {
    border-left: 0;
  }

  .pool-health-stat:nth-child(n + 4) {
    border-top: 0;
  }

  .pool-health-stat:nth-child(even) {
    border-left: 1px solid rgba(177, 160, 130, 0.25);
  }

  .pool-health-stat:nth-child(n + 3) {
    border-top: 1px solid rgba(177, 160, 130, 0.25);
  }
}

@media (max-width: 460px) {
  .pool-health-paper-grid {
    grid-template-columns: 1fr;
  }

  .pool-health-stat:nth-child(even) {
    border-left: 0;
  }

  .pool-health-stat + .pool-health-stat {
    border-top: 1px solid rgba(177, 160, 130, 0.25);
  }
}
</style>
