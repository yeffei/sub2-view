<template>
  <AppLayout>
    <div class="payment-result-shell mx-auto w-full max-w-2xl space-y-6 py-8">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>
      <template v-else>
        <!-- Status Icon -->
        <section class="payment-result-hero text-center">
          <div class="payment-result-kicker">{{ t('payment.result.receiptKicker') }}</div>
          <div v-if="isSuccess"
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
            <Icon name="check" size="xl" class="text-green-500" :stroke-width="2" />
          </div>
          <div v-else-if="isPending"
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900/30">
            <div class="h-10 w-10 animate-spin rounded-full border-4 border-yellow-500 border-t-transparent"></div>
          </div>
          <div v-else
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
            <Icon name="x" size="xl" class="text-red-500" :stroke-width="2" />
          </div>
          <h2 class="mt-4 font-serif text-3xl font-semibold text-zen-ink dark:text-zen-paper">
            {{ statusTitle }}
          </h2>
          <p class="payment-result-note mt-3 text-sm text-zen-mist dark:text-zen-stone">
            {{ statusNote }}
          </p>
        </section>
        <!-- Order Info -->
        <div v-if="order" class="payment-result-card">
          <div class="space-y-3 text-sm">
            <div class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">#{{ order.id }}</span>
            </div>
            <div v-if="order.out_trade_no" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderNo') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ order.out_trade_no }}</span>
            </div>
            <div class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.baseAmount') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatGatewayAmount(baseAmount) }}</span>
            </div>
            <div v-if="order.fee_rate > 0" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.fee') }} ({{ order.fee_rate }}%)</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatGatewayAmount(feeAmount) }}</span>
            </div>
            <div class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') }}</span>
              <span class="font-bold text-primary-600 dark:text-primary-400">{{ formatGatewayAmount(order.pay_amount) }}</span>
            </div>
            <div v-if="order.amount !== order.pay_amount" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.creditedAmount') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ order.order_type === 'balance' ? '$' + order.amount.toFixed(2) : formatGatewayAmount(order.amount) }}</span>
            </div>
            <div class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ t(paymentMethodI18nKey(order.payment_type), normalizedOrderPaymentType(order.payment_type)) }}</span>
            </div>
            <div class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.status') }}</span>
              <OrderStatusBadge :status="order.status" />
            </div>
          </div>
        </div>
        <!-- EasyPay return info (when no order loaded) -->
        <div v-else-if="returnInfo" class="payment-result-card">
          <div class="space-y-3 text-sm">
            <div v-if="returnInfo.outTradeNo" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ returnInfo.outTradeNo }}</span>
            </div>
            <div v-if="returnInfo.money" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatGatewayAmount(Number(returnInfo.money) || 0) }}</span>
            </div>
            <div v-if="returnInfo.type" class="payment-result-row flex justify-between gap-4">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ t(paymentMethodI18nKey(returnInfo.type), normalizedOrderPaymentType(returnInfo.type)) }}</span>
            </div>
          </div>
        </div>
        <!-- Actions -->
        <div class="payment-result-actions flex gap-3">
          <button class="btn btn-secondary flex-1" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') }}</button>
          <button class="btn btn-primary flex-1" @click="router.push('/orders')">{{ t('payment.result.viewOrders') }}</button>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import {
  PAYMENT_RECOVERY_STORAGE_KEY,
  clearPaymentRecoverySnapshot,
  readPaymentRecoverySnapshot,
} from '@/components/payment/paymentFlow'
import { usePaymentStore } from '@/stores/payment'
import { paymentAPI } from '@/api/payment'
import type { PaymentOrder } from '@/types/payment'
import { formatPaymentAmount, normalizePaymentCurrency } from '@/components/payment/currency'
import { normalizePaymentMethodForDisplay, paymentMethodI18nKey } from './paymentUx'

const i18n = useI18n()
const { t } = i18n
const route = useRoute()
const router = useRouter()
const paymentStore = usePaymentStore()

const order = ref<PaymentOrder | null>(null)
const loading = ref(true)
const currency = ref('CNY')

interface ReturnInfo {
  outTradeNo: string
  money: string
  type: string
  tradeStatus: string
}
const returnInfo = ref<ReturnInfo | null>(null)

const SUCCESS_STATUSES = new Set(['COMPLETED', 'PAID', 'RECHARGING'])
const PENDING_STATUSES = new Set(['PENDING', 'CREATED', 'WAITING', 'PROCESSING'])
const STATUS_REFRESH_INTERVAL_MS = 2000
const STATUS_REFRESH_MAX_ATTEMPTS = 15

let statusRefreshTimer: ReturnType<typeof setTimeout> | null = null
const refreshAttempts = ref(0)

/** 充值金额 = pay_amount / (1 + fee_rate/100)，fee_rate=0 时等于 pay_amount */
const baseAmount = computed(() => {
  if (!order.value) return 0
  const feeRate = Number(order.value.fee_rate) || 0
  if (feeRate <= 0) return order.value.pay_amount ?? 0
  return Math.round((order.value.pay_amount / (1 + feeRate / 100)) * 100) / 100
})

