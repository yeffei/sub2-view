<template>
  <PublicPageLayout
    tone="faq"
    eyebrow="问答"
    title="常见问题"
    intro="这一页不再承载公开说明。"
    description="原有 FAQ 已收起；多数问题都应回到价格、开通边界与控制台实际配置中确认。"
    :highlights="['保留旧链接', '不再堆产品话术', '以实际配置为准']"
  >
    <template #aside>
      <div class="space-y-6">
        <div>
          <div class="text-xs uppercase tracking-[0.36em] text-zen-mist dark:text-zen-stone">说明</div>
          <div class="mt-3 font-serif text-3xl leading-tight text-zen-ink dark:text-zen-paper">公开问答已收起，旧链接仍保留。</div>
        </div>

        <div class="rounded-[1.2rem] border border-zen-paperLine/75 bg-white/32 px-5 py-5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/54">
          <div class="flex items-start gap-3">
            <span class="inline-flex h-10 w-10 items-center justify-center rounded-full bg-zen-seal/10 text-zen-seal dark:bg-zen-seal/16">
              <Icon name="questionCircle" size="md" />
            </span>
            <div>
              <div class="font-medium text-zen-ink dark:text-zen-paper">不再单独写 FAQ</div>
              <p class="mt-2 text-sm leading-7 text-zen-mist dark:text-zen-stone">这类页面容易累积泛化文案，信息密度低，也容易偏离当前真实配置。</p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <section class="grid gap-6 lg:grid-cols-[minmax(0,1.08fr)_minmax(18rem,0.92fr)]">
      <article class="card-glass rounded-[1.2rem] px-6 py-6 sm:px-8">
        <div class="max-w-2xl">
          <div class="text-xs uppercase tracking-[0.3em] text-zen-mist dark:text-zen-stone">当前处理</div>
          <h2 class="mt-3 font-serif text-2xl text-zen-ink dark:text-zen-paper">如果你是从旧入口进来，这里只保留一句话。</h2>
          <p class="mt-4 text-sm leading-7 text-zen-mist dark:text-zen-stone">
            公开问答已删除。是否可用、如何计量、怎么开通，应以价格页、控制台与实际接口行为为准，而不是再读一层解释性文案。
          </p>
        </div>

        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <section v-for="item in nextSteps" :key="item.title" class="rounded-[1.08rem] border border-zen-paperLine/75 bg-white/34 px-5 py-5 shadow-paper-sm dark:border-zen-nightLine dark:bg-zen-nightPanel/54">
            <div class="flex items-center gap-3 text-zen-ink dark:text-zen-paper">
              <Icon :name="item.icon" size="sm" class="text-zen-seal" />
              <h3 class="font-medium">{{ item.title }}</h3>
            </div>
            <p class="mt-3 text-sm leading-7 text-zen-mist dark:text-zen-stone">{{ item.copy }}</p>
          </section>
        </div>
      </article>

      <article class="card-glass rounded-[1.2rem] px-6 py-6">
        <h2 class="font-serif text-xl text-zen-ink dark:text-zen-paper">现在看什么</h2>
        <div class="mt-5 space-y-3">
          <RouterLink
            v-for="item in links"
            :key="item.to"
            :to="item.to"
            class="flex items-start gap-3 rounded-[1.08rem] border border-zen-paperLine/75 bg-white/32 px-4 py-4 text-left transition hover:border-zen-seal/30 hover:bg-white/52 dark:border-zen-nightLine dark:bg-zen-nightPanel/54 dark:hover:bg-zen-nightPanel/70"
          >
            <Icon :name="item.icon" size="md" class="mt-0.5 text-zen-seal" />
            <div>
              <div class="font-medium text-zen-ink dark:text-zen-paper">{{ item.title }}</div>
              <p class="mt-2 text-sm leading-7 text-zen-mist dark:text-zen-stone">{{ item.copy }}</p>
            </div>
          </RouterLink>
        </div>
      </article>
    </section>
  </PublicPageLayout>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'
import PublicPageLayout from '@/components/layout/PublicPageLayout.vue'

