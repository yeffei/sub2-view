<template>
  <div>
    <!-- 铃铛按钮 -->
    <button
      v-if="!props.buttonless"
      @click="openModal"
      class="announcement-bell-trigger relative flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition-all hover:bg-gray-100 hover:scale-105 dark:text-gray-400 dark:hover:bg-dark-800"
      :class="{ 'text-blue-600 dark:text-blue-400': unreadCount > 0 }"
      :aria-label="t('announcements.title')"
    >
      <Icon name="bell" size="md" />
      <!-- 未读红点 -->
      <span
        v-if="unreadCount > 0"
        class="absolute right-1 top-1 flex h-2 w-2"
      >
        <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-red-500 opacity-75"></span>
        <span class="relative inline-flex h-2 w-2 rounded-full bg-red-500"></span>
      </span>
    </button>

    <!-- 公告列表 Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="isModalOpen"
          class="announcement-overlay fixed inset-0 z-[100] flex items-start justify-center overflow-y-auto p-4 pt-[8vh] backdrop-blur-md"
          @click="closeModal"
        >
          <div
            class="announcement-modal announcement-modal-list w-full max-w-[620px] overflow-hidden rounded-[1.4rem]"
            @click.stop
          >
            <div class="announcement-header relative overflow-hidden px-6 py-5">
              <div class="relative z-10 flex items-start justify-between">
                <div>
                  <div class="flex items-center gap-2">
                    <div class="announcement-header-icon flex h-8 w-8 items-center justify-center rounded-lg">
                      <Icon name="bell" size="sm" />
                    </div>
                    <h2 class="announcement-title text-lg font-semibold">
                      {{ t('announcements.title') }}
                    </h2>
                  </div>
                  <p v-if="unreadCount > 0" class="announcement-subtitle mt-2 text-sm">
                    <span class="announcement-subtitle-count font-medium">{{ unreadCount }}</span>
                    {{ t('announcements.unread') }}
                  </p>
                </div>
                <div class="flex items-center gap-2">
                  <button
                    v-if="unreadCount > 0"
                    @click="markAllAsRead"
                    :disabled="loading"
                    class="announcement-primary-action rounded-lg px-4 py-2 text-xs font-medium transition-all disabled:opacity-50"
                  >
                    {{ t('announcements.markAllRead') }}
                  </button>
                  <button
                    @click="closeModal"
                    class="announcement-close flex h-9 w-9 items-center justify-center rounded-lg transition-all"
                    :aria-label="t('common.close')"
                  >
                    <Icon name="x" size="sm" />
                  </button>
                </div>
              </div>
              <div class="announcement-header-glow absolute right-0 top-0 h-full w-48"></div>
            </div>

            <!-- Body -->
            <div class="max-h-[65vh] overflow-y-auto">
              <!-- Loading -->
              <div v-if="loading" class="flex items-center justify-center py-16">
                <div class="relative">
                  <div class="announcement-spinner h-12 w-12 animate-spin rounded-full border-4"></div>
                  <div class="announcement-spinner-glow absolute inset-0 h-12 w-12 animate-pulse rounded-full border-4"></div>
                </div>
              </div>

              <!-- Announcements List -->
              <div v-else-if="announcements.length > 0">
                <div
                  v-for="item in announcements"
                  :key="item.id"
                  class="announcement-row group relative flex items-center gap-4 px-6 py-4 transition-all"
                  :class="{ 'is-unread': !item.read_at }"
                  style="min-height: 72px"
                  @click="openDetail(item)"
                >
                  <!-- Status Indicator -->
                  <div class="flex h-10 w-10 flex-shrink-0 items-center justify-center">
                    <div
                      v-if="!item.read_at"
                      class="announcement-status-unread relative flex h-10 w-10 items-center justify-center rounded-xl"
                    >
                      <span class="announcement-status-ping absolute inline-flex h-full w-full animate-ping rounded-xl"></span>
                      <svg class="relative z-10 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <div
                      v-else
                      class="announcement-status-read flex h-10 w-10 items-center justify-center rounded-xl"
                    >
                      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                  </div>

                  <!-- Content -->
                  <div class="flex min-w-0 flex-1 items-center justify-between gap-4">
                    <div class="min-w-0 flex-1">
                      <h3 class="announcement-row-title truncate text-sm font-medium">
                        {{ item.title }}
                      </h3>
                      <div class="mt-1 flex items-center gap-2">
                        <time class="announcement-row-time text-xs">
                          {{ formatRelativeTime(item.created_at) }}
                        </time>
                        <span
                          v-if="!item.read_at"
                          class="announcement-row-badge inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 text-xs font-medium"
                        >
                          <span class="relative flex h-1.5 w-1.5">
                            <span class="announcement-row-badge-ping absolute inline-flex h-full w-full animate-ping rounded-full"></span>
                            <span class="announcement-row-badge-dot relative inline-flex h-1.5 w-1.5 rounded-full"></span>
                          </span>
                          {{ t('announcements.unread') }}
                        </span>
                      </div>
                    </div>

                    <div class="flex-shrink-0">
                      <svg
                        class="announcement-row-arrow h-5 w-5 transition-transform group-hover:translate-x-1"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                      </svg>
                    </div>
                  </div>

                  <!-- Unread indicator bar -->
                  <div
                    v-if="!item.read_at"
                    class="announcement-row-bar absolute left-0 top-0 h-full w-1"
                  ></div>
                </div>
              </div>

              <!-- Empty State -->
              <div v-else class="announcement-empty flex flex-col items-center justify-center py-16">
                <div class="relative mb-4">
                  <div class="announcement-empty-icon flex h-20 w-20 items-center justify-center rounded-full">
                    <Icon name="inbox" size="xl" class="announcement-empty-icon-mark" />
                  </div>
                  <div class="announcement-empty-ok absolute -right-1 -top-1 flex h-6 w-6 items-center justify-center rounded-full text-white">
                    <svg class="h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <p class="announcement-empty-title text-sm font-medium">{{ t('announcements.empty') }}</p>
                <p class="announcement-empty-copy mt-1 text-xs">{{ t('announcements.emptyDescription') }}</p>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- 公告详情 Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div
          v-if="detailModalOpen && selectedAnnouncement"
          class="announcement-overlay fixed inset-0 z-[110] flex items-start justify-center overflow-y-auto p-4 pt-[6vh] backdrop-blur-md"
          @click="closeDetail"
        >
          <div
            class="announcement-modal announcement-modal-detail w-full max-w-[780px] overflow-hidden rounded-[1.5rem]"
            @click.stop
          >
            <div class="announcement-header announcement-detail-header relative overflow-hidden px-8 py-6">
              <div class="relative z-10 flex items-start justify-between gap-4">
                <div class="flex-1 min-w-0">
                  <div class="mb-3 flex items-center gap-2">
                    <div class="announcement-header-icon flex h-10 w-10 items-center justify-center rounded-xl">
                      <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <div class="flex items-center gap-2">
                      <span class="announcement-detail-tag rounded-lg px-2.5 py-1 text-xs font-medium">
                        {{ t('announcements.title') }}
                      </span>
                      <span
                        v-if="!selectedAnnouncement.read_at"
                        class="announcement-detail-badge inline-flex items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs font-medium"
                      >
                        <span class="relative flex h-2 w-2">
                          <span class="announcement-detail-badge-ping absolute inline-flex h-full w-full animate-ping rounded-full"></span>
                          <span class="announcement-detail-badge-dot relative inline-flex h-2 w-2 rounded-full"></span>
                        </span>
                        {{ t('announcements.unread') }}
                      </span>
                    </div>
                  </div>

                  <h2 class="announcement-detail-title mb-3 text-2xl font-bold leading-tight">
                    {{ selectedAnnouncement.title }}
                  </h2>

                  <div class="announcement-detail-meta flex items-center gap-4 text-sm">
                    <div class="flex items-center gap-1.5">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <time>{{ formatRelativeWithDateTime(selectedAnnouncement.created_at) }}</time>
                    </div>
                    <div class="flex items-center gap-1.5">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      <span>{{ selectedAnnouncement.read_at ? t('announcements.read') : t('announcements.unread') }}</span>
                    </div>
                  </div>
                </div>

                <!-- Close button -->
                <button
                  @click="closeDetail"
                  class="announcement-close flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-xl transition-all"
                  :aria-label="t('common.close')"
                >
                  <Icon name="x" size="md" />
                </button>
              </div>
            </div>

            <div class="announcement-body max-h-[60vh] overflow-y-auto px-8 py-8">
              <div class="relative">
                <div class="announcement-body-bar absolute left-0 top-0 bottom-0 w-1 rounded-full"></div>

                <div class="pl-6">
                  <div
                    class="markdown-body prose prose-sm max-w-none dark:prose-invert"
                    v-html="renderMarkdown(selectedAnnouncement.content)"
                  ></div>
                </div>
              </div>
            </div>

            <div class="announcement-footer px-8 py-5">
              <div class="flex items-center justify-between">
                <div class="announcement-footer-note flex items-center gap-2 text-xs">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span>{{ selectedAnnouncement.read_at ? t('announcements.readStatus') : t('announcements.markReadHint') }}</span>
                </div>
                <div class="flex items-center gap-3">
                  <button
                    @click="closeDetail"
                    class="announcement-secondary-action rounded-xl px-5 py-2.5 text-sm font-medium transition-all"
                  >
                    {{ t('common.close') }}
                  </button>
                  <button
                    v-if="!selectedAnnouncement.read_at"
                    @click="markAsReadAndClose(selectedAnnouncement.id)"
                    class="announcement-primary-action rounded-xl px-5 py-2.5 text-sm font-medium transition-all"
                  >
                    <span class="flex items-center gap-2">
                      <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                      {{ t('announcements.markRead') }}
                    </span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores/app'
