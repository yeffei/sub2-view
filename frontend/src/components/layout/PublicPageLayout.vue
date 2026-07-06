<template>
  <div
    class="sst-public min-h-screen overflow-hidden bg-zen-paper dark:bg-zen-night"
    :class="[toneClass, { 'is-dark': isDark }]"
    :style="backgroundStyle"
  >
    <div class="sst-public-bg" aria-hidden="true"></div>
    <div class="sst-public-paper" aria-hidden="true"></div>
    <div class="sst-public-fibers" aria-hidden="true"></div>
    <div class="sst-public-wash" aria-hidden="true"></div>

    <PublicSiteHeader :nav-items="publicNavItems" />

    <main class="relative z-10 px-5 pb-10 sm:px-8 sm:pb-12">
      <section class="mx-auto grid max-w-7xl gap-12 py-12 lg:grid-cols-[minmax(0,1.03fr)_minmax(21rem,0.97fr)] lg:gap-16 lg:py-18">
        <div class="public-copy-block max-w-3xl">
          <div class="mb-7 flex items-center gap-4">
            <span class="h-px w-14 bg-gradient-to-r from-transparent via-zen-paperLine/90 to-zen-paperLine/30 dark:via-zen-nightLine dark:to-transparent"></span>
            <span class="text-xs uppercase tracking-[0.42em] text-zen-mist dark:text-zen-stone">
              {{ eyebrow }}
            </span>
          </div>

          <h1 class="public-display font-serif text-[clamp(2.9rem,7.8vw,5.8rem)] font-semibold leading-[0.98] text-zen-ink dark:text-zen-paper">
            {{ title }}
          </h1>
          <p v-if="intro" class="public-intro mt-7 max-w-2xl font-serif text-[1.9rem] leading-[1.22] text-zen-inkSoft dark:text-zen-paper sm:text-[2.15rem]">
            {{ intro }}
          </p>
          <p v-if="description" class="mt-6 max-w-[34rem] text-[1.02rem] leading-8 text-zen-mist dark:text-zen-stone sm:text-[1.08rem]">
            {{ description }}
          </p>

          <div v-if="highlights.length" class="mt-9 flex flex-wrap gap-3">
            <span
              v-for="item in highlights"
              :key="item"
              class="inline-flex items-center rounded-full border border-zen-paperLine/80 bg-white/58 px-4 py-2 text-sm text-zen-inkSoft shadow-paper-sm backdrop-blur-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/72 dark:text-zen-paper"
            >
              {{ item }}
            </span>
          </div>
        </div>

        <aside class="public-hero-panel card-glass relative overflow-hidden px-6 py-7 sm:px-8 sm:py-9 lg:sticky lg:top-8">
          <div class="public-hero-mark" aria-hidden="true"></div>
          <div class="public-hero-seal" aria-hidden="true"></div>
          <div class="public-hero-axis" aria-hidden="true"></div>
          <div class="relative z-10">
            <slot name="aside" />
          </div>
        </aside>
      </section>

      <section class="mx-auto max-w-7xl">
        <slot />
      </section>

      <section v-if="showCta" class="public-cta-section mx-auto mt-8 max-w-7xl sm:mt-10">
        <div class="public-cta card-glass relative overflow-hidden rounded-[1.6rem] px-6 py-7 sm:px-8 sm:py-9">
          <div class="public-cta-mark" aria-hidden="true"></div>
          <div class="relative z-10 flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
            <div class="max-w-2xl">
              <div class="text-xs uppercase tracking-[0.32em] text-zen-mist dark:text-zen-stone">SST</div>
              <h2 class="mt-3 font-serif text-3xl leading-tight text-zen-ink dark:text-zen-paper">若已看清路径，就直接入庭。</h2>
              <p class="mt-3 text-sm leading-7 text-zen-mist dark:text-zen-stone sm:text-base">
                Key、用量、订单与账册都在同一入口。
              </p>
            </div>
            <div class="flex flex-col gap-3 sm:flex-row">
              <RouterLink
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex items-center justify-center rounded-zen bg-zen-ink px-6 py-3 text-sm font-medium text-zen-paper shadow-seal transition hover:bg-zen-seal dark:bg-zen-paper dark:text-zen-night dark:hover:bg-zen-seal dark:hover:text-white"
              >
                {{ isAuthenticated ? authenticatedActionLabel : '前往开通' }}
              </RouterLink>
              <RouterLink
                :to="secondaryAction.to"
                class="inline-flex items-center justify-center rounded-zen border border-zen-paperLine bg-white/48 px-6 py-3 text-sm font-medium text-zen-ink transition hover:border-zen-stone hover:bg-white/72 dark:border-zen-nightLine dark:bg-zen-nightPanel/68 dark:text-zen-paper"
              >
                {{ secondaryAction.label }}
              </RouterLink>
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
import { RouterLink, useRoute } from 'vue-router'
import PublicSiteFooter from '@/components/layout/PublicSiteFooter.vue'
import PublicSiteHeader from '@/components/layout/PublicSiteHeader.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { useThemeState } from '@/utils/theme'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

