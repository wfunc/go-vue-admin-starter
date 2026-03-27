<template>
  <div class="space-y-6">
    <PageSection :eyebrow="t('profile.eyebrow')" :title="t('profile.title')" :description="t('profile.description')">
      <div class="grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
        <div class="space-y-6">
          <section class="section-card">
            <div class="mb-5 flex items-center gap-4 border-b border-gray-200 pb-4">
              <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary-600 text-lg font-semibold uppercase text-white">
                {{ initials }}
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900">{{ authStore.user?.nickname || authStore.user?.username }}</h2>
                <p class="mt-1 text-sm text-gray-500">{{ authStore.user?.role_name }} · {{ authStore.user?.status }}</p>
              </div>
            </div>

            <dl class="grid gap-4 md:grid-cols-2">
              <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                <dt class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">{{ t('profile.username') }}</dt>
                <dd class="mt-2 text-sm font-medium text-gray-900">{{ authStore.user?.username }}</dd>
              </div>
              <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                <dt class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">{{ t('profile.role') }}</dt>
                <dd class="mt-2 text-sm font-medium text-gray-900">{{ authStore.user?.role_name }}</dd>
              </div>
              <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                <dt class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">{{ t('profile.status') }}</dt>
                <dd class="mt-2 text-sm font-medium text-gray-900">{{ statusLabel(authStore.user?.status) }}</dd>
              </div>
              <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
                <dt class="text-xs font-semibold uppercase tracking-[0.18em] text-gray-500">{{ t('profile.permissionCount') }}</dt>
                <dd class="mt-2 text-sm font-medium text-gray-900">{{ authStore.user?.permissions.length || 0 }}</dd>
              </div>
            </dl>
          </section>

          <section class="section-card">
            <div class="mb-5 border-b border-gray-200 pb-4">
              <p class="section-eyebrow">{{ t('profile.profileEyebrow') }}</p>
              <h2 class="mt-1 text-xl font-semibold text-gray-900">{{ t('profile.profileTitle') }}</h2>
              <p class="mt-1 text-sm text-gray-500">{{ t('profile.profileDescription') }}</p>
            </div>

            <form class="space-y-4" @submit.prevent="submitProfile">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('profile.nickname') }}</label>
                <input v-model="profileForm.nickname" class="ui-input" :placeholder="t('profile.nicknamePlaceholder')" />
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('profile.email') }}</label>
                <input v-model="profileForm.email" class="ui-input" :placeholder="t('profile.emailPlaceholder')" />
              </div>
              <div class="flex gap-3 pt-2">
                <button :disabled="profileSubmitting" class="btn-primary">{{ profileSubmitting ? t('profile.savingProfile') : t('profile.saveProfile') }}</button>
                <button type="button" class="btn-secondary" @click="resetProfileForm">{{ t('common.reset') }}</button>
              </div>
            </form>
          </section>
        </div>

        <section class="section-card">
          <div class="mb-5 border-b border-gray-200 pb-4">
            <p class="section-eyebrow">{{ t('profile.securityEyebrow') }}</p>
            <h2 class="mt-1 text-xl font-semibold text-gray-900">{{ t('profile.securityTitle') }}</h2>
            <p class="mt-1 text-sm text-gray-500">{{ t('profile.securityDescription') }}</p>
          </div>

          <form class="space-y-4" @submit.prevent="submitPassword">
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('profile.currentPassword') }}</label>
              <input v-model="passwordForm.current_password" type="password" class="ui-input" :placeholder="t('profile.currentPasswordPlaceholder')" autocomplete="current-password" />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('profile.newPassword') }}</label>
              <input v-model="passwordForm.new_password" type="password" class="ui-input" :placeholder="t('profile.newPasswordPlaceholder')" autocomplete="new-password" />
            </div>
            <div>
              <label class="mb-1 block text-sm font-medium text-gray-700">{{ t('profile.confirmPassword') }}</label>
              <input v-model="passwordForm.confirm_password" type="password" class="ui-input" :placeholder="t('profile.confirmPasswordPlaceholder')" autocomplete="new-password" />
            </div>
            <div class="flex gap-3 pt-2">
              <button :disabled="passwordSubmitting" class="btn-primary">{{ passwordSubmitting ? t('profile.updatingPassword') : t('profile.updatePassword') }}</button>
              <button type="button" class="btn-secondary" @click="resetPasswordForm">{{ t('common.clear') }}</button>
            </div>
          </form>
        </section>
      </div>
    </PageSection>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import PageSection from '@/components/base/PageSection.vue'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const authStore = useAuthStore()
const appStore = useAppStore()
const router = useRouter()
const { t } = useI18n()

const profileSubmitting = ref(false)
const passwordSubmitting = ref(false)

const profileForm = reactive({
  nickname: authStore.user?.nickname || '',
  email: authStore.user?.email || ''
})

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const initials = computed(() => (authStore.user?.nickname || authStore.user?.username || 'AD').slice(0, 2).toUpperCase())

watch(
  () => authStore.user,
  () => {
    resetProfileForm()
  },
  { immediate: true }
)

function resetProfileForm() {
  profileForm.nickname = authStore.user?.nickname || ''
  profileForm.email = authStore.user?.email || ''
}

function resetPasswordForm() {
  passwordForm.current_password = ''
  passwordForm.new_password = ''
  passwordForm.confirm_password = ''
}

function statusLabel(status?: string) {
  if (status === 'active') {
    return t('userManagement.active')
  }
  if (status === 'disabled') {
    return t('userManagement.disabled')
  }
  return status || t('common.unknown')
}

async function submitProfile() {
  profileSubmitting.value = true
  try {
    await authStore.updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email
    })
    appStore.notify(t('profile.profileUpdated'))
  } finally {
    profileSubmitting.value = false
  }
}

async function submitPassword() {
  if (passwordForm.new_password.length < 8) {
    appStore.notify(t('profile.passwordTooShort'))
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    appStore.notify(t('profile.passwordMismatch'))
    return
  }

  passwordSubmitting.value = true
  try {
    await authStore.changePassword({
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password,
      confirm_password: passwordForm.confirm_password
    })
    authStore.logout()
    appStore.notify(t('profile.passwordUpdated'))
    resetPasswordForm()
    await router.push({ name: 'login' })
  } finally {
    passwordSubmitting.value = false
  }
}
</script>
