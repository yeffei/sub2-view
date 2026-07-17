<template>
  <AppLayout>
    <div class="sst-dashboard" :class="{ 'sst-dashboard-en': isEnglishDashboard }">
      <AnnouncementBell buttonless />
      <div v-if="loading" class="home-state" :class="isDark ? 'is-night' : 'is-paper'">
        <div class="seal-mark"><img src="/logo.png" alt="" aria-hidden="true" /></div>
        <LoadingSpinner />
        <div>
          <h1>{{ dashboardCopy.loadingTitle }}</h1>
          <p>{{ dashboardCopy.loadingCopy }}</p>
        </div>
      </div>

      <template v-else-if="stats">
        <section class="courtyard-console" :aria-label="dashboardCopy.consoleAria">
          <header class="console-head">
            <router-link to="/home" class="brand-lockup" :aria-label="t('userShell.backHomeAria')">
              <span class="seal-mark" aria-hidden="true"><img src="/logo.png" alt="" /></span>
              <div class="brand-copy">
                <span>{{ t('brand.defaultName') }}</span>
                <h1>{{ t('publicSite.tagline') }}</h1>
              </div>
            </router-link>
            <div
              ref="accountMenuRef"
              class="account-mark account-menu"
              :class="{ 'is-open': isAccountMenuOpen }"
            >
              <div>
                <span>{{ currentDateLabel }}</span>
                <strong>{{ user?.email || t('userShell.defaultAccount') }}</strong>
                <small>{{ healthLabel }} · {{ todayRequestLabel }}</small>
              </div>
              <button
                type="button"
                class="account-menu-trigger"
                :aria-label="t('userShell.accountMenu')"
                aria-haspopup="menu"
                :aria-expanded="isAccountMenuOpen ? 'true' : 'false'"
                @click="toggleAccountMenu"
              >
                {{ t('userShell.accountMark') }}
              </button>
              <div class="account-menu-dropdown">
                <router-link to="/profile" @click="closeAccountMenu">{{ t('userShell.profileDocument') }}</router-link>
                <div class="account-theme" :aria-label="t('userShell.themeSettings')">
                  <span>{{ t('userShell.themeSettings') }}</span>
                  <div class="account-theme-options">
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
                <router-link to="/dashboard" @click="closeAccountMenu">{{ t('userShell.backDashboard') }}</router-link>
                <button type="button" @click="handleLogout">{{ t('nav.logout') }}</button>
              </div>
            </div>
          </header>

          <section class="courtyard-stage" :aria-label="dashboardCopy.stageAria">
            <div class="courtyard-map watch-desk">
              <div class="ink-wash ink-wash-map" aria-hidden="true"></div>

              <div class="watch-main">
                <div class="watch-status" :class="`status-${healthTone}`">
                  <div class="center-seal watch-score">
                    <span>{{ watchSealLabel }}</span>
                    <strong>{{ watchSealValue }}</strong>
                    <small>{{ watchSealNote }}</small>
                  </div>

                  <div class="watch-copy">
                    <span>{{ dashboardCopy.watchStatus }}</span>
                    <h2>{{ statusTitle }}</h2>
                    <p>{{ statusSummary }}</p>
                  </div>
                </div>

                <nav class="gate-grid" :aria-label="dashboardCopy.gateAria">
                  <router-link
                    v-for="gate in courtyardGates"
                    :key="gate.to"
                    :to="gate.to"
                    class="gate-link"
                    :class="gate.position"
                  >
                    <Icon :name="gate.icon" size="sm" />
                    <span>{{ gate.label }}</span>
                    <small>{{ gate.note }}</small>
                  </router-link>
                </nav>

              </div>

              <div class="watch-reasons" :aria-label="dashboardCopy.reasonsAria">
                <div class="watch-advice" v-if="watchAdvices.length">
                  <span>{{ dashboardCopy.nextStep }}</span>
                  <router-link
                    v-for="advice in watchAdvices"
                    :key="advice.title"
                    :to="advice.to"
                    :class="`tone-${advice.tone}`"
                  >
                    <strong>{{ advice.title }}</strong>
                    <small>{{ advice.detail }}</small>
                  </router-link>
                </div>
                <router-link
                  v-for="reason in statusReasons"
                  :key="reason.label"
                  :to="reason.to"
                  class="reason-item"
                  :class="`tone-${reason.tone}`"
                >
                  <Icon :name="reason.icon" size="sm" />
                  <span>
                    <strong>{{ reason.label }}</strong>
                    <small>{{ reason.detail }}</small>
                  </span>
                  <em>{{ reasonActionLabel(reason.tone) }}</em>
                </router-link>
                <div v-if="!statusReasons.length" class="reason-item tone-calm">
                  <Icon name="chart" size="sm" />
                  <span>
                    <strong>{{ dashboardCopy.stableReasonTitle }}</strong>
                    <small>{{ dashboardCopy.stableReasonDetail }}</small>
                  </span>
                  <em>{{ dashboardCopy.calmAction }}</em>
                </div>
              </div>

              <section class="health-check-sheet watch-health-sheet" :aria-label="dashboardCopy.healthAria">
                <div class="section-mark">
                  <span>{{ dashboardCopy.healthTitle }}</span>
                  <router-link to="/keys?panel=connection-test">{{ dashboardCopy.fullCheck }}</router-link>
                </div>
                <div class="health-check-grid">
                  <article
                    v-for="item in connectionCheckItems"
                    :key="item.key"
                    class="check-card"
                    :class="`tone-${item.tone}`"
                  >
                    <div class="check-card-head">
                      <span>{{ item.label }}</span>
                      <strong>{{ item.status }}</strong>
                    </div>
                    <p>{{ item.detail }}</p>
                    <router-link :to="item.to">{{ item.action }}</router-link>
                  </article>
                </div>
              </section>
            </div>

            <aside class="ledger-slips" :aria-label="dashboardCopy.ledgerAria">
              <div class="slips-head">
                <span>{{ dashboardCopy.ledgerTitle }}</span>
                <small>{{ ledgerSummary }}</small>
              </div>
              <router-link
                v-for="item in gardenEntries"
                :key="item.to + item.label"
                :to="item.to"
                class="ledger-slip"
                :class="{ 'is-alert': item.alert }"
              >
                <Icon :name="item.icon" size="sm" />
                <div class="ledger-slip-copy">
                  <div class="ledger-slip-main">
                    <span>{{ item.label }}</span>
                    <strong>{{ item.value }}</strong>
                  </div>
                  <small>{{ item.note }}</small>
                </div>
              </router-link>
            </aside>
          </section>

          <section
            ref="requestsFocusRef"
            class="water-ledger"
            :class="{ 'focus-requests': isRequestsFocusActive }"
            :aria-label="dashboardCopy.waterLedgerAria"
          >
            <div class="usage-scroll" :class="{ 'focus-surface': isRequestsFocusActive }">
              <div class="section-mark">
                <span>{{ dashboardCopy.callsTitle }}</span>
                <router-link to="/usage">{{ dashboardCopy.allUsage }}</router-link>
              </div>
              <div v-if="requestsFocusEnabled" class="focus-note" role="status">
                <div>
                  <span>{{ dashboardCopy.requestFocus }}</span>
                  <strong>{{ requestsFocusTitle }}</strong>
                  <p>{{ requestsFocusDetail }}</p>
                </div>
                <div class="focus-note-actions">
                  <router-link to="/usage?tab=errors&category=rate_limit">{{ dashboardCopy.rateLimitOnly }}</router-link>
                  <router-link :to="recentUsage.length ? '/usage' : '/keys?panel=connection-test'">
                    {{ recentUsage.length ? dashboardCopy.reviewAllLedger : dashboardCopy.doConnectionCheck }}
                  </router-link>
                </div>
              </div>
              <div v-if="loadingUsage" class="mini-state">{{ dashboardCopy.loadingUsage }}</div>
              <div v-else-if="!recentUsage.length" class="empty-note">
                <strong>{{ dashboardCopy.noRecentUsageTitle }}</strong>
                <span>{{ dashboardCopy.noRecentUsageCopy }}</span>
                <div class="empty-actions">
                  <router-link to="/keys">{{ dashboardCopy.manageKeys }}</router-link>
                  <router-link to="/monitor">{{ dashboardCopy.checkService }}</router-link>
                  <router-link to="/usage">{{ dashboardCopy.viewLedger }}</router-link>
                </div>
              </div>
              <ol v-else class="call-list">
                <li v-for="log in recentUsage" :key="log.id">
                  <div>
                    <strong>{{ log.model }}</strong>
                    <span>{{ formatDateTime(log.created_at) }}</span>
                  </div>
                  <div>
                    <strong>${{ formatCost(log.actual_cost) }}</strong>
                    <span>{{ formatTokens(log.input_tokens + log.output_tokens) }} tokens</span>
                  </div>
                </li>
              </ol>
              <div v-if="recentUsage.length" class="call-list-footer">
                <div>
                  <small>{{ recentUsageFooter }}</small>
                </div>
                <div class="call-list-footer-actions">
                  <router-link to="/usage">{{ dashboardCopy.fullLedger }}</router-link>
                  <router-link to="/usage?tab=errors">{{ dashboardCopy.errorLedger }}</router-link>
                </div>
              </div>
            </div>

            <div class="ledger-side">
              <button
                type="button"
                class="notice-strip"
                :aria-label="dashboardCopy.openAnnouncements"
                @click="openAnnouncementCenter"
              >
                <div class="notice-strip-copy">
                  <span>{{ dashboardCopy.announcements }}</span>
                  <strong>{{ announcementSummary.title }}</strong>
                  <small>{{ announcementSummary.note }}</small>
                </div>
                <div class="notice-strip-meta">
                  <em>{{ announcementSummary.badge }}</em>
                  <small>{{ dashboardCopy.viewAllAnnouncements }}</small>
                </div>
              </button>

              <div class="folio flow-folio">
                <div class="section-mark">
                  <span>{{ dashboardCopy.lastSevenDays }}</span>
                  <strong>{{ trendPeakLabel }}</strong>
                </div>
                <div class="waterline" :aria-label="dashboardCopy.tokenTrendAria">
                  <span
                    v-for="point in trendBars"
                    :key="point.date"
                    :style="{ height: point.height }"
                    :class="{ 'is-peak': point.isPeak, 'is-empty': point.empty }"
                    :title="`${point.date}: ${point.label}`"
                  ></span>
                </div>
                <div class="waterline-axis" aria-hidden="true">
                  <span v-for="point in trendBars" :key="`${point.date}-axis`">{{ point.weekday }}</span>
                </div>
              </div>

              <div class="folio quota-folio">
                <div class="section-mark">
                  <span>{{ dashboardCopy.window }}</span>
                  <span>{{ dashboardCopy.accountOverview }}</span>
                </div>
                <div class="quota-list">
                  <div
                    v-for="(quota, index) in quotaSummary"
                    :key="quota.label"
                    :class="{ 'quota-inline-row': index === 0 }"
                  >
                    <span>{{ quota.label }}</span>
                    <strong>{{ quota.value }}</strong>
                    <small v-if="quota.note">{{ quota.note }}</small>
                  </div>
                </div>
              </div>

              <div class="folio models-folio">
                <div class="section-mark">
                  <span>{{ dashboardCopy.models }}</span>
                  <router-link to="/monitor">{{ dashboardCopy.channelStatus }}</router-link>
                </div>
                <div v-if="loadingCharts" class="mini-state">{{ dashboardCopy.loadingModels }}</div>
                <div v-else-if="!modelPreview.length" class="mini-state">{{ dashboardCopy.noModelDistribution }}</div>
                <div v-else class="model-river">
                  <div v-for="model in modelPreview" :key="model.model">
                    <span>{{ model.model }}</span>
                    <strong>{{ formatTokens(model.total_tokens) }}</strong>
                    <i :style="{ width: modelShare(model.total_tokens) }"></i>
                  </div>
                </div>
              </div>

              <div class="folio platforms-folio">
                <div class="section-mark">
                  <span>{{ dashboardCopy.platforms }}</span>
                  <router-link to="/monitor">{{ dashboardCopy.status }}</router-link>
                </div>
                <div v-if="!platformPreview.length" class="mini-state">{{ dashboardCopy.noPlatformRecords }}</div>
                <div v-else class="platform-list">
                  <div v-for="platform in platformPreview" :key="platform.platform">
                    <span>{{ platformLabel(platform.platform) }}</span>
                    <strong>${{ formatCost(platform.total_actual_cost) }}</strong>
                    <small>{{ platformUsageLabel(platform.total_requests, platform.total_tokens) }}</small>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </section>
      </template>

      <div v-else class="home-state home-state-error" :class="isDark ? 'is-night' : 'is-paper'">
        <div class="seal-mark"><img src="/logo.png" alt="" aria-hidden="true" /></div>
        <h1>{{ dashboardCopy.errorTitle }}</h1>
        <p>{{ errorMessage || dashboardCopy.errorCopy }}</p>
        <button type="button" class="btn btn-primary" @click="refreshAll">{{ dashboardCopy.retry }}</button>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useAnnouncementStore } from '@/stores/announcements'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { keysAPI } from '@/api'
import { usageAPI, type ApiKeyWorkbenchSummary, type UserDashboardStats as UserStatsType } from '@/api/usage'
import { getMyPlatformQuotas } from '@/api/user'
import AppLayout from '@/components/layout/AppLayout.vue'
import AnnouncementBell from '@/components/common/AnnouncementBell.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateLocalInput, formatDateTime } from '@/utils/format'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { IMAGE_WORKSHOP_MENU_ID, findImageWorkshopMenuItem } from '@/utils/imageWorkshop'
import { buildWorkbenchModelInsight } from '@/utils/apiKeyWorkbench'
import { setThemePreference, useThemePreference, useThemeState, type ThemePreference } from '@/utils/theme'
import type { ApiKey, ModelStat, PlatformQuotaItem, TrendDataPoint, UsageLog } from '@/types'

