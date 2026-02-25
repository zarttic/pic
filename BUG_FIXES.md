# Bug 修复记录

本文档记录了项目中所有重要 Bug 的修复过程和解决方案。

---

## 目录

1. [缩略图路径错误](#1-缩略图路径错误)
2. [后台管理界面鼠标隐藏问题](#2-后台管理界面鼠标隐藏问题)
3. [精选作品缩略图不显示](#3-精选作品缩略图不显示)
4. [全屏查看器关闭按钮无效](#4-全屏查看器关闭按钮无效)
5. [相册功能问题修复](#5-相册功能问题修复)
6. [CORS 配置问题](#6-cors-配置问题)
7. [相册密码验证问题](#7-相册密码验证问题)
8. [相册密码输入界面优化](#8-相册密码输入界面优化)

---

## 1. 缩略图路径错误

### 问题描述

数据库中存储的缩略图路径格式错误，导致前端无法正确加载缩略图。

**错误路径**: `/./uploads/xxx_thumb.jpg`
**正确路径**: `/uploads/xxx_thumb.jpg`

### 根本原因

在 `backend/internal/handlers/photo.go` 的 Create 方法中，路径拼接逻辑有误：

```go
// 错误代码
FilePath:      "/" + uploadPath,      // uploadPath = "./uploads/xxx.jpg"
ThumbnailPath: "/" + thumbnailPath,   // 结果: "/./uploads/xxx.jpg" ❌
```

### 解决方案

**文件**: `backend/internal/handlers/photo.go`

```go
// 修复代码
FilePath:      "/" + strings.TrimPrefix(uploadPath, "./"),
ThumbnailPath: "/" + strings.TrimPrefix(thumbnailPath, "./"),
// 结果: "/uploads/xxx.jpg" ✅
```

### 数据修复

对于已存在的照片，需要修复数据库记录：

**工具**: `backend/cmd/fix-paths/main.go`

```go
// 批量修复路径
if strings.HasPrefix(photo.FilePath, "/./uploads/") {
    correctPath := strings.Replace(photo.FilePath, "/./uploads/", "/uploads/", 1)
    // 更新数据库
}
```

### 验证结果

```bash
curl http://localhost:9421/api/photos?pageSize=1
# 返回: {"file_path":"/uploads/xxx.jpg","thumbnail_path":"/uploads/xxx_thumb.jpg"}
```

---

## 2. 后台管理界面鼠标隐藏问题

### 问题描述

后台管理界面的鼠标被隐藏，无法正常操作。

### 根本原因

前端 CSS 中强制隐藏了 Vuetify 组件的系统光标，用于前台自定义光标效果：

```css
/* 错误代码 - 应用于所有 Vuetify 组件 */
.v-btn, .v-card, .v-list-item, .v-navigation-drawer, .v-app-bar, .v-tab, .v-chip {
  cursor: none !important;
}
```

### 解决方案

**文件**: `frontend/src/styles/main.css`

移除强制隐藏光标的样式，后台管理界面使用默认系统光标：

```css
/* Vuetify 兼容性：后台管理界面使用默认光标 */
/* 前台自定义光标效果已禁用 */
```

### 影响范围

- ✅ 后台管理界面鼠标正常显示
- ⚠️ 前台自定义光标效果被禁用（如需要可单独实现）

---

## 3. 精选作品缩略图不显示

### 问题描述

Gallery 页面的"精选作品"无法显示缩略图。

### 根本原因

`v-lazyload` 指令被错误地应用到了 Vuetify 的 `<v-img>` 组件上。

**错误用法**:
```vue
<!-- ❌ v-lazyload 为原生 img 标签设计，不能用于 Vue 组件 -->
<v-img
  v-lazyload="{
    src: getImageUrl(photo.thumbnail_path || photo.file_path),
    placeholder: '/placeholder.jpg'
  }"
/>
```

### 解决方案

**文件**: `frontend/src/views/Gallery.vue`

使用 Vuetify `<v-img>` 组件内置的加载机制：

```vue
<!-- ✅ 正确用法 -->
<v-img
  :src="getImageUrl(photo.thumbnail_path || photo.file_path)"
  :alt="photo.title"
  height="250"
  cover
>
  <template v-slot:placeholder>
    <div class="d-flex align-center justify-center fill-height">
      <v-progress-circular
        color="grey-lighten-4"
        indeterminate
      ></v-progress-circular>
    </div>
  </template>
</v-img>
```

### 技术要点

- `v-lazyload` 指令适用于原生 `<img>` 标签
- Vuetify 的 `<v-img>` 组件内置懒加载功能
- 通过 `placeholder` slot 显示加载状态

---

## 4. 全屏查看器关闭按钮无效

### 问题描述

Gallery 页面全屏查看器的右上角关闭按钮无法关闭对话框。

### 根本原因

1. 点击事件可能被父元素拦截
2. 工具栏缺少明确的 z-index，可能被图片遮挡
3. 缺少事件传播控制

### 解决方案

**文件**: `frontend/src/views/Gallery.vue`

#### 1. 多种关闭方式

```vue
<!-- ESC 键关闭 -->
<v-dialog @keydown.esc="closeViewer">

<!-- 点击背景关闭 -->
<v-card @click="closeViewer">

<!-- 点击按钮关闭 -->
<v-btn @click.stop="closeViewer">
```

#### 2. 事件传播控制

```vue
<!-- 阻止事件冒泡 -->
<v-toolbar @click.stop>
  <v-btn @click.stop="closeViewer">
</v-toolbar>

<v-card-text @click.stop>
  <!-- 图片内容 -->
</v-card-text>
```

#### 3. Z-Index 层级管理

```vue
<v-toolbar style="z-index: 10;">
  <v-btn style="background: rgba(0,0,0,0.5);">
```

### 用户交互

| 操作 | 效果 |
|------|------|
| ✅ 点击关闭按钮 | 关闭对话框 |
| ✅ 按 ESC 键 | 关闭对话框 |
| ✅ 点击图片外区域 | 关闭对话框 |
| ✅ 点击图片本身 | 不关闭 |

---

## 5. 相册功能问题修复

### 问题描述

1. 管理员无法查看密码保护的相册内容
2. 前台访问密码保护相册后显示异常
3. 缺少相册照片管理页面

### 解决方案

#### 1. 创建相册照片管理页面

**新文件**: `frontend/src/views/admin/AlbumPhotoManagement.vue`

- 支持查看相册中的所有照片
- 支持添加新照片到相册
- 支持从相册移除照片
- 支持拖拽排序

#### 2. 修复管理员访问逻辑

**文件**: `backend/internal/handlers/album.go`

```go
// GetByID - 管理员访问（需要 JWT 认证）
func (h *AlbumHandler) GetByID(c *gin.Context) {
    // 管理员可以查看所有相册，包括密码保护的
    // 返回完整的相册信息
}

// GetPublicAlbum - 前台访问（可能需要相册 token）
func (h *AlbumHandler) GetPublicAlbum(c *gin.Context) {
    // 前台访问需要验证密码（如果相册受保护）
    // 返回公开的相册信息
}
```

#### 3. 修复前台密码验证逻辑

**文件**: `frontend/src/views/AlbumDetail.vue`

- 分离管理员和前台访问逻辑
- 管理员使用 JWT token
- 前台使用相册 token

### 路由配置

```javascript
{
  path: '/admin/albums/:id/photos',
  component: AlbumPhotoManagement,
  meta: { requiresAuth: true }
}
```

---

## 6. CORS 配置问题

### 问题描述

前端使用 `X-Album-Token` header 发送请求时，被 CORS 策略阻止。

```
Request header field x-album-token is not allowed by Access-Control-Allow-Headers in preflight response.
```

### 解决方案

**文件**: `backend/internal/middleware/cors.go`

```go
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Album-Token"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

### 关键变更

添加 `X-Album-Token` 到 `AllowHeaders` 列表：

```go
AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Album-Token"},
```

### 允许的自定义 Headers

- `Authorization` - JWT token (管理后台认证)
- `X-Album-Token` - Album token (相册密码验证)
- `Content-Type` - 请求体类型
- `Origin` - 请求来源

---

## 7. 相册密码验证问题

### 问题描述

1. 输入错误密码后跳转到相册列表页面
2. 输入正确密码后停留在原地，没有进入相册

### 根本原因

**问题1**: 验证失败时调用 `closePasswordDialog()`，该方法会检查相册是否需要认证并跳转

```javascript
// 错误逻辑
const closePasswordDialog = () => {
  showPasswordDialog.value = false
  if (albumStore.currentAlbum?.require_auth) {
    router.push('/albums')  // ❌ 总是跳转
  }
}
```

**问题2**: 验证成功后没有正确处理相册数据获取

```javascript
// 错误逻辑
await albumStore.verifyPassword(route.params.id, passwordInput.value)
// ❌ 没有重新获取相册数据
```

### 解决方案

**文件**: `frontend/src/views/AlbumDetail.vue`

#### 修复验证失败处理

```javascript
const handleVerifyPassword = async () => {
  try {
    await albumStore.verifyPassword(route.params.id, passwordInput.value)
    // 验证成功，重新获取相册数据
    const token = localStorage.getItem(`album_token_${route.params.id}`)
    await fetchAlbumWithToken(route.params.id, token)
    // 成功后关闭对话框（不调用 closePasswordDialog，避免跳转）
    showPasswordDialog.value = false
    passwordInput.value = ''
    passwordError.value = ''
  } catch (error) {
    // ❌ 不调用 closePasswordDialog，保持在当前页面
    passwordError.value = '密码错误，请重试'
  }
}
```

#### 改进对话框关闭逻辑

```javascript
const closePasswordDialog = () => {
  showPasswordDialog.value = false
  passwordInput.value = ''
  passwordError.value = ''
  // 只有在相册需要认证时才跳转
  if (albumStore.currentAlbum?.require_auth) {
    router.push('/albums')
  }
}
```

---

## 8. 相册密码输入界面优化

### 问题描述

密码输入对话框样式简陋，不符合摄影网站的艺术风格。

### 旧设计问题

- ❌ 样式过于简单，缺乏设计感
- ❌ 使用原生表单元素，不符合网站整体风格
- ❌ 交互体验单调，缺少动画和反馈
- ❌ 没有显示/隐藏密码功能
- ❌ 错误提示不明显

### 新设计特点

#### 1. 视觉升级

- 渐变背景 + 装饰性光晕效果
- Vuetify 输入框组件（带图标）
- 精美的锁图标（80px）
- 金色主题色贯穿整体

#### 2. 交互改进

| 功能 | 旧版 | 新版 |
|------|------|------|
| 显示/隐藏密码 | ❌ | ✅ 眼睛图标切换 |
| 自动聚焦 | ✅ | ✅ |
| 加载状态 | 文字"验证中..." | ✅ 按钮旋转动画 |
| 错误提示 | 红色文字 | ✅ 可关闭的Alert组件 |
| 按钮状态 | 禁用 | ✅ 禁用 + 加载状态 |
| 帮助提示 | ❌ | ✅ 底部说明文字 |

#### 3. 动画效果

- 对话框从底部滑入
- 错误提示平滑展开
- 装饰元素金色光晕

### 实现代码

**文件**: `frontend/src/views/AlbumDetail.vue`

```vue
<v-dialog v-model="showPasswordDialog" max-width="500">
  <v-card class="password-card">
    <!-- 装饰背景 -->
    <div class="card-decoration">
      <div class="decoration-circle"></div>
    </div>

    <!-- 锁图标 -->
    <v-icon size="80" color="primary">mdi-lock-outline</v-icon>

    <!-- 标题 -->
    <h2 class="dialog-title">私密相册</h2>

    <!-- 密码输入框 -->
    <v-text-field
      v-model="passwordInput"
      :type="showPassword ? 'text' : 'password'"
      prepend-inner-icon="mdi-key-outline"
      :append-inner-icon="showPassword ? 'mdi-eye-outline' : 'mdi-eye-off-outline'"
      @click:append-inner="showPassword = !showPassword"
      variant="outlined"
      density="comfortable"
    />

    <!-- 按钮 -->
    <v-btn color="primary" @click="handleVerifyPassword" :loading="verifying">
      验证访问
    </v-btn>
  </v-card>
</v-dialog>
```

### 样式设计

```css
.password-card {
  background: linear-gradient(135deg, rgba(20, 20, 20, 0.98) 0%, rgba(10, 10, 10, 0.98) 100%);
  border: 1px solid rgba(201, 169, 98, 0.2);
  padding: 3rem 2.5rem;
}

.decoration-circle {
  background: radial-gradient(circle, rgba(201, 169, 98, 0.08) 0%, transparent 70%);
  filter: blur(40px);
}
```

---

## 总结

### 修复统计

| 类别 | 数量 |
|------|------|
| 路径问题 | 1 |
| UI/UX 问题 | 3 |
| 功能问题 | 2 |
| 配置问题 | 1 |
| 交互问题 | 1 |
| **总计** | **8** |

### 涉及文件

**后端**:
- `backend/internal/handlers/photo.go`
- `backend/internal/handlers/album.go`
- `backend/internal/middleware/cors.go`
- `backend/cmd/fix-paths/main.go` (新增)

**前端**:
- `frontend/src/styles/main.css`
- `frontend/src/views/Gallery.vue`
- `frontend/src/views/AlbumDetail.vue`
- `frontend/src/views/admin/AlbumPhotoManagement.vue` (新增)

### 经验教训

1. **路径处理要标准化**：避免 `./` 和 `/` 混用
2. **Vue 指令与组件的区别**：指令操作 DOM，组件通过 props
3. **事件传播要控制**：使用 `.stop` 修饰符防止冒泡
4. **用户体验要重视**：加载状态、错误提示、动画效果
5. **文档要完整**：详细记录问题和解决方案
