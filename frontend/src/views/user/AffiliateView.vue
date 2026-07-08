<template>
  <AppLayout>
    <div class="affiliate-page space-y-6">
      <div v-if="loading" class="affiliate-loading-state flex justify-center py-14">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-zen-seal border-t-transparent"></div>
      </div>

      <section
        v-else-if="detail"
        class="affiliate-hero overflow-hidden rounded-zen border border-zen-paperLine/75 bg-white/45 p-6 shadow-paper dark:border-zen-nightLine dark:bg-zen-nightPanel/72 lg:p-7"
      >
        <div class="affiliate-hero-grid grid gap-6 xl:grid-cols-[minmax(0,1.35fr)_minmax(18rem,0.9fr)]">
          <div class="space-y-6">
            <div>
              <div class="mb-4 flex items-center gap-4">
                <span class="h-px w-14 bg-zen-paperLine dark:bg-zen-nightLine"></span>
                <span class="font-mono text-xs uppercase tracking-[0.34em] text-zen-mist dark:text-zen-stone">{{ affiliateCopy.kicker }}</span>
              </div>
              <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div class="max-w-2xl">
                  <h1 class="font-serif text-3xl font-semibold text-zen-ink dark:text-zen-paper sm:text-4xl">
                    {{ t('affiliate.title') }}
                  </h1>
                  <p class="mt-3 text-sm leading-7 text-zen-mist dark:text-zen-stone sm:text-[0.96rem]">
                    {{ t('affiliate.description') }}
                  </p>
                </div>

                <div class="affiliate-rate-pill self-start lg:self-auto">
                  <span>{{ t('affiliate.stats.rebateRate') }}</span>
                  <strong>{{ formattedRebateRate }}%</strong>
                </div>
              </div>
            </div>

            <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.12fr)]">
              <div class="affiliate-share-block">
                <div class="affiliate-share-copy">
                  <span>{{ affiliateCopy.codeKicker }}</span>
                  <strong>{{ t('affiliate.yourCode') }}</strong>
                </div>
                <div class="affiliate-share-row">
                  <code>{{ detail.aff_code }}</code>
                  <button class="btn btn-secondary btn-sm" @click="copyCode">
                    <Icon name="copy" size="sm" />
                    <span>{{ t('affiliate.copyCode') }}</span>
                  </button>
                </div>
              </div>

              <div class="affiliate-share-block affiliate-share-block-link">
                <div class="affiliate-share-copy">
                  <span>{{ affiliateCopy.linkKicker }}</span>
                  <strong>{{ t('affiliate.inviteLink') }}</strong>
                </div>
                <div class="affiliate-share-row affiliate-share-row-link">
                  <code>{{ inviteLink }}</code>
                  <button class="btn btn-secondary btn-sm" @click="copyInviteLink">
                    <Icon name="copy" size="sm" />
                    <span>{{ t('affiliate.copyLink') }}</span>
                  </button>
                </div>
              </div>
            </div>

            <div class="affiliate-rules">
              <div class="affiliate-rules-head">
                <span class="affiliate-rules-mark">{{ affiliateCopy.rulesMark }}</span>
                <div>
                  <p class="text-sm font-medium text-zen-ink dark:text-zen-paper">{{ t('affiliate.tips.title') }}</p>
                </div>
              </div>
              <ul class="affiliate-rules-list">
                <li>1. {{ t('affiliate.tips.line1') }}</li>
                <li>2. {{ t('affiliate.tips.line2', { rate: `${formattedRebateRate}%` }) }}</li>
                <li>3. {{ t('affiliate.tips.line3') }}</li>
                <li v-if="detail.aff_frozen_quota > 0">4. {{ t('affiliate.tips.line4') }}</li>
              </ul>
            </div>
          </div>

          <aside class="affiliate-summary-panel">
            <div class="affiliate-summary-shell">
              <div class="affiliate-summary-lead">
                <div class="affiliate-summary-card affiliate-summary-card-primary">
                  <span class="affiliate-summary-label">{{ affiliateCopy.availableKicker }}</span>
                  <strong>{{ formatCurrency(detail.aff_quota) }}</strong>
                  <p>
                    {{ detail.aff_quota > 0 ? t('affiliate.transfer.description') : t('affiliate.transfer.empty') }}
                  </p>
                </div>

                <button
                  class="btn btn-primary affiliate-transfer-button"
                  :disabled="transferring || detail.aff_quota <= 0"
                  @click="transferQuota"
                >
                  <Icon v-if="transferring" name="refresh" size="sm" class="animate-spin" />
                  <Icon v-else name="dollar" size="sm" />
                  <span>{{ transferring ? t('affiliate.transfer.transferring') : t('affiliate.transfer.button') }}</span>
                </button>
              </div>

              <div class="affiliate-summary-metrics">
                <div class="affiliate-summary-metric">
                  <span class="affiliate-summary-label">{{ t('affiliate.stats.invitedUsers') }}</span>
                  <strong>{{ formatCount(detail.aff_count) }}</strong>
                  <p>{{ affiliateCopy.invitedUsersHint }}</p>
                </div>
                <div class="affiliate-summary-metric">
                  <span class="affiliate-summary-label">{{ t('affiliate.stats.totalQuota') }}</span>
                  <strong>{{ formatCurrency(detail.aff_history_quota) }}</strong>
                  <p>{{ affiliateCopy.totalQuotaHint }}</p>
                </div>
                <div class="affiliate-summary-metric" :class="{ 'affiliate-summary-metric-accent': detail.aff_frozen_quota > 0 }">
                  <span class="affiliate-summary-label">{{ t('affiliate.stats.frozenQuota') }}</span>
                  <strong>{{ formatCurrency(detail.aff_frozen_quota) }}</strong>
                  <p>
                    {{ detail.aff_frozen_quota > 0 ? t('affiliate.tips.line4') : t('affiliate.stats.frozenQuotaHint') }}
                  </p>
                </div>
              </div>
            </div>
          </aside>
        </div>
      </section>

      <section v-else class="affiliate-empty-state card p-7 text-center sm:p-8">
        <div class="mx-auto flex max-w-md flex-col items-center gap-3">
          <span class="affiliate-empty-kicker">{{ affiliateCopy.kicker }}</span>
          <h2 class="font-serif text-2xl font-semibold text-zen-ink dark:text-zen-paper">{{ affiliateCopy.emptyTitle }}</h2>
          <p class="text-sm leading-7 text-zen-mist dark:text-zen-stone">
            {{ affiliateCopy.emptyCopy }}
          </p>
          <button class="btn btn-secondary" @click="reloadAffiliateDetail">
            <Icon name="refresh" size="sm" />
            <span>{{ t('common.refresh') }}</span>
          </button>
        </div>
      </section>

    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import userAPI from '@/api/user'