/** 手续费 = pay_amount - baseAmount */
const feeAmount = computed(() => {
  if (!order.value) return 0
  const feeRate = Number(order.value.fee_rate) || 0
  if (feeRate <= 0) return 0
  return Math.round((order.value.pay_amount - baseAmount.value) * 100) / 100
})

const localeCode = computed(() => {
  const raw = i18n.locale as unknown
  if (typeof raw === 'string') return raw
  if (raw && typeof raw === 'object' && 'value' in raw) {
    return String((raw as { value?: string }).value || '')
  }
  return undefined
})

const isSuccess = computed(() => {
  return isSuccessStatus(order.value?.status)
})

const isPending = computed(() => {
  return isPendingStatus(order.value?.status)
})

const statusTitle = computed(() => {
  if (isSuccess.value) {
    return t('payment.result.success')
  }
  if (isPending.value) {
    return t('payment.result.processing')
  }
  return t('payment.result.failed')
})

const statusNote = computed(() => {
  if (isSuccess.value) {
    return t('payment.result.successNote')
  }
  if (isPending.value) {
    return t('payment.result.processingHint')
  }
  return t('payment.result.failedNote')
})

function normalizedOrderPaymentType(paymentType: string): string {
  return normalizePaymentMethodForDisplay(paymentType) || paymentType
}

function formatGatewayAmount(value: number): string {
  return formatPaymentAmount(value, currency.value, localeCode.value)
}

function setResolvedOrder(nextOrder: PaymentOrder | null): void {
  order.value = nextOrder
  if (nextOrder?.currency) {
    currency.value = normalizePaymentCurrency(nextOrder.currency)
  }
}

function normalizeOrderStatus(status: string | null | undefined): string {
  return String(status || '').trim().toUpperCase()
}

function isSuccessStatus(status: string | null | undefined): boolean {
  return SUCCESS_STATUSES.has(normalizeOrderStatus(status))
}

function isPendingStatus(status: string | null | undefined): boolean {
  return PENDING_STATUSES.has(normalizeOrderStatus(status))
}

function readRouteQueryString(key: string): string {
  const value = route.query[key]
  if (Array.isArray(value)) {
    return typeof value[0] === 'string' ? value[0] : ''
  }
  return typeof value === 'string' ? value : ''
}

function restoreRecoverySnapshot(context: {
  resumeToken: string
  routeOrderId: number
  routeOutTradeNo: string
}) {
  if (typeof window === 'undefined') {
    return null
  }

  const rawSnapshot = window.localStorage.getItem(PAYMENT_RECOVERY_STORAGE_KEY)
  if (!rawSnapshot) {
    return null
  }

  if (context.resumeToken) {
    return readPaymentRecoverySnapshot(rawSnapshot, {
      resumeToken: context.resumeToken,
    })
  }

  if (!context.routeOrderId && !context.routeOutTradeNo) {
    return null
  }

  const restored = readPaymentRecoverySnapshot(rawSnapshot)
  if (!restored) {
    return null
  }

  if (context.routeOrderId > 0 && restored.orderId !== context.routeOrderId) {
    return null
  }

  if (context.routeOutTradeNo && restored.outTradeNo !== context.routeOutTradeNo) {
    return null
  }

  return restored
}

async function resolveOrderFromResumeToken(resumeToken: string): Promise<PaymentOrder | null> {
  try {
    const result = await paymentAPI.resolveOrderPublicByResumeToken(resumeToken)
    return result.data
  } catch (_err: unknown) {
    return null
  }
}

async function resolveOrderFromOutTradeNo(outTradeNo: string): Promise<PaymentOrder | null> {
  try {
    const result = await paymentAPI.verifyOrder(outTradeNo)
    return result.data
  } catch (_err: unknown) {
    try {
      const result = await paymentAPI.verifyOrderPublic(outTradeNo)
      return result.data
    } catch (_innerErr: unknown) {
      return null
    }
  }
}

function clearStatusRefreshTimer(): void {
  if (statusRefreshTimer !== null) {
    clearTimeout(statusRefreshTimer)
    statusRefreshTimer = null
  }
}

function clearRecoverySnapshot(): void {
  if (typeof window === 'undefined') return
  clearPaymentRecoverySnapshot(window.localStorage, PAYMENT_RECOVERY_STORAGE_KEY)
}

function clearRecoverySnapshotForTerminalStatus(status: string | null | undefined): void {
  if (!status) return
  if (!isPendingStatus(status)) {
    clearRecoverySnapshot()
  }
}

function scheduleStatusRefresh(refreshOrder: (() => Promise<PaymentOrder | null>) | null): void {
  clearStatusRefreshTimer()
  if (!refreshOrder || !isPending.value || refreshAttempts.value >= STATUS_REFRESH_MAX_ATTEMPTS) {
    return
  }

  statusRefreshTimer = setTimeout(async () => {
    refreshAttempts.value += 1
    const refreshedOrder = await refreshOrder()
    if (refreshedOrder) {
      setResolvedOrder(refreshedOrder)
      clearRecoverySnapshotForTerminalStatus(refreshedOrder.status)
    }

    if (isPendingStatus(order.value?.status)) {
      scheduleStatusRefresh(refreshOrder)
    }
  }, STATUS_REFRESH_INTERVAL_MS)
}

