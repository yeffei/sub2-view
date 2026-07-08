<template>
  <footer class="relative z-10 px-5 pb-2 sm:px-8 sm:pb-2.5">
    <div class="public-site-footer-shell mx-auto flex max-w-7xl flex-col gap-3 border-t border-zen-paperLine/70 pt-4 text-[0.92rem] text-zen-mist dark:border-zen-nightLine dark:text-zen-stone sm:flex-row sm:items-center sm:justify-between">
      <div class="footer-note">{{ brandName }} · {{ t('publicSite.tagline') }}</div>
      <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
        <RouterLink to="/home" class="footer-link">{{ t('publicSite.nav.home') }}</RouterLink>
        <RouterLink to="/pricing" class="footer-link">{{ t('publicSite.nav.pricing') }}</RouterLink>
        <RouterLink to="/docs" class="footer-link">{{ t('publicSite.nav.docs') }}</RouterLink>
        <RouterLink to="/faq" class="footer-link">{{ t('publicSite.nav.faq') }}</RouterLink>
        <RouterLink v-if="publicContact" to="/faq#public-contact" class="footer-link">{{ t('publicSite.nav.contact') }}</RouterLink>
        <RouterLink to="/terms" class="footer-link">{{ t('publicSite.nav.terms') }}</RouterLink>
        <RouterLink to="/privacy" class="footer-link">{{ t('publicSite.nav.privacy') }}</RouterLink>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { resolvePublicContact } from '@/utils/contact'

const appStore = useAppStore()
const { t } = useI18n()

const configuredName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '')
const brandName = computed(() => configuredName.value && configuredName.value !== 'Sub2API' ? configuredName.value : '山枢庭')
const publicContact = computed(() => resolvePublicContact(appStore.cachedPublicSettings?.contact_info || appStore.contactInfo))
</script>

<style scoped>
.footer-note {
  letter-spacing: 0.01em;
}

.footer-link {
  font-size: 0.88rem;
  letter-spacing: 0.08em;
  transition: color 180ms ease;
}

.footer-link:hover {
  color: #a73a2a;
}

</style>

<style>
html.dark .sst-public.sst-public-tone-pricing .public-site-footer-shell,
html.dark .sst-public.sst-public-tone-docs .public-site-footer-shell {
  border-top-color: rgba(173, 145, 104, 0.2);
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.025);
}

html.dark .sst-public.sst-public-tone-pricing .footer-note,
html.dark .sst-public.sst-public-tone-docs .footer-note {
  color: #b4a58d;
}

html.dark .sst-public.sst-public-tone-pricing .footer-link,
html.dark .sst-public.sst-public-tone-docs .footer-link {
  color: #c9b28d;
}

html.dark .sst-public.sst-public-tone-pricing .footer-link:hover,
html.dark .sst-public.sst-public-tone-docs .footer-link:hover {
  color: #f0c487;
}

html.dark .sst-public-tone-legal .footer-note {
  color: #9f8e74;
}

html.dark .sst-public-tone-legal .footer-link {
  color: #bfa983;
}

html.dark .sst-public-tone-legal .footer-link:hover {
  color: #f0be79;
}

html.dark .sst-home .footer-note {
  color: #aa9980;
}

html.dark .sst-home .public-site-footer-shell {
  border-top-color: rgba(173, 145, 104, 0.18);
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.02);
}

html.dark .sst-home .footer-link {
  color: #bda17a;
}

html.dark .sst-home .footer-link:hover {
  color: #efc183;
}

html.dark .sst-public .public-site-footer-shell,
html.dark .sst-public.sst-public-tone-pricing .public-site-footer-shell,
html.dark .sst-public.sst-public-tone-docs .public-site-footer-shell,
html.dark .sst-public.sst-public-tone-legal .public-site-footer-shell {
  border-top-color: rgba(173, 145, 104, 0.18);
  box-shadow: inset 0 1px 0 rgba(255, 243, 223, 0.02);
}

html.dark .sst-public .footer-note,
html.dark .sst-public.sst-public-tone-pricing .footer-note,
html.dark .sst-public.sst-public-tone-docs .footer-note,
html.dark .sst-public.sst-public-tone-legal .footer-note {
  color: #aa9980;
}

html.dark .sst-public .footer-link,
html.dark .sst-public.sst-public-tone-pricing .footer-link,
html.dark .sst-public.sst-public-tone-docs .footer-link,
html.dark .sst-public.sst-public-tone-legal .footer-link {
  color: #bda17a;
}

html.dark .sst-public .footer-link:hover,
html.dark .sst-public.sst-public-tone-pricing .footer-link:hover,
html.dark .sst-public.sst-public-tone-docs .footer-link:hover,
html.dark .sst-public.sst-public-tone-legal .footer-link:hover {
  color: #efc183;
}

html:not(.dark) .sst-public .public-site-footer-shell {
  border-top-color: rgba(154, 128, 92, 0.18);
  box-shadow: inset 0 1px 0 rgba(255, 249, 239, 0.52);
}

html:not(.dark) .sst-public .footer-note {
  color: #7a6a53;
}

html:not(.dark) .sst-public .footer-link {
  color: #8f6f43;
}

html:not(.dark) .sst-public .footer-link:hover {
  color: #a73a2a;
}
</style>
