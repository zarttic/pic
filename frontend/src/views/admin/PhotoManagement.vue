<template>
  <div class="photo-management">
    <div class="page-header">
      <h2 class="page-title">照片管理</h2>
      <div class="header-actions">
        <button
          v-if="selectedPhotos.length > 0"
          class="btn-danger"
          @click="handleBatchDelete"
        >
          批量删除 ({{ selectedPhotos.length }})
        </button>
        <button
          v-if="selectedPhotos.length > 0"
          class="btn-secondary"
          @click="clearSelection"
        >
          取消选择
        </button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filter-section">
      <div class="search-box">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索照片..."
          @input="handleSearch"
        />
      </div>
      <div class="filter-controls">
        <select v-model="filterYear" @change="handleFilter">
          <option value="">所有年份</option>
          <option v-for="year in uniqueYears" :key="year" :value="year">
            {{ year }}
          </option>
        </select>
        <select v-model="filterFeatured" @change="handleFilter">
          <option value="">全部</option>
          <option value="true">精选</option>
          <option value="false">非精选</option>
        </select>
      </div>
    </div>

    <!-- 上传表单 -->
    <div class="upload-section">
      <h3 class="subsection-title">上传新照片</h3>
      <form @submit.prevent="handleUpload" class="upload-form">
        <div class="form-group">
          <label for="title">标题</label>
          <input
            id="title"
            v-model="uploadForm.title"
            type="text"
            required
            placeholder="照片标题"
          />
        </div>

        <div class="form-group">
          <label for="location">拍摄地点</label>
          <input
            id="location"
            v-model="uploadForm.location"
            type="text"
            placeholder="拍摄地点"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="year">年份</label>
            <input
              id="year"
              v-model="uploadForm.year"
              type="number"
              placeholder="2024"
            />
          </div>

          <div class="form-group">
            <label for="camera">相机</label>
            <input
              id="camera"
              v-model="uploadForm.camera_model"
              type="text"
              placeholder="相机型号"
            />
          </div>
        </div>

        <div class="form-group">
          <label for="description">描述</label>
          <textarea
            id="description"
            v-model="uploadForm.description"
            rows="3"
            placeholder="照片描述"
          ></textarea>
        </div>

        <div class="form-group">
          <label for="tags">标签</label>
          <input
            id="tags"
            v-model="uploadForm.tags"
            type="text"
            placeholder="标签，用逗号分隔"
          />
        </div>

        <div class="form-group">
          <label for="file">选择照片</label>
          <input
            id="file"
            type="file"
            accept="image/*"
            @change="handleFileSelect"
            required
          />
        </div>

        <div class="form-group">
          <label class="checkbox-label">
            <input type="checkbox" v-model="uploadForm.is_featured" />
            设为精选
          </label>
        </div>

        <button type="submit" class="btn-primary" :disabled="uploading">
          {{ uploading ? '上传中...' : '上传照片' }}
        </button>
      </form>
    </div>

    <!-- 照片列表 -->
    <div class="photos-list">
      <div class="list-header">
        <h3 class="subsection-title">已上传照片 ({{ filteredPhotos.length }})</h3>
        <label class="select-all-label">
          <input
            type="checkbox"
            :checked="isAllSelected"
            @change="toggleSelectAll"
          />
          全选
        </label>
      </div>

      <PhotoGridSkeleton v-if="photoStore.loading" :count="8" />

      <div v-else-if="filteredPhotos.length === 0" class="empty">
        暂无照片
      </div>

      <div v-else class="photos-grid">
        <div
          v-for="photo in filteredPhotos"
          :key="photo.id"
          class="photo-item"
          :class="{ selected: selectedPhotos.includes(photo.id) }"
        >
          <div class="photo-checkbox">
            <input
              type="checkbox"
              :checked="selectedPhotos.includes(photo.id)"
              @change="togglePhotoSelection(photo.id)"
            />
          </div>
          <img
            v-lazyload="{
              src: getImageUrl(photo.thumbnail_path || photo.file_path)
            }"
            :alt="photo.title"
          />
          <div class="photo-info">
            <h4>{{ photo.title }}</h4>
            <p>{{ photo.location }}</p>
            <div class="photo-meta">
              <span v-if="photo.year">{{ photo.year }}</span>
              <span v-if="photo.view_count > 0" class="view-count">
                {{ photo.view_count }} 次浏览
              </span>
              <span v-if="photo.is_featured" class="featured-badge">精选</span>
            </div>
            <div v-if="photo.tags && photo.tags.length" class="photo-tags">
              <span v-for="tag in photo.tags" :key="tag" class="tag">
                {{ tag }}
              </span>
            </div>
          </div>
          <div class="photo-actions">
            <button class="btn-icon" @click="openEditDialog(photo)" title="编辑">
              ✏️
            </button>
            <button class="btn-icon btn-icon-danger" @click="handleDelete(photo.id)" title="删除">
              🗑️
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑对话框 -->
    <div v-if="editDialogVisible" class="dialog-overlay" @click="closeEditDialog">
      <div class="dialog-content" @click.stop>
        <h3 class="dialog-title">编辑照片</h3>
        <form @submit.prevent="handleEdit" class="edit-form">
          <div class="form-group">
            <label for="edit-title">标题</label>
            <input
              id="edit-title"
              v-model="editForm.title"
              type="text"
              required
            />
          </div>

          <div class="form-group">
            <label for="edit-location">拍摄地点</label>
            <input
              id="edit-location"
              v-model="editForm.location"
              type="text"
            />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label for="edit-year">年份</label>
              <input
                id="edit-year"
                v-model="editForm.year"
                type="number"
              />
            </div>

            <div class="form-group">
              <label for="edit-camera">相机</label>
              <input
                id="edit-camera"
                v-model="editForm.camera_model"
                type="text"
              />
            </div>
          </div>

          <div class="form-group">
            <label for="edit-description">描述</label>
            <textarea
              id="edit-description"
              v-model="editForm.description"
              rows="3"
            ></textarea>
          </div>

          <div class="form-group">
            <label for="edit-tags">标签</label>
            <input
              id="edit-tags"
              v-model="editForm.tagsInput"
              type="text"
              placeholder="标签，用逗号分隔"
            />
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="editForm.is_featured" />
              设为精选
            </label>
          </div>

          <div class="dialog-actions">
            <button type="button" class="btn-secondary" @click="closeEditDialog">
              取消
            </button>
            <button type="submit" class="btn-primary">
              保存
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { usePhotoStore } from '../../stores/photos'
import { useNotificationStore } from '../../stores/notification'
import { getImageUrl } from '../../utils/index'
import PhotoGridSkeleton from '../../components/PhotoGridSkeleton.vue'

