<template>
  <div>
    <!-- 页面标题 -->
    <v-row class="mb-4" align="center">
      <v-col cols="auto">
        <v-btn
          variant="text"
          prepend-icon="mdi-arrow-left"
          @click="goBack"
        >
          返回
        </v-btn>
      </v-col>
      <v-col>
        <h1 class="text-h4 font-weight-light">{{ album?.name || '相册照片管理' }}</h1>
        <p class="text-body-2 text-medium-emphasis mb-0">{{ album?.description }}</p>
      </v-col>
      <v-col cols="auto">
        <v-btn color="primary" @click="showAddPhotoDialog = true">
          <v-icon start>mdi-plus</v-icon>
          添加照片
        </v-btn>
      </v-col>
    </v-row>

    <!-- 加载状态 -->
    <v-skeleton-loader v-if="loading" type="card" height="400"></v-skeleton-loader>

    <!-- 空状态 -->
    <v-card v-else-if="!albumPhotos || albumPhotos.length === 0" class="text-center pa-8">
      <v-icon size="64" color="grey-lighten-1">mdi-image-off</v-icon>
      <p class="text-h6 mt-4">相册中暂无照片</p>
      <v-btn
        color="primary"
        class="mt-4"
        @click="showAddPhotoDialog = true"
      >
        添加第一张照片
      </v-btn>
    </v-card>

    <!-- 照片网格 -->
    <v-row v-else>
      <v-col
        v-for="(albumPhoto, index) in albumPhotos"
        :key="albumPhoto.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card>
          <v-img
            :src="getImageUrl(albumPhoto.thumbnail_path || albumPhoto.file_path)"
            :alt="albumPhoto.title"
            height="200"
            cover
          >
            <!-- 悬停覆盖层 -->
            <v-overlay
              absolute
              contained
              class="align-center justify-center"
            >
              <div class="text-center">
                <v-btn
                  icon="mdi-close"
                  size="small"
                  color="error"
                  variant="flat"
                  @click="removePhoto(albumPhoto.id)"
                  class="mb-2"
                ></v-btn>
                <div class="text-white text-subtitle-1">{{ albumPhoto.title }}</div>
                <div class="text-white text-caption">{{ albumPhoto.location }}</div>
              </div>
            </v-overlay>
          </v-img>

          <!-- 操作按钮 -->
          <v-card-actions>
            <v-btn
              icon="mdi-arrow-up"
              size="small"
              variant="text"
              :disabled="index === 0"
              @click="movePhoto(albumPhoto.id, 'up')"
            ></v-btn>
            <v-btn
              icon="mdi-arrow-down"
              size="small"
              variant="text"
              :disabled="index === albumPhotos.length - 1"
              @click="movePhoto(albumPhoto.id, 'down')"
            ></v-btn>
            <v-spacer></v-spacer>
            <v-btn
              icon="mdi-star"
              size="small"
              variant="text"
              :color="album?.cover_photo_id === albumPhoto.id ? 'warning' : 'default'"
              :disabled="album?.cover_photo_id === albumPhoto.id"
              @click="setAsCover(albumPhoto.id)"
            ></v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 添加照片对话框 -->
    <v-dialog v-model="showAddPhotoDialog" max-width="900" persistent>
      <v-card>
        <v-card-title>添加照片到相册</v-card-title>

        <v-card-text>
          <!-- 搜索框 -->
          <v-text-field
            v-model="searchQuery"
            prepend-inner-icon="mdi-magnify"
            label="搜索照片..."
            variant="outlined"
            density="compact"
            hide-details
            class="mb-4"
          ></v-text-field>

          <!-- 可选照片列表 -->
          <v-row v-if="availablePhotos.length > 0">
            <v-col
              v-for="photo in availablePhotos"
              :key="photo.id"
              cols="6"
              sm="4"
              md="3"
            >
              <v-card
                :color="selectedPhotoIds.includes(photo.id) ? 'primary' : 'default'"
                :variant="selectedPhotoIds.includes(photo.id) ? 'outlined' : 'elevated'"
                hover
                @click="togglePhotoSelection(photo.id)"
              >
                <v-img
                  :src="getImageUrl(photo.thumbnail_path || photo.file_path)"
                  :alt="photo.title"
                  height="120"
                  cover
                >
                  <v-overlay
                    v-if="selectedPhotoIds.includes(photo.id)"
                    absolute
                    contained
                    class="align-center justify-center"
                  >
                    <v-icon size="48" color="white">mdi-check-circle</v-icon>
                  </v-overlay>
                </v-img>
                <v-card-text class="pa-2">
                  <div class="text-caption text-truncate">{{ photo.title }}</div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <v-empty-state
            v-else
            icon="mdi-image-off"
            text="没有可添加的照片"
          ></v-empty-state>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="showAddPhotoDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            :disabled="selectedPhotoIds.length === 0"
            @click="addPhotosToAlbum"
          >
            添加 {{ selectedPhotoIds.length }} 张照片
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAlbumStore } from '@/stores/albums'
import { usePhotoStore } from '@/stores/photos'
import { useNotificationStore } from '@/stores/notification'
import { useConfirm } from '@/composables/useConfirm'
import { getImageUrl } from '@/utils/index'

