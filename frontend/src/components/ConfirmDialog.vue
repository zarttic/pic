<template>
  <v-dialog
    v-model="visible"
    max-width="480"
    persistent
  >
    <v-card>
      <!-- 卡片头部 -->
      <v-card-title class="d-flex align-center pa-6">
        <v-avatar
          :color="getIconColor"
          size="48"
          class="mr-4"
        >
          <v-icon :icon="getIcon" color="white"></v-icon>
        </v-avatar>
        <span class="text-h6">{{ title }}</span>
      </v-card-title>

      <!-- 卡片内容 -->
      <v-card-text class="px-6 pb-6">
        <p class="text-body-1 text-medium-emphasis mb-0">{{ message }}</p>
      </v-card-text>

      <!-- 卡片操作 -->
      <v-card-actions class="pa-4 bg-grey-lighten-4">
        <v-spacer></v-spacer>
        <v-btn
          variant="outlined"
          @click="cancel"
          class="mr-2"
        >
          {{ cancelText }}
        </v-btn>
        <v-btn
          :color="getButtonColor"
          variant="flat"
          @click="confirm"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { mdiAlert, mdiAlertCircle, mdiInformation } from '@mdi/js'

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

const getIcon = computed(() => {
  const icons = {
    danger: mdiAlertCircle,
    warning: mdiAlert,
    info: mdiInformation
  }
  return icons[props.type] || mdiAlert
})

const getIconColor = computed(() => {
  const colors = {
    danger: 'error',
    warning: 'warning',
    info: 'info'
  }
  return colors[props.type] || 'warning'
})

const getButtonColor = computed(() => {
  const colors = {
    danger: 'error',
    warning: 'warning',
    info: 'primary'
  }
  return colors[props.type] || 'warning'
})

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
/* Vuetify 已经处理了响应式和动画 */
</style>