const authStore = useAuthStore()
const appStore = useAppStore()
const announcementStore = useAnnouncementStore()
const OPEN_ANNOUNCEMENT_CENTER_EVENT = 'sst-open-announcement-center'
const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const isAccountMenuOpen = ref(false)
const accountMenuRef = ref<HTMLElement | null>(null)
const requestsFocusRef = ref<HTMLElement | null>(null)
const user = computed(() => authStore.user)
const stats = ref<UserStatsType | null>(null)
const loading = ref(false)
const loadingUsage = ref(false)
const loadingCharts = ref(false)
const isRequestsFocusActive = ref(false)
const errorMessage = ref('')
const trendData = ref<TrendDataPoint[]>([])
const modelStats = ref<ModelStat[]>([])
const recentUsage = ref<UsageLog[]>([])
const platformQuotas = ref<PlatformQuotaItem[]>([])
const dashboardKeys = ref<ApiKey[]>([])
const dashboardWorkbenchStats = ref<Record<string, ApiKeyWorkbenchSummary>>({})
const themePreference = useThemePreference()
const isDark = useThemeState()
let requestsFocusTimer: ReturnType<typeof setTimeout> | null = null

type IconName = 'key' | 'chart' | 'user' | 'globe' | 'dollar' | 'clock' | 'database' | 'image'
type StatusTone = 'calm' | 'notice' | 'alert'
type LedgerEntry = { to: string; label: string; icon: IconName; value: string; note: string; alert: boolean }
type QuotaSummaryItem = { label: string; value: string; note: string }

const zhDashboardCopy = {
  loadingTitle: '正在整理今日账册',
  loadingCopy: '余额、调用、配额与近时段流转正在收束。',
  consoleAria: '山枢庭用户后台首页',
  stageAria: '用户值守台',
  watchStatus: '值守状态',
  gateAria: '常用入口',
  reasonsAria: '值守判断',
  nextStep: '下一步',
  stableReasonTitle: '今日值守稳定',
  stableReasonDetail: '密钥、余额与响应均未触发异常。',
  calmAction: '安稳',
  healthAria: '接入体检单',
  healthTitle: '接入体验',
  fullCheck: '完整体检',
  ledgerAria: '值守案牍',
  ledgerTitle: '值守案牍',
  waterLedgerAria: '水文账册',
  callsTitle: '调用',
  allUsage: '全部用量',
  requestFocus: '请求聚焦',
  rateLimitOnly: '只看限流错误',
  reviewAllLedger: '回看全部账册',
  doConnectionCheck: '做一次接入体检',
  loadingUsage: '正在归拢调用记录...',
  noRecentUsageTitle: '暂无最近调用',
  noRecentUsageCopy: '新的请求会在这里形成最近账册。可以先检查密钥、服务状态或充值与兑换。',
  manageKeys: '管理 API 密钥',
  checkService: '检查服务状态',
  viewLedger: '查看用量账册',
  fullLedger: '完整账册',
  errorLedger: '错误账册',
  openAnnouncements: '打开庭讯列表',
  announcements: '庭讯',
  viewAllAnnouncements: '点击查看全部庭讯',
  lastSevenDays: '近七日',
  tokenTrendAria: '近七日 token 趋势',
  window: '窗口',
  accountOverview: '账户概览',
  models: '模型',
  channelStatus: '通道状态',
  loadingModels: '正在校准模型流向...',
  noModelDistribution: '暂无模型分布',
  platforms: '平台',
  status: '状态',
  noPlatformRecords: '暂无平台记录',
  errorTitle: '暂未取到账册',
  errorCopy: '当前无法读取概览数据，请稍后刷新或检查服务连接。',
  retry: '重新整理',
  health: { calm: '安稳', notice: '留意', alert: '待处理' },
  seal: {
    normalLabel: '今日账册',
    alertLabel: '待处置',
    noticeLabel: '待核项',
    normalValue: '安',
    normalNote: '无碍',
    alertNote: '先处置',
    noticeNote: '待核对',
  },
  reasonActions: { alert: '处置', notice: '待核', calm: '前往' },
  statusTitles: {
    calm: '今日值守安稳',
    notice: '{count} 项需要核对',
    alert: '{count} 项需要先处置',
    summaryCalm: '账户、密钥与响应保持稳定，今日账册无明显异动。',
  },
  statusReasons: {
    zeroBalance: ['账户余量为零', '充值或调整账户后再放量更稳。'],
    noKeys: ['尚未创建密钥', '先建立 API Key，再接入调用。'],
    noActiveKeys: ['暂无启用密钥', '已有密钥，但当前没有可用入口。'],
    slowAlert: ['响应偏慢', '平均响应 {duration}，建议检查通道状态。'],
    slowNotice: ['响应需留意', '平均响应 {duration}，可观察近时段波动。'],
    noTraffic: ['今日未起流', '若已接入，请检查密钥、服务状态与调用记录。'],
  },
  advices: {
    runwayCritical: ['余量即将见底', '按近期速度约可用 {days} 天，建议先充值。'],
    runwayWeek: ['余额续航不足一周', '按当前速度预计还能运行 {days}，可提前充值或降速。'],
    fewKeys: ['密钥启用偏少', '建议清理停用 Key，保留生产入口更清晰。'],
    slow: ['响应需要观察', '平均响应 {duration}，可查看服务状态。'],
    noTraffic: ['今日尚未起流', '可用接入体检确认 Key 与地址是否可用。'],
    model: '先处理模型入口',
    calm: ['今日无碍', '密钥、用量和响应均保持可用，可继续观察账册。'],
  },
  labels: {
    requestCount: '{count} 请求',
    noTrafficToday: '今日未起流',
    enabled: '启用',
    accountSettings: '账户设置',
    imageWorkshop: '图像工坊',
    imageWorkshopNote: '外接创作入口',
    serviceStatus: '服务状态',
    apiKeys: 'API 密钥',
    usageLedger: '用量账册',
    recharge: '充值与兑换',
    key: '密钥',
    noEnabled: '暂无启用',
    todayRequests: '今日请求',
    totalPrefix: '累计 ',
    todayCost: '今日消耗',
    standardPrefix: '标准 $',
    response: '响应',
    ledger: '账册',
    accountBalance: '账户余量',
    ledgerSummaryAlert: '{count} 项待核',
    ledgerSummaryCalm: '今日账页',
    none: '暂无',
    noAnnouncements: '暂无庭讯',
    noAnnouncementsNote: '庭讯通道已接入，等待后台发出庭讯。',
    noAnnouncementsBadge: '待发庭讯',
    read: '已读',
    unread: '未读',
    unreadCount: '{count} 条未读',
    synced: '已同步',
    simpleRunwayValue: '按当前站点模式运行',
    simpleRunwayNote: '当前不以余额续航作为主要约束。',
    zeroRunwayValue: '已见底',
    zeroRunwayNote: '当前余额已经无法继续覆盖后续消耗。',
    noSampleRunwayValue: '暂无样本',
    noSampleRunwayNote: '近期缺少稳定消耗样本，续航会在产生真实用量后更新。',
    aboutDays: '约 {days}',
    callCredential: '调用凭证',
    traffic: '起流情况',
    modelAvailability: '模型可用性',
    stability: '请求稳定性',
    quota: '余额与额度',
    pending: '待处理',
    confirm: '待确认',
    notice: '留意',
    calm: '安稳',
    platformWindow: '平台窗口',
    open: '开放',
    noPlatformLimit: '暂无平台配额限制',
    windowUnused: '窗口未使用',
    maxUsage: '最高用量 $',
    todayToken: '今日 Token',
    inputOutput: '输入 {input} / 输出 {output}',
    totalCost: '累计消耗',
    balanceRunway: '余额续航',
  },
  checks: {
    noKeys: '当前还没有 API Key，外部请求还无法进入山枢庭。',
    createKey: '去创建密钥',
    noActiveKeys: '已有密钥但没有启用入口，建议先恢复至少一个生产 Key。',
    enableKey: '去启用密钥',
    partialKeys: '当前仅有 {active}/{total} 个密钥处于启用状态，建议确认生产入口是否清晰。',
    checkKeys: '检查密钥状态',
    keysReady: '密钥入口已就绪，当前启用中的 Key 可以直接承接请求。',
    doCheck: '做一次接入体检',
    noTraffic: '首页还没有近时段调用记录。若你刚完成接入，建议立即发起一次真实请求确认起流。',
    startCheck: '开始接入体检',
    trafficNoSuccess: '已有最近调用，但暂时没有看到明确成功账单，建议回看错误账册确认是否连续失败。',
    viewErrors: '查看错误账册',
    trafficReady: '近期已经有调用进入账册，可继续观察模型分布与费用变化。',
    viewRecent: '查看最近用量',
    slowAlert: '平均响应已超过 {duration}，更像持续拥塞而不是偶发波动。',
    checkChannels: '检查通道状态',
    slowNotice: '近期响应偏慢或刚经历限流，建议结合最近请求与错误账册判断是否已经恢复。',
    reviewRecent: '回看最近请求',
    stable: '最近响应速度和请求节奏都较平稳，可继续正常放量。',
    viewStatus: '查看服务状态',
    quotaZero: '账户余额已经为 0，请先充值或调整额度后再继续放量。',
    recharge: '去充值与兑换',
    quotaRunway: '按当前消耗速度预计还能运行 {days}，建议提前充值或先降速。',
    runwayEstimate: '去看续航估算',
    quotaWindow: '{platform} 窗口已使用约 {percent}%，建议提前留出缓冲。',
    viewWindow: '查看账户窗口',
    quotaReady: '当前账户余额和平台窗口尚可继续承接请求，没有明显逼近阈值。',
    viewQuota: '查看额度摘要',
  },
  focus: {
    noUsageTitle: '先确认请求是否已真正到达山枢庭',
    slowTitle: '先看近时段请求与响应节奏',
    normalTitle: '这里收拢了最近一段时间的调用账册',
    noUsageDetail: '当前首页还没有近时段调用记录。若你刚处理过限流或失败请求，先做一次接入体检，再回到错误账册确认是否仍在发生。',
    slowDetail: '若刚刚遇到限流或重试，这里可以先回看最近 7 笔调用，再结合错误账册判断是瞬时拥塞还是持续异常。',
    normalDetail: '若错误详情刚引导你回到首页，这一段会优先展示最近请求，方便对照失败时间、请求密度与后续恢复情况。',
  },
  footerUsage: '今日共 {total} 次请求，近列成功入账 {success} 条。',
  platformUsage: '{requests} 请求 · {tokens}',
  days: '{days} 天',
} as const