import { useAnnouncementStore } from '@/stores/announcements'
import { formatRelativeTime, formatRelativeWithDateTime } from '@/utils/format'
import type { UserAnnouncement } from '@/types'
import Icon from '@/components/icons/Icon.vue'

const props = withDefaults(defineProps<{
  buttonless?: boolean
}>(), {
  buttonless: false,
})

const { t } = useI18n()
const appStore = useAppStore()
const announcementStore = useAnnouncementStore()
const OPEN_ANNOUNCEMENT_CENTER_EVENT = 'sst-open-announcement-center'

// Configure marked
marked.setOptions({
  breaks: true,
  gfm: true,
})

// Use store state (storeToRefs for reactivity)
const { announcements, loading } = storeToRefs(announcementStore)
const unreadCount = computed(() => announcementStore.unreadCount)

// Local modal state
const isModalOpen = ref(false)
const detailModalOpen = ref(false)
const selectedAnnouncement = ref<UserAnnouncement | null>(null)

// Methods
function renderMarkdown(content: string): string {
  if (!content) return ''
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
}

function openModal() {
  isModalOpen.value = true
}

function closeModal() {
  isModalOpen.value = false
}

function openDetail(announcement: UserAnnouncement) {
  selectedAnnouncement.value = announcement
  detailModalOpen.value = true
  if (!announcement.read_at) {
    markAsRead(announcement.id)
  }
}

