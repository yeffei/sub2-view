<template>
  <aside
    class="sidebar"
    :class="[
      sidebarCollapsed ? 'w-[72px]' : 'w-64',
      { '-translate-x-full lg:translate-x-0': !mobileOpen }
    ]"
  >
    <!-- Logo/Brand -->
    <div class="sidebar-header" :class="{ 'sidebar-header-collapsed': sidebarCollapsed }">
      <!-- Custom Logo or Default Logo -->
      <div class="sidebar-logo flex h-9 w-9 items-center justify-center overflow-hidden rounded-zen shadow-seal">
        <img v-if="settingsLoaded" :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
      </div>
      <div class="sidebar-brand" :class="{ 'sidebar-brand-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">
        <span class="sidebar-brand-title font-serif text-lg font-semibold text-zen-ink dark:text-zen-paper">
          {{ siteName }}
        </span>
        <div class="sidebar-version-wrap">
          <VersionBadge :version="siteVersion" />
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="sidebar-nav scrollbar-hide">
      <!-- Admin View: Admin menu first, then personal menu -->
      <template v-if="isAdmin">
        <!-- Admin Section -->
        <div class="sidebar-section">
          <template v-for="item in adminNavItems" :key="item.path">
            <!-- Collapsible group (has children) -->
            <template v-if="item.children?.length">
              <button
                type="button"
                class="sidebar-link mb-1 w-full"
                :class="{
                  'sidebar-link-active': isGroupActive(item) && !isGroupExpanded(item),
                  'sidebar-link-collapsed': sidebarCollapsed
                }"
                :title="sidebarCollapsed ? item.label : undefined"
                @click="handleGroupClick(item)"
              >
                <Icon v-if="item.icon" :name="item.icon" size="md" class="flex-shrink-0" />
                <span
                  class="sidebar-label sidebar-label-flex"
                  :class="{ 'sidebar-label-collapsed': sidebarCollapsed }"
                  :aria-hidden="sidebarCollapsed ? 'true' : 'false'"
                >
                  <span class="min-w-0 truncate">{{ item.label }}</span>
                  <Icon
                    name="chevronDown"
                    size="sm"
                    class="flex-shrink-0 transition-transform duration-200"
                    :class="isGroupExpanded(item) ? 'rotate-180' : ''"
                  />
                </span>
              </button>
              <!-- Children -->
              <div v-if="!sidebarCollapsed && isGroupExpanded(item)" class="mb-1 ml-4 border-l border-gray-200 pl-2 dark:border-dark-600">
                <router-link
                  v-for="child in item.children"
                  :key="child.path"
                  :to="child.path"
                  class="sidebar-link mb-0.5 py-1.5 text-sm"
                  :class="{ 'sidebar-link-active': route.path === child.path }"
                  @click="handleMenuItemClick(child.path)"
                >
                  <Icon v-if="child.icon" :name="child.icon" size="sm" class="flex-shrink-0" />
                  <span>{{ child.label }}</span>
                </router-link>
              </div>
            </template>
            <!-- Normal item (no children) -->
            <router-link
              v-else
              :to="item.path"
              class="sidebar-link mb-1"
              :class="{ 'sidebar-link-active': isActive(item.path), 'sidebar-link-collapsed': sidebarCollapsed }"
              :title="sidebarCollapsed ? item.label : undefined"
              :id="
                item.path === '/admin/accounts'
                  ? 'sidebar-channel-manage'
                  : item.path === '/admin/groups'
                    ? 'sidebar-group-manage'
                    : item.path === '/admin/redeem'
                      ? 'sidebar-wallet'
                      : undefined
              "
              @click="handleMenuItemClick(item.path)"
            >
              <span v-if="item.iconSvg" class="h-5 w-5 flex-shrink-0 sidebar-svg-icon" v-html="sanitizeSvg(item.iconSvg)"></span>
              <Icon v-else-if="item.icon" :name="item.icon" size="md" class="flex-shrink-0" />
              <span class="sidebar-label" :class="{ 'sidebar-label-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">{{ item.label }}</span>
            </router-link>
          </template>
        </div>

        <!-- Personal Section for Admin (hidden in simple mode) -->
        <div v-if="!authStore.isSimpleMode" class="sidebar-section">
          <div class="sidebar-section-title" :class="{ 'sidebar-section-title-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">
            <span class="sidebar-section-title-text" :class="{ 'sidebar-section-title-text-collapsed': sidebarCollapsed }">
              {{ t('nav.myAccount') }}
            </span>
          </div>

          <router-link
            v-for="item in personalNavItems"
            :key="item.path"
            :to="item.path"
            class="sidebar-link mb-1"
            :class="{ 'sidebar-link-active': isActive(item.path), 'sidebar-link-collapsed': sidebarCollapsed }"
            :title="sidebarCollapsed ? item.label : undefined"
            :data-tour="item.path === '/keys' ? 'sidebar-my-keys' : undefined"
            @click="handleMenuItemClick(item.path)"
          >
            <span v-if="item.iconSvg" class="h-5 w-5 flex-shrink-0 sidebar-svg-icon" v-html="sanitizeSvg(item.iconSvg)"></span>
            <Icon v-else-if="item.icon" :name="item.icon" size="md" class="flex-shrink-0" />
            <span class="sidebar-label" :class="{ 'sidebar-label-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">{{ item.label }}</span>
          </router-link>
        </div>
      </template>

      <!-- Regular User View -->
      <template v-else-if="!appStore.backendModeEnabled">
        <div v-for="section in userNavSections" :key="section.title" class="sidebar-section">
          <div class="sidebar-section-title" :class="{ 'sidebar-section-title-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">
            <span class="sidebar-section-title-text" :class="{ 'sidebar-section-title-text-collapsed': sidebarCollapsed }">
              {{ section.title }}
            </span>
          </div>
          <router-link
            v-for="item in section.items"
            :key="item.path"
            :to="item.path"
            class="sidebar-link mb-1"
            :class="{ 'sidebar-link-active': isActive(item.path), 'sidebar-link-collapsed': sidebarCollapsed }"
            :title="sidebarCollapsed ? item.label : undefined"
            :data-tour="item.path === '/keys' ? 'sidebar-my-keys' : undefined"
            @click="handleMenuItemClick(item.path)"
          >
            <span v-if="item.iconSvg" class="h-5 w-5 flex-shrink-0 sidebar-svg-icon" v-html="sanitizeSvg(item.iconSvg)"></span>
            <Icon v-else-if="item.icon" :name="item.icon" size="md" class="flex-shrink-0" />
            <span class="sidebar-label" :class="{ 'sidebar-label-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">{{ item.label }}</span>
          </router-link>
        </div>
      </template>
    </nav>

    <!-- Bottom Section -->
    <div class="mt-auto border-t border-gray-100 p-3 dark:border-dark-800">
      <!-- Theme Toggle -->
      <button
        @click="toggleTheme"
        class="sidebar-link mb-2 w-full"
        :class="{ 'sidebar-link-collapsed': sidebarCollapsed }"
        :title="sidebarCollapsed ? (isDark ? t('nav.lightMode') : t('nav.darkMode')) : undefined"
      >
        <Icon v-if="isDark" name="sun" size="md" class="flex-shrink-0 text-amber-500" />
        <Icon v-else name="moon" size="md" class="flex-shrink-0" />
        <span class="sidebar-label" :class="{ 'sidebar-label-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">{{
          isDark ? t('nav.lightMode') : t('nav.darkMode')
        }}</span>
      </button>

      <!-- Collapse Button -->
      <button
        @click="toggleSidebar"
        class="sidebar-link w-full"
        :class="{ 'sidebar-link-collapsed': sidebarCollapsed }"
        :title="sidebarCollapsed ? t('nav.expand') : t('nav.collapse')"
      >
        <Icon v-if="!sidebarCollapsed" name="chevronLeft" size="md" class="flex-shrink-0" :stroke-width="1.8" />
        <Icon v-else name="chevronRight" size="md" class="flex-shrink-0" :stroke-width="1.8" />
        <span class="sidebar-label" :class="{ 'sidebar-label-collapsed': sidebarCollapsed }" :aria-hidden="sidebarCollapsed ? 'true' : 'false'">{{ t('nav.collapse') }}</span>
      </button>
    </div>
  </aside>

  <!-- Mobile Overlay -->
  <transition name="fade">
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-30 bg-black/50 lg:hidden"
      @click="closeMobile"
    ></div>
  </transition>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAdminSettingsStore, useAppStore, useAuthStore, useOnboardingStore } from '@/stores'
