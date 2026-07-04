import type { ApiKeyWorkbenchLatestError } from '@/api/usage'

export type WorkbenchModelInsight = {
  status: string
  detail: string
  action: string
  to: string
  tone: 'calm' | 'notice' | 'alert'
}

export function buildWorkbenchLatestErrorLabel(
  latest: ApiKeyWorkbenchLatestError | null | undefined,
  reasonLabels: Record<string, string>,
  categoryLabels: Record<string, string>,
): string {
  if (!latest) return '当前无失败记录'
  const label = latest.reason_code
    ? reasonLabels[latest.reason_code] || categoryLabels[latest.category] || '最近失败'
    : categoryLabels[latest.category] || '最近失败'
  return `最近失败 ${label}`
}

export function buildWorkbenchModelHints(latest: ApiKeyWorkbenchLatestError | null | undefined): string[] {
  if (!latest) return []

  const requested = String(latest.requested_model || '').trim()
  const upstream = String(latest.upstream_model || '').trim()
  const hints: string[] = []

  switch (latest.reason_code) {
    case 'request_model_not_supported':
      if (requested) hints.push(`模型 ${requested} 当前分组未开放`)
      break
    case 'service_model_not_available':
      if (requested) hints.push(`模型 ${requested} 当前线路未开放`)
      break
    case 'service_model_rate_limited':
      if (requested) hints.push(`模型 ${requested} 当前线路限流中`)
      break
    case 'service_no_route_available':
      if (requested) hints.push(`模型 ${requested} 当前无可用线路`)
      break
  }

  if (requested && upstream && requested !== upstream) {
    hints.push(`请求 ${requested} -> 实际上游 ${upstream}`)
  }

  return Array.from(new Set(hints))
}

export function buildWorkbenchModelInsight(
  latest: ApiKeyWorkbenchLatestError | null | undefined,
  keyName?: string,
): WorkbenchModelInsight | null {
  if (!latest) return null

  const requested = String(latest.requested_model || '').trim()
  const upstream = String(latest.upstream_model || '').trim()
  const keyPrefix = keyName ? `密钥 ${keyName} ` : '这把 Key '

  switch (latest.reason_code) {
    case 'request_model_not_supported':
      if (!requested) return null
      return {
        status: '待处理',
        detail: `${keyPrefix}最近请求的模型 ${requested} 当前分组未开放，先核对模型名、分组和客户端配置。`,
        action: '去做接入体检',
        to: '/keys?panel=connection-test',
        tone: 'alert',
      }
    case 'service_model_not_available':
      if (!requested) return null
      return {
        status: '待处理',
        detail: `${keyPrefix}最近请求的模型 ${requested} 当前线路未开放，更像路由覆盖不到该模型而不是瞬时波动。`,
        action: '去做接入体检',
        to: '/keys?panel=connection-test',
        tone: 'alert',
      }
    case 'service_model_rate_limited':
      if (!requested) return null
      return {
        status: '留意',
        detail: `${keyPrefix}最近请求的模型 ${requested} 当前线路限流中，可先错峰重试或临时换模型。`,
        action: '查看错误账册',
        to: '/usage?tab=errors&category=service_unavailable',
        tone: 'notice',
      }
    case 'service_no_route_available':
      if (!requested) return null
      return {
        status: '留意',
        detail: `${keyPrefix}最近请求的模型 ${requested} 当前没有可用线路，建议先做体检确认模型列表与路由是否已恢复。`,
        action: '去做接入体检',
        to: '/keys?panel=connection-test',
        tone: 'notice',
      }
  }

  if (requested && upstream && requested !== upstream) {
    return {
      status: '留意',
      detail: `${keyPrefix}最近一次请求以 ${requested} 进入，但实际上游模型是 ${upstream}，如果这不是你的预期，优先检查分组映射。`,
      action: '查看最近失败',
      to: '/usage?tab=errors',
      tone: 'notice',
    }
  }

  return null
}
