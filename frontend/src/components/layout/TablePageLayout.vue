<template>
  <div class="table-page-layout" :class="{ 'mobile-mode': isMobile }">
    <!-- 固定区域：操作按钮 -->
    <div v-if="$slots.actions" class="layout-section-fixed">
      <slot name="actions" />
    </div>

    <!-- 固定区域：搜索和过滤器 -->
    <div v-if="$slots.filters" class="layout-section-fixed">
      <slot name="filters" />
    </div>

    <!-- 滚动区域：表格 -->
    <div class="layout-section-scrollable">
      <div class="card table-scroll-container">
        <slot name="table" />
      </div>
    </div>

    <!-- 固定区域：分页器 -->
    <div v-if="$slots.pagination" class="layout-section-fixed">
      <slot name="pagination" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const isMobile = ref(false)

const checkMobile = () => {
  isMobile.value = window.innerWidth < 1024
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
/* 桌面端：Flexbox 布局 */
.table-page-layout {
  @apply flex flex-col;
  gap: 0.9rem;
  max-width: 1520px;
  margin: 0 auto;
  height: calc(100vh - 64px - 4rem); /* 减去 header + lg:p-8 的上下padding */
}

.layout-section-fixed {
  @apply flex-shrink-0;
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 10px;
  background:
    linear-gradient(90deg, rgba(167, 58, 42, 0.028), transparent 30%),
    rgba(250, 247, 239, 0.48);
  padding: 0.9rem;
}

.layout-section-scrollable {
  @apply flex-1 min-h-0 min-w-0 flex flex-col;
}

/* 表格滚动容器 - 增强版表体滚动方案 */
.table-scroll-container {
  @apply flex h-full min-w-0 flex-col overflow-hidden;
  border: 1px solid rgba(198, 184, 157, 0.44);
  border-radius: 10px;
  background: rgba(250, 247, 239, 0.52);
  box-shadow: 0 18px 46px -38px rgba(31, 35, 32, 0.24);
}

.table-scroll-container :deep(.table-wrapper) {
  @apply flex-1 overflow-x-auto overflow-y-auto;
  /* 确保横向滚动条显示在最底部 */
  scrollbar-gutter: stable;
}

.table-scroll-container :deep(table) {
  @apply w-full;
  min-width: max-content; /* 关键：确保表格宽度根据内容撑开，从而触发横向滚动 */
  display: table; /* 使用标准 table 布局以支持 sticky 列 */
}

.table-scroll-container :deep(thead) {
  background: rgba(237, 229, 212, 0.72);
  backdrop-filter: blur(8px);
}

.table-scroll-container :deep(tbody) {
  /* 保持默认 table-row-group 显示，不使用 block */
}

.table-scroll-container :deep(th) {
  @apply px-5 py-4 text-left text-sm font-medium;
  border-bottom: 1px solid rgba(198, 184, 157, 0.54);
  color: #59645a;
}

.table-scroll-container :deep(td) {
  @apply px-5 py-4 text-sm;
  border-bottom: 1px solid rgba(198, 184, 157, 0.32);
  color: #38413a;
}

/* 移动端：恢复正常滚动 */
.table-page-layout.mobile-mode .table-scroll-container {
  @apply h-auto overflow-visible shadow-none;
}

.table-page-layout.mobile-mode .layout-section-scrollable {
  @apply flex-none min-h-fit;
}

.table-page-layout.mobile-mode .table-scroll-container :deep(.table-wrapper) {
  @apply overflow-visible;
}

.table-page-layout.mobile-mode .table-scroll-container :deep(table) {
  @apply flex-none;
  display: table;
  min-width: 100%;
}

</style>
<style>
.dark .layout-section-fixed,
.dark .table-scroll-container {
  border-color: rgba(48, 52, 43, 0.95);
  background: rgba(24, 26, 21, 0.72);
}

.dark .table-scroll-container thead {
  background: rgba(17, 19, 15, 0.62);
}

.dark .table-scroll-container th,
.dark .table-scroll-container td {
  border-color: rgba(48, 52, 43, 0.82);
}

.dark .table-scroll-container th {
  color: #879186;
}

.dark .table-scroll-container td {
  color: #d9d0be;
}
</style>
