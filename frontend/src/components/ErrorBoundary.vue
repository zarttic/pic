<template>
  <slot v-if="!hasError" />
  <div v-else class="error-boundary">
    <div class="error-content">
      <div class="error-icon">⚠️</div>
      <h2 class="error-title">出错了</h2>
      <p class="error-message">{{ errorMessage }}</p>
      <div class="error-actions">
        <button class="error-button" @click="resetError">
          重试
        </button>
        <button class="error-button error-button-secondary" @click="goHome">
          返回首页
        </button>
      </div>
      <details class="error-details" v-if="showDetails && errorDetails">
        <summary>技术详情</summary>
        <pre>{{ errorDetails }}</pre>
      </details>
    </div>
  </div>
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationStore } from '@/stores/notification'

const props = defineProps({
  showDetails: {
    type: Boolean,
    default: process.env.NODE_ENV === 'development'
  },
  fallbackMessage: {
    type: String,
    default: '页面加载失败，请重试'
  }
})

const router = useRouter()
const notificationStore = useNotificationStore()

const hasError = ref(false)
const errorMessage = ref('')
const errorDetails = ref('')

// 捕获子组件错误
onErrorCaptured((error, instance, info) => {
  hasError.value = true
  errorMessage.value = error.message || props.fallbackMessage
  errorDetails.value = `${error.stack}\n\nComponent: ${info}`

  // 显示错误通知
  notificationStore.error(props.fallbackMessage, '页面错误')

  // 阻止错误继续向上传播
  return false
})

// 重置错误状态
const resetError = () => {
  hasError.value = false
  errorMessage.value = ''
  errorDetails.value = ''
}

// 返回首页
const goHome = () => {
  resetError()
  router.push('/')
}

// 暴露方法供外部调用
defineExpose({
  resetError
})
</script>

<style scoped>
.error-boundary {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 40px 20px;
}

.error-content {
  text-align: center;
  max-width: 500px;
}

.error-icon {
  font-size: 64px;
  margin-bottom: 24px;
}

.error-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.error-message {
  font-size: 16px;
  color: #6b7280;
  margin-bottom: 32px;
  line-height: 1.6;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 24px;
}

.error-button {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  background: #667eea;
  color: white;
}

.error-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.error-button-secondary {
  background: white;
  color: #667eea;
  border: 2px solid #667eea;
}

.error-button-secondary:hover {
  background: #f5f5f5;
}

.error-details {
  margin-top: 24px;
  text-align: left;
  background: #f9fafb;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.error-details summary {
  cursor: pointer;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 12px;
}

.error-details pre {
  font-size: 12px;
  color: #374151;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* Responsive */
@media (max-width: 768px) {
  .error-actions {
    flex-direction: column;
  }

  .error-button {
    width: 100%;
  }
}
</style>
