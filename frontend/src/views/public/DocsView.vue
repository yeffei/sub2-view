<template>
  <PublicPageLayout
    class="docs-page"
    tone="docs"
    :eyebrow="t('publicDocs.eyebrow')"
    :title="t('publicDocs.title')"
    intro=""
    description=""
    :show-cta="false"
  >
      <section class="docs-main-grid mt-10 grid gap-8 lg:gap-8">
        <aside class="self-start">
          <div class="sticky top-5 h-fit">
            <div class="mb-3 flex items-end justify-between gap-3">
              <div class="text-xs uppercase tracking-[0.24em] text-zen-mist dark:text-zen-stone">{{ t('publicDocs.toc') }}</div>
              <div class="docs-toc-caption">{{ t('publicDocs.tocCaption') }}</div>
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
            <div class="docs-kicker">{{ t('publicDocs.quickstart.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.quickstart.title') }}</h2>
            <ol class="docs-list docs-quickstart-list mt-5">
              <li>{{ t('publicDocs.quickstart.steps.login') }}</li>
              <li>{{ t('publicDocs.quickstart.steps.configureBefore') }} <code>base_url</code> {{ t('publicDocs.quickstart.steps.or') }} <code>baseURL</code> {{ t('publicDocs.quickstart.steps.configureAfter') }} <code>{{ sdkBaseUrl }}</code>{{ t('publicDocs.quickstart.steps.authHeader') }} <code>Authorization: Bearer YOUR_API_KEY</code>.</li>
              <li>{{ t('publicDocs.quickstart.steps.models') }}</li>
            </ol>
            <div class="docs-notice mt-6">
              <div>
                <div class="docs-notice-title">{{ t('publicDocs.quickstart.noticeTitle') }}</div>
                <p class="docs-notice-copy">{{ t('publicDocs.quickstart.noticeBefore') }} <code>api_key</code> {{ t('publicDocs.quickstart.steps.or') }} <code>base_url</code>{{ t('publicDocs.quickstart.noticeAfter') }}</p>
              </div>
            </div>
          </section>

          <section v-show="activeSection === 'authentication'" id="authentication">
            <div class="docs-kicker">{{ t('publicDocs.authentication.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.authentication.title') }}</h2>
            <p class="docs-copy">{{ t('publicDocs.authentication.copy') }}</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>Authorization: Bearer YOUR_API_KEY</code></pre>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article class="docs-card">
                <div class="docs-card-title">{{ t('publicDocs.authentication.recommendedTitle') }}</div>
                <p class="docs-card-copy">{{ t('publicDocs.authentication.recommendedBefore') }} <code>SST_API_KEY</code>{{ t('publicDocs.authentication.recommendedAfter') }}</p>
              </article>
              <article class="docs-card">
                <div class="docs-card-title">{{ t('publicDocs.authentication.failureTitle') }}</div>
                <p class="docs-card-copy"><code>Bearer</code> {{ t('publicDocs.authentication.failureCopy') }}</p>
              </article>
            </div>
          </section>

          <section v-show="activeSection === 'models'" id="models">
            <div class="docs-kicker">{{ t('publicDocs.models.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.models.title') }}</h2>
            <p class="docs-copy">{{ t('publicDocs.models.copyBefore') }} <code>id</code>{{ t('publicDocs.models.copyMiddle') }} <code>chat/completions</code>{{ t('publicDocs.models.copyAfter') }} <code>model</code> {{ t('publicDocs.models.copyEnd') }}</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>{{ modelsExample }}</code></pre>
          </section>

          <section v-show="activeSection === 'examples'" id="examples">
            <div class="flex flex-wrap items-end justify-between gap-4">
              <div>
                <div class="docs-kicker">{{ t('publicDocs.examples.kicker') }}</div>
                <h2 class="docs-title">{{ t('publicDocs.examples.title') }}</h2>
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
            <div class="docs-kicker">{{ t('publicDocs.streaming.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.streaming.titleBefore') }} <code>stream</code> {{ t('publicDocs.streaming.titleAfter') }} <code>true</code>.</h2>
            <p class="docs-copy">{{ t('publicDocs.streaming.copy') }}</p>
            <pre class="docs-code mt-5 overflow-x-auto"><code>{{ streamExample }}</code></pre>
          </section>

          <section v-show="activeSection === 'parameters'" id="parameters">
            <div class="docs-kicker">{{ t('publicDocs.parameters.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.parameters.title') }}</h2>
            <div class="docs-table mt-5 overflow-hidden rounded-[1rem] border border-zen-paperLine/70 dark:border-zen-nightLine">
              <div class="docs-table-head grid grid-cols-[11rem_7rem_minmax(0,1fr)]">
                <span class="px-4 py-3">{{ t('publicDocs.parameters.table.name') }}</span>
                <span class="px-4 py-3">{{ t('publicDocs.parameters.table.type') }}</span>
                <span class="px-4 py-3">{{ t('publicDocs.parameters.table.description') }}</span>
              </div>
              <div v-for="item in parameters" :key="item.name" class="docs-table-row grid grid-cols-[11rem_7rem_minmax(0,1fr)] border-t border-zen-paperLine/60 bg-white/26 text-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/45">
                <code class="px-4 py-4 text-zen-ink dark:text-zen-paper">{{ item.name }}</code>
                <span class="docs-table-copy px-4 py-4">{{ item.type }}</span>
                <p class="docs-table-copy px-4 py-4">{{ item.copy }}</p>
              </div>
            </div>
          </section>

          <section v-show="activeSection === 'errors'" id="errors">
            <div class="docs-kicker">{{ t('publicDocs.errors.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.errors.title') }}</h2>
            <div class="mt-6 grid gap-4 md:grid-cols-2">
              <article v-for="item in errors" :key="item.title" class="docs-card">
                <div class="docs-card-title">{{ item.title }}</div>
                <p class="docs-card-copy">{{ item.copy }}</p>
              </article>
            </div>
          </section>

          <section v-show="activeSection === 'notes'" id="notes">
            <div class="docs-kicker">{{ t('publicDocs.notes.kicker') }}</div>
            <h2 class="docs-title">{{ t('publicDocs.notes.title') }}</h2>
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
import { useI18n } from 'vue-i18n'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'
import { useAppStore } from '@/stores'

const appStore = useAppStore()
const { t, tm } = useI18n()


const apiBaseUrl = computed(() => (appStore.cachedPublicSettings?.api_base_url || window.location.origin).replace(/\/$/, ''))
const sdkBaseUrl = computed(() => `${apiBaseUrl.value}/v1`)

const sections = computed(() => tm('publicDocs.sections') as unknown as Array<{
  id: 'quickstart' | 'authentication' | 'models' | 'examples' | 'streaming' | 'parameters' | 'errors' | 'notes'
  label: string
}>)

type SectionId = 'quickstart' | 'authentication' | 'models' | 'examples' | 'streaming' | 'parameters' | 'errors' | 'notes'

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

const parameters = computed(() => tm('publicDocs.parameters.items') as unknown as Array<{ name: string, type: string, copy: string }>)
const errors = computed(() => tm('publicDocs.errors.items') as unknown as Array<{ title: string, copy: string }>)
const notes = computed(() => tm('publicDocs.notes.items') as unknown as Array<{ title: string, copy: string }>)
const topicLinks = computed(() => tm('publicDocs.topicLinks') as unknown as Array<{ to: string, kicker: string, title: string, copy: string }>)

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

