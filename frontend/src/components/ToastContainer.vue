<template>
  <Teleport to="body">
    <div class="toast-container">
      <v-snackbar
        v-for="notification in notifications"
        :key="notification.id"
        v-model="notification.visible"
        :color="getSnackbarColor(notification.type)"
        :timeout="notification.duration"
        location="top right"
        :multi-line="!!notification.title"
        @update:model-value="handleClose(notification.id)"
      >
        <div class="d-flex align-center">
          <v-icon start :icon="getIcon(notification.type)"></v-icon>
          <div>
            <div v-if="notification.title" class="font-weight-bold mb-1">
              {{ notification.title }}
            </div>
            <div>{{ notification.message }}</div>
          </div>
        </div>

        <template v-slot:actions>
          <v-btn
            variant="text"
            :icon="mdiClose"
            @click="removeNotification(notification.id)"
          ></v-btn>
        </template>
      </v-snackbar>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import { mdiCheckCircle, mdiAlertCircle, mdiAlert, mdiInformation, mdiClose } from '@mdi/js'

const notificationStore = useNotificationStore()
const notifications = computed(() => notificationStore.notifications)

const getSnackbarColor = (type) => {
  const colors = {
    success: 'success',
    error: 'error',
    warning: 'warning',
    info: 'info'
  }
  return colors[type] || 'info'
}

const getIcon = (type) => {
  const icons = {
    success: mdiCheckCircle,
    error: mdiAlertCircle,
    warning: mdiAlert,
    info: mdiInformation
  }
  return icons[type] || mdiInformation
}

const removeNotification = (id) => {
  notificationStore.remove(id)
}

const handleClose = (id) => {
  // 当 snackbar 自动关闭时，从 store 中移除
  removeNotification(id)
}

// 设置通知为可见
onMounted(() => {
  notifications.value.forEach(n => {
    n.visible = true
  })
})
</script>

<style scoped>
.toast-container {
  position: fixed;
  z-index: 9999;
}
</style>
