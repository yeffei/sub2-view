import { apiClient } from './client'
import type { MonitorStatus } from './admin/channelMonitor'

export type { MonitorStatus } from './admin/channelMonitor'

export interface PoolHealthTimelinePoint {
  status: MonitorStatus
  latency_ms: number | null
  ping_latency_ms: number | null
  checked_at: string
}

export interface PoolHealthView {
  id: number
  name: string
  status: MonitorStatus
  capacity_status: 'ample' | 'observe' | 'tight' | 'queueing'
  availability_7d: number
  best_latency_ms: number | null
  best_ping_latency_ms: number | null
  timeline: PoolHealthTimelinePoint[]
}

export interface PoolHealthDetail {
  id: number
  name: string
  status: MonitorStatus
  capacity_status: 'ample' | 'observe' | 'tight' | 'queueing'
  availability_7d: number
  availability_15d: number
  availability_30d: number
  best_latency_ms: number | null
  best_ping_latency_ms: number | null
  timeline: PoolHealthTimelinePoint[]
}

export interface PoolHealthListResponse {
  items: PoolHealthView[]
}

export async function list(options?: { signal?: AbortSignal }): Promise<PoolHealthListResponse> {
  const { data } = await apiClient.get<PoolHealthListResponse>('/upstream-pools/health', {
    signal: options?.signal,
  })
  return data
}

export async function status(id: number): Promise<PoolHealthDetail> {
  const { data } = await apiClient.get<PoolHealthDetail>(`/upstream-pools/${id}/health`)
  return data
}

export const poolHealthAPI = {
  list,
  status,
}

export default poolHealthAPI
