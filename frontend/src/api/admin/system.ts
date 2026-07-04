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
 * Rollback to previous version
 */
export async function rollback(idempotencyKey?: string): Promise<UpdateResult> {
  const { data } = await apiClient.post<UpdateResult>('/admin/system/rollback', undefined, {
    headers: idempotencyKey ? { 'Idempotency-Key': idempotencyKey } : undefined,
    timeout: 5 * 60 * 1000
  })
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
  rollback,
  restartService
}

export default systemAPI
