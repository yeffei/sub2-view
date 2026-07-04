<template>
  <BaseDialog :show="show" :title="t('usage.errors.detail.title')" width="wide" @close="emit('update:show', false)">
    <div v-if="loading" class="flex justify-center py-10">
      <svg class="h-7 w-7 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
    </div>

    <div v-else-if="loadError" class="py-8 text-center text-sm text-red-500">
      {{ t('usage.errors.detail.loadFailed') }}
    </div>

    <div v-else-if="detail" class="space-y-4 text-sm">
      <div class="grid grid-cols-2 gap-x-6 gap-y-3">
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.time') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ formatDateTime(detail.created_at) }}</p>
        </div>
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.model') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ detail.model || '-' }}</p>
        </div>
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.endpoint') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ detail.inbound_endpoint || '-' }}</p>
        </div>
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.status') }}</span>
          <p class="mt-0.5">
            <span class="badge" :class="statusClass(detail.status_code)">{{ detail.status_code || '-' }}</span>
          </p>
        </div>
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.category') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ t('usage.errors.categories.' + detail.category) }}</p>
        </div>
        <div>
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.platform') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ detail.platform || '-' }}</p>
        </div>
        <div v-if="detail.upstream_status_code != null">
          <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.detail.upstreamStatus') }}</span>
          <p class="mt-0.5 text-gray-900 dark:text-dark-100">{{ detail.upstream_status_code }}</p>
        </div>
      </div>

      <div v-if="explanationSummary || explanationAdvice.length > 0" class="rounded-xl border border-amber-200/80 bg-amber-50/70 p-4 dark:border-amber-900/50 dark:bg-amber-950/10">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="min-w-0 flex-1">
            <span class="font-medium text-amber-800 dark:text-amber-300">{{ t('usage.errors.detail.explanationTitle') }}</span>
            <p v-if="explanationSummary" class="mt-1 text-sm leading-6 text-amber-950/85 dark:text-amber-100/90">{{ explanationSummary }}</p>
            <ul v-if="explanationAdvice.length > 0" class="mt-2 space-y-1 text-sm leading-6 text-amber-900/85 dark:text-amber-100/85">
              <li v-for="item in explanationAdvice" :key="item">{{ item }}</li>
            </ul>
            <div v-if="modelTrace" class="mt-3 rounded-lg border border-amber-200/80 bg-white/70 p-3 dark:border-amber-900/60 dark:bg-amber-950/20">
              <span class="text-[11px] font-semibold uppercase tracking-[0.18em] text-amber-700/80 dark:text-amber-300/80">
                {{ t('usage.errors.detail.modelTrace.title') }}
              </span>
              <div class="mt-2 grid gap-3 sm:grid-cols-2">
                <div>
                  <span class="text-xs text-amber-700/75 dark:text-amber-200/75">{{ t('usage.errors.detail.modelTrace.requested') }}</span>
                  <p class="mt-1 break-all text-sm font-medium text-amber-950 dark:text-amber-50">{{ modelTrace.requested }}</p>
                </div>
                <div v-if="modelTrace.upstream">
                  <span class="text-xs text-amber-700/75 dark:text-amber-200/75">{{ t('usage.errors.detail.modelTrace.upstream') }}</span>
                  <p class="mt-1 break-all text-sm font-medium text-amber-950 dark:text-amber-50">{{ modelTrace.upstream }}</p>
                </div>
              </div>
              <p v-if="modelTrace.hint" class="mt-2 text-xs leading-6 text-amber-900/85 dark:text-amber-100/85">{{ modelTrace.hint }}</p>
            </div>
          </div>
          <button
            type="button"
            class="inline-flex items-center rounded-full border border-amber-300 bg-white/80 px-3 py-1.5 text-xs font-medium text-amber-900 transition-colors hover:bg-white dark:border-amber-800 dark:bg-amber-950/10 dark:text-amber-100 dark:hover:bg-amber-950/20"
            @click="copyDiagnosticSummary"
          >
            {{ t('usage.errors.detail.copySummary') }}
          </button>
        </div>
        <button
          v-if="nextAction"
          type="button"
          class="mt-3 inline-flex items-center rounded-full border border-amber-300 bg-white/80 px-3 py-1.5 text-xs font-medium text-amber-900 transition-colors hover:bg-white dark:border-amber-800 dark:bg-amber-950/10 dark:text-amber-100 dark:hover:bg-amber-950/20"
          @click="goToNextAction"
        >
          {{ nextAction.label }}
        </button>
      </div>

      <div v-if="recoveryGuide" class="rounded-xl border border-stone-200/80 bg-stone-50/75 p-4 dark:border-dark-700 dark:bg-dark-900/50">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <span class="font-medium text-stone-800 dark:text-dark-100">{{ t('usage.errors.detail.recoveryGuide.title') }}</span>
            <p class="mt-1 text-sm leading-6 text-stone-700 dark:text-dark-300">{{ recoveryGuide.summary }}</p>
          </div>
          <span class="rounded-full border border-stone-200 bg-white/80 px-2.5 py-1 text-xs text-stone-500 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300">
            {{ t('usage.errors.detail.recoveryGuide.badge') }}
          </span>
        </div>

        <ol class="mt-4 space-y-3">
          <li v-for="step in recoveryGuide.steps" :key="step.title" class="rounded-lg border border-stone-200/80 bg-white/80 px-3 py-3 dark:border-dark-700 dark:bg-dark-800/70">
            <div class="flex items-start gap-3">
              <span class="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full border border-stone-200 bg-stone-100 text-xs font-semibold text-stone-700 dark:border-dark-600 dark:bg-dark-700 dark:text-dark-100">
                {{ step.index }}
              </span>
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <strong class="text-sm text-stone-900 dark:text-dark-100">{{ step.title }}</strong>
                  <span v-if="step.anchorLabel" class="rounded-full bg-stone-100 px-2 py-0.5 text-[11px] text-stone-600 dark:bg-dark-700 dark:text-dark-300">{{ step.anchorLabel }}</span>
                </div>
                <p class="mt-1 text-sm leading-6 text-stone-700 dark:text-dark-300">{{ step.detail }}</p>
              </div>
            </div>
          </li>
        </ol>

        <div v-if="recoveryGuide.checks.length" class="mt-4 rounded-lg border border-stone-200/80 bg-white/75 px-3 py-3 dark:border-dark-700 dark:bg-dark-800/60">
          <span class="font-medium text-stone-800 dark:text-dark-100">{{ t('usage.errors.detail.recoveryGuide.checksTitle') }}</span>
          <ul class="mt-2 space-y-1 text-sm leading-6 text-stone-700 dark:text-dark-300">
            <li v-for="item in recoveryGuide.checks" :key="item">{{ item }}</li>
          </ul>
        </div>
      </div>

      <div class="rounded-xl border border-slate-200/80 bg-slate-50/70 p-4 dark:border-dark-700 dark:bg-dark-900/40">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <span class="font-medium text-slate-800 dark:text-dark-100">{{ t('usage.errors.detail.timeline.title') }}</span>
            <p v-if="timelineSummary" class="mt-1 text-sm leading-6 text-slate-700 dark:text-dark-300">{{ timelineSummary }}</p>
          </div>
          <span class="rounded-full border border-slate-200 bg-white/80 px-2.5 py-1 text-xs text-slate-500 dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300">
            {{ t('usage.errors.detail.timeline.windowLabel') }}
          </span>
        </div>

        <div v-if="timelineLoading" class="mt-3 text-sm text-gray-500 dark:text-dark-400">
          {{ t('usage.errors.detail.timeline.loading') }}
        </div>
        <div v-else-if="timelineLoadFailed" class="mt-3 text-sm text-red-500">
          {{ t('usage.errors.detail.timeline.loadFailed') }}
        </div>
        <div v-else-if="!timelineItems.length" class="mt-3 text-sm text-gray-500 dark:text-dark-400">
          {{ t('usage.errors.detail.timeline.empty') }}
        </div>
        <ol v-else class="mt-4 space-y-3">
          <li
            v-for="item in timelineItems"
            :key="item.id"
            class="rounded-lg border px-3 py-3"
            :class="timelineToneClass(item.tone)"
          >
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <strong class="text-sm">{{ item.title }}</strong>
                  <span class="rounded-full px-2 py-0.5 text-[11px] font-medium" :class="timelineBadgeClass(item.tone)">{{ item.badge }}</span>
                </div>
                <p class="mt-1 text-sm leading-6 text-slate-700 dark:text-dark-300">{{ item.detail }}</p>
              </div>
              <span class="whitespace-nowrap text-xs text-slate-500 dark:text-dark-400">{{ formatDateTime(item.at) }}</span>
            </div>
          </li>
        </ol>
      </div>

      <div v-if="detail.message">
        <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.message') }}</span>
        <p class="mt-0.5 break-all text-gray-900 dark:text-dark-100">{{ detail.message }}</p>
      </div>

      <div v-if="detail.error_body">
        <span class="font-medium text-gray-500 dark:text-dark-400">{{ t('usage.errors.detail.responseBody') }}</span>
        <pre class="mt-1 max-h-[40vh] overflow-auto whitespace-pre-wrap break-all rounded-lg border border-gray-200 bg-gray-50 p-3 text-xs text-gray-800 dark:border-dark-700 dark:bg-dark-900 dark:text-dark-200">{{ detail.error_body }}</pre>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { getMyErrorDetail, listMyErrorRequests, query as queryUsage } from '@/api/usage'
