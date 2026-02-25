# PicSite - 摄影作品展示网站

[![Test](https://github.com/zarttic/pic/workflows/Test/badge.svg)](https://github.com/zarttic/pic/actions/workflows/test.yml)
[![Deploy](https://github.com/zarttic/pic/workflows/Deploy/badge.svg)](https://github.com/zarttic/pic/actions/workflows/deploy.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个现代化、功能完善的摄影作品展示网站，基于 Vue 3 + Go 构建，支持照片管理、相册管理、加密相册、访问统计等功能。

![PicSite Screenshot](docs/screenshot.png)

## ✨ 特性

### 核心功能
- 📸 **照片管理** - 上传、编辑、删除、批量操作
- 📁 **相册管理** - 创建、编辑、加密相册
- 🔐 **加密相册** - 密码保护的私密相册
- 📊 **访问统计** - 照片浏览次数统计
- 🔍 **智能搜索** - 多字段搜索和筛选
- 🎯 **EXIF 提取** - 自动提取相机参数

### 安全特性
- 🔒 **JWT 认证** - 安全的管理员认证系统
- 🛡️ **bcrypt 加密** - 密码使用 bcrypt 哈希
- ✅ **文件验证** - 类型和大小的严格验证
- 🎫 **随机 Token** - 加密安全的会话管理

### 用户体验
- ⚡ **懒加载** - 图片按需加载，提升性能
- 🎨 **骨架屏** - 优雅的加载状态
- 🔔 **Toast 通知** - 实时的操作反馈
- 📱 **响应式设计** - 完美适配移动端
- 🎭 **错误边界** - 友好的错误处理

### 技术栈

#### 后端
- **框架**: Go 1.24 + Gin
- **ORM**: GORM
- **数据库**: SQLite（可迁移至 PostgreSQL）
- **认证**: JWT (HS256)
- **加密**: bcrypt

#### 前端
- **框架**: Vue 3 + Vite
- **状态管理**: Pinia
- **路由**: Vue Router
- **样式**: CSS Variables
- **字体**: Cormorant Garamond + Outfit

## 🚀 快速开始

### 方式一：Docker 部署（推荐）

**前置要求：**
- Docker 20.10+
- Docker Compose 2.0+

**部署步骤：**

```bash
# 1. 克隆仓库
git clone https://github.com/zarttic/pic.git
cd pic

# 2. 配置环境变量
cp backend/.env.example backend/.env
# 编辑 backend/.env，设置 JWT_SECRET

# 3. 启动服务
docker-compose up -d

# 4. 初始化管理员
docker-compose exec backend ./init-admin

# 5. 访问应用
# 前端: http://localhost
# 后端: http://localhost:9421
# 管理后台: http://localhost/admin
```

### 方式二：本地开发

**后端：**

```bash
cd backend

# 安装依赖
go mod download

# 配置环境
cp .env.example .env
# 编辑 .env，设置 JWT_SECRET

# 初始化管理员
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your-password
go run cmd/init-admin/main.go

# 启动服务
go run cmd/server/main.go
```

**前端：**

```bash
cd frontend

# 安装依赖
npm install

# 配置环境
cp .env.example .env
# 编辑 .env，设置 API 地址

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 📖 文档

- **[BUG_FIXES.md](./BUG_FIXES.md)** - Bug 修复记录（汇总所有问题修复）
- **[PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)** - 项目迭代总结
- **[TESTING_GUIDE.md](./TESTING_GUIDE.md)** - 测试指南
- **[CHANGELOG.md](./CHANGELOG.md)** - 版本更新日志
- [后端开发文档](./backend/README.md)

## 🗂️ 项目结构

```
.
├── backend/                 # Go 后端
│   ├── cmd/                # 命令行工具
│   │   ├── server/        # 主程序
│   │   └── init-admin/    # 管理员初始化
│   ├── internal/
│   │   ├── config/        # 配置管理
│   │   ├── handlers/      # HTTP 处理器
│   │   ├── middleware/    # 中间件
│   │   ├── models/        # 数据模型
│   │   ├── services/      # 业务逻辑
│   │   └── utils/         # 工具函数
│   ├── Dockerfile
│   └── README.md
│
├── frontend/               # Vue 前端
│   ├── src/
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── stores/        # 状态管理
│   │   ├── router/        # 路由配置
│   │   ├── api/           # API 客户端
│   │   └── directives/    # 自定义指令
│   ├── Dockerfile
│   └── nginx.conf
│
├── .github/
│   └── workflows/         # CI/CD 配置
│       ├── test.yml
│       └── deploy.yml
│
├── docker-compose.yml     # Docker 编排
└── README.md             # 本文件
```

## 🔧 配置

### 后端环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | 9421 | 服务器端口 |
| `DB_PATH` | ./picsite.db | 数据库路径 |
| `UPLOAD_PATH` | ./uploads | 上传目录 |
| `JWT_SECRET` | *必填* | JWT 密钥 |

### 前端环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `VITE_API_URL` | http://localhost:9421/api | API 地址 |

## 📚 API 文档

### 公开接口

#### 照片
- `GET /api/photos` - 获取照片列表
- `GET /api/photos/:id` - 获取单张照片
- `POST /api/photos/:id/view` - 增加浏览次数

#### 相册
- `GET /api/albums` - 获取相册列表
- `GET /api/albums/:id` - 获取相册详情
- `POST /api/albums/:id/verify` - 验证相册密码

### 认证接口
- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `POST /api/auth/refresh` - 刷新令牌
- `GET /api/me` - 获取当前用户

### 管理接口（需认证）

#### 照片管理
- `POST /api/photos` - 创建照片
- `PUT /api/photos/:id` - 更新照片
- `DELETE /api/photos/:id` - 删除照片
- `DELETE /api/photos/batch` - 批量删除
- `PATCH /api/photos/batch/tags` - 批量更新标签
- `PATCH /api/photos/batch/featured` - 批量设置精选

#### 相册管理
- `POST /api/albums` - 创建相册
- `PUT /api/albums/:id` - 更新相册
- `DELETE /api/albums/:id` - 删除相册
- `POST /api/albums/:id/photos` - 添加照片到相册
- `DELETE /api/albums/:id/photos/:photo_id` - 移除照片
- `POST /api/albums/:id/password` - 设置密码

完整 API 文档请查看 [API.md](./docs/API.md)

## 🧪 测试

### 后端测试

```bash
cd backend
go test ./... -v
```

### 前端测试

```bash
cd frontend
npm run test
```

## 🚢 部署

### 生产环境部署建议

1. **使用 HTTPS**
   - 配置 SSL 证书
   - 强制 HTTPS 重定向

2. **数据库优化**
   - 迁移到 PostgreSQL
   - 定期备份

3. **文件存储**
   - 使用对象存储（AWS S3、阿里云 OSS）
   - 配置 CDN 加速

4. **性能优化**
   - 启用 Gzip 压缩
   - 配置浏览器缓存
   - 使用 Redis 缓存会话

5. **安全加固**
   - 设置防火墙规则
   - 定期更新依赖
   - 启用日志审计

### Docker 生产部署

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  backend:
    image: your-registry/picsite-backend:latest
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - DB_PATH=/app/data/picsite.db
    volumes:
      - ./data:/app/data
      - ./uploads:/app/uploads
    restart: always

  frontend:
    image: your-registry/picsite-frontend:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    restart: always
```

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md)

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📝 版本历史

### v0.6.0 (当前)
- ✅ Phase 4: 功能增强与代码优化
- 批量操作（选择、删除、标签、精选）
- 实时搜索和多条件筛选
- 标签管理系统
- 图片预览和上传进度组件
- Gallery 页面性能优化

### v0.5.0
- ✅ Phase 3: 用户体验优化
- 骨架屏加载
- 图片懒加载
- Toast 通知
- 错误边界
- 移动端响应式

### v0.4.0
- ✅ Phase 1: 安全加固与认证系统
- ✅ Phase 2: 功能完善
- JWT 认证
- 批量操作
- EXIF 提取
- 搜索功能

### v0.3.0
- 相册功能
- 加密相册
- 访问统计

查看 [CHANGELOG.md](./CHANGELOG.md) 获取完整历史

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Vue.js](https://vuejs.org/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- 所有贡献者

## 📮 联系方式

项目地址: [https://github.com/zarttic/pic](https://github.com/zarttic/pic)

问题反馈: [Issues](https://github.com/zarttic/pic/issues)

---

**⭐ 如果这个项目对你有帮助，请给一个 Star！**
