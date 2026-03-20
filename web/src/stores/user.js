import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { authApi } from '@/api/request';
export const useUserStore = defineStore('user', () => {
    const accessToken = ref(localStorage.getItem('access_token') || '');
    const refreshToken = ref(localStorage.getItem('refresh_token') || '');
    const userInfo = ref(localStorage.getItem('user_info') ? JSON.parse(localStorage.getItem('user_info')) : null);
    const isLoggedIn = computed(() => !!accessToken.value);
    const role = computed(() => userInfo.value?.role ?? '');
    // 供 useAuth / Header 使用的 profile 对象（与后端字段对齐）
    const profile = computed(() => userInfo.value ? {
        user_id: userInfo.value.userId,
        username: userInfo.value.username,
        real_name: userInfo.value.realName,
        role_name: userInfo.value.roleName,
        role: userInfo.value.role,
        department_id: userInfo.value.departmentId,
        permissions: userInfo.value.permissions ?? [],
    } : null);
    function hasRole(...roles) {
        return roles.includes(role.value);
    }
    function hasPermission(permission) {
        return userInfo.value?.permissions?.includes(permission) ?? false;
    }
    async function login(form) {
        const res = await authApi.login(form.username, form.password);
        accessToken.value = res.access_token;
        refreshToken.value = res.refresh_token;
        userInfo.value = {
            userId: res.user_id,
            username: res.username,
            realName: res.real_name,
            role: res.role,
            roleName: res.role_name ?? res.role,
            permissions: res.permissions ?? [],
            departmentId: res.department_id,
        };
        localStorage.setItem('access_token', res.access_token);
        localStorage.setItem('refresh_token', res.refresh_token);
        localStorage.setItem('user_info', JSON.stringify(userInfo.value));
    }
    async function logout() {
        try {
            await authApi.logout();
        }
        finally {
            accessToken.value = '';
            refreshToken.value = '';
            userInfo.value = null;
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('user_info');
        }
    }
    function setTokens(at, rt) {
        accessToken.value = at;
        refreshToken.value = rt;
        localStorage.setItem('access_token', at);
        localStorage.setItem('refresh_token', rt);
    }
    return { accessToken, refreshToken, userInfo, profile, isLoggedIn, role, hasRole, hasPermission, login, logout, setTokens };
});
