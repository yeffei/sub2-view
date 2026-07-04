<template>
  <AppLayout>
    <div class="orders-shell mx-auto max-w-[76rem] space-y-4 sm:space-y-5">
      <section class="orders-brief card p-6 lg:p-7">
        <div class="orders-brief-copy">
          <span>往来账册</span>
          <h2>查看充值、退款与处理记录</h2>
          <p>站内会保留每次充值、退款与状态变动，便于你回看、核对与补记余额流转。</p>
        </div>
        <button class="btn btn-primary inline-flex items-center gap-2" @click="router.push('/purchase')">
          <Icon name="wallet" size="sm" />
          <span>充值与兑换</span>
        </button>
      </section>

      <!-- Filters -->
      <div class="orders-filter card p-5">
        <div class="flex flex-wrap items-center gap-3">
          <div class="orders-filter-copy">
            <span>筛选</span>
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
      <OrderTable :orders="orders" :loading="loading" class="orders-table orders-surface">
        <template #empty>
          <div class="orders-empty">
            <span>暂无往来</span>
            <strong>{{ currentFilter ? '当前筛选下暂无订单' : '还没有充值或退款记录' }}</strong>
            <p>需要补充余额时，可从充值与兑换发起；待支付或已完成的记录都会在这里留痕。</p>
            <button class="btn btn-primary inline-flex items-center gap-2" @click="router.push('/purchase')">
              <Icon name="wallet" size="sm" />
              <span>充值与兑换</span>
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

      <!-- Pagination -->
      <Pagination
        v-if="pagination.total > 0"
        :page="pagination.page"
        :total="pagination.total"
        :page-size="pagination.page_size"
        @update:page="handlePageChange"
        @update:pageSize="handlePageSizeChange"
      />
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

const { t } = useI18n()
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
}

.orders-brief-copy {
  max-width: 38rem;
}

.orders-brief span {
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
.orders-surface :deep(.table-wrapper) {
  border-color: rgba(198, 184, 157, 0.46);
  background: rgba(250, 247, 239, 0.52);
}

.orders-filter-copy {
  display: grid;
  gap: 0.16rem;
  min-width: 9rem;
}

.orders-filter-copy span,
.orders-empty span {
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

.dark .orders-brief h2 {
  color: #f4efe4;
}

.dark .orders-brief p {
  color: #879186;
}

.dark .orders-filter,
.dark .orders-table :deep(.table-wrapper) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .orders-filter-copy span,
.dark .orders-empty span,
.dark .orders-empty p {
  color: #879186;
}

.dark .orders-filter-copy strong,
.dark .orders-empty strong {
  color: #f4efe4;
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
