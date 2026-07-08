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

      <div class="order-2 flex items-center gap-2 sm:order-3 sm:gap-2.5">
        <LocaleSwitcher variant="public" />
        <RouterLink
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="public-site-cta public-site-cta-mobile inline-flex sm:hidden"
        >
          {{ t('publicSite.enter') }}
        </RouterLink>
        <button
          type="button"
          @click="toggleTheme"
          class="public-site-tool public-site-tool-public"
          :title="isDark ? t('publicSite.theme.toLight') : t('publicSite.theme.toDark')"
        >
          <span class="public-site-tool-mark" aria-hidden="true">
            <Icon v-if="isDark" name="sun" size="sm" :stroke-width="1.55" />
            <Icon v-else name="courtyardMoon" size="sm" :stroke-width="1.55" />
          </span>
          <span class="public-site-tool-label">{{ isDark ? t('publicSite.theme.lightShort') : t('publicSite.theme.darkShort') }}</span>
        </button>
        <RouterLink
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="public-site-cta hidden sm:inline-flex"
        >
          {{ t('publicSite.enter') }}
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
import { useI18n } from 'vue-i18n'
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
const { t } = useI18n()
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
    { to: '/pricing', label: t('publicSite.nav.pricing') },
    { to: '/docs', label: t('publicSite.nav.docs') },
  ])

function toggleTheme() {
  isDark.value = toggleDocumentTheme(isDark.value)
}
</script>

<style scoped>
.public-site-header {
  --sst-public-tool-height: 2.28rem;
  --sst-public-tool-gap: 0.42rem;
  --sst-public-tool-border: rgba(142, 124, 95, 0.12);
  --sst-public-tool-hover-border: rgba(167, 58, 42, 0.14);
  --sst-public-tool-bg-top: rgba(255, 252, 247, 0.6);
  --sst-public-tool-bg-bottom: rgba(245, 238, 228, 0.74);
  --sst-public-tool-hover-top: rgba(255, 250, 244, 0.84);
  --sst-public-tool-hover-bottom: rgba(242, 233, 220, 0.9);
  --sst-public-tool-fg: #343831;
  --sst-public-tool-kicker: #9d7852;
  --sst-public-tool-label: #31352f;
  --sst-public-tool-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.54),
    0 6px 14px rgba(74, 56, 33, 0.045);
  --sst-public-tool-hover-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.58),
    0 8px 16px rgba(108, 74, 41, 0.06);
  --sst-public-tool-mark-size: 1.62rem;
  --sst-public-tool-mark-radius: 0.58rem;
  --sst-public-tool-mark-top: rgba(244, 238, 229, 0.94);
  --sst-public-tool-mark-bottom: rgba(234, 224, 208, 0.9);
  --sst-public-tool-mark-fg: #826446;
  --sst-public-tool-mark-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.56),
    inset 0 0 0 1px rgba(182, 155, 122, 0.11);
  --sst-public-tool-dot-size: 0.16rem;
  --sst-public-tool-dot: rgba(167, 58, 42, 0.72);
  --sst-public-tool-dot-shadow: 0 0 0 0.13rem rgba(167, 58, 42, 0.06);
}

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
  transition: color 180ms ease, background-color 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.public-site-tool:hover {
  color: #a73a2a;
}

.public-site-tool-public {
  min-height: calc(var(--sst-public-tool-height) - 0.04rem);
  gap: calc(var(--sst-public-tool-gap) - 0.05rem);
  border-radius: 999px;
  border: 1px solid rgba(154, 144, 129, 0.13);
  padding: 0.2rem 0.54rem 0.2rem 0.3rem;
  background:
    linear-gradient(180deg, rgba(252, 249, 244, 0.54), rgba(240, 235, 227, 0.7));
  color: #373a35;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.48),
    0 5px 12px rgba(70, 58, 41, 0.038);
}

.public-site-tool-public:hover {
  color: #3f413b;
  border-color: rgba(160, 149, 134, 0.17);
  background:
    linear-gradient(180deg, rgba(254, 251, 246, 0.7), rgba(243, 238, 231, 0.82));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.54),
    0 7px 14px rgba(70, 58, 41, 0.05);
}

