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
    <div class="sst-paper-fibers" aria-hidden="true"></div>

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
            {{ t('publicSite.tagline') }}
          </p>
          <p class="home-copy mt-6 max-w-[40rem] text-base leading-8 text-zen-mist dark:text-zen-stone sm:text-lg">
            {{ heroSubtitle }}
          </p>

          <div class="home-cta-row mt-10 flex flex-col gap-3 sm:flex-row">
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="home-primary-cta inline-flex items-center justify-center rounded-zen px-7 py-3 text-sm font-medium"
            >
              {{ t('publicSite.enter') }}
              <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
            </router-link>
            <component
              :is="docsLinkTag"
              v-bind="docsLinkProps"
              class="home-secondary-cta inline-flex items-center justify-center rounded-zen px-7 py-3 text-sm font-medium"
            >
              {{ t('publicSite.nav.docs') }}
            </component>
          </div>
        </div>

        <div class="hero-visual relative lg:min-h-[34rem]" aria-hidden="true">
          <div class="hero-logo-stage">
            <GoldfishScene :dark="isDark" />
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
        <div class="home-notice-ribbon" :aria-label="t('publicHome.notice.aria')">
          <div class="notice-intro">
            <span class="notice-kicker">{{ t('publicHome.notice.kicker') }}</span>
            <p>
              <span v-for="line in noticeLines" :key="line" class="notice-intro-line">{{ line }}</span>
            </p>
          </div>

          <div class="notice-list" :aria-label="t('publicHome.notice.listAria')">
            <component
              :is="item.to ? 'router-link' : 'article'"
              v-for="item in trustAnchors"
              :key="item.title"
              :to="item.to"
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
            </component>
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
              {{ t('publicHome.capabilityCopy') }}
            </p>
          </div>
          <div class="providers-meta">
            <div class="providers-kicker">{{ t('publicHome.providersKicker') }}</div>
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
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import GoldfishScene from '@/components/brand/GoldfishScene.vue'
import PublicSiteFooter from '@/components/layout/PublicSiteFooter.vue'
import PublicSiteHeader from '@/components/layout/PublicSiteHeader.vue'
import { IMAGE_WORKSHOP_MENU_ID, findImageWorkshopMenuItem } from '@/utils/imageWorkshop'
import { useThemeState } from '@/utils/theme'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

const authStore = useAuthStore()
const appStore = useAppStore()
const { t, tm, locale } = useI18n()
const isDark = useThemeState()

const configuredName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
const brandName = computed(() => configuredName.value && configuredName.value !== 'Sub2API' ? configuredName.value : '山枢庭')
const heroSubtitle = computed(() => {
  const subtitle = appStore.cachedPublicSettings?.site_subtitle?.trim()
  const legacySubtitles = new Set([
    'AI API Gateway Platform',
    'Subscription to API Conversion Platform',
    '统一入口，安静流转。',
  ])

  if (locale.value.startsWith('en') && subtitle === '统一入口，安静流转。') {
    return t('publicHome.heroSubtitle')
  }

  return subtitle && !legacySubtitles.has(subtitle)
    ? subtitle
    : t('publicHome.heroSubtitle')
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
const imageWorkshopMenuItem = computed(() => findImageWorkshopMenuItem(appStore.cachedPublicSettings?.custom_menu_items))
const imageWorkshopPath = computed(() => `/custom/${IMAGE_WORKSHOP_MENU_ID}`)

const quietWords = computed(() => tm('publicHome.quietWords') as unknown as string[])
const noticeLines = computed(() => tm('publicHome.notice.lines') as unknown as string[])

interface TrustAnchor {
  index: string
  title: string
  copy: string
  copyClass: string
  copyLines: string[]
  to?: string
}

const baseTrustAnchors = computed<TrustAnchor[]>(() => tm('publicHome.trustAnchors') as unknown as TrustAnchor[])

const trustAnchors = computed(() => {
  const anchors = [...baseTrustAnchors.value]
  if (imageWorkshopMenuItem.value) {
    anchors.splice(2, 0, {
      index: '丙',
      title: t('publicHome.imageWorkshop.title'),
      copy: t('publicHome.imageWorkshop.copy'),
      copyClass: 'notice-copy-two-line',
      copyLines: tm('publicHome.imageWorkshop.copyLines') as unknown as string[],
      to: imageWorkshopPath.value,
    })
  }
  return anchors.map((item, index) => ({
    ...item,
    index: ['甲', '乙', '丙', '丁', '戊'][index] ?? item.index,
  }))
})

const valueCards = computed(() => tm('publicHome.valueCards') as unknown as Array<{
  index: string
  title: string
  copy: string
  copyClass: string
}>)

const connectedProviders = computed(() => [
  'Anthropic',
  'OpenAI',
  ...(imageWorkshopMenuItem.value ? [t('publicHome.imageWorkshop.title')] : []),
])

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
  background:
    radial-gradient(circle at 18% 8%, rgba(83, 59, 33, 0.16), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(106, 54, 28, 0.12), transparent 30%),
    #0c0d0b;
}

.sst-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    radial-gradient(ellipse at 18% 10%, rgba(255, 253, 247, 0.46), transparent 36%),
    radial-gradient(ellipse at 82% 20%, rgba(226, 209, 176, 0.22), transparent 34%),
    linear-gradient(90deg, rgba(244, 239, 228, 0.94) 0%, rgba(244, 239, 228, 0.58) 42%, rgba(244, 239, 228, 0.14) 100%),
    linear-gradient(180deg, rgba(251, 246, 236, 0.2) 0%, rgba(244, 239, 228, 0.18) 55%, rgba(231, 219, 200, 0.44) 100%),
    var(--sst-home-bg);
  background-size: cover, cover, cover, cover, cover;
  background-position: center, center, center, center, center bottom;
  opacity: 1;
}

