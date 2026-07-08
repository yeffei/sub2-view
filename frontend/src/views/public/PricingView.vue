<template>
  <PublicPageLayout
    class="pricing-page"
    tone="pricing"
    :eyebrow="t('publicPricing.eyebrow')"
    :title="t('publicPricing.title')"
    intro=""
    description=""
    :authenticated-action-label="t('publicSite.enter')"
  >
    <div class="space-y-4">
      <div class="pricing-platform-tabs" role="tablist" :aria-label="t('publicPricing.platformSwitch')">
        <button
          v-for="platform in platforms"
          :key="platform.key"
          type="button"
          class="pricing-platform-tab"
          :class="[activePlatform === platform.key ? 'is-active' : '', platformTabClass(platform.key)]"
          :aria-pressed="activePlatform === platform.key"
          @click="selectPlatform(platform.key)"
        >
          <span class="pricing-platform-icon">{{ platform.icon }}</span>
          <span class="pricing-platform-label">{{ platform.label }}</span>
        </button>
      </div>

      <article class="pricing-ledger overflow-hidden rounded-[1.2rem] border border-zen-paperLine/70 bg-[linear-gradient(180deg,rgba(255,255,255,0.66)_0%,rgba(251,245,236,0.78)_100%)] shadow-paper dark:border-[rgba(86,92,80,0.58)] dark:bg-[linear-gradient(180deg,rgba(13,15,13,0.98)_0%,rgba(24,21,18,0.96)_100%)] dark:shadow-[0_24px_60px_rgba(0,0,0,0.3)]">
        <div
          class="pricing-group-grid px-4 py-4 sm:px-5 sm:py-5"
          :class="{ 'pricing-group-grid--single': activeGroups.length === 1 }"
        >
          <template v-if="activeGroups.length > 0">
            <section
              v-for="group in activeGroups"
              :key="group.key"
              class="pricing-group-card"
              :class="[group.key === activeGroupKey ? 'is-active' : '', groupCardClass(group.key)]"
              @click="activeGroupKey = group.key"
            >
              <div class="pricing-group-head">
                <div class="pricing-group-title-row">
                  <div class="pricing-group-name">{{ group.name }}</div>
                  <div class="pricing-group-note">{{ group.subtitle }}</div>
                </div>
                <div class="pricing-group-badge">{{ group.rateText }}</div>
              </div>
            </section>
          </template>

          <section v-else class="pricing-empty-card">
            {{ t('publicPricing.empty') }}
          </section>
        </div>

        <div
          v-if="activeHasPricingRows"
          class="pricing-table-wrap border-t border-zen-paperLine/70 bg-[linear-gradient(180deg,rgba(255,252,247,0.40)_0%,rgba(249,241,231,0.58)_100%)] dark:border-[rgba(72,78,69,0.78)] dark:bg-[linear-gradient(180deg,rgba(17,19,16,0.82)_0%,rgba(33,26,22,0.72)_100%)]"
        >
          <div class="overflow-x-auto">
            <table class="pricing-table w-full">
              <thead>
                <tr>
                  <th>{{ t('publicPricing.table.modelId') }}</th>
                  <th>{{ t('publicPricing.table.input') }}</th>
                  <th>{{ t('publicPricing.table.output') }}</th>
                  <th v-if="activePlatform === 'claude'">{{ t('publicPricing.table.cacheWrite') }}</th>
                  <th>{{ t('publicPricing.table.cacheRead') }}</th>
                </tr>
              </thead>
              <tbody v-show="activePlatform === 'codex'">
                <tr v-for="item in activeModels" :key="item.name">
                  <td :data-label="t('publicPricing.table.model')">
                    <div class="pricing-model-cell">
                      <div class="pricing-model-name-row">
                        <span class="pricing-model-name-text">{{ item.name }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.input')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.inputPrice, activeGroup.multiplier)) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.output')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.outputPrice, activeGroup.multiplier)) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.cacheRead')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.cacheReadPrice, activeGroup.multiplier)) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                </tr>
              </tbody>
              <tbody v-show="activePlatform === 'claude'">
                <tr v-for="item in claudeCodeModels" :key="item.name">
                  <td :data-label="t('publicPricing.table.model')">
                    <div class="pricing-model-cell">
                      <div class="pricing-model-name-row">
                        <span class="pricing-model-name-text">{{ item.name }}</span>
                      </div>
                    </div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.input')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.inputPrice, activeGroup.multiplier), 2) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.output')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.outputPrice, activeGroup.multiplier), 2) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.cacheWrite')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.cacheWritePrice, activeGroup.multiplier), 2) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                  <td class="pricing-value-cell" :data-label="t('publicPricing.table.cacheRead')">
                    <div class="pricing-price-main text-zen-seal dark:text-[#f0a25f]">{{ formatUsdPrice(scalePrice(item.cacheReadPrice, activeGroup.multiplier), 2) }} <span class="pricing-price-unit dark:text-[#8f8572]">/ 1M tokens</span></div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </article>

      <div v-if="publicContact" class="pricing-contact-note">
        <span>{{ t('publicPricing.contactPrefix') }}</span>
        <a
          v-if="publicContact.href"
          :href="publicContact.href"
          :target="publicContact.external ? '_blank' : undefined"
          :rel="publicContact.external ? 'noopener noreferrer' : undefined"
        >
          {{ t('publicPricing.contactAction') }}
        </a>
        <strong v-else>{{ t('publicPricing.contactAction') }}: {{ publicContact.label }}</strong>
      </div>
    </div>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'
