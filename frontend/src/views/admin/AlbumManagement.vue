<template>
  <div>
    <!-- 页面标题 -->
    <v-row class="mb-4" align="center">
      <v-col>
        <h1 class="text-h4 font-weight-light">相册管理</h1>
      </v-col>
      <v-col cols="auto">
        <v-btn color="primary" @click="showCreateDialog = true">
          <v-icon start>mdi-plus</v-icon>
          创建相册
        </v-btn>
      </v-col>
    </v-row>

    <!-- 加载骨架屏 -->
    <v-row v-if="albumStore.loading">
      <v-col v-for="i in 4" :key="i" cols="12" sm="6" md="4" lg="3">
        <v-skeleton-loader type="card" height="300"></v-skeleton-loader>
      </v-col>
    </v-row>

    <!-- 空状态 -->
    <v-card v-else-if="albumStore.albums.length === 0" class="text-center pa-8">
      <v-icon size="64" color="grey-lighten-1">mdi-folder-multiple</v-icon>
      <p class="text-h6 mt-4">暂无相册</p>
      <p class="text-body-2 text-medium-emphasis">点击上方按钮创建第一个相册</p>
    </v-card>

    <!-- 相册网格 -->
    <v-row v-else>
      <v-col
        v-for="album in albumStore.albums"
        :key="album.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card hover>
          <!-- 封面图片 -->
          <v-img
            :src="getCoverPhoto(album)"
            height="200"
            cover
            class="cursor-pointer"
            @click="openAlbumDetail(album)"
          >
            <template v-slot:placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-avatar size="80" color="grey-lighten-2">
                  <span class="text-h3">{{ album.name.charAt(0) }}</span>
                </v-avatar>
              </v-row>
            </template>

            <!-- 悬停覆盖层 -->
            <v-overlay
              absolute
              contained
              class="align-center justify-center"
            >
              <v-chip color="white" variant="flat">
                <v-icon start>mdi-image-multiple</v-icon>
                {{ album.photos?.length || 0 }} 张照片
              </v-chip>
            </v-overlay>
          </v-img>

          <!-- 卡片内容 -->
          <v-card-title>{{ album.name }}</v-card-title>
          <v-card-text>
            <p class="text-body-2 text-medium-emphasis mb-2">
              {{ album.description || '暂无描述' }}
            </p>
            <v-chip
              :color="album.is_protected ? 'warning' : 'success'"
              size="small"
            >
              <v-icon start size="small">
                {{ album.is_protected ? 'mdi-lock' : 'mdi-lock-open' }}
              </v-icon>
              {{ album.is_protected ? '已加密' : '公开' }}
            </v-chip>
          </v-card-text>

          <!-- 操作按钮 -->
          <v-card-actions>
            <v-btn
              icon="mdi-pencil"
              size="small"
              variant="text"
              color="primary"
              @click="editAlbum(album)"
            ></v-btn>
            <v-btn
              :icon="album.is_protected ? 'mdi-lock-open' : 'mdi-lock'"
              size="small"
              variant="text"
              color="warning"
              @click="togglePassword(album)"
            ></v-btn>
            <v-btn
              icon="mdi-image-multiple"
              size="small"
              variant="text"
              color="info"
              @click="openPhotoManager(album)"
            ></v-btn>
            <v-btn
              icon="mdi-delete"
              size="small"
              variant="text"
              color="error"
              @click="deleteAlbum(album.id)"
            ></v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 创建/编辑相册对话框 -->
    <v-dialog v-model="showCreateDialog" max-width="500">
      <v-card>
        <v-card-title>{{ editMode ? '编辑相册' : '创建相册' }}</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="handleCreate">
            <v-text-field
              v-model="createForm.name"
              label="相册名称"
              required
              class="mb-4"
            ></v-text-field>
            <v-textarea
              v-model="createForm.description"
              label="描述"
              rows="3"
            ></v-textarea>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showCreateDialog = false">取消</v-btn>
          <v-btn color="primary" @click="handleCreate">
            {{ editMode ? '保存' : '创建' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 设置密码对话框 -->
    <v-dialog v-model="showPasswordDialog" max-width="400">
      <v-card>
        <v-card-title>设置相册密码</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="passwordForm.password"
            type="password"
            label="密码"
            required
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showPasswordDialog = false">取消</v-btn>
          <v-btn color="primary" @click="handleSetPassword">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlbumStore } from '@/stores/albums'
import { useNotificationStore } from '@/stores/notification'
import { useConfirm } from '@/composables/useConfirm'
import { getImageUrl } from '@/utils/index'

const router = useRouter()
const albumStore = useAlbumStore()
const notification = useNotificationStore()
const confirm = useConfirm()

const showCreateDialog = ref(false)
const showPasswordDialog = ref(false)
const editMode = ref(false)
const editAlbumId = ref(null)
const passwordAlbumId = ref(null)

const createForm = ref({
  name: '',
  description: ''
})

const passwordForm = ref({
  password: ''
})

onMounted(() => {
  albumStore.fetchAlbums()
})

const getCoverPhoto = (album) => {
  if (album.cover_photo_id && album.photos) {
    const photo = album.photos.find(p => p.id === album.cover_photo_id)
    if (photo) return getImageUrl(photo.thumbnail_path || photo.file_path)
  }
  return null
}

const openAlbumDetail = (album) => {
  router.push(`/albums/${album.id}`)
}

const editAlbum = (album) => {
  editMode.value = true
  editAlbumId.value = album.id
  createForm.value = {
    name: album.name,
    description: album.description || ''
  }
  showCreateDialog.value = true
}

const handleCreate = async () => {
  if (!createForm.value.name) {
    notification.error('请输入相册名称')
    return
  }

  try {
    if (editMode.value) {
      await albumStore.updateAlbum(editAlbumId.value, createForm.value)
      notification.success('相册更新成功！')
    } else {
      await albumStore.createAlbum(createForm.value)
      notification.success('相册创建成功！')
    }
    showCreateDialog.value = false
    createForm.value = { name: '', description: '' }
    editMode.value = false
  } catch (error) {
    notification.error('操作失败：' + error.message)
  }
}

const togglePassword = async (album) => {
  if (album.is_protected) {
    const result = await confirm({
      type: 'warning',
      title: '移除密码',
      message: '确定要移除相册密码吗？移除后相册将公开访问。',
      confirmText: '确定'
    })
    if (result) {
      try {
        await albumStore.setAlbumPassword(album.id, '')
        notification.success('密码已移除')
      } catch (error) {
        notification.error('操作失败：' + error.message)
      }
    }
  } else {
    passwordAlbumId.value = album.id
    passwordForm.value.password = ''
    showPasswordDialog.value = true
  }
}

const handleSetPassword = async () => {
  if (!passwordForm.value.password) {
    notification.error('请输入密码')
    return
  }

  try {
    await albumStore.setAlbumPassword(passwordAlbumId.value, passwordForm.value.password)
    notification.success('密码设置成功！')
    showPasswordDialog.value = false
  } catch (error) {
    notification.error('设置失败：' + error.message)
  }
}

const openPhotoManager = (album) => {
  router.push(`/admin/albums/${album.id}/photos`)
}

const deleteAlbum = async (id) => {
  const result = await confirm({
    type: 'danger',
    title: '删除相册',
    message: '确定要删除这个相册吗？此操作不可撤销。',
    confirmText: '删除'
  })

  if (result) {
    try {
      await albumStore.deleteAlbum(id)
      notification.success('相册删除成功！')
    } catch (error) {
      notification.error('删除失败：' + error.message)
    }
  }
}
</script>

<style scoped>
/* Vuetify 已处理大部分样式 */
</style>
