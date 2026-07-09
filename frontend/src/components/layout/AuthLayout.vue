<template>
  <div
    class="auth-shell relative min-h-screen overflow-hidden"
    :class="{ 'auth-shell-dark': isDark }"
    :style="authBackgroundStyle"
  >
    <div class="auth-bg" aria-hidden="true"></div>
    <div class="auth-paper" aria-hidden="true"></div>

    <header class="auth-brand-bar absolute inset-x-0 top-0 z-20 px-5 py-5 sm:px-8">
      <div class="mx-auto flex max-w-7xl items-center border-b border-zen-paperLine/60 pb-4 dark:border-zen-nightLine/80">
        <router-link
          to="/home"
          class="auth-brand-link inline-flex items-center gap-3"
          aria-label="Back to home"
        >
          <span class="auth-brand-mark" aria-hidden="true">
            <img :src="siteLogo" alt="" class="h-full w-full object-contain" />
          </span>
          <div class="leading-tight">
            <div class="font-serif text-base font-semibold tracking-[0.18em] text-zen-ink dark:text-zen-paper">
              {{ siteName }}
            </div>
            <div class="text-[10px] uppercase tracking-[0.36em] text-zen-mist dark:text-zen-stone">
              SST
            </div>
          </div>
        </router-link>
      </div>
    </header>

    <div class="relative z-10 mx-auto grid min-h-[calc(100vh-116px)] w-full max-w-[84rem] px-5 pb-12 pt-28 sm:px-8 sm:pt-32 lg:grid-cols-[minmax(0,1fr)_minmax(24rem,29rem)] lg:items-center lg:gap-10 lg:px-10 lg:pt-16 xl:gap-16">
      <section class="auth-hero hidden lg:flex lg:items-center">
        <slot name="hero" :siteName="siteName" :siteSubtitle="siteSubtitle" :isDark="isDark">
          <div class="mb-8 flex items-center gap-4">
            <span class="h-px w-16 bg-zen-paperLine dark:bg-zen-nightLine"></span>
            <span class="auth-hero-kicker text-xs uppercase tracking-[0.42em] text-zen-mist dark:text-zen-stone">SST</span>
          </div>
          <h1 class="auth-hero-title font-serif text-[clamp(4.6rem,7.2vw,7.8rem)] font-semibold leading-none text-zen-ink dark:text-zen-paper">
            {{ siteName }}
          </h1>
          <p class="auth-hero-lead mt-7 max-w-xl font-serif text-4xl leading-tight text-zen-inkSoft dark:text-zen-paper">
            {{ t('publicSite.tagline') }}
          </p>
          <p class="auth-hero-copy mt-5 max-w-md text-sm leading-7 text-zen-mist dark:text-zen-stone">
            {{ siteSubtitle }}
          </p>
          <div class="auth-hero-marks mt-14 flex items-center gap-6 text-sm text-zen-mist dark:text-zen-stone">
            <span>{{ t('authBrand.marks.stable') }}</span>
            <span class="h-px w-16 bg-zen-paperLine dark:bg-zen-nightLine"></span>
            <span>{{ t('authBrand.marks.ledger') }}</span>
            <span class="h-px w-16 bg-zen-paperLine dark:bg-zen-nightLine"></span>
            <span>{{ t('authBrand.marks.access') }}</span>
          </div>
        </slot>
      </section>

      <section class="auth-form-shell mx-auto w-full max-w-[29rem] lg:justify-self-end">
        <div class="auth-card rounded-zen p-7 sm:p-8">
          <div class="auth-card-head mb-7 flex items-start justify-between gap-6">
            <div>
              <div class="auth-card-kicker text-xs uppercase tracking-[0.36em] text-zen-mist dark:text-zen-stone">
                {{ t('authBrand.cardKicker') }}
              </div>
              <div class="auth-card-title mt-2 font-serif text-2xl font-semibold text-zen-ink dark:text-zen-paper">
                {{ t('authBrand.cardTitle') }}
              </div>
            </div>
            <span class="auth-card-mark" aria-hidden="true">
              <img :src="siteLogo" alt="" class="h-full w-full object-contain" />
            </span>
          </div>
          <slot />
        </div>

        <div class="auth-shell-footer mt-6 text-center text-sm">
          <slot name="footer" />
        </div>

        <div class="auth-shell-copyright mt-8 text-center text-xs text-zen-mist dark:text-zen-stone">
          &copy; {{ currentYear }} {{ siteName }} · {{ t('authBrand.copyright') }}
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { useThemeState } from '@/utils/theme'
import { sanitizeUrl } from '@/utils/url'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

