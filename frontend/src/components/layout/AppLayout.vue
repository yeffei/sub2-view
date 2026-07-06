<template>
  <div class="app-shell min-h-screen bg-zen-paper/80 dark:bg-zen-night">
    <div class="app-paper pointer-events-none fixed inset-0"></div>

    <!-- Sidebar -->
    <AppSidebar v-if="!immersiveUserExperience" />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="[immersiveUserExperience ? 'lg:ml-0' : (sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64')]"
    >
      <!-- Header -->
      <AppHeader v-if="!immersiveUserExperience" />

      <!-- Main Content -->
      <main :class="immersiveUserExperience ? 'p-0' : 'p-4 md:p-6 lg:p-8'">
        <section v-if="userWorkbenchShell" class="sst-user-workbench">
          <div class="sst-user-frame">
            <header class="sst-user-head">
              <router-link to="/home" class="sst-brand-lockup" aria-label="返回首页">
                <span class="sst-seal" aria-hidden="true"><img src="/logo.png" alt="" /></span>
                <span class="sst-brand-copy">
                  <small>山枢庭</small>
                  <strong>{{ currentPageTitle }}</strong>
                </span>
              </router-link>

              <div
                ref="accountMenuRef"
                class="sst-account-menu"
                :class="{ 'is-open': isAccountMenuOpen }"
              >
                <div>
                  <span>{{ currentDateLabel }}</span>
                  <strong>{{ authStore.user?.email || '山枢庭账户' }}</strong>
                  <small>统一入口，安静流转。</small>
                </div>
                <button
                  type="button"
                  class="sst-account-trigger"
                  aria-label="账户菜单"
                  aria-haspopup="menu"
                  :aria-expanded="isAccountMenuOpen ? 'true' : 'false'"
                  @click="toggleAccountMenu"
                >
                  身
                </button>
                <div class="sst-account-dropdown">
                  <router-link to="/profile" @click="closeAccountMenu">身份文书</router-link>
                  <div class="sst-account-theme" aria-label="外观设置">
                    <span>外观设置</span>
                    <div class="sst-account-theme-options">
                      <button
                        v-for="option in themeOptions"
                        :key="option.value"
                        type="button"
                        :class="{ 'is-selected': themePreference === option.value }"
                        @click="chooseTheme(option.value)"
                      >
                        {{ option.label }}
                      </button>
                    </div>
                  </div>
                  <router-link to="/dashboard" @click="closeAccountMenu">返回庭院</router-link>
                  <button type="button" @click="handleLogout">退出登录</button>
                </div>
              </div>
            </header>

            <nav class="sst-user-nav" aria-label="用户功能导航">
              <template v-for="item in userNavItems" :key="item.path">
                <router-link
                  v-if="!isNavActive(item.path)"
                  :to="item.path"
                  class="sst-user-nav-link"
                >
                  <span>{{ item.mark }}</span>
                  <strong>{{ item.label }}</strong>
                </router-link>
                <span
                  v-else
                  class="sst-user-nav-link is-active"
                  aria-current="page"
                >
                  <span>{{ item.mark }}</span>
                  <strong>{{ item.label }}</strong>
                </span>
              </template>
            </nav>

            <section class="sst-user-content">
              <slot />
            </section>
          </div>
        </section>

        <section v-else-if="adminWorkbenchShell" class="sst-admin-workbench">
          <div class="sst-admin-stage">
            <slot />
          </div>
        </section>

        <slot v-else />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { setThemePreference, themePreferenceLabels, useThemePreference, type ThemePreference } from '@/utils/theme'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const isAccountMenuOpen = ref(false)
const accountMenuRef = ref<HTMLElement | null>(null)
const themePreference = useThemePreference()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')
const immersiveUserRoutePrefixes = [
  '/dashboard',
  '/keys',
  '/usage',
  '/monitor',
  '/profile',
  '/purchase',
  '/orders',
  '/redeem',
  '/affiliate',
  '/payment/qrcode',
  '/payment/result',
  '/payment/stripe',
  '/payment/airwallex',
  '/custom/',
]
const isImmersiveUserRoute = computed(() => immersiveUserRoutePrefixes.some(path => route.path === path || route.path.startsWith(path)))
const immersiveUserExperience = computed(() => isImmersiveUserRoute.value && !route.path.startsWith('/admin'))
const userWorkbenchShell = computed(() => immersiveUserExperience.value && route.path !== '/dashboard')
const adminWorkbenchShell = computed(() => route.path.startsWith('/admin'))
const currentDateLabel = computed(() => new Intl.DateTimeFormat('zh-CN', {
  timeZone: 'Asia/Shanghai',
  year: 'numeric',
  month: '2-digit',
  day: '2-digit',
}).format(new Date()).replace(/\//g, '-'))

interface UserShellNavItem {
  path: string
  label: string
  mark: string
  hideInSimpleMode?: boolean
  visible?: boolean
}

const userNavItems = computed(() => {
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const items: UserShellNavItem[] = [
    { path: '/dashboard', label: '今日概览', mark: '庭' },
    { path: '/keys', label: 'API 密钥', mark: '钥' },
    { path: '/usage', label: '用量账册', mark: '账', hideInSimpleMode: true },
    { path: '/monitor', label: '服务状态', mark: '候', visible: isFeatureFlagEnabled(FeatureFlags.channelMonitor) },
    { path: '/purchase', label: '充值与兑换', mark: '财', hideInSimpleMode: true, visible: paymentEnabled || !authStore.isSimpleMode },
    { path: '/affiliate', label: '团队引荐', mark: '荐', hideInSimpleMode: true, visible: isFeatureFlagEnabled(FeatureFlags.affiliate) },
    { path: '/profile', label: '账户安全', mark: '身' },
  ]

  const customItems = (appStore.cachedPublicSettings?.custom_menu_items ?? [])
    .filter(item => item.visibility === 'user')
    .sort((a, b) => a.sort_order - b.sort_order)
    .map((item): UserShellNavItem => ({ path: `/custom/${item.id}`, label: item.label, mark: '册' }))

  return [...items, ...customItems]
    .filter(item => item.visible !== false)
    .filter(item => !authStore.isSimpleMode || !item.hideInSimpleMode)
})

const billingRoutePaths = ['/purchase', '/orders', '/redeem']
const themeOptions: Array<{ value: ThemePreference, label: string }> = [
  { value: 'system', label: themePreferenceLabels.system },
  { value: 'light', label: themePreferenceLabels.light },
  { value: 'dark', label: themePreferenceLabels.dark },
]

const currentPageTitle = computed(() => {
  if (route.path === '/profile') {
    return '账户安全'
  }
  const matched = userNavItems.value.find(item => isNavActive(item.path))
  return matched?.label || '用户后台'
})

function isNavActive(path: string) {
  if ((path === '/purchase' || path === '/redeem') && billingRoutePaths.includes(route.path)) {
    return true
  }
  return route.path === path || (path.startsWith('/custom/') && route.path.startsWith(path))
}

function chooseTheme(preference: ThemePreference) {
  setThemePreference(preference)
}


function closeAccountMenu() {
  isAccountMenuOpen.value = false
}

function openAccountMenu() {
  isAccountMenuOpen.value = true
}

function toggleAccountMenu() {
  if (isAccountMenuOpen.value) closeAccountMenu()
  else openAccountMenu()
}


function handleAccountPointerDown(event: PointerEvent) {
  if (!isAccountMenuOpen.value) return
  const target = event.target as Node | null
  if (target && accountMenuRef.value?.contains(target)) return
  closeAccountMenu()
}

function handleAccountKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') closeAccountMenu()
}

async function handleLogout() {
  closeAccountMenu()
  try {
    await authStore.logout()
  } catch (error) {
    console.error('Logout error:', error)
  }
  await router.push('/login')
}

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: false
})

