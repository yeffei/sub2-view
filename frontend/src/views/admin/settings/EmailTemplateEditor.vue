<template>
  <div class="sst-email-template-editor">
    <div
      class="sst-email-template-header"
    >
      <div class="sst-email-template-header-copy">
        <span class="sst-email-template-kicker">山枢庭 · 通知文书</span>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ t("admin.settings.emailTemplates.title") }}
        </h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t("admin.settings.emailTemplates.description") }}
        </p>
      </div>
      <div class="sst-email-template-header-actions">
        <button
          type="button"
          class="btn btn-secondary btn-sm"
          :disabled="loadingTemplate || previewing || !canPreview"
          @click="refreshPreview"
        >
          {{ previewing ? t("admin.settings.emailTemplates.previewing") : t("admin.settings.emailTemplates.preview") }}
        </button>
        <button
          type="button"
          class="btn btn-secondary btn-sm"
          :disabled="loadingTemplate || restoring || !selectedEvent || !selectedLocale"
          @click="restoreOfficial"
        >
          {{ restoring ? t("admin.settings.emailTemplates.restoring") : t("admin.settings.emailTemplates.restoreOfficial") }}
        </button>
        <button
          type="button"
          class="btn btn-primary btn-sm"
          :disabled="loadingTemplate || saving || !canSave"
          @click="saveTemplate"
        >
          {{ saving ? t("admin.settings.emailTemplates.saving") : t("admin.settings.emailTemplates.save") }}
        </button>
      </div>
    </div>

    <div class="space-y-6 p-6">
      <div
        v-if="loadingList"
        class="sst-email-template-loading"
      >
        <span
          class="h-4 w-4 animate-spin rounded-full border-b-2 border-primary-600"
        ></span>
        {{ t("common.loading") }}
      </div>

      <template v-else>
        <div class="sst-email-template-toolbar">
          <div class="sst-email-template-toolbar-copy">
            <span class="sst-email-template-toolbar-label">编排范围</span>
            <p>{{ selectionSummary }}</p>
          </div>
          <div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:min-w-[32rem]">
            <div>
            <label class="input-label" for="email-template-event">
              {{ t("admin.settings.emailTemplates.event") }}
            </label>
            <select
              id="email-template-event"
              v-model="selectedEvent"
              class="input"
              :disabled="loadingTemplate || eventOptions.length === 0"
            >
              <option
                v-for="option in eventOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ formatEventOptionLabel(option) }}
              </option>
            </select>
            </div>
            <div>
            <label class="input-label" for="email-template-locale">
              {{ t("admin.settings.emailTemplates.locale") }}
            </label>
            <select
              id="email-template-locale"
              v-model="selectedLocale"
              class="input"
              :disabled="loadingTemplate || localeOptions.length === 0"
            >
              <option
                v-for="localeOption in localeOptions"
                :key="localeOption"
                :value="localeOption"
              >
                {{ formatLocale(localeOption) }}
              </option>
            </select>
            </div>
          </div>
        </div>

        <div
          v-if="selectedEventMeta"
          class="sst-email-template-meta"
        >
          <div class="flex flex-wrap items-center gap-2">
            <div class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ selectedEventMeta.label }}
            </div>
            <span
              class="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-gray-600 shadow-sm ring-1 ring-gray-200 dark:bg-dark-800 dark:text-gray-300 dark:ring-dark-600"
            >
              {{ selectedEventMeta.categoryLabel }}
            </span>
            <span
              class="rounded-full px-2.5 py-1 text-xs font-medium"
              :class="
                selectedEventMeta.optional
                  ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300'
                  : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300'
              "
            >
              {{ selectedEventMeta.optional ? localText("可退订通知", "Optional") : localText("事务邮件", "Transactional") }}
            </span>
          </div>
          <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-gray-300">
            {{ selectedEventMeta.timing }}
          </p>
          <p
            v-if="selectedEventDescription"
            class="mt-1 text-xs text-gray-500 dark:text-gray-400"
          >
            {{ selectedEventDescription }}
          </p>
        </div>

        <div
          v-if="eventOptions.length && localeOptions.length"
          class="sst-email-template-status-grid"
        >
          <section class="sst-email-template-status-card">
            <span class="sst-email-template-status-label">当前版本</span>
            <strong>{{ currentTemplateStateLabel }}</strong>
            <p>{{ currentTemplateStateHint }}</p>
          </section>
          <section class="sst-email-template-status-card">
            <span class="sst-email-template-status-label">最近改动</span>
            <strong>{{ currentTemplateUpdatedAtText }}</strong>
            <p>{{ templateInventorySummary }}</p>
          </section>
          <section class="sst-email-template-status-card">
            <span class="sst-email-template-status-label">变量目录</span>
            <strong>{{ placeholderList.length }}</strong>
            <p>{{ placeholderSummary }}</p>
          </section>
        </div>

        <div
          v-if="!eventOptions.length || !localeOptions.length"
          class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-700 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-300"
        >
          {{ t("admin.settings.emailTemplates.empty") }}
        </div>

        <div v-else class="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <div class="space-y-4">
            <section class="sst-email-template-panel">
              <div class="sst-email-template-panel-heading">
                <h3>主题与正文</h3>
                <p>按事件与语言维护主题及 HTML 内容，保留占位符即可复用系统变量。</p>
              </div>

              <div class="space-y-4">
                <div>
              <label class="input-label" for="email-template-subject">
                {{ t("admin.settings.emailTemplates.subject") }}
              </label>
              <input
                id="email-template-subject"
                v-model="subject"
                type="text"
                class="input"
                :disabled="loadingTemplate"
                :placeholder="t('admin.settings.emailTemplates.subjectPlaceholder')"
              />
                </div>

                <div>
              <label class="input-label" for="email-template-html">
                {{ t("admin.settings.emailTemplates.html") }}
              </label>
              <textarea
                id="email-template-html"
                v-model="html"
                rows="18"
                class="input min-h-[28rem] resize-y font-mono text-sm leading-6"
                :disabled="loadingTemplate"
                :placeholder="t('admin.settings.emailTemplates.htmlPlaceholder')"
              ></textarea>
                </div>
              </div>
            </section>

            <div
              class="sst-email-template-panel sst-email-template-placeholder-panel"
            >
              <div class="sst-email-template-panel-heading">
                <h3>{{ t("admin.settings.emailTemplates.placeholders") }}</h3>
                <p>{{ t("admin.settings.emailTemplates.placeholdersHelp") }}</p>
              </div>
              <div class="sr-only text-sm font-medium text-gray-900 dark:text-white">
                {{ t("admin.settings.emailTemplates.placeholders") }}
              </div>
              <div class="mt-3 flex flex-wrap gap-2">
                <button
                  v-for="placeholder in placeholderList"
                  :key="placeholder"
                  type="button"
                  class="rounded-full border border-gray-200 bg-white px-3 py-1 font-mono text-xs text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-200 dark:hover:border-primary-500 dark:hover:text-primary-300"
                  @click="copyPlaceholder(placeholder)"
                >
                  {{ placeholder }}
                </button>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <div class="sst-email-template-preview-shell">
              <div class="sst-email-template-preview-header">
                <div>
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ t("admin.settings.emailTemplates.livePreview") }}
                  </div>
                  <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                    {{ previewSubject || t("admin.settings.emailTemplates.noPreview") }}
                  </div>
                </div>
                <span
                  v-if="isCustomTemplate"
                  class="rounded-full bg-primary-50 px-2.5 py-1 text-xs font-medium text-primary-700 dark:bg-primary-900/30 dark:text-primary-300"
                >
                  {{ t("admin.settings.emailTemplates.customized") }}
                </span>
              </div>
              <div class="sst-email-template-preview-frame">
                <iframe
                  class="h-[36rem] w-full rounded-md border border-gray-200 bg-white dark:border-dark-700"
                  sandbox=""
                  :srcdoc="previewHtml"
                  :title="t('admin.settings.emailTemplates.livePreview')"
                ></iframe>
              </div>
            </div>

            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ t("admin.settings.emailTemplates.previewSecurityHint") }}
            </p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { adminAPI } from "@/api";