const appStore = useAppStore()
const { t, locale } = useI18n()
const isDark = useThemeState()

const rawName = computed(() => appStore.siteName || '山枢庭')
const siteName = computed(() => rawName.value === 'Sub2API' ? '山枢庭' : rawName.value)
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '/logo.png', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => {
  const subtitle = appStore.cachedPublicSettings?.site_subtitle
  const legacySubtitles = new Set([
    'Subscription to API Conversion Platform',
    '统一入口，安静流转。',
  ])

  if (locale.value.startsWith('en') && subtitle === '统一入口，安静流转。') {
    return t('authBrand.defaultSubtitle')
  }

  return subtitle && !legacySubtitles.has(subtitle)
    ? subtitle
    : t('authBrand.defaultSubtitle')
})
const currentYear = computed(() => new Date().getFullYear())
const authBackgroundStyle = computed(() => ({
  '--sst-auth-bg': `url(${paperInkBg})`,
}))

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>

<style scoped>
.auth-shell {
  background: #f4efe4;
  color: #1f2320;
}

.auth-shell-dark {
  background: #13110f;
}

:global(html.dark) .auth-shell {
  background: #13110f;
}

.auth-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image:
    linear-gradient(90deg, rgba(244, 239, 228, 0.93) 0%, rgba(244, 239, 228, 0.62) 48%, rgba(244, 239, 228, 0.26) 100%),
    linear-gradient(180deg, rgba(244, 239, 228, 0.08) 0%, rgba(244, 239, 228, 0.26) 100%),
    var(--sst-auth-bg);
  background-position: center, center, center bottom;
  background-size: cover, cover, cover;
}

.auth-shell-dark .auth-bg {
  background-image:
    linear-gradient(180deg, rgba(19, 17, 15, 0.74) 0%, rgba(19, 17, 15, 0.56) 36%, rgba(19, 17, 15, 0.82) 100%),
    linear-gradient(90deg, rgba(15, 13, 12, 0.68) 0%, rgba(15, 13, 12, 0.26) 54%, rgba(15, 13, 12, 0.72) 100%),
    radial-gradient(circle at 18% 22%, rgba(244, 226, 190, 0.08), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.18), transparent 26%),
    var(--sst-auth-bg);
  background-position: center, center, center, center, center bottom;
  background-size: cover, cover, cover, cover, cover;
  opacity: 1;
  filter: grayscale(0.16) brightness(0.72) contrast(0.98) saturate(0.92);
}

:global(html.dark) .auth-bg {
  background-image:
    linear-gradient(180deg, rgba(19, 17, 15, 0.74) 0%, rgba(19, 17, 15, 0.56) 36%, rgba(19, 17, 15, 0.82) 100%),
    linear-gradient(90deg, rgba(15, 13, 12, 0.68) 0%, rgba(15, 13, 12, 0.26) 54%, rgba(15, 13, 12, 0.72) 100%),
    radial-gradient(circle at 18% 22%, rgba(244, 226, 190, 0.08), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(194, 129, 61, 0.18), transparent 26%),
    var(--sst-auth-bg);
  background-position: center, center, center, center, center bottom;
  background-size: cover, cover, cover, cover, cover;
  opacity: 1;
  filter: grayscale(0.16) brightness(0.72) contrast(0.98) saturate(0.92);
}

