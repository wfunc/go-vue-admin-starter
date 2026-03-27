<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('customerService.tickets.eyebrow')" :title="t('customerService.tickets.title')" :description="t('customerService.tickets.description')">
      <template #actions>
        <button class="btn-secondary">{{ t('customerService.tickets.createTicket') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_320px]">
        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('customerService.tickets.columns.no') }}</th>
                <th>{{ t('customerService.tickets.columns.subject') }}</th>
                <th>{{ t('customerService.tickets.columns.customer') }}</th>
                <th>{{ t('customerService.tickets.columns.priority') }}</th>
                <th>{{ t('customerService.tickets.columns.status') }}</th>
                <th>{{ t('customerService.tickets.columns.assignee') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in tickets" :key="item.no">
                <td class="font-mono text-xs text-gray-500">{{ item.no }}</td>
                <td>
                  <p class="font-medium text-gray-900">{{ item.subject }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.category }}</p>
                </td>
                <td>{{ item.customer }}</td>
                <td><StatusBadge :label="t(`customerService.status.${item.priority}`)" :tone="item.priority === 'urgent' || item.priority === 'high' ? 'warning' : 'neutral'" /></td>
                <td><StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'resolved' ? 'success' : item.status === 'pending' ? 'warning' : 'neutral'" /></td>
                <td>{{ item.assignee }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <section class="space-y-4">
          <div class="section-card">
            <p class="section-eyebrow">Queue</p>
            <h3 class="mt-1 text-lg font-semibold text-gray-900">Workload Snapshot</h3>
            <div class="mt-4 space-y-3">
              <div v-for="item in queue" :key="item.name" class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                <div class="flex items-center justify-between gap-4">
                  <div>
                    <p class="text-sm font-semibold text-gray-900">{{ item.name }}</p>
                    <p class="mt-1 text-sm text-gray-500">{{ item.meta }}</p>
                  </div>
                  <span class="text-2xl font-semibold text-gray-900">{{ item.value }}</span>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'

const { t } = useI18n()

const tickets = [
  { no: 'CS-2026-031', subject: 'Payment callback failed after renewal', customer: 'BlueRiver Retail', priority: 'urgent', status: 'processing', assignee: 'Lena', category: 'Billing' },
  { no: 'CS-2026-028', subject: 'Invoice format needs additional company field', customer: 'Northwind Labs', priority: 'high', status: 'pending', assignee: 'Mason', category: 'Finance' },
  { no: 'CS-2026-024', subject: 'Need migration guidance for legacy data import', customer: 'Aster Health', priority: 'medium', status: 'open', assignee: 'Unassigned', category: 'Onboarding' },
  { no: 'CS-2026-017', subject: 'Login occasionally redirects to expired callback', customer: 'Lattice Care', priority: 'low', status: 'resolved', assignee: 'Mia', category: 'Authentication' }
]

const queue = [
  { name: 'Awaiting first reply', meta: '3 tickets within 30 minutes', value: '3' },
  { name: 'Pending customer reply', meta: 'Waiting on external confirmation', value: '6' },
  { name: 'Escalated to engineering', meta: 'Cross-team follow-up required', value: '2' }
]
</script>
