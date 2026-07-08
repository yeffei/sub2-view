<template>
  <AppLayout>
    <div
      data-testid="profile-shell"
      class="profile-shell mx-auto max-w-[960px] space-y-5"
    >
      <section class="profile-brief">
        <div>
          <span>{{ profileCopy.hero.kicker }}</span>
          <h1>{{ profileCopy.hero.title }}</h1>
          <p>{{ profileCopy.hero.copy }}</p>
        </div>
        <small>{{ user?.email || profileCopy.hero.accountFallback }}</small>
      </section>

      <details class="profile-fold" open>
        <summary :data-expand="profileCopy.fold.expand" :data-collapse="profileCopy.fold.collapse">
          <span>{{ profileCopy.summary.kicker }}</span>
          <strong>{{ profileCopy.summary.title }}</strong>
        </summary>
        <ProfileInfoCard
          :user="user"
          compact-only
          :linuxdo-enabled="linuxdoOAuthEnabled"
          :dingtalk-enabled="dingtalkOAuthEnabled"
          :oidc-enabled="oidcOAuthEnabled"
          :oidc-provider-name="oidcOAuthProviderName"
          :wechat-enabled="wechatOAuthEnabled"
          :wechat-open-enabled="wechatOAuthOpenEnabled"
          :wechat-mp-enabled="wechatOAuthMPEnabled"
        />
      </details>

      <section class="runway-card">
        <div>
          <span>{{ profileCopy.runway.kicker }}</span>
          <strong>{{ balanceRunwayValue }}</strong>
          <p>{{ balanceRunwayNote }}</p>
        </div>
        <dl>
          <div>
            <dt>{{ profileCopy.runway.baseline }}</dt>
            <dd>{{ balanceRunwayBaseline }}</dd>
          </div>
          <div>
            <dt>{{ profileCopy.runway.action }}</dt>
            <dd>{{ balanceRunwayAction }}</dd>
          </div>
        </dl>
      </section>

      <details
        ref="securityFoldRef"
        class="profile-fold"
        :class="{ 'profile-fold-focus': isBalanceNotifyFocusActive }"
        :open="isBalanceNotifyFocusEnabled"
      >
        <summary :data-expand="profileCopy.fold.expand" :data-collapse="profileCopy.fold.collapse">
          <span>{{ profileCopy.security.kicker }}</span>
          <strong>{{ profileCopy.security.title }}</strong>
        </summary>
        <div class="profile-security-grid">
          <ProfilePasswordForm />

          <ProfileTotpCard />

          <div
            id="balance-notify-section"
            ref="balanceNotifySectionRef"
            class="balance-notify-section"
            :class="{ 'balance-notify-section-focus': isBalanceNotifyFocusActive }"
          >
            <div v-if="isBalanceNotifyFocusEnabled" class="balance-notify-note" role="status">
              <div>
                <span>{{ profileCopy.notify.kicker }}</span>
                <strong>{{ balanceNotifyFocusTitle }}</strong>
                <p>{{ balanceNotifyFocusDetail }}</p>
              </div>
              <router-link to="/dashboard">
                {{ profileCopy.notify.backDashboard }}
              </router-link>
            </div>
            <ProfileBalanceNotifyCard
              v-if="user && balanceLowNotifyEnabled"
              :enabled="user.balance_notify_enabled ?? true"
              :threshold="user.balance_notify_threshold"
              :extra-emails="user.balance_notify_extra_emails ?? []"
              :system-default-threshold="systemDefaultThreshold"
              :user-email="user.email"
            />
            <div v-else class="balance-notify-unavailable card p-5">
              <span>{{ profileCopy.notify.unavailableKicker }}</span>
              <strong>{{ profileCopy.notify.unavailableTitle }}</strong>
              <p>{{ profileCopy.notify.unavailableCopy }}</p>
            </div>
          </div>
        </div>
      </details>

      <div
        v-if="contactInfo"
        class="profile-contact card p-5"
      >
        <Icon name="chat" size="md" />
        <span>{{ t('common.contactSupport') }}</span>
        <strong>{{ contactInfo }}</strong>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usageAPI, type UserDashboardStats as UserStatsType } from '@/api/usage'
