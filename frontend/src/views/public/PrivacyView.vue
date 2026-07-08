<template>
  <PublicPageLayout
    class="privacy-page"
    tone="legal"
    :eyebrow="privacyContent.eyebrow"
    :title="privacyContent.title"
    intro=""
    :show-cta="false"
  >
      <section class="privacy-main-grid mt-10 grid gap-8 lg:gap-8">
        <aside class="self-start">
          <div class="sticky top-5 h-fit">
            <div class="mb-4 text-xs uppercase tracking-[0.24em] text-zen-mist dark:text-zen-stone">{{ privacyContent.toc }}</div>
            <nav class="grid gap-1 text-sm text-zen-ink dark:text-zen-paper">
              <button
                v-for="item in sections"
                :key="item.id"
                type="button"
                class="privacy-toc-link"
                :class="activeSection === item.id ? 'privacy-toc-link-active' : ''"
                @click="activeSection = item.id"
              >
                {{ item.label }}
              </button>
            </nav>
          </div>
        </aside>

        <article class="privacy-article rounded-[1.35rem] border border-zen-paperLine/70 bg-white/62 p-5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/76 sm:p-6 lg:p-6">
          <section v-if="activeSection === 'overview'" id="overview">
            <div class="privacy-kicker">{{ privacyContent.overview.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.overview.title }}</h2>
            <p class="privacy-copy">
              {{ privacyContent.overview.copy }}
            </p>
            <div class="privacy-notice mt-5">
              <Icon name="shield" size="md" class="mt-0.5 text-zen-seal" />
              <div>
                <div class="privacy-notice-title">{{ privacyContent.overview.noticeTitle }}</div>
                <p class="privacy-notice-copy">{{ privacyContent.overview.noticeCopy }}</p>
              </div>
            </div>
          </section>

          <section v-else-if="activeSection === 'region'" id="region">
            <div class="privacy-kicker">{{ privacyContent.region.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.region.title }}</h2>
            <div class="mt-6 space-y-4">
              <article v-for="item in regionItems" :key="item.title" class="privacy-card">
                <div class="privacy-card-title">{{ item.title }}</div>
                <p class="privacy-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'collect'" id="collect">
            <div class="privacy-kicker">{{ privacyContent.collect.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.collect.title }}</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in collectItems" :key="item.title" class="privacy-card">
                <div class="flex items-start gap-3">
                  <span class="privacy-card-icon"><Icon :name="item.icon" size="md" /></span>
                  <div>
                    <div class="privacy-card-title">{{ item.title }}</div>
                    <p class="privacy-card-copy">{{ item.copy }}</p>
                  </div>
                </div>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'usage'" id="usage">
            <div class="privacy-kicker">{{ privacyContent.usage.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.usage.title }}</h2>
            <div class="privacy-table mt-5 overflow-hidden rounded-[1rem] border border-zen-paperLine/70 dark:border-zen-nightLine">
              <div class="grid grid-cols-[12rem_minmax(0,1fr)] bg-white/32 text-xs uppercase tracking-[0.12em] text-zen-mist dark:bg-zen-nightPanel/55 dark:text-zen-stone">
                <span class="px-4 py-3">{{ privacyContent.usage.tablePurpose }}</span>
                <span class="px-4 py-3">{{ privacyContent.usage.tableDescription }}</span>
              </div>
              <div v-for="item in usageItems" :key="item.title" class="privacy-table-row grid grid-cols-[12rem_minmax(0,1fr)] border-t border-zen-paperLine/60 bg-white/26 text-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/45">
                <div class="px-4 py-4 font-medium text-zen-ink dark:text-zen-paper">{{ item.title }}</div>
                <p class="px-4 py-4 leading-7 text-zen-mist dark:text-zen-stone">{{ item.copy }}</p>
              </div>
            </div>
          </section>

          <section v-else-if="activeSection === 'sharing'" id="sharing">
            <div class="privacy-kicker">{{ privacyContent.sharing.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.sharing.title }}</h2>
            <div class="mt-6 space-y-4">
              <article v-for="item in sharingItems" :key="item.title" class="privacy-card">
                <div class="privacy-card-title">{{ item.title }}</div>
                <p class="privacy-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'retention'" id="retention">
            <div class="privacy-kicker">{{ privacyContent.retention.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.retention.title }}</h2>
            <ul class="mt-5 grid gap-3 text-sm leading-8 text-zen-mist dark:text-zen-stone sm:text-base">
              <li v-for="item in retentionItems" :key="item.title"><strong class="text-zen-ink dark:text-zen-paper">{{ item.title }}{{ privacyContent.colon }}</strong>{{ item.copy }}</li>
            </ul>
          </section>

          <section v-else-if="activeSection === 'rights'" id="rights">
            <div class="privacy-kicker">{{ privacyContent.rights.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.rights.title }}</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in rightsItems" :key="item.title" class="privacy-card">
                <div class="privacy-card-title">{{ item.title }}</div>
                <p class="privacy-card-copy">{{ item.copy }}</p>
              </article>
            </div>
            <div class="privacy-notice mt-6">
              <Icon name="mail" size="md" class="mt-0.5 text-zen-seal" />
              <div>
                <div class="privacy-notice-title">{{ privacyContent.rights.noticeTitle }}</div>
                <p class="privacy-notice-copy">{{ privacyContent.rights.noticeBefore }} {{ contactInfoLabel }} {{ privacyContent.rights.noticeAfter }}</p>
              </div>
            </div>
          </section>

          <section v-else id="updates">
            <div class="privacy-kicker">{{ privacyContent.updates.kicker }}</div>
            <h2 class="privacy-title">{{ privacyContent.updates.title }}</h2>
            <div class="mt-5 space-y-4">
              <article class="privacy-card">
                <div class="privacy-card-title">{{ privacyContent.updates.policyTitle }}</div>
                <p class="privacy-card-copy">{{ privacyContent.updates.policyCopy }}</p>
              </article>
              <article class="privacy-card">
                <div class="privacy-card-title">{{ privacyContent.updates.minorsTitle }}</div>
                <p class="privacy-card-copy">{{ privacyContent.updates.minorsCopy }}</p>
              </article>
            </div>
          </section>
        </article>
      </section>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'

