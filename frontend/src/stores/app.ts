import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    pending: 0,
    toast: '' as string
  }),
  getters: {
    loading: (state) => state.pending > 0
  },
  actions: {
    startLoading() {
      this.pending += 1
    },
    stopLoading() {
      this.pending = Math.max(0, this.pending - 1)
    },
    notify(message: string) {
      this.toast = message
      window.clearTimeout((this as { toastTimer?: number }).toastTimer)
      ;(this as { toastTimer?: number }).toastTimer = window.setTimeout(() => {
        this.toast = ''
      }, 3200)
    }
  }
})
