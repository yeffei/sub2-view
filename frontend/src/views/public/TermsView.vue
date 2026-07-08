<template>
  <PublicPageLayout
    class="terms-page"
    tone="legal"
    :eyebrow="termsContent.eyebrow"
    :title="termsContent.title"
    intro=""
    :show-cta="false"
  >
      <section class="terms-main-grid mt-10 grid gap-8 lg:gap-8">
        <aside class="self-start">
          <div class="sticky top-5 h-fit">
            <div class="mb-4 text-xs uppercase tracking-[0.24em] text-zen-mist dark:text-zen-stone">{{ termsContent.toc }}</div>
            <nav class="grid gap-1 text-sm text-zen-ink dark:text-zen-paper">
              <button
                v-for="item in sections"
                :key="item.id"
                type="button"
                class="terms-toc-link"
                :class="activeSection === item.id ? 'terms-toc-link-active' : ''"
                @click="activeSection = item.id"
              >
                {{ item.label }}
              </button>
            </nav>
          </div>
        </aside>

        <article class="terms-article rounded-[1.35rem] border border-zen-paperLine/70 bg-white/62 p-5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/76 sm:p-6 lg:p-6">
          <section v-if="activeSection === 'overview'" id="overview">
            <div class="terms-kicker">{{ termsContent.overview.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.overview.title }}</h2>
            <p class="terms-copy">
              {{ termsContent.overview.copy }}
            </p>
            <div class="terms-notice mt-5">
              <Icon name="shield" size="md" class="mt-0.5 text-zen-seal" />
              <div>
                <div class="terms-notice-title">{{ termsContent.overview.noticeTitle }}</div>
                <p class="terms-notice-copy">{{ termsContent.overview.noticeCopy }}</p>
              </div>
            </div>
          </section>

          <section v-else-if="activeSection === 'eligibility'" id="eligibility">
            <div class="terms-kicker">{{ termsContent.eligibility.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.eligibility.title }}</h2>
            <div class="mt-6 space-y-4">
              <article v-for="item in eligibilityItems" :key="item.title" class="terms-card">
                <div class="terms-card-title">{{ item.title }}</div>
                <p class="terms-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'account'" id="account">
            <div class="terms-kicker">{{ termsContent.account.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.account.title }}</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in accountItems" :key="item.title" class="terms-card">
                <div class="flex items-start gap-3">
                  <span class="terms-card-icon"><Icon :name="item.icon" size="md" /></span>
                  <div>
                    <div class="terms-card-title">{{ item.title }}</div>
                    <p class="terms-card-copy">{{ item.copy }}</p>
                  </div>
                </div>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'billing'" id="billing">
            <div class="terms-kicker">{{ termsContent.billing.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.billing.title }}</h2>
            <div class="terms-table mt-5 overflow-hidden rounded-[1rem] border border-zen-paperLine/70 dark:border-zen-nightLine">
              <div class="grid grid-cols-[12rem_minmax(0,1fr)] bg-white/32 text-xs uppercase tracking-[0.12em] text-zen-mist dark:bg-zen-nightPanel/55 dark:text-zen-stone">
                <span class="px-4 py-3">{{ termsContent.billing.tableItem }}</span>
                <span class="px-4 py-3">{{ termsContent.billing.tableDescription }}</span>
              </div>
              <div v-for="item in billingItems" :key="item.title" class="terms-table-row grid grid-cols-[12rem_minmax(0,1fr)] border-t border-zen-paperLine/60 bg-white/26 text-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/45">
                <div class="px-4 py-4 font-medium text-zen-ink dark:text-zen-paper">{{ item.title }}</div>
                <p class="px-4 py-4 leading-7 text-zen-mist dark:text-zen-stone">{{ item.copy }}</p>
              </div>
            </div>
          </section>

          <section v-else-if="activeSection === 'conduct'" id="conduct">
            <div class="terms-kicker">{{ termsContent.conduct.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.conduct.title }}</h2>
            <div class="mt-6 space-y-4">
              <article v-for="item in conductItems" :key="item.title" class="terms-card">
                <div class="terms-card-title">{{ item.title }}</div>
                <p class="terms-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'thirdparty'" id="thirdparty">
            <div class="terms-kicker">{{ termsContent.thirdParty.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.thirdParty.title }}</h2>
            <div class="mt-6 space-y-4">
              <article v-for="item in thirdPartyItems" :key="item.title" class="terms-card">
                <div class="terms-card-title">{{ item.title }}</div>
                <p class="terms-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else-if="activeSection === 'availability'" id="availability">
            <div class="terms-kicker">{{ termsContent.availability.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.availability.title }}</h2>
            <ul class="mt-5 grid gap-3 text-sm leading-8 text-zen-mist dark:text-zen-stone sm:text-base">
              <li v-for="item in availabilityItems" :key="item.title"><strong class="text-zen-ink dark:text-zen-paper">{{ item.title }}{{ termsContent.colon }}</strong>{{ item.copy }}</li>
            </ul>
          </section>

          <section v-else-if="activeSection === 'termination'" id="termination">
            <div class="terms-kicker">{{ termsContent.termination.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.termination.title }}</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in terminationItems" :key="item.title" class="terms-card">
                <div class="terms-card-title">{{ item.title }}</div>
                <p class="terms-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-else id="updates">
            <div class="terms-kicker">{{ termsContent.updates.kicker }}</div>
            <h2 class="terms-title">{{ termsContent.updates.title }}</h2>
            <div class="mt-5 space-y-4">
              <article class="terms-card">
                <div class="terms-card-title">{{ termsContent.updates.termsTitle }}</div>
                <p class="terms-card-copy">{{ termsContent.updates.termsCopy }}</p>
              </article>
              <article class="terms-card">
                <div class="terms-card-title">{{ termsContent.updates.contactTitle }}</div>
                <p class="terms-card-copy">{{ termsContent.updates.contactBefore }} {{ contactInfoLabel }} {{ termsContent.updates.contactAfter }}</p>
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