type PrivacySectionId = 'overview' | 'region' | 'collect' | 'usage' | 'sharing' | 'retention' | 'rights' | 'updates'
type IconName = InstanceType<typeof Icon>['$props']['name']
type SectionItem = { id: PrivacySectionId; label: string }
type IconCardItem = { title: string; copy: string; icon: IconName }

const appStore = useAppStore()
const { locale } = useI18n()

const zhPrivacyContent = {
  eyebrow: '隐私',
  title: '隐私政策',
  toc: '目录',
  colon: '：',
  overview: {
    kicker: '总则',
    title: '本政策适用于站点公开页面、账户体系与控制台服务。',
    copy: '当你浏览首页、登录、注册、购买、创建 Key、查看用量或联系管理员时，系统会处理提供服务所必需的账户、账册、请求与安全信息。除法律法规另有要求外，我们遵循最少收集、目的明确、期限合理与按需可见的处理原则。',
    noticeTitle: '第三方边界',
    noticeCopy: '如果你通过本服务访问第三方模型、支付渠道或其他外部服务，相关第三方仍将按照其自身规则处理必要数据；本政策不替代第三方的隐私条款。',
  },
  region: { kicker: '适用地域与数据位置', title: '本服务的适用用户范围、服务器位置与数据处理边界如下。' },
  collect: { kicker: '收集范围', title: '我们可能收集的主要信息类型如下。' },
  usage: { kicker: '使用目的', title: '收集到的信息仅用于与服务运行相关的明确目的。', tablePurpose: '用途', tableDescription: '说明' },
  sharing: { kicker: '共享与外部处理', title: '我们仅在提供服务、履行义务或处理争议所必需时对外提供信息。' },
  retention: { kicker: '保存期限', title: '不同类型的信息会按照其处理目的保留至合理期限届满。' },
  rights: {
    kicker: '你的选择',
    title: '在适用法律允许的范围内，你可以申请访问、更正、删除或限制处理。',
    noticeTitle: '请求方式',
    noticeBefore: '如需提交隐私请求，请通过',
    noticeAfter: '联系我们。为保护账户安全，我们可能会在处理删除、导出或修正请求前要求完成身份核验。',
  },
  updates: {
    kicker: '更新与未成年人',
    title: '当政策内容发生实质变化时，我们会更新本页；未成年人不应独立使用本服务。',
    policyTitle: '政策更新',
    policyCopy: '当收集范围、处理目的、共享方式或你的权利发生实质变化时，我们会更新本页，并同步更新页面顶端日期。更新后的政策自公布之日起生效。',
    minorsTitle: '未成年人',
    minorsCopy: '若你未达到所在地法律要求的使用年龄，请不要独立注册、购买或提交个人信息；如确需使用，应在监护人同意与监督下进行。',
  },
} as const

const enPrivacyContent = {
  eyebrow: 'Privacy',
  title: 'Privacy Policy',
  toc: 'Contents',
  colon: ': ',
  overview: {
    kicker: 'Overview',
    title: 'This policy applies to the public site, account system, and console services.',
    copy: 'When you browse public pages, sign in, register, purchase, create Keys, review usage, or contact support, the system processes account, ledger, request, and security information needed to provide the service. Unless applicable law requires otherwise, we follow principles of minimal collection, clear purpose, reasonable retention, and need-to-know access.',
    noticeTitle: 'Third-party boundary',
    noticeCopy: 'If you use this service to access third-party models, payment channels, or other external services, those third parties process necessary data under their own rules. This policy does not replace third-party privacy terms.',
  },
  region: { kicker: 'Region and data location', title: 'The eligible user scope, server location, and processing boundaries are as follows.' },
  collect: { kicker: 'Information collected', title: 'The main categories of information we may collect are listed below.' },
  usage: { kicker: 'Use purposes', title: 'Collected information is used only for clear purposes related to service operation.', tablePurpose: 'Purpose', tableDescription: 'Description' },
  sharing: { kicker: 'Sharing and external processing', title: 'We provide information externally only when necessary to deliver service, meet obligations, or handle disputes.' },
  retention: { kicker: 'Retention', title: 'Different information types are retained for reasonable periods according to their processing purpose.' },
  rights: {
    kicker: 'Your choices',
    title: 'Where applicable law allows, you may request access, correction, deletion, or processing restrictions.',
    noticeTitle: 'Request method',
    noticeBefore: 'To submit a privacy request, contact us through',
    noticeAfter: 'We may require identity verification before processing deletion, export, or correction requests to protect account security.',
  },
  updates: {
    kicker: 'Updates and minors',
    title: 'We update this page when material policy changes occur; minors should not use the service independently.',
    policyTitle: 'Policy updates',
    policyCopy: 'When collection scope, processing purposes, sharing methods, or your rights materially change, we will update this page and the date at the top. Updated policies take effect when published.',
    minorsTitle: 'Minors',
    minorsCopy: 'If you have not reached the legal age required in your location, do not register, purchase, or submit personal information independently. If use is necessary, it should occur with guardian consent and supervision.',
  },
} as const

