<template>
  <PublicPageLayout
    class="docs-topic-page"
    tone="docs"
    :eyebrow="topic.eyebrow"
    :title="topic.title"
    :intro="topic.intro"
    :description="topic.description"
    :highlights="topic.highlights"
    authenticated-action-label="入庭"
  >
    <template #aside>
      <div class="space-y-5">
        <div>
          <div class="text-xs uppercase tracking-[0.36em] text-zen-mist dark:text-zen-stone">Docs</div>
          <h2 class="mt-3 font-serif text-3xl leading-tight text-zen-ink dark:text-zen-paper">{{ topic.asideTitle }}</h2>
          <p class="mt-3 text-sm leading-7 text-zen-mist dark:text-zen-stone">{{ topic.asideCopy }}</p>
        </div>

        <nav class="grid gap-2" aria-label="专题文档">
          <RouterLink
            v-for="item in topics"
            :key="item.path"
            :to="item.path"
            class="docs-topic-link"
            :class="route.path === item.path ? 'is-active' : ''"
          >
            <span>{{ item.navLabel }}</span>
          </RouterLink>
        </nav>
      </div>
    </template>

    <article class="docs-topic-article">
      <section v-for="section in topic.sections" :key="section.heading" class="docs-topic-section">
        <div class="docs-topic-kicker">{{ section.kicker }}</div>
        <h2>{{ section.heading }}</h2>
        <p v-for="paragraph in section.paragraphs" :key="paragraph">{{ paragraph }}</p>
        <pre v-if="section.code" class="docs-topic-code"><code>{{ section.code }}</code></pre>
        <ul v-if="section.points?.length" class="docs-topic-list">
          <li v-for="point in section.points" :key="point">{{ point }}</li>
        </ul>
      </section>
    </article>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'

interface TopicSection {
  kicker: string
  heading: string
  paragraphs: string[]
  code?: string
  points?: string[]
}

interface TopicContent {
  path: string
  navLabel: string
  eyebrow: string
  title: string
  intro: string
  description: string
  highlights: string[]
  asideTitle: string
  asideCopy: string
  sections: TopicSection[]
}

