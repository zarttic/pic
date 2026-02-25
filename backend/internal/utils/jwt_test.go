package utils

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		role     string
		secret   string
		wantErr  bool
	}{
		{
			name:     "valid token generation",
			userID:   1,
			username: "admin",
			role:     "admin",
			secret:   "test-secret-key",
			wantErr:  false,
		},
		{
			name:     "empty secret",
			userID:   1,
			username: "admin",
			role:     "admin",
			secret:   "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username, tt.role, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		role     string
		secret   string
		wantErr  bool
	}{
		{
			name:     "valid refresh token generation",
			userID:   1,
			username: "admin",
			role:     "admin",
			secret:   "test-secret-key",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tt.userID, tt.username, tt.role, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateRefreshToken() returned empty token")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key"

	// 生成有效 token
	validToken, err := GenerateToken(1, "admin", "admin", secret)
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	// 生成刷新 token
	refreshToken, err := GenerateRefreshToken(1, "admin", "admin", secret)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secret    string
		wantErr   bool
		errType   error
		checkUser func(*Claims) bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			secret:  secret,
			wantErr: false,
			checkUser: func(c *Claims) bool {
				return c.UserID == 1 && c.Username == "admin" && c.Role == "admin"
			},
		},
		{
			name:    "valid refresh token",
			token:   refreshToken,
			secret:  secret,
			wantErr: false,
			checkUser: func(c *Claims) bool {
				return c.UserID == 1 && c.Username == "admin"
			},
		},
		{
			name:    "invalid token format",
			token:   "invalid-token-string",
			secret:  secret,
			wantErr: true,
			errType: ErrInvalidToken,
		},
		{
			name:    "empty token",
			token:   "",
			secret:  secret,
			wantErr: true,
			errType: ErrInvalidToken,
		},
		{
			name:    "wrong secret",
			token:   validToken,
			secret:  "wrong-secret",
			wantErr: true,
			errType: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errType != nil && err != tt.errType {
					t.Errorf("ParseToken() error type = %v, want %v", err, tt.errType)
				}
				return
			}

			if tt.checkUser != nil && !tt.checkUser(claims) {
				t.Errorf("ParseToken() claims validation failed: %+v", claims)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	secret := "test-secret-key"

	t.Run("token should have correct expiration", func(t *testing.T) {
		token, err := GenerateToken(1, "admin", "admin", secret)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		claims, err := ParseToken(token, secret)
		if err != nil {
			t.Fatalf("Failed to parse token: %v", err)
		}

		// 验证过期时间在 24 小时左右
		expectedExpiry := time.Now().Add(24 * time.Hour)
		actualExpiry := claims.ExpiresAt.Time

		// 允许 1 分钟误差
		diff := actualExpiry.Sub(expectedExpiry)
		if diff < -time.Minute || diff > time.Minute {
			t.Errorf("Token expiration time incorrect. Diff: %v", diff)
		}
	})

	t.Run("refresh token should have correct expiration", func(t *testing.T) {
		token, err := GenerateRefreshToken(1, "admin", "admin", secret)
		if err != nil {
			t.Fatalf("Failed to generate refresh token: %v", err)
		}

		claims, err := ParseToken(token, secret)
		if err != nil {
			t.Fatalf("Failed to parse refresh token: %v", err)
		}

		// 验证过期时间在 7 天左右
		expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
		actualExpiry := claims.ExpiresAt.Time

		// 允许 1 分钟误差
		diff := actualExpiry.Sub(expectedExpiry)
		if diff < -time.Minute || diff > time.Minute {
			t.Errorf("Refresh token expiration time incorrect. Diff: %v", diff)
		}
	})
}
