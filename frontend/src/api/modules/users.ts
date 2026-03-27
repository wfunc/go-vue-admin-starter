import http, { unwrap } from '@/api/http'
import type { PaginatedData, User } from '@/types'

export function listUsers(params: Record<string, unknown>) {
  return unwrap<PaginatedData<User>>(http.get('/users', { params }))
}

export function createUser(payload: Record<string, unknown>) {
  return unwrap<User>(http.post('/users', payload))
}

export function updateUser(id: number, payload: Record<string, unknown>) {
  return unwrap<User>(http.put(`/users/${id}`, payload))
}

export function deleteUser(id: number) {
  return unwrap<{ message: string }>(http.delete(`/users/${id}`))
}