const props = withDefaults(defineProps<{
  eyebrow?: string
  title: string
  intro: string
  description?: string
  highlights?: string[]
  tone?: 'default' | 'pricing' | 'faq' | 'docs' | 'legal'
  authenticatedActionLabel?: string
  showCta?: boolean
}>(), {
  eyebrow: 'SST',
  description: '',
  highlights: () => [],
  tone: 'default',
  authenticatedActionLabel: '进入控制台',
  showCta: true,
})

const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()
const isDark = useThemeState()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const toneClass = computed(() => `sst-public-tone-${props.tone}`)
const backgroundStyle = computed(() => ({
  '--sst-public-bg': `url(${paperInkBg})`,
}))
const secondaryAction = computed(() => {
  if (route.path === '/pricing') {
    return { to: '/home', label: '返回首页' }
  }

  return { to: '/pricing', label: '查看价目' }
})

const publicNavItems = [
  { to: '/pricing', label: '价目' },
  { to: '/docs', label: '文档' },
]

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.sst-public {
  position: relative;
  display: flex;
  flex-direction: column;
  background-color: var(--sst-public-page-bg);
  color: var(--sst-public-page-fg);
  --sst-public-page-bg: #f4efe4;
  --sst-public-page-fg: #1f2320;
  --sst-public-panel-bg: rgba(255, 252, 246, 0.64);
  --sst-public-panel-border: rgba(216, 205, 185, 0.72);
  --sst-public-panel-shadow: 0 24px 64px rgba(58, 48, 32, 0.06);
  --sst-public-chip-bg: rgba(255, 255, 255, 0.58);
  --sst-public-chip-border: rgba(216, 205, 185, 0.8);
  --sst-public-chip-text: #4f5750;
  --sst-public-cta-bg: linear-gradient(140deg, rgba(255, 252, 247, 0.94) 0%, rgba(247, 239, 228, 0.96) 58%, rgba(241, 232, 220, 0.92) 100%);
  --sst-public-cta-border: rgba(176, 150, 118, 0.54);
  --sst-public-cta-shadow: 0 18px 46px rgba(66, 52, 34, 0.1);
  --sst-public-cta-primary-bg: linear-gradient(135deg, rgba(33, 37, 31, 0.96), rgba(51, 38, 29, 0.94));
  --sst-public-cta-primary-text: #f7eddf;
  --sst-public-cta-secondary-bg: rgba(255, 255, 255, 0.42);
  --sst-public-cta-secondary-border: rgba(176, 150, 118, 0.6);
  --sst-public-cta-secondary-text: #2b241d;
  --sst-public-wash-a: rgba(167, 58, 42, 0.08);
  --sst-public-wash-b: rgba(126, 112, 87, 0.08);
  --sst-public-wash-c: rgba(255, 255, 255, 0.16);
  --sst-public-hero-band: rgba(167, 58, 42, 0.07);
  --sst-public-axis-strong: rgba(216, 205, 185, 0.72);
  --sst-public-axis-soft: rgba(216, 205, 185, 0.06);
  --sst-public-seal-border: rgba(167, 58, 42, 0.12);
  --sst-public-seal-core: rgba(167, 58, 42, 0.09);
  --sst-public-cta-radial: rgba(167, 58, 42, 0.1);
  --sst-public-dark-glow-a: rgba(229, 218, 202, 0.12);
  --sst-public-dark-glow-b: rgba(167, 92, 48, 0.14);
}

.sst-public > main {
  flex: 0 0 auto;
}

.sst-public > :deep(footer) {
  flex: 0 0 auto;
  margin-top: auto;
}

