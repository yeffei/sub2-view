/**
 * Admin upstream pool API endpoints
 * Handles upstream pool routing metadata for administrators
 */

import { apiClient } from '../client'
import type {
  UpstreamPool,
  UpstreamPoolMember,
  UpstreamPoolBinding,
  UpstreamAccountSet,
  UpstreamAccountSetMember,
  UpstreamPoolMemberSet,
  UpstreamPoolMemberSyncMode,
  UpstreamPoolMemberSyncResult
} from '@/types'

interface ApiListResponse<T> {
  items: T[]
  total?: number
  page?: number
  page_size?: number
  pages?: number
}

export async function list(): Promise<UpstreamPool[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamPool> | UpstreamPool[]>('/admin/upstream-pools')
  return Array.isArray(data) ? data : data.items || []
}

export async function getById(id: number): Promise<UpstreamPool> {
  const { data } = await apiClient.get<UpstreamPool>(`/admin/upstream-pools/${id}`)
  return data
}

export async function getMembers(id: number): Promise<UpstreamPoolMember[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamPoolMember> | UpstreamPoolMember[]>(
    `/admin/upstream-pools/${id}/members`
  )
  return Array.isArray(data) ? data : data.items || []
}

export async function previewMemberSync(id: number, payload: {
  mode: UpstreamPoolMemberSyncMode
}): Promise<UpstreamPoolMemberSyncResult> {
  const { data } = await apiClient.post<UpstreamPoolMemberSyncResult>(
    `/admin/upstream-pools/${id}/member-sync/preview`,
    payload
  )
  return data
}

export async function applyMemberSync(id: number, payload: {
  mode: UpstreamPoolMemberSyncMode
}): Promise<UpstreamPoolMemberSyncResult> {
  const { data } = await apiClient.post<UpstreamPoolMemberSyncResult>(
    `/admin/upstream-pools/${id}/member-sync/apply`,
    payload
  )
  return data
}

export async function getBindings(): Promise<UpstreamPoolBinding[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamPoolBinding> | UpstreamPoolBinding[]>(
    '/admin/upstream-pools/bindings'
  )
  return Array.isArray(data) ? data : data.items || []
}

export async function getAccountSets(): Promise<UpstreamAccountSet[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamAccountSet> | UpstreamAccountSet[]>(
    '/admin/upstream-pools/account-sets'
  )
  return Array.isArray(data) ? data : data.items || []
}

export async function createAccountSet(payload: {
  name: string
  platform: string
  description?: string
  enabled?: boolean
  code?: string
}): Promise<UpstreamAccountSet> {
  const { data } = await apiClient.post<UpstreamAccountSet>('/admin/upstream-pools/account-sets', payload)
  return data
}

export async function updateAccountSet(setId: number, payload: {
  name?: string
  code?: string
  platform?: string
  description?: string
  enabled?: boolean
}): Promise<UpstreamAccountSet> {
  const { data } = await apiClient.put<UpstreamAccountSet>(`/admin/upstream-pools/account-sets/${setId}`, payload)
  return data
}

export async function removeAccountSet(setId: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/upstream-pools/account-sets/${setId}`)
  return data
}

export async function getAccountSetMembers(setId: number): Promise<UpstreamAccountSetMember[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamAccountSetMember> | UpstreamAccountSetMember[]>(
    `/admin/upstream-pools/account-sets/${setId}/members`
  )
  return Array.isArray(data) ? data : data.items || []
}

export async function addAccountSetMembers(setId: number, payload: {
  account_ids: number[]
}): Promise<{ message: string }> {
  const { data } = await apiClient.post<{ message: string }>(
    `/admin/upstream-pools/account-sets/${setId}/members`,
    payload
  )
  return data
}

export async function removeAccountSetMember(setId: number, accountId: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(
    `/admin/upstream-pools/account-sets/${setId}/members/${accountId}`
  )
  return data
}

export async function getMemberSets(poolId: number): Promise<UpstreamPoolMemberSet[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamPoolMemberSet> | UpstreamPoolMemberSet[]>(
    `/admin/upstream-pools/${poolId}/member-sets`
  )
  return Array.isArray(data) ? data : data.items || []
}

export async function createMemberSet(poolId: number, payload: {
  set_id: number
  enabled?: boolean
  notes?: string | null
}): Promise<UpstreamPoolMemberSet> {
  const { data } = await apiClient.post<UpstreamPoolMemberSet>(`/admin/upstream-pools/${poolId}/member-sets`, payload)
  return data
}

export async function updateMemberSet(memberSetId: number, payload: {
  enabled?: boolean | null
  notes?: string | null
}): Promise<UpstreamPoolMemberSet> {
  const { data } = await apiClient.put<UpstreamPoolMemberSet>(
    `/admin/upstream-pools/member-sets/${memberSetId}`,
    payload
  )
  return data
}

export async function removeMemberSet(memberSetId: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/upstream-pools/member-sets/${memberSetId}`)
  return data
}

