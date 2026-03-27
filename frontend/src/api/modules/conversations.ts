import http, { unwrap } from '@/api/http'
import type {
  ConversationDetail,
  ConversationItem,
  ConversationSummary,
  PaginatedData
} from '@/types'

export function listConversations(params: Record<string, unknown>) {
  return unwrap<PaginatedData<ConversationItem>>(http.get('/conversations', { params }))
}

export function getConversationSummary() {
  return unwrap<ConversationSummary>(http.get('/conversations/summary'))
}

export function getConversationDetail(id: number) {
  return unwrap<ConversationDetail>(http.get(`/conversations/${id}`))
}

export function replyConversation(id: number, payload: { content: string; internal_note: boolean }) {
  return unwrap<{ message: string }>(http.post(`/conversations/${id}/reply`, payload))
}

export function transferConversation(id: number, payload: { assignee?: string; queue?: string; note?: string }) {
  return unwrap<{ message: string }>(http.post(`/conversations/${id}/transfer`, payload))
}

export function resolveConversation(id: number, payload: { note?: string }) {
  return unwrap<{ message: string }>(http.post(`/conversations/${id}/resolve`, payload))
}