import type { UserAffiliateDetail } from '@/types'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useClipboard } from '@/composables/useClipboard'
import { formatCurrency } from '@/utils/format'
import { extractApiErrorMessage } from '@/utils/apiError'

const { t, locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { copyToClipboard } = useClipboard()

const loading = ref(true)
const transferring = ref(false)
const detail = ref<UserAffiliateDetail | null>(null)

const zhAffiliateCopy = {
  kicker: '团队引荐',
  codeKicker: '引荐凭引',
  linkKicker: '入庭路径',
  rulesMark: '则',
  availableKicker: '当前可转',
  invitedUsersHint: '已完成绑定邀请码并注册的用户数量。',
  totalQuotaHint: '截至目前累计形成的全部返利额度。',
  emptyTitle: '邀请返利暂未载入',
  emptyCopy: '这次没有拿到返利信息。你可以重新加载页面；如果问题持续，再检查接口或登录状态。'
}

const enAffiliateCopy = {
  kicker: 'Team referrals',
  codeKicker: 'Referral code',
  linkKicker: 'Invite path',
  rulesMark: 'Rule',
  availableKicker: 'Available',
  invitedUsersHint: 'Users who registered with your invitation code.',
  totalQuotaHint: 'Total rebate quota accumulated so far.',
  emptyTitle: 'Affiliate rebates are not loaded',
  emptyCopy: 'Affiliate data could not be loaded this time. Refresh the page, or check the API and login state if it continues.'
}

const affiliateCopy = computed(() => locale.value === 'zh' ? zhAffiliateCopy : enAffiliateCopy)

const inviteLink = computed(() => {
  if (!detail.value) return ''
  if (typeof window === 'undefined') return `/register?aff=${encodeURIComponent(detail.value.aff_code)}`
  return `${window.location.origin}/register?aff=${encodeURIComponent(detail.value.aff_code)}`
})

// Rebate rate is a percentage in the range [0, 100]; backend already clamps it.
// We trim trailing zeros (e.g. 20.00 → "20", 12.50 → "12.5") for a cleaner UI.
const formattedRebateRate = computed(() => {
  const v = detail.value?.effective_rebate_rate_percent ?? 0
  const rounded = Math.round(v * 100) / 100
  return Number.isInteger(rounded) ? String(rounded) : rounded.toString()
})

function formatCount(value: number): string {
  return value.toLocaleString()
}

async function loadAffiliateDetail(silent = false): Promise<void> {
  if (!silent) {
    loading.value = true
  }
  try {
    detail.value = await userAPI.getAffiliateDetail()
  } catch (error) {
    detail.value = null
    appStore.showError(extractApiErrorMessage(error, t('affiliate.loadFailed')))
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

async function reloadAffiliateDetail(): Promise<void> {
  await loadAffiliateDetail()
}

async function copyCode(): Promise<void> {
  if (!detail.value?.aff_code) return
  await copyToClipboard(detail.value.aff_code, t('affiliate.codeCopied'))
}

async function copyInviteLink(): Promise<void> {
  if (!inviteLink.value) return
  await copyToClipboard(inviteLink.value, t('affiliate.linkCopied'))
}

async function transferQuota(): Promise<void> {
  if (!detail.value || detail.value.aff_quota <= 0 || transferring.value) return
  transferring.value = true
  try {
    const resp = await userAPI.transferAffiliateQuota()
    appStore.showSuccess(t('affiliate.transfer.success', { amount: formatCurrency(resp.transferred_quota) }))
    await Promise.all([
      loadAffiliateDetail(true),
      authStore.refreshUser().catch(() => undefined),
    ])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.transferFailed')))
  } finally {
    transferring.value = false
  }
}

onMounted(() => {
  void loadAffiliateDetail()
})
</script>

<style scoped>
.affiliate-page {
  color: #38413a;
}

.affiliate-hero {
  position: relative;
  isolation: isolate;
}

.affiliate-hero::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 16% 18%, rgba(167, 58, 42, 0.08), transparent 30%),
    radial-gradient(circle at 88% 12%, rgba(128, 115, 92, 0.08), transparent 22%),
    linear-gradient(180deg, rgba(255, 252, 246, 0.3), transparent 55%);
  pointer-events: none;
  z-index: -1;
}

.affiliate-rate-pill {
  display: inline-grid;
  gap: 0.28rem;
  min-width: 11rem;
  padding: 0.85rem 1rem;
  border: 1px solid rgba(167, 58, 42, 0.16);
  border-radius: 6px;
  background: linear-gradient(180deg, rgba(255, 250, 245, 0.94), rgba(246, 238, 227, 0.78));
  box-shadow: 0 18px 38px -34px rgba(86, 66, 44, 0.42);
  text-align: left;
}

.affiliate-rate-pill span,
.affiliate-empty-kicker,
.affiliate-summary-label,
.affiliate-rules-mark {
  color: #766148;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.affiliate-rate-pill strong,
.affiliate-empty-state h2 {
  color: #1f2320;
}

.affiliate-rate-pill strong {
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.7rem;
  font-weight: 600;
  line-height: 1.15;
}

.affiliate-share-block,
.affiliate-rules,
.affiliate-summary-shell,
.affiliate-summary-card {
  border: 1px solid rgba(216, 205, 185, 0.76);
  border-radius: 6px;
  background: rgba(250, 247, 239, 0.72);
  box-shadow: 0 18px 42px -34px rgba(31, 35, 32, 0.28);
}

.affiliate-share-block,
.affiliate-rules,
.affiliate-summary-card,
.affiliate-summary-shell {
  padding: 1rem;
}

.affiliate-share-copy {
  display: grid;
  gap: 0.24rem;
  margin-bottom: 0.8rem;
}

.affiliate-share-copy strong,
.affiliate-summary-card strong,
.affiliate-summary-metric strong {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
}

.affiliate-share-copy strong {
  font-size: 1rem;
  font-weight: 600;
}

.affiliate-share-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 0;
  padding: 0.72rem 0.8rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 6px;
  background: rgba(255, 252, 246, 0.78);
}

.affiliate-share-row code {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  color: #1f2320;
  font-size: 0.84rem;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.affiliate-share-row-link code {
  color: #4f5b50;
  font-weight: 500;
}

.affiliate-rules {
  display: grid;
  gap: 0.95rem;
}

.affiliate-rules-head {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 0.9rem;
  align-items: start;
}

.affiliate-rules-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2.2rem;
  height: 2.2rem;
  border: 1px solid rgba(167, 58, 42, 0.16);
  border-radius: 999px;
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
  letter-spacing: 0.22em;
}

.affiliate-rules-list {
  display: grid;
  gap: 0.55rem;
  color: #4a564b;
  font-size: 0.9rem;
  line-height: 1.9;
}

.affiliate-summary-panel {
  align-content: start;
}

.affiliate-summary-shell {
  display: grid;
  gap: 1.1rem;
  padding: 1.15rem;
  background:
    linear-gradient(180deg, rgba(255, 250, 245, 0.95), rgba(244, 236, 224, 0.78)),
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.08), transparent 32%);
}

