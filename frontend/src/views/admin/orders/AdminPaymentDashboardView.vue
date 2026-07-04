<template>
  <AppLayout>
    <div class="sst-admin-page">
      <div class="sst-payment-dashboard space-y-6">
        <section class="sst-payment-toolbar">
          <div class="sst-payment-toolbar-copy">
            <span class="sst-payment-toolbar-label">观测区间</span>
            <p>按近 7 / 30 / 90 日切换营收与支付分布，快速比对近期波动与高价值用户。</p>
          </div>
          <div class="sst-payment-toolbar-actions">
            <div class="sst-payment-range-switcher">
              <button
                v-for="d in DAYS_OPTIONS"
                :key="d"
                type="button"
                class="sst-payment-range-button"
                :class="days === d
                  ? 'sst-payment-range-button-active'
                  : 'sst-payment-range-button-idle'"
                @click="days = d"
              >
                {{ d }}{{ t('payment.admin.daySuffix') }}
              </button>
            </div>
            <button @click="loadDashboard" :disabled="loading" class="btn btn-secondary inline-flex items-center gap-2" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
              <span class="hidden sm:inline">{{ t('common.refresh') }}</span>
            </button>
          </div>
        </section>

      <!-- Dashboard Content -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>
      <template v-else-if="stats">
        <OrderStatsCards :stats="stats" />
        <DailyRevenueChart :data="stats.daily_series || []" :loading="loading" />
          <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <section class="card sst-payment-panel p-5">
              <div class="sst-payment-panel-heading">
                <h3>{{ t('payment.admin.paymentDistribution') }}</h3>
                <p>观察支付方式金额占比与单量分布，辅助对账与通道判断。</p>
              </div>
              <EmptyState
                v-if="!stats.payment_methods?.length"
                class="sst-payment-panel-empty"
                :title="t('payment.admin.noData')"
                description="当前区间暂无支付方式分布记录，可切换时间范围后再观察。"
              />
              <div v-else class="space-y-3">
                <div v-for="method in stats.payment_methods" :key="method.type" class="sst-payment-list-item">
                  <div class="flex items-center gap-2">
                    <span :class="['inline-block h-3 w-3 rounded-full', methodColor(method.type)]"></span>
                    <span class="text-sm text-gray-700 dark:text-gray-300">{{ t('payment.methods.' + method.type, method.type) }}</span>
                  </div>
                  <div class="text-right">
                    <span class="text-sm font-medium text-gray-900 dark:text-white">&yen;{{ method.amount.toFixed(2) }}</span>
                    <span class="ml-2 text-xs text-gray-500 dark:text-gray-400">({{ method.count }})</span>
                  </div>
                </div>
              </div>
            </section>
            <section class="card sst-payment-panel p-5">
              <div class="sst-payment-panel-heading">
                <h3>{{ t('payment.admin.topUsers') }}</h3>
                <p>快速识别近期高价值用户，便于运营跟进与异常复核。</p>
              </div>
              <EmptyState
                v-if="!stats.top_users?.length"
                class="sst-payment-panel-empty"
                :title="t('payment.admin.noData')"
                description="当前区间暂无高价值用户记录，可扩大观察区间后再查看。"
              />
              <div v-else class="space-y-2">
                <div v-for="(user, idx) in stats.top_users" :key="user.user_id" class="sst-payment-user-row">
                  <div class="flex items-center gap-3">
                    <span :class="['flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold', rankClass(idx)]">{{ idx + 1 }}</span>
                    <span class="text-sm text-gray-700 dark:text-gray-300">{{ user.email }}</span>
                  </div>
                  <span class="text-sm font-medium text-gray-900 dark:text-white">&yen;{{ user.amount.toFixed(2) }}</span>
                </div>
              </div>
            </section>
          </div>
      </template>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { adminPaymentAPI } from '@/api/admin/payment'
import { extractI18nErrorMessage } from '@/utils/apiError'
import type { DashboardStats } from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import OrderStatsCards from '@/components/admin/payment/OrderStatsCards.vue'
import DailyRevenueChart from '@/components/admin/payment/DailyRevenueChart.vue'

const { t } = useI18n()
const appStore = useAppStore()

const DAYS_OPTIONS = [7, 30, 90] as const
const days = ref<number>(30)
const loading = ref(false)
const stats = ref<DashboardStats | null>(null)

