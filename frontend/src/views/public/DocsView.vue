<template>
  <PublicPageLayout
    class="docs-page"
    tone="docs"
    eyebrow="文档"
    title="API 文档"
    intro=""
    description=""
    :show-cta="false"
  >
      <section class="docs-main-grid mt-10 grid gap-8 lg:gap-8">
        <aside class="self-start">
          <div class="sticky top-5 h-fit">
            <div class="mb-3 flex items-end justify-between gap-3">
              <div class="text-xs uppercase tracking-[0.24em] text-zen-mist dark:text-zen-stone">目录</div>
              <div class="docs-toc-caption">按常用顺序阅读</div>
            </div>
            <nav class="grid gap-1 text-sm text-zen-ink dark:text-zen-paper">
              <button
                v-for="item in sections"
                :key="item.id"
                type="button"
                class="docs-toc-link"
                :class="activeSection === item.id ? 'docs-toc-link-active' : ''"
                @click="activeSection = item.id"
              >
                {{ item.label }}
              </button>
            </nav>
          </div>
        </aside>

        <article class="docs-article rounded-[1.35rem] border border-zen-paperLine/70 bg-white/62 p-5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/76 sm:p-6 lg:p-6">
          <section v-show="activeSection === 'quickstart'" id="quickstart">
            <div class="docs-kicker">快速开始</div>
            <h2 class="docs-title">按 OpenAI 兼容方式配置 Key 与基础地址后，即可直接开始调用。</h2>
            <ol class="docs-list docs-quickstart-list mt-5">
              <li>登录控制台并创建 API Key。</li>
              <li>把 SDK 的 <code>base_url</code> 或 <code>baseURL</code> 设置为 <code>{{ sdkBaseUrl }}</code>，并带上 <code>Authorization: Bearer YOUR_API_KEY</code>。</li>
              <li>先调用 <code>GET /v1/models</code> 确认可用模型，再发起对话请求。</li>
            </ol>
            <div class="docs-notice mt-6">
              <div>
                <div class="docs-notice-title">接入提示</div>
                <p class="docs-notice-copy">大多数客户端只需要改 <code>api_key</code> 与 <code>base_url</code>，模型名以模型列表接口返回值为准。</p>
              </div>
            </div>
          </section>

          <section v-show="activeSection === 'authentication'" id="authentication">
            <div class="docs-kicker">认证方式</div>
            <h2 class="docs-title">所有接口请求都需要在 Header 中携带有效的 API Key。</h2>
            <p class="docs-copy">Key 属于账户凭据，不要直接暴露在前端公开代码、客户端 App 或公开仓库中。推荐将 Key 保存在服务端环境变量中，由后端代理完成调用。</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>Authorization: Bearer YOUR_API_KEY</code></pre>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article class="docs-card">
                <div class="docs-card-title">推荐做法</div>
                <p class="docs-card-copy">把 Key 放在服务端环境变量中，例如 <code>SST_API_KEY</code>，并由后端统一拼接鉴权头。</p>
              </article>
              <article class="docs-card">
                <div class="docs-card-title">常见失败</div>
                <p class="docs-card-copy"><code>Bearer</code> 前缀缺失、Key 前后多出空格、使用了已停用 Key，都可能返回 401。</p>
              </article>
            </div>
          </section>

          <section v-show="activeSection === 'models'" id="models">
            <div class="docs-kicker">模型查询</div>
            <h2 class="docs-title">接入前先查询模型列表，可以避免模型名错误或权限未开通。</h2>
            <p class="docs-copy">不同账户、分组、渠道或权限范围返回的模型可能不同。返回结果中的 <code>id</code>，通常就是后续请求 <code>chat/completions</code> 时应填写的 <code>model</code> 值。</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>{{ modelsExample }}</code></pre>
          </section>

          <section v-show="activeSection === 'examples'" id="examples">
            <div class="flex flex-wrap items-end justify-between gap-4">
              <div>
                <div class="docs-kicker">请求示例</div>
                <h2 class="docs-title">下面提供最常用的三种调用方式，可直接按需替换模型与消息内容。</h2>
              </div>
              <div class="docs-example-tabs flex flex-wrap gap-2">
                <button
                  v-for="tab in exampleTabs"
                  :key="tab.key"
                  type="button"
                  class="docs-example-tab"
                  :class="activeExample === tab.key ? 'docs-example-tab-active' : ''"
                  @click="activeExample = tab.key"
                >
                  {{ tab.label }}
                </button>
              </div>
            </div>
            <pre class="docs-code mt-5 overflow-x-auto"><code>{{ activeExampleCode }}</code></pre>
          </section>

          <section v-show="activeSection === 'streaming'" id="streaming">
            <div class="docs-kicker">流式输出</div>
            <h2 class="docs-title">需要边生成边展示时，将 <code>stream</code> 设置为 <code>true</code>。</h2>
            <p class="docs-copy">兼容客户端会按 SSE 事件逐段返回内容。适合聊天窗口、逐字展示或需要更快首字响应的场景。</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>{{ streamExample }}</code></pre>
          </section>

          <section v-show="activeSection === 'parameters'" id="parameters">
            <div class="docs-kicker">常用参数</div>
            <h2 class="docs-title">以下参数是最常见的请求字段，具体支持情况以模型与上游能力为准。</h2>
            <div class="docs-table mt-5 overflow-hidden rounded-[1rem] border border-zen-paperLine/70 dark:border-zen-nightLine">
              <div class="docs-table-head grid grid-cols-[11rem_7rem_minmax(0,1fr)]">
                <span class="px-4 py-3">参数</span>
                <span class="px-4 py-3">类型</span>
                <span class="px-4 py-3">说明</span>
              </div>
              <div v-for="item in parameters" :key="item.name" class="docs-table-row grid grid-cols-[11rem_7rem_minmax(0,1fr)] border-t border-zen-paperLine/60 bg-white/26 text-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/45">
                <code class="px-4 py-4 text-zen-ink dark:text-zen-paper">{{ item.name }}</code>
                <span class="docs-table-copy px-4 py-4">{{ item.type }}</span>
                <p class="docs-table-copy px-4 py-4">{{ item.copy }}</p>
              </div>
            </div>
          </section>

          <section v-show="activeSection === 'errors'" id="errors">
            <div class="docs-kicker">响应说明</div>
            <h2 class="docs-title">当接口返回异常状态时，通常可以先从以下几类原因排查。</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in errors" :key="item.title" class="docs-card">
                <div class="docs-card-title">{{ item.title }}</div>
                <p class="docs-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-show="activeSection === 'notes'" id="notes">
            <div class="docs-kicker">注意事项</div>
            <h2 class="docs-title">排查问题时，优先保留模型名、时间、错误码与请求上下文。</h2>
            <ul class="docs-list mt-5">
              <li v-for="item in notes" :key="item.title"><strong>{{ item.title }}：</strong>{{ item.copy }}</li>
            </ul>
          </section>
        </article>
      </section>

      <section class="docs-topic-grid mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        <RouterLink
          v-for="item in topicLinks"
          :key="item.to"
          :to="item.to"
          class="docs-topic-card"
        >
          <div class="docs-topic-card-kicker">{{ item.kicker }}</div>
          <h2>{{ item.title }}</h2>
          <p>{{ item.copy }}</p>
        </RouterLink>
      </section>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'
