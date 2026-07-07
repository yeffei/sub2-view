<template>
  <PublicPageLayout
    class="faq-page"
    tone="faq"
    eyebrow="问答"
    title="常见问题"
    intro="入庭之前，先把边界问清楚。"
    :highlights="['开通边界', '接入口径', '账册核对', '凭据安全']"
  >
    <template #aside>
      <div class="faq-aside">
        <div>
          <div class="faq-aside-kicker">SST</div>
          <div class="faq-aside-title">统一入口，安静流转。</div>
          <p class="faq-aside-copy">
            账户、分组、模型与额度，以控制台当前状态为准。
          </p>
        </div>

        <div class="faq-aside-rules" aria-label="问答范围">
          <div v-for="item in asideRules" :key="item.title" class="faq-aside-rule">
            <Icon :name="item.icon" size="sm" />
            <span>{{ item.title }}</span>
          </div>
        </div>

        <div class="faq-aside-links">
          <RouterLink v-for="item in quickLinks" :key="item.to" :to="item.to" class="faq-aside-link">
            <span>{{ item.label }}</span>
            <Icon name="arrowRight" size="sm" />
          </RouterLink>
        </div>

        <div v-if="publicContact" id="public-contact" class="faq-contact">
          <a
            v-if="publicContact.href"
            class="faq-contact-link"
            :href="publicContact.href"
            :target="publicContact.external ? '_blank' : undefined"
            :rel="publicContact.external ? 'noopener noreferrer' : undefined"
          >
            <span>联系庭务</span>
            <strong>{{ publicContact.label }}</strong>
          </a>
          <div v-else class="faq-contact-link" role="note">
            <span>联系庭务</span>
            <strong>{{ publicContact.label }}</strong>
          </div>
        </div>
      </div>
    </template>

    <section class="faq-overview" aria-label="核心问题">
      <article v-for="item in overviewItems" :key="item.title" class="faq-overview-item">
        <div class="faq-overview-index">{{ item.index }}</div>
        <h2>{{ item.title }}</h2>
        <p>{{ item.copy }}</p>
      </article>
    </section>

    <section class="faq-main-grid">
      <aside class="faq-toc" aria-label="问答目录">
        <div class="faq-toc-inner">
          <div class="faq-toc-kicker">目录</div>
          <button
            v-for="group in faqGroups"
            :key="group.id"
            type="button"
            class="faq-toc-link"
            :class="{ 'is-active': activeGroupId === group.id }"
            :aria-pressed="activeGroupId === group.id"
            @click="selectGroup(group.id)"
          >
            <span class="faq-toc-index">{{ group.index }}</span>
            <span>{{ group.title }}</span>
          </button>
        </div>
      </aside>

      <div class="faq-groups">
        <Transition name="faq-panel" mode="out-in">
          <section :id="activeFaqGroup.id" :key="activeFaqGroup.id" class="faq-group">
            <div class="faq-group-heading">
              <span>{{ activeFaqGroup.index }}</span>
              <h2>{{ activeFaqGroup.title }}</h2>
            </div>

            <div class="faq-question-list">
              <details v-for="item in activeFaqGroup.items" :key="item.question" class="faq-question">
                <summary>
                  <span>{{ item.question }}</span>
                  <Icon name="chevronDown" size="sm" />
                </summary>
                <p>{{ item.answer }}</p>
              </details>
            </div>
          </section>
        </Transition>
      </div>
    </section>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'
import { useAppStore } from '@/stores'
import { resolvePublicContact } from '@/utils/contact'

const appStore = useAppStore()
const publicContact = computed(() => resolvePublicContact(appStore.cachedPublicSettings?.contact_info || appStore.contactInfo))

const asideRules = [
  { title: '价格看价目', icon: 'dollar' },
  { title: '接入看文档', icon: 'book' },
  { title: '状态看控制台', icon: 'grid' },
] as const

