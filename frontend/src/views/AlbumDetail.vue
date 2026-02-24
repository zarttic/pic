<template>
  <div class="album-detail-page">
    <section class="album-section">
      <div v-if="albumStore.loading" class="loading">
        加载中...
      </div>

      <template v-else-if="albumStore.currentAlbum">
        <div class="album-header">
          <button class="back-button" @click="goBack">
            ← 返回
          </button>
          <h2 class="album-title">{{ albumStore.currentAlbum.name }}</h2>
          <p v-if="albumStore.currentAlbum.description" class="album-description">
            {{ albumStore.currentAlbum.description }}
          </p>
          <div class="section-divider"></div>
        </div>

        <div v-if="albumStore.currentAlbum.photos?.length === 0" class="empty">
          相册中暂无照片
        </div>

        <div v-else class="photos-grid">
          <div
            v-for="photo in albumStore.currentAlbum.photos"
            :key="photo.id"
            class="photo-item"
            @click="openViewer(photo)"
          >
            <img :src="photo.thumbnail_path || photo.file_path" :alt="photo.title" />
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
        <img :src="selectedPhoto.file_path" :alt="selectedPhoto.title" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAlbumStore } from '../stores/albums'

const route = useRoute()
const router = useRouter()
const albumStore = useAlbumStore()
const selectedPhoto = ref(null)

onMounted(() => {
  const albumId = route.params.id
  albumStore.fetchAlbum(albumId)
})

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
</script>

<style scoped>
.album-detail-page {
  min-height: 100vh;
  padding-top: 100px;
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