const route = useRoute()
const router = useRouter()
const albumStore = useAlbumStore()
const photoStore = usePhotoStore()
const notification = useNotificationStore()
const confirm = useConfirm()

const albumId = route.params.id
const loading = ref(true)
const showAddPhotoDialog = ref(false)
const searchQuery = ref('')
const selectedPhotoIds = ref([])

const album = computed(() => albumStore.albums.find(a => a.id === parseInt(albumId)))
const albumPhotos = computed(() => album.value?.photos || [])

const availablePhotos = computed(() => {
  const albumPhotoIds = albumPhotos.value.map(p => p.id)
  let photos = photoStore.photos.filter(p => !albumPhotoIds.includes(p.id))

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    photos = photos.filter(p =>
      p.title.toLowerCase().includes(query) ||
      (p.location && p.location.toLowerCase().includes(query))
    )
  }

  return photos
})

onMounted(async () => {
  await Promise.all([
    albumStore.fetchAlbums(),
    photoStore.fetchPhotos()
  ])
  loading.value = false
})

const goBack = () => {
  router.push('/admin/albums')
}

const togglePhotoSelection = (photoId) => {
  const index = selectedPhotoIds.value.indexOf(photoId)
  if (index > -1) {
    selectedPhotoIds.value.splice(index, 1)
  } else {
    selectedPhotoIds.value.push(photoId)
  }
}

const addPhotosToAlbum = async () => {
  try {
    await albumStore.addPhotosToAlbum(albumId, selectedPhotoIds.value)
    notification.success(`成功添加 ${selectedPhotoIds.value.length} 张照片`)
    showAddPhotoDialog.value = false
    selectedPhotoIds.value = []
  } catch (error) {
    notification.error('添加失败：' + error.message)
  }
}

const removePhoto = async (photoId) => {
  const result = await confirm({
    type: 'warning',
    title: '移除照片',
    message: '确定要从相册中移除这张照片吗？',
    confirmText: '移除'
  })

  if (result) {
    try {
      await albumStore.removePhotoFromAlbum(albumId, photoId)
      notification.success('照片已从相册移除')
    } catch (error) {
      notification.error('移除失败：' + error.message)
    }
  }
}

const movePhoto = async (photoId, direction) => {
  try {
    await albumStore.movePhotoInAlbum(albumId, photoId, direction)
  } catch (error) {
    notification.error('移动失败：' + error.message)
  }
}

const setAsCover = async (photoId) => {
  try {
    await albumStore.setAlbumCover(albumId, photoId)
    notification.success('封面设置成功')
  } catch (error) {
    notification.error('设置失败：' + error.message)
  }
}
</script>

<style scoped>
/* Vuetify 已处理样式 */
</style>