type TermsSectionId = 'overview' | 'eligibility' | 'account' | 'billing' | 'conduct' | 'thirdparty' | 'availability' | 'termination' | 'updates'
type IconName = InstanceType<typeof Icon>['$props']['name']
type SectionItem = { id: TermsSectionId; label: string }
type IconCardItem = { title: string; copy: string; icon: IconName }

const appStore = useAppStore()
const { locale } = useI18n()

const zhTermsContent = {
  eyebrow: '条款',
  title: '服务条款',
  toc: '目录',
  colon: '：',
  overview: {
    kicker: '总述',
    title: '访问、注册、购买或继续使用本服务，即表示你同意遵守本条款。',
    copy: '本条款适用于站点公开页面、账户体系、控制台、支付与订单相关页面，以及通过本服务发起的 API 调用。若你不同意本条款的任何内容，请不要注册、购买或继续使用本服务。除非适用法律另有要求，本条款与本网站展示的价格、说明、风险提示及隐私政策共同构成对外规则基础。',
    noticeTitle: '与隐私政策配套适用',
    noticeCopy: '关于数据如何收集、保存、共享与保护，请同时参阅隐私政策；本页主要说明使用规则、计费关系、责任边界与终止条件。',
  },
  eligibility: { kicker: '使用资格', title: '本服务并非面向所有地区或所有主体开放，使用前应先确认自身资格。' },
  account: { kicker: '账户与凭证', title: '账号、登录方式、API Key 与控制台权限均由你自行保管并承担相应责任。' },
  billing: { kicker: '计费与订单', title: '价格展示、倍率折算、订单结果与账户账册记录共同决定最终计费结果。', tableItem: '事项', tableDescription: '说明' },
  conduct: { kicker: '禁止行为', title: '你不得利用本服务从事违法、滥用、规避规则或危害系统稳定性的行为。' },
  thirdParty: { kicker: '上游与第三方', title: '本服务可能依赖第三方模型、支付通道、云资源与安全基础设施，相关部分不由我们单方完全控制。' },
  availability: { kicker: '服务可用性', title: '服务能力、模型可用性、价格、入口与风控策略可能随时调整，持续可用并不当然构成承诺。' },
  termination: { kicker: '暂停与终止', title: '如出现违规、争议、风险或运营需要，我们可以限制、暂停或终止全部或部分服务。' },
  updates: {
    kicker: '更新与联系',
    title: '条款会随服务结构、价格体系、风险策略或合规要求变化而更新；争议应先联系管理员处理。',
    termsTitle: '条款更新',
    termsCopy: '当服务范围、计费方式、资格边界、禁止行为、责任限制或终止规则发生实质变化时，我们会更新本页，并同步更新页面顶部日期。更新后的条款自公布之日起生效。',
    contactTitle: '联系与争议',
    contactBefore: '如你对订单、账册、限制措施、资格判断或其他条款适用问题有异议，请先通过',
    contactAfter: '联系我们。双方将优先通过沟通与核验解决；如适用法律另有要求，则按适用法律处理。',
  },
} as const

const enTermsContent = {
  eyebrow: 'Terms',
  title: 'Terms of Service',
  toc: 'Contents',
  colon: ': ',
  overview: {
    kicker: 'Overview',
    title: 'By accessing, registering, purchasing, or continuing to use the service, you agree to these terms.',
    copy: 'These terms apply to public pages, account systems, the console, payment and order pages, and API calls made through this service. If you disagree with any part, do not register, purchase, or continue using the service. Unless applicable law requires otherwise, these terms, displayed prices, explanations, risk notices, and the privacy policy form the external rule basis.',
    noticeTitle: 'Read with the privacy policy',
    noticeCopy: 'For how data is collected, retained, shared, and protected, also read the privacy policy. This page explains usage rules, billing relationships, responsibility boundaries, and termination conditions.',
  },
  eligibility: { kicker: 'Eligibility', title: 'The service is not open to every region or entity; confirm your eligibility before use.' },
  account: { kicker: 'Accounts and credentials', title: 'You are responsible for safeguarding your account, sign-in methods, API Keys, and console permissions.' },
  billing: { kicker: 'Billing and orders', title: 'Displayed prices, rate conversion, order results, and account ledger records together determine final charges.', tableItem: 'Item', tableDescription: 'Description' },
  conduct: { kicker: 'Prohibited conduct', title: 'You must not use the service for unlawful, abusive, evasive, or stability-harming behavior.' },
  thirdParty: { kicker: 'Upstreams and third parties', title: 'The service may depend on third-party models, payment channels, cloud resources, and security infrastructure not fully controlled by us.' },
  availability: { kicker: 'Service availability', title: 'Capabilities, model availability, prices, entrances, and risk-control policies may change at any time; continued availability is not guaranteed.' },
  termination: { kicker: 'Suspension and termination', title: 'For violations, disputes, risks, or operational needs, we may limit, suspend, or terminate all or part of the service.' },
  updates: {
    kicker: 'Updates and contact',
    title: 'Terms may change with service structure, pricing, risk policies, or compliance requirements. Disputes should first be raised with support.',
    termsTitle: 'Terms updates',
    termsCopy: 'When service scope, billing methods, eligibility boundaries, prohibited conduct, liability limits, or termination rules materially change, we will update this page and the date at the top. Updated terms take effect when published.',
    contactTitle: 'Contact and disputes',
    contactBefore: 'If you dispute orders, ledgers, restrictions, eligibility decisions, or any other application of these terms, contact us first through',
    contactAfter: 'Both parties will prioritize communication and verification. If applicable law requires otherwise, applicable law controls.',
  },
} as const

