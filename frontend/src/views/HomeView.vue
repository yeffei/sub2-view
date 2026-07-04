<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="trimmedHomeContent"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- SST Default Home Page -->
  <div v-else class="sst-home min-h-screen overflow-hidden" :class="{ 'is-dark': isDark }" :style="homeBackgroundStyle">
    <div class="sst-bg" aria-hidden="true"></div>
    <div class="sst-paper" aria-hidden="true"></div>

    <PublicSiteHeader />

    <main class="relative z-10">
      <section class="home-hero-shell mx-auto grid max-w-7xl gap-10 px-5 pb-10 pt-10 sm:px-8 lg:min-h-[calc(100vh-84px)] lg:grid-cols-[0.94fr_1.06fr] lg:items-center lg:gap-6 lg:pt-0">
        <div class="max-w-3xl">
          <div class="mb-7 flex items-center gap-4">
            <span class="h-px w-14 bg-zen-paperLine/80 dark:bg-zen-nightLine"></span>
            <span class="home-kicker text-xs uppercase tracking-[0.42em] text-zen-mist dark:text-zen-stone">
              SST
            </span>
          </div>

          <h1 class="home-title font-serif text-[clamp(3.2rem,8.35vw,7.08rem)] font-semibold leading-[0.94] tracking-normal text-zen-ink dark:text-zen-paper">
            {{ brandName }}
          </h1>
          <p class="home-lead mt-7 max-w-2xl font-serif text-[1.9rem] leading-tight text-zen-inkSoft dark:text-zen-paper sm:text-4xl">
            统一入口，安静流转。
          </p>
          <p class="home-copy mt-6 max-w-[40rem] text-base leading-8 text-zen-mist dark:text-zen-stone sm:text-lg">
            {{ heroSubtitle }}
          </p>

          <div class="home-cta-row mt-10 flex flex-col gap-3 sm:flex-row">
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="home-primary-cta inline-flex items-center justify-center rounded-zen px-7 py-3 text-sm font-medium"
            >
              入庭
              <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
            </router-link>
            <component
              :is="docsLinkTag"
              v-bind="docsLinkProps"
              class="home-secondary-cta inline-flex items-center justify-center rounded-zen px-7 py-3 text-sm font-medium"
            >
              接入文档
            </component>
          </div>
        </div>

        <div class="hero-visual relative lg:min-h-[34rem]" aria-hidden="true">
          <div class="hero-courtyard-flow"></div>
          <div class="hero-logo-stage">
            <img src="/logo.png" alt="" class="hero-logo-image" />
          </div>
          <div class="hero-seal-imprint">
            <img src="/sst-seal.svg" alt="" />
          </div>

          <div class="hero-glyph-notes">
            <span v-for="word in quietWords" :key="word" class="hero-glyph-note">{{ word }}</span>
          </div>
        </div>
      </section>

      <section class="home-notice-shell mx-auto max-w-7xl px-5 pb-14 sm:px-8">
        <div class="home-notice-ribbon" aria-label="庭前告示">
          <div class="notice-intro">
            <span class="notice-kicker">庭前告示</span>
            <p>
              <span class="notice-intro-line">不是只把多家能力并排陈列，</span>
              <span class="notice-intro-line">而是把接入、值守、计量与准入</span>
              <span class="notice-intro-line">整理成一套可长期使用的入口秩序。</span>
            </p>
          </div>

          <div class="notice-list" aria-label="首页能力锚点">
            <article
              v-for="item in trustAnchors"
              :key="item.title"
              class="notice-item"
            >
              <div class="notice-item-mark">
                <span class="notice-item-seal"></span>
                <span class="notice-item-index">{{ item.index }}</span>
              </div>
              <div class="notice-item-copy">
                <h2>{{ item.title }}</h2>
                <p :class="item.copyClass">
                  <template v-if="item.copyLines.length">
                    <span v-for="line in item.copyLines" :key="line" class="notice-copy-line">{{ line }}</span>
                  </template>
                  <template v-else>{{ item.copy }}</template>
                </p>
              </div>
            </article>
          </div>
        </div>
      </section>

      <section class="home-values-shell mx-auto max-w-7xl px-5 pb-12 sm:px-8">
        <div class="home-values-board">
          <article
            v-for="item in valueCards"
            :key="item.index"
            :class="['value-item', item.index === '01' ? 'value-item-primary' : '']"
          >
            <div class="value-index-row">
              <span class="value-dot"></span>
              <span class="value-index">{{ item.index }}</span>
            </div>
            <h2>{{ item.title }}</h2>
            <p :class="item.copyClass">{{ item.copy }}</p>
          </article>
        </div>
      </section>

      <section class="home-providers-shell mx-auto max-w-7xl px-5 pb-14 sm:px-8">
        <div class="home-provider-panel">
          <div class="home-capability-note">
            <p class="capability-copy">
              <span class="capability-dot" aria-hidden="true"></span>
              自研负载均衡、自动故障转移与智能路由，守住高并发下的稳定与低延迟。
            </p>
          </div>
          <div class="providers-meta">
            <div class="providers-kicker">已接入</div>
            <div class="providers-list">
              <span
                v-for="provider in connectedProviders"
                :key="provider"
                class="provider-item"
              >
                {{ provider }}
              </span>
            </div>
          </div>
        </div>
      </section>
    </main>

    <PublicSiteFooter />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import PublicSiteFooter from '@/components/layout/PublicSiteFooter.vue'