import VersionBadge from '@/components/common/VersionBadge.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeSvg } from '@/utils/sanitize'
import { FeatureFlags, makeSidebarFlag } from '@/utils/featureFlags'
import { IMAGE_WORKSHOP_MENU_ID, findImageWorkshopMenuItem, isImageWorkshopMenuItem } from '@/utils/imageWorkshop'
import { initTheme, toggleTheme as toggleDocumentTheme, useThemeState } from '@/utils/theme'

type SidebarIconName = InstanceType<typeof Icon>['$props']['name']

interface NavItem {
  path: string
  label: string
  icon: SidebarIconName | null
  iconSvg?: string
  hideInSimpleMode?: boolean
  children?: NavItem[]
  /**
   * When true, the parent item only toggles the expand/collapse state and
   * does NOT navigate to its `path`. The `path` is purely a stable key.
   */
  expandOnly?: boolean
  /**
   * 可选的功能开关 getter。返回 false 时菜单项被隐藏；返回 undefined/true 时显示。
   * 宽容策略（undefined → 显示）避免 public settings 未加载完成时菜单闪烁消失。
   * Getter 里访问的 reactive 来源（store / composable）会被 computed 自动追踪，
   * 开关切换时菜单自动更新。
   */
  featureFlag?: () => boolean | undefined
}

