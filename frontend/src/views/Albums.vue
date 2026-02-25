<template>
  <v-container fluid class="albums">
    <!-- 页面标题 -->
    <v-row class="mb-8">
      <v-col cols="12" class="text-center">
        <h1 class="text-h3 font-weight-light" style="font-family: 'Cormorant Garamond', serif;">
          相册集
        </h1>
        <v-divider length="100" thickness="2" color="primary" class="my-4 mx-auto"></v-divider>
      </v-col>
    </v-row>

    <!-- 加载状态 -->
    <v-row v-if="albumStore.loading">
      <v-col v-for="i in 4" :key="i" cols="12" sm="6" md="4">
        <v-skeleton-loader type="image, article" height="300"></v-skeleton-loader>
      </v-col>
    </v-row>

    <!-- 空状态 -->
    <v-row v-else-if="albumStore.albums.length === 0">
      <v-col cols="12">
        <v-card class="text-center pa-8">
          <v-icon size="64" color="grey-lighten-1">mdi-folder-multiple</v-icon>
          <p class="text-h6 mt-4">暂无相册</p>
        </v-card>
      </v-col>
    </v-row>

    <!-- 相册网格 -->
    <v-row v-else>
      <v-col
        v-for="album in albumStore.albums"
        :key="album.id"
        cols="12"
        sm="6"
        md="4"
      >
        <v-card hover @click="viewAlbum(album.id)">
          <v-img
            :src="getCoverPhoto(album)"
            height="200"
            cover
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

            <!-- 加密标识 -->
            <v-chip
              v-if="album.is_protected"
              color="warning"
              variant="flat"
              absolute
              top
              right
              class="ma-2"
            >
              <v-icon start>mdi-lock</v-icon>
              加密
            </v-chip>
          </v-img>

          <v-card-title>{{ album.name }}</v-card-title>
          <v-card-text>
            <p class="text-body-2 text-medium-emphasis mb-0">
              {{ album.description || '暂无描述' }}
            </p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAlbumStore } from '@/stores/albums'
import { getImageUrl } from '@/utils/index'

const router = useRouter()
const albumStore = useAlbumStore()

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

const viewAlbum = (id) => {
  router.push(`/albums/${id}`)
}
</script>

<style scoped>
.albums {
  padding-top: var(--spacing-xl);
  padding-bottom: var(--spacing-xl);
}
</style>
