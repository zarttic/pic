<template>
  <div class="admin-page">
    <section class="admin-section">
      <div class="admin-header">
        <h2 class="section-title">照片管理</h2>
        <div class="section-divider"></div>
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
            <label for="file">选择照片</label>
            <input
              id="file"
              type="file"
              accept="image/*"
              @change="handleFileSelect"
              required
            />
          </div>

          <button type="submit" class="btn-primary" :disabled="uploading">
            {{ uploading ? '上传中...' : '上传照片' }}
          </button>
        </form>
      </div>

      <!-- 照片列表 -->
      <div class="photos-list">
        <h3 class="subsection-title">已上传照片</h3>
        <PhotoGridSkeleton v-if="photoStore.loading" :count="6" />
        <div v-else-if="photoStore.photos.length === 0" class="empty">
          暂无照片
        </div>
        <div v-else class="photos-grid">
          <div
            v-for="photo in photoStore.photos"
            :key="photo.id"
            class="photo-item"
          >
            <img :src="photo.file_path" :alt="photo.title" />
            <div class="photo-info">
              <h4>{{ photo.title }}</h4>
              <p>{{ photo.location }}</p>
              <p v-if="photo.view_count > 0" class="view-count">
                {{ photo.view_count }} 次浏览
              </p>
            </div>
            <button
              class="btn-delete"
              @click="handleDelete(photo.id)"
            >
              删除
            </button>
            <button
              class="btn-edit"
              @click="openEditDialog(photo)"
            >
              编辑
            </button>
          </div>
        </div>
      </div>
    </section>

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
            <label>
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
import { ref, onMounted } from 'vue'
import { usePhotoStore } from '../stores/photos'
import { useNotificationStore } from '../stores/notification'
import { useConfirm } from '../composables/useConfirm'
import PhotoGridSkeleton from '../components/PhotoGridSkeleton.vue'

const photoStore = usePhotoStore()
const notification = useNotificationStore()
const confirm = useConfirm()
const uploading = ref(false)
const editDialogVisible = ref(false)

const uploadForm = ref({
  title: '',
  description: '',
  location: '',
  year: new Date().getFullYear(),
  camera_model: '',
  lens: '',
  aperture: '',
  shutter_speed: '',
  iso: null
})

const editForm = ref({
  id: null,
  title: '',
  description: '',
  location: '',
  year: null,
  camera_model: '',
  is_featured: false
})

const selectedFile = ref(null)

onMounted(() => {
  photoStore.fetchPhotos()
})

const handleFileSelect = (event) => {
  selectedFile.value = event.target.files[0]
}

const handleUpload = async () => {
  if (!selectedFile.value) return

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
      iso: null
    }
    selectedFile.value = null
    document.getElementById('file').value = ''
    notification.success('上传成功！')
  } catch (error) {
    notification.error('上传失败：' + error.message)
  } finally {
    uploading.value = false
  }
}

const handleDelete = async (id) => {
  const result = await confirm({
    type: 'danger',
    title: '删除照片',
    message: '确定要删除这张照片吗？此操作不可撤销。',
    confirmText: '删除'
  })

  if (result) {
    try {
      await photoStore.deletePhoto(id)
      notification.success('删除成功！')
    } catch (error) {
      notification.error('删除失败：' + error.message)
    }
  }
}

const openEditDialog = (photo) => {
  editForm.value = {
    id: photo.id,
    title: photo.title,
    description: photo.description || '',
    location: photo.location || '',
    year: photo.year,
    camera_model: photo.camera_model || '',
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
    is_featured: false
  }
}

const handleEdit = async () => {
  try {
    await photoStore.updatePhoto(editForm.value.id, editForm.value)
    notification.success('更新成功！')
    closeEditDialog()
  } catch (error) {
    notification.error('更新失败：' + error.message)
  }
}
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  padding-top: 100px;
}

.admin-section {
  padding: var(--spacing-xl) var(--spacing-lg);
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.section-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(2.5rem, 6vw, 5rem);
  font-weight: 300;
  letter-spacing: 0.15em;
  margin-bottom: var(--spacing-sm);
}

.section-divider {
  width: 60px;
  height: 1px;
  background: var(--accent-gold);
  margin: var(--spacing-md) auto;
}

.upload-section,
.photos-list {
  background: var(--bg-secondary);
  padding: var(--spacing-lg);
  border-radius: 8px;
  margin-bottom: var(--spacing-lg);
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

.btn-primary {
  background: var(--accent-gold);
  color: var(--bg-primary);
  border: none;
  padding: 1rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: 0.1em;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-warm);
  transform: translateY(-2px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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

.view-count {
  font-size: 0.75rem;
  color: var(--accent-gold);
  margin-top: var(--spacing-xs);
}

.btn-delete {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  background: rgba(10, 10, 10, 0.9);
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-delete:hover {
  background: rgba(255, 0, 0, 0.8);
  border-color: rgba(255, 0, 0, 0.8);
}

.btn-edit {
  position: absolute;
  top: var(--spacing-sm);
  right: 60px;
  background: rgba(10, 10, 10, 0.9);
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
  padding: 0.5rem 1rem;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-edit:hover {
  background: var(--accent-gold);
  border-color: var(--accent-gold);
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

.btn-secondary {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  letter-spacing: 0.1em;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-secondary:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
