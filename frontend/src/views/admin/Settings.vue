<template>
  <div class="settings-page">
    <div class="page-header">
      <h2 class="page-title">系统设置</h2>
    </div>

    <div class="settings-sections">
      <div class="settings-card">
        <h3 class="section-title">网站信息</h3>
        <form @submit.prevent="saveSiteSettings" class="settings-form">
          <div class="form-group">
            <label>网站名称</label>
            <input v-model="siteSettings.title" type="text" placeholder="PicSite" />
          </div>
          <div class="form-group">
            <label>网站描述</label>
            <textarea v-model="siteSettings.description" rows="3" placeholder="一个精美的摄影作品展示平台"></textarea>
          </div>
          <div class="form-group">
            <label>作者名称</label>
            <input v-model="siteSettings.author" type="text" placeholder="摄影师姓名" />
          </div>
          <button type="submit" class="btn-primary">保存设置</button>
        </form>
      </div>

      <div class="settings-card">
        <h3 class="section-title">上传设置</h3>
        <form class="settings-form">
          <div class="form-group">
            <label>最大上传大小 (MB)</label>
            <input type="number" value="10" disabled />
            <span class="hint">服务端限制，需修改后端配置</span>
          </div>
          <div class="form-group">
            <label>支持格式</label>
            <input type="text" value="JPG, PNG, WebP" disabled />
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" checked disabled />
              自动生成缩略图
            </label>
          </div>
          <div class="form-group">
            <label>
              <input type="checkbox" checked disabled />
              自动提取 EXIF 信息
            </label>
          </div>
        </form>
      </div>

      <div class="settings-card">
        <h3 class="section-title">系统信息</h3>
        <div class="info-list">
          <div class="info-item">
            <span class="info-label">版本</span>
            <span class="info-value">v0.3.0</span>
          </div>
          <div class="info-item">
            <span class="info-label">前端框架</span>
            <span class="info-value">Vue 3 + Vite</span>
          </div>
          <div class="info-item">
            <span class="info-label">后端框架</span>
            <span class="info-value">Go + Gin</span>
          </div>
          <div class="info-item">
            <span class="info-label">数据库</span>
            <span class="info-value">SQLite</span>
          </div>
        </div>
      </div>

      <div class="settings-card danger-zone">
        <h3 class="section-title">危险操作</h3>
        <p class="warning-text">以下操作不可逆，请谨慎操作</p>
        <div class="danger-actions">
          <button class="btn-danger" @click="clearCache">
            清空缓存
          </button>
          <button class="btn-danger" @click="exportData">
            导出数据
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const siteSettings = ref({
  title: 'PicSite',
  description: '一个精美的摄影作品展示平台',
  author: ''
})

const saveSiteSettings = () => {
  localStorage.setItem('siteSettings', JSON.stringify(siteSettings.value))
  alert('设置已保存！')
}

const clearCache = () => {
  if (confirm('确定要清空缓存吗？这不会删除照片和相册数据。')) {
    localStorage.clear()
    alert('缓存已清空！')
  }
}

const exportData = () => {
  alert('数据导出功能开发中...')
}
</script>

<style scoped>
.settings-page {
  max-width: 800px;
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

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.settings-card {
  background: var(--bg-secondary);
  padding: var(--spacing-xl);
  border-radius: 8px;
}

.section-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.5rem;
  font-weight: 300;
  margin-bottom: var(--spacing-lg);
  letter-spacing: 0.1em;
}

.settings-form {
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
  letter-spacing: 0.05em;
}

.form-group input,
.form-group textarea {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 1rem;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--accent-gold);
}

.form-group input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.hint {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-style: italic;
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
  margin-top: var(--spacing-md);
}

.btn-primary:hover {
  background: var(--accent-warm);
  transform: translateY(-2px);
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: var(--spacing-md);
  background: var(--bg-primary);
  border-radius: 4px;
}

.info-label {
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
  font-weight: 500;
}

.danger-zone {
  border: 1px solid rgba(255, 0, 0, 0.3);
}

.danger-zone .section-title {
  color: #ff6b6b;
}

.warning-text {
  color: var(--text-secondary);
  margin-bottom: var(--spacing-md);
  font-size: 0.9rem;
}

.danger-actions {
  display: flex;
  gap: var(--spacing-md);
}

.btn-danger {
  background: transparent;
  color: #ff6b6b;
  border: 1px solid #ff6b6b;
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.btn-danger:hover {
  background: rgba(255, 107, 107, 0.1);
}
</style>
