<template>
  <AppLayout>
    <div class="payment-shell mx-auto max-w-[76rem] space-y-6">
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>

      <template v-else>
        <section class="payment-shop card p-6 lg:p-7">
          <div class="payment-shop-copy">
            <span>{{ t('payment.externalPurchaseKicker') }}</span>
            <h3>{{ t('payment.externalPurchaseTitle') }}</h3>
          </div>
          <div class="payment-shop-actions">
            <button class="btn btn-secondary" @click="router.push('/orders')">
              {{ t('payment.viewOrders') }}
            </button>
            <button class="btn btn-primary" @click="openExternalPurchase">
              {{ t('payment.externalPurchaseAction') }}
            </button>
            <button class="btn btn-secondary" @click="scrollToRedeem">
              {{ t('payment.goToRedeem') }}
            </button>
          </div>
        </section>

        <div v-if="checkoutUnavailable" class="payment-unavailable card p-8 text-center">
          <Icon
            name="creditCard"
            size="xl"
            class="mx-auto mb-4 text-gray-300 dark:text-dark-600"
          />
          <h2>{{ t('payment.unavailableTitle') }}</h2>
          <p>{{ t('payment.unavailableWithExternalPurchase') }}</p>
          <div>
            <button class="btn btn-primary" @click="openExternalPurchase">
              {{ t('payment.externalPurchaseAction') }}
            </button>
            <button class="btn btn-secondary" @click="scrollToRedeem">
              {{ t('payment.useRedeemCode') }}
            </button>
            <button class="btn btn-secondary" @click="router.push('/orders')">
              {{ t('payment.viewOrders') }}
            </button>
          </div>
        </div>

        <template v-else>
          <div
            v-if="errorMessage && paymentPhase === 'select'"
            class="payment-inline-error"
          >
            <Icon name="exclamationTriangle" size="md" />
            <div>
              <strong>{{ errorMessage }}</strong>
              <p v-if="errorHintMessage">{{ errorHintMessage }}</p>
              <p v-else>{{ t('payment.retryOrCheckOrders') }}</p>
            </div>
            <button class="btn btn-secondary" @click="router.push('/orders')">
              {{ t('payment.viewOrders') }}
            </button>
          </div>

          <template v-if="paymentPhase === 'paying'">
            <PaymentStatusPanel
              :order-id="paymentState.orderId"
              :qr-code="paymentState.qrCode"
              :expires-at="paymentState.expiresAt"
              :payment-type="paymentState.paymentType"
              :pay-url="paymentState.payUrl"
              :order-type="paymentState.orderType"
              :currency="paymentState.currency || selectedCurrency"
              @done="onPaymentDone"
              @success="onPaymentSuccess"
              @settled="onPaymentSettled"
            />
          </template>

          <template v-else>
            <div v-if="rechargeUnavailable" class="card py-16 text-center">
              <p class="text-gray-500 dark:text-gray-400">
                {{ t('payment.notAvailable') }}
              </p>
            </div>

            <template v-else>
              <div class="card p-6">
                <div class="payment-recharge-head">
                  <div>
                    <p class="text-xs font-medium text-gray-400 dark:text-gray-500">
                      {{ t('payment.rechargeAccount') }}
                    </p>
                    <p class="mt-1 text-base font-semibold text-gray-900 dark:text-white">
                      {{ user?.username || '' }}
                    </p>
                  </div>
                  <div class="payment-balance-chip">
                    <span>{{ t('payment.currentBalance') }}</span>
                    <strong>{{ user?.balance?.toFixed(2) || '0.00' }}</strong>
                  </div>
                </div>
                <AmountInput
                  v-model="amount"
                  :amounts="[10, 20, 50, 100, 200, 500, 1000, 2000, 5000]"
                  :min="globalMinAmount"
                  :max="globalMaxAmount"
                />
                <p
                  v-if="amountError"
                  class="mt-2 text-xs text-amber-600 dark:text-amber-300"
                >
                  {{ amountError }}
                </p>
              </div>

              <div class="card p-6">
                <PaymentMethodSelector
                  :methods="methodOptions"
                  :selected="selectedMethod"
                  @select="selectedMethod = $event"
                />
              </div>

              <div v-if="validAmount > 0" class="card p-6">
                <div class="space-y-2 text-sm">
                  <div class="flex justify-between">
                    <span class="text-gray-500 dark:text-gray-400">
                      {{ t('payment.paymentAmount') }}
                    </span>
                    <span class="text-gray-900 dark:text-white">
                      {{ formatSelectedPaymentAmount(validAmount) }}
                    </span>
                  </div>
                  <div v-if="feeRate > 0" class="flex justify-between">
                    <span class="text-gray-500 dark:text-gray-400">
                      {{ t('payment.fee') }} ({{ feeRate }}%)
                    </span>
                    <span class="text-gray-900 dark:text-white">
                      {{ formatSelectedPaymentAmount(feeAmount) }}
                    </span>
                  </div>
                  <div
                    v-if="feeRate > 0"
                    class="flex justify-between border-t border-gray-200 pt-2 dark:border-dark-600"
                  >
                    <span class="font-medium text-gray-700 dark:text-gray-300">
                      {{ t('payment.actualPay') }}
                    </span>
                    <span class="text-lg font-bold text-primary-600 dark:text-primary-400">
                      {{ formatSelectedPaymentAmount(totalAmount) }}
                    </span>
                  </div>
                  <div
                    v-if="balanceRechargeMultiplier !== 1"
                    class="flex justify-between"
                    :class="{
                      'border-t border-gray-200 pt-2 dark:border-dark-600': feeRate <= 0,
                    }"
                  >
                    <span class="text-gray-500 dark:text-gray-400">
                      {{ t('payment.creditedBalance') }}
                    </span>
                    <span class="text-gray-900 dark:text-white">
                      ${{ creditedAmount.toFixed(2) }}
                    </span>
                  </div>
                  <p
                    v-if="balanceRechargeMultiplier !== 1"
                    class="border-t border-gray-200 pt-2 text-xs text-gray-500 dark:border-dark-600 dark:text-gray-400"
                  >
                    {{ t('payment.rechargeRatePreview', { usd: balanceRechargeMultiplier.toFixed(2) }) }}
                  </p>
                  <template v-if="campaignEnabled">
                    <div class="flex justify-between border-t border-gray-200 pt-2 dark:border-dark-600">
                      <span class="text-gray-500 dark:text-gray-400">
                        {{ t('payment.rechargeCampaignBonus') }}
                      </span>
                      <span class="text-gray-900 dark:text-white">
                        +${{ campaignBonusAmount.toFixed(2) }}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="font-medium text-gray-700 dark:text-gray-300">
                        {{ t('payment.rechargeCampaignTotal') }}
                      </span>
                      <span class="text-lg font-bold text-primary-600 dark:text-primary-400">
                        ${{ totalCreditedAmount.toFixed(2) }}
                      </span>
                    </div>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      {{ t('payment.rechargeCampaignPreview', { amount: formatSelectedPaymentAmount(campaignThreshold), rate: campaignBonusRate.toFixed(2), bonus: campaignBonusAmount.toFixed(2) }) }}
                    </p>
                  </template>
                </div>
              </div>

              <button
                :class="['btn w-full py-3 text-base font-medium', paymentButtonClass]"
                :disabled="!canSubmit || submitting"
                @click="handleSubmitRecharge"
              >
                <span v-if="submitting" class="flex items-center justify-center gap-2">
                  <span class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
                  {{ t('common.processing') }}
                </span>
                <span v-else>
                  {{ t('payment.createOrder') }}
                  {{ formatSelectedPaymentAmount(totalAmount) }}
                </span>
              </button>
            </template>
          </template>

          <section
            v-if="paymentPhase === 'select'"
            id="purchase-redeem-section"
            class="payment-redeem-section"
          >
            <div class="payment-redeem-heading">
              <div>
                <span>{{ t('redeem.panelKicker') }}</span>
                <h3>{{ t('redeem.redeemCodeDeposit') }}</h3>
              </div>
            </div>
            <RedeemPanel :embedded="true" :show-intro="false" />
          </section>

          <div
            v-if="(checkout.help_text || checkout.help_image_url) && paymentPhase === 'select'"
            class="card p-4"
          >
            <div class="flex flex-col items-center gap-3">
              <img
                v-if="checkout.help_image_url"
                :src="checkout.help_image_url"
                alt=""
                class="h-40 max-w-full cursor-pointer rounded-lg object-contain transition-opacity hover:opacity-80"
                @click="previewImage = checkout.help_image_url"
              />
              <p
                v-if="checkout.help_text"
                class="text-center text-sm text-gray-500 dark:text-gray-400"
              >
                {{ checkout.help_text }}
              </p>
            </div>
          </div>
        </template>
      </template>
    </div>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="previewImage"
          class="fixed inset-0 z-[60] flex items-center justify-center bg-black/70 backdrop-blur-sm"
          @click="previewImage = ''"
        >
          <img
            :src="previewImage"
            alt=""
            class="max-h-[85vh] max-w-[90vw] rounded-xl object-contain shadow-2xl"
          />
        </div>
      </Transition>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usePaymentStore } from '@/stores/payment'
