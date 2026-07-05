<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      :disabled="switching"
      class="locale-switcher-btn"
      :class="variant === 'public' ? 'locale-switcher-btn-public' : 'locale-switcher-btn-default'"
      :title="currentLocale?.name"
    >
      <span v-if="variant !== 'public'" class="text-base">{{ currentLocale?.flag }}</span>
      <template v-if="variant === 'public'">
        <span class="locale-switcher-public-mark" aria-hidden="true">
          <Icon name="paperScroll" size="sm" :stroke-width="1.55" />
        </span>
        <span class="locale-switcher-public-name">{{ triggerLabel }}</span>
      </template>
      <span v-else>{{ triggerLabel }}</span>
      <Icon
        name="chevronDown"
        size="xs"
        class="locale-switcher-chevron"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <transition name="dropdown">
      <div
        v-if="isOpen"
        class="absolute right-0 z-50 mt-1 overflow-hidden rounded-lg shadow-lg"
        :class="variant === 'public'
          ? 'w-36 border border-zen-paperLine/80 bg-[#fbf7ef] dark:border-zen-nightLine dark:bg-[#1b1d18]'
          : 'w-32 border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800'"
      >
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          :disabled="switching"
          @click="selectLocale(locale.code)"
          class="flex w-full items-center gap-2 px-3 py-2 text-sm transition-colors"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-900/20 dark:text-primary-400':
              locale.code === currentLocaleCode && variant !== 'public',
            'bg-[rgba(216,205,185,0.28)] text-[#a73a2a] dark:bg-[rgba(216,205,185,0.08)] dark:text-[#ffd8bb]':
              locale.code === currentLocaleCode && variant === 'public',
            'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-dark-700':
              variant !== 'public',
            'text-[#38413a] hover:bg-[rgba(216,205,185,0.18)] dark:text-[#d8cdb9] dark:hover:bg-[rgba(216,205,185,0.08)]':
              variant === 'public'
          }"
        >
          <span v-if="variant !== 'public'" class="text-base">{{ locale.flag }}</span>
          <span>{{ locale.name }}</span>
          <Icon v-if="locale.code === currentLocaleCode" name="check" size="sm" class="ml-auto text-primary-500" />
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { setLocale, availableLocales } from '@/i18n'

const props = withDefaults(defineProps<{
  variant?: 'default' | 'public'
}>(), {
  variant: 'default',
})

const { locale } = useI18n()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const switching = ref(false)

const currentLocaleCode = computed(() => locale.value)
const currentLocale = computed(() => availableLocales.find((l) => l.code === locale.value))
const triggerLabel = computed(() => {
  if (!currentLocale.value) return ''
  return props.variant === 'public'
    ? currentLocale.value.name
    : currentLocale.value.code.toUpperCase()
})

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

async function selectLocale(code: string) {
  if (switching.value || code === currentLocaleCode.value) {
    isOpen.value = false
    return
  }
  switching.value = true
  try {
    await setLocale(code)
    isOpen.value = false
  } finally {
    switching.value = false
  }
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.locale-switcher-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  transition: color 180ms ease, background-color 180ms ease;
}

.locale-switcher-btn-default {
  border-radius: 0.5rem;
  padding: 0.375rem 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(75 85 99);
}

.locale-switcher-btn-default:hover {
  background: rgb(243 244 246);
}