import { useAppStore } from '@/stores'

const appStore = useAppStore()


const apiBaseUrl = computed(() => (appStore.cachedPublicSettings?.api_base_url || window.location.origin).replace(/\/$/, ''))
const sdkBaseUrl = computed(() => `${apiBaseUrl.value}/v1`)

const sections = [
  { id: 'quickstart', label: '快速开始' },
  { id: 'authentication', label: '认证方式' },
  { id: 'models', label: '获取模型列表' },
  { id: 'examples', label: '请求示例' },
  { id: 'streaming', label: '流式输出' },
  { id: 'parameters', label: '常用参数' },
  { id: 'errors', label: '常见响应说明' },
  { id: 'notes', label: '注意事项' },
] as const

type SectionId = (typeof sections)[number]['id']

const exampleTabs = [
  { key: 'curl', label: 'cURL' },
  { key: 'python', label: 'Python' },
  { key: 'node', label: 'Node.js' },
] as const

const activeExample = ref<(typeof exampleTabs)[number]['key']>('curl')
const activeSection = ref<SectionId>('quickstart')

const modelsExample = computed(() => `curl ${apiBaseUrl.value}/v1/models \\
  -H "Authorization: Bearer YOUR_API_KEY"`)

const curlExample = computed(() => `curl ${apiBaseUrl.value}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -d '{
    "model": "gpt-4o",
    "messages": [
      {"role": "user", "content": "你好"}
    ]
  }'`)

