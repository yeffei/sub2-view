<template>
  <header class="app-header-shell glass sticky top-0 z-30 border-b border-zen-paperLine/70 dark:border-zen-nightLine">
    <div class="app-header-inner flex h-16 items-center justify-between gap-4 px-4 md:px-6">
      <div class="flex min-w-0 items-center gap-3 md:gap-4">
        <button
          @click="toggleMobileSidebar"
          class="btn-ghost btn-icon lg:hidden"
          aria-label="Toggle Menu"
        >
          <Icon name="menu" size="md" />
        </button>

        <router-link to="/home" class="app-header-brand" aria-label="返回首页">
          <span class="app-header-brand-seal" aria-hidden="true">
            <img src="/logo.png" alt="" class="h-full w-full object-contain" />
          </span>
          <span class="app-header-brand-copy">
            <small>{{ brandName }}</small>
            <strong>{{ pageTitle }}</strong>
          </span>
        </router-link>

        <div v-if="pageDescription" class="hidden min-w-0 xl:block">
          <p class="app-header-description truncate text-xs text-zen-mist dark:text-zen-stone">
            {{ pageDescription }}
          </p>
        </div>
      </div>

      <div class="app-header-actions flex items-center gap-3">
        <div v-if="!authStore.isAdmin" class="hidden items-center gap-1 lg:flex">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="app-header-utility-link"
          >
            {{ t('nav.docs') }}
          </a>
          <router-link
            v-else
            to="/docs"
            class="app-header-utility-link"
            :class="{ 'is-active': route.path === '/docs' }"
          >
            {{ t('nav.docs') }}
          </router-link>
        </div>

        <!-- Announcement Bell -->
        <AnnouncementBell v-if="user" />

        <!-- Language Switcher -->
        <LocaleSwitcher />

        <!-- Subscription Progress (for users with active subscriptions) -->
        <SubscriptionProgressMini v-if="user" />

        <!-- Balance Display -->
        <div
          v-if="user"
          class="app-header-balance hidden items-center gap-2 rounded-zen border border-zen-paperLine/80 bg-white/45 px-3 py-1.5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-night/60 sm:flex"
        >
          <Icon name="wallet" size="sm" class="app-header-balance-icon text-zen-seal dark:text-[#e9a092]" />
          <span class="app-header-balance-value font-mono text-sm font-semibold text-zen-ink dark:text-zen-paper">
            {{ formatHeaderMoney(availableBalance) }}
          </span>
          <span
            v-if="frozenBalance > 0"
            class="rounded-full bg-amber-100 px-1.5 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/40 dark:text-amber-200"
          >
            {{ balanceFrozenLabel }}
          </span>
          <div
            class="pointer-events-none absolute right-0 top-full mt-2 hidden w-56 rounded-lg border border-gray-200 bg-white p-3 text-xs shadow-lg group-hover:block dark:border-dark-700 dark:bg-dark-800"
          >
            <div class="flex items-center justify-between">
              <span class="text-gray-500 dark:text-dark-400">{{ balanceAvailableText }}</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ formatHeaderMoney(availableBalance) }}</span>
            </div>
            <div class="mt-2 flex items-center justify-between">
              <span class="text-gray-500 dark:text-dark-400">{{ balanceFrozenText }}</span>
              <span class="font-medium text-amber-700 dark:text-amber-200">{{ formatHeaderMoney(frozenBalance) }}</span>
            </div>
            <div class="mt-2 border-t border-gray-100 pt-2 dark:border-dark-700">
              <div class="flex items-center justify-between">
                <span class="text-gray-500 dark:text-dark-400">{{ balanceTotalText }}</span>
                <span class="font-semibold text-gray-900 dark:text-white">{{ formatHeaderMoney(totalBalance) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- User Dropdown -->
        <div v-if="user" class="relative" ref="dropdownRef">
          <button
            type="button"
            @click="toggleDropdown"
            class="app-header-user-trigger flex items-center gap-2 rounded-zen p-1.5 transition-colors hover:bg-zen-paperDeep/70 dark:hover:bg-zen-night"
            aria-label="User Menu"
            :aria-expanded="dropdownOpen"
            aria-haspopup="menu"
          >
            <div class="user-avatar flex h-8 w-8 items-center justify-center overflow-hidden rounded-zen text-sm font-medium shadow-seal">
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                :alt="displayName"
                class="h-full w-full object-cover"
              >
              <span v-else>{{ userInitials }}</span>
            </div>
            <div class="hidden text-left md:block">
              <div class="app-header-user-name text-sm font-medium text-zen-ink dark:text-zen-paper">
                {{ displayName }}
              </div>
              <div class="app-header-user-role text-xs capitalize text-zen-mist dark:text-zen-stone">
                {{ user.role }}
              </div>
            </div>
            <Icon name="chevronDown" size="sm" class="app-header-user-chevron hidden text-zen-stone md:block" />
          </button>

          <!-- Dropdown Menu -->
          <transition name="dropdown">
            <div v-if="dropdownOpen" class="dropdown right-0 mt-2 w-56" role="menu">
              <!-- User Info -->
              <div class="border-b border-zen-paperLine/70 px-4 py-3 dark:border-zen-nightLine">
                <div class="text-sm font-medium text-zen-ink dark:text-zen-paper">
                  {{ displayName }}
                </div>
                <div class="text-xs text-zen-mist dark:text-zen-stone">{{ user.email }}</div>
              </div>

              <!-- Balance (mobile only) -->
              <div class="border-b border-zen-paperLine/70 px-4 py-2 dark:border-zen-nightLine sm:hidden">
                <div class="text-xs text-zen-mist dark:text-zen-stone">
                  {{ t('common.balance') }}
                </div>
                <div class="font-mono text-sm font-semibold text-zen-seal dark:text-[#e9a092]">
                  {{ formatHeaderMoney(availableBalance) }}
                </div>
                <div v-if="frozenBalance > 0" class="mt-1 text-xs text-amber-600 dark:text-amber-300">
                  {{ balanceFrozenText }} {{ formatHeaderMoney(frozenBalance) }}
                </div>
              </div>

              <div class="py-1">
                <router-link to="/profile" @click="closeDropdown" class="dropdown-item">
                  <Icon name="user" size="sm" />
                  {{ t('nav.profile') }}
                </router-link>

                <router-link to="/keys" @click="closeDropdown" class="dropdown-item">
                  <Icon name="key" size="sm" />
                  {{ t('nav.apiKeys') }}
                </router-link>

                <a
                  v-if="authStore.isAdmin"
                  href="https://github.com/Wei-Shaw/sub2api"
                  target="_blank"
                  rel="noopener noreferrer"
                  @click="closeDropdown"
                  class="dropdown-item"
                >
                  <Icon name="github" size="sm" />
                  {{ t('nav.github') }}
                </a>

              </div>

              <!-- Contact Support (only show if configured) -->
              <div
                v-if="contactInfo"
                class="border-t border-zen-paperLine/70 px-4 py-2.5 dark:border-zen-nightLine"
              >
                <div class="flex items-center gap-2 text-xs text-zen-mist dark:text-zen-stone">
                  <Icon name="chat" size="xs" class="h-3.5 w-3.5 flex-shrink-0" />
                  <span>{{ t('common.contactSupport') }}:</span>
                  <span class="font-medium text-zen-inkSoft dark:text-zen-paper">{{
                    contactInfo
                  }}</span>
                </div>
              </div>

              <div v-if="showOnboardingButton" class="border-t border-zen-paperLine/70 py-1 dark:border-zen-nightLine">
                <button @click="handleReplayGuide" class="dropdown-item w-full">
                  <Icon name="questionCircle" size="sm" />
                  {{ $t('onboarding.restartTour') }}
                </button>
              </div>

              <div class="border-t border-zen-paperLine/70 py-1 dark:border-zen-nightLine">
                <button
                  @click="handleLogout"
                  class="dropdown-item w-full text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
                >
                  <Icon name="logout" size="sm" />
                  {{ t('nav.logout') }}
                </button>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore, useOnboardingStore } from '@/stores'
import { useAdminSettingsStore } from '@/stores/adminSettings'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import SubscriptionProgressMini from '@/components/common/SubscriptionProgressMini.vue'
import AnnouncementBell from '@/components/common/AnnouncementBell.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const adminSettingsStore = useAdminSettingsStore()
const onboardingStore = useOnboardingStore()

const user = computed(() => authStore.user)
const dropdownOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const contactInfo = computed(() => appStore.contactInfo)
const docUrl = computed(() => sanitizeUrl(appStore.docUrl))
const avatarUrl = computed(() => user.value?.avatar_url?.trim() || '')
const brandName = computed(() => {
  const configuredName = appStore.cachedPublicSettings?.site_name || appStore.siteName || ''
  return configuredName && configuredName !== 'Sub2API' ? configuredName : '山枢庭'
})
const availableBalance = computed(() => Number(user.value?.balance || 0))
const frozenBalance = computed(() => Number(user.value?.frozen_balance || 0))
const totalBalance = computed(() => availableBalance.value + frozenBalance.value)
const balanceAvailableText = computed(() => t('common.availableBalance') === 'common.availableBalance' ? '可用余额' : t('common.availableBalance'))
const balanceFrozenText = computed(() => t('common.frozenBalance') === 'common.frozenBalance' ? '冻结金额' : t('common.frozenBalance'))
const balanceTotalText = computed(() => t('common.totalBalance') === 'common.totalBalance' ? '总余额' : t('common.totalBalance'))
const balanceFrozenLabel = computed(() => `${balanceFrozenText.value} ${formatHeaderMoney(frozenBalance.value)}`)

// 只在标准模式的管理员下显示新手引导按钮
const showOnboardingButton = computed(() => {
  return false
})

const userInitials = computed(() => {
  if (!user.value) return ''
  // Prefer username, fallback to email
  if (user.value.username) {
    return user.value.username.substring(0, 2).toUpperCase()
  }
  if (user.value.email) {
    // Get the part before @ and take first 2 chars
    const localPart = user.value.email.split('@')[0]
    return localPart.substring(0, 2).toUpperCase()
  }
  return ''
})

const displayName = computed(() => {
  if (!user.value) return ''
  return user.value.username || user.value.email?.split('@')[0] || ''
})

const pageTitle = computed(() => {
  // For custom pages, use the menu item's label instead of generic "自定义页面"
  if (route.name === 'CustomPage') {
    const id = route.params.id as string
    const publicItems = appStore.cachedPublicSettings?.custom_menu_items ?? []
    const menuItem = publicItems.find((item) => item.id === id)
      ?? (authStore.isAdmin ? adminSettingsStore.customMenuItems.find((item) => item.id === id) : undefined)
    if (menuItem?.label) return menuItem.label
  }
  const titleKey = route.meta.titleKey as string
  if (titleKey) {
    return t(titleKey)
  }
  return (route.meta.title as string) || ''
})

const pageDescription = computed(() => {
  const descKey = route.meta.descriptionKey as string
  if (descKey) {
    return t(descKey)
  }
  return (route.meta.description as string) || ''
})

function toggleMobileSidebar() {
  appStore.toggleMobileSidebar()
}

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function closeDropdown() {
  dropdownOpen.value = false
}

async function handleLogout() {
  closeDropdown()
  try {
    await authStore.logout()
  } catch (error) {
    // Ignore logout errors - still redirect to login
    console.error('Logout error:', error)
  }
  await router.push('/login')
}

function handleReplayGuide() {
  closeDropdown()
  onboardingStore.replay()
}

function formatHeaderMoney(value: number) {
  if (!Number.isFinite(value)) return '$0.00'
  return `$${value.toFixed(2)}`
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    closeDropdown()
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
.app-header-shell {
  position: sticky;
  overflow: visible;
  isolation: isolate;
}

.app-header-shell::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 7% 24%, rgba(255, 255, 255, 0.4), transparent 16rem),
    linear-gradient(90deg, rgba(167, 58, 42, 0.035), transparent 26%, transparent 78%, rgba(155, 129, 85, 0.045));
  opacity: 0.9;
}

