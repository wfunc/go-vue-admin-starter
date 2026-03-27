<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-10">
    <div class="grid w-full max-w-5xl gap-6 lg:grid-cols-[1.1fr_0.9fr]">
      <section class="card p-8 lg:p-10">
        <p class="section-eyebrow">{{ t('login.framework') }}</p>
        <h1 class="mt-3 text-3xl font-semibold text-gray-900 lg:text-4xl">{{ t('login.headline') }}</h1>
        <p class="mt-4 max-w-2xl text-sm leading-7 text-gray-500 lg:text-base">
          {{ t('login.description') }}
        </p>
        <div class="mt-8 grid gap-4 md:grid-cols-3">
          <article class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Backend</p>
            <h2 class="mt-2 text-lg font-semibold text-gray-900">{{ t('login.backendTitle') }}</h2>
            <p class="mt-2 text-sm text-gray-500">{{ t('login.backendDesc') }}</p>
          </article>
          <article class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Frontend</p>
            <h2 class="mt-2 text-lg font-semibold text-gray-900">{{ t('login.frontendTitle') }}</h2>
            <p class="mt-2 text-sm text-gray-500">{{ t('login.frontendDesc') }}</p>
          </article>
          <article class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">Default</p>
            <h2 class="mt-2 text-lg font-semibold text-gray-900">{{ t('login.defaultTitle') }}</h2>
            <p class="mt-2 text-sm text-gray-500">{{ t('login.defaultDesc') }}</p>
          </article>
        </div>
      </section>

      <section class="card p-8 lg:p-10">
        <p class="section-eyebrow">{{ t('login.access') }}</p>
        <h2 class="mt-3 text-3xl font-semibold text-gray-900">{{ t('login.title') }}</h2>
        <p class="mt-2 text-sm text-gray-500">{{ t('login.subtitle') }}</p>

        <form class="mt-8 space-y-5" @submit.prevent="handleSubmit">
          <label class="block space-y-2">
            <span class="text-sm font-medium text-gray-700">{{ t('login.username') }}</span>
            <input v-model="form.username" class="ui-input" :placeholder="t('login.usernamePlaceholder')" autocomplete="username" />
          </label>
          <label class="block space-y-2">
            <span class="text-sm font-medium text-gray-700">{{ t('login.password') }}</span>
            <input v-model="form.password" type="password" class="ui-input" :placeholder="t('login.passwordPlaceholder')" autocomplete="current-password" />
          </label>
          <button :disabled="submitting" class="btn-primary w-full justify-center py-3">
            {{ submitting ? t('login.submitting') : t('login.submit') }}
          </button>
        </form>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const submitting = ref(false)
const form = reactive({
  username: 'admin',
  password: 'admin123456'
})

async function handleSubmit() {
  submitting.value = true
  try {
    await authStore.login(form)
    await router.push(String(route.query.redirect || '/dashboard'))
  } finally {
    submitting.value = false
  }
}
</script>