import { useAppStore } from '@/stores/app'
import { formatDateTime } from '@/utils/format'
import type { UsageLog, UserErrorRequest, UserErrorRequestDetail } from '@/types'

const props = defineProps<{
  show: boolean
  errorId: number | null
}>()

const emit = defineEmits<{
  (e: 'update:show', v: boolean): void
}>()

type TimelineTone = 'calm' | 'notice' | 'alert'

type TimelineItem = {
  id: string
  at: string
  title: string
  detail: string
  badge: string
  tone: TimelineTone
  distanceMs: number
  kind: 'success' | 'error' | 'current'
}

type RecoveryGuideStep = {
  index: number
  title: string
  detail: string
  anchorLabel?: string
}

type RecoveryGuide = {
  summary: string
  steps: RecoveryGuideStep[]
  checks: string[]
}

const { t } = useI18n()
const router = useRouter()
const appStore = useAppStore()

const loading = ref(false)
const loadError = ref(false)
const detail = ref<UserErrorRequestDetail | null>(null)
const timelineLoading = ref(false)
const timelineLoadFailed = ref(false)
const timelineItems = ref<TimelineItem[]>([])
const timelineSummary = ref('')

function resolveDiagnosisSummary(reasonCode: string, endpoint: string) {
  if (!reasonCode) return ''
  const key = `usage.errors.detail.reasonExplanations.${reasonCode}.summary`
  const value = t(key, { endpoint })
  return typeof value === 'string' && value !== key ? value : ''
}