const onboardingStore = useOnboardingStore()

watch(() => route.fullPath, closeAccountMenu)

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
  document.addEventListener('pointerdown', handleAccountPointerDown)
  document.addEventListener('keydown', handleAccountKeydown)
})

defineExpose({ replayTour })

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleAccountPointerDown)
  document.removeEventListener('keydown', handleAccountKeydown)
})
</script>

<style scoped>
.app-shell {
  position: relative;
}

.app-paper {
  background:
    radial-gradient(circle at 9% 4%, rgba(167, 58, 42, 0.045), transparent 20rem),
    radial-gradient(circle at 82% 8%, rgba(155, 129, 85, 0.07), transparent 22rem),
    linear-gradient(180deg, rgba(250, 247, 239, 0.88), rgba(237, 229, 212, 0.84)),
    linear-gradient(rgba(31, 35, 32, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.014) 1px, transparent 1px);
  background-size: auto, 112px 112px, 112px 112px;
}

.sst-admin-workbench {
  --sst-admin-line: rgba(198, 184, 157, 0.54);
  --sst-admin-line-soft: rgba(216, 205, 185, 0.38);
  --sst-admin-paper: rgba(250, 247, 239, 0.86);
  --sst-admin-paper-strong: rgba(255, 252, 245, 0.94);
  --sst-admin-paper-muted: rgba(244, 239, 228, 0.62);
  --sst-admin-ink: #1f2320;
  --sst-admin-ink-soft: #445046;
  --sst-admin-muted: #6f7a70;
  --sst-admin-seal: #a73a2a;
  position: relative;
  min-height: calc(100vh - 7.5rem);
}

.sst-admin-stage {
  position: relative;
  max-width: 1480px;
  min-width: 0;
  margin: 0 auto;
  padding: clamp(0.6rem, 1.2vw, 1rem);
  overflow-x: hidden;
  overflow-y: visible;
  border: 1px solid var(--sst-admin-line);
  border-radius: 14px;
  background:
    radial-gradient(circle at 92% 6%, rgba(167, 58, 42, 0.045), transparent 16rem),
    radial-gradient(circle at 8% 2%, rgba(255, 252, 245, 0.88), transparent 18rem),
    linear-gradient(180deg, rgba(255, 252, 245, 0.9), rgba(246, 241, 231, 0.76)),
    var(--sst-admin-paper);
  box-shadow:
    0 26px 70px -54px rgba(31, 35, 32, 0.32),
    inset 0 1px 0 rgba(255, 255, 255, 0.68);
}

.sst-admin-stage::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.035), transparent 14%, transparent 84%, rgba(155, 129, 85, 0.05)),
    linear-gradient(rgba(31, 35, 32, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.015) 1px, transparent 1px);
  background-size: auto, 120px 120px, 120px 120px;
  pointer-events: none;
}

.sst-admin-stage > :deep(*) {
  position: relative;
  z-index: 1;
  min-width: 0;
  max-width: 100%;
}

.sst-admin-stage :deep(*) {
  letter-spacing: 0;
}

.sst-admin-stage :deep(h1),
.sst-admin-stage :deep(h2),
.sst-admin-stage :deep(h3) {
  color: var(--sst-admin-ink);
  line-height: 1.28;
}

.sst-admin-stage :deep(h1) {
  font-size: clamp(1.42rem, 1.6vw, 1.72rem);
  font-weight: 680;
}

.sst-admin-stage :deep(h2) {
  font-size: clamp(1.18rem, 1.25vw, 1.38rem);
  font-weight: 660;
}

.sst-admin-stage :deep(h3),
.sst-admin-stage :deep(.card-title),
.sst-admin-stage :deep(.modal-title) {
  font-size: 1rem;
  font-weight: 650;
}

.sst-admin-stage :deep(p),
.sst-admin-stage :deep(label),
.sst-admin-stage :deep(.text-sm) {
  line-height: 1.62;
}

.sst-admin-stage :deep(.text-xs) {
  line-height: 1.48;
}

.sst-admin-stage :deep(.btn),
.sst-admin-stage :deep(button),
.sst-admin-stage :deep(input),
.sst-admin-stage :deep(textarea),
.sst-admin-stage :deep(select) {
  font-size: 0.875rem;
}

.sst-admin-stage :deep(.card),
.sst-admin-stage :deep(.bg-white),
.sst-admin-stage :deep(.bg-white\/70),
.sst-admin-stage :deep(.bg-white\/80),
.sst-admin-stage :deep(.bg-gray-50),
.sst-admin-stage :deep(.bg-gray-100),
.sst-admin-stage :deep(.bg-stone-50),
.sst-admin-stage :deep(.bg-stone-100),
.sst-admin-stage :deep(.bg-stone-100\/35),
.sst-admin-stage :deep(.dark\:bg-dark-800),
.sst-admin-stage :deep(.dark\:bg-dark-700),
.sst-admin-stage :deep(.dark\:bg-dark-900\/40),
.sst-admin-stage :deep(.rounded-2xl),
.sst-admin-stage :deep(.rounded-xl),
.sst-admin-stage :deep(.sst-admin-panel) {
  border-color: var(--sst-admin-line);
  background-color: var(--sst-admin-paper);
  box-shadow:
    0 18px 48px -42px rgba(31, 35, 32, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.58);
}