function methodColor(type: string): string {
  const c: Record<string, string> = {
    alipay: 'bg-blue-500', wxpay: 'bg-green-500',
    alipay_direct: 'bg-blue-400', wxpay_direct: 'bg-green-400',
    stripe: 'bg-purple-500',
  }
  return c[type] || 'bg-gray-400'
}

function rankClass(idx: number): string {
  if (idx === 0) return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
  if (idx === 1) return 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
  if (idx === 2) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  return 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400'
}

async function loadDashboard() {
  loading.value = true
  try {
    const res = await adminPaymentAPI.getDashboard(days.value)
    stats.value = res.data
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally {
    loading.value = false
  }
}

watch(days, () => loadDashboard())
onMounted(() => loadDashboard())
</script>

<style scoped>
.sst-payment-toolbar {
  @apply flex flex-col gap-4 rounded-2xl border px-4 py-4 sm:px-5;
  border-color: rgba(198, 184, 157, 0.5);
  background:
    linear-gradient(135deg, rgba(167, 58, 42, 0.05), rgba(255, 255, 255, 0) 38%),
    rgba(250, 247, 239, 0.58);
  box-shadow: 0 18px 38px -34px rgba(58, 48, 34, 0.34);
}

.sst-payment-toolbar-copy {
  @apply space-y-2;
}

.sst-payment-toolbar-label {
  @apply inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-medium tracking-[0.18em] uppercase;
  color: #8b5e3c;
  background: rgba(167, 58, 42, 0.08);
}

.sst-payment-toolbar-copy p {
  @apply text-sm leading-6;
  color: #5f6257;
}

.sst-payment-toolbar-actions {
  @apply flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between;
}

.sst-payment-range-switcher {
  @apply inline-flex w-full flex-wrap rounded-xl border p-1 sm:w-auto;
  border-color: rgba(198, 184, 157, 0.52);
  background: rgba(255, 252, 246, 0.82);
}

.sst-payment-range-button {
  @apply rounded-lg px-3 py-1.5 text-xs font-medium transition-colors;
}

.sst-payment-range-button-active {
  background: #a73a2a;
  color: #fff7f2;
}

.sst-payment-range-button-idle {
  color: #5d655b;
}

.sst-payment-range-button-idle:hover {
  background: rgba(167, 58, 42, 0.08);
}

.sst-payment-panel {
  border-color: rgba(198, 184, 157, 0.44);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.66), rgba(248, 243, 231, 0.9));
  box-shadow: 0 24px 44px -40px rgba(58, 48, 34, 0.38);
}

.sst-payment-panel-heading {
  @apply mb-4 space-y-1;
}

.sst-payment-panel-heading h3 {
  @apply text-sm font-semibold text-gray-900 dark:text-white;
}

.sst-payment-panel-heading p {
  @apply text-xs leading-5 text-gray-500 dark:text-gray-400;
}

.sst-payment-panel-empty {
  @apply min-h-[12rem] px-4 py-5;
}

.sst-payment-panel-empty:deep(.empty-state) {
  @apply h-full justify-center;
}

.sst-payment-panel-empty:deep(.empty-state-title) {
  @apply text-base;
}

.sst-payment-panel-empty:deep(.empty-state-description) {
  @apply text-sm leading-6;
}

.sst-payment-list-item {
  @apply flex items-center justify-between rounded-xl px-3 py-2.5;
  background: rgba(255, 255, 255, 0.54);
}

.sst-payment-user-row {
  @apply flex items-center justify-between rounded-xl px-3 py-2.5 transition-colors;
  background: rgba(255, 255, 255, 0.5);
}

.sst-payment-user-row:hover {
  background: rgba(167, 58, 42, 0.06);
}

</style>
<style>
.dark .sst-payment-toolbar,
.dark .sst-payment-panel {
  border-color: rgba(58, 61, 54, 0.96);
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.9), rgba(18, 20, 16, 0.94));
}

.dark .sst-payment-toolbar-label {
  color: #e7b58e;
  background: rgba(167, 58, 42, 0.22);
}

.dark .sst-payment-toolbar-copy p {
  color: #9ea49a;
}

.dark .sst-payment-range-switcher {
  border-color: rgba(58, 61, 54, 0.96);
  background: rgba(16, 18, 14, 0.82);
}

.dark .sst-payment-range-button-idle {
  color: #cfc6b4;
}

.dark .sst-payment-range-button-idle:hover,
.dark .sst-payment-user-row:hover {
  background: rgba(167, 58, 42, 0.16);
}

.dark .sst-payment-list-item,
.dark .sst-payment-user-row {
  background: rgba(29, 32, 27, 0.75);
}
</style>

