<template>
  <div class="profile-overview-shell space-y-6">
    <section
      data-testid="profile-overview-hero"
      class="profile-overview-hero card overflow-hidden border border-zen-paperLine/70 bg-white/72 dark:border-zen-nightLine dark:bg-zen-nightPanel/72"
    >
      <div class="px-6 py-6 md:px-8">
        <div class="flex flex-col gap-6 lg:flex-row lg:items-start">
          <div
            class="flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-[1.75rem] bg-zen-seal text-2xl font-bold text-zen-paper shadow-paper-sm"
          >
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              :alt="displayName"
              class="h-full w-full object-cover"
            >
            <span v-else>{{ avatarInitial }}</span>
          </div>

          <div class="profile-overview-copy min-w-0 flex-1 space-y-5">
            <div class="space-y-3">
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="truncate text-2xl font-semibold text-gray-900 dark:text-white">
                  {{ displayName }}
                </h2>
                <span :class="['badge', user?.role === 'admin' ? 'badge-primary' : 'badge-gray']">
                  {{ user?.role === 'admin' ? t('profile.administrator') : t('profile.user') }}
                </span>
                <span
                  :class="['badge', user?.status === 'active' ? 'badge-success' : 'badge-danger']"
                >
                  {{
                    user?.status === 'active'
                      ? t('common.active')
                      : t('common.disabled')
                  }}
                </span>
              </div>

              <div class="space-y-1">
                <p class="truncate text-sm text-zen-mist dark:text-zen-stone">
                  {{ primaryEmailDisplay }}
                </p>
                <div
                  v-if="sourceHints.length"
                  class="flex flex-wrap gap-2 text-xs text-zen-mist dark:text-zen-stone"
                >
                  <span
                    v-for="hint in sourceHints"
                    :key="hint.key"
                    class="inline-flex items-center gap-1 rounded-full bg-white/80 px-3 py-1 ring-1 ring-zen-paperLine dark:bg-zen-nightPanel/70 dark:ring-zen-nightLine"
                  >
                    <Icon name="link" size="sm" />
                    {{ hint.text }}
                  </span>
                </div>
              </div>
            </div>

            <div class="profile-metrics-grid grid gap-3 sm:grid-cols-3">
              <div
                data-testid="profile-overview-metric-balance"
                class="profile-metric-card rounded-2xl bg-white/80 px-4 py-3 shadow-sm ring-1 ring-zen-paperLine/70 dark:bg-zen-nightPanel/60 dark:ring-zen-nightLine"
              >
                <p class="text-xs font-medium uppercase tracking-[0.16em] text-zen-mist dark:text-zen-stone">
                  {{ t('profile.accountBalance') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-zen-ink dark:text-zen-paper">
                  {{ formatCurrency(user?.balance || 0) }}
                </p>
              </div>
              <div
                data-testid="profile-overview-metric-concurrency"
                class="profile-metric-card rounded-2xl bg-white/80 px-4 py-3 shadow-sm ring-1 ring-zen-paperLine/70 dark:bg-zen-nightPanel/60 dark:ring-zen-nightLine"
              >
                <p class="text-xs font-medium uppercase tracking-[0.16em] text-zen-mist dark:text-zen-stone">
                  {{ t('profile.concurrencyLimit') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-zen-ink dark:text-zen-paper">
                  {{ user?.concurrency || 0 }}
                </p>
              </div>
              <div
                data-testid="profile-overview-metric-member-since"
                class="profile-metric-card rounded-2xl bg-white/80 px-4 py-3 shadow-sm ring-1 ring-zen-paperLine/70 dark:bg-zen-nightPanel/60 dark:ring-zen-nightLine"
              >
                <p class="text-xs font-medium uppercase tracking-[0.16em] text-zen-mist dark:text-zen-stone">
                  {{ t('profile.memberSince') }}
                </p>
                <p class="mt-1 text-lg font-semibold text-zen-ink dark:text-zen-paper">
                  {{ memberSinceLabel }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <div v-if="!compactOnly" class="space-y-6">
      <div data-testid="profile-main-column" class="space-y-6">
        <section
          data-testid="profile-basics-panel"
          class="card border border-gray-100 bg-white/90 p-6 dark:border-dark-700 dark:bg-dark-900/50"
        >
          <div class="mb-5 flex items-start justify-between gap-4">
            <div>
              <span class="profile-section-kicker">{{ t('profile.basicsKicker') }}</span>
              <h3 class="mt-2 text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('profile.basicsTitle') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('profile.basicsDescription') }}
              </p>
            </div>
          </div>

          <div class="grid gap-6 sm:grid-cols-1 md:grid-cols-2">
            <div class="rounded-3xl border border-zen-paperLine/70 bg-white/62 p-5 dark:border-zen-nightLine dark:bg-zen-nightPanel/40">
              <ProfileAvatarCard
                :user="user"
                embedded
              />
            </div>

            <div class="rounded-3xl border border-zen-paperLine/70 bg-white/62 p-5 dark:border-zen-nightLine dark:bg-zen-nightPanel/40">
              <ProfileEditForm
                :initial-username="user?.username || ''"
                embedded
              />
            </div>
          </div>
        </section>

        <section
          data-testid="profile-auth-bindings-panel"
          class="card border border-zen-paperLine/70 bg-white/78 p-6 dark:border-zen-nightLine dark:bg-zen-nightPanel/54"
        >
          <div class="mb-5">
            <span class="profile-section-kicker">{{ t('profile.bindingsKicker') }}</span>
            <h3 class="mt-2 text-lg font-semibold text-zen-ink dark:text-zen-paper">{{ t('profile.bindingsHeading') }}</h3>
          </div>
          <ProfileIdentityBindingsSection
            :user="user"
            :linuxdo-enabled="linuxdoEnabled"
            :dingtalk-enabled="dingtalkEnabled"
            :oidc-enabled="oidcEnabled"
            :oidc-provider-name="oidcProviderName"
            :wechat-enabled="wechatEnabled"
            :wechat-open-enabled="wechatOpenEnabled"
            :wechat-mp-enabled="wechatMpEnabled"
            embedded
            compact
          />
        </section>
      </div>

      <div data-testid="profile-side-column" class="space-y-6">
        <section
          v-if="sourceHints.length"
          class="card border border-zen-paperLine/70 bg-white/78 p-6 dark:border-zen-nightLine dark:bg-zen-nightPanel/54"
        >
          <span class="profile-section-kicker">{{ t('profile.sourceKicker') }}</span>
          <h3 class="mt-2 text-lg font-semibold text-zen-ink dark:text-zen-paper">
            {{ t('profile.linkedProfileSources') }}
          </h3>

          <div class="mt-5 grid gap-3">
            <div
              v-for="hint in sourceHints"
              :key="hint.key"
              class="flex items-start gap-3 rounded-2xl border border-zen-paperLine/70 bg-white/70 px-4 py-3 text-sm text-zen-mist dark:border-zen-nightLine dark:bg-zen-nightPanel/40 dark:text-zen-stone"
            >
              <Icon name="link" size="sm" class="mt-0.5 text-zen-mist dark:text-zen-stone" />
              <span>{{ hint.text }}</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import ProfileAvatarCard from '@/components/user/profile/ProfileAvatarCard.vue'
import ProfileEditForm from '@/components/user/profile/ProfileEditForm.vue'
import ProfileIdentityBindingsSection from '@/components/user/profile/ProfileIdentityBindingsSection.vue'
import type { User, UserAuthBindingStatus, UserAuthProvider, UserProfileSourceContext } from '@/types'

const props = withDefaults(defineProps<{
  user: User | null
  linuxdoEnabled?: boolean
  dingtalkEnabled?: boolean
  oidcEnabled?: boolean
  oidcProviderName?: string
  wechatEnabled?: boolean
  wechatOpenEnabled?: boolean
  wechatMpEnabled?: boolean
  compactOnly?: boolean
}>(), {
  linuxdoEnabled: false,
  dingtalkEnabled: false,
  oidcEnabled: false,
  oidcProviderName: 'OIDC',
  wechatEnabled: false,
  wechatOpenEnabled: undefined,
  wechatMpEnabled: undefined,
  compactOnly: false,
})

const { t } = useI18n()

function normalizeBindingStatus(binding: boolean | UserAuthBindingStatus | undefined): boolean | null {
  if (typeof binding === 'boolean') {
    return binding
  }
  if (!binding) {
    return null
  }
  if (typeof binding.bound === 'boolean') {
    return binding.bound
  }
  return Boolean(binding.provider_subject || binding.issuer || binding.provider_key)
}

function isEmailBound(user: User | null | undefined): boolean {
  if (typeof user?.email_bound === 'boolean') {
    return user.email_bound
  }

  const nested = user?.auth_bindings?.email ?? user?.identity_bindings?.email
  const normalized = normalizeBindingStatus(nested)
  return normalized ?? false
}

const avatarUrl = computed(() => props.user?.avatar_url?.trim() || '')
const displayName = computed(() => props.user?.username?.trim() || props.user?.email?.trim() || t('profile.user'))
const primaryEmailDisplay = computed(() => {
  const email = props.user?.email?.trim() || ''
  if (!email) {
    return ''
  }
  if (email.endsWith('.invalid') && !isEmailBound(props.user)) {
    return ''
  }
  return email
})
const avatarInitial = computed(() => displayName.value.charAt(0).toUpperCase() || 'U')
const memberSinceLabel = computed(() => {
  const raw = props.user?.created_at?.trim()
  if (!raw) {
    return '-'
  }

  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) {
    return '-'
  }

  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
  }).format(date)
})

const providerLabels = computed<Record<UserAuthProvider, string>>(() => ({
  email: t('profile.authBindings.providers.email'),
  linuxdo: t('profile.authBindings.providers.linuxdo'),
  dingtalk: t('profile.authBindings.providers.dingtalk'),
  oidc: t('profile.authBindings.providers.oidc', { providerName: props.oidcProviderName }),
  wechat: t('profile.authBindings.providers.wechat'),
  github: 'GitHub',
  google: 'Google'
}))

function formatCurrency(value: number): string {
  return `$${value.toFixed(2)}`
}

function normalizeProvider(value: string): UserAuthProvider | null {
  const normalized = value.trim().toLowerCase()
  if (
    normalized === 'email' ||
    normalized === 'linuxdo' ||
    normalized === 'wechat' ||
    normalized === 'github' ||
    normalized === 'google'
  ) {
    return normalized
  }
  if (normalized === 'oidc' || normalized.startsWith('oidc:') || normalized.startsWith('oidc/')) {
    return 'oidc'
  }
  return null
}

function readObjectString(source: Record<string, unknown>, ...keys: string[]): string {
  for (const key of keys) {
    const value = source[key]
    if (typeof value === 'string' && value.trim()) {
      return value.trim()
    }
  }
  return ''
}

function resolveThirdPartySource(
  rawSource: string | UserProfileSourceContext | null | undefined
): { provider: UserAuthProvider; label: string } | null {
  if (!rawSource) {
    return null
  }

  if (typeof rawSource === 'string') {
    const provider = normalizeProvider(rawSource)
    if (!provider || provider === 'email') {
      return null
    }
    return {
      provider,
      label: providerLabels.value[provider]
    }
  }

  const sourceRecord = rawSource as Record<string, unknown>
  const provider = normalizeProvider(
    readObjectString(sourceRecord, 'provider', 'source', 'provider_type', 'auth_provider')
  )
  if (!provider || provider === 'email') {
    return null
  }

  const explicitLabel = readObjectString(
    sourceRecord,
    'provider_label',
    'label',
    'provider_name',
    'providerName'
  )

  return {
    provider,
    label: explicitLabel || providerLabels.value[provider]
  }
}

const sourceHints = computed(() => {
  const currentUser = props.user
  if (!currentUser) {
    return []
  }

  const hints: Array<{ key: string; text: string }> = []
  const avatarSource = resolveThirdPartySource(
    currentUser.profile_sources?.avatar ?? currentUser.avatar_source
  )
  const usernameSource = resolveThirdPartySource(
    currentUser.profile_sources?.username ??
      currentUser.profile_sources?.display_name ??
      currentUser.profile_sources?.nickname ??
      currentUser.display_name_source ??
      currentUser.username_source ??
      currentUser.nickname_source
  )

  if (avatarSource) {
    hints.push({
      key: 'avatar',
      text: t('profile.authBindings.source.avatar', { providerName: avatarSource.label })
    })
  }

  if (usernameSource) {
    hints.push({
      key: 'username',
      text: t('profile.authBindings.source.username', { providerName: usernameSource.label })
    })
  }

  return hints
})
</script>

<style scoped>
.profile-section-kicker {
  color: #8c7a61;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.66rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.profile-overview-hero {
  border-color: rgba(198, 184, 157, 0.44);
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.84), rgba(246, 241, 231, 0.74)),
    rgba(250, 247, 239, 0.56);
}

