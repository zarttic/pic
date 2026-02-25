<template>
  <div class="settings-skeleton">
    <div class="skeleton-card">
      <div class="skeleton-title"></div>
      <div class="skeleton-form">
        <div class="skeleton-field">
          <div class="skeleton-label"></div>
          <div class="skeleton-input"></div>
        </div>
        <div class="skeleton-field">
          <div class="skeleton-label"></div>
          <div class="skeleton-textarea"></div>
        </div>
        <div class="skeleton-field">
          <div class="skeleton-label"></div>
          <div class="skeleton-input"></div>
        </div>
        <div class="skeleton-button"></div>
      </div>
    </div>

    <div class="skeleton-card">
      <div class="skeleton-title"></div>
      <div class="skeleton-form">
        <div class="skeleton-field">
          <div class="skeleton-label"></div>
          <div class="skeleton-input skeleton-input-short"></div>
        </div>
        <div class="skeleton-field">
          <div class="skeleton-label"></div>
          <div class="skeleton-input"></div>
        </div>
        <div class="skeleton-checkbox"></div>
        <div class="skeleton-checkbox"></div>
      </div>
    </div>

    <div class="skeleton-card">
      <div class="skeleton-title"></div>
      <div class="skeleton-info-list">
        <div class="skeleton-info-item">
          <div class="skeleton-info-label"></div>
          <div class="skeleton-info-value"></div>
        </div>
        <div class="skeleton-info-item">
          <div class="skeleton-info-label"></div>
          <div class="skeleton-info-value"></div>
        </div>
        <div class="skeleton-info-item">
          <div class="skeleton-info-label"></div>
          <div class="skeleton-info-value"></div>
        </div>
        <div class="skeleton-info-item">
          <div class="skeleton-info-label"></div>
          <div class="skeleton-info-value"></div>
        </div>
      </div>
    </div>

    <div class="skeleton-card skeleton-danger">
      <div class="skeleton-title skeleton-title-danger"></div>
      <div class="skeleton-warning-text"></div>
      <div class="skeleton-danger-actions">
        <div class="skeleton-button skeleton-button-danger"></div>
        <div class="skeleton-button skeleton-button-danger"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
// 无需脚本逻辑
</script>

<style scoped>
.settings-skeleton {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  max-width: 800px;
  margin: 0 auto;
}

.skeleton-card {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
}

.skeleton-danger {
  border: 1px solid rgba(255, 107, 107, 0.3);
}

.skeleton-title {
  height: 32px;
  width: 150px;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-bottom: var(--spacing-lg);
}

.skeleton-title-danger {
  width: 120px;
}

.skeleton-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.skeleton-field {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.skeleton-label {
  height: 18px;
  width: 100px;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-input {
  height: 42px;
  width: 100%;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-input-short {
  width: 200px;
}

.skeleton-textarea {
  height: 80px;
  width: 100%;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-checkbox {
  height: 24px;
  width: 180px;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-button {
  height: 42px;
  width: 120px;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-top: var(--spacing-md);
}

.skeleton-button-danger {
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(255, 107, 107, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
}

.skeleton-info-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.skeleton-info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--spacing-md);
  background: var(--bg-primary);
  border-radius: 4px;
}

.skeleton-info-label {
  height: 20px;
  width: 80px;
  background: linear-gradient(90deg, var(--bg-secondary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-secondary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-info-value {
  height: 20px;
  width: 100px;
  background: linear-gradient(90deg, var(--bg-secondary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-secondary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

.skeleton-warning-text {
  height: 20px;
  width: 200px;
  background: linear-gradient(90deg, var(--bg-primary) 0%, rgba(201, 169, 98, 0.1) 50%, var(--bg-primary) 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
  margin-bottom: var(--spacing-md);
}

.skeleton-danger-actions {
  display: flex;
  gap: var(--spacing-md);
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

@media (max-width: 768px) {
  .skeleton-danger-actions {
    flex-direction: column;
  }

  .skeleton-button-danger {
    width: 100%;
  }
}
</style>
