<template>
  <div class="min-h-screen bg-gray-50">
    <div class="mx-auto flex min-h-screen max-w-[1600px] flex-col lg:flex-row">
      <aside class="border-b border-gray-200 bg-white lg:sticky lg:top-0 lg:h-screen lg:w-64 lg:flex-none lg:border-b-0 lg:border-r">
        <SidebarNav :menus="menuItems" />
      </aside>
      <div class="flex min-h-screen flex-1 flex-col">
        <div class="px-4 py-4 md:px-6 lg:px-8">
          <AppTopbar
            :section="currentSection"
            :title="currentTitle"
            :username="authStore.user?.nickname || authStore.user?.username || 'Admin'"
            :role-name="authStore.user?.role_name || 'Role'"
            @logout="handleLogout"
          />
          <main class="mt-6">
            <RouterView />
          </main>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute, useRouter } from 'vue-router'
import AppTopbar from '@/components/business/AppTopbar.vue'
import SidebarNav from '@/components/business/SidebarNav.vue'
import { syncDocumentTitle } from '@/i18n'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()

const menuItems = computed(() => authStore.menus.filter((item) => !item.hidden))
const currentTitle = computed(() => t(String(route.meta.titleKey || 'app.name')))
const currentSection = computed(() => {
  if (route.path.startsWith('/customer-service')) return t('section.customerService')
  if (route.path.startsWith('/audit')) return t('section.audit')
  if (route.path.startsWith('/system') || route.path.startsWith('/roles') || route.path.startsWith('/menus')) return t('section.system')
  if (route.path.startsWith('/users')) return t('section.user')
  if (route.path.startsWith('/profile')) return t('section.account')
  return t('section.dashboard')
})

function handleLogout() {
  authStore.logout()
  router.push({ name: 'login' })
}

watch(
  [() => route.fullPath, () => locale.value],
  () => {
    syncDocumentTitle(route.meta.titleKey as string | undefined)
  },
  { immediate: true }
)
</script>
