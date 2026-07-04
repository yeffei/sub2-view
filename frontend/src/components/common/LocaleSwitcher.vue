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
      <span>{{ triggerLabel }}</span>
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
  min-height: 2.25rem;
  border-radius: 999px;
  padding: 0.3rem 0.7rem;
  font-size: 0.9rem;
  color: #3b433d;
}

.locale-switcher-btn-public:hover {
  color: #a73a2a;
  background: rgba(216, 205, 185, 0.18);
}

.locale-switcher-chevron {
  color: rgb(156 163 175);
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
  color: #ece2d2;
  text-shadow: 0 1px 0 rgba(0, 0, 0, 0.18);
}

html.dark .sst-public-tone-legal .locale-switcher-btn-public {
  color: #f0dfc3 !important;
  text-shadow: 0 1px 0 rgba(0, 0, 0, 0.22);
}

html.dark .public-site-header .locale-switcher-btn-public:hover {
  color: #ffd8bb;
  background: rgba(216, 205, 185, 0.08);
}

html.dark .sst-public-tone-legal .locale-switcher-btn-public:hover {
  color: #f1c27c;
  background: rgba(205, 163, 103, 0.1);
}

html.dark .public-site-header .locale-switcher-chevron {
  color: #d9c8ae;
}

html.dark .sst-public-tone-legal .locale-switcher-chevron {
  color: #e2c89e !important;
}
</style>
