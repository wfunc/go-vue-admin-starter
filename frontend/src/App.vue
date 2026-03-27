<template>
  <div class="min-h-screen text-ink">
    <div v-if="appStore.loading" class="fixed inset-x-0 top-0 z-50 h-1 overflow-hidden bg-transparent">
      <div class="h-full w-1/3 animate-[loading_1.1s_ease-in-out_infinite] rounded-full bg-primary-500"></div>
    </div>
    <RouterView />
    <transition name="toast">
      <div v-if="appStore.toast" class="fixed bottom-6 right-6 z-50 max-w-sm rounded-2xl bg-dark-900 px-4 py-3 text-sm text-white shadow-panel">
        {{ appStore.toast }}
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
</script>