function resolveDiagnosisAdvice(reasonCode: string) {
  if (!reasonCode) return [] as string[]
  const key = `usage.errors.detail.reasonExplanations.${reasonCode}.advice`
  const value = t(key)
  if (Array.isArray(value)) return value.filter((item) => typeof item === 'string' && item.trim().length > 0) as string[]
  return []
}

const explanationSummary = computed(() => {
  if (!detail.value) return ''
  const category = detail.value.category
  const endpoint = detail.value.inbound_endpoint || t('usage.errors.detail.thisRequest')
  const diagnosisSummary = resolveDiagnosisSummary(detail.value.diagnosis?.reason_code || '', endpoint)
  if (diagnosisSummary) return diagnosisSummary

  if (category === 'rate_limit') return t('usage.errors.detail.explanations.rate_limit.summary', { endpoint })
  if (category === 'quota') return t('usage.errors.detail.explanations.quota.summary', { endpoint })
  if (category === 'invalid_request') return t('usage.errors.detail.explanations.invalid_request.summary', { endpoint })
  if (category === 'service_unavailable') return t('usage.errors.detail.explanations.service_unavailable.summary', { endpoint })
  if (category === 'upstream') return t('usage.errors.detail.explanations.upstream.summary', { endpoint })
  if (category === 'internal') return t('usage.errors.detail.explanations.internal.summary', { endpoint })
  if (category === 'cyber') return t('usage.errors.detail.explanations.cyber.summary', { endpoint })
  if (category === 'auth') return t('usage.errors.detail.explanations.auth.summary', { endpoint })
  return t('usage.errors.detail.explanations.other.summary', { endpoint })
})