import PublicSiteHeader from '@/components/layout/PublicSiteHeader.vue'
import { useThemeState } from '@/utils/theme'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

const authStore = useAuthStore()
const appStore = useAppStore()
const isDark = useThemeState()

const configuredName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
const brandName = computed(() => configuredName.value && configuredName.value !== 'Sub2API' ? configuredName.value : '山枢庭')
const heroSubtitle = computed(() => {
  const subtitle = appStore.cachedPublicSettings?.site_subtitle?.trim()
  const legacySubtitles = new Set([
    'AI API Gateway Platform',
    'Subscription to API Conversion Platform',
  ])

  return subtitle && !legacySubtitles.has(subtitle)
    ? subtitle
    : '为长期使用而设，收束多源能力，保持稳定供给、清晰计量与审慎准入。'
})
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const resolvedDocsRoute = '/docs'
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const trimmedHomeContent = computed(() => homeContent.value.trim())
const hasHomeContent = computed(() => trimmedHomeContent.value.length > 0)
const isHomeContentUrl = computed(() => {
  const content = trimmedHomeContent.value
  return content.startsWith('http://') || content.startsWith('https://')
})
const docsLinkTag = computed(() => docUrl.value ? 'a' : 'router-link')
const docsLinkProps = computed(() => docUrl.value
  ? { href: docUrl.value, target: '_blank', rel: 'noopener noreferrer' }
  : { to: resolvedDocsRoute })

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const homeBackgroundStyle = computed(() => ({
  '--sst-home-bg': `url(${paperInkBg})`,
}))

const quietWords = [
  '稳定供给',
  '清晰计量',
  '审慎准入',
  '长期维护',
]

const trustAnchors = [
  {
    index: '甲',
    title: '故障转移',
    copy: '上游波动时先稳住入口连续性，把切换动作收在系统内部完成。',
    copyClass: '',
    copyLines: [],
  },
  {
    index: '乙',
    title: '统一账册',
    copy: '请求、消耗、余额与分组计量归在同一账册内，便于统一核算。',
    copyClass: '',
    copyLines: [],
  },
  {
    index: '丙',
    title: '准入分层',
    copy: '分组、密钥与风控边界分层落位，让权限与调用路径各守其序。',
    copyClass: 'notice-copy-two-line',
    copyLines: ['分组、密钥与风控边界分层落位，', '让权限与调用路径各守其序。'],
  },
  {
    index: '丁',
    title: '文档兼容',
    copy: '以文档与统一入口承接 OpenAI、Anthropic 等常见调用方式。',
    copyClass: '',
    copyLines: [],
  },
] as const

const valueCards = [
  {
    index: '01',
    title: '稳定供给',
    copy: '多源能力被收束成一个稳定、持续可用的统一入口。',
    copyClass: 'value-copy-single-line',
  },
  {
    index: '02',
    title: '核对清楚',
    copy: '余额、用量与调用明细对应明确，日常对账更省心。',
    copyClass: '',
  },
  {
    index: '03',
    title: '权限安静',
    copy: '权限、渠道与风控边界清楚，减少无效与噪声调用。',
    copyClass: '',
  },
  {
    index: '04',
    title: '可长期托付',
    copy: '链路可追踪、系统可维护，适合放在长期运行里使用。',
    copyClass: '',
  },
] as const

const connectedProviders = [
  'Anthropic',
  'OpenAI',
] as const

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.sst-home {
  position: relative;
  background: #f4efe4;
}

:global(html.dark) .sst-home,
.sst-home.is-dark {
  background: #0b0d0c;
}