.sst-admin-stage :deep(.input),
.sst-admin-stage :deep(input),
.sst-admin-stage :deep(textarea),
.sst-admin-stage :deep(select),
.sst-admin-stage :deep(.select-trigger) {
  border-color: rgba(198, 184, 157, 0.56);
  background-color: var(--sst-admin-paper-strong);
  color: var(--sst-admin-ink);
}

.sst-admin-stage :deep(.input:focus),
.sst-admin-stage :deep(input:focus),
.sst-admin-stage :deep(textarea:focus),
.sst-admin-stage :deep(select:focus),
.sst-admin-stage :deep(.select-trigger:focus-within) {
  border-color: rgba(167, 58, 42, 0.38);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.1);
  outline: none;
}

.sst-admin-stage :deep(table) {
  width: 100%;
  font-variant-numeric: tabular-nums;
  border-collapse: separate;
  border-spacing: 0;
  color: #38413a;
}

.sst-admin-stage :deep(thead) {
  background: rgba(237, 229, 212, 0.72);
}

.sst-admin-stage :deep(table th) {
  height: 2.75rem !important;
  border-bottom: 1px solid rgba(198, 184, 157, 0.54) !important;
  background:
    linear-gradient(180deg, rgba(237, 229, 212, 0.84), rgba(231, 220, 201, 0.72)) !important;
  color: #6d634f !important;
  font-size: 0.76rem !important;
  font-weight: 680 !important;
  line-height: 1.35 !important;
  white-space: nowrap;
}

.sst-admin-stage :deep(table td) {
  min-height: 3rem !important;
  border-bottom: 1px solid rgba(198, 184, 157, 0.28) !important;
  color: #38413a !important;
  font-size: 0.84rem !important;
  line-height: 1.52 !important;
  vertical-align: middle;
}

.sst-admin-stage :deep(table tbody tr:last-child td) {
  border-bottom-color: transparent !important;
}

.sst-admin-stage :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.036);
}

.sst-admin-stage :deep(.table-container),
.sst-admin-stage :deep(.table-wrapper),
.sst-admin-stage :deep(.overflow-x-auto:has(table)) {
  border-color: rgba(198, 184, 157, 0.48);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.62);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.52);
}

.sst-admin-stage :deep(.btn-primary) {
  border-color: rgba(167, 58, 42, 0.34);
  background: #a73a2a;
  color: #faf7ef;
}

.sst-admin-stage :deep(.btn-primary:hover) {
  background: #8f2f23;
}

.sst-admin-stage :deep(.btn-secondary) {
  border-color: rgba(198, 184, 157, 0.62);
  background: var(--sst-admin-paper-strong);
  color: var(--sst-admin-ink-soft);
}

.sst-admin-stage :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.26);
  background: rgba(250, 247, 239, 0.98);
  color: var(--sst-admin-seal);
}

.sst-admin-stage :deep(.text-primary-600),
.sst-admin-stage :deep(.text-primary-700),
.sst-admin-stage :deep(.text-blue-600),
.sst-admin-stage :deep(.text-indigo-600),
.sst-admin-stage :deep(.text-purple-600),
.sst-admin-stage :deep(.text-violet-600),
.sst-admin-stage :deep(.text-rose-600) {
  color: var(--sst-admin-seal);
}

.sst-admin-stage :deep(.bg-primary-50),
.sst-admin-stage :deep(.bg-blue-100),
.sst-admin-stage :deep(.bg-indigo-100),
.sst-admin-stage :deep(.bg-purple-100),
.sst-admin-stage :deep(.bg-violet-100),
.sst-admin-stage :deep(.bg-rose-100) {
  background-color: rgba(167, 58, 42, 0.1);
}

.sst-admin-stage :deep(.bg-green-100),
.sst-admin-stage :deep(.bg-emerald-100),
.sst-admin-stage :deep(.bg-teal-100) {
  background-color: rgba(81, 98, 79, 0.1);
}

.sst-admin-stage :deep(.bg-amber-100),
.sst-admin-stage :deep(.bg-orange-100),
.sst-admin-stage :deep(.bg-yellow-100) {
  background-color: rgba(155, 129, 85, 0.13);
}

.sst-admin-stage :deep(.sst-admin-page) {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: clamp(0.2rem, 0.6vw, 0.4rem);
}

.sst-admin-stage :deep(.sst-admin-panel) {
  border: 1px solid rgba(198, 184, 157, 0.48);
  border-radius: 12px;
  background: rgba(250, 247, 239, 0.72);
  box-shadow: 0 18px 48px -42px rgba(31, 35, 32, 0.2);
}

.sst-admin-stage :deep(.sst-admin-soft-button) {
  color: #38413a;
  background: rgba(255, 252, 245, 0.82);
}

.sst-admin-stage :deep(.modal-content),
.sst-admin-stage :deep(.dialog-container),
.sst-admin-stage :deep(.dropdown) {
  border-color: rgba(198, 184, 157, 0.66);
  background: rgba(255, 252, 245, 0.98);
  color: var(--sst-admin-ink);
  box-shadow: 0 26px 72px -48px rgba(31, 35, 32, 0.34);
}

.sst-admin-stage :deep(.modal-header),
.sst-admin-stage :deep(.modal-footer),
.sst-admin-stage :deep(.dialog-header),
.sst-admin-stage :deep(.dialog-footer),
.sst-admin-stage :deep(.card-header),
.sst-admin-stage :deep(.card-footer) {
  border-color: var(--sst-admin-line-soft);
  background-color: rgba(250, 247, 239, 0.62);
}

.sst-admin-stage :deep(.badge-gray),
.sst-admin-stage :deep(.bg-gray-200),
.sst-admin-stage :deep(.bg-slate-100) {
  background-color: rgba(216, 205, 185, 0.42);
  color: #445046;
}