const pythonExample = computed(() => `from openai import OpenAI

client = OpenAI(api_key="YOUR_API_KEY", base_url="${sdkBaseUrl.value}")

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "你好"}],
)

print(response.choices[0].message.content)`)

const nodeExample = computed(() => `import OpenAI from "openai";

const client = new OpenAI({ apiKey: "YOUR_API_KEY", baseURL: "${sdkBaseUrl.value}" });

const completion = await client.chat.completions.create({
  model: "gpt-4o",
  messages: [{ role: "user", content: "你好" }]
});

console.log(completion.choices[0].message.content);`)

const streamExample = computed(() => `curl ${apiBaseUrl.value}/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -d '{
    "model": "gpt-4o",
    "stream": true,
    "messages": [
      {"role": "user", "content": "写一段简短介绍"}
    ]
  }'`)

const activeExampleCode = computed(() => {
  switch (activeExample.value) {
    case 'python':
      return pythonExample.value
    case 'node':
      return nodeExample.value
    default:
      return curlExample.value
  }
})

const parameters = [
  { name: 'model', type: 'string', copy: '要调用的模型名称，建议先通过模型列表接口确认具体值。' },
  { name: 'messages', type: 'array', copy: '对话消息数组，按 OpenAI Chat Completions 标准结构传入。' },
  { name: 'temperature', type: 'number', copy: '可选参数，用于控制生成随机性。' },
  { name: 'stream', type: 'boolean', copy: '可选参数，设为 true 时返回流式输出。' },
  { name: 'max_tokens', type: 'number', copy: '可选参数，用于限制本次生成的最大 token 数。' },
  { name: 'top_p', type: 'number', copy: '可选参数，用于控制采样范围，通常不建议和 temperature 同时大幅调整。' },
] as const

const errors = [
  { title: '401 Unauthorized', copy: 'API Key 缺失、无效、被停用，或请求头没有使用 Bearer 格式。' },
  { title: '403 Forbidden', copy: '当前账号、分组或渠道没有权限调用该模型。' },
  { title: '404 Model Not Found', copy: '模型名不存在，或没有先按模型列表接口返回的 id 填写。' },
  { title: '429 Rate Limited', copy: '请求过于频繁、额度不足，或上游限流。稍后重试或切换可用模型。' },
  { title: '400 Bad Request', copy: '请求体字段格式不正确，常见于 messages 结构、stream 类型或 JSON 格式错误。' },
  { title: '5xx Upstream Error', copy: '通常表示上游服务异常或网络波动，可以稍后重试并保留请求时间用于排查。' },
] as const

const notes = [
  { title: '先查模型', copy: '不同账户、分组或渠道的可用模型可能不同，接入前先调用模型列表接口。' },
  { title: '服务端保存 Key', copy: '不要把 API Key 暴露在浏览器、移动端客户端或公开仓库中。' },
  { title: '地址使用 /v1', copy: 'OpenAI SDK 的 base_url 或 baseURL 通常填写到 /v1 这一层。' },
  { title: '记录请求信息', copy: '排查问题时保留时间、模型名、错误码和 request id，有助于定位上游或权限问题。' },
] as const

const topicLinks = [
  {
    to: '/docs/openai-compatible-api',
    kicker: 'OpenAI API',
    title: 'OpenAI 兼容 API 接入',
    copy: '用现有 SDK 替换 base_url 与 API Key 后接入统一入口。',
  },
  {
    to: '/docs/base-url',
    kicker: 'Base URL',
    title: 'base_url 配置',
    copy: '确认 /v1 基础地址、endpoint 拼接和生产域名边界。',
  },
  {
    to: '/docs/api-key',
    kicker: 'API Key',
    title: 'API Key 使用',
    copy: '创建、保存、鉴权和轮换 Key，避免凭据泄露。',
  },
  {
    to: '/docs/streaming',
    kicker: 'Streaming',
    title: '流式输出',
    copy: '用 stream=true 和 SSE 方式处理逐段响应。',
  },
  {
    to: '/docs/codex',
    kicker: 'Codex',
    title: 'Codex 客户端接入',
    copy: '配置 Codex 类客户端的入口、模型和账册核对。',
  },
  {
    to: '/docs/claude-code',
    kicker: 'Claude Code',
    title: 'Claude Code 接入',
    copy: '统一管理 Claude Code 类客户端的 Key、权限和计量。',
  },
] as const