import { useAppStore } from '@/stores'
import { resolvePublicContact } from '@/utils/contact'

interface ModelRow {
  name: string
  inputPrice: number | null
  outputPrice: number | null
  cacheReadPrice: number | null
}

interface GroupRow {
  key: string
  name: string
  subtitle: string
  rateText: string
  multiplier: number
}

type PlatformKey = 'claude' | 'codex'

const activePlatform = ref<PlatformKey>('codex')
const activeGroupKey = ref('codex-plus')
const perMillionScale = 1_000_000
const appStore = useAppStore()
const { t } = useI18n()

const platforms = [
  { key: 'claude', label: 'Claude Code', icon: '✳' },
  { key: 'codex', label: 'Codex', icon: '◈' },
] as const

const codexModels: ModelRow[] = [
  { name: 'gpt-5.5', inputPrice: 5e-6, outputPrice: 30e-6, cacheReadPrice: 0.5e-6 },
  { name: 'gpt-5.4', inputPrice: 2.5e-6, outputPrice: 15e-6, cacheReadPrice: 0.25e-6 },
  { name: 'gpt-5.4-mini', inputPrice: 7.5e-7, outputPrice: 4.5e-6, cacheReadPrice: 7.5e-8 },
  { name: 'gpt-4o', inputPrice: 2.5e-6, outputPrice: 1e-5, cacheReadPrice: 1.25e-6 },
]

const codexGroups: GroupRow[] = [
  { key: 'codex-plus', name: 'Codex Plus', subtitle: '0.16x', rateText: '0.16x', multiplier: 0.16 },
  { key: 'codex-pro', name: 'Codex Pro', subtitle: '0.45x', rateText: '0.45x', multiplier: 0.45 },
]

const claudeCodeGroups: GroupRow[] = [
  { key: 'claude-code-max', name: 'CC MAX', subtitle: '1.8x', rateText: '1.8x', multiplier: 1.8 },
]

const claudeCodeModels = [
  { name: 'claude-fable-5', inputPrice: 10e-6, outputPrice: 50e-6, cacheWritePrice: 12.5e-6, cacheReadPrice: 1e-6 },
  { name: 'claude-opus-4-8', inputPrice: 5e-6, outputPrice: 25e-6, cacheWritePrice: 6.25e-6, cacheReadPrice: 0.5e-6 },
  { name: 'claude-opus-4-7', inputPrice: 5e-6, outputPrice: 25e-6, cacheWritePrice: 6.25e-6, cacheReadPrice: 0.5e-6 },
  { name: 'claude-opus-4-6', inputPrice: 5e-6, outputPrice: 25e-6, cacheWritePrice: 6.25e-6, cacheReadPrice: 0.5e-6 },
]

const groupsByPlatform: Record<PlatformKey, GroupRow[]> = {
  claude: claudeCodeGroups,
  codex: codexGroups,
}

const modelsByPlatform: Record<PlatformKey, ModelRow[]> = {
  claude: [],
  codex: codexModels,
}

const activeGroups = computed(() => groupsByPlatform[activePlatform.value])
const activeModels = computed(() => modelsByPlatform[activePlatform.value])
const activeGroup = computed(() => activeGroups.value.find(group => group.key === activeGroupKey.value) || activeGroups.value[0] || codexGroups[0])
const activeHasPricingRows = computed(() => activePlatform.value === 'codex' ? activeModels.value.length > 0 : claudeCodeModels.length > 0)
const publicContact = computed(() => resolvePublicContact(appStore.cachedPublicSettings?.contact_info || appStore.contactInfo))

function selectPlatform(key: PlatformKey) {
  activePlatform.value = key
  activeGroupKey.value = groupsByPlatform[key][0]?.key || ''
}

function platformTabClass(key: PlatformKey) {
  if (activePlatform.value === key) {
    return 'border-[rgba(188,93,31,0.35)] bg-[rgba(188,93,31,0.08)] text-zen-ink shadow-[0_8px_18px_rgba(167,58,42,0.08)] dark:border-[rgba(188,93,31,0.42)] dark:bg-[linear-gradient(135deg,rgba(56,41,29,0.94),rgba(28,31,24,0.96))] dark:text-zen-paper dark:shadow-[0_12px_28px_rgba(0,0,0,0.2)]'
  }

  return 'text-zen-inkSoft hover:border-[rgba(188,93,31,0.22)] hover:bg-[rgba(188,93,31,0.04)] dark:border-[rgba(78,84,73,0.68)] dark:bg-[rgba(21,24,19,0.74)] dark:text-[#c8c0ad] dark:hover:border-[rgba(188,93,31,0.3)] dark:hover:bg-[rgba(53,39,28,0.82)] dark:hover:text-zen-paper'
}