const nextSteps = [
  {
    title: '先看价格',
    copy: '公开页里真正长期有效的，是开通方式、计量口径和边界说明。',
    icon: 'dollar',
  },
  {
    title: '再进控制台',
    copy: '涉及 Key、余额、订单、用量与实际能力时，以登录后的真实配置为准。',
    icon: 'grid',
  },
] as const

const links = [
  {
    to: '/pricing',
    title: '价格与开通',
    copy: '确认开通方式、计量边界和对外保留的信息。',
    icon: 'dollar',
  },
  {
    to: '/login',
    title: '进入控制台',
    copy: '如需继续判断可用性或配置细节，直接进入实际系统。',
    icon: 'login',
  },
] as const
</script>

<style>
html.dark .sst-public.sst-public-tone-faq .card-glass,
html.dark .sst-public.sst-public-tone-faq .rounded-\[1\.08rem\],
html.dark .sst-public.sst-public-tone-faq .rounded-\[1\.2rem\] {
  border-color: rgba(114, 103, 88, 0.54);
  background:
    linear-gradient(180deg, rgba(22, 24, 20, 0.9), rgba(31, 27, 23, 0.82));
  box-shadow:
    inset 0 1px 0 rgba(255, 244, 224, 0.05),
    0 18px 42px rgba(0, 0, 0, 0.2);
}

html.dark .sst-public.sst-public-tone-faq .text-zen-ink,
html.dark .sst-public.sst-public-tone-faq h2,
html.dark .sst-public.sst-public-tone-faq h3 {
  color: #efe3d1;
}

html.dark .sst-public.sst-public-tone-faq .text-zen-mist,
html.dark .sst-public.sst-public-tone-faq p {
  color: #bcae98;
}

html.dark .sst-public.sst-public-tone-faq .hover\:bg-white\/52:hover,
html.dark .sst-public.sst-public-tone-faq .dark\:hover\:bg-zen-nightPanel\/70:hover {
  background: linear-gradient(180deg, rgba(44, 34, 28, 0.9), rgba(27, 30, 24, 0.9));
}

html.dark .sst-public.sst-public-tone-faq .card-glass,
html.dark .sst-public.sst-public-tone-faq .rounded-\[1\.08rem\],
html.dark .sst-public.sst-public-tone-faq .rounded-\[1\.2rem\] {
  border-color: rgba(155, 126, 86, 0.24) !important;
  background:
    linear-gradient(180deg, rgba(24, 27, 22, 0.88), rgba(34, 29, 23, 0.78)),
    repeating-linear-gradient(0deg, transparent 0 33px, rgba(230, 194, 142, 0.025) 33px 34px) !important;
  box-shadow:
    0 22px 48px rgba(0, 0, 0, 0.22),
    inset 0 1px 0 rgba(245, 225, 194, 0.055) !important;
}

html.dark .sst-public.sst-public-tone-faq .text-zen-ink,
html.dark .sst-public.sst-public-tone-faq h2,
html.dark .sst-public.sst-public-tone-faq h3 {
  color: #f6e8d2 !important;
}

html.dark .sst-public.sst-public-tone-faq .text-zen-mist,
html.dark .sst-public.sst-public-tone-faq p {
  color: #d0baa0 !important;
}

html:not(.dark) .sst-public.sst-public-tone-faq .card-glass,
html:not(.dark) .sst-public.sst-public-tone-faq .rounded-\[1\.08rem\],
html:not(.dark) .sst-public.sst-public-tone-faq .rounded-\[1\.2rem\] {
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

html:not(.dark) .sst-public.sst-public-tone-faq .text-zen-ink,
html:not(.dark) .sst-public.sst-public-tone-faq h2,
html:not(.dark) .sst-public.sst-public-tone-faq h3 {
  color: #2f281f !important;
}

html:not(.dark) .sst-public.sst-public-tone-faq .text-zen-mist,
html:not(.dark) .sst-public.sst-public-tone-faq p {
  color: #5f685c !important;
}
</style>