const termsContent = computed(() => locale.value.startsWith('zh') ? zhTermsContent : enTermsContent)

const zhSections: SectionItem[] = [
  { id: 'overview', label: '总述' },
  { id: 'eligibility', label: '使用资格' },
  { id: 'account', label: '账户与凭证' },
  { id: 'billing', label: '计费与订单' },
  { id: 'conduct', label: '禁止行为' },
  { id: 'thirdparty', label: '上游与第三方' },
  { id: 'availability', label: '服务可用性' },
  { id: 'termination', label: '暂停与终止' },
  { id: 'updates', label: '更新与联系' },
] as const

const enSections: SectionItem[] = [
  { id: 'overview', label: 'Overview' },
  { id: 'eligibility', label: 'Eligibility' },
  { id: 'account', label: 'Accounts' },
  { id: 'billing', label: 'Billing' },
  { id: 'conduct', label: 'Conduct' },
  { id: 'thirdparty', label: 'Third parties' },
  { id: 'availability', label: 'Availability' },
  { id: 'termination', label: 'Termination' },
  { id: 'updates', label: 'Updates' },
]

const sections = computed(() => locale.value.startsWith('zh') ? zhSections : enSections)

type SectionId = TermsSectionId

const activeSection = ref<SectionId>('overview')

const contactInfoLabel = computed(() => appStore.cachedPublicSettings?.contact_info?.trim() || (locale.value.startsWith('zh') ? '站点管理员联系入口' : 'site administrator contact channel'))

const zhEligibilityItems = [
  {
    title: '适用地域',
    copy: '本服务目前仅面向中国大陆以外的地区和用户提供。中国大陆用户不得注册、购买或使用本服务；若你的访问或使用行为受中国大陆相关限制，请不要继续使用。',
  },
  {
    title: '使用年龄与资格',
    copy: '你应具备适用法律要求的民事行为能力；若你代表企业、团队或其他组织使用本服务，你应确保自己具有代表该主体接受本条款并承担相应责任的权限。',
  },
  {
    title: '服务器与处理位置',
    copy: '本服务的服务器和主要数据处理设施不位于中国大陆，与服务有关的运营、技术支持、日志分析、账务处理与安全监测均在中国大陆以外进行。',
  },
] as const
const enEligibilityItems = [
  { title: 'Applicable regions', copy: 'This service is currently provided only to regions and users outside mainland China. Mainland China users must not register, purchase, or use this service. If your access or use is subject to mainland China restrictions, do not continue.' },
  { title: 'Age and authority', copy: 'You must have the civil capacity required by applicable law. If you use the service on behalf of a company, team, or organization, you must have authority to accept these terms and assume responsibilities for that entity.' },
  { title: 'Server and processing location', copy: 'The service servers and primary data processing facilities are not located in mainland China. Operations, support, log analysis, accounting, and security monitoring related to the service occur outside mainland China.' },
]
const eligibilityItems = computed(() => locale.value.startsWith('zh') ? zhEligibilityItems : enEligibilityItems)

const zhAccountItems: IconCardItem[] = [
  {
    title: '账户信息真实有效',
    copy: '注册、登录、绑定第三方身份或找回账号时，你应提供真实、完整且持续有效的信息；若信息失真、失效或存在冒用风险，我们可以限制相关账户功能。',
    icon: 'userCircle',
  },
  {
    title: '凭证自行保管',
    copy: '账号密码、邮箱验证码、第三方登录状态、API Key 及其他访问凭证由你自行保管。因泄露、转借、共享、截图传播、上传到公开仓库或其他保管不当造成的后果，由你自行承担。',
    icon: 'key',
  },
  {
    title: '控制台权限边界',
    copy: '控制台中的 Key、用量、余额、订单、订阅、渠道或其他功能，仅在当前账户、当前分组和当前权限范围内开放，并可能随风控、支付、合规或运营策略调整。',
    icon: 'grid',
  },
  {
    title: '异常访问处理',
    copy: '如系统识别到异常登录、批量注册、脚本滥用、共享凭证、代理转售或其他高风险行为，我们可以要求二次核验、临时冻结、撤销 Key、暂停支付或直接终止服务。',
    icon: 'shield',
  },
] as const
const enAccountItems: IconCardItem[] = [
  { title: 'Accurate account information', copy: 'When registering, signing in, binding third-party identities, or recovering an account, you must provide truthful, complete, and valid information. If information is false, invalid, or impersonation risk exists, we may restrict account functions.', icon: 'userCircle' },
  { title: 'Credential protection', copy: 'You are responsible for safeguarding passwords, email codes, third-party sessions, API Keys, and other access credentials. Consequences caused by leakage, lending, sharing, screenshots, public repository uploads, or poor handling are your responsibility.', icon: 'key' },
  { title: 'Console permission boundaries', copy: 'Keys, usage, balance, orders, subscriptions, channels, and other console functions are available only within your current account, group, and permission scope, and may change with risk-control, payment, compliance, or operational policies.', icon: 'grid' },
  { title: 'Abnormal access handling', copy: 'If abnormal sign-ins, bulk registrations, script abuse, shared credentials, proxy resale, or other high-risk behavior is detected, we may require secondary verification, temporarily freeze access, revoke Keys, suspend payment, or terminate service.', icon: 'shield' },
]
const accountItems = computed(() => locale.value.startsWith('zh') ? zhAccountItems : enAccountItems)

