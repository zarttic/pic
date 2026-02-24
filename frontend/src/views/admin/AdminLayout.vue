<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="admin-logo">
        <h2>PicSite 管理</h2>
      </div>
      <nav class="admin-nav">
        <router-link to="/admin/photos" class="nav-item">
          <span class="icon">📸</span>
          <span>照片管理</span>
        </router-link>
        <router-link to="/admin/albums" class="nav-item">
          <span class="icon">📁</span>
          <span>相册管理</span>
        </router-link>
        <router-link to="/admin/statistics" class="nav-item">
          <span class="icon">📊</span>
          <span>统计数据</span>
        </router-link>
        <router-link to="/admin/settings" class="nav-item">
          <span class="icon">⚙️</span>
          <span>系统设置</span>
        </router-link>
      </nav>
      <div class="admin-footer">
        <router-link to="/" class="back-link">
          ← 返回前台
        </router-link>
      </div>
    </aside>
    <main class="admin-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

.admin-sidebar {
  width: 260px;
  background: var(--bg-secondary);
  border-right: 1px solid rgba(201, 169, 98, 0.1);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  overflow-y: auto;
}

.admin-logo {
  padding: var(--spacing-xl) var(--spacing-lg);
  border-bottom: 1px solid rgba(201, 169, 98, 0.1);
}

.admin-logo h2 {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  color: var(--accent-gold);
  letter-spacing: 0.1em;
}

.admin-nav {
  flex: 1;
  padding: var(--spacing-md) 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-lg);
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
}

.nav-item:hover {
  background: rgba(201, 169, 98, 0.05);
  color: var(--text-primary);
}

.nav-item.router-link-active {
  background: rgba(201, 169, 98, 0.1);
  color: var(--accent-gold);
  border-left-color: var(--accent-gold);
}

.nav-item .icon {
  font-size: 1.2rem;
}

.admin-footer {
  padding: var(--spacing-lg);
  border-top: 1px solid rgba(201, 169, 98, 0.1);
}

.back-link {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.3s ease;
}

.back-link:hover {
  color: var(--accent-gold);
}

.admin-content {
  flex: 1;
  margin-left: 260px;
  padding: var(--spacing-xl);
  overflow-y: auto;
}

@media (max-width: 1024px) {
  .admin-sidebar {
    width: 200px;
  }

  .admin-content {
    margin-left: 200px;
  }
}

@media (max-width: 768px) {
  .admin-sidebar {
    width: 60px;
  }

  .admin-logo h2,
  .nav-item span:not(.icon),
  .back-link {
    display: none;
  }

  .admin-content {
    margin-left: 60px;
  }

  .nav-item {
    justify-content: center;
    padding: var(--spacing-md);
  }
}
</style>
