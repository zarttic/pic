<template>
  <div class="upload-progress">
    <div class="progress-info">
      <span class="progress-text">{{ fileName }}</span>
      <span class="progress-percent">{{ progress }}%</span>
    </div>
    <div class="progress-bar">
      <div
        class="progress-fill"
        :style="{ width: `${progress}%` }"
        :class="{ error: hasError }"
      ></div>
    </div>
    <div v-if="hasError" class="error-message">
      上传失败：{{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
defineProps({
  fileName: {
    type: String,
    default: ''
  },
  progress: {
    type: Number,
    default: 0
  },
  hasError: {
    type: Boolean,
    default: false
  },
  errorMessage: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.upload-progress {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  border-radius: 8px;
  padding: var(--spacing-md);
  margin-top: var(--spacing-sm);
}

.progress-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
  font-size: 0.9rem;
}

.progress-text {
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.progress-percent {
  color: var(--accent-gold);
  font-weight: 600;
  margin-left: var(--spacing-sm);
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: rgba(201, 169, 98, 0.2);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--accent-gold);
  transition: width 0.3s ease;
}

.progress-fill.error {
  background: rgba(220, 53, 69, 0.9);
}

.error-message {
  color: rgba(220, 53, 69, 1);
  font-size: 0.85rem;
  margin-top: var(--spacing-xs);
}
</style>
