# 摄影网站后端

基于 Go + Gin + GORM + SQLite 构建的摄影作品展示网站后端服务。

## 功能特性

- ✅ 照片管理（上传、编辑、删除、浏览）
- ✅ 相册管理（创建、编辑、删除、加密）
- ✅ 管理员认证系统（JWT）
- ✅ 文件上传验证
- ✅ 密码加密（bcrypt）
- ✅ 访问统计

## 技术栈

- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite
- **认证**: JWT (HS256)
- **密码加密**: bcrypt

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

复制环境变量模板：

```bash
cp .env.example .env
```

编辑 `.env` 文件，**必须修改 `JWT_SECRET`**：

```bash
# 生成安全的密钥
openssl rand -base64 32
```

### 3. 初始化管理员账户

```bash
# 设置管理员信息
export ADMIN_USERNAME=admin
export ADMIN_PASSWORD=your-secure-password
export ADMIN_EMAIL=admin@example.com

# 运行初始化脚本
go run cmd/init-admin/main.go
```

### 4. 启动服务器

```bash
go run cmd/server/main.go
```

服务器将在 `http://localhost:9421` 启动。

## API 接口

### 公开接口

#### 照片

- `GET /api/photos` - 获取照片列表
- `GET /api/photos/:id` - 获取单张照片
- `POST /api/photos/:id/view` - 增加浏览次数

#### 相册

- `GET /api/albums` - 获取相册列表
- `GET /api/albums/:id` - 获取单个相册
- `POST /api/albums/:id/verify` - 验证相册密码

### 认证接口

- `POST /api/auth/login` - 登录
- `POST /api/auth/logout` - 登出
- `POST /api/auth/refresh` - 刷新令牌

### 管理接口（需要认证）

所有管理接口需要在请求头中携带 `Authorization: Bearer <token>`。

#### 照片管理

- `POST /api/photos` - 创建照片
- `PUT /api/photos/:id` - 更新照片
- `DELETE /api/photos/:id` - 删除照片
- `POST /api/upload` - 上传文件

#### 相册管理

- `POST /api/albums` - 创建相册
- `PUT /api/albums/:id` - 更新相册
- `DELETE /api/albums/:id` - 删除相册
- `POST /api/albums/:id/photos` - 添加照片到相册
- `DELETE /api/albums/:id/photos/:photo_id` - 从相册移除照片
- `POST /api/albums/:id/password` - 设置相册密码
- `DELETE /api/albums/:id/password` - 移除相册密码

#### 用户信息

- `GET /api/me` - 获取当前用户信息

## 安全特性

### 密码加密

- 使用 bcrypt 进行密码哈希
- 自动处理 salt
- 相册密码同样使用 bcrypt 加密

### JWT 认证

- Token 有效期：24小时
- Refresh Token 有效期：7天
- 使用 HS256 签名算法
- **生产环境必须设置强随机 JWT_SECRET**

### 会话管理

- 相册访问使用随机生成的 UUID Token
- 会话有效期：24小时
- 内存存储（生产环境建议使用 Redis）

### 文件上传验证

- 允许的文件类型：jpg, jpeg, png, webp
- 文件大小限制：10MB
- 扩展名白名单验证

## 项目结构

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go          # 主程序入口
│   └── init-admin/
│       └── main.go          # 初始化管理员脚本
├── internal/
│   ├── config/
│   │   └── config.go        # 配置管理
│   ├── handlers/
│   │   ├── auth.go          # 认证处理器
│   │   ├── photo.go         # 照片处理器
│   │   └── album.go         # 相册处理器
│   ├── middleware/
│   │   ├── auth.go          # JWT 认证中间件
│   │   ├── album_auth.go    # 相册访问认证
│   │   └── cors.go          # CORS 中间件
│   ├── models/
│   │   └── models.go        # 数据模型
│   ├── services/
│   │   ├── database.go      # 数据库服务
│   │   └── thumbnail.go     # 缩略图服务
│   └── utils/
│       ├── password.go      # 密码工具
│       └── jwt.go           # JWT 工具
├── uploads/                  # 上传文件存储
├── .env.example              # 环境变量模板
├── go.mod
└── go.sum
```

## 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| SERVER_PORT | 9421 | 服务器端口 |
| DB_PATH | ./picsite.db | 数据库文件路径 |
| UPLOAD_PATH | ./uploads | 上传文件存储路径 |
| JWT_SECRET | *需设置* | JWT 签名密钥（**生产环境必须修改**）|
| ADMIN_USERNAME | admin | 初始管理员用户名 |
| ADMIN_PASSWORD | admin123 | 初始管理员密码 |
| ADMIN_EMAIL | admin@example.com | 初始管理员邮箱 |

## 开发指南

### 运行测试

```bash
go test ./...
```

### 构建生产版本

```bash
go build -o picsite-server cmd/server/main.go
```

### 数据库迁移

GORM 会自动迁移，无需手动操作。

## 注意事项

⚠️ **生产环境安全提示**：

1. 必须设置强随机的 `JWT_SECRET`
2. 修改默认管理员密码
3. 使用 HTTPS
4. 配置防火墙规则
5. 定期备份数据库
6. 考虑使用 Redis 存储会话
7. 添加速率限制

## 迁移指南

### 从 SHA256 密码迁移

如果您之前的相册密码使用 SHA256 加密，需要进行迁移：

1. 所有相册密码需要重新设置
2. 用户密码会自动使用 bcrypt

## License

MIT
