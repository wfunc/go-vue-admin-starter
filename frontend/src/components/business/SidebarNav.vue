<template>
  <div class="nav-shell px-4 py-5">
    <div class="border-b border-gray-200 pb-4">
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary-600 text-sm font-bold text-white">
          GS
        </div>
        <div>
          <p class="text-xs uppercase tracking-[0.32em] text-gray-400">{{ t('app.shortName') }}</p>
          <h1 class="text-lg font-semibold text-gray-900">{{ t('app.name') }}</h1>
          <p class="mt-1 text-sm text-gray-500">{{ t('app.description') }}</p>
        </div>
      </div>
    </div>

    <nav class="mt-4 space-y-2 overflow-y-auto pb-4">
      <div v-for="item in menus" :key="item.id" class="space-y-2">
        <button
          v-if="hasChildren(item)"
          type="button"
          class="nav-link w-full justify-between text-left"
          :class="isActive(item) ? 'nav-link-active' : ''"
          @click="toggleExpand(item.id)"
        >
          <div class="flex items-center gap-3">
            <span class="nav-icon">{{ iconText(item) }}</span>
            <span>{{ menuTitle(item) }}</span>
          </div>
          <div class="flex items-center gap-3">
            <span class="text-[11px] uppercase tracking-[0.16em] text-gray-400">{{ item.name }}</span>
            <span class="text-xs text-gray-400">{{ isExpanded(item) ? 'v' : '>' }}</span>
          </div>
        </button>

        <RouterLink
          v-else
          :to="resolveTarget(item)"
          class="nav-link justify-between"
          :class="isActive(item) ? 'nav-link-active' : ''"
        >
          <div class="flex items-center gap-3">
            <span class="nav-icon">{{ iconText(item) }}</span>
            <span>{{ menuTitle(item) }}</span>
          </div>
          <span class="text-[11px] uppercase tracking-[0.16em] text-gray-400">{{ item.name }}</span>
        </RouterLink>

        <div v-if="hasChildren(item) && isExpanded(item)" class="space-y-1 pl-4">
          <RouterLink
            v-for="child in item.children"
            :key="child.id"
            :to="child.path"
            class="nav-link py-2 text-sm"
            :class="isActive(child) ? 'nav-link-active' : ''"
          >
            <span class="nav-icon h-8 w-8 text-[11px]">{{ iconText(child) }}</span>
            <span>{{ menuTitle(child) }}</span>
          </RouterLink>
        </div>
      </div>
    </nav>

    <div class="mt-auto rounded-2xl border border-gray-200 bg-gray-50 p-4 text-sm text-gray-500">
      <p class="font-medium text-gray-700">{{ t('sidebar.demoAccountTitle') }}</p>
      <p class="mt-1">{{ t('sidebar.demoAccountValue') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { resolveMenuTitle } from '@/i18n'
import type { MenuItem } from '@/types'

const props = defineProps<{ menus: MenuItem[] }>()
const route = useRoute()
const { t } = useI18n()
const expandedIds = ref<number[]>([])

function resolveTarget(item: MenuItem) {
  if (item.children?.length) {
    return item.children[0].path
  }
  return item.path
}

function iconText(item: MenuItem) {
  return (item.icon || item.title || 'M').slice(0, 2).toUpperCase()
}

function menuTitle(item: MenuItem) {
  return resolveMenuTitle(item.name, item.path, item.title)
}

function hasChildren(item: MenuItem) {
  return Boolean(item.children?.length)
}

function isActive(item: MenuItem) {
  if (route.path === item.path) {
    return true
  }
  return item.children?.some((child) => route.path === child.path) ?? false
}

function isExpanded(item: MenuItem) {
  return isActive(item) || expandedIds.value.includes(item.id)
}

function toggleExpand(id: number) {
  expandedIds.value = expandedIds.value.includes(id)
    ? expandedIds.value.filter((itemId) => itemId !== id)
    : [...expandedIds.value, id]
}

watch(
  () => [props.menus, route.path],
  () => {
    const activeParentIds = props.menus.filter((item) => hasChildren(item) && isActive(item)).map((item) => item.id)
    expandedIds.value = Array.from(new Set([...expandedIds.value, ...activeParentIds]))
  },
  { immediate: true, deep: true }
)
</script>
