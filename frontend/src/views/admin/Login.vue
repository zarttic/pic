<template>
  <v-container fluid class="login-container fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="login-card" elevation="12">
          <!-- 卡片头部 -->
          <v-card-title class="text-center pb-2">
            <h1 class="text-h4 font-weight-medium">管理后台登录</h1>
          </v-card-title>

          <v-card-subtitle class="text-center pb-6">
            请输入您的账号密码
          </v-card-subtitle>

          <v-card-text>
            <v-form @submit.prevent="handleLogin">
              <v-text-field
                v-model="form.username"
                label="用户名"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                :disabled="loading"
                required
                class="mb-4"
              ></v-text-field>

              <v-text-field
                v-model="form.password"
                label="密码"
                prepend-inner-icon="mdi-lock"
                type="password"
                variant="outlined"
                :disabled="loading"
                required
                class="mb-4"
              ></v-text-field>

              <v-alert
                v-if="error"
                type="error"
                variant="tonal"
                class="mb-4"
              >
                {{ error }}
              </v-alert>

              <v-btn
                type="submit"
                color="primary"
                size="large"
                block
                :loading="loading"
                :disabled="loading"
              >
                登录
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="justify-center pb-4">
            <v-btn
              to="/"
              variant="text"
              color="primary"
            >
              返回首页
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const theme = useTheme()

const form = ref({
  username: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    error.value = '请输入用户名和密码'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const result = await authStore.login(form.value.username, form.value.password)

    if (result.success) {
      // 登录成功，跳转到管理后台
      router.push('/admin')
    } else {
      error.value = result.error
    }
  } catch (err) {
    error.value = '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 切换到后台浅色主题
onMounted(() => {
  theme.global.name.value = 'adminLightTheme'
})
</script>

<style scoped>
.login-container {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  border-radius: 16px;
}
</style>