import type {
  EmailTemplateEventOption,
  EmailTemplateOption,
  EmailTemplateSummary,
} from "@/api/admin/settings";
import { useAppStore } from "@/stores";
import { extractApiErrorMessage } from "@/utils/apiError";

const { t, locale } = useI18n();
const appStore = useAppStore();

const fallbackPlaceholders = [
  "{{site_name}}",
  "{{recipient_name}}",
  "{{recipient_email}}",
  "{{verification_code}}",
  "{{expires_in_minutes}}",
  "{{reset_url}}",
  "{{subscription_group}}",
  "{{subscription_days}}",
  "{{expiry_time}}",
  "{{days_remaining}}",
  "{{current_balance}}",
  "{{threshold}}",
  "{{recharge_url}}",
  "{{recharge_amount}}",
  "{{order_id}}",
  "{{unsubscribe_url}}",
  "{{account_id}}",
  "{{account_name}}",
  "{{platform}}",
  "{{quota_dimension}}",
  "{{quota_used}}",
  "{{quota_limit}}",
  "{{quota_remaining}}",
  "{{quota_threshold}}",
  "{{triggered_at}}",
  "{{group_name}}",
  "{{moderation_category}}",
  "{{moderation_score}}",
  "{{violation_count}}",
  "{{ban_threshold}}",
  "{{rule_name}}",
  "{{severity}}",
  "{{alert_status}}",
  "{{metric_type}}",
  "{{operator}}",
  "{{metric_value}}",
  "{{threshold_value}}",
  "{{alert_description}}",
  "{{report_name}}",
  "{{report_type}}",
  "{{report_start_time}}",
  "{{report_end_time}}",
  "{{report_html}}",
  "{{model}}",
  "{{upstream_message}}",
];

