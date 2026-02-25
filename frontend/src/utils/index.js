/**
 * 获取完整的图片URL
 * @param {string} path - 图片路径（如 /uploads/xxx.jpg）
 * @returns {string} 完整的图片URL
 */
export function getImageUrl(path) {
  if (!path) {
    return ''
  }

  // 如果已经是完整URL，直接返回
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }

  // 获取API基础URL（去掉 /api 部分）
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:9421/api'
  const baseUrl = apiUrl.replace('/api', '')

  // 确保路径以 / 开头
  const normalizedPath = path.startsWith('/') ? path : `/${path}`

  return `${baseUrl}${normalizedPath}`
}
