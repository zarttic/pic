package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AlbumAccessSession 相册访问会话
type AlbumAccessSession struct {
	AlbumID   uint
	ExpiresAt time.Time
}

// AlbumSessionManager 相册会话管理器
type AlbumSessionManager struct {
	sessions map[string]*AlbumAccessSession
	mu       sync.RWMutex
}

var SessionManagerInstance = &AlbumSessionManager{
	sessions: make(map[string]*AlbumAccessSession),
}

// GenerateSessionToken 生成会话token (简化版)
func GenerateSessionToken(albumID uint, password string) string {
	return strconv.Itoa(int(albumID)) + "_" + password + "_" + time.Now().Format("20060102")
}

// SetSession 设置会话
func (m *AlbumSessionManager) SetSession(token string, albumID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[token] = &AlbumAccessSession{
		AlbumID:   albumID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}
}

// GetSession 获取会话
func (m *AlbumSessionManager) GetSession(token string) *AlbumAccessSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, exists := m.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil
	}
	return session
}

// CleanExpiredSessions 清理过期会话
func (m *AlbumSessionManager) CleanExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for token, session := range m.sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(m.sessions, token)
		}
	}
}

// AlbumAuthMiddleware 相册访问权限中间件
func AlbumAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		albumID := c.Param("id")
		if albumID == "" {
			c.Next()
			return
		}

		// 从 header 或 cookie 获取 token
		token := c.GetHeader("X-Album-Token")
		if token == "" {
			token, _ = c.Cookie("album_token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要验证密码"})
			c.Abort()
			return
		}

		// 验证会话
		session := SessionManagerInstance.GetSession(token)
		if session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "会话已过期，请重新验证"})
			c.Abort()
			return
		}

		c.Next()
	}
}