const loadingList = ref(true);
const loadingTemplate = ref(false);
const saving = ref(false);
const previewing = ref(false);
const restoring = ref(false);
const eventOptions = ref<EmailTemplateOption[]>([]);
const localeOptions = ref<string[]>([]);
const templateSummaries = ref<EmailTemplateSummary[]>([]);
const selectedEvent = ref("");
const selectedLocale = ref("");
const subject = ref("");
const html = ref("");
const isCustomTemplate = ref(false);
const placeholders = ref<string[]>([]);
const previewSubject = ref("");
const previewHtml = ref("");
const initializingSelection = ref(false);

interface EventDisplayMeta {
  label: string;
  timing: string;
  categoryLabel: string;
}

function localText(zh: string, en: string): string {
  return locale.value.toLowerCase().startsWith("zh") ? zh : en;
}

const eventDisplayMeta: Record<string, EventDisplayMeta> = {
  "auth.verify_code": {
    label: "邮箱验证码",
    timing: "注册、绑定邮箱、OAuth 补全邮箱或 TOTP 邮箱校验时发送。",
    categoryLabel: "认证安全",
  },
  "auth.password_reset": {
    label: "密码重置",
    timing: "用户请求密码重置链接时发送。",
    categoryLabel: "认证安全",
  },
  "notification_email.verify_code": {
    label: "通知邮箱验证码",
    timing: "用户添加并验证额外通知邮箱时发送。",
    categoryLabel: "认证安全",
  },
  "subscription.purchase_success": {
    label: "订阅开通成功",
    timing: "订阅订单完成支付并成功开通或续期后发送。",
    categoryLabel: "订阅",
  },
  "subscription.expiry_reminder": {
    label: "订阅到期提醒",
    timing: "后台任务在订阅仍有效且距离到期剩余 7 天、3 天、1 天时各发送一次，可通过邮件设置中的开关关闭。",
    categoryLabel: "订阅",
  },
  "balance.low": {
    label: "余额不足提醒",
    timing: "用户余额低于全局或个人配置的提醒阈值时发送。",
    categoryLabel: "计费",
  },
  "balance.recharge_success": {
    label: "余额充值成功",
    timing: "余额充值订单支付完成并入账后发送。",
    categoryLabel: "计费",
  },
  "account.quota_alert": {
    label: "账号限额告警",
    timing: "上游账号的用量达到配置的额度告警阈值时发送给管理员通知邮箱。",
    categoryLabel: "管理告警",
  },
  "content_moderation.violation_notice": {
    label: "内容审计违规提醒",
    timing: "用户请求命中内容审计或风控规则、但尚未被禁用时发送。",
    categoryLabel: "风控",
  },
  "content_moderation.account_disabled": {
    label: "内容审计禁用账号",
    timing: "内容审计违规次数达到封禁阈值并自动禁用用户账号时发送。",
    categoryLabel: "风控",
  },
  "content_moderation.cyber_policy_notice": {
    label: "网络安全策略拦截提醒",
    timing: "请求被上游网络安全策略拦截时发送，便于说明触发时间、分组、模型与上游提示。",
    categoryLabel: "风控",
  },
  "ops.alert": {
    label: "运维告警",
    timing: "运维监控规则触发告警并满足邮件通知配置时发送给运维收件人。",
    categoryLabel: "运维",
  },
  "ops.scheduled_report": {
    label: "运维定时报表",
    timing: "运维日报、周报、错误摘要或账号健康报表到达配置的发送时间时发送。",
    categoryLabel: "运维",
  },
};

