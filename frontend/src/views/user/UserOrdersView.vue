<template>
  <AppLayout>
    <div class="orders-shell mx-auto max-w-[76rem] space-y-4 sm:space-y-5">
      <section class="orders-brief card p-6 lg:p-7">
        <div class="orders-brief-copy">
          <span>{{ ordersCopy.kicker }}</span>
          <h2>{{ ordersCopy.title }}</h2>
          <p>{{ ordersCopy.copy }}</p>
        </div>
        <button class="btn btn-primary inline-flex items-center gap-2" @click="router.push('/purchase')">
          <Icon name="wallet" size="sm" />
          <span>{{ t('payment.goToRedeem') }}</span>
        </button>
      </section>

      <!-- Filters -->
      <div class="orders-filter card p-5">
        <div class="flex flex-wrap items-center gap-3">
          <div class="orders-filter-copy">
            <span>{{ t('common.filter') }}</span>
            <strong>{{ currentFilterLabel }}</strong>
          </div>
          <Select v-model="currentFilter" :options="statusFilters" class="w-36" @change="fetchOrders" />
          <div class="flex flex-1 items-center justify-end gap-2">
            <button @click="fetchOrders" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <div class="orders-surface">
        <OrderTable :orders="orders" :loading="loading">
          <template #empty>
            <div class="orders-empty">
              <span>{{ ordersCopy.emptyKicker }}</span>
              <strong>{{ currentFilter ? ordersCopy.emptyFiltered : ordersCopy.emptyAll }}</strong>
              <p>{{ ordersCopy.emptyCopy }}</p>
              <button class="btn btn-primary inline-flex items-center gap-2" @click="router.push('/purchase')">
                <Icon name="wallet" size="sm" />
                <span>{{ t('payment.goToRedeem') }}</span>
              </button>
            </div>
          </template>
          <template #actions="{ row }">
            <div class="flex items-center gap-2">
              <button v-if="row.status === 'PENDING'" @click="handleCancel(row.id)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-yellow-600 hover:bg-yellow-50 dark:text-yellow-400 dark:hover:bg-yellow-900/20">
                <Icon name="x" size="sm" />
                <span>{{ t('payment.orders.cancel') }}</span>
              </button>
              <button v-if="canRequestRefund(row)" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-[#a73a2a] hover:bg-[#a73a2a]/8 dark:text-[#f0b4a8] dark:hover:bg-[#a73a2a]/18">
                <Icon name="dollar" size="sm" />
                <span>{{ t('payment.orders.requestRefund') }}</span>
              </button>
            </div>
          </template>
        </OrderTable>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total > 0" class="orders-pagination">
        <Pagination
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </div>
    </div>

    <!-- Cancel Confirm Dialog -->
    <BaseDialog :show="!!cancelTargetId" :title="t('payment.orders.cancel')" width="narrow" @close="cancelTargetId = null">
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('payment.confirmCancel') }}</p>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="cancelTargetId = null">{{ t('common.cancel') }}</button>
          <button class="btn btn-danger" :disabled="actionLoading" @click="confirmCancel">{{ actionLoading ? t('common.processing') : t('payment.orders.cancel') }}</button>
        </div>
      </template>
    </BaseDialog>

    <!-- Refund Dialog -->
    <BaseDialog :show="!!refundTarget" :title="t('payment.orders.requestRefund')" @close="refundTarget = null">
      <div v-if="refundTarget" class="space-y-4">
        <div class="rounded-lg border border-stone-200/70 bg-stone-100/35 p-4 dark:border-dark-700 dark:bg-dark-900/40">
          <div class="flex justify-between text-sm">
            <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') }}</span>
            <span class="font-mono text-gray-900 dark:text-white">#{{ refundTarget.id }}</span>
          </div>
          <div class="mt-2 flex justify-between text-sm">
            <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.amount') }}</span>
            <span class="text-gray-900 dark:text-white">${{ refundTarget.amount.toFixed(2) }}</span>
          </div>
        </div>
        <div>
          <label class="input-label">{{ t('payment.refundReason') }}</label>
          <textarea v-model="refundReason" rows="3" class="input mt-1 w-full" :placeholder="t('payment.refundReasonPlaceholder')" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="refundTarget = null">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="actionLoading || !refundReason.trim()" @click="confirmRefund">{{ actionLoading ? t('common.processing') : t('payment.orders.requestRefund') }}</button>
        </div>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores'
