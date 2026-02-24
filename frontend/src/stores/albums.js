import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

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
      console.error('Error fetching albums:', err)
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
      console.error('Error fetching album:', err)
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
      console.error('Error creating album:', err)
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
      console.error('Error updating album:', err)
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
      console.error('Error deleting album:', err)
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
      console.error('Error adding photo to album:', err)
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
      console.error('Error removing photo from album:', err)
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
    createAlbum,
    updateAlbum,
    deleteAlbum,
    addPhotoToAlbum,
    removePhotoFromAlbum
  }
})