const zhBillingItems = [
  {
    title: '价格展示',
    copy: '公开页、控制台、订单页、分组倍率、模型价格表、活动说明与补充提示，可能共同构成对某次调用或某项权益的价格说明；若显示内容不一致，以实际下单页、支付页和账户账册记录为准。',
  },
  {
    title: '用量统计',
    copy: '模型调用产生的 token、请求次数、缓存读取、图片输出、倍率折算、渠道定价或其他计量口径，可能因模型类型、账务规则、分组策略与站点配置不同而有所差异。',
  },
  {
    title: '订单与支付',
    copy: '当你购买、充值、订阅或处理订单时，支付结果、订单状态、账册入账与权益开通通常需要经过第三方支付服务与站内结算流程确认；支付成功并不当然代表争议已最终解决。',
  },
  {
    title: '退款与争议',
    copy: '除适用法律另有要求外，已消耗的用量、已实际开通并可使用的权益、因上游模型实际调用而产生的成本、以及因你自身原因导致的错误购买，一般不当然构成退款义务。',
  },
  {
    title: '记录优先',
    copy: '如你对价格、扣费、倍率、订单状态或权益期限有异议，应以账户账册、订单记录、时间戳、请求记录、支付状态及站点配置记录为核验基础。',
  },
] as const
const enBillingItems = [
  { title: 'Price display', copy: 'Public pages, console pages, order pages, group rates, model price tables, campaigns, and supplemental notices may together describe prices for a call or entitlement. If displays conflict, actual order pages, payment pages, and ledger records prevail.' },
  { title: 'Usage metering', copy: 'Tokens, request counts, cache reads, image outputs, rate conversion, channel pricing, and other metering dimensions may differ by model type, accounting rules, group strategy, and site configuration.' },
  { title: 'Orders and payment', copy: 'Purchases, recharges, subscriptions, and order handling usually require confirmation through third-party payment services and site settlement flows. Successful payment does not automatically mean every dispute is finally resolved.' },
  { title: 'Refunds and disputes', copy: 'Unless applicable law requires otherwise, consumed usage, opened and usable entitlements, upstream model costs already incurred, and purchases made in error for your own reasons generally do not automatically create refund obligations.' },
  { title: 'Records prevail', copy: 'If you dispute prices, deductions, rates, order status, or entitlement periods, account ledgers, order records, timestamps, request records, payment status, and site configuration records are the basis for verification.' },
]
const billingItems = computed(() => locale.value.startsWith('zh') ? zhBillingItems : enBillingItems)

const zhConductItems = [
  {
    title: '违法或侵权使用',
    copy: '你不得利用本服务从事违反适用法律法规、侵犯知识产权、侵犯隐私、散播恶意内容、欺诈、洗钱、绕过制裁或其他违法违规活动。',
  },
  {
    title: '滥用系统与资源',
    copy: '你不得通过批量脚本、撞库、扫描、压测、漏洞利用、接口泛洪、恶意重试、刷量、薅羊毛、绕过频控或其他方式影响服务稳定性、他人权益或上游资源分配。',
  },
  {
    title: '转售与共享风险',
    copy: '未经明确许可，你不得将账户、API Key、控制台权限、余额、分组权益或已开通能力转借、转售、出租、拼单、代充后倒卖，或以其他方式作为二次分发服务。',
  },
  {
    title: '规避风控与限制',
    copy: '你不得通过多账号、多设备、伪造身份、伪造地区、伪造支付信息、虚构交易、套壳请求、逆向规则或其他方式规避站点风控、资格判断、计费逻辑与使用限制。',
  },
] as const
const enConductItems = [
  { title: 'Unlawful or infringing use', copy: 'You must not use the service for activities that violate applicable law, infringe intellectual property or privacy rights, spread malicious content, commit fraud, launder money, evade sanctions, or otherwise break rules.' },
  { title: 'System and resource abuse', copy: 'You must not affect service stability, others rights, or upstream allocation through bulk scripts, credential stuffing, scanning, load testing, exploit attempts, floods, malicious retries, fake volume, abuse of promotions, or bypassing rate controls.' },
  { title: 'Resale and sharing risk', copy: 'Without explicit permission, you must not lend, resell, rent, share, recharge for resale, or otherwise redistribute accounts, API Keys, console permissions, balance, group entitlements, or opened capabilities.' },
  { title: 'Evasion of controls and limits', copy: 'You must not evade site risk controls, eligibility decisions, billing logic, or use limits through multiple accounts, devices, false identities, false regions, false payment information, fake transactions, wrapped requests, or reverse-engineered rules.' },
]
const conductItems = computed(() => locale.value.startsWith('zh') ? zhConductItems : enConductItems)