.affiliate-summary-lead {
  display: grid;
  gap: 0.9rem;
}

.affiliate-summary-card {
  display: grid;
  gap: 0.45rem;
  padding: 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

.affiliate-summary-card strong {
  font-size: 2rem;
  font-weight: 600;
  line-height: 1.15;
}

.affiliate-summary-card p,
.affiliate-empty-state p,
.affiliate-summary-metric p {
  color: #526053;
  font-size: 0.86rem;
  line-height: 1.8;
}

.affiliate-summary-card-primary {
  gap: 0.5rem;
}

.affiliate-summary-card-primary strong {
  color: #2e7161;
}

.affiliate-summary-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.7rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(167, 58, 42, 0.14);
}

.affiliate-summary-metric {
  display: grid;
  gap: 0.38rem;
  min-width: 0;
  padding: 0.9rem 0.85rem 0;
  border-left: 1px solid rgba(216, 205, 185, 0.7);
}

.affiliate-summary-metric:first-child {
  padding-left: 0;
  border-left: 0;
}

.affiliate-summary-metric strong {
  font-size: 1.42rem;
  font-weight: 600;
  line-height: 1.15;
}

.affiliate-summary-metric-accent strong {
  color: #9a6a21;
}

.affiliate-transfer-button {
  width: 100%;
  margin-top: 0.1rem;
  justify-content: center;
  min-height: 2.7rem;
}

