import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref([])
  let notificationId = 0

  // 获取所有通知
  const allNotifications = computed(() => notifications.value)

  // 添加通知
  const add = (notification) => {
    const id = notificationId++
    const newNotification = {
      id,
      type: notification.type || 'info',
      title: notification.title || '',
      message: notification.message || '',
      duration: notification.duration !== undefined ? notification.duration : 3000
    }

    notifications.value.push(newNotification)

    // 自动移除
    if (newNotification.duration > 0) {
      setTimeout(() => {
        remove(id)
      }, newNotification.duration)
    }

    return id
  }

  // 移除通知
  const remove = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  // 清空所有通知
  const clear = () => {
    notifications.value = []
  }

  // 快捷方法
  const success = (message, title = '') => {
    return add({ type: 'success', message, title })
  }

  const error = (message, title = '错误') => {
    return add({ type: 'error', message, title, duration: 5000 })
  }

  const warning = (message, title = '警告') => {
    return add({ type: 'warning', message, title, duration: 4000 })
  }

  const info = (message, title = '') => {
    return add({ type: 'info', message, title })
  }

  return {
    notifications: allNotifications,
    add,
    remove,
    clear,
    success,
    error,
    warning,
    info
  }
})
