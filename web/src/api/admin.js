import http from './request';
export const adminApi = {
    // 用户管理
    listUsers: (params) => http.get('/admin/users', { params }),
    createUser: (data) => http.post('/admin/users', data),
    updateUser: (id, data) => http.put(`/admin/users/${id}`, data),
    deleteUser: (id) => http.delete(`/admin/users/${id}`),
    resetPassword: (id, newPassword) => http.post(`/admin/users/${id}/reset-password`, { new_password: newPassword }),
    // 角色管理
    listRoles: () => http.get('/admin/roles'),
    createRole: (data) => http.post('/admin/roles', data),
    updateRole: (id, permissions) => http.put(`/admin/roles/${id}`, { permissions }),
    deleteRole: (id) => http.delete(`/admin/roles/${id}`),
};
