import http, { type PageResult } from './request'

export interface UserInfo {
  id: string
  username: string
  real_name: string
  role_id: string
  role_name: string
  department_id?: string
  status: string
  last_login_at?: string
  created_at: string
}

export interface RoleInfo {
  id: string
  name: string
  permissions: string[]
  is_preset: boolean
  created_at: string
}

export const adminApi = {
  // 用户管理
  listUsers: (params?: Record<string, unknown>) =>
    http.get<any, PageResult<UserInfo>>('/admin/users', { params }),
  createUser: (data: { username: string, password: string, real_name?: string, role_id: string, department_id?: string }) =>
    http.post<any, UserInfo>('/admin/users', data),
  updateUser: (id: string, data: { real_name?: string, role_id?: string, department_id?: string, status?: string }) =>
    http.put<any, void>(`/admin/users/${id}`, data),
  deleteUser: (id: string) =>
    http.delete<any, void>(`/admin/users/${id}`),
  resetPassword: (id: string, newPassword: string) =>
    http.post<any, void>(`/admin/users/${id}/reset-password`, { new_password: newPassword }),

  // 角色管理
  listRoles: () =>
    http.get<any, { items: RoleInfo[], total: number }>('/admin/roles'),
  createRole: (data: { name: string, permissions: string[] }) =>
    http.post<any, RoleInfo>('/admin/roles', data),
  updateRole: (id: string, permissions: string[]) =>
    http.put<any, void>(`/admin/roles/${id}`, { permissions }),
  deleteRole: (id: string) =>
    http.delete<any, void>(`/admin/roles/${id}`),
}