const eventDisplayMetaEn: Record<string, EventDisplayMeta> = {
  "auth.verify_code": {
    label: "Email Verification Code",
    timing: "Sent for registration, email binding, OAuth pending email completion, or TOTP email verification.",
    categoryLabel: "Auth",
  },
  "auth.password_reset": {
    label: "Password Reset",
    timing: "Sent when a user requests a password reset link.",
    categoryLabel: "Auth",
  },
  "notification_email.verify_code": {
    label: "Notification Email Verification",
    timing: "Sent when a user adds and verifies an extra notification email address.",
    categoryLabel: "Auth",
  },
  "subscription.purchase_success": {
    label: "Subscription Activated",
    timing: "Sent after a subscription order is paid and the subscription is activated or extended.",
    categoryLabel: "Subscription",
  },
  "subscription.expiry_reminder": {
    label: "Subscription Expiry Reminder",
    timing: "Sent by the background job when an active subscription has 7, 3, or 1 day remaining. It can be disabled in Email settings.",
    categoryLabel: "Subscription",
  },
  "balance.low": {
    label: "Low Balance Alert",
    timing: "Sent when a user's balance drops below the global or personal reminder threshold.",
    categoryLabel: "Billing",
  },
  "balance.recharge_success": {
    label: "Balance Recharge Success",
    timing: "Sent after a balance recharge order is paid and credited.",
    categoryLabel: "Billing",
  },
  "account.quota_alert": {
    label: "Account Quota Alert",
    timing: "Sent to admin notification emails when an upstream account reaches the configured quota alert threshold.",
    categoryLabel: "Admin",
  },
  "content_moderation.violation_notice": {
    label: "Risk Control Violation Notice",
    timing: "Sent when a user request triggers content moderation or risk-control rules but the account is not disabled yet.",
    categoryLabel: "Risk Control",
  },
  "content_moderation.account_disabled": {
    label: "Risk Control Account Disabled",
    timing: "Sent when content moderation reaches the ban threshold and automatically disables the user account.",
    categoryLabel: "Risk Control",
  },
  "content_moderation.cyber_policy_notice": {
    label: "Cyber Policy Notice",
    timing: "Sent when an upstream cyber-security policy blocks a request, with the trigger time, group, model, and upstream explanation.",
    categoryLabel: "Risk Control",
  },
  "ops.alert": {
    label: "Ops Alert",
    timing: "Sent to ops recipients when an ops monitoring rule fires and email notification settings allow it.",
    categoryLabel: "Ops",
  },
  "ops.scheduled_report": {
    label: "Ops Scheduled Report",
    timing: "Sent when a configured daily, weekly, error digest, or account health report reaches its scheduled send time.",
    categoryLabel: "Ops",
  },
};

