<template>
  <div class="skeleton-wrapper">
    <!-- Photo Card Skeleton -->
    <div v-if="type === 'photo-card'" class="skeleton-photo-card">
      <div class="skeleton-image"></div>
      <div class="skeleton-content">
        <div class="skeleton-title"></div>
        <div class="skeleton-text skeleton-text-short"></div>
      </div>
    </div>

    <!-- Album Card Skeleton -->
    <div v-else-if="type === 'album-card'" class="skeleton-album-card">
      <div class="skeleton-cover"></div>
      <div class="skeleton-info">
        <div class="skeleton-title"></div>
        <div class="skeleton-text"></div>
      </div>
    </div>

    <!-- List Item Skeleton -->
    <div v-else-if="type === 'list-item'" class="skeleton-list-item">
      <div class="skeleton-avatar"></div>
      <div class="skeleton-content">
        <div class="skeleton-title"></div>
        <div class="skeleton-text"></div>
      </div>
    </div>

    <!-- Table Row Skeleton -->
    <div v-else-if="type === 'table-row'" class="skeleton-table-row">
      <div class="skeleton-cell"></div>
      <div class="skeleton-cell"></div>
      <div class="skeleton-cell"></div>
      <div class="skeleton-cell"></div>
    </div>

    <!-- Text Skeleton -->
    <div v-else-if="type === 'text'" class="skeleton-text-wrapper">
      <div class="skeleton-text"></div>
      <div class="skeleton-text"></div>
      <div class="skeleton-text skeleton-text-short"></div>
    </div>

    <!-- Default Skeleton -->
    <div v-else class="skeleton-default"></div>
  </div>
</template>

<script setup>
defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['photo-card', 'album-card', 'list-item', 'table-row', 'text', 'default'].includes(value)
  }
})
</script>

<style scoped>
.skeleton-wrapper {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* Photo Card Skeleton */
.skeleton-photo-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.skeleton-image {
  width: 100%;
  height: 200px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-content {
  padding: 16px;
}

.skeleton-title {
  height: 20px;
  width: 70%;
  background: #e0e0e0;
  border-radius: 4px;
  margin-bottom: 12px;
}

.skeleton-text {
  height: 16px;
  width: 100%;
  background: #e0e0e0;
  border-radius: 4px;
  margin-bottom: 8px;
}

.skeleton-text-short {
  width: 60%;
}

/* Album Card Skeleton */
.skeleton-album-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.skeleton-cover {
  width: 100%;
  height: 280px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-info {
  padding: 20px;
}

/* List Item Skeleton */
.skeleton-list-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: white;
  border-radius: 8px;
  margin-bottom: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.skeleton-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #e0e0e0;
  margin-right: 16px;
  flex-shrink: 0;
}

/* Table Row Skeleton */
.skeleton-table-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 16px;
  background: white;
  border-bottom: 1px solid #f0f0f0;
}

.skeleton-cell {
  height: 20px;
  background: #e0e0e0;
  border-radius: 4px;
}

/* Text Skeleton */
.skeleton-text-wrapper {
  padding: 16px;
  background: white;
  border-radius: 8px;
}

/* Default Skeleton */
.skeleton-default {
  height: 100px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  border-radius: 8px;
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

/* Responsive */
@media (max-width: 768px) {
  .skeleton-image {
    height: 150px;
  }

  .skeleton-cover {
    height: 200px;
  }

  .skeleton-table-row {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}
</style>
