<template>
  <div class="auth-artwork-hero" :class="frameClass">
    <div class="auth-artwork-hero__frame">
      <img
        :src="src"
        alt=""
        class="auth-artwork-hero__image"
        :style="{ objectPosition }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  src: string
  objectPosition?: string
  variant?: 'register' | 'login'
}>(), {
  objectPosition: 'center center',
  variant: 'register'
})

const frameClass = computed(() => `auth-artwork-hero--${props.variant}`)
</script>

<style scoped>
.auth-artwork-hero {
  position: relative;
  width: min(100%, 28.75rem);
  margin: 1.25rem auto 0 0;
  padding: 0.95rem 1.15rem 0.35rem 0.95rem;
}

.auth-artwork-hero::before {
  content: '';
  position: absolute;
  inset: 2.65rem 0 0.5rem 2.05rem;
  border: 1px solid rgba(182, 164, 136, 0.14);
  border-radius: 1.65rem 0.85rem 1.65rem 0.85rem;
  background: linear-gradient(180deg, rgba(255, 252, 247, 0.46), rgba(242, 233, 220, 0.14));
  box-shadow: 0 18px 38px -34px rgba(53, 44, 34, 0.22);
  pointer-events: none;
}

.auth-artwork-hero__frame {
  position: relative;
  z-index: 1;
  overflow: hidden;
  padding: 0.68rem;
  border: 1px solid rgba(176, 158, 131, 0.28);
  border-radius: 1.75rem 0.8rem 1.75rem 0.8rem;
  background:
    linear-gradient(180deg, rgba(255, 253, 249, 0.9), rgba(245, 237, 225, 0.76)),
    rgba(247, 240, 228, 0.82);
  box-shadow:
    0 24px 64px -44px rgba(42, 34, 24, 0.34),
    0 1px 0 rgba(255, 255, 255, 0.58) inset;
}

.auth-artwork-hero__frame::before {
  content: '';
  position: absolute;
  inset: 0.68rem;
  z-index: 1;
  border: 1px solid rgba(255, 249, 240, 0.44);
  border-radius: 1.28rem 0.52rem 1.28rem 0.52rem;
  pointer-events: none;
}

.auth-artwork-hero__frame::after {
  content: '';
  position: absolute;
  inset: 0;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.16), transparent 20%),
    radial-gradient(circle at 84% 18%, rgba(168, 69, 53, 0.1), transparent 18%);
  pointer-events: none;
}

.auth-artwork-hero__image {
  display: block;
  width: 100%;
  aspect-ratio: 5 / 6.35;
  border-radius: 1.15rem 0.42rem 1.15rem 0.42rem;
  filter: saturate(0.94) contrast(0.98) brightness(1.01);
  object-fit: cover;
}

.auth-artwork-hero--login .auth-artwork-hero__image {
  aspect-ratio: 5 / 6.45;
}

:global(html.dark) .auth-artwork-hero::before {
  border-color: rgba(191, 161, 118, 0.12);
  background: linear-gradient(180deg, rgba(29, 26, 22, 0.44), rgba(16, 14, 12, 0.18));
  box-shadow: 0 24px 54px -44px rgba(0, 0, 0, 0.66);
}

:global(html.dark) .auth-artwork-hero__frame {
  border-color: rgba(191, 161, 118, 0.24);
  background:
    linear-gradient(180deg, rgba(28, 25, 21, 0.86), rgba(20, 18, 16, 0.92)),
    rgba(22, 18, 16, 0.82);
  box-shadow:
    inset 0 1px 0 rgba(248, 231, 200, 0.05),
    0 28px 84px -46px rgba(0, 0, 0, 0.58);
}

:global(html.dark) .auth-artwork-hero__frame::before {
  border-color: rgba(248, 231, 200, 0.08);
}

:global(html.dark) .auth-artwork-hero__frame::after {
  background:
    linear-gradient(180deg, rgba(255, 248, 235, 0.05), transparent 22%),
    radial-gradient(circle at 84% 18%, rgba(194, 129, 61, 0.12), transparent 18%);
}

@media (max-width: 1279px) {
  .auth-artwork-hero {
    width: min(100%, 27.5rem);
    margin-top: 0.5rem;
    padding-right: 0.9rem;
  }
}
</style>