.sst-admin-stage :deep(.text-gray-900),
.sst-admin-stage :deep(.text-gray-800),
.sst-admin-stage :deep(.text-gray-700) {
  color: var(--sst-admin-ink);
}

.sst-admin-stage :deep(.text-gray-600),
.sst-admin-stage :deep(.text-gray-500),
.sst-admin-stage :deep(.text-gray-400) {
  color: var(--sst-admin-muted);
}

.sst-user-workbench {
  --sst-paper: #faf7ef;
  --sst-paper-soft: rgba(250, 247, 239, 0.76);
  --sst-ink: #1f2320;
  --sst-ink-soft: #38413a;
  --sst-mute: #667066;
  --sst-line: rgba(198, 184, 157, 0.62);
  --sst-line-soft: rgba(198, 184, 157, 0.32);
  --sst-seal: #a73a2a;
  min-height: 100vh;
  padding: clamp(0.75rem, 1.8vw, 1.6rem);
  background:
    linear-gradient(90deg, rgba(119, 104, 78, 0.08) 1px, transparent 1px) 0 0 / 4rem 4rem,
    radial-gradient(circle at 14% 12%, rgba(167, 58, 42, 0.055), transparent 22rem),
    linear-gradient(180deg, rgba(244, 239, 228, 0.38), rgba(237, 229, 212, 0.56)),
    #f4efe4;
  color: var(--sst-ink);
}

.sst-user-frame {
  position: relative;
  max-width: 1360px;
  min-height: calc(100vh - clamp(1.5rem, 3.6vw, 3.2rem));
  margin: 0 auto;
  overflow: hidden;
  border: 1px solid rgba(198, 184, 157, 0.82);
  border-radius: 14px;
  background:
    radial-gradient(circle at 86% 8%, rgba(167, 58, 42, 0.045), transparent 18rem),
    linear-gradient(135deg, rgba(250, 247, 239, 0.96), rgba(244, 239, 228, 0.88)),
    url('@/assets/brand/sst-paper-ink-bg.png') center/cover;
  box-shadow: 0 26px 70px -52px rgba(31, 35, 32, 0.44);
}

.sst-user-frame::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.07), transparent 16%, transparent 84%, rgba(155, 129, 85, 0.08)),
    linear-gradient(180deg, transparent, rgba(244, 239, 228, 0.18));
  pointer-events: none;
}

.sst-user-head,
.sst-user-nav,
.sst-user-content {
  position: relative;
}

.sst-user-head {
  z-index: 30;
}

.sst-user-nav {
  z-index: 10;
}

.sst-user-content {
  z-index: 1;
}

.sst-user-head {
  display: flex;
  min-height: 5.25rem;
  align-items: center;
  justify-content: space-between;
  gap: 1.2rem;
  border-bottom: 1px solid var(--sst-line-soft);
  padding: clamp(0.85rem, 1.65vw, 1.25rem) clamp(1rem, 2.2vw, 1.75rem);
}

.sst-brand-lockup {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 1.05rem;
  color: inherit;
}

.sst-seal {
  display: inline-grid;
  width: 2.55rem;
  height: 2.55rem;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 0.42rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 14px 28px -24px rgba(31, 35, 32, 0.45);
}

.sst-seal img {
  display: block;
  width: 100%;
  height: 100%;
}

.sst-brand-copy {
  display: grid;
  min-width: 0;
  gap: 0.24rem;
}

.sst-brand-copy small,
.sst-account-menu span,
.sst-user-nav-link span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.sst-brand-copy strong {
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.28rem, 1.68vw, 1.68rem);
  font-weight: 600;
  line-height: 1.15;
}

.sst-account-menu {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.9rem;
  min-width: min(22rem, 42vw);
  border-right: 2px solid var(--sst-seal);
  padding-right: 0.9rem;
  text-align: right;
}

.sst-account-menu strong,
.sst-account-menu small {
  display: block;
}

.sst-account-menu strong {
  margin-top: 0.34rem;
  overflow-wrap: anywhere;
  color: var(--sst-ink-soft);
  font-size: 0.86rem;
}

.sst-account-menu small {
  margin-top: 0.42rem;
  color: var(--sst-mute);
  font-size: 0.74rem;
}

.sst-account-trigger {
  display: inline-grid;
  width: 2.28rem;
  height: 2.28rem;
  flex: 0 0 auto;
  border: 1px solid rgba(167, 58, 42, 0.24);
  border-radius: 0.42rem;
  place-items: center;
  background:
    linear-gradient(135deg, rgba(255, 244, 226, 0.12), transparent 48%),
    rgba(167, 58, 42, 0.1);
  color: var(--sst-seal);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 700;
  transition: background-color 180ms ease, border-color 180ms ease, transform 180ms ease;
}

.sst-account-trigger:hover,
.sst-account-menu.is-open .sst-account-trigger {
  border-color: rgba(167, 58, 42, 0.42);
  background-color: rgba(167, 58, 42, 0.14);
  transform: translateY(-1px);
}

.sst-account-dropdown {
  position: absolute;
  top: calc(100% - 0.1rem);
  right: 0.9rem;
  z-index: 10000;
  display: grid;
  min-width: 9.5rem;
  border: 1px solid rgba(198, 184, 157, 0.68);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.96);
  box-shadow: 0 20px 46px -32px rgba(31, 35, 32, 0.4);
  opacity: 0;
  visibility: hidden;
  padding: 0.35rem;
  pointer-events: none;
  transform: translateY(-0.35rem);
  transition: opacity 160ms ease, transform 160ms ease, visibility 0s linear 160ms;
}


.sst-account-menu.is-open .sst-account-dropdown {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: translateY(0);
  transition-delay: 0s;
}

.sst-account-dropdown a,
.sst-account-dropdown button {
  border-radius: 7px;
  padding: 0.62rem 0.72rem;
  color: var(--sst-ink-soft);
  font-size: 0.82rem;
  font-weight: 650;
  text-align: left;
  position: relative;
  z-index: 1;
  pointer-events: auto;
  transition: background-color 160ms ease, color 160ms ease;
}

.sst-account-theme {
  display: grid;
  gap: 0.45rem;
  border-top: 1px solid rgba(198, 184, 157, 0.44);
  border-bottom: 1px solid rgba(198, 184, 157, 0.44);
  margin: 0.2rem 0;
  padding: 0.62rem 0.5rem;
}