onMounted(() => {
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.docs-page :deep(main > section:first-child) {
  grid-template-columns: minmax(0, 1fr);
  gap: 0.35rem;
  padding-top: 1.2rem;
  padding-bottom: 0.45rem;
}

.docs-page :deep(.public-copy-block) {
  max-width: 52rem;
}

.docs-page :deep(.public-display) {
  font-size: clamp(2.4rem, 5.3vw, 4.25rem);
  font-weight: 600;
  line-height: 1.02;
  letter-spacing: 0.01em;
}

.docs-page :deep(.public-intro) {
  margin-top: 1rem;
  max-width: 44rem;
  font-family: inherit;
  font-size: clamp(0.94rem, 0.98vw, 1.02rem);
  line-height: 1.8;
  color: #5f685c;
}

.docs-page :deep(.public-hero-panel),
.docs-page :deep(.public-cta) {
  display: none;
}

.docs-main-grid {
  max-width: 60rem;
  grid-template-columns: minmax(0, 11.4rem) minmax(0, 1fr);
}

.docs-main-grid > *,
.docs-article {
  min-width: 0;
}

.docs-lead {
  font-size: clamp(0.96rem, 1.02vw, 1.04rem);
  line-height: 1.88;
}

.docs-lead code,
.docs-copy code,
.docs-list code,
.docs-card-copy code,
.docs-notice-copy code {
  font-size: 0.94em;
}

.docs-meta-strip {
  border-top: 1px solid rgba(216, 205, 185, 0.72);
  padding-top: 0.95rem;
}

.docs-meta-item {
  display: grid;
  gap: 0.2rem;
  align-content: start;
  min-height: 5rem;
  padding: 0.9rem 1rem;
  border: 1px solid rgba(216, 205, 185, 0.56);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.24);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.24);
}

.docs-meta-label {
  font-size: 0.72rem;
  line-height: 1.5;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #8f7d63;
}

.docs-meta-value {
  font-size: 0.92rem;
  line-height: 1.54;
  font-weight: 500;
  color: #1f2320;
}

.docs-toc-caption {
  font-size: 0.72rem;
  line-height: 1.5;
  color: #8f7d63;
}

.docs-toc-link {
  border-left: 1px solid transparent;
  padding: 0.56rem 0 0.56rem 0.7rem;
  text-align: left;
  font-size: 0.95rem;
  line-height: 1.45;
  color: #4a524a;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.docs-toc-link:hover {
  color: #b95d1f;
}

.docs-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.8);
  background: rgba(185, 93, 31, 0.05);
  color: #b95d1f;
}

.docs-kicker {
  font-size: 0.7rem;
  letter-spacing: 0.22em;
  text-transform: uppercase;
  color: #7b6a53;
  font-weight: 500;
}

.docs-title {
  margin-top: 0.85rem;
  width: 100%;
  max-width: none;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.16rem, 1.28vw, 1.42rem);
  line-height: 1.58;
  letter-spacing: 0.01em;
  color: #1f2320;
}

.docs-copy {
  margin-top: 1rem;
  max-width: 52rem;
  font-size: 0.97rem;
  line-height: 1.86;
  color: #5f685c;
}

.docs-list {
  display: grid;
  gap: 0.8rem;
  padding-left: 1.25rem;
  font-size: 0.97rem;
  line-height: 1.86;
  color: #5f685c;
}

.docs-list strong {
  color: #1f2320;
}

.docs-card,
.docs-notice {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.34);
}

.docs-card {
  padding: 1.05rem 1.05rem 1.15rem;
}

.docs-card-title,
.docs-notice-title {
  color: #1f2320;
  font-size: 0.98rem;
  font-weight: 600;
  line-height: 1.45;
}

.docs-card-copy,
.docs-notice-copy,
.docs-table-copy {
  margin-top: 0.55rem;
  font-size: 0.94rem;
  line-height: 1.82;
  color: #5f685c;
}

