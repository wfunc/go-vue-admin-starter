import axios from 'axios'
import type { ApiEnvelope } from '@/types'
import i18n from '@/i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000
})

http.interceptors.request.use((config) => {
  const authStore = useAuthStore()
  const appStore = useAppStore()
  appStore.startLoading()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => {
    useAppStore().stopLoading()
    return response
  },
  (error) => {
    const appStore = useAppStore()
    appStore.stopLoading()
    const message = error.response?.data?.message || error.message || i18n.global.t('common.requestFailed')
    appStore.notify(message)
    if (error.response?.status === 401) {
      useAuthStore().logout()
    }
    return Promise.reject(error)
  }
)

export async function unwrap<T>(promise: Promise<{ data: ApiEnvelope<T> }>) {
  const response = await promise
  return response.data.data
}

export default http