import { useI18n } from 'vue-i18n'
import { Icon } from '@/components/icons'
import AppLayout from '@/components/layout/AppLayout.vue'
import ProfileBalanceNotifyCard from '@/components/user/profile/ProfileBalanceNotifyCard.vue'
import ProfileInfoCard from '@/components/user/profile/ProfileInfoCard.vue'
import ProfilePasswordForm from '@/components/user/profile/ProfilePasswordForm.vue'
import ProfileTotpCard from '@/components/user/profile/ProfileTotpCard.vue'
import { isWeChatWebOAuthEnabled } from '@/api/auth'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t, locale } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()
const user = computed(() => authStore.user)
const securityFoldRef = ref<HTMLDetailsElement | null>(null)
const balanceNotifySectionRef = ref<HTMLElement | null>(null)

const contactInfo = ref('')
const dashboardStats = ref<UserStatsType | null>(null)
const balanceLowNotifyEnabled = ref(false)
const systemDefaultThreshold = ref(0)
const linuxdoOAuthEnabled = ref(false)
const dingtalkOAuthEnabled = ref(false)
const wechatOAuthEnabled = ref(false)
const wechatOAuthOpenEnabled = ref<boolean | undefined>(undefined)
const wechatOAuthMPEnabled = ref<boolean | undefined>(undefined)
const oidcOAuthEnabled = ref(false)
const oidcOAuthProviderName = ref('OIDC')
const isBalanceNotifyFocusActive = ref(false)
const isBalanceNotifyFocusEnabled = computed(() => route.query.focus === 'balance-notify')
let balanceNotifyFocusTimer: ReturnType<typeof setTimeout> | null = null

const zhProfileCopy = {
  hero: {
    kicker: '账户安全',
    title: '基础资料、验证方式与余额提醒',
    copy: '身份信息、登录绑定和提醒阈值都在这里整理归位，方便你集中核对。',
    accountFallback: '山枢庭账户'
  },
  summary: { kicker: '账户摘要', title: '基础资料与绑定方式' },
  runway: { kicker: '余额续航', baseline: '估算基线', action: '建议动作' },
  security: { kicker: '安全设置', title: '密码、二次验证与余额提醒' },
  notify: {
    kicker: '余额提醒',
    backDashboard: '回首页查看账户状态',
    unavailableKicker: '当前未开放',
    unavailableTitle: '站点暂未启用余额提醒',
    unavailableCopy: '若你是从余额或额度提示跳到这里，可先回首页查看余量，或补充余额后再决定是否联系管理员开启提醒。'
  },
  fold: { expand: '展开', collapse: '收起' },
  runwayText: {
    notBalanceMode: '当前未以余额模式计费',
    insufficientSample: '近期样本不足',
    baselinePrefix: '近 14 天约 $',
    baselineSuffix: ' / 天',
    siteMode: '按站点当前计费方式运行',
    noSample: '暂无样本',
    exhausted: '已见底',
    about: '约 ',
    notBalanceNote: '当前站点不以余额为主要限制，这里主要用于查看提醒设置。',
    noSampleNote: '近时段样本不足，续航会在形成真实消耗后自动更新。',
    criticalNote: '按近期消耗估算，余额已接近见底，建议优先充值或放缓请求节奏。',
    weekNote: '按近期速度估算，余额续航已进入一周窗口，适合提前安排充值或削峰。',
    calmNote: '当前余额对近期请求仍有缓冲，可以继续观察消耗与提醒阈值。',
    keepEmail: '保持提醒邮箱可用',
    waitSample: '先保持提醒开启，等待样本形成',
    rechargeNow: '优先充值，并开启余额提醒',
    checkThreshold: '校对提醒阈值，提前安排充值',
    keepNotify: '保持提醒开关和邮箱状态正常',
    day: ' 天'
  },
  focus: {
    location: '这里是余额提醒所在位置',
    enableFirst: '可以先开启提醒，再设定阈值',
    checkSettings: '可以在这里核对阈值和通知邮箱',
    unavailable: '当前站点尚未开放余额提醒，所以这次先把你带到最接近的账户安全位置。',
    enableDetail: '若你刚遇到额度或余额相关提示，建议先开启提醒，再补上阈值和接收邮箱，避免再次无声耗尽。',
    checkDetail: '若你是从余额提示跳到这里，可直接核对阈值、邮箱列表与开关状态，确保下次临近阈值时能及时收到通知。'
  }
}

