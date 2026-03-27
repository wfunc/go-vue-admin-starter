<template>
  <div class="space-y-6">
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
      <article v-for="card in cards" :key="card.label" class="stat-card">
        <p class="stat-card-label">{{ card.label }}</p>
        <p class="stat-card-value">{{ card.value }}</p>
        <p class="stat-card-description">{{ card.description }}</p>
      </article>
    </section>

    <PageSection :eyebrow="t('dashboard.eyebrow')" :title="t('dashboard.recentTitle')" :description="t('dashboard.recentDescription')">
      <div class="table-card">
        <table class="table-shell">
          <thead>
            <tr>
              <th>{{ t('dashboard.columns.username') }}</th>
              <th>{{ t('dashboard.columns.module') }}</th>
              <th>{{ t('dashboard.columns.action') }}</th>
              <th>{{ t('dashboard.columns.path') }}</th>
              <th>{{ t('dashboard.columns.time') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in recentLogs" :key="item.id">
              <td>{{ item.username }}</td>
              <td>{{ item.module }}</td>
              <td>{{ item.action }}</td>
              <td class="font-mono text-xs text-gray-500">{{ item.path }}</td>
              <td>{{ formatDate(item.created_at) }}</td>
            </tr>
            <tr v-if="!recentLogs.length">
              <td colspan="5" class="py-10 text-center text-sm text-gray-500">{{ t('dashboard.emptyLogs') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import { listAuditLogs } from '@/api/modules/audit'
import { getDashboardSummary } from '@/api/modules/system'
import { formatDate } from '@/utils/format'
import type { AuditLog, DashboardSummary } from '@/types'

const summary = ref<DashboardSummary>({
  user_count: 0,
  role_count: 0,
  menu_count: 0,
  config_count: 0,
  audit_log_count: 0
})
const recentLogs = ref<AuditLog[]>([])
const { t } = useI18n()

const cards = computed(() => [
  { label: t('dashboard.cards.usersLabel'), value: summary.value.user_count, description: t('dashboard.cards.usersDesc') },
  { label: t('dashboard.cards.rolesLabel'), value: summary.value.role_count, description: t('dashboard.cards.rolesDesc') },
  { label: t('dashboard.cards.menusLabel'), value: summary.value.menu_count, description: t('dashboard.cards.menusDesc') },
  { label: t('dashboard.cards.configsLabel'), value: summary.value.config_count, description: t('dashboard.cards.configsDesc') },
  { label: t('dashboard.cards.logsLabel'), value: summary.value.audit_log_count, description: t('dashboard.cards.logsDesc') }
])

onMounted(async () => {
  summary.value = await getDashboardSummary()
  recentLogs.value = (await listAuditLogs({ page: 1, page_size: 5 })).items
})
</script>