const zhThirdPartyItems = [
  {
    title: '上游模型服务',
    copy: '你通过本服务发起的请求，可能会被转发至相应上游模型、通道方或基础设施服务。模型输出质量、响应时延、内容风格、限流策略、拒答行为与临时下线，不由本服务单方控制。',
  },
  {
    title: '第三方支付与账务',
    copy: '支付、退款、扣款风控、结算失败、回调延迟、拒付与争议处理，可能受第三方支付服务商、银行、卡组织或账务处理方规则影响。',
  },
  {
    title: '规则独立适用',
    copy: '第三方模型、支付渠道、邮件服务、验证码、安全服务或云基础设施，均可能有各自独立的使用规则、隐私条款、可用性限制与地区限制；本条款不替代第三方条款。',
  },
] as const
const enThirdPartyItems = [
  { title: 'Upstream model services', copy: 'Requests sent through the service may be forwarded to the relevant upstream model, channel provider, or infrastructure service. Output quality, latency, style, rate limits, refusals, and temporary downtime are not solely controlled by this service.' },
  { title: 'Third-party payment and accounting', copy: 'Payment, refunds, deduction risk controls, settlement failures, callback delays, chargebacks, and dispute handling may be affected by payment providers, banks, card networks, or accounting processors.' },
  { title: 'Independent rules apply', copy: 'Third-party models, payment channels, email services, captcha, security services, or cloud infrastructure may have their own usage rules, privacy terms, availability limits, and regional restrictions. These terms do not replace third-party terms.' },
]
const thirdPartyItems = computed(() => locale.value.startsWith('zh') ? zhThirdPartyItems : enThirdPartyItems)

const zhAvailabilityItems = [
  {
    title: '服务可能变更',
    copy: '我们可以基于运营、上游变化、合规、支付、风控、容量或产品策略，对模型列表、价格结构、分组规则、入口展示、订单流程、支付方式、限速或权限进行调整。',
  },
  {
    title: '不保证持续可用',
    copy: '除适用法律另有明确要求外，本服务不保证任何模型、路由、功能、渠道、价格或能力会长期持续提供，也不保证绝对无中断、无错误、无波动或完全符合你的特定目的。',
  },
  {
    title: '维护与中断',
    copy: '系统维护、风控处置、账务核验、上游波动、网络故障、机房事件、政策调整、支付异常或其他不可控原因，均可能导致部分或全部服务临时不可用。',
  },
] as const
const enAvailabilityItems = [
  { title: 'Service may change', copy: 'We may adjust model lists, pricing structures, group rules, entry displays, order flows, payment methods, rate limits, or permissions due to operations, upstream changes, compliance, payment, risk control, capacity, or product strategy.' },
  { title: 'No continuous availability guarantee', copy: 'Unless applicable law clearly requires otherwise, the service does not guarantee that any model, route, function, channel, price, or capability will remain available long term, or be uninterrupted, error-free, fluctuation-free, or fit for a specific purpose.' },
  { title: 'Maintenance and interruptions', copy: 'System maintenance, risk-control actions, accounting checks, upstream fluctuations, network failures, data-center incidents, policy changes, payment abnormalities, or other uncontrollable causes may make some or all services temporarily unavailable.' },
]
const availabilityItems = computed(() => locale.value.startsWith('zh') ? zhAvailabilityItems : enAvailabilityItems)

const zhTerminationItems = [
  {
    title: '限制措施',
    copy: '如出现未支付订单、异常退款、共享 Key、批量注册、滥用调用、争议未决、投诉举报、风险命中或其他违反本条款的情况，我们可以采取限流、禁用 Key、冻结订单、暂停登录、限制支付或降低权限等措施。',
  },
  {
    title: '服务终止',
    copy: '若违规情节严重、争议无法解决、风控风险过高、第三方明确要求停止服务或适用法律要求终止处理，我们可以中止或终止你对全部或部分服务的访问，而无需继续维持既有能力。',
  },
  {
    title: '终止后的记录',
    copy: '账户被注销、冻结或终止后，订单、账册、日志、争议处理记录、安全审计记录及依法必须保留的信息，仍可能在合理期限内继续保存。',
  },
  {
    title: '责任边界',
    copy: '在适用法律允许的范围内，对因你违反本条款、凭证保管不当、错误配置、误购误用、依赖第三方服务失败或超出合理控制范围的中断与损失，我们不承担超出法律强制要求之外的额外责任。',
  },
] as const
const enTerminationItems = [
  { title: 'Restrictive measures', copy: 'For unpaid orders, abnormal refunds, shared Keys, bulk registration, abusive calls, unresolved disputes, complaints, risk-control hits, or other violations, we may limit rate, disable Keys, freeze orders, suspend sign-in, restrict payment, or reduce permissions.' },
  { title: 'Service termination', copy: 'If violations are serious, disputes cannot be resolved, risk is too high, third parties require discontinuation, or applicable law requires termination, we may suspend or terminate your access to all or part of the service without maintaining previous capabilities.' },
  { title: 'Records after termination', copy: 'After an account is cancelled, frozen, or terminated, orders, ledgers, logs, dispute records, security audit records, and legally required information may still be retained for a reasonable period.' },
  { title: 'Responsibility boundary', copy: 'To the extent allowed by applicable law, we do not assume additional liability beyond mandatory legal requirements for interruptions or losses caused by your breach, poor credential handling, misconfiguration, mistaken purchase or use, third-party service failure, or events beyond reasonable control.' },
]
const terminationItems = computed(() => locale.value.startsWith('zh') ? zhTerminationItems : enTerminationItems)

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.terms-page :deep(main > section:first-child) {
  grid-template-columns: minmax(0, 1fr);
  gap: 0.35rem;
  padding-top: 1.2rem;
  padding-bottom: 0.45rem;
}

