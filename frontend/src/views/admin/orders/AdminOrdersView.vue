<template>
  <AppLayout>
    <div class="admin-orders-shell space-y-4">
      <section class="sst-payment-toolbar">
        <div class="sst-payment-toolbar-copy">
          <span class="sst-payment-toolbar-label">交易概览</span>
          <p>按近 7 / 30 / 90 日查看营收、支付分布和高价值用户，下面保留完整订单明细。</p>
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
          <button @click="loadDashboard" :disabled="dashboardLoading" class="btn btn-secondary inline-flex items-center gap-2" :title="t('common.refresh')">
            <Icon name="refresh" size="md" :class="dashboardLoading ? 'animate-spin' : ''" />
            <span class="hidden sm:inline">{{ t('common.refresh') }}</span>
          </button>
        </div>
      </section>

      <div v-if="dashboardLoading" class="flex items-center justify-center py-8">
        <LoadingSpinner />
      </div>
      <template v-else-if="stats">
        <OrderStatsCards :stats="stats" />
        <DailyRevenueChart :data="stats.daily_series || []" :loading="dashboardLoading" />
        <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
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

      <!-- Filters -->
      <div class="card p-4">
        <div class="flex flex-wrap items-center gap-3">
          <div class="flex-1 sm:max-w-64">
            <input v-model="orderSearch" type="text" :placeholder="t('payment.admin.searchOrders')" class="input" @input="debounceLoadOrders" />
          </div>
          <Select v-model="orderFilters.status" :options="statusFilterOptions" class="w-36" @change="loadOrders" />
          <Select v-model="orderFilters.payment_type" :options="paymentTypeFilterOptions" class="w-40" @change="loadOrders" />
          <Select v-model="orderFilters.order_type" :options="orderTypeFilterOptions" class="w-36" @change="loadOrders" />
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <button @click="loadOrders" :disabled="ordersLoading" class="btn btn-secondary" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="ordersLoading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <OrderTable :orders="orders" :loading="ordersLoading" show-user>
        <template #actions="{ row }">
          <div class="flex items-center gap-1">
            <button @click="showOrderDetail(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-600">
              <Icon name="eye" size="sm" />
              {{ t('common.view') }}
            </button>
            <button v-if="row.status === 'PENDING'" @click="handleCancelOrder(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-yellow-600 hover:bg-yellow-50 dark:text-yellow-400 dark:hover:bg-yellow-900/20">
              <Icon name="x" size="sm" />
              {{ t('payment.orders.cancel') }}
            </button>
            <button v-if="row.status === 'FAILED'" @click="handleRetryOrder(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-blue-900/20">
              <Icon name="refresh" size="sm" />
              {{ t('payment.admin.retry') }}
            </button>
            <template v-if="row.status === 'REFUND_REQUESTED'">
              <span v-if="row.refund_amount" class="rounded-full bg-purple-100 px-1.5 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/30 dark:text-purple-300">{{ row.order_type === 'balance' ? '$' : '¥' }}{{ row.refund_amount.toFixed(2) }}</span>
              <button @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-purple-600 hover:bg-purple-50 dark:text-purple-400 dark:hover:bg-purple-900/20">
                <Icon name="check" size="sm" />
                {{ t('payment.admin.approveRefund') }}
              </button>
            </template>
            <button v-else-if="row.status === 'REFUND_FAILED'" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-purple-600 hover:bg-purple-50 dark:text-purple-400 dark:hover:bg-purple-900/20">
              <Icon name="refresh" size="sm" />
              {{ t('payment.admin.retryRefund') }}
            </button>
            <button v-else-if="row.status === 'COMPLETED' || row.status === 'PARTIALLY_REFUNDED'" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20">
              <Icon name="dollar" size="sm" />
              {{ t('payment.admin.refund') }}
            </button>
          </div>
        </template>
      </OrderTable>
      <Pagination v-if="orderPagination.total > 0" :page="orderPagination.page" :total="orderPagination.total" :page-size="orderPagination.page_size" @update:page="handleOrderPageChange" @update:pageSize="handleOrderPageSizeChange" />
    </div>

    <!-- Order Detail Dialog -->
    <BaseDialog :show="showDetailDialog" :title="t('payment.admin.orderDetail')" width="wide" @close="showDetailDialog = false">
      <div v-if="selectedOrder" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') }}</p><p class="font-mono text-sm font-medium text-gray-900 dark:text-white">#{{ selectedOrder.id }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderNo') }}</p><p class="text-sm font-medium text-gray-900 dark:text-white">{{ selectedOrder.out_trade_no }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.status') }}</p><OrderStatusBadge :status="selectedOrder.status" /></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.amount') }}</p><p class="text-sm font-medium text-gray-900 dark:text-white">{{ selectedOrder.order_type === 'balance' ? '$' : '¥' }}{{ selectedOrder.amount.toFixed(2) }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') }}</p><p class="text-sm font-medium text-gray-900 dark:text-white">¥{{ selectedOrder.pay_amount.toFixed(2) }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ t('payment.methods.' + selectedOrder.payment_type, selectedOrder.payment_type) }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.feeRate') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.fee_rate }}%</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.createdAt') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.created_at) }}</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.expiresAt') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.expires_at) }}</p></div>
          <div v-if="selectedOrder.paid_at"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.paidAt') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.paid_at) }}</p></div>
          <div v-if="selectedOrder.refund_amount"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundAmount') }}</p><p class="text-sm font-medium text-red-600 dark:text-red-400">{{ selectedOrder.order_type === 'balance' ? '$' : '¥' }}{{ selectedOrder.refund_amount.toFixed(2) }}</p></div>
          <div v-if="selectedOrder.refund_reason" class="col-span-2"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundReason') }}</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.refund_reason }}</p></div>
          <!-- Refund request info -->
          <div v-if="selectedOrder.refund_requested_at" class="col-span-2 border-t border-gray-200 pt-3 dark:border-dark-600">
            <p class="mb-2 text-xs font-medium text-purple-600 dark:text-purple-400">{{ t('payment.admin.refundRequestInfo') }}</p>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestedAt') }}</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.refund_requested_at) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestedBy') }}</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">#{{ selectedOrder.refund_requested_by }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestReason') }}</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.refund_request_reason }}</p>
              </div>
            </div>
          </div>
        </div>
        <!-- Audit Logs -->
        <div v-if="orderAuditLogs.length > 0" class="border-t border-gray-200 pt-4 dark:border-dark-600">
          <p class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.auditLogs') }}</p>
          <div class="max-h-48 space-y-2 overflow-y-auto">
            <div v-for="log in orderAuditLogs" :key="log.id" class="rounded-lg border border-gray-100 bg-gray-50 p-2.5 dark:border-dark-600 dark:bg-dark-800">
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ log.action }}</span>
                <span class="text-xs text-gray-400">{{ formatDateTime(log.created_at) }}</span>
              </div>
              <div v-if="log.detail" class="mt-1 break-all text-xs text-gray-500 dark:text-gray-400">{{ log.detail }}</div>
              <div v-if="log.operator" class="mt-1 text-xs text-gray-400">{{ t('payment.admin.operator') }}: {{ log.operator }}</div>
            </div>
          </div>
        </div>
      </div>
    </BaseDialog>

    <AdminRefundDialog :show="showRefundDialog" :order="selectedOrder" :submitting="refundSubmitting" @confirm="handleRefund" @cancel="showRefundDialog = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { adminPaymentAPI } from '@/api/admin/payment'