interface NavSection {
  title: string
  items: NavItem[]
}

// applyFeatureFlags 递归过滤掉 featureFlag() === false 的节点（含子节点）。
// 使用 `!== false` 宽容语义：undefined（设置未加载）或 true 都视为显示。
function applyFeatureFlags(items: NavItem[]): NavItem[] {
  const out: NavItem[] = []
  for (const item of items) {
    if (item.featureFlag && item.featureFlag() === false) continue
    if (item.children) {
      out.push({ ...item, children: applyFeatureFlags(item.children) })
    } else {
      out.push(item)
    }
  }
  return out
}

const { t } = useI18n()

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const onboardingStore = useOnboardingStore()
const adminSettingsStore = useAdminSettingsStore()

const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const mobileOpen = computed(() => appStore.mobileOpen)
const isAdmin = computed(() => authStore.isAdmin)
const isDark = useThemeState()

// Track which parent nav groups are expanded
const expandedGroups = ref<Set<string>>(new Set())

// Site settings from appStore (cached, no flicker)
const siteName = computed(() => appStore.siteName === 'Sub2API' ? '山枢庭' : appStore.siteName)
const siteLogo = computed(() => appStore.siteLogo)
const siteVersion = computed(() => appStore.siteVersion)
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

// Public-settings flags go through the registry in utils/featureFlags.ts,
// which handles the opt-in vs opt-out fallback when settings haven't loaded
// yet. Admin-only flags (not in public settings) stay inline below.
const flagChannelMonitor = makeSidebarFlag(FeatureFlags.channelMonitor)
const flagPayment = makeSidebarFlag(FeatureFlags.payment)
const flagAffiliate = makeSidebarFlag(FeatureFlags.affiliate)
const flagRiskControl = makeSidebarFlag(FeatureFlags.riskControl)
const flagOpsMonitoring = () => adminSettingsStore.opsMonitoringEnabled
const flagAdminPayment = () => adminSettingsStore.paymentEnabled

// buildSelfNavItems 构造用户自己的导航项（用户端主菜单和管理员的"我的账户"子菜单共享这组声明）。
// withDashboard=true 时包含仪表盘（用户端），false 时不含（管理员的个人区已经有独立仪表盘入口）。
//
// 条目顺序：密钥 → 用量 → 可用渠道 → 渠道状态 → 订阅/支付 → 兑换/资料。
// 可用渠道紧挨渠道状态之上，让用户"先看自己能用什么、再看对应状态"。
function buildSelfNavItems(withDashboard: boolean): NavItem[] {
  const items: NavItem[] = []
  if (withDashboard) {
    items.push({ path: '/dashboard', label: t('nav.dashboard'), icon: 'grid' })
  }
  items.push(
    { path: '/keys', label: t('nav.apiKeys'), icon: 'key' },
    { path: '/usage', label: '用量账册', icon: 'chart', hideInSimpleMode: true },
    { path: '/monitor', label: '服务状态', icon: 'signal', featureFlag: flagChannelMonitor },
    { path: '/purchase', label: '充值与兑换', icon: 'wallet', hideInSimpleMode: true, featureFlag: flagPayment },
    { path: '/orders', label: '往来订单', icon: 'document', hideInSimpleMode: true, featureFlag: flagPayment },
    { path: '/affiliate', label: '团队引荐', icon: 'users', hideInSimpleMode: true, featureFlag: flagAffiliate },
    { path: '/profile', label: t('nav.profile'), icon: 'user' },
    ...customMenuItemsForUser.value.map((item): NavItem => ({
      path: `/custom/${item.id}`,
      label: item.label,
      icon: null,
      iconSvg: item.icon_svg,
    })),
  )
  return items
}

