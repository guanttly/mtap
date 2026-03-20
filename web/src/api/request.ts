import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { message } from 'ant-design-vue'

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
    const { code, message: msg, data } = response.data
    if (code !== 0) {
      message.error(msg || '请求失败')
      return Promise.reject(new Error(msg))
    }
    return data as any
  },
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }
    if (error.response?.status === 401 && !originalRequest._retry) {
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

    const msg = error.response?.data?.message || error.message || '网络错误'
    message.error(msg)
    return Promise.reject(error)
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