import { paymentAPI } from '@/api/payment'
import { extractI18nErrorMessage } from '@/utils/apiError'
import type { PaymentOrder } from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import OrderTable from '@/components/payment/OrderTable.vue'

const { t, locale } = useI18n()
const router = useRouter()
const appStore = useAppStore()

const loading = ref(false)
const actionLoading = ref(false)
const orders = ref<PaymentOrder[]>([])
const refundEligibleProviders = ref<Set<string>>(new Set())
const currentFilter = ref('')
const cancelTargetId = ref<number | null>(null)
const refundTarget = ref<PaymentOrder | null>(null)
const refundReason = ref('')
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const zhOrdersCopy = {
  kicker: '往来账册',
  title: '查看充值、退款与处理记录',
  copy: '站内会保留每次充值、退款与状态变动，便于你回看、核对与补记余额流转。',
  emptyKicker: '暂无往来',
  emptyFiltered: '当前筛选下暂无订单',
  emptyAll: '还没有充值或退款记录',
  emptyCopy: '需要补充余额时，可从充值与兑换发起；待支付或已完成的记录都会在这里留痕。'
}

const enOrdersCopy = {
  kicker: 'Order ledger',
  title: 'Review top-ups, refunds, and processing records',
  copy: 'Each recharge, refund, and status change is kept here so you can reconcile account balance movement later.',
  emptyKicker: 'No records',
  emptyFiltered: 'No orders match the current filter',
  emptyAll: 'No top-up or refund records yet',
  emptyCopy: 'Start from recharge and redeem when you need more balance. Pending and completed records will stay here.'
}

const ordersCopy = computed(() => locale.value === 'zh' ? zhOrdersCopy : enOrdersCopy)

const statusFilters = computed(() => [
  { value: '', label: t('common.all') },
  { value: 'PENDING', label: t('payment.status.pending') },
  { value: 'COMPLETED', label: t('payment.status.completed') },
  { value: 'FAILED', label: t('payment.status.failed') },
  { value: 'REFUNDED', label: t('payment.status.refunded') },
])

const currentFilterLabel = computed(() =>
  statusFilters.value.find((item) => item.value === currentFilter.value)?.label || t('common.all')
)

async function fetchOrders() {
  loading.value = true
  try {
    const res = await paymentAPI.getMyOrders({
      page: pagination.page,
      page_size: pagination.page_size,
      status: currentFilter.value || undefined,
    })
    orders.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) { pagination.page = page; fetchOrders() }
function handlePageSizeChange(size: number) { pagination.page_size = size; pagination.page = 1; fetchOrders() }

function handleCancel(orderId: number) { cancelTargetId.value = orderId }

async function confirmCancel() {
  if (!cancelTargetId.value) return
  actionLoading.value = true
  try {
    await paymentAPI.cancelOrder(cancelTargetId.value)
    appStore.showSuccess(t('common.success'))
    cancelTargetId.value = null
    await fetchOrders()
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally {
    actionLoading.value = false
  }
}

function openRefundDialog(order: PaymentOrder) { refundTarget.value = order; refundReason.value = '' }

async function confirmRefund() {
  if (!refundTarget.value || !refundReason.value.trim()) return
  actionLoading.value = true
  try {
    await paymentAPI.requestRefund(refundTarget.value.id, { reason: refundReason.value.trim() })
    appStore.showSuccess(t('common.success'))
    refundTarget.value = null
    refundReason.value = ''
    await fetchOrders()
  } catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  } finally {
    actionLoading.value = false
  }
}

function canRequestRefund(order: PaymentOrder): boolean {
  if (order.status !== 'COMPLETED') return false
  if (!order.provider_instance_id) return false
  return refundEligibleProviders.value.has(order.provider_instance_id)
}

async function loadRefundEligibility() {
  try {
    const res = await paymentAPI.getRefundEligibleProviders()
    refundEligibleProviders.value = new Set(res.data.provider_instance_ids || [])
  } catch { /* ignore — default to hiding refund button */ }
}

onMounted(() => { fetchOrders(); loadRefundEligibility() })
</script>

<style scoped>
.orders-shell {
  position: relative;
}

.orders-brief {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  overflow: hidden;
  border-color: rgba(198, 184, 157, 0.54);
  background:
    radial-gradient(circle at 4% 18%, rgba(255, 252, 245, 0.92), transparent 18rem),
    linear-gradient(100deg, rgba(167, 58, 42, 0.052), transparent 34%),
    linear-gradient(180deg, rgba(255, 252, 245, 0.78), rgba(244, 239, 228, 0.56));
  box-shadow:
    0 22px 54px -46px rgba(31, 35, 32, 0.34),
    inset 0 1px 0 rgba(255, 255, 255, 0.72);
}

.orders-brief-copy {
  max-width: 38rem;
}

.orders-brief-copy > span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.orders-brief h2 {
  margin-top: 0.32rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.12rem, 1.6vw, 1.42rem);
  font-weight: 600;
}