:global(html.dark) .sst-bg,
.sst-home.is-dark .sst-bg {
  background-image:
    radial-gradient(ellipse at 18% 12%, rgba(222, 180, 116, 0.1), transparent 34%),
    radial-gradient(ellipse at 82% 12%, rgba(164, 82, 42, 0.13), transparent 30%),
    linear-gradient(90deg, rgba(8, 10, 9, 0.96) 0%, rgba(10, 12, 10, 0.72) 42%, rgba(10, 12, 10, 0.34) 100%),
    linear-gradient(180deg, rgba(13, 15, 12, 0.22) 0%, rgba(13, 15, 12, 0.32) 52%, rgba(8, 9, 8, 0.7) 100%),
    var(--sst-home-bg);
  background-size: cover, cover, cover, cover, cover;
  background-position: center, center, center, center, center bottom;
  filter: grayscale(0.08) brightness(0.58) sepia(0.12) saturate(0.95);
}

.sst-paper {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.22;
  background-image:
    linear-gradient(rgba(90, 77, 57, 0.026) 1px, transparent 1px),
    linear-gradient(90deg, rgba(90, 77, 57, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, transparent 0 49.6%, rgba(144, 113, 76, 0.04) 49.8% 50%, transparent 50.2% 100%);
  background-size: 132px 132px, 132px 132px, min(78rem, 92vw) 100%;
  background-position: center top, center top, center top;
}

.sst-paper::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.3;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(92, 74, 48, 0.07) 0 0.8px, transparent 1.4px),
    radial-gradient(circle at 72% 42%, rgba(92, 74, 48, 0.05) 0 0.9px, transparent 1.5px),
    linear-gradient(90deg, rgba(87, 67, 42, 0.055), transparent 5%, transparent 95%, rgba(87, 67, 42, 0.045)),
    linear-gradient(180deg, rgba(255, 254, 250, 0.24), transparent 15%, transparent 76%, rgba(126, 99, 63, 0.08));
  background-size: 34px 41px, 48px 57px, 100% 100%, 100% 100%;
}

.sst-paper-fibers {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.18;
  mix-blend-mode: multiply;
  background-image:
    repeating-linear-gradient(7deg, transparent 0 12px, rgba(113, 93, 62, 0.035) 12px 13px, transparent 13px 31px),
    repeating-linear-gradient(96deg, transparent 0 18px, rgba(255, 255, 255, 0.24) 18px 19px, transparent 19px 44px);
  background-size: 240px 180px, 180px 220px;
}

:global(html.dark) .sst-paper,
.sst-home.is-dark .sst-paper {
  opacity: 0.11;
  background-image:
    linear-gradient(rgba(239, 217, 179, 0.024) 1px, transparent 1px),
    linear-gradient(90deg, rgba(239, 217, 179, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, transparent 0 49.6%, rgba(217, 165, 94, 0.06) 49.8% 50%, transparent 50.2% 100%);
}

:global(html.dark) .sst-paper::after,
.sst-home.is-dark .sst-paper::after {
  opacity: 0.13;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(241, 217, 177, 0.06) 0 0.8px, transparent 1.4px),
    radial-gradient(circle at 72% 42%, rgba(241, 217, 177, 0.04) 0 0.9px, transparent 1.5px),
    linear-gradient(90deg, rgba(255, 226, 176, 0.045), transparent 5%, transparent 95%, rgba(255, 226, 176, 0.036)),
    linear-gradient(180deg, rgba(255, 236, 202, 0.05), transparent 16%, transparent 76%, rgba(0, 0, 0, 0.26));
}