.sst-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(90deg, rgba(244, 239, 228, 0.9) 0%, rgba(244, 239, 228, 0.48) 42%, rgba(244, 239, 228, 0.08) 100%),
    linear-gradient(180deg, rgba(244, 239, 228, 0.04) 0%, rgba(244, 239, 228, 0.18) 55%, rgba(244, 239, 228, 0.38) 100%),
    var(--sst-home-bg);
  background-size: cover, cover, cover;
  background-position: center, center, center bottom;
  opacity: 1;
}

:global(html.dark) .sst-bg,
.sst-home.is-dark .sst-bg {
  background-image:
    linear-gradient(90deg, rgba(8, 10, 9, 0.92) 0%, rgba(10, 12, 11, 0.58) 42%, rgba(10, 12, 11, 0.28) 100%),
    linear-gradient(180deg, rgba(10, 12, 11, 0.18) 0%, rgba(10, 12, 11, 0.24) 52%, rgba(10, 12, 11, 0.48) 100%),
    radial-gradient(circle at 16% 22%, rgba(176, 132, 81, 0.08), transparent 28%),
    radial-gradient(circle at 82% 18%, rgba(88, 59, 29, 0.08), transparent 24%),
    var(--sst-home-bg);
  background-size: cover, cover, cover;
  background-position: center, center, center bottom;
  filter: grayscale(0.18) brightness(0.62) sepia(0.08) saturate(0.82);
}

.sst-paper {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.16;
  background-image:
    linear-gradient(rgba(31, 35, 32, 0.022) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.018) 1px, transparent 1px);
  background-size: 128px 128px, 128px 128px;
}

.sst-paper::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.22;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(31, 35, 32, 0.06) 0 1px, transparent 1.5px),
    radial-gradient(circle at 72% 42%, rgba(31, 35, 32, 0.045) 0 1px, transparent 1.5px);
  background-size: 34px 41px, 48px 57px;
}

:global(html.dark) .sst-paper,
.sst-home.is-dark .sst-paper {
  opacity: 0.05;
}

:global(html.dark) .sst-paper::after,
.sst-home.is-dark .sst-paper::after {
  opacity: 0.06;
}

.home-title {
  max-width: 5.3em;
  text-shadow: 0 1px 0 rgba(244, 239, 228, 0.42);
}

.home-copy {
  max-width: 40rem;
  text-wrap: pretty;
}

:global(html.dark) .home-title,
.sst-home.is-dark .home-title {
  color: #f4e8d7;
  text-shadow: 0 1px 0 rgba(255, 243, 223, 0.06), 0 0 18px rgba(140, 92, 43, 0.06);
}

:global(html.dark) .home-lead,
.sst-home.is-dark .home-lead {
  color: #e6d8c2;
}

:global(html.dark) .home-copy,
:global(html.dark) .quiet-words,
.sst-home.is-dark .home-copy,
.sst-home.is-dark .quiet-words {
  color: #cdbca4;
}

.home-primary-cta {
  min-width: 8.8rem;
  min-height: 3.22rem;
  border: 1px solid rgba(235, 214, 182, 0.44);
  background: linear-gradient(135deg, rgba(246, 231, 208, 0.98), rgba(234, 218, 193, 0.96));
  color: #241c16;
  box-shadow: 0 18px 34px rgba(84, 57, 31, 0.14), inset 0 1px 0 rgba(255, 246, 232, 0.44);
}

.home-primary-cta:hover {
  background: linear-gradient(135deg, rgba(249, 238, 219, 1), rgba(239, 225, 202, 0.98));
}

.home-secondary-cta {
  min-width: 7.2rem;
  min-height: 3.22rem;
  border: 1px solid rgba(98, 84, 63, 0.22);
  background: rgba(54, 45, 35, 0.96);
  color: #f2e6d0;
  box-shadow: 0 16px 30px rgba(31, 35, 32, 0.12), inset 0 1px 0 rgba(255, 248, 236, 0.05);
}

.home-secondary-cta:hover {
  background: rgba(68, 55, 42, 0.98);
}

:global(html.dark) .home-primary-cta,
.sst-home.is-dark .home-primary-cta {
  border-color: rgba(214, 176, 122, 0.12);
  background: linear-gradient(135deg, rgba(236, 224, 204, 0.94), rgba(204, 183, 157, 0.9));
  color: #241c16;
}

:global(html.dark) .home-secondary-cta,
.sst-home.is-dark .home-secondary-cta {
  border-color: rgba(92, 84, 70, 0.5);
  background: rgba(15, 18, 16, 0.92);
  color: #ece0cb;
}

.hero-visual {
  isolation: isolate;
}

.hero-visual::before,
.hero-visual::after {
  content: '';
  position: absolute;
  display: none;
  pointer-events: none;
  z-index: 0;
}

.hero-courtyard-flow,
.hero-seal-imprint {
  display: none;
}

