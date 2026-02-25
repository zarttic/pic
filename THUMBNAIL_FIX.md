# 缩略图不显示问题修复

## 问题描述
前端 Gallery 页面的"精选作品"不显示缩略图。

## 根本原因

### 1. 指令使用错误 ❌
`v-lazyload` 指令是为原生 `<img>` 标签设计的，通过设置 `img.src` 属性来加载图片。

但在 Gallery.vue 中，它被错误地应用到了 Vuetify 的 `<v-img>` 组件上：

```vue
<!-- 错误用法 ❌ -->
<v-img
  v-lazyload="{
    src: getImageUrl(photo.thumbnail_path || photo.file_path),
    placeholder: '/placeholder.jpg'
  }"
/>
```

**为什么不起作用？**
- `<v-img>` 不是原生 `<img>` 标签
- `v-lazyload` 指令设置的是 `el.src`，但对 Vue 组件来说，这不会触发 Vue 的响应式更新
- Vuetify 的 `<v-img>` 组件有自己的 props 系统，需要通过 `:src` prop 传递图片地址

### 2. 正确的做法 ✅

Vuetify 的 `<v-img>` 组件本身就支持懒加载和占位符功能：

```vue
<!-- 正确用法 ✅ -->
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

## 修复方案

### 修改文件
`frontend/src/views/Gallery.vue`

### 修改内容

**修改前**:
```vue
<v-img
  v-lazyload="{
    src: getImageUrl(photo.thumbnail_path || photo.file_path),
    placeholder: '/placeholder.jpg'
  }"
  :alt="photo.title"
  height="250"
  cover
  class="cursor-pointer"
>
  <!-- 悬停覆盖层 -->
  <v-overlay ...>
```

**修改后**:
```vue
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
  <v-overlay ...>
```

## Vuetify v-img 的优势

1. **内置懒加载**: Vuetify 的 `<v-img>` 自动处理图片加载
2. **占位符支持**: 通过 `placeholder` slot 可以显示加载中的状态
3. **响应式**: 与 Vue 的响应式系统完美集成
4. **性能优化**: 自动处理图片尺寸、加载状态等

## 测试验证

1. 访问 Gallery 页面 (`http://localhost:5173/gallery`)
2. 检查是否显示缩略图
3. 检查加载过程中是否显示加载动画
4. 检查悬停效果是否正常

## 注意事项

### 何时使用 v-lazyload 指令？

`v-lazyload` 指令适用于原生 `<img>` 标签：

```vue
<!-- 适用场景 -->
<img v-lazyload="imageUrl" />
<img v-lazyload="{ src: imageUrl, placeholder: placeholderUrl }" />
```

### 何时使用 v-img？

对于 Vuetify 项目，推荐使用 `<v-img>` 组件：

```vue
<!-- 推荐用法 -->
<v-img
  :src="imageUrl"
  height="250"
  cover
>
  <template v-slot:placeholder>
    <v-progress-circular indeterminate />
  </template>
</v-img>
```

## 其他可能的问题

如果缩略图仍然不显示，检查：

1. **路径格式**: 确保 `thumbnail_path` 格式正确 (`/uploads/xxx.jpg`)
2. **文件存在**: 确保缩略图文件实际存在于 `backend/uploads/` 目录
3. **CORS 配置**: 确保后端 CORS 允许前端域名
4. **静态文件服务**: 确保后端正确配置了 `/uploads` 静态文件路由

## 相关文件

- `frontend/src/views/Gallery.vue` - 图片展示页面
- `frontend/src/directives/lazyload.js` - 懒加载指令（仅用于原生 img 标签）
- `frontend/src/utils/index.js` - getImageUrl 工具函数
- `backend/internal/handlers/photo.go` - 图片路径生成逻辑
