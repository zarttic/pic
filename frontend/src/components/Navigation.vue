<template>
  <nav class="navigation">
    <div class="logo">光迹</div>
    <ul class="nav-links">
      <li>
        <router-link to="/" @click="closeMobileMenu">首页</router-link>
      </li>
      <li>
        <router-link to="/gallery" @click="closeMobileMenu">作品</router-link>
      </li>
      <li>
        <router-link to="/albums" @click="closeMobileMenu">相册</router-link>
      </li>
      <li>
        <router-link to="/about" @click="closeMobileMenu">关于</router-link>
      </li>
      <li>
        <router-link to="/admin/photos" @click="closeMobileMenu">管理</router-link>
      </li>
    </ul>
  </nav>
</template>

<script setup>
import { ref } from 'vue'

const mobileMenuOpen = ref(false)

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}
</script>

<style scoped>
.navigation {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: var(--spacing-md) var(--spacing-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(to bottom, rgba(10, 10, 10, 0.95) 0%, transparent 100%);
  backdrop-filter: blur(10px);
}

.logo {
  font-family: 'Cormorant Garamond', serif;
  font-size: 1.8rem;
  font-weight: 300;
  letter-spacing: 0.3em;
  color: var(--accent-gold);
}

.nav-links {
  display: flex;
  gap: var(--spacing-md);
  list-style: none;
}

.nav-links a {
  color: var(--text-secondary);
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 200;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  transition: color 0.3s ease;
  position: relative;
}

.nav-links a::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 0;
  height: 1px;
  background: var(--accent-gold);
  transition: width 0.3s ease;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  color: var(--text-primary);
}

.nav-links a:hover::after,
.nav-links a.router-link-active::after {
  width: 100%;
}

@media (max-width: 768px) {
  .navigation {
    padding: var(--spacing-sm) var(--spacing-md);
  }

  .nav-links {
    gap: var(--spacing-sm);
  }

  .nav-links a {
    font-size: 0.75rem;
  }
}
</style>
