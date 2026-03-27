<template>
  <div class="space-y-6">
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-6">
      <article v-for="card in cards" :key="card.label" class="stat-card">
        <p class="stat-card-label">{{ card.label }}</p>
        <p class="stat-card-value">{{ card.value }}</p>
        <p class="stat-card-description">{{ card.description }}</p>
      </article>
    </section>

    <div class="grid gap-6 xl:grid-cols-[minmax(0,1.25fr)_minmax(0,0.95fr)]">
      <PageSection :eyebrow="t('customerService.dashboard.eyebrow')" :title="t('customerService.dashboard.panels.myQueue')" :description="t('customerService.dashboard.description')">
        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('customerService.tickets.columns.no') }}</th>
                <th>{{ t('customerService.tickets.columns.subject') }}</th>
                <th>{{ t('customerService.tickets.columns.customer') }}</th>
                <th>{{ t('customerService.tickets.columns.priority') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in summary.queue_items" :key="item.id">
                <td class="font-mono text-xs text-gray-500">{{ item.ticket_no }}</td>
                <td>
                  <p class="font-medium text-gray-900">{{ item.subject }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.channel }}</p>
                </td>
                <td>{{ item.customer }}</td>
                <td><StatusBadge :label="t(`customerService.status.${item.priority}`)" :tone="item.priority === 'urgent' || item.priority === 'high' ? 'warning' : 'neutral'" /></td>
              </tr>
              <tr v-if="!summary.queue_items.length">
                <td colspan="4" class="py-10 text-center text-sm text-gray-500">{{ t('common.noData') }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </PageSection>

      <div class="space-y-6">
        <PageSection :eyebrow="t('customerService.dashboard.eyebrow')" :title="t('customerService.dashboard.panels.focus')" :description="t('customerService.dashboard.panels.recent')">
          <div class="space-y-4">
            <div v-for="item in summary.recent_transfers" :key="item.id" class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
              <div class="flex items-start justify-between gap-4">
                <div>
                  <p class="text-sm font-semibold text-gray-900">{{ item.actor }}</p>
                  <p class="mt-1 text-sm text-gray-500">{{ item.summary }}</p>
                </div>
                <StatusBadge :label="formatDate(item.time)" tone="neutral" />
              </div>
            </div>
            <div v-if="!summary.recent_transfers.length" class="rounded-2xl border border-dashed border-gray-200 bg-gray-50 p-6 text-sm text-gray-500">
              {{ t('common.noData') }}
            </div>
          </div>
        </PageSection>

        <PageSection :eyebrow="t('customerService.dashboard.eyebrow')" :title="t('customerService.dashboard.panels.agentStatus')" :description="t('customerService.dashboard.panels.agentDescription')">
          <div class="space-y-3">
            <div v-for="item in summary.agent_statuses" :key="item.name" class="rounded-2xl border border-gray-200 bg-white p-4">
              <div class="flex items-center justify-between gap-4">
                <div>
                  <p class="text-sm font-semibold text-gray-900">{{ item.name }}</p>
                  <p class="mt-1 text-sm text-gray-500">{{ item.queue }}</p>
                </div>
                <div class="text-right">
                  <StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'online' ? 'success' : 'neutral'" />
                  <p class="mt-2 text-xs text-gray-500">{{ t('customerService.dashboard.activeLoad') }}: {{ item.active_load }}</p>
                </div>
              </div>
            </div>
            <div v-if="!summary.agent_statuses.length" class="rounded-2xl border border-dashed border-gray-200 bg-gray-50 p-6 text-sm text-gray-500">
              {{ t('common.noData') }}
            </div>
          </div>
        </PageSection>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'
import { getConversationSummary } from '@/api/modules/conversations'
import { formatDate } from '@/utils/format'
import type { ConversationSummary } from '@/types'

const { t } = useI18n()

const summary = ref<ConversationSummary>({
  waiting_count: 0,
  processing_count: 0,
  sla_risk_count: 0,
  resolved_today: 0,
  online_agent_count: 0,
  transfer_count: 0,
  queue_items: [],
  agent_statuses: [],
  recent_transfers: []
})

const cards = computed(() => [
  { label: t('customerService.dashboard.cards.newTickets'), value: summary.value.waiting_count, description: t('customerService.dashboard.cardDescriptions.waiting') },
  { label: t('customerService.dashboard.cards.pendingTickets'), value: summary.value.processing_count, description: t('customerService.dashboard.cardDescriptions.processing') },
  { label: t('customerService.dashboard.cards.processingTickets'), value: summary.value.online_agent_count, description: t('customerService.dashboard.cardDescriptions.onlineAgents') },
  { label: t('customerService.dashboard.cards.resolvedTickets'), value: summary.value.resolved_today, description: t('customerService.dashboard.cardDescriptions.resolved') },
  { label: t('customerService.dashboard.cards.urgentTickets'), value: summary.value.sla_risk_count, description: t('customerService.dashboard.cardDescriptions.risk') },
  { label: t('customerService.dashboard.cards.satisfaction'), value: summary.value.transfer_count, description: t('customerService.dashboard.cardDescriptions.transfers') }
])

onMounted(async () => {
  summary.value = await getConversationSummary()
})
</script>