const quickLinks = [
  { to: '/pricing', label: '查看价目' },
  { to: '/docs', label: '接入文档' },
  { to: '/login', label: '进入控制台' },
] as const

const overviewItems = [
  {
    index: '甲',
    title: '开通之前',
    copy: '先确认套餐、倍率、模型与分组，再决定接入方式。',
  },
  {
    index: '乙',
    title: '接入之后',
    copy: 'Key、base_url、模型列表与调用明细在同一入口中核对。',
  },
  {
    index: '丙',
    title: '遇到波动',
    copy: '保留时间、模型、错误码与 request id，便于定位线路或权限问题。',
  },
] as const

const faqGroups = [
  {
    id: 'entry',
    index: '一',
    title: '开通与准入',
    items: [
      {
        question: '山枢庭适合什么使用方式？',
        answer: '更适合长期、稳定、可核对的 API 使用场景，例如开发工具、自动化任务、团队内部服务和需要统一账册的调用链路。',
      },
      {
        question: '注册后是否一定可以立即调用所有模型？',
        answer: '不一定。可用模型与权限会受账户状态、分组、套餐和后台配置影响；最终以控制台显示的模型列表与实际接口返回为准。',
      },
      {
        question: '开通前最先确认什么？',
        answer: '先确认可用套餐、倍率、模型范围、分组权限和付款方式；如果需要团队使用，再确认 Key 管理、账册核对和额度提醒。',
      },
    ],
  },
  {
    id: 'access',
    index: '二',
    title: '接入与调用',
    items: [
      {
        question: '现有 OpenAI SDK 能不能直接接入？',
        answer: '通常可以。把 SDK 的 API Key 换成山枢庭控制台生成的 Key，并把 base_url 或 baseURL 配到文档指定的 /v1 地址即可。',
      },
      {
        question: '模型名应该从哪里确认？',
        answer: '优先调用模型列表接口或查看控制台可用模型，不建议凭记忆填写。不同分组、账户和上游能力可能返回不同模型范围。',
      },
      {
        question: '流式输出是否支持？',
        answer: '支持 OpenAI 兼容的 stream 调用方式。客户端按 SSE 事件处理逐段返回，具体行为仍以所选模型和上游能力为准。',
      },
    ],
  },
  {
    id: 'billing',
    index: '三',
    title: '计量与账册',
    items: [
      {
        question: '余额、用量和订单在哪里看？',
        answer: '登录后进入控制台查看。余额、订单、调用记录和用量明细会归到同一套账册中，便于按时间和模型核对。',
      },
      {
        question: '价格页上的倍率和实际扣费是什么关系？',
        answer: '价格页用于说明公开计量口径；实际扣费还会结合账户分组、套餐、模型价格和后台配置，最终以账册记录为准。',
      },
      {
        question: '如果发现用量和预期不一致怎么办？',
        answer: '先按时间、Key、模型和请求类型筛选调用记录；如果仍有疑问，保留订单号、请求时间和相关记录再联系管理员核对。',
      },
    ],
  },
  {
    id: 'stability',
    index: '四',
    title: '稳定性与排查',
    items: [
      {
        question: '上游波动时会发生什么？',
        answer: '系统会尽量通过路由、故障转移和可用账号调度维持入口连续性；但第三方上游异常、限流或权限变化仍可能影响单次请求。',
      },
      {
        question: '接口返回 401、403 或 429 时先查什么？',
        answer: '401 先查 Key 和 Bearer 格式；403 先查模型权限与分组范围；429 先查请求频率、余额、套餐限制和上游限流窗口。',
      },
      {
        question: '反馈问题时需要提供哪些信息？',
        answer: '建议提供请求时间、模型名、错误码、request id、使用的 Key 前缀和简要调用场景，不要发送完整 API Key 或敏感请求内容。',
      },
    ],
  },
  {
    id: 'security',
    index: '五',
    title: '凭据与边界',
    items: [
      {
        question: 'API Key 应该放在哪里？',
        answer: '应放在服务端环境变量或后端配置中，不要写入浏览器前端、移动端包、公开仓库、日志截图或可分享的配置文件。',
      },
      {
        question: 'Key 泄露后应该怎么处理？',
        answer: '立即停用或删除旧 Key，重新创建新 Key，并检查近期调用记录、余额变化和异常请求来源。',
      },
      {
        question: '请求内容会被当作长期资料保存吗？',
        answer: '山枢庭作为统一入口处理转发、计量和必要审计信息。具体保存范围以隐私政策、服务条款和后台实际配置为准。',
      },
    ],
  },
] as const

