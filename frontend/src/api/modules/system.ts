import http, { unwrap } from '@/api/http'
import type { DashboardSummary, PaginatedData, SystemConfigItem } from '@/types'

export function listSystemConfigs(params: Record<string, unknown>) {
  return unwrap<PaginatedData<SystemConfigItem>>(http.get('/system-configs', { params }))
}

export function createSystemConfig(payload: Record<string, unknown>) {
  return unwrap<SystemConfigItem>(http.post('/system-configs', payload))
}

export function updateSystemConfig(id: number, payload: Record<string, unknown>) {
  return unwrap<SystemConfigItem>(http.put(`/system-configs/${id}`, payload))
}

export function deleteSystemConfig(id: number) {
  return unwrap<{ message: string }>(http.delete(`/system-configs/${id}`))
}

export function getDashboardSummary() {
  return unwrap<DashboardSummary>(http.get('/system-configs/summary'))
}

export function getPublicSettings() {
  return unwrap<Record<string, string>>(http.get('/system-configs/public'))
}