:global(html.dark) .sst-paper-fibers,
.sst-home.is-dark .sst-paper-fibers {
  opacity: 0.09;
  mix-blend-mode: screen;
  background-image:
    repeating-linear-gradient(7deg, transparent 0 12px, rgba(236, 203, 154, 0.04) 12px 13px, transparent 13px 31px),
    repeating-linear-gradient(96deg, transparent 0 18px, rgba(255, 238, 210, 0.05) 18px 19px, transparent 19px 44px);
}

.home-title {
  max-width: 5.3em;
  text-shadow: 0 1px 0 rgba(244, 239, 228, 0.42);
}

.home-copy {
  max-width: 48rem;
  text-wrap: pretty;
  overflow-wrap: anywhere;
}

@media (min-width: 1024px) {
  .home-copy {
    max-width: 44rem;
  }
}

:global(html.dark) .home-title,
.sst-home.is-dark .home-title {
  color: #f7ead4;
  text-shadow:
    0 1px 0 rgba(255, 243, 223, 0.08),
    0 18px 44px rgba(0, 0, 0, 0.38),
    0 0 22px rgba(176, 120, 57, 0.08);
}

:global(html.dark) .home-lead,
.sst-home.is-dark .home-lead {
  color: #eadbc4;
}

:global(html.dark) .home-copy,
:global(html.dark) .quiet-words,
.sst-home.is-dark .home-copy,
.sst-home.is-dark .quiet-words {
  color: #d4c4ad;
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
  border-color: rgba(236, 208, 166, 0.22);
  background:
    linear-gradient(135deg, rgba(242, 229, 207, 0.98), rgba(208, 185, 154, 0.94)),
    radial-gradient(circle at 82% 20%, rgba(180, 104, 49, 0.14), transparent 28%);
  color: #241c16;
  box-shadow:
    0 18px 38px rgba(0, 0, 0, 0.26),
    0 0 0 1px rgba(255, 238, 210, 0.04),
    inset 0 1px 0 rgba(255, 248, 236, 0.3);
}

:global(html.dark) .home-secondary-cta,
.sst-home.is-dark .home-secondary-cta {
  border-color: rgba(142, 118, 86, 0.48);
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.94), rgba(14, 16, 13, 0.96)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.1), transparent 26%);
  color: #efe2cf;
  box-shadow:
    inset 0 1px 0 rgba(255, 238, 210, 0.06),
    0 14px 30px rgba(0, 0, 0, 0.2);
}

.hero-visual {
  isolation: isolate;
}

.hero-seal-imprint {
  display: none;
}

.hero-logo-stage {
  position: absolute;
  right: 8%;
  top: 8%;
  width: min(100%, 31.5rem);
  aspect-ratio: 1.48;
  transform: translateY(-0.25rem);
  overflow: visible;
  pointer-events: none;
}

.hero-seal-imprint img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
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
  .hero-logo-stage {
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
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    linear-gradient(90deg, rgba(144, 113, 76, 0.038), transparent 18%, rgba(144, 113, 76, 0.024) 82%, transparent),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28);
  box-shadow:
    0 14px 34px rgba(84, 57, 31, 0.05),
    inset 0 1px 0 rgba(255, 249, 239, 0.6),
    inset 0 -1px 0 rgba(140, 111, 76, 0.07),
    inset 0 0 0 1px rgba(255, 255, 255, 0.22);
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
    radial-gradient(circle at 78% 18%, rgba(173, 134, 78, 0.08), transparent 18%),
    radial-gradient(circle at 8% 12%, rgba(255, 255, 255, 0.2), transparent 22%);
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
  max-width: 100%;
  overflow-wrap: anywhere;
  text-wrap: balance;
}

.notice-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(9.5rem, 1fr));
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
  max-width: 100%;
  overflow-wrap: anywhere;
  text-wrap: balance;
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
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.16), transparent 42%),
    linear-gradient(90deg, rgba(144, 113, 76, 0.026), transparent 30%, transparent 70%, rgba(144, 113, 76, 0.018));
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
  font-size: clamp(0.74rem, 0.88vw, 0.82rem);
  letter-spacing: 0;
  max-width: none;
  overflow-wrap: anywhere;
  text-wrap: pretty;
}

