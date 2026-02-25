<template>
  <div class="album-detail-page">
    <!-- 密码验证对话框 -->
    <div v-if="showPasswordDialog" class="password-dialog-overlay" @click="closePasswordDialog">
      <div class="password-dialog" @click.stop>
        <h3 class="dialog-title">🔒 相册需要密码访问</h3>
        <p class="dialog-subtitle">{{ albumStore.currentAlbum?.name }}</p>
        <form @submit.prevent="handleVerifyPassword" class="password-form">
          <div class="form-group">
            <label for="password">请输入访问密码</label>
            <input
              id="password"
              v-model="passwordInput"
              type="password"
              placeholder="输入密码"
              required
              autofocus
            />
          </div>
          <div v-if="passwordError" class="error-message">
            {{ passwordError }}
          </div>
          <div class="dialog-actions">
            <button type="button" class="btn-secondary" @click="closePasswordDialog">
              取消
            </button>
            <button type="submit" class="btn-primary" :disabled="verifying">
              {{ verifying ? '验证中...' : '确认' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <section class="album-section">
      <div v-if="albumStore.loading" class="loading">
        加载中...
      </div>

      <template v-else-if="albumStore.currentAlbum">
        <div class="album-header">
          <button class="back-button" @click="goBack">
            ← 返回
          </button>
          <h2 class="album-title">
            {{ albumStore.currentAlbum.name }}
            <span v-if="albumStore.currentAlbum.is_protected" class="protected-badge">🔒</span>
          </h2>
          <p v-if="albumStore.currentAlbum.description" class="album-description">
            {{ albumStore.currentAlbum.description }}
          </p>
          <div class="section-divider"></div>
        </div>

        <div v-if="albumStore.currentAlbum.require_auth" class="empty">
          <p>🔒 此相册需要密码访问</p>
          <button class="btn-primary" @click="openPasswordDialog">
            输入密码
          </button>
        </div>

        <div v-else-if="!albumStore.currentAlbum.photos || albumStore.currentAlbum.photos.length === 0" class="empty">
          相册中暂无照片
        </div>

        <div v-else class="photos-grid">
          <div
            v-for="photo in albumStore.currentAlbum.photos"
            :key="photo.id"
            class="photo-item"
            @click="openViewer(photo)"
          >
            <img :src="getImageUrl(photo.thumbnail_path || photo.file_path)" :alt="photo.title" />
            <div class="photo-overlay">
              <h3 class="photo-title">{{ photo.title }}</h3>
              <p class="photo-meta">{{ photo.location }} · {{ photo.year }}</p>
            </div>
          </div>
        </div>
      </template>
    </section>

    <!-- 全屏查看器 -->
    <div v-if="selectedPhoto" class="viewer active" @click="closeViewer">
      <button class="viewer-close" @click="closeViewer">×</button>
      <div class="viewer-content" @click.stop>
        <img :src="getImageUrl(selectedPhoto.file_path)" :alt="selectedPhoto.title" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAlbumStore } from '../stores/albums'
import { getImageUrl } from '../utils/index'

const route = useRoute()
const router = useRouter()
const albumStore = useAlbumStore()
const selectedPhoto = ref(null)
const showPasswordDialog = ref(false)
const passwordInput = ref('')
const passwordError = ref('')
const verifying = ref(false)

onMounted(async () => {
  const albumId = route.params.id
  await albumStore.fetchAlbum(albumId)

  // 如果需要验证，检查是否有保存的 token
  if (albumStore.currentAlbum?.require_auth) {
    const token = localStorage.getItem(`album_token_${albumId}`)
    if (token) {
      // 尝试使用保存的 token 获取相册数据
      await fetchAlbumWithToken(albumId, token)
    } else {
      showPasswordDialog.value = true
    }
  }
})

const fetchAlbumWithToken = async (albumId, token) => {
  try {
    const response = await fetch(`/api/albums/${albumId}`, {
      headers: {
        'X-Album-Token': token
      }
    })
    const data = await response.json()
    if (!data.require_auth) {
      albumStore.currentAlbum = data
    } else {
      showPasswordDialog.value = true
    }
  } catch (error) {
    console.error('Error fetching album with token:', error)
    showPasswordDialog.value = true
  }
}

const goBack = () => {
  router.push('/albums')
}

const openViewer = (photo) => {
  selectedPhoto.value = photo
  document.body.style.overflow = 'hidden'
}

const closeViewer = () => {
  selectedPhoto.value = null
  document.body.style.overflow = ''
}

const openPasswordDialog = () => {
  showPasswordDialog.value = true
  passwordInput.value = ''
  passwordError.value = ''
}

const closePasswordDialog = () => {
  showPasswordDialog.value = false
  passwordInput.value = ''
  passwordError.value = ''
  if (albumStore.currentAlbum?.require_auth) {
    router.push('/albums')
  }
}

const handleVerifyPassword = async () => {
  if (!passwordInput.value) return

  verifying.value = true
  passwordError.value = ''

  try {
    await albumStore.verifyPassword(route.params.id, passwordInput.value)
    // 验证成功，重新获取相册数据
    const token = localStorage.getItem(`album_token_${route.params.id}`)
    await fetchAlbumWithToken(route.params.id, token)
    closePasswordDialog()
  } catch (error) {
    passwordError.value = '密码错误，请重试'
  } finally {
    verifying.value = false
  }
}
</script>

<style scoped>
.album-detail-page {
  min-height: 100vh;
  padding-top: 100px;
}

/* 密码对话框 */
.password-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.password-dialog {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
  max-width: 400px;
  width: 90%;
}

.dialog-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem;
  font-weight: 300;
  margin-bottom: var(--spacing-sm);
  text-align: center;
}