.orders-brief p {
  margin-top: 0.5rem;
  color: #59645a;
  font-size: 0.88rem;
  line-height: 1.72;
}

.orders-filter,
.orders-surface {
  border: 1px solid rgba(198, 184, 157, 0.5);
  border-radius: 10px;
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.62), rgba(244, 239, 228, 0.44));
  box-shadow:
    0 20px 50px -44px rgba(31, 35, 32, 0.28),
    inset 0 1px 0 rgba(255, 255, 255, 0.62);
}

.orders-filter {
  border-color: rgba(198, 184, 157, 0.48);
  background:
    linear-gradient(180deg, rgba(250, 247, 239, 0.72), rgba(237, 229, 212, 0.42));
}

.orders-filter :deep(.select-trigger) {
  border-color: rgba(198, 184, 157, 0.56);
  background: rgba(255, 252, 245, 0.74);
  color: #1f2320;
}

.orders-surface {
  overflow: hidden;
}

.orders-surface :deep(.table-wrapper) {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.orders-surface :deep(table) {
  color: #2f3831;
}

.orders-surface :deep(thead),
.orders-surface :deep(.table-header),
.orders-surface :deep(.sticky-header-cell) {
  background:
    linear-gradient(180deg, rgba(237, 229, 212, 0.86), rgba(225, 213, 193, 0.68)) !important;
}

.orders-surface :deep(th) {
  border-bottom-color: rgba(198, 184, 157, 0.56) !important;
  color: #756850 !important;
  font-weight: 680 !important;
}

.orders-surface :deep(tbody),
.orders-surface :deep(.table-body) {
  background: rgba(255, 252, 245, 0.46) !important;
}

.orders-surface :deep(td) {
  border-bottom-color: rgba(198, 184, 157, 0.28) !important;
}

.orders-surface :deep(tbody tr:last-child td) {
  border-bottom-color: transparent !important;
}

.orders-surface :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.04) !important;
}

.orders-surface :deep(tbody .sticky-col) {
  background: rgba(255, 252, 245, 0.82) !important;
}

.orders-surface :deep(tbody tr:hover .sticky-col) {
  background: rgba(247, 239, 224, 0.96) !important;
}

.orders-surface :deep(> .space-y-3 > div) {
  border-color: rgba(198, 184, 157, 0.48);
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.72), rgba(244, 239, 228, 0.52));
  box-shadow: 0 16px 34px -30px rgba(31, 35, 32, 0.28);
}

.orders-surface :deep(> .space-y-3 > div > .space-y-3 > div) {
  border-color: rgba(198, 184, 157, 0.32);
}

.orders-surface :deep(.border-gray-200) {
  border-color: rgba(198, 184, 157, 0.32) !important;
}

.orders-surface :deep(.bg-gray-200) {
  background-color: rgba(198, 184, 157, 0.34) !important;
}