const explanationAdvice = computed(() => {
  if (!detail.value) return [] as string[]
  const diagnosisAdvice = resolveDiagnosisAdvice(detail.value.diagnosis?.reason_code || '')
  if (diagnosisAdvice.length > 0) return diagnosisAdvice
  const category = detail.value.category
  const key = 'usage.errors.detail.explanations.' + category + '.advice'
  const value = t(key)
  if (Array.isArray(value)) return value.filter((item) => typeof item === 'string' && item.trim().length > 0) as string[]
  return []
})

const modelTrace = computed(() => {
  if (!detail.value) return null as null | { requested: string; upstream: string; hint: string }
  const requested = String(detail.value.diagnosis?.requested_model || detail.value.model || '').trim()
  const upstream = String(detail.value.diagnosis?.upstream_model || '').trim()
  if (!requested && !upstream) return null

  let hint = ''
  if (requested && upstream && requested !== upstream) {
    hint = t('usage.errors.detail.modelTrace.mapped', { requested, upstream })
  } else if (requested && ['request_model_not_supported', 'service_model_not_available'].includes(detail.value.diagnosis?.reason_code || '')) {
    hint = t('usage.errors.detail.modelTrace.unavailable', { requested })
  } else if (requested) {
    hint = t('usage.errors.detail.modelTrace.requestedOnly', { requested })
  }

  return {
    requested,
    upstream,
    hint,
  }
})

const diagnosticSummary = computed(() => {
  if (!detail.value) return ''

  const lines = [
    t('usage.errors.detail.summaryLabels.title'),
    `${t('usage.errors.time')}: ${formatDateTime(detail.value.created_at)}`,
    `${t('usage.errors.model')}: ${detail.value.model || '-'}`,
    `${t('usage.errors.endpoint')}: ${detail.value.inbound_endpoint || '-'}`,
    `${t('usage.errors.status')}: ${detail.value.status_code || '-'}`,
    `${t('usage.errors.category')}: ${t('usage.errors.categories.' + detail.value.category)}`,
    `${t('usage.errors.platform')}: ${detail.value.platform || '-'}`,
  ]

  if (detail.value.upstream_status_code != null) {
    lines.push(`${t('usage.errors.detail.upstreamStatus')}: ${detail.value.upstream_status_code}`)
  }

  if (detail.value.message) {
    lines.push(`${t('usage.errors.message')}: ${detail.value.message}`)
  }

  if (modelTrace.value?.requested) {
    lines.push(`${t('usage.errors.detail.modelTrace.requested')}: ${modelTrace.value.requested}`)
  }

  if (modelTrace.value?.upstream) {
    lines.push(`${t('usage.errors.detail.modelTrace.upstream')}: ${modelTrace.value.upstream}`)
  }

  if (explanationSummary.value) {
    lines.push(`${t('usage.errors.detail.summaryLabels.explanation')}: ${explanationSummary.value}`)
  }

  if (explanationAdvice.value.length > 0) {
    lines.push(`${t('usage.errors.detail.summaryLabels.advice')}: ${explanationAdvice.value.join('；')}`)
  }

  if (timelineSummary.value) {
    lines.push(`${t('usage.errors.detail.summaryLabels.timeline')}: ${timelineSummary.value}`)
  }

  if (nextAction.value) {
    lines.push(`${t('usage.errors.detail.summaryLabels.nextAction')}: ${nextAction.value.label}`)
  }

  return lines.join('\n')
})

