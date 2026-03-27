<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('systemConfig.eyebrow')" :title="t('systemConfig.title')" :description="t('systemConfig.description')">
      <template #actions>
        <button class="btn-secondary" @click="resetForm">{{ t('systemConfig.newConfig') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <form class="form-panel space-y-3" @submit.prevent="submitForm">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('systemConfig.key') }}</label>
            <input v-model="form.key" :disabled="Boolean(form.id)" class="ui-input" :placeholder="t('systemConfig.keyPlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('systemConfig.value') }}</label>
            <input v-model="form.value" class="ui-input" :placeholder="t('systemConfig.valuePlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('systemConfig.category') }}</label>
            <input v-model="form.category" class="ui-input" :placeholder="t('systemConfig.categoryPlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('systemConfig.descriptionLabel') }}</label>
            <textarea v-model="form.description" rows="3" class="ui-textarea" :placeholder="t('systemConfig.descriptionPlaceholder')"></textarea>
          </div>
          <label class="ui-checkbox-row">
            <input v-model="form.is_public" type="checkbox" class="h-4 w-4" />
            {{ t('systemConfig.publicRead') }}
          </label>
          <div class="flex gap-3 pt-2">
            <button :disabled="running" class="btn-primary flex-1">{{ running ? t('common.submitting') : form.id ? t('systemConfig.updateConfig') : t('systemConfig.createConfig') }}</button>
            <button type="button" class="btn-secondary" @click="resetForm">{{ t('common.reset') }}</button>
          </div>
        </form>

        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('systemConfig.key') }}</th>
                <th>{{ t('systemConfig.value') }}</th>
                <th>{{ t('systemConfig.category') }}</th>
                <th>{{ t('systemConfig.public') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in configs" :key="item.id">
                <td>
                  <p class="font-medium text-gray-900">{{ item.key }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.description || '-' }}</p>
                </td>
                <td class="font-mono text-xs text-gray-500">{{ item.value || '-' }}</td>
                <td>{{ item.category || '-' }}</td>
                <td>
                  <StatusBadge :label="item.is_public ? t('systemConfig.public') : t('systemConfig.private')" :tone="item.is_public ? 'success' : 'neutral'" />
                </td>
                <td>
                  <div class="flex justify-end gap-4">
                    <button class="action-link" @click="editConfig(item)">{{ t('common.edit') }}</button>
                    <button class="btn-danger-text" @click="removeConfig(item.id)">{{ t('common.delete') }}</button>
                  </div>
                </td>
              </tr>
              <tr v-if="!configs.length">
                <td colspan="5" class="py-10 text-center text-sm text-gray-500">{{ t('systemConfig.empty') }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'
import { createSystemConfig, deleteSystemConfig, listSystemConfigs, updateSystemConfig } from '@/api/modules/system'
import { useAsyncAction } from '@/composables/useAsyncAction'
import type { SystemConfigItem } from '@/types'

const { running, run } = useAsyncAction()
const configs = ref<SystemConfigItem[]>([])
const { t } = useI18n()
const form = reactive({
  id: 0,
  key: '',
  value: '',
  category: '',
  description: '',
  is_public: false
})

async function loadData() {
  configs.value = (await listSystemConfigs({ page: 1, page_size: 100 })).items
}

function resetForm() {
  form.id = 0
  form.key = ''
  form.value = ''
  form.category = ''
  form.description = ''
  form.is_public = false
}

function editConfig(item: SystemConfigItem) {
  form.id = item.id
  form.key = item.key
  form.value = item.value
  form.category = item.category
  form.description = item.description
  form.is_public = item.is_public
}

async function submitForm() {
  const payload = {
    key: form.key,
    value: form.value,
    category: form.category,
    description: form.description,
    is_public: form.is_public
  }

  await run(async () => {
    if (form.id) {
      await updateSystemConfig(form.id, payload)
    } else {
      await createSystemConfig(payload)
    }
    await loadData()
    resetForm()
  })
}

async function removeConfig(id: number) {
  await deleteSystemConfig(id)
  await loadData()
}

onMounted(async () => {
  await loadData()
})
</script>
