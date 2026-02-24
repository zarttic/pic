<template>
  <div class="home">
    <section class="hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <h1 class="hero-title">光迹</h1>
        <p class="hero-subtitle">捕捉瞬间永恒之美</p>
      </div>
      <div class="scroll-indicator">
        <span class="scroll-text">探索</span>
        <div class="scroll-line"></div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  // 添加淡入动画
  const elements = document.querySelectorAll('.hero-title, .hero-subtitle, .scroll-indicator')
  elements.forEach((el, index) => {
    el.style.opacity = '0'
    setTimeout(() => {
      el.style.transition = 'opacity 1s ease'
      el.style.opacity = '1'
    }, index * 300)
  })
})
</script>

<style scoped>
.home {
  min-height: 100vh;
}

.hero {
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse at 20% 50%, rgba(201, 169, 98, 0.08) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 50%, rgba(201, 169, 98, 0.05) 0%, transparent 50%);
  filter: blur(100px);
}

.hero-content {
  text-align: center;
  z-index: 2;
}

.hero-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(4rem, 12vw, 10rem);
  font-weight: 300;
  line-height: 1;
  letter-spacing: 0.2em;
  margin-bottom: var(--spacing-md);
}

.hero-subtitle {
  font-size: 1rem;
  font-weight: 200;
  letter-spacing: 0.4em;
  text-transform: uppercase;
  color: var(--text-secondary);
}

.scroll-indicator {
  position: absolute;
  bottom: var(--spacing-lg);
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
}

.scroll-line {
  width: 1px;
  height: 60px;
  background: linear-gradient(to bottom, var(--accent-gold), transparent);
  animation: scrollPulse 2s ease-in-out infinite;
}

.scroll-text {
  font-size: 0.75rem;
  letter-spacing: 0.2em;
  color: var(--text-secondary);
  text-transform: uppercase;
}

@keyframes scrollPulse {
  0%, 100% {
    opacity: 1;
    transform: translateY(0);
  }
  50% {
    opacity: 0.5;
    transform: translateY(10px);
  }
}
</style>
