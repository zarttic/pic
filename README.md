# 摄影作品集网站

一个基于 Vue 3 + Go 的全栈摄影作品展示与管理平台。

## 项目结构

```
.
├── frontend/          # Vue 3 前端项目
│   ├── src/
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── router/        # 路由
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── api/           # API 接口
│   │   └── styles/        # 样式
│   └── package.json
│
├── backend/           # Go 后端项目
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
├── photography-website-backup.html  # 原始单文件版本
└── README.md
```

## 技术栈

### 前端
- Vue 3 + Composition API
- Vite
- Vue Router
- Pinia
- Axios

### 后端
- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- SQLite (数据库)

## 快速开始

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:5173

### 后端开发

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

API 服务运行在 http://localhost:8080

## API 文档

### 照片管理

- `GET /api/photos` - 获取照片列表
  - 查询参数: `page`, `page_size`, `featured`, `tag`

- `GET /api/photos/:id` - 获取单张照片详情

- `POST /api/photos` - 上传新照片
  - Content-Type: `multipart/form-data`
  - 字段: `file`, `title`, `description`, `location`, `year`, `camera_model`, `lens`, `aperture`, `shutter_speed`, `iso`

- `PUT /api/photos/:id` - 更新照片信息
  - Content-Type: `application/json`

- `DELETE /api/photos/:id` - 删除照片

- `POST /api/photos/:id/view` - 增加浏览次数

### 文件上传

- `POST /api/upload` - 上传图片文件
  - Content-Type: `multipart/form-data`
  - 字段: `file`

## 开发进度

### v0.1.0 - 项目初始化 ✅
- [x] 初始化 Vue 3 + Vite 前端项目
- [x] 初始化 Go 后端项目
- [x] 配置项目结构
- [x] 实现基础路由
- [x] 配置 CORS
- [x] 实现数据库连接

### v0.2.0 - 照片管理功能（进行中）
- [x] 照片数据模型
- [x] 照片 CRUD API
- [x] 文件上传功能
- [x] 前端 Gallery 页面
- [x] 前端 Admin 管理页面
- [ ] 图片缩略图生成
- [ ] 照片编辑功能优化

### v0.3.0 - 相册功能（计划中）
- [ ] 相册数据模型
- [ ] 相册管理 API
- [ ] 前端相册页面
- [ ] 照片分类和标签

### v0.4.0 - 后台管理与认证（计划中）
- [ ] 用户认证系统
- [ ] JWT 中间件
- [ ] 权限控制
- [ ] 统计数据

### v1.0.0 - 部署上线（计划中）
- [ ] 生产环境配置
- [ ] PostgreSQL 数据库迁移
- [ ] 性能优化
- [ ] 部署文档

## 环境变量

### 前端 (.env)
```
VITE_API_URL=http://localhost:8080/api
```

### 后端 (.env)
```
SERVER_PORT=8080
DB_PATH=./photography.db
UPLOAD_PATH=./uploads
JWT_SECRET=your-secret-key-change-in-production
```

## 测试功能

1. 启动后端服务
2. 启动前端开发服务器
3. 访问前端首页
4. 导航到"管理"页面上传照片
5. 导航到"作品"页面查看照片

## 注意事项

- 首次运行后端时会自动创建 SQLite 数据库
- 上传的照片保存在 `backend/uploads/` 目录
- 开发环境使用 SQLite，生产环境建议切换到 PostgreSQL
- 生产环境需要更改 JWT_SECRET

## Git 工作流

当前使用简单的 master 分支进行迭代开发：

```
master
  └── v0.1.0 (项目初始化)
  └── v0.2.0 (照片管理)
  └── ...
```

提交规范：
- `feat:` 新功能
- `fix:` 修复 bug
- `refactor:` 重构代码
- `docs:` 文档更新
- `style:` 代码格式调整

## License

MIT