:global(html.dark) .sst-public,
.sst-public.is-dark {
  background:
    radial-gradient(circle at 18% 8%, rgba(83, 59, 33, 0.16), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(106, 54, 28, 0.12), transparent 30%),
    #0c0d0b;
  color: #f2e7d7;
  --sst-public-panel-bg:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px);
  --sst-public-panel-border: rgba(155, 126, 86, 0.26);
  --sst-public-panel-shadow: 0 22px 48px rgba(0, 0, 0, 0.24), inset 0 1px 0 rgba(245, 225, 194, 0.055);
  --sst-public-chip-bg: rgba(24, 27, 22, 0.74);
  --sst-public-chip-border: rgba(155, 126, 86, 0.24);
  --sst-public-chip-text: #d4c4ad;
  --sst-public-cta-bg:
    linear-gradient(135deg, rgba(24, 27, 22, 0.84), rgba(35, 29, 23, 0.66)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px),
    rgba(13, 15, 13, 0.5);
  --sst-public-cta-border: rgba(155, 126, 86, 0.28);
  --sst-public-cta-shadow: 0 24px 54px rgba(0, 0, 0, 0.26), inset 0 1px 0 rgba(245, 225, 194, 0.055);
  --sst-public-cta-primary-bg: linear-gradient(135deg, rgba(240, 228, 208, 0.96), rgba(214, 193, 165, 0.92));
  --sst-public-cta-primary-text: #241c16;
  --sst-public-cta-secondary-bg: rgba(24, 27, 22, 0.78);
  --sst-public-cta-secondary-border: rgba(123, 104, 81, 0.82);
  --sst-public-cta-secondary-text: #e7dbc8;
}

.sst-public-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(180deg, rgba(244, 239, 228, 0.9) 0%, rgba(244, 239, 228, 0.52) 38%, rgba(244, 239, 228, 0.16) 100%),
    linear-gradient(90deg, rgba(244, 239, 228, 0.94) 0%, rgba(244, 239, 228, 0.3) 55%, rgba(244, 239, 228, 0.08) 100%),
    var(--sst-public-bg);
  background-size: cover, cover, cover;
  background-position: center, center, center bottom;
}

:global(html.dark) .sst-public-bg,
.sst-public.is-dark .sst-public-bg {
  opacity: 1;
  filter: grayscale(0.08) brightness(0.58) sepia(0.12) saturate(0.95);
  background-image:
    radial-gradient(ellipse at 18% 12%, rgba(222, 180, 116, 0.1), transparent 34%),
    radial-gradient(ellipse at 82% 12%, rgba(164, 82, 42, 0.13), transparent 30%),
    linear-gradient(90deg, rgba(8, 10, 9, 0.96) 0%, rgba(10, 12, 10, 0.72) 42%, rgba(10, 12, 10, 0.34) 100%),
    linear-gradient(180deg, rgba(13, 15, 12, 0.22) 0%, rgba(13, 15, 12, 0.32) 52%, rgba(8, 9, 8, 0.7) 100%),
    var(--sst-public-bg);
  background-size: cover, cover, cover, cover, cover;
  background-position: center, center, center, center, center bottom;
}

.sst-public-paper {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.14;
  background-image:
    linear-gradient(rgba(31, 35, 32, 0.022) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.016) 1px, transparent 1px);
  background-size: 128px 128px, 128px 128px;
}

:global(html.dark) .sst-public-paper,
.sst-public.is-dark .sst-public-paper {
  opacity: 0.11;
  background-image:
    linear-gradient(rgba(239, 217, 179, 0.024) 1px, transparent 1px),
    linear-gradient(90deg, rgba(239, 217, 179, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, transparent 0 49.6%, rgba(217, 165, 94, 0.06) 49.8% 50%, transparent 50.2% 100%);
  background-size: 132px 132px, 132px 132px, min(78rem, 92vw) 100%;
  background-position: center top, center top, center top;
}

.sst-public-paper::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.18;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(31, 35, 32, 0.06) 0 1px, transparent 1.5px),
    radial-gradient(circle at 72% 42%, rgba(31, 35, 32, 0.045) 0 1px, transparent 1.5px);
  background-size: 34px 41px, 48px 57px;
}