.terms-page :deep(.public-copy-block) {
  max-width: 52rem;
}

.terms-page :deep(.public-display) {
  font-size: clamp(2.4rem, 5.3vw, 4.25rem);
  font-weight: 600;
  line-height: 1.02;
  letter-spacing: 0.01em;
}

.terms-page :deep(.public-intro) {
  margin-top: 1.25rem;
  max-width: 47rem;
  font-family: inherit;
  font-size: clamp(0.96rem, 1.02vw, 1.04rem);
  line-height: 1.88;
  color: #5f685c;
}

.terms-page :deep(.public-hero-panel),
.terms-page :deep(.public-cta) {
  display: none;
}

.terms-main-grid {
  max-width: 60rem;
  grid-template-columns: minmax(0, 11.4rem) minmax(0, 1fr);
}

.terms-main-grid > *,
.terms-article {
  min-width: 0;
}

.terms-lead {
  font-size: clamp(0.96rem, 1.02vw, 1.04rem);
  line-height: 1.88;
}

.terms-card,
.terms-notice {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.34);
}

.terms-meta-strip {
  border-top: 1px solid rgba(216, 205, 185, 0.72);
  padding-top: 1rem;
}

.terms-meta-item {
  display: grid;
  gap: 0.35rem;
  align-content: start;
}

.terms-meta-label {
  font-size: 0.72rem;
  line-height: 1.5;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #8f7d63;
}

.terms-meta-value {
  font-size: 0.92rem;
  line-height: 1.68;
  font-weight: 500;
  color: #1f2320;
}

.terms-kicker {
  font-size: 0.7rem;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: #7b6a53;
  font-weight: 500;
}

.terms-toc-link {
  border-left: 1px solid transparent;
  padding: 0.56rem 0 0.56rem 0.7rem;
  text-align: left;
  font-size: 0.95rem;
  line-height: 1.45;
  color: #4a524a;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.terms-toc-link:hover {
  color: #b95d1f;
}

.terms-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.8);
  background: rgba(185, 93, 31, 0.05);
  color: #b95d1f;
}

.terms-title {
  margin-top: 0.85rem;
  width: 100%;
  max-width: none;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.12rem, 1.22vw, 1.34rem);
  line-height: 1.54;
  letter-spacing: 0.01em;
  color: #1f2320;
}

.terms-copy {
  margin-top: 1rem;
  max-width: 52rem;
  font-size: 0.96rem;
  line-height: 1.84;
  color: #5f685c;
}

.terms-card {
  padding: 1.05rem 1.05rem 1.15rem;
}

.terms-card-icon {
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

.terms-card-title,
.terms-notice-title {
  color: #1f2320;
  font-size: 0.98rem;
  font-weight: 600;
  line-height: 1.45;
}

.terms-card-copy,
.terms-notice-copy {
  margin-top: 0.55rem;
  font-size: 0.93rem;
  line-height: 1.8;
  color: #5f685c;
}

.terms-notice {
  display: flex;
  align-items: flex-start;
  gap: 0.9rem;
  padding: 1rem 1rem 1.05rem;
  background: linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58));
}

:global(html.dark) .terms-meta-strip {
  border-top-color: rgba(82, 87, 76, 0.84);
}

:global(html.dark) .terms-meta-label {
  color: #b8af9a;
}

:global(html.dark) .terms-meta-value {
  color: #efe5d2;
}

:global(html.dark) .terms-table > :first-child span {
  color: #b8af9a;
}

:global(html.dark) .terms-table-row > :first-child {
  color: #efe5d2;
}

:global(html.dark) .terms-card,
:global(html.dark) .terms-notice {
  border-color: rgba(82, 87, 76, 0.84);
  background: linear-gradient(180deg, rgba(20, 22, 18, 0.92), rgba(14, 16, 14, 0.96));
}

:global(html.dark) .terms-title,
:global(html.dark) .terms-card-title,
:global(html.dark) .terms-notice-title {
  color: #efe5d2;
}

:global(html.dark) .terms-copy,
:global(html.dark) .terms-card-copy,
:global(html.dark) .terms-notice-copy {
  color: #cdbfab;
}

:global(html.dark) .terms-kicker {
  color: #b8af9a;
}

:global(html.dark) .terms-toc-link {
  color: #d8cfbe;
}

:global(html.dark) .terms-toc-link:hover,
:global(html.dark) .terms-toc-link-active {
  color: #ffd8bb;
}

:global(html.dark) .terms-toc-link-active {
  border-left-color: rgba(255, 216, 187, 0.72);
  background: rgba(185, 93, 31, 0.1);
}

:global(html.dark) .terms-card-icon {
  background: rgba(188, 93, 31, 0.16);
  color: #ffd8bb;
}