type FaqGroupId = (typeof faqGroups)[number]['id']

const activeGroupId = ref<FaqGroupId>('entry')
const activeFaqGroup = computed(() => faqGroups.find(group => group.id === activeGroupId.value) ?? faqGroups[0])

function selectGroup(id: FaqGroupId) {
  activeGroupId.value = id
}
</script>

<style scoped>
.faq-page :deep(main > section:first-child) {
  padding-top: 1.4rem;
  padding-bottom: 1.2rem;
}

.faq-page :deep(.public-display) {
  font-size: clamp(2.8rem, 6.2vw, 5.1rem);
}

.faq-page :deep(.public-intro) {
  max-width: 34rem;
}

.faq-page :deep(.public-hero-panel) {
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.84), rgba(244, 235, 220, 0.64)),
    radial-gradient(circle at 88% 10%, rgba(183, 129, 70, 0.08), transparent 22%),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.025) 33px 34px);
  box-shadow:
    0 18px 48px rgba(77, 59, 37, 0.07),
    inset 0 1px 0 rgba(255, 251, 244, 0.72),
    inset 0 -1px 0 rgba(139, 107, 68, 0.08);
}

.faq-aside {
  display: grid;
  gap: 1.45rem;
}

.faq-aside-kicker,
.faq-toc-kicker {
  color: #9b7a52;
  font-size: 0.7rem;
  letter-spacing: 0.32em;
  text-transform: uppercase;
}

.faq-aside-title {
  margin-top: 0.78rem;
  color: #2f281f;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.7rem, 2.1vw, 2.3rem);
  line-height: 1.18;
}

.faq-aside-copy {
  margin-top: 0.9rem;
  max-width: 24rem;
  color: #5f685c;
  font-size: 0.95rem;
  line-height: 1.85;
}

.faq-aside-rules {
  display: grid;
  gap: 0.7rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(154, 128, 92, 0.16);
}

.faq-aside-rule {
  display: flex;
  align-items: center;
  gap: 0.62rem;
  color: #3f473f;
  font-size: 0.94rem;
  line-height: 1.5;
}

.faq-aside-rule svg {
  color: #a73a2a;
}

.faq-aside-links {
  display: grid;
  gap: 0.62rem;
}

.faq-aside-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.8rem;
  border: 1px solid rgba(190, 168, 134, 0.42);
  border-radius: 1rem;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(244, 235, 220, 0.34)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.06), transparent 26%);
  padding: 0.82rem 0.95rem;
  color: #2f281f;
  font-size: 0.94rem;
  transition: border-color 160ms ease, background-color 160ms ease, color 160ms ease, transform 160ms ease;
}

.faq-aside-link:hover {
  border-color: rgba(196, 136, 68, 0.32);
  color: #9c4c26;
  transform: translateY(-1px);
}

.faq-contact {
  border-top: 1px solid rgba(154, 128, 92, 0.16);
  padding-top: 0.95rem;
}

.faq-contact-link {
  display: grid;
  gap: 0.36rem;
  color: #4f5750;
  font-size: 0.9rem;
  line-height: 1.5;
}

.faq-contact-link span {
  color: #9b7a52;
  font-size: 0.72rem;
  letter-spacing: 0.28em;
}

