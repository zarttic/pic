<template>
  <v-form @submit.prevent="handleSubmit">
    <v-text-field
      v-model="formData.title"
      label="标题"
      required
      :error-messages="errors.title"
      @blur="validateField('title')"
      class="mb-4"
    ></v-text-field>

    <v-text-field
      v-model="formData.location"
      label="拍摄地点"
      class="mb-4"
    ></v-text-field>

    <v-row>
      <v-col cols="6">
        <v-text-field
          v-model.number="formData.year"
          label="年份"
          type="number"
          :min="1900"
          :max="currentYear"
          :error-messages="errors.year"
          @blur="validateField('year')"
        ></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field
          v-model="formData.camera_model"
          label="相机型号"
        ></v-text-field>
      </v-col>
    </v-row>

    <v-textarea
      v-model="formData.description"
      label="描述"
      rows="3"
      :error-messages="errors.description"
      @blur="validateField('description')"
      class="mb-4"
    ></v-textarea>

    <v-combobox
      v-model="formData.tags"
      label="标签"
      multiple
      chips
      closable-chips
      hint="按 Enter 添加标签"
      class="mb-4"
    ></v-combobox>

    <!-- 文件上传（仅创建模式） -->
    <v-file-input
      v-if="mode === 'create'"
      v-model="fileInput"
      label="选择照片"
      accept="image/*"
      prepend-icon="mdi-camera"
      show-size
      :error-messages="errors.file"
      @change="handleFileSelect"
      class="mb-4"
    ></v-file-input>

    <!-- 图片预览 -->
    <v-card v-if="previewUrl" class="mb-4">
      <v-img :src="previewUrl" max-height="300" contain></v-img>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text color="error" @click="clearFile">移除图片</v-btn>
      </v-card-actions>
    </v-card>

    <!-- 精选设置 -->
    <v-checkbox
      v-model="formData.is_featured"
      label="设为精选照片"
      color="primary"
      class="mb-4"
    ></v-checkbox>

    <!-- 操作按钮 -->
    <v-card-actions class="pa-0">
      <v-btn
        v-if="mode === 'edit'"
        text
        @click="$emit('cancel')"
      >
        取消
      </v-btn>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        type="submit"
        :loading="isSubmitting"
        :disabled="isSubmitting"
      >
        {{ submitText || (mode === 'create' ? '上传' : '保存') }}
      </v-btn>
    </v-card-actions>
  </v-form>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'

const props = defineProps({
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  },
  initialData: {
    type: Object,
    default: null
  },
  isSubmitting: {
    type: Boolean,
    default: false
  },
  submitText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['submit', 'cancel'])

const currentYear = new Date().getFullYear()

const formData = ref({
  title: '',
  location: '',
  year: currentYear,
  camera_model: '',
  description: '',
  tags: [],
  is_featured: false,
  file: null
})

const fileInput = ref([])
const previewUrl = ref('')
const errors = ref({})

// 初始化数据（编辑模式）
onMounted(() => {
  if (props.initialData) {
    formData.value = {
      ...formData.value,
      ...props.initialData,
      tags: props.initialData.tags || []
    }
  }
})

// 监听 initialData 变化
watch(() => props.initialData, (newData) => {
  if (newData) {
    formData.value = {
      ...formData.value,
      ...newData,
      tags: newData.tags || []
    }
  }
}, { deep: true })

const handleFileSelect = (event) => {
  const file = event.target?.files?.[0] || fileInput.value
  if (file instanceof File) {
    formData.value.file = file

    // 创建预览
    const reader = new FileReader()
    reader.onload = (e) => {
      previewUrl.value = e.target.result
    }
    reader.readAsDataURL(file)

    // 清除错误
    if (errors.value.file) {
      errors.value.file = ''
    }
  }
}

const clearFile = () => {
  fileInput.value = []
  formData.value.file = null
  previewUrl.value = ''
}

const validateField = (field) => {
  switch (field) {
    case 'title':
      if (!formData.value.title) {
        errors.value.title = '请输入标题'
      } else {
        errors.value.title = ''
      }
      break
    case 'year':
      if (!formData.value.year || formData.value.year < 1900 || formData.value.year > currentYear) {
        errors.value.year = `请输入有效年份 (1900-${currentYear})`
      } else {
        errors.value.year = ''
      }
      break
  }
}

const validateForm = () => {
  let isValid = true

  // 标题验证
  if (!formData.value.title) {
    errors.value.title = '请输入标题'
    isValid = false
  }

  // 年份验证
  if (!formData.value.year || formData.value.year < 1900 || formData.value.year > currentYear) {
    errors.value.year = `请输入有效年份 (1900-${currentYear})`
    isValid = false
  }

  // 文件验证（仅创建模式）
  if (props.mode === 'create' && !formData.value.file) {
    errors.value.file = '请选择照片文件'
    isValid = false
  }

  return isValid
}

const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  emit('submit', formData.value)
}

const resetForm = () => {
  formData.value = {
    title: '',
    location: '',
    year: currentYear,
    camera_model: '',
    description: '',
    tags: [],
    is_featured: false,
    file: null
  }
  fileInput.value = []
  previewUrl.value = ''
  errors.value = {}
}

defineExpose({
  resetForm
})
</script>

<style scoped>
/* Vuetify 已处理样式 */
</style>
