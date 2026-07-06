<template>
  <AppLayout>
    <div class="sst-admin-page">
      <section class="affiliate-tabs" aria-label="团队引荐记录分类">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          type="button"
          class="affiliate-tab"
          :class="{ 'affiliate-tab-active': activeTab === tab.value }"
          @click="setActiveTab(tab.value)"
        >
          <Icon :name="tab.icon" size="sm" />
          <span>{{ tab.label }}</span>
        </button>
      </section>

      <AdminAffiliateRecordsTable :type="activeTab" embedded />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import AdminAffiliateRecordsTable from './AdminAffiliateRecordsTable.vue'

type AffiliateTab = 'invites' | 'rebates' | 'transfers'
type IconName = InstanceType<typeof Icon>['$props']['name']

const route = useRoute()
const router = useRouter()

const tabs: Array<{ value: AffiliateTab; label: string; icon: IconName }> = [
  { value: 'invites', label: '邀请记录', icon: 'users' },
  { value: 'rebates', label: '返利记录', icon: 'document' },
  { value: 'transfers', label: '转入记录', icon: 'creditCard' },
]

const activeTab = computed<AffiliateTab>(() => {
  const tab = route.query.tab
  return tab === 'rebates' || tab === 'transfers' ? tab : 'invites'
})

function setActiveTab(tab: AffiliateTab): void {
  if (tab === activeTab.value) return
  router.replace({ path: '/admin/affiliates', query: { ...route.query, tab } })
}
</script>

<style scoped>
.affiliate-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
  padding: 0.35rem;
  border: 1px solid rgba(198, 184, 157, 0.5);
  border-radius: 0.9rem;
  background: rgba(250, 247, 239, 0.62);
}

.affiliate-tab {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  min-height: 2.35rem;
  padding: 0.48rem 0.85rem;
  border-radius: 0.65rem;
  color: #5d655b;
  font-size: 0.875rem;
  font-weight: 500;
  transition:
    background-color 0.18s ease,
    color 0.18s ease,
    box-shadow 0.18s ease;
}

.affiliate-tab:hover {
  background: rgba(167, 58, 42, 0.08);
}

.affiliate-tab-active {
  background: #a73a2a;
  color: #fff7f2;
  box-shadow: 0 10px 22px -18px rgba(167, 58, 42, 0.5);
}
</style>
<style>
.dark .affiliate-tabs {
  border-color: rgba(58, 61, 54, 0.96);
  background: rgba(16, 18, 14, 0.82);
}

.dark .affiliate-tab {
  color: #cfc6b4;
}

.dark .affiliate-tab:hover {
  background: rgba(167, 58, 42, 0.16);
}

.dark .affiliate-tab-active {
  color: #fff7f2;
}
</style>