function normalizeEventOption(option: EmailTemplateEventOption): EmailTemplateOption {
  if (typeof option === "string") {
    return { value: option };
  }
  return option;
}

function eventMetaFor(option?: EmailTemplateOption | null) {
  if (!option) return null;
  const displayMeta = (
    locale.value.toLowerCase().startsWith("zh")
      ? eventDisplayMeta
      : eventDisplayMetaEn
  )[option.value];
  const label = displayMeta?.label || option.label || option.value;
  const timing = displayMeta?.timing || option.description || "";
  const categoryLabel =
    displayMeta?.categoryLabel || formatCategory(option.category || "");
  return {
    label,
    timing,
    categoryLabel,
    optional: option.optional === true,
  };
}

function formatEventOptionLabel(option: EmailTemplateOption): string {
  const meta = eventMetaFor(option);
  if (!meta) return option.label || option.value;
  return meta.label;
}

function formatCategory(category: string): string {
  const normalized = category.trim().toLowerCase();
  if (!normalized) return localText("通知", "Notification");
  const labels: Record<string, { zh: string; en: string }> = {
    auth: { zh: "认证安全", en: "Auth" },
    subscription: { zh: "订阅", en: "Subscription" },
    billing: { zh: "计费", en: "Billing" },
    admin: { zh: "管理告警", en: "Admin" },
    risk_control: { zh: "风控", en: "Risk Control" },
    ops: { zh: "运维", en: "Ops" },
  };
  const item = labels[normalized];
  return item ? localText(item.zh, item.en) : category;
}

const selectedEventOption = computed(() => {
  return (
    eventOptions.value.find((option) => option.value === selectedEvent.value) ||
    null
  );
});

const selectedEventMeta = computed(() => eventMetaFor(selectedEventOption.value));

const selectedEventDescription = computed(() => {
  return (
    selectedEventOption.value?.description || ""
  );
});

const totalTemplateCount = computed(() => templateSummaries.value.length);

const customTemplateCount = computed(
  () => templateSummaries.value.filter((template) => template.is_custom).length,
);

const currentTemplateSummary = computed(() => {
  return (
    templateSummaries.value.find(
      (template) =>
        template.event === selectedEvent.value &&
        template.locale === selectedLocale.value,
    ) || null
  );
});

const currentTemplateStateLabel = computed(() => {
  return isCustomTemplate.value
    ? localText("已自定义", "Customized")
    : localText("官方默认", "Official");
});

const currentTemplateStateHint = computed(() => {
  return isCustomTemplate.value
    ? localText("当前事件与语言正在使用覆盖版本。", "This event and locale currently use an overridden template.")
    : localText("当前事件与语言正在使用系统内置模板。", "This event and locale currently use the built-in template.");
});

const currentTemplateUpdatedAtText = computed(() => {
  const updatedAt = currentTemplateSummary.value?.updated_at;
  if (!updatedAt) {
    return localText("未单独修改", "No custom update yet");
  }
  return formatTimestamp(updatedAt);
});

const templateInventorySummary = computed(() => {
  if (!totalTemplateCount.value) {
    return localText("等待加载可编辑模板目录。", "Loading the editable template catalog.");
  }
  return localText(
    `当前共 ${totalTemplateCount.value} 份模板，其中 ${customTemplateCount.value} 份为自定义版本。`,
    `${totalTemplateCount.value} templates are available, with ${customTemplateCount.value} customized versions.`,
  );
});

