<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('roleManagement.eyebrow')" :title="t('roleManagement.title')" :description="t('roleManagement.description')">
      <template #actions>
        <button class="btn-secondary" @click="resetForm">{{ t('roleManagement.newRole') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[340px_minmax(0,1fr)]">
        <form class="form-panel space-y-3" @submit.prevent="submitForm">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('roleManagement.name') }}</label>
            <input v-model="form.name" class="ui-input" :placeholder="t('roleManagement.name')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('roleManagement.code') }}</label>
            <input v-model="form.code" :disabled="Boolean(form.id)" class="ui-input" :placeholder="t('roleManagement.code')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('roleManagement.descriptionLabel') }}</label>
            <textarea v-model="form.description" rows="3" class="ui-textarea" :placeholder="t('roleManagement.descriptionLabel')"></textarea>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('roleManagement.permissions') }}</label>
            <input v-model="permissionsInput" class="ui-input" :placeholder="t('roleManagement.permissionsPlaceholder')" />
          </div>
          <label class="ui-checkbox-row">
            <input v-model="form.is_system" type="checkbox" class="h-4 w-4" />
            {{ t('roleManagement.systemRole') }}
          </label>
          <div class="flex gap-3 pt-2">
            <button :disabled="running" class="btn-primary flex-1">{{ running ? t('common.submitting') : form.id ? t('roleManagement.updateRole') : t('roleManagement.createRole') }}</button>
            <button type="button" class="btn-secondary" @click="resetForm">{{ t('common.reset') }}</button>
          </div>
        </form>

        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('roleManagement.name') }}</th>
                <th>{{ t('roleManagement.permissionCount') }}</th>
                <th>{{ t('roleManagement.type') }}</th>
                <th>{{ t('roleManagement.updatedAt') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in roles" :key="item.id">
                <td>
                  <p class="font-medium text-gray-900">{{ item.name }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.code }}</p>
                </td>
                <td>{{ item.permissions.length }}</td>
                <td>
                  <StatusBadge :label="item.is_system ? t('roleManagement.system') : t('roleManagement.custom')" :tone="item.is_system ? 'success' : 'neutral'" />
                </td>
                <td>{{ formatDate(item.updated_at) }}</td>
                <td>
                  <div class="flex justify-end gap-4">
                    <button class="action-link" @click="editRole(item)">{{ t('common.edit') }}</button>
                    <button class="btn-danger-text" :disabled="item.is_system" @click="removeRole(item)">{{ t('common.delete') }}</button>
                  </div>
                </td>
              </tr>
              <tr v-if="!roles.length">
                <td colspan="5" class="py-10 text-center text-sm text-gray-500">{{ t('roleManagement.empty') }}</td>
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
import { createRole, deleteRole, listRoles, updateRole } from '@/api/modules/roles'
import { useAsyncAction } from '@/composables/useAsyncAction'
import { formatDate } from '@/utils/format'
import type { Role } from '@/types'

const { running, run } = useAsyncAction()
const roles = ref<Role[]>([])
const permissionsInput = ref('')
const { t } = useI18n()
const form = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  is_system: false
})

async function loadData() {
  roles.value = (await listRoles({ page: 1, page_size: 100 })).items
}

function resetForm() {
  form.id = 0
  form.name = ''
  form.code = ''
  form.description = ''
  form.is_system = false
  permissionsInput.value = ''
}

function editRole(item: Role) {
  form.id = item.id
  form.name = item.name
  form.code = item.code
  form.description = item.description
  form.is_system = item.is_system
  permissionsInput.value = item.permissions.join(', ')
}

function getPermissions() {
  return permissionsInput.value
    .split(',')
    .map((value) => value.trim())
    .filter(Boolean)
}

async function submitForm() {
  const payload = {
    name: form.name,
    code: form.code,
    description: form.description,
    permissions: getPermissions(),
    is_system: form.is_system
  }

  await run(async () => {
    if (form.id) {
      await updateRole(form.id, payload)
    } else {
      await createRole(payload)
    }
    await loadData()
    resetForm()
  })
}

async function removeRole(item: Role) {
  if (item.is_system) {
    return
  }
  await deleteRole(item.id)
  await loadData()
}

onMounted(async () => {
  await loadData()
})
</script>