.auth-paper {
  position: absolute;
  inset: 0;
  opacity: 0.1;
  background-image:
    linear-gradient(rgba(31, 35, 32, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.024) 1px, transparent 1px);
  background-size: 128px 128px;
}

.auth-paper::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.18;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(31, 35, 32, 0.05) 0 1px, transparent 1.5px),
    radial-gradient(circle at 72% 42%, rgba(31, 35, 32, 0.035) 0 1px, transparent 1.5px);
  background-size: 34px 41px, 48px 57px;
}

.auth-shell-dark .auth-paper {
  opacity: 0.06;
}

.auth-shell-dark .auth-paper::after {
  opacity: 0.06;
}

:global(html.dark) .auth-paper {
  opacity: 0.06;
}

:global(html.dark) .auth-paper::after {
  opacity: 0.06;
}

.auth-card {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(216, 205, 185, 0.9);
  background: linear-gradient(180deg, rgba(252, 249, 242, 0.88), rgba(247, 239, 228, 0.8));
  box-shadow: 0 28px 74px -52px rgba(31, 35, 32, 0.52), 0 1px 0 rgba(255, 255, 255, 0.5) inset;
  backdrop-filter: blur(14px);
}

.auth-card::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.2), transparent 24%),
    radial-gradient(circle at 82% 0%, rgba(173, 111, 51, 0.08), transparent 22%);
}

.auth-shell-dark .auth-card {
  border-color: rgba(194, 165, 117, 0.24);
  background: linear-gradient(180deg, rgba(27, 23, 19, 0.9), rgba(21, 18, 16, 0.94));
  box-shadow: inset 0 1px 0 rgba(248, 231, 200, 0.05), 0 24px 64px rgba(0, 0, 0, 0.28);
}

.auth-shell-dark .auth-card::before {
  background:
    linear-gradient(180deg, rgba(255, 247, 235, 0.04), transparent 28%),
    linear-gradient(90deg, transparent 0, rgba(176, 120, 57, 0.08) 48%, transparent 100%),
    radial-gradient(circle at 82% 0%, rgba(188, 115, 47, 0.14), transparent 24%);
}

.auth-shell-dark .auth-card-head {
  border-bottom-color: rgba(194, 165, 117, 0.24);
}

.auth-shell-dark .auth-card-head::after {
  background: linear-gradient(90deg, rgba(194, 147, 89, 0.2), transparent 72%);
}

.auth-card-head {
  position: relative;
  padding-bottom: 1.15rem;
  border-bottom: 1px solid rgba(216, 205, 185, 0.46);
}

.auth-card-head::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: -1px;
  height: 1px;
  background: linear-gradient(90deg, rgba(173, 111, 51, 0.18), transparent 72%);
}

:global(html.dark) .auth-card {
  border-color: rgba(194, 165, 117, 0.24);
  background: linear-gradient(180deg, rgba(27, 23, 19, 0.9), rgba(21, 18, 16, 0.94));
  box-shadow: inset 0 1px 0 rgba(248, 231, 200, 0.05), 0 24px 64px rgba(0, 0, 0, 0.28);
}

:global(html.dark) .auth-card::before {
  background:
    linear-gradient(180deg, rgba(255, 247, 235, 0.04), transparent 28%),
    linear-gradient(90deg, transparent 0, rgba(176, 120, 57, 0.08) 48%, transparent 100%),
    radial-gradient(circle at 82% 0%, rgba(188, 115, 47, 0.14), transparent 24%);
}

:global(html.dark) .auth-card-head {
  border-bottom-color: rgba(194, 165, 117, 0.24);
}

:global(html.dark) .auth-card-head::after {
  background: linear-gradient(90deg, rgba(194, 147, 89, 0.2), transparent 72%);
}

.auth-brand-bar {
  backdrop-filter: blur(3px);
}

.auth-shell-dark .auth-brand-bar > div {
  border-bottom-color: rgba(194, 165, 117, 0.22);
}

:global(html.dark) .auth-brand-bar > div {
  border-bottom-color: rgba(194, 165, 117, 0.22);
}