const enDashboardCopy = {
  loadingTitle: 'Preparing today\'s ledger',
  loadingCopy: 'Balance, calls, quotas, and recent flow are being gathered.',
  consoleAria: 'SST user console home',
  stageAria: 'User watch desk',
  watchStatus: 'Watch status',
  gateAria: 'Common entries',
  reasonsAria: 'Watch reasoning',
  nextStep: 'Next step',
  stableReasonTitle: 'Today is stable',
  stableReasonDetail: 'Keys, balance, and response status have not triggered issues.',
  calmAction: 'Stable',
  healthAria: 'Integration check sheet',
  healthTitle: 'Integration check',
  fullCheck: 'Full check',
  ledgerAria: 'Watch ledger slips',
  ledgerTitle: 'Watch ledger',
  waterLedgerAria: 'Usage ledger',
  callsTitle: 'Calls',
  allUsage: 'All usage',
  requestFocus: 'Request focus',
  rateLimitOnly: 'Rate-limit errors only',
  reviewAllLedger: 'Review full ledger',
  doConnectionCheck: 'Run an integration check',
  loadingUsage: 'Collecting call records...',
  noRecentUsageTitle: 'No recent calls',
  noRecentUsageCopy: 'New requests will form the recent ledger here. You can check Keys, service status, or recharge and redeem first.',
  manageKeys: 'Manage API Keys',
  checkService: 'Check service status',
  viewLedger: 'View usage ledger',
  fullLedger: 'Full ledger',
  errorLedger: 'Error ledger',
  openAnnouncements: 'Open announcements',
  announcements: 'Notices',
  viewAllAnnouncements: 'View all notices',
  lastSevenDays: 'Last 7 days',
  tokenTrendAria: '7-day token trend',
  window: 'Window',
  accountOverview: 'Account overview',
  models: 'Models',
  channelStatus: 'Channel status',
  loadingModels: 'Calibrating model flow...',
  noModelDistribution: 'No model distribution',
  platforms: 'Platforms',
  status: 'Status',
  noPlatformRecords: 'No platform records',
  errorTitle: 'Ledger unavailable',
  errorCopy: 'Overview data cannot be read right now. Refresh later or check the service connection.',
  retry: 'Refresh',
  health: { calm: 'Stable', notice: 'Notice', alert: 'Needs action' },
  seal: {
    normalLabel: 'Today ledger',
    alertLabel: 'Action',
    noticeLabel: 'Review',
    normalValue: 'OK',
    normalNote: 'Clear',
    alertNote: 'Act first',
    noticeNote: 'Review',
  },
  reasonActions: { alert: 'Fix', notice: 'Review', calm: 'Open' },
  statusTitles: {
    calm: 'Today is stable',
    notice: '{count} items need review',
    alert: '{count} items need action',
    summaryCalm: 'Account, Keys, and responses are stable; no obvious ledger anomalies today.',
  },
  statusReasons: {
    zeroBalance: ['Account balance is zero', 'Recharge or adjust the account before increasing traffic.'],
    noKeys: ['No Keys created', 'Create an API Key before sending calls.'],
    noActiveKeys: ['No active Keys', 'Keys exist, but no entrance is currently usable.'],
    slowAlert: ['Responses are slow', 'Average response is {duration}; check channel status.'],
    slowNotice: ['Response needs attention', 'Average response is {duration}; watch recent fluctuations.'],
    noTraffic: ['No traffic today', 'If already integrated, check Keys, service status, and call records.'],
  },
  advices: {
    runwayCritical: ['Balance nearly depleted', 'At recent speed, about {days} days remain. Recharge first.'],
    runwayWeek: ['Runway under one week', 'At current speed, about {days} remain. Recharge early or slow down.'],
    fewKeys: ['Few active Keys', 'Clean up inactive Keys and keep the production entrance clear.'],
    slow: ['Response needs observation', 'Average response is {duration}; check service status.'],
    noTraffic: ['No traffic yet today', 'Use the integration check to confirm the Key and address.'],
    model: 'Handle the model entrance first',
    calm: ['Clear for today', 'Keys, usage, and responses are available. Continue watching the ledger.'],
  },
  labels: {
    requestCount: '{count} requests',
    noTrafficToday: 'No traffic today',
    enabled: 'enabled',
    accountSettings: 'Account settings',
    imageWorkshop: 'Image Studio',
    imageWorkshopNote: 'External creation entrance',
    serviceStatus: 'Service status',
    apiKeys: 'API Keys',
    usageLedger: 'Usage ledger',
    recharge: 'Recharge and redeem',
    key: 'Keys',
    noEnabled: 'None active',
    todayRequests: 'Today requests',
    totalPrefix: 'Total ',
    todayCost: 'Today cost',
    standardPrefix: 'Standard $',
    response: 'Response',
    ledger: 'Ledger',
    accountBalance: 'Account balance',
    ledgerSummaryAlert: '{count} to review',
    ledgerSummaryCalm: 'Today ledger',
    none: 'None',
    noAnnouncements: 'No notices',
    noAnnouncementsNote: 'Notice channel is connected; waiting for admin notices.',
    noAnnouncementsBadge: 'Waiting',
    read: 'Read',
    unread: 'Unread',
    unreadCount: '{count} unread',
    synced: 'Synced',
    simpleRunwayValue: 'Current site mode',
    simpleRunwayNote: 'Balance runway is not the main constraint in this mode.',
    zeroRunwayValue: 'Depleted',
    zeroRunwayNote: 'Current balance can no longer cover future consumption.',
    noSampleRunwayValue: 'No sample',
    noSampleRunwayNote: 'Recent stable usage samples are missing; runway updates after real usage appears.',
    aboutDays: 'about {days}',
    callCredential: 'Call credentials',
    traffic: 'Traffic',
    modelAvailability: 'Model availability',
    stability: 'Request stability',
    quota: 'Balance and quota',
    pending: 'Needs action',
    confirm: 'To confirm',
    notice: 'Notice',
    calm: 'Stable',
    platformWindow: 'Platform window',
    open: 'Open',
    noPlatformLimit: 'No platform quota limit',
    windowUnused: 'Window unused',
    maxUsage: 'Max usage $',
    todayToken: 'Today tokens',
    inputOutput: 'Input {input} / output {output}',
    totalCost: 'Total cost',
    balanceRunway: 'Balance runway',
  },
  checks: {
    noKeys: 'There is no API Key yet, so external requests cannot enter SST.',
    createKey: 'Create a Key',
    noActiveKeys: 'Keys exist but none are active. Restore at least one production Key first.',
    enableKey: 'Enable a Key',
    partialKeys: 'Only {active}/{total} Keys are active. Confirm the production entrance is clear.',
    checkKeys: 'Check Key status',
    keysReady: 'Key entrance is ready; active Keys can receive requests.',
    doCheck: 'Run integration check',
    noTraffic: 'The home page has no recent call records. If you just integrated, send a real request to confirm traffic.',
    startCheck: 'Start integration check',
    trafficNoSuccess: 'Recent calls exist, but no clear successful ledger entry is visible. Review error records for repeated failures.',
    viewErrors: 'View error ledger',
    trafficReady: 'Recent calls have entered the ledger. Continue watching model distribution and cost changes.',
    viewRecent: 'View recent usage',
    slowAlert: 'Average response exceeds {duration}; this looks more like sustained congestion than a temporary fluctuation.',
    checkChannels: 'Check channel status',
    slowNotice: 'Recent responses are slow or just experienced rate limits. Compare recent requests with error records to confirm recovery.',
    reviewRecent: 'Review recent requests',
    stable: 'Recent response speed and request rhythm are steady. Normal traffic can continue.',
    viewStatus: 'View service status',
    quotaZero: 'Account balance is 0. Recharge or adjust quota before continuing traffic.',
    recharge: 'Recharge and redeem',
    quotaRunway: 'At current consumption speed, about {days} remain. Recharge early or slow down.',
    runwayEstimate: 'View runway estimate',
    quotaWindow: '{platform} window is about {percent}% used. Leave buffer in advance.',
    viewWindow: 'View account windows',
    quotaReady: 'Current balance and platform windows can still handle requests; no threshold pressure is visible.',
    viewQuota: 'View quota summary',
  },
  focus: {
    noUsageTitle: 'Confirm whether requests have reached SST',
    slowTitle: 'Check recent request and response rhythm first',
    normalTitle: 'Recent call ledger is gathered here',
    noUsageDetail: 'There are no recent call records on the home page. If you just handled rate limits or failed requests, run an integration check first, then return to the error ledger.',
    slowDetail: 'If you just hit rate limits or retries, review the latest 7 calls here and compare them with the error ledger to determine whether congestion is temporary or persistent.',
    normalDetail: 'If error details brought you back home, this section highlights recent requests so you can compare failure time, request density, and later recovery.',
  },
  footerUsage: '{total} requests today; {success} recent entries posted successfully.',
  platformUsage: '{requests} requests · {tokens}',
  days: '{days} days',
} as const

const dashboardCopy = computed(() => locale.value.startsWith('zh') ? zhDashboardCopy : enDashboardCopy)
const localeCode = computed(() => locale.value.startsWith('zh') ? 'zh-CN' : 'en-US')
const isEnglishDashboard = computed(() => !locale.value.startsWith('zh'))
const interpolate = (template: string, values: Record<string, string | number>) =>
  template.replace(/\{(\w+)\}/g, (_, key) => String(values[key] ?? ''))

const themeOptions = computed<Array<{ value: ThemePreference, label: string }>>(() => [
  { value: 'system', label: t('userShell.theme.system') },
  { value: 'light', label: t('userShell.theme.light') },
  { value: 'dark', label: t('userShell.theme.dark') },
])