:global(html.dark) .sst-public-paper::after,
.sst-public.is-dark .sst-public-paper::after {
  opacity: 0.13;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(241, 217, 177, 0.06) 0 0.8px, transparent 1.4px),
    radial-gradient(circle at 72% 42%, rgba(241, 217, 177, 0.04) 0 0.9px, transparent 1.5px),
    linear-gradient(90deg, rgba(255, 226, 176, 0.045), transparent 5%, transparent 95%, rgba(255, 226, 176, 0.036)),
    linear-gradient(180deg, rgba(255, 236, 202, 0.05), transparent 16%, transparent 76%, rgba(0, 0, 0, 0.26));
  background-size: 34px 41px, 48px 57px, 100% 100%, 100% 100%;
}

.sst-public-fibers {
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

:global(html.dark) .sst-public-fibers,
.sst-public.is-dark .sst-public-fibers {
  opacity: 0.09;
  mix-blend-mode: screen;
  background-image:
    repeating-linear-gradient(7deg, transparent 0 12px, rgba(236, 203, 154, 0.04) 12px 13px, transparent 13px 31px),
    repeating-linear-gradient(96deg, transparent 0 18px, rgba(255, 238, 210, 0.05) 18px 19px, transparent 19px 44px);
}

.sst-public-wash {
  position: absolute;
  inset: auto 0 0 0;
  height: 24rem;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 40%, var(--sst-public-wash-a), transparent 34%),
    radial-gradient(circle at 78% 12%, var(--sst-public-wash-b), transparent 28%),
    linear-gradient(180deg, transparent, var(--sst-public-wash-c));
}

:global(html.dark) .sst-public-wash,
.sst-public.is-dark .sst-public-wash {
  height: 26rem;
  background:
    radial-gradient(circle at 14% 18%, rgba(222, 180, 116, 0.055), transparent 30%),
    radial-gradient(circle at 82% 12%, rgba(164, 82, 42, 0.12), transparent 24%),
    linear-gradient(180deg, transparent 18%, rgba(12, 13, 11, 0.42) 100%);
}

.sst-public-tone-pricing {
  --sst-public-wash-a: rgba(167, 58, 42, 0.11);
  --sst-public-wash-b: rgba(146, 118, 79, 0.1);
  --sst-public-hero-band: rgba(167, 58, 42, 0.1);
  --sst-public-seal-border: rgba(167, 58, 42, 0.16);
  --sst-public-seal-core: rgba(167, 58, 42, 0.11);
  --sst-public-cta-radial: rgba(167, 58, 42, 0.14);
}

.sst-public-tone-faq {
  --sst-public-wash-a: rgba(98, 94, 83, 0.06);
  --sst-public-wash-b: rgba(133, 124, 104, 0.06);
  --sst-public-hero-band: rgba(104, 97, 83, 0.06);
  --sst-public-axis-strong: rgba(188, 177, 158, 0.62);
  --sst-public-axis-soft: rgba(188, 177, 158, 0.08);
  --sst-public-seal-border: rgba(121, 109, 95, 0.12);
  --sst-public-seal-core: rgba(121, 109, 95, 0.08);
  --sst-public-cta-radial: rgba(121, 109, 95, 0.09);
}

.sst-public-tone-docs {
  --sst-public-wash-a: rgba(86, 117, 117, 0.08);
  --sst-public-wash-b: rgba(126, 112, 87, 0.06);
  --sst-public-hero-band: rgba(77, 126, 126, 0.08);
  --sst-public-axis-strong: rgba(180, 198, 192, 0.54);
  --sst-public-axis-soft: rgba(180, 198, 192, 0.08);
  --sst-public-seal-border: rgba(77, 126, 126, 0.14);
  --sst-public-seal-core: rgba(77, 126, 126, 0.1);
  --sst-public-cta-radial: rgba(77, 126, 126, 0.12);
}

.sst-public-tone-legal {
  --sst-public-wash-a: rgba(160, 103, 49, 0.1);
  --sst-public-wash-b: rgba(132, 104, 67, 0.08);
  --sst-public-wash-c: rgba(255, 244, 224, 0.14);
  --sst-public-hero-band: rgba(176, 120, 57, 0.08);
  --sst-public-axis-strong: rgba(194, 165, 117, 0.6);
  --sst-public-axis-soft: rgba(194, 165, 117, 0.08);
  --sst-public-seal-border: rgba(176, 120, 57, 0.16);
  --sst-public-seal-core: rgba(176, 120, 57, 0.12);
  --sst-public-cta-radial: rgba(176, 120, 57, 0.14);
  --sst-public-dark-glow-a: rgba(236, 223, 197, 0.14);
  --sst-public-dark-glow-b: rgba(173, 111, 51, 0.18);
}

