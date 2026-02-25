<template>
  <div>
    <!-- 页面标题 -->
    <v-row class="mb-4" align="center">
      <v-col>
        <h1 class="text-h4 font-weight-light">照片管理</h1>
      </v-col>
      <v-col cols="auto">
        <v-btn
          v-if="selected.length > 0"
          color="error"
          variant="outlined"
          class="mr-2"
          @click="handleBatchDelete"
        >
          批量删除 ({{ selected.length }})
        </v-btn>
        <v-btn
          v-if="selected.length > 0"
          variant="outlined"
          @click="selected = []"
        >
          取消选择
        </v-btn>
      </v-col>
    </v-row>

    <!-- 搜索和筛选 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="search"
              prepend-inner-icon="mdi-magnify"
              label="搜索照片..."
              variant="outlined"
              density="compact"
              hide-details
            ></v-text-field>
          </v-col>
          <v-col cols="6" md="3">
            <v-select
              v-model="filterYear"
              :items="yearOptions"
              label="年份"
              variant="outlined"
              density="compact"
              hide-details
            ></v-select>
          </v-col>
          <v-col cols="6" md="3">
            <v-select
              v-model="filterFeatured"
              :items="featuredOptions"
              label="精选状态"
              variant="outlined"
              density="compact"
              hide-details
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 上传表单 -->
    <v-card class="mb-4">
      <v-card-title>上传新照片</v-card-title>
      <v-card-text>
        <PhotoForm
          ref="uploadFormRef"
          mode="create"
          :is-submitting="uploading"
          submit-text="上传中..."
          @submit="handleUpload"
        />
      </v-card-text>
    </v-card>

    <!-- 照片列表 -->
    <v-card>
      <v-card-title>
        已上传照片 ({{ filteredPhotos.length }})
      </v-card-title>

      <v-data-table
        v-model="selected"
        :headers="headers"
        :items="filteredPhotos"
        :loading="photoStore.loading"
        show-select
        item-value="id"
        class="elevation-1"
      >
        <!-- 缩略图列 -->
        <template v-slot:item.thumbnail="{ item }">
          <v-img
            :src="getImageUrl(item.thumbnail_path || item.file_path)"
            :alt="item.title"
            width="80"
            height="60"
            cover
            class="my-2"
          ></v-img>
        </template>

        <!-- 标题列 -->
        <template v-slot:item.title="{ item }">
          <div>
            <div class="font-weight-medium">{{ item.title }}</div>
            <div v-if="item.location" class="text-caption text-medium-emphasis">
              {{ item.location }}
            </div>
          </div>
        </template>

        <!-- 年份列 -->
        <template v-slot:item.year="{ item }">
          <v-chip size="small" v-if="item.year">{{ item.year }}</v-chip>
        </template>

        <!-- 浏览量列 -->
        <template v-slot:item.view_count="{ item }">
          <v-chip size="small" color="primary" variant="outlined">
            {{ item.view_count }} 次
          </v-chip>
        </template>

        <!-- 精选列 -->
        <template v-slot:item.is_featured="{ item }">
          <v-chip
            size="small"
            :color="item.is_featured ? 'warning' : 'default'"
            :variant="item.is_featured ? 'flat' : 'outlined'"
          >
            {{ item.is_featured ? '精选' : '普通' }}
          </v-chip>
        </template>

        <!-- 操作列 -->
        <template v-slot:item.actions="{ item }">
          <v-btn
            icon="mdi-pencil"
            size="small"
            variant="text"
            color="primary"
            @click="openEditDialog(item)"
          ></v-btn>
          <v-btn
            icon="mdi-delete"
            size="small"
            variant="text"
            color="error"
            @click="handleDelete(item.id)"
          ></v-btn>
        </template>
      </v-data-table>
    </v-card>

    <!-- 编辑对话框 -->
    <v-dialog v-model="editDialogVisible" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h5 font-weight-light">编辑照片</v-card-title>
        <v-card-text>
          <PhotoForm
            v-if="editDialogVisible"
            ref="editFormRef"
            mode="edit"
            :initial-data="editPhotoData"
            :is-submitting="updating"
            submit-text="保存中..."
            @submit="handleEdit"
            @cancel="closeEditDialog"
          />
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { usePhotoStore } from '@/stores/photos'
import { useNotificationStore } from '@/stores/notification'
import { useConfirm } from '@/composables/useConfirm'
import { getImageUrl } from '@/utils/index'
import PhotoForm from '@/components/PhotoForm.vue'

