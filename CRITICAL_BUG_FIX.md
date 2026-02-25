# 相册密码验证核心Bug修复

## 问题描述

用户报告前台访问密码保护相册时，输入正确密码后仍然无法进入相册查看照片。

## 根本原因分析

### 问题1: JWT Token 干扰
**现象**: 前台用户访问密码保护相册时，后端误判为管理员访问。

**根本原因**:
1. 用户之前登录过管理后台，localStorage 中保存了 JWT token
2. 前端 axios 请求拦截器会自动为所有请求添加 `Authorization: Bearer <JWT>`
3. 后端的判断逻辑:
   ```go
   // 检查是否是管理员访问(通过JWT token)
   isAdmin := false
   authHeader := c.GetHeader("Authorization")
   if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
       isAdmin = true  // ❌ 这里有 JWT 就认为是管理员
   }
   ```
4. 结果: 后端认为有 JWT token 就是管理员，直接返回完整数据，跳过了密码验证

### 问题2: 密码验证后仍显示 require_auth
**现象**: 密码验证成功后，再次请求相册数据时仍然返回 `require_auth: true`。

**根本原因**:
1. 密码验证成功后，前端使用 `api.get()` 再次请求相册数据
2. `api.get()` 会自动添加 JWT token（问题1）
3. 后端检测到 JWT token，认为用户是管理员
4. 但是管理员访问时，后端会检查 `IsProtected` 字段
5. 如果相册有密码保护且用户有 JWT，后端返回完整数据（管理员权限）
6. 但此时逻辑有问题...

**等等，让我重新理解后端逻辑**:

```go
if album.IsProtected && !isAdmin {
    // 需要验证 album token
    session := middleware.SessionManagerInstance.GetSession(token)
    if session == nil || session.AlbumID != album.ID {
        // 返回 require_auth
    }
}
// 返回完整数据
```

**真正的问题**:
- 如果 `isAdmin = true`（有 JWT），则跳过整个验证逻辑
- 直接返回完整数据，包括照片
- **但是前台用户也有 JWT token（之前登录过管理后台）！**

## 解决方案

### 方案设计
区分前台访问和后台访问：
- **后台访问**: 使用 `api.get()`，自动添加 JWT token
- **前台访问**: 使用 `fetch()` 或专用方法，不添加 JWT token

### 实现细节

#### 1. 修改 albums store
**文件**: `frontend/src/stores/albums.js`

新增 `fetchAlbumPublic` 方法，专门用于前台访问:

```javascript
async function fetchAlbumPublic(id) {
  loading.value = true
  error.value = null
  try {
    // 使用原生 fetch，不发送 JWT token
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
    const response = await fetch(`${apiUrl}/albums/${id}`)

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to fetch album')
    }

    const data = await response.json()
    currentAlbum.value = data
    return data
  } catch (err) {
    error.value = err.message
    handleError(err, '获取相册详情')
    throw err
  } finally {
    loading.value = false
  }
}
```

#### 2. 修改 AlbumDetail.vue
**文件**: `frontend/src/views/AlbumDetail.vue`

**修改点1**: 使用 `fetchAlbumPublic` 替代 `fetchAlbum`:
```javascript
onMounted(async () => {
  const albumId = route.params.id
  // 使用 fetchAlbumPublic 避免发送 JWT token
  await albumStore.fetchAlbumPublic(albumId)

  if (albumStore.currentAlbum?.require_auth) {
    const token = localStorage.getItem(`album_token_${albumId}`)
    if (token) {
      await fetchAlbumWithToken(albumId, token)
    } else {
      showPasswordDialog.value = true
    }
  }
})
```