.hero-logo-stage {
  position: absolute;
  right: 8%;
  top: 8%;
  width: min(100%, 29.8rem);
  aspect-ratio: 2048 / 1490;
  transform: translateY(-0.25rem);
  overflow: visible;
  pointer-events: none;
}

.hero-logo-image {
  --hero-logo-opacity: 0.94;
  width: 100%;
  height: 100%;
  display: block;
  position: relative;
  z-index: 2;
  object-fit: contain;
  opacity: var(--hero-logo-opacity);
  filter:
    brightness(0.93)
    contrast(1.08)
    sepia(0.045)
    saturate(0.9)
    drop-shadow(0 1px 0 rgba(83, 76, 63, 0.18))
    drop-shadow(0 18px 32px rgba(63, 45, 28, 0.06));
}

.hero-seal-imprint img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.sst-home.is-dark .hero-logo-image {
  --hero-logo-opacity: 0.86;
  filter:
    brightness(0.88)
    contrast(0.98)
    sepia(0.08)
    saturate(0.86)
    drop-shadow(0 1px 0 rgba(246, 227, 195, 0.08))
    drop-shadow(0 20px 34px rgba(0, 0, 0, 0.18));
}

.hero-glyph-notes {
  position: absolute;
  right: 0.2rem;
  top: -0.15rem;
  display: grid;
  justify-items: end;
  gap: 0.64rem;
  z-index: 3;
}

.hero-glyph-note {
  color: #b8a781;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.82rem;
  letter-spacing: 0.06em;
  writing-mode: vertical-rl;
  opacity: 0.9;
}

.home-cta-row :deep(svg) {
  flex: 0 0 auto;
}

:global(html.dark) .hero-glyph-note,
.sst-home.is-dark .hero-glyph-note {
  color: #d9c3a1;
}

@media (min-width: 1024px) {
  .hero-visual::before {
    display: block;
    left: 3%;
    right: 8%;
    top: 55%;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(161, 139, 106, 0.26) 18%, rgba(167, 58, 42, 0.16) 54%, transparent);
    opacity: 0;
    transform: scaleX(0.2);
    transform-origin: left center;
    animation: heroAxisReveal 1500ms cubic-bezier(0.19, 1, 0.22, 1) 220ms forwards;
  }

  .hero-visual::after {
    display: block;
    top: 8%;
    bottom: 8%;
    left: 58%;
    width: 1px;
    background: linear-gradient(180deg, transparent, rgba(161, 139, 106, 0.18) 18%, rgba(167, 58, 42, 0.09) 50%, transparent);
    opacity: 0;
    transform: scaleY(0.18);
    transform-origin: top center;
    animation: heroAxisRevealY 1800ms cubic-bezier(0.19, 1, 0.22, 1) 360ms forwards;
  }

  .hero-courtyard-flow {
    position: absolute;
    inset: 0;
    display: block;
    pointer-events: none;
    z-index: 1;
  }

  .hero-courtyard-flow::before {
    content: '';
    position: absolute;
    left: 14%;
    top: calc(55% - 0.18rem);
    width: 0.36rem;
    height: 0.36rem;
    border-radius: 999px;
    background: #a73a2a;
    box-shadow: 0 0 0 4px rgba(167, 58, 42, 0.032), 0 0 14px rgba(167, 58, 42, 0.08);
    opacity: 0;
    animation: heroCourtyardFlow 9.6s ease-in-out 1800ms infinite;
  }

  .hero-logo-stage {
    animation: heroLogoDrift 15s ease-in-out infinite;
    will-change: transform;
  }

  .hero-logo-image {
    animation: heroLogoSettle 1100ms cubic-bezier(0.19, 1, 0.22, 1) both;
  }

  .hero-seal-imprint {
    position: absolute;
    right: 8.8%;
    bottom: 7.5%;
    z-index: 4;
    display: block;
    width: 3.18rem;
    height: 3.18rem;
    opacity: 0;
    transform: translate3d(0.42rem, -0.28rem, 0) scale(1.12) rotate(-7deg);
    filter: drop-shadow(0 10px 18px rgba(86, 37, 24, 0.12));
    animation: heroSealImprint 920ms cubic-bezier(0.18, 0.86, 0.24, 1) 760ms forwards;
  }

  .hero-glyph-notes {
    animation: heroLogoNotes 12s ease-in-out infinite;
  }

  .hero-glyph-note {
    opacity: 0;
    animation: heroGlyphIn 900ms ease-out forwards;
  }

  .hero-glyph-note:nth-child(1) {
    animation-delay: 520ms;
  }

  .hero-glyph-note:nth-child(2) {
    animation-delay: 720ms;
  }

  .hero-glyph-note:nth-child(3) {
    animation-delay: 920ms;
  }

  .hero-glyph-note:nth-child(4) {
    animation-delay: 1120ms;
  }

  .sst-home.is-dark .hero-courtyard-flow::before {
    background: #b87f42;
    box-shadow: 0 0 0 5px rgba(184, 127, 66, 0.055), 0 0 18px rgba(184, 127, 66, 0.15);
  }
}

