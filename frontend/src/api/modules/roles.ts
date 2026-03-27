import http, { unwrap } from '@/api/http'
import type { PaginatedData, Role } from '@/types'

export function listRoles(params: Record<string, unknown>) {
  return unwrap<PaginatedData<Role>>(http.get('/roles', { params }))
}

export function createRole(payload: Record<string, unknown>) {
  return unwrap<Role>(http.post('/roles', payload))
}

export function updateRole(id: number, payload: Record<string, unknown>) {
  return unwrap<Role>(http.put(`/roles/${id}`, payload))
}

export function deleteRole(id: number) {
  return unwrap<{ message: string }>(http.delete(`/roles/${id}`))
}
