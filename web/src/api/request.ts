import { message } from 'ant-design-vue'
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { getErrorMessage } from '@/utils/errorCode'

// 统一响应结构
interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

// 分页结果
export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

const http: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

// 请求拦截器：注入 JWT Token
http.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token)
    config.headers.Authorization = `Bearer ${token}`
  return config
})

// 响应拦截器：统一错误处理 + Token 刷新
let isRefreshing = false
let pendingQueue: Array<(token: string) => void> = []

http.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, data } = response.data
    if (code !== 0) {
      const detail = typeof data === 'string' ? data : undefined
      const errMsg = getErrorMessage(code, detail)
      message.error(errMsg)
      return Promise.reject(new Error(errMsg))
    }
    return data as any
  },
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    // 登录接口本身不走 token 刷新逻辑，直接抛出错误让调用方处理
    const isLoginRequest = originalRequest.url?.includes('/auth/login')

    if (error.response?.status === 401 && !originalRequest._retry && !isLoginRequest) {
      const refreshToken = localStorage.getItem('refresh_token')
      if (!refreshToken) {
        localStorage.clear()
        window.location.href = '/login'
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve) => {
          pendingQueue.push((token: string) => {
            originalRequest.headers = { ...originalRequest.headers, Authorization: `Bearer ${token}` }
            resolve(http(originalRequest))
          })
        })
      }

      isRefreshing = true
      originalRequest._retry = true
      try {
        const res = await axios.post<ApiResponse<{ access_token: string, refresh_token: string }>>('/api/v1/auth/refresh', { refresh_token: refreshToken })
        const newToken = res.data.data.access_token
        const newRefresh = res.data.data.refresh_token
        localStorage.setItem('access_token', newToken)
        localStorage.setItem('refresh_token', newRefresh)
        pendingQueue.forEach(cb => cb(newToken))
        pendingQueue = []
        originalRequest.headers = { ...originalRequest.headers, Authorization: `Bearer ${newToken}` }
        return http(originalRequest)
      }
      catch {
        localStorage.clear()
        window.location.href = '/login'
        return Promise.reject(error)
      }
      finally {
        isRefreshing = false
      }
    }

    const resData = error.response?.data
    const code: number | undefined = resData?.code
    const detail = typeof resData?.data === 'string' ? resData.data : undefined
    // 按错误码查前端码表，找不到则用 message 字段，最后兜底网络错误
    const msg = code
      ? getErrorMessage(code, detail)
      : (resData?.message ?? error.message ?? '网络错误')
    // 登录请求由调用方自己处理错误提示，避免重复弹 toast
    if (!isLoginRequest) {
      message.error(msg)
    }
    return Promise.reject(new Error(msg))
  },
)

export default http

// Auth API
export const authApi = {
  login: (username: string, password: string) =>
    http.post<any, any>('/auth/login', { username, password }),
  refresh: (refreshToken: string) =>
    http.post<any, any>('/auth/refresh', { refresh_token: refreshToken }),
  logout: () => http.post<any, any>('/auth/logout'),
  profile: () => http.get<any, any>('/auth/profile'),
}
