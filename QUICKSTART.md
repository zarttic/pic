# 摄影作品集网站 - 快速开始指南

## 🚀 快速启动

### 方式一：使用启动脚本（Windows）
双击运行 `start-dev.bat`，会自动启动前后端服务器。

### 方式二：手动启动

#### 启动后端
```bash
cd backend
go run cmd/server/main.go
```
后端将运行在 http://localhost:9421

#### 启动前端
```bash
cd frontend
npm run dev
```
前端将运行在 http://localhost:5173

## 📋 功能清单

### ✅ 已完成（v0.1.0）
- [x] 项目初始化和基础架构
- [x] Vue 3 + Vite 前端配置
- [x] Go + Gin 后端配置
- [x] SQLite 数据库集成
- [x] 照片 CRUD API
- [x] 文件上传功能
- [x] 前端路由配置
- [x] Home 页面（首页展示）
- [x] Gallery 页面（照片展示）
- [x] About 页面（关于页面）
- [x] Admin 页面（照片管理）

### 🔄 进行中（v0.2.0）
- [ ] 照片编辑功能完善
- [ ] 图片缩略图生成
- [ ] 照片标签功能
- [ ] 分页优化

### 📅 计划中（v0.3.0+）
- [ ] 相册管理功能
- [ ] 用户认证系统
- [ ] 后台统计功能
- [ ] 生产环境部署

## 🎨 设计特色

- **现代画廊风格**：深色基调 + 暖金色点缀
- **自定义光标**：光晕跟随效果
- **流畅动画**：滚动渐入、悬停效果
- **响应式设计**：适配各种设备

## 📚 API 文档

### 照片管理

#### 获取照片列表
```
GET /api/photos?page=1&page_size=20
```

#### 获取单张照片
```
GET /api/photos/:id
```

#### 上传照片
```
POST /api/photos
Content-Type: multipart/form-data

字段：
- file: 图片文件（必需）
- title: 标题（必需）
- description: 描述
- location: 地点
- year: 年份
- camera_model: 相机型号
- lens: 镜头
- aperture: 光圈
- shutter_speed: 快门速度
- iso: ISO
```

#### 更新照片
```
PUT /api/photos/:id
Content-Type: application/json
```

#### 删除照片
```
DELETE /api/photos/:id
```

## 🧪 测试步骤

1. 启动后端服务器
2. 启动前端开发服务器
3. 访问 http://localhost:5173
4. 导航到"管理"页面
5. 上传一张照片（填写标题和选择文件）
6. 导航到"作品"页面查看上传的照片
7. 点击照片查看全屏预览

## 📂 项目结构

```
.
├── frontend/              # Vue 3 前端
│   ├── src/
│   │   ├── components/    # 可复用组件
│   │   ├── views/         # 页面组件
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── api/           # API 封装
│   │   └── styles/        # 全局样式
│   └── package.json
│
├── backend/               # Go 后端
│   ├── cmd/server/        # 入口文件
│   ├── internal/
│   │   ├── models/        # 数据模型
│   │   ├── handlers/      # HTTP 处理器
│   │   ├── middleware/    # 中间件
│   │   ├── services/      # 业务逻辑
│   │   └── config/        # 配置
│   ├── uploads/           # 上传文件存储
│   └── go.mod
│
└── start-dev.bat          # Windows 启动脚本
```

## ⚙️ 技术栈详情

### 前端
- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 下一代前端构建工具
- **Vue Router** - 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Axios** - HTTP 客户端

### 后端
- **Go 1.26** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 库
- **SQLite** - 嵌入式数据库（纯 Go 实现）
- **CORS** - 跨域资源共享

## 🔧 环境要求

- Node.js 18+
- Go 1.21+
- 现代浏览器（支持 ES6+）

## 📝 开发规范

### Git 提交格式
- `feat:` 新功能
- `fix:` 修复 bug
- `refactor:` 重构
- `docs:` 文档
- `style:` 格式
- `test:` 测试

### 分支策略
```
master (主分支)
  └── v0.1.0 (项目初始化)
  └── v0.2.0 (照片管理)
  └── ...
```

## 🐛 常见问题

### Q: 后端启动失败？
A: 确保端口 9421 未被占用，或修改 `backend/internal/config/config.go` 中的端口配置。

### Q: 前端无法连接后端？
A: 检查 `frontend/.env` 中的 `VITE_API_URL` 是否正确指向后端地址。

### Q: 上传照片失败？
A: 确保 `backend/uploads/` 目录存在且有写入权限。

### Q: 数据库在哪里？
A: SQLite 数据库文件位于 `backend/photography.db`。

## 📈 性能优化建议

1. **图片优化**：添加图片压缩和缩略图生成
2. **懒加载**：实现虚拟滚动优化大列表
3. **缓存**：添加 HTTP 缓存头
4. **CDN**：生产环境使用 CDN 加速静态资源

## 🔐 安全建议

生产环境部署时：
1. 更改 JWT_SECRET
2. 使用 HTTPS
3. 切换到 PostgreSQL
4. 添加用户认证
5. 实施文件上传限制

## 📄 License

MIT

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

---

**当前版本**: v0.1.0
**最后更新**: 2026-02-24
