<template>
  <div class="image-preview">
    <div v-if="previewUrl" class="preview-container">
      <img :src="previewUrl" :alt="alt" class="preview-image" />
      <div class="preview-overlay">
        <button class="btn-remove" @click="removeImage" title="移除图片">
          ×
        </button>
        <div class="file-info">
          <p class="file-name">{{ fileName }}</p>
          <p class="file-size">{{ formatFileSize(fileSize) }}</p>
        </div>
      </div>
    </div>
    <div v-else class="upload-placeholder" @click="triggerUpload">
      <input
        ref="fileInput"
        type="file"
        :accept="accept"
        @change="handleFileChange"
        hidden
      />
      <div class="placeholder-content">
        <span class="upload-icon">📷</span>
        <p class="placeholder-text">点击选择图片</p>
        <p class="placeholder-hint">支持 JPG、PNG、WebP 格式</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: File,
    default: null
  },
  accept: {
    type: String,
    default: 'image/*'
  },
  alt: {
    type: String,
    default: '预览图片'
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB
  }
})

const emit = defineEmits(['update:modelValue', 'error'])

const fileInput = ref(null)
const previewUrl = ref('')
const fileName = ref('')
const fileSize = ref(0)

watch(() => props.modelValue, (newFile) => {
  if (newFile) {
    generatePreview(newFile)
  } else {
    revokePreview()
  }
}, { immediate: true })

const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (!file) return

  // 验证文件大小
  if (file.size > props.maxSize) {
    emit('error', `文件大小不能超过 ${formatFileSize(props.maxSize)}`)
    return
  }

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    emit('error', '请选择图片文件')
    return
  }

  emit('update:modelValue', file)
}

const generatePreview = (file) => {
  revokePreview()
  previewUrl.value = URL.createObjectURL(file)
  fileName.value = file.name
  fileSize.value = file.size
}

const revokePreview = () => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
    previewUrl.value = ''
    fileName.value = ''
    fileSize.value = 0
  }
}

const removeImage = () => {
  revokePreview()
  emit('update:modelValue', null)
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.image-preview {
  width: 100%;
}

.preview-container {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid var(--accent-gold);
}

.preview-image {
  width: 100%;
  height: auto;
  max-height: 400px;
  object-fit: cover;
  display: block;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8) 0%, transparent 50%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: var(--spacing-md);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.preview-container:hover .preview-overlay {
  opacity: 1;
}

.btn-remove {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  width: 32px;
  height: 32px;
  background: rgba(220, 53, 69, 0.9);
  border: none;
  border-radius: 50%;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  line-height: 1;
}

.btn-remove:hover {
  background: rgba(220, 53, 69, 1);
  transform: scale(1.1);
}

.file-info {
  color: white;
}

.file-name {
  font-weight: 500;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 0.85rem;
  opacity: 0.8;
}

.upload-placeholder {
  border: 2px dashed rgba(201, 169, 98, 0.5);
  border-radius: 8px;
  padding: var(--spacing-xl);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: center;
}

.upload-placeholder:hover {
  border-color: var(--accent-gold);
  background: rgba(201, 169, 98, 0.05);
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
}

.upload-icon {
  font-size: 3rem;
  opacity: 0.6;
}

.placeholder-text {
  font-size: 1.1rem;
  color: var(--text-primary);
  margin: 0;
}

.placeholder-hint {
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin: 0;
}
</style>
