<template>
  <div
    class="sst-public min-h-screen overflow-hidden bg-zen-paper dark:bg-zen-night"
    :class="[toneClass, { 'is-dark': isDark }]"
    :style="backgroundStyle"
  >
    <div class="sst-public-bg" aria-hidden="true"></div>
    <div class="sst-public-paper" aria-hidden="true"></div>
    <div class="sst-public-wash" aria-hidden="true"></div>

    <PublicSiteHeader :nav-items="publicNavItems" />

    <main class="relative z-10 px-5 pb-20 sm:px-8">
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

      <section class="mx-auto mt-10 max-w-7xl sm:mt-12">
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
}>(), {
  eyebrow: 'SST',
  description: '',
  highlights: () => [],
  tone: 'default',
  authenticatedActionLabel: '进入控制台',
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

:global(html.dark) .sst-public,
.sst-public.is-dark {
  background: #111310;
  color: #f2e7d7;
  --sst-public-panel-bg: linear-gradient(180deg, rgba(22, 24, 20, 0.9), rgba(30, 26, 22, 0.84));
  --sst-public-panel-border: rgba(122, 104, 81, 0.56);
  --sst-public-panel-shadow: 0 26px 68px rgba(0, 0, 0, 0.24);
  --sst-public-chip-bg: rgba(26, 29, 24, 0.78);
  --sst-public-chip-border: rgba(111, 95, 74, 0.62);
  --sst-public-chip-text: #dfd1be;
  --sst-public-cta-bg: linear-gradient(140deg, rgba(24, 26, 22, 0.94) 0%, rgba(34, 30, 25, 0.96) 58%, rgba(45, 34, 27, 0.94) 100%);
  --sst-public-cta-border: rgba(123, 104, 81, 0.62);
  --sst-public-cta-shadow: 0 28px 72px rgba(0, 0, 0, 0.28);
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
  filter: grayscale(0.12) brightness(0.88) contrast(1) saturate(0.9);
  background-image:
    linear-gradient(180deg, rgba(16, 18, 15, 0.64) 0%, rgba(16, 18, 15, 0.42) 36%, rgba(16, 18, 15, 0.7) 100%),
    linear-gradient(90deg, rgba(15, 18, 15, 0.5) 0%, rgba(15, 18, 15, 0.16) 54%, rgba(15, 18, 15, 0.48) 100%),
    radial-gradient(circle at 18% 22%, rgba(226, 214, 195, 0.13), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(171, 101, 52, 0.18), transparent 26%),
    var(--sst-public-bg);
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
  opacity: 0.06;
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
  opacity: 0.06;
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
    radial-gradient(circle at 14% 18%, rgba(226, 211, 191, 0.045), transparent 30%),
    radial-gradient(circle at 82% 12%, rgba(171, 101, 52, 0.12), transparent 24%),
    linear-gradient(180deg, transparent 18%, rgba(12, 13, 11, 0.34) 100%);
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
  background: #12100e;
}

:global(html.dark) .sst-public-tone-legal .sst-public-bg {
  filter: grayscale(0.1) brightness(0.86) contrast(1) saturate(0.92);
  background-image:
    linear-gradient(180deg, rgba(18, 16, 14, 0.66) 0%, rgba(18, 16, 14, 0.44) 36%, rgba(18, 16, 14, 0.74) 100%),
    linear-gradient(90deg, rgba(15, 13, 12, 0.56) 0%, rgba(15, 13, 12, 0.18) 54%, rgba(15, 13, 12, 0.56) 100%),
    radial-gradient(circle at 18% 22%, rgba(239, 220, 190, 0.13), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.22), transparent 26%),
    var(--sst-public-bg);
}

:global(html.dark) .sst-public-tone-legal .sst-public-paper {
  opacity: 0.06;
}

:global(html.dark) .sst-public-tone-legal .sst-public-paper::after {
  opacity: 0.06;
}

:global(html.dark) .sst-public-tone-legal .sst-public-wash {
  background:
    radial-gradient(circle at 16% 18%, rgba(236, 219, 189, 0.05), transparent 26%),
    radial-gradient(circle at 82% 10%, rgba(193, 129, 62, 0.14), transparent 22%),
    linear-gradient(180deg, transparent 16%, rgba(12, 11, 10, 0.34) 100%);
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
  background: #161915 !important;
  color: #f3eadb;
}

html.dark .sst-public.sst-public-tone-legal {
  background: #13110f !important;
}

html.dark .sst-public .sst-public-bg {
  opacity: 1;
  filter: grayscale(0.14) brightness(0.96) contrast(0.98) saturate(0.94);
  background-image:
    linear-gradient(180deg, rgba(22, 25, 21, 0.5) 0%, rgba(22, 25, 21, 0.32) 36%, rgba(22, 25, 21, 0.58) 100%),
    linear-gradient(90deg, rgba(19, 22, 18, 0.42) 0%, rgba(19, 22, 18, 0.12) 54%, rgba(19, 22, 18, 0.44) 100%),
    radial-gradient(circle at 18% 22%, rgba(229, 218, 202, 0.15), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(167, 92, 48, 0.17), transparent 26%),
    var(--sst-public-bg);
}

html.dark .sst-public.sst-public-tone-legal .sst-public-bg {
  filter: grayscale(0.1) brightness(0.92) contrast(1) saturate(0.96);
  background-image:
    linear-gradient(180deg, rgba(19, 17, 15, 0.58) 0%, rgba(19, 17, 15, 0.38) 36%, rgba(19, 17, 15, 0.66) 100%),
    linear-gradient(90deg, rgba(15, 13, 12, 0.5) 0%, rgba(15, 13, 12, 0.16) 54%, rgba(15, 13, 12, 0.52) 100%),
    radial-gradient(circle at 18% 22%, rgba(244, 226, 190, 0.13), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.2), transparent 26%),
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
    filter: grayscale(0.14) brightness(0.96) contrast(0.98) saturate(0.94);
    background-image:
      linear-gradient(180deg, rgba(22, 25, 21, 0.54) 0%, rgba(22, 25, 21, 0.3) 48%, rgba(22, 25, 21, 0.6) 100%),
      radial-gradient(circle at 18% 22%, rgba(229, 218, 202, 0.14), transparent 34%),
      radial-gradient(circle at 82% 18%, rgba(167, 92, 48, 0.16), transparent 26%),
      var(--sst-public-bg);
    background-position: center top, 58% bottom, 42% top, 58% bottom;
    background-size: cover, cover, cover, auto 92%;
  }

  html.dark .sst-public.sst-public-tone-legal .sst-public-bg {
    filter: grayscale(0.1) brightness(0.92) contrast(1) saturate(0.96);
    background-image:
      linear-gradient(180deg, rgba(19, 17, 15, 0.6) 0%, rgba(19, 17, 15, 0.34) 48%, rgba(19, 17, 15, 0.64) 100%),
      radial-gradient(circle at 18% 22%, rgba(244, 226, 190, 0.14), transparent 34%),
      radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.2), transparent 26%),
      var(--sst-public-bg);
  }
}
</style>