const privacyContent = computed(() => locale.value.startsWith('zh') ? zhPrivacyContent : enPrivacyContent)

const zhSections: SectionItem[] = [
  { id: 'overview', label: '总述' },
  { id: 'region', label: '适用地域' },
  { id: 'collect', label: '收集范围' },
  { id: 'usage', label: '使用目的' },
  { id: 'sharing', label: '共享边界' },
  { id: 'retention', label: '保存期限' },
  { id: 'rights', label: '你的选择' },
  { id: 'updates', label: '更新与未成年人' },
] as const

const enSections: SectionItem[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'region', label: 'Regions' },
  { id: 'collect', label: 'Collection' },
  { id: 'usage', label: 'Use purposes' },
  { id: 'sharing', label: 'Sharing' },
  { id: 'retention', label: 'Retention' },
  { id: 'rights', label: 'Your choices' },
  { id: 'updates', label: 'Updates and minors' },
]

const sections = computed(() => locale.value.startsWith('zh') ? zhSections : enSections)

type SectionId = PrivacySectionId

const activeSection = ref<SectionId>('overview')

const contactInfoLabel = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || (locale.value.startsWith('zh') ? '站点管理员联系入口' : 'site administrator contact channel'))

const zhRegionItems = [
  {
    title: '适用用户范围',
    copy: '本服务目前仅面向中国大陆以外的地区和用户提供。中国大陆用户不得注册、购买或使用本服务；若你的访问、注册、购买或使用行为受中国大陆相关限制，请不要继续使用。',
  },
  {
    title: '服务器与数据中心位置',
    copy: '本服务的服务器和主要数据处理设施不位于中国大陆。你访问本服务并提交相关信息，即表示你理解这些信息可能在中国大陆以外的地区存储、传输和处理。',
  },
  {
    title: '业务运营与数据处理',
    copy: '与本服务有关的业务运营、技术支持、账务处理、日志分析与安全监测均在中国大陆以外进行。若接入第三方模型、支付或基础设施服务，相关数据还可能依据其自身规则继续在境外处理。',
  },
] as const

const enRegionItems = [
  { title: 'Eligible user scope', copy: 'This service is currently provided only to regions and users outside mainland China. Mainland China users must not register, purchase, or use this service. If your access, registration, purchase, or use is subject to mainland China restrictions, do not continue.' },
  { title: 'Server and data center location', copy: 'The service servers and primary data processing facilities are not located in mainland China. By accessing the service and submitting information, you understand that information may be stored, transmitted, and processed outside mainland China.' },
  { title: 'Operations and processing', copy: 'Business operations, technical support, accounting, log analysis, and security monitoring related to the service are performed outside mainland China. Third-party model, payment, or infrastructure providers may also process related data under their own rules.' },
]
const regionItems = computed(() => locale.value.startsWith('zh') ? zhRegionItems : enRegionItems)

const zhCollectItems: IconCardItem[] = [
  {
    title: '账户与身份信息',
    copy: '当你注册、登录、绑定第三方身份或找回账号时，系统会处理邮箱、用户名、身份来源及必要的验证状态。',
    icon: 'userCircle',
  },
  {
    title: '账册与订单信息',
    copy: '当你购买、充值、开票或查看账单时，系统会处理订单号、支付结果、金额、订阅或权益状态。',
    icon: 'creditCard',
  },
  {
    title: '密钥与调用记录',
    copy: '当你创建 API Key、调用接口、查看用量时，系统会处理 key 标识、请求时间、模型、消耗、错误状态等运行数据。',
    icon: 'key',
  },
  {
    title: '设备与日志信息',
    copy: '站点会记录基础访问日志、IP、浏览器信息、来源页与错误日志，用于保障安全、排障和性能稳定。',
    icon: 'server',
  },
  {
    title: '联系与反馈信息',
    copy: '当你联系管理员、提交工单或反馈问题时，我们会处理你主动提供的联系方式与问题描述。',
    icon: 'mail',
  },
  {
    title: '你主动提交的内容',
    copy: '你在输入框、备注、工单、工单附件或可视化文档中主动填写的内容，也可能包含个人信息，请按需提交。',
    icon: 'document',
  },
] as const

