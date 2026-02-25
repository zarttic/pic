<template>
  <div class="statistics-page">
    <div class="page-header">
      <h2 class="page-title">统计数据</h2>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📸</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalPhotos }}</div>
          <div class="stat-label">总照片数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalAlbums }}</div>
          <div class="stat-label">总相册数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">👁️</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.totalViews }}</div>
          <div class="stat-label">总浏览量</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🔒</div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.protectedAlbums }}</div>
          <div class="stat-label">加密相册</div>
        </div>
      </div>
    </div>

    <div class="charts-section">
      <div class="chart-card">
        <h3 class="chart-title">热门照片 TOP 10</h3>
        <div v-if="topPhotos.length === 0" class="empty-chart">
          暂无数据
        </div>
        <div v-else class="photo-list">
          <div v-for="(photo, index) in topPhotos" :key="photo.id" class="photo-item">
            <div class="rank">{{ index + 1 }}</div>
            <img :src="getImageUrl(photo.thumbnail_path || photo.file_path)" :alt="photo.title" class="photo-thumb" />
            <div class="photo-info">
              <div class="photo-title">{{ photo.title }}</div>
              <div class="photo-location">{{ photo.location || '未知地点' }}</div>
            </div>
            <div class="photo-views">
              <span class="view-count">{{ photo.view_count }}</span>
              <span class="view-label">次浏览</span>
            </div>
          </div>
        </div>
      </div>

      <div class="chart-card">
        <h3 class="chart-title">最近上传</h3>
        <div v-if="recentPhotos.length === 0" class="empty-chart">
          暂无数据
        </div>
        <div v-else class="photo-list">
          <div v-for="photo in recentPhotos" :key="photo.id" class="photo-item">
            <img :src="getImageUrl(photo.thumbnail_path || photo.file_path)" :alt="photo.title" class="photo-thumb" />
            <div class="photo-info">
              <div class="photo-title">{{ photo.title }}</div>
              <div class="photo-date">{{ formatDate(photo.created_at) }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { usePhotoStore } from '../../stores/photos'
import { useAlbumStore } from '../../stores/albums'
import { getImageUrl } from '../../utils/index'

const photoStore = usePhotoStore()
const albumStore = useAlbumStore()

onMounted(() => {
  photoStore.fetchPhotos()
  albumStore.fetchAlbums()
})

const stats = computed(() => {
  const photos = photoStore.photos || []
  const albums = albumStore.albums || []

  return {
    totalPhotos: photos.length,
    totalAlbums: albums.length,
    totalViews: photos.reduce((sum, p) => sum + (p.view_count || 0), 0),
    protectedAlbums: albums.filter(a => a.is_protected).length
  }
})

const topPhotos = computed(() => {
  const photos = [...(photoStore.photos || [])]
  return photos
    .sort((a, b) => (b.view_count || 0) - (a.view_count || 0))
    .slice(0, 10)
})

const recentPhotos = computed(() => {
  const photos = [...(photoStore.photos || [])]
  return photos
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    .slice(0, 10)
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}
</script>

<style scoped>
.statistics-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 2rem;
  font-weight: 300;
  letter-spacing: 0.1em;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-icon {
  font-size: 3rem;
  opacity: 0.8;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 300;
  color: var(--accent-gold);
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  letter-spacing: 0.1em;
}

.charts-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: var(--spacing-lg);
}

.chart-card {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
}

.chart-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.5rem;
  font-weight: 300;
  margin-bottom: var(--spacing-lg);
  letter-spacing: 0.1em;
}

.empty-chart {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-secondary);
}

.photo-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.photo-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--bg-primary);
  border-radius: 8px;
  transition: transform 0.2s ease;
}

.photo-item:hover {
  transform: translateX(4px);
}

.rank {
  font-size: 1.5rem;
  font-weight: 300;
  color: var(--accent-gold);
  min-width: 30px;
  text-align: center;
}

.photo-thumb {
  width: 60px;
  height: 45px;
  object-fit: cover;
  border-radius: 4px;
}

.photo-info {
  flex: 1;
}

.photo-title {
  font-size: 1rem;
  margin-bottom: var(--spacing-xs);
}

.photo-location,
.photo-date {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.photo-views {
  text-align: right;
}

.view-count {
  display: block;
  font-size: 1.2rem;
  color: var(--accent-gold);
  font-weight: 500;
}

.view-label {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .charts-section {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
