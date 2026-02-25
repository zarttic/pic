# 快速测试指南

本指南帮助您快速测试新实现的功能。

## 前置准备

### 1. 启动后端服务

```bash
cd backend

# 设置环境变量（Windows PowerShell）
$env:JWT_SECRET="test-secret-key-for-development-only"
$env:ADMIN_USERNAME="admin"
$env:ADMIN_PASSWORD="admin123"

# 初始化管理员账户
go run cmd/init-admin/main.go

# 启动服务器
go run cmd/server/main.go
```

后端将在 `http://localhost:9421` 运行。

### 2. 启动前端服务

```bash
cd frontend

# 安装依赖（如果尚未安装）
npm install

# 启动开发服务器
npm run dev
```

前端将在 `http://localhost:5173` 运行。

## 功能测试

### 一、认证系统测试

#### 1.1 登录测试

**步骤：**
1. 访问 `http://localhost:5173/admin`
2. 应自动跳转到登录页面
3. 输入凭据：
   - 用户名：`admin`
   - 密码：`admin123`
4. 点击登录

**预期结果：**
- ✅ 登录成功，跳转到管理后台
- ✅ 浏览器开发者工具 → Application → Local Storage 中有 token 和 refreshToken

#### 1.2 Token 刷新测试

**步骤：**
1. 打开浏览器开发者工具 → Application → Local Storage
2. 手动修改 token（使其无效）
3. 刷新页面或访问需要认证的 API

**预期结果：**
- ✅ 系统自动使用 refreshToken 刷新 token
- ✅ 无需重新登录

#### 1.3 登出测试

**步骤：**
1. 在管理后台找到登出按钮
2. 点击登出

**预期结果：**
- ✅ Token 被清除
- ✅ 跳转到登录页面

### 二、批量操作测试

#### 2.1 批量删除照片

**使用 Postman 或 curl：**

```bash
# 先上传几张照片
curl -X POST http://localhost:9421/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test1.jpg" \
  -F "title=Test Photo 1"

curl -X POST http://localhost:9421/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test2.jpg" \
  -F "title=Test Photo 2"

# 批量删除
curl -X DELETE http://localhost:9421/api/photos/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ids": [1, 2]}'
```

**预期结果：**
- ✅ 返回 `{"message": "批量删除成功", "deleted": 2}`
- ✅ 照片记录和文件都被删除

#### 2.2 批量更新标签

```bash
curl -X PATCH http://localhost:9421/api/photos/batch/tags \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "ids": [3, 4, 5],
    "tags": "[\"风景\",\"旅行\"]"
  }'
```

**预期结果：**
- ✅ 指定照片的标签被更新

#### 2.3 批量设置精选

```bash
curl -X PATCH http://localhost:9421/api/photos/batch/featured \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "ids": [3, 4],
    "is_featured": true
  }'
```

**预期结果：**
- ✅ 指定照片被标记为精选

### 三、EXIF 自动提取测试

**准备：**
需要一张包含 EXIF 信息的照片（手机拍摄的照片通常都有）

**步骤：**
1. 访问管理后台上传照片
2. 选择一张照片上传（不要手动填写相机参数）
3. 上传成功后，查看照片详情

**预期结果：**
- ✅ 相机型号自动填充
- ✅ 镜头信息自动填充（如果有）
- ✅ 光圈、快门、ISO 自动填充
- ✅ 拍摄日期自动填充

**验证方式：**
- 右键照片 → 查看详情
- 或使用在线 EXIF 查看器对比

### 四、搜索功能测试

#### 4.1 关键词搜索

**使用 API：**

```bash
# 搜索标题或描述包含"风景"的照片
curl "http://localhost:9421/api/photos?search=风景"
```

**前端测试：**
1. 在画廊页面使用搜索框
2. 输入关键词（如照片标题、地点、标签）
3. 按回车或点击搜索按钮

**预期结果：**
- ✅ 返回匹配的照片列表
- ✅ 支持中文搜索
- ✅ 搜索标题、描述、地点、标签

#### 4.2 组合筛选

```bash
# 搜索 2024 年用 Canon 拍摄的精选照片
curl "http://localhost:9421/api/photos?year=2024&camera=Canon&featured=true"
```

**前端测试：**
1. 使用 SearchBar 组件的筛选器
2. 选择年份
3. 选择相机
4. 勾选"仅显示精选"