const enProfileCopy = {
  hero: {
    kicker: 'Account security',
    title: 'Profile, sign-in methods, and balance alerts',
    copy: 'Identity details, login bindings, and alert thresholds are gathered here for quick review.',
    accountFallback: 'SST account'
  },
  summary: { kicker: 'Account summary', title: 'Profile details and sign-in bindings' },
  runway: { kicker: 'Balance runway', baseline: 'Estimate baseline', action: 'Suggested action' },
  security: { kicker: 'Security settings', title: 'Password, two-factor auth, and balance alerts' },
  notify: {
    kicker: 'Balance alert',
    backDashboard: 'Back to account status',
    unavailableKicker: 'Not enabled',
    unavailableTitle: 'Balance alerts are not enabled for this site',
    unavailableCopy: 'If you arrived here from a balance or quota notice, check the dashboard first, then top up or ask an administrator to enable alerts.'
  },
  fold: { expand: 'Expand', collapse: 'Collapse' },
  runwayText: {
    notBalanceMode: 'This site is not currently billing mainly by balance',
    insufficientSample: 'Not enough recent samples',
    baselinePrefix: 'About $',
    baselineSuffix: ' / day over the last 14 days',
    siteMode: 'Running under the current site billing mode',
    noSample: 'No sample yet',
    exhausted: 'Exhausted',
    about: 'About ',
    notBalanceNote: 'Balance is not the main limiter for this site right now; this section mainly helps review alert settings.',
    noSampleNote: 'Recent samples are not enough yet. Runway will update after real usage is available.',
    criticalNote: 'Based on recent spend, the balance is nearly exhausted. Top up first or slow request volume.',
    weekNote: 'At the recent pace, runway is within a week. It is a good time to top up or smooth traffic.',
    calmNote: 'The current balance still has room for recent request volume. Keep watching spend and alert thresholds.',
    keepEmail: 'Keep alert email reachable',
    waitSample: 'Keep alerts on while samples accumulate',
    rechargeNow: 'Top up first and enable balance alerts',
    checkThreshold: 'Review the alert threshold and plan a top-up',
    keepNotify: 'Keep alert switches and email status healthy',
    day: ' days'
  },
  focus: {
    location: 'This is where balance alerts live',
    enableFirst: 'Enable alerts first, then set a threshold',
    checkSettings: 'Review thresholds and notification emails here',
    unavailable: 'Balance alerts are not enabled for this site, so this is the closest account security location.',
    enableDetail: 'If you just hit a balance or quota notice, enable alerts first, then add a threshold and recipient email.',
    checkDetail: 'If you came here from a balance notice, review thresholds, email recipients, and the alert switch.'
  }
}

const profileCopy = computed(() => locale.value === 'zh' ? zhProfileCopy : enProfileCopy)

const paymentEnabled = computed(() => {
  const settings = appStore.cachedPublicSettings
  return settings?.payment_enabled ?? true
})

const balanceRunwayDays = computed(() => {
  const balance = user.value?.balance ?? 0
  const avgDailyCost = Math.max(dashboardStats.value?.today_actual_cost || 0, (dashboardStats.value?.total_actual_cost || 0) / 14)
  if (balance <= 0 || avgDailyCost <= 0) return avgDailyCost <= 0 ? null : 0
  return balance / avgDailyCost
})

const balanceRunwayBaseline = computed(() => {
  const avgDailyCost = Math.max(dashboardStats.value?.today_actual_cost || 0, (dashboardStats.value?.total_actual_cost || 0) / 14)
  if (!paymentEnabled.value) return profileCopy.value.runwayText.notBalanceMode
  if (avgDailyCost <= 0) return profileCopy.value.runwayText.insufficientSample
  return profileCopy.value.runwayText.baselinePrefix + avgDailyCost.toFixed(4) + profileCopy.value.runwayText.baselineSuffix
})

const balanceRunwayValue = computed(() => {
  if (!paymentEnabled.value) return profileCopy.value.runwayText.siteMode
  const days = balanceRunwayDays.value
  if (days === null) return profileCopy.value.runwayText.noSample
  if (days <= 0) return profileCopy.value.runwayText.exhausted
  return profileCopy.value.runwayText.about + formatRunwayDays(days)
})