.home-notice-shell {
  position: relative;
}

.home-notice-ribbon {
  position: relative;
  display: grid;
  grid-template-columns: minmax(13.5rem, 0.62fr) minmax(0, 2.38fr);
  align-items: stretch;
  border: 1px solid rgba(154, 128, 92, 0.15);
  border-radius: 1.02rem;
  background:
    linear-gradient(180deg, rgba(252, 248, 240, 0.92), rgba(244, 235, 220, 0.72)),
    linear-gradient(90deg, rgba(144, 113, 76, 0.022), transparent 18%, rgba(144, 113, 76, 0.015) 82%, transparent),
    rgba(255, 255, 255, 0.28);
  box-shadow:
    0 10px 24px rgba(84, 57, 31, 0.035),
    inset 0 1px 0 rgba(255, 249, 239, 0.52),
    inset 0 -1px 0 rgba(140, 111, 76, 0.05);
  overflow: hidden;
}

.home-notice-ribbon::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    linear-gradient(180deg, rgba(167, 58, 42, 0.028), transparent 22%),
    repeating-linear-gradient(90deg, transparent 0 28px, rgba(128, 98, 66, 0.018) 28px 29px),
    radial-gradient(circle at 78% 18%, rgba(173, 134, 78, 0.08), transparent 18%);
}

.notice-intro,
.notice-list {
  position: relative;
  z-index: 1;
}

.notice-intro {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 0.52rem;
  padding: 0.98rem 1.42rem 1.04rem 1.48rem;
  border-right: 1px solid rgba(161, 139, 106, 0.11);
}

.notice-kicker {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: #9a6a36;
  font-size: 0.66rem;
  letter-spacing: 0.34em;
  text-transform: uppercase;
}

.notice-kicker::before {
  content: '';
  width: 0.34rem;
  height: 0.34rem;
  border-radius: 999px;
  background: #a73a2a;
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.08);
}

.notice-intro p {
  max-width: none;
  color: #5d5549;
  font-size: 0.86rem;
  line-height: 1.82;
  letter-spacing: 0.01em;
}

.notice-intro-line {
  display: block;
  white-space: nowrap;
}

.notice-list {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.notice-item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  align-items: center;
  gap: 0.74rem;
  min-height: 100%;
  padding: 0.92rem 1rem 0.96rem;
  border-left: 1px solid rgba(161, 139, 106, 0.11);
  transition: background-color 180ms ease, box-shadow 180ms ease;
}

.notice-item:first-child {
  border-left: 0;
}

.notice-item:hover {
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.28), rgba(247, 239, 226, 0.14));
  box-shadow:
    inset 0 1px 0 rgba(255, 252, 246, 0.34),
    inset 0 -1px 0 rgba(167, 58, 42, 0.06);
}

.notice-item-mark {
  display: grid;
  justify-items: center;
  gap: 0.28rem;
  min-width: 1.18rem;
  padding-top: 0.02rem;
}

.notice-item-seal {
  width: 0.44rem;
  height: 0.44rem;
  border-radius: 0.08rem;
  background:
    linear-gradient(135deg, rgba(186, 72, 54, 0.96), rgba(145, 39, 24, 0.96));
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.055);
}

.notice-item-index {
  color: #af7a43;
  font-family: 'Geist Mono', 'JetBrains Mono', monospace;
  font-size: 0.62rem;
  letter-spacing: 0.16em;
}

.notice-item-copy h2 {
  color: #27221c;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 0.92rem;
  line-height: 1.14;
  letter-spacing: 0.02em;
  text-wrap: balance;
}

.notice-item-copy p {
  margin-top: 0.22rem;
  max-width: 17.2em;
  color: #6b6154;
  font-size: 0.74rem;
  line-height: 1.66;
  letter-spacing: 0.01em;
  text-wrap: pretty;
}

.notice-copy-two-line {
  max-width: none;
}

.notice-copy-line {
  display: block;
  white-space: nowrap;
}

.home-values-shell,
.home-providers-shell {
  position: relative;
}

