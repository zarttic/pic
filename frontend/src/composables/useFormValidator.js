import { ref, reactive } from 'vue'
import { validateField, validateForm, hasErrors } from '../utils/validators'

/**
 * 表单验证组合式函数
 * 提供表单验证状态管理和验证方法
 *
 * @example
 * const { errors, validate, validateOnBlur, clearErrors } = useFormValidator(rules)
 *
 * const rules = {
 *   title: commonRules.title,
 *   year: commonRules.year
 * }
 *
 * // 验证整个表单
 * const isValid = await validate(formData)
 *
 * // 字段 blur 时验证
 * validateOnBlur('title', formData.title)
 */
export function useFormValidator(rulesMap) {
  // 错误状态
  const errors = reactive({})
  // 是否已提交（用于显示所有错误）
  const submitted = ref(false)
  // 正在验证的字段
  const validating = ref(false)

  /**
   * 验证单个字段
   * @param {string} field - 字段名
   * @param {any} value - 字段值
   * @returns {string|null} 错误消息或 null
   */
  const validateSingleField = (field, value) => {
    const rules = rulesMap[field]
    if (!rules) {
      return null
    }

    const error = validateField(value, rules)

    if (error) {
      errors[field] = error
    } else {
      delete errors[field]
    }

    return error
  }

  /**
   * 验证整个表单
   * @param {Object} formData - 表单数据
   * @returns {Promise<boolean>} 表单是否有效
   */
  const validate = async (formData) => {
    submitted.value = true
    validating.value = true

    // 清空之前错误
    Object.keys(errors).forEach(key => delete errors[key])

    // 执行验证
    const newErrors = validateForm(formData, rulesMap)

    // 合并错误
    Object.assign(errors, newErrors)

    validating.value = false

    return !hasErrors(errors)
  }

  /**
   * 在字段 blur 时验证
   * @param {string} field - 字段名
   * @param {any} value - 字段值
   */
  const validateOnBlur = (field, value) => {
    // 只在提交后或字段已有错误时进行实时验证
    if (submitted.value || errors[field]) {
      validateSingleField(field, value)
    }
  }

  /**
   * 清空所有错误
   */
  const clearErrors = () => {
    Object.keys(errors).forEach(key => delete errors[key])
    submitted.value = false
  }

  /**
   * 清空单个字段错误
   * @param {string} field - 字段名
   */
  const clearFieldError = (field) => {
    delete errors[field]
  }

  /**
   * 设置字段错误
   * @param {string} field - 字段名
   * @param {string} message - 错误消息
   */
  const setFieldError = (field, message) => {
    errors[field] = message
  }

  /**
   * 获取字段错误
   * @param {string} field - 字段名
   * @returns {string|null} 错误消息
   */
  const getFieldError = (field) => {
    return errors[field] || null
  }

  /**
   * 字段是否有错误
   * @param {string} field - 字段名
   * @returns {boolean}
   */
  const hasError = (field) => {
    return !!errors[field]
  }

  return {
    errors,
    submitted,
    validating,
    validate,
    validateOnBlur,
    validateSingleField,
    clearErrors,
    clearFieldError,
    setFieldError,
    getFieldError,
    hasError,
    hasErrors: () => hasErrors(errors)
  }
}
