import { createApp, h } from 'vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'

/**
 * 确认对话框组合式函数
 * 提供异步确认对话框功能
 *
 * @example
 * const confirm = useConfirm()
 * const result = await confirm({
 *   type: 'danger',
 *   title: '删除照片',
 *   message: '确定要删除这张照片吗？此操作不可撤销。',
 *   confirmText: '删除',
 *   cancelText: '取消'
 * })
 * if (result) {
 *   // 用户确认
 * }
 */
export function useConfirm() {
  /**
   * 显示确认对话框
   * @param {Object} options - 对话框选项
   * @param {string} options.type - 对话框类型: 'danger' | 'warning' | 'info'
   * @param {string} options.title - 对话框标题
   * @param {string} options.message - 对话框消息
   * @param {string} options.confirmText - 确认按钮文本
   * @param {string} options.cancelText - 取消按钮文本
   * @returns {Promise<boolean>} 用户是否确认
   */
  return (options) => {
    return new Promise((resolve) => {
      // 创建容器元素
      const container = document.createElement('div')
      document.body.appendChild(container)

      // 创建应用实例
      const app = createApp({
        render() {
          return h(ConfirmDialog, {
            type: options.type || 'warning',
            title: options.title || '确认操作',
            message: options.message,
            confirmText: options.confirmText || '确认',
            cancelText: options.cancelText || '取消',
            onConfirm: (result) => {
              resolve(result)
              cleanup()
            },
            onCancel: (result) => {
              resolve(result)
              cleanup()
            }
          })
        }
      })

      // 挂载应用
      const instance = app.mount(container)

      // 清理函数
      const cleanup = () => {
        // 等待动画完成后再清理
        setTimeout(() => {
          app.unmount()
          if (container.parentNode) {
            container.parentNode.removeChild(container)
          }
        }, 300)
      }
    })
  }
}
