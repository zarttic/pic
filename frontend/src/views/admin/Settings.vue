<template>
  <div>
    <!-- 页面标题 -->
    <h1 class="text-h4 font-weight-light mb-6">系统设置</h1>

    <v-row>
      <!-- 网站信息 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>网站信息</v-card-title>
          <v-card-text>
            <v-form @submit.prevent="saveSiteSettings">
              <v-text-field
                v-model="siteSettings.title"
                label="网站名称"
                class="mb-4"
              ></v-text-field>

              <v-textarea
                v-model="siteSettings.description"
                label="网站描述"
                rows="3"
                class="mb-4"
              ></v-textarea>

              <v-text-field
                v-model="siteSettings.author"
                label="作者名称"
                class="mb-4"
              ></v-text-field>

              <v-btn color="primary" type="submit">保存设置</v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 上传设置 -->
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>上传设置</v-card-title>
          <v-card-text>
            <v-list>
              <v-list-item>
                <v-list-item-title>最大上传大小</v-list-item-title>
                <v-list-item-subtitle>10 MB（服务端限制）</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <v-list-item-title>支持格式</v-list-item-title>
                <v-list-item-subtitle>JPG, PNG, WebP</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <v-list-item-action>
                  <v-switch
                    v-model="uploadSettings.autoThumbnail"
                    color="primary"
                    disabled
                  ></v-switch>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>自动生成缩略图</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item>
                <v-list-item-action>
                  <v-switch
                    v-model="uploadSettings.extractExif"
                    color="primary"
                    disabled
                  ></v-switch>
                </v-list-item-action>
                <v-list-item-content>
                  <v-list-item-title>自动提取 EXIF 信息</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 显示设置 -->
      <v-col cols="12">
        <v-card>
          <v-card-title>显示设置</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-switch
                  v-model="displaySettings.showViewCount"
                  label="显示浏览量统计"
                  color="primary"
                ></v-switch>
              </v-col>
              <v-col cols="12" md="6">
                <v-switch
                  v-model="displaySettings.showTags"
                  label="显示照片标签"
                  color="primary"
                ></v-switch>
              </v-col>
              <v-col cols="12" md="6">
                <v-switch
                  v-model="displaySettings.showExif"
                  label="显示 EXIF 信息"
                  color="primary"
                ></v-switch>
              </v-col>
              <v-col cols="12" md="6">
                <v-switch
                  v-model="displaySettings.enableLazyLoad"
                  label="启用图片懒加载"
                  color="primary"
                ></v-switch>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useNotificationStore } from '@/stores/notification'

const notification = useNotificationStore()

const siteSettings = ref({
  title: 'PicSite',
  description: '一个精美的摄影作品展示平台',
  author: '摄影师'
})

const uploadSettings = ref({
  autoThumbnail: true,
  extractExif: true
})

const displaySettings = ref({
  showViewCount: true,
  showTags: true,
  showExif: false,
  enableLazyLoad: true
})

const saveSiteSettings = () => {
  notification.success('设置已保存！')
}
</script>

<style scoped>
/* Vuetify 已处理样式 */
</style>