.orders-surface :deep(.table-wrapper::-webkit-scrollbar-track) {
  background-color: rgba(198, 184, 157, 0.18) !important;
}

.orders-surface :deep(.table-wrapper::-webkit-scrollbar-thumb) {
  background-color: rgba(123, 106, 83, 0.56) !important;
}

.orders-filter-copy {
  display: grid;
  gap: 0.16rem;
  min-width: 9rem;
}

.orders-filter-copy > span,
.orders-empty > span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.orders-filter-copy strong {
  color: #1f2320;
  font-size: 0.95rem;
}

.orders-empty {
  display: grid;
  justify-items: center;
  gap: 0.55rem;
  padding: 2.5rem 1rem;
  text-align: center;
}

.orders-empty strong {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.18rem;
  font-weight: 600;
}

.orders-empty p {
  max-width: 28rem;
  color: #59645a;
  font-size: 0.86rem;
  line-height: 1.7;
}

.orders-brief .btn-primary span,
.orders-empty .btn-primary span {
  color: inherit;
  font-family: inherit;
  font-size: inherit;
  letter-spacing: 0;
}

.orders-pagination {
  overflow: hidden;
  border: 1px solid rgba(198, 184, 157, 0.46);
  border-radius: 10px;
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.58), rgba(244, 239, 228, 0.4));
  box-shadow: 0 18px 42px -38px rgba(31, 35, 32, 0.24);
}

.orders-pagination :deep(.pagination-root) {
  border-top: 0;
  background: transparent;
}

.orders-pagination :deep(button) {
  border-color: rgba(198, 184, 157, 0.54);
  background-color: rgba(255, 252, 245, 0.7);
}

.orders-pagination :deep(button:hover:not(:disabled)) {
  background-color: rgba(250, 247, 239, 0.96);
  color: #a73a2a;
}

.orders-pagination :deep(button[aria-current="page"]) {
  border-color: rgba(167, 58, 42, 0.34);
  background-color: rgba(167, 58, 42, 0.1);
  color: #a73a2a;
}

html.dark .orders-brief {
  border-color: rgba(55, 59, 49, 0.96);
  background:
    radial-gradient(circle at 5% 16%, rgba(244, 239, 228, 0.055), transparent 18rem),
    linear-gradient(100deg, rgba(167, 58, 42, 0.105), transparent 34%),
    linear-gradient(180deg, rgba(29, 31, 25, 0.92), rgba(19, 21, 17, 0.86));
  box-shadow:
    0 24px 60px -46px rgba(0, 0, 0, 0.82),
    inset 0 1px 0 rgba(244, 239, 228, 0.045);
}

html.dark .orders-brief-copy > span {
  color: #b99a78;
}

html.dark .orders-brief h2 {
  color: #f4efe4;
}

html.dark .orders-brief p {
  color: #a8a091;
}

html.dark .orders-filter,
html.dark .orders-surface {
  border-color: rgba(55, 59, 49, 0.96);
  background:
    linear-gradient(180deg, rgba(27, 30, 24, 0.9), rgba(16, 18, 14, 0.76));
  box-shadow:
    0 22px 52px -44px rgba(0, 0, 0, 0.82),
    inset 0 1px 0 rgba(244, 239, 228, 0.035);
}

html.dark .orders-filter :deep(.select-trigger) {
  border-color: rgba(68, 71, 58, 0.92);
  background: rgba(17, 19, 15, 0.66);
  color: #f4efe4;
}

html.dark .orders-filter-copy > span,
html.dark .orders-empty > span {
  color: #b99a78;
}

html.dark .orders-filter-copy strong,
html.dark .orders-empty strong {
  color: #f4efe4;
}

html.dark .orders-empty p {
  color: #a8a091;
}

html.dark .orders-surface :deep(.table-wrapper) {
  background: transparent;
}

html.dark .orders-surface :deep(table) {
  color: #d8cfbf;
}

html.dark .orders-surface :deep(thead),
html.dark .orders-surface :deep(.table-header),
html.dark .orders-surface :deep(.sticky-header-cell) {
  background:
    linear-gradient(180deg, rgba(54, 48, 38, 0.9), rgba(31, 33, 26, 0.82)) !important;
}

