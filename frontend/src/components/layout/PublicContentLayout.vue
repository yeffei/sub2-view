<template>
  <div class="sst-public min-h-screen" :style="bgStyle">
    <div class="sst-public__bg" aria-hidden="true"></div>
    <PublicSiteHeader :nav-items="resolvedNavItems" />

    <main class="relative z-10 mx-auto grid max-w-7xl gap-8 px-5 pb-16 pt-4 sm:px-8 lg:grid-cols-[0.85fr_1.15fr] lg:gap-14 lg:pb-24 lg:pt-10">
      <section class="max-w-2xl lg:pt-8">
        <div class="mb-6 flex items-center gap-4">
          <span class="h-px w-14 bg-zen-paperLine/80 dark:bg-zen-nightLine"></span>
          <span class="text-xs uppercase tracking-[0.42em] text-zen-mist dark:text-zen-stone">{{ eyebrow }}</span>
        </div>
        <h1 class="font-serif text-[clamp(2.5rem,5vw,4.8rem)] font-semibold leading-[1.02] text-zen-ink dark:text-zen-paper">{{ title }}</h1>
        <p v-if="lead" class="mt-6 max-w-xl text-base leading-8 text-zen-mist dark:text-zen-stone sm:text-lg">{{ lead }}</p>
        <div v-if="$slots.summary" class="mt-8 hidden lg:block"><slot name="summary" /></div>
      </section>

      <section class="panel rounded-[28px] border border-zen-paperLine/70 bg-white/60 p-6 dark:border-zen-nightLine dark:bg-zen-nightPanel/76 sm:p-8 lg:p-10">
        <slot />
      </section>
    </main>

    <div v-if="$slots.summary" class="relative z-10 mx-auto max-w-7xl px-5 pb-16 sm:px-8 lg:hidden"><slot name="summary" /></div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore, useAuthStore } from '@/stores'
import PublicSiteHeader from '@/components/layout/PublicSiteHeader.vue'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

interface NavItem {
  to: string
  label: string
}

interface PublicContentLayoutProps {
  title?: string
  lead?: string
  eyebrow?: string
  navItems?: NavItem[]
}

const props = withDefaults(defineProps<PublicContentLayoutProps>(), {
  title: '',
  lead: '',
  eyebrow: 'SST',
  navItems: () => []
})

const appStore = useAppStore()
const authStore = useAuthStore()
const bgStyle = computed(() => ({ '--sst-public-bg': `url(${paperInkBg})` }))
const resolvedNavItems = computed(() => props.navItems.length > 0
  ? props.navItems
  : [
    { to: '/pricing', label: '价目' },
    { to: '/docs', label: '文档' },
  ])

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
})
</script>

<style scoped>
.sst-public { position: relative; overflow: hidden; background: #f4efe4; }
.dark .sst-public { background: #11130f; }
.sst-public__bg {
  position: absolute; inset: 0; pointer-events: none;
  background-image: linear-gradient(90deg, rgba(244,239,228,.92) 0%, rgba(244,239,228,.56) 40%, rgba(244,239,228,.16) 100%), var(--sst-public-bg);
  background-size: cover, cover; background-position: center, center bottom;
}
.dark .sst-public__bg { opacity: .24; filter: grayscale(.9) invert(.92); }
.panel { position: relative; box-shadow: 0 24px 70px rgba(58,48,32,.08); backdrop-filter: blur(2px); }
.panel::before {
  content: ''; position: absolute; inset: 1rem; border: 1px solid rgba(126,112,87,.16); border-radius: 22px; pointer-events: none;
}
.dark .panel::before { border-color: rgba(216,205,185,.1); }
</style>