function groupCardClass(key: string) {
  if (activeGroupKey.value === key) {
    return 'border-[rgba(188,93,31,0.35)] bg-[rgba(188,93,31,0.06)] shadow-[0_10px_22px_rgba(167,58,42,0.06)] dark:border-[rgba(188,93,31,0.42)] dark:bg-[linear-gradient(135deg,rgba(51,37,28,0.94),rgba(26,30,23,0.95))] dark:shadow-[0_14px_32px_rgba(0,0,0,0.2)]'
  }

  return 'hover:border-[rgba(188,93,31,0.24)] hover:bg-[rgba(188,93,31,0.03)] dark:border-[rgba(78,84,73,0.64)] dark:bg-[rgba(23,25,20,0.84)] dark:hover:border-[rgba(188,93,31,0.28)] dark:hover:bg-[rgba(41,31,24,0.84)]'
}

function formatUsdPrice(value: number | null, minimumFractionDigits = 0) {
  if (value == null) return '-'
  const formatted = new Intl.NumberFormat('en-US', {
    minimumFractionDigits,
    maximumFractionDigits: 2,
  }).format(value * perMillionScale)
  return `$${formatted}`
}

function scalePrice(value: number | null, multiplier: number): number | null {
  return value == null ? null : value * multiplier
}
</script>

<style scoped>
.pricing-page :deep(main > section:first-child) {
  gap: 0.35rem;
  padding-top: 1.2rem;
  padding-bottom: 0.45rem;
  grid-template-columns: minmax(0, 1fr);
}

.pricing-page :deep(.sst-public-wash) { display: none; }
.pricing-page :deep(.public-copy-block) { max-width: 100%; }
.pricing-page :deep(.public-display) { font-size: clamp(2.2rem, 4.6vw, 3.4rem); font-weight: 400; }
.pricing-page :deep(.public-intro) { display: none; }
.pricing-page :deep(.public-copy-block > p:last-of-type) { margin-top: 0.35rem; max-width: 34rem; font-size: 0.88rem; line-height: 1.7; font-weight: 400; color: rgba(123, 106, 83, 0.92); }
.pricing-page :deep(.public-hero-panel) { position: relative; top: auto; min-height: 0; padding: 0; background: transparent; box-shadow: none; }
.pricing-page :deep(.public-hero-panel::before), .pricing-page :deep(.public-hero-mark), .pricing-page :deep(.public-hero-axis), .pricing-page :deep(.public-hero-seal) { display: none; }
.pricing-page :deep(.public-copy-block), .pricing-page :deep(.public-hero-panel) { animation: none; }
.pricing-page :deep(.public-display) { color: #2f281f; text-shadow: 0 1px 0 rgba(255, 248, 238, 0.55); }
.pricing-page :deep(.public-copy-block > div:first-child span:last-child) { color: #9b7a52; }
.pricing-page :deep(.public-cta) {
  padding: clamp(1.4rem, 2vw, 2rem) clamp(1.35rem, 2.6vw, 2.35rem);
  border: 1px solid rgba(176, 150, 118, 0.54);
  background:
    linear-gradient(140deg, rgba(255, 252, 247, 0.96) 0%, rgba(247, 239, 228, 0.98) 58%, rgba(241, 232, 220, 0.94) 100%);
  box-shadow: 0 16px 40px rgba(111, 87, 56, 0.1);
}
.pricing-page :deep(.public-cta-mark) {
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.36), transparent 28%),
    radial-gradient(circle at 82% 24%, rgba(196, 136, 68, 0.12), transparent 24%),
    linear-gradient(180deg, rgba(255, 248, 238, 0.35), transparent 62%);
}
.pricing-page :deep(.public-cta > .relative) {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) auto;
  align-items: center;
  gap: 1.6rem;
}
.pricing-page :deep(.public-cta > .relative > div:first-child) {
  max-width: 40rem;
}
.pricing-page :deep(.public-cta > .relative > div:first-child > div:first-child) {
  letter-spacing: 0.42em;
  color: #9b7a52;
}
.pricing-page :deep(.public-cta h2) {
  margin-top: 0.55rem;
  max-width: 13ch;
  font-size: clamp(2rem, 3.2vw, 2.7rem);
  line-height: 1.12;
  letter-spacing: 0.01em;
  color: #2d261e;
}
.pricing-page :deep(.public-cta p) {
  margin-top: 0.8rem;
  max-width: 28rem;
  font-size: 0.96rem;
  line-height: 1.9;
  color: rgba(102, 85, 63, 0.88);
}
.pricing-page :deep(.public-cta > .relative > div:last-child) {
  position: relative;
  align-self: center;
  justify-self: end;
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem;
  padding-left: 1.4rem;
}
.pricing-page :deep(.public-cta > .relative > div:last-child::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 0.2rem;
  bottom: 0.2rem;
  width: 1px;
  background: linear-gradient(180deg, transparent, rgba(166, 136, 96, 0.72), transparent);
}
.pricing-page :deep(.public-cta > .relative > div:last-child a) {
  min-width: 10.5rem;
  padding: 0.9rem 1.45rem;
  border-radius: 1rem;
  font-size: 0.95rem;
  letter-spacing: 0.02em;
}
.pricing-page :deep(.public-cta > .relative > div:last-child a:first-child) {
  box-shadow: 0 12px 24px rgba(150, 101, 47, 0.16);
}
.pricing-page :deep(.public-cta > .relative > div:last-child a:last-child) {
  border-color: rgba(190, 157, 118, 0.54);
  background: rgba(255, 251, 245, 0.8);
  color: #7a6041;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.55);
}
.pricing-page :deep(.public-cta > .relative > div:last-child a:last-child:hover) {
  border-color: rgba(196, 136, 68, 0.44);
  background: rgba(248, 239, 226, 0.96);
  color: #6f5330;
}