import { useAppStore } from '@/stores'
import { paymentAPI } from '@/api/payment'
import { extractApiErrorMessage, extractI18nErrorMessage } from '@/utils/apiError'
import { isMobileDevice } from '@/utils/device'
import type { CheckoutInfoResponse, CreateOrderResult, OrderType } from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import AmountInput from '@/components/payment/AmountInput.vue'
import PaymentMethodSelector from '@/components/payment/PaymentMethodSelector.vue'
import { METHOD_ORDER, getPaymentPopupFeatures } from '@/components/payment/providerConfig'
import {
  PAYMENT_RECOVERY_STORAGE_KEY,
  buildCreateOrderPayload,
  clearPaymentRecoverySnapshot,
  decidePaymentLaunch,
  getVisibleMethods,
  normalizeVisibleMethod,
  readPaymentRecoverySnapshot,
  type PaymentRecoverySnapshot,
  writePaymentRecoverySnapshot,
} from '@/components/payment/paymentFlow'
import PaymentStatusPanel from '@/components/payment/PaymentStatusPanel.vue'
import Icon from '@/components/icons/Icon.vue'
import RedeemPanel from '@/components/user/RedeemPanel.vue'
import { formatPaymentAmount, normalizePaymentCurrency } from '@/components/payment/currency'
import type { PaymentMethodOption } from '@/components/payment/PaymentMethodSelector.vue'
import { buildPaymentErrorToastMessage, describePaymentScenarioError } from './paymentUx'
import { hasWechatResumeQuery, parseWechatResumeRoute, stripWechatResumeQuery } from './paymentWechatResume'

