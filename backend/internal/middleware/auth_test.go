package middleware

import (
	"net/http"
	"net/http/httptest"
	"picsite/internal/utils"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouterForMiddleware() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-key"

	// 生成有效的测试 token
	validToken, err := utils.GenerateToken(1, "admin", "admin", secret)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// 测试用例
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name:           "invalid authorization format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name:           "invalid bearer prefix",
			authHeader:     "Basic token123",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token-string",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 为每个测试创建新的 router
			router := setupTestRouterForMiddleware()
			abortCount := 0
			router.GET("/protected", AuthMiddleware(secret), func(c *gin.Context) {
				abortCount++
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.shouldAbort && abortCount != 0 {
				t.Error("Expected request to be aborted but handler was called")
			}

			if !tt.shouldAbort && abortCount != 1 {
				t.Error("Expected handler to be called but request was aborted")
			}
		})
	}
}

func TestAuthMiddleware_SetsContext(t *testing.T) {
	router := setupTestRouterForMiddleware()
	secret := "test-secret-key"

	// 生成有效的测试 token
	validToken, err := utils.GenerateToken(1, "admin", "admin", secret)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// 测试中间件是否正确设置上下文
	var contextUserID interface{}
	var contextUsername interface{}
	var contextRole interface{}

	router.GET("/protected", AuthMiddleware(secret), func(c *gin.Context) {
		contextUserID, _ = c.Get("userID")
		contextUsername, _ = c.Get("username")
		contextRole, _ = c.Get("role")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if contextUserID.(uint) != 1 {
		t.Errorf("Expected userID 1, got %v", contextUserID)
	}
	if contextUsername.(string) != "admin" {
		t.Errorf("Expected username 'admin', got %v", contextUsername)
	}
	if contextRole.(string) != "admin" {
		t.Errorf("Expected role 'admin', got %v", contextRole)
	}
}
