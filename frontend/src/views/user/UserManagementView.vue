<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('userManagement.eyebrow')" :title="t('userManagement.title')" :description="t('userManagement.description')">
      <template #actions>
        <button class="btn-secondary" @click="resetForm">{{ t('userManagement.newUser') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[340px_minmax(0,1fr)]">
        <form class="form-panel space-y-3" @submit.prevent="submitForm">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.username') }}</label>
            <input v-model="form.username" :disabled="Boolean(form.id)" class="ui-input" :placeholder="t('login.usernamePlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.email') }}</label>
            <input v-model="form.email" class="ui-input" :placeholder="t('profile.emailPlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.nickname') }}</label>
            <input v-model="form.nickname" class="ui-input" :placeholder="t('profile.nicknamePlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.role') }}</label>
            <select v-model.number="form.role_id" class="ui-select">
              <option :value="0">{{ t('userManagement.selectRole') }}</option>
              <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.status') }}</label>
            <select v-model="form.status" class="ui-select">
              <option value="active">{{ t('userManagement.active') }}</option>
              <option value="disabled">{{ t('userManagement.disabled') }}</option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('userManagement.password') }}</label>
            <input v-model="form.password" type="password" class="ui-input" :placeholder="form.id ? t('userManagement.passwordPlaceholderUpdate') : t('userManagement.passwordPlaceholderCreate')" />
          </div>
          <div class="flex gap-3 pt-2">
            <button class="btn-primary flex-1">{{ form.id ? t('userManagement.updateUser') : t('userManagement.createUser') }}</button>
            <button type="button" class="btn-secondary" @click="resetForm">{{ t('common.reset') }}</button>
          </div>
        </form>

        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('userManagement.username') }}</th>
                <th>{{ t('userManagement.role') }}</th>
                <th>{{ t('userManagement.status') }}</th>
                <th>{{ t('userManagement.createdAt') }}</th>
                <th class="text-right">{{ t('userManagement.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in users" :key="item.id">
                <td>
                  <p class="font-medium text-gray-900">{{ item.username }}</p>
                  <p class="mt-1 text-xs text-gray-500">{{ item.email || '-' }}</p>
                </td>
                <td>{{ item.role_name }}</td>
                <td><StatusBadge :label="statusLabel(item.status)" :tone="item.status === 'active' ? 'success' : 'warning'" /></td>
                <td>{{ formatDate(item.created_at) }}</td>
                <td>
                  <div class="flex justify-end gap-4">
                    <button class="action-link" @click="editUser(item)">{{ t('common.edit') }}</button>
                    <button class="btn-danger-text" @click="removeUser(item.id)">{{ t('common.delete') }}</button>
                  </div>
                </td>
              </tr>
              <tr v-if="!users.length">
                <td colspan="5" class="py-10 text-center text-sm text-gray-500">{{ t('userManagement.empty') }}</td>
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
import { createUser, deleteUser, listUsers, updateUser } from '@/api/modules/users'
import { listRoles } from '@/api/modules/roles'
import { formatDate } from '@/utils/format'
import type { Role, User } from '@/types'

const users = ref<User[]>([])
const roles = ref<Role[]>([])
const { t } = useI18n()
const form = reactive({ id: 0, username: '', email: '', nickname: '', role_id: 0, status: 'active', password: '' })

async function loadData() {
  users.value = (await listUsers({ page: 1, page_size: 100 })).items
  roles.value = (await listRoles({ page: 1, page_size: 100 })).items
  if (!form.role_id && roles.value[0]) form.role_id = roles.value[0].id
}

function resetForm() {
  form.id = 0
  form.username = ''
  form.email = ''
  form.nickname = ''
  form.status = 'active'
  form.password = ''
  form.role_id = roles.value[0]?.id || 0
}

function editUser(item: User) {
  form.id = item.id
  form.username = item.username
  form.email = item.email || ''
  form.nickname = item.nickname
  form.status = item.status
  form.role_id = item.role_id
  form.password = ''
}

function statusLabel(status: string) {
  return status === 'active' ? t('userManagement.active') : t('userManagement.disabled')
}

async function submitForm() {
  const payload = { username: form.username, email: form.email, nickname: form.nickname, status: form.status, role_id: form.role_id, password: form.password }
  if (form.id) {
    await updateUser(form.id, payload)
  } else {
    await createUser(payload)
  }
  await loadData()
  resetForm()
}

async function removeUser(id: number) {
  await deleteUser(id)
  await loadData()
}

onMounted(async () => {
  await loadData()
  resetForm()
})
</script>