const nextAction = computed(() => {
  if (!detail.value) return null as null | { to: string; query?: Record<string, string>; label: string }
  switch (detail.value.diagnosis?.action_code) {
    case 'keys_connection_test':
      return { to: '/keys', query: { panel: 'connection-test' }, label: t('usage.errors.detail.actions.auth') }
    case 'profile_balance_notify':
      return { to: '/profile', query: { focus: 'balance-notify' }, label: t('usage.errors.detail.actions.quota') }
    case 'usage_review_payload':
      return { to: '/usage', query: { tab: 'errors', category: 'invalid_request' }, label: t('usage.errors.detail.actions.invalid_request') }
    case 'dashboard_requests_focus':
      return { to: '/dashboard', query: { focus: 'requests' }, label: t('usage.errors.detail.actions.rate_limit') }
    case 'usage_review_cyber':
      return { to: '/usage', query: { tab: 'errors', category: 'cyber' }, label: t('usage.errors.detail.actions.cyber') }
    case 'usage_retry_later':
      return { to: '/usage', query: { tab: 'errors', category: detail.value.category }, label: t('usage.errors.detail.actions.retry_later') }
  }
  const category = detail.value.category
  if (category === 'auth') return { to: '/keys', query: { panel: 'connection-test' }, label: t('usage.errors.detail.actions.auth') }
  if (category === 'quota') return { to: '/profile', query: { focus: 'balance-notify' }, label: t('usage.errors.detail.actions.quota') }
  if (category === 'invalid_request') return { to: '/usage', query: { tab: 'errors', category: 'invalid_request' }, label: t('usage.errors.detail.actions.invalid_request') }
  if (category === 'rate_limit') return { to: '/dashboard', query: { focus: 'requests' }, label: t('usage.errors.detail.actions.rate_limit') }
  if (category === 'service_unavailable' || category === 'upstream' || category === 'internal' || category === 'other') {
    return { to: '/usage', query: { tab: 'errors', category }, label: t('usage.errors.detail.actions.retry_later') }
  }
  if (category === 'cyber') return { to: '/usage', query: { tab: 'errors', category: 'cyber' }, label: t('usage.errors.detail.actions.cyber') }
  return null
})

const recoveryGuide = computed<RecoveryGuide | null>(() => {
  if (!detail.value) return null

  const category = detail.value.category
  const actionLabel = nextAction.value?.label
  const focusLabel = actionLabel ? t('usage.errors.detail.recoveryGuide.anchor', { label: actionLabel }) : ''

  const withAction = (steps: Array<Omit<RecoveryGuideStep, 'index'>>) => steps.map((step, idx) => ({
    index: idx + 1,
    ...step,
  }))

  if (category === 'quota') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.quota'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.quota.checkBalance.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.quota.checkBalance.detail'),
          anchorLabel: focusLabel,
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.quota.confirmRefresh.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.quota.confirmRefresh.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.quota.retest.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.quota.retest.detail'),
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.quota.balanceVisible'),
        t('usage.errors.detail.recoveryGuide.checks.quota.noNewQuotaErrors'),
      ],
    }
  }

  if (category === 'auth') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.auth'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.auth.verifyKey.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.auth.verifyKey.detail'),
          anchorLabel: focusLabel,
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.auth.testConnection.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.auth.testConnection.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.auth.retest.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.auth.retest.detail'),
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.auth.connectionPasses'),
        t('usage.errors.detail.recoveryGuide.checks.auth.noRepeatedAuthErrors'),
      ],
    }
  }

  if (category === 'invalid_request') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.invalid_request'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.invalid_request.reviewPayload.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.invalid_request.reviewPayload.detail'),
          anchorLabel: focusLabel,
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.invalid_request.matchModel.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.invalid_request.matchModel.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.invalid_request.retest.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.invalid_request.retest.detail'),
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.invalid_request.statusNormal'),
        t('usage.errors.detail.recoveryGuide.checks.invalid_request.noSameCategory'),
      ],
    }
  }

  if (category === 'rate_limit') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.rate_limit'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.rate_limit.slowDown.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.rate_limit.slowDown.detail'),
          anchorLabel: focusLabel,
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.rate_limit.observeWindow.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.rate_limit.observeWindow.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.rate_limit.retry.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.rate_limit.retry.detail'),
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.rate_limit.requestPaceDropped'),
        t('usage.errors.detail.recoveryGuide.checks.rate_limit.successReturned'),
      ],
    }
  }

  if (category === 'service_unavailable' || category === 'upstream' || category === 'internal' || category === 'other') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.service'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.service.pauseRetry.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.service.pauseRetry.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.service.observeTimeline.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.service.observeTimeline.detail'),
          anchorLabel: t('usage.errors.detail.recoveryGuide.anchorTimeline'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.service.retryLater.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.service.retryLater.detail'),
          anchorLabel: focusLabel,
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.service.timelineRecovered'),
        t('usage.errors.detail.recoveryGuide.checks.service.errorClusterShrank'),
      ],
    }
  }

  if (category === 'cyber') {
    return {
      summary: t('usage.errors.detail.recoveryGuide.summaries.cyber'),
      steps: withAction([
        {
          title: t('usage.errors.detail.recoveryGuide.steps.cyber.reviewContent.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.cyber.reviewContent.detail'),
          anchorLabel: focusLabel,
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.cyber.adjustRequest.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.cyber.adjustRequest.detail'),
        },
        {
          title: t('usage.errors.detail.recoveryGuide.steps.cyber.retest.title'),
          detail: t('usage.errors.detail.recoveryGuide.steps.cyber.retest.detail'),
        },
      ]),
      checks: [
        t('usage.errors.detail.recoveryGuide.checks.cyber.contentAdjusted'),
        t('usage.errors.detail.recoveryGuide.checks.cyber.categoryCleared'),
      ],
    }
  }

  return null
})

