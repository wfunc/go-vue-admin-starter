<template>
  <div class="space-y-6">
    <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <article v-for="card in cards" :key="card.label" class="stat-card">
        <p class="stat-card-label">{{ card.label }}</p>
        <p class="stat-card-value">{{ card.value }}</p>
        <p class="stat-card-description">{{ card.description }}</p>
      </article>
    </section>

    <div class="grid gap-6 xl:grid-cols-[320px_minmax(0,1fr)_320px]">
      <PageSection
        :eyebrow="t('customerService.conversations.eyebrow')"
        :title="t('customerService.conversations.queueTitle')"
        :description="t('customerService.conversations.queueDescription')"
      >
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-2">
            <button
              v-for="filter in filters"
              :key="filter.value"
              type="button"
              class="rounded-xl border px-3 py-2 text-sm font-medium transition"
              :class="activeFilter === filter.value ? 'border-primary-200 bg-primary-50 text-primary-700' : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-50'"
              @click="activeFilter = filter.value"
            >
              {{ filter.label }}
            </button>
          </div>

          <div v-if="conversations.length" class="space-y-3">
            <button
              v-for="item in conversations"
              :key="item.id"
              type="button"
              class="w-full rounded-2xl border p-4 text-left transition"
              :class="selectedConversationId === item.id ? 'border-primary-200 bg-primary-50/70' : 'border-gray-200 bg-white hover:bg-gray-50'"
              @click="selectedConversationId = item.id"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <p class="truncate text-sm font-semibold text-gray-900">{{ item.customer.name }}</p>
                    <span v-if="item.unread" class="rounded-full bg-rose-500 px-2 py-0.5 text-[11px] font-semibold text-white">
                      {{ item.unread }}
                    </span>
                  </div>
                  <p class="mt-1 truncate text-sm text-gray-500">{{ item.subject }}</p>
                </div>
                <StatusBadge :label="t(`customerService.status.${item.status}`)" :tone="item.status === 'waiting' ? 'warning' : item.status === 'closed' ? 'success' : 'neutral'" />
              </div>
              <div class="mt-3 flex items-center justify-between text-xs text-gray-500">
                <span>{{ item.channel }}</span>
                <span>{{ formatDate(item.last_active) }}</span>
              </div>
              <p class="mt-3 line-clamp-2 text-sm text-gray-600">{{ item.preview }}</p>
            </button>
          </div>

          <div v-else class="rounded-2xl border border-dashed border-gray-200 bg-gray-50 p-6 text-sm text-gray-500">
            {{ t('customerService.conversations.empty') }}
          </div>
        </div>
      </PageSection>

      <PageSection
        :eyebrow="t('customerService.conversations.eyebrow')"
        :title="t('customerService.conversations.detailTitle')"
        :description="t('customerService.conversations.detailDescription')"
      >
        <div v-if="selectedConversation" class="space-y-5">
          <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
              <div>
                <p class="text-lg font-semibold text-gray-900">{{ selectedConversation.subject }}</p>
                <p class="mt-1 text-sm text-gray-500">{{ selectedConversation.ticket_no }} · {{ selectedConversation.queue }}</p>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <StatusBadge :label="t(`customerService.status.${selectedConversation.status}`)" :tone="selectedConversation.status === 'waiting' ? 'warning' : selectedConversation.status === 'closed' ? 'success' : 'neutral'" />
                <StatusBadge :label="t(`customerService.status.${selectedConversation.priority}`)" :tone="selectedConversation.priority === 'urgent' || selectedConversation.priority === 'high' ? 'warning' : 'neutral'" />
              </div>
            </div>
            <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
              <div class="rounded-xl border border-gray-200 bg-white p-3">
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.channel') }}</p>
                <p class="mt-2 text-sm font-medium text-gray-900">{{ selectedConversation.channel }}</p>
              </div>
              <div class="rounded-xl border border-gray-200 bg-white p-3">
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.assignee') }}</p>
                <p class="mt-2 text-sm font-medium text-gray-900">{{ selectedConversation.assignee || '-' }}</p>
              </div>
              <div class="rounded-xl border border-gray-200 bg-white p-3">
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.priority') }}</p>
                <p class="mt-2 text-sm font-medium text-gray-900">{{ t(`customerService.status.${selectedConversation.priority}`) }}</p>
              </div>
              <div class="rounded-xl border border-gray-200 bg-white p-3">
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.sla') }}</p>
                <p class="mt-2 text-sm font-medium text-gray-900">{{ selectedConversation.sla }}</p>
              </div>
            </div>
          </div>

          <div class="space-y-4 rounded-2xl border border-gray-200 bg-white p-4">
            <div
              v-for="message in selectedConversation.messages"
              :key="message.id"
              class="flex"
              :class="message.sender === 'agent' ? 'justify-end' : 'justify-start'"
            >
              <div
                class="max-w-[85%] rounded-2xl px-4 py-3"
                :class="message.sender === 'agent' ? 'bg-primary-600 text-white' : message.type === 'note' ? 'bg-amber-50 text-amber-800' : 'bg-gray-100 text-gray-800'"
              >
                <div class="flex items-center gap-2 text-xs opacity-80">
                  <span>{{ message.author }}</span>
                  <span>{{ formatDate(message.time) }}</span>
                </div>
                <p class="mt-2 whitespace-pre-line text-sm leading-6">{{ message.content }}</p>
              </div>
            </div>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <div class="flex items-center justify-between gap-3">
              <p class="text-sm font-semibold text-gray-900">{{ t('customerService.conversations.composer.shortcuts') }}</p>
              <label class="flex items-center gap-2 text-sm text-gray-600">
                <input v-model="internalNote" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                <span>{{ t('customerService.conversations.composer.internalNote') }}</span>
              </label>
            </div>
            <div class="mt-3 flex flex-wrap gap-2">
              <button
                v-for="reply in quickReplies"
                :key="reply.key"
                type="button"
                class="rounded-full border border-gray-200 bg-white px-3 py-1.5 text-sm text-gray-600 transition hover:bg-gray-100"
                @click="appendQuickReply(reply.text)"
              >
                {{ reply.label }}
              </button>
            </div>
            <textarea
              v-model="draftMessage"
              rows="4"
              class="ui-textarea mt-4"
              :placeholder="t('customerService.conversations.composer.placeholder')"
            />
            <div class="mt-4 flex flex-wrap items-center justify-between gap-3">
              <div class="text-sm text-gray-500">
                {{ t('customerService.conversations.labels.notes') }}: {{ internalNote ? t('customerService.conversations.composer.internalNote') : t('common.no') }}
              </div>
              <div class="flex flex-wrap gap-2">
                <button type="button" class="btn-secondary" :disabled="actionRunning || !selectedConversation" @click="handleTransfer">
                  {{ t('customerService.conversations.composer.transfer') }}
                </button>
                <button type="button" class="btn-secondary" :disabled="actionRunning || !selectedConversation" @click="handleResolve">
                  {{ t('customerService.conversations.composer.resolve') }}
                </button>
                <button type="button" class="btn-primary" :disabled="actionRunning || !draftMessage.trim() || !selectedConversation" @click="handleSend">
                  {{ t('customerService.conversations.composer.send') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </PageSection>

      <PageSection
        :eyebrow="t('customerService.conversations.eyebrow')"
        :title="t('customerService.conversations.profileTitle')"
        :description="t('customerService.conversations.profileDescription')"
      >
        <div v-if="selectedConversation" class="space-y-4">
          <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4">
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="text-lg font-semibold text-gray-900">{{ selectedConversation.customer.name }}</p>
                <p class="mt-1 text-sm text-gray-500">{{ selectedConversation.customer.company }}</p>
              </div>
              <StatusBadge :label="t(`customerService.status.${selectedConversation.customer.presence}`)" :tone="selectedConversation.customer.presence === 'online' ? 'success' : 'neutral'" />
            </div>
            <div class="mt-4 space-y-3 text-sm text-gray-600">
              <div>
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.contact') }}</p>
                <p class="mt-1 font-medium text-gray-900">{{ selectedConversation.customer.contact }}</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-[0.18em] text-gray-400">{{ t('customerService.conversations.labels.tags') }}</p>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span v-for="tag in selectedConversation.customer.tags" :key="tag" class="rounded-full bg-white px-3 py-1 text-xs font-medium text-gray-600">
                    {{ tag }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-4">
            <div class="space-y-3 text-sm">
              <div class="flex items-center justify-between gap-3">
                <span class="text-gray-500">{{ t('customerService.conversations.labels.lastOrder') }}</span>
                <span class="font-medium text-gray-900">{{ selectedConversation.customer.last_order }}</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-gray-500">{{ t('customerService.conversations.labels.openTickets') }}</span>
                <span class="font-medium text-gray-900">{{ selectedConversation.customer.open_tickets }}</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-gray-500">{{ t('customerService.conversations.labels.satisfaction') }}</span>
                <span class="font-medium text-gray-900">{{ selectedConversation.customer.satisfaction }}</span>
              </div>
            </div>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-4">
            <p class="text-sm font-semibold text-gray-900">{{ t('customerService.conversations.labels.history') }}</p>
            <div class="mt-4 space-y-4">
              <div v-for="event in selectedConversation.history" :key="event.id" class="border-l-2 border-gray-200 pl-4">
                <p class="text-sm font-medium text-gray-900">{{ event.title }}</p>
                <p class="mt-1 text-xs text-gray-500">{{ formatDate(event.time) }}</p>
                <p class="mt-2 text-sm text-gray-600">{{ event.description }}</p>
              </div>
            </div>
          </div>

          <div class="rounded-2xl border border-gray-200 bg-white p-4">
            <p class="text-sm font-semibold text-gray-900">{{ t('customerService.conversations.labels.actions') }}</p>
            <div class="mt-4 space-y-3">
              <button type="button" class="btn-secondary w-full justify-center" :disabled="actionRunning || !selectedConversation" @click="handleTransfer">
                {{ t('customerService.conversations.composer.transfer') }}
              </button>
              <button type="button" class="btn-secondary w-full justify-center" :disabled="actionRunning || !selectedConversation" @click="handleResolve">
                {{ t('customerService.conversations.composer.resolve') }}
              </button>
            </div>
          </div>
        </div>
      </PageSection>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import PageSection from '@/components/base/PageSection.vue'
import StatusBadge from '@/components/base/StatusBadge.vue'
import { useAsyncAction } from '@/composables/useAsyncAction'
import {
  getConversationDetail,
  getConversationSummary,
  listConversations,
  replyConversation,
  resolveConversation,
  transferConversation
} from '@/api/modules/conversations'
import { formatDate } from '@/utils/format'
import type { ConversationDetail, ConversationItem, ConversationSummary } from '@/types'
import { useAppStore } from '@/stores/app'

type FilterValue = 'all' | 'waiting' | 'processing' | 'vip'

const { t } = useI18n()
const appStore = useAppStore()
const { running: actionRunning, run } = useAsyncAction()

const activeFilter = ref<FilterValue>('all')
const selectedConversationId = ref<number>(0)
const conversations = ref<ConversationItem[]>([])
const selectedConversation = ref<ConversationDetail | null>(null)
const summary = ref<ConversationSummary>({
  waiting_count: 0,
  processing_count: 0,
  sla_risk_count: 0,
  resolved_today: 0,
  online_agent_count: 0,
  transfer_count: 0,
  queue_items: [],
  agent_statuses: [],
  recent_transfers: []
})
const draftMessage = ref('')
const internalNote = ref(false)

const filters = computed(() => [
  { value: 'all' as const, label: t('customerService.conversations.filters.all') },
  { value: 'waiting' as const, label: t('customerService.conversations.filters.waiting') },
  { value: 'processing' as const, label: t('customerService.conversations.filters.processing') },
  { value: 'vip' as const, label: t('customerService.conversations.filters.vip') }
])

const cards = computed(() => [
  { label: t('customerService.conversations.cards.waiting'), value: summary.value.waiting_count, description: t('customerService.conversations.cardDescriptions.waiting') },
  { label: t('customerService.conversations.cards.active'), value: summary.value.processing_count, description: t('customerService.conversations.cardDescriptions.active') },
  { label: t('customerService.conversations.cards.slaRisk'), value: summary.value.sla_risk_count, description: t('customerService.conversations.cardDescriptions.risk') },
  { label: t('customerService.conversations.cards.resolved'), value: summary.value.resolved_today, description: t('customerService.conversations.cardDescriptions.resolved') }
])

const quickReplies = computed(() => [
  { key: 'greeting', label: t('customerService.conversations.composer.greetingLabel'), text: t('customerService.conversations.quickReplies.greeting') },
  { key: 'investigate', label: t('customerService.conversations.composer.updateLabel'), text: t('customerService.conversations.quickReplies.investigate') },
  { key: 'escalation', label: t('customerService.conversations.composer.escalationLabel'), text: t('customerService.conversations.quickReplies.escalation') }
])

watch(activeFilter, async () => {
  await loadConversations()
})

watch(selectedConversationId, async (value) => {
  if (!value) {
    selectedConversation.value = null
    return
  }
  selectedConversation.value = await getConversationDetail(value)
})

onMounted(async () => {
  await Promise.all([loadSummary(), loadConversations()])
})

function filterParams() {
  if (activeFilter.value === 'vip') {
    return { page: 1, page_size: 20, tier: 'vip' }
  }
  if (activeFilter.value === 'all') {
    return { page: 1, page_size: 20 }
  }
  return { page: 1, page_size: 20, status: activeFilter.value }
}

async function loadSummary() {
  summary.value = await getConversationSummary()
}

async function loadConversations() {
  const response = await listConversations(filterParams())
  conversations.value = response.items
  if (!conversations.value.length) {
    selectedConversationId.value = 0
    return
  }
  if (!conversations.value.some((item) => item.id === selectedConversationId.value)) {
    selectedConversationId.value = conversations.value[0].id
  } else if (selectedConversationId.value) {
    selectedConversation.value = await getConversationDetail(selectedConversationId.value)
  }
}

function appendQuickReply(text: string) {
  draftMessage.value = draftMessage.value ? `${draftMessage.value}\n${text}` : text
}

async function handleSend() {
  if (!selectedConversationId.value || !draftMessage.value.trim()) {
    return
  }
  await run(async () => {
    await replyConversation(selectedConversationId.value, {
      content: draftMessage.value,
      internal_note: internalNote.value
    })
    appStore.notify(t('customerService.feedback.replySuccess'))
    draftMessage.value = ''
    internalNote.value = false
    await Promise.all([loadSummary(), loadConversations()])
  })
}

async function handleTransfer() {
  if (!selectedConversationId.value) {
    return
  }
  await run(async () => {
    await transferConversation(selectedConversationId.value, {
      assignee: 'Tier-2 Queue',
      queue: 'Escalation Queue',
      note: t('customerService.feedback.transferNote')
    })
    appStore.notify(t('customerService.feedback.transferSuccess'))
    await Promise.all([loadSummary(), loadConversations()])
  })
}

async function handleResolve() {
  if (!selectedConversationId.value) {
    return
  }
  await run(async () => {
    await resolveConversation(selectedConversationId.value, {
      note: t('customerService.feedback.resolveNote')
    })
    appStore.notify(t('customerService.feedback.resolveSuccess'))
    await Promise.all([loadSummary(), loadConversations()])
  })
}
</script>