.public-site-tool-public:focus-visible {
  outline: none;
  border-color: rgba(168, 156, 140, 0.2);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.52),
    0 0 0 1px rgba(244, 239, 231, 0.42),
    0 7px 14px rgba(70, 58, 41, 0.045);
}

.public-site-tool-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: calc(var(--sst-public-tool-mark-size) - 0.02rem);
  height: calc(var(--sst-public-tool-mark-size) - 0.02rem);
  border-radius: var(--sst-public-tool-mark-radius);
  background: linear-gradient(180deg, rgba(244, 239, 231, 0.88), rgba(234, 228, 219, 0.84));
  color: #8a6f4f;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.48),
    inset 0 0 0 1px rgba(184, 171, 151, 0.09);
}

.public-site-tool-mark::after {
  content: '';
  position: absolute;
  left: 0.34rem;
  top: 0.34rem;
  width: calc(var(--sst-public-tool-dot-size) - 0.01rem);
  height: calc(var(--sst-public-tool-dot-size) - 0.01rem);
  border-radius: 999px;
  background: rgba(179, 97, 78, 0.58);
  box-shadow: 0 0 0 0.12rem rgba(179, 97, 78, 0.045);
}

.public-site-tool-label {
  min-width: 1.72rem;
  font-size: 0.88rem;
  line-height: 1;
  letter-spacing: 0.02em;
  color: #353833;
  text-align: left;
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
  .public-site-header {
    --sst-public-tool-height: 2.14rem;
    --sst-public-tool-gap: 0.36rem;
    --sst-public-tool-mark-size: 1.48rem;
    --sst-public-tool-mark-radius: 0.52rem;
  }

  .public-site-tool-public {
    padding: 0.18rem 0.44rem 0.18rem 0.26rem;
  }

  .public-site-tool-label {
    min-width: 1.52rem;
    font-size: 0.8rem;
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
}

:global(html.dark) .public-site-header {
  --sst-public-tool-border: rgba(122, 108, 85, 0.28);
  --sst-public-tool-hover-border: rgba(155, 125, 88, 0.34);
  --sst-public-tool-bg-top: rgba(30, 31, 27, 0.82);
  --sst-public-tool-bg-bottom: rgba(23, 25, 21, 0.92);
  --sst-public-tool-hover-top: rgba(38, 36, 31, 0.88);
  --sst-public-tool-hover-bottom: rgba(29, 28, 24, 0.94);
  --sst-public-tool-fg: #dccfbc;
  --sst-public-tool-kicker: #b99567;
  --sst-public-tool-label: #ece0cd;
  --sst-public-tool-shadow:
    inset 0 1px 0 rgba(255, 243, 223, 0.035),
    0 8px 18px rgba(0, 0, 0, 0.18);
  --sst-public-tool-hover-shadow:
    inset 0 1px 0 rgba(255, 243, 223, 0.045),
    0 10px 20px rgba(0, 0, 0, 0.2);
  --sst-public-tool-mark-top: rgba(60, 52, 42, 0.78);
  --sst-public-tool-mark-bottom: rgba(42, 38, 31, 0.86);
  --sst-public-tool-mark-fg: #d1ac7d;
  --sst-public-tool-mark-shadow:
    inset 0 1px 0 rgba(255, 243, 223, 0.045),
    inset 0 0 0 1px rgba(201, 167, 122, 0.08);
  --sst-public-tool-dot: rgba(177, 95, 72, 0.7);
  --sst-public-tool-dot-shadow: 0 0 0 0.13rem rgba(167, 58, 42, 0.08);
}

:global(html.dark) .public-site-tool-public:hover {
  color: #f0e4d2;
}

:global(html.dark) .public-site-tool-public {
  border-color: rgba(112, 101, 81, 0.28);
  background:
    linear-gradient(180deg, rgba(31, 32, 28, 0.82), rgba(22, 24, 20, 0.92));
  color: #d9cbb8;
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.04),
    0 5px 12px rgba(0, 0, 0, 0.18);
}

:global(html.dark) .public-site-tool-public:hover {
  border-color: rgba(144, 118, 86, 0.34);
  background:
    linear-gradient(180deg, rgba(38, 36, 31, 0.88), rgba(29, 28, 24, 0.94));
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 7px 14px rgba(0, 0, 0, 0.22);
}