.app-header-shell::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 1px;
  pointer-events: none;
  background: linear-gradient(90deg, transparent, rgba(167, 58, 42, 0.18), transparent);
}

.app-header-inner {
  position: relative;
  z-index: 1;
}

.app-header-brand {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 0.85rem;
}

.app-header-brand-seal {
  display: inline-flex;
  height: 2.4rem;
  width: 2.4rem;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 0.85rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 12px 28px -18px rgba(31, 35, 32, 0.42);
}

.app-header-description {
  color: #7f776a;
}

.app-header-brand-copy {
  display: flex;
  min-width: 0;
  flex-direction: column;
  line-height: 1.15;
}

.app-header-brand-copy small {
  font-size: 0.68rem;
  letter-spacing: 0.22em;
  color: #7b6a53;
}

.app-header-brand-copy strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.98rem;
  font-weight: 600;
  color: #1f2320;
}

.app-header-utility-link {
  display: inline-flex;
  align-items: center;
  min-height: 2.25rem;
  border-radius: 0.9rem;
  border: 1px solid rgba(216, 205, 185, 0.78);
  padding: 0.32rem 0.78rem;
  background: rgba(255, 252, 245, 0.58);
  color: #4f5a51;
  font-size: 0.9rem;
  font-weight: 500;
  transition: color 180ms ease, background-color 180ms ease, border-color 180ms ease;
}

