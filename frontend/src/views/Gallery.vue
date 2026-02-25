<template>
  <v-container fluid class="gallery">
    <!-- 页面标题 -->
    <v-row class="mb-8">
      <v-col cols="12" class="text-center">
        <h1 class="text-h3 font-weight-light" style="font-family: 'Cormorant Garamond', serif;">
          精选作品
        </h1>
        <v-divider length="100" thickness="2" color="primary" class="my-4 mx-auto"></v-divider>
      </v-col>
    </v-row>

    <!-- 加载骨架屏 -->
    <v-row v-if="photoStore.loading">
      <v-col v-for="i in 12" :key="i" cols="12" sm="6" md="4" lg="3">
        <v-skeleton-loader type="image, article" height="300"></v-skeleton-loader>
      </v-col>
    </v-row>

    <!-- 错误提示 -->
    <v-row v-else-if="photoStore.error">
      <v-col cols="12">
        <v-alert type="error" variant="tonal">
          {{ photoStore.error }}
        </v-alert>
      </v-col>
    </v-row>

    <!-- 照片网格 -->
    <v-row v-else>
      <v-col
        v-for="photo in photoStore.photos"
        :key="photo.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card hover @click="openViewer(photo)">
          <v-img
            :src="getImageUrl(photo.thumbnail_path || photo.file_path)"
            :alt="photo.title"
            height="250"
            cover
            class="cursor-pointer"
          >
            <template v-slot:placeholder>
              <div class="d-flex align-center justify-center fill-height">
                <v-progress-circular
                  color="grey-lighten-4"
                  indeterminate
                ></v-progress-circular>
              </div>
            </template>
            <!-- 悬停覆盖层 -->
            <v-overlay
              absolute
              contained
              class="align-end"
            >
              <div class="pa-4" style="background: linear-gradient(to top, rgba(0,0,0,0.8), transparent);">
                <h3 class="text-h6 text-white mb-1">{{ photo.title }}</h3>
                <p class="text-caption text-white mb-0">
                  {{ photo.location }} · {{ photo.year }}
                </p>
                <p v-if="photo.view_count > 0" class="text-caption text-primary mt-1 mb-0">
                  <v-icon size="small" color="primary">mdi-eye</v-icon>
                  {{ photo.view_count }} 次浏览
                </p>
              </div>
            </v-overlay>
          </v-img>
        </v-card>
      </v-col>
    </v-row>

    <!-- 全屏查看器 -->
    <v-dialog
      v-model="viewerOpen"
      fullscreen
      :scrim="false"
      transition="dialog-bottom-transition"
      @keydown.esc="closeViewer"
    >
      <v-card v-if="selectedPhoto" color="black" @click="closeViewer">
        <!-- 顶部工具栏 -->
        <v-toolbar
          dark
          color="transparent"
          flat
          absolute
          top
          width="100%"
          style="z-index: 10;"
          @click.stop
        >
          <v-spacer></v-spacer>
          <v-btn
            icon
            dark
            @click.stop="closeViewer"
            style="background: rgba(0,0,0,0.5);"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-toolbar>

        <!-- 图片内容 -->
        <v-card-text class="pa-0 fill-height d-flex align-center justify-center" @click.stop>
          <v-img
            :src="getImageUrl(selectedPhoto.file_path)"
            :alt="selectedPhoto.title"
            contain
            max-height="100vh"
          >
            <template v-slot:placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
              </v-row>
            </template>
          </v-img>
        </v-card-text>

        <!-- 底部信息 -->
        <v-card-actions class="pa-4" style="background: rgba(0,0,0,0.7);" @click.stop>
          <div>
            <h3 class="text-h6 text-white">{{ selectedPhoto.title }}</h3>
            <p class="text-caption text-grey-lighten-1 mb-0">
              {{ selectedPhoto.location }} · {{ selectedPhoto.year }}
            </p>
          </div>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { usePhotoStore } from '@/stores/photos'
import { getImageUrl } from '@/utils/index'

const photoStore = usePhotoStore()
const selectedPhoto = ref(null)
const viewerOpen = ref(false)

onMounted(() => {
  photoStore.fetchPhotos()
})

const openViewer = (photo) => {
  selectedPhoto.value = photo
  viewerOpen.value = true
  // 增加浏览量
  photoStore.incrementViewCount(photo.id)
}

const closeViewer = () => {
  viewerOpen.value = false
  selectedPhoto.value = null
}
</script>

<style scoped>
.gallery {
  padding-top: var(--spacing-xl);
  padding-bottom: var(--spacing-xl);
}

/* 保持前台深色主题 */
.gallery :deep(.v-card) {
  background-color: var(--bg-secondary);
}
</style>