const enCollectItems: IconCardItem[] = [
  { title: 'Account and identity information', copy: 'When you register, sign in, bind third-party identities, or recover an account, the system processes email, username, identity source, and necessary verification status.', icon: 'userCircle' },
  { title: 'Ledger and order information', copy: 'When you purchase, recharge, request invoices, or view bills, the system processes order numbers, payment results, amounts, subscriptions, or entitlement status.', icon: 'creditCard' },
  { title: 'Keys and call records', copy: 'When you create API Keys, call APIs, or review usage, the system processes Key identifiers, request time, model, consumption, error status, and similar operational data.', icon: 'key' },
  { title: 'Device and log information', copy: 'The site records basic access logs, IP address, browser information, referrer, and error logs for security, troubleshooting, and performance stability.', icon: 'server' },
  { title: 'Contact and feedback information', copy: 'When you contact administrators, submit tickets, or report issues, we process the contact details and issue descriptions you provide.', icon: 'mail' },
  { title: 'Content you submit', copy: 'Content you enter in forms, notes, tickets, attachments, or visual documents may contain personal information. Submit only what is needed.', icon: 'document' },
]
const collectItems = computed(() => locale.value.startsWith('zh') ? zhCollectItems : enCollectItems)

const zhUsageItems = [
  {
    title: '提供服务',
    copy: '用于完成注册、登录、支付、创建 Key、展示可用模型、记录用量、返回接口结果以及维护账户状态。',
  },
  {
    title: '风控与安全',
    copy: '用于识别异常登录、滥用调用、支付争议、恶意流量、凭证泄露风险及违反站点规则的行为。',
  },
  {
    title: '运维与排障',
    copy: '用于定位错误、查看链路状态、恢复订单异常、分析接口可用性、容量压力与服务稳定性。',
  },
  {
    title: '账册与通知',
    copy: '用于发送支付结果、余额提醒、配额提醒、账号验证以及与服务连续性直接相关的必要通知。',
  },
  {
    title: '合规与留档',
    copy: '在法律、监管或争议处理要求下，用于保留必要记录、配合审计、处理申诉与履行法定义务。',
  },
] as const

const enUsageItems = [
  { title: 'Providing service', copy: 'Used to complete registration, sign-in, payment, Key creation, model availability display, usage records, API responses, and account status maintenance.' },
  { title: 'Risk control and security', copy: 'Used to identify abnormal sign-ins, abusive calls, payment disputes, malicious traffic, credential leakage risk, and violations of site rules.' },
  { title: 'Operations and troubleshooting', copy: 'Used to locate errors, review link status, recover abnormal orders, and analyze API availability, capacity pressure, and service stability.' },
  { title: 'Ledger and notifications', copy: 'Used to send payment results, balance reminders, quota reminders, account verification, and other necessary notices directly related to service continuity.' },
  { title: 'Compliance and records', copy: 'Used to retain necessary records, support audits, handle appeals, and meet legal obligations when required by law, regulation, or dispute handling.' },
]
const usageItems = computed(() => locale.value.startsWith('zh') ? zhUsageItems : enUsageItems)

const zhSharingItems = [
  {
    title: 'Token 与密钥数据',
    copy: '本服务不会以出售、交易或向无关第三方提供的方式处理用户的 Token、API Key 或同类凭证数据。相关数据仅在完成鉴权校验、请求转发、计费统计、安全审计与必要排障的范围内按需处理。',
  },
  {
    title: '上游模型与渠道',
    copy: '当你发起模型请求时，请求内容、模型标识及相关技术元数据可能会传递至对应上游服务或通道方，以完成实际推理。',
  },
  {
    title: '支付与账务服务',
    copy: '当你付款、申请退款或处理账务争议时，必要的订单与支付信息会提供给支付服务商或账务处理方。',
  },
  {
    title: '基础设施与安全服务',
    copy: '站点可能使用云主机、对象存储、邮件发送、验证码、风控或日志服务；上述服务仅能在其职责范围内处理提供服务所必需的数据。',
  },
  {
    title: '法律或争议处理',
    copy: '当法律要求、监管调查、安全事件、执法协助或交易争议处理确有必要时，我们可能依法披露相关信息。',
  },
] as const

const enSharingItems = [
  { title: 'Tokens and credential data', copy: 'This service does not sell, trade, or provide user Tokens, API Keys, or similar credentials to unrelated third parties. Such data is processed only as needed for authentication, forwarding, metering, security audit, and necessary troubleshooting.' },
  { title: 'Upstream models and channels', copy: 'When you send a model request, request content, model identifiers, and related technical metadata may be passed to the corresponding upstream service or channel provider to complete inference.' },
  { title: 'Payment and accounting services', copy: 'When you pay, request refunds, or handle accounting disputes, necessary order and payment information is provided to payment providers or accounting processors.' },
  { title: 'Infrastructure and security services', copy: 'The site may use cloud hosting, object storage, email delivery, captcha, risk control, or logging services. These services may process only the data necessary for their assigned responsibilities.' },
  { title: 'Legal or dispute handling', copy: 'When legally required or necessary for regulatory inquiries, security incidents, law-enforcement assistance, or transaction disputes, we may disclose relevant information according to law.' },
]
const sharingItems = computed(() => locale.value.startsWith('zh') ? zhSharingItems : enSharingItems)