async function copyDiagnosticSummary() {
  if (!diagnosticSummary.value) return
  try {
    await navigator.clipboard.writeText(diagnosticSummary.value)
    appStore.showSuccess(t('usage.errors.detail.copySummarySuccess'))
  } catch (error) {
    console.error('[UserErrorDetailModal] Failed to copy diagnostic summary:', error)
    appStore.showError(t('usage.errors.detail.copySummaryFailed'))
  }
}

async function goToNextAction() {
  if (!nextAction.value) return
  emit('update:show', false)
  await router.push({ path: nextAction.value.to, query: nextAction.value.query })
}

watch(
  () => [props.show, props.errorId] as const,
  ([show, id]) => {
    if (show && id != null) {
      void fetchDetail(id)
    } else if (!show) {
      detail.value = null
      loadError.value = false
      timelineItems.value = []
      timelineSummary.value = ''
      timelineLoadFailed.value = false
    }
  }
)

async function fetchDetail(id: number) {
  loading.value = true
  loadError.value = false
  detail.value = null
  timelineItems.value = []
  timelineSummary.value = ''
  timelineLoadFailed.value = false
  try {
    const loaded = await getMyErrorDetail(id)
    detail.value = loaded
    await fetchTimeline(loaded)
  } catch (e) {
    console.error('[UserErrorDetailModal] Failed to load error detail:', e)
    loadError.value = true
  } finally {
    loading.value = false
  }
}

async function fetchTimeline(currentDetail: UserErrorRequestDetail) {
  timelineLoading.value = true
  timelineLoadFailed.value = false
  try {
    const date = toQueryDate(currentDetail.created_at)
    const [usageRes, errorRes] = await Promise.all([
      queryUsage({
        start_date: date,
        end_date: date,
        page: 1,
        page_size: 60,
        model: currentDetail.model,
        sort_by: 'created_at',
        sort_order: 'desc',
      }),
      listMyErrorRequests({
        start_date: date,
        end_date: date,
        page: 1,
        page_size: 40,
        model: currentDetail.model,
      }),
    ])

    const built = buildTimelineItems(currentDetail, usageRes.items || [], errorRes.items || [])
    timelineItems.value = built.items
    timelineSummary.value = built.summary
  } catch (error) {
    console.error('[UserErrorDetailModal] Failed to build timeline:', error)
    timelineItems.value = []
    timelineSummary.value = ''
    timelineLoadFailed.value = true
  } finally {
    timelineLoading.value = false
  }
}

