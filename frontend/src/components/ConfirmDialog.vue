<template>
  <Transition name="dialog-fade">
    <div v-if="visible" class="confirm-dialog-overlay" @click.self="cancel">
      <div class="confirm-dialog" :class="`dialog-${type}`">
        <div class="dialog-header">
          <div class="dialog-icon">
            <span v-if="type === 'danger'">⚠</span>
            <span v-else-if="type === 'warning'">⚡</span>
            <span v-else-if="type === 'info'">ℹ</span>
          </div>
          <h3 class="dialog-title">{{ title }}</h3>
        </div>

        <div class="dialog-body">
          <p class="dialog-message">{{ message }}</p>
        </div>

        <div class="dialog-footer">
          <button class="btn btn-cancel" @click="cancel">
            {{ cancelText }}
          </button>
          <button class="btn btn-confirm" :class="`btn-${type}`" @click="confirm">
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'warning',
    validator: (value) => ['danger', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: '确认操作'
  },
  message: {
    type: String,
    required: true
  },
  confirmText: {
    type: String,
    default: '确认'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  onConfirm: {
    type: Function,
    default: null
  },
  onCancel: {
    type: Function,
    default: null
  }
})

const visible = ref(false)

const show = () => {
  visible.value = true
}

const hide = () => {
  visible.value = false
}

const confirm = () => {
  hide()
  if (props.onConfirm) {
    props.onConfirm(true)
  }
}

const cancel = () => {
  hide()
  if (props.onCancel) {
    props.onCancel(false)
  }
}

const handleKeydown = (event) => {
  if (!visible.value) return

  if (event.key === 'Escape') {
    cancel()
  } else if (event.key === 'Enter') {
    confirm()
  }
}

onMounted(() => {
  show()
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

defineExpose({
  show,
  hide,
  confirm,
  cancel
})
</script>

<style scoped>
.confirm-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.confirm-dialog {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 480px;
  width: 90%;
  overflow: hidden;
}

/* Dialog Types */
.dialog-danger {
  border-top: 4px solid #ef4444;
}

.dialog-warning {
  border-top: 4px solid #f59e0b;
}

.dialog-info {
  border-top: 4px solid #3b82f6;
}

/* Header */
.dialog-header {
  padding: 24px 24px 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.dialog-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.dialog-danger .dialog-icon {
  background: #fee2e2;
  color: #ef4444;
}

.dialog-warning .dialog-icon {
  background: #fef3c7;
  color: #f59e0b;
}

.dialog-info .dialog-icon {
  background: #dbeafe;
  color: #3b82f6;
}

.dialog-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

/* Body */
.dialog-body {
  padding: 0 24px 24px;
}

.dialog-message {
  margin: 0;
  color: #6b7280;
  line-height: 1.6;
  font-size: 15px;
}

/* Footer */
.dialog-footer {
  padding: 16px 24px;
  background: #f9fafb;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Buttons */
.btn {
  padding: 10px 24px;
  border-radius: 8px;
  border: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel {
  background: white;
  color: #6b7280;
  border: 1px solid #d1d5db;
}

.btn-cancel:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

.btn-confirm {
  color: white;
}

.btn-danger {
  background: #ef4444;
}

.btn-danger:hover {
  background: #dc2626;
}

.btn-warning {
  background: #f59e0b;
}

.btn-warning:hover {
  background: #d97706;
}

.btn-info {
  background: #3b82f6;
}

.btn-info:hover {
  background: #2563eb;
}

/* Animations */
.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: all 0.3s ease;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}

.dialog-fade-enter-from .confirm-dialog,
.dialog-fade-leave-to .confirm-dialog {
  transform: scale(0.9) translateY(-20px);
}

/* Responsive */
@media (max-width: 768px) {
  .confirm-dialog {
    width: 95%;
    margin: 20px;
  }

  .dialog-header {
    padding: 20px 20px 12px;
  }

  .dialog-body {
    padding: 0 20px 20px;
  }

  .dialog-footer {
    padding: 12px 20px;
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
    padding: 12px 24px;
  }
}
</style>
