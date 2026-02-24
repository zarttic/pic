package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBPath     string
	UploadPath string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "./photography.db"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