.sst-account-theme > span {
  padding: 0 0.22rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.sst-account-theme-options {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.25rem;
}

.sst-account-dropdown .sst-account-theme-options button {
  justify-content: center;
  padding: 0.42rem 0.36rem;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 7px;
  background: rgba(255, 252, 245, 0.54);
  color: #59645a;
  font-size: 0.74rem;
  font-weight: 650;
  text-align: center;
}

.sst-account-dropdown .sst-account-theme-options button:hover,
.sst-account-dropdown .sst-account-theme-options button:focus-visible,
.sst-account-dropdown .sst-account-theme-options button.is-selected {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.09);
  color: var(--sst-seal);
}

.sst-account-dropdown a:hover,
.sst-account-dropdown button:hover,
.sst-account-dropdown a:focus-visible,
.sst-account-dropdown button:focus-visible {
  background: rgba(167, 58, 42, 0.08);
  color: var(--sst-seal);
  outline: none;
}

.sst-user-nav {
  position: relative;
  display: flex;
  gap: 0.28rem;
  overflow-x: auto;
  border-bottom: 1px solid rgba(198, 184, 157, 0.22);
  padding: 0.48rem 1rem 0.5rem;
  scrollbar-width: thin;
}

.sst-user-nav::before {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 2;
  width: 1.4rem;
  background: linear-gradient(90deg, rgba(250, 247, 239, 0.96), transparent);
  pointer-events: none;
}

.sst-user-nav-link {
  position: relative;
  display: grid;
  min-width: 7.1rem;
  gap: 0.18rem;
  border: 1px solid transparent;
  border-radius: 9px;
  padding: 0.58rem 0.78rem 0.54rem;
  color: var(--sst-ink-soft);
  transition: background-color 180ms ease, color 180ms ease, transform 180ms ease;
}

.sst-user-nav-link:first-child {
  border-left-color: transparent;
}

.sst-user-nav-link strong {
  font-size: 0.84rem;
  font-weight: 650;
  white-space: nowrap;
}

.sst-user-nav-link::after {
  content: '';
  position: absolute;
  right: 0.78rem;
  bottom: 0.32rem;
  left: 0.78rem;
  height: 2px;
  background: var(--sst-seal);
  opacity: 0;
  transform: scaleX(0.4);
  transform-origin: left;
  transition: opacity 180ms ease, transform 180ms ease;
}

.sst-user-nav-link:hover,
.sst-user-nav-link.is-active {
  border-color: rgba(198, 184, 157, 0.32);
  background: rgba(250, 247, 239, 0.46);
  color: var(--sst-seal);
}

.sst-user-nav-link.is-active {
  cursor: default;
  user-select: none;
}

.sst-user-nav-link.is-active::after {
  opacity: 0.9;
  transform: scaleX(1);
}

.sst-user-content {
  min-width: 0;
  padding: clamp(0.85rem, 1.5vw, 1.15rem);
}

.sst-user-content > :deep(*) {
  min-width: 0;
}

.sst-user-content :deep(.keys-page),
.sst-user-content :deep(.usage-page) {
  max-width: none;
}

.sst-user-content :deep(.console-workbench) {
  max-width: none;
  min-height: min(48rem, calc(100vh - 13.4rem));
}

.sst-user-content :deep(.table-page-layout) {
  max-width: none;
  height: min(48rem, calc(100vh - 13.4rem));
}

.sst-user-content :deep(.table-wrapper) {
  max-width: 100%;
}

/* Respect each user page's authored content width. The previous global
   max-width overrides made night mode feel optically zoomed because many
   page sections were stretched beyond their intended reading measure. */

.sst-user-content :deep(.keys-hero),
.sst-user-content :deep(.usage-hero),
.sst-user-content :deep(.profile-hero) {
  display: none;
}

.sst-user-content :deep(.card),
.sst-user-content :deep(.bg-white),
.sst-user-content :deep(.dark\:bg-dark-800),
.sst-user-content :deep(.rounded-2xl),
.sst-user-content :deep(.rounded-xl) {
  border-color: rgba(198, 184, 157, 0.42);
  background-color: rgba(250, 247, 239, 0.52);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.sst-user-content :deep(.orders-brief),
.sst-user-content :deep(.redeem-brief),
.sst-user-content :deep(.payment-unavailable),
.sst-user-content :deep(.profile-fold),
.sst-user-content :deep(.profile-contact) {
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.028), transparent 28%),
    rgba(250, 247, 239, 0.54);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.sst-user-content :deep(.orders-brief span),
.sst-user-content :deep(.redeem-brief span),
.sst-user-content :deep(.profile-fold summary span) {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.sst-user-content :deep(.orders-brief h2),
.sst-user-content :deep(.redeem-brief h2),
.sst-user-content :deep(.payment-unavailable h2),
.sst-user-content :deep(.profile-fold summary strong) {
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-weight: 600;
}

.sst-user-content :deep(.rounded-2xl),
.sst-user-content :deep(.rounded-xl) {
  border-radius: 10px;
}

.sst-user-content :deep(.shadow-lg),
.sst-user-content :deep(.shadow-xl),
.sst-user-content :deep(.shadow-2xl),
.sst-user-content :deep(.shadow-card),
.sst-user-content :deep(.shadow-card-hover),
.sst-user-content :deep(.shadow-primary-500\/20) {
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.sst-user-content :deep([data-testid="profile-overview-hero"]) {
  border-color: rgba(198, 184, 157, 0.42);
  background:
    radial-gradient(circle at 12% 18%, rgba(167, 58, 42, 0.05), transparent 18rem),
    linear-gradient(135deg, rgba(250, 247, 239, 0.72), rgba(237, 229, 212, 0.48));
}

.sst-user-content :deep([data-testid="profile-overview-hero"] .rounded-\[1\.75rem\]),
.sst-user-content :deep(.profile-contact-icon) {
  border: 1px solid rgba(167, 58, 42, 0.26);
  border-radius: 0.42rem;
  background:
    linear-gradient(135deg, rgba(255, 244, 226, 0.12), transparent 48%),
    var(--sst-seal);
  color: rgba(244, 239, 228, 0.92);
}

.sst-user-content :deep([data-testid="profile-overview-hero"] .rounded-\[1\.75rem\] .text-white),
.sst-user-content :deep([data-testid="profile-overview-hero"] .rounded-\[1\.75rem\] span) {
  color: rgba(244, 239, 228, 0.92);
}

.sst-user-content :deep(.backdrop-blur-xl) {
  backdrop-filter: blur(8px);
}

.sst-user-content :deep(.group.text-left.min-h-\[280px\]),
.sst-user-content :deep(.p-5.rounded-2xl.min-h-\[280px\]) {
  min-height: 15.5rem;
  border-color: rgba(198, 184, 157, 0.42);
  background:
    linear-gradient(180deg, rgba(250, 247, 239, 0.62), rgba(244, 239, 228, 0.36));
  box-shadow: none;
}

.sst-user-content :deep(.group.text-left.min-h-\[280px\]:hover) {
  border-color: rgba(167, 58, 42, 0.28);
  transform: translateY(-1px);
}

.sst-user-content :deep(.font-mono) {
  font-variant-numeric: tabular-nums;
}

.sst-user-content :deep(.input),
.sst-user-content :deep(input),
.sst-user-content :deep(textarea),
.sst-user-content :deep(select) {
  border-color: rgba(198, 184, 157, 0.54);
  background-color: rgba(250, 247, 239, 0.66);
}

.sst-user-content :deep(.input:focus),
.sst-user-content :deep(input:focus),
.sst-user-content :deep(textarea:focus),
.sst-user-content :deep(select:focus) {
  border-color: rgba(167, 58, 42, 0.44);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.1);
  outline: none;
}

.sst-user-content :deep(table) {
  font-variant-numeric: tabular-nums;
}

.sst-user-content :deep(thead) {
  background: rgba(237, 229, 212, 0.72);
}

.sst-user-content :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.035);
}

