import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { setupLazyload } from './directives/lazyload'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 注册全局指令
setupLazyload(app)

app.mount('#app')

// 自定义光标跟随
const cursor = document.querySelector('.cursor')
const cursorGlow = document.querySelector('.cursor-glow')
let mouseX = 0,
  mouseY = 0
let cursorX = 0,
  cursorY = 0

document.addEventListener('mousemove', (e) => {
  mouseX = e.clientX
  mouseY = e.clientY
  if (cursorGlow) {
    cursorGlow.style.left = mouseX + 'px'
    cursorGlow.style.top = mouseY + 'px'
  }
})

function animateCursor() {
  cursorX += (mouseX - cursorX) * 0.1
  cursorY += (mouseY - cursorY) * 0.1
  if (cursor) {
    cursor.style.left = cursorX + 'px'
    cursor.style.top = cursorY + 'px'
  }
  requestAnimationFrame(animateCursor)
}
animateCursor()

// Hover效果
document.addEventListener('mouseover', (e) => {
  if (e.target.matches('a, button, .gallery-item')) {
    if (cursor) {
      cursor.style.transform = 'translate(-50%, -50%) scale(1.5)'
      cursor.style.borderColor = '#d4af37'
    }
  }
})

document.addEventListener('mouseout', (e) => {
  if (e.target.matches('a, button, .gallery-item')) {
    if (cursor) {
      cursor.style.transform = 'translate(-50%, -50%) scale(1)'
      cursor.style.borderColor = '#c9a962'
    }
  }
})