.affiliate-empty-state {
  display: grid;
  gap: 0.38rem;
}

.affiliate-empty-kicker {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.affiliate-empty-state {
  border-color: rgba(216, 205, 185, 0.82);
  background: rgba(250, 247, 239, 0.54);
}

.affiliate-empty-kicker {
  color: #a73a2a;
  letter-spacing: 0.24em;
}

.affiliate-page :deep(.btn-secondary) {
  border-color: rgba(216, 205, 185, 0.92);
  background: rgba(255, 255, 255, 0.44);
  color: #38413a;
}

.affiliate-page :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(255, 255, 255, 0.72);
}

.affiliate-page :deep(.btn-primary) {
  border-color: #a73a2a;
  background-color: #a73a2a;
  color: #f4efe4;
}

.affiliate-page :deep(.btn-primary:hover) {
  border-color: #8f3024;
  background-color: #8f3024;
}

.affiliate-page :deep(.btn-primary:disabled) {
  cursor: not-allowed;
  border-color: rgba(216, 205, 185, 0.86);
  background: rgba(250, 247, 239, 0.78);
  color: #38413a;
  opacity: 1;
}

.affiliate-page :deep(.card) {
  border-color: rgba(216, 205, 185, 0.78);
  border-radius: 6px;
  background: rgba(250, 247, 239, 0.62);
  box-shadow: 0 14px 38px -32px rgba(31, 35, 32, 0.26);
}

.affiliate-page :deep(.rounded),
.affiliate-page :deep(.rounded-lg),
.affiliate-page :deep(.rounded-xl),
.affiliate-page :deep(.rounded-2xl) {
  border-radius: 6px;
}

.dark .affiliate-page {
  color: #f4efe4;
}