.auth-brand-link {
  transition: opacity 180ms ease, transform 180ms ease;
}

.auth-brand-mark,
.auth-card-mark {
  display: inline-grid;
  place-items: center;
  overflow: hidden;
  border-radius: 0.85rem;
  background: linear-gradient(180deg, rgba(31, 35, 32, 0.96), rgba(47, 42, 35, 0.92));
  box-shadow: 0 10px 24px rgba(31, 35, 32, 0.16);
}

.auth-brand-mark {
  width: 2.25rem;
  height: 2.25rem;
}

.auth-card-mark {
  width: 3rem;
  height: 3rem;
}

.auth-brand-link:hover {
  opacity: 0.82;
  transform: translateY(-1px);
}

.auth-shell-dark .auth-brand-link {
  opacity: 0.98;
}

.auth-shell-dark .auth-brand-link .font-serif {
  color: #f6ead7;
  text-shadow: 0 1px 0 rgba(255, 243, 223, 0.04);
}

.auth-shell-dark .auth-brand-link .text-\[10px\] {
  color: #cfb184;
}

:global(html.dark) .auth-brand-link {
  opacity: 0.98;
}

:global(html.dark) .auth-brand-link .font-serif {
  color: #f6ead7;
  text-shadow: 0 1px 0 rgba(255, 243, 223, 0.04);
}

:global(html.dark) .auth-brand-link .text-\[10px\] {
  color: #cfb184;
}

.auth-hero {
  max-width: 40rem;
}

.auth-hero-title {
  text-wrap: balance;
}

.auth-hero-copy {
  max-width: 30rem;
}

.auth-shell-dark .auth-hero-kicker,
.auth-shell-dark .auth-hero-copy,
.auth-shell-dark .auth-hero-marks {
  color: #c8b597;
}

.auth-shell-dark .auth-hero-marks .h-px,
.auth-shell-dark .auth-hero .mb-8 .h-px {
  background: linear-gradient(90deg, rgba(194, 165, 117, 0.08), rgba(194, 165, 117, 0.56), rgba(194, 165, 117, 0.08));
}

.auth-shell-dark .auth-hero-title,
.auth-shell-dark .auth-hero-lead {
  color: #f6ead7;
}

.auth-shell-dark .auth-hero-title {
  text-shadow: 0 1px 0 rgba(255, 243, 223, 0.05);
}

.auth-shell-dark .auth-hero-lead {
  color: #d6c4a7;
}

:global(html.dark) .auth-hero-kicker,
:global(html.dark) .auth-hero-copy,
:global(html.dark) .auth-hero-marks {
  color: #c8b597;
}

:global(html.dark) .auth-hero-marks .h-px,
:global(html.dark) .auth-hero .mb-8 .h-px {
  background: linear-gradient(90deg, rgba(194, 165, 117, 0.08), rgba(194, 165, 117, 0.56), rgba(194, 165, 117, 0.08));
}

:global(html.dark) .auth-hero-title,
:global(html.dark) .auth-hero-lead {
  color: #f6ead7;
}

:global(html.dark) .auth-hero-title {
  text-shadow: 0 1px 0 rgba(255, 243, 223, 0.05);
}

:global(html.dark) .auth-hero-lead {
  color: #d6c4a7;
}

:deep(.auth-form-title) {
  text-wrap: balance;
}

:deep(.auth-form-subtitle) {
  max-width: 24rem;
  margin-left: auto;
  margin-right: auto;
}