:global(html.dark) .sst-public.sst-public-tone-legal {
  background:
    radial-gradient(circle at 18% 8%, rgba(83, 59, 33, 0.18), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(126, 70, 34, 0.13), transparent 30%),
    #0c0d0b;
}

:global(html.dark) .sst-public-tone-legal .sst-public-bg {
  filter: grayscale(0.08) brightness(0.58) sepia(0.12) saturate(0.95);
  background-image:
    radial-gradient(ellipse at 18% 12%, rgba(230, 185, 118, 0.11), transparent 34%),
    radial-gradient(ellipse at 82% 12%, rgba(184, 94, 42, 0.14), transparent 30%),
    linear-gradient(90deg, rgba(8, 10, 9, 0.96) 0%, rgba(10, 12, 10, 0.72) 42%, rgba(10, 12, 10, 0.34) 100%),
    linear-gradient(180deg, rgba(13, 15, 12, 0.22) 0%, rgba(13, 15, 12, 0.32) 52%, rgba(8, 9, 8, 0.72) 100%),
    var(--sst-public-bg);
}

:global(html.dark) .sst-public-tone-legal .sst-public-paper {
  opacity: 0.11;
}

:global(html.dark) .sst-public-tone-legal .sst-public-paper::after {
  opacity: 0.13;
}

:global(html.dark) .sst-public-tone-legal .sst-public-wash {
  background:
    radial-gradient(circle at 16% 18%, rgba(230, 185, 118, 0.06), transparent 26%),
    radial-gradient(circle at 82% 10%, rgba(193, 129, 62, 0.13), transparent 22%),
    linear-gradient(180deg, transparent 16%, rgba(12, 13, 11, 0.42) 100%);
}

:global(html.dark) .sst-public-tone-legal .public-display {
  color: #f6ead7;
  text-shadow: 0 1px 0 rgba(255, 240, 218, 0.05);
}

:global(html.dark) .sst-public-tone-legal .public-intro {
  color: #d6c4a7;
}

:global(html.dark) .sst-public-tone-legal .public-copy-block > div:first-child span:last-child,
:global(html.dark) .sst-public-tone-legal .public-cta > .relative > div:first-child > div:first-child {
  color: #cfb184;
}

.public-hero-panel {
  min-height: 100%;
  background: var(--sst-public-panel-bg);
  border: 1px solid var(--sst-public-panel-border);
  box-shadow: var(--sst-public-panel-shadow);
}

:global(html.dark) .public-hero-panel {
  backdrop-filter: blur(18px);
  box-shadow: var(--sst-public-panel-shadow);
}

.public-hero-mark {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.18), transparent 36%),
    linear-gradient(90deg, transparent 0, var(--sst-public-hero-band) 48%, transparent 100%);
  opacity: 0.95;
}

.public-hero-axis {
  position: absolute;
  inset: 1.1rem auto 1.1rem 1.1rem;
  width: 1px;
  background: linear-gradient(180deg, var(--sst-public-axis-strong), var(--sst-public-axis-soft));
  opacity: 0.75;
}

.public-copy-block,
.public-hero-panel,
.public-cta {
  animation: publicRise 560ms ease-out both;
}

.public-hero-panel {
  animation-delay: 70ms;
}

.public-cta {
  animation-delay: 120ms;
}

.public-hero-seal {
  position: absolute;
  right: 1.4rem;
  top: 1.2rem;
  width: 4.5rem;
  height: 4.5rem;
  border-radius: 999px;
  border: 1px solid var(--sst-public-seal-border);
  background:
    radial-gradient(circle at center, var(--sst-public-seal-core), transparent 62%),
    radial-gradient(circle at center, rgba(255, 255, 255, 0.22), transparent 72%);
  opacity: 0.9;
  filter: blur(0.2px);
}

.public-cta-mark {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.16), transparent 34%),
    radial-gradient(circle at 82% 18%, var(--sst-public-cta-radial), transparent 24%);
}

.public-cta {
  border: 1px solid var(--sst-public-cta-border);
  background: var(--sst-public-cta-bg);
  box-shadow: var(--sst-public-cta-shadow);
}

:global(html.dark) .public-cta {
  backdrop-filter: blur(18px);
}

.sst-public :deep(.card-glass) {
  background: var(--sst-public-panel-bg);
  border-color: var(--sst-public-panel-border);
  box-shadow: var(--sst-public-panel-shadow);
}