**修改点2**: `fetchAlbumWithToken` 使用原生 fetch:
```javascript
const fetchAlbumWithToken = async (albumId, token) => {
  try {
    // 前台访问相册时，使用原生 fetch，不使用 JWT token
    const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080/api'}/albums/${albumId}`, {
      headers: {
        'X-Album-Token': token
      }
    })
    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.error || 'Failed to fetch album')
    }

    if (!data.require_auth) {
      albumStore.currentAlbum = data
      showPasswordDialog.value = false
      return true
    } else {
      showPasswordDialog.value = true
      return false
    }
  } catch (error) {
    console.error('Error fetching album with token:', error)
    showPasswordDialog.value = true
    return false
  }
}
```

## 完整流程

### 前台访问密码保护相册的流程

1. **用户访问相册页面**
   ```
   GET /albums/:id
   Headers: (无 Authorization)
   ```

2. **后端判断**
   ```
   - album.IsProtected = true
   - isAdmin = false (没有 JWT)
   - 没有 album_token
   → 返回 { require_auth: true, ... }
   ```

3. **前端显示密码对话框**

4. **用户输入密码**
   ```
   POST /albums/:id/verify
   Body: { password: "xxx" }
   ```

5. **后端验证密码**
   ```
   - 密码正确
   - 生成 session token
   - 返回 { token: "xxx" }
   ```

6. **前端保存 token**
   ```
   localStorage.setItem('album_token_' + id, token)
   ```

7. **前端再次请求相册数据**
   ```
   GET /albums/:id
   Headers: {
     'X-Album-Token': token
     (注意: 没有 Authorization header)
   }
   ```

8. **后端验证 album token**
   ```
   - album.IsProtected = true
   - isAdmin = false (没有 JWT)
   - 有 album_token
   - session.AlbumID = album.ID
   → 返回完整数据（包括照片）
   ```

9. **前端显示照片** ✅

### 管理员访问密码保护相册的流程

1. **管理员访问后台**
   ```
   GET /admin/albums/:id/photos
   Headers: {
     Authorization: Bearer <JWT>
   }
   ```

2. **后端判断**
   ```
   - album.IsProtected = true
   - isAdmin = true (有 JWT)
   → 返回完整数据（跳过密码验证）
   ```

3. **前端显示照片** ✅

## 关键改进

### 1. 前后台访问分离
- **后台**: 使用 `api.get()` → 自动添加 JWT
- **前台**: 使用 `fetch()` 或 `fetchAlbumPublic()` → 不添加 JWT

### 2. 清晰的权限模型
- **管理员**: JWT token → 完全访问权限
- **普通用户**: 无 JWT + album token → 访问特定相册
- **未认证**: 无 JWT + 无 album token → require_auth

### 3. Token 类型
- **JWT Token**: 用于管理后台认证，存储在 `localStorage.token`
- **Album Token**: 用于相册密码验证，存储在 `localStorage.album_token_{id}`

## 测试场景

### 场景1: 普通用户从未登录过管理后台
1. 访问密码保护相册
2. 输入正确密码
3. ✅ 应该能看到照片

### 场景2: 普通用户之前登录过管理后台（关键场景）
1. localStorage 中有 JWT token
2. 访问密码保护相册
3. 输入正确密码
4. ✅ 应该能看到照片（不会因为 JWT 被误判）

### 场景3: 管理员访问密码保护相册
1. 管理后台访问
2. ✅ 应该直接看到照片，无需密码

### 场景4: 密码错误
1. 输入错误密码
2. ✅ 显示错误提示
3. ✅ 不跳转到管理登录页

## 文件修改总结

### 修改的文件
1. `frontend/src/stores/albums.js` - 新增 `fetchAlbumPublic` 方法
2. `frontend/src/views/AlbumDetail.vue` - 使用 `fetchAlbumPublic` 和原生 `fetch`

### 核心原则
**前台访问相册时，绝对不能发送 JWT token！**

## 技术要点

### 1. 为什么不能用 axios 拦截器跳过 JWT？
```javascript
// ❌ 这样不行
api.interceptors.request.use(config => {
  if (config.url.includes('/albums/')) {
    delete config.headers.Authorization
  }
})
```
因为：
- 管理后台访问相册时也需要 JWT
- 无法区分前台访问和后台访问

### 2. 为什么后台访问相册需要 JWT？
管理员应该能够：
- 查看所有相册（包括密码保护的）
- 管理相册中的照片
- 设置封面、密码等

### 3. 为什么前台访问相册不能有 JWT？
因为会导致后端误判，认为用户是管理员，跳过密码验证。但实际上：
- 用户可能只是之前登录过管理后台
- 现在是以普通用户身份访问前台
- 应该遵守密码保护相册的验证规则

## 总结

这个Bug的根本原因是**JWT token 的自动注入干扰了相册密码验证流程**。

解决方案是通过**前后台访问分离**，确保：
- 前台访问时不发送 JWT → 正常进行密码验证
- 后台访问时发送 JWT → 获得管理员权限

这样既保证了普通用户的密码验证流程，又保证了管理员的便捷操作。