.auth-shell-dark .auth-card :deep(.auth-form-title),
.auth-shell-dark .auth-card :deep(.auth-form-subtitle),
.auth-shell-dark .auth-card :deep(.auth-footer-copy),
.auth-shell-dark .auth-card :deep(.auth-divider-label),
.auth-shell-dark .auth-card :deep(.auth-status-card),
.auth-shell-dark .auth-card :deep(.input-label),
.auth-shell-dark .auth-card :deep(.input),
.auth-shell-dark .auth-card :deep(.input::placeholder),
.auth-shell-dark .auth-card :deep(.input-hint),
.auth-shell-dark .auth-card :deep(.btn-primary),
.auth-shell-dark .auth-card :deep(.btn-primary:hover),
.auth-shell-dark .auth-card :deep(.btn-secondary),
.auth-shell-dark .auth-card :deep(.btn-secondary:hover),
.auth-shell-dark .auth-card :deep(a),
.auth-shell-dark .auth-card :deep(a:hover),
.auth-shell-dark .auth-card :deep(.text-gray-500),
.auth-shell-dark .auth-card :deep(.dark\:text-dark-400),
.auth-shell-dark .auth-card :deep(.bg-gray-200),
.auth-shell-dark .auth-card :deep(.dark\:bg-dark-700) {
  transition: none;
}

.auth-shell-dark .auth-card :deep(.auth-form-title),
.auth-shell-dark .auth-card-title,
.auth-shell-dark .auth-card .font-serif.text-2xl {
  color: #f1e6d3;
}

.auth-shell-dark .auth-card-kicker,
.auth-shell-dark .auth-card :deep(.auth-divider-label) {
  color: #b79a70;
}

.auth-shell-dark .auth-card :deep(.auth-form-subtitle),
.auth-shell-dark .auth-card :deep(.auth-footer-copy),
.auth-shell-dark .auth-card :deep(.text-gray-500),
.auth-shell-dark .auth-card :deep(.dark\:text-dark-400),
.auth-shell-dark .auth-card :deep(.input-hint) {
  color: #b4a48d;
}

.auth-shell-dark .auth-card :deep(.input-label) {
  color: #dccdb7;
}

.auth-shell-dark .auth-card :deep(.input) {
  border-color: rgba(145, 118, 81, 0.62);
  background: rgba(23, 20, 18, 0.84);
  color: #f0e4d2;
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.04);
}

.auth-shell-dark .auth-card :deep(.input::placeholder) {
  color: #a1937f;
}

.auth-shell-dark .auth-card :deep(.input:focus) {
  border-color: rgba(194, 129, 61, 0.8);
  box-shadow: 0 0 0 3px rgba(194, 129, 61, 0.16);
}

.auth-shell-dark .auth-card :deep(.btn-primary) {
  border: 1px solid rgba(214, 176, 122, 0.16);
  background: linear-gradient(135deg, rgba(244, 230, 205, 0.97), rgba(221, 198, 165, 0.94));
  color: #201914;
  box-shadow: 0 16px 34px rgba(76, 49, 24, 0.24), inset 0 1px 0 rgba(255, 246, 232, 0.3);
}

.auth-shell-dark .auth-card :deep(.btn-primary:hover) {
  background: linear-gradient(135deg, rgba(248, 236, 218, 0.99), rgba(229, 208, 177, 0.96));
  color: #1a140f;
}

.auth-shell-dark .auth-card :deep(.btn-secondary) {
  border-color: rgba(130, 108, 79, 0.56);
  background: rgba(23, 20, 18, 0.68);
  color: #e9dcc8;
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.04);
}

.auth-shell-dark .auth-card :deep(.btn-secondary:hover) {
  border-color: rgba(194, 147, 89, 0.38);
  background: rgba(43, 34, 27, 0.84);
  color: #f4e9d7;
}

.auth-shell-dark .auth-card :deep(a) {
  color: #d7bb90;
}

.auth-shell-dark .auth-card :deep(a:hover) {
  color: #f1d8b1;
}

.auth-shell-dark .auth-card :deep(.bg-gray-200),
.auth-shell-dark .auth-card :deep(.dark\:bg-dark-700) {
  background: linear-gradient(90deg, rgba(110, 94, 72, 0.2), rgba(194, 147, 89, 0.38), rgba(110, 94, 72, 0.2));
}

.auth-shell-dark .auth-card :deep(.auth-status-card) {
  border-color: rgba(176, 128, 76, 0.34);
  background: rgba(74, 43, 18, 0.16);
}

