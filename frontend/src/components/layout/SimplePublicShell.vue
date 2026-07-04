<template>
  <div
    class="simple-public-shell legal-public-shell flex min-h-screen flex-col overflow-hidden"
    :style="backgroundStyle"
  >
    <div class="legal-public-bg" aria-hidden="true"></div>
    <div class="legal-public-paper" aria-hidden="true"></div>
    <div class="legal-public-wash" aria-hidden="true"></div>

    <PublicSiteHeader />
    <div class="relative z-10 flex-1">
      <slot />
    </div>
    <div class="relative z-10">
      <PublicSiteFooter />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import PublicSiteFooter from '@/components/layout/PublicSiteFooter.vue'
import PublicSiteHeader from '@/components/layout/PublicSiteHeader.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { useThemeState } from '@/utils/theme'
import paperInkBg from '@/assets/brand/sst-paper-ink-bg.png'

const appStore = useAppStore()
const authStore = useAuthStore()
const isDark = useThemeState()
const backgroundStyle = computed(() => ({
  '--sst-public-bg': `url(${paperInkBg})`,
  backgroundColor: isDark.value ? '#161915' : '#f4efe4',
  color: isDark.value ? '#f3eadb' : '#1f2320',
}))

onMounted(() => {
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.simple-public-shell {
  position: relative;
  min-height: 100vh;
}

.legal-public-bg {
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

:global(.dark) .legal-public-bg {
  display: none;
}

.legal-public-paper {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.14;
  background-image:
    linear-gradient(rgba(31, 35, 32, 0.022) 1px, transparent 1px),
    linear-gradient(90deg, rgba(31, 35, 32, 0.016) 1px, transparent 1px);
  background-size: 128px 128px, 128px 128px;
}

:global(.dark) .legal-public-paper {
  opacity: 0.04;
}

.legal-public-paper::after {
  content: '';
  position: absolute;
  inset: 0;
  opacity: 0.18;
  background-image:
    radial-gradient(circle at 14% 18%, rgba(31, 35, 32, 0.06) 0 1px, transparent 1.5px),
    radial-gradient(circle at 72% 42%, rgba(31, 35, 32, 0.045) 0 1px, transparent 1.5px);
  background-size: 34px 41px, 48px 57px;
}

:global(.dark) .legal-public-paper::after {
  opacity: 0.04;
}

.legal-public-wash {
  position: absolute;
  inset: auto 0 0 0;
  height: 24rem;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 40%, rgba(77, 126, 126, 0.08), transparent 34%),
    radial-gradient(circle at 78% 12%, rgba(126, 112, 87, 0.08), transparent 28%),
    linear-gradient(180deg, transparent, rgba(255, 255, 255, 0.16));
}

:global(.dark) .legal-public-wash {
  display: none;
}

@media (max-width: 768px) {
  .legal-public-bg {
    background-image:
      linear-gradient(180deg, rgba(244, 239, 228, 0.9) 0%, rgba(244, 239, 228, 0.34) 54%, rgba(244, 239, 228, 0.14) 100%),
      var(--sst-public-bg);
    background-position: center top, 58% bottom;
    background-size: cover, auto 92%;
  }
}
</style>


