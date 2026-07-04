<template>
  <span
    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
    :class="statusClass"
  >
    {{ statusLabel }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { OrderStatus } from '@/types/payment'

const props = defineProps<{
  status: OrderStatus
}>()

const { t } = useI18n()

const statusMap: Record<OrderStatus, { key: string; class: string }> = {
  PENDING: { key: 'payment.status.pending', class: 'border border-[#9b8155]/25 bg-[#9b8155]/12 text-[#7b6a53] dark:bg-[#9b8155]/20 dark:text-amber-300' },
  PAID: { key: 'payment.status.paid', class: 'border border-[#51624f]/25 bg-[#51624f]/10 text-[#51624f] dark:bg-[#51624f]/20 dark:text-emerald-300' },
  RECHARGING: { key: 'payment.status.recharging', class: 'border border-[#51624f]/25 bg-[#51624f]/10 text-[#51624f] dark:bg-[#51624f]/20 dark:text-emerald-300' },
  COMPLETED: { key: 'payment.status.completed', class: 'border border-[#51624f]/25 bg-[#51624f]/10 text-[#51624f] dark:bg-[#51624f]/20 dark:text-emerald-300' },
  EXPIRED: { key: 'payment.status.expired', class: 'border border-stone-300 bg-stone-100 text-[#667066] dark:border-dark-600 dark:bg-dark-800 dark:text-gray-400' },
  CANCELLED: { key: 'payment.status.cancelled', class: 'border border-stone-300 bg-stone-100 text-[#667066] dark:border-dark-600 dark:bg-dark-800 dark:text-gray-400' },
  FAILED: { key: 'payment.status.failed', class: 'border border-[#a73a2a]/25 bg-[#a73a2a]/10 text-[#a73a2a] dark:bg-[#a73a2a]/20 dark:text-red-300' },
  REFUND_REQUESTED: { key: 'payment.status.refund_requested', class: 'border border-[#9b8155]/25 bg-[#9b8155]/12 text-[#7b6a53] dark:bg-[#9b8155]/20 dark:text-amber-300' },
  REFUNDING: { key: 'payment.status.refunding', class: 'border border-[#9b8155]/25 bg-[#9b8155]/12 text-[#7b6a53] dark:bg-[#9b8155]/20 dark:text-amber-300' },
  REFUNDED: { key: 'payment.status.refunded', class: 'border border-[#51624f]/25 bg-[#51624f]/10 text-[#51624f] dark:bg-[#51624f]/20 dark:text-emerald-300' },
  PARTIALLY_REFUNDED: { key: 'payment.status.partially_refunded', class: 'border border-[#9b8155]/25 bg-[#9b8155]/12 text-[#7b6a53] dark:bg-[#9b8155]/20 dark:text-amber-300' },
  REFUND_FAILED: { key: 'payment.status.refund_failed', class: 'border border-[#a73a2a]/25 bg-[#a73a2a]/10 text-[#a73a2a] dark:bg-[#a73a2a]/20 dark:text-red-300' },
}

const statusLabel = computed(() => {
  const entry = statusMap[props.status]
  return entry ? t(entry.key) : props.status
})

const statusClass = computed(() => {
  const entry = statusMap[props.status]
  return entry?.class ?? 'border border-stone-300 bg-stone-100 text-[#667066] dark:border-dark-600 dark:bg-dark-800 dark:text-gray-400'
})
</script>