const photoStore = usePhotoStore()
const notification = useNotificationStore()
const uploading = ref(false)
const editDialogVisible = ref(false)

const selectedPhotos = ref([])
const searchQuery = ref('')
const filterYear = ref('')
const filterFeatured = ref('')

const uploadForm = ref({
  title: '',
  description: '',
  location: '',
  year: new Date().getFullYear(),
  camera_model: '',
  lens: '',
  aperture: '',
  shutter_speed: '',
  iso: null,
  tags: '',
  is_featured: false
})

const editForm = ref({
  id: null,
  title: '',
  description: '',
  location: '',
  year: null,
  camera_model: '',
  tagsInput: '',
  is_featured: false
})

const selectedFile = ref(null)

// 计算属性
const uniqueYears = computed(() => {
  const years = photoStore.photos
    .map(p => p.year)
    .filter(year => year)
  return [...new Set(years)].sort((a, b) => b - a)
})

const filteredPhotos = computed(() => {
  let photos = photoStore.photos

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    photos = photos.filter(p =>
      p.title.toLowerCase().includes(query) ||
      (p.location && p.location.toLowerCase().includes(query)) ||
      (p.description && p.description.toLowerCase().includes(query)) ||
      (p.tags && p.tags.some(tag => tag.toLowerCase().includes(query)))
    )
  }

  // 年份过滤
  if (filterYear.value) {
    photos = photos.filter(p => p.year === parseInt(filterYear.value))
  }

  // 精选过滤
  if (filterFeatured.value !== '') {
    const isFeatured = filterFeatured.value === 'true'
    photos = photos.filter(p => p.is_featured === isFeatured)
  }

  return photos
})

