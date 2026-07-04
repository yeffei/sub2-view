<template>
  <header class="public-site-header relative z-20 px-5 py-5 sm:px-8">
    <nav class="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-x-4 gap-y-3 border-b border-zen-paperLine/60 pb-4 dark:border-zen-nightLine/80">
      <div class="flex min-w-0 items-center gap-5 sm:gap-8">
        <RouterLink to="/home" class="flex items-center gap-3" aria-label="SST Home">
          <span class="public-site-logo-shell" aria-hidden="true">
            <img src="/logo.png" alt="" class="h-full w-full object-contain" />
          </span>
          <div class="leading-tight">
            <div class="public-site-brand font-serif text-base font-semibold tracking-[0.18em] text-zen-ink dark:text-zen-paper">
              {{ brandName }}
            </div>
            <div class="public-site-brand-sub text-[10px] uppercase tracking-[0.36em] text-zen-mist dark:text-zen-stone">
              SST
            </div>
          </div>
        </RouterLink>

        <div class="hidden min-w-0 items-center gap-1 sm:flex">
          <RouterLink
            v-for="item in resolvedNavItems"
            :key="item.to"
            :to="item.to"
            class="public-site-nav"
            :class="route.path === item.to ? 'is-active' : ''"
          >
            {{ item.label }}
          </RouterLink>
        </div>
      </div>

      <div class="order-2 flex items-center gap-1.5 sm:order-3 sm:gap-2">
        <LocaleSwitcher variant="public" />
        <RouterLink
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="public-site-cta public-site-cta-mobile inline-flex sm:hidden"
        >
          入庭
        </RouterLink>
        <button
          type="button"
          @click="toggleTheme"
          class="public-site-tool public-site-tool-icon"
          :title="isDark ? '切换到纸面' : '切换到夜庭'"
        >
          <Icon v-if="isDark" name="sun" size="md" />
          <Icon v-else name="moon" size="md" />
        </button>
        <RouterLink
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="public-site-cta hidden sm:inline-flex"
        >
          入庭
        </RouterLink>
      </div>

      <div class="order-3 flex w-full items-center gap-1 overflow-x-auto pb-1 sm:hidden">
        <RouterLink
          v-for="item in resolvedNavItems"
          :key="`${item.to}-mobile`"
          :to="item.to"
          class="public-site-nav"
          :class="route.path === item.to ? 'is-active' : ''"
        >
          {{ item.label }}
        </RouterLink>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { toggleTheme as toggleDocumentTheme, useThemeState } from '@/utils/theme'

interface NavItem {
  to: string
  label: string
}

const props = withDefaults(defineProps<{
  navItems?: NavItem[]
}>(), {
  navItems: () => [],
})

const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const isDark = useThemeState()
const configuredName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
const brandName = computed(() => configuredName.value && configuredName.value !== 'Sub2API' ? configuredName.value : '山枢庭')
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const resolvedNavItems = computed(() => props.navItems.length > 0
  ? props.navItems
  : [
    { to: '/pricing', label: '价目' },
    { to: '/docs', label: '文档' },
  ])

function toggleTheme() {
  isDark.value = toggleDocumentTheme(isDark.value)
}
</script>

<style scoped>
.public-site-logo-shell {
  display: inline-grid;
  width: 2.5rem;
  height: 2.5rem;
  place-items: center;
  overflow: hidden;
  border-radius: 0.85rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 10px 24px rgba(31, 35, 32, 0.16);
}

.public-site-nav {
  position: relative;
  display: inline-flex;
  align-items: center;
  min-height: 2.25rem;
  padding: 0.32rem 0.55rem 0.38rem;
  color: #5f685c;
  font-size: 0.95rem;
  letter-spacing: 0.01em;
  white-space: nowrap;
  transition: color 180ms ease, background-color 180ms ease;
}

.public-site-nav::after {
  content: '';
  position: absolute;
  left: 0.55rem;
  right: 0.55rem;
  bottom: 0.28rem;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(167, 58, 42, 0.7), transparent);
  opacity: 0;
  transition: opacity 180ms ease;
}

.public-site-nav:hover {
  color: #1f2320;
  background: rgba(216, 205, 185, 0.08);
}

.public-site-nav.is-active {
  color: #1f2320;
  background: transparent;
}

.public-site-nav:hover::after,
.public-site-nav.is-active::after {
  opacity: 1;
}

.public-site-tool {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 2.25rem;
  border-radius: 999px;
  color: #3b433d;
  transition: color 180ms ease, background-color 180ms ease;
}

.public-site-tool:hover {
  color: #a73a2a;
  background: rgba(216, 205, 185, 0.18);
}

.public-site-tool-text {
  padding: 0.32rem 0.72rem;
  font-size: 0.9rem;
}

.public-site-tool-icon {
  width: 2.25rem;
}

.public-site-cta {
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  border: 1px solid rgba(31, 35, 32, 0.1);
  background: #1f2320;
  padding: 0.78rem 1.35rem;
  color: #f4efe4;
  font-size: 0.95rem;
  font-weight: 500;
  letter-spacing: 0.01em;
  box-shadow: 0 10px 28px rgba(31, 35, 32, 0.14);
  transition: background-color 180ms ease, transform 180ms ease, box-shadow 180ms ease;
}

.public-site-cta:hover {
  background: #a73a2a;
  transform: translateY(-1px);
  box-shadow: 0 14px 32px rgba(167, 58, 42, 0.18);
}

