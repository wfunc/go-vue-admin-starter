export interface ApiEnvelope<T> {
  code: string
  message: string
  data: T
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface Role {
  id: number
  name: string
  code: string
  description: string
  permissions: string[]
  is_system: boolean
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  username: string
  email?: string
  nickname: string
  status: string
  role_id: number
  role_name: string
  role_code: string
  created_at: string
  updated_at: string
}

export interface MenuItem {
  id: number
  title: string
  name: string
  path: string
  component: string
  icon: string
  menu_type: string
  permission: string
  sort: number
  hidden: boolean
  parent_id?: number
  children?: MenuItem[]
}

export interface SystemConfigItem {
  id: number
  key: string
  value: string
  category: string
  description: string
  is_public: boolean
  created_at: string
  updated_at: string
}

export interface AuditLog {
  id: number
  username: string
  module: string
  action: string
  method: string
  path: string
  status_code: number
  ip: string
  detail: string
  created_at: string
}

export interface Profile {
  user_id: number
  username: string
  nickname: string
  email?: string
  status: string
  role_id: number
  role_name: string
  role_code: string
  permissions: string[]
}

export interface UpdateProfilePayload {
  nickname: string
  email?: string
}

export interface ChangePasswordPayload {
  current_password: string
  new_password: string
  confirm_password: string
}

export interface Session {
  access_token: string
  expires_at: string
  user: Profile
  menus: MenuItem[]
}

export interface DashboardSummary {
  user_count: number
  role_count: number
  menu_count: number
  config_count: number
  audit_log_count: number
}