.dark .affiliate-hero {
  border-color: rgba(48, 52, 43, 0.95) !important;
  background: rgba(24, 26, 21, 0.72) !important;
  box-shadow: 0 14px 38px -32px rgba(0, 0, 0, 0.44);
}

.dark .affiliate-hero::before {
  background:
    radial-gradient(circle at 16% 18%, rgba(167, 58, 42, 0.08), transparent 30%),
    radial-gradient(circle at 88% 12%, rgba(184, 156, 116, 0.04), transparent 22%),
    linear-gradient(180deg, rgba(17, 19, 15, 0.08), transparent 55%);
}

.dark .affiliate-rate-pill,
.dark .affiliate-share-block,
.dark .affiliate-rules,
.dark .affiliate-summary-shell,
.dark .affiliate-empty-state,
.dark .affiliate-page :deep(.card) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.74);
}

.dark .affiliate-summary-shell {
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.82), rgba(17, 19, 15, 0.78)),
    radial-gradient(circle at top right, rgba(167, 58, 42, 0.08), transparent 34%);
}

.dark .affiliate-rate-pill {
  background: rgba(24, 26, 21, 0.72);
}

.dark .affiliate-share-row {
  border-color: rgba(48, 52, 43, 0.92);
  background: rgba(17, 19, 15, 0.28);
}

.dark .affiliate-rate-pill strong,
.dark .affiliate-share-copy strong,
.dark .affiliate-summary-card strong,
.dark .affiliate-empty-state h2,
.dark .affiliate-share-row code,
.dark .affiliate-summary-metric strong {
  color: #f4efe4;
}

.dark .affiliate-rate-pill span,
.dark .affiliate-empty-kicker,
.dark .affiliate-summary-label,
.dark .affiliate-rules-mark {
  color: #879186;
}

.dark .affiliate-rules-mark {
  border-color: rgba(167, 58, 42, 0.22);
  background: rgba(167, 58, 42, 0.1);
  color: #f0b4a8;
}

.dark .affiliate-summary-card-primary {
  padding: 0.95rem 1rem;
  border: 1px solid rgba(167, 58, 42, 0.36);
  border-radius: 6px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 38%),
    rgba(24, 26, 21, 0.78);
}

.dark .affiliate-summary-card-primary strong {
  color: #84c4b4;
}

.dark .affiliate-summary-metric-accent strong {
  color: #e2bc77;
}

.dark .affiliate-summary-card p,
.dark .affiliate-summary-metric p,
.dark .affiliate-empty-state p,
.dark .affiliate-rules-list,
.dark .affiliate-share-row-link code {
  color: #d8cdb9;
}

.dark .affiliate-summary-metrics {
  border-top-color: rgba(167, 58, 42, 0.2);
}

.dark .affiliate-summary-metric {
  border-left-color: rgba(48, 52, 43, 0.95);
}

.dark .affiliate-page :deep(.btn-secondary) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
}

.dark .affiliate-page :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  color: #f0b4a8;
}

.dark .affiliate-page :deep(.btn-primary:disabled) {
  border-color: rgba(216, 205, 185, 0.24);
  background: rgba(216, 205, 185, 0.16);
  color: #f4efe4;
}

@media (max-width: 768px) {
  .affiliate-page {
    gap: 1.1rem;
  }

  .affiliate-hero,
  .affiliate-empty-state {
    padding-left: 1rem !important;
    padding-right: 1rem !important;
  }

  .affiliate-share-row {
    flex-direction: column;
    align-items: stretch;
  }

  .affiliate-summary-shell {
    gap: 1rem;
    padding: 1rem;
  }

  .affiliate-summary-metrics {
    grid-template-columns: 1fr;
    gap: 0;
    padding-top: 0.85rem;
  }

  .affiliate-summary-metric {
    padding: 0.8rem 0 0;
    border-top: 1px solid rgba(216, 205, 185, 0.58);
    border-left: 0;
  }

  .affiliate-summary-metric:first-child {
    padding-top: 0;
    border-top: 0;
  }

  .affiliate-share-row .btn,
  .affiliate-transfer-button,
  .affiliate-empty-state .btn {
    width: 100%;
    justify-content: center;
  }

  .dark .affiliate-summary-metric {
    border-top-color: rgba(72, 76, 65, 0.68);
  }
}
</style>