.sst-user-content :deep(.bg-gradient-to-br) {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.045), transparent 62%),
    linear-gradient(135deg, rgba(250, 247, 239, 0.88), rgba(237, 229, 212, 0.58));
  color: var(--sst-ink);
}

.sst-user-content :deep(.bg-gradient-to-br .text-white),
.sst-user-content :deep(.bg-gradient-to-br .text-indigo-200) {
  color: var(--sst-ink);
}

.sst-user-content :deep(.bg-primary-50),
.sst-user-content :deep(.bg-primary-100),
.sst-user-content :deep(.bg-blue-50),
.sst-user-content :deep(.bg-blue-100),
.sst-user-content :deep(.bg-indigo-50),
.sst-user-content :deep(.bg-indigo-100),
.sst-user-content :deep(.bg-violet-100),
.sst-user-content :deep(.bg-purple-100) {
  background-color: rgba(167, 58, 42, 0.08);
}

.sst-user-content :deep(.bg-green-50),
.sst-user-content :deep(.bg-green-100),
.sst-user-content :deep(.bg-emerald-50),
.sst-user-content :deep(.bg-emerald-100),
.sst-user-content :deep(.bg-teal-50) {
  background-color: rgba(81, 98, 79, 0.1);
}

.sst-user-content :deep(.bg-yellow-100),
.sst-user-content :deep(.bg-orange-100),
.sst-user-content :deep(.bg-amber-100) {
  background-color: rgba(155, 129, 85, 0.14);
}

.sst-user-content :deep(.text-primary-600),
.sst-user-content :deep(.text-primary-700),
.sst-user-content :deep(.text-primary-800),
.sst-user-content :deep(.text-blue-500),
.sst-user-content :deep(.text-blue-600),
.sst-user-content :deep(.text-indigo-500),
.sst-user-content :deep(.text-indigo-600),
.sst-user-content :deep(.text-violet-500),
.sst-user-content :deep(.text-violet-600),
.sst-user-content :deep(.text-purple-600) {
  color: var(--sst-seal);
}

.sst-user-content :deep(.text-green-500),
.sst-user-content :deep(.text-green-600),
.sst-user-content :deep(.text-emerald-500),
.sst-user-content :deep(.text-emerald-600),
.sst-user-content :deep(.text-emerald-700),
.sst-user-content :deep(.text-teal-600) {
  color: #51624f;
}

.sst-user-content :deep(.border-primary-200),
.sst-user-content :deep(.border-blue-200),
.sst-user-content :deep(.border-indigo-200),
.sst-user-content :deep(.border-violet-200),
.sst-user-content :deep(.border-purple-200) {
  border-color: rgba(167, 58, 42, 0.22);
}

.sst-user-content :deep(.border-green-200),
.sst-user-content :deep(.border-emerald-200),
.sst-user-content :deep(.border-teal-200) {
  border-color: rgba(81, 98, 79, 0.22);
}

.sst-user-content :deep(.bg-primary-500),
.sst-user-content :deep(.bg-primary-600),
.sst-user-content :deep(.bg-blue-500),
.sst-user-content :deep(.bg-indigo-500),
.sst-user-content :deep(.bg-violet-500),
.sst-user-content :deep(.bg-purple-500) {
  background-color: var(--sst-seal);
}

.sst-user-content :deep(.bg-green-500),
.sst-user-content :deep(.bg-emerald-500),
.sst-user-content :deep(.bg-teal-500) {
  background-color: #51624f;
}

.sst-user-content :deep(.bg-orange-500),
.sst-user-content :deep(.bg-yellow-500),
.sst-user-content :deep(.bg-amber-500) {
  background-color: #9b8155;
}

.sst-user-content :deep(.animate-spin.border-primary-500),
.sst-user-content :deep(.animate-spin.border-emerald-500) {
  border-color: rgba(167, 58, 42, 0.74);
  border-top-color: transparent;
}

.sst-user-content :deep(.btn-primary) {
  border-color: rgba(167, 58, 42, 0.36);
  background: #a73a2a;
  color: #faf7ef;
  box-shadow: 0 14px 28px -24px rgba(167, 58, 42, 0.62);
}

.sst-user-content :deep(.btn-primary:hover) {
  background: #8f2f23;
}

.dark .sst-user-workbench {
  --sst-paper: rgba(24, 26, 21, 0.9);
  --sst-paper-soft: rgba(24, 26, 21, 0.82);
  --sst-ink: #f4efe4;
  --sst-ink-soft: #d9d0be;
  --sst-mute: #879186;
  --sst-line: rgba(48, 52, 43, 0.95);
  --sst-line-soft: rgba(48, 52, 43, 0.72);
  background: #11130f;
}