.app-header-utility-link:hover,
.app-header-utility-link.is-active {
  color: #1f2320;
  border-color: rgba(167, 58, 42, 0.24);
  background: rgba(237, 229, 212, 0.72);
}

.app-header-balance {
  min-height: 3rem;
  padding-inline: 0.95rem;
  border-radius: 0.9rem;
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.72), rgba(244, 239, 228, 0.62));
  box-shadow:
    0 16px 32px -28px rgba(31, 35, 32, 0.28),
    inset 0 1px 0 rgba(255, 255, 255, 0.64);
}

.app-header-balance-icon {
  opacity: 0.9;
}

.app-header-balance-value {
  letter-spacing: 0.01em;
}

.app-header-user-trigger {
  min-height: 3rem;
  padding-inline: 0.45rem;
}

.app-header-user-name {
  line-height: 1.1;
}

.app-header-user-role {
  margin-top: 0.1rem;
}

.app-header-user-chevron {
  opacity: 0.82;
}

.user-avatar {
  background: #a73a2a;
  color: #f4efe4;
  letter-spacing: 0.04em;
}

:deep(.dropdown) {
  border: 1px solid rgba(216, 205, 185, 0.86);
  border-radius: 6px;
  background: rgba(250, 247, 239, 0.96);
  box-shadow: 0 24px 70px -42px rgba(31, 35, 32, 0.36);
  backdrop-filter: blur(18px);
}

