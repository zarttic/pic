<template>
  <div class="album-management">
    <div class="page-header">
      <h2 class="page-title">相册管理</h2>
      <button class="btn-primary" @click="showCreateDialog = true">
        创建相册
      </button>
    </div>

    <PhotoGridSkeleton v-if="albumStore.loading" :count="4" />

    <div v-else-if="albumStore.albums.length === 0" class="empty">
      暂无相册，点击上方按钮创建第一个相册
    </div>

    <div v-else class="albums-grid">
      <div
        v-for="album in albumStore.albums"
        :key="album.id"
        class="album-card"
      >
        <div class="album-cover" @click="openAlbumDetail(album)">
          <img
            v-if="album.cover_photo_id && getCoverPhoto(album)"
            v-lazyload="{
              src: getCoverPhoto(album)
            }"
            :alt="album.name"
          />
          <div v-else class="cover-placeholder">
            {{ album.name.charAt(0) }}
          </div>
          <div class="album-overlay">
            <span class="photo-count">{{ album.photos?.length || 0 }} 张照片</span>
          </div>
        </div>
        <div class="album-info">
          <h3 class="album-name">{{ album.name }}</h3>
          <p class="album-description">{{ album.description || '暂无描述' }}</p>
          <div class="album-badges">
            <span :class="['badge', album.is_protected ? 'badge-protected' : 'badge-public']">
              {{ album.is_protected ? '🔒 已加密' : '公开' }}
            </span>
          </div>
        </div>
        <div class="album-actions">
          <button class="btn-icon" @click="editAlbum(album)" title="编辑">
            ✏️
          </button>
          <button
            class="btn-icon"
            @click="togglePassword(album)"
            :title="album.is_protected ? '移除密码' : '设置密码'"
          >
            {{ album.is_protected ? '🔓' : '🔒' }}
          </button>
          <button class="btn-icon" @click="openPhotoManager(album)" title="管理照片">
            📷
          </button>
          <button class="btn-icon btn-danger" @click="deleteAlbum(album.id)" title="删除">
            🗑️
          </button>
        </div>
      </div>
    </div>

    <!-- 创建相册对话框 -->
    <div v-if="showCreateDialog" class="dialog-overlay" @click="showCreateDialog = false">
      <div class="dialog" @click.stop>
        <h3 class="dialog-title">创建相册</h3>
        <form @submit.prevent="handleCreate">
          <div class="form-group">
            <label>相册名称</label>
            <input v-model="createForm.name" type="text" required placeholder="相册名称" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="createForm.description" rows="3" placeholder="相册描述"></textarea>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="createForm.is_protected" />
              启用密码保护
            </label>
          </div>
          <div v-if="createForm.is_protected" class="form-group">
            <label>访问密码</label>
            <input v-model="createForm.password" type="password" placeholder="输入访问密码" />
          </div>
          <div class="dialog-actions">
            <button type="button" class="btn-secondary" @click="showCreateDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">创建</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑相册对话框 -->
    <div v-if="showEditDialog" class="dialog-overlay" @click="showEditDialog = false">
      <div class="dialog" @click.stop>
        <h3 class="dialog-title">编辑相册</h3>
        <form @submit.prevent="handleEdit">
          <div class="form-group">
            <label>相册名称</label>
            <input v-model="editForm.name" type="text" required />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea v-model="editForm.description" rows="3"></textarea>
          </div>
          <div class="dialog-actions">
            <button type="button" class="btn-secondary" @click="showEditDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">保存</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 设置密码对话框 -->
    <div v-if="showPasswordDialog" class="dialog-overlay" @click="showPasswordDialog = false">
      <div class="dialog" @click.stop>
        <h3 class="dialog-title">{{ passwordForm.is_protected ? '移除密码' : '设置密码' }}</h3>
        <form @submit.prevent="handleSetPassword">
          <div v-if="!passwordForm.is_protected" class="form-group">
            <label>访问密码</label>
            <input v-model="passwordForm.password" type="password" required placeholder="输入访问密码" />
          </div>
          <p v-else>确定要移除相册的密码保护吗？</p>
          <div class="dialog-actions">
            <button type="button" class="btn-secondary" @click="showPasswordDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              {{ passwordForm.is_protected ? '移除' : '设置' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlbumStore } from '../../stores/albums'
import { useNotificationStore } from '../../stores/notification'
import { useConfirm } from '../../composables/useConfirm'
import { getImageUrl } from '../../utils/index'
import PhotoGridSkeleton from '../../components/PhotoGridSkeleton.vue'

const router = useRouter()
const albumStore = useAlbumStore()
const notification = useNotificationStore()
const confirm = useConfirm()
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showPasswordDialog = ref(false)

const createForm = ref({
  name: '',
  description: '',
  is_protected: false,
  password: ''
})

const editForm = ref({
  id: null,
  name: '',
  description: ''
})

const passwordForm = ref({
  id: null,
  is_protected: false,
  password: ''
})

onMounted(() => {
  albumStore.fetchAlbums()
})