const topics: TopicContent[] = [
  {
    path: '/docs/openai-compatible-api',
    navLabel: 'OpenAI 兼容接口',
    eyebrow: 'OpenAI API',
    title: 'OpenAI 兼容 API 接入',
    intro: '客户端通常只需要替换 base_url 与 API Key。',
    description: '山枢庭 SST 保持 OpenAI 兼容调用方式，便于现有 SDK、CLI 与服务端应用接入统一入口。',
    highlights: ['OpenAI SDK', 'chat/completions', 'models'],
    asideTitle: '迁移重点',
    asideCopy: '保留客户端结构，先确认 base_url、Authorization 与模型 ID，再进入业务调用。',
    sections: [
      {
        kicker: 'Base URL',
        heading: '把 SDK 的基础地址指向 SST 的 /v1 入口。',
        paragraphs: ['大多数 OpenAI 兼容 SDK 都支持自定义 base_url 或 baseURL。接入时保持原有消息结构，先替换基础地址和 Key。'],
        code: `base_url = "https://your-domain.example/v1"`,
      },
      {
        kicker: 'Auth',
        heading: '请求使用 Bearer API Key 鉴权。',
        paragraphs: ['API Key 属于账户凭据，不应放入公开前端代码。推荐由服务端读取环境变量，再向 SST 发起调用。'],
        code: `Authorization: Bearer YOUR_API_KEY`,
      },
      {
        kicker: 'Check',
        heading: '先查询模型列表，再发起正式请求。',
        paragraphs: ['不同分组、权限或上游池可见的模型可能不同。先调用模型列表可以减少模型名错误和权限未开通造成的失败。'],
        code: `curl https://your-domain.example/v1/models \\
  -H "Authorization: Bearer YOUR_API_KEY"`,
      },
    ],
  },
  {
    path: '/docs/base-url',
    navLabel: 'base_url 配置',
    eyebrow: 'Base URL',
    title: 'base_url 与 endpoint 配置',
    intro: '基础地址决定请求进入哪一个统一入口。',
    description: '配置 SST base_url 时应保留 /v1 前缀，避免把业务路径、网关路径和模型路径混在一起。',
    highlights: ['base_url', '/v1', 'endpoint'],
    asideTitle: '地址原则',
    asideCopy: '基础地址只到 /v1，具体资源路径由 SDK 或客户端自动拼接。',
    sections: [
      {
        kicker: 'Format',
        heading: '推荐把 base_url 写到 /v1。',
        paragraphs: ['如果 SDK 会自动拼接 chat/completions、models 等路径，base_url 应保持在 /v1 这一层。'],
        code: `https://your-domain.example/v1`,
      },
      {
        kicker: 'Avoid',
        heading: '不要把完整业务路径写进 base_url。',
        paragraphs: ['把 /v1/chat/completions 写入 base_url 可能导致 SDK 再次拼接路径，形成重复或错误 URL。'],
        points: ['正确：/v1', '错误：/v1/chat/completions', '错误：缺少 /v1 但客户端又不自动补齐'],
      },
      {
        kicker: 'Runtime',
        heading: '生产环境优先使用自有域名。',
        paragraphs: ['上线后 sitemap、canonical、回调和客户端配置都应指向同一个公开域名，避免 Google 与用户看到不同入口。'],
      },
    ],
  },
  {
    path: '/docs/api-key',
    navLabel: 'API Key 使用',
    eyebrow: 'API Key',
    title: 'API Key 创建与使用',
    intro: 'Key 是调用统一入口的账户凭据。',
    description: '了解 SST API Key 的创建、保存、调用、轮换和权限边界，避免在公开代码或客户端泄露凭据。',
    highlights: ['Authorization', 'Bearer', 'Key rotation'],
    asideTitle: '凭据边界',
    asideCopy: 'Key 用于服务端调用，不应暴露在公开仓库、浏览器代码或客户端安装包中。',
    sections: [
      {
        kicker: 'Create',
        heading: '登录控制台后创建 API Key。',
        paragraphs: ['创建 Key 后应立即保存到服务端环境变量或密钥管理系统中，避免通过聊天记录、截图或公开文档传播。'],
      },
      {
        kicker: 'Use',
        heading: '请求头使用 Authorization Bearer。',
        paragraphs: ['所有 API 请求都应携带有效 Key。前后空格、缺少 Bearer 前缀或使用停用 Key 都会导致鉴权失败。'],
        code: `Authorization: Bearer YOUR_API_KEY`,
      },
      {
        kicker: 'Operate',
        heading: '定期轮换和停用不再使用的 Key。',
        paragraphs: ['当成员离职、服务迁移或怀疑泄露时，应创建新 Key、切换服务端配置，再停用旧 Key。'],
      },
    ],
  },
  {
    path: '/docs/streaming',
    navLabel: '流式输出',
    eyebrow: 'Streaming',
    title: '流式输出与 SSE 调用',
    intro: '需要更快首字响应时，将 stream 设置为 true。',
    description: 'SST 支持 OpenAI 兼容的流式输出方式，适用于聊天窗口、终端客户端和需要逐段展示的应用。',
    highlights: ['stream=true', 'SSE', '逐段返回'],
    asideTitle: '流式场景',
    asideCopy: '聊天、终端和长文本生成通常更适合流式输出，服务端需要按事件逐段读取。',
    sections: [
      {
        kicker: 'Request',
        heading: '请求体中启用 stream。',
        paragraphs: ['流式输出会把响应拆成连续事件，客户端需要持续读取连接直到结束。'],
        code: `"stream": true`,
      },
      {
        kicker: 'Client',
        heading: '客户端应按 SSE 或 SDK 的流接口处理。',
        paragraphs: ['不要把流式响应当成一次性 JSON 解析。不同 SDK 的流接口名称不同，但核心都是逐段读取增量内容。'],
      },
      {
        kicker: 'Fallback',
        heading: '排障时先切回非流式请求。',
        paragraphs: ['如果客户端无法正确读取流，可以先关闭 stream，确认模型、Key、base_url 和基础请求都正常，再恢复流式处理。'],
      },
    ],
  },
  {
    path: '/docs/codex',
    navLabel: 'Codex 接入',
    eyebrow: 'Codex',
    title: 'Codex 客户端接入',
    intro: '将 Codex 客户端请求接入 SST 统一入口。',
    description: '使用 SST 为 Codex 类客户端配置统一 API 入口、Key、模型与账册，便于集中计量和维护。',
    highlights: ['Codex', '统一入口', '账册'],
    asideTitle: '接入顺序',
    asideCopy: '先确认客户端支持自定义 base_url，再配置 Key 和模型名，最后查看用量记录。',
    sections: [
      {
        kicker: 'Config',
        heading: '客户端配置应指向 SST 的 OpenAI 兼容地址。',
        paragraphs: ['支持 OpenAI 兼容接口的 Codex 类客户端，通常可以通过环境变量或配置文件替换 base_url 与 API Key。'],
      },
      {
        kicker: 'Model',
        heading: '模型名以当前账号可见列表为准。',
        paragraphs: ['不要只复制公开示例里的模型名。实际可用模型取决于管理员配置、分组权限和上游可用状态。'],
      },
      {
        kicker: 'Usage',
        heading: '调用后在控制台核对用量和账册。',
        paragraphs: ['接入完成后，应发起一次小请求，并在控制台检查请求、消耗、模型和分组是否符合预期。'],
      },
    ],
  },
  {
    path: '/docs/claude-code',
    navLabel: 'Claude Code 接入',
    eyebrow: 'Claude Code',
    title: 'Claude Code 接入说明',
    intro: '将 Claude Code 类调用纳入统一入口和计量。',
    description: '通过 SST 统一管理 Claude Code 类客户端的入口、Key、模型权限、倍率和调用账册。',
    highlights: ['Claude Code', '模型权限', '计量'],
    asideTitle: '配置重点',
    asideCopy: '先确认客户端的接口兼容模式，再配置入口、Key、模型与权限边界。',
    sections: [
      {
        kicker: 'Endpoint',
        heading: '使用管理员提供的 API 入口。',
        paragraphs: ['不同部署可能启用不同上游池和模型映射。客户端应使用当前站点公开或管理员提供的 endpoint。'],
      },
      {
        kicker: 'Auth',
        heading: '用 SST API Key 统一鉴权。',
        paragraphs: ['Key 与账户、分组、并发和账册关联。不要把上游账号凭据直接分发给终端用户。'],
      },
      {
        kicker: 'Billing',
        heading: '价格和倍率以价格页与控制台为准。',
        paragraphs: ['公开价格页用于说明主要计量口径，实际订单、余额和用量明细以登录后的控制台记录为准。'],
      },
    ],
  },
]