.home-values-board {
  display: grid;
  grid-template-columns: minmax(0, 2.28fr) repeat(3, minmax(0, 1fr));
  gap: 0;
  border-top: 1px solid rgba(161, 139, 106, 0.08);
}

.value-item {
  padding: 1.18rem 1.2rem 1.2rem;
  text-align: left;
}

.value-item-primary {
  position: relative;
  padding-right: 1.08rem;
}

.value-item-primary::after {
  content: '';
  position: absolute;
  left: 1.2rem;
  right: 2rem;
  bottom: 0;
  height: 1px;
  background: linear-gradient(90deg, rgba(167, 58, 42, 0.18), rgba(161, 139, 106, 0.02));
  opacity: 0.55;
}

.value-item + .value-item {
  border-left: 1px solid rgba(161, 139, 106, 0.08);
}

.value-index-row {
  display: flex;
  align-items: center;
  gap: 0.48rem;
}

.value-dot {
  width: 0.38rem;
  height: 0.38rem;
  border-radius: 999px;
  background: #a73a2a;
}

.value-index {
  color: #af7840;
  font-family: 'Geist Mono', 'JetBrains Mono', monospace;
  font-size: 0.68rem;
  letter-spacing: 0.12em;
}

.value-item h2 {
  margin-top: 0.56rem;
  color: #2a241e;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(1.2rem, 1.3vw, 1.48rem);
  line-height: 1.12;
  letter-spacing: 0.01em;
}

.value-item p {
  margin-top: 0.38rem;
  max-width: 16.8em;
  color: #71675c;
  font-size: 0.82rem;
  line-height: 1.68;
  letter-spacing: 0.01em;
  text-wrap: pretty;
}

.value-item-primary h2 {
  font-size: clamp(1.74rem, 2.3vw, 2.18rem);
  line-height: 1.02;
  letter-spacing: 0.015em;
}

.value-item-primary p {
  margin-top: 0.52rem;
  max-width: 22em;
  font-size: 0.94rem;
  line-height: 1.78;
  color: #675d52;
}

.value-item-primary p.value-copy-single-line {
  font-size: clamp(0.72rem, 0.88vw, 0.78rem);
  letter-spacing: 0;
  max-width: none;
  white-space: nowrap;
  text-wrap: nowrap;
}

.home-provider-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.9rem 1.6rem;
  padding: 0.96rem 0.1rem 0;
  border-top: 1px solid rgba(161, 139, 106, 0.08);
}

.home-capability-note {
  margin: 0;
  text-align: left;
}

.capability-copy {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: #6c6054;
  font-size: 0.81rem;
  line-height: 1.72;
  letter-spacing: 0.012em;
  text-wrap: pretty;
}

.capability-dot {
  width: 0.32rem;
  height: 0.32rem;
  flex: 0 0 auto;
  border-radius: 999px;
  background: #a73a2a;
}

.providers-meta {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.7rem 1rem;
  flex-wrap: wrap;
}

.providers-kicker {
  color: #a07a4f;
  font-size: 0.64rem;
  letter-spacing: 0.38em;
  text-transform: uppercase;
}

.providers-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 0.6rem 1.2rem;
}

.provider-item {
  color: #6f6357;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: clamp(0.9rem, 0.96vw, 1.02rem);
  line-height: 1.12;
  letter-spacing: 0.02em;
  white-space: nowrap;
}

.provider-item::before {
  content: '◦';
  margin-right: 0.46rem;
  color: #a96a31;
}

:global(html.dark) .home-notice-ribbon,
.sst-home.is-dark .home-notice-ribbon {
  border-color: rgba(125, 100, 66, 0.18);
  background:
    linear-gradient(135deg, rgba(18, 20, 17, 0.72), rgba(28, 24, 20, 0.52)),
    rgba(16, 18, 16, 0.44);
  box-shadow:
    0 18px 40px rgba(0, 0, 0, 0.14),
    inset 0 1px 0 rgba(245, 225, 194, 0.035);
}

:global(html.dark) .home-notice-ribbon::before,
.sst-home.is-dark .home-notice-ribbon::before {
  background:
    linear-gradient(180deg, rgba(167, 58, 42, 0.04), transparent 24%),
    radial-gradient(circle at 82% 18%, rgba(176, 124, 58, 0.08), transparent 20%);
}

:global(html.dark) .notice-intro,
.sst-home.is-dark .notice-intro,
:global(html.dark) .notice-item,
.sst-home.is-dark .notice-item {
  border-color: rgba(125, 100, 66, 0.18);
}

:global(html.dark) .notice-kicker,
.sst-home.is-dark .notice-kicker,
:global(html.dark) .notice-item-index,
.sst-home.is-dark .notice-item-index {
  color: #cfa86a;
}