.sst-public :deep(.shadow-paper-sm) {
  box-shadow: 0 14px 34px rgba(54, 43, 29, 0.08);
}

:global(html.dark) .sst-public :deep(.shadow-paper-sm) {
  box-shadow: 0 18px 38px rgba(0, 0, 0, 0.18);
}

.sst-public :deep(.public-copy-block .inline-flex.rounded-full) {
  background: var(--sst-public-chip-bg);
  border-color: var(--sst-public-chip-border);
  color: var(--sst-public-chip-text);
}

.sst-public :deep(.public-cta a:first-child) {
  background: var(--sst-public-cta-primary-bg);
  color: var(--sst-public-cta-primary-text);
  border: 1px solid rgba(244, 232, 214, 0.14);
  box-shadow: 0 18px 36px rgba(84, 57, 31, 0.18);
}

.sst-public :deep(.public-cta a:first-child:hover) {
  filter: brightness(1.03);
}

.sst-public :deep(.public-cta a:last-child) {
  background: var(--sst-public-cta-secondary-bg);
  border-color: var(--sst-public-cta-secondary-border);
  color: var(--sst-public-cta-secondary-text);
}

.sst-public :deep(.public-cta a:last-child:hover) {
  border-color: rgba(188, 93, 31, 0.34);
  color: #f3e8d5;
}

.public-display {
  text-wrap: balance;
}

.public-intro {
  text-wrap: balance;
}