**预期结果：**
- ✅ 返回符合所有条件的照片
- ✅ 分页正常工作

## 安全测试

### 1. 未认证访问测试

**步骤：**
```bash
# 尝试无 Token 访问管理接口
curl -X POST http://localhost:9421/api/photos \
  -H "Content-Type: application/json" \
  -d '{"title": "Unauthorized"}'
```

**预期结果：**
- ✅ 返回 401 Unauthorized
- ✅ 错误信息：`{"error": "未提供认证令牌"}`

### 2. 无效 Token 测试

**步骤：**
```bash
curl -X GET http://localhost:9421/api/me \
  -H "Authorization: Bearer invalid-token"
```

**预期结果：**
- ✅ 返回 401 Unauthorized
- ✅ 错误信息：`{"error": "无效或过期的认证令牌"}`

### 3. 文件类型验证测试

**步骤：**
```bash
# 尝试上传不允许的文件类型
curl -X POST http://localhost:9421/api/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@test.pdf"
```

**预期结果：**
- ✅ 返回 400 Bad Request
- ✅ 错误信息：`{"error": "不支持的文件类型，仅支持 jpg, jpeg, png, webp"}`

### 4. 文件大小限制测试

**步骤：**
尝试上传超过 10MB 的照片

**预期结果：**
- ✅ 返回 400 Bad Request
- ✅ 错误信息：`{"error": "文件大小超过限制 (最大 10MB)"}`

## 性能测试

### 1. 批量操作性能

**测试批量删除 100 张照片：**

```bash
# 准备 ID 列表
curl -X DELETE http://localhost:9421/api/photos/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"ids": [1,2,3,...,100]}'
```

**预期结果：**
- ✅ 操作在合理时间内完成
- ✅ 返回正确的删除数量

### 2. 搜索性能

**测试搜索响应时间：**

```bash
# 测试大量数据的搜索
time curl "http://localhost:9420/api/photos?search=test&year=2024"
```

**预期结果：**
- ✅ 响应时间 < 500ms（取决于数据量）

## 故障排查

### 问题 1：无法登录

**检查：**
- 后端服务是否正常运行
- 数据库是否已初始化
- 管理员账户是否已创建
- JWT_SECRET 是否设置

**解决：**
```bash
# 重新初始化管理员
go run cmd/init-admin/main.go
```

### 问题 2：Token 刷新失败

**检查：**
- refreshToken 是否有效
- 后端日志中的错误信息

**解决：**
- 清除 LocalStorage 重新登录

### 问题 3：EXIF 未提取

**检查：**
- 照片是否包含 EXIF 信息
- 照片格式是否支持

**解决：**
- 使用手机拍摄的照片测试
- 检查照片是否被压缩（微信等会删除 EXIF）

### 问题 4：批量操作失败

**检查：**
- Token 是否有效
- ID 数量是否超过 100
- 照片是否存在

## 日志查看

### 后端日志

后端服务会在控制台输出日志，包括：
- 请求路径和方法
- 错误信息
- 文件操作结果

### 前端日志

打开浏览器开发者工具 → Console 查看：
- API 请求错误
- 认证状态变化
- 组件错误

## 测试检查清单

### 认证系统
- [ ] 登录成功
- [ ] 登录失败提示正确
- [ ] Token 自动刷新
- [ ] 登出清除 Token
- [ ] 路由守卫工作正常

### 批量操作
- [ ] 批量删除成功
- [ ] 批量更新标签成功
- [ ] 批量设置精选成功
- [ ] 超过 100 张限制

### EXIF 提取
- [ ] 自动提取相机信息
- [ ] 自动提取拍摄参数
- [ ] 无 EXIF 不报错

### 搜索功能
- [ ] 关键词搜索
- [ ] 筛选器组合
- [ ] 分页正常

### 安全验证
- [ ] 未认证访问被拒绝
- [ ] 无效 Token 被拒绝
- [ ] 文件类型验证
- [ ] 文件大小验证

## 下一步

测试完成后，可以：
1. 开始前端界面的批量操作 UI 开发
2. 添加更多筛选维度
3. 实现照片-相册关联 UI
4. 开始 Phase 3 的性能优化

## 反馈问题

如发现问题，请记录：
- 问题描述
- 复现步骤
- 预期结果 vs 实际结果
- 错误日志
- 环境信息（浏览器、操作系统）