const i18n = useI18n()
const { t } = i18n
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const paymentStore = usePaymentStore()
const appStore = useAppStore()
const fallbackExternalPurchaseUrl = 'https://catfk.com/shop/SST'

const user = computed(() => authStore.user)
const externalPurchaseUrl = computed(() => {
  const configured = appStore.cachedPublicSettings?.purchase_subscription_url?.trim()
  return configured || fallbackExternalPurchaseUrl
})

const loading = ref(true)
const submitting = ref(false)
const checkoutUnavailable = ref(false)
const errorMessage = ref('')
const errorHintMessage = ref('')
const amount = ref<number | null>(null)
const selectedMethod = ref('')
const previewImage = ref('')
const paymentPhase = ref<'select' | 'paying'>('select')

interface CreateOrderOptions {
  openid?: string
  wechatResumeToken?: string
  paymentType?: string
  isResume?: boolean
  mobileQrFallbackAttempted?: boolean
}

interface WeixinJSBridgeLike {
  invoke(
    action: string,
    payload: Record<string, unknown>,
    callback: (result: Record<string, unknown>) => void,
  ): void
}

function emptyPaymentState(): PaymentRecoverySnapshot {
  return {
    orderId: 0,
    amount: 0,
    qrCode: '',
    expiresAt: '',
    paymentType: '',
    payUrl: '',
    outTradeNo: '',
    clientSecret: '',
    intentId: '',
    currency: '',
    countryCode: '',
    paymentEnv: '',
    payAmount: 0,
    orderType: '',
    paymentMode: '',
    resumeToken: '',
    createdAt: 0,
  }
}

function getWeixinJSBridge(): WeixinJSBridgeLike | undefined {
  return (window as Window & { WeixinJSBridge?: WeixinJSBridgeLike }).WeixinJSBridge
}

function waitForWeixinJSBridge(timeoutMs = 4000): Promise<WeixinJSBridgeLike | null> {
  const existing = getWeixinJSBridge()
  if (existing) return Promise.resolve(existing)

  return new Promise((resolve) => {
    let settled = false
    const finish = (bridge: WeixinJSBridgeLike | null) => {
      if (settled) return
      settled = true
      document.removeEventListener('WeixinJSBridgeReady', handleReady)
      document.removeEventListener('onWeixinJSBridgeReady', handleReady)
      window.clearTimeout(timer)
      resolve(bridge)
    }
    const handleReady = () => finish(getWeixinJSBridge() ?? null)
    const timer = window.setTimeout(
      () => finish(getWeixinJSBridge() ?? null),
      timeoutMs,
    )
    document.addEventListener('WeixinJSBridgeReady', handleReady, false)
    document.addEventListener('onWeixinJSBridgeReady', handleReady, false)
  })
}

async function invokeWechatJsapiPayment(
  payload: Record<string, unknown>,
): Promise<Record<string, unknown>> {
  const bridge = await waitForWeixinJSBridge()
  if (!bridge) {
    throw new Error('WECHAT_JSAPI_UNAVAILABLE')
  }
  return new Promise((resolve) => {
    bridge.invoke('getBrandWCPayRequest', payload, (result) => resolve(result || {}))
  })
}

const paymentState = ref<PaymentRecoverySnapshot>(emptyPaymentState())

function persistRecoverySnapshot(snapshot: PaymentRecoverySnapshot) {
  if (typeof window === 'undefined' || !snapshot.orderId) return
  writePaymentRecoverySnapshot(
    window.localStorage,
    snapshot,
    PAYMENT_RECOVERY_STORAGE_KEY,
  )
}

function removeRecoverySnapshot() {
  if (typeof window === 'undefined') return
  clearPaymentRecoverySnapshot(window.localStorage, PAYMENT_RECOVERY_STORAGE_KEY)
}

function resetPayment() {
  paymentPhase.value = 'select'
  paymentState.value = emptyPaymentState()
  removeRecoverySnapshot()
}

async function redirectToPaymentResult(
  state: PaymentRecoverySnapshot,
): Promise<void> {
  const query: Record<string, string | undefined> = {}
  if (state.orderId > 0) {
    query.order_id = String(state.orderId)
  }
  if (state.outTradeNo) {
    query.out_trade_no = state.outTradeNo
  }
  if (state.resumeToken) {
    query.resume_token = state.resumeToken
  }
  await router.push({
    path: '/payment/result',
    query,
  })
}

