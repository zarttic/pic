package services

import (
	"picsite/internal/models"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestInitDB(t *testing.T) {
	t.Run("initialize database successfully", func(t *testing.T) {
		// 使用内存数据库测试
		var err error
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("Failed to initialize database: %v", err)
		}

		// 测试自动迁移
		err = DB.AutoMigrate(&models.Photo{}, &models.Album{}, &models.User{}, &models.AlbumPhoto{})
		if err != nil {
			t.Fatalf("Failed to migrate database: %v", err)
		}

		// 验证表是否存在
		if !DB.Migrator().HasTable(&models.Photo{}) {
			t.Error("Expected Photos table to exist")
		}
		if !DB.Migrator().HasTable(&models.Album{}) {
			t.Error("Expected Albums table to exist")
		}
		if !DB.Migrator().HasTable(&models.User{}) {
			t.Error("Expected Users table to exist")
		}
		if !DB.Migrator().HasTable(&models.AlbumPhoto{}) {
			t.Error("Expected AlbumPhotos table to exist")
		}
	})
}

func TestGetDB(t *testing.T) {
	t.Run("get database instance", func(t *testing.T) {
		// 先初始化数据库
		var err error
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("Failed to initialize database: %v", err)
		}

		// 执行迁移
		err = DB.AutoMigrate(&models.Photo{}, &models.Album{}, &models.User{}, &models.AlbumPhoto{})
		if err != nil {
			t.Fatalf("Failed to migrate database: %v", err)
		}

		// 获取数据库实例
		db := GetDB()
		if db == nil {
			t.Error("Expected non-nil database instance")
		}

		// 验证可以执行查询
		var count int64
		err = db.Model(&models.Photo{}).Count(&count).Error
		if err != nil {
			t.Errorf("Failed to execute query: %v", err)
		}
	})
}