:deep(.dropdown-item) {
  color: #38413a;
}

:deep(.dropdown-item:hover) {
  background: rgba(237, 229, 212, 0.72);
  color: #1f2320;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}

:deep(.locale-switcher-btn-default) {
  min-height: 2.5rem;
  border: 1px solid rgba(216, 205, 185, 0.52);
  border-radius: 999px;
  padding-inline: 0.72rem;
  background: rgba(255, 252, 245, 0.42);
  color: #5d665d;
}

:deep(.locale-switcher-btn-default:hover) {
  background: rgba(255, 252, 245, 0.72);
  color: #1f2320;
}

:deep(.announcement-bell-trigger) {
  border: 1px solid transparent;
  border-radius: 999px;
}
</style>
<style>
.dark .app-header-shell {
  border-bottom-color: rgba(54, 59, 50, 0.94) !important;
  background:
    radial-gradient(circle at 12% 12%, rgba(255, 244, 223, 0.05), transparent 18rem),
    radial-gradient(circle at 86% 0%, rgba(167, 58, 42, 0.12), transparent 15rem),
    linear-gradient(180deg, rgba(18, 20, 16, 0.96), rgba(23, 26, 20, 0.92)) !important;
  box-shadow:
    inset 0 -1px 0 rgba(244, 239, 228, 0.02),
    0 18px 44px -36px rgba(0, 0, 0, 0.72) !important;
}

.dark .app-header-shell::before {
  background:
    radial-gradient(circle at 8% 32%, rgba(234, 220, 196, 0.05), transparent 16rem),
    linear-gradient(90deg, rgba(167, 58, 42, 0.07), transparent 24%, transparent 82%, rgba(214, 188, 141, 0.04));
}

.dark .app-header-shell::after {
  background: linear-gradient(90deg, transparent, rgba(167, 58, 42, 0.22), transparent);
}

.dark .app-header-brand-copy small {
  color: #93896f;
}

.dark .app-header-brand-copy strong {
  color: #f4efe4;
}

.dark .app-header-description {
  color: #8d8577;
}

.dark .app-header-utility-link {
  color: #c9c0ac;
  border-color: rgba(78, 84, 73, 0.82);
  background: rgba(24, 26, 21, 0.72);
}

.dark .app-header-utility-link:hover,
.dark .app-header-utility-link.is-active {
  color: #f4efe4;
  border-color: rgba(188, 93, 31, 0.34);
  background: rgba(48, 39, 31, 0.74);
}

.dark .app-header-balance {
  border-color: rgba(61, 66, 56, 0.9);
  background:
    linear-gradient(180deg, rgba(53, 57, 49, 0.92), rgba(39, 42, 35, 0.9));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.05),
    0 16px 36px -30px rgba(0, 0, 0, 0.6);
}

