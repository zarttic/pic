# 修复说明

## 问题 1: 鼠标隐藏问题 ✅ 已修复

### 问题描述
前端 CSS 中隐藏了 Vuetify 组件的系统光标，导致后台管理界面鼠标不可见。

### 原因
`frontend/src/styles/main.css` 第 28-31 行设置了自定义光标效果，强制隐藏系统光标：
```css
.v-btn, .v-card, .v-list-item, .v-navigation-drawer, .v-app-bar, .v-tab, .v-chip {
  cursor: none !important;
}
```

### 解决方案
移除了强制隐藏光标的 CSS 规则。

**修改文件**: `frontend/src/styles/main.css`

**修改后**:
```css
/* Vuetify 兼容性：后台管理界面使用默认光标 */
/* 前台自定义光标效果已禁用 */
```

---

## 问题 2: 缩略图路径错误 ✅ 已修复

### 问题描述
前端无法正确加载缩略图，显示空白或加载失败。

### 原因
数据库中存储的缩略图路径格式错误：
- 错误路径: `/./uploads/xxx_thumb.jpg`
- 正确路径: `/uploads/xxx_thumb.jpg`

路径生成时，在 `./uploads/xxx.jpg` 前面又加了 `/`，导致变成 `/./uploads/`。

### 解决方案
在保存路径到数据库前，移除 `./` 前缀。

**修改文件**: `backend/internal/handlers/photo.go`

**修改位置**: 第 189-202 行（Create 方法中）

**修改前**:
```go
photo := models.Photo{
    FilePath:      "/" + uploadPath,  // 结果: /./uploads/xxx.jpg ❌
    ThumbnailPath: "/" + thumbnailPath, // 结果: /./uploads/xxx_thumb.jpg ❌
    // ...
}
```

**修改后**:
```go
photo := models.Photo{
    FilePath:      "/" + strings.TrimPrefix(uploadPath, "./"),      // 结果: /uploads/xxx.jpg ✅
    ThumbnailPath: "/" + strings.TrimPrefix(thumbnailPath, "./"),   // 结果: /uploads/xxx_thumb.jpg ✅
    // ...
}
```

### 影响
- ✅ 新上传的照片将使用正确的路径格式
- ⚠️ 已存在的照片需要修复数据库记录（见下方修复脚本）

---

## 已存在照片的修复

对于已存在的照片，需要更新数据库中的路径：

```sql
-- 修复 file_path
UPDATE photos
SET file_path = REPLACE(file_path, '/./uploads/', '/uploads/')
WHERE file_path LIKE '/./uploads/%';

-- 修复 thumbnail_path
UPDATE photos
SET thumbnail_path = REPLACE(thumbnail_path, '/./uploads/', '/uploads/')
WHERE thumbnail_path LIKE '/./uploads/%';
```

### 执行修复脚本
```bash
# 使用 Go 程序修复
go run cmd/fix-paths/main.go
```

或手动执行 SQL（需要 sqlite3 工具）。

---

## 验证修复

### 1. 验证鼠标显示
- 访问后台管理页面：`http://localhost:5173/admin`
- 鼠标应该正常显示，不再隐藏

### 2. 验证缩略图显示
- 上传新照片
- 检查数据库路径格式是否正确
- 前端应该能正常显示缩略图

### 3. 检查 API 响应
```bash
curl http://localhost:9421/api/photos?pageSize=1 | grep -E "file_path|thumbnail_path"
```

期望输出：
```json
"file_path":"/uploads/xxx.jpg",
"thumbnail_path":"/uploads/xxx_thumb.jpg"
```

---

## 后续建议

1. **缩略图生成失败处理**
   - 当前如果缩略图生成失败，只是打印日志但不阻止上传
   - 建议添加重试机制或使用占位图

2. **路径标准化**
   - 建议创建一个工具函数统一处理路径格式
   - 避免在不同地方重复处理路径拼接

3. **前端容错**
   - 前端已经有 fallback 机制：`photo.thumbnail_path || photo.file_path`
   - 可以添加默认占位图作为最终 fallback

---

## 修改文件清单

1. ✅ `frontend/src/styles/main.css` - 移除鼠标隐藏样式
2. ✅ `backend/internal/handlers/photo.go` - 修复路径生成逻辑
3. ✅ `FIXES.md` - 本文档

---

## 测试清单

- [ ] 后台鼠标正常显示
- [ ] 上传新照片能正确生成缩略图
- [ ] 前端能正常显示缩略图
- [ ] 已存在照片需要修复路径（执行 SQL 脚本）