.faq-contact-link strong {
  color: #2f281f;
  font-weight: 500;
  overflow-wrap: anywhere;
}

a.faq-contact-link:hover strong {
  color: #a73a2a;
}

.faq-overview {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0;
  margin-top: 0.5rem;
  border-top: 1px solid rgba(154, 128, 92, 0.16);
  border-bottom: 1px solid rgba(154, 128, 92, 0.16);
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.16), transparent 42%),
    linear-gradient(90deg, rgba(144, 113, 76, 0.026), transparent 30%, transparent 70%, rgba(144, 113, 76, 0.018));
}

.faq-overview-item {
  min-width: 0;
  padding: 1.35rem 1.25rem;
}

.faq-overview-item + .faq-overview-item {
  border-left: 1px solid rgba(154, 128, 92, 0.14);
}

.faq-overview-index {
  color: #af7840;
  font-family: 'Geist Mono', 'JetBrains Mono', monospace;
  font-size: 0.68rem;
  letter-spacing: 0.12em;
}

.faq-overview-item h2 {
  margin-top: 0.56rem;
  color: #2a241e;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.16rem, 1.35vw, 1.45rem);
  line-height: 1.2;
}

.faq-overview-item p {
  margin-top: 0.52rem;
  color: #5f685c;
  font-size: 0.9rem;
  line-height: 1.78;
}

.faq-main-grid {
  display: grid;
  grid-template-columns: minmax(0, 11.4rem) minmax(0, 1fr);
  gap: 2rem;
  margin-top: 2rem;
  max-width: 62rem;
}

.faq-main-grid > * {
  min-width: 0;
}

.faq-toc-inner {
  position: sticky;
  top: 1.25rem;
  display: grid;
  gap: 0.2rem;
}

.faq-toc-kicker {
  margin-bottom: 0.45rem;
}

.faq-toc-link {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.58rem;
  border-left: 1px solid transparent;
  background: transparent;
  padding: 0.56rem 0 0.56rem 0.7rem;
  text-align: left;
  color: #4a524a;
  font: inherit;
  font-size: 0.95rem;
  line-height: 1.45;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease, transform 160ms ease;
}

.faq-toc-link:hover,
.faq-toc-link.is-active {
  border-left-color: rgba(185, 93, 31, 0.56);
  background: rgba(185, 93, 31, 0.045);
  color: #b95d1f;
}

.faq-toc-link.is-active {
  transform: translateX(2px);
}

.faq-toc-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  flex: 0 0 auto;
  border-radius: 0.28rem;
  background: rgba(185, 93, 31, 0.075);
  color: #9c4c26;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.72rem;
}

.faq-groups {
  display: grid;
  gap: 1rem;
}

.faq-group {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  min-height: 17.5rem;
  border: 1px solid rgba(154, 128, 92, 0.16);
  border-radius: 1.18rem;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    linear-gradient(90deg, rgba(144, 113, 76, 0.038), transparent 18%, rgba(144, 113, 76, 0.024) 82%, transparent),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28);
  box-shadow:
    0 14px 34px rgba(84, 57, 31, 0.05),
    inset 0 1px 0 rgba(255, 249, 239, 0.6),
    inset 0 -1px 0 rgba(140, 111, 76, 0.07),
    inset 0 0 0 1px rgba(255, 255, 255, 0.22);
}

.faq-group::before {
  content: '';
  position: absolute;
  inset: 0;
  z-index: -1;
  pointer-events: none;
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.2), transparent 28%),
    radial-gradient(circle at 92% 8%, rgba(167, 58, 42, 0.055), transparent 22%),
    linear-gradient(180deg, rgba(255, 252, 246, 0.18), transparent 42%);
}

.faq-group-heading {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  padding: 1.05rem 1.18rem;
  border-bottom: 1px solid rgba(154, 128, 92, 0.14);
}