import { extractI18nErrorMessage } from '@/utils/apiError'
import { formatOrderDateTime } from '@/components/payment/orderUtils'
import type { PaymentOrder } from '@/types/payment'
import type { DashboardStats } from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import AdminRefundDialog from '@/components/admin/payment/AdminRefundDialog.vue'
import DailyRevenueChart from '@/components/admin/payment/DailyRevenueChart.vue'
import OrderStatsCards from '@/components/admin/payment/OrderStatsCards.vue'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import OrderTable from '@/components/payment/OrderTable.vue'

interface AuditLog {
  id: number
  action: string
  detail: string | null
  operator: string | null
  created_at: string
}

const { t } = useI18n()
const appStore = useAppStore()

const DAYS_OPTIONS = [7, 30, 90] as const
const days = ref<number>(30)
const dashboardLoading = ref(false)
const stats = ref<DashboardStats | null>(null)
const ordersLoading = ref(false)
const orders = ref<PaymentOrder[]>([])
const orderSearch = ref('')
const orderFilters = reactive({ status: '', payment_type: '', order_type: '' })
const orderPagination = reactive({ page: 1, page_size: 20, total: 0 })
const selectedOrder = ref<PaymentOrder | null>(null)
const showDetailDialog = ref(false)
const showRefundDialog = ref(false)
const refundSubmitting = ref(false)
const orderAuditLogs = ref<AuditLog[]>([])

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function debounceLoadOrders() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => loadOrders(), 300)
}

