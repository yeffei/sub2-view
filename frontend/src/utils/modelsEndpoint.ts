export function safeParseJson(text: string): unknown | null {
  const raw = String(text || '').trim()
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch {
    return null
  }
}

export function extractModelIdsFromModelsPayload(payload: unknown): string[] {
  if (!payload || typeof payload !== 'object') return []
  const data = (payload as { data?: unknown }).data
  if (!Array.isArray(data)) return []

  const ids = data
    .map((item) => {
      if (!item || typeof item !== 'object') return ''
      return String((item as { id?: unknown }).id || '').trim()
    })
    .filter((item) => item.length > 0)

  return Array.from(new Set(ids))
}

export function extractErrorTextFromPayload(payload: unknown): string {
  if (!payload || typeof payload !== 'object') return ''

  const root = payload as { message?: unknown; error?: unknown }
  if (typeof root.message === 'string' && root.message.trim()) return root.message.trim()
  if (typeof root.error === 'string' && root.error.trim()) return root.error.trim()

  if (root.error && typeof root.error === 'object') {
    const err = root.error as { message?: unknown; detail?: unknown; type?: unknown }
    if (typeof err.message === 'string' && err.message.trim()) return err.message.trim()
    if (typeof err.detail === 'string' && err.detail.trim()) return err.detail.trim()
    if (typeof err.type === 'string' && err.type.trim()) return err.type.trim()
  }

  return ''
}
