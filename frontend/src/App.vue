<script setup>
import { onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import Navigation from './components/Navigation.vue'
import ToastContainer from './components/ToastContainer.vue'

const route = useRoute()

// 判断是否为后台路由
const isAdminRoute = computed(() => route.path.startsWith('/admin'))

// 自定义光标跟随（仅前台）
let mouseX = 0
let mouseY = 0
let cursorX = 0
let cursorY = 0
let animationId = null

const handleMouseMove = (e) => {
  mouseX = e.clientX
  mouseY = e.clientY
  const cursorGlow = document.querySelector('.cursor-glow')
  if (cursorGlow) {
    cursorGlow.style.left = mouseX + 'px'
    cursorGlow.style.top = mouseY + 'px'
  }
}

const animateCursor = () => {
  cursorX += (mouseX - cursorX) * 0.1
  cursorY += (mouseY - cursorY) * 0.1
  const cursor = document.querySelector('.cursor')
  if (cursor) {
    cursor.style.left = cursorX + 'px'
    cursor.style.top = cursorY + 'px'
  }
  animationId = requestAnimationFrame(animateCursor)
}

const handleMouseOver = (e) => {
  if (e.target.matches('a, button, .gallery-item, .v-btn, .v-card')) {
    const cursor = document.querySelector('.cursor')
    if (cursor) {
      cursor.style.transform = 'translate(-50%, -50%) scale(1.5)'
      cursor.style.borderColor = '#d4af37'
    }
  }
}

const handleMouseOut = (e) => {
  if (e.target.matches('a, button, .gallery-item, .v-btn, .v-card')) {
    const cursor = document.querySelector('.cursor')
    if (cursor) {
      cursor.style.transform = 'translate(-50%, -50%) scale(1)'
      cursor.style.borderColor = '#c9a962'
    }
  }
}

onMounted(() => {
  // 仅在前台路由启用自定义光标
  if (!isAdminRoute.value) {
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseover', handleMouseOver)
    document.addEventListener('mouseout', handleMouseOut)
    animateCursor()
  }
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseover', handleMouseOver)
  document.removeEventListener('mouseout', handleMouseOut)
})
</script>

<template>
  <v-app>
    <!-- 前台自定义光标 -->
    <div v-if="!isAdminRoute" class="cursor"></div>
    <div v-if="!isAdminRoute" class="cursor-glow"></div>

    <!-- 前台导航栏（仅在非后台路由显示） -->
    <Navigation v-if="!isAdminRoute" />

    <!-- 主内容区域 -->
    <v-main>
      <router-view />
    </v-main>

    <!-- Toast 通知 -->
    <ToastContainer />
  </v-app>
</template>

<style>
@import './styles/main.css';
@import url('https://fonts.googleapis.com/css2?family=Cormorant+Garamond:wght@300;400;600&family=Outfit:wght@200;300;400&display=swap');

#app {
  min-height: 100vh;
}
</style>