const placeholderSummary = computed(() => {
  return localText(
    `当前可用 ${placeholderList.value.length} 个占位符，可直接复制到主题或 HTML 中。`,
    `${placeholderList.value.length} placeholders are available to copy into the subject or HTML.`,
  );
});

const selectionSummary = computed(() => {
  const eventLabel = selectedEventMeta.value?.label || localText("未选择事件", "No event selected");
  const localeLabel = selectedLocale.value ? formatLocale(selectedLocale.value) : localText("未选择语言", "No locale selected");
  return localText(
    `当前正在编排 ${eventLabel} 的 ${localeLabel} 模板，并同步预览邮件主题与正文。${templateInventorySummary.value}`,
    `You are editing the ${localeLabel} template for ${eventLabel}, with the subject and HTML preview kept in sync. ${templateInventorySummary.value}`,
  );
});

const placeholderList = computed(() => {
  const combined = [...placeholders.value, ...fallbackPlaceholders];
  return Array.from(
    new Set(
      combined
        .map((item) => formatPlaceholder(item))
        .filter((item) => item.length > 0),
    ),
  );
});

function formatPlaceholder(placeholder: string): string {
  const trimmed = placeholder.trim();
  if (!trimmed) return "";
  if (trimmed.startsWith("{{") && trimmed.endsWith("}}")) return trimmed;
  return `{{${trimmed}}}`;
}

const canSave = computed(
  () =>
    Boolean(selectedEvent.value && selectedLocale.value) &&
    subject.value.trim().length > 0 &&
    html.value.trim().length > 0,
);

const canPreview = computed(
  () => Boolean(selectedEvent.value && selectedLocale.value) && html.value.trim().length > 0,
);

function formatLocale(locale: string): string {
  const lower = locale.toLowerCase();
  if (lower === "zh" || lower.startsWith("zh-")) {
    return t("admin.settings.emailTemplates.localeZh");
  }
  if (lower === "en" || lower.startsWith("en-")) {
    return t("admin.settings.emailTemplates.localeEn");
  }
  return locale;
}

function selectInitialLocale(locales: string[]): string {
  const currentLocale = locale.value.toLowerCase();
  const exactMatch = locales.find(
    (availableLocale) => availableLocale.toLowerCase() === currentLocale,
  );
  if (exactMatch) return exactMatch;

  const currentLanguage = currentLocale.split("-")[0];
  const languageMatch = locales.find(
    (availableLocale) => availableLocale.toLowerCase().split("-")[0] === currentLanguage,
  );
  if (languageMatch) return languageMatch;

  return locales[0] || "";
}

