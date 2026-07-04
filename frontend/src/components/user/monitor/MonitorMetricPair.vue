<template>
  <div class="mt-4 grid grid-cols-2 gap-2">
    <div
      class="monitor-metric rounded-lg border border-stone-200/70 bg-stone-100/35 p-2.5 dark:border-[#44473a]/80"
    >
      <div
        class="flex items-center gap-1.5 text-[10px] font-medium text-gray-500"
      >
        <Icon :name="primaryIcon" size="xs" />
        <span>{{ primaryLabel }}</span>
      </div>
      <div class="mt-1.5 text-[1.05rem] font-semibold font-mono tabular-nums text-gray-900 dark:text-gray-100">
        {{ primaryValue }}<span class="text-[11px] font-normal text-gray-400 ml-0.5">{{ primaryUnit }}</span>
      </div>
    </div>
    <div
      class="monitor-metric rounded-lg border border-stone-200/70 bg-stone-100/35 p-2.5 dark:border-[#44473a]/80"
    >
      <div
        class="flex items-center gap-1.5 text-[10px] font-medium text-gray-500"
      >
        <Icon :name="secondaryIcon" size="xs" />
        <span>{{ secondaryLabel }}</span>
      </div>
      <div class="mt-1.5 text-[1.05rem] font-semibold font-mono tabular-nums text-gray-900 dark:text-gray-100">
        {{ secondaryValue }}<span class="text-[11px] font-normal text-gray-400 ml-0.5">{{ secondaryUnit }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/icons/Icon.vue'

defineProps<{
  primaryLabel: string
  primaryValue: string
  primaryUnit: string
  primaryIcon: 'bolt' | 'globe' | 'clock' | 'link'
  secondaryLabel: string
  secondaryValue: string
  secondaryUnit: string
  secondaryIcon: 'bolt' | 'globe' | 'clock' | 'link'
}>()
</script>

<style scoped>
.dark .monitor-metric {
  background: rgba(24, 26, 21, 0.72);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.025);
}
</style>
