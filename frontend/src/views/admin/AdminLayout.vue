<template>
  <v-app>
    <!-- 移动端顶部栏 -->
    <v-app-bar
      v-if="mobile"
      density="comfortable"
      elevation="1"
    >
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>PicSite 管理</v-toolbar-title>
    </v-app-bar>

    <!-- 侧边导航栏 -->
    <v-navigation-drawer
      v-model="drawer"
      :permanent="!mobile"
      :temporary="mobile"
      width="260"
    >
      <!-- Logo -->
      <v-list-item class="pa-6">
        <v-list-item-title class="text-h5 font-weight-light">
          PicSite 管理
        </v-list-item-title>
      </v-list-item>

      <v-divider></v-divider>

      <!-- 导航菜单 -->
      <v-list nav density="comfortable">
        <v-list-item
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          :prepend-icon="item.icon"
          :title="item.text"
          @click="handleNavClick"
        ></v-list-item>
      </v-list>

      <v-spacer></v-spacer>

      <v-divider></v-divider>

      <!-- 底部操作 -->
      <v-list nav density="comfortable">
        <v-list-item
          to="/"
          prepend-icon="mdi-arrow-left"
          title="返回前台"
        ></v-list-item>

        <v-list-item
          v-if="isAuthenticated"
          prepend-icon="mdi-logout"
          title="退出登录"
          @click="handleLogout"
        ></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 主内容区域 -->
    <v-main>
      <v-container fluid class="pa-8">
        <router-view />
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTheme } from 'vuetify'
import { useAuthStore } from '@/stores/auth'
import { useDisplay } from 'vuetify'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const theme = useTheme()
const { mobile } = useDisplay()

const drawer = ref(!mobile.value)
const isAuthenticated = computed(() => authStore.isAuthenticated)

const navItems = [
  { path: '/admin/photos', text: '照片管理', icon: 'mdi-image-multiple' },
  { path: '/admin/albums', text: '相册管理', icon: 'mdi-folder-multiple' },
  { path: '/admin/statistics', text: '统计数据', icon: 'mdi-chart-bar' },
  { path: '/admin/settings', text: '系统设置', icon: 'mdi-cog' },
]

const handleNavClick = () => {
  if (mobile.value) {
    drawer.value = false
  }
}

const handleLogout = async () => {
  await authStore.logout()
  if (mobile.value) {
    drawer.value = false
  }
  router.push('/admin/login')
}

// 切换到后台浅色主题
onMounted(() => {
  theme.global.name.value = 'adminLightTheme'
})
</script>

<style scoped>
/* Vuetify 已经处理了响应式布局，这里只需要少量自定义样式 */
</style>