html.dark .orders-surface :deep(th) {
  border-bottom-color: rgba(72, 69, 56, 0.92) !important;
  color: #c0ad91 !important;
}

html.dark .orders-surface :deep(tbody),
html.dark .orders-surface :deep(.table-body) {
  background: rgba(17, 19, 15, 0.55) !important;
}

html.dark .orders-surface :deep(td) {
  border-bottom-color: rgba(55, 59, 49, 0.72) !important;
}

html.dark .orders-surface :deep(tbody tr:last-child td) {
  border-bottom-color: transparent !important;
}

html.dark .orders-surface :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.095) !important;
}

html.dark .orders-surface :deep(tbody .sticky-col) {
  background: rgba(18, 20, 16, 0.96) !important;
}

html.dark .orders-surface :deep(tbody tr:hover .sticky-col) {
  background: rgba(35, 31, 25, 0.98) !important;
}

html.dark .orders-surface :deep(> .space-y-3 > div) {
  border-color: rgba(55, 59, 49, 0.96);
  background:
    linear-gradient(180deg, rgba(27, 30, 24, 0.92), rgba(17, 19, 15, 0.84));
  box-shadow: 0 18px 38px -30px rgba(0, 0, 0, 0.72);
}

html.dark .orders-surface :deep(> .space-y-3 > div > .space-y-3 > div) {
  border-color: rgba(55, 59, 49, 0.82);
}

html.dark .orders-surface :deep(.border-gray-200),
html.dark .orders-surface :deep(.dark\:border-dark-700) {
  border-color: rgba(55, 59, 49, 0.82) !important;
}

html.dark .orders-surface :deep(.bg-gray-200),
html.dark .orders-surface :deep(.dark\:bg-dark-700) {
  background-color: rgba(68, 71, 58, 0.72) !important;
}

html.dark .orders-surface :deep(.table-wrapper::-webkit-scrollbar-track) {
  background-color: rgba(244, 239, 228, 0.055) !important;
}

html.dark .orders-surface :deep(.table-wrapper::-webkit-scrollbar-thumb) {
  background-color: rgba(185, 154, 120, 0.54) !important;
}

html.dark .orders-surface :deep(.text-gray-900),
html.dark .orders-surface :deep(.dark\:text-white),
html.dark .orders-surface :deep(.dark\:text-gray-100) {
  color: #f4efe4 !important;
}

html.dark .orders-surface :deep(.text-gray-700),
html.dark .orders-surface :deep(.text-gray-500),
html.dark .orders-surface :deep(.dark\:text-gray-300),
html.dark .orders-surface :deep(.dark\:text-gray-400) {
  color: #a8a091 !important;
}

html.dark .orders-pagination {
  border-color: rgba(55, 59, 49, 0.96);
  background:
    linear-gradient(180deg, rgba(27, 30, 24, 0.86), rgba(17, 19, 15, 0.74));
  box-shadow: 0 20px 46px -38px rgba(0, 0, 0, 0.74);
}

html.dark .orders-pagination :deep(.pagination-root) {
  border-top: 0;
  background: transparent;
}

html.dark .orders-pagination :deep(.text-gray-700),
html.dark .orders-pagination :deep(.dark\:text-gray-300) {
  color: #a8a091 !important;
}

html.dark .orders-pagination :deep(button) {
  border-color: rgba(68, 71, 58, 0.9);
  background-color: rgba(17, 19, 15, 0.66);
  color: #c9c0ac;
}

html.dark .orders-pagination :deep(button:hover:not(:disabled)) {
  border-color: rgba(167, 58, 42, 0.34);
  background-color: rgba(167, 58, 42, 0.11);
  color: #f0b4a8;
}

html.dark .orders-pagination :deep(button[aria-current="page"]) {
  border-color: rgba(167, 58, 42, 0.4);
  background-color: rgba(167, 58, 42, 0.18);
  color: #f0b4a8;
}

@media (max-width: 640px) {
  .orders-brief {
    align-items: flex-start;
    flex-direction: column;
  }

  .orders-brief .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
