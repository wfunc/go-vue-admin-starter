<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('customerService.quickReplies.eyebrow')" :title="t('customerService.quickReplies.title')" :description="t('customerService.quickReplies.description')">
      <template #actions>
        <button class="btn-secondary">{{ t('customerService.quickReplies.newReply') }}</button>
      </template>

      <div class="grid gap-4 xl:grid-cols-2">
        <article v-for="item in replies" :key="item.title" class="section-card">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-lg font-semibold text-gray-900">{{ item.title }}</p>
              <p class="mt-1 text-sm text-gray-500">{{ item.category }}</p>
            </div>
            <StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'enabled' ? 'success' : 'neutral'" />
          </div>
          <p class="mt-4 rounded-2xl border border-gray-200 bg-gray-50 p-4 text-sm leading-6 text-gray-600">{{ item.content }}</p>
          <div class="mt-4 flex justify-end">
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

const replies = [
  { title: 'Refund request acknowledgement', category: 'Billing', status: 'enabled', content: 'We have received your refund request and assigned it to the billing queue. Our team will update you after the payment channel confirms the original transaction.' },
  { title: 'Integration troubleshooting follow-up', category: 'Technical Support', status: 'enabled', content: 'To speed up diagnosis, please share your request timestamp, environment, and the exact response payload. We will compare it with the server-side audit log.' },
  { title: 'Migration onboarding confirmation', category: 'Onboarding', status: 'disabled', content: 'Your migration checklist is ready. We can schedule a guided session once your source data fields are confirmed.' }
]
</script>
