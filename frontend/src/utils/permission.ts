import type { MenuItem } from '@/types'

export function hasPermission(permissions: string[], permission?: string) {
  if (!permission) return true
  return permissions.includes('*') || permissions.includes(permission)
}

export function flattenMenus(items: MenuItem[], depth = 0): Array<MenuItem & { depth: number }> {
  return items.flatMap((item) => {
    const current = [{ ...item, depth }]
    const children = item.children ? flattenMenus(item.children, depth + 1) : []
    return [...current, ...children]
  })
}