onMounted(async () => {
  const resumeToken = readRouteQueryString('resume_token')
  const routeOrderId = Number(readRouteQueryString('order_id')) || 0
  let outTradeNo = readRouteQueryString('out_trade_no')
  let orderId = 0
  let resumeTokenLookupFailed = false

  const restored = restoreRecoverySnapshot({
    resumeToken,
    routeOrderId,
    routeOutTradeNo: outTradeNo,
  })
  if (restored?.orderId) {
    orderId = restored.orderId
  }
  if (restored?.currency) {
    currency.value = normalizePaymentCurrency(restored.currency)
  }
  if (!outTradeNo && restored?.outTradeNo) {
    outTradeNo = restored.outTradeNo
  }

  if (resumeToken) {
    const resolvedOrder = await resolveOrderFromResumeToken(resumeToken)
    if (resolvedOrder) {
      setResolvedOrder(resolvedOrder)
      if (!orderId) {
        orderId = resolvedOrder.id
      }
    } else if (routeOrderId > 0) {
      resumeTokenLookupFailed = true
      orderId = routeOrderId
    } else {
      resumeTokenLookupFailed = true
    }
  } else if (routeOrderId > 0) {
    orderId = routeOrderId
  }

  const hasLegacyFallbackContext = readRouteQueryString('trade_status').trim() !== ''
  const shouldUsePublicOutTradeNo = outTradeNo !== '' && (hasLegacyFallbackContext || routeOrderId > 0 || orderId > 0)

  if (!order.value && orderId && (!resumeToken || routeOrderId > 0)) {
    try {
      setResolvedOrder(await paymentStore.pollOrderStatus(orderId))
    } catch (_err: unknown) {
      // Order lookup failed, will try legacy fallback below when possible.
    }
  }

  if (!order.value && shouldUsePublicOutTradeNo && (!resumeToken || resumeTokenLookupFailed)) {
    const legacyOrder = await resolveOrderFromOutTradeNo(outTradeNo)
    if (legacyOrder) {
      setResolvedOrder(legacyOrder)
      if (!orderId) {
        orderId = legacyOrder.id
      }
    }
  }

  if (!order.value && !orderId && outTradeNo && hasLegacyFallbackContext) {
    returnInfo.value = {
      outTradeNo,
      money: String(route.query.money || ''),
      type: String(route.query.type || ''),
      tradeStatus: String(route.query.trade_status || ''),
    }
  }

  const refreshOrder = async (): Promise<PaymentOrder | null> => {
    if (resumeToken) {
      const resolvedOrder = await resolveOrderFromResumeToken(resumeToken)
      if (resolvedOrder) {
        return resolvedOrder
      }
    }

    if (orderId) {
      try {
        return await paymentStore.pollOrderStatus(orderId)
      } catch (_err: unknown) {
        // Fall through to legacy public verification when order polling is unavailable.
      }
    }

    if (shouldUsePublicOutTradeNo) {
      return await resolveOrderFromOutTradeNo(outTradeNo)
    }

    return null
  }

  if (isPendingStatus(order.value?.status)) {
    scheduleStatusRefresh(refreshOrder)
  } else if (order.value) {
    clearRecoverySnapshotForTerminalStatus(order.value.status)
  } else if (returnInfo.value) {
    clearRecoverySnapshot()
  }
  loading.value = false
})

onBeforeUnmount(() => {
  clearStatusRefreshTimer()
})
</script>

<style scoped>
.payment-result-hero,
.payment-result-card {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 0.9rem;
  background: rgba(250, 247, 239, 0.58);
  box-shadow: 0 20px 48px -38px rgba(31, 35, 32, 0.26);
}

.payment-result-hero {
  padding: 1.75rem 1.5rem;
}

.payment-result-card {
  padding: 1.25rem;
}

.payment-result-kicker {
  margin-bottom: 1rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.payment-result-note {
  max-width: 28rem;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.72;
}

.payment-result-row {
  padding-bottom: 0.72rem;
  border-bottom: 1px solid rgba(216, 205, 185, 0.4);
}

.payment-result-row:last-child {
  padding-bottom: 0;
  border-bottom: 0;
}

.payment-result-actions {
  flex-wrap: wrap;
}

.dark .payment-result-hero,
.dark .payment-result-card {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .payment-result-kicker {
  color: #879186;
}

.dark .payment-result-row {
  border-bottom-color: rgba(48, 52, 43, 0.78);
}

@media (max-width: 640px) {
  .payment-result-shell {
    padding-top: 1rem;
    padding-bottom: 1.25rem;
  }

  .payment-result-actions > * {
    width: 100%;
  }
}
</style>
