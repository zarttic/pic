<template>
  <Transition name="toast-fade">
    <div v-if="visible" :class="['toast', `toast-${type}`]" @click="close">
      <div class="toast-icon">
        <span v-if="type === 'success'">✓</span>
        <span v-else-if="type === 'error'">✕</span>
        <span v-else-if="type === 'warning'">⚠</span>
        <span v-else-if="type === 'info'">ℹ</span>
      </div>
      <div class="toast-content">
        <div v-if="title" class="toast-title">{{ title }}</div>
        <div class="toast-message">{{ message }}</div>
      </div>
      <button class="toast-close" @click.stop="close">×</button>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  },
  duration: {
    type: Number,
    default: 3000
  },
  onClose: {
    type: Function,
    default: null
  }
})

const visible = ref(false)

const show = () => {
  visible.value = true
  if (props.duration > 0) {
    setTimeout(() => {
      close()
    }, props.duration)
  }
}

const close = () => {
  visible.value = false
  if (props.onClose) {
    props.onClose()
  }
}

onMounted(() => {
  show()
})

defineExpose({
  show,
  close
})
</script>

<style scoped>
.toast {
  display: flex;
  align-items: center;
  min-width: 300px;
  max-width: 500px;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.toast:hover {
  transform: translateX(-4px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

/* Types */
.toast-success {
  background: #f0fdf4;
  border-left: 4px solid #10b981;
}

.toast-error {
  background: #fef2f2;
  border-left: 4px solid #ef4444;
}

.toast-warning {
  background: #fffbeb;
  border-left: 4px solid #f59e0b;
}

.toast-info {
  background: #eff6ff;
  border-left: 4px solid #3b82f6;
}

/* Icon */
.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  margin-right: 12px;
  font-size: 18px;
  font-weight: bold;
  flex-shrink: 0;
}

.toast-success .toast-icon {
  background: #10b981;
  color: white;
}

.toast-error .toast-icon {
  background: #ef4444;
  color: white;
}

.toast-warning .toast-icon {
  background: #f59e0b;
  color: white;
}

.toast-info .toast-icon {
  background: #3b82f6;
  color: white;
}

/* Content */
.toast-content {
  flex: 1;
}

.toast-title {
  font-weight: 600;
  margin-bottom: 4px;
  color: #1f2937;
}

.toast-message {
  color: #4b5563;
  font-size: 14px;
}

/* Close Button */
.toast-close {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  font-size: 20px;
  color: #9ca3af;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  margin-left: 12px;
  flex-shrink: 0;
}

.toast-close:hover {
  color: #6b7280;
}

/* Animations */
.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.3s ease;
}

.toast-fade-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-fade-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

/* Responsive */
@media (max-width: 768px) {
  .toast {
    min-width: 280px;
    max-width: 90vw;
  }
}
</style>
