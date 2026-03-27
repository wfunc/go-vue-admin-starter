import http, { unwrap } from '@/api/http'
import type { MenuItem } from '@/types'

export function listMenus(params: Record<string, unknown>) {
  return unwrap<MenuItem[]>(http.get('/menus', { params }))
}

export function createMenu(payload: Record<string, unknown>) {
  return unwrap<MenuItem>(http.post('/menus', payload))
}

export function updateMenu(id: number, payload: Record<string, unknown>) {
  return unwrap<MenuItem>(http.put(`/menus/${id}`, payload))
}

export function deleteMenu(id: number) {
  return unwrap<{ message: string }>(http.delete(`/menus/${id}`))
}