function buildTimelineItems(currentDetail: UserErrorRequestDetail, usageItems: UsageLog[], errorItems: UserErrorRequest[]) {
  const center = new Date(currentDetail.created_at).getTime()
  const windowMs = 20 * 60 * 1000

  const currentItem: TimelineItem = {
    id: `current-${currentDetail.id}`,
    at: currentDetail.created_at,
    title: t('usage.errors.detail.timeline.currentTitle'),
    detail: currentDetail.message || currentDetail.inbound_endpoint || currentDetail.model || '-',
    badge: t('usage.errors.detail.timeline.badges.current'),
    tone: 'alert',
    distanceMs: 0,
    kind: 'current',
  }

  const successItems = usageItems
    .filter((item) => item.model === currentDetail.model)
    .map((item) => ({
      id: `success-${item.id}`,
      at: item.created_at,
      title: t('usage.errors.detail.timeline.successTitle'),
      detail: `${item.inbound_endpoint || item.model || '-'} · ${formatUsageSummary(item)}`,
      badge: t('usage.errors.detail.timeline.badges.success'),
      tone: 'calm' as TimelineTone,
      distanceMs: Math.abs(new Date(item.created_at).getTime() - center),
      kind: 'success' as const,
    }))

  const siblingErrors = errorItems
    .filter((item) => item.id !== currentDetail.id)
    .map((item) => ({
      id: `error-${item.id}`,
      at: item.created_at,
      title: t('usage.errors.detail.timeline.errorTitle', { category: t('usage.errors.categories.' + item.category) }),
      detail: `${item.inbound_endpoint || item.model || '-'} · ${item.message || t('usage.errors.detail.timeline.badges.error')}`,
      badge: t('usage.errors.detail.timeline.badges.error'),
      tone: 'notice' as TimelineTone,
      distanceMs: Math.abs(new Date(item.created_at).getTime() - center),
      kind: 'error' as const,
    }))

  const nearby = [...successItems, ...siblingErrors]
    .filter((item) => item.distanceMs <= windowMs)
    .sort((a, b) => a.distanceMs - b.distanceMs)

  const selected = (nearby.length ? nearby : [...successItems, ...siblingErrors].sort((a, b) => a.distanceMs - b.distanceMs))
    .slice(0, 4)

  const items = [...selected, currentItem].sort((a, b) => new Date(a.at).getTime() - new Date(b.at).getTime())
  const summary = buildTimelineSummary(center, items)

  return { items, summary }
}

function buildTimelineSummary(center: number, items: TimelineItem[]) {
  const around = items.filter((item) => item.kind !== 'current')
  const successAfter = around.some((item) => item.kind === 'success' && new Date(item.at).getTime() > center && Math.abs(new Date(item.at).getTime() - center) <= 10 * 60 * 1000)
  const successBefore = around.some((item) => item.kind === 'success' && new Date(item.at).getTime() < center && Math.abs(new Date(item.at).getTime() - center) <= 10 * 60 * 1000)
  const nearbyErrors = around.filter((item) => item.kind === 'error' && Math.abs(new Date(item.at).getTime() - center) <= 10 * 60 * 1000).length

  if (successAfter) return t('usage.errors.detail.timeline.summaries.recovered')
  if (nearbyErrors >= 2) return t('usage.errors.detail.timeline.summaries.continuous')
  if (successBefore) return t('usage.errors.detail.timeline.summaries.isolated')
  if (!around.length) return t('usage.errors.detail.timeline.summaries.sparse')
  return t('usage.errors.detail.timeline.summaries.observe')
}

function formatUsageSummary(item: UsageLog) {
  const tokenCount = item.input_tokens + item.output_tokens
  return `${formatTokens(tokenCount)} tokens · $${item.actual_cost.toFixed(4)}`
}

function formatTokens(value: number) {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1000) return `${(value / 1000).toFixed(1)}K`
  return String(value)
}

function toQueryDate(value: string) {
  const date = new Date(value)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function timelineToneClass(tone: TimelineTone) {
  if (tone === 'alert') return 'border-amber-300 bg-amber-50/70 dark:border-amber-800 dark:bg-amber-950/10'
  if (tone === 'notice') return 'border-slate-200 bg-white/80 dark:border-dark-600 dark:bg-dark-800/70'
  return 'border-emerald-200 bg-emerald-50/70 dark:border-emerald-900/60 dark:bg-emerald-950/10'
}

function timelineBadgeClass(tone: TimelineTone) {
  if (tone === 'alert') return 'bg-amber-100 text-amber-900 dark:bg-amber-900/40 dark:text-amber-100'
  if (tone === 'notice') return 'bg-slate-100 text-slate-700 dark:bg-dark-700 dark:text-dark-100'
  return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-100'
}

function statusClass(code: number) {
  if (code >= 500) return 'badge-danger'
  if (code === 429) return 'badge-warning'
  return 'badge-gray'
}
</script>