.dark .app-header-balance-icon {
  color: #d08e7f !important;
}

.dark .app-header-balance-value {
  color: #f4efe4;
}

.dark .app-header-user-trigger {
  border: 1px solid rgba(54, 59, 50, 0.72);
  background: rgba(23, 26, 20, 0.54);
}

.dark .app-header-user-trigger:hover {
  background: rgba(33, 37, 29, 0.88) !important;
  border-color: rgba(71, 77, 65, 0.86);
}

.dark .app-header-user-name {
  color: #f4efe4;
}

.dark .app-header-user-role,
.dark .app-header-user-chevron {
  color: #9d9585;
}

.dark .dropdown {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.96);
}

.dark .dropdown-item {
  color: #d7d0c2;
}

.dark .dropdown-item:hover {
  background: rgba(17, 19, 15, 0.88);
  color: #f4efe4;
}

.dark .announcement-bell-trigger {
  border-color: rgba(54, 59, 50, 0.7);
  background: rgba(23, 26, 20, 0.5);
  color: #b8b09f !important;
}

.dark .announcement-bell-trigger:hover {
  background: rgba(33, 37, 29, 0.82) !important;
  color: #f4efe4 !important;
}

.dark .app-header-actions .locale-switcher-btn-default {
  border-color: rgba(54, 59, 50, 0.74);
  background: rgba(23, 26, 20, 0.42);
  color: #b8b09f;
}

.dark .app-header-actions .locale-switcher-btn-default:hover {
  background: rgba(33, 37, 29, 0.82);
  color: #f4efe4;
}

.dark .app-header-actions .locale-switcher-chevron {
  color: #8d8577;
}
</style>