const isAllSelected = computed(() => {
  return filteredPhotos.value.length > 0 &&
         filteredPhotos.value.every(p => selectedPhotos.value.includes(p.id))
})

onMounted(() => {
  photoStore.fetchPhotos()
})

const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const handleUpload = async () => {
  if (!selectedFile.value) {
    notification.error('请选择照片文件')
    return
  }

  uploading.value = true
  const formData = new FormData()
  formData.append('file', selectedFile.value)
  formData.append('title', uploadForm.value.title)
  formData.append('description', uploadForm.value.description || '')
  formData.append('location', uploadForm.value.location || '')
  formData.append('year', uploadForm.value.year || new Date().getFullYear())
  formData.append('camera_model', uploadForm.value.camera_model || '')
  formData.append('lens', uploadForm.value.lens || '')
  formData.append('aperture', uploadForm.value.aperture || '')
  formData.append('shutter_speed', uploadForm.value.shutter_speed || '')
  formData.append('iso', uploadForm.value.iso || 0)
  formData.append('tags', uploadForm.value.tags || '')
  formData.append('is_featured', uploadForm.value.is_featured)

  try {
    await photoStore.uploadPhoto(formData)
    // 重置表单
    uploadForm.value = {
      title: '',
      description: '',
      location: '',
      year: new Date().getFullYear(),
      camera_model: '',
      lens: '',
      aperture: '',
      shutter_speed: '',
      iso: null,
      tags: '',
      is_featured: false
    }
    selectedFile.value = null
    document.getElementById('file').value = ''
    notification.success('照片上传成功！')
  } catch (error) {
    notification.error('上传失败：' + error.message)
  } finally {
    uploading.value = false
  }
}

const handleDelete = async (id) => {
  if (confirm('确定要删除这张照片吗？')) {
    try {
      await photoStore.deletePhoto(id)
      notification.success('删除成功！')
    } catch (error) {
      notification.error('删除失败：' + error.message)
    }
  }
}

const handleBatchDelete = async () => {
  if (!confirm(`确定要删除选中的 ${selectedPhotos.value.length} 张照片吗？`)) {
    return
  }

  try {
    await photoStore.batchDelete(selectedPhotos.value)
    notification.success(`成功删除 ${selectedPhotos.value.length} 张照片`)
    selectedPhotos.value = []
  } catch (error) {
    notification.error('批量删除失败：' + error.message)
  }
}

const togglePhotoSelection = (id) => {
  const index = selectedPhotos.value.indexOf(id)
  if (index > -1) {
    selectedPhotos.value.splice(index, 1)
  } else {
    selectedPhotos.value.push(id)
  }
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedPhotos.value = []
  } else {
    selectedPhotos.value = filteredPhotos.value.map(p => p.id)
  }
}

const clearSelection = () => {
  selectedPhotos.value = []
}

const handleSearch = () => {
  // 搜索时清空选择
  selectedPhotos.value = []
}

const handleFilter = () => {
  // 筛选时清空选择
  selectedPhotos.value = []
}

const openEditDialog = (photo) => {
  editForm.value = {
    id: photo.id,
    title: photo.title,
    description: photo.description || '',
    location: photo.location || '',
    year: photo.year,
    camera_model: photo.camera_model || '',
    tagsInput: photo.tags ? photo.tags.join(', ') : '',
    is_featured: photo.is_featured || false
  }
  editDialogVisible.value = true
}

const closeEditDialog = () => {
  editDialogVisible.value = false
  editForm.value = {
    id: null,
    title: '',
    description: '',
    location: '',
    year: null,
    camera_model: '',
    tagsInput: '',
    is_featured: false
  }
}

