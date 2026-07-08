import type { ApiKeyWorkbenchLatestError } from '@/api/usage'

export type WorkbenchModelInsight = {
  status: string
  detail: string
  action: string
  to: string
  tone: 'calm' | 'notice' | 'alert'
}

type WorkbenchLocale = 'zh' | 'en'

const workbenchText = {
  zh: {
    noFailures: '当前无失败记录',
    recentFailure: '最近失败',
    label: (label: string) => `最近失败 ${label}`,
    hints: {
      unsupported: (model: string) => `模型 ${model} 当前分组未开放`,
      unavailable: (model: string) => `模型 ${model} 当前线路未开放`,
      rateLimited: (model: string) => `模型 ${model} 当前线路限流中`,
      noRoute: (model: string) => `模型 ${model} 当前无可用线路`,
      mapped: (requested: string, upstream: string) => `请求 ${requested} -> 实际上游 ${upstream}`,
    },
    insight: {
      keyPrefix: (keyName?: string) => keyName ? `密钥 ${keyName} ` : '这把 Key ',
      pending: '待处理',
      notice: '留意',
      connectionCheck: '去做接入体检',
      errorLedger: '查看错误账册',
      recentFailures: '查看最近失败',
      unsupported: (prefix: string, model: string) => `${prefix}最近请求的模型 ${model} 当前分组未开放，先核对模型名、分组和客户端配置。`,
      unavailable: (prefix: string, model: string) => `${prefix}最近请求的模型 ${model} 当前线路未开放，更像路由覆盖不到该模型而不是瞬时波动。`,
      rateLimited: (prefix: string, model: string) => `${prefix}最近请求的模型 ${model} 当前线路限流中，可先错峰重试或临时换模型。`,
      noRoute: (prefix: string, model: string) => `${prefix}最近请求的模型 ${model} 当前没有可用线路，建议先做体检确认模型列表与路由是否已恢复。`,
      mapped: (prefix: string, requested: string, upstream: string) => `${prefix}最近一次请求以 ${requested} 进入，但实际上游模型是 ${upstream}，如果这不是你的预期，优先检查分组映射。`,
    }
  },
  en: {
    noFailures: 'No recent failures',
    recentFailure: 'Recent failure',
    label: (label: string) => `Recent failure: ${label}`,
    hints: {
      unsupported: (model: string) => `Model ${model} is not enabled for the current group`,
      unavailable: (model: string) => `Model ${model} is not enabled on the current route`,
      rateLimited: (model: string) => `Model ${model} is currently rate-limited on this route`,
      noRoute: (model: string) => `Model ${model} has no available route right now`,
      mapped: (requested: string, upstream: string) => `Requested ${requested} -> actual upstream ${upstream}`,
    },
    insight: {
      keyPrefix: (keyName?: string) => keyName ? `Key ${keyName} ` : 'This key ',
      pending: 'Action needed',
      notice: 'Check',
      connectionCheck: 'Run connection check',
      errorLedger: 'View error ledger',
      recentFailures: 'View recent failures',
      unsupported: (prefix: string, model: string) => `${prefix}recently requested model ${model}, but it is not enabled for the current group. Check the model name, group, and client config first.`,
      unavailable: (prefix: string, model: string) => `${prefix}recently requested model ${model}, but it is not enabled on the current route. This looks more like route coverage than a transient issue.`,
      rateLimited: (prefix: string, model: string) => `${prefix}recently requested model ${model}, and the current route is rate-limited. Try off-peak retry or switch models temporarily.`,
      noRoute: (prefix: string, model: string) => `${prefix}recently requested model ${model}, but no route is available. Run a connection check to confirm whether the model list and route recovered.`,
      mapped: (prefix: string, requested: string, upstream: string) => `${prefix}entered as ${requested}, but the actual upstream model was ${upstream}. If that is unexpected, check group mapping first.`,
    }
  }
}

const normalizeWorkbenchLocale = (locale?: string): WorkbenchLocale => locale?.startsWith('en') ? 'en' : 'zh'

export function buildWorkbenchLatestErrorLabel(
  latest: ApiKeyWorkbenchLatestError | null | undefined,
  reasonLabels: Record<string, string>,
  categoryLabels: Record<string, string>,
  locale?: string,
): string {
  const copy = workbenchText[normalizeWorkbenchLocale(locale)]
  if (!latest) return copy.noFailures
  const label = latest.reason_code
    ? reasonLabels[latest.reason_code] || categoryLabels[latest.category] || copy.recentFailure
    : categoryLabels[latest.category] || copy.recentFailure
  return copy.label(label)
}

export function buildWorkbenchModelHints(latest: ApiKeyWorkbenchLatestError | null | undefined, locale?: string): string[] {
  if (!latest) return []
  const copy = workbenchText[normalizeWorkbenchLocale(locale)].hints

  const requested = String(latest.requested_model || '').trim()
  const upstream = String(latest.upstream_model || '').trim()
  const hints: string[] = []

  switch (latest.reason_code) {
    case 'request_model_not_supported':
      if (requested) hints.push(copy.unsupported(requested))
      break
    case 'service_model_not_available':
      if (requested) hints.push(copy.unavailable(requested))
      break
    case 'service_model_rate_limited':
      if (requested) hints.push(copy.rateLimited(requested))
      break
    case 'service_no_route_available':
      if (requested) hints.push(copy.noRoute(requested))
      break
  }

  if (requested && upstream && requested !== upstream) {
    hints.push(copy.mapped(requested, upstream))
  }

  return Array.from(new Set(hints))
}

export function buildWorkbenchModelInsight(
  latest: ApiKeyWorkbenchLatestError | null | undefined,
  keyName?: string,
  locale?: string,
): WorkbenchModelInsight | null {
  if (!latest) return null
  const copy = workbenchText[normalizeWorkbenchLocale(locale)].insight

  const requested = String(latest.requested_model || '').trim()
  const upstream = String(latest.upstream_model || '').trim()
  const keyPrefix = copy.keyPrefix(keyName)

  switch (latest.reason_code) {
    case 'request_model_not_supported':
      if (!requested) return null
      return {
        status: copy.pending,
        detail: copy.unsupported(keyPrefix, requested),
        action: copy.connectionCheck,
        to: '/keys?panel=connection-test',
        tone: 'alert',
      }
    case 'service_model_not_available':
      if (!requested) return null
      return {
        status: copy.pending,
        detail: copy.unavailable(keyPrefix, requested),
        action: copy.connectionCheck,
        to: '/keys?panel=connection-test',
        tone: 'alert',
      }
    case 'service_model_rate_limited':
      if (!requested) return null
      return {
        status: copy.notice,
        detail: copy.rateLimited(keyPrefix, requested),
        action: copy.errorLedger,
        to: '/usage?tab=errors&category=service_unavailable',
        tone: 'notice',
      }
    case 'service_no_route_available':
      if (!requested) return null
      return {
        status: copy.notice,
        detail: copy.noRoute(keyPrefix, requested),
        action: copy.connectionCheck,
        to: '/keys?panel=connection-test',
        tone: 'notice',
      }
  }

  if (requested && upstream && requested !== upstream) {
    return {
      status: copy.notice,
      detail: copy.mapped(keyPrefix, requested, upstream),
      action: copy.recentFailures,
      to: '/usage?tab=errors',
      tone: 'notice',
    }
  }

  return null
}