.pricing-ledger,
.pricing-table-wrap,
.pricing-group-card,
.pricing-platform-tab {
  position: relative;
  isolation: isolate;
}

.pricing-ledger {
  backdrop-filter: blur(14px);
}

.pricing-ledger::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.16), transparent 28%);
  opacity: 0.78;
}

.pricing-platform-tabs { display: flex; flex-wrap: wrap; gap: 0.45rem; }
.pricing-platform-tab { display: inline-flex; align-items: center; justify-content: center; gap: 0.38rem; border: 1px solid rgba(216, 205, 185, 0.7); border-radius: 999px; padding: 0.42rem 0.72rem; background: rgba(255, 252, 247, 0.34); color: #273027; box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.34); backdrop-filter: blur(12px); transition: border-color 180ms ease, background-color 180ms ease, transform 180ms ease, box-shadow 180ms ease, color 180ms ease; }
.pricing-platform-tab.is-active { border-color: rgba(187, 126, 59, 0.5); background: linear-gradient(135deg, rgba(255, 249, 239, 0.98), rgba(242, 228, 210, 0.96)); box-shadow: 0 14px 28px rgba(138, 104, 61, 0.14), inset 0 1px 0 rgba(255, 255, 255, 0.62); color: #2b241c; transform: translateY(-1px); }
.pricing-platform-icon { width: 1rem; display: inline-flex; justify-content: center; color: #a07a49; font-size: 0.88rem; }
.pricing-platform-label { font-size: 0.88rem; font-weight: 400; white-space: nowrap; }
.pricing-platform-tab.is-active .pricing-platform-icon,
.pricing-platform-tab.is-active .pricing-platform-label { color: #8f5f2e; }

.pricing-ledger-kicker { color: #a88a60; font-size: 0.72rem; letter-spacing: 0.22em; text-transform: uppercase; font-weight: 400; }

:global(html.dark) .pricing-ledger-kicker {
  color: #a89f8c;
}

.pricing-ledger { position: relative; }
.pricing-group-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr)); gap: 0.7rem; }
.pricing-group-grid--single { grid-template-columns: minmax(15rem, min(100%, 30rem)); }
.pricing-group-card { border: 1px solid rgba(216, 205, 185, 0.62); border-radius: 1rem; background: rgba(255, 252, 247, 0.22); padding: 0.72rem 0.92rem; cursor: pointer; box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3); backdrop-filter: blur(12px); transition: border-color 180ms ease, background-color 180ms ease, transform 180ms ease, box-shadow 180ms ease; }
.pricing-group-card.is-active { border-color: rgba(196, 136, 68, 0.38); background: linear-gradient(135deg, rgba(255, 249, 241, 0.95), rgba(243, 234, 222, 0.92)); box-shadow: 0 12px 28px rgba(136, 103, 62, 0.1), inset 0 1px 0 rgba(255, 255, 255, 0.5); transform: translateY(-1px); }
.pricing-group-head { display: flex; align-items: center; justify-content: space-between; gap: 0.8rem; }
.pricing-group-title-row { display: flex; align-items: center; gap: 0.45rem; min-width: 0; flex-wrap: nowrap; }
.pricing-group-name { color: #2b241c; font-family: 'Noto Serif SC', 'Source Han Serif SC', serif; font-size: 1rem; font-weight: 400; white-space: nowrap; }
.pricing-group-note { color: #97795a; font-size: 0.8rem; white-space: nowrap; font-weight: 400; }
.pricing-group-badge { flex-shrink: 0; border-radius: 999px; border: 1px solid rgba(190, 150, 92, 0.28); background: rgba(205, 170, 112, 0.12); color: #9a7344; padding: 0.22rem 0.5rem; font-size: 0.76rem; font-weight: 400; }

.pricing-table { border-collapse: separate; border-spacing: 0; min-width: 720px; }
.pricing-table:has(th:nth-child(5)) { min-width: 900px; }
.pricing-table th, .pricing-table td { padding: 1rem 1.1rem; border-top: 1px solid rgba(216, 205, 185, 0.38); text-align: left; vertical-align: top; }
.pricing-table thead th { border-top: 1px solid rgba(216, 205, 185, 0.38); color: #8f734c; font-size: 0.76rem; font-weight: 500; letter-spacing: 0.18em; text-transform: uppercase; }
.pricing-model-cell { min-width: 10rem; }
.pricing-model-name-row { display: inline-flex; align-items: center; gap: 0.55rem; }
.pricing-model-name-text { color: #26211b; font-size: 1rem; font-weight: 500; letter-spacing: 0.01em; }
.pricing-value-cell { min-width: 10rem; }
.pricing-price-main { color: #c98a43; font-family: 'Geist Mono', 'JetBrains Mono', monospace; font-size: 1rem; font-weight: 600; white-space: nowrap; }
.pricing-price-unit { color: #9f8662; font-size: 0.82rem; font-weight: 400; }
.pricing-empty-card { padding: 1rem 1rem 1.2rem; color: #7b6a53; font-size: 0.92rem; }

.pricing-contact-note {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.25rem;
  padding: 0.15rem 0.1rem 0;
  color: #6f6658;
  font-size: 0.9rem;
  line-height: 1.7;
}

.pricing-contact-note a,
.pricing-contact-note strong {
  color: #9c4c26;
  font-weight: 500;
  overflow-wrap: anywhere;
}

.pricing-contact-note a:hover {
  color: #a73a2a;
}

.pricing-table tbody tr {
  transition: background-color 180ms ease;
}

.pricing-table tbody tr:hover {
  background: rgba(188, 93, 31, 0.035);
}

.pricing-page :deep(.public-cta > .relative > div:last-child a:first-child) {
  background: linear-gradient(135deg, #28231c, #403227);
  color: #f7f0e4;
  box-shadow: 0 14px 26px rgba(84, 57, 31, 0.24), inset 0 1px 0 rgba(255, 241, 220, 0.12);
}

.pricing-page :deep(.public-cta > .relative > div:last-child a:first-child:hover) {
  background: linear-gradient(135deg, #3a2d23, #59432f);
  color: #fff8ee;
  box-shadow: 0 16px 30px rgba(84, 57, 31, 0.28), inset 0 1px 0 rgba(255, 241, 220, 0.16);
}

@media (max-width: 1023px) {
  .pricing-page :deep(main > section:first-child) { padding-top: 1rem; padding-bottom: 0.45rem; }
  .pricing-page :deep(main > section:nth-of-type(3)) { margin-top: 1.4rem; }
  .pricing-page :deep(.public-cta) { padding: 1.05rem; }
  .pricing-page :deep(.public-cta > .relative) { grid-template-columns: 1fr; gap: 1.15rem; }
  .pricing-page :deep(.public-cta h2) { max-width: none; }
  .pricing-page :deep(.public-cta p) { margin-top: 0.6rem; line-height: 1.72; }
  .pricing-page :deep(.public-cta > .relative > div:last-child) {
    justify-self: start;
    padding-left: 0;
    padding-top: 0.8rem;
  }
  .pricing-page :deep(.public-cta > .relative > div:last-child::before) {
    left: 0;
    right: 0;
    top: 0;
    bottom: auto;
    width: auto;
    height: 1px;
    background: linear-gradient(90deg, rgba(216, 205, 185, 0.9), transparent 72%);
  }
}

@media (max-width: 767px) {
  .pricing-platform-tabs { display: grid; grid-template-columns: 1fr; }
  .pricing-group-grid { grid-template-columns: 1fr; }
  .pricing-page :deep(.public-cta > .relative > div:last-child) { width: 100%; }
  .pricing-page :deep(.public-cta > .relative > div:last-child a) { flex: 1 1 100%; min-width: 0; }
  .pricing-page :deep(.public-cta > .relative > div:first-child > div:first-child) { letter-spacing: 0.3em; }
  .pricing-page :deep(.public-cta h2) { margin-top: 0.45rem; font-size: clamp(1.8rem, 8vw, 2.25rem); line-height: 1.14; }
  .pricing-page :deep(.public-cta p) { font-size: 0.92rem; }
  .pricing-table { min-width: 0; }
  .pricing-table thead { display: none; }
  .pricing-table,
  .pricing-table tbody,
  .pricing-table tr,
  .pricing-table td { display: block; width: 100%; }
  .pricing-table tbody { display: grid; gap: 0.85rem; padding: 0.9rem; }
  .pricing-table tbody tr {
    overflow: hidden;
    border: 1px solid rgba(216, 205, 185, 0.32);
    border-radius: 1rem;
    background: rgba(255, 252, 247, 0.08);
  }
  .pricing-table td {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.95rem;
    min-width: 0;
    padding: 0.9rem 1rem;
    border-top: 1px solid rgba(216, 205, 185, 0.22);
  }
  .pricing-table td:first-child {
    border-top: 0;
    background: rgba(255, 255, 255, 0.05);
  }
  .pricing-table td::before {
    content: attr(data-label);
    flex-shrink: 0;
    padding-top: 0.18rem;
    font-size: 0.72rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: #9b7a52;
  }
  .pricing-model-cell,
  .pricing-value-cell { min-width: 0; }
  .pricing-model-name-row,
  .pricing-price-main {
    width: 100%;
    justify-content: flex-end;
    text-align: right;
  }
  .pricing-price-main {
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-end;
    gap: 0.3rem 0.45rem;
  }
  .pricing-price-unit { white-space: nowrap; }
}
</style>

<style>
html.dark .pricing-page .public-copy-block > p:last-of-type,
html.dark .pricing-page .pricing-summary-copy,
.pricing-page.is-dark .public-copy-block > p:last-of-type,
.pricing-page.is-dark .pricing-summary-copy {
  color: #a79d89;
}

html.dark .pricing-page .public-display,
.pricing-page.is-dark .public-display {
  color: #f2e8d7;
  text-shadow: 0 1px 0 rgba(255, 244, 226, 0.06);
}

html.dark .pricing-page .public-copy-block > div:first-child span:last-child,
.pricing-page.is-dark .public-copy-block > div:first-child span:last-child {
  color: #a88456;
}

html.dark .pricing-page .public-cta,
.pricing-page.is-dark .public-cta {
  border-color: rgba(119, 105, 85, 0.56);
  background:
    linear-gradient(135deg, rgba(27, 30, 25, 0.9) 0%, rgba(33, 36, 30, 0.9) 38%, rgba(40, 31, 25, 0.88) 100%),
    radial-gradient(circle at 82% 36%, rgba(165, 84, 36, 0.16), transparent 24%),
    linear-gradient(180deg, rgba(255, 247, 235, 0.02), transparent 32%);
  box-shadow: inset 0 1px 0 rgba(255, 245, 228, 0.05), 0 24px 52px rgba(0, 0, 0, 0.16);
}

html.dark .pricing-page .public-cta-mark,
.pricing-page.is-dark .public-cta-mark {
  background:
    linear-gradient(90deg, rgba(255, 247, 235, 0.035), transparent 30%),
    radial-gradient(circle at 76% 28%, rgba(188, 93, 31, 0.18), transparent 22%),
    linear-gradient(180deg, rgba(255, 247, 235, 0.025), transparent 64%);
}

html.dark .pricing-page .public-cta h2,
.pricing-page.is-dark .public-cta h2 {
  color: #f2e8d7;
}

html.dark .pricing-page .public-cta p,
.pricing-page.is-dark .public-cta p {
  color: #a79d89;
}

html.dark .pricing-page .public-cta > .relative > div:last-child::before,
.pricing-page.is-dark .public-cta > .relative > div:last-child::before {
  background: linear-gradient(180deg, transparent, rgba(123, 105, 84, 0.82), transparent);
}

html.dark .pricing-page .public-cta > .relative > div:last-child a:first-child,
.pricing-page.is-dark .public-cta > .relative > div:last-child a:first-child {
  border: 1px solid rgba(244, 232, 214, 0.14);
  background: linear-gradient(180deg, rgba(240, 229, 210, 0.95), rgba(221, 204, 179, 0.92));
  color: #281f18;
  box-shadow: 0 18px 36px rgba(118, 54, 19, 0.26);
}

html.dark .pricing-page .public-cta > .relative > div:last-child a:first-child:hover,
.pricing-page.is-dark .public-cta > .relative > div:last-child a:first-child:hover {
  background: linear-gradient(180deg, rgba(245, 235, 217, 0.98), rgba(229, 212, 186, 0.95));
}

html.dark .pricing-page .public-cta > .relative > div:last-child a:last-child,
.pricing-page.is-dark .public-cta > .relative > div:last-child a:last-child {
  border-color: rgba(119, 102, 82, 0.86);
  background: rgba(25, 27, 21, 0.78);
  color: #ddd3c1;
}

html.dark .pricing-page .public-cta > .relative > div:last-child a:last-child:hover,
.pricing-page.is-dark .public-cta > .relative > div:last-child a:last-child:hover {
  border-color: rgba(188, 93, 31, 0.42);
  background: rgba(37, 30, 24, 0.88);
  color: #f2e8d7;
}

html.dark .pricing-page .pricing-summary-panel,
html.dark .pricing-page .pricing-ledger {
  backdrop-filter: blur(14px);
}

html.dark .pricing-page .pricing-summary-panel::before,
html.dark .pricing-page .pricing-ledger::before {
  background:
    linear-gradient(180deg, rgba(255, 247, 235, 0.04), transparent 18%),
    radial-gradient(circle at 82% 0%, rgba(188, 93, 31, 0.08), transparent 20%);
  opacity: 1;
}

html.dark .pricing-page .pricing-platform-tab {
  color: #cfc5b2;
  border-color: rgba(74, 79, 69, 0.82);
  background: linear-gradient(180deg, rgba(20, 22, 18, 0.9), rgba(14, 16, 14, 0.96));
  box-shadow: inset 0 1px 0 rgba(255, 248, 240, 0.05), inset 0 0 0 1px rgba(255, 248, 240, 0.01);
}

html.dark .pricing-page .pricing-platform-tab:hover {
  color: #f4efe4;
  border-color: rgba(188, 93, 31, 0.3);
  background: linear-gradient(180deg, rgba(42, 33, 27, 0.92), rgba(22, 25, 20, 0.96));
}

html.dark .pricing-page .pricing-platform-tab.is-active {
  color: #f4efe4;
  border-color: rgba(188, 93, 31, 0.48);
  background: linear-gradient(135deg, rgba(58, 40, 28, 0.98), rgba(22, 27, 21, 0.98));
  box-shadow: inset 0 1px 0 rgba(255, 248, 240, 0.06), 0 12px 28px rgba(0, 0, 0, 0.22);
}

html.dark .pricing-page .pricing-platform-tab.is-active .pricing-platform-icon,
html.dark .pricing-page .pricing-platform-tab.is-active .pricing-platform-label {
  color: #ffd8bb;
}

html.dark .pricing-page .pricing-group-card {
  border-color: rgba(74, 79, 69, 0.84);
  background: linear-gradient(180deg, rgba(18, 20, 17, 0.96), rgba(12, 14, 12, 0.98));
  box-shadow: inset 0 1px 0 rgba(255, 248, 240, 0.04);
}

html.dark .pricing-page .pricing-group-card:hover {
  border-color: rgba(188, 93, 31, 0.32);
  background: linear-gradient(180deg, rgba(34, 27, 22, 0.94), rgba(16, 19, 16, 0.98));
}

html.dark .pricing-page .pricing-group-card.is-active {
  border-color: rgba(188, 93, 31, 0.52);
  background: linear-gradient(135deg, rgba(56, 38, 27, 0.98), rgba(16, 19, 16, 0.98));
  box-shadow: inset 0 1px 0 rgba(255, 248, 240, 0.06), 0 18px 38px rgba(0, 0, 0, 0.26);
}

html.dark .pricing-page .pricing-group-name,
html.dark .pricing-page .pricing-model-name-text {
  color: #e9dec9;
}

html.dark .pricing-page .pricing-group-note,
html.dark .pricing-page .pricing-table thead th {
  color: #b8af9a;
}

html.dark .pricing-page .pricing-group-badge {
  border-color: rgba(188, 93, 31, 0.18);
  background: rgba(188, 93, 31, 0.16);
  color: #ffd8bb;
}

html.dark .pricing-page .pricing-group-card.is-active .pricing-group-name,
html.dark .pricing-page .pricing-group-card.is-active .pricing-group-note,
html.dark .pricing-page .pricing-group-card.is-active .pricing-group-badge {
  color: #f4efe4;
}

html.dark .pricing-page .pricing-table th,
html.dark .pricing-page .pricing-table td {
  border-top-color: rgba(70, 74, 65, 0.68);
}

html.dark .pricing-page .pricing-table thead th {
  color: #998a70;
}

html.dark .pricing-page .pricing-price-main {
  color: #efab69;
  text-shadow: 0 0 16px rgba(188, 93, 31, 0.08);
}

html.dark .pricing-page .pricing-price-unit {
  color: #928570;
}

html.dark .pricing-page .pricing-contact-note {
  color: #d0baa0;
}

html.dark .pricing-page .pricing-contact-note a,
html.dark .pricing-page .pricing-contact-note strong {
  color: #efab69;
}

html.dark .pricing-page .pricing-contact-note a:hover {
  color: #efc183;
}

html.dark .pricing-page .pricing-table tbody tr:hover {
  background: rgba(188, 93, 31, 0.05);
}

@media (max-width: 1023px) {
  html.dark .pricing-page .public-cta > .relative > div:last-child::before {
    background: linear-gradient(90deg, rgba(172, 134, 83, 0.72), transparent 72%);
  }
}

@media (max-width: 767px) {
  html.dark .pricing-page .pricing-table tbody tr {
    border-color: rgba(74, 79, 69, 0.72);
    background: rgba(17, 19, 16, 0.76);
  }

  html.dark .pricing-page .pricing-table td {
    border-top-color: rgba(70, 74, 65, 0.56);
  }

  html.dark .pricing-page .pricing-table td:first-child {
    background: rgba(255, 247, 235, 0.03);
  }

  html.dark .pricing-page .pricing-table td::before {
    color: #b8af9a;
  }
}

html.dark .pricing-page .pricing-ledger {
  border-color: rgba(155, 126, 86, 0.26) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .pricing-page .pricing-table-wrap {
  border-top-color: rgba(155, 126, 86, 0.22) !important;
  background:
    linear-gradient(180deg, rgba(255, 226, 184, 0.035), transparent 44%),
    linear-gradient(90deg, rgba(206, 151, 82, 0.025), transparent 28%, transparent 72%, rgba(206, 151, 82, 0.02)) !important;
}

html.dark .pricing-page .pricing-platform-tab,
html.dark .pricing-page .pricing-group-card {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(23, 26, 21, 0.88), rgba(14, 16, 13, 0.94)),
    radial-gradient(circle at 84% 14%, rgba(174, 102, 45, 0.08), transparent 26%) !important;
  color: #d4c4ad !important;
  box-shadow: inset 0 1px 0 rgba(255, 238, 210, 0.05) !important;
}

html.dark .pricing-page .pricing-platform-tab:hover,
html.dark .pricing-page .pricing-group-card:hover {
  border-color: rgba(194, 126, 74, 0.36) !important;
  background:
    linear-gradient(180deg, rgba(39, 32, 26, 0.94), rgba(18, 21, 17, 0.96)),
    radial-gradient(circle at 84% 14%, rgba(194, 126, 74, 0.12), transparent 26%) !important;
  color: #f3e1c7 !important;
}

html.dark .pricing-page .pricing-platform-tab.is-active,
html.dark .pricing-page .pricing-group-card.is-active {
  border-color: rgba(194, 126, 74, 0.5) !important;
  background:
    linear-gradient(135deg, rgba(58, 40, 28, 0.96), rgba(20, 24, 19, 0.96)),
    radial-gradient(circle at 84% 14%, rgba(194, 126, 74, 0.14), transparent 26%) !important;
  color: #f7ead4 !important;
  box-shadow:
    0 18px 38px rgba(0, 0, 0, 0.26),
    inset 0 1px 0 rgba(255, 238, 210, 0.06) !important;
}

html.dark .pricing-page .pricing-group-name,
html.dark .pricing-page .pricing-model-name-text {
  color: #f3e1c7 !important;
}

html.dark .pricing-page .pricing-group-note,
html.dark .pricing-page .pricing-table thead th {
  color: #d0baa0 !important;
}

html.dark .pricing-page .pricing-group-badge {
  border-color: rgba(194, 126, 74, 0.24) !important;
  background: rgba(194, 126, 74, 0.13) !important;
  color: #efc183 !important;
}

html.dark .pricing-page .pricing-price-main {
  color: #efab69 !important;
}

html.dark .pricing-page .public-cta {
  border-color: rgba(155, 126, 86, 0.28) !important;
  background:
    linear-gradient(135deg, rgba(24, 27, 22, 0.84), rgba(35, 29, 23, 0.66)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px),
    rgba(13, 15, 13, 0.5) !important;
  box-shadow:
    0 24px 54px rgba(0, 0, 0, 0.26),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html:not(.dark) .pricing-page .pricing-ledger {
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

html:not(.dark) .pricing-page .pricing-table-wrap {
  border-top-color: rgba(154, 128, 92, 0.14) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.16), transparent 42%),
    linear-gradient(90deg, rgba(144, 113, 76, 0.026), transparent 30%, transparent 70%, rgba(144, 113, 76, 0.018)) !important;
}

html:not(.dark) .pricing-page .pricing-platform-tab,
html:not(.dark) .pricing-page .pricing-group-card {
  border-color: rgba(190, 168, 134, 0.42) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.6), rgba(244, 235, 220, 0.36)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.06), transparent 26%) !important;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.42) !important;
}

html:not(.dark) .pricing-page .pricing-platform-tab:hover,
html:not(.dark) .pricing-page .pricing-group-card:hover {
  border-color: rgba(196, 136, 68, 0.32) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.82), rgba(244, 235, 220, 0.48)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.08), transparent 26%) !important;
}

html:not(.dark) .pricing-page .pricing-platform-tab.is-active,
html:not(.dark) .pricing-page .pricing-group-card.is-active {
  border-color: rgba(196, 136, 68, 0.42) !important;
  background:
    linear-gradient(135deg, rgba(255, 249, 241, 0.95), rgba(243, 234, 222, 0.84)),
    radial-gradient(circle at 84% 14%, rgba(196, 136, 68, 0.1), transparent 26%) !important;
  box-shadow: 0 12px 28px rgba(136, 103, 62, 0.1), inset 0 1px 0 rgba(255, 255, 255, 0.5) !important;
}

html:not(.dark) .pricing-page .pricing-table th,
html:not(.dark) .pricing-page .pricing-table td {
  border-top-color: rgba(190, 168, 134, 0.32) !important;
}

html:not(.dark) .pricing-page .pricing-table tbody tr:hover {
  background: rgba(188, 93, 31, 0.035) !important;
}

html:not(.dark) .pricing-page .public-cta {
  border-color: rgba(176, 150, 118, 0.44) !important;
  background:
    linear-gradient(180deg, rgba(255, 252, 246, 0.78), rgba(244, 235, 220, 0.58)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(139, 107, 68, 0.022) 33px 34px),
    rgba(255, 255, 255, 0.28) !important;
  box-shadow:
    0 14px 34px rgba(84, 57, 31, 0.06),
    inset 0 1px 0 rgba(255, 249, 239, 0.62) !important;
}
</style>