// finalizeNav 合并三重过滤：featureFlag 过滤 + simple 模式过滤。
function finalizeNav(items: NavItem[]): NavItem[] {
  const visible = applyFeatureFlags(items)
  if (!authStore.isSimpleMode) return visible
  const filterSimple = (list: NavItem[]): NavItem[] => list
    .filter(item => !item.hideInSimpleMode)
    .map(item => item.children ? { ...item, children: filterSimple(item.children) } : item)
    .filter(item => !item.children || item.children.length > 0)
  return filterSimple(visible)
}

function finalizeSectionItems(items: NavItem[]): NavItem[] {
  return finalizeNav(items)
}

// User navigation sections (for regular users)
const userNavSections = computed((): NavSection[] => {
  const customItems = customMenuItemsForUser.value
    .filter((item) => !isImageWorkshopMenuItem(item))
    .map((item): NavItem => ({
      path: `/custom/${item.id}`,
      label: item.label,
      icon: null,
      iconSvg: item.icon_svg,
    }))
  const imageWorkshopItem = findImageWorkshopMenuItem(customMenuItemsForUser.value)

  return [
    {
      title: '庭中概览',
      items: finalizeSectionItems([
        { path: '/dashboard', label: '今日概览', icon: 'grid' },
      ]),
    },
    {
      title: '密钥与调用',
      items: finalizeSectionItems([
        { path: '/keys', label: t('nav.apiKeys'), icon: 'key' },
      ]),
    },
    {
      title: '用量与账册',
      items: finalizeSectionItems([
        { path: '/usage', label: '用量账册', icon: 'chart', hideInSimpleMode: true },
        { path: '/purchase', label: '充值与兑换', icon: 'wallet', hideInSimpleMode: true, featureFlag: flagPayment },
        { path: '/orders', label: '往来订单', icon: 'document', hideInSimpleMode: true, featureFlag: flagPayment },
      ]),
    },
    {
      title: '能力与状态',
      items: finalizeSectionItems([
        ...(imageWorkshopItem ? [{ path: `/custom/${IMAGE_WORKSHOP_MENU_ID}`, label: '图像工坊', icon: 'image' as const }] : []),
        { path: '/monitor', label: '服务状态', icon: 'signal', featureFlag: flagChannelMonitor },
      ]),
    },
    {
      title: '账户与规则',
      items: finalizeSectionItems([
        { path: '/profile', label: t('nav.profile'), icon: 'user' },
        { path: '/affiliate', label: '团队引荐', icon: 'users', hideInSimpleMode: true, featureFlag: flagAffiliate },
        ...customItems,
      ]),
    },
  ].filter(section => section.items.length > 0)
})

// Personal navigation items (for admin's "My Account" section, without Dashboard).
// Admins access 可用渠道 from this section just like regular users — there is no
// separate admin entry, since the page is purely a user-facing view.
const personalNavItems = computed((): NavItem[] => finalizeNav(buildSelfNavItems(false)))

// Custom menu items filtered by visibility
const customMenuItemsForUser = computed(() => {
  const items = appStore.cachedPublicSettings?.custom_menu_items ?? []
  return items
    .filter((item) => item.visibility === 'user')
    .sort((a, b) => a.sort_order - b.sort_order)
})

const customMenuItemsForAdmin = computed(() => {
  return adminSettingsStore.customMenuItems
    .filter((item) => item.visibility === 'admin')
    .sort((a, b) => a.sort_order - b.sort_order)
})