const route = useRoute()
const topic = computed(() => topics.find((item) => item.path === route.path) ?? topics[0])
</script>

<style scoped>
.docs-topic-article {
  display: grid;
  gap: 1rem;
}

.docs-topic-section {
  border: 1px solid rgba(216, 205, 185, 0.7);
  border-radius: 1.2rem;
  background: rgba(255, 252, 246, 0.62);
  padding: 1.35rem;
  box-shadow: 0 14px 34px rgba(84, 57, 31, 0.05);
}

.docs-topic-section h2 {
  margin-top: 0.5rem;
  font-family: var(--font-serif);
  font-size: 1.55rem;
  line-height: 1.25;
  color: #252920;
}

.docs-topic-section p {
  margin-top: 0.8rem;
  color: #64705f;
  line-height: 1.8;
}

.docs-topic-kicker {
  font-size: 0.72rem;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: #9b7b54;
}

.docs-topic-code {
  margin-top: 1rem;
  overflow-x: auto;
  border-radius: 0.95rem;
  border: 1px solid rgba(126, 112, 87, 0.18);
  background: rgba(31, 35, 32, 0.94);
  padding: 1rem;
  color: #f4efe4;
  font-size: 0.86rem;
  line-height: 1.65;
}

.docs-topic-list {
  margin-top: 1rem;
  display: grid;
  gap: 0.55rem;
  color: #64705f;
  line-height: 1.7;
}

.docs-topic-list li {
  position: relative;
  padding-left: 1rem;
}

.docs-topic-list li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0.72em;
  width: 0.34rem;
  height: 0.34rem;
  border-radius: 999px;
  background: #a73a2a;
}

.docs-topic-link {
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 0.85rem;
  background: rgba(255, 252, 246, 0.42);
  padding: 0.72rem 0.85rem;
  color: #485247;
  font-size: 0.92rem;
  transition: border-color 160ms ease, background-color 160ms ease, color 160ms ease;
}

.docs-topic-link:hover,
.docs-topic-link.is-active {
  border-color: rgba(167, 58, 42, 0.28);
  background: rgba(167, 58, 42, 0.06);
  color: #1f2320;
}

:global(html.dark) .docs-topic-section {
  border-color: rgba(173, 145, 104, 0.2);
  background: rgba(21, 24, 19, 0.78);
  box-shadow: 0 22px 48px rgba(0, 0, 0, 0.22);
}

:global(html.dark) .docs-topic-section h2 {
  color: #f4efe4;
}

:global(html.dark) .docs-topic-section p,
:global(html.dark) .docs-topic-list {
  color: #bda98f;
}

:global(html.dark) .docs-topic-link {
  border-color: rgba(173, 145, 104, 0.2);
  background: rgba(21, 24, 19, 0.72);
  color: #c9b28d;
}

:global(html.dark) .docs-topic-link:hover,
:global(html.dark) .docs-topic-link.is-active {
  border-color: rgba(188, 93, 31, 0.38);
  background: rgba(53, 39, 28, 0.82);
  color: #f4efe4;
}
</style>