@keyframes publicRise {
  from {
    opacity: 0.82;
    transform: translateY(14px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (min-width: 1024px) {
  .public-hero-panel::before {
    content: '';
    position: absolute;
    inset: 1.1rem;
    border: 1px solid rgba(216, 205, 185, 0.22);
    border-radius: 1.2rem;
    pointer-events: none;
  }

  :global(html.dark) .public-hero-panel::before {
    border-color: rgba(99, 104, 93, 0.22);
  }
}

@media (max-width: 768px) {
  .sst-public-bg {
    background-image:
      linear-gradient(180deg, rgba(244, 239, 228, 0.9) 0%, rgba(244, 239, 228, 0.34) 54%, rgba(244, 239, 228, 0.14) 100%),
      var(--sst-public-bg);
    background-position: center top, 58% bottom;
    background-size: cover, auto 92%;
  }

  :global(html.dark) .sst-public-bg {
    opacity: 1;
    filter: grayscale(0.14) brightness(0.96) contrast(0.98) saturate(0.94);
    background-image:
      linear-gradient(180deg, rgba(22, 25, 21, 0.54) 0%, rgba(22, 25, 21, 0.3) 48%, rgba(22, 25, 21, 0.6) 100%),
      radial-gradient(circle at 18% 22%, rgba(229, 218, 202, 0.14), transparent 34%),
      radial-gradient(circle at 82% 18%, rgba(167, 92, 48, 0.16), transparent 26%),
      var(--sst-public-bg);
    background-position: center top, 58% bottom, 42% top, 58% bottom;
    background-size: cover, cover, cover, auto 92%;
  }

  :global(html.dark) .sst-public-tone-legal .sst-public-bg {
    filter: grayscale(0.1) brightness(0.92) contrast(1) saturate(0.96);
    background-image:
      linear-gradient(180deg, rgba(19, 17, 15, 0.6) 0%, rgba(19, 17, 15, 0.34) 48%, rgba(19, 17, 15, 0.64) 100%),
      radial-gradient(circle at 18% 22%, rgba(244, 226, 190, 0.14), transparent 34%),
      radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.2), transparent 26%),
      var(--sst-public-bg);
  }

  .public-hero-seal {
    width: 3.8rem;
    height: 3.8rem;
  }

  .public-hero-axis {
    display: none;
  }

  .public-cta {
    border-radius: 1.35rem;
  }

}
</style>

<style>
html.dark .sst-public {
  background:
    radial-gradient(circle at 18% 8%, rgba(83, 59, 33, 0.16), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(106, 54, 28, 0.12), transparent 30%),
    #0c0d0b !important;
  color: #f3eadb;
}

html.dark .sst-public.sst-public-tone-legal {
  background:
    radial-gradient(circle at 18% 8%, rgba(83, 59, 33, 0.18), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(126, 70, 34, 0.13), transparent 30%),
    #0c0d0b !important;
}

html.dark .sst-public .sst-public-bg {
  opacity: 1;
  filter: grayscale(0.08) brightness(0.58) sepia(0.12) saturate(0.95);
  background-image:
    radial-gradient(ellipse at 18% 12%, rgba(222, 180, 116, 0.1), transparent 34%),
    radial-gradient(ellipse at 82% 12%, rgba(164, 82, 42, 0.13), transparent 30%),
    linear-gradient(90deg, rgba(8, 10, 9, 0.96) 0%, rgba(10, 12, 10, 0.72) 42%, rgba(10, 12, 10, 0.34) 100%),
    linear-gradient(180deg, rgba(13, 15, 12, 0.22) 0%, rgba(13, 15, 12, 0.32) 52%, rgba(8, 9, 8, 0.7) 100%),
    var(--sst-public-bg);
  background-size: cover, cover, cover, cover, cover;
  background-position: center, center, center, center, center bottom;
}

html.dark .sst-public.sst-public-tone-legal .sst-public-bg {
  filter: grayscale(0.08) brightness(0.58) sepia(0.12) saturate(0.95);
  background-image:
    radial-gradient(ellipse at 18% 12%, rgba(230, 185, 118, 0.11), transparent 34%),
    radial-gradient(ellipse at 82% 12%, rgba(184, 94, 42, 0.14), transparent 30%),
    linear-gradient(90deg, rgba(8, 10, 9, 0.96) 0%, rgba(10, 12, 10, 0.72) 42%, rgba(10, 12, 10, 0.34) 100%),
    linear-gradient(180deg, rgba(13, 15, 12, 0.22) 0%, rgba(13, 15, 12, 0.32) 52%, rgba(8, 9, 8, 0.72) 100%),
    var(--sst-public-bg);
}

html.dark .sst-public.sst-public-tone-legal .public-display {
  color: #f6ead7;
  text-shadow: 0 1px 0 rgba(255, 240, 218, 0.05);
}

html.dark .sst-public.sst-public-tone-legal .public-intro {
  color: #d6c4a7;
}

html.dark .sst-public.sst-public-tone-legal .public-copy-block > div:first-child span:last-child,
html.dark .sst-public.sst-public-tone-legal .public-cta > .relative > div:first-child > div:first-child {
  color: #cfb184;
}

@media (max-width: 768px) {
  html.dark .sst-public .sst-public-bg {
    filter: grayscale(0.08) brightness(0.6) sepia(0.12) saturate(0.95);
    background-image:
      radial-gradient(ellipse at 18% 12%, rgba(222, 180, 116, 0.1), transparent 34%),
      radial-gradient(ellipse at 82% 12%, rgba(164, 82, 42, 0.13), transparent 30%),
      linear-gradient(180deg, rgba(12, 14, 12, 0.82) 0%, rgba(13, 15, 12, 0.52) 48%, rgba(8, 9, 8, 0.74) 100%),
      var(--sst-public-bg);
    background-position: center top, center top, center top, 58% bottom;
    background-size: cover, cover, cover, auto 92%;
  }

  html.dark .sst-public.sst-public-tone-legal .sst-public-bg {
    filter: grayscale(0.08) brightness(0.6) sepia(0.12) saturate(0.95);
    background-image:
      radial-gradient(ellipse at 18% 12%, rgba(230, 185, 118, 0.11), transparent 34%),
      radial-gradient(ellipse at 82% 12%, rgba(184, 94, 42, 0.14), transparent 30%),
      linear-gradient(180deg, rgba(12, 14, 12, 0.82) 0%, rgba(13, 15, 12, 0.52) 48%, rgba(8, 9, 8, 0.76) 100%),
      var(--sst-public-bg);
  }
}

html:not(.dark) .sst-public {
  background:
    radial-gradient(circle at 18% 8%, rgba(255, 253, 247, 0.48), transparent 34%),
    radial-gradient(circle at 86% 18%, rgba(197, 151, 83, 0.1), transparent 30%),
    #f4efe4 !important;
  color: #1f2320;
  --sst-public-panel-bg:
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    linear-gradient(90deg, rgba(144, 113, 76, 0.038), transparent 18%, rgba(144, 113, 76, 0.024) 82%, transparent),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28);
  --sst-public-panel-border: rgba(154, 128, 92, 0.16);
  --sst-public-panel-shadow:
    0 14px 34px rgba(84, 57, 31, 0.05),
    inset 0 1px 0 rgba(255, 249, 239, 0.6),
    inset 0 -1px 0 rgba(140, 111, 76, 0.07),
    inset 0 0 0 1px rgba(255, 255, 255, 0.22);
  --sst-public-chip-bg: rgba(255, 252, 246, 0.62);
  --sst-public-chip-border: rgba(190, 168, 134, 0.42);
  --sst-public-chip-text: #4f5750;
  --sst-public-cta-bg:
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28);
  --sst-public-cta-border: rgba(176, 150, 118, 0.44);
  --sst-public-cta-shadow:
    0 14px 34px rgba(84, 57, 31, 0.06),
    inset 0 1px 0 rgba(255, 249, 239, 0.62);
  --sst-public-cta-primary-bg: linear-gradient(135deg, #28231c, #403227);
  --sst-public-cta-primary-text: #f7f0e4;
  --sst-public-cta-secondary-bg: rgba(255, 251, 245, 0.72);
  --sst-public-cta-secondary-border: rgba(190, 157, 118, 0.54);
  --sst-public-cta-secondary-text: #6f5330;
}

