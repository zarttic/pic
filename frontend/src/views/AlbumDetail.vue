<template>
  <div class="album-detail-page">
    <!-- 密码验证对话框 -->
    <v-dialog
      v-model="showPasswordDialog"
      persistent
      max-width="500"
      transition="dialog-bottom-transition"
    >
      <v-card class="password-card">
        <!-- 装饰背景 -->
        <div class="card-decoration">
          <div class="decoration-circle"></div>
          <div class="decoration-circle"></div>
        </div>

        <!-- 图标 -->
        <div class="lock-icon-wrapper">
          <v-icon size="80" color="primary" class="lock-icon">mdi-lock-outline</v-icon>
        </div>

        <!-- 标题区域 -->
        <v-card-title class="text-center pa-0 mb-2">
          <h2 class="dialog-title">私密相册</h2>
          <p class="dialog-subtitle">{{ albumStore.currentAlbum?.name }}</p>
        </v-card-title>

        <!-- 表单区域 -->
        <v-card-text class="pa-0">
          <v-form @submit.prevent="handleVerifyPassword" class="password-form">
            <v-text-field
              v-model="passwordInput"
              :type="showPassword ? 'text' : 'password'"
              label="访问密码"
              placeholder="请输入相册访问密码"
              variant="outlined"
              density="comfortable"
              prepend-inner-icon="mdi-key-outline"
              :append-inner-icon="showPassword ? 'mdi-eye-outline' : 'mdi-eye-off-outline'"
              @click:append-inner="showPassword = !showPassword"
              :error-messages="passwordError ? [passwordError] : []"
              :disabled="verifying"
              autofocus
              hide-details="auto"
              class="password-input"
            ></v-text-field>

            <v-expand-transition>
              <v-alert
                v-if="passwordError"
                type="error"
                variant="tonal"
                density="compact"
                class="mt-3"
                closable
                @click:close="passwordError = ''"
              >
                {{ passwordError }}
              </v-alert>
            </v-expand-transition>
          </v-form>
        </v-card-text>

        <!-- 操作按钮 -->
        <v-card-actions class="pa-0 mt-6">
          <v-spacer></v-spacer>
          <v-btn
            variant="outlined"
            color="grey"
            @click="closePasswordDialog"
            :disabled="verifying"
            class="mr-2"
          >
            返回
          </v-btn>
          <v-btn
            variant="flat"
            color="primary"
            @click="handleVerifyPassword"
            :loading="verifying"
            :disabled="!passwordInput"
            class="submit-btn"
          >
            验证访问
          </v-btn>
        </v-card-actions>

        <!-- 提示文字 -->
        <p class="hint-text">
          <v-icon size="small" color="grey">mdi-information-outline</v-icon>
          此相册受密码保护，请联系摄影师获取访问权限
        </p>
      </v-card>
    </v-dialog>

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
import api from '../api'

const route = useRoute()
const router = useRouter()
const albumStore = useAlbumStore()
const selectedPhoto = ref(null)
const showPasswordDialog = ref(false)
const passwordInput = ref('')
const passwordError = ref('')
const verifying = ref(false)
const showPassword = ref(false)

onMounted(async () => {
  const albumId = route.params.id
  // 使用 fetchAlbumPublic 避免发送 JWT token
  await albumStore.fetchAlbumPublic(albumId)

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
    // 前台访问相册时，创建一个新的 axios 实例，不使用 JWT token
    const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080/api'}/albums/${albumId}`, {
      headers: {
        'X-Album-Token': token
      }
    })
    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Failed to fetch album')
    }

    // 如果返回的数据不包含 require_auth 或 require_auth 为 false,说明验证成功
    if (!data.require_auth) {
      albumStore.currentAlbum = data
      showPasswordDialog.value = false
      return true
    } else {
      showPasswordDialog.value = true
      return false
    }
  } catch (error) {
    console.error('Error fetching album with token:', error)
    showPasswordDialog.value = true
    return false
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
    // 成功获取数据后关闭对话框(不调用 closePasswordDialog,避免跳转)
    showPasswordDialog.value = false
    passwordInput.value = ''
    passwordError.value = ''
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

/* 密码对话框样式 */
.password-card {
  background: linear-gradient(135deg, rgba(20, 20, 20, 0.98) 0%, rgba(10, 10, 10, 0.98) 100%) !important;
  border: 1px solid rgba(201, 169, 98, 0.2);
  padding: 3rem 2.5rem;
  position: relative;
  overflow: hidden;
}

/* 装饰背景 */
.card-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  pointer-events: none;
  overflow: hidden;
}

.decoration-circle {
  position: absolute;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(201, 169, 98, 0.08) 0%, transparent 70%);
  filter: blur(40px);
}

.decoration-circle:first-child {
  top: -150px;
  left: -100px;
}

.decoration-circle:last-child {
  top: -100px;
  right: -100px;
}

/* 锁图标 */
.lock-icon-wrapper {
  text-align: center;
  margin-bottom: 1.5rem;
  position: relative;
  z-index: 1;
}

.lock-icon {
  opacity: 0.9;
  filter: drop-shadow(0 4px 12px rgba(201, 169, 98, 0.3));
}

/* 标题样式 */
.dialog-title {
  font-family: 'Cormorant Garamond', serif !important;
  font-size: 2.5rem !important;
  font-weight: 300 !important;
  letter-spacing: 0.15em;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.dialog-subtitle {
  font-size: 1rem;
  color: var(--accent-gold);
  opacity: 0.8;
  letter-spacing: 0.1em;
  margin: 0;
}

/* 表单样式 */
.password-form {
  position: relative;
  z-index: 1;
}

.password-input :deep(.v-field) {
  background: rgba(10, 10, 10, 0.6);
  border: 1px solid rgba(201, 169, 98, 0.3);
  transition: all 0.3s ease;
}

.password-input :deep(.v-field:hover) {
  border-color: rgba(201, 169, 98, 0.5);
}

.password-input :deep(.v-field--focused) {
  border-color: var(--accent-gold);
  box-shadow: 0 0 0 2px rgba(201, 169, 98, 0.1);
}

.password-input :deep(.v-label) {
  color: var(--text-secondary);
  letter-spacing: 0.05em;
}

.password-input :deep(input) {
  color: var(--text-primary);
  letter-spacing: 0.1em;
}

.password-input :deep(input::placeholder) {
  color: var(--text-secondary);
  opacity: 0.5;
}

/* 按钮样式 */
.submit-btn {
  min-width: 120px;
  letter-spacing: 0.1em;
}

/* 提示文字 */
.hint-text {
  text-align: center;
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin-top: 1.5rem;
  margin-bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  opacity: 0.7;
}

/* 相册区域样式 */.album-section {
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