.auth-card :deep(.input-label) {
  color: #38413a;
}

:global(html.dark) .auth-card :deep(.input-label) {
  color: #dccdb7;
}

.auth-card :deep(.input) {
  min-height: 2.75rem;
  border-color: rgba(126, 112, 87, 0.34);
  background: rgba(255, 255, 255, 0.6);
}

:global(html.dark) .auth-card :deep(.input) {
  border-color: rgba(145, 118, 81, 0.62);
  background: rgba(23, 20, 18, 0.84);
  color: #f0e4d2;
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.04);
}

:global(html.dark) .auth-card :deep(.input::placeholder) {
  color: #a1937f;
}

.auth-card :deep(.input:focus) {
  border-color: #a73a2a;
  box-shadow: 0 0 0 3px rgba(167, 58, 42, 0.13);
}

:global(html.dark) .auth-card :deep(.input:focus) {
  border-color: rgba(194, 129, 61, 0.8);
  box-shadow: 0 0 0 3px rgba(194, 129, 61, 0.16);
}

.auth-card :deep(.btn-primary) {
  min-height: 2.85rem;
}

.auth-card :deep(.btn-secondary) {
  border-color: rgba(184, 163, 132, 0.48);
  background: rgba(255, 252, 247, 0.78);
  color: #2d261e;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.55);
}

.auth-card :deep(.btn-secondary:hover) {
  border-color: rgba(173, 111, 51, 0.34);
  background: rgba(248, 240, 228, 0.92);
  color: #221c16;
}

:global(html.dark) .auth-card :deep(.btn-primary) {
  border: 1px solid rgba(214, 176, 122, 0.16);
  background: linear-gradient(135deg, rgba(244, 230, 205, 0.97), rgba(221, 198, 165, 0.94));
  color: #201914;
  box-shadow: 0 16px 34px rgba(76, 49, 24, 0.24), inset 0 1px 0 rgba(255, 246, 232, 0.3);
}

:global(html.dark) .auth-card :deep(.btn-primary:hover) {
  background: linear-gradient(135deg, rgba(248, 236, 218, 0.99), rgba(229, 208, 177, 0.96));
  color: #1a140f;
}

:global(html.dark) .auth-card :deep(.btn-secondary) {
  border-color: rgba(130, 108, 79, 0.56);
  background: rgba(23, 20, 18, 0.68);
  color: #e9dcc8;
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.04);
}

:global(html.dark) .auth-card :deep(.btn-secondary:hover) {
  border-color: rgba(194, 147, 89, 0.38);
  background: rgba(43, 34, 27, 0.84);
  color: #f4e9d7;
}

.auth-card :deep(a) {
  color: #8a5f34;
}

.auth-card :deep(a:hover) {
  color: #a73a2a;
}

:global(html.dark) .auth-card :deep(a) {
  color: #d7bb90;
}

:global(html.dark) .auth-card :deep(a:hover) {
  color: #f1d8b1;
}

.auth-card :deep(.bg-gray-200) {
  background-color: rgba(216, 205, 185, 0.78);
}

.auth-divider-label {
  color: #9b7a52;
  letter-spacing: 0.22em;
}

.auth-footer-copy {
  color: #7f6b4d;
}

.auth-shell-footer :deep(.auth-footer-copy) {
  color: #6f6556;
  letter-spacing: 0.04em;
}

.auth-shell-footer :deep(.auth-inline-link) {
  color: #2f6f5e;
  text-decoration: underline;
  text-decoration-color: rgba(47, 111, 94, 0.28);
  text-underline-offset: 0.22em;
}

.auth-shell-footer :deep(.auth-inline-link:hover) {
  color: #a73a2a;
  text-decoration-color: rgba(167, 58, 42, 0.3);
}

.auth-shell-copyright {
  color: #7a715f;
  letter-spacing: 0.12em;
}

.auth-card :deep(.text-amber-600),
.auth-card :deep(.text-amber-700) {
  color: #96602d;
}