:global(html.dark) .notice-intro p,
.sst-home.is-dark .notice-intro p,
:global(html.dark) .notice-item-copy p,
.sst-home.is-dark .notice-item-copy p {
  color: #c8b59b;
}

:global(html.dark) .notice-item-copy h2,
.sst-home.is-dark .notice-item-copy h2 {
  color: #f4e7d6;
}

:global(html.dark) .notice-item-seal,
.sst-home.is-dark .notice-item-seal {
  background: #b87f42;
  box-shadow: 0 0 0 3px rgba(184, 127, 66, 0.08);
}

:global(html.dark) .notice-item:hover,
.sst-home.is-dark .notice-item:hover {
  background: linear-gradient(180deg, rgba(245, 225, 194, 0.04), rgba(255, 247, 235, 0.01));
  box-shadow: inset 0 -1px 0 rgba(184, 127, 66, 0.12);
}

:global(html.dark) .value-item,
.sst-home.is-dark .value-item,
:global(html.dark) .home-values-board,
.sst-home.is-dark .home-values-board,
:global(html.dark) .home-provider-panel,
.sst-home.is-dark .home-provider-panel {
  border-color: rgba(125, 100, 66, 0.16);
}

:global(html.dark) .value-index,
.sst-home.is-dark .value-index,
:global(html.dark) .providers-kicker,
.sst-home.is-dark .providers-kicker {
  color: #cfa86a;
}

:global(html.dark) .value-item h2,
.sst-home.is-dark .value-item h2,
:global(html.dark) .provider-item,
.sst-home.is-dark .provider-item {
  color: #f0dfc7;
}

:global(html.dark) .value-item-primary::after,
.sst-home.is-dark .value-item-primary::after {
  background: linear-gradient(90deg, rgba(184, 127, 66, 0.28), rgba(184, 127, 66, 0.02));
  opacity: 0.75;
}

:global(html.dark) .value-item p,
.sst-home.is-dark .value-item p,
:global(html.dark) .capability-copy,
.sst-home.is-dark .capability-copy {
  color: #c8b59b;
}

:global(html.dark) .value-dot,
.sst-home.is-dark .value-dot,
:global(html.dark) .capability-dot,
.sst-home.is-dark .capability-dot,
:global(html.dark) .provider-item::before,
.sst-home.is-dark .provider-item::before {
  background: #b87f42;
  color: #b87f42;
}

@keyframes heroAxisReveal {
  0% {
    opacity: 0;
    transform: scaleX(0.2);
  }

  100% {
    opacity: 1;
    transform: scaleX(1);
  }
}

@keyframes heroAxisRevealY {
  0% {
    opacity: 0;
    transform: scaleY(0.18);
  }

  100% {
    opacity: 1;
    transform: scaleY(1);
  }
}

@keyframes heroCourtyardFlow {
  0%,
  18% {
    opacity: 0;
    transform: translate3d(0, 0, 0) scale(0.72);
  }

  28% {
    opacity: 0.42;
    transform: translate3d(6.5rem, 0, 0) scale(1);
  }

  70% {
    opacity: 0.3;
    transform: translate3d(23rem, 0, 0) scale(0.88);
  }

  100% {
    opacity: 0;
    transform: translate3d(31rem, 0, 0) scale(0.72);
  }
}

@keyframes heroLogoSettle {
  0% {
    transform: translate3d(0.8rem, 0.7rem, 0) scale(0.985);
  }

  100% {
    transform: translate3d(0, 0, 0) scale(1);
  }
}

@keyframes heroLogoDrift {
  0%,
  100% {
    transform: translate3d(0, -0.25rem, 0) rotate(-0.18deg);
  }

  50% {
    transform: translate3d(-0.28rem, 0.12rem, 0) rotate(0.24deg);
  }
}

@keyframes heroSealImprint {
  0% {
    opacity: 0;
    transform: translate3d(0.42rem, -0.28rem, 0) scale(1.12) rotate(-7deg);
  }

  62% {
    opacity: 0.9;
    transform: translate3d(0, 0, 0) scale(0.96) rotate(-7deg);
  }

  100% {
    opacity: 0.78;
    transform: translate3d(0, 0, 0) scale(1) rotate(-7deg);
  }
}

@keyframes heroGlyphIn {
  0% {
    opacity: 0;
    transform: translate3d(0, 0.42rem, 0);
  }

  100% {
    opacity: 0.9;
    transform: translate3d(0, 0, 0);
  }
}

@keyframes heroLogoNotes {
  0%,
  100% {
    transform: translate3d(0, 0, 0);
  }

  50% {
    transform: translate3d(0, 0.38rem, 0);
  }
}

