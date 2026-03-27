import { defineStore } from 'pinia'
import type { ChangePasswordPayload, MenuItem, Profile, Session, UpdateProfilePayload } from '@/types'
import { changePassword as changePasswordApi, getMenus, getProfile, login as loginApi, updateProfile as updateProfileApi } from '@/api/modules/auth'

const TOKEN_KEY = 'starter_token'
const USER_KEY = 'starter_user'
const MENUS_KEY = 'starter_menus'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '' as string,
    user: null as Profile | null,
    menus: [] as MenuItem[],
    bootstrapped: false
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    permissions: (state) => state.user?.permissions ?? []
  },
  actions: {
    hydrate() {
      this.token = localStorage.getItem(TOKEN_KEY) || ''
      const user = localStorage.getItem(USER_KEY)
      const menus = localStorage.getItem(MENUS_KEY)
      this.user = user ? (JSON.parse(user) as Profile) : null
      this.menus = menus ? (JSON.parse(menus) as MenuItem[]) : []
      this.bootstrapped = false
    },
    setSession(session: Session) {
      this.token = session.access_token
      this.user = session.user
      this.menus = session.menus
      this.bootstrapped = true
      localStorage.setItem(TOKEN_KEY, this.token)
      localStorage.setItem(USER_KEY, JSON.stringify(this.user))
      localStorage.setItem(MENUS_KEY, JSON.stringify(this.menus))
    },
    setProfile(profile: Profile) {
      this.user = profile
      localStorage.setItem(USER_KEY, JSON.stringify(profile))
    },
    async login(payload: { username: string; password: string }) {
      const session = await loginApi(payload)
      this.setSession(session)
      return session
    },
    async fetchProfile() {
      const profile = await getProfile()
      this.setProfile(profile)
      return profile
    },
    async updateProfile(payload: UpdateProfilePayload) {
      const profile = await updateProfileApi(payload)
      this.setProfile(profile)
      return profile
    },
    async changePassword(payload: ChangePasswordPayload) {
      return changePasswordApi(payload)
    },
    async fetchMenus() {
      const menus = await getMenus()
      this.menus = menus
      localStorage.setItem(MENUS_KEY, JSON.stringify(menus))
      return menus
    },
    async refreshSession() {
      await Promise.all([this.fetchProfile(), this.fetchMenus()])
      this.bootstrapped = true
    },
    async ensureBootstrapped() {
      if (!this.isAuthenticated || this.bootstrapped) {
        return
      }
      await this.refreshSession()
    },
    logout() {
      this.token = ''
      this.user = null
      this.menus = []
      this.bootstrapped = false
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
      localStorage.removeItem(MENUS_KEY)
    },
    hasPermission(permission?: string) {
      if (!permission) return true
      return this.permissions.includes('*') || this.permissions.includes(permission)
    }
  }
})
