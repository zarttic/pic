<template>
  <div class="admin-layout">
    <!-- Mobile Header -->
    <header class="mobile-header">
      <button class="menu-toggle" @click="toggleSidebar">
        <span :class="['hamburger', { active: isSidebarOpen }]">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </button>
      <h2 class="mobile-title">PicSite 管理</h2>
      <div class="mobile-spacer"></div>
    </header>

    <!-- Overlay for mobile -->
    <div v-if="isSidebarOpen" class="sidebar-overlay" @click="closeSidebar"></div>

    <!-- Sidebar -->
    <aside :class="['admin-sidebar', { open: isSidebarOpen }]">
      <div class="admin-logo">
        <h2>PicSite 管理</h2>
      </div>
      <nav class="admin-nav">
        <router-link to="/admin/photos" class="nav-item" @click="closeSidebarOnMobile">
          <span class="icon">📸</span>
          <span class="nav-text">照片管理</span>
        </router-link>
        <router-link to="/admin/albums" class="nav-item" @click="closeSidebarOnMobile">
          <span class="icon">📁</span>
          <span class="nav-text">相册管理</span>
        </router-link>
        <router-link to="/admin/statistics" class="nav-item" @click="closeSidebarOnMobile">
          <span class="icon">📊</span>
          <span class="nav-text">统计数据</span>
        </router-link>
        <router-link to="/admin/settings" class="nav-item" @click="closeSidebarOnMobile">
          <span class="icon">⚙️</span>
          <span class="nav-text">系统设置</span>
        </router-link>
      </nav>
      <div class="admin-footer">
        <router-link to="/" class="back-link" @click="closeSidebarOnMobile">
          ← 返回前台
        </router-link>
        <button v-if="isAuthenticated" class="logout-button" @click="handleLogout">
          退出登录
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="admin-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isSidebarOpen = ref(false)
const isAuthenticated = computed(() => authStore.isAuthenticated)

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}

const closeSidebar = () => {
  isSidebarOpen.value = false
}

const closeSidebarOnMobile = () => {
  if (window.innerWidth < 1024) {
    closeSidebar()
  }
}

const handleLogout = async () => {
  await authStore.logout()
  closeSidebarOnMobile()
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--bg-primary);
}

/* Mobile Header */
.mobile-header {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: var(--bg-secondary);
  border-bottom: 1px solid rgba(201, 169, 98, 0.1);
  z-index: 1000;
  padding: 0 16px;
  align-items: center;
  gap: 16px;
}

.menu-toggle {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hamburger {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 24px;
}

.hamburger span {
  display: block;
  height: 2px;
  background: var(--accent-gold);
  transition: all 0.3s ease;
  transform-origin: center;
}

.hamburger.active span:nth-child(1) {
  transform: rotate(45deg) translate(4px, 4px);
}

.hamburger.active span:nth-child(2) {
  opacity: 0;
}

.hamburger.active span:nth-child(3) {
  transform: rotate(-45deg) translate(4px, -4px);
}

.mobile-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.4rem;
  font-weight: 300;
  color: var(--accent-gold);
  letter-spacing: 0.1em;
  flex: 1;
}

.mobile-spacer {
  width: 40px;
}

/* Overlay */
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Sidebar */
.admin-sidebar {
  width: 260px;
  background: var(--bg-secondary);
  border-right: 1px solid rgba(201, 169, 98, 0.1);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  overflow-y: auto;
  transition: transform 0.3s ease;
  z-index: 1001;
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
  display: flex;
  flex-direction: column;
  gap: 12px;
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

.logout-button {
  width: 100%;
  padding: 10px 16px;
  border: 1px solid var(--accent-gold);
  background: transparent;
  color: var(--accent-gold);
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.logout-button:hover {
  background: var(--accent-gold);
  color: var(--bg-primary);
}

.admin-content {
  flex: 1;
  margin-left: 260px;
  padding: var(--spacing-xl);
  overflow-y: auto;
}

/* Desktop */
@media (min-width: 1025px) {
  .mobile-header,
  .sidebar-overlay {
    display: none;
  }

  .admin-sidebar {
    transform: translateX(0);
  }
}

/* Tablet */
@media (max-width: 1024px) {
  .mobile-header {
    display: flex;
  }

  .sidebar-overlay {
    display: block;
  }

  .admin-sidebar {
    width: 260px;
    transform: translateX(-100%);
  }

  .admin-sidebar.open {
    transform: translateX(0);
  }

  .admin-content {
    margin-left: 0;
    margin-top: 60px;
    padding: var(--spacing-lg);
  }
}

/* Mobile */
@media (max-width: 768px) {
  .admin-content {
    padding: var(--spacing-md);
  }

  .admin-logo {
    padding: var(--spacing-lg);
  }

  .nav-item {
    padding: var(--spacing-md) var(--spacing-lg);
  }
}
</style>