.profile-overview-copy h2 {
  color: #1f2320;
  letter-spacing: -0.01em;
}

.profile-overview-copy p,
.profile-source-hints,
.profile-metric-card p:first-child {
  color: #6f6454;
}

.profile-metrics-grid {
  margin-top: 0.15rem;
}

.profile-metric-card {
  border: 1px solid rgba(198, 184, 157, 0.44);
  background:
    linear-gradient(180deg, rgba(255, 252, 245, 0.76), rgba(246, 241, 231, 0.64)),
    rgba(250, 247, 239, 0.66);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.36), 0 16px 36px -32px rgba(31, 35, 32, 0.22);
}

.profile-metric-card p:last-child {
  color: #1f2320;
}

.profile-source-hints {
  color: #6f6454;
}

.profile-source-hints span {
  border: 1px solid rgba(198, 184, 157, 0.58);
  background: rgba(255, 252, 245, 0.7);
  color: #6f6454;
}

.dark .profile-section-kicker {
  color: #a89a80;
}

.dark .profile-overview-hero {
  border-color: rgba(48, 52, 43, 0.95);
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.9), rgba(17, 19, 15, 0.86)),
    rgba(17, 19, 15, 0.74);
}

.dark .profile-overview-copy h2,
.dark .profile-metric-card p:last-child {
  color: #f4efe4;
}

.dark .profile-overview-copy p,
.dark .profile-source-hints,
.dark .profile-metric-card p:first-child,
.dark .profile-source-hints span {
  color: #9b8f79;
}

.dark .profile-metric-card {
  border-color: rgba(48, 52, 43, 0.92);
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.88), rgba(17, 19, 15, 0.82)),
    rgba(17, 19, 15, 0.72);
  box-shadow: inset 0 1px 0 rgba(255, 247, 235, 0.04), 0 16px 34px -34px rgba(0, 0, 0, 0.56);
}

.dark .profile-source-hints span {
  border-color: rgba(48, 52, 43, 0.9);
  background: rgba(17, 19, 15, 0.7);
}
</style>
