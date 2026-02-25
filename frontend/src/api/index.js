import axios from 'axios'
import { handleError } from '../utils/errorHandler'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (error.response) {
      // 开发环境记录错误
      if (import.meta.env.DEV) {
        console.error(`API Error [${error.response.status}]:`, error.response.data)
      }

      // 检查是否是相册密码验证错误(不跳转登录页)
      const isAlbumVerifyError = originalRequest.url?.includes('/albums/') && originalRequest.url?.includes('/verify')

      switch (error.response.status) {
        case 401:
          // 如果是相册密码验证失败,不跳转登录页
          if (isAlbumVerifyError) {
            break
          }

          // 如果是刷新令牌失败，直接清除认证信息
          if (originalRequest.url === '/auth/refresh') {
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            window.location.href = '/admin/login'
            break
          }

          // 尝试刷新令牌
          const refreshToken = localStorage.getItem('refreshToken')
          if (refreshToken && !originalRequest._retry) {
            originalRequest._retry = true

            try {
              const response = await api.post('/auth/refresh', {
                refresh_token: refreshToken
              })

              const { token } = response.data
              localStorage.setItem('token', token)

              // 重试原请求
              originalRequest.headers.Authorization = `Bearer ${token}`
              return api(originalRequest)
            } catch (refreshError) {
              // 刷新失败，跳转登录页
              localStorage.removeItem('token')
              localStorage.removeItem('refreshToken')
              window.location.href = '/admin/login'
              break
            }
          } else {
            // 没有刷新令牌或已重试过，跳转登录页
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            window.location.href = '/admin/login'
          }
          break
      }
    }
    return Promise.reject(error)
  }
)

export default api