// Admin navigation items
const adminNavItems = computed((): NavItem[] => {
  const baseItems: NavItem[] = [
    { path: '/admin/dashboard', label: t('nav.dashboard'), icon: 'grid' },
    { path: '/admin/ops', label: t('nav.ops'), icon: 'chart', featureFlag: flagOpsMonitoring },
    { path: '/admin/users', label: t('nav.users'), icon: 'users', hideInSimpleMode: true },
    { path: '/admin/groups', label: t('nav.groups'), icon: 'folder', hideInSimpleMode: true },
    { path: '/admin/risk-control', label: t('nav.riskControl'), icon: 'shield', hideInSimpleMode: true, featureFlag: flagRiskControl },
    { path: '/admin/accounts', label: t('nav.accounts'), icon: 'globe', hideInSimpleMode: true },
    { path: '/admin/channels/upstream-pools', label: t('nav.upstreamPools'), icon: 'server', hideInSimpleMode: true },
    { path: '/admin/channels/monitor', label: t('nav.channelMonitor'), icon: 'signal', hideInSimpleMode: true, featureFlag: flagChannelMonitor },
    { path: '/admin/proxies', label: t('nav.proxies'), icon: 'server', hideInSimpleMode: true },
    { path: '/admin/usage', label: '计量记录', icon: 'chart' },
    { path: '/admin/affiliates', label: '团队引荐', icon: 'users', hideInSimpleMode: true, featureFlag: flagAffiliate },
    {
      path: '/admin/orders',
      label: '计量与交易',
      icon: 'document',
      hideInSimpleMode: true,
      featureFlag: flagAdminPayment,
    },
    { path: '/admin/announcements', label: t('nav.announcements'), icon: 'bell' },
    { path: '/admin/redeem', label: t('nav.redeemCodes'), icon: 'ticket', hideInSimpleMode: true },
    { path: '/admin/settings', label: t('nav.settings'), icon: 'cog' },
  ]

  const visible = applyFeatureFlags(baseItems)

  // 简单模式下，在系统设置前插入 API密钥
  if (authStore.isSimpleMode) {
    const filtered = visible.filter(item => !item.hideInSimpleMode)
    filtered.push({ path: '/keys', label: t('nav.apiKeys'), icon: 'key' })
    filtered.push({ path: '/admin/settings', label: t('nav.settings'), icon: 'cog' })
    for (const cm of customMenuItemsForAdmin.value) {
      filtered.push({ path: `/custom/${cm.id}`, label: cm.label, icon: null, iconSvg: cm.icon_svg })
    }
    return filtered
  }

  for (const cm of customMenuItemsForAdmin.value) {
    visible.push({ path: `/custom/${cm.id}`, label: cm.label, icon: null, iconSvg: cm.icon_svg })
  }
  return visible
})

function toggleSidebar() {
  appStore.toggleSidebar()
}

function toggleTheme() {
  isDark.value = toggleDocumentTheme(isDark.value)
}

function closeMobile() {
  appStore.setMobileOpen(false)
}

function handleMenuItemClick(itemPath: string) {
  if (mobileOpen.value) {
    setTimeout(() => {
      appStore.setMobileOpen(false)
    }, 150)
  }

  // Map paths to tour selectors
  const pathToSelector: Record<string, string> = {
    '/admin/groups': '#sidebar-group-manage',
    '/admin/accounts': '#sidebar-channel-manage',
    '/keys': '[data-tour="sidebar-my-keys"]'
  }

  const selector = pathToSelector[itemPath]
  if (selector && onboardingStore.isCurrentStep(selector)) {
    onboardingStore.nextStep(500)
  }
}

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

function isGroupActive(item: NavItem): boolean {
  if (!item.children) return false
  return item.children.some(child => route.path === child.path)
}

function isGroupExpanded(item: NavItem): boolean {
  return expandedGroups.value.has(item.path) || isGroupActive(item)
}

function toggleGroup(item: NavItem) {
  if (expandedGroups.value.has(item.path)) {
    expandedGroups.value.delete(item.path)
  } else {
    expandedGroups.value.add(item.path)
  }
}

/**
 * Click handler for collapsible parent items.
 * - When sidebar is collapsed: do nothing (children are not visible).
 * - When `expandOnly` is true: only toggle expand state.
 * - Otherwise (default, e.g. /admin/orders): navigate to the parent path
 *   (router-link semantics) and ensure the group is expanded.
 */
function handleGroupClick(item: NavItem) {
  if (sidebarCollapsed.value) return
  if (item.expandOnly) {
    toggleGroup(item)
    return
  }
  // Push to path and ensure expanded
  if (route.path !== item.path) {
    router.push(item.path)
  }
  if (!expandedGroups.value.has(item.path)) {
    expandedGroups.value.add(item.path)
  }
}

// Fetch admin settings (for feature-gated nav items like Ops).
watch(
  isAdmin,
  (v) => {
    if (v) {
      adminSettingsStore.fetch()
    }
  },
  { immediate: true }
)

onMounted(() => {
  initTheme()
  if (isAdmin.value) {
    adminSettingsStore.fetch()
  }
})
</script>