.auth-card :deep(.bg-green-50),
.auth-card :deep(.dark\:bg-green-900\/20) {
  border: 1px solid rgba(99, 145, 93, 0.18);
  background: rgba(236, 248, 233, 0.88);
}

.auth-card :deep(.text-green-700),
.auth-card :deep(.text-green-600) {
  color: #3f6f42;
}

:global(html.dark) .auth-card :deep(.bg-gray-200),
:global(html.dark) .auth-card :deep(.dark\:bg-dark-700) {
  background: linear-gradient(90deg, rgba(110, 94, 72, 0.2), rgba(194, 147, 89, 0.38), rgba(110, 94, 72, 0.2));
}

:global(html.dark) .auth-card :deep(.text-gray-500),
:global(html.dark) .auth-card :deep(.dark\:text-dark-400),
:global(html.dark) .auth-card :deep(.auth-footer-copy),
:global(html.dark) .auth-card :deep(.auth-form-subtitle),
:global(html.dark) .auth-card :deep(.input-hint) {
  color: #b4a48d;
}

:global(html.dark) .auth-card :deep(.auth-form-title),
:global(html.dark) .auth-card-title,
:global(html.dark) .auth-card .font-serif.text-2xl {
  color: #f1e6d3;
}

:global(html.dark) .auth-card-kicker {
  color: #b79a70;
}

:global(html.dark) .auth-card :deep(.auth-status-card) {
  border-color: rgba(176, 128, 76, 0.34);
  background: rgba(74, 43, 18, 0.16);
}

:global(html.dark) .auth-card :deep(.text-amber-600),
:global(html.dark) .auth-card :deep(.text-amber-700),
:global(html.dark) .auth-card :deep(.dark\:text-amber-400) {
  color: #e0b177;
}

:global(html.dark) .auth-card :deep(.bg-green-50),
:global(html.dark) .auth-card :deep(.dark\:bg-green-900\/20) {
  border: 1px solid rgba(110, 148, 90, 0.24);
  background: rgba(34, 54, 33, 0.36);
}

:global(html.dark) .auth-card :deep(.text-green-700),
:global(html.dark) .auth-card :deep(.text-green-600),
:global(html.dark) .auth-card :deep(.dark\:text-green-400) {
  color: #9fd2a0;
}

.auth-shell-dark .auth-shell-footer :deep(.auth-footer-copy) {
  color: #d6c4a7;
  text-shadow: 0 1px 0 rgba(12, 14, 10, 0.3);
}

.auth-shell-dark .auth-shell-footer :deep(.auth-inline-link) {
  color: #d7bb90;
  font-weight: 600;
  text-decoration-color: rgba(215, 187, 144, 0.42);
  text-shadow: 0 1px 0 rgba(12, 14, 10, 0.28);
}

.auth-shell-dark .auth-shell-footer :deep(.auth-inline-link:hover) {
  color: #f1d8b1;
  text-decoration-color: rgba(241, 216, 177, 0.56);
}

.auth-shell-dark .auth-shell-copyright {
  color: #bda98c;
  text-shadow: 0 1px 0 rgba(12, 14, 10, 0.28);
}

