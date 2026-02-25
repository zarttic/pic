package main

import (
	"fmt"
	"log"
	"strings"

	"picsite/internal/config"
	"picsite/internal/models"
	"picsite/internal/services"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := services.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 查询所有照片
	var photos []models.Photo
	if err := services.DB.Find(&photos).Error; err != nil {
		log.Fatalf("Failed to fetch photos: %v", err)
	}

	fmt.Printf("找到 %d 张照片\n", len(photos))

	fixedCount := 0

	for _, photo := range photos {
		needsUpdate := false
		updates := make(map[string]interface{})

		// 检查并修复 file_path
		if strings.HasPrefix(photo.FilePath, "/./uploads/") {
			correctPath := strings.Replace(photo.FilePath, "/./uploads/", "/uploads/", 1)
			updates["file_path"] = correctPath
			needsUpdate = true
			fmt.Printf("照片 %d - 修复 file_path: %s -> %s\n", photo.ID, photo.FilePath, correctPath)
		}

		// 检查并修复 thumbnail_path
		if strings.HasPrefix(photo.ThumbnailPath, "/./uploads/") {
			correctPath := strings.Replace(photo.ThumbnailPath, "/./uploads/", "/uploads/", 1)
			updates["thumbnail_path"] = correctPath
			needsUpdate = true
			fmt.Printf("照片 %d - 修复 thumbnail_path: %s -> %s\n", photo.ID, photo.ThumbnailPath, correctPath)
		}

		// 执行更新
		if needsUpdate {
			if err := services.DB.Model(&photo).Updates(updates).Error; err != nil {
				log.Printf("❌ 更新照片 %d 失败: %v\n", photo.ID, err)
			} else {
				fixedCount++
				fmt.Printf("✅ 照片 %d 已更新\n", photo.ID)
			}
		}
	}

	fmt.Printf("\n================================\n")
	fmt.Printf("总计: %d 张照片\n", len(photos))
	fmt.Printf("修复: %d 张照片\n", fixedCount)
	fmt.Printf("无需修复: %d 张照片\n", len(photos)-fixedCount)
	fmt.Printf("================================\n")
}