const getCoverPhoto = (album) => {
  if (!album.photos || !Array.isArray(album.photos)) {
    return null
  }
  const photo = album.photos.find(p => p.id === album.cover_photo_id)
  if (!photo) {
    return null
  }
  // 返回缩略图或原图路径（转换为完整URL）
  return getImageUrl(photo.thumbnail_path || photo.file_path || '')
}

const handleCreate = async () => {
  try {
    await albumStore.createAlbum(createForm.value)
    if (createForm.value.is_protected && createForm.value.password) {
      const album = albumStore.albums[0]
      await albumStore.setPassword(album.id, createForm.value.password)
    }
    showCreateDialog.value = false
    createForm.value = {
      name: '',
      description: '',
      is_protected: false,
      password: ''
    }
    notification.success('相册创建成功！')
  } catch (error) {
    notification.error('创建失败：' + error.message)
  }
}

const editAlbum = (album) => {
  editForm.value = {
    id: album.id,
    name: album.name,
    description: album.description || ''
  }
  showEditDialog.value = true
}

const handleEdit = async () => {
  try {
    await albumStore.updateAlbum(editForm.value.id, editForm.value)
    showEditDialog.value = false
    notification.success('更新成功！')
  } catch (error) {
    notification.error('更新失败：' + error.message)
  }
}

const openAlbumDetail = (album) => {
  router.push(`/admin/albums/${album.id}`)
}

const openPhotoManager = (album) => {
  router.push(`/admin/albums/${album.id}/photos`)
}

const togglePassword = (album) => {
  passwordForm.value = {
    id: album.id,
    is_protected: album.is_protected,
    password: ''
  }
  showPasswordDialog.value = true
}

const handleSetPassword = async () => {
  try {
    if (passwordForm.value.is_protected) {
      await albumStore.removePassword(passwordForm.value.id)
      notification.success('密码已移除！')
    } else {
      await albumStore.setPassword(passwordForm.value.id, passwordForm.value.password)
      notification.success('密码设置成功！')
    }
    showPasswordDialog.value = false
    await albumStore.fetchAlbums()
  } catch (error) {
    notification.error('操作失败：' + error.message)
  }
}

const deleteAlbum = async (id) => {
  const result = await confirm({
    type: 'danger',
    title: '删除相册',
    message: '确定要删除这个相册吗？相册内的照片不会被删除，只会移除关联关系。',
    confirmText: '删除'
  })

  if (result) {
    try {
      await albumStore.deleteAlbum(id)
      notification.success('删除成功！')
    } catch (error) {
      notification.error('删除失败：' + error.message)
    }
  }
}
</script>

<style scoped>
.album-management {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem;
  font-weight: 300;
  letter-spacing: 0.1em;
}

.btn-primary {
  background: var(--accent-gold);
  color: var(--bg-primary);
  border: none;
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: 0.1em;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-primary:hover {
  background: var(--accent-warm);
  transform: translateY(-2px);
}

.btn-secondary {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-secondary:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
}

.loading,
.empty {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-secondary);
}

.albums-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.album-card {
  background: var(--bg-secondary);
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
  border: 1px solid rgba(201, 169, 98, 0.1);
}

.album-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  border-color: rgba(201, 169, 98, 0.3);
}

.album-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 16/10;
  overflow: hidden;
  cursor: pointer;
}

.album-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s ease;
}

.album-card:hover .album-cover img {
  transform: scale(1.1);
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-gold);
  font-size: 3rem;
  font-family: 'Cormorant Garamond', serif;
}

.album-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7) 0%, transparent 60%);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: var(--spacing-md);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.album-card:hover .album-overlay {
  opacity: 1;
}

.photo-count {
  color: white;
  font-size: 0.9rem;
  font-weight: 500;
}

.album-info {
  padding: var(--spacing-md);
}

.album-name {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.5rem;
  font-weight: 300;
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.album-description {
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.album-badges {
  display: flex;
  gap: var(--spacing-xs);
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.badge-protected {
  background: rgba(201, 169, 98, 0.2);
  color: var(--accent-gold);
}

.badge-public {
  background: rgba(100, 200, 100, 0.2);
  color: #64c864;
}

.album-actions {
  display: flex;
  gap: var(--spacing-xs);
  padding: var(--spacing-sm) var(--spacing-md) var(--spacing-md);
  border-top: 1px solid rgba(201, 169, 98, 0.1);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-xs);
}

.btn-icon {
  background: transparent;
  border: 1px solid rgba(201, 169, 98, 0.2);
  padding: 0.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.btn-icon:hover {
  background: rgba(201, 169, 98, 0.1);
  border-color: var(--accent-gold);
}

.btn-danger:hover {
  background: rgba(255, 0, 0, 0.2);
  border-color: #ff0000;
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

.dialog {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
}

.dialog-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  margin-bottom: var(--spacing-lg);
}

.form-group {
  margin-bottom: var(--spacing-md);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-xs);
  color: var(--text-secondary);
  font-size: 0.9rem;
  letter-spacing: 0.05em;
}

.form-group input,
.form-group textarea {
  width: 100%;
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--accent-gold);
}

.dialog-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-lg);
}

@media (max-width: 768px) {
  .albums-grid {
    grid-template-columns: 1fr;
  }

  .album-cover {
    aspect-ratio: 16/9;
  }
}
</style>
