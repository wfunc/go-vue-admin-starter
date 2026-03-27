<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('customerService.categories.eyebrow')" :title="t('customerService.categories.title')" :description="t('customerService.categories.description')">
      <template #actions>
        <button class="btn-secondary">{{ t('customerService.categories.newCategory') }}</button>
      </template>

      <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <article v-for="item in categories" :key="item.name" class="section-card">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-lg font-semibold text-gray-900">{{ item.name }}</p>
              <p class="mt-2 text-sm leading-6 text-gray-500">{{ item.description }}</p>
            </div>
            <StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'enabled' ? 'success' : 'neutral'" />
          </div>
          <div class="mt-5 flex items-center justify-between text-sm text-gray-500">
            <span>{{ item.count }} tickets</span>
            <button class="action-link">{{ t('common.edit') }}</button>
          </div>
        </article>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'

const { t } = useI18n()

const categories = [
  { name: 'Billing', description: 'Renewal, refund, invoice, and payment callback cases.', status: 'enabled', count: 18 },
  { name: 'Technical Support', description: 'API issues, integration troubleshooting, and incident follow-up.', status: 'enabled', count: 32 },
  { name: 'Onboarding', description: 'Migration help, setup guidance, and enablement requests.', status: 'enabled', count: 11 },
  { name: 'Complaints', description: 'Escalation-sensitive feedback and service dissatisfaction tracking.', status: 'disabled', count: 3 }
]
</script>
