<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('audit.eyebrow')" :title="t('audit.title')" :description="t('audit.description')">
      <div class="table-card">
        <table class="table-shell">
          <thead>
            <tr>
              <th>{{ t('audit.columns.username') }}</th>
              <th>{{ t('audit.columns.module') }}</th>
              <th>{{ t('audit.columns.action') }}</th>
              <th>{{ t('audit.columns.request') }}</th>
              <th>{{ t('audit.columns.status') }}</th>
              <th>{{ t('audit.columns.time') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in logs" :key="item.id">
              <td>{{ item.username }}</td>
              <td>{{ item.module }}</td>
              <td>{{ item.action }}</td>
              <td>
                <p class="font-mono text-xs text-gray-600">{{ item.method }} {{ item.path }}</p>
                <p class="mt-1 text-xs text-gray-500">{{ item.ip }}</p>
              </td>
              <td>
                <StatusBadge :label="String(item.status_code)" :tone="item.status_code >= 400 ? 'warning' : 'success'" />
              </td>
              <td>{{ formatDate(item.created_at) }}</td>
            </tr>
            <tr v-if="!logs.length">
              <td colspan="6" class="py-10 text-center text-sm text-gray-500">{{ t('audit.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { listAuditLogs } from '@/api/modules/audit'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'
import { formatDate } from '@/utils/format'
import type { AuditLog } from '@/types'

const logs = ref<AuditLog[]>([])
const { t } = useI18n()

onMounted(async () => {
  logs.value = (await listAuditLogs({ page: 1, page_size: 100 })).items
})
</script>