function buildWechatOAuthAuthorizeUrl(
  authorizeUrl: string,
  context: {
    paymentType: string
    orderType: OrderType
    planId?: number
    orderAmount: number
  },
): string {
  const normalizedUrl = authorizeUrl.trim()
  if (!normalizedUrl || typeof window === 'undefined') {
    return normalizedUrl
  }

  try {
    const targetUrl = new URL(normalizedUrl, window.location.origin)
    const redirectPath = targetUrl.searchParams.get('redirect') || '/purchase'
    const redirectUrl = new URL(redirectPath, window.location.origin)
    const paymentType =
      normalizeVisibleMethod(context.paymentType) ||
      context.paymentType.trim() ||
      'wxpay'

    redirectUrl.searchParams.set('payment_type', paymentType)
    redirectUrl.searchParams.set('order_type', context.orderType)

    if (context.planId) {
      redirectUrl.searchParams.set('plan_id', String(context.planId))
    } else {
      redirectUrl.searchParams.delete('plan_id')
    }

    if (context.orderAmount > 0) {
      redirectUrl.searchParams.set('amount', String(context.orderAmount))
    } else {
      redirectUrl.searchParams.delete('amount')
    }

    targetUrl.searchParams.set(
      'redirect',
      `${redirectUrl.pathname}${redirectUrl.search}`,
    )
    return targetUrl.toString()
  } catch {
    return normalizedUrl
  }
}

function onPaymentDone() {
  resetPayment()
}

function onPaymentSuccess() {
  removeRecoverySnapshot()
  authStore.refreshUser()
}

function onPaymentSettled() {
  removeRecoverySnapshot()
}

const checkout = ref<CheckoutInfoResponse>({
  methods: {},
  global_min: 0,
  global_max: 0,
  plans: [],
  balance_disabled: false,
  balance_recharge_multiplier: 1,
  recharge_fee_rate: 0,
  recharge_campaign_enabled: false,
  recharge_campaign_amount: 100,
  recharge_campaign_bonus_rate: 10,
  help_text: '',
  help_image_url: '',
  stripe_publishable_key: '',
})