function closeDetail() {
  detailModalOpen.value = false
  selectedAnnouncement.value = null
}

async function markAsRead(id: number) {
  try {
    await announcementStore.markAsRead(id)
  } catch (err: any) {
    appStore.showError(err?.message || t('common.unknownError'))
  }
}

async function markAsReadAndClose(id: number) {
  await markAsRead(id)
  appStore.showSuccess(t('announcements.markedAsRead'))
  closeDetail()
}

async function markAllAsRead() {
  try {
    await announcementStore.markAllAsRead()
    appStore.showSuccess(t('announcements.allMarkedAsRead'))
  } catch (err: any) {
    appStore.showError(err?.message || t('common.unknownError'))
  }
}

function handleEscape(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (detailModalOpen.value) {
      closeDetail()
    } else if (isModalOpen.value) {
      closeModal()
    }
  }
}

function handleOpenAnnouncementCenter() {
  openModal()
}

onMounted(() => {
  document.addEventListener('keydown', handleEscape)
  window.addEventListener(OPEN_ANNOUNCEMENT_CENTER_EVENT, handleOpenAnnouncementCenter)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleEscape)
  window.removeEventListener(OPEN_ANNOUNCEMENT_CENTER_EVENT, handleOpenAnnouncementCenter)
  document.body.style.overflow = ''
})

watch(
  [isModalOpen, detailModalOpen, () => announcementStore.currentPopup],
  ([modal, detail, popup]) => {
    document.body.style.overflow = (modal || detail || popup) ? 'hidden' : ''
  }
)
</script>