const zhRetentionItems = [
  {
    title: '账户信息',
    copy: '在账号持续有效期间保留；注销后，将在完成安全核验、争议处理、审计要求与法定义务后按需删除、去标识化或匿名化。',
  },
  {
    title: '订单与账务记录',
    copy: '会按照财务、税务、支付结算与争议处理要求保留合理期限，不会因页面删除或账户退出而立即消失。',
  },
  {
    title: '访问日志与错误日志',
    copy: '通常仅按运维、安全、风控和审计需要保留；超过必要期限后会按轮转策略清理或覆盖。',
  },
  {
    title: '请求内容',
    copy: '是否保留及保留期限取决于具体功能、上游服务、排障状态与站点配置；原则上仅在提供服务、排障、安全与争议处理所必需的范围内处理。',
  },
] as const

const enRetentionItems = [
  { title: 'Account information', copy: 'Retained while the account remains active. After cancellation, it may be deleted, de-identified, or anonymized after security verification, dispute handling, audit needs, and legal obligations are completed.' },
  { title: 'Orders and accounting records', copy: 'Retained for reasonable periods required by finance, tax, payment settlement, and dispute handling, and will not disappear immediately because a page is deleted or an account signs out.' },
  { title: 'Access and error logs', copy: 'Usually retained only as needed for operations, security, risk control, and audit. Data beyond the necessary period is cleaned or overwritten according to rotation policies.' },
  { title: 'Request content', copy: 'Whether and how long request content is retained depends on the feature, upstream service, troubleshooting state, and site configuration. In principle, it is processed only as needed for service delivery, troubleshooting, security, and dispute handling.' },
]
const retentionItems = computed(() => locale.value.startsWith('zh') ? zhRetentionItems : enRetentionItems)

const zhRightsItems = [
  {
    title: '查看与更正',
    copy: '你可以要求查看我们持有的与你账户相关的信息，并对其中明显错误、过期或不完整的部分提出修正请求。',
  },
  {
    title: '删除与注销',
    copy: '在不违反法定义务、财务留档、安全调查、审计要求或争议处理需求的前提下，你可以申请删除账号与相关数据。',
  },
  {
    title: '停止营销类通知',
    copy: '如果未来存在非必要通知，你可以选择取消；但与安全、支付、账号验证、权益变化及服务状态直接相关的通知通常无法完全关闭。',
  },
  {
    title: '减少敏感提交',
    copy: '你始终可以选择不在输入内容、备注、工单或附件中提交与服务无关的敏感个人信息。',
  },
] as const

const enRightsItems = [
  { title: 'Access and correction', copy: 'You may request access to information related to your account and ask us to correct parts that are clearly wrong, outdated, or incomplete.' },
  { title: 'Deletion and account cancellation', copy: 'Where legal obligations, financial records, security investigations, audit requirements, or dispute needs do not prevent it, you may request deletion of your account and related data.' },
  { title: 'Stop marketing notices', copy: 'If non-essential notices exist in the future, you may opt out. Notices directly related to security, payment, account verification, entitlement changes, and service status usually cannot be fully disabled.' },
  { title: 'Reduce sensitive submissions', copy: 'You can always choose not to submit sensitive personal information unrelated to the service in prompts, notes, tickets, or attachments.' },
]
const rightsItems = computed(() => locale.value.startsWith('zh') ? zhRightsItems : enRightsItems)

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.privacy-page :deep(main > section:first-child) {
  grid-template-columns: minmax(0, 1fr);
  gap: 0.35rem;
  padding-top: 1.2rem;
  padding-bottom: 0.45rem;
}

.privacy-page :deep(.public-copy-block) {
  max-width: 52rem;
}

.privacy-page :deep(.public-display) {
  font-size: clamp(2.4rem, 5.3vw, 4.25rem);
  font-weight: 600;
  line-height: 1.02;
  letter-spacing: 0.01em;
}

.privacy-page :deep(.public-intro) {
  margin-top: 1.25rem;
  max-width: 47rem;
  font-family: inherit;
  font-size: clamp(0.96rem, 1.02vw, 1.04rem);
  line-height: 1.88;
  color: #5f685c;
}

.privacy-page :deep(.public-hero-panel),
.privacy-page :deep(.public-cta) {
  display: none;
}

.privacy-main-grid {
  max-width: 60rem;
  grid-template-columns: minmax(0, 11.4rem) minmax(0, 1fr);
}

.privacy-main-grid > *,
.privacy-article {
  min-width: 0;
}

.privacy-lead {
  font-size: clamp(0.96rem, 1.02vw, 1.04rem);
  line-height: 1.88;
}

.privacy-card,
.privacy-notice {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.34);
}

.privacy-meta-strip {
  border-top: 1px solid rgba(216, 205, 185, 0.72);
  padding-top: 1rem;
}

.privacy-meta-item {
  display: grid;
  gap: 0.35rem;
  align-content: start;
}

.privacy-meta-label {
  font-size: 0.72rem;
  line-height: 1.5;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #8f7d63;
}