.public-site-cta-mobile {
  border-radius: 999px;
  border-color: rgba(31, 35, 32, 0.08);
  background: rgba(244, 239, 228, 0.86);
  padding: 0.52rem 0.9rem;
  color: #2d332e;
  font-size: 0.82rem;
  letter-spacing: 0.08em;
  box-shadow: 0 6px 16px rgba(31, 35, 32, 0.08);
  backdrop-filter: blur(10px);
}

.public-site-cta-mobile:hover {
  background: rgba(232, 221, 201, 0.96);
  color: #1f2320;
  box-shadow: 0 10px 20px rgba(31, 35, 32, 0.12);
}

@media (max-width: 767px) {
  .public-site-tool-icon {
    display: none;
  }
}

:global(html.dark) .public-site-nav {
  color: #c9c0ac;
}

:global(html.dark) .public-site-nav:hover {
  color: #f4efe4;
  background: rgba(216, 205, 185, 0.04);
}

:global(html.dark) .public-site-nav.is-active {
  color: #f4efe4;
  background: transparent;
}

:global(html.dark) .public-site-tool {
  color: #d8cdb9;
}

:global(html.dark) .public-site-tool:hover {
  color: #ffd8bb;
  background: rgba(216, 205, 185, 0.08);
}

:global(html.dark) .public-site-cta {
  border-color: rgba(244, 239, 228, 0.12);
  background: #f4efe4;
  color: #171a16;
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.24);
}

:global(html.dark) .public-site-cta:hover {
  background: #b95d1f;
  color: #fff8f0;
}

:global(html.dark) .public-site-cta-mobile {
  border-color: rgba(244, 239, 228, 0.12);
  background: rgba(27, 31, 27, 0.72);
  color: #ece1cf;
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.04), 0 8px 18px rgba(0, 0, 0, 0.18);
}

:global(html.dark) .public-site-cta-mobile:hover {
  background: rgba(57, 42, 30, 0.88);
  color: #fff1de;
}
</style>

<style>
html.dark .public-site-nav {
  color: #c5bcaa;
}

html.dark .sst-public-tone-legal .public-site-nav {
  color: #ccb892;
}

html.dark .public-site-nav::after {
  background: linear-gradient(90deg, transparent, rgba(188, 93, 31, 0.9), transparent);
}

html.dark .sst-public-tone-legal .public-site-nav::after {
  background: linear-gradient(90deg, transparent, rgba(203, 147, 76, 0.96), transparent);
}

html.dark .public-site-nav:hover {
  color: #f2e7d6;
  background: rgba(255, 247, 235, 0.06);
}

html.dark .sst-public-tone-legal .public-site-nav:hover {
  color: #f5e7cf;
  background: rgba(205, 163, 103, 0.08);
}

html.dark .public-site-nav.is-active {
  color: #fbf2e4;
  background: transparent;
}

html.dark .sst-public-tone-legal .public-site-nav.is-active {
  color: #f1c27c;
}

html.dark .sst-public-tone-legal .public-site-brand {
  color: #f6ead7;
  text-shadow: 0 1px 0 rgba(255, 239, 210, 0.05);
}

html.dark .sst-public-tone-legal .public-site-brand-sub {
  color: #bca37d;
}

html.dark .public-site-tool {
  color: #bfb5a2;
}

html.dark .sst-public-tone-legal .public-site-tool {
  color: #cdb892;
}

html.dark .public-site-tool:hover {
  color: #f5ead8;
  background: rgba(255, 247, 235, 0.08);
}

html.dark .sst-public-tone-legal .public-site-tool:hover {
  color: #f1c27c;
  background: rgba(205, 163, 103, 0.1);
}

html.dark .public-site-cta {
  border-color: rgba(246, 236, 221, 0.14);
  background: rgba(31, 34, 28, 0.88);
  color: #f7eddf;
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.05), 0 12px 28px rgba(0, 0, 0, 0.18);
}

html.dark .sst-public-tone-legal .public-site-cta {
  border-color: rgba(188, 143, 84, 0.16);
  background: linear-gradient(180deg, rgba(31, 28, 24, 0.94), rgba(22, 20, 18, 0.96));
  color: #f7ebd8;
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.06), 0 12px 28px rgba(0, 0, 0, 0.22);
}

html.dark .public-site-cta:hover {
  background: rgba(59, 42, 31, 0.94);
  color: #fff6e9;
}

html.dark .public-site-cta-mobile {
  border-color: rgba(246, 236, 221, 0.12);
  background: rgba(24, 28, 24, 0.72);
  color: #f0e1ca;
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.04), 0 8px 18px rgba(0, 0, 0, 0.18);
}

html.dark .public-site-cta-mobile:hover {
  background: rgba(60, 43, 28, 0.9);
  color: #fff1de;
}

html.dark .sst-public-tone-legal .public-site-cta:hover {
  background: linear-gradient(180deg, rgba(67, 47, 29, 0.96), rgba(39, 31, 24, 0.98));
  color: #fff1de;
}

html.dark .sst-public-tone-legal .public-site-cta-mobile {
  border-color: rgba(188, 143, 84, 0.18);
  background: rgba(34, 28, 23, 0.76);
  color: #f5e5cc;
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.05), 0 8px 18px rgba(0, 0, 0, 0.2);
}

html.dark .sst-public-tone-legal .public-site-cta-mobile:hover {
  background: rgba(70, 50, 31, 0.92);
  color: #fff1de;
}

html.dark .sst-home .public-site-header nav {
  border-bottom-color: rgba(173, 145, 104, 0.22);
  box-shadow: inset 0 -1px 0 rgba(255, 243, 223, 0.02);
}

html.dark .sst-home .public-site-brand {
  color: #f0e2cc;
  text-shadow: 0 1px 0 rgba(255, 242, 220, 0.05);
}
</style>
