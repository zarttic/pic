/**
 * 统一错误处理工具
 * 提供错误分类、格式化和日志记录功能
 */

/**
 * 错误类型枚举
 */
export const ErrorType = {
  NETWORK: 'NETWORK_ERROR',
  API: 'API_ERROR',
  VALIDATION: 'VALIDATION_ERROR',
  AUTH: 'AUTH_ERROR',
  UNKNOWN: 'UNKNOWN_ERROR'
}

/**
 * 错误分类器
 * @param {Error} error - 错误对象
 * @returns {Object} 分类后的错误信息
 */
export function classifyError(error) {
  // Axios 错误
  if (error.response) {
    const { status, data } = error.response

    // 认证错误
    if (status === 401) {
      return {
        type: ErrorType.AUTH,
        message: '登录已过期，请重新登录',
        code: status,
        details: data
      }
    }

    // 权限错误
    if (status === 403) {
      return {
        type: ErrorType.AUTH,
        message: '没有权限执行此操作',
        code: status,
        details: data
      }
    }

    // 资源不存在
    if (status === 404) {
      return {
        type: ErrorType.API,
        message: '请求的资源不存在',
        code: status,
        details: data
      }
    }

    // 验证错误
    if (status === 400 || status === 422) {
      return {
        type: ErrorType.VALIDATION,
        message: data.message || '数据验证失败',
        code: status,
        details: data
      }
    }

    // 服务器错误
    if (status >= 500) {
      return {
        type: ErrorType.API,
        message: '服务器错误，请稍后重试',
        code: status,
        details: data
      }
    }

    // 其他 API 错误
    return {
      type: ErrorType.API,
      message: data.message || '请求失败',
      code: status,
      details: data
    }
  }

  // 网络错误
  if (error.request) {
    return {
      type: ErrorType.NETWORK,
      message: '网络连接失败，请检查网络设置',
      code: 0,
      details: null
    }
  }

  // 请求配置错误
  if (error.message) {
    return {
      type: ErrorType.UNKNOWN,
      message: error.message,
      code: 0,
      details: null
    }
  }

  // 完全未知的错误
  return {
    type: ErrorType.UNKNOWN,
    message: '发生未知错误',
    code: 0,
    details: null
  }
}

/**
 * 错误处理器
 * @param {Error} error - 错误对象
 * @param {string} context - 错误上下文描述
 * @param {Object} options - 处理选项
 * @returns {Object} 处理后的错误信息
 */
export function handleError(error, context = '', options = {}) {
  const { silent = false, showToast = true } = options

  // 分类错误
  const classified = classifyError(error)

  // 构建完整的错误信息
  const fullContext = context ? `[${context}]` : ''
  const errorMessage = `${fullContext} ${classified.message}`.trim()

  // 开发环境输出详细日志
  if (import.meta.env.DEV) {
    console.group(`🔴 Error ${fullContext}`)
    console.error('Message:', errorMessage)
    console.error('Type:', classified.type)
    console.error('Code:', classified.code)
    if (classified.details) {
      console.error('Details:', classified.details)
    }
    console.error('Stack:', error.stack)
    console.groupEnd()
  }

  // 生产环境可以上报到错误追踪服务
  if (import.meta.env.PROD) {
    // TODO: 集成错误追踪服务（如 Sentry）
    // sendToSentry(error, { context, ...classified })
  }

  return {
    ...classified,
    message: errorMessage,
    originalError: error
  }
}

/**
 * 获取用户友好的错误消息
 * @param {Error} error - 错误对象
 * @returns {string} 友好的错误消息
 */
export function getUserMessage(error) {
  const classified = classifyError(error)
  return classified.message
}