.locale-switcher-btn-public {
  min-height: var(--sst-public-tool-height, 2.28rem);
  gap: var(--sst-public-tool-gap, 0.42rem);
  border-radius: 999px;
  border: 1px solid var(--sst-public-tool-border, rgba(145, 128, 99, 0.14));
  padding: 0.22rem 0.58rem 0.22rem 0.3rem;
  background:
    linear-gradient(180deg, var(--sst-public-tool-bg-top, rgba(255, 252, 247, 0.7)), var(--sst-public-tool-bg-bottom, rgba(244, 237, 226, 0.82)));
  color: var(--sst-public-tool-fg, #343b35);
  box-shadow: var(--sst-public-tool-shadow, inset 0 1px 0 rgba(255, 255, 255, 0.58), 0 8px 18px rgba(74, 56, 33, 0.06));
}

.locale-switcher-btn-public:hover {
  color: #9f3d2f;
  border-color: var(--sst-public-tool-hover-border, rgba(167, 58, 42, 0.18));
  background:
    linear-gradient(180deg, var(--sst-public-tool-hover-top, rgba(255, 251, 245, 0.94)), var(--sst-public-tool-hover-bottom, rgba(242, 232, 218, 0.98)));
  box-shadow: var(--sst-public-tool-hover-shadow, inset 0 1px 0 rgba(255, 255, 255, 0.62), 0 10px 20px rgba(167, 58, 42, 0.08));
}

.locale-switcher-public-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: calc(var(--sst-public-tool-mark-size, 1.62rem) + 0.02rem);
  height: calc(var(--sst-public-tool-mark-size, 1.62rem) + 0.02rem);
  border-radius: var(--sst-public-tool-mark-radius, 0.58rem);
  background:
    linear-gradient(180deg, var(--sst-public-tool-mark-top, rgba(244, 238, 229, 0.98)), var(--sst-public-tool-mark-bottom, rgba(233, 222, 205, 0.94)));
  color: var(--sst-public-tool-mark-fg, #8d6845);
  box-shadow: var(--sst-public-tool-mark-shadow, inset 0 1px 0 rgba(255, 255, 255, 0.6), inset 0 0 0 1px rgba(189, 161, 126, 0.12));
}

.locale-switcher-public-mark::after {
  content: '';
  position: absolute;
  left: 0.34rem;
  top: 0.34rem;
  width: var(--sst-public-tool-dot-size, 0.16rem);
  height: var(--sst-public-tool-dot-size, 0.16rem);
  border-radius: 999px;
  background: var(--sst-public-tool-dot, rgba(167, 58, 42, 0.72));
  box-shadow: var(--sst-public-tool-dot-shadow, 0 0 0 0.13rem rgba(167, 58, 42, 0.06));
}

.locale-switcher-public-name {
  min-width: 1.88rem;
  font-size: 0.9rem;
  line-height: 1;
  letter-spacing: 0.02em;
  color: var(--sst-public-tool-label, #2f332e);
  text-align: left;
}

@media (max-width: 767px) {
  .locale-switcher-btn-public {
    padding: 0.18rem 0.46rem 0.18rem 0.26rem;
  }

  .locale-switcher-public-name {
    min-width: 1.56rem;
    font-size: 0.8rem;
  }
}

.locale-switcher-chevron {
  color: var(--sst-public-tool-chevron, #aa9274);
  transition: transform 0.2s ease;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>

<style>
html.dark .locale-switcher-btn-default {
  color: rgb(209 213 219);
}

html.dark .locale-switcher-btn-default:hover {
  background: rgb(55 65 81 / 0.5);
}

html.dark .public-site-header .locale-switcher-btn-public {
  color: #e6dac7;
  border-color: rgba(112, 101, 81, 0.28);
  background: linear-gradient(180deg, rgba(31, 32, 28, 0.82), rgba(22, 24, 20, 0.92));
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.04),
    0 5px 12px rgba(0, 0, 0, 0.18);
  text-shadow: none;
}

html.dark .sst-public-tone-legal .locale-switcher-btn-public {
  color: #f0dfc3 !important;
  text-shadow: 0 1px 0 rgba(0, 0, 0, 0.22);
}

html.dark .public-site-header .locale-switcher-btn-public:hover {
  color: #f0e4d2;
  border-color: rgba(144, 118, 86, 0.34);
  background: linear-gradient(180deg, rgba(38, 36, 31, 0.88), rgba(29, 28, 24, 0.94));
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 7px 14px rgba(0, 0, 0, 0.22);
}

html.dark .sst-public-tone-legal .locale-switcher-btn-public:hover {
  color: #f1c27c;
  background: rgba(205, 163, 103, 0.1);
}

html.dark .public-site-header .locale-switcher-chevron {
  color: #b99868;
}

html.dark .public-site-header .locale-switcher-public-mark {
  background: linear-gradient(180deg, rgba(58, 50, 41, 0.8), rgba(42, 38, 31, 0.88));
  color: #d0ab7a;
  box-shadow:
    inset 0 1px 0 rgba(255, 251, 243, 0.06),
    inset 0 0 0 1px rgba(201, 167, 122, 0.08);
}

html.dark .public-site-header .locale-switcher-public-mark::after {
  background: rgba(177, 95, 72, 0.68);
  box-shadow: 0 0 0 0.12rem rgba(177, 95, 72, 0.06);
}

html.dark .public-site-header .locale-switcher-public-name {
  color: #e8dcc9;
}

html.dark .sst-public-tone-legal .locale-switcher-chevron {
  color: #e2c89e !important;
}
</style>
