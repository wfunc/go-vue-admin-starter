<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-10">
    <div class="empty-state-card w-full">
      <p class="section-eyebrow">403</p>
      <h1 class="mt-3 text-3xl font-semibold text-gray-900">{{ t('exception.forbiddenTitle') }}</h1>
      <p class="mt-3 text-sm leading-7 text-gray-500">{{ t('exception.forbiddenDescription') }}</p>
      <div class="mt-8 flex flex-col justify-center gap-3 sm:flex-row">
        <RouterLink to="/dashboard" class="btn-primary">{{ t('exception.goDashboard') }}</RouterLink>
        <RouterLink to="/roles" class="btn-secondary">{{ t('exception.goRoles') }}</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

const { t } = useI18n()
</script>