:global(html.dark) .public-site-tool-public:focus-visible {
  border-color: rgba(154, 128, 95, 0.38);
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 0 0 1px rgba(181, 144, 100, 0.12),
    0 7px 14px rgba(0, 0, 0, 0.22);
}

:global(html.dark) .public-site-tool-mark {
  background: linear-gradient(180deg, rgba(58, 50, 41, 0.8), rgba(42, 38, 31, 0.88));
  color: #d0ab7a;
  box-shadow:
    inset 0 1px 0 rgba(255, 251, 243, 0.06),
    inset 0 0 0 1px rgba(201, 167, 122, 0.08);
}

:global(html.dark) .public-site-tool-mark::after {
  background: rgba(177, 95, 72, 0.68);
  box-shadow: 0 0 0 0.12rem rgba(177, 95, 72, 0.06);
}

:global(html.dark) .public-site-tool-label {
  color: #e8dcc9;
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

html.dark .public-site-tool-public {
  border-color: rgba(112, 101, 81, 0.28);
  background: linear-gradient(180deg, rgba(31, 32, 28, 0.82), rgba(22, 24, 20, 0.92));
  color: #d9cbb8;
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.04),
    0 5px 12px rgba(0, 0, 0, 0.18);
}

html.dark .public-site-tool-public:hover {
  color: #f0e4d2;
  border-color: rgba(144, 118, 86, 0.34);
  background: linear-gradient(180deg, rgba(38, 36, 31, 0.88), rgba(29, 28, 24, 0.94));
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 7px 14px rgba(0, 0, 0, 0.22);
}

html.dark .public-site-tool-public:focus-visible {
  border-color: rgba(154, 128, 95, 0.38);
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 0 0 1px rgba(181, 144, 100, 0.12),
    0 7px 14px rgba(0, 0, 0, 0.22);
}

html.dark .public-site-tool-mark {
  background: linear-gradient(180deg, rgba(58, 50, 41, 0.8), rgba(42, 38, 31, 0.88));
  color: #d0ab7a;
  box-shadow:
    inset 0 1px 0 rgba(255, 251, 243, 0.06),
    inset 0 0 0 1px rgba(201, 167, 122, 0.08);
}

html.dark .public-site-tool-mark::after {
  background: rgba(177, 95, 72, 0.68);
  box-shadow: 0 0 0 0.12rem rgba(177, 95, 72, 0.06);
}

html.dark .public-site-tool-label {
  color: #e8dcc9;
}

html.dark .public-site-tool:not(.public-site-tool-public):hover {
  color: #f5ead8;
  background: rgba(255, 247, 235, 0.08);
}

html.dark .sst-public-tone-legal .public-site-tool:not(.public-site-tool-public):hover {
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

html.dark .sst-public .public-site-header nav {
  border-bottom-color: rgba(173, 145, 104, 0.22);
  box-shadow: inset 0 -1px 0 rgba(255, 243, 223, 0.02);
}

html.dark .sst-public .public-site-brand,
html.dark .sst-public.sst-public-tone-legal .public-site-brand {
  color: #f0e2cc;
  text-shadow: 0 1px 0 rgba(255, 242, 220, 0.05);
}

html.dark .sst-public .public-site-brand-sub,
html.dark .sst-public.sst-public-tone-legal .public-site-brand-sub {
  color: #aa9980;
}

html.dark .sst-public .public-site-nav,
html.dark .sst-public.sst-public-tone-legal .public-site-nav {
  color: #c9c0ac;
}

html.dark .sst-public .public-site-nav::after,
html.dark .sst-public.sst-public-tone-legal .public-site-nav::after {
  background: linear-gradient(90deg, transparent, rgba(188, 93, 31, 0.86), transparent);
}

html.dark .sst-public .public-site-nav:hover,
html.dark .sst-public.sst-public-tone-legal .public-site-nav:hover,
html.dark .sst-public .public-site-nav.is-active,
html.dark .sst-public.sst-public-tone-legal .public-site-nav.is-active {
  color: #f4efe4;
  background: rgba(255, 247, 235, 0.045);
}

html.dark .sst-public .public-site-tool,
html.dark .sst-public.sst-public-tone-legal .public-site-tool {
  color: #d8cdb9;
}

html.dark .sst-public .public-site-tool-public,
html.dark .sst-public.sst-public-tone-legal .public-site-tool-public {
  border-color: rgba(142, 118, 86, 0.36);
  background:
    linear-gradient(180deg, rgba(31, 32, 28, 0.84), rgba(22, 24, 20, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.1), transparent 26%);
  color: #e8dcc9;
  box-shadow:
    inset 0 1px 0 rgba(255, 247, 235, 0.05),
    0 8px 18px rgba(0, 0, 0, 0.2);
}

html.dark .sst-public .public-site-tool-public:hover,
html.dark .sst-public.sst-public-tone-legal .public-site-tool-public:hover {
  border-color: rgba(174, 136, 91, 0.44);
  background:
    linear-gradient(180deg, rgba(41, 37, 31, 0.9), rgba(29, 28, 24, 0.96)),
    radial-gradient(circle at 84% 14%, rgba(194, 126, 74, 0.13), transparent 26%);
  color: #f0e4d2;
}

html.dark .sst-public .public-site-cta,
html.dark .sst-public.sst-public-tone-legal .public-site-cta {
  border-color: rgba(142, 118, 86, 0.48);
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.94), rgba(14, 16, 13, 0.96)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.1), transparent 26%);
  color: #efe2cf;
  box-shadow:
    inset 0 1px 0 rgba(255, 238, 210, 0.06),
    0 14px 30px rgba(0, 0, 0, 0.2);
}