.docs-notice {
  padding: 1rem 1rem 1.05rem;
  background: linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58));
}

.docs-quickstart-list {
  gap: 0.7rem;
}

.docs-code {
  max-width: 100%;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.42);
  padding: 1rem;
  font-size: 0.92rem;
  line-height: 1.72;
  color: #1f2320;
}

.docs-table-head {
  background: rgba(255, 255, 255, 0.32);
  font-size: 0.72rem;
  line-height: 1.5;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #8f7d63;
}

.docs-table-copy {
  margin-top: 0;
}

.docs-topic-grid {
  max-width: 60rem;
}

.docs-topic-card {
  display: grid;
  gap: 0.65rem;
  min-width: 0;
  border: 1px solid rgba(216, 205, 185, 0.7);
  border-radius: 1rem;
  background: rgba(255, 252, 246, 0.58);
  padding: 1.05rem;
  color: inherit;
  box-shadow: 0 12px 28px rgba(84, 57, 31, 0.045);
  transition: border-color 160ms ease, background-color 160ms ease;
}

.docs-topic-card:hover {
  border-color: rgba(185, 93, 31, 0.28);
  background: rgba(255, 252, 246, 0.82);
}

.docs-topic-card-kicker {
  color: #8f6f43;
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.docs-topic-card h2 {
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.1rem;
  line-height: 1.45;
}

.docs-topic-card p {
  color: #5f685c;
  font-size: 0.92rem;
  line-height: 1.75;
}

.docs-example-tab {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.32);
  padding: 0.52rem 0.95rem;
  font-size: 0.88rem;
  line-height: 1.3;
  color: #6a7267;
  transition: color 160ms ease, border-color 160ms ease, background-color 160ms ease;
}

.docs-example-tab:hover {
  color: #b95d1f;
}

.docs-example-tab-active {
  border-color: rgba(185, 93, 31, 0.32);
  background: rgba(185, 93, 31, 0.08);
  color: #b95d1f;
}

