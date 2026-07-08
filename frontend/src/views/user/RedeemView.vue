<template>
  <AppLayout>
    <div class="redeem-redirect-shell mx-auto max-w-3xl">
      <div class="redeem-redirect-card">
        <div class="redeem-redirect-icon">
          <Icon name="wallet" size="lg" />
        </div>
        <div class="redeem-redirect-copy">
          <span>{{ t('redeem.redirect.kicker') }}</span>
          <strong>{{ t('redeem.redirect.title') }}</strong>
          <p>{{ t('redeem.redirect.copy') }}</p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'

const router = useRouter()
const { t } = useI18n()

onMounted(() => {
  void router.replace({ path: '/purchase', query: { redeem: '1' } })
})
</script>

<style scoped>
.redeem-redirect-shell {
  padding-top: 1rem;
}

.redeem-redirect-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  border: 1px solid rgba(216, 205, 185, 0.72);
  border-radius: 0.9rem;
  background: rgba(250, 247, 239, 0.58);
  box-shadow: 0 20px 48px -38px rgba(31, 35, 32, 0.26);
  padding: 1.2rem 1.25rem;
}

.redeem-redirect-copy {
  min-width: 0;
}

.redeem-redirect-icon {
  display: flex;
  height: 3rem;
  width: 3rem;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(167, 58, 42, 0.08);
  color: #a73a2a;
}

.redeem-redirect-card span {
  display: block;
  color: #7b6a53;
  font-size: 0.68rem;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

.redeem-redirect-card strong {
  display: block;
  margin-top: 0.3rem;
  color: #1f2320;
  font-family: 'Noto Serif SC', 'Source Han Serif SC', serif;
  font-size: 1.2rem;
  font-weight: 600;
}

.redeem-redirect-card p {
  margin-top: 0.45rem;
  color: #59645a;
  font-size: 0.92rem;
  line-height: 1.75;
}

.dark .redeem-redirect-card {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .redeem-redirect-card strong {
  color: #f4efe4;
}

.dark .redeem-redirect-card p {
  color: #a6afa4;
}

.dark .redeem-redirect-card span {
  color: #879186;
}

@media (max-width: 640px) {
  .redeem-redirect-card {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