html:not(.dark) .sst-public .sst-public-bg {
  filter: none;
  background-image:
    radial-gradient(ellipse at 18% 10%, rgba(255, 253, 247, 0.46), transparent 36%),
    radial-gradient(ellipse at 82% 20%, rgba(226, 209, 176, 0.22), transparent 34%),
    linear-gradient(90deg, rgba(244, 239, 228, 0.94) 0%, rgba(244, 239, 228, 0.58) 42%, rgba(244, 239, 228, 0.14) 100%),
    linear-gradient(180deg, rgba(251, 246, 236, 0.2) 0%, rgba(244, 239, 228, 0.18) 55%, rgba(231, 219, 200, 0.44) 100%),
    var(--sst-public-bg);
  background-size: cover, cover, cover, cover, cover;
  background-position: center, center, center, center, center bottom;
}

html:not(.dark) .sst-public .sst-public-paper {
  opacity: 0.22;
  background-image:
    linear-gradient(rgba(90, 77, 57, 0.026) 1px, transparent 1px),
    linear-gradient(90deg, rgba(90, 77, 57, 0.018) 1px, transparent 1px),
    linear-gradient(90deg, transparent 0 49.6%, rgba(144, 113, 76, 0.04) 49.8% 50%, transparent 50.2% 100%);
  background-size: 132px 132px, 132px 132px, min(78rem, 92vw) 100%;
  background-position: center top, center top, center top;
}

html:not(.dark) .sst-public .sst-public-paper::after {
  opacity: 0.3;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(92, 74, 48, 0.07) 0 0.8px, transparent 1.4px),
    radial-gradient(circle at 72% 42%, rgba(92, 74, 48, 0.05) 0 0.9px, transparent 1.5px),
    linear-gradient(90deg, rgba(87, 67, 42, 0.055), transparent 5%, transparent 95%, rgba(87, 67, 42, 0.045)),
    linear-gradient(180deg, rgba(255, 254, 250, 0.24), transparent 15%, transparent 76%, rgba(126, 99, 63, 0.08));
  background-size: 34px 41px, 48px 57px, 100% 100%, 100% 100%;
}

html:not(.dark) .sst-public .sst-public-fibers {
  opacity: 0.18;
  mix-blend-mode: multiply;
}

html:not(.dark) .sst-public .public-display {
  color: #2f281f;
  text-shadow: 0 1px 0 rgba(255, 248, 238, 0.55);
}

html:not(.dark) .sst-public .public-intro {
  color: #4a3f32;
}

html:not(.dark) .sst-public .public-copy-block > div:first-child span:last-child,
html:not(.dark) .sst-public .public-cta > .relative > div:first-child > div:first-child {
  color: #9b7a52;
}

@media (max-width: 768px) {
  html:not(.dark) .sst-public .sst-public-bg {
    background-image:
      radial-gradient(ellipse at 18% 10%, rgba(255, 253, 247, 0.44), transparent 36%),
      linear-gradient(180deg, rgba(244, 239, 228, 0.88) 0%, rgba(244, 239, 228, 0.36) 54%, rgba(244, 239, 228, 0.18) 100%),
      var(--sst-public-bg);
    background-position: center top, center top, 58% bottom;
    background-size: cover, cover, auto 92%;
  }
}
</style>