const currentDateLabel = computed(() => new Intl.DateTimeFormat(localeCode.value, {
  timeZone: 'Asia/Shanghai',
  year: 'numeric',
  month: '2-digit',
  day: '2-digit'
}).format(new Date()).replace(/\//g, '-'))

const activeKeyRate = computed(() => {
  const total = stats.value?.total_api_keys || 0
  if (!total) return 0
  return Math.round(((stats.value?.active_api_keys || 0) / total) * 100)
})

const healthScore = computed(() => {
  if (!stats.value) return 0
  let score = 88
  if ((stats.value.total_api_keys || 0) > 0 && activeKeyRate.value === 0) score -= 30
  if (!authStore.isSimpleMode && (user.value?.balance ?? 0) <= 0) score -= 24
  if ((stats.value.average_duration_ms || 0) > 6000) score -= 20
  else if ((stats.value.average_duration_ms || 0) > 3000) score -= 10
  if ((stats.value.today_requests || 0) > 0) score += 4
  return Math.max(0, Math.min(100, score))
})

const healthLabel = computed(() => {
  if (healthScore.value >= 82) return dashboardCopy.value.health.calm
  if (healthScore.value >= 60) return dashboardCopy.value.health.notice
  return dashboardCopy.value.health.alert
})

const healthTone = computed<StatusTone>(() => {
  if (healthScore.value >= 82) return 'calm'
  if (healthScore.value >= 60) return 'notice'
  return 'alert'
})

const statusIssueCount = computed(() => statusReasons.value.length)

const watchSealLabel = computed(() => {
  if (!statusIssueCount.value) return dashboardCopy.value.seal.normalLabel
  return healthTone.value === 'alert' ? dashboardCopy.value.seal.alertLabel : dashboardCopy.value.seal.noticeLabel
})

const watchSealValue = computed(() => {
  if (!statusIssueCount.value) return dashboardCopy.value.seal.normalValue
  return String(statusIssueCount.value).padStart(2, '0')
})

const watchSealNote = computed(() => {
  if (!statusIssueCount.value) return dashboardCopy.value.seal.normalNote
  return healthTone.value === 'alert' ? dashboardCopy.value.seal.alertNote : dashboardCopy.value.seal.noticeNote
})

const statusReasons = computed<Array<{ label: string, detail: string, to: string, icon: IconName, tone: StatusTone }>>(() => {
  if (!stats.value) return []
  const reasons: Array<{ label: string, detail: string, to: string, icon: IconName, tone: StatusTone }> = []
  const totalKeys = stats.value.total_api_keys || 0
  const activeKeys = stats.value.active_api_keys || 0
  const balance = user.value?.balance ?? 0
  const duration = stats.value.average_duration_ms || 0
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const copy = dashboardCopy.value

  if (!authStore.isSimpleMode && paymentEnabled && balance <= 0) {
    reasons.push({ label: copy.statusReasons.zeroBalance[0], detail: copy.statusReasons.zeroBalance[1], to: '/purchase', icon: 'database', tone: 'alert' })
  }
  if (!totalKeys) {
    reasons.push({ label: copy.statusReasons.noKeys[0], detail: copy.statusReasons.noKeys[1], to: '/keys', icon: 'key', tone: 'notice' })
  } else if (!activeKeys) {
    reasons.push({ label: copy.statusReasons.noActiveKeys[0], detail: copy.statusReasons.noActiveKeys[1], to: '/keys', icon: 'key', tone: 'alert' })
  }
  if (duration > 6000) {
    reasons.push({ label: copy.statusReasons.slowAlert[0], detail: interpolate(copy.statusReasons.slowAlert[1], { duration: formatDuration(duration) }), to: '/monitor', icon: 'clock', tone: 'alert' })
  } else if (duration > 3000) {
    reasons.push({ label: copy.statusReasons.slowNotice[0], detail: interpolate(copy.statusReasons.slowNotice[1], { duration: formatDuration(duration) }), to: '/usage', icon: 'clock', tone: 'notice' })
  }
  if ((stats.value.today_requests || 0) === 0) {
    reasons.push({ label: copy.statusReasons.noTraffic[0], detail: copy.statusReasons.noTraffic[1], to: '/usage', icon: 'chart', tone: 'notice' })
  }

  return reasons.slice(0, 3)
})

const statusTitle = computed(() => {
  const copy = dashboardCopy.value.statusTitles
  if (healthTone.value === 'calm') return copy.calm
  if (healthTone.value === 'notice') return interpolate(copy.notice, { count: statusIssueCount.value })
  return interpolate(copy.alert, { count: statusIssueCount.value })
})

const statusSummary = computed(() => {
  if (!statusReasons.value.length) return dashboardCopy.value.statusTitles.summaryCalm
  return statusReasons.value.map((reason) => reason.label).join('、')
})

function reasonActionLabel(tone: StatusTone) {
  return dashboardCopy.value.reasonActions[tone]
}

const dashboardModelInsight = computed(() => {
  const entries = dashboardKeys.value
    .filter((key) => key.status === 'active')
    .map((key) => {
      const summary = dashboardWorkbenchStats.value[key.id] || dashboardWorkbenchStats.value[String(key.id)]
      return {
        key,
        latestError: summary?.latest_error || null,
      }
    })
    .filter((item) => item.latestError)
    .sort((a, b) => new Date(b.latestError!.created_at).getTime() - new Date(a.latestError!.created_at).getTime())

  for (const entry of entries) {
    const insight = buildWorkbenchModelInsight(entry.latestError, entry.key.name, locale.value)
    if (insight) return insight
  }

  return null
})

const watchAdvices = computed<Array<{ title: string, detail: string, to: string, tone: StatusTone }>>(() => {
  if (!stats.value) return []
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const advices: Array<{ title: string, detail: string, to: string, tone: StatusTone }> = []
  const balance = user.value?.balance ?? 0
  const averageDailyCost = Math.max(stats.value.today_actual_cost || 0, (stats.value.total_actual_cost || 0) / 14)
  const balanceDays = averageDailyCost > 0 ? Math.floor(balance / averageDailyCost) : null
  const copy = dashboardCopy.value

  if (!authStore.isSimpleMode && paymentEnabled && balanceDays !== null && balanceDays <= 3) {
    advices.push({ title: copy.advices.runwayCritical[0], detail: interpolate(copy.advices.runwayCritical[1], { days: Math.max(balanceDays, 0) }), to: '/purchase', tone: 'alert' })
  } else if (!authStore.isSimpleMode && paymentEnabled && balanceRunway.value.days !== null && balanceRunway.value.days <= 7) {
    advices.push({ title: copy.advices.runwayWeek[0], detail: interpolate(copy.advices.runwayWeek[1], { days: formatRunwayDays(balanceRunway.value.days) }), to: '/purchase', tone: 'notice' })
  }
  if ((stats.value.total_api_keys || 0) > 0 && activeKeyRate.value < 50) {
    advices.push({ title: copy.advices.fewKeys[0], detail: copy.advices.fewKeys[1], to: '/keys', tone: 'notice' })
  }
  if ((stats.value.average_duration_ms || 0) > 3000) {
    advices.push({ title: copy.advices.slow[0], detail: interpolate(copy.advices.slow[1], { duration: formatDuration(stats.value.average_duration_ms || 0) }), to: '/monitor', tone: healthTone.value })
  }
  if ((stats.value.today_requests || 0) === 0 && (stats.value.active_api_keys || 0) > 0) {
    advices.push({ title: copy.advices.noTraffic[0], detail: copy.advices.noTraffic[1], to: '/keys?panel=connection-test', tone: 'notice' })
  }
  if (dashboardModelInsight.value) {
    advices.unshift({
      title: copy.advices.model,
      detail: dashboardModelInsight.value.detail,
      to: dashboardModelInsight.value.to,
      tone: dashboardModelInsight.value.tone,
    })
  }
  if (!advices.length) {
    advices.push({ title: copy.advices.calm[0], detail: copy.advices.calm[1], to: '/usage', tone: 'calm' })
  }

  return advices.slice(0, 3)
})

const todayRequestLabel = computed(() => {
  const count = stats.value?.today_requests || 0
  return count ? interpolate(dashboardCopy.value.labels.requestCount, { count: formatNumber(count) }) : dashboardCopy.value.labels.noTrafficToday
})

const imageWorkshopMenuItem = computed(() => findImageWorkshopMenuItem(appStore.cachedPublicSettings?.custom_menu_items))
const imageWorkshopPath = computed(() => `/custom/${IMAGE_WORKSHOP_MENU_ID}`)

const courtyardGates = computed<Array<{ to: string, label: string, note: string, icon: IconName, position: string }>>(() => {
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const copy = dashboardCopy.value.labels
  const gates: Array<{ to: string, label: string, note: string, icon: IconName, position: string }> = [
    { to: '/keys', label: copy.apiKeys, note: `${stats.value?.active_api_keys || 0}/${stats.value?.total_api_keys || 0} ${copy.enabled}`, icon: 'key', position: 'gate-north' },
    { to: '/usage', label: copy.usageLedger, note: todayRequestLabel.value, icon: 'chart', position: 'gate-east' },
    { to: !authStore.isSimpleMode && paymentEnabled ? '/purchase' : '/usage', label: !authStore.isSimpleMode && paymentEnabled ? copy.recharge : copy.usageLedger, note: authStore.isSimpleMode || !paymentEnabled ? copy.accountSettings : '$' + formatMoney(user.value?.balance || 0), icon: 'user', position: 'gate-south' }
  ]
  if (imageWorkshopMenuItem.value) {
    gates.splice(1, 0, { to: imageWorkshopPath.value, label: copy.imageWorkshop, note: copy.imageWorkshopNote, icon: 'image', position: 'gate-workshop' })
  }
  if (!authStore.isSimpleMode) gates.splice(2, 0, { to: '/monitor', label: copy.serviceStatus, note: healthLabel.value, icon: 'globe', position: 'gate-west' })
  return gates
})

const gardenEntries = computed(() => {
  const copy = dashboardCopy.value.labels
  const entries: LedgerEntry[] = [
    { to: '/keys', label: copy.key, icon: 'key' as IconName, value: (stats.value?.active_api_keys || 0) + '/' + (stats.value?.total_api_keys || 0), note: activeKeyRate.value ? activeKeyRate.value + `% ${copy.enabled}` : copy.noEnabled, alert: (stats.value?.total_api_keys || 0) === 0 || (stats.value?.active_api_keys || 0) === 0 },
    { to: '/usage', label: copy.todayRequests, icon: 'chart' as IconName, value: formatNumber(stats.value?.today_requests || 0), note: copy.totalPrefix + formatNumber(stats.value?.total_requests || 0), alert: false },
    { to: '/usage', label: copy.todayCost, icon: 'dollar' as IconName, value: '$' + formatCost(stats.value?.today_actual_cost || 0), note: copy.standardPrefix + formatCost(stats.value?.today_cost || 0), alert: false },
    { to: '/usage', label: copy.response, icon: 'clock' as IconName, value: formatDuration(stats.value?.average_duration_ms || 0), note: formatTokens(stats.value?.rpm || 0) + ' RPM', alert: (stats.value?.average_duration_ms || 0) > 3000 }
  ]

  if (!authStore.isSimpleMode) {
    const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
    entries.unshift({ to: paymentEnabled ? '/purchase' : '/usage', label: copy.ledger, icon: 'database' as IconName, value: '$' + formatMoney(user.value?.balance || 0), note: copy.accountBalance, alert: (user.value?.balance ?? 0) <= 0 })
  }

  return entries
})

const ledgerSummary = computed(() => {
  const alertCount = gardenEntries.value.filter((item) => item.alert).length
  if (alertCount > 0) return interpolate(dashboardCopy.value.labels.ledgerSummaryAlert, { count: alertCount })
  return dashboardCopy.value.labels.ledgerSummaryCalm
})

const modelPreview = computed(() => modelStats.value.slice(0, 5))
const platformPreview = computed(() => [...(stats.value?.by_platform ?? [])]
  .sort((a, b) => b.total_actual_cost - a.total_actual_cost)
  .slice(0, 4))
const maxModelTokens = computed(() => Math.max(1, ...modelPreview.value.map((item) => item.total_tokens || 0)))
const modelShare = (tokens: number) => `${Math.max(8, Math.round((tokens / maxModelTokens.value) * 100))}%`

const trendBars = computed(() => {
  const weekdayFormatter = new Intl.DateTimeFormat(localeCode.value, { weekday: 'short' })
  const recentDays = Array.from({ length: 7 }, (_, index) => {
    const date = new Date(Date.now() - (6 - index) * 86400000)
    return formatDateLocalInput(date)
  })
  const pointMap = new Map(trendData.value.map((point) => [point.date, point.total_tokens || 0]))
  const values = recentDays.map((date) => pointMap.get(date) ?? 0)
  const max = Math.max(1, ...values)

  return recentDays.map((date) => {
    const tokens = pointMap.get(date) ?? 0
    const weekday = weekdayFormatter.format(new Date(`${date}T00:00:00`))
    return {
      date,
      label: formatTokens(tokens),
      weekday,
      total_tokens: tokens,
      height: `${Math.max(10, Math.round((tokens / max) * 100))}%`,
      isPeak: tokens === max && tokens > 0,
      empty: tokens === 0,
    }
  })
})

const trendPeakLabel = computed(() => {
  if (!trendBars.value.length) return dashboardCopy.value.labels.none
  const peak = trendBars.value.reduce((max, point) => (point.isPeak ? point : max), trendBars.value[0])
  if (!peak || !peak.label) return dashboardCopy.value.labels.none
  return `${peak.date} / ${peak.label}`
})

const requestsFocusEnabled = computed(() => route.query.focus === 'requests')

const announcementSummary = computed(() => {
  const items = announcementStore.announcements
  const unreadCount = items.filter((item) => !item.read_at).length
  const latest = items[0]
  const copy = dashboardCopy.value.labels
  if (!latest) {
    return {
      unreadCount: 0,
      title: copy.noAnnouncements,
      note: copy.noAnnouncementsNote,
      badge: copy.noAnnouncementsBadge,
    }
  }
  return {
    unreadCount,
    title: latest.title,
    note: latest.read_at ? `${copy.read} · ${formatDateTime(latest.created_at)}` : `${copy.unread} · ${formatDateTime(latest.created_at)}`,
    badge: unreadCount > 0 ? interpolate(copy.unreadCount, { count: unreadCount }) : copy.synced,
  }
})

const requestsFocusTitle = computed(() => {
  if (!recentUsage.value.length) return dashboardCopy.value.focus.noUsageTitle
  if ((stats.value?.average_duration_ms || 0) > 3000) return dashboardCopy.value.focus.slowTitle
  return dashboardCopy.value.focus.normalTitle
})

const requestsFocusDetail = computed(() => {
  if (!recentUsage.value.length) {
    return dashboardCopy.value.focus.noUsageDetail
  }

  if ((stats.value?.average_duration_ms || 0) > 3000) {
    return dashboardCopy.value.focus.slowDetail
  }

  return dashboardCopy.value.focus.normalDetail
})

const recentSuccessCount = computed(() => recentUsage.value.filter((item) => (item.actual_cost || 0) > 0).length)
const recentUsageFooter = computed(() => interpolate(dashboardCopy.value.footerUsage, {
  total: formatNumber(stats.value?.today_requests || 0),
  success: recentSuccessCount.value,
}))
const platformUsageLabel = (requests: number, tokens: number) => interpolate(dashboardCopy.value.platformUsage, {
  requests: formatNumber(requests),
  tokens: formatTokens(tokens),
})
const hasRecentTraffic = computed(() => recentUsage.value.length > 0 || (stats.value?.today_requests || 0) > 0)
const balanceRunway = computed(() => {
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const balance = user.value?.balance ?? 0
  const averageDailyCost = Math.max(stats.value?.today_actual_cost || 0, (stats.value?.total_actual_cost || 0) / 14)

  if (authStore.isSimpleMode || !paymentEnabled) {
    return {
      enabled: false,
      days: null as number | null,
      value: dashboardCopy.value.labels.simpleRunwayValue,
      note: '',
      tone: 'calm' as StatusTone,
    }
  }

  if (balance <= 0) {
    return {
      enabled: true,
      days: 0,
      value: dashboardCopy.value.labels.zeroRunwayValue,
      note: dashboardCopy.value.labels.zeroRunwayNote,
      tone: 'alert' as StatusTone,
    }
  }

  if (averageDailyCost <= 0) {
    return {
      enabled: true,
      days: null,
      value: dashboardCopy.value.labels.noSampleRunwayValue,
      note: dashboardCopy.value.labels.noSampleRunwayNote,
      tone: 'notice' as StatusTone,
    }
  }

  const days = balance / averageDailyCost
  if (days <= 3) {
    return {
      enabled: true,
      days,
      value: interpolate(dashboardCopy.value.labels.aboutDays, { days: formatRunwayDays(days) }),
      note: '',
      tone: 'alert' as StatusTone,
    }
  }

  if (days <= 7) {
    return {
      enabled: true,
      days,
      value: interpolate(dashboardCopy.value.labels.aboutDays, { days: formatRunwayDays(days) }),
      note: '',
      tone: 'notice' as StatusTone,
    }
  }

  return {
    enabled: true,
    days,
    value: interpolate(dashboardCopy.value.labels.aboutDays, { days: formatRunwayDays(days) }),
    note: '',
    tone: 'calm' as StatusTone,
  }
})
const strongestQuotaUsage = computed(() => {
  if (!platformQuotas.value.length) return null
  return platformQuotas.value.reduce((max, quota) => {
    const percent = Math.max(
      quota.daily_limit_usd ? (quota.daily_usage_usd / quota.daily_limit_usd) * 100 : 0,
      quota.weekly_limit_usd ? (quota.weekly_usage_usd / quota.weekly_limit_usd) * 100 : 0,
      quota.monthly_limit_usd ? (quota.monthly_usage_usd / quota.monthly_limit_usd) * 100 : 0,
    )
    if (!max || percent > max.percent) {
      return { quota, percent }
    }
    return max
  }, null as null | { quota: PlatformQuotaItem; percent: number })
})

const connectionCheckItems = computed<Array<{ key: string, label: string, status: string, detail: string, action: string, to: string, tone: StatusTone }>>(() => {
  const items: Array<{ key: string, label: string, status: string, detail: string, action: string, to: string, tone: StatusTone }> = []
  const totalKeys = stats.value?.total_api_keys || 0
  const activeKeys = stats.value?.active_api_keys || 0
  const balance = user.value?.balance ?? 0
  const paymentEnabled = isFeatureFlagEnabled(FeatureFlags.payment)
  const topQuota = strongestQuotaUsage.value
  const labels = dashboardCopy.value.labels
  const checks = dashboardCopy.value.checks

  if (!totalKeys) {
    items.push({ key: 'keys', label: labels.callCredential, status: labels.pending, detail: checks.noKeys, action: checks.createKey, to: '/keys', tone: 'alert' })
  } else if (!activeKeys) {
    items.push({ key: 'keys', label: labels.callCredential, status: labels.pending, detail: checks.noActiveKeys, action: checks.enableKey, to: '/keys', tone: 'alert' })
  } else if (activeKeys < totalKeys) {
    items.push({ key: 'keys', label: labels.callCredential, status: labels.notice, detail: interpolate(checks.partialKeys, { active: activeKeys, total: totalKeys }), action: checks.checkKeys, to: '/keys', tone: 'notice' })
  } else {
    items.push({ key: 'keys', label: labels.callCredential, status: labels.calm, detail: checks.keysReady, action: checks.doCheck, to: '/keys?panel=connection-test', tone: 'calm' })
  }

  if (!hasRecentTraffic.value) {
    items.push({ key: 'traffic', label: labels.traffic, status: labels.confirm, detail: checks.noTraffic, action: checks.startCheck, to: '/keys?panel=connection-test', tone: 'notice' })
  } else if (!recentSuccessCount.value && recentUsage.value.length) {
    items.push({ key: 'traffic', label: labels.traffic, status: labels.notice, detail: checks.trafficNoSuccess, action: checks.viewErrors, to: '/usage?tab=errors', tone: 'notice' })
  } else {
    items.push({ key: 'traffic', label: labels.traffic, status: labels.calm, detail: checks.trafficReady, action: checks.viewRecent, to: '/usage', tone: 'calm' })
  }

  if (dashboardModelInsight.value) {
    items.push({
      key: 'models',
      label: labels.modelAvailability,
      status: dashboardModelInsight.value.status,
      detail: dashboardModelInsight.value.detail,
      action: dashboardModelInsight.value.action,
      to: dashboardModelInsight.value.to,
      tone: dashboardModelInsight.value.tone,
    })
  }

  if ((stats.value?.average_duration_ms || 0) > 6000) {
    items.push({ key: 'stability', label: labels.stability, status: labels.pending, detail: interpolate(checks.slowAlert, { duration: formatDuration(stats.value?.average_duration_ms || 0) }), action: checks.checkChannels, to: '/monitor', tone: 'alert' })
  } else if ((stats.value?.average_duration_ms || 0) > 3000 || requestsFocusEnabled.value) {
    items.push({ key: 'stability', label: labels.stability, status: labels.notice, detail: checks.slowNotice, action: checks.reviewRecent, to: '/dashboard?focus=requests', tone: 'notice' })
  } else {
    items.push({ key: 'stability', label: labels.stability, status: labels.calm, detail: checks.stable, action: checks.viewStatus, to: '/monitor', tone: 'calm' })
  }

  if (!authStore.isSimpleMode && paymentEnabled && balance <= 0) {
    items.push({ key: 'quota', label: labels.quota, status: labels.pending, detail: checks.quotaZero, action: checks.recharge, to: '/purchase', tone: 'alert' })
  } else if (!authStore.isSimpleMode && paymentEnabled && balanceRunway.value.days !== null && balanceRunway.value.days <= 7) {
    items.push({ key: 'quota', label: labels.quota, status: balanceRunway.value.days <= 3 ? labels.pending : labels.notice, detail: interpolate(checks.quotaRunway, { days: formatRunwayDays(balanceRunway.value.days) }), action: checks.runwayEstimate, to: '/profile?focus=balance-notify', tone: balanceRunway.value.days <= 3 ? 'alert' : 'notice' })
  } else if (topQuota && topQuota.percent >= 85) {
    items.push({ key: 'quota', label: labels.quota, status: labels.notice, detail: interpolate(checks.quotaWindow, { platform: platformLabel(topQuota.quota.platform), percent: Math.round(topQuota.percent) }), action: checks.viewWindow, to: '/profile', tone: 'notice' })
  } else {
    items.push({ key: 'quota', label: labels.quota, status: labels.calm, detail: checks.quotaReady, action: checks.viewQuota, to: '/profile', tone: 'calm' })
  }

  return items
})

const quotaSummary = computed(() => {
  const labels = dashboardCopy.value.labels
  const base: QuotaSummaryItem[] = [
    {
      label: labels.todayToken,
      value: formatTokens(stats.value?.today_tokens || 0),
      note: interpolate(labels.inputOutput, { input: formatTokens(stats.value?.today_input_tokens || 0), output: formatTokens(stats.value?.today_output_tokens || 0) })
    },
    {
      label: labels.totalCost,
      value: '$' + formatCost(stats.value?.total_actual_cost || 0),
      note: labels.standardPrefix + formatCost(stats.value?.total_cost || 0)
    },
    {
      label: labels.balanceRunway,
      value: balanceRunway.value.value,
      note: balanceRunway.value.note
    }
  ]

  const quota = platformQuotas.value[0]
  if (quota) {
    base.push({
      label: platformLabel(quota.platform),
      value: quotaLimitLabel(quota),
      note: quotaUsageLabel(quota)
    })
  } else {
    base.push({ label: labels.platformWindow, value: labels.open, note: labels.noPlatformLimit })
  }

  return base
})

const PLATFORM_LABELS: Record<string, string> = {
  anthropic: 'Claude',
  openai: 'OpenAI',
  gemini: 'Gemini',
  antigravity: 'Antigravity'
}

const platformLabel = (platform: string) => PLATFORM_LABELS[platform] ?? platform

const quotaLimitLabel = (quota: PlatformQuotaItem) => {
  const limits = [quota.daily_limit_usd, quota.weekly_limit_usd, quota.monthly_limit_usd].filter((limit): limit is number => limit != null)
  if (!limits.length) return dashboardCopy.value.labels.open
  return '$' + formatCost(Math.max(...limits))
}

const quotaUsageLabel = (quota: PlatformQuotaItem) => {
  const usages = [quota.daily_usage_usd, quota.weekly_usage_usd, quota.monthly_usage_usd].filter((usage): usage is number => usage != null)
  if (!usages.length) return dashboardCopy.value.labels.windowUnused
  return dashboardCopy.value.labels.maxUsage + formatCost(Math.max(...usages))
}
const startDate = ref(formatDateLocalInput(new Date(Date.now() - 6 * 86400000)))
const endDate = ref(formatDateLocalInput(new Date()))
const granularity = ref<'day' | 'hour'>('day')


const closeAccountMenu = () => {
  isAccountMenuOpen.value = false
}

const openAccountMenu = () => {
  isAccountMenuOpen.value = true
}

const toggleAccountMenu = () => {
  if (isAccountMenuOpen.value) closeAccountMenu()
  else openAccountMenu()
}


const handleAccountPointerDown = (event: PointerEvent) => {
  if (!isAccountMenuOpen.value) return
  const target = event.target as Node | null
  if (target && accountMenuRef.value?.contains(target)) return
  closeAccountMenu()
}

const handleAccountKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') closeAccountMenu()
}

const chooseTheme = (preference: ThemePreference) => {
  setThemePreference(preference)
}

const openAnnouncementCenter = () => {
  if (typeof window === 'undefined') return
  window.dispatchEvent(new Event(OPEN_ANNOUNCEMENT_CENTER_EVENT))
}

watch(() => route.fullPath, closeAccountMenu)

const clearRequestsFocusTimer = () => {
  if (requestsFocusTimer) {
    clearTimeout(requestsFocusTimer)
    requestsFocusTimer = null
  }
}

const activateRequestsFocus = async () => {
  if (!requestsFocusEnabled.value || !stats.value) return
  await nextTick()
  requestsFocusRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  isRequestsFocusActive.value = false
  requestAnimationFrame(() => {
    isRequestsFocusActive.value = true
  })
  clearRequestsFocusTimer()
  requestsFocusTimer = setTimeout(() => {
    isRequestsFocusActive.value = false
    requestsFocusTimer = null
  }, 2600)
}

watch(
  [requestsFocusEnabled, stats],
  ([enabled, dashboardStats]) => {
    if (!enabled || !dashboardStats) {
      clearRequestsFocusTimer()
      isRequestsFocusActive.value = false
      return
    }
    activateRequestsFocus()
  },
  { immediate: true }
)

const handleLogout = async () => {
  closeAccountMenu()
  await authStore.logout()
  router.push('/login')
}

const loadStats = async () => {
  loading.value = true
  errorMessage.value = dashboardCopy.value.errorCopy
  try {
    await authStore.refreshUser()
    stats.value = await usageAPI.getDashboardStats()
  } catch (error) {
    console.error('Failed to load dashboard stats:', error)
    errorMessage.value = dashboardCopy.value.errorCopy
  } finally {
    loading.value = false
  }
}

const loadCharts = async () => {
  loadingCharts.value = true
  try {
    const [trendResult, modelsResult] = await Promise.allSettled([
      usageAPI.getDashboardTrend({ start_date: startDate.value, end_date: endDate.value, granularity: granularity.value }),
      usageAPI.getDashboardModels({ start_date: startDate.value, end_date: endDate.value })
    ])

    if (trendResult.status === 'fulfilled') {
      trendData.value = trendResult.value.trend || []
    } else {
      console.error('Failed to load dashboard trend:', trendResult.reason)
    }

    if (modelsResult.status === 'fulfilled') {
      modelStats.value = modelsResult.value.models || []
    } else {
      console.error('Failed to load dashboard models:', modelsResult.reason)
    }
  } catch (error) {
    console.error('Failed to load charts:', error)
  } finally {
    loadingCharts.value = false
  }
}

const loadRecent = async () => {
  loadingUsage.value = true
  try {
    const res = await usageAPI.getByDateRange(startDate.value, endDate.value)
    recentUsage.value = res.items.slice(0, 7)
  } catch (error) {
    console.error('Failed to load recent usage:', error)
  } finally {
    loadingUsage.value = false
  }
}

const loadPlatformQuotas = async () => {
  try {
    const data = await getMyPlatformQuotas()
    platformQuotas.value = data.platform_quotas ?? []
  } catch (error) {
    console.warn('Failed to load platform quotas:', error)
    platformQuotas.value = []
  }
}

const loadDashboardKeyWorkbench = async () => {
  try {
    const response = await keysAPI.list(1, 100, {
      sort_by: 'created_at',
      sort_order: 'desc',
    })
    dashboardKeys.value = response.items || []
    const ids = dashboardKeys.value.map((key) => key.id)
    if (!ids.length) {
      dashboardWorkbenchStats.value = {}
      return
    }
    const workbench = await usageAPI.getDashboardApiKeysWorkbench(ids)
    dashboardWorkbenchStats.value = workbench.stats || {}
  } catch (error) {
    console.warn('Failed to load dashboard key workbench:', error)
    dashboardKeys.value = []
    dashboardWorkbenchStats.value = {}
  }
}

const refreshAll = () => {
  loadStats()
  loadCharts()
  loadRecent()
  loadPlatformQuotas()
  loadDashboardKeyWorkbench()
}

const formatNumber = (n: number) => n.toLocaleString(localeCode.value)
const formatMoney = (n: number) => n.toFixed(2)
const formatCost = (n: number) => n.toFixed(4)
const formatTokens = (t: number) => {
  if (t >= 1_000_000) return `${(t / 1_000_000).toFixed(1)}M`
  if (t >= 1000) return `${(t / 1000).toFixed(1)}K`
  return t.toString()
}
const formatDuration = (ms: number) => ms >= 1000 ? `${(ms / 1000).toFixed(2)}s` : `${ms.toFixed(0)}ms`
const formatRunwayDays = (days: number) => interpolate(dashboardCopy.value.days, {
  days: days >= 10 ? Math.round(days) : days.toFixed(1),
})

onMounted(() => {
  document.addEventListener('pointerdown', handleAccountPointerDown)
  document.addEventListener('keydown', handleAccountKeydown)
  refreshAll()
})

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', handleAccountPointerDown)
  document.removeEventListener('keydown', handleAccountKeydown)
  clearRequestsFocusTimer()
})
</script>

