<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('menuManagement.eyebrow')" :title="t('menuManagement.title')" :description="t('menuManagement.description')">
      <template #actions>
        <button class="btn-secondary" @click="resetForm">{{ t('menuManagement.newMenu') }}</button>
      </template>

      <div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
        <form class="form-panel space-y-3" @submit.prevent="submitForm">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.titleLabel') }}</label>
            <input v-model="form.title" class="ui-input" :placeholder="t('menuManagement.titleLabel')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.name') }}</label>
            <input v-model="form.name" :disabled="Boolean(form.id)" class="ui-input" :placeholder="t('menuManagement.name')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.path') }}</label>
            <input v-model="form.path" class="ui-input" :placeholder="t('menuManagement.pathPlaceholder')" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.component') }}</label>
            <input v-model="form.component" class="ui-input" :placeholder="t('menuManagement.componentPlaceholder')" />
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.icon') }}</label>
              <input v-model="form.icon" class="ui-input" :placeholder="t('menuManagement.iconPlaceholder')" />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.sort') }}</label>
              <input v-model.number="form.sort" type="number" class="ui-input" :placeholder="t('menuManagement.sortPlaceholder')" />
            </div>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.parent') }}</label>
            <select v-model="parentValue" class="ui-select">
              <option value="">{{ t('menuManagement.noParent') }}</option>
              <option v-for="option in parentOptions" :key="option.id" :value="String(option.id)">{{ displayTitle(option) }}</option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.permission') }}</label>
            <input v-model="form.permission" class="ui-input" :placeholder="t('menuManagement.permissionPlaceholder')" />
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('menuManagement.type') }}</label>
              <select v-model="form.menu_type" class="ui-select">
                <option value="menu">{{ t('menuManagement.menuTypeMenu') }}</option>
                <option value="button">{{ t('menuManagement.menuTypeButton') }}</option>
              </select>
            </div>
            <div class="flex items-end">
              <label class="ui-checkbox-row w-full">
                <input v-model="form.hidden" type="checkbox" class="h-4 w-4" />
                {{ t('menuManagement.hidden') }}
              </label>
            </div>
          </div>
          <div class="flex gap-3 pt-2">
            <button :disabled="running" class="btn-primary flex-1">{{ running ? t('common.submitting') : form.id ? t('menuManagement.updateMenu') : t('menuManagement.createMenu') }}</button>
            <button type="button" class="btn-secondary" @click="resetForm">{{ t('common.reset') }}</button>
          </div>
        </form>

        <div class="table-card">
          <table class="table-shell">
            <thead>
              <tr>
                <th>{{ t('menuManagement.titleLabel') }}</th>
                <th>{{ t('menuManagement.path') }}</th>
                <th>{{ t('menuManagement.permission') }}</th>
                <th>{{ t('menuManagement.sort') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="item in menus" :key="item.id">
                <tr>
                  <td>
                    <p class="font-medium text-gray-900">{{ displayTitle(item) }}</p>
                    <p class="mt-1 text-xs text-gray-500">{{ item.name }}</p>
                  </td>
                  <td class="font-mono text-xs text-gray-500">{{ item.path }}</td>
                  <td>{{ item.permission || '-' }}</td>
                  <td>{{ item.sort }}</td>
                  <td>
                    <div class="flex justify-end gap-4">
                      <button class="action-link" @click="editMenu(item)">{{ t('common.edit') }}</button>
                      <button class="btn-danger-text" @click="removeMenu(item.id)">{{ t('common.delete') }}</button>
                    </div>
                  </td>
                </tr>
                <tr v-for="child in item.children || []" :key="child.id" class="bg-gray-50/80">
                  <td class="pl-10">
                    <p class="font-medium text-gray-900">{{ displayTitle(child) }}</p>
                    <p class="mt-1 text-xs text-gray-500">{{ child.name }}</p>
                  </td>
                  <td class="font-mono text-xs text-gray-500">{{ child.path }}</td>
                  <td>{{ child.permission || '-' }}</td>
                  <td>{{ child.sort }}</td>
                  <td>
                    <div class="flex justify-end gap-4">
                      <button class="action-link" @click="editMenu(child)">{{ t('common.edit') }}</button>
                      <button class="btn-danger-text" @click="removeMenu(child.id)">{{ t('common.delete') }}</button>
                    </div>
                  </td>
                </tr>
              </template>
              <tr v-if="!menus.length">
                <td colspan="5" class="py-10 text-center text-sm text-gray-500">{{ t('menuManagement.empty') }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import { createMenu, deleteMenu, listMenus, updateMenu } from '@/api/modules/menus'
import { useAsyncAction } from '@/composables/useAsyncAction'
import { resolveMenuTitle } from '@/i18n'
import { useAuthStore } from '@/stores/auth'
import type { MenuItem } from '@/types'

const { running, run } = useAsyncAction()
const menus = ref<MenuItem[]>([])
const parentValue = ref('')
const { t } = useI18n()
const authStore = useAuthStore()
const form = reactive({
  id: 0,
  title: '',
  name: '',
  path: '',
  component: '',
  icon: '',
  menu_type: 'menu',
  permission: '',
  sort: 10,
  hidden: false
})

const parentOptions = computed(() => menus.value.filter((item) => !item.parent_id))

async function loadData() {
  menus.value = await listMenus({})
}

function resetForm() {
  form.id = 0
  form.title = ''
  form.name = ''
  form.path = ''
  form.component = ''
  form.icon = ''
  form.menu_type = 'menu'
  form.permission = ''
  form.sort = 10
  form.hidden = false
  parentValue.value = ''
}

function editMenu(item: MenuItem) {
  form.id = item.id
  form.title = item.title
  form.name = item.name
  form.path = item.path
  form.component = item.component
  form.icon = item.icon
  form.menu_type = item.menu_type
  form.permission = item.permission
  form.sort = item.sort
  form.hidden = item.hidden
  parentValue.value = item.parent_id ? String(item.parent_id) : ''
}

function displayTitle(item: MenuItem) {
  return resolveMenuTitle(item.name, item.path, item.title)
}

async function submitForm() {
  const payload = {
    title: form.title,
    name: form.name,
    path: form.path,
    component: form.component,
    icon: form.icon,
    menu_type: form.menu_type,
    permission: form.permission,
    sort: form.sort,
    hidden: form.hidden,
    parent_id: parentValue.value ? Number(parentValue.value) : null
  }

  await run(async () => {
    if (form.id) {
      await updateMenu(form.id, payload)
    } else {
      await createMenu(payload)
    }
    await Promise.all([loadData(), authStore.fetchMenus()])
    resetForm()
  })
}

async function removeMenu(id: number) {
  await deleteMenu(id)
  await Promise.all([loadData(), authStore.fetchMenus()])
}

onMounted(async () => {
  await loadData()
})
</script>