function formatTimestamp(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(locale.value, {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function upsertTemplateSummary(summary: EmailTemplateSummary) {
  const next = [...templateSummaries.value];
  const index = next.findIndex(
    (template) =>
      template.event === summary.event &&
      template.locale === summary.locale,
  );
  if (index >= 0) {
    next[index] = summary;
  } else {
    next.push(summary);
  }
  templateSummaries.value = next;
}

function applyTemplate(template: {
  event?: string;
  locale?: string;
  subject: string;
  html: string;
  is_custom?: boolean;
  updated_at?: string;
  placeholders?: string[];
}) {
  subject.value = template.subject;
  html.value = template.html;
  isCustomTemplate.value = template.is_custom === true;
  placeholders.value = template.placeholders || [];

  const eventValue = template.event || selectedEvent.value;
  const localeValue = template.locale || selectedLocale.value;
  if (eventValue && localeValue) {
    upsertTemplateSummary({
      event: eventValue,
      locale: localeValue,
      subject: template.subject,
      is_custom: template.is_custom,
      updated_at: template.updated_at,
    });
  }
}

async function loadTemplate() {
  if (!selectedEvent.value || !selectedLocale.value) return;
  loadingTemplate.value = true;
  try {
    const template = await adminAPI.settings.getEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
    );
    applyTemplate(template);
    await refreshPreview();
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  } finally {
    loadingTemplate.value = false;
  }
}

async function loadTemplateList() {
  loadingList.value = true;
  try {
    const response = await adminAPI.settings.getEmailTemplates();
    eventOptions.value = response.events.map(normalizeEventOption);
    localeOptions.value = response.locales;
    templateSummaries.value = response.templates || [];
    placeholders.value = response.placeholders || [];
    initializingSelection.value = true;
    selectedEvent.value = eventOptions.value[0]?.value || "";
    selectedLocale.value = selectInitialLocale(response.locales);
    await loadTemplate();
    initializingSelection.value = false;
  } catch (err: unknown) {
    initializingSelection.value = false;
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  } finally {
    loadingList.value = false;
  }
}

async function saveTemplate() {
  if (!canSave.value) {
    appStore.showError(t("admin.settings.emailTemplates.validationRequired"));
    return;
  }
  saving.value = true;
  try {
    const template = await adminAPI.settings.updateEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
      {
        subject: subject.value,
        html: html.value,
      },
    );
    applyTemplate(template);
    await refreshPreview();
    appStore.showSuccess(t("admin.settings.emailTemplates.saveSuccess"));
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  } finally {
    saving.value = false;
  }
}

async function refreshPreview() {
  if (!canPreview.value) {
    previewSubject.value = "";
    previewHtml.value = "";
    return;
  }
  previewing.value = true;
  try {
    const preview = await adminAPI.settings.previewEmailTemplate({
      event: selectedEvent.value,
      locale: selectedLocale.value,
      subject: subject.value,
      html: html.value,
    });
    previewSubject.value = preview.subject;
    previewHtml.value = preview.html;
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  } finally {
    previewing.value = false;
  }
}

async function restoreOfficial() {
  if (!selectedEvent.value || !selectedLocale.value) return;
  if (!window.confirm(t("admin.settings.emailTemplates.restoreConfirm"))) return;

  restoring.value = true;
  try {
    const template = await adminAPI.settings.restoreOfficialEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
    );
    applyTemplate(template);
    await refreshPreview();
    appStore.showSuccess(t("admin.settings.emailTemplates.restoreSuccess"));
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  } finally {
    restoring.value = false;
  }
}

async function copyPlaceholder(placeholder: string) {
  try {
    await navigator.clipboard.writeText(placeholder);
    appStore.showSuccess(t("admin.settings.emailTemplates.placeholderCopied"));
  } catch {
    appStore.showError(t("common.error"));
  }
}

watch([selectedEvent, selectedLocale], ([eventValue, localeValue], [oldEvent, oldLocale]) => {
  if (initializingSelection.value) return;
  if (!eventValue || !localeValue) return;
  if (eventValue === oldEvent && localeValue === oldLocale) return;
  void loadTemplate();
});

onMounted(() => {
  void loadTemplateList();
});
</script>

<style scoped>
.sst-email-template-editor {
  border: 1px solid rgba(198, 184, 157, 0.46);
  border-radius: 22px;
  background:
    radial-gradient(circle at top left, rgba(167, 58, 42, 0.06), transparent 28%),
    linear-gradient(180deg, rgba(252, 249, 243, 0.98), rgba(248, 243, 232, 0.92));
  box-shadow: 0 24px 52px -42px rgba(58, 48, 34, 0.38);
}

.sst-email-template-header {
  @apply flex flex-col gap-4 border-b px-6 py-5 lg:flex-row lg:items-start lg:justify-between;
  border-color: rgba(198, 184, 157, 0.42);
}

.sst-email-template-header-copy {
  @apply space-y-1;
}

.sst-email-template-kicker {
  @apply inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-medium uppercase tracking-[0.18em];
  color: #8b5e3c;
  background: rgba(167, 58, 42, 0.08);
}

.sst-email-template-header-actions {
  @apply flex flex-wrap gap-2;
}