<style scoped>
.sst-dashboard {
  --sst-paper: #faf7ef;
  --sst-paper-soft: rgba(250, 247, 239, 0.78);
  --sst-paper-deep: rgba(244, 239, 228, 0.92);
  --sst-ink: #1f2320;
  --sst-ink-soft: #38413a;
  --sst-mute: #59645a;
  --sst-line: rgba(216, 205, 185, 0.78);
  --sst-line-soft: rgba(216, 205, 185, 0.46);
  --sst-seal: #a73a2a;
  --sst-brass: #9b8155;
  --sst-wash: rgba(99, 111, 94, 0.1);
  min-height: calc(100vh - 2rem);
  padding: clamp(0.75rem, 1.8vw, 1.6rem) clamp(0.75rem, 1.8vw, 1.6rem) clamp(1.35rem, 2.4vw, 2.2rem);
  background:
    linear-gradient(90deg, rgba(119, 104, 78, 0.08) 1px, transparent 1px) 0 0 / 4rem 4rem,
    linear-gradient(180deg, rgba(244, 239, 228, 0.28), rgba(237, 229, 212, 0.52)),
    #f4efe4;
  color: var(--sst-ink);
}

.home-state {
  display: grid;
  min-height: 26rem;
  place-items: center;
  gap: 1rem;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 14px;
  text-align: center;
  transition: background-color 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.home-state.is-paper {
  border-color: rgba(198, 184, 157, 0.46);
  background:
    radial-gradient(circle at 50% 22%, rgba(167, 58, 42, 0.07), transparent 20%),
    radial-gradient(circle at 50% 76%, rgba(155, 129, 85, 0.06), transparent 18%),
    linear-gradient(180deg, rgba(250, 246, 238, 0.98), rgba(242, 234, 219, 0.96));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.62), 0 22px 48px -38px rgba(91, 70, 41, 0.22);
}

.home-state.is-night {
  border-color: rgba(198, 184, 157, 0.16);
  background:
    radial-gradient(circle at 50% 24%, rgba(167, 58, 42, 0.08), transparent 22%),
    radial-gradient(circle at 50% 72%, rgba(155, 129, 85, 0.05), transparent 18%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.95), rgba(17, 19, 15, 0.92));
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.03), 0 22px 48px -38px rgba(0, 0, 0, 0.52);
}

.home-state-error {
  gap: 1.1rem;
  padding: 1.5rem 1.2rem;
}

.home-state .seal-mark {
  width: 2.45rem;
  height: 2.45rem;
  box-shadow: 0 14px 28px -24px rgba(0, 0, 0, 0.6);
}

.home-state.is-paper .seal-mark {
  background: linear-gradient(180deg, rgba(241, 232, 216, 0.96), rgba(228, 216, 197, 0.94));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.68), 0 14px 28px -24px rgba(91, 70, 41, 0.32);
}

.home-state.is-night .seal-mark {
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
}

.home-state.is-paper .spinner {
  color: #b36a31;
}

.home-state.is-night .spinner {
  color: #21b8b1;
}

.home-state h1 {
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.36rem;
  font-weight: 620;
}

.home-state.is-paper h1 {
  color: #2a241e;
}

.home-state.is-night h1 {
  color: #f4efe4;
}

.home-state p,
.mini-state {
  font-size: 0.84rem;
}

.home-state.is-paper p,
.home-state.is-paper .mini-state {
  color: #746858;
}

.home-state.is-night p,
.home-state.is-night .mini-state {
  color: #9b8f79;
}

.mini-state {
  display: grid;
  min-height: 6.8rem;
  place-items: center;
  border: 1px dashed rgba(198, 184, 157, 0.28);
  border-radius: 10px;
  background: rgba(24, 26, 21, 0.35);
  padding: 0.8rem 1rem;
  text-align: center;
}

.home-state p {
  max-width: 30rem;
  line-height: 1.6;
}

.seal-mark {
  display: inline-grid;
  width: 2.18rem;
  height: 2.18rem;
  place-items: center;
  border-radius: 0.42rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 10px 22px -22px rgba(31, 35, 32, 0.4);
}

.seal-mark img {
  display: block;
  width: 100%;
  height: 100%;
}

.courtyard-console {
  position: relative;
  max-width: 1360px;
  margin: 0 auto;
  overflow: hidden;
  border: 1px solid rgba(198, 184, 157, 0.86);
  border-radius: 14px;
  background:
    radial-gradient(circle at 12% 18%, rgba(167, 58, 42, 0.045), transparent 24%),
    linear-gradient(135deg, rgba(250, 247, 239, 0.97), rgba(244, 239, 228, 0.88)),
    url('@/assets/brand/sst-paper-ink-bg.png') center/cover;
  box-shadow: 0 26px 70px -52px rgba(31, 35, 32, 0.44);
}

.courtyard-console::before {
  display: none;
}

.courtyard-console::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 18%, transparent 82%, rgba(155, 129, 85, 0.08)),
    radial-gradient(circle at 78% 18%, rgba(167, 58, 42, 0.045), transparent 22%),
    linear-gradient(180deg, transparent, rgba(244, 239, 228, 0.22));
  pointer-events: none;
}

.console-head,
.courtyard-stage {
  padding: 0.58rem;
  padding-top: 0.2rem;
}

.health-check-sheet {
  padding: 0;
}

.health-check-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.58rem;
}

.check-card {
  display: grid;
  gap: 0.55rem;
  min-width: 0;
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 10px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.22), transparent 24%),
    rgba(250, 247, 239, 0.42);
  padding: 0.9rem;
}

.check-card-head {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 0.75rem;
  min-width: 0;
}

.check-card-head span {
  min-width: 0;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  overflow-wrap: anywhere;
}

.check-card-head strong {
  min-width: 0;
  color: var(--sst-ink);
  font-size: 0.82rem;
  font-weight: 650;
  overflow-wrap: anywhere;
  text-align: right;
}

.sst-dashboard-en .check-card-head span {
  flex: 1 1 auto;
}

.sst-dashboard-en .check-card-head strong {
  flex: 0 0 auto;
  font-size: 0.76rem;
  line-height: 1.1;
  white-space: nowrap;
}

.check-card p {
  color: var(--sst-mute);
  font-size: 0.78rem;
  line-height: 1.66;
  overflow-wrap: anywhere;
}

.check-card a {
  width: fit-content;
  max-width: 100%;
  border-bottom: 1px solid rgba(167, 58, 42, 0.22);
  color: var(--sst-seal);
  font-size: 0.76rem;
  font-weight: 650;
  padding-bottom: 0.14rem;
  overflow-wrap: anywhere;
}

.check-card.tone-calm {
  border-color: rgba(117, 138, 104, 0.32);
}

.check-card.tone-notice {
  border-color: rgba(155, 129, 85, 0.34);
  background: linear-gradient(90deg, rgba(155, 129, 85, 0.06), transparent 76%), rgba(250, 247, 239, 0.46);
}

.check-card.tone-alert {
  border-color: rgba(167, 58, 42, 0.28);
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 72%), rgba(250, 247, 239, 0.5);
}

.water-ledger {
  position: relative;
}