const photoStore = usePhotoStore()
const notification = useNotificationStore()
const confirm = useConfirm()

const uploadFormRef = ref(null)
const editFormRef = ref(null)
const uploading = ref(false)
const updating = ref(false)
const editDialogVisible = ref(false)

const selected = ref([])
const search = ref('')
const filterYear = ref(null)
const filterFeatured = ref(null)
const editPhotoData = ref(null)

// 数据表头
const headers = [
  { title: '', key: 'thumbnail', width: '100px', sortable: false },
  { title: '标题', key: 'title', sortable: true },
  { title: '年份', key: 'year', width: '100px', sortable: true },
  { title: '浏览量', key: 'view_count', width: '120px', sortable: true },
  { title: '状态', key: 'is_featured', width: '100px', sortable: true },
  { title: '操作', key: 'actions', width: '120px', sortable: false },
]

// 年份选项
const yearOptions = computed(() => {
  const years = photoStore.photos
    .map(p => p.year)
    .filter(year => year)
  const uniqueYears = [...new Set(years)].sort((a, b) => b - a)
  return [{ title: '所有年份', value: null }, ...uniqueYears.map(y => ({ title: y, value: y }))]
})

// 精选选项
const featuredOptions = [
  { title: '全部', value: null },
  { title: '精选', value: true },
  { title: '非精选', value: false },
]

// 过滤后的照片列表
const filteredPhotos = computed(() => {
  let photos = photoStore.photos

  // 搜索过滤
  if (search.value) {
    const query = search.value.toLowerCase()
    photos = photos.filter(p =>
      p.title.toLowerCase().includes(query) ||
      (p.location && p.location.toLowerCase().includes(query)) ||
      (p.description && p.description.toLowerCase().includes(query))
    )
  }

  // 年份过滤
  if (filterYear.value !== null) {
    photos = photos.filter(p => p.year === filterYear.value)
  }

  // 精选过滤
  if (filterFeatured.value !== null) {
    photos = photos.filter(p => p.is_featured === filterFeatured.value)
  }

  return photos
})

onMounted(() => {
  photoStore.fetchPhotos()
})

// 处理上传
const handleUpload = async (formData) => {
  uploading.value = true

  const uploadData = new FormData()
  uploadData.append('file', formData.file)
  uploadData.append('title', formData.title)
  uploadData.append('description', formData.description || '')
  uploadData.append('location', formData.location || '')
  uploadData.append('year', formData.year || new Date().getFullYear())
  uploadData.append('camera_model', formData.camera_model || '')
  uploadData.append('tags', formData.tags.join(','))
  uploadData.append('is_featured', formData.is_featured)

  try {
    await photoStore.uploadPhoto(uploadData)
    uploadFormRef.value?.resetForm()
    notification.success('照片上传成功！')
  } catch (error) {
    notification.error('上传失败：' + error.message)
  } finally {
    uploading.value = false
  }
}

// 处理删除
const handleDelete = async (id) => {
  const result = await confirm({
    type: 'danger',
    title: '删除照片',
    message: '确定要删除这张照片吗？此操作不可撤销。',
    confirmText: '删除'
  })

  if (result) {
    try {
      await photoStore.deletePhoto(id)
      notification.success('删除成功！')
    } catch (error) {
      notification.error('删除失败：' + error.message)
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  const result = await confirm({
    type: 'danger',
    title: '批量删除照片',
    message: `确定要删除选中的 ${selected.value.length} 张照片吗？此操作不可撤销。`,
    confirmText: '删除'
  })

  if (!result) return

  try {
    await photoStore.batchDelete(selected.value)
    notification.success(`成功删除 ${selected.value.length} 张照片`)
    selected.value = []
  } catch (error) {
    notification.error('批量删除失败：' + error.message)
  }
}

// 打开编辑对话框
const openEditDialog = (photo) => {
  editPhotoData.value = photo
  editDialogVisible.value = true
}

// 关闭编辑对话框
const closeEditDialog = () => {
  editDialogVisible.value = false
  editPhotoData.value = null
}

// 处理编辑
const handleEdit = async (formData) => {
  updating.value = true

  try {
    await photoStore.updatePhoto(editPhotoData.value.id, formData)
    notification.success('更新成功！')
    closeEditDialog()
  } catch (error) {
    notification.error('更新失败：' + error.message)
  } finally {
    updating.value = false
  }
}
</script>

<style scoped>
/* Vuetify 已经处理了大部分样式 */
</style>