@media (max-width: 1023px) {
  .docs-main-grid {
    grid-template-columns: 1fr;
    gap: 1.15rem;
  }

  .docs-main-grid > aside .sticky {
    position: static;
  }

  .docs-main-grid > aside nav {
    display: flex;
    gap: 0.65rem;
    overflow-x: auto;
    padding-bottom: 0.2rem;
    scrollbar-width: none;
  }

  .docs-main-grid > aside nav::-webkit-scrollbar {
    display: none;
  }

  .docs-toc-link {
    flex: 0 0 auto;
    border: 1px solid rgba(216, 205, 185, 0.64);
    border-left-width: 1px;
    border-radius: 999px;
    padding: 0.58rem 0.88rem;
    font-size: 0.9rem;
    line-height: 1.2;
    white-space: nowrap;
    background: rgba(255, 255, 255, 0.2);
  }

  .docs-title {
    font-size: clamp(1.14rem, 3.7vw, 1.38rem);
    line-height: 1.48;
  }

  .docs-kicker {
    letter-spacing: 0.18em;
  }

  .docs-meta-item {
    min-height: 4.6rem;
    padding: 0.82rem 0.9rem;
  }

  .docs-meta-value,
  .docs-copy,
  .docs-list,
  .docs-card-copy,
  .docs-notice-copy,
  .docs-lead,
  .docs-code,
  .docs-table-copy {
    font-size: 0.95rem;
  }

  .docs-article {
    padding: 0.95rem;
  }

  .docs-table-head,
  .docs-table-row {
    grid-template-columns: 1fr;
  }

  .docs-table-head span,
  .docs-table-row > * {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .docs-table-head span:last-child {
    padding-top: 0;
  }

  .docs-table-row > :first-child {
    padding-bottom: 0.35rem;
  }

  .docs-table-row > :last-child {
    padding-top: 0;
  }

  .docs-toc-link-active {
    border-left-color: rgba(185, 93, 31, 0.8);
  }
}

@media (max-width: 639px) {
  .docs-page :deep(main > section:first-child) {
    gap: 0.2rem;
    padding-top: 0.95rem;
    padding-bottom: 0.3rem;
  }

  .docs-page :deep(.public-display) {
    font-size: clamp(2.22rem, 12vw, 3.3rem);
    line-height: 0.98;
  }

  .docs-page :deep(.public-intro) {
    margin-top: 0.8rem;
    font-size: 0.98rem;
    line-height: 1.74;
  }

  .docs-page :deep(main > section:first-child) {
    gap: 0.12rem;
  }

  .docs-meta-strip {
    gap: 0.75rem;
    padding-top: 0.85rem;
  }

  .docs-meta-label,
  .docs-toc-caption {
    font-size: 0.68rem;
    letter-spacing: 0.14em;
  }

  .docs-meta-value {
    font-size: 0.88rem;
    line-height: 1.48;
  }

  .docs-main-grid > aside nav {
    gap: 0.55rem;
  }

  .docs-copy,
  .docs-list,
  .docs-card-copy,
  .docs-notice-copy,
  .docs-code,
  .docs-table-copy {
    line-height: 1.74;
  }

  .docs-quickstart-list {
    gap: 0.58rem;
  }

  .docs-notice {
    padding: 0.86rem 0.9rem 0.9rem;
  }
}
</style>

<style>
html.dark .docs-page .public-intro,
.docs-page.is-dark .public-intro {
  color: #ead9bd;
}

html.dark .docs-page .public-display,
.docs-page.is-dark .public-display {
  color: #fff4dd;
  text-shadow: 0 1px 0 rgba(255, 240, 218, 0.08), 0 18px 48px rgba(0, 0, 0, 0.28);
}

html.dark .docs-page .public-copy-block > div:first-child span:last-child,
html.dark .docs-page .docs-meta-label,
html.dark .docs-page .docs-kicker,
html.dark .docs-page .docs-table-head,
html.dark .docs-page .docs-main-grid > aside > div > div:first-child,
.docs-page.is-dark .public-copy-block > div:first-child span:last-child,
.docs-page.is-dark .docs-meta-label,
.docs-page.is-dark .docs-kicker,
.docs-page.is-dark .docs-table-head,
.docs-page.is-dark .docs-main-grid > aside > div > div:first-child {
  color: #cdb387;
}

html.dark .docs-page .docs-meta-value,
html.dark .docs-page .docs-title,
html.dark .docs-page .docs-card-title,
html.dark .docs-page .docs-notice-title,
html.dark .docs-page .docs-table-row > :first-child,
html.dark .docs-page .docs-list strong,
.docs-page.is-dark .docs-meta-value,
.docs-page.is-dark .docs-title,
.docs-page.is-dark .docs-card-title,
.docs-page.is-dark .docs-notice-title,
.docs-page.is-dark .docs-table-row > :first-child,
.docs-page.is-dark .docs-list strong {
  color: #fff0d5;
}

html.dark .docs-page .docs-toc-link:hover,
html.dark .docs-page .docs-toc-link-active,
.docs-page.is-dark .docs-toc-link:hover,
.docs-page.is-dark .docs-toc-link-active {
  color: #f3c786;
}

html.dark .docs-page .docs-toc-link-active,
.docs-page.is-dark .docs-toc-link-active {
  border-left-color: rgba(212, 153, 80, 0.8);
  background: linear-gradient(90deg, rgba(176, 120, 57, 0.16), rgba(176, 120, 57, 0.04));
}

html.dark .docs-page .docs-card,
html.dark .docs-page .docs-notice,
html.dark .docs-page .docs-code,
html.dark .docs-page .docs-table-row,
.docs-page.is-dark .docs-meta-item,
.docs-page.is-dark .docs-card,
.docs-page.is-dark .docs-notice,
.docs-page.is-dark .docs-code,
.docs-page.is-dark .docs-table-row {
  border-color: rgba(141, 109, 72, 0.56) !important;
}

html.dark .docs-page .docs-meta-item,
.docs-page.is-dark .docs-meta-item {
  background:
    linear-gradient(180deg, rgba(29, 31, 26, 0.9), rgba(22, 24, 20, 0.94)),
    radial-gradient(circle at 84% 18%, rgba(173, 102, 46, 0.08), transparent 24%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 242, 219, 0.04),
    0 10px 24px rgba(0, 0, 0, 0.12) !important;
}

html.dark .docs-article {
  border-color: rgba(120, 109, 90, 0.54) !important;
  background:
    linear-gradient(180deg, rgba(35, 37, 31, 0.95), rgba(24, 26, 21, 0.96)),
    radial-gradient(circle at top right, rgba(166, 97, 45, 0.1), transparent 28%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 244, 224, 0.07),
    0 18px 40px rgba(0, 0, 0, 0.2) !important;
}

html.dark .docs-lead,
html.dark .docs-copy,
html.dark .docs-list,
html.dark .docs-card-copy,
html.dark .docs-notice-copy,
html.dark .docs-table-row p,
html.dark .docs-table-copy {
  color: #e2d6c4;
}

html.dark .docs-meta-value,
html.dark .docs-title,
html.dark .docs-card-title,
html.dark .docs-notice-title,
html.dark .docs-copy strong,
html.dark .docs-list strong,
html.dark .docs-table-row > :first-child {
  color: #fff0da;
}

html.dark .docs-meta-label,
html.dark .docs-kicker,
html.dark .docs-table-head,
html.dark .docs-main-grid > aside > div > div:first-child {
  color: #caba9f;
}

html.dark .docs-meta-strip {
  border-top-color: rgba(102, 95, 79, 0.68);
}

html.dark .docs-toc-link {
  color: #f0e3cc;
}

html.dark .docs-toc-link:hover,
html.dark .docs-toc-link-active {
  color: #ffe0bf;
}

html.dark .docs-card,
html.dark .docs-page .docs-card,
html.dark .docs-main-grid .docs-card,
html.dark .docs-page .docs-main-grid .docs-card {
  border-color: rgba(118, 106, 87, 0.58) !important;
  background:
    linear-gradient(180deg, rgba(36, 38, 31, 0.94), rgba(27, 29, 24, 0.95)),
    radial-gradient(circle at 88% 18%, rgba(163, 97, 45, 0.09), transparent 24%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 242, 219, 0.06),
    0 14px 28px rgba(0, 0, 0, 0.18) !important;
}

html.dark .docs-card-title,
html.dark .docs-page .docs-card-title,
html.dark .docs-main-grid .docs-card-title,
html.dark .docs-page .docs-main-grid .docs-card-title {
  color: #f4ead7 !important;
}

html.dark .docs-card-copy,
html.dark .docs-page .docs-card-copy,
html.dark .docs-main-grid .docs-card-copy,
html.dark .docs-page .docs-main-grid .docs-card-copy {
  color: #d9cdbb !important;
}

html.dark .docs-notice,
html.dark .docs-page .docs-notice,
html.dark .docs-main-grid .docs-notice,
html.dark .docs-page .docs-main-grid .docs-notice {
  border-color: rgba(136, 115, 88, 0.52) !important;
  background:
    linear-gradient(145deg, rgba(42, 38, 31, 0.95), rgba(30, 27, 23, 0.96)),
    radial-gradient(circle at 84% 34%, rgba(173, 89, 36, 0.12), transparent 30%) !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 241, 220, 0.07),
    0 16px 32px rgba(0, 0, 0, 0.24) !important;
}

