import { createI18n } from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

export const LOCALE_STORAGE_KEY = 'go-vue-admin-locale'
export const supportedLocales = ['zh', 'en'] as const

export type AppLocale = (typeof supportedLocales)[number]

function normalizeLocale(value?: string | null): AppLocale {
  if (!value) {
    return 'zh'
  }

  const lower = value.toLowerCase()
  if (lower.startsWith('en')) {
    return 'en'
  }
  return 'zh'
}

export function getStoredLocale(): AppLocale {
  if (typeof window === 'undefined') {
    return 'zh'
  }

  const saved = window.localStorage.getItem(LOCALE_STORAGE_KEY)
  if (saved) {
    return normalizeLocale(saved)
  }

  return normalizeLocale(window.navigator.language)
}

const i18n = createI18n({
  legacy: false,
  locale: getStoredLocale(),
  fallbackLocale: 'en',
  messages: { zh, en }
})

export function setLocale(locale: AppLocale) {
  i18n.global.locale.value = locale
  if (typeof window !== 'undefined') {
    window.localStorage.setItem(LOCALE_STORAGE_KEY, locale)
  }
}

export function resolveRouteTitle(titleKey?: string) {
  if (!titleKey) {
    return i18n.global.t('app.name')
  }
  return i18n.global.t(titleKey)
}

export function syncDocumentTitle(titleKey?: string) {
  document.title = `${resolveRouteTitle(titleKey)} | ${i18n.global.t('app.documentSuffix')}`
}

const menuKeyMap: Record<string, string> = {
  dashboard: 'nav.dashboard',
  'customer-service': 'nav.customerService',
  'cs-dashboard': 'nav.csDashboard',
  conversations: 'nav.conversations',
  customers: 'nav.customers',
  tickets: 'nav.tickets',
  'ticket-categories': 'nav.ticketCategories',
  'quick-replies': 'nav.quickReplies',
  'knowledge-articles': 'nav.knowledgeArticles',
  system: 'nav.system',
  users: 'nav.users',
  roles: 'nav.roles',
  menus: 'nav.menus',
  'system-configs': 'nav.systemConfigs',
  'audit-logs': 'nav.auditLogs',
  profile: 'nav.profile',
  '/dashboard': 'nav.dashboard',
  '/customer-service': 'nav.customerService',
  '/customer-service/dashboard': 'nav.csDashboard',
  '/customer-service/conversations': 'nav.conversations',
  '/customer-service/customers': 'nav.customers',
  '/customer-service/tickets': 'nav.tickets',
  '/customer-service/categories': 'nav.ticketCategories',
  '/customer-service/quick-replies': 'nav.quickReplies',
  '/customer-service/knowledge': 'nav.knowledgeArticles',
  '/roles': 'nav.roles',
  '/menus': 'nav.menus',
  '/users': 'nav.users',
  '/system-configs': 'nav.systemConfigs',
  '/audit-logs': 'nav.auditLogs',
  '/profile': 'nav.profile'
}

export function resolveMenuTitle(name?: string, path?: string, fallback?: string) {
  const titleKey = menuKeyMap[name || ''] || menuKeyMap[path || '']
  if (titleKey && i18n.global.te(titleKey)) {
    return i18n.global.t(titleKey)
  }
  return fallback || name || path || ''
}

export default i18n
