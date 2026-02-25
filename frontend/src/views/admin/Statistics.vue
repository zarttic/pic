<template>
  <div>
    <!-- 页面标题 -->
    <h1 class="text-h4 font-weight-light mb-6">统计数据</h1>

    <!-- 统计卡片网格 -->
    <v-row class="mb-6">
      <v-col cols="12" sm="6" md="3">
        <v-card>
          <v-card-text class="text-center">
            <v-icon size="48" color="primary" class="mb-2">mdi-image-multiple</v-icon>
            <div class="text-h4 font-weight-bold">{{ stats.totalPhotos }}</div>
            <div class="text-body-2 text-medium-emphasis">总照片数</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-card>
          <v-card-text class="text-center">
            <v-icon size="48" color="secondary" class="mb-2">mdi-folder-multiple</v-icon>
            <div class="text-h4 font-weight-bold">{{ stats.totalAlbums }}</div>
            <div class="text-body-2 text-medium-emphasis">总相册数</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-card>
          <v-card-text class="text-center">
            <v-icon size="48" color="success" class="mb-2">mdi-eye</v-icon>
            <div class="text-h4 font-weight-bold">{{ stats.totalViews }}</div>
            <div class="text-body-2 text-medium-emphasis">总浏览量</div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="3">
        <v-card>
          <v-card-text class="text-center">
            <v-icon size="48" color="warning" class="mb-2">mdi-lock</v-icon>
            <div class="text-h4 font-weight-bold">{{ stats.protectedAlbums }}</div>
            <div class="text-body-2 text-medium-emphasis">加密相册</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 热门照片 TOP 10 -->
    <v-card>
      <v-card-title>热门照片 TOP 10</v-card-title>
      <v-card-text>
        <v-list v-if="topPhotos.length > 0">
          <v-list-item
            v-for="(photo, index) in topPhotos"
            :key="photo.id"
          >
            <template v-slot:prepend>
              <v-avatar color="primary" size="32">
                <span class="text-caption font-weight-bold">{{ index + 1 }}</span>
              </v-avatar>
            </template>

            <v-list-item-title>{{ photo.title }}</v-list-item-title>
            <v-list-item-subtitle>{{ photo.location || '未知地点' }}</v-list-item-subtitle>

            <template v-slot:append>
              <v-chip color="primary" variant="flat">
                <v-icon start>mdi-eye</v-icon>
                {{ photo.view_count }}
              </v-chip>
            </template>
          </v-list-item>
        </v-list>

        <v-empty-state
          v-else
          icon="mdi-image-off"
          text="暂无数据"
        ></v-empty-state>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePhotoStore } from '@/stores/photos'
import { useAlbumStore } from '@/stores/albums'
import { getImageUrl } from '@/utils/index'

const photoStore = usePhotoStore()
const albumStore = useAlbumStore()

const stats = computed(() => ({
  totalPhotos: photoStore.photos.length,
  totalAlbums: albumStore.albums.length,
  totalViews: photoStore.photos.reduce((sum, p) => sum + (p.view_count || 0), 0),
  protectedAlbums: albumStore.albums.filter(a => a.is_protected).length,
}))

const topPhotos = computed(() => {
  return [...photoStore.photos]
    .sort((a, b) => (b.view_count || 0) - (a.view_count || 0))
    .slice(0, 10)
})

onMounted(() => {
  photoStore.fetchPhotos()
  albumStore.fetchAlbums()
})
</script>

<style scoped>
/* Vuetify 已处理样式 */
</style>