.console-head {
  z-index: 30;
}

.courtyard-stage {
  z-index: 10;
}

.water-ledger {
  z-index: 1;
}

.console-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1.5rem;
  min-height: 5.35rem;
  border-bottom: 1px solid rgba(198, 184, 157, 0.42);
  padding: clamp(0.85rem, 1.65vw, 1.25rem) clamp(1rem, 2.2vw, 1.75rem);
}

.brand-lockup {
  display: flex;
  min-width: 0;
  align-items: flex-start;
  gap: 0.82rem;
}

.brand-copy {
  display: grid;
  gap: 0.18rem;
  padding-top: 0.06rem;
}

.brand-copy > span {
  font-size: 0.62rem;
  letter-spacing: 0.22em;
  color: rgba(123, 106, 83, 0.86);
}

.brand-lockup span,
.account-mark span,
.slips-head span,
.center-seal span,
.section-mark span,
.ledger-slip span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.brand-lockup h1 {
  margin-top: 0;
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.38rem, 1.68vw, 1.66rem);
  font-weight: 620;
  letter-spacing: 0;
  line-height: 1.08;
  max-width: none;
  text-wrap: balance;
}

.account-mark {
  min-width: min(22rem, 42vw);
  border-right: 2px solid var(--sst-seal);
  padding-right: 0.9rem;
  text-align: right;
}

.account-mark strong {
  display: block;
  margin-top: 0.34rem;
  overflow-wrap: anywhere;
  color: var(--sst-ink-soft);
  font-size: 0.86rem;
}

.account-mark small {
  display: block;
  margin-top: 0.42rem;
  color: var(--sst-mute);
  font-size: 0.74rem;
}

.account-menu {
  position: relative;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.account-menu-trigger {
  display: inline-grid;
  width: 2.35rem;
  height: 2.35rem;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid rgba(167, 58, 42, 0.25);
  border-radius: 0.42rem;
  background: rgba(167, 58, 42, 0.08);
  color: var(--sst-seal);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-weight: 700;
  transition: background-color 160ms ease, border-color 160ms ease, transform 160ms ease;
}

.account-menu-trigger:hover,
.account-menu-trigger:focus-visible,
.account-menu.is-open .account-menu-trigger {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.14);
  transform: translateY(-1px);
  outline: none;
}

.account-menu-dropdown {
  position: absolute;
  top: calc(100% + 0.45rem);
  right: 0;
  z-index: 10000;
  display: grid;
  min-width: 11rem;
  gap: 0.12rem;
  padding: 0.35rem;
  border: 1px solid rgba(216, 205, 185, 0.86);
  border-radius: 0.75rem;
  background: rgba(250, 247, 239, 0.98);
  box-shadow: 0 24px 70px -42px rgba(31, 35, 32, 0.36);
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transform: translateY(-0.35rem);
  transition: opacity 160ms ease, transform 160ms ease;
}
.account-menu.is-open .account-menu-dropdown {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
  transform: translateY(0);
}

.account-menu-dropdown a,
.account-menu-dropdown button {
  border-radius: 0.44rem;
  padding: 0.62rem 0.72rem;
  color: var(--sst-ink-soft);
  font-size: 0.82rem;
  font-weight: 650;
  position: relative;
  z-index: 1;
  pointer-events: none;
  text-align: left;
  transition: background-color 160ms ease, color 160ms ease;
}
.account-menu.is-open .account-menu-dropdown a,
.account-menu.is-open .account-menu-dropdown button {
  pointer-events: auto;
}

.account-menu-dropdown a:hover,
.account-menu-dropdown button:hover,
.account-menu-dropdown a:focus-visible,
.account-menu-dropdown button:focus-visible {
  background: rgba(167, 58, 42, 0.08);
  color: var(--sst-seal);
  outline: none;
}

.account-menu-dropdown a {
  display: flex;
  align-items: center;
}

.account-menu-dropdown button {
  width: 100%;
}

.account-theme {
  display: grid;
  gap: 0.45rem;
  border-top: 1px solid rgba(198, 184, 157, 0.44);
  border-bottom: 1px solid rgba(198, 184, 157, 0.44);
  margin: 0.2rem 0;
  padding: 0.62rem 0.5rem;
}

.account-theme > span {
  padding: 0 0.22rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.14em;
}

.account-theme-options {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.25rem;
}

.account-menu-dropdown .account-theme-options button {
  justify-content: center;
  padding: 0.42rem 0.36rem;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 0.44rem;
  background: rgba(255, 252, 245, 0.54);
  color: #59645a;
  font-size: 0.74rem;
  font-weight: 650;
  text-align: center;
}

.account-menu-dropdown .account-theme-options button:hover,
.account-menu-dropdown .account-theme-options button:focus-visible,
.account-menu-dropdown .account-theme-options button.is-selected {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.09);
  color: var(--sst-seal);
}

.courtyard-stage {
  display: grid;
  grid-template-columns: minmax(34rem, 1fr) minmax(13.5rem, 0.36fr);
  align-items: stretch;
  gap: 0.62rem;
  padding: 0.58rem;
  border-bottom: 1px solid rgba(198, 184, 157, 0.28);
}

.courtyard-stage::after {
  display: none;
}

.courtyard-map,
.ledger-slips,
.usage-scroll,
.folio {
  position: relative;
}

.courtyard-map,
.ledger-slip,
.usage-scroll,
.folio {
  border: 1px solid rgba(198, 184, 157, 0.34);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.42);
}

.courtyard-map {
  position: relative;
  height: 100%;
  min-height: 0;
  overflow: hidden;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.052), transparent 9rem),
    radial-gradient(circle at 72% 20%, rgba(99, 111, 94, 0.08), transparent 22%),
    linear-gradient(180deg, rgba(252, 250, 245, 0.95), rgba(246, 241, 231, 0.76));
}

.watch-desk {
  display: grid;
  grid-template-columns: minmax(0, 1.08fr) minmax(18rem, 0.82fr);
  align-items: start;
  gap: 0.55rem 0.95rem;
  padding: clamp(0.72rem, 1.2vw, 0.9rem);
}

.courtyard-map::before,
.courtyard-map::after {
  display: none;
}

.ink-wash {
  position: absolute;
  pointer-events: none;
}

.ink-wash-map {
  inset: 0;
  opacity: 0.13;
  background:
    radial-gradient(circle at 18% 22%, rgba(167, 58, 42, 0.1), transparent 16%),
    radial-gradient(circle at 76% 62%, rgba(90, 93, 82, 0.08), transparent 18%),
    url('@/assets/brand/sst-paper-ink-bg.png') center/cover no-repeat;
  mix-blend-mode: multiply;
}

.ink-wash-side {
  display: none;
}

.watch-main,
.watch-status,
.watch-reasons,
.gate-grid {
  position: relative;
  z-index: 2;
}

.watch-main {
  display: grid;
  gap: 0.68rem;
  align-content: start;
  min-width: 0;
}

.watch-status {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-self: start;
  align-items: center;
  gap: 0.82rem;
  padding: 0.18rem 0.15rem 0;
}

.watch-copy {
  display: grid;
  gap: 0.32rem;
  max-width: 44rem;
}

.watch-copy span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
}

.watch-copy h2 {
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.18rem, 1.54vw, 1.52rem);
  font-weight: 600;
  line-height: 1.2;
}

.watch-copy p {
  color: var(--sst-mute);
  font-size: 0.8rem;
  line-height: 1.6;
}

.watch-reasons {
  grid-column: 2;
  display: grid;
  align-content: start;
  min-height: 0;
  border-left: 0;
  padding: 0.45rem 0 0.45rem 0.8rem;
  gap: 0;
}

.watch-advice {
  display: grid;
  gap: 0;
  margin-bottom: 0.35rem;
  border-bottom: 1px solid rgba(198, 184, 157, 0.34);
  padding-bottom: 0.35rem;
}

.watch-advice > span {
  margin-bottom: 0.28rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
}

.watch-advice a {
  display: grid;
  gap: 0.16rem;
  border-left: 2px solid rgba(155, 129, 85, 0.5);
  padding: 0.42rem 0 0.42rem 0.62rem;
  color: var(--sst-ink-soft);
}

.watch-advice a:hover,
.watch-advice a:focus-visible {
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.045), transparent 76%);
  color: var(--sst-seal);
  outline: none;
}

.watch-advice strong,
.watch-advice small {
  display: block;
  overflow-wrap: anywhere;
}

.watch-advice strong {
  color: var(--sst-ink);
  font-size: 0.82rem;
}

.watch-advice small {
  color: var(--sst-mute);
  font-size: 0.72rem;
  line-height: 1.42;
}

.watch-advice .tone-calm {
  border-left-color: #51624f;
}

.watch-advice .tone-notice {
  border-left-color: #9b8155;
}

.watch-advice .tone-alert {
  border-left-color: var(--sst-seal);
}

.watch-notice {
  display: grid;
  gap: 0.22rem;
  border-top: 1px solid rgba(198, 184, 157, 0.34);
  margin-top: 0.3rem;
  padding-top: 0.45rem;
}

.watch-notice > span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.watch-notice a {
  display: grid;
  gap: 0.16rem;
  border-left: 2px solid rgba(167, 58, 42, 0.48);
  padding-left: 0.62rem;
  color: var(--sst-ink-soft);
}

.watch-notice a:hover,
.watch-notice a:focus-visible {
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.045), transparent 76%);
  color: var(--sst-seal);
  outline: none;
}

.watch-notice strong,
.watch-notice small {
  display: block;
}

.watch-notice strong {
  color: var(--sst-ink);
  font-size: 0.82rem;
}

.watch-notice small {
  color: var(--sst-mute);
  font-size: 0.72rem;
  line-height: 1.42;
}

.reason-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 0.68rem;
  border-top: 1px solid rgba(198, 184, 157, 0.32);
  padding: 0.48rem 0;
  color: var(--sst-ink-soft);
}

.reason-item svg {
  color: var(--sst-seal);
}

.reason-item.tone-calm svg {
  color: #51624f;
}

.reason-item.tone-notice svg {
  color: var(--sst-brass);
}

.reason-item.tone-alert svg {
  color: var(--sst-seal);
}

.reason-item strong,
.reason-item small {
  display: block;
  min-width: 0;
  overflow-wrap: anywhere;
}

.reason-item strong {
  color: var(--sst-ink);
  font-size: 0.84rem;
}

.reason-item small {
  margin-top: 0.2rem;
  color: var(--sst-mute);
  font-size: 0.72rem;
  line-height: 1.42;
}

.reason-item em {
  border-left: 1px solid rgba(198, 184, 157, 0.58);
  padding-left: 0.72rem;
  color: var(--sst-seal);
  font-size: 0.74rem;
  font-style: normal;
  font-weight: 650;
  overflow-wrap: anywhere;
  text-align: right;
}

.tone-calm em,
.status-calm .watch-score strong {
  color: #51624f;
}

.tone-notice em,
.status-notice .watch-score strong {
  color: var(--sst-brass);
}

.tone-alert {
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.06), transparent 68%);
}

.tone-alert em,
.status-alert .watch-score strong {
  color: var(--sst-seal);
}

.gate-grid {
  align-self: start;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  grid-auto-flow: dense;
  gap: 0.38rem;
  padding-top: 0.42rem;
}

.watch-health-sheet {
  grid-column: 1 / -1;
  border-top: 1px solid rgba(198, 184, 157, 0.34);
  padding-top: 0.72rem;
}

.watch-health-sheet .section-mark {
  margin-bottom: 0.58rem;
}

.watch-health-sheet .health-check-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.52rem;
}

.watch-health-sheet .check-card {
  padding: 0.84rem 0.9rem;
}

.gate-link {
  position: relative;
  z-index: 2;
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-items: center;
  gap: 0.12rem 0.48rem;
  min-width: 0;
  border: 0;
  border-top: 1px solid rgba(198, 184, 157, 0.3);
  border-radius: 0;
  background: transparent;
  min-height: 3.25rem;
  padding: 0.48rem 0.2rem;
  color: var(--sst-ink-soft);
  font-size: 0.84rem;
  font-weight: 650;
  box-shadow: none;
  transition: border-color 180ms ease, background-color 180ms ease, color 180ms ease, transform 180ms ease, box-shadow 180ms ease;
}

.gate-link:hover {
  border-color: rgba(167, 58, 42, 0.34);
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.045), transparent 82%);
  color: var(--sst-seal);
  transform: translateX(2px);
  box-shadow: none;
}

.gate-link small {
  grid-column: 2;
  color: var(--sst-mute);
  font-size: 0.69rem;
  font-weight: 400;
  line-height: 1.34;
  overflow-wrap: anywhere;
}

.gate-link:focus-visible,
.ledger-slip:focus-visible,
.section-mark a:focus-visible {
  outline: 2px solid rgba(167, 58, 42, 0.5);
  outline-offset: 2px;
}

.gate-north {
  left: auto;
  top: auto;
  transform: none;
}

.gate-east {
  right: auto;
  top: auto;
  transform: none;
}

.gate-south {
  left: auto;
  bottom: auto;
  transform: none;
}

.gate-west {
  left: auto;
  top: auto;
  transform: none;
}

.gate-workshop {
  grid-column: 1 / -1;
  min-height: 3.5rem;
  border-color: rgba(167, 58, 42, 0.28);
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.055), rgba(255, 252, 245, 0.18) 58%, transparent),
    linear-gradient(180deg, rgba(255, 252, 246, 0.22), transparent);
}

.gate-workshop svg {
  color: var(--sst-seal);
}

.center-seal {
  position: relative;
  left: auto;
  top: auto;
  z-index: 2;
  display: grid;
  width: 5.6rem;
  height: 5.6rem;
  place-items: center;
  border: 1px solid rgba(167, 58, 42, 0.32);
  border-radius: 0.42rem;
  background:
    linear-gradient(135deg, rgba(167, 58, 42, 0.08), transparent 52%),
    rgba(250, 247, 239, 0.78);
  box-shadow: inset 0 0 0 1px rgba(167, 58, 42, 0.08), 0 20px 36px -30px rgba(31, 35, 32, 0.36);
  transform: none;
}

.center-seal strong {
  color: var(--sst-seal);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 1.76rem;
  font-weight: 660;
  line-height: 1;
  letter-spacing: 0;
  font-variant-numeric: tabular-nums;
}

.center-seal small {
  color: #7b6a53;
  font-size: 0.82rem;
}

