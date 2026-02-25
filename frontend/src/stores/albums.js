import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'
import { handleError } from '../utils/errorHandler'

export const useAlbumStore = defineStore('albums', () => {
  const albums = ref([])
  const currentAlbum = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchAlbums() {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/albums')
      albums.value = response.data.data || response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '获取相册列表')
    } finally {
      loading.value = false
    }
  }

  async function fetchAlbum(id) {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/albums/${id}`)
      currentAlbum.value = response.data
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '获取相册详情')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchAlbumPublic(id) {
    loading.value = true
    error.value = null
    try {
      // 使用原生 fetch，不发送 JWT token
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
      const response = await fetch(`${apiUrl}/albums/${id}`)

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || 'Failed to fetch album')
      }

      const data = await response.json()
      currentAlbum.value = data
      return data
    } catch (err) {
      error.value = err.message
      handleError(err, '获取相册详情')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function verifyPassword(id, password) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post(`/albums/${id}/verify`, { password })
      // 保存 token 到 localStorage
      localStorage.setItem(`album_token_${id}`, response.data.token)
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '验证密码')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function setPassword(id, password) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post(`/albums/${id}/password`, { password })
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '设置密码')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function removePassword(id) {
    loading.value = true
    error.value = null
    try {
      const response = await api.delete(`/albums/${id}/password`)
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '移除密码')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createAlbum(data) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/albums', data)
      albums.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '创建相册')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateAlbum(id, data) {
    loading.value = true
    error.value = null
    try {
      const response = await api.put(`/albums/${id}`, data)
      const index = albums.value.findIndex(a => a.id === id)
      if (index !== -1) {
        albums.value[index] = response.data
      }
      return response.data
    } catch (err) {
      error.value = err.message
      handleError(err, '更新相册')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteAlbum(id) {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/albums/${id}`)
      albums.value = albums.value.filter(a => a.id !== id)
    } catch (err) {
      error.value = err.message
      handleError(err, '删除相册')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function addPhotoToAlbum(albumId, photoId, sortOrder = 0) {
    loading.value = true
    error.value = null
    try {
      await api.post(`/albums/${albumId}/photos`, {
        photo_id: photoId,
        sort_order: sortOrder
      })
    } catch (err) {
      error.value = err.message
      handleError(err, '添加照片到相册')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function removePhotoFromAlbum(albumId, photoId) {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/albums/${albumId}/photos/${photoId}`)
    } catch (err) {
      error.value = err.message
      handleError(err, '从相册移除照片')
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    albums,
    currentAlbum,
    loading,
    error,
    fetchAlbums,
    fetchAlbum,
    fetchAlbumPublic,
    createAlbum,
    updateAlbum,
    deleteAlbum,
    addPhotoToAlbum,
    removePhotoFromAlbum,
    verifyPassword,
    setPassword,
    removePassword
  }
})
