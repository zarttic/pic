import { handleError, getUserMessage } from '../utils/errorHandler'
import { useNotificationStore } from '../stores/notification'

/**
 * 错误处理组合式函数
 * 提供统一的错误处理和用户通知功能
 *
 * @example
 * const { handleAsyncError, showError } = useErrorHandler()
 *
 * // 方式1: 包装异步函数
 * await handleAsyncError(
 *   () => photoStore.fetchPhotos(),
 *   '加载照片失败'
 * )
 *
 * // 方式2: 直接显示错误
 * try {
 *   await someAsyncOperation()
 * } catch (error) {
 *   showError(error, '操作失败')
 * }
 */
export function useErrorHandler() {
  const notification = useNotificationStore()

  /**
   * 处理异步操作错误
   * @param {Function} asyncFn - 异步函数
   * @param {string} context - 错误上下文描述
   * @param {Object} options - 处理选项
   * @returns {Promise<any>} 异步函数结果，错误时返回 null
   */
  async function handleAsyncError(asyncFn, context = '', options = {}) {
    const { showToast = true, rethrow = false } = options

    try {
      return await asyncFn()
    } catch (error) {
      const processedError = handleError(error, context, { showToast })

      if (showToast) {
        notification.error(processedError.message)
      }

      if (rethrow) {
        throw error
      }

      return null
    }
  }

  /**
   * 显示错误通知
   * @param {Error} error - 错误对象
   * @param {string} context - 错误上下文描述
   * @param {Object} options - 处理选项
   */
  function showError(error, context = '', options = {}) {
    const { showToast = true } = options
    const processedError = handleError(error, context, { showToast })

    if (showToast) {
      notification.error(processedError.message)
    }
  }

  /**
   * 处理同步操作错误
   * @param {Function} fn - 同步函数
   * @param {string} context - 错误上下文描述
   * @param {Object} options - 处理选项
   * @returns {any} 函数结果，错误时返回 null
   */
  function handleSyncError(fn, context = '', options = {}) {
    const { showToast = true } = options

    try {
      return fn()
    } catch (error) {
      const processedError = handleError(error, context, { showToast })

      if (showToast) {
        notification.error(processedError.message)
      }

      return null
    }
  }

  /**
   * 创建错误处理包装器
   * @param {string} defaultContext - 默认错误上下文
   * @returns {Function} 错误处理包装函数
   */
  function createErrorHandler(defaultContext = '') {
    return (error, context = '', options = {}) => {
      const fullContext = context || defaultContext
      return showError(error, fullContext, options)
    }
  }

  return {
    handleAsyncError,
    handleSyncError,
    showError,
    getUserMessage,
    createErrorHandler,
    handleError
  }
}
