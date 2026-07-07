import type { RouteLocationNormalizedLoaded } from 'vue-router'

export const DEFAULT_SEO_DESCRIPTION = '山枢庭 SST 提供统一的 AI API 入口，将接入、计量、账册与稳定供给整理在同一套长期可维护的服务中。'
export const DEFAULT_SITE_NAME = '山枢庭'

const PUBLIC_CANONICAL_PATHS = new Set([
  '/home',
  '/pricing',
  '/docs',
  '/docs/openai-compatible-api',
  '/docs/base-url',
  '/docs/api-key',
  '/docs/streaming',
  '/docs/codex',
  '/docs/claude-code',
  '/faq',
  '/privacy',
  '/terms',
])

function normalizeSiteName(siteName?: string): string {
  const trimmed = siteName?.trim()
  return trimmed && trimmed !== 'Sub2API' ? trimmed : DEFAULT_SITE_NAME
}

function normalizePath(path: string): string {
  if (!path || path === '/') return '/home'
  return path.split(/[?#]/)[0] || '/home'
}

function buildCanonicalUrl(path: string): string {
  const canonicalPath = normalizePath(path)
  return `${window.location.origin}${canonicalPath}`
}

function upsertMetaByName(name: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[name="${name}"]`)
  if (!meta) {
    meta = document.createElement('meta')
    meta.setAttribute('name', name)
    document.head.appendChild(meta)
  }
  meta.setAttribute('content', content)
}

function upsertMetaByProperty(property: string, content: string) {
  let meta = document.head.querySelector<HTMLMetaElement>(`meta[property="${property}"]`)
  if (!meta) {
    meta = document.createElement('meta')
    meta.setAttribute('property', property)
    document.head.appendChild(meta)
  }
  meta.setAttribute('content', content)
}

function upsertCanonical(href: string) {
  let link = document.head.querySelector<HTMLLinkElement>('link[rel="canonical"]')
  if (!link) {
    link = document.createElement('link')
    link.setAttribute('rel', 'canonical')
    document.head.appendChild(link)
  }
  link.setAttribute('href', href)
}

function resolveRobots(route: RouteLocationNormalizedLoaded, canonicalPath: string): string {
  if (typeof route.meta.robots === 'string' && route.meta.robots.trim()) {
    return route.meta.robots.trim()
  }

  if (route.name === 'NotFound') {
    return 'noindex,nofollow'
  }

  if (route.meta.requiresAuth !== false) {
    return 'noindex,nofollow'
  }

  return PUBLIC_CANONICAL_PATHS.has(canonicalPath) ? 'index,follow' : 'noindex,nofollow'
}

export function applyRouteSeo(route: RouteLocationNormalizedLoaded, documentTitle: string, siteName?: string) {
  const normalizedSiteName = normalizeSiteName(siteName)
  const canonicalPath = normalizePath(
    typeof route.meta.canonicalPath === 'string' && route.meta.canonicalPath.trim()
      ? route.meta.canonicalPath
      : route.path,
  )
  const description = typeof route.meta.seoDescription === 'string' && route.meta.seoDescription.trim()
    ? route.meta.seoDescription.trim()
    : DEFAULT_SEO_DESCRIPTION
  const canonicalUrl = buildCanonicalUrl(canonicalPath)
  const robots = resolveRobots(route, canonicalPath)

  upsertMetaByName('description', description)
  upsertMetaByName('robots', robots)
  upsertCanonical(canonicalUrl)

  upsertMetaByProperty('og:site_name', normalizedSiteName)
  upsertMetaByProperty('og:title', documentTitle)
  upsertMetaByProperty('og:description', description)
  upsertMetaByProperty('og:type', 'website')
  upsertMetaByProperty('og:url', canonicalUrl)
  upsertMetaByProperty('og:image', `${window.location.origin}/logo.png`)

  upsertMetaByName('twitter:card', 'summary_large_image')
  upsertMetaByName('twitter:title', documentTitle)
  upsertMetaByName('twitter:description', description)
  upsertMetaByName('twitter:image', `${window.location.origin}/logo.png`)
}