.dialog-subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);
}

.password-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-group label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  letter-spacing: 0.1em;
}

.form-group input {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-gold);
}

.error-message {
  color: #ff6b6b;
  font-size: 0.9rem;
  text-align: center;
}

.dialog-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

.btn-primary {
  background: var(--accent-gold);
  color: var(--bg-primary);
  border: none;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
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

.album-section {
  padding: var(--spacing-xl) var(--spacing-lg);
  min-height: 100vh;
}

.album-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
}

.back-button {
  background: transparent;
  border: 1px solid var(--text-secondary);
  color: var(--text-primary);
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  margin-bottom: var(--spacing-md);
}

.back-button:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
}

.album-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(2.5rem, 6vw, 5rem);
  font-weight: 300;
  letter-spacing: 0.15em;
  margin-bottom: var(--spacing-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
}

.protected-badge {
  font-size: 2rem;
}

.album-description {
  font-size: 1.1rem;
  color: var(--text-secondary);
  max-width: 600px;
  margin: 0 auto var(--spacing-md);
  line-height: 1.8;
}

.section-divider {
  width: 60px;
  height: 1px;
  background: var(--accent-gold);
  margin: var(--spacing-md) auto;
}

.loading,
.empty {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-secondary);
}

.empty p {
  margin-bottom: var(--spacing-md);
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: var(--spacing-md);
  max-width: 1600px;
  margin: 0 auto;
}

.photo-item {
  position: relative;
  aspect-ratio: 3/2;
  overflow: hidden;
  cursor: pointer;
}

.photo-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s ease;
}

.photo-item:hover img {
  transform: scale(1.08);
}

.photo-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(10, 10, 10, 0.9) 0%, transparent 60%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: var(--spacing-md);
  opacity: 0;
  transition: opacity 0.4s ease;
}

.photo-item:hover .photo-overlay {
  opacity: 1;
}

.photo-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  margin-bottom: var(--spacing-xs);
}

.photo-meta {
  font-size: 0.85rem;
  color: var(--text-secondary);
  font-weight: 200;
  letter-spacing: 0.1em;
}

/* 全屏查看器 */
.viewer {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.97);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.4s ease, visibility 0.4s ease;
}

.viewer.active {
  opacity: 1;
  visibility: visible;
}

.viewer-content {
  max-width: 90vw;
  max-height: 90vh;
  position: relative;
  transform: scale(0.9);
  transition: transform 0.4s ease;
}

.viewer.active .viewer-content {
  transform: scale(1);
}

.viewer-content img {
  max-width: 100%;
  max-height: 90vh;
  display: block;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.viewer-close {
  position: absolute;
  top: var(--spacing-lg);
  right: var(--spacing-lg);
  width: 50px;
  height: 50px;
  border: 1px solid var(--text-secondary);
  background: transparent;
  color: var(--text-primary);
  font-size: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.viewer-close:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
  transform: rotate(90deg);
}

@media (max-width: 768px) {
  .photos-grid {
    grid-template-columns: 1fr;
  }
}
</style>
