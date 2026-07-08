import { beforeEach, describe, expect, it, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

const {
  appStore,
  authStore,
  checkUpdatePreflight,
  performUpdate,
  restartService,
  rollback,
} = vi.hoisted(() => ({
  appStore: {
    versionLoading: false,
    currentVersion: '0.1.137',
    latestVersion: '0.1.143',
    hasUpdate: true,
    releaseRepo: 'yeffei/sub2-view',
    releaseInfo: { name: 'v0.1.143', body: '', published_at: '', html_url: 'https://github.com/Wei-Shaw/sub2api/releases/tag/v0.1.143' },
    buildType: 'release',
    fetchVersion: vi.fn(),
    clearVersionCache: vi.fn(),
    showInfo: vi.fn(),
    showSuccess: vi.fn(),
  },
  authStore: {
    isAdmin: true,
  },
  checkUpdatePreflight: vi.fn(),
  performUpdate: vi.fn(),
  restartService: vi.fn(),
  rollback: vi.fn(),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) => {
      if (!params) return key
      return `${key} ${JSON.stringify(params)}`
    },
  }),
}))

vi.mock('@/stores', () => ({
  useAppStore: () => appStore,
  useAuthStore: () => authStore,
}))

vi.mock('@/api/admin/system', () => ({
  checkUpdatePreflight,
  performUpdate,
  restartService,
  rollback,
}))

import VersionBadge from '@/components/common/VersionBadge.vue'

function resetStore() {
  appStore.versionLoading = false
  appStore.currentVersion = '0.1.137'
  appStore.latestVersion = '0.1.143'
  appStore.hasUpdate = true
  appStore.releaseRepo = 'yeffei/sub2-view'
  appStore.releaseInfo = {
    name: 'v0.1.143',
    body: '',
    published_at: '',
    html_url: 'https://github.com/Wei-Shaw/sub2api/releases/tag/v0.1.143',
  }
  appStore.buildType = 'release'
  appStore.fetchVersion.mockResolvedValue(null)
  appStore.clearVersionCache.mockReset()
  appStore.showInfo.mockReset()
  appStore.showSuccess.mockReset()
  authStore.isAdmin = true
  checkUpdatePreflight.mockReset()
  performUpdate.mockReset()
  restartService.mockReset()
  rollback.mockReset()
}

async function openDropdown() {
  const wrapper = mount(VersionBadge)
  await wrapper.find('button').trigger('click')
  await flushPromises()
  return wrapper
}

function findButton(wrapper: ReturnType<typeof mount>, text: string) {
  const button = wrapper.findAll('button').find((item) => item.text().includes(text))
  expect(button, `button containing ${text}`).toBeTruthy()
  return button!
}

describe('VersionBadge release update flow', () => {
  beforeEach(() => {
    resetStore()
    vi.spyOn(globalThis.crypto, 'randomUUID').mockReturnValue('00000000-0000-4000-8000-000000000001')
  })

  it('blocks update when preflight fails', async () => {
    checkUpdatePreflight.mockResolvedValue({
      current_version: '0.1.137',
      latest_version: '0.1.143',
      has_update: true,
      can_update: false,
      build_type: 'source',
      archive_name: 'windows_amd64',
      checks: [],
      blocking_reasons: ['source build must be upgraded with git/worktree workflow'],
      warnings: [],
    })

    const wrapper = await openDropdown()
    await findButton(wrapper, 'version.updateNow').trigger('click')
    await flushPromises()

    expect(checkUpdatePreflight).toHaveBeenCalledWith(true)
    expect(performUpdate).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('source build must be upgraded with git/worktree workflow')
    expect(wrapper.text()).toContain('yeffei/sub2-view')
  })

  it('runs update only after preflight passes', async () => {
    checkUpdatePreflight.mockResolvedValue({
      current_version: '0.1.137',
      latest_version: '0.1.143',
      has_update: true,
      can_update: true,
      build_type: 'release',
      archive_name: 'windows_amd64',
      checks: [],
      warnings: [],
    })
    performUpdate.mockResolvedValue({
      message: 'Update completed. Please restart the service.',
      need_restart: true,
    })

    const wrapper = await openDropdown()
    await findButton(wrapper, 'version.updateNow').trigger('click')
    await flushPromises()

    expect(performUpdate).toHaveBeenCalledWith('system-update-00000000-0000-4000-8000-000000000001')
    expect(appStore.clearVersionCache).toHaveBeenCalled()
    expect(wrapper.text()).toContain('version.restartNow')
    expect(wrapper.text()).toContain('version.rollbackNow')
  })

  it('restarts and rolls back with idempotency keys after update', async () => {
    checkUpdatePreflight.mockResolvedValue({
      current_version: '0.1.137',
      latest_version: '0.1.143',
      has_update: true,
      can_update: true,
      build_type: 'release',
      archive_name: 'windows_amd64',
      checks: [],
      warnings: [],
    })
    performUpdate.mockResolvedValue({ message: 'Update completed.', need_restart: true })
    restartService.mockResolvedValue({ message: 'Service restart initiated' })
    rollback.mockResolvedValue({ message: 'Rollback completed.', need_restart: true })
    vi.spyOn(window, 'confirm').mockReturnValue(true)

    const wrapper = await openDropdown()
    await findButton(wrapper, 'version.updateNow').trigger('click')
    await flushPromises()

    await findButton(wrapper, 'version.rollbackNow').trigger('click')
    await flushPromises()
    expect(rollback).toHaveBeenCalledWith('system-rollback-00000000-0000-4000-8000-000000000001')

    await findButton(wrapper, 'version.restartNow').trigger('click')
    await flushPromises()
    expect(restartService).toHaveBeenCalledWith('system-restart-00000000-0000-4000-8000-000000000001')
  })
})
