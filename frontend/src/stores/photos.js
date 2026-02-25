import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const usePhotoStore = defineStore('photos', () => {
  const photos = ref([])
  const currentPhoto = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchPhotos() {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/photos')
      photos.value = response.data.data || response.data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching photos:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchPhoto(id) {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/photos/${id}`)
      currentPhoto.value = response.data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching photo:', err)
    } finally {
      loading.value = false
    }
  }

  async function uploadPhoto(formData) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/photos', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      photos.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.message
      console.error('Error uploading photo:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updatePhoto(id, data) {
    loading.value = true
    error.value = null
    try {
      const response = await api.put(`/photos/${id}`, data)
      const index = photos.value.findIndex(p => p.id === id)
      if (index !== -1) {
        photos.value[index] = response.data
      }
      return response.data
    } catch (err) {
      error.value = err.message
      console.error('Error updating photo:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deletePhoto(id) {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/photos/${id}`)
      photos.value = photos.value.filter(p => p.id !== id)
    } catch (err) {
      error.value = err.message
      console.error('Error deleting photo:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function incrementViewCount(id) {
    try {
      await api.post(`/photos/${id}/view`)
      const photo = photos.value.find(p => p.id === id)
      if (photo) {
        photo.view_count++
      }
    } catch (err) {
      console.error('Error incrementing view count:', err)
    }
  }

  async function batchDelete(ids) {
    loading.value = true
    error.value = null
    try {
      await api.delete('/photos/batch', { data: { ids } })
      photos.value = photos.value.filter(p => !ids.includes(p.id))
    } catch (err) {
      error.value = err.message
      console.error('Error batch deleting photos:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function batchUpdateTags(ids, tags) {
    loading.value = true
    error.value = null
    try {
      await api.patch('/photos/batch/tags', { ids, tags })
      photos.value = photos.value.map(p => {
        if (ids.includes(p.id)) {
          return { ...p, tags }
        }
        return p
      })
    } catch (err) {
      error.value = err.message
      console.error('Error batch updating tags:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function batchSetFeatured(ids, is_featured) {
    loading.value = true
    error.value = null
    try {
      await api.patch('/photos/batch/featured', { ids, is_featured })
      photos.value = photos.value.map(p => {
        if (ids.includes(p.id)) {
          return { ...p, is_featured }
        }
        return p
      })
    } catch (err) {
      error.value = err.message
      console.error('Error batch setting featured:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    photos,
    currentPhoto,
    loading,
    error,
    fetchPhotos,
    fetchPhoto,
    uploadPhoto,
    updatePhoto,
    deletePhoto,
    incrementViewCount,
    batchDelete,
    batchUpdateTags,
    batchSetFeatured
  }
})
