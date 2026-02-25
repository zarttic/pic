import { createApp } from 'vue'

/**
 * 图片懒加载指令
 * 使用 Intersection Observer API 实现高性能懒加载
 *
 * 使用示例：
 * <img v-lazyload="imageUrl" />
 * <img v-lazyload="{ src: imageUrl, placeholder: placeholderUrl }" />
 */

// 默认占位图（1x1 透明像素）
const defaultPlaceholder = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7'

// 懒加载配置
const lazyloadConfig = {
  rootMargin: '50px 0px', // 提前 50px 开始加载
  threshold: 0.01 // 当 1% 可见时触发
}

// Intersection Observer 实例
let observer = null

// 初始化 Observer
function initObserver() {
  if (observer) return observer

  observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const img = entry.target
        const src = img.dataset.lazySrc

        if (src) {
          // 加载图片
          loadImage(img, src)
          // 加载后停止观察
          observer.unobserve(img)
        }
      }
    })
  }, lazyloadConfig)

  return observer
}

// 加载图片
function loadImage(el, src) {
  // 创建新图片对象预加载
  const img = new Image()

  img.onload = () => {
    // 加载成功，替换 src
    el.src = src
    el.classList.remove('lazy-loading')
    el.classList.add('lazy-loaded')
  }

  img.onerror = () => {
    // 加载失败
    el.classList.remove('lazy-loading')
    el.classList.add('lazy-error')
    console.error('Failed to load image:', src)
  }

  img.src = src
}

// 指令定义
export const lazyload = {
  mounted(el, binding) {
    const value = binding.value
    let src = ''
    let placeholder = defaultPlaceholder

    // 解析绑定值
    if (typeof value === 'string') {
      src = value
    } else if (typeof value === 'object') {
      src = value.src || ''
      placeholder = value.placeholder || defaultPlaceholder
    }

    // 设置占位图
    el.src = placeholder
    el.dataset.lazySrc = src

    // 添加加载状态类
    el.classList.add('lazy-loading')

    // 初始化 Observer 并观察元素
    const obs = initObserver()
    obs.observe(el)
  },

  updated(el, binding) {
    const value = binding.value
    let src = ''

    if (typeof value === 'string') {
      src = value
    } else if (typeof value === 'object') {
      src = value.src || ''
    }

    // 如果 src 改变，重新观察
    if (src && el.dataset.lazySrc !== src) {
      el.dataset.lazySrc = src
      el.classList.remove('lazy-loaded', 'lazy-error')
      el.classList.add('lazy-loading')

      const obs = initObserver()
      obs.observe(el)
    }
  },

  unmounted(el) {
    // 元素卸载时停止观察
    if (observer) {
      observer.unobserve(el)
    }
  }
}

// 注册全局指令的函数
export function setupLazyload(app) {
  app.directive('lazyload', lazyload)
}

export default lazyload
