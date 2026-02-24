package handlers

import (
	"net/http"
	"strconv"

	"picsite/internal/models"
	"picsite/internal/services"

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

func parseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 32)
	return uint(val)
}
