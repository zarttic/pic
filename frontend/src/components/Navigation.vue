<template>
  <v-app-bar
    flat
    :height="80"
    class="navigation"
  >
    <v-toolbar-title class="logo">光迹</v-toolbar-title>

    <v-spacer></v-spacer>

    <v-btn
      v-for="link in navLinks"
      :key="link.path"
      :to="link.path"
      variant="text"
      class="nav-link"
      :ripple="false"
    >
      {{ link.text }}
    </v-btn>
  </v-app-bar>
</template>

<script setup>
const navLinks = [
  { path: '/', text: '首页' },
  { path: '/gallery', text: '作品' },
  { path: '/albums', text: '相册' },
  { path: '/about', text: '关于' },
  { path: '/admin/photos', text: '管理' },
]
</script>

<style scoped>
.navigation {
  position: fixed !important;
  background: linear-gradient(to bottom, rgba(10, 10, 10, 0.95) 0%, transparent 100%) !important;
  backdrop-filter: blur(10px);
}

.logo {
  font-family: 'Cormorant Garamond', serif !important;
  font-size: 1.8rem !important;
  font-weight: 300 !important;
  letter-spacing: 0.3em;
  color: var(--accent-gold) !important;
}

.nav-link {
  color: var(--text-secondary) !important;
  font-size: 0.9rem !important;
  font-weight: 200 !important;
  letter-spacing: 0.15em !important;
  text-transform: uppercase !important;
  position: relative;
  transition: color 0.3s ease;
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 5px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 1px;
  background: var(--accent-gold);
  transition: width 0.3s ease;
}

.nav-link:hover,
.nav-link.v-btn--active {
  color: var(--text-primary) !important;
}

.nav-link:hover::after,
.nav-link.v-btn--active::after {
  width: 80%;
}

@media (max-width: 768px) {
  .navigation {
    padding: 0 var(--spacing-md) !important;
  }

  .nav-link {
    font-size: 0.75rem !important;
  }
}
</style>