@media (max-width: 1023px) {
  .home-notice-ribbon {
    grid-template-columns: 1fr;
  }

  .notice-intro {
    border-right: 0;
    border-bottom: 1px solid rgba(161, 139, 106, 0.12);
  }

  .notice-list,
  .home-values-board {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .value-item-primary {
    grid-column: 1 / -1;
    padding-right: 0.2rem;
  }

  .value-item-primary::after {
    right: 0.2rem;
  }

  .notice-item:nth-child(2n + 1),
  .value-item:nth-child(2n + 1) {
    border-left: 0;
  }

  .notice-item:nth-child(n + 3),
  .value-item:nth-child(n + 3) {
    border-top: 1px solid rgba(161, 139, 106, 0.08);
  }

  .hero-visual {
    min-height: 17.5rem;
  }

  .hero-logo-stage {
    right: 4%;
    top: 12%;
    width: min(88%, 18.8rem);
    transform: none;
  }

  .hero-glyph-notes {
    right: 0.2rem;
    top: 0.6rem;
    gap: 0.55rem;
  }

  .hero-glyph-note {
    font-size: 0.72rem;
  }
}

@media (max-width: 640px) {
  .home-hero-shell {
    gap: 2.75rem;
    padding-top: 1.8rem;
    padding-bottom: 1.2rem;
  }

  .home-notice-shell {
    padding-bottom: 1rem;
  }

  .home-title {
    max-width: none;
  }

  .home-lead {
    margin-top: 1.05rem;
    font-size: 1.72rem;
    line-height: 1.12;
  }

  .home-copy {
    margin-top: 0.92rem;
    font-size: 0.96rem;
    line-height: 1.7;
    max-width: 20rem;
  }

  .home-primary-cta,
  .home-secondary-cta {
    width: 100%;
  }

  .hero-visual {
    min-height: 11rem;
    margin-top: 0.35rem;
  }

  .home-notice-ribbon {
    border-radius: 1.12rem;
  }

  .notice-intro {
    gap: 0.58rem;
    padding: 0.94rem 0.98rem 0.98rem 1.02rem;
  }

  .notice-intro p {
    max-width: none;
    font-size: 0.84rem;
    line-height: 1.66;
    text-wrap: pretty;
  }

  .notice-intro-line,
  .notice-copy-line {
    white-space: normal;
  }

  .notice-list {
    grid-template-columns: 1fr;
  }

  .home-values-board {
    grid-template-columns: 1fr;
  }

  .notice-item {
    border-left: 0;
    padding: 0.8rem 0.94rem 0.84rem;
    gap: 0.62rem;
  }

  .notice-item:nth-child(n + 2),
  .value-item:nth-child(n + 2) {
    border-top: 1px solid rgba(161, 139, 106, 0.12);
  }

  .notice-item-copy h2 {
    font-size: 0.9rem;
  }

  .notice-item-copy p {
    max-width: none;
    font-size: 0.76rem;
    line-height: 1.58;
  }

  .value-item {
    border-left: 0;
    padding: 0.94rem 0.2rem 0.98rem 0;
  }

  .value-item-primary::after {
    left: 0;
    right: 0;
  }

  .value-item h2 {
    font-size: 1.14rem;
  }

  .value-item-primary h2 {
    font-size: 1.52rem;
  }

  .value-item-primary p {
    font-size: 0.9rem;
    line-height: 1.72;
  }

  .value-item p,
  .capability-copy {
    max-width: none;
  }

  .notice-copy-two-line {
    white-space: normal;
    text-wrap: pretty;
  }

  .value-item-primary p.value-copy-single-line {
    font-size: clamp(0.72rem, 3.25vw, 0.82rem);
    white-space: nowrap;
    text-wrap: nowrap;
  }

  .home-provider-panel,
  .providers-meta {
    align-items: flex-start;
    flex-direction: column;
  }

  .hero-logo-stage {
    left: 50%;
    right: auto;
    top: auto;
    bottom: 0;
    width: min(72vw, 13rem);
    transform: translateX(-50%);
  }

  .hero-glyph-notes {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-visual::before,
  .hero-visual::after,
  .hero-courtyard-flow::before,
  .hero-logo-stage,
  .hero-logo-image,
  .hero-seal-imprint,
  .hero-glyph-note,
  .hero-glyph-notes {
    animation: none;
  }

  .hero-visual::before,
  .hero-visual::after {
    opacity: 0.72;
    transform: none;
  }

  .hero-seal-imprint,
  .hero-glyph-note {
    opacity: 0.86;
  }
}
</style>