<style scoped>
/* Modal Animations */
.modal-fade-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-fade-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 1, 1);
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-from > div {
  transform: scale(0.94) translateY(-12px);
  opacity: 0;
}

.modal-fade-leave-to > div {
  transform: scale(0.96) translateY(-8px);
  opacity: 0;
}

/* Scrollbar Styling */
.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #cbd5e1, #94a3b8);
  border-radius: 4px;
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb {
  background: linear-gradient(to bottom, #4b5563, #374151);
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #94a3b8, #64748b);
}

.dark .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(to bottom, #6b7280, #4b5563);
}

.announcement-overlay {
  background:
    linear-gradient(180deg, rgba(42, 33, 23, 0.22), rgba(42, 33, 23, 0.3)),
    rgba(90, 73, 53, 0.14);
}

.announcement-modal {
  border: 1px solid rgba(171, 145, 108, 0.2);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.68), rgba(248, 243, 233, 0.96)),
    #f7f1e6;
  box-shadow:
    0 28px 60px -34px rgba(82, 62, 39, 0.28),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
  color: #2a241e;
}

.announcement-header {
  border-bottom: 1px solid rgba(188, 166, 133, 0.24);
  background:
    linear-gradient(135deg, rgba(247, 239, 226, 0.92), rgba(242, 233, 217, 0.82));
}

.announcement-header-glow {
  background: linear-gradient(270deg, rgba(183, 107, 50, 0.08), transparent);
}

