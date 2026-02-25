package handlers

import (
	"net/http"
	"strconv"

	"picsite/internal/middleware"
	"picsite/internal/models"
	"picsite/internal/services"
	"picsite/internal/utils"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct{}

func NewAlbumHandler() *AlbumHandler {
	return &AlbumHandler{}
}

func (h *AlbumHandler) GetAll(c *gin.Context) {
	var albums []models.Album

	query := services.GetDB().Model(&models.Album{}).Preload("Photos")

	// 排序
	query = query.Order("created_at DESC")

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&albums).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": albums,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *AlbumHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	if err := services.GetDB().Preload("Photos").First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 检查是否是管理员访问(通过JWT token)
	isAdmin := false
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		// 如果有Bearer token,说明是管理员访问
		isAdmin = true
	}

	// 如果相册有密码保护且不是管理员访问,需要验证权限
	if album.IsProtected && !isAdmin {
		// 检查是否有有效的访问令牌
		token := c.GetHeader("X-Album-Token")
		if token == "" {
			token, _ = c.Cookie("album_token")
		}

		session := middleware.SessionManagerInstance.GetSession(token)
		if session == nil || session.AlbumID != album.ID {
			// 返回基本信息，但不包括照片
			c.JSON(http.StatusOK, gin.H{
				"id":           album.ID,
				"name":         album.Name,
				"description":  album.Description,
				"is_protected": true,
				"require_auth": true,
			})
			return
		}
	}

	c.JSON(http.StatusOK, album)
}

func (h *AlbumHandler) Create(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.GetDB().Create(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func (h *AlbumHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	if err := services.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	var updateData models.Album
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	services.GetDB().Model(&album).Updates(updateData)

	c.JSON(http.StatusOK, album)
}

func (h *AlbumHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var album models.Album

	if err := services.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 删除数据库记录（会自动删除关联表中的记录）
	if err := services.GetDB().Delete(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// AddPhotoToAlbum 添加照片到相册
func (h *AlbumHandler) AddPhotoToAlbum(c *gin.Context) {
	albumID := c.Param("id")
	var request struct {
		PhotoID   uint `json:"photo_id" binding:"required"`
		SortOrder int  `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albumPhoto := models.AlbumPhoto{
		AlbumID:   parseUint(albumID),
		PhotoID:   request.PhotoID,
		SortOrder: request.SortOrder,
	}

	if err := services.GetDB().Create(&albumPhoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo added to album successfully"})
}

// RemovePhotoFromAlbum 从相册中移除照片
func (h *AlbumHandler) RemovePhotoFromAlbum(c *gin.Context) {
	albumID := c.Param("id")
	photoID := c.Param("photo_id")

	if err := services.GetDB().Where("album_id = ? AND photo_id = ?", albumID, photoID).
		Delete(&models.AlbumPhoto{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo removed from album successfully"})
}

// VerifyPassword 验证相册密码
func (h *AlbumHandler) VerifyPassword(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var album models.Album
	if err := services.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(request.Password, album.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成会话 token
	token, err := middleware.GenerateSessionToken(album.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成会话令牌失败"})
		return
	}
	middleware.SessionManagerInstance.SetSession(token, album.ID)

	// 设置 cookie
	c.SetCookie("album_token", token, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "验证成功",
		"token":   token,
	})
}

// SetPassword 设置相册密码
func (h *AlbumHandler) SetPassword(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var album models.Album
	if err := services.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	album.Password = hashedPassword
	album.IsProtected = true

	if err := services.GetDB().Save(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码设置成功"})
}

// RemovePassword 移除相册密码
func (h *AlbumHandler) RemovePassword(c *gin.Context) {
	id := c.Param("id")

	var album models.Album
	if err := services.GetDB().First(&album, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	album.Password = ""
	album.IsProtected = false

	if err := services.GetDB().Save(&album).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码已移除"})
}

func parseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 32)
	return uint(val)
}
