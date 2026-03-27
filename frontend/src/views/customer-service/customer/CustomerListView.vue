<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('customerService.customers.eyebrow')" :title="t('customerService.customers.title')" :description="t('customerService.customers.description')">
      <template #actions>
        <button class="btn-secondary">{{ t('customerService.customers.newCustomer') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <section class="form-panel space-y-4">
          <div class="rounded-2xl border border-gray-200 bg-white p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Filters</p>
            <div class="mt-4 space-y-3">
              <input class="ui-input" placeholder="Search customer / email / company" />
              <select class="ui-select">
                <option>All sources</option>
                <option>Website</option>
                <option>Email</option>
                <option>Partner referral</option>
              </select>
              <select class="ui-select">
                <option>All status</option>
                <option>{{ t('customerService.status.active') }}</option>
                <option>{{ t('customerService.status.vip') }}</option>
                <option>{{ t('customerService.status.atRisk') }}</option>
              </select>
            </div>
          </div>

          <div class="rounded-2xl border border-dashed border-primary-200 bg-primary-50 p-4">
            <p class="text-sm font-semibold text-primary-700">Demo Note</p>
            <p class="mt-2 text-sm leading-6 text-primary-700/80">
              This list is the first-phase skeleton for customer profiles. The next iteration will connect ticket history and customer detail drawers.
            </p>
          </div>
        </section>

        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('customerService.customers.columns.name') }}</th>
                <th>{{ t('customerService.customers.columns.source') }}</th>
                <th>{{ t('customerService.customers.columns.owner') }}</th>
                <th>{{ t('customerService.customers.columns.status') }}</th>
                <th>{{ t('customerService.customers.columns.lastTicket') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in customers" :key="item.name">
                <td>
                  <p class="font-medium text-gray-900">{{ item.name }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.contact }}</p>
                </td>
                <td>{{ item.source }}</td>
                <td>{{ item.owner }}</td>
                <td><StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'vip' ? 'success' : item.status === 'atRisk' ? 'warning' : 'neutral'" /></td>
                <td>{{ item.lastTicket }}</td>
                <td>
                  <div class="flex justify-end gap-4">
                    <button class="action-link">{{ t('common.edit') }}</button>
                    <button class="action-link">Timeline</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'

const { t } = useI18n()

const customers = [
  { name: 'BlueRiver Retail', contact: 'ops@blueriver.io', source: 'Website', owner: 'Lena', status: 'vip', lastTicket: 'Payment callback failed after renewal' },
  { name: 'Northwind Labs', contact: 'finance@northwind.dev', source: 'Partner referral', owner: 'Mason', status: 'active', lastTicket: 'Invoice format needs company field' },
  { name: 'Aster Health', contact: 'pm@asterhealth.cn', source: 'Email', owner: 'Luna', status: 'atRisk', lastTicket: 'Data import delayed after onboarding' }
]
</script>