.faq-group-heading span {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.7rem;
  height: 1.7rem;
  border-radius: 0.38rem;
  background: linear-gradient(135deg, rgba(186, 72, 54, 0.96), rgba(145, 39, 24, 0.96));
  color: #fff8ef;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.88rem;
  box-shadow: 0 0 0 4px rgba(167, 58, 42, 0.045);
}

.faq-group-heading h2 {
  color: #2f281f;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.12rem, 1.3vw, 1.38rem);
  line-height: 1.24;
}

.faq-question-list {
  display: grid;
}

.faq-question + .faq-question {
  border-top: 1px solid rgba(154, 128, 92, 0.14);
}

.faq-question summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  cursor: pointer;
  list-style: none;
  padding: 1.05rem 1.18rem;
  color: #2d261e;
  font-size: 1rem;
  line-height: 1.55;
  transition: color 160ms ease, background-color 160ms ease;
}

.faq-question summary::-webkit-details-marker {
  display: none;
}

.faq-question summary:hover {
  background: rgba(255, 252, 246, 0.34);
  color: #9c4c26;
}

.faq-question summary svg {
  flex: 0 0 auto;
  color: #a07a49;
  transition: transform 180ms ease;
}

.faq-question[open] summary svg {
  transform: rotate(180deg);
}

.faq-question p {
  margin: -0.2rem 1.18rem 1.1rem;
  max-width: 44rem;
  color: #5f685c;
  font-size: 0.95rem;
  line-height: 1.88;
}

.faq-panel-enter-active,
.faq-panel-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}

.faq-panel-enter-from,
.faq-panel-leave-to {
  opacity: 0;
  transform: translateY(6px);
}

@media (max-width: 1023px) {
  .faq-page :deep(main > section:first-child) {
    padding-top: 1.05rem;
  }

  .faq-overview {
    grid-template-columns: 1fr;
  }

  .faq-overview-item + .faq-overview-item {
    border-top: 1px solid rgba(154, 128, 92, 0.14);
    border-left: 0;
  }

  .faq-main-grid {
    grid-template-columns: 1fr;
    gap: 1.1rem;
  }

  .faq-toc-inner {
    position: static;
    display: flex;
    gap: 0.6rem;
    overflow-x: auto;
    padding-bottom: 0.15rem;
    scrollbar-width: none;
  }

  .faq-toc-inner::-webkit-scrollbar {
    display: none;
  }

  .faq-toc-kicker {
    display: none;
  }

  .faq-toc-link {
    flex: 0 0 auto;
    width: auto;
    border: 1px solid rgba(190, 168, 134, 0.42);
    border-radius: 999px;
    background: rgba(255, 252, 246, 0.34);
    padding: 0.58rem 0.88rem;
    white-space: nowrap;
  }

  .faq-toc-link.is-active {
    border-color: rgba(185, 93, 31, 0.34);
    background: rgba(185, 93, 31, 0.08);
    transform: none;
  }
}

@media (max-width: 639px) {
  .faq-page :deep(.public-display) {
    font-size: clamp(2.55rem, 13vw, 3.5rem);
  }

  .faq-page :deep(.public-intro) {
    font-size: 1.42rem;
    line-height: 1.24;
  }

  .faq-aside-link,
  .faq-question summary {
    min-height: 3rem;
  }

  .faq-group-heading,
  .faq-question summary {
    padding-left: 0.95rem;
    padding-right: 0.95rem;
  }

  .faq-question p {
    margin-left: 0.95rem;
    margin-right: 0.95rem;
    font-size: 0.92rem;
    line-height: 1.78;
  }
}
</style>

<style>
html.dark .faq-page .faq-aside-kicker,
html.dark .faq-page .faq-toc-kicker,
html.dark .faq-page .faq-toc-index,
html.dark .faq-page .faq-overview-index,
html.dark .faq-page .public-copy-block > div:first-child span:last-child {
  color: #d8b171 !important;
}