html.dark .sst-public .public-site-cta:hover,
html.dark .sst-public.sst-public-tone-legal .public-site-cta:hover,
html.dark .sst-public .public-site-cta-mobile:hover,
html.dark .sst-public.sst-public-tone-legal .public-site-cta-mobile:hover {
  border-color: rgba(194, 126, 74, 0.44);
  background: linear-gradient(180deg, rgba(51, 39, 30, 0.96), rgba(28, 24, 20, 0.98));
  color: #fff1de;
}

html:not(.dark) .sst-public .public-site-header nav {
  border-bottom-color: rgba(154, 128, 92, 0.18);
  box-shadow: inset 0 -1px 0 rgba(255, 249, 239, 0.52);
}

html:not(.dark) .sst-public .public-site-brand {
  color: #2f281f;
}

html:not(.dark) .sst-public .public-site-brand-sub {
  color: #8f7d63;
}

html:not(.dark) .sst-public .public-site-nav {
  color: #5f685c;
}

html:not(.dark) .sst-public .public-site-nav:hover,
html:not(.dark) .sst-public .public-site-nav.is-active {
  color: #2f281f;
  background: rgba(144, 113, 76, 0.055);
}

html:not(.dark) .sst-public .public-site-nav::after {
  background: linear-gradient(90deg, transparent, rgba(167, 58, 42, 0.7), transparent);
}

html:not(.dark) .sst-public .public-site-tool-public {
  border-color: rgba(154, 128, 92, 0.16);
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.68), rgba(244, 235, 220, 0.74)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.08), transparent 26%);
  color: #373a35;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.58),
    0 7px 16px rgba(84, 57, 31, 0.05);
}

html:not(.dark) .sst-public .public-site-tool-public:hover {
  border-color: rgba(167, 58, 42, 0.18);
  background:
    linear-gradient(180deg, rgba(255, 252, 247, 0.84), rgba(242, 233, 220, 0.9)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.11), transparent 26%);
  color: #2f281f;
}

html:not(.dark) .sst-public .public-site-cta {
  border-color: rgba(31, 35, 32, 0.1);
  background: linear-gradient(135deg, #28231c, #403227);
  color: #f7f0e4;
  box-shadow: 0 14px 26px rgba(84, 57, 31, 0.18), inset 0 1px 0 rgba(255, 241, 220, 0.12);
}

html:not(.dark) .sst-public .public-site-cta:hover {
  border-color: rgba(167, 58, 42, 0.24);
  background: linear-gradient(135deg, #3a2d23, #59432f);
  color: #fff8ee;
}

html:not(.dark) .sst-public .public-site-cta-mobile {
  border-color: rgba(154, 128, 92, 0.18);
  background: rgba(255, 252, 246, 0.74);
  color: #2f281f;
  box-shadow: 0 7px 16px rgba(84, 57, 31, 0.07);
}
</style>
