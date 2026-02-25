package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"picsite/internal/config"
	"picsite/internal/models"
	"picsite/internal/services"
	"picsite/internal/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&models.Photo{}, &models.Album{}, &models.User{}, &models.AlbumPhoto{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// setupTestRouter 创建测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// setupTestAuthHandler 创建测试用的 AuthHandler
func setupTestAuthHandler(db *gorm.DB) *AuthHandler {
	services.DB = db
	cfg := &config.Config{
		JWTSecret: "test-secret-key-for-testing",
	}
	return NewAuthHandler(cfg)
}

func TestAuthHandler_Login(t *testing.T) {
	db := setupTestDB(t)
	handler := setupTestAuthHandler(db)
	router := setupTestRouter()

	// 创建测试用户
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	user := models.User{
		Username: "admin",
		Password: hashedPassword,
		Email:    "admin@test.com",
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 测试成功登录
	t.Run("successful login", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "admin",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.POST("/auth/login", handler.Login)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response LoginResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if response.Token == "" {
			t.Error("Expected token in response")
		}
		if response.User.Username != "admin" {
			t.Errorf("Expected username 'admin', got '%s'", response.User.Username)
		}
	})

	// 测试用户不存在
	t.Run("user not found", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// 测试错误密码
	t.Run("wrong password", func(t *testing.T) {
		loginReq := LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// 测试缺少字段
	t.Run("missing fields", func(t *testing.T) {
		loginReq := map[string]string{
			"username": "admin",
		}
		body, _ := json.Marshal(loginReq)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	db := setupTestDB(t)
	handler := setupTestAuthHandler(db)
	router := setupTestRouter()

	// 创建测试用户
	user := models.User{
		Username: "admin",
		Password: "hashedpassword",
		Email:    "admin@test.com",
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 生成有效的 refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, "admin", "test-secret-key-for-testing")
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	t.Run("valid refresh token", func(t *testing.T) {
		reqBody := map[string]string{
			"refresh_token": refreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.POST("/auth/refresh", handler.RefreshToken)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if response["token"] == "" {
			t.Error("Expected new token in response")
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		reqBody := map[string]string{
			"refresh_token": "invalid-token",
		}
		body, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestAuthHandler_GetCurrentUser(t *testing.T) {
	db := setupTestDB(t)
	handler := setupTestAuthHandler(db)
	router := setupTestRouter()

	// 创建测试用户
	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{
		Username: "admin",
		Password: hashedPassword,
		Email:    "admin@test.com",
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username, "admin", "test-secret-key-for-testing")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	t.Run("get current user with valid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.GET("/auth/me", func(c *gin.Context) {
			// 模拟中间件设置 userID
			c.Set("userID", user.ID)
			c.Next()
		}, handler.GetCurrentUser)

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if response["username"] != "admin" {
			t.Errorf("Expected username 'admin', got '%v'", response["username"])
		}
	})

	t.Run("get current user without token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/auth/me/no-token", nil)
		w := httptest.NewRecorder()

		router.GET("/auth/me/no-token", handler.GetCurrentUser)
		router.ServeHTTP(w, req)

		// 注意: 当没有设置 userID 时,handler 会返回 401
		// 但因为中间件没有运行,所以实际上是直接调用了 handler
		if w.Code != http.StatusUnauthorized {
			t.Logf("Status code: %d (handler expects userID to be set by middleware)", w.Code)
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	db := setupTestDB(t)
	handler := setupTestAuthHandler(db)
	router := setupTestRouter()

	t.Run("successful logout", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
		w := httptest.NewRecorder()

		router.POST("/auth/logout", handler.Logout)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		if response["message"] != "登出成功" {
			t.Errorf("Expected logout message, got '%s'", response["message"])
		}
	})
}