.sst-email-template-loading {
  @apply flex items-center gap-2 rounded-xl border px-4 py-3 text-sm;
  border-color: rgba(198, 184, 157, 0.38);
  background: rgba(255, 252, 246, 0.75);
  color: #6a6a63;
}

.sst-email-template-toolbar {
  @apply flex flex-col gap-4 rounded-2xl border px-4 py-4 xl:flex-row xl:items-start xl:justify-between;
  border-color: rgba(198, 184, 157, 0.4);
  background: rgba(255, 252, 246, 0.72);
}

.sst-email-template-toolbar-copy {
  @apply space-y-2;
}

.sst-email-template-toolbar-label {
  @apply inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-medium uppercase tracking-[0.18em];
  color: #8b5e3c;
  background: rgba(167, 58, 42, 0.08);
}

.sst-email-template-toolbar-copy p {
  @apply text-sm leading-6;
  color: #5f6257;
}

.sst-email-template-meta,
.sst-email-template-panel,
.sst-email-template-preview-shell {
  border: 1px solid rgba(198, 184, 157, 0.42);
  border-radius: 20px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.9), rgba(248, 243, 231, 0.84));
  box-shadow: 0 18px 42px -38px rgba(58, 48, 34, 0.32);
}

.sst-email-template-meta {
  @apply p-4;
}

.sst-email-template-status-grid {
  @apply grid grid-cols-1 gap-3 md:grid-cols-3;
}

.sst-email-template-status-card {
  @apply rounded-2xl border px-4 py-4;
  border-color: rgba(198, 184, 157, 0.38);
  background: rgba(255, 252, 246, 0.76);
}

.sst-email-template-status-card strong {
  @apply mt-2 block text-base font-semibold text-gray-900 dark:text-white;
}

.sst-email-template-status-card p {
  @apply mt-2 text-xs leading-6 text-gray-500 dark:text-gray-400;
}

.sst-email-template-status-label {
  @apply text-[11px] font-medium uppercase tracking-[0.18em];
  color: #8b5e3c;
}

.sst-email-template-panel {
  @apply p-4 sm:p-5;
}

.sst-email-template-panel-heading {
  @apply mb-4 space-y-1;
}

.sst-email-template-panel-heading h3 {
  @apply text-sm font-semibold text-gray-900 dark:text-white;
}

.sst-email-template-panel-heading p {
  @apply text-xs leading-6 text-gray-500 dark:text-gray-400;
}

.sst-email-template-placeholder-panel :deep(button) {
  box-shadow: 0 10px 18px -18px rgba(58, 48, 34, 0.42);
}

.sst-email-template-preview-header {
  @apply flex items-center justify-between gap-3 border-b px-4 py-4 sm:px-5;
  border-color: rgba(198, 184, 157, 0.36);
}

.sst-email-template-preview-frame {
  padding: 14px;
  background:
    linear-gradient(180deg, rgba(239, 232, 218, 0.72), rgba(247, 242, 233, 0.82));
}

</style>
<style>
.dark .sst-email-template-editor,
.dark .sst-email-template-toolbar,
.dark .sst-email-template-meta,
.dark .sst-email-template-status-card,
.dark .sst-email-template-panel,
.dark .sst-email-template-preview-shell,
.dark .sst-email-template-loading {
  border-color: rgba(58, 61, 54, 0.96);
  background:
    linear-gradient(180deg, rgba(24, 26, 21, 0.92), rgba(18, 20, 16, 0.96));
}

.dark .sst-email-template-header {
  border-color: rgba(58, 61, 54, 0.96);
}

.dark .sst-email-template-kicker,
.dark .sst-email-template-toolbar-label,
.dark .sst-email-template-status-label {
  color: #e7b58e;
  background: rgba(167, 58, 42, 0.22);
}

.dark .sst-email-template-toolbar-copy p {
  color: #9ea49a;
}

.dark .sst-email-template-preview-frame {
  background: rgba(16, 18, 14, 0.76);
}
</style>