const handleEdit = async () => {
  try {
    const updateData = {
      ...editForm.value,
      tags: editForm.value.tagsInput
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag)
    }

    await photoStore.updatePhoto(editForm.value.id, updateData)
    notification.success('更新成功！')
    closeEditDialog()
  } catch (error) {
    notification.error('更新失败：' + error.message)
  }
}
</script>

<style scoped>
.photo-management {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: var(--spacing-xl);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.page-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem;
  font-weight: 300;
  letter-spacing: 0.1em;
}

.header-actions {
  display: flex;
  gap: var(--spacing-sm);
}

/* 筛选区域 */
.filter-section {
  background: var(--bg-secondary);
  padding: var(--spacing-md);
  border-radius: 8px;
  margin-bottom: var(--spacing-lg);
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 250px;
}

.search-box input {
  width: 100%;
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-family: inherit;
}

.filter-controls {
  display: flex;
  gap: var(--spacing-sm);
}

.filter-controls select {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  cursor: pointer;
}

.upload-section,
.photos-list {
  background: var(--bg-secondary);
  padding: var(--spacing-lg);
  border-radius: 8px;
  margin-bottom: var(--spacing-lg);
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.subsection-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  margin-bottom: var(--spacing-md);
  letter-spacing: 0.1em;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
}

label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  text-transform: none;
}

.select-all-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  text-transform: none;
}

input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

input,
textarea {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-family: inherit;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--accent-gold);
}

input[type="file"] {
  cursor: pointer;
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: 0.1em;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  border: none;
}

.btn-primary {
  background: var(--accent-gold);
  color: var(--bg-primary);
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-warm);
  transform: translateY(-2px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
}

.btn-secondary:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
}

.btn-danger {
  background: rgba(220, 53, 69, 0.9);
  color: white;
}

.btn-danger:hover {
  background: rgba(220, 53, 69, 1);
}

.empty {
  text-align: center;
  padding: var(--spacing-lg);
  color: var(--text-secondary);
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: var(--spacing-md);
}

.photo-item {
  background: var(--bg-primary);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  border: 2px solid transparent;
  transition: border-color 0.3s ease;
}

.photo-item.selected {
  border-color: var(--accent-gold);
}

.photo-checkbox {
  position: absolute;
  top: var(--spacing-sm);
  left: var(--spacing-sm);
  z-index: 10;
  background: rgba(10, 10, 10, 0.8);
  padding: 4px;
  border-radius: 4px;
}

.photo-item img {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.photo-info {
  padding: var(--spacing-sm);
}

.photo-info h4 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.2rem;
  font-weight: 300;
  margin-bottom: var(--spacing-xs);
}

.photo-info p {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.photo-meta {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-xs);
  flex-wrap: wrap;
  font-size: 0.75rem;
}

.view-count {
  color: var(--accent-gold);
}

.featured-badge {
  background: var(--accent-gold);
  color: var(--bg-primary);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.7rem;
  font-weight: 600;
}

.photo-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-xs);
}

.tag {
  background: rgba(201, 169, 98, 0.2);
  color: var(--accent-gold);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 0.7rem;
}

.photo-actions {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  display: flex;
  gap: var(--spacing-xs);
}

.btn-icon {
  width: 32px;
  height: 32px;
  background: rgba(10, 10, 10, 0.8);
  border: 1px solid rgba(201, 169, 98, 0.3);
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.btn-icon:hover {
  background: var(--accent-gold);
  border-color: var(--accent-gold);
}

.btn-icon-danger:hover {
  background: rgba(220, 53, 69, 0.9);
  border-color: rgba(220, 53, 69, 0.9);
}

/* 对话框样式 */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog-content {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.dialog-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem;
  font-weight: 300;
  margin-bottom: var(--spacing-lg);
  text-align: center;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.dialog-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .filter-section {
    flex-direction: column;
  }

  .filter-controls {
    width: 100%;
  }

  .filter-controls select {
    flex: 1;
  }

  .photos-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
}
</style>