.ledger-slips {
  display: flex;
  flex-direction: column;
  height: 100%;
  align-self: stretch;
  border: 1px solid rgba(198, 184, 157, 0.32);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.026), transparent 22%),
    rgba(250, 247, 239, 0.36);
  padding: 0.7rem 0.86rem;
}

.ledger-slips::after {
  content: '';
  display: block;
  height: 0.4rem;
}

.slips-head {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.8rem;
  padding: 0.08rem 0.15rem 0.5rem;
}

.slips-head small {
  color: var(--sst-mute);
  font-size: 0.7rem;
  white-space: nowrap;
}

.ledger-slip {
  display: grid;
  min-width: 0;
  grid-template-columns: auto minmax(0, 1fr);
  column-gap: 0.65rem;
  border-width: 0;
  border-top: 1px solid rgba(198, 184, 157, 0.3);
  border-radius: 0;
  background: transparent;
  flex: 1 1 0;
  align-items: start;
  padding: 0.58rem 0;
  color: var(--sst-ink-soft);
  box-shadow: none;
  transition: border-color 180ms ease, background-color 180ms ease, color 180ms ease, transform 180ms ease, box-shadow 180ms ease;
}

.ledger-slip:first-of-type {
  border-top: 0;
}

.ledger-slip:hover {
  border-color: rgba(167, 58, 42, 0.28);
  background: transparent;
  color: var(--sst-seal);
  transform: translateX(2px);
  box-shadow: none;
}

.ledger-slip.is-alert {
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.052), transparent 72%);
}

.ledger-slip svg {
  margin-top: 0.12rem;
  color: var(--sst-seal);
}

.ledger-slip-copy {
  display: grid;
  align-content: space-between;
  gap: 0.28rem;
  min-width: 0;
  min-height: 100%;
}

.ledger-slip-main {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.8rem;
  min-width: 0;
}

.ledger-slip strong {
  flex: 0 0 auto;
  color: var(--sst-ink);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 1.02rem;
  line-height: 1.1;
  overflow-wrap: anywhere;
  font-variant-numeric: tabular-nums;
  text-align: right;
}

.ledger-slip small {
  color: var(--sst-mute);
  font-size: 0.74rem;
  line-height: 1.45;
}

.water-ledger {
  display: grid;
  grid-template-columns: minmax(28rem, 1.12fr) minmax(28rem, 0.88fr);
  align-items: start;
  gap: 0.58rem;
  padding: 0.58rem;
  padding-top: 0.52rem;
}

.usage-scroll,
.folio {
  min-width: 0;
  padding: 0.9rem;
  box-shadow: none;
}

.usage-scroll {
  display: flex;
  flex-direction: column;
  min-height: 13.2rem;
}

.sst-dashboard-en .usage-scroll {
  align-self: start;
  min-height: 0;
}

.sst-dashboard-en .call-list-footer {
  margin-top: 1rem;
}

.focus-note {
  display: grid;
  gap: 0.72rem;
  margin-bottom: 0.9rem;
  border: 1px solid rgba(167, 58, 42, 0.18);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 72%),
    rgba(250, 247, 239, 0.78);
  padding: 0.78rem 0.9rem;
}

.focus-note span {
  display: block;
  margin-bottom: 0.18rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.focus-note strong {
  display: block;
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 600;
}

.focus-note p {
  margin-top: 0.34rem;
  color: var(--sst-mute);
  font-size: 0.78rem;
  line-height: 1.66;
}

.focus-note-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.45rem;
}

.focus-note-actions a {
  border: 1px solid rgba(198, 184, 157, 0.48);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.54);
  padding: 0.38rem 0.72rem;
  color: var(--sst-ink-soft);
  font-size: 0.75rem;
  font-weight: 650;
  transition: border-color 180ms ease, background-color 180ms ease, color 180ms ease, transform 180ms ease;
}

.focus-note-actions a:hover,
.focus-note-actions a:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.08);
  color: var(--sst-seal);
  outline: none;
  transform: translateX(1px);
}

.focus-requests {
  animation: request-focus-glow 2.4s ease;
}

.focus-surface {
  border-color: rgba(167, 58, 42, 0.42);
  box-shadow: 0 18px 42px -36px rgba(167, 58, 42, 0.32);
}

@keyframes request-focus-glow {
  0% {
    transform: translateY(0.65rem);
    opacity: 0.7;
  }
  35% {
    transform: translateY(0);
    opacity: 1;
  }
  100% {
    transform: translateY(0);
    opacity: 1;
  }
}

.ledger-side {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.62rem;
}

.notice-strip {
  grid-column: 1 / -1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  width: 100%;
  padding: 0.05rem 0 0.35rem;
  margin-bottom: 0.08rem;
  border-bottom: 1px solid rgba(198, 184, 157, 0.2);
  text-align: left;
  transition: border-color 180ms ease, opacity 180ms ease, transform 180ms ease;
}

.notice-strip:hover,
.notice-strip:focus-visible {
  border-bottom-color: rgba(167, 58, 42, 0.32);
}

.notice-strip:focus-visible {
  outline: none;
  opacity: 0.96;
}

.notice-strip-copy {
  min-width: 0;
}

.notice-strip-copy span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.notice-strip-copy strong {
  display: block;
  margin-top: 0.22rem;
  color: var(--sst-ink);
  font-size: 0.88rem;
  font-weight: 620;
  overflow-wrap: anywhere;
  text-wrap: pretty;
}

.notice-strip-copy small {
  display: block;
  margin-top: 0.16rem;
  color: var(--sst-mute);
  font-size: 0.7rem;
  line-height: 1.5;
  overflow-wrap: anywhere;
}

.notice-strip-meta {
  flex: 0 0 auto;
  text-align: right;
}

.notice-strip-meta em {
  display: block;
  color: var(--sst-seal);
  font-size: 0.72rem;
  font-style: normal;
  font-weight: 650;
  white-space: nowrap;
}

.notice-strip-meta small {
  display: block;
  margin-top: 0.14rem;
  color: var(--sst-mute);
  font-size: 0.66rem;
  line-height: 1.4;
  overflow-wrap: anywhere;
}

.folio {
  min-height: 11.8rem;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.2), transparent 26%),
    linear-gradient(90deg, rgba(99, 111, 94, 0.035), transparent 42%),
    rgba(250, 247, 239, 0.42);
}

.flow-folio {
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
}

.section-mark {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 0.64rem;
}

.section-mark strong {
  color: #a73a2a;
  font-size: 0.76rem;
  font-weight: 620;
  text-align: right;
}

.section-mark a,
.section-mark span:last-child {
  color: var(--sst-mute);
  font-size: 0.74rem;
  font-weight: 600;
  text-align: right;
}

.section-mark strong {
  font-size: 0.76rem;
  color: #d96a55;
}

.waterline {
  display: flex;
  align-items: end;
  height: auto;
  min-height: 9.2rem;
  gap: 0.34rem;
  padding: 0.24rem 0 0;
}

.waterline span {
  flex: 1;
  min-width: 0.5rem;
  border-radius: 0.34rem 0.34rem 0 0;
  background:
    linear-gradient(180deg, rgba(167, 58, 42, 0.96), rgba(155, 129, 85, 0.9)),
    linear-gradient(180deg, var(--sst-seal), #8f7d5f);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.08);
  opacity: 0.88;
  transform-origin: bottom center;
  transition: transform 180ms ease, opacity 180ms ease;
}

.waterline span:hover {
  transform: translateY(-2px);
  opacity: 1;
}

.waterline span.is-peak {
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12), 0 0 0 1px rgba(167, 58, 42, 0.1);
}

.waterline span.is-empty {
  opacity: 0.28;
  background: linear-gradient(180deg, rgba(155, 129, 85, 0.28), rgba(155, 129, 85, 0.18));
}

.waterline-axis {
  display: flex;
  gap: 0.34rem;
  padding-top: 0.18rem;
}

.waterline-axis span {
  flex: 1;
  min-width: 0.5rem;
  color: var(--sst-mute);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.62rem;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
}

.call-list,
.model-river,
.quota-list,
.platform-list {
  display: grid;
  gap: 0.62rem;
}

.call-list li,
.platform-list > div {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  list-style: none;
  border-top: 1px solid rgba(198, 184, 157, 0.26);
  padding-top: 0.62rem;
}

.quota-list > div {
  display: grid;
  grid-template-columns: minmax(7.5rem, 0.72fr) minmax(0, 1fr);
  align-items: baseline;
  gap: 0.28rem 0.9rem;
  min-width: 0;
  list-style: none;
  border-top: 1px solid rgba(198, 184, 157, 0.26);
  padding-top: 0.62rem;
}

.quota-list > div > span,
.quota-list > div > strong,
.quota-list > div > small {
  min-width: 0;
  overflow-wrap: anywhere;
}

.quota-list > div > strong {
  justify-self: end;
  text-align: right;
}

.quota-list > div > small {
  grid-column: 1 / -1;
  line-height: 1.5;
}

.platform-list > div:first-child {
  border-top: 0;
  padding-top: 0;
}

.call-list {
  counter-reset: usage-item;
}

.call-list li {
  position: relative;
  padding-left: 1.45rem;
}

.call-list li::before {
  content: counter(usage-item, decimal-leading-zero);
  counter-increment: usage-item;
  position: absolute;
  left: 0;
  top: 0.74rem;
  color: #8f7d5f;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.7rem;
  letter-spacing: 0.08em;
}

.call-list li > div:first-child {
  min-width: 0;
}

.call-list li > div:first-child strong,
.model-river span {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.call-list li > div:last-child,
.platform-list strong {
  flex: 0 0 auto;
  text-align: right;
}

.call-list strong,
.quota-list strong,
.model-river strong,
.platform-list strong {
  color: var(--sst-ink);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.84rem;
  font-variant-numeric: tabular-nums;
}

.call-list span,
.quota-list span,
.quota-list small,
.model-river span,
.platform-list span,
.platform-list small {
  display: block;
  margin-top: 0.16rem;
  color: var(--sst-mute);
  font-size: 0.72rem;
}

.call-list-footer {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 1rem;
  border-top: 1px solid rgba(198, 184, 157, 0.26);
  padding-top: 1rem;
  color: var(--sst-mute);
}

.call-list-footer span {
  display: block;
  color: #8f7d5f;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.62rem;
  letter-spacing: 0.18em;
}

.call-list-footer strong {
  display: block;
  margin-top: 0.28rem;
  color: var(--sst-ink);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.9rem;
  font-weight: 600;
}

.call-list-footer small {
  display: block;
  margin-top: 0.2rem;
  font-size: 0.72rem;
  line-height: 1.5;
}

.call-list-footer-actions {
  display: flex;
  flex: 0 0 auto;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.5rem;
}

.call-list-footer-actions a {
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 999px;
  padding: 0.42rem 0.72rem;
  color: #8f3024;
  font-size: 0.72rem;
  font-weight: 650;
}

.call-list-footer-actions a:hover,
.call-list-footer-actions a:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  outline: none;
}

.quota-list > div:first-child {
  border-top: 0;
  padding-top: 0;
}

.quota-list > div.quota-inline-row {
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: baseline;
  gap: 0.28rem 0.72rem;
}

.quota-list > div.quota-inline-row span,
.quota-list > div.quota-inline-row strong,
.quota-list > div.quota-inline-row small {
  display: block;
  grid-row: 1;
  margin-top: 0;
}

.quota-list > div.quota-inline-row span,
.quota-list > div.quota-inline-row strong {
  white-space: normal;
}

.quota-list > div.quota-inline-row small {
  grid-column: 2;
  min-width: 0;
  justify-self: stretch;
  text-align: left;
  white-space: nowrap;
}

.quota-list > div.quota-inline-row strong {
  grid-column: 3;
  justify-self: end;
}

.quota-list > div:nth-child(3) {
  gap: 0.28rem 0.9rem;
}

.quota-list > div:nth-child(3) strong {
  font-size: 0.82rem;
  line-height: 1.15;
  white-space: nowrap;
}

.quota-list > div:nth-child(3) small {
  max-width: none;
}

.quota-list > div:nth-child(3) span,
.quota-list > div:nth-child(3) strong {
  display: block;
}

.quota-list > div:nth-child(3) {
  align-items: baseline;
  grid-template-columns: minmax(5.8rem, 0.5fr) minmax(0, 1fr);
}

.quota-list > div:nth-child(3) strong {
  margin-left: 0;
}

.sst-dashboard-en .quota-list > div {
  grid-template-columns: minmax(8.3rem, 0.78fr) minmax(0, 1fr);
}

.sst-dashboard-en .quota-list > div > strong {
  font-size: 0.82rem;
  line-height: 1.15;
  white-space: nowrap;
}

.sst-dashboard-en .quota-list > div:nth-child(3) strong {
  font-size: 0.88rem;
}

.model-river > div {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 0.8rem;
  padding-bottom: 0.42rem;
}

.model-river i {
  position: absolute;
  left: 0;
  bottom: 0;
  height: 3px;
  border-radius: 999px;
  background: var(--sst-seal);
}

.empty-note {
  display: grid;
  min-height: 8.6rem;
  align-content: end;
  justify-items: start;
  color: var(--sst-mute);
  padding-bottom: 0.25rem;
}

.empty-note strong {
  color: var(--sst-ink-soft);
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.08rem;
  font-weight: 580;
}

.empty-note span {
  margin-top: 0.4rem;
  font-size: 0.84rem;
}

.empty-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 0.82rem;
}

.empty-actions a {
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 0.35rem;
  background: rgba(250, 247, 239, 0.34);
  padding: 0.38rem 0.54rem;
  color: var(--sst-ink-soft);
  font-size: 0.78rem;
  font-weight: 650;
  transition: border-color 180ms ease, color 180ms ease, transform 180ms ease;
}

.empty-actions a:hover {
  border-color: rgba(167, 58, 42, 0.36);
  background: rgba(167, 58, 42, 0.045);
  color: var(--sst-seal);
  transform: translateX(1px);
}

.home-state-error .btn-primary {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.14);
  color: #f4efe4;
}

.home-state-error .btn-primary:hover,
.home-state-error .btn-primary:focus-visible {
  border-color: rgba(167, 58, 42, 0.56);
  background: rgba(167, 58, 42, 0.22);
}