const balanceRunwayNote = computed(() => {
  if (!paymentEnabled.value) return profileCopy.value.runwayText.notBalanceNote
  const days = balanceRunwayDays.value
  if (days === null) return profileCopy.value.runwayText.noSampleNote
  if (days <= 3) return profileCopy.value.runwayText.criticalNote
  if (days <= 7) return profileCopy.value.runwayText.weekNote
  return profileCopy.value.runwayText.calmNote
})

const balanceRunwayAction = computed(() => {
  if (!paymentEnabled.value) return profileCopy.value.runwayText.keepEmail
  const days = balanceRunwayDays.value
  if (days === null) return profileCopy.value.runwayText.waitSample
  if (days <= 3) return profileCopy.value.runwayText.rechargeNow
  if (days <= 7) return profileCopy.value.runwayText.checkThreshold
  return profileCopy.value.runwayText.keepNotify
})

const balanceNotifyFocusTitle = computed(() => {
  if (!balanceLowNotifyEnabled.value) return profileCopy.value.focus.location
  if (user.value?.balance_notify_enabled === false) return profileCopy.value.focus.enableFirst
  return profileCopy.value.focus.checkSettings
})
const balanceNotifyFocusDetail = computed(() => {
  if (!balanceLowNotifyEnabled.value) {
    return profileCopy.value.focus.unavailable
  }
  if (user.value?.balance_notify_enabled === false) {
    return profileCopy.value.focus.enableDetail
  }
  return profileCopy.value.focus.checkDetail
})

const clearBalanceNotifyFocusTimer = () => {
  if (balanceNotifyFocusTimer) {
    clearTimeout(balanceNotifyFocusTimer)
    balanceNotifyFocusTimer = null
  }
}

const activateBalanceNotifyFocus = async () => {
  if (!isBalanceNotifyFocusEnabled.value) return
  await nextTick()
  securityFoldRef.value?.setAttribute('open', '')
  balanceNotifySectionRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  isBalanceNotifyFocusActive.value = false
  requestAnimationFrame(() => {
    isBalanceNotifyFocusActive.value = true
  })
  clearBalanceNotifyFocusTimer()
  balanceNotifyFocusTimer = setTimeout(() => {
    isBalanceNotifyFocusActive.value = false
    balanceNotifyFocusTimer = null
  }, 2600)
}

watch(
  [isBalanceNotifyFocusEnabled, balanceLowNotifyEnabled, user],
  ([enabled]) => {
    if (!enabled) {
      clearBalanceNotifyFocusTimer()
      isBalanceNotifyFocusActive.value = false
      return
    }
    activateBalanceNotifyFocus()
  },
  { immediate: true }
)

onMounted(async () => {
  const profileRefresh = authStore.refreshUser().catch((error) => {
    console.error('Failed to refresh profile:', error)
  })

  const dashboardStatsLoad = usageAPI.getDashboardStats()
    .then((stats) => {
      dashboardStats.value = stats
    })
    .catch((error) => {
      console.error('Failed to load dashboard stats for profile:', error)
      dashboardStats.value = null
    })

  const settingsLoad = appStore.fetchPublicSettings()
    .then((settings) => {
      if (!settings) {
        return
      }
      contactInfo.value = settings.contact_info || ''
      balanceLowNotifyEnabled.value = settings.balance_low_notify_enabled ?? false
      systemDefaultThreshold.value = settings.balance_low_notify_threshold ?? 0
      linuxdoOAuthEnabled.value = settings.linuxdo_oauth_enabled ?? false
      dingtalkOAuthEnabled.value = settings.dingtalk_oauth_enabled ?? false
      wechatOAuthEnabled.value = isWeChatWebOAuthEnabled(settings)
      wechatOAuthOpenEnabled.value = typeof settings.wechat_oauth_open_enabled === 'boolean'
        ? settings.wechat_oauth_open_enabled
        : undefined
      wechatOAuthMPEnabled.value = typeof settings.wechat_oauth_mp_enabled === 'boolean'
        ? settings.wechat_oauth_mp_enabled
        : undefined
      oidcOAuthEnabled.value = settings.oidc_oauth_enabled ?? false
      oidcOAuthProviderName.value = settings.oidc_oauth_provider_name || 'OIDC'
    })
    .catch((error) => {
      console.error('Failed to load settings:', error)
    })

  await Promise.all([profileRefresh, settingsLoad, dashboardStatsLoad])
})

