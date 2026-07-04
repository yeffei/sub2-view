/**
 * Admin upstream pool API endpoints
 * Handles upstream pool routing metadata for administrators
 */

import { apiClient } from '../client'
import type { UpstreamPool, UpstreamPoolMember, UpstreamPoolBinding } from '@/types'

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

export async function getBindings(): Promise<UpstreamPoolBinding[]> {
  const { data } = await apiClient.get<ApiListResponse<UpstreamPoolBinding> | UpstreamPoolBinding[]>(
    '/admin/upstream-pools/bindings'
  )
  return Array.isArray(data) ? data : data.items || []
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

const upstreamPoolsAPI = { list, getById, getMembers, getBindings, create, update, remove, createMember, updateMember, removeMember, createBinding, updateBinding, removeBinding }
export default upstreamPoolsAPI