.home-provider-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.9rem 1.6rem;
  padding: 0.96rem 0.1rem 0;
  border-top: 1px solid rgba(161, 139, 106, 0.08);
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.12), transparent);
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
  border-color: rgba(155, 126, 86, 0.26);
  background:
    linear-gradient(135deg, rgba(24, 27, 22, 0.82), rgba(35, 29, 23, 0.62)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px),
    rgba(13, 15, 13, 0.5);
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055),
    inset 0 -1px 0 rgba(130, 88, 46, 0.12);
}

:global(html.dark) .home-notice-ribbon::before,
.sst-home.is-dark .home-notice-ribbon::before {
  background:
    linear-gradient(180deg, rgba(188, 93, 31, 0.055), transparent 24%),
    repeating-linear-gradient(90deg, transparent 0 28px, rgba(224, 188, 132, 0.018) 28px 29px),
    radial-gradient(circle at 82% 18%, rgba(198, 139, 70, 0.12), transparent 20%),
    radial-gradient(circle at 8% 12%, rgba(255, 235, 196, 0.045), transparent 22%);
}

:global(html.dark) .notice-intro,
.sst-home.is-dark .notice-intro,
:global(html.dark) .notice-item,
.sst-home.is-dark .notice-item {
  border-color: rgba(155, 126, 86, 0.24);
}

:global(html.dark) .notice-kicker,
.sst-home.is-dark .notice-kicker,
:global(html.dark) .notice-item-index,
.sst-home.is-dark .notice-item-index {
  color: #d9b372;
}

:global(html.dark) .notice-intro p,
.sst-home.is-dark .notice-intro p,
:global(html.dark) .notice-item-copy p,
.sst-home.is-dark .notice-item-copy p {
  color: #d2bea2;
}

:global(html.dark) .notice-item-copy h2,
.sst-home.is-dark .notice-item-copy h2 {
  color: #f6e8d2;
}

:global(html.dark) .notice-item-seal,
.sst-home.is-dark .notice-item-seal {
  background: linear-gradient(135deg, #c27e4a, #9e452d);
  box-shadow: 0 0 0 3px rgba(194, 126, 74, 0.09), 0 0 14px rgba(194, 126, 74, 0.1);
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
  border-color: rgba(155, 126, 86, 0.22);
  background:
    linear-gradient(180deg, rgba(255, 226, 184, 0.035), transparent 44%),
    linear-gradient(90deg, rgba(206, 151, 82, 0.025), transparent 28%, transparent 72%, rgba(206, 151, 82, 0.02));
}

:global(html.dark) .value-index,
.sst-home.is-dark .value-index,
:global(html.dark) .providers-kicker,
.sst-home.is-dark .providers-kicker {
  color: #d8b171;
}

:global(html.dark) .value-item h2,
.sst-home.is-dark .value-item h2,
:global(html.dark) .provider-item,
.sst-home.is-dark .provider-item {
  color: #f3e1c7;
}

:global(html.dark) .value-item-primary::after,
.sst-home.is-dark .value-item-primary::after {
  background: linear-gradient(90deg, rgba(198, 139, 70, 0.36), rgba(198, 139, 70, 0.02));
  opacity: 0.75;
}

:global(html.dark) .value-item p,
.sst-home.is-dark .value-item p,
:global(html.dark) .capability-copy,
.sst-home.is-dark .capability-copy {
  color: #d0baa0;
}

:global(html.dark) .value-dot,
.sst-home.is-dark .value-dot,
:global(html.dark) .capability-dot,
.sst-home.is-dark .capability-dot,
:global(html.dark) .provider-item::before,
.sst-home.is-dark .provider-item::before {
  background: #c27e4a;
  color: #c27e4a;
}

@keyframes heroLogoSettle {
  0% {
    transform: translate3d(0.8rem, 0.7rem, 0) scale(0.985);
  }

  100% {
    transform: translate3d(0, 0, 0) scale(1);
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
    white-space: normal;
    text-wrap: pretty;
  }

  .home-provider-panel,
  .providers-meta {
    align-items: flex-start;
    flex-direction: column;
  }

  .hero-logo-stage {
    display: none;
  }

  .hero-glyph-notes {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-logo-stage,
  .hero-seal-imprint,
  .hero-glyph-note,
  .hero-glyph-notes {
    animation: none;
  }

  .hero-seal-imprint,
  .hero-glyph-note {
    opacity: 0.86;
  }
}
</style>
