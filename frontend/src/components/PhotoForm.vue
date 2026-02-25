<template>
  <form @submit.prevent="handleSubmit" class="photo-form">
    <div class="form-group">
      <label for="title">标题</label>
      <input
        id="title"
        v-model="formData.title"
        type="text"
        required
        placeholder="照片标题"
        :class="{ 'input-error': errors.title }"
        @blur="validateField('title')"
      />
      <span v-if="errors.title" class="error-message">{{ errors.title }}</span>
    </div>

    <div class="form-group">
      <label for="location">拍摄地点</label>
      <input
        id="location"
        v-model="formData.location"
        type="text"
        placeholder="拍摄地点"
      />
    </div>

    <div class="form-row">
      <div class="form-group">
        <label for="year">年份</label>
        <input
          id="year"
          v-model.number="formData.year"
          type="number"
          :min="1900"
          :max="currentYear"
          placeholder="2024"
          :class="{ 'input-error': errors.year }"
          @blur="validateField('year')"
        />
        <span v-if="errors.year" class="error-message">{{ errors.year }}</span>
      </div>

      <div class="form-group">
        <label for="camera">相机</label>
        <input
          id="camera"
          v-model="formData.camera_model"
          type="text"
          placeholder="相机型号"
        />
      </div>
    </div>

    <div class="form-group">
      <label for="description">描述</label>
      <textarea
        id="description"
        v-model="formData.description"
        rows="3"
        placeholder="照片描述"
        :class="{ 'input-error': errors.description }"
        @blur="validateField('description')"
      ></textarea>
      <span v-if="errors.description" class="error-message">{{ errors.description }}</span>
    </div>

    <div class="form-group">
      <label for="tags">标签</label>
      <input
        id="tags"
        v-model="tagsInput"
        type="text"
        placeholder="标签，用逗号分隔"
      />
      <div v-if="parsedTags.length > 0" class="tags-preview">
        <span v-for="tag in parsedTags" :key="tag" class="tag">
          {{ tag }}
        </span>
      </div>
    </div>

    <!-- 文件上传（仅创建模式） -->
    <div v-if="mode === 'create'" class="form-group">
      <label for="file">选择照片</label>
      <input
        id="file"
        type="file"
        accept="image/*"
        @change="handleFileSelect"
        required
        :class="{ 'input-error': errors.file }"
      />
      <span v-if="errors.file" class="error-message">{{ errors.file }}</span>

      <!-- 图片预览 -->
      <div v-if="previewUrl" class="image-preview">
        <img :src="previewUrl" alt="预览" />
        <button type="button" class="btn-remove-preview" @click="clearFile">
          ×
        </button>
      </div>
    </div>

    <div class="form-group">
      <label class="checkbox-label">
        <input type="checkbox" v-model="formData.is_featured" />
        设为精选
      </label>
    </div>

    <div class="form-actions">
      <button type="button" class="btn-secondary" @click="handleCancel">
        取消
      </button>
      <button type="submit" class="btn-primary" :disabled="isSubmitting">
        {{ isSubmitting ? submitText : submitButtonText }}
      </button>
    </div>
  </form>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { useNotificationStore } from '../stores/notification'

const props = defineProps({
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  isSubmitting: {
    type: Boolean,
    default: false
  },
  submitText: {
    type: String,
    default: '保存中...'
  }
})

const emit = defineEmits(['submit', 'cancel'])

const notification = useNotificationStore()

const currentYear = new Date().getFullYear()

const defaultFormData = {
  title: '',
  description: '',
  location: '',
  year: currentYear,
  camera_model: '',
  is_featured: false
}

const formData = ref({ ...defaultFormData })
const tagsInput = ref('')
const selectedFile = ref(null)
const previewUrl = ref(null)
const errors = ref({})

// 解析标签
const parsedTags = computed(() => {
  return tagsInput.value
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag)
})

// 提交按钮文本
const submitButtonText = computed(() => {
  return props.mode === 'create' ? '上传照片' : '保存'
})

