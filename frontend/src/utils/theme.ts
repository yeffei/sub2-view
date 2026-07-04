import { onBeforeUnmount, onMounted, ref } from 'vue'

const THEME_KEY = 'theme'
const THEME_EVENT = 'sst-theme-change'
const THEME_PREFERENCE_EVENT = 'sst-theme-preference-change'
const LIGHT_BACKGROUND = '#f4efe4'
const DARK_BACKGROUND = '#11130f'
const LIGHT_TEXT = '#1f2320'
const DARK_TEXT = '#f4efe4'

export type ThemePreference = 'light' | 'dark' | 'system'

export const themePreferenceLabels: Record<ThemePreference, string> = {
  system: '跟随系统',
  light: '纸面',
  dark: '夜庭',
}

function readStoredTheme(): ThemePreference | null {
  try {
    const savedTheme = window.localStorage.getItem(THEME_KEY)
    if (savedTheme === 'light' || savedTheme === 'dark' || savedTheme === 'system') {
      return savedTheme
    }
    return null
  } catch {
    return null
  }
}

export function getThemePreference(): ThemePreference {
  if (typeof window === 'undefined') return 'system'
  return readStoredTheme() ?? 'system'
}

function prefersDarkMode(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

function syncDocumentSurface(isDark: boolean) {
  const root = document.documentElement
  const body = document.body
  const backgroundColor = isDark ? DARK_BACKGROUND : LIGHT_BACKGROUND
  const textColor = isDark ? DARK_TEXT : LIGHT_TEXT
  const colorScheme = isDark ? 'dark' : 'light'

  root.classList.toggle('dark', isDark)
  root.style.backgroundColor = backgroundColor
  root.style.backgroundImage = 'none'
  root.style.colorScheme = colorScheme
  root.style.color = textColor
  root.style.opacity = '1'
  root.style.filter = 'none'
  root.style.boxShadow = 'none'
  root.style.textShadow = 'none'
  root.style.height = ''

  if (!body) return

  body.style.backgroundColor = backgroundColor
  body.style.backgroundImage = 'none'
  body.style.colorScheme = colorScheme
  body.style.color = textColor
  body.style.opacity = '1'
  body.style.filter = 'none'
  body.style.boxShadow = 'none'
  body.style.textShadow = 'none'
}

function broadcastThemeChange(isDark: boolean) {
  window.dispatchEvent(new CustomEvent<boolean>(THEME_EVENT, { detail: isDark }))
}

export function resolveThemeIsDark(): boolean {
  if (typeof window === 'undefined') return false

  const savedTheme = readStoredTheme()
  return savedTheme === 'dark' || ((!savedTheme || savedTheme === 'system') && prefersDarkMode())
}

export function applyTheme(isDark: boolean, persist: false | ThemePreference = false): boolean {
  if (typeof document === 'undefined' || typeof window === 'undefined') return isDark

  syncDocumentSurface(isDark)

  if (persist) {
    try {
      window.localStorage.setItem(THEME_KEY, persist)
    } catch {
      // Ignore storage failures and still honor the in-memory theme choice.
    }
    window.dispatchEvent(new CustomEvent<ThemePreference>(THEME_PREFERENCE_EVENT, { detail: persist }))
  }

  broadcastThemeChange(isDark)
  return isDark
}

export function initTheme(): boolean {
  return applyTheme(resolveThemeIsDark())
}

export function setThemePreference(preference: ThemePreference): boolean {
  const isDark = preference === 'dark' || (preference === 'system' && prefersDarkMode())
  return applyTheme(isDark, preference)
}

export function setTheme(isDark: boolean): boolean {
  return applyTheme(isDark, isDark ? 'dark' : 'light')
}

export function toggleTheme(current?: boolean): boolean {
  const nextTheme = typeof current === 'boolean' ? !current : !resolveThemeIsDark()
  return setTheme(nextTheme)
}

export function useThemeState() {
  const isDark = ref(resolveThemeIsDark())
  let cleanup = () => {}

  onMounted(() => {
    isDark.value = document.documentElement.classList.contains('dark')

    const handleThemeChange = (event: Event) => {
      const customEvent = event as CustomEvent<boolean>
      isDark.value = typeof customEvent.detail === 'boolean'
        ? customEvent.detail
        : document.documentElement.classList.contains('dark')
    }

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const handleSystemThemeChange = () => {
      if (getThemePreference() === 'system') {
        isDark.value = applyTheme(media.matches)
      }
    }

    window.addEventListener(THEME_EVENT, handleThemeChange as EventListener)
    media.addEventListener?.('change', handleSystemThemeChange)
    cleanup = () => {
      window.removeEventListener(THEME_EVENT, handleThemeChange as EventListener)
      media.removeEventListener?.('change', handleSystemThemeChange)
    }
  })

  onBeforeUnmount(() => {
    cleanup()
  })

  return isDark
}

export function useThemePreference() {
  const preference = ref<ThemePreference>(getThemePreference())
  let cleanup = () => {}

  onMounted(() => {
    preference.value = getThemePreference()

    const handlePreferenceChange = (event: Event) => {
      const customEvent = event as CustomEvent<ThemePreference>
      preference.value = customEvent.detail ?? getThemePreference()
    }

    window.addEventListener(THEME_PREFERENCE_EVENT, handlePreferenceChange as EventListener)
    cleanup = () => window.removeEventListener(THEME_PREFERENCE_EVENT, handlePreferenceChange as EventListener)
  })

  onBeforeUnmount(() => {
    cleanup()
  })

  return preference
}
