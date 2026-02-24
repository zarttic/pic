<template>
  <div class="gallery">
    <section class="gallery-section">
      <div class="section-header">
        <h2 class="section-title">精选作品</h2>
        <div class="section-divider"></div>
      </div>

      <div v-if="photoStore.loading" class="loading">
        加载中...
      </div>

      <div v-else-if="photoStore.error" class="error">
        {{ photoStore.error }}
      </div>

      <div v-else class="gallery-grid">
        <div
          v-for="photo in photoStore.photos"
          :key="photo.id"
          class="gallery-item"
          @click="openViewer(photo)"
        >
          <img :src="photo.file_path" :alt="photo.title" />
          <div class="gallery-overlay">
            <h3 class="photo-title">{{ photo.title }}</h3>
            <p class="photo-meta">{{ photo.location }} · {{ photo.year }}</p>
          </div>
        </div>
      </div>
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
import { usePhotoStore } from '../stores/photos'

const photoStore = usePhotoStore()
const selectedPhoto = ref(null)

onMounted(() => {
  photoStore.fetchPhotos()
  // 触发滚动动画
  observeElements()
})

const openViewer = (photo) => {
  selectedPhoto.value = photo
  document.body.style.overflow = 'hidden'
}

const closeViewer = () => {
  selectedPhoto.value = null
  document.body.style.overflow = ''
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
    document.querySelectorAll('.section-header, .gallery-item').forEach((el) => {
      observer.observe(el)
    })
  }, 100)
}
</script>

<style scoped>
.gallery {
  min-height: 100vh;
  padding-top: 100px;
}

.gallery-section {
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
.error {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-secondary);
}

.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: var(--spacing-md);
  max-width: 1600px;
  margin: 0 auto;
}

.gallery-item {
  position: relative;
  aspect-ratio: 3/2;
  overflow: hidden;
  cursor: pointer;
  opacity: 0;
  transform: translateY(40px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.gallery-item.visible {
  opacity: 1;
  transform: translateY(0);
}

.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s ease;
}

.gallery-item:hover img {
  transform: scale(1.08);
}

.gallery-overlay {
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

.gallery-item:hover .gallery-overlay {
  opacity: 1;
}

.photo-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  margin-bottom: var(--spacing-xs);
  transform: translateY(20px);
  transition: transform 0.4s ease;
}

.gallery-item:hover .photo-title {
  transform: translateY(0);
}

.photo-meta {
  font-size: 0.85rem;
  color: var(--text-secondary);
  font-weight: 200;
  letter-spacing: 0.1em;
  transform: translateY(20px);
  transition: transform 0.4s ease 0.1s;
}

.gallery-item:hover .photo-meta {
  transform: translateY(0);
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
  .gallery-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-sm);
  }
}
</style>