const formatRunwayDays = (days: number) => (days >= 10 ? Math.round(days) : days.toFixed(1)) + profileCopy.value.runwayText.day

onBeforeUnmount(() => {
  clearBalanceNotifyFocusTimer()
})
</script>

<style scoped>
.profile-shell {
  color: #1f2320;
}

.profile-brief {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 1rem;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.035), transparent 28%),
    rgba(250, 247, 239, 0.5);
  padding: 1.1rem 1.2rem;
}

.profile-brief span {
  display: block;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.profile-brief h1 {
  margin-top: 0.35rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.28rem, 1.85vw, 1.7rem);
  font-weight: 600;
}

.profile-brief p {
  max-width: 34rem;
  margin-top: 0.48rem;
  color: #59645a;
  font-size: 0.84rem;
  line-height: 1.72;
}

.profile-brief small {
  color: #59645a;
  font-size: 0.82rem;
  overflow-wrap: anywhere;
  text-align: right;
}

.profile-shell :deep(.card),
.profile-contact {
  border-color: rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.52);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.runway-card {
  display: flex;
  flex-wrap: wrap;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.05), transparent 26%),
    rgba(250, 247, 239, 0.54);
  padding: 1rem 1.1rem;
}

.runway-card span,
.runway-card dt {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.runway-card strong {
  display: block;
  margin-top: 0.3rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.2rem;
  font-weight: 600;
}

.runway-card p {
  max-width: 34rem;
  margin-top: 0.35rem;
  color: #59645a;
  font-size: 0.82rem;
  line-height: 1.7;
}

.runway-card dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(9rem, 1fr));
  gap: 0.8rem;
  min-width: min(100%, 24rem);
}

.runway-card dd {
  margin-top: 0.28rem;
  color: #1f2320;
  font-size: 0.84rem;
  font-weight: 600;
  line-height: 1.5;
}

.profile-fold {
  overflow: hidden;
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.46);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.profile-fold summary {
  display: flex;
  cursor: pointer;
  list-style: none;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.15rem;
  border-bottom: 1px solid transparent;
}

.profile-fold[open] summary {
  border-bottom-color: rgba(198, 184, 157, 0.28);
}

.profile-fold summary::-webkit-details-marker {
  display: none;
}

.profile-fold summary::after {
  content: attr(data-expand);
  flex: 0 0 auto;
  border: 1px solid rgba(167, 58, 42, 0.22);
  border-radius: 999px;
  padding: 0.28rem 0.62rem;
  color: #a73a2a;
  font-size: 0.72rem;
  font-weight: 650;
}

.profile-fold[open] summary::after {
  content: attr(data-collapse);
}

.profile-fold summary span {
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
}

.profile-fold summary strong {
  display: block;
  margin-top: 0.28rem;
  color: #1f2320;
  font-size: 0.96rem;
}

.profile-fold > :deep(.card),
.profile-security-grid {
  margin: 1rem;
}

.profile-security-grid {
  display: grid;
  gap: 1rem;
}

.balance-notify-section {
  display: grid;
  gap: 0.8rem;
  scroll-margin-top: 1.25rem;
}

.balance-notify-note {
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 1rem;
  border: 1px solid rgba(167, 58, 42, 0.18);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 72%),
    rgba(250, 247, 239, 0.7);
  padding: 0.85rem 0.95rem;
}

.balance-notify-note span,
.balance-notify-unavailable span {
  display: block;
  margin-bottom: 0.18rem;
  color: #7b6a53;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.64rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.balance-notify-note strong,
.balance-notify-unavailable strong {
  display: block;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1rem;
  font-weight: 600;
}

.balance-notify-note p,
.balance-notify-unavailable p {
  margin-top: 0.3rem;
  color: #59645a;
  font-size: 0.78rem;
  line-height: 1.7;
}

.balance-notify-note a {
  flex: 0 0 auto;
  border: 1px solid rgba(198, 184, 157, 0.52);
  border-radius: 999px;
  background: rgba(255, 252, 246, 0.56);
  padding: 0.42rem 0.76rem;
  color: #38413a;
  font-size: 0.76rem;
  font-weight: 650;
  transition: border-color 180ms ease, background-color 180ms ease, color 180ms ease, transform 180ms ease;
}

.balance-notify-note a:hover,
.balance-notify-note a:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
  outline: none;
  transform: translateX(1px);
}