<style scoped>
.sidebar-logo {
  display: grid;
  flex: 0 0 2.25rem;
  min-width: 2.25rem;
  place-items: center;
  overflow: hidden;
  border-radius: 0.7rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 12px 24px -18px rgba(31, 35, 32, 0.42);
}

.sidebar-version-wrap {
  margin-top: 0.45rem;
}

.sidebar-header-collapsed {
  gap: 0;
  padding-left: 1.125rem;
  padding-right: 1.125rem;
}

.sidebar-brand {
  min-width: 0;
  flex: 1 1 auto;
  white-space: nowrap;
  transition:
    max-width 0.22s ease,
    opacity 0.14s ease,
    transform 0.14s ease;
  max-width: 12rem;
}

.sidebar-brand-collapsed {
  max-width: 0;
  overflow: hidden;
  opacity: 0;
  transform: translateX(-4px);
  pointer-events: none;
}

.sidebar-brand-title {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.02;
  letter-spacing: 0.01em;
}

.sidebar-link-collapsed {
  gap: 0;
  padding-left: 0.875rem;
  padding-right: 0.875rem;
}

.sidebar-section-title {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 1.25rem;
  overflow: hidden;
  white-space: nowrap;
}

.sidebar-section-title-text {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition:
    opacity 0.16s ease,
    transform 0.16s ease;
}

.sidebar-section-title::after {
  content: '';
  position: absolute;
  left: 0.75rem;
  right: 0.75rem;
  top: 50%;
  height: 1px;
  background: rgb(229 231 235);
  opacity: 0;
  transform: translateY(-50%);
  transition: opacity 0.18s ease;
}

.sidebar-section-title-text-collapsed {
  opacity: 0;
  transform: translateX(-4px);
}

.sidebar-section-title-collapsed::after {
  opacity: 1;
  transition-delay: 0.08s;
}

.sidebar-label {
  display: block;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition:
    max-width 0.2s ease,
    opacity 0.12s ease,
    transform 0.12s ease;
  max-width: 12rem;
}

.sidebar-label-flex {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.sidebar-label-collapsed {
  max-width: 0;
  opacity: 0;
  transform: translateX(-4px);
  pointer-events: none;
}

/* Custom SVG icon in sidebar: constrain size without overriding uploaded SVG colors */
.sidebar-svg-icon {
  color: currentColor;
}

.sidebar-svg-icon :deep(svg) {
  display: block;
  width: 1.25rem;
  height: 1.25rem;
}

.sidebar-version-wrap :deep(.version-badge-trigger) {
  min-height: 2.05rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 0.95rem;
  padding: 0.38rem 0.82rem;
  background:
    linear-gradient(180deg, rgba(255, 249, 238, 0.92), rgba(242, 232, 214, 0.88));
  color: #8c4a27 !important;
  box-shadow:
    0 14px 28px -24px rgba(82, 62, 39, 0.34),
    inset 0 1px 0 rgba(255, 255, 255, 0.7);
}

.sidebar-version-wrap :deep(.version-badge-trigger:hover) {
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.96), rgba(245, 236, 220, 0.92));
}
</style>
<style>
.dark .sidebar-logo {
  border: 1px solid rgba(54, 59, 50, 0.9);
  background:
    radial-gradient(circle at 34% 30%, rgba(255, 247, 228, 0.05), transparent 58%),
    linear-gradient(180deg, rgba(24, 27, 22, 0.98), rgba(31, 34, 28, 0.94));
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.04),
    0 16px 30px -22px rgba(0, 0, 0, 0.62);
}

.dark .sidebar-brand-title {
  color: #f1ebdf;
  text-shadow: 0 1px 0 rgba(0, 0, 0, 0.28);
}

.dark .sidebar-version-wrap :deep(.version-badge-trigger) {
  border-color: rgba(96, 72, 39, 0.78);
  background:
    linear-gradient(180deg, rgba(93, 72, 48, 0.88), rgba(72, 53, 34, 0.92));
  color: #ffbf2f !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 240, 214, 0.08),
    0 16px 28px -22px rgba(0, 0, 0, 0.58);
}

.dark .sidebar-version-wrap :deep(.version-badge-trigger:hover) {
  background:
    linear-gradient(180deg, rgba(105, 80, 53, 0.9), rgba(79, 58, 37, 0.94));
}

.dark .sidebar-version-wrap :deep(.version-badge-trigger .version-badge-value) {
  color: #ffbf2f;
}

.dark .sidebar-section-title::after {
  background: rgb(55 65 81);
}
</style>