export async function create(payload: Partial<UpstreamPool> & { name: string; code: string; platform: string }): Promise<UpstreamPool> {
  const { data } = await apiClient.post<UpstreamPool>('/admin/upstream-pools', payload)
  return data
}

export async function update(id: number, payload: Partial<UpstreamPool>): Promise<UpstreamPool> {
  const { data } = await apiClient.put<UpstreamPool>(`/admin/upstream-pools/${id}`, payload)
  return data
}

export async function remove(id: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/upstream-pools/${id}`)
  return data
}

export async function createMember(poolId: number, payload: {
  account_id: number
  enabled?: boolean
  schedulable_override?: boolean | null
  manual_drained?: boolean
  weight?: number
  priority_override?: number | null
  max_concurrency_override?: number | null
  notes?: string | null
}): Promise<UpstreamPoolMember> {
  const { data } = await apiClient.post<UpstreamPoolMember>(`/admin/upstream-pools/${poolId}/members`, payload)
  return data
}

export async function updateMember(memberId: number, payload: {
  enabled?: boolean | null
  schedulable_override?: boolean | null
  manual_drained?: boolean | null
  weight?: number | null
  priority_override?: number | null
  max_concurrency_override?: number | null
  notes?: string | null
}): Promise<UpstreamPoolMember> {
  const { data } = await apiClient.put<UpstreamPoolMember>(`/admin/upstream-pools/members/${memberId}`, payload)
  return data
}

export async function removeMember(memberId: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/upstream-pools/members/${memberId}`)
  return data
}

export async function createBinding(payload: {
  group_id: number
  pool_id: number
  platform: string
  models?: string[]
  request_path_scope?: string[]
  priority?: number
  enabled?: boolean
}): Promise<UpstreamPoolBinding> {
  const { data } = await apiClient.post<UpstreamPoolBinding>('/admin/upstream-pools/bindings', payload)
  return data
}

export async function updateBinding(bindingId: number, payload: {
  group_id?: number | null
  pool_id?: number | null
  platform?: string | null
  models?: string[]
  request_path_scope?: string[]
  priority?: number | null
  enabled?: boolean | null
}): Promise<UpstreamPoolBinding> {
  const { data } = await apiClient.put<UpstreamPoolBinding>(`/admin/upstream-pools/bindings/${bindingId}`, payload)
  return data
}

export async function removeBinding(bindingId: number): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(`/admin/upstream-pools/bindings/${bindingId}`)
  return data
}

const upstreamPoolsAPI = {
  list,
  getById,
  getMembers,
  previewMemberSync,
  applyMemberSync,
  getBindings,
  getAccountSets,
  getAccountSetMembers,
  getMemberSets,
  create,
  update,
  remove,
  createMember,
  updateMember,
  removeMember,
  createBinding,
  updateBinding,
  removeBinding,
  createAccountSet,
  updateAccountSet,
  removeAccountSet,
  addAccountSetMembers,
  removeAccountSetMember,
  createMemberSet,
  updateMemberSet,
  removeMemberSet
}
export default upstreamPoolsAPI