async function loadOrders() {
  ordersLoading.value = true
  try {
    const res = await adminPaymentAPI.getOrders({
      page: orderPagination.page, page_size: orderPagination.page_size,
      keyword: orderSearch.value || undefined, status: orderFilters.status || undefined,
      payment_type: orderFilters.payment_type || undefined, order_type: orderFilters.order_type || undefined,
    })
    orders.value = res.data.items || []
    orderPagination.total = res.data.total || 0
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally { ordersLoading.value = false }
}

async function loadDashboard() {
  dashboardLoading.value = true
  try {
    const res = await adminPaymentAPI.getDashboard(days.value)
    stats.value = res.data
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally {
    dashboardLoading.value = false
  }
}

function methodColor(type: string): string {
  const colors: Record<string, string> = {
    alipay: 'bg-blue-500',
    wxpay: 'bg-green-500',
    alipay_direct: 'bg-blue-400',
    wxpay_direct: 'bg-green-400',
    stripe: 'bg-purple-500',
    airwallex: 'bg-cyan-500',
  }
  return colors[type] || 'bg-gray-400'
}

function rankClass(idx: number): string {
  if (idx === 0) return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
  if (idx === 1) return 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300'
  if (idx === 2) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
  return 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400'
}

function handleOrderPageChange(page: number) { orderPagination.page = page; loadOrders() }
function handleOrderPageSizeChange(size: number) { orderPagination.page_size = size; orderPagination.page = 1; loadOrders() }

const statusFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allStatuses') },
  { value: 'PENDING', label: t('payment.status.pending') },
  { value: 'PAID', label: t('payment.status.paid') },
  { value: 'COMPLETED', label: t('payment.status.completed') },
  { value: 'EXPIRED', label: t('payment.status.expired') },
  { value: 'CANCELLED', label: t('payment.status.cancelled') },
  { value: 'FAILED', label: t('payment.status.failed') },
  { value: 'REFUNDED', label: t('payment.status.refunded') },
  { value: 'REFUND_REQUESTED', label: t('payment.status.refund_requested') },
  { value: 'REFUND_FAILED', label: t('payment.status.refund_failed') },
])

const paymentTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allPaymentTypes') },
  { value: 'alipay', label: t('payment.methods.alipay') },
  { value: 'wxpay', label: t('payment.methods.wxpay') },
  { value: 'stripe', label: t('payment.methods.stripe') },
  { value: 'airwallex', label: t('payment.methods.airwallex') },
])

const orderTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allOrderTypes') },
  { value: 'balance', label: t('payment.admin.balanceOrder') },
  { value: 'subscription', label: t('payment.admin.subscriptionOrder') },
])

async function showOrderDetail(order: PaymentOrder) {
  selectedOrder.value = order
  orderAuditLogs.value = []
  showDetailDialog.value = true
  try {
    const res = await adminPaymentAPI.getOrder(order.id)
    const data = res.data as unknown as Record<string, unknown>
    if (data.order) selectedOrder.value = data.order as PaymentOrder
    orderAuditLogs.value = ((data.auditLogs || data.audit_logs || []) as unknown) as AuditLog[]
  } catch (_err: unknown) { /* keep cached order data */ }
}

async function handleCancelOrder(order: PaymentOrder) {
  try { await adminPaymentAPI.cancelOrder(order.id); appStore.showSuccess(t('payment.admin.orderCancelled')); loadOrders() }
  catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) }
}

async function handleRetryOrder(order: PaymentOrder) {
  try { await adminPaymentAPI.retryRecharge(order.id); appStore.showSuccess(t('payment.admin.retrySuccess')); loadOrders() }
  catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) }
}

function openRefundDialog(order: PaymentOrder) { selectedOrder.value = order; showRefundDialog.value = true }

async function handleRefund(data: { amount: number; reason: string; deduct_balance: boolean; force: boolean }) {
  if (!selectedOrder.value) return
  refundSubmitting.value = true
  try {
    await adminPaymentAPI.refundOrder(selectedOrder.value.id, { amount: data.amount, reason: data.reason, deduct_balance: data.deduct_balance, force: data.force })
    appStore.showSuccess(t('payment.admin.refundSuccess')); showRefundDialog.value = false; loadOrders()
  } catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) }
  finally { refundSubmitting.value = false }
}

function formatDateTime(dateStr: string): string { return formatOrderDateTime(dateStr) }

watch(days, () => loadDashboard())

onMounted(() => {
  loadOrders()
  loadDashboard()
})
</script>

<style scoped>
.admin-orders-shell {
  padding: clamp(0.2rem, 0.6vw, 0.4rem);
}

.admin-page-hero {
  display: flex;
  flex-wrap: wrap;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
  padding: 1.15rem 1.25rem 1.05rem;
  border: 1px solid rgba(198, 184, 157, 0.52);
  border-radius: 12px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.055), transparent 30%),
    linear-gradient(180deg, rgba(255, 252, 245, 0.9), rgba(246, 241, 231, 0.78));
}

.admin-page-kicker {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.admin-page-hero h2 {
  margin-top: 0.45rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.4rem, 1.8vw, 1.82rem);
  font-weight: 600;
}

.admin-page-hero p {
  max-width: 38rem;
  margin-top: 0.5rem;
  color: #5f675d;
  font-size: 0.94rem;
  line-height: 1.7;
}

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
.dark .admin-page-hero {
  border-color: rgba(48, 52, 43, 0.92);
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 32%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.92), rgba(17, 19, 15, 0.84));
}

.dark .admin-page-hero h2 {
  color: #f4efe4;
}

.dark .admin-page-hero p {
  color: #bdb5a8;
}

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