.privacy-meta-value {
  font-size: 0.92rem;
  line-height: 1.68;
  font-weight: 500;
  color: #1f2320;
}

.privacy-kicker {
  font-size: 0.7rem;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: #7b6a53;
}

.privacy-kicker {
  font-weight: 500;
}

.privacy-toc-link {
  border-left: 1px solid transparent;
  padding: 0.56rem 0 0.56rem 0.7rem;
  text-align: left;
  font-size: 0.95rem;
  line-height: 1.45;
  color: #4a524a;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.privacy-toc-link:hover {
  color: #b95d1f;
}

.privacy-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.8);
  background: rgba(185, 93, 31, 0.05);
  color: #b95d1f;
}

.privacy-title {
  margin-top: 0.85rem;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  width: 100%;
  max-width: none;
  font-size: clamp(1.12rem, 1.22vw, 1.34rem);
  line-height: 1.54;
  letter-spacing: 0.01em;
  color: #1f2320;
}

.privacy-copy {
  margin-top: 1rem;
  max-width: 52rem;
  font-size: 0.96rem;
  line-height: 1.84;
  color: #5f685c;
}

.privacy-card {
  padding: 1.05rem 1.05rem 1.15rem;
}

.privacy-card + .privacy-card {
  margin-top: 0;
}

.privacy-card-icon {
  display: inline-flex;
  height: 2.4rem;
  width: 2.4rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(188, 93, 31, 0.1);
  color: #b95d1f;
}

.privacy-card-title,
.privacy-notice-title {
  color: #1f2320;
  font-size: 0.98rem;
  font-weight: 600;
  line-height: 1.45;
}

.privacy-card-copy,
.privacy-notice-copy {
  margin-top: 0.55rem;
  font-size: 0.93rem;
  line-height: 1.8;
  color: #5f685c;
}

:global(html.dark) .privacy-meta-strip {
  border-top-color: rgba(82, 87, 76, 0.84);
}

:global(html.dark) .privacy-meta-label {
  color: #b8af9a;
}

:global(html.dark) .privacy-meta-value {
  color: #efe5d2;
}

:global(html.dark) .privacy-table > :first-child span {
  color: #b8af9a;
}

:global(html.dark) .privacy-table-row > :first-child {
  color: #efe5d2;
}

.privacy-notice {
  display: flex;
  align-items: flex-start;
  gap: 0.9rem;
  padding: 1rem 1rem 1.05rem;
  background: linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58));
}

:global(html.dark) .privacy-card,
:global(html.dark) .privacy-notice {
  border-color: rgba(82, 87, 76, 0.84);
  background: linear-gradient(180deg, rgba(20, 22, 18, 0.92), rgba(14, 16, 14, 0.96));
}

:global(html.dark) .privacy-title,
:global(html.dark) .privacy-card-title,
:global(html.dark) .privacy-notice-title {
  color: #efe5d2;
}

:global(html.dark) .privacy-copy,
:global(html.dark) .privacy-card-copy,
:global(html.dark) .privacy-notice-copy {
  color: #cdbfab;
}

:global(html.dark) .privacy-kicker {
  color: #b8af9a;
}

:global(html.dark) .privacy-toc-link {
  color: #d8cfbe;
}

:global(html.dark) .privacy-toc-link:hover,
:global(html.dark) .privacy-toc-link-active {
  color: #ffd8bb;
}

:global(html.dark) .privacy-toc-link-active {
  border-left-color: rgba(255, 216, 187, 0.72);
  background: rgba(185, 93, 31, 0.1);
}

:global(html.dark) .privacy-card-icon {
  background: rgba(188, 93, 31, 0.16);
  color: #ffd8bb;
}

:global(html.dark) .privacy-notice {
  border-color: rgba(138, 116, 90, 0.38);
  background:
    linear-gradient(145deg, rgba(29, 32, 26, 0.9), rgba(24, 22, 19, 0.9)),
    radial-gradient(circle at 84% 34%, rgba(173, 89, 36, 0.12), transparent 28%);
  box-shadow:
    inset 0 1px 0 rgba(255, 241, 220, 0.06),
    0 16px 32px rgba(0, 0, 0, 0.28);
}

:global(html.dark) .privacy-notice-title {
  color: #f2e8d5;
}

:global(html.dark) .privacy-notice-copy {
  color: #d5cab9;
}