.balance-notify-unavailable {
  border-style: dashed;
}

.balance-notify-section-focus {
  animation: balance-notify-focus-glow 2.4s ease;
}

.profile-fold-focus {
  border-color: rgba(167, 58, 42, 0.32);
  box-shadow: 0 18px 42px -36px rgba(167, 58, 42, 0.28);
}

@keyframes balance-notify-focus-glow {
  0% {
    transform: translateY(0.5rem);
    opacity: 0.72;
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

.profile-security-grid {
  display: grid;
  gap: 1rem;
}

.profile-contact {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #59645a;
}

.profile-contact span {
  color: #7b6a53;
  font-size: 0.78rem;
}

.profile-contact strong {
  color: #1f2320;
  font-size: 0.9rem;
  overflow-wrap: anywhere;
}

.dark .profile-shell,
.dark .profile-brief h1,
.dark .profile-contact h3 {
  color: #f4efe4;
}

.dark .profile-brief {
  box-shadow: 0 14px 38px -32px rgba(0, 0, 0, 0.44);
}

.dark .profile-brief,
.dark .profile-shell :deep(.card),
.dark .profile-contact,
.dark .profile-fold,
.dark .runway-card {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .runway-card {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 26%),
    rgba(24, 26, 21, 0.78);
}

.dark .profile-shell :deep([data-testid="profile-overview-hero"]),
.dark .profile-shell :deep([data-testid="profile-overview-hero"] .bg-white),
.dark .profile-shell :deep([data-testid="profile-overview-hero"] .bg-gray-50),
.dark .profile-shell :deep([data-testid="profile-overview-hero"] .bg-gray-100) {
  border-color: rgba(68, 71, 58, 0.88);
  background:
    linear-gradient(180deg, rgba(255, 247, 235, 0.02), transparent 28%),
    rgba(17, 19, 15, 0.72);
}

.dark .profile-fold[open] summary {
  border-bottom-color: rgba(48, 52, 43, 0.82);
}

.dark .profile-fold summary strong,
.dark .profile-contact strong,
.dark .runway-card strong,
.dark .runway-card dd {
  color: #f4efe4;
}

.dark .profile-brief span,
.dark .profile-brief p,
.dark .profile-brief small,
.dark .profile-fold summary span,
.dark .profile-contact,
.dark .profile-contact span,
.dark .balance-notify-note p,
.dark .balance-notify-unavailable p,
.dark .balance-notify-unavailable span,
.dark .runway-card p,
.dark .runway-card span,
.dark .runway-card dt {
  color: #879186;
}

.dark .balance-notify-note,
.dark .balance-notify-unavailable,
.dark .profile-fold-focus {
  border-color: rgba(48, 52, 43, 0.95);
}

.dark .balance-notify-note {
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.08), transparent 72%),
    rgba(24, 26, 21, 0.84);
}

.dark .balance-notify-note strong,
.dark .balance-notify-unavailable strong {
  color: #f4efe4;
}

.dark .balance-notify-note a {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
  color: #d8cdb9;
}

.dark .balance-notify-note a:hover,
.dark .balance-notify-note a:focus-visible {
  border-color: rgba(167, 58, 42, 0.34);
  background: rgba(167, 58, 42, 0.06);
  color: #f0b4a8;
}

.dark .profile-fold-focus {
  box-shadow: 0 22px 48px -38px rgba(0, 0, 0, 0.52);
}

@media (max-width: 640px) {
  .profile-brief {
    align-items: start;
    flex-direction: column;
  }

  .runway-card,
  .runway-card dl {
    grid-template-columns: 1fr;
  }

  .profile-brief small {
    text-align: left;
  }

  .balance-notify-note {
    flex-direction: column;
  }

  .balance-notify-note a {
    width: fit-content;
  }
}
</style>
