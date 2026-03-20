// 核心目的：认证组合式函数
// 模块功能：登录/登出、Token管理、用户权限判断
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

export function useAuth() {
  const userStore = useUserStore()
  const router = useRouter()

  const isLoggedIn = computed(() => userStore.isLoggedIn)
  const currentUser = computed(() => userStore.profile)
  const userRole = computed(() => userStore.profile?.role_name ?? '')

  async function login(form: { username: string, password: string }) {
    await userStore.login(form)
    await router.push('/')
  }

  async function logout() {
    await userStore.logout()
    await router.push('/login')
  }

  function hasRole(...roles: string[]) {
    return roles.includes(userRole.value)
  }

  function hasPermission(permission: string) {
    return userStore.hasPermission(permission)
  }

  return { isLoggedIn, currentUser, userRole, login, logout, hasRole, hasPermission }
}
