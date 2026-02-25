/**
 * 表单验证工具
 * 提供常用的表单验证规则和验证器
 */

/**
 * 验证规则类型
 */
export const ValidatorType = {
  REQUIRED: 'required',
  MIN_LENGTH: 'minLength',
  MAX_LENGTH: 'maxLength',
  MIN: 'min',
  MAX: 'max',
  PATTERN: 'pattern',
  EMAIL: 'email',
  URL: 'url',
  YEAR: 'year',
  FILE_SIZE: 'fileSize',
  FILE_TYPE: 'fileType'
}

/**
 * 验证器函数映射
 */
const validators = {
  // 必填验证
  [ValidatorType.REQUIRED]: (value, message = '此字段为必填项') => {
    if (value === null || value === undefined || value === '') {
      return message
    }
    if (Array.isArray(value) && value.length === 0) {
      return message
    }
    return null
  },

  // 最小长度验证
  [ValidatorType.MIN_LENGTH]: (value, min, message) => {
    if (!value || value.length < min) {
      return message || `长度不能少于 ${min} 个字符`
    }
    return null
  },

  // 最大长度验证
  [ValidatorType.MAX_LENGTH]: (value, max, message) => {
    if (value && value.length > max) {
      return message || `长度不能超过 ${max} 个字符`
    }
    return null
  },

  // 最小值验证
  [ValidatorType.MIN]: (value, min, message) => {
    if (value !== null && value !== undefined && value < min) {
      return message || `值不能小于 ${min}`
    }
    return null
  },

  // 最大值验证
  [ValidatorType.MAX]: (value, max, message) => {
    if (value !== null && value !== undefined && value > max) {
      return message || `值不能大于 ${max}`
    }
    return null
  },

  // 正则表达式验证
  [ValidatorType.PATTERN]: (value, pattern, message = '格式不正确') => {
    if (value && !pattern.test(value)) {
      return message
    }
    return null
  },

  // 邮箱验证
  [ValidatorType.EMAIL]: (value, message = '请输入有效的邮箱地址') => {
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (value && !emailPattern.test(value)) {
      return message
    }
    return null
  },

  // URL 验证
  [ValidatorType.URL]: (value, message = '请输入有效的 URL') => {
    try {
      if (value) {
        new URL(value)
      }
    } catch {
      return message
    }
    return null
  },

  // 年份验证
  [ValidatorType.YEAR]: (value, message) => {
    const currentYear = new Date().getFullYear()
    if (value && (value < 1900 || value > currentYear)) {
      return message || `年份应在 1900-${currentYear} 之间`
    }
    return null
  },

  // 文件大小验证
  [ValidatorType.FILE_SIZE]: (file, maxSizeMB, message) => {
    if (file && file.size > maxSizeMB * 1024 * 1024) {
      return message || `文件大小不能超过 ${maxSizeMB}MB`
    }
    return null
  },

  // 文件类型验证
  [ValidatorType.FILE_TYPE]: (file, allowedTypes, message) => {
    if (file && !allowedTypes.includes(file.type)) {
      return message || `文件类型不支持，仅支持：${allowedTypes.join(', ')}`
    }
    return null
  }
}

/**
 * 创建验证规则
 * @param {string} type - 验证类型
 * @param {any} options - 验证选项
 * @param {string} message - 自定义错误消息
 * @returns {Function} 验证函数
 */
export function createRule(type, options, message) {
  return (value) => {
    const validator = validators[type]
    if (!validator) {
      console.warn(`Unknown validator type: ${type}`)
      return null
    }

    // 根据验证器类型传递参数
    switch (type) {
      case ValidatorType.MIN_LENGTH:
      case ValidatorType.MIN:
      case ValidatorType.MAX:
      case ValidatorType.FILE_SIZE:
        return validator(value, options, message)

      case ValidatorType.MAX_LENGTH:
        return validator(value, options, message)

      case ValidatorType.PATTERN:
        return validator(value, options, message)

      case ValidatorType.FILE_TYPE:
        return validator(value, options, message)

      default:
        return validator(value, message)
    }
  }
}

/**
 * 验证单个字段
 * @param {any} value - 字段值
 * @param {Array<Function>} rules - 验证规则数组
 * @returns {string|null} 第一个错误消息或 null
 */
export function validateField(value, rules) {
  for (const rule of rules) {
    const error = rule(value)
    if (error) {
      return error
    }
  }
  return null
}

/**
 * 验证整个表单
 * @param {Object} formData - 表单数据
 * @param {Object} rulesMap - 字段规则映射 { fieldName: [rules] }
 * @returns {Object} 错误对象 { fieldName: errorMessage }
 */
export function validateForm(formData, rulesMap) {
  const errors = {}

  for (const [field, rules] of Object.entries(rulesMap)) {
    const error = validateField(formData[field], rules)
    if (error) {
      errors[field] = error
    }
  }

  return errors
}

/**
 * 预定义的验证规则集
 */
export const commonRules = {
  // 标题验证规则
  title: [
    createRule(ValidatorType.REQUIRED, null, '请输入标题'),
    createRule(ValidatorType.MAX_LENGTH, 100, '标题不能超过 100 个字符')
  ],

  // 描述验证规则
  description: [
    createRule(ValidatorType.MAX_LENGTH, 500, '描述不能超过 500 个字符')
  ],

  // 年份验证规则
  year: [
    createRule(ValidatorType.YEAR, null, '请输入有效的年份')
  ],

  // 图片文件验证规则
  imageFile: [
    createRule(ValidatorType.FILE_SIZE, 10, '图片大小不能超过 10MB'),
    createRule(ValidatorType.FILE_TYPE, ['image/jpeg', 'image/png', 'image/webp'], '仅支持 JPG、PNG、WebP 格式')
  ],

  // 必填项验证
  required: [
    createRule(ValidatorType.REQUIRED)
  ]
}

/**
 * 检查表单是否有错误
 * @param {Object} errors - 错误对象
 * @returns {boolean} 是否有错误
 */
export function hasErrors(errors) {
  return Object.keys(errors).length > 0
}
