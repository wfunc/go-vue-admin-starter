<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('customerService.knowledge.eyebrow')" :title="t('customerService.knowledge.title')" :description="t('customerService.knowledge.description')">
      <template #actions>
        <button class="btn-secondary">{{ t('customerService.knowledge.newArticle') }}</button>
      </template>

      <div class="grid gap-4 xl:grid-cols-3">
        <article v-for="item in articles" :key="item.title" class="section-card">
          <div class="flex items-start justify-between gap-4">
            <div>
              <p class="text-lg font-semibold text-gray-900">{{ item.title }}</p>
              <p class="mt-2 text-sm leading-6 text-gray-500">{{ item.summary }}</p>
            </div>
            <StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'published' ? 'success' : 'neutral'" />
          </div>
          <div class="mt-5 flex items-center justify-between text-sm text-gray-500">
            <span>{{ item.author }}</span>
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

const articles = [
  { title: 'Billing callback troubleshooting SOP', summary: 'Standard process for tracing payment notifications, validating signatures, and confirming retries.', status: 'published', author: 'Operations Team' },
  { title: 'Customer data import checklist', summary: 'A reusable onboarding checklist for mapping legacy customer fields into the platform starter.', status: 'published', author: 'Enablement Team' },
  { title: 'Escalation and severity policy', summary: 'Defines when a ticket should move from frontline handling to engineering or management review.', status: 'draft', author: 'Support Lead' }
]
</script>
