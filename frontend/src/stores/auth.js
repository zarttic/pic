import { defineStore } from 'pinia'
import api from '../api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token'),
    refreshToken: localStorage.getItem('refreshToken'),
    isAuthenticated: !!localStorage.getItem('token')
  }),

  actions: {
    async login(username, password) {
      try {
        const response = await api.post('/auth/login', {
          username,
          password
        })

        const { token, refresh_token, user } = response.data

        this.token = token
        this.refreshToken = refresh_token
        this.user = user
        this.isAuthenticated = true

        localStorage.setItem('token', token)
        localStorage.setItem('refreshToken', refresh_token)

        return { success: true }
      } catch (error) {
        const message = error.response?.data?.error || '登录失败'
        return { success: false, error: message }
      }
    },

    async logout() {
      try {
        await api.post('/auth/logout')
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.clearAuth()
      }
    },

    async refreshAccessToken() {
      try {
        const response = await api.post('/auth/refresh', {
          refresh_token: this.refreshToken
        })

        const { token } = response.data
        this.token = token
        localStorage.setItem('token', token)

        return { success: true }
      } catch (error) {
        this.clearAuth()
        return { success: false }
      }
    },

    async fetchCurrentUser() {
      try {
        const response = await api.get('/me')
        this.user = response.data
        return { success: true }
      } catch (error) {
        this.clearAuth()
        return { success: false }
      }
    },

    clearAuth() {
      this.token = null
      this.refreshToken = null
      this.user = null
      this.isAuthenticated = false

      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
    },

    checkAuth() {
      return this.isAuthenticated && !!this.token
    }
  }
})