html.dark .docs-notice-title,
html.dark .docs-page .docs-notice-title,
html.dark .docs-main-grid .docs-notice-title,
html.dark .docs-page .docs-main-grid .docs-notice-title {
  color: #f7ebd8 !important;
}

html.dark .docs-notice-copy,
html.dark .docs-page .docs-notice-copy,
html.dark .docs-main-grid .docs-notice-copy,
html.dark .docs-page .docs-main-grid .docs-notice-copy {
  color: #dfd2be !important;
}

html.dark .docs-code,
html.dark .docs-page .docs-code,
html.dark .docs-main-grid .docs-code,
html.dark .docs-page .docs-main-grid .docs-code {
  border-color: rgba(130, 117, 94, 0.52) !important;
  background:
    linear-gradient(180deg, rgba(24, 26, 22, 0.98), rgba(17, 19, 16, 0.99)),
    radial-gradient(circle at 88% 18%, rgba(163, 97, 45, 0.09), transparent 22%) !important;
  color: #f2e6d4 !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 242, 219, 0.05),
    0 14px 28px rgba(0, 0, 0, 0.12) !important;
}

html.dark .docs-table-head {
  background: rgba(35, 32, 27, 0.92);
  color: #d1c2a8;
}

html.dark .docs-table-row {
  background: rgba(29, 31, 26, 0.72) !important;
  border-top-color: rgba(96, 89, 75, 0.58) !important;
}

