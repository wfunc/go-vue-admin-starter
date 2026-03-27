<template>
  <span :class="badgeClass" class="inline-flex rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-wide">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  label: string
  tone?: 'success' | 'warning' | 'neutral'
}>()

const badgeClass = computed(() => {
  if (props.tone === 'success') return 'bg-emerald-100 text-emerald-700'
  if (props.tone === 'warning') return 'bg-amber-100 text-amber-700'
  return 'bg-gray-100 text-gray-600'
})
</script>
