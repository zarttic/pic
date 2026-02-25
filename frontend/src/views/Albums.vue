<template>
  <div class="albums-page">
    <section class="albums-section">
      <div class="section-header">
        <h2 class="section-title">相册集</h2>
        <div class="section-divider"></div>
      </div>

      <div v-if="albumStore.loading" class="loading">
        加载中...
      </div>

      <div v-else-if="albumStore.error" class="error">
        {{ albumStore.error }}
      </div>

      <div v-else-if="albumStore.albums.length === 0" class="empty">
        暂无相册
      </div>

      <div v-else class="albums-grid">
        <div
          v-for="album in albumStore.albums"
          :key="album.id"
          class="album-card"
          @click="viewAlbum(album.id)"
        >
          <div class="album-cover">
            <img
              v-if="album.cover_photo_id && getCoverPhoto(album)"
              :src="getCoverPhoto(album)"
              :alt="album.name"
            />
            <div v-else class="album-cover-placeholder">
              <span>{{ album.name.charAt(0) }}</span>
            </div>
            <div class="album-overlay">
              <h3 class="album-name">
                {{ album.name }}
                <span v-if="album.is_protected" class="lock-icon">🔒</span>
              </h3>
              <p class="album-count">{{ album.photos?.length || 0 }} 张照片</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlbumStore } from '../stores/albums'
import { getImageUrl } from '../utils/index'

const router = useRouter()
const albumStore = useAlbumStore()

onMounted(() => {
  albumStore.fetchAlbums()
  observeElements()
})

const viewAlbum = (id) => {
  router.push(`/albums/${id}`)
}

const getCoverPhoto = (album) => {
  if (!album.photos || !Array.isArray(album.photos)) {
    return null
  }
  const photo = album.photos.find(p => p.id === album.cover_photo_id)
  if (!photo) {
    return null
  }
  return getImageUrl(photo.thumbnail_path || photo.file_path || '')
}

const observeElements = () => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    {
      threshold: 0.1,
      rootMargin: '0px 0px -100px 0px'
    }
  )

  setTimeout(() => {
    document.querySelectorAll('.section-header, .album-card').forEach((el) => {
      observer.observe(el)
    })
  }, 100)
}
</script>

<style scoped>
.albums-page {
  min-height: 100vh;
  padding-top: 100px;
}

.albums-section {
  padding: var(--spacing-xl) var(--spacing-lg);
  min-height: 100vh;
}

.section-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  opacity: 0;
  transform: translateY(30px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.section-header.visible {
  opacity: 1;
  transform: translateY(0);
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

.loading,
.error,
.empty {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-secondary);
}

.albums-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--spacing-lg);
  max-width: 1400px;
  margin: 0 auto;
}

.album-card {
  cursor: pointer;
  opacity: 0;
  transform: translateY(40px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.album-card.visible {
  opacity: 1;
  transform: translateY(0);
}

.album-cover {
  position: relative;
  aspect-ratio: 4/3;
  overflow: hidden;
  border-radius: 8px;
  background: var(--bg-secondary);
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

.album-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
  font-family: 'Cormorant Garamond', serif;
  font-size: 5rem;
  color: var(--accent-gold);
  opacity: 0.3;
}

.album-overlay {
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

.album-card:hover .album-overlay {
  opacity: 1;
}

.album-name {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  margin-bottom: var(--spacing-xs);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.lock-icon {
  font-size: 1.2rem;
}

.album-count {
  font-size: 0.9rem;
  color: var(--text-secondary);
  letter-spacing: 0.1em;
}

@media (max-width: 768px) {
  .albums-grid {
    grid-template-columns: 1fr;
  }
}
</style>