.dark .sst-user-frame {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    linear-gradient(180deg, rgba(17, 19, 15, 0.95), rgba(24, 26, 21, 0.9)),
    url('@/assets/brand/sst-paper-ink-bg.png') center/cover;
}

.dark .sst-user-frame::after {
  background:
    radial-gradient(circle at top left, rgba(255, 255, 255, 0.03), transparent 24%),
    radial-gradient(circle at 78% 18%, rgba(167, 58, 42, 0.08), transparent 22%),
    linear-gradient(180deg, transparent, rgba(0, 0, 0, 0.16));
}

.dark .sst-user-nav {
  border-bottom-color: rgba(48, 52, 43, 0.82);
}

.dark .sst-user-nav::before {
  background: linear-gradient(90deg, rgba(17, 19, 15, 0.96), transparent);
}

.dark .sst-user-nav-link:hover,
.dark .sst-user-nav-link.is-active {
  border-color: rgba(167, 58, 42, 0.22);
  background: rgba(167, 58, 42, 0.065);
}

.dark .sst-brand-copy strong,
.dark .sst-account-menu strong,
.dark .sst-user-nav-link strong {
  color: #f4efe4;
}

.dark .sst-account-menu small {
  color: #879186;
}

.dark .sst-account-dropdown {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.98);
}

.dark .sst-account-theme {
  border-color: rgba(48, 52, 43, 0.9);
}

.dark .sst-account-theme > span {
  color: #8f8168;
}

.dark .sst-account-dropdown .sst-account-theme-options button {
  border-color: rgba(48, 52, 43, 0.86);
  background: rgba(17, 19, 15, 0.64);
  color: #c9c0ac;
}

.dark .sst-account-dropdown .sst-account-theme-options button:hover,
.dark .sst-account-dropdown .sst-account-theme-options button:focus-visible,
.dark .sst-account-dropdown .sst-account-theme-options button.is-selected {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.14);
  color: #f0b4a8;
}

.dark .sst-user-nav-link:hover strong,
.dark .sst-user-nav-link.is-active strong {
  color: #f0b4a8;
}

.dark .sst-user-content :deep(.card),
.dark .sst-user-content :deep(.bg-white),
.dark .sst-user-content :deep(.bg-white\/70),
.dark .sst-user-content :deep(.bg-gray-50),
.dark .sst-user-content :deep(.bg-gray-100),
.dark .sst-user-content :deep(.bg-stone-50),
.dark .sst-user-content :deep(.bg-stone-100),
.dark .sst-user-content :deep(.bg-stone-100\/35),
.dark .sst-user-content :deep(.dark\:bg-dark-800),
.dark .sst-user-content :deep(.dark\:bg-dark-800\/60),
.dark .sst-user-content :deep(.dark\:bg-dark-900\/40) {
  border-color: rgba(48, 52, 43, 0.95);
  background-color: rgba(24, 26, 21, 0.78);
}

.dark .sst-user-content :deep(.bg-emerald-100),
.dark .sst-user-content :deep(.bg-green-100),
.dark .sst-user-content :deep(.bg-teal-100) {
  background-color: rgba(81, 98, 79, 0.18);
}

.dark .sst-user-content :deep(.bg-yellow-100),
.dark .sst-user-content :deep(.bg-orange-100),
.dark .sst-user-content :deep(.bg-amber-100) {
  background-color: rgba(155, 129, 85, 0.2);
}

.dark .sst-user-content :deep(.bg-red-100),
.dark .sst-user-content :deep(.bg-rose-100),
.dark .sst-user-content :deep(.bg-primary-100) {
  background-color: rgba(167, 58, 42, 0.18);
}

.dark .sst-user-content :deep(.text-gray-900),
.dark .sst-user-content :deep(.text-gray-800),
.dark .sst-user-content :deep(.text-gray-700) {
  color: #f4efe4;
}

.dark .sst-user-content :deep(.text-gray-600),
.dark .sst-user-content :deep(.text-gray-500),
.dark .sst-user-content :deep(.text-gray-400) {
  color: #9f9787;
}

.dark .sst-user-content :deep(.input),
.dark .sst-user-content :deep(input),
.dark .sst-user-content :deep(textarea),
.dark .sst-user-content :deep(select) {
  border-color: rgba(68, 71, 58, 0.92);
  background-color: rgba(17, 19, 15, 0.68);
  color: #f4efe4;
}

.dark .sst-user-content :deep(.input::placeholder),
.dark .sst-user-content :deep(input::placeholder),
.dark .sst-user-content :deep(textarea::placeholder) {
  color: rgba(201, 192, 172, 0.58);
}

.dark .sst-user-content :deep(.input:focus),
.dark .sst-user-content :deep(input:focus),
.dark .sst-user-content :deep(textarea:focus),
.dark .sst-user-content :deep(select:focus) {
  border-color: rgba(167, 58, 42, 0.48);
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.16);
}

.dark .sst-user-content :deep(thead) {
  background: rgba(54, 50, 39, 0.72);
}

.dark .sst-user-content :deep(tbody tr:hover) {
  background: rgba(167, 58, 42, 0.08);
}

.dark .sst-user-content :deep(.btn-secondary) {
  border-color: rgba(68, 71, 58, 0.92);
  background: rgba(17, 19, 15, 0.62);
  color: #d9d0be;
}

.dark .sst-user-content :deep(.btn-secondary:hover) {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.1);
  color: #f0b4a8;
}

.dark .sst-user-content :deep(.orders-brief),
.dark .sst-user-content :deep(.redeem-brief),
.dark .sst-user-content :deep(.payment-unavailable),
.dark .sst-user-content :deep(.profile-fold),
.dark .sst-user-content :deep(.profile-contact) {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.07), transparent 32%),
    rgba(24, 26, 21, 0.76);
}

