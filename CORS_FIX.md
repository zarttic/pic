# CORS 问题修复

## 问题描述
前端使用 `X-Album-Token` header 发送请求时，被 CORS 策略阻止：
```
Request header field x-album-token is not allowed by Access-Control-Allow-Headers in preflight response.
```

## 解决方案
在后端 CORS 配置中添加 `X-Album-Token` 到允许的 headers 列表。

## 修改文件
`backend/internal/middleware/cors.go`

### 修改前
```go
AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
```

### 修改后
```go
AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Album-Token"},
```

## 重启服务器

修改 CORS 配置后，需要重启后端服务器：

```bash
# 停止当前运行的后端服务器 (Ctrl+C)

# 重新启动
cd backend
go run cmd/server/main.go
```

## 验证修复

重启后端后，在前台测试：
1. 访问密码保护的相册
2. 输入正确密码
3. ✅ 应该能够成功看到照片

## 完整的 CORS 配置

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

## 允许的自定义 Headers
- `Authorization` - JWT token (管理后台认证)
- `X-Album-Token` - Album token (相册密码验证)
- `Content-Type` - 请求体类型
- `Origin` - 请求来源

## 注意事项
1. CORS 配置修改后必须重启服务器
2. 浏览器可能会缓存 CORS preflight 响应，可以清除浏览器缓存
3. 开发环境允许 localhost，生产环境需要配置实际域名
