import http, { unwrap } from '@/api/http'
import type { ChangePasswordPayload, MenuItem, Profile, Session, UpdateProfilePayload } from '@/types'

export function login(payload: { username: string; password: string }) {
  return unwrap<Session>(http.post('/auth/login', payload))
}

export function getProfile() {
  return unwrap<Profile>(http.get('/auth/me'))
}

export function updateProfile(payload: UpdateProfilePayload) {
  return unwrap<Profile>(http.put('/auth/profile', payload))
}

export function changePassword(payload: ChangePasswordPayload) {
  return unwrap<{ message: string }>(http.post('/auth/change-password', payload))
}

export function getMenus() {
  return unwrap<MenuItem[]>(http.get('/auth/menus'))
}

export function logout() {
  return unwrap<{ message: string }>(http.post('/auth/logout'))
}