@media (max-width: 1023px) {
  .privacy-main-grid {
    grid-template-columns: 1fr;
  }

  .privacy-toc-link {
    padding-left: 0.65rem;
    font-size: 0.94rem;
  }

  .privacy-title {
    max-width: none;
    font-size: clamp(1.1rem, 3.5vw, 1.28rem);
    line-height: 1.45;
  }

  .privacy-meta-value,
  .privacy-copy,
  .privacy-card-copy,
  .privacy-notice-copy,
  .privacy-lead {
    font-size: 0.95rem;
  }

  .privacy-table > :first-child,
  .privacy-table-row {
    grid-template-columns: 1fr;
  }

  .privacy-table > :first-child span,
  .privacy-table-row > * {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .privacy-table > :first-child span:last-child {
    padding-top: 0;
  }

  .privacy-table-row > :first-child {
    padding-bottom: 0.35rem;
  }

  .privacy-table-row > :last-child {
    padding-top: 0;
  }
}
</style>

<style>
html.dark .privacy-page .public-intro {
  color: #ead9bd;
}

html.dark .privacy-page .public-display {
  color: #fff4dd;
  text-shadow: 0 1px 0 rgba(255, 240, 218, 0.08), 0 18px 48px rgba(0, 0, 0, 0.28);
}

html.dark .privacy-page .public-copy-block > div:first-child span:last-child,
html.dark .privacy-page .privacy-meta-label,
html.dark .privacy-page .privacy-kicker,
html.dark .privacy-page .privacy-table > :first-child span,
html.dark .privacy-page .privacy-main-grid > aside > div > div:first-child {
  color: #cdb387;
}

html.dark .privacy-page .privacy-meta-value,
html.dark .privacy-page .privacy-title,
html.dark .privacy-page .privacy-card-title,
html.dark .privacy-page .privacy-notice-title,
html.dark .privacy-page .privacy-table-row > :first-child,
html.dark .privacy-page .privacy-main-grid li strong {
  color: #fff0d5;
}

html.dark .privacy-page .privacy-toc-link:hover,
html.dark .privacy-page .privacy-toc-link-active {
  color: #f3c786;
}

html.dark .privacy-page .privacy-toc-link-active {
  border-left-color: rgba(212, 153, 80, 0.8);
  background: linear-gradient(90deg, rgba(176, 120, 57, 0.16), rgba(176, 120, 57, 0.04));
}

html.dark .privacy-page .privacy-card,
html.dark .privacy-page .privacy-notice,
html.dark .privacy-page .privacy-table-row {
  border-color: rgba(141, 109, 72, 0.56) !important;
}

html.dark .privacy-article {
  border-color: rgba(120, 109, 90, 0.54) !important;
  background:
    linear-gradient(180deg, rgba(35, 37, 31, 0.95), rgba(24, 26, 21, 0.96)),
    radial-gradient(circle at top right, rgba(166, 97, 45, 0.1), transparent 28%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 244, 224, 0.07),
    0 18px 40px rgba(0, 0, 0, 0.2) !important;
}

html.dark .privacy-lead,
html.dark .privacy-copy,
html.dark .privacy-card-copy,
html.dark .privacy-notice-copy,
html.dark .privacy-main-grid li,
html.dark .privacy-table-row p {
  color: #e2d6c4;
}

html.dark .privacy-meta-value,
html.dark .privacy-title,
html.dark .privacy-card-title,
html.dark .privacy-notice-title,
html.dark .privacy-table-row > :first-child,
html.dark .privacy-copy strong,
html.dark .privacy-copy b,
html.dark .privacy-copy em,
html.dark .privacy-copy a,
html.dark .privacy-card-copy strong,
html.dark .privacy-notice-copy strong,
html.dark .privacy-main-grid li strong {
  color: #fff0da;
}

html.dark .privacy-meta-label,
html.dark .privacy-kicker,
html.dark .privacy-table > :first-child span,
html.dark .privacy-main-grid > aside > div > div:first-child {
  color: #caba9f;
}

html.dark .privacy-meta-strip {
  border-top-color: rgba(102, 95, 79, 0.68);
}

html.dark .privacy-toc-link {
  color: #f0e3cc;
}

html.dark .privacy-toc-link:hover,
html.dark .privacy-toc-link-active {
  color: #ffe0bf;
}

html.dark .privacy-card,
html.dark .privacy-page .privacy-card,
html.dark .privacy-main-grid .privacy-card,
html.dark .privacy-page .privacy-main-grid .privacy-card {
  border-color: rgba(118, 106, 87, 0.58) !important;
  background:
    linear-gradient(180deg, rgba(36, 38, 31, 0.94), rgba(27, 29, 24, 0.95)),
    radial-gradient(circle at 88% 18%, rgba(163, 97, 45, 0.09), transparent 24%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 242, 219, 0.06),
    0 14px 28px rgba(0, 0, 0, 0.18) !important;
}

html.dark .privacy-card-title,
html.dark .privacy-page .privacy-card-title,
html.dark .privacy-main-grid .privacy-card-title,
html.dark .privacy-page .privacy-main-grid .privacy-card-title {
  color: #f4ead7 !important;
}

html.dark .privacy-card-copy,
html.dark .privacy-page .privacy-card-copy,
html.dark .privacy-main-grid .privacy-card-copy,
html.dark .privacy-page .privacy-main-grid .privacy-card-copy {
  color: #d9cdbb !important;
}

html.dark .privacy-notice,
html.dark .privacy-page .privacy-notice,
html.dark .privacy-main-grid .privacy-notice,
html.dark .privacy-page .privacy-main-grid .privacy-notice {
  border-color: rgba(136, 115, 88, 0.52) !important;
  background:
    linear-gradient(145deg, rgba(42, 38, 31, 0.95), rgba(30, 27, 23, 0.96)),
    radial-gradient(circle at 84% 34%, rgba(173, 89, 36, 0.12), transparent 30%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 241, 220, 0.07),
    0 16px 32px rgba(0, 0, 0, 0.24) !important;
}

html.dark .privacy-notice-title,
html.dark .privacy-page .privacy-notice-title,
html.dark .privacy-main-grid .privacy-notice-title,
html.dark .privacy-page .privacy-main-grid .privacy-notice-title {
  color: #f7ebd8 !important;
}

html.dark .privacy-notice-copy,
html.dark .privacy-page .privacy-notice-copy,
html.dark .privacy-main-grid .privacy-notice-copy,
html.dark .privacy-page .privacy-main-grid .privacy-notice-copy {
  color: #dfd2be !important;
}

html.dark .privacy-table > :first-child {
  background: rgba(35, 32, 27, 0.92) !important;
}

html.dark .privacy-table-row {
  background: rgba(29, 31, 26, 0.72) !important;
  border-top-color: rgba(96, 89, 75, 0.58) !important;
}

html.dark .privacy-article {
  border-color: rgba(155, 126, 86, 0.26) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .privacy-page .privacy-card,
html.dark .privacy-page .privacy-notice,
html.dark .privacy-page .privacy-table,
html.dark .privacy-page .privacy-table-row {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.88), rgba(14, 16, 13, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.08), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 238, 210, 0.05) !important;
}

