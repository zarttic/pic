<template>
  <div class="about-page">
    <section class="about-section">
      <div class="about-bg"></div>
      <div class="about-content">
        <div class="about-image">
          <img src="https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800" alt="摄影师" />
        </div>
        <div class="about-text">
          <h2 class="about-title">关于我</h2>
          <p class="about-description">
            用光影书写故事，以镜头捕捉永恒。每一帧画面都是时间的切片，
            承载着瞬间的情感与永恒的记忆。我相信摄影不仅是记录，
            更是一种与世界对话的方式。
          </p>
          <p class="about-description">
            从城市街角到自然旷野，从人文故事到抽象艺术，
            我寻找那些转瞬即逝的光芒，将它们凝固在时间的河流中。
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    {
      threshold: 0.1
    }
  )

  setTimeout(() => {
    document.querySelectorAll('.about-image, .about-text').forEach((el) => {
      observer.observe(el)
    })
  }, 100)
})
</script>

<style scoped>
.about-page {
  min-height: 100vh;
  padding-top: 100px;
}

.about-section {
  min-height: 100vh;
  display: flex;
  align-items: center;
  padding: var(--spacing-xl) var(--spacing-lg);
  background: var(--bg-secondary);
  position: relative;
  overflow: hidden;
}

.about-bg {
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse at 30% 50%, rgba(201, 169, 98, 0.03) 0%, transparent 50%);
}

.about-content {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: var(--spacing-xl);
  align-items: center;
  position: relative;
  z-index: 2;
}

.about-image {
  position: relative;
  opacity: 0;
  transform: translateX(-50px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.about-image.visible {
  opacity: 1;
  transform: translateX(0);
}

.about-image img {
  width: 100%;
  filter: grayscale(100%);
  transition: filter 0.5s ease;
}

.about-image:hover img {
  filter: grayscale(0%);
}

.about-image::after {
  content: '';
  position: absolute;
  top: var(--spacing-md);
  left: var(--spacing-md);
  right: calc(-1 * var(--spacing-md));
  bottom: calc(-1 * var(--spacing-md));
  border: 1px solid var(--accent-gold);
  opacity: 0.3;
  z-index: -1;
}

.about-text {
  opacity: 0;
  transform: translateX(50px);
  transition: opacity 0.8s ease, transform 0.8s ease 0.2s;
}

.about-text.visible {
  opacity: 1;
  transform: translateX(0);
}

.about-title {
  font-family: 'Cormorant Garamond', serif;
  font-size: clamp(2.5rem, 5vw, 4rem);
  font-weight: 300;
  letter-spacing: 0.1em;
  margin-bottom: var(--spacing-md);
}

.about-description {
  font-size: 1.1rem;
  line-height: 1.9;
  color: var(--text-secondary);
  font-weight: 200;
  margin-bottom: var(--spacing-md);
}

@media (max-width: 768px) {
  .about-content {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
  }
}
</style>