html.dark .faq-page .faq-aside-title,
html.dark .faq-page .faq-overview-item h2,
html.dark .faq-page .faq-group-heading h2,
html.dark .faq-page .faq-question summary {
  color: #f6e8d2 !important;
}

html.dark .faq-page .faq-aside-copy,
html.dark .faq-page .faq-overview-item p,
html.dark .faq-page .faq-question p,
html.dark .faq-page .faq-aside-rule {
  color: #d0baa0 !important;
}

html.dark .faq-page .faq-aside-rules,
html.dark .faq-page .faq-contact,
html.dark .faq-page .faq-overview,
html.dark .faq-page .faq-overview-item + .faq-overview-item,
html.dark .faq-page .faq-group-heading,
html.dark .faq-page .faq-question + .faq-question {
  border-color: rgba(155, 126, 86, 0.22) !important;
}

html.dark .faq-page .faq-overview {
  background:
    linear-gradient(180deg, rgba(255, 226, 184, 0.035), transparent 44%),
    linear-gradient(90deg, rgba(206, 151, 82, 0.025), transparent 28%, transparent 72%, rgba(206, 151, 82, 0.02)) !important;
}

html.dark .faq-page .public-hero-panel {
  border-color: rgba(155, 126, 86, 0.28) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.92), rgba(34, 29, 23, 0.82)),
    radial-gradient(circle at 86% 12%, rgba(194, 126, 74, 0.13), transparent 24%),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.028) 33px 34px) !important;
  box-shadow:
    0 26px 56px rgba(0, 0, 0, 0.28),
    inset 0 1px 0 rgba(245, 225, 194, 0.06),
    inset 0 -1px 0 rgba(194, 126, 74, 0.08) !important;
}

html.dark .faq-page .faq-group {
  border-color: rgba(155, 126, 86, 0.26) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .faq-page .faq-group::before {
  background:
    linear-gradient(90deg, rgba(255, 238, 210, 0.035), transparent 30%),
    radial-gradient(circle at 92% 8%, rgba(194, 126, 74, 0.12), transparent 24%),
    linear-gradient(180deg, rgba(255, 226, 184, 0.032), transparent 46%) !important;
}

html.dark .faq-page .faq-aside-link,
html.dark .faq-page .faq-toc-link {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.88), rgba(14, 16, 13, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.08), transparent 26%) !important;
  color: #d4c4ad !important;
}

html.dark .faq-page .faq-contact-link {
  color: #d0baa0 !important;
}

html.dark .faq-page .faq-contact-link span {
  color: #d8b171 !important;
}

html.dark .faq-page .faq-contact-link strong {
  color: #f3e1c7 !important;
}

html.dark .faq-page a.faq-contact-link:hover strong {
  color: #efc183 !important;
}

html.dark .faq-page .faq-toc-index {
  background: rgba(194, 126, 74, 0.12) !important;
}

html.dark .faq-page .faq-aside-link:hover,
html.dark .faq-page .faq-toc-link:hover,
html.dark .faq-page .faq-toc-link.is-active,
html.dark .faq-page .faq-question summary:hover {
  color: #efc183 !important;
}

html.dark .faq-page .faq-toc-link.is-active {
  border-color: rgba(194, 126, 74, 0.5) !important;
  background:
    linear-gradient(90deg, rgba(194, 126, 74, 0.18), rgba(194, 126, 74, 0.05)),
    linear-gradient(180deg, rgba(24, 27, 22, 0.9), rgba(14, 16, 13, 0.94)) !important;
}

html.dark .faq-page .faq-question summary:hover {
  background: rgba(194, 126, 74, 0.08) !important;
}

html.dark .faq-page .faq-aside-rule svg,
html.dark .faq-page .faq-question summary svg {
  color: #efab69 !important;
}
</style>