@media (max-width: 768px) {
  .auth-bg {
    background-image:
      linear-gradient(180deg, rgba(244, 239, 228, 0.9) 0%, rgba(244, 239, 228, 0.5) 58%, rgba(244, 239, 228, 0.24) 100%),
      var(--sst-auth-bg);
    background-position: center top, 58% bottom;
    background-size: cover, auto 92%;
  }

  :global(html.dark) .auth-bg {
    background-image:
      linear-gradient(180deg, rgba(19, 17, 15, 0.82) 0%, rgba(19, 17, 15, 0.58) 58%, rgba(19, 17, 15, 0.88) 100%),
      radial-gradient(circle at 16% 18%, rgba(236, 219, 189, 0.05), transparent 26%),
      radial-gradient(circle at 82% 10%, rgba(193, 129, 62, 0.14), transparent 22%),
      var(--sst-auth-bg);
    background-position: center top, center, center, 58% bottom;
    background-size: cover, cover, cover, auto 92%;
    opacity: 1;
    filter: grayscale(0.16) brightness(0.7) contrast(0.98) saturate(0.9);
  }

  :global(html.dark) .auth-shell {
    background: #13110f;
  }

  :global(html.dark) .auth-paper,
  :global(html.dark) .auth-paper::after {
    opacity: 0.06;
  }

  :global(html.dark) .auth-brand-link .font-serif,
  :global(html.dark) .auth-hero-title,
  :global(html.dark) .auth-hero-lead,
  :global(html.dark) .auth-card :deep(.auth-form-title),
  :global(html.dark) .auth-card-title,
  :global(html.dark) .auth-card .font-serif.text-2xl {
    color: #f6ead7;
    text-shadow: 0 1px 0 rgba(255, 243, 223, 0.05);
  }

  :global(html.dark) .auth-brand-link .text-\[10px\],
  :global(html.dark) .auth-hero-kicker,
  :global(html.dark) .auth-hero-copy,
  :global(html.dark) .auth-hero-marks,
  :global(html.dark) .auth-card-kicker,
  :global(html.dark) .auth-card :deep(.text-gray-500),
  :global(html.dark) .auth-card :deep(.dark\:text-dark-400),
  :global(html.dark) .auth-card :deep(.auth-footer-copy),
  :global(html.dark) .auth-card :deep(.auth-form-subtitle),
  :global(html.dark) .auth-card :deep(.input-hint) {
    color: #c8b597;
  }

  :global(html.dark) .auth-hero-marks .h-px,
  :global(html.dark) .auth-hero .mb-8 .h-px,
  :global(html.dark) .auth-card :deep(.bg-gray-200),
  :global(html.dark) .auth-card :deep(.dark\:bg-dark-700) {
    background: linear-gradient(90deg, rgba(194, 165, 117, 0.08), rgba(194, 165, 117, 0.56), rgba(194, 165, 117, 0.08));
  }

  :global(html.dark) .auth-card {
    border-color: rgba(194, 165, 117, 0.24);
    background: linear-gradient(180deg, rgba(27, 23, 19, 0.92), rgba(21, 18, 16, 0.95));
    box-shadow: inset 0 1px 0 rgba(248, 231, 200, 0.05), 0 24px 64px rgba(0, 0, 0, 0.28);
  }

  :global(html.dark) .auth-card :deep(.input-label) {
    color: #dccdb7;
  }

  :global(html.dark) .auth-card :deep(.input) {
    border-color: rgba(145, 118, 81, 0.62);
    background: rgba(23, 20, 18, 0.84);
    color: #f0e4d2;
    box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.04);
  }

  :global(html.dark) .auth-card :deep(.input::placeholder) {
    color: #a1937f;
  }

  :global(html.dark) .auth-card :deep(.input:focus) {
    border-color: rgba(194, 129, 61, 0.8);
    box-shadow: 0 0 0 3px rgba(194, 129, 61, 0.16);
  }

  :global(html.dark) .auth-card :deep(.btn-primary) {
    border: 1px solid rgba(214, 176, 122, 0.16);
    background: linear-gradient(135deg, rgba(244, 230, 205, 0.97), rgba(221, 198, 165, 0.94));
    color: #201914;
    box-shadow: 0 16px 34px rgba(76, 49, 24, 0.24), inset 0 1px 0 rgba(255, 246, 232, 0.3);
  }

  :global(html.dark) .auth-card :deep(.btn-primary:hover) {
    background: linear-gradient(135deg, rgba(248, 236, 218, 0.99), rgba(229, 208, 177, 0.96));
    color: #1a140f;
  }

  :global(html.dark) .auth-card :deep(a) {
    color: #d7bb90;
  }

  :global(html.dark) .auth-card :deep(a:hover) {
    color: #f1d8b1;
  }
}
</style>
