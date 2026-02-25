package main

import (
	"log"
	"picsite/internal/config"
	"picsite/internal/handlers"
	"picsite/internal/middleware"
	"picsite/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 验证 JWT Secret
	if cfg.JWTSecret == "your-secret-key-change-in-production" {
		log.Println("警告: JWT_SECRET 未设置，使用了默认值！生产环境必须设置 JWT_SECRET 环境变量")
	}

	// 初始化数据库
	if err := services.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API 路由
	api := r.Group("/api")
	{
		photoHandler := handlers.NewPhotoHandler()
		albumHandler := handlers.NewAlbumHandler()
		authHandler := handlers.NewAuthHandler(cfg)

		// 认证路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// 照片相关路由（公开）
		photos := api.Group("/photos")
		{
			photos.GET("", photoHandler.GetAll)
			photos.GET("/:id", photoHandler.GetByID)
			photos.POST("/:id/view", photoHandler.IncrementView)
		}

		// 照片管理路由（需要认证）
		photosAdmin := api.Group("/photos")
		photosAdmin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			photosAdmin.POST("", photoHandler.Create)
			photosAdmin.PUT("/:id", photoHandler.Update)
			photosAdmin.DELETE("/:id", photoHandler.Delete)
			photosAdmin.DELETE("/batch", photoHandler.BatchDelete)
			photosAdmin.PATCH("/batch/tags", photoHandler.BatchUpdateTags)
			photosAdmin.PATCH("/batch/featured", photoHandler.BatchUpdateFeatured)
		}

		// 相册相关路由（公开）
		albums := api.Group("/albums")
		{
			albums.GET("", albumHandler.GetAll)
			albums.GET("/:id", albumHandler.GetByID)
			albums.POST("/:id/verify", albumHandler.VerifyPassword)
		}

		// 相册管理路由（需要认证）
		albumsAdmin := api.Group("/albums")
		albumsAdmin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			albumsAdmin.POST("", albumHandler.Create)
			albumsAdmin.PUT("/:id", albumHandler.Update)
			albumsAdmin.DELETE("/:id", albumHandler.Delete)
			albumsAdmin.POST("/:id/photos", albumHandler.AddPhotoToAlbum)
			albumsAdmin.DELETE("/:id/photos/:photo_id", albumHandler.RemovePhotoFromAlbum)
			albumsAdmin.POST("/:id/password", albumHandler.SetPassword)
			albumsAdmin.DELETE("/:id/password", albumHandler.RemovePassword)
		}

		// 文件上传（需要认证）
		api.POST("/upload", middleware.AuthMiddleware(cfg.JWTSecret), photoHandler.UploadFile)

		// 当前用户信息（需要认证）
		api.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetCurrentUser)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 启动服务器
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