async function scrollToRedeem() {
  await router.replace({
    path: '/purchase',
    query: { ...route.query, redeem: '1' },
  })
  await nextTick()
  document
    .getElementById('purchase-redeem-section')
    ?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function openExternalPurchase() {
  if (typeof window === 'undefined' || !externalPurchaseUrl.value) {
    return
  }
  window.open(externalPurchaseUrl.value, '_blank', 'noopener,noreferrer')
}

const visibleMethods = computed(() => getVisibleMethods(checkout.value.methods))
const enabledMethods = computed(() => Object.keys(visibleMethods.value))
const rechargeUnavailable = computed(
  () => checkout.value.balance_disabled || enabledMethods.value.length === 0,
)
const validAmount = computed(() => amount.value ?? 0)
const balanceRechargeMultiplier = computed(() => {
  const multiplier = checkout.value.balance_recharge_multiplier
  return multiplier > 0 ? multiplier : 1
})
const creditedAmount = computed(() =>
  Math.round(validAmount.value * balanceRechargeMultiplier.value * 100) / 100,
)
const campaignEnabled = computed(() => checkout.value.recharge_campaign_enabled)
const campaignThreshold = computed(() =>
  checkout.value.recharge_campaign_amount > 0
    ? checkout.value.recharge_campaign_amount
    : 100,
)
const campaignBonusRate = computed(() =>
  checkout.value.recharge_campaign_bonus_rate >= 0
    ? checkout.value.recharge_campaign_bonus_rate
    : 10,
)
const campaignBonusAmount = computed(() => {
  if (!campaignEnabled.value || validAmount.value < campaignThreshold.value) {
    return 0
  }
  return (
    Math.round(
      (creditedAmount.value * campaignBonusRate.value) / 100 * 100,
    ) / 100
  )
})
const totalCreditedAmount = computed(
  () => Math.round((creditedAmount.value + campaignBonusAmount.value) * 100) / 100,
)

function amountFitsMethod(amt: number, methodType: string): boolean {
  if (amt <= 0) return true
  const limit = visibleMethods.value[methodType]
  if (!limit) return false
  if (limit.single_min > 0 && amt < limit.single_min) return false
  if (limit.single_max > 0 && amt > limit.single_max) return false
  return true
}

const globalMinAmount = computed(() => {
  const limits = Object.values(visibleMethods.value)
  if (limits.length === 0) return 0
  if (limits.some((limit) => limit.single_min <= 0)) return 0
  return Math.min(...limits.map((limit) => limit.single_min))
})

const globalMaxAmount = computed(() => {
  const limits = Object.values(visibleMethods.value)
  if (limits.length === 0) return 0
  if (limits.some((limit) => limit.single_max <= 0)) return 0
  return Math.max(...limits.map((limit) => limit.single_max))
})

const selectedLimit = computed(() => visibleMethods.value[selectedMethod.value])
const selectedCurrency = computed(() =>
  normalizePaymentCurrency(selectedLimit.value?.currency),
)
const localeCode = computed(() => {
  const raw = i18n.locale as unknown
  if (typeof raw === 'string') return raw
  if (raw && typeof raw === 'object' && 'value' in raw) {
    return String((raw as { value?: string }).value || '')
  }
  return undefined
})

function formatSelectedPaymentAmount(value: number): string {
  return formatPaymentAmount(value, selectedCurrency.value, localeCode.value)
}

const methodOptions = computed<PaymentMethodOption[]>(() =>
  enabledMethods.value.map((type) => {
    const limit = visibleMethods.value[type]
    return {
      type,
      fee_rate: limit?.fee_rate ?? 0,
      available:
        !rechargeUnavailable.value &&
        limit?.available !== false &&
        amountFitsMethod(validAmount.value, type),
    }
  }),
)

const feeRate = computed(() => checkout.value.recharge_fee_rate ?? 0)
const feeAmount = computed(() =>
  feeRate.value > 0 && validAmount.value > 0
    ? Math.ceil((validAmount.value * feeRate.value) / 100 * 100) / 100
    : 0,
)
const totalAmount = computed(() =>
  feeRate.value > 0 && validAmount.value > 0
    ? Math.round((validAmount.value + feeAmount.value) * 100) / 100
    : validAmount.value,
)

const amountError = computed(() => {
  if (validAmount.value <= 0) return ''
  if (!enabledMethods.value.some((method) => amountFitsMethod(validAmount.value, method))) {
    return t('payment.amountNoMethod')
  }
  const limit = selectedLimit.value
  if (limit) {
    if (limit.single_min > 0 && validAmount.value < limit.single_min) {
      return t('payment.amountTooLow', {
        min: formatSelectedPaymentAmount(limit.single_min),
      })
    }
    if (limit.single_max > 0 && validAmount.value > limit.single_max) {
      return t('payment.amountTooHigh', {
        max: formatSelectedPaymentAmount(limit.single_max),
      })
    }
  }
  return ''
})

const canSubmit = computed(
  () =>
    !rechargeUnavailable.value &&
    validAmount.value > 0 &&
    amountFitsMethod(validAmount.value, selectedMethod.value) &&
    selectedLimit.value?.available !== false,
)

watch(
  () => [validAmount.value, selectedMethod.value] as const,
  ([amt, method]) => {
    if (amt <= 0 || amountFitsMethod(amt, method)) return
    const available = enabledMethods.value.find((m) => amountFitsMethod(amt, m))
    if (available) selectedMethod.value = available
  },
)

const paymentButtonClass = computed(() => {
  const method = selectedMethod.value
  if (!method) return 'btn-primary'
  if (method.includes('alipay')) return 'btn-alipay'
  if (method.includes('wxpay')) return 'btn-wxpay'
  if (method === 'stripe') return 'btn-stripe'
  if (method === 'airwallex') return 'btn-airwallex'
  return 'btn-primary'
})

async function handleSubmitRecharge() {
  if (!canSubmit.value || submitting.value) return
  await createOrder(validAmount.value, 'balance')
}

async function createOrder(
  orderAmount: number,
  orderType: OrderType,
  planId?: number,
  options: CreateOrderOptions = {},
) {
  submitting.value = true
  errorMessage.value = ''
  errorHintMessage.value = ''
  const requestType =
    normalizeVisibleMethod(options.paymentType || selectedMethod.value) ||
    options.paymentType ||
    selectedMethod.value

  try {
    const payload = buildCreateOrderPayload({
      amount: orderAmount,
      paymentType: requestType,
      orderType,
      planId,
      origin: typeof window !== 'undefined' ? window.location.origin : '',
      isMobile: isMobileDevice(),
      isWechatBrowser:
        typeof window !== 'undefined' &&
        /MicroMessenger/i.test(window.navigator.userAgent),
      forceQRCode:
        !!(
          checkout.value.alipay_force_qrcode &&
          normalizeVisibleMethod(requestType) === 'alipay'
        ),
    })
    if (options.openid) payload.openid = options.openid
    if (options.wechatResumeToken) {
      payload.wechat_resume_token = options.wechatResumeToken
    }

    const result = (await paymentStore.createOrder(
      payload,
    )) as CreateOrderResult & { resume_token?: string }

    const openWindow = (url: string) => {
      const win = window.open(url, 'paymentPopup', getPaymentPopupFeatures())
      if (!win || win.closed) {
        window.location.href = url
      }
    }

    const visibleMethod = normalizeVisibleMethod(requestType) || requestType
    const stripeMethod =
      visibleMethod === 'stripe'
        ? ''
        : visibleMethod === 'wxpay'
          ? 'wechat_pay'
          : 'alipay'

    const stripeRouteUrl =
      result.client_secret && visibleMethod !== 'airwallex'
        ? router.resolve({
            path: '/payment/stripe',
            query: {
              order_id: String(result.order_id),
              client_secret: result.client_secret,
              method: stripeMethod || undefined,
              resume_token: result.resume_token || undefined,
            },
          }).href
        : ''

    const airwallexRouteUrl =
      result.client_secret && result.intent_id
        ? router.resolve({
            path: '/payment/airwallex',
            query: {
              order_id: String(result.order_id),
              out_trade_no: result.out_trade_no || undefined,
              resume_token: result.resume_token || undefined,
            },
          }).href
        : ''

    const decision = decidePaymentLaunch(result, {
      visibleMethod,
      orderType,
      isMobile: isMobileDevice(),
      isWechatBrowser:
        typeof window !== 'undefined' &&
        /MicroMessenger/i.test(window.navigator.userAgent),
      forceQRCode:
        !!(
          checkout.value.alipay_force_qrcode && visibleMethod === 'alipay'
        ),
      stripePopupUrl: stripeRouteUrl,
      stripeRouteUrl,
      airwallexRouteUrl,
    })

    if (decision.kind === 'wechat_oauth' && decision.oauth?.authorize_url) {
      window.location.href = buildWechatOAuthAuthorizeUrl(
        decision.oauth.authorize_url,
        {
          paymentType: visibleMethod,
          orderType,
          planId,
          orderAmount,
        },
      )
      return
    }

    if (decision.kind === 'unhandled') {
      applyScenarioError({ reason: 'UNHANDLED_PAYMENT_SCENARIO' }, visibleMethod)
      return
    }

    paymentState.value = decision.paymentState
    paymentPhase.value = 'paying'
    persistRecoverySnapshot(decision.recovery)

    if (decision.kind === 'stripe_popup') {
      openWindow(decision.paymentState.payUrl)
      return
    }
    if (decision.kind === 'stripe_route' || decision.kind === 'airwallex_route') {
      window.location.href = decision.paymentState.payUrl
      return
    }
    if (decision.kind === 'wechat_jsapi' && decision.jsapi) {
      try {
        const jsapiResult = await invokeWechatJsapiPayment(
          decision.jsapi as Record<string, unknown>,
        )
        const errMsg = String(jsapiResult.err_msg || '').toLowerCase()
        if (errMsg.includes('cancel')) {
          appStore.showInfo(t('payment.qr.cancelled'))
          resetPayment()
        } else if (errMsg && !errMsg.includes('ok')) {
          resetPayment()
          const fallbackApplied = await attemptMobileQrFallback(
            { reason: 'WECHAT_JSAPI_FAILED', message: errMsg },
            {
              orderAmount,
              orderType,
              planId,
              paymentType: visibleMethod,
              attempted: options.mobileQrFallbackAttempted === true,
            },
          )
          if (!fallbackApplied) {
            applyScenarioError(
              { reason: 'WECHAT_JSAPI_FAILED', message: errMsg },
              visibleMethod,
            )
          }
        } else {
          const resultState = { ...decision.paymentState }
          resetPayment()
          await redirectToPaymentResult(resultState)
        }
      } catch (err: unknown) {
        resetPayment()
        const fallbackApplied = await attemptMobileQrFallback(err, {
          orderAmount,
          orderType,
          planId,
          paymentType: visibleMethod,
          attempted: options.mobileQrFallbackAttempted === true,
        })
        if (!fallbackApplied) {
          throw err
        }
      }
      return
    }
    if (decision.kind === 'redirect_waiting' && decision.paymentState.payUrl) {
      if (isMobileDevice()) {
        window.location.href = decision.paymentState.payUrl
        return
      }
      openWindow(decision.paymentState.payUrl)
    }
  } catch (err: unknown) {
    const apiErr = err as Record<string, unknown>
    if (apiErr.reason === 'TOO_MANY_PENDING') {
      const metadata = apiErr.metadata as Record<string, unknown> | undefined
      errorMessage.value = t('payment.errors.tooManyPending', {
        max: metadata?.max || '',
      })
      errorHintMessage.value = ''
    } else if (apiErr.reason === 'CANCEL_RATE_LIMITED') {
      errorMessage.value = t('payment.errors.cancelRateLimited')
      errorHintMessage.value = ''
    } else if (
      await attemptMobileQrFallback(err, {
        orderAmount,
        orderType,
        planId,
        paymentType: requestType,
        attempted: options.mobileQrFallbackAttempted === true,
      })
    ) {
      return
    } else {
      const handled = applyScenarioError(
        err,
        normalizeVisibleMethod(options.paymentType || selectedMethod.value) ||
          selectedMethod.value,
      )
      if (!handled) {
        errorMessage.value = extractI18nErrorMessage(
          err,
          t,
          'payment.errors',
          extractApiErrorMessage(err, t('payment.result.failed')),
        )
        errorHintMessage.value = ''
      }
      if (handled) return
    }
    appStore.showError(
      buildPaymentErrorToastMessage(errorMessage.value, errorHintMessage.value),
    )
  } finally {
    submitting.value = false
  }
}

interface MobileQrFallbackContext {
  orderAmount: number
  orderType: OrderType
  planId?: number
  paymentType: string
  attempted: boolean
}

function shouldFallbackToDesktopQr(
  err: unknown,
  paymentMethod: string,
  attempted: boolean,
): boolean {
  if (attempted || !isMobileDevice()) {
    return false
  }

  const normalizedMethod = normalizeVisibleMethod(paymentMethod) || paymentMethod
  const reason =
    typeof err === 'object' &&
    err &&
    'reason' in err &&
    typeof err.reason === 'string'
      ? err.reason
      : ''
  const message =
    err instanceof Error
      ? err.message
      : typeof err === 'object' &&
          err &&
          'message' in err &&
          typeof err.message === 'string'
        ? err.message
        : ''
  const normalizedMessage = message.toLowerCase()

  if (normalizedMethod === 'wxpay') {
    return (
      reason === 'WECHAT_H5_NOT_AUTHORIZED' ||
      reason === 'WECHAT_PAYMENT_MP_NOT_CONFIGURED' ||
      reason === 'WECHAT_JSAPI_FAILED' ||
      reason === 'PAYMENT_GATEWAY_ERROR' ||
      reason === 'UNHANDLED_PAYMENT_SCENARIO' ||
      normalizedMessage.includes('weixinjsbridge is unavailable') ||
      normalizedMessage.includes('wechat_jsapi_unavailable')
    )
  }

  if (normalizedMethod === 'alipay') {
    return (
      reason === 'PAYMENT_GATEWAY_ERROR' ||
      reason === 'UNHANDLED_PAYMENT_SCENARIO'
    )
  }

  return false
}

async function attemptMobileQrFallback(
  err: unknown,
  context: MobileQrFallbackContext,
): Promise<boolean> {
  if (!shouldFallbackToDesktopQr(err, context.paymentType, context.attempted)) {
    return false
  }

  try {
    const visibleMethod =
      normalizeVisibleMethod(context.paymentType) || context.paymentType
    const payload = buildCreateOrderPayload({
      amount: context.orderAmount,
      paymentType: visibleMethod,
      orderType: context.orderType,
      planId: context.planId,
      origin: typeof window !== 'undefined' ? window.location.origin : '',
      isMobile: false,
      isWechatBrowser: false,
    })
    const result = (await paymentStore.createOrder(
      payload,
    )) as CreateOrderResult & { resume_token?: string }
    const stripeMethod = visibleMethod === 'wxpay' ? 'wechat_pay' : 'alipay'
    const stripeRouteUrl = result.client_secret
      ? router.resolve({
          path: '/payment/stripe',
          query: {
            order_id: String(result.order_id),
            client_secret: result.client_secret,
            method: stripeMethod,
            resume_token: result.resume_token || undefined,
          },
        }).href
      : ''
    const decision = decidePaymentLaunch(result, {
      visibleMethod,
      orderType: context.orderType,
      isMobile: false,
      isWechatBrowser: false,
      stripePopupUrl: stripeRouteUrl,
      stripeRouteUrl,
    })

    if (decision.kind !== 'qr_waiting' || !decision.paymentState.qrCode) {
      return false
    }

    errorMessage.value = ''
    errorHintMessage.value = ''
    paymentState.value = decision.paymentState
    paymentPhase.value = 'paying'
    persistRecoverySnapshot(decision.recovery)
    appStore.showWarning(t('payment.errors.mobilePaymentFallbackToQr'))
    return true
  } catch {
    return false
  }
}

function applyScenarioError(err: unknown, paymentMethod: string): boolean {
  const descriptor = describePaymentScenarioError(err, {
    paymentMethod,
    isMobile: isMobileDevice(),
    isWechatBrowser:
      typeof window !== 'undefined' &&
      /MicroMessenger/i.test(window.navigator.userAgent),
  })
  if (!descriptor) {
    errorMessage.value = ''
    errorHintMessage.value = ''
    return false
  }
  errorMessage.value = t(descriptor.messageKey)
  errorHintMessage.value = descriptor.hintKey ? t(descriptor.hintKey) : ''
  appStore.showError(
    buildPaymentErrorToastMessage(errorMessage.value, errorHintMessage.value),
  )
  return true
}

async function resumeWechatPaymentFromQuery() {
  const resume = parseWechatResumeRoute(route.query, checkout.value.plans, validAmount.value)
  if (!resume) {
    return
  }

  if (resume.orderType !== 'balance') {
    await router.replace({
      path: route.path,
      query: stripWechatResumeQuery(route.query),
    })
    return
  }

  selectedMethod.value = resume.paymentType
  if (resume.orderType === 'balance' && resume.orderAmount > 0) {
    amount.value = resume.orderAmount
  }

  await router.replace({
    path: route.path,
    query: stripWechatResumeQuery(route.query),
  })

  if (resume.wechatResumeToken) {
    await createOrder(0, resume.orderType, resume.planId, {
      wechatResumeToken: resume.wechatResumeToken,
      paymentType: resume.paymentType,
      isResume: true,
    })
    return
  }

  if (resume.orderAmount > 0 && resume.openid) {
    await createOrder(resume.orderAmount, resume.orderType, resume.planId, {
      openid: resume.openid,
      paymentType: resume.paymentType,
      isResume: true,
    })
  }
}

onMounted(async () => {
  try {
    await appStore.fetchPublicSettings().catch(() => null)
    checkoutUnavailable.value = false
    const res = await paymentAPI.getCheckoutInfo()
    checkout.value = res.data

    if (enabledMethods.value.length) {
      const methodOrder = METHOD_ORDER as readonly string[]
      const sorted = [...enabledMethods.value].sort((a, b) => {
        const ai = methodOrder.indexOf(a)
        const bi = methodOrder.indexOf(b)
        return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
      })
      selectedMethod.value = sorted[0]
    }

    if (typeof window !== 'undefined') {
      if (hasWechatResumeQuery(route.query)) {
        removeRecoverySnapshot()
      }
      const routeResumeToken =
        typeof route.query.resume_token === 'string'
          ? route.query.resume_token
          : typeof route.query.wechat_resume_token === 'string'
            ? route.query.wechat_resume_token
            : undefined
      const restored = readPaymentRecoverySnapshot(
        window.localStorage.getItem(PAYMENT_RECOVERY_STORAGE_KEY),
        { resumeToken: routeResumeToken },
      )
      if (restored) {
        paymentState.value = restored
        paymentPhase.value = 'paying'
        const restoredMethod = normalizeVisibleMethod(restored.paymentType)
        if (restoredMethod) {
          selectedMethod.value = restoredMethod
        }
      } else {
        removeRecoverySnapshot()
      }
    }

    await resumeWechatPaymentFromQuery()

    if (route.query.redeem === '1') {
      await nextTick()
      document
        .getElementById('purchase-redeem-section')
        ?.scrollIntoView({ behavior: 'auto', block: 'start' })
    }
  } catch (err: unknown) {
    checkoutUnavailable.value = true
    appStore.showError(
      extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')),
    )
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 180ms ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.payment-shell {
  position: relative;
}
.payment-redeem-heading span,
.payment-shop-copy span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.payment-redeem-heading h3,
.payment-shop-copy h3,
.payment-unavailable h2 {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-weight: 600;
}

.payment-shop {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.payment-shop-copy h3 {
  margin-top: 0.32rem;
  font-size: clamp(1.04rem, 1.45vw, 1.28rem);
}

.payment-shop-copy p {
  margin-top: 0.45rem;
  max-width: 40rem;
  color: #59645a;
  font-size: 0.92rem;
  line-height: 1.7;
}

.payment-shop-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.payment-recharge-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.1rem;
}

.payment-balance-chip {
  display: grid;
  gap: 0.18rem;
  min-width: 10rem;
  padding: 0.78rem 0.95rem;
  border: 1px solid rgba(198, 184, 157, 0.46);
  border-radius: 12px;
  background: rgba(255, 252, 245, 0.72);
  text-align: right;
}

.payment-balance-chip span {
  color: #7b6a53;
  font-size: 0.72rem;
  letter-spacing: 0.06em;
}

.payment-balance-chip strong {
  color: #1f2320;
  font-size: 1.2rem;
  font-weight: 700;
}

.payment-unavailable {
  display: grid;
  justify-items: center;
  gap: 0.75rem;
}

.payment-unavailable h2 {
  font-size: clamp(1.28rem, 2vw, 1.72rem);
}

.payment-unavailable p {
  max-width: 30rem;
  color: #59645a;
  font-size: 0.9rem;
  line-height: 1.7;
}

.payment-unavailable div {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.6rem;
  margin-top: 0.5rem;
}

.payment-inline-error {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 0.85rem;
  border: 1px solid rgba(167, 58, 42, 0.22);
  border-radius: 10px;
  background: rgba(167, 58, 42, 0.055);
  color: #7b2f25;
  padding: 0.95rem 1rem;
}

.payment-inline-error strong {
  display: block;
  color: #7b2f25;
  font-size: 0.9rem;
}

.payment-inline-error p {
  margin-top: 0.18rem;
  color: #7b6a53;
  font-size: 0.78rem;
  line-height: 1.55;
}

.payment-redeem-section {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(198, 184, 157, 0.32);
}

.payment-redeem-heading {
  margin-bottom: 0.85rem;
}

.payment-redeem-heading h3 {
  margin-top: 0.28rem;
  font-size: clamp(1.02rem, 1.45vw, 1.28rem);
}

.dark .payment-redeem-heading h3,
.dark .payment-shop-copy h3,
.dark .payment-unavailable h2 {
  color: #f4efe4;
}

.dark .payment-redeem-heading span,
.dark .payment-shop-copy span {
  color: #879186;
}

.dark .payment-shop-copy p,
.dark .payment-unavailable p {
  color: #879186;
}

.dark .payment-balance-chip {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .payment-balance-chip span {
  color: #879186;
}

.dark .payment-balance-chip strong {
  color: #f4efe4;
}

.dark .payment-inline-error {
  border-color: rgba(167, 58, 42, 0.38);
  background: rgba(167, 58, 42, 0.12);
  color: #f0b4a8;
}

.dark .payment-inline-error strong {
  color: #f0b4a8;
}

.dark .payment-inline-error p {
  color: #879186;
}

.dark .payment-redeem-section {
  border-top-color: rgba(48, 52, 43, 0.95);
}

@media (max-width: 768px) {
  .payment-shop {
    flex-direction: column;
    align-items: flex-start;
  }

  .payment-recharge-head {
    flex-direction: column;
  }

  .payment-balance-chip {
    width: 100%;
    min-width: 0;
    text-align: left;
  }
}

@media (max-width: 640px) {
  .payment-inline-error {
    align-items: flex-start;
    grid-template-columns: 1fr;
  }
}
</style>
