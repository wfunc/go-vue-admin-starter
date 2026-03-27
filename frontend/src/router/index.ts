import { createRouter, createWebHistory } from 'vue-router'
import { syncDocumentTitle } from '@/i18n'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { titleKey: 'routes.login', public: true }
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminLayout.vue'),
      children: [
        {
          path: '',
          redirect: '/dashboard'
        },
        {
          path: 'system',
          redirect: '/roles'
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue'),
          meta: { titleKey: 'routes.dashboard', permission: 'dashboard:view' }
        },
        {
          path: 'customer-service',
          redirect: '/customer-service/dashboard'
        },
        {
          path: 'customer-service/dashboard',
          name: 'cs-dashboard',
          component: () => import('@/views/customer-service/dashboard/CustomerServiceDashboardView.vue'),
          meta: { titleKey: 'routes.csDashboard', permission: 'cs_dashboard:view' }
        },
        {
          path: 'customer-service/conversations',
          name: 'conversations',
          component: () => import('@/views/customer-service/conversation/ConversationCenterView.vue'),
          meta: { titleKey: 'routes.conversations', permission: 'conversation:view' }
        },
        {
          path: 'customer-service/customers',
          name: 'customers',
          component: () => import('@/views/customer-service/customer/CustomerListView.vue'),
          meta: { titleKey: 'routes.customers', permission: 'customer:view' }
        },
        {
          path: 'customer-service/tickets',
          name: 'tickets',
          component: () => import('@/views/customer-service/ticket/TicketListView.vue'),
          meta: { titleKey: 'routes.tickets', permission: 'ticket:view' }
        },
        {
          path: 'customer-service/categories',
          name: 'ticket-categories',
          component: () => import('@/views/customer-service/category/TicketCategoryView.vue'),
          meta: { titleKey: 'routes.ticketCategories', permission: 'ticket_category:view' }
        },
        {
          path: 'customer-service/quick-replies',
          name: 'quick-replies',
          component: () => import('@/views/customer-service/quick-reply/QuickReplyView.vue'),
          meta: { titleKey: 'routes.quickReplies', permission: 'quick_reply:view' }
        },
        {
          path: 'customer-service/knowledge',
          name: 'knowledge-articles',
          component: () => import('@/views/customer-service/knowledge/KnowledgeArticleView.vue'),
          meta: { titleKey: 'routes.knowledgeArticles', permission: 'knowledge:view' }
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/user/ProfileView.vue'),
          meta: { titleKey: 'routes.profile' }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/user/UserManagementView.vue'),
          meta: { titleKey: 'routes.users', permission: 'user:view' }
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('@/views/system/RoleManagementView.vue'),
          meta: { titleKey: 'routes.roles', permission: 'role:view' }
        },
        {
          path: 'menus',
          name: 'menus',
          component: () => import('@/views/system/MenuManagementView.vue'),
          meta: { titleKey: 'routes.menus', permission: 'menu:view' }
        },
        {
          path: 'system-configs',
          name: 'system-configs',
          component: () => import('@/views/system/SystemConfigView.vue'),
          meta: { titleKey: 'routes.systemConfigs', permission: 'system_config:view' }
        },
        {
          path: 'audit-logs',
          name: 'audit-logs',
          component: () => import('@/views/audit/AuditLogView.vue'),
          meta: { titleKey: 'routes.auditLogs', permission: 'audit_log:view' }
        }
      ]
    },
    {
      path: '/forbidden',
      name: 'forbidden',
      component: () => import('@/views/exception/ForbiddenView.vue'),
      meta: { titleKey: 'routes.forbidden', public: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/exception/NotFoundView.vue'),
      meta: { titleKey: 'routes.notFound', public: true }
    }
  ]
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()
  const isPublic = Boolean(to.meta.public)

  if (!isPublic && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (authStore.isAuthenticated) {
    try {
      await authStore.ensureBootstrapped()
    } catch {
      authStore.logout()
      return { name: 'login' }
    }
  }

  const permission = to.meta.permission as string | undefined
  if (!isPublic && permission && !authStore.hasPermission(permission)) {
    return { name: 'forbidden' }
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return { name: 'dashboard' }
  }

  syncDocumentTitle(to.meta.titleKey as string | undefined)
  return true
})

export default router