@media (max-width: 1023px) {
  html.dark .docs-page .docs-toc-link {
    border-color: rgba(118, 106, 87, 0.58);
    background: rgba(24, 27, 22, 0.76);
  }

  html.dark .docs-page .docs-toc-link-active {
    border-color: rgba(212, 153, 80, 0.72);
    background: linear-gradient(90deg, rgba(176, 120, 57, 0.16), rgba(176, 120, 57, 0.05));
  }
}

html.dark .docs-article {
  border-color: rgba(155, 126, 86, 0.26) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .docs-page .docs-card,
html.dark .docs-page .docs-notice,
html.dark .docs-page .docs-code,
html.dark .docs-page .docs-table,
html.dark .docs-page .docs-table-row {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.88), rgba(14, 16, 13, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.08), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 238, 210, 0.05) !important;
}

html.dark .docs-page .docs-notice {
  background:
    linear-gradient(135deg, rgba(35, 29, 23, 0.88), rgba(24, 27, 22, 0.84)),
    radial-gradient(circle at 84% 18%, rgba(194, 126, 74, 0.13), transparent 28%) !important;
}

html.dark .docs-page .docs-table-head {
  background:
    linear-gradient(180deg, rgba(39, 32, 26, 0.9), rgba(24, 27, 22, 0.88)) !important;
  color: #d0baa0 !important;
}

html.dark .docs-page .docs-title,
html.dark .docs-page .docs-card-title,
html.dark .docs-page .docs-notice-title,
html.dark .docs-page .docs-table-row > :first-child,
html.dark .docs-page .docs-list strong {
  color: #f6e8d2 !important;
}

html.dark .docs-page .docs-copy,
html.dark .docs-page .docs-list,
html.dark .docs-page .docs-card-copy,
html.dark .docs-page .docs-notice-copy,
html.dark .docs-page .docs-table-copy,
html.dark .docs-page .docs-table-row p {
  color: #d0baa0 !important;
}

html.dark .docs-page .docs-kicker,
html.dark .docs-page .docs-main-grid > aside > div > div:first-child,
html.dark .docs-page .docs-toc-caption,
html.dark .docs-page .public-copy-block > div:first-child span:last-child {
  color: #d8b171 !important;
}

html.dark .docs-page .docs-toc-link {
  color: #d4c4ad !important;
}

html.dark .docs-page .docs-toc-link:hover,
html.dark .docs-page .docs-toc-link-active {
  color: #efc183 !important;
}

html.dark .docs-page .docs-toc-link-active {
  border-left-color: rgba(194, 126, 74, 0.72) !important;
  background: linear-gradient(90deg, rgba(194, 126, 74, 0.16), rgba(194, 126, 74, 0.04)) !important;
}

html:not(.dark) .docs-article {
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

html:not(.dark) .docs-page .docs-card,
html:not(.dark) .docs-page .docs-notice,
html:not(.dark) .docs-page .docs-code,
html:not(.dark) .docs-page .docs-table,
html:not(.dark) .docs-page .docs-table-row {
  border-color: rgba(190, 168, 134, 0.42) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.58), rgba(244, 235, 220, 0.34)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.06), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.42) !important;
}

html:not(.dark) .docs-page .docs-notice {
  background:
    linear-gradient(135deg, rgba(255, 252, 247, 0.78), rgba(249, 240, 230, 0.58)),
    radial-gradient(circle at 84% 18%, rgba(196, 136, 68, 0.08), transparent 28%) !important;
}

html:not(.dark) .docs-page .docs-table-head {
  background: rgba(255, 252, 246, 0.48) !important;
  color: #8f7d63 !important;
}

html:not(.dark) .docs-page .docs-toc-link {
  color: #4a524a !important;
}

html:not(.dark) .docs-page .docs-toc-link:hover,
html:not(.dark) .docs-page .docs-toc-link-active {
  color: #b95d1f !important;
}

html:not(.dark) .docs-page .docs-toc-link-active {
  border-left-color: rgba(185, 93, 31, 0.74) !important;
  background: linear-gradient(90deg, rgba(185, 93, 31, 0.07), rgba(185, 93, 31, 0.02)) !important;
}
</style>