:global(html.dark) .terms-notice {
  border-color: rgba(138, 116, 90, 0.38);
  background:
    linear-gradient(145deg, rgba(29, 32, 26, 0.9), rgba(24, 22, 19, 0.9)),
    radial-gradient(circle at 84% 34%, rgba(173, 89, 36, 0.12), transparent 28%);
  box-shadow:
    inset 0 1px 0 rgba(255, 241, 220, 0.06),
    0 16px 32px rgba(0, 0, 0, 0.28);
}

:global(html.dark) .terms-notice-title {
  color: #f2e8d5;
}

:global(html.dark) .terms-notice-copy {
  color: #d5cab9;
}

@media (max-width: 1023px) {
  .terms-main-grid {
    grid-template-columns: 1fr;
  }

  .terms-toc-link {
    padding-left: 0.65rem;
    font-size: 0.94rem;
  }

  .terms-title {
    font-size: clamp(1.1rem, 3.5vw, 1.28rem);
    line-height: 1.45;
  }

  .terms-meta-value,
  .terms-copy,
  .terms-card-copy,
  .terms-notice-copy,
  .terms-lead {
    font-size: 0.95rem;
  }

  .terms-table > :first-child,
  .terms-table-row {
    grid-template-columns: 1fr;
  }

  .terms-table > :first-child span,
  .terms-table-row > * {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .terms-table > :first-child span:last-child {
    padding-top: 0;
  }

  .terms-table-row > :first-child {
    padding-bottom: 0.35rem;
  }

  .terms-table-row > :last-child {
    padding-top: 0;
  }
}
</style>

<style>
html.dark .terms-page .public-intro {
  color: #ead9bd;
}

html.dark .terms-page .public-display {
  color: #fff4dd;
  text-shadow: 0 1px 0 rgba(255, 240, 218, 0.08), 0 18px 48px rgba(0, 0, 0, 0.28);
}

html.dark .terms-page .public-copy-block > div:first-child span:last-child,
html.dark .terms-page .terms-meta-label,
html.dark .terms-page .terms-kicker,
html.dark .terms-page .terms-table > :first-child span,
html.dark .terms-page .terms-main-grid > aside > div > div:first-child {
  color: #cdb387;
}

html.dark .terms-page .terms-meta-value,
html.dark .terms-page .terms-title,
html.dark .terms-page .terms-card-title,
html.dark .terms-page .terms-notice-title,
html.dark .terms-page .terms-table-row > :first-child,
html.dark .terms-page .terms-main-grid li strong {
  color: #fff0d5;
}

html.dark .terms-page .terms-toc-link:hover,
html.dark .terms-page .terms-toc-link-active {
  color: #f3c786;
}

html.dark .terms-page .terms-toc-link-active {
  border-left-color: rgba(212, 153, 80, 0.8);
  background: linear-gradient(90deg, rgba(176, 120, 57, 0.16), rgba(176, 120, 57, 0.04));
}

html.dark .terms-page .terms-card,
html.dark .terms-page .terms-notice,
html.dark .terms-page .terms-table-row {
  border-color: rgba(141, 109, 72, 0.56) !important;
}

html.dark .terms-article {
  border-color: rgba(120, 109, 90, 0.54) !important;
  background:
    linear-gradient(180deg, rgba(35, 37, 31, 0.95), rgba(24, 26, 21, 0.96)),
    radial-gradient(circle at top right, rgba(166, 97, 45, 0.1), transparent 28%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 244, 224, 0.07),
    0 18px 40px rgba(0, 0, 0, 0.2) !important;
}

html.dark .terms-lead,
html.dark .terms-copy,
html.dark .terms-card-copy,
html.dark .terms-notice-copy,
html.dark .terms-main-grid li,
html.dark .terms-table-row p {
  color: #e2d6c4;
}

html.dark .terms-meta-value,
html.dark .terms-title,
html.dark .terms-card-title,
html.dark .terms-notice-title,
html.dark .terms-table-row > :first-child,
html.dark .terms-copy strong,
html.dark .terms-copy b,
html.dark .terms-copy em,
html.dark .terms-copy a,
html.dark .terms-card-copy strong,
html.dark .terms-notice-copy strong,
html.dark .terms-main-grid li strong {
  color: #fff0da;
}

html.dark .terms-meta-label,
html.dark .terms-kicker,
html.dark .terms-table > :first-child span,
html.dark .terms-main-grid > aside > div > div:first-child {
  color: #caba9f;
}

html.dark .terms-meta-strip {
  border-top-color: rgba(102, 95, 79, 0.68);
}

html.dark .terms-toc-link {
  color: #f0e3cc;
}

html.dark .terms-toc-link:hover,
html.dark .terms-toc-link-active {
  color: #ffe0bf;
}

html.dark .terms-card,
html.dark .terms-page .terms-card,
html.dark .terms-main-grid .terms-card,
html.dark .terms-page .terms-main-grid .terms-card {
  border-color: rgba(118, 106, 87, 0.58) !important;
  background:
    linear-gradient(180deg, rgba(36, 38, 31, 0.94), rgba(27, 29, 24, 0.95)),
    radial-gradient(circle at 88% 18%, rgba(163, 97, 45, 0.09), transparent 24%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 242, 219, 0.06),
    0 14px 28px rgba(0, 0, 0, 0.18) !important;
}

html.dark .terms-card-title,
html.dark .terms-page .terms-card-title,
html.dark .terms-main-grid .terms-card-title,
html.dark .terms-page .terms-main-grid .terms-card-title {
  color: #f4ead7 !important;
}

html.dark .terms-card-copy,
html.dark .terms-page .terms-card-copy,
html.dark .terms-main-grid .terms-card-copy,
html.dark .terms-page .terms-main-grid .terms-card-copy {
  color: #d9cdbb !important;
}

html.dark .terms-notice,
html.dark .terms-page .terms-notice,
html.dark .terms-main-grid .terms-notice,
html.dark .terms-page .terms-main-grid .terms-notice {
  border-color: rgba(136, 115, 88, 0.52) !important;
  background:
    linear-gradient(145deg, rgba(42, 38, 31, 0.95), rgba(30, 27, 23, 0.96)),
    radial-gradient(circle at 84% 34%, rgba(173, 89, 36, 0.12), transparent 30%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 241, 220, 0.07),
    0 16px 32px rgba(0, 0, 0, 0.24) !important;
}

html.dark .terms-notice-title,
html.dark .terms-page .terms-notice-title,
html.dark .terms-main-grid .terms-notice-title,
html.dark .terms-page .terms-main-grid .terms-notice-title {
  color: #f7ebd8 !important;
}

html.dark .terms-notice-copy,
html.dark .terms-page .terms-notice-copy,
html.dark .terms-main-grid .terms-notice-copy,
html.dark .terms-page .terms-main-grid .terms-notice-copy {
  color: #dfd2be !important;
}

html.dark .terms-table > :first-child {
  background: rgba(35, 32, 27, 0.92) !important;
}

html.dark .terms-table-row {
  background: rgba(29, 31, 26, 0.72) !important;
  border-top-color: rgba(96, 89, 75, 0.58) !important;
}

html.dark .terms-article {
  border-color: rgba(155, 126, 86, 0.26) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .terms-page .terms-card,
html.dark .terms-page .terms-notice,
html.dark .terms-page .terms-table,
html.dark .terms-page .terms-table-row {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.88), rgba(14, 16, 13, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.08), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 238, 210, 0.05) !important;
}

