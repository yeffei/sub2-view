export interface PublicContactEntry {
  label: string
  href: string
  external: boolean
}

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const phonePattern = /^\+?[\d\s().-]{6,}$/

export function resolvePublicContact(raw: string | null | undefined): PublicContactEntry | null {
  const value = raw?.trim()

  if (!value) {
    return null
  }

  if (/^https?:\/\//i.test(value)) {
    return {
      label: value.replace(/^https?:\/\//i, '').replace(/\/$/, ''),
      href: value,
      external: true,
    }
  }

  if (emailPattern.test(value)) {
    return {
      label: value,
      href: `mailto:${value}`,
      external: false,
    }
  }

  if (phonePattern.test(value)) {
    const dialValue = value.replace(/[^\d+]/g, '')

    return {
      label: value,
      href: `tel:${dialValue}`,
      external: false,
    }
  }

  return {
    label: value,
    href: '',
    external: false,
  }
}