@media (max-width: 760px) {
  .sst-admin-stage {
    padding: 0.55rem;
    border-radius: 12px;
  }

  .sst-user-workbench {
    padding: 0.5rem;
  }

  .sst-user-head {
    align-items: start;
    flex-direction: column;
  }

  .sst-account-menu {
    width: 100%;
    min-width: 0;
    align-items: flex-start;
    border-left: 2px solid var(--sst-seal);
    border-right: 0;
    padding-left: 0.8rem;
    padding-right: 0;
    text-align: left;
  }

  .sst-account-trigger {
    margin-top: 0.1rem;
  }

  .sst-account-dropdown {
    right: 0;
  }

  .sst-user-nav {
    padding: 0 0.45rem;
    -webkit-overflow-scrolling: touch;
  }

  .sst-user-nav-link {
    min-width: 6.65rem;
    padding-inline: 0.7rem;
  }

  .sst-user-content {
    padding: 0.65rem;
  }

  .sst-user-content :deep(.console-workbench) {
    min-height: 0;
  }

  .sst-user-content :deep(.table-page-layout) {
    height: auto;
    min-height: 0;
  }

  .sst-user-content :deep(.card) {
    overflow-wrap: anywhere;
  }
}
</style>
<style>
.dark .app-paper {
  background:
    linear-gradient(180deg, rgba(17, 19, 15, 0.94), rgba(24, 26, 21, 0.96)),
    linear-gradient(rgba(244, 239, 228, 0.022) 1px, transparent 1px),
    linear-gradient(90deg, rgba(244, 239, 228, 0.016) 1px, transparent 1px);
  background-size: auto, 112px 112px, 112px 112px;
}

.dark .sst-admin-workbench {
  --sst-admin-line: rgba(48, 52, 43, 0.9);
  --sst-admin-ink: #f4efe4;
  --sst-admin-ink-soft: #d7d0c2;
}

.dark .sst-admin-stage {
  border-color: var(--sst-admin-line);
  background:
    radial-gradient(circle at 92% 6%, rgba(167, 58, 42, 0.07), transparent 16rem),
    radial-gradient(circle at 12% 0%, rgba(244, 239, 228, 0.035), transparent 18rem),
    linear-gradient(180deg, rgba(17, 19, 15, 0.97), rgba(22, 25, 20, 0.93)),
    rgba(17, 19, 15, 0.92);
  box-shadow:
    0 28px 76px -58px rgba(0, 0, 0, 0.92),
    inset 0 1px 0 rgba(244, 239, 228, 0.035);
}

.dark .sst-admin-stage::before {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 14%, transparent 84%, rgba(255, 255, 255, 0.03)),
    linear-gradient(rgba(244, 239, 228, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(244, 239, 228, 0.018) 1px, transparent 1px);
  background-size: auto, 120px 120px, 120px 120px;
}

.dark .sst-admin-stage .card,
.dark .sst-admin-stage .bg-white,
.dark .sst-admin-stage .bg-white\/70,
.dark .sst-admin-stage .bg-white\/80,
.dark .sst-admin-stage .bg-gray-50,
.dark .sst-admin-stage .bg-gray-100,
.dark .sst-admin-stage .bg-stone-50,
.dark .sst-admin-stage .bg-stone-100,
.dark .sst-admin-stage .dark\:bg-dark-800,
.dark .sst-admin-stage .dark\:bg-dark-800\/60,
.dark .sst-admin-stage .dark\:bg-dark-700,
.dark .sst-admin-stage .dark\:bg-dark-900\/40,
.dark .sst-admin-stage .sst-admin-panel {
  border-color: var(--sst-admin-line);
  background-color: rgba(24, 26, 21, 0.78);
  box-shadow: 0 18px 48px -42px rgba(0, 0, 0, 0.72);
}

.dark .sst-admin-stage .btn-secondary {
  border-color: rgba(48, 52, 43, 0.92);
  background: rgba(17, 19, 15, 0.84);
  color: #d7d0c2;
}

.dark .sst-admin-stage h1,
.dark .sst-admin-stage h2,
.dark .sst-admin-stage h3,
.dark .sst-admin-stage .card-title,
.dark .sst-admin-stage .modal-title,
.dark .sst-admin-stage .text-gray-900,
.dark .sst-admin-stage .text-gray-800,
.dark .sst-admin-stage .text-gray-700 {
  color: #f4efe4;
}

.dark .sst-admin-stage .text-gray-600,
.dark .sst-admin-stage .text-gray-500,
.dark .sst-admin-stage .text-gray-400 {
  color: #9f9787;
}

.dark .sst-admin-stage .input,
.dark .sst-admin-stage input,
.dark .sst-admin-stage textarea,
.dark .sst-admin-stage select,
.dark .sst-admin-stage .select-trigger {
  border-color: rgba(68, 71, 58, 0.92);
  background-color: rgba(17, 19, 15, 0.72);
  color: #f4efe4;
}

.dark .sst-admin-stage .input::placeholder,
.dark .sst-admin-stage input::placeholder,
.dark .sst-admin-stage textarea::placeholder {
  color: rgba(201, 192, 172, 0.56);
}

.dark .sst-admin-stage thead {
  background:
    linear-gradient(180deg, rgba(54, 50, 39, 0.78), rgba(36, 36, 29, 0.72));
}

.dark .sst-admin-stage table {
  color: #d7d0c2;
}

.dark .sst-admin-stage table th {
  border-bottom-color: rgba(68, 71, 58, 0.9) !important;
  background:
    linear-gradient(180deg, rgba(54, 50, 39, 0.82), rgba(36, 36, 29, 0.74)) !important;
  color: #b5aa94 !important;
}

.dark .sst-admin-stage table td {
  border-bottom-color: rgba(48, 52, 43, 0.72) !important;
  color: #d7d0c2 !important;
}

.dark .sst-admin-stage tbody tr:hover {
  background: rgba(167, 58, 42, 0.08);
}

.dark .sst-admin-stage .table-container,
.dark .sst-admin-stage .table-wrapper,
.dark .sst-admin-stage .overflow-x-auto:has(table) {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.56);
  box-shadow: inset 0 1px 0 rgba(244, 239, 228, 0.025);
}

.dark .sst-admin-stage .modal-content,
.dark .sst-admin-stage .dialog-container,
.dark .sst-admin-stage .dropdown {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.98);
  color: #f4efe4;
}

.dark .sst-admin-stage .modal-header,
.dark .sst-admin-stage .modal-footer,
.dark .sst-admin-stage .dialog-header,
.dark .sst-admin-stage .dialog-footer {
  border-color: rgba(48, 52, 43, 0.95);
  background-color: rgba(17, 19, 15, 0.72);
}

.dark .sst-admin-stage .btn-secondary:hover {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.1);
  color: #f0b4a8;
}
</style>