// 监听初始数据变化（编辑模式）
watch(() => props.initialData, (newData) => {
  if (props.mode === 'edit' && newData) {
    formData.value = {
      title: newData.title || '',
      description: newData.description || '',
      location: newData.location || '',
      year: newData.year || currentYear,
      camera_model: newData.camera_model || '',
      is_featured: newData.is_featured || false
    }
    tagsInput.value = newData.tags ? newData.tags.join(', ') : ''
  }
}, { immediate: true, deep: true })

// 文件选择处理
const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    errors.value.file = '请选择图片文件'
    return
  }

  // 验证文件大小（最大 10MB）
  const maxSize = 10 * 1024 * 1024
  if (file.size > maxSize) {
    errors.value.file = '文件大小不能超过 10MB'
    return
  }

  delete errors.value.file
  selectedFile.value = file

  // 创建预览
  const reader = new FileReader()
  reader.onload = (e) => {
    previewUrl.value = e.target.result
  }
  reader.readAsDataURL(file)
}

// 清除文件
const clearFile = () => {
  selectedFile.value = null
  previewUrl.value = null
  const fileInput = document.getElementById('file')
  if (fileInput) {
    fileInput.value = ''
  }
}

// 验证字段
const validateField = (field) => {
  switch (field) {
    case 'title':
      if (!formData.value.title.trim()) {
        errors.value.title = '请输入照片标题'
      } else if (formData.value.title.length > 100) {
        errors.value.title = '标题不能超过 100 个字符'
      } else {
        delete errors.value.title
      }
      break

    case 'year':
      if (formData.value.year && (formData.value.year < 1900 || formData.value.year > currentYear)) {
        errors.value.year = `年份应在 1900-${currentYear} 之间`
      } else {
        delete errors.value.year
      }
      break

    case 'description':
      if (formData.value.description && formData.value.description.length > 500) {
        errors.value.description = '描述不能超过 500 个字符'
      } else {
        delete errors.value.description
      }
      break
  }
}

// 验证整个表单
const validateForm = () => {
  validateField('title')
  validateField('year')
  validateField('description')

  if (props.mode === 'create' && !selectedFile.value) {
    errors.value.file = '请选择照片文件'
  }

  return Object.keys(errors.value).length === 0
}

// 提交表单
const handleSubmit = () => {
  if (!validateForm()) {
    notification.error('请修正表单中的错误')
    return
  }

  const submitData = {
    ...formData.value,
    tags: parsedTags.value
  }

  if (props.mode === 'create') {
    submitData.file = selectedFile.value
  }

  emit('submit', submitData)
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

// 重置表单
const resetForm = () => {
  formData.value = { ...defaultFormData }
  tagsInput.value = ''
  selectedFile.value = null
  previewUrl.value = null
  errors.value = {}
}

// 清理预览 URL
onUnmounted(() => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
})

// 暴露方法
defineExpose({
  resetForm,
  validateForm
})
</script>

<style scoped>
.photo-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--spacing-md);
}

label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
  text-transform: none;
}

input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

input,
textarea {
  background: var(--bg-primary);
  border: 1px solid rgba(201, 169, 98, 0.3);
  color: var(--text-primary);
  padding: 0.75rem;
  border-radius: 4px;
  font-family: inherit;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--accent-gold);
}

input[type="file"] {
  cursor: pointer;
}

.input-error {
  border-color: #ef4444;
}

.error-message {
  color: #ef4444;
  font-size: 0.85rem;
}

/* 标签预览 */
.tags-preview {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-xs);
}

.tag {
  background: rgba(201, 169, 98, 0.2);
  color: var(--accent-gold);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.85rem;
}

/* 图片预览 */
.image-preview {
  position: relative;
  margin-top: var(--spacing-sm);
  border-radius: 8px;
  overflow: hidden;
  background: var(--bg-primary);
}

.image-preview img {
  width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.btn-remove-preview {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.3s ease;
}

.btn-remove-preview:hover {
  background: rgba(239, 68, 68, 0.9);
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-md);
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: 0.1em;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
  border: none;
}

.btn-primary {
  background: var(--accent-gold);
  color: var(--bg-primary);
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-warm);
  transform: translateY(-2px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--text-secondary);
}

.btn-secondary:hover {
  border-color: var(--accent-gold);
  color: var(--accent-gold);
}

@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
  }
}
</style>
