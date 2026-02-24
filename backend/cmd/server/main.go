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

		// 照片相关路由
		photos := api.Group("/photos")
		{
			photos.GET("", photoHandler.GetAll)
			photos.GET("/:id", photoHandler.GetByID)
			photos.POST("", photoHandler.Create)
			photos.PUT("/:id", photoHandler.Update)
			photos.DELETE("/:id", photoHandler.Delete)
			photos.POST("/:id/view", photoHandler.IncrementView)
		}

		// 相册相关路由
		albums := api.Group("/albums")
		{
			albums.GET("", albumHandler.GetAll)
			albums.GET("/:id", albumHandler.GetByID)
			albums.POST("", albumHandler.Create)
			albums.PUT("/:id", albumHandler.Update)
			albums.DELETE("/:id", albumHandler.Delete)
			albums.POST("/:id/photos", albumHandler.AddPhotoToAlbum)
			albums.DELETE("/:id/photos/:photo_id", albumHandler.RemovePhotoFromAlbum)
		}

		// 文件上传
		api.POST("/upload", photoHandler.UploadFile)
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