.dark .sst-dashboard {
  --sst-paper: rgba(24, 26, 21, 0.9);
  --sst-paper-soft: rgba(24, 26, 21, 0.82);
  --sst-paper-deep: rgba(17, 19, 15, 0.95);
  --sst-ink: #f4efe4;
  --sst-ink-soft: #d9d0be;
  --sst-mute: #879186;
  --sst-line: rgba(48, 52, 43, 0.95);
  --sst-line-soft: rgba(48, 52, 43, 0.72);
  background:
    radial-gradient(circle at 12% 0%, rgba(167, 58, 42, 0.08), transparent 18%),
    radial-gradient(circle at 88% 10%, rgba(155, 129, 85, 0.08), transparent 16%),
    linear-gradient(180deg, #10120e 0%, #151812 48%, #10120e 100%);
  color: #f4efe4;
}

.dark .courtyard-console,
.dark .courtyard-map,
.dark .ledger-slips,
.dark .ledger-slip,
.dark .usage-scroll,
.dark .folio,
.dark .center-seal,
.dark .gate-link,
.dark .empty-actions a {
  border-color: rgba(48, 52, 43, 0.95);
  background-color: rgba(24, 26, 21, 0.82);
}

.dark .courtyard-console {
  background:
    radial-gradient(circle at 8% 12%, rgba(167, 58, 42, 0.08), transparent 20%),
    radial-gradient(circle at 86% 12%, rgba(155, 129, 85, 0.08), transparent 18%),
    linear-gradient(180deg, rgba(18, 20, 16, 0.98), rgba(24, 26, 21, 0.92)),
    url('@/assets/brand/sst-paper-ink-bg.png') center/cover;
  box-shadow: 0 28px 80px -58px rgba(0, 0, 0, 0.7);
}

.dark .courtyard-console::after {
  background:
    radial-gradient(circle at top left, rgba(255, 255, 255, 0.02), transparent 24%),
    radial-gradient(circle at 78% 18%, rgba(167, 58, 42, 0.1), transparent 22%),
    linear-gradient(180deg, transparent, rgba(0, 0, 0, 0.24));
}

.dark .courtyard-map {
  background:
    radial-gradient(circle at 16% 18%, rgba(167, 58, 42, 0.1), transparent 18%),
    radial-gradient(circle at 76% 22%, rgba(155, 129, 85, 0.08), transparent 20%),
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 10rem),
    linear-gradient(180deg, rgba(24, 26, 21, 0.96), rgba(17, 19, 15, 0.88));
}

.dark .ledger-slips,
.dark .usage-scroll,
.dark .folio {
  background:
    radial-gradient(circle at 10% 10%, rgba(167, 58, 42, 0.055), transparent 18%),
    linear-gradient(180deg, rgba(255, 247, 235, 0.022), transparent 30%),
    linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 44%),
    rgba(17, 19, 15, 0.8);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.03);
}

.dark .ledger-slip,
.dark .gate-link,
.dark .empty-actions a {
  background: transparent;
}

.dark .center-seal {
  background:
    radial-gradient(circle at 50% 30%, rgba(167, 58, 42, 0.15), transparent 56%),
    linear-gradient(135deg, rgba(167, 58, 42, 0.12), transparent 56%),
    rgba(17, 19, 15, 0.84);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.04);
}

.dark .ledger-slips .ink-wash-side {
  opacity: 0.025;
  mix-blend-mode: normal;
}

.dark .check-card {
  border-color: rgba(68, 71, 58, 0.9);
  background:
    radial-gradient(circle at 16% 12%, rgba(167, 58, 42, 0.06), transparent 16%),
    linear-gradient(180deg, rgba(255, 247, 235, 0.038), transparent 28%),
    rgba(17, 19, 15, 0.76);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.03);
}

.dark .check-card.tone-calm {
  border-color: rgba(81, 98, 79, 0.42);
  background:
    radial-gradient(circle at 16% 12%, rgba(81, 98, 79, 0.12), transparent 20%),
    linear-gradient(90deg, rgba(81, 98, 79, 0.11), transparent 76%),
    rgba(17, 19, 15, 0.72);
}

.dark .check-card.tone-notice {
  border-color: rgba(155, 129, 85, 0.42);
  background:
    radial-gradient(circle at 16% 12%, rgba(155, 129, 85, 0.14), transparent 20%),
    linear-gradient(90deg, rgba(155, 129, 85, 0.12), transparent 76%),
    rgba(17, 19, 15, 0.72);
}

.dark .check-card.tone-alert {
  border-color: rgba(167, 58, 42, 0.42);
  background:
    radial-gradient(circle at 16% 12%, rgba(167, 58, 42, 0.16), transparent 20%),
    linear-gradient(90deg, rgba(167, 58, 42, 0.15), transparent 74%),
    rgba(17, 19, 15, 0.72);
}

.dark .ink-wash-map,
.dark .ink-wash-side {
  opacity: 0.06;
  mix-blend-mode: screen;
  filter: grayscale(0.2);
}

.dark .section-mark {
  border-bottom-color: rgba(93, 88, 73, 0.68);
}

.dark .brand-lockup span,
.dark .account-mark span,
.dark .slips-head span,
.dark .center-seal span,
.dark .section-mark span,
.dark .ledger-slip span,
.dark .watch-copy span,
.dark .watch-advice > span,
.dark .account-theme > span {
  color: #9b8f79;
}

.dark .call-list li,
.dark .call-list-footer,
.dark .quota-list > div,
.dark .platform-list > div,
.dark .reason-item,
.dark .watch-advice,
.dark .gate-link,
.dark .ledger-slip {
  border-color: rgba(93, 88, 73, 0.52);
}

.dark .section-mark a,
.dark .section-mark strong,
.dark .check-card a,
.dark .reason-item em {
  color: #f0b4a8;
}

.dark .waterline span {
  background: linear-gradient(180deg, #a73a2a, rgba(155, 129, 85, 0.72));
  opacity: 0.78;
}

.dark .model-river i {
  background: #a73a2a;
  box-shadow: 0 0 0 1px rgba(167, 58, 42, 0.18);
}

.dark .brand-lockup h1,
.dark .watch-copy h2,
.dark .reason-item strong,
.dark .ledger-slip strong,
.dark .call-list strong,
.dark .call-list-footer strong,
.dark .quota-list strong,
.dark .model-river strong,
.dark .platform-list strong,
.dark .empty-note strong {
  color: #f4efe4;
}

.dark .call-list-footer span {
  color: #9b8f79;
}

.dark .call-list-footer-actions a {
  border-color: rgba(93, 88, 73, 0.58);
  color: #f0b4a8;
}

.dark .call-list-footer-actions a:hover,
.dark .call-list-footer-actions a:focus-visible {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.12);
}

.dark .focus-note {
  border-color: rgba(167, 58, 42, 0.28);
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.12), transparent 72%),
    rgba(24, 26, 21, 0.84);
}

.dark .focus-note strong {
  color: #f4efe4;
}

.dark .focus-note-actions a {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(17, 19, 15, 0.42);
  color: #d9d0be;
}

.dark .focus-note-actions a:hover,
.dark .focus-note-actions a:focus-visible {
  border-color: rgba(216, 205, 185, 0.34);
  background: rgba(216, 205, 185, 0.06);
  color: #f4efe4;
}

.dark .focus-surface {
  border-color: rgba(167, 58, 42, 0.34);
  box-shadow: 0 22px 48px -38px rgba(0, 0, 0, 0.52);
}

.dark .account-menu-dropdown {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    radial-gradient(circle at 18% 12%, rgba(167, 58, 42, 0.08), transparent 18%),
    rgba(24, 26, 21, 0.98);
}

.dark .account-menu-dropdown a,
.dark .account-menu-dropdown button {
  color: #d9d0be;
}

.dark .account-menu-dropdown a:hover,
.dark .account-menu-dropdown a:focus-visible,
.dark .account-menu-dropdown button:hover,
.dark .account-menu-dropdown button:focus-visible {
  background: rgba(167, 58, 42, 0.12);
  color: #f4efe4;
}

.dark .account-theme {
  border-color: rgba(48, 52, 43, 0.9);
}

.dark .account-menu-dropdown .account-theme-options button {
  border-color: rgba(48, 52, 43, 0.86);
  background: rgba(17, 19, 15, 0.64);
  color: #c9c0ac;
}

.dark .account-menu-dropdown .account-theme-options button:hover,
.dark .account-menu-dropdown .account-theme-options button:focus-visible,
.dark .account-menu-dropdown .account-theme-options button.is-selected {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.14);
  color: #f0b4a8;
}

.dark .home-state {
  border-color: rgba(68, 71, 58, 0.92);
  background:
    radial-gradient(circle at 50% 22%, rgba(167, 58, 42, 0.09), transparent 22%),
    radial-gradient(circle at 50% 72%, rgba(155, 129, 85, 0.06), transparent 18%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.96), rgba(17, 19, 15, 0.94));
}

.dark .home-state .spinner {
  color: #21b8b1;
}

.dark .mini-state {
  border-color: rgba(68, 71, 58, 0.68);
  background: rgba(17, 19, 15, 0.72);
  color: #9b8f79;
}

.dark .empty-note {
  color: #9b8f79;
}

.dark .empty-note strong {
  color: #f4efe4;
}

.dark .empty-actions a {
  border-color: rgba(48, 52, 43, 0.9);
  background: rgba(17, 19, 15, 0.42);
  color: #d9d0be;
}

.dark .empty-actions a:hover {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.12);
  color: #f4efe4;
}

.dark .home-state-error {
  border-color: rgba(167, 58, 42, 0.28);
  background:
    radial-gradient(circle at 50% 22%, rgba(167, 58, 42, 0.08), transparent 22%),
    linear-gradient(180deg, rgba(24, 26, 21, 0.96), rgba(17, 19, 15, 0.94));
}

.dark .home-state-error .btn-primary {
  border-color: rgba(167, 58, 42, 0.42);
  background: rgba(167, 58, 42, 0.14);
  color: #f4efe4;
}

.dark .home-state-error .btn-primary:hover,
.dark .home-state-error .btn-primary:focus-visible {
  border-color: rgba(167, 58, 42, 0.56);
  background: rgba(167, 58, 42, 0.22);
}

.dark .account-mark strong,
.dark .account-mark small,
.dark .watch-copy p,
.dark .watch-notice small,
.dark .notice-summary small,
.dark .reason-item small,
.dark .ledger-slip small,
.dark .focus-note p,
.dark .call-list span,
.dark .call-list-footer small,
.dark .quota-list span,
.dark .quota-list small,
.dark .model-river span,
.dark .platform-list span,
.dark .platform-list small,
.dark .empty-note {
  color: #879186;
}

.dark .watch-notice {
  border-top-color: rgba(68, 71, 58, 0.72);
}

.dark .watch-notice > span {
  color: #9b8f79;
}

.dark .watch-notice a {
  color: #d9d0be;
}

.dark .watch-notice a:hover,
.dark .watch-notice a:focus-visible {
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 76%);
  color: #f4efe4;
}

.dark .notice-strip {
  border-bottom-color: rgba(93, 88, 73, 0.68);
}

.dark .notice-strip-copy strong {
  color: #f4efe4;
}

.dark .notice-strip-copy small,
.dark .notice-strip-meta small {
  color: #879186;
}

.dark .notice-strip-meta em {
  color: #f0b4a8;
}

.dark .waterline span.is-empty {
  opacity: 0.22;
}

.dark .waterline-axis span {
  color: #9b8f79;
}

@media (min-width: 1121px) and (max-width: 1320px) {
  .courtyard-stage,
  .water-ledger {
    grid-template-columns: 1fr;
  }

  .ledger-slips {
    height: auto;
  }

  .watch-health-sheet .health-check-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 1120px) {
  .courtyard-stage,
  .water-ledger,
  .health-check-grid {
    grid-template-columns: 1fr;
  }

  .watch-desk {
    grid-template-columns: 1fr;
  }

  .watch-reasons {
    min-height: 0;
    border-left: 0;
    border-top: 1px solid rgba(198, 184, 157, 0.5);
    padding-left: 0;
    padding-top: 0.3rem;
  }

  .gate-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .ledger-slips {
    height: auto;
  }

  .slips-head {
    padding-bottom: 0.42rem;
  }

  .ledger-slip {
    flex: none;
  }
}

@media (max-width: 760px) {
  .sst-dashboard {
    padding: 0.5rem;
  }

  .console-head {
    align-items: start;
    flex-direction: column;
  }

  .account-mark {
    width: 100%;
    min-width: 0;
    border-left: 2px solid var(--sst-seal);
    border-right: 0;
    padding-left: 0.8rem;
    padding-right: 0;
    text-align: left;
  }

  .brand-lockup h1 {
    font-size: 1.42rem;
  }

  .courtyard-stage,
  .water-ledger {
    padding: 0.65rem;
  }

  .watch-desk {
    padding: 0.9rem;
  }

  .watch-status {
    grid-template-columns: 1fr;
    align-items: start;
  }

  .gate-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .watch-health-sheet .health-check-grid {
    grid-template-columns: 1fr;
  }

  .center-seal {
    width: 5.8rem;
    height: 5.8rem;
  }

  .ink-wash-side {
    width: 7rem;
    height: 7rem;
  }

  .ledger-slips,
  .ledger-side,
  .health-check-grid {
    grid-template-columns: 1fr;
  }

  .focus-note-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .reason-item {
    grid-template-columns: auto minmax(0, 1fr);
  }

  .reason-item em {
    grid-column: 2;
    width: fit-content;
    border-left: 0;
    padding-left: 0;
  }
}

@media (max-width: 520px) {
  .brand-lockup {
    align-items: start;
  }

  .seal-mark {
    width: 1.92rem;
    height: 1.92rem;
  }

  .gate-grid {
    grid-template-columns: 1fr;
  }

  .call-list li,
  .quota-list > div,
  .platform-list > div {
    display: grid;
  }

  .quota-list > div {
    grid-template-columns: 1fr;
    gap: 0.18rem;
  }

  .quota-list > div > strong {
    justify-self: start;
    text-align: left;
  }

  .quota-list > div.quota-inline-row {
    grid-template-columns: 1fr;
    gap: 0.18rem;
  }

  .quota-list > div.quota-inline-row span,
  .quota-list > div.quota-inline-row strong,
  .quota-list > div.quota-inline-row small {
    display: block;
    white-space: normal;
    text-align: left;
  }

  .call-list li > div:last-child,
  .platform-list strong {
    text-align: left;
  }
}
</style>