.announcement-header-icon {
  background: linear-gradient(135deg, #b66c35, #8f4b28);
  color: #fff8f1;
  box-shadow: 0 14px 28px -18px rgba(143, 75, 40, 0.5);
}

.announcement-title,
.announcement-detail-title {
  color: #2a241e;
}

.announcement-subtitle,
.announcement-detail-meta,
.announcement-footer-note {
  color: #746858;
}

.announcement-subtitle-count {
  color: #8f4b28;
}

.announcement-primary-action {
  border: 1px solid rgba(143, 75, 40, 0.16);
  background: linear-gradient(135deg, #b66c35, #944c2a);
  color: #fffaf4;
  box-shadow: 0 14px 28px -20px rgba(143, 75, 40, 0.48);
}

.announcement-primary-action:hover {
  filter: brightness(1.04);
}

.announcement-secondary-action,
.announcement-close {
  border: 1px solid rgba(188, 166, 133, 0.28);
  background: rgba(255, 251, 244, 0.74);
  color: #746858;
}

.announcement-secondary-action:hover,
.announcement-close:hover {
  background: rgba(255, 255, 255, 0.92);
  color: #2a241e;
}

.announcement-spinner {
  border-color: rgba(204, 190, 167, 0.84);
  border-top-color: #a55832;
}

.announcement-spinner-glow {
  border-color: rgba(182, 108, 53, 0.18);
}

.announcement-row {
  border-bottom: 1px solid rgba(198, 184, 157, 0.2);
}

.announcement-row:hover {
  background: rgba(191, 146, 83, 0.07);
}

.announcement-row.is-unread {
  background: rgba(167, 58, 42, 0.03);
}

.announcement-row-title {
  color: #2a241e;
}

.announcement-row-time,
.announcement-empty-copy {
  color: #746858;
}

.announcement-status-unread {
  background: linear-gradient(135deg, #b66c35, #944c2a);
  color: #fff8f1;
  box-shadow: 0 14px 28px -20px rgba(143, 75, 40, 0.48);
}

.announcement-status-ping {
  background: rgba(182, 108, 53, 0.28);
}

.announcement-status-read {
  background: rgba(236, 229, 216, 0.84);
  color: #9a8b75;
}

.announcement-row-badge {
  background: rgba(182, 108, 53, 0.1);
  color: #8f4b28;
}

.announcement-row-badge-ping,
.announcement-row-badge-dot,
.announcement-detail-badge-ping,
.announcement-detail-badge-dot {
  background: #a55832;
}

.announcement-row-arrow {
  color: #a69780;
}

.announcement-row-bar,
.announcement-body-bar {
  background: linear-gradient(180deg, #b66c35, #8f4b28);
}

.announcement-empty-icon {
  background: linear-gradient(135deg, rgba(241, 231, 214, 0.92), rgba(232, 220, 199, 0.86));
}

.announcement-empty-icon-mark {
  color: #ad9a81;
}

.announcement-empty-ok {
  background: #8f4b28;
}

.announcement-empty-title {
  color: #2a241e;
}

.announcement-detail-tag {
  background: rgba(182, 108, 53, 0.1);
  color: #8f4b28;
}

.announcement-detail-badge {
  background: linear-gradient(135deg, #b66c35, #944c2a);
  color: #fff8f1;
}

.announcement-body {
  background: rgba(251, 247, 240, 0.94);
}

.announcement-footer {
  border-top: 1px solid rgba(198, 184, 157, 0.2);
  background: rgba(247, 241, 231, 0.9);
}

.dark .announcement-overlay {
  background: linear-gradient(to bottom right, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.6), rgba(0, 0, 0, 0.7));
}

.dark .announcement-modal {
  border-color: rgba(255, 255, 255, 0.08);
  background: #1f2328;
  box-shadow: 0 24px 56px -28px rgba(0, 0, 0, 0.6);
  color: #fff;
}

.dark .announcement-header {
  border-bottom-color: rgba(255, 255, 255, 0.08);
  background: linear-gradient(135deg, rgba(37, 99, 235, 0.12), rgba(79, 70, 229, 0.08));
}

.dark .announcement-header-glow {
  background: linear-gradient(270deg, rgba(99, 102, 241, 0.12), transparent);
}

.dark .announcement-header-icon {
  background: linear-gradient(135deg, #3b82f6, #4f46e5);
  color: #fff;
  box-shadow: 0 14px 28px -18px rgba(59, 130, 246, 0.36);
}

.dark .announcement-title,
.dark .announcement-detail-title,
.dark .announcement-row-title,
.dark .announcement-empty-title {
  color: #fff;
}

.dark .announcement-subtitle,
.dark .announcement-detail-meta,
.dark .announcement-footer-note,
.dark .announcement-row-time,
.dark .announcement-empty-copy {
  color: #9ca3af;
}

.dark .announcement-subtitle-count {
  color: #60a5fa;
}

.dark .announcement-primary-action {
  border-color: transparent;
  background: linear-gradient(135deg, #2563eb, #4f46e5);
  color: #fff;
  box-shadow: 0 14px 28px -20px rgba(59, 130, 246, 0.4);
}

.dark .announcement-secondary-action,
.dark .announcement-close {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: #9ca3af;
}

.dark .announcement-secondary-action:hover,
.dark .announcement-close:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.dark .announcement-spinner {
  border-color: #4b5563;
  border-top-color: #60a5fa;
}

.dark .announcement-spinner-glow {
  border-color: rgba(96, 165, 250, 0.2);
}

.dark .announcement-row {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}

.dark .announcement-row:hover {
  background: rgba(255, 255, 255, 0.03);
}

.dark .announcement-row.is-unread {
  background: rgba(37, 99, 235, 0.08);
}

.dark .announcement-status-unread {
  background: linear-gradient(135deg, #3b82f6, #4f46e5);
  color: #fff;
  box-shadow: 0 14px 28px -20px rgba(59, 130, 246, 0.36);
}

.dark .announcement-status-ping {
  background: rgba(96, 165, 250, 0.3);
}

.dark .announcement-status-read {
  background: rgba(255, 255, 255, 0.06);
  color: #6b7280;
}

.dark .announcement-row-badge {
  background: rgba(37, 99, 235, 0.18);
  color: #93c5fd;
}

.dark .announcement-row-badge-ping,
.dark .announcement-row-badge-dot,
.dark .announcement-detail-badge-ping,
.dark .announcement-detail-badge-dot {
  background: #60a5fa;
}

.dark .announcement-row-arrow {
  color: #6b7280;
}

.dark .announcement-row-bar,
.dark .announcement-body-bar {
  background: linear-gradient(180deg, #3b82f6, #4f46e5);
}

.dark .announcement-empty-icon {
  background: linear-gradient(135deg, rgba(55, 65, 81, 0.9), rgba(31, 41, 55, 0.9));
}

.dark .announcement-empty-icon-mark {
  color: #6b7280;
}

.dark .announcement-empty-ok {
  background: #16a34a;
}

.dark .announcement-detail-tag {
  background: rgba(37, 99, 235, 0.18);
  color: #93c5fd;
}

.dark .announcement-detail-badge {
  background: linear-gradient(135deg, #3b82f6, #4f46e5);
  color: #fff;
}

.dark .announcement-body {
  background: #1f2328;
}

.dark .announcement-footer {
  border-top-color: rgba(255, 255, 255, 0.08);
  background: rgba(17, 24, 39, 0.32);
}
</style>

<style>
/* Enhanced Markdown Styles */
.markdown-body {
  @apply text-[15px] leading-[1.75];
  @apply text-gray-700 dark:text-gray-300;
}

.markdown-body h1 {
  @apply mb-6 mt-8 border-b border-gray-200 pb-3 text-3xl font-bold text-gray-900 dark:border-dark-600 dark:text-white;
}

.markdown-body h2 {
  @apply mb-4 mt-7 border-b border-gray-100 pb-2 text-2xl font-bold text-gray-900 dark:border-dark-700 dark:text-white;
}

.markdown-body h3 {
  @apply mb-3 mt-6 text-xl font-semibold text-gray-900 dark:text-white;
}

.markdown-body h4 {
  @apply mb-2 mt-5 text-lg font-semibold text-gray-900 dark:text-white;
}

.markdown-body p {
  @apply mb-4 leading-relaxed;
}

.markdown-body a {
  @apply font-medium text-blue-600 underline decoration-blue-600/30 decoration-2 underline-offset-2 transition-all hover:decoration-blue-600 dark:text-blue-400 dark:decoration-blue-400/30 dark:hover:decoration-blue-400;
}

.markdown-body ul,
.markdown-body ol {
  @apply mb-4 ml-6 space-y-2;
}

.markdown-body ul {
  @apply list-disc;
}

.markdown-body ol {
  @apply list-decimal;
}

.markdown-body li {
  @apply leading-relaxed;
  @apply pl-2;
}

.markdown-body li::marker {
  @apply text-blue-600 dark:text-blue-400;
}

.markdown-body blockquote {
  @apply relative my-5 border-l-4 border-blue-500 bg-blue-50/50 py-3 pl-5 pr-4 italic text-gray-700 dark:border-blue-400 dark:bg-blue-900/10 dark:text-gray-300;
}

.markdown-body blockquote::before {
  content: '"';
  @apply absolute -left-1 top-0 text-5xl font-serif text-blue-500/20 dark:text-blue-400/20;
}

.markdown-body code {
  @apply rounded-lg bg-gray-100 px-2 py-1 text-[13px] font-mono text-pink-600 dark:bg-dark-700 dark:text-pink-400;
}

.markdown-body pre {
  @apply my-5 overflow-x-auto rounded-xl border border-gray-200 bg-gray-50 p-5 dark:border-dark-600 dark:bg-dark-900/50;
}

.markdown-body pre code {
  @apply bg-transparent p-0 text-[13px] text-gray-800 dark:text-gray-200;
}

.markdown-body hr {
  @apply my-8 border-0 border-t-2 border-gray-200 dark:border-dark-700;
}

.markdown-body table {
  @apply mb-5 w-full overflow-hidden rounded-lg border border-gray-200 dark:border-dark-600;
}

.markdown-body th,
.markdown-body td {
  @apply border-r border-b border-gray-200 px-4 py-3 text-left dark:border-dark-600;
}

.markdown-body th:last-child,
.markdown-body td:last-child {
  @apply border-r-0;
}

.markdown-body tr:last-child td {
  @apply border-b-0;
}

.markdown-body th {
  @apply bg-gradient-to-br from-blue-50 to-indigo-50 font-semibold text-gray-900 dark:from-blue-900/20 dark:to-indigo-900/10 dark:text-white;
}

.markdown-body tbody tr {
  @apply transition-colors hover:bg-gray-50 dark:hover:bg-dark-700/30;
}

.markdown-body img {
  @apply my-5 max-w-full rounded-xl border border-gray-200 shadow-md dark:border-dark-600;
}

.markdown-body strong {
  @apply font-semibold text-gray-900 dark:text-white;
}

.markdown-body em {
  @apply italic text-gray-600 dark:text-gray-400;
}
</style>