html.dark .terms-page .terms-notice {
  background:
    linear-gradient(135deg, rgba(35, 29, 23, 0.88), rgba(24, 27, 22, 0.84)),
    radial-gradient(circle at 84% 18%, rgba(194, 126, 74, 0.13), transparent 28%) !important;
}

html.dark .terms-page .terms-table > :first-child {
  background:
    linear-gradient(180deg, rgba(39, 32, 26, 0.9), rgba(24, 27, 22, 0.88)) !important;
}

html.dark .terms-page .terms-title,
html.dark .terms-page .terms-card-title,
html.dark .terms-page .terms-notice-title,
html.dark .terms-page .terms-table-row > :first-child,
html.dark .terms-page .terms-main-grid li strong {
  color: #f6e8d2 !important;
}

html.dark .terms-page .terms-copy,
html.dark .terms-page .terms-card-copy,
html.dark .terms-page .terms-notice-copy,
html.dark .terms-page .terms-main-grid li,
html.dark .terms-page .terms-table-row p {
  color: #d0baa0 !important;
}

html.dark .terms-page .terms-kicker,
html.dark .terms-page .terms-table > :first-child span,
html.dark .terms-page .terms-main-grid > aside > div > div:first-child,
html.dark .terms-page .public-copy-block > div:first-child span:last-child {
  color: #d8b171 !important;
}

html.dark .terms-page .terms-toc-link {
  color: #d4c4ad !important;
}

html.dark .terms-page .terms-toc-link:hover,
html.dark .terms-page .terms-toc-link-active {
  color: #efc183 !important;
}

html.dark .terms-page .terms-toc-link-active {
  border-left-color: rgba(194, 126, 74, 0.72) !important;
  background: linear-gradient(90deg, rgba(194, 126, 74, 0.16), rgba(194, 126, 74, 0.04)) !important;
}

html.dark .terms-page .terms-card-icon {
  background: rgba(194, 126, 74, 0.13) !important;
  color: #efc183 !important;
}

html:not(.dark) .terms-article {
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

html:not(.dark) .terms-page .terms-card,
html:not(.dark) .terms-page .terms-notice,
html:not(.dark) .terms-page .terms-table,
html:not(.dark) .terms-page .terms-table-row {
  border-color: rgba(190, 168, 134, 0.42) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(244, 235, 220, 0.34)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.06), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.42) !important;
}

html:not(.dark) .terms-page .terms-notice {
  background:
    linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58)),
    radial-gradient(circle at 84% 18%, rgba(196, 136, 68, 0.08), transparent 28%) !important;
}

html:not(.dark) .terms-page .terms-table > :first-child {
  background: rgba(255, 252, 246, 0.48) !important;
}

html:not(.dark) .terms-page .terms-toc-link {
  color: #4a524a !important;
}

html:not(.dark) .terms-page .terms-toc-link:hover,
html:not(.dark) .terms-page .terms-toc-link-active {
  color: #b95d1f !important;
}

html:not(.dark) .terms-page .terms-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.74) !important;
  background: linear-gradient(90deg, rgba(185, 93, 31, 0.07), rgba(185, 93, 31, 0.02)) !important;
}

html:not(.dark) .terms-page .terms-card-icon {
  background: rgba(188, 93, 31, 0.1) !important;
  color: #b95d1f !important;
}
</style>

