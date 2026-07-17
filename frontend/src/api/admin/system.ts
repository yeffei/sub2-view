/**
 * System API endpoints for admin operations
 */

import { apiClient } from '../client'

export interface ReleaseInfo {
  name: string
  body: string
  published_at: string
  html_url: string
}

export interface VersionInfo {
  current_version: string
  latest_version: string
  has_update: boolean
  release_info?: ReleaseInfo
  cached: boolean
  warning?: string
  build_type: string // "source" for manual builds, "release" for CI builds
  release_repo: string
}

export type UpdatePreflightStatus = 'pass' | 'warn' | 'fail'

export interface UpdatePreflightCheck {
  key: string
  label: string
  status: UpdatePreflightStatus
  message: string
}

export interface UpdatePreflightInfo {
  current_version: string
  latest_version: string
  has_update: boolean
  can_update: boolean
  build_type: string
  release_repo: string
  archive_name: string
  download_asset?: {
    name: string
    download_url: string
    size: number
  }
  checksum_asset?: {
    name: string
    download_url: string
    size: number
  }
  executable_path?: string
  backup_path?: string
  checks: UpdatePreflightCheck[]
  blocking_reasons?: string[]
  warnings?: string[]
}

/**
 * Get current version
 */
export async function getVersion(): Promise<{ version: string }> {
  const { data } = await apiClient.get<{ version: string }>('/admin/system/version')
  return data
}

/**
 * Check for updates
 * @param force - Force refresh from GitHub API
 */
export async function checkUpdates(force = false): Promise<VersionInfo> {
  const { data } = await apiClient.get<VersionInfo>('/admin/system/check-updates', {
    params: force ? { force: 'true' } : undefined
  })
  return data
}

export async function checkUpdatePreflight(force = false): Promise<UpdatePreflightInfo> {
  const { data } = await apiClient.get<UpdatePreflightInfo>('/admin/system/update/preflight', {
    params: force ? { force: 'true' } : undefined
  })
  return data
}

export interface UpdateResult {
  message: string
  need_restart?: boolean
  already_up_to_date?: boolean
  current_version?: string
  latest_version?: string
  operation_id?: string
}

export interface RollbackVersionInfo {
  version: string
  published_at: string
  html_url: string
}

/**
 * Get versions available for rollback (up to 3 versions older than current)
 */
export async function getRollbackVersions(): Promise<{ versions: RollbackVersionInfo[] }> {
  const { data } = await apiClient.get<{ versions: RollbackVersionInfo[] }>(
    '/admin/system/rollback-versions'
  )
  return data
}

/**
 * Perform system update
 * Downloads and applies the latest version
 */
export async function performUpdate(idempotencyKey?: string): Promise<UpdateResult> {
  const { data } = await apiClient.post<UpdateResult>('/admin/system/update', undefined, {
    headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : undefined,
    timeout: 10 * 60 * 1000
  })
  return data
}

/**
 * Rollback to a previous version
 * @param version - Target version (e.g. "0.1.146"); omit to restore the local backup binary.
 *   A non-version first argument remains supported as the legacy SST idempotency key.
 * @param idempotencyKey - Optional idempotency key for a versioned rollback.
 */
export async function rollback(version?: string, idempotencyKey?: string): Promise<UpdateResult> {
  const looksLikeVersion = typeof version === 'string' && /^v?\d+\.\d+\.\d+(?:[-+].+)?$/.test(version)
  const targetVersion = looksLikeVersion ? version : undefined
  const legacyIdempotencyKey = looksLikeVersion ? undefined : version
  const body = targetVersion ? { version: targetVersion } : undefined
  const requestIdempotencyKey = idempotencyKey || legacyIdempotencyKey
  const response = requestIdempotencyKey
    ? await apiClient.post<UpdateResult>('/admin/system/rollback', body, {
        headers: { 'Idempotency-Key': requestIdempotencyKey },
        timeout: 5 * 60 * 1000
      })
    : await apiClient.post<UpdateResult>('/admin/system/rollback', body)
  const { data } = response
  return data
}

/**
 * Restart the service
 */
export async function restartService(idempotencyKey?: string): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>('/admin/system/restart', undefined, {
    headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : undefined,
    timeout: 60 * 1000
  })
  return data
}

export const systemAPI = {
  getVersion,
  checkUpdates,
  checkUpdatePreflight,
  performUpdate,
  getRollbackVersions,
  rollback,
  restartService
}

export default systemAPI