html.dark .privacy-page .privacy-notice {
  background:
    linear-gradient(135deg, rgba(35, 29, 23, 0.88), rgba(24, 27, 22, 0.84)),
    radial-gradient(circle at 84% 18%, rgba(194, 126, 74, 0.13), transparent 28%) !important;
}

html.dark .privacy-page .privacy-table > :first-child {
  background:
    linear-gradient(180deg, rgba(39, 32, 26, 0.9), rgba(24, 27, 22, 0.88)) !important;
}

html.dark .privacy-page .privacy-title,
html.dark .privacy-page .privacy-card-title,
html.dark .privacy-page .privacy-notice-title,
html.dark .privacy-page .privacy-table-row > :first-child,
html.dark .privacy-page .privacy-main-grid li strong {
  color: #f6e8d2 !important;
}

html.dark .privacy-page .privacy-copy,
html.dark .privacy-page .privacy-card-copy,
html.dark .privacy-page .privacy-notice-copy,
html.dark .privacy-page .privacy-main-grid li,
html.dark .privacy-page .privacy-table-row p {
  color: #d0baa0 !important;
}

html.dark .privacy-page .privacy-kicker,
html.dark .privacy-page .privacy-table > :first-child span,
html.dark .privacy-page .privacy-main-grid > aside > div > div:first-child,
html.dark .privacy-page .public-copy-block > div:first-child span:last-child {
  color: #d8b171 !important;
}

html.dark .privacy-page .privacy-toc-link {
  color: #d4c4ad !important;
}

html.dark .privacy-page .privacy-toc-link:hover,
html.dark .privacy-page .privacy-toc-link-active {
  color: #efc183 !important;
}

html.dark .privacy-page .privacy-toc-link-active {
  border-left-color: rgba(194, 126, 74, 0.72) !important;
  background: linear-gradient(90deg, rgba(194, 126, 74, 0.16), rgba(194, 126, 74, 0.04)) !important;
}

html.dark .privacy-page .privacy-card-icon {
  background: rgba(194, 126, 74, 0.13) !important;
  color: #efc183 !important;
}

html:not(.dark) .privacy-article {
  border-color: rgba(154, 128, 92, 0.16) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    linear-gradient(90deg, rgba(144, 113, 76, 0.038), transparent 18%, rgba(144, 113, 76, 0.024) 82%, transparent),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28) !important;
  box-shadow:
    0 14px 34px rgba(84, 57, 31, 0.05),
    inset 0 1px 0 rgba(255, 249, 239, 0.6),
    inset 0 -1px 0 rgba(140, 111, 76, 0.07),
    inset 0 0 0 1px rgba(255, 255, 255, 0.22) !important;
}

html:not(.dark) .privacy-page .privacy-card,
html:not(.dark) .privacy-page .privacy-notice,
html:not(.dark) .privacy-page .privacy-table,
html:not(.dark) .privacy-page .privacy-table-row {
  border-color: rgba(190, 168, 134, 0.42) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(244, 235, 220, 0.34)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.06), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.42) !important;
}

html:not(.dark) .privacy-page .privacy-notice {
  background:
    linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58)),
    radial-gradient(circle at 84% 18%, rgba(196, 136, 68, 0.08), transparent 28%) !important;
}

html:not(.dark) .privacy-page .privacy-table > :first-child {
  background: rgba(255, 252, 246, 0.48) !important;
}

html:not(.dark) .privacy-page .privacy-toc-link {
  color: #4a524a !important;
}

html:not(.dark) .privacy-page .privacy-toc-link:hover,
html:not(.dark) .privacy-page .privacy-toc-link-active {
  color: #b95d1f !important;
}

html:not(.dark) .privacy-page .privacy-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.74) !important;
  background: linear-gradient(90deg, rgba(185, 93, 31, 0.07), rgba(185, 93, 31, 0.02)) !important;
}

html:not(.dark) .privacy-page .privacy-card-icon {
  background: rgba(188, 93, 31, 0.1) !important;
  color: #b95d1f !important;
}
</style>

