package main

import (
	"fmt"
	"log"
	"os"

	"picsite/internal/config"
	"picsite/internal/models"
	"picsite/internal/services"
	"picsite/internal/utils"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := services.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 获取管理员信息
	username := getEnv("ADMIN_USERNAME", "admin")
	password := getEnv("ADMIN_PASSWORD", "admin123")
	email := getEnv("ADMIN_EMAIL", "admin@example.com")

	// 检查是否已存在
	var existingUser models.User
	result := services.DB.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		fmt.Printf("管理员用户 '%s' 已存在\n", username)
		return
	}

	// 创建管理员用户
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Role:     "admin",
	}

	if err := services.DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("✓ 管理员用户创建成功!\n")
	fmt.Printf("  用户名: %s\n", username)
	fmt.Printf("  密码: %s\n", password)
	fmt.Printf("  邮箱: %s\n", email)
	fmt.Println("\n⚠️  请在生产环境中立即修改默认密码!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
