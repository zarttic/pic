# 相册密码验证Bug修复

## 修复的两个Bug

### Bug 1: 密码错误会跳转到管理登录页面
**问题原因:**
- API 响应拦截器对所有 401 错误都会跳转到 `/admin/login`
- 相册密码验证失败返回 401,但不应该跳转到管理登录页

**修复方案:**
在 `frontend/src/api/index.js` 的响应拦截器中添加特殊判断:
```javascript
// 检查是否是相册密码验证错误(不跳转登录页)
const isAlbumVerifyError = originalRequest.url?.includes('/albums/') && originalRequest.url?.includes('/verify')

switch (error.response.status) {
  case 401:
    // 如果是相册密码验证失败,不跳转登录页
    if (isAlbumVerifyError) {
      break
    }
    // ... 其他 401 处理逻辑
}
```

### Bug 2: 输入正确密码后原地不动
**问题原因:**
1. `closePasswordDialog()` 函数会检查 `require_auth`,如果为 true 就跳转回相册列表
2. 但密码验证成功后,数据还没有及时更新,导致判断错误

**修复方案:**
在 `frontend/src/views/AlbumDetail.vue` 中优化密码验证成功后的逻辑:

```javascript
const handleVerifyPassword = async () => {
  if (!passwordInput.value) return

  verifying.value = true
  passwordError.value = ''

  try {
    await albumStore.verifyPassword(route.params.id, passwordInput.value)
    // 验证成功，重新获取相册数据
    const token = localStorage.getItem(`album_token_${route.params.id}`)
    await fetchAlbumWithToken(route.params.id, token)
    // 成功获取数据后关闭对话框(不调用 closePasswordDialog,避免跳转)
    showPasswordDialog.value = false
    passwordInput.value = ''
    passwordError.value = ''
  } catch (error) {
    passwordError.value = '密码错误，请重试'
  } finally {
    verifying.value = false
  }
}
```

同时优化 `fetchAlbumWithToken` 函数,确保正确更新数据:
```javascript
const fetchAlbumWithToken = async (albumId, token) => {
  try {
    const response = await api.get(`/albums/${albumId}`, {
      headers: {
        'X-Album-Token': token
      }
    })
    const data = response.data
    // 如果返回的数据不包含 require_auth 或 require_auth 为 false,说明验证成功
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

## 修改的文件

1. **frontend/src/api/index.js**
   - 添加相册密码验证错误的特殊处理
   - 防止 401 错误时跳转到管理登录页

2. **frontend/src/views/AlbumDetail.vue**
   - 优化密码验证成功后的逻辑
   - 直接关闭对话框而不调用 `closePasswordDialog`
   - 改进 `fetchAlbumWithToken` 函数的返回值和错误处理

## 测试步骤

### 测试场景 1: 密码错误不跳转
1. 前台访问密码保护的相册
2. 输入错误的密码
3. 点击"确认"
4. ✅ 应该显示"密码错误,请重试"
5. ✅ **不应该**跳转到管理登录页面

### 测试场景 2: 密码正确成功进入
1. 在密码对话框中输入正确的密码
2. 点击"确认"
3. ✅ 对话框应该关闭
4. ✅ 应该显示相册中的照片
5. ✅ **不应该**停留在原地或跳转到其他页面

### 测试场景 3: Token 持久化
1. 成功输入密码后,刷新页面
2. ✅ 应该直接显示照片,无需再次输入密码
3. 点击浏览器返回,再重新进入该相册
4. ✅ 应该直接显示照片

### 测试场景 4: 取消输入
1. 点击密码对话框的"取消"按钮
2. ✅ 应该返回相册列表
3. 再次点击该相册
4. ✅ 应该重新弹出密码对话框

## 技术要点

### 1. 区分不同类型的 401 错误
- **管理认证失败**: JWT token 无效或过期 → 跳转登录页
- **相册密码验证失败**: 密码错误 → 显示错误提示,不跳转

### 2. 正确的状态管理
- 密码验证成功后立即更新 `albumStore.currentAlbum`
- 关闭对话框后不需要额外判断 `require_auth`
- 避免在数据更新前进行跳转判断

### 3. 用户体验优化
- 错误提示清晰明确
- 对话框关闭逻辑简洁
- Token 自动保存和重用

## 预期结果

修复后,用户应该能够:
1. ✅ 输入错误密码时看到错误提示,不会被跳转
2. ✅ 输入正确密码后立即看到相册照片
3. ✅ Token 持久化提供流畅的访问体验
4. ✅ 取消输入时正常返回相册列表
