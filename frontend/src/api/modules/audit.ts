import http, { unwrap } from '@/api/http'
import type { AuditLog, PaginatedData } from '@/types'

export function listAuditLogs(params: Record<string, unknown>) {
  return unwrap<PaginatedData<AuditLog>>(http.get('/audit-logs', { params }))
}
