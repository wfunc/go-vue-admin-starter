<template>
  <header class="topbar-shell">
    <div>
      <p class="section-eyebrow">{{ section }}</p>
      <h1 class="mt-1 text-2xl font-semibold text-gray-900">{{ title }}</h1>
    </div>
    <div class="flex flex-wrap items-center justify-end gap-3">
      <RouterLink to="/profile" class="topbar-chip transition hover:border-primary-200 hover:bg-primary-50" :class="route.path === '/profile' ? 'border-primary-200 bg-primary-50' : ''">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-primary-600 text-xs font-semibold uppercase text-white">
            {{ initials }}
          </div>
          <div class="text-left">
            <p class="text-sm font-semibold text-gray-900">{{ username }}</p>
            <p class="text-xs text-gray-500">{{ roleName }}</p>
          </div>
        </div>
      </RouterLink>
      <div class="inline-flex rounded-xl border border-gray-200 bg-white p-1 text-sm shadow-sm">
        <button
          v-for="item in localeOptions"
          :key="item.value"
          type="button"
          class="rounded-lg px-3 py-1.5 font-medium transition"
          :class="locale === item.value ? 'bg-primary-600 text-white' : 'text-gray-600 hover:bg-gray-100'"
          @click="switchLocale(item.value)"
        >
          {{ item.label }}
        </button>
      </div>
      <RouterLink to="/profile" class="btn-secondary">{{ t('topbar.profile') }}</RouterLink>
      <button class="btn-secondary" @click="$emit('logout')">{{ t('topbar.logout') }}</button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { setLocale, type AppLocale } from '@/i18n'

const props = defineProps<{
  section: string
  title: string
  username: string
  roleName: string
}>()

defineEmits<{ logout: [] }>()

const route = useRoute()
const { t, locale } = useI18n()
const initials = computed(() => props.username.slice(0, 2).toUpperCase())
const localeOptions = computed(() => [
  { value: 'zh' as AppLocale, label: t('topbar.langZh') },
  { value: 'en' as AppLocale, label: t('topbar.langEn') }
])

function switchLocale(value: AppLocale) {
  setLocale(value)
}
</script>
