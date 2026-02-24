package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"photography-website/internal/models"
	"photography-website/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct{}

func NewPhotoHandler() *PhotoHandler {
	return &PhotoHandler{}
}

func (h *PhotoHandler) GetAll(c *gin.Context) {
	var photos []models.Photo

	query := services.GetDB().Model(&models.Photo{})

	// 筛选
	if featured := c.Query("featured"); featured != "" {
		query = query.Where("is_featured = ?", featured == "true")
	}

	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	// 排序
	query = query.Order("created_at DESC")

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": photos,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *PhotoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo

	if err := services.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, photo)
}

func (h *PhotoHandler) Create(c *gin.Context) {
	// 解析表单数据
	title := c.PostForm("title")
	description := c.PostForm("description")
	location := c.PostForm("location")
	year, _ := strconv.Atoi(c.PostForm("year"))
	cameraModel := c.PostForm("camera_model")
	lens := c.PostForm("lens")
	aperture := c.PostForm("aperture")
	shutterSpeed := c.PostForm("shutter_speed")
	iso, _ := strconv.Atoi(c.PostForm("iso"))

	// 处理文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	uploadPath := "./uploads/" + filename

	// 确保上传目录存在
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 创建照片记录
	photo := models.Photo{
		Title:        title,
		Description:  description,
		FilePath:     "/" + uploadPath,
		Location:     location,
		Year:         year,
		CameraModel:  cameraModel,
		Lens:         lens,
		Aperture:     aperture,
		ShutterSpeed: shutterSpeed,
		ISO:          iso,
	}

	if err := services.GetDB().Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

func (h *PhotoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo

	if err := services.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	var updateData models.Photo
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	services.GetDB().Model(&photo).Updates(updateData)

	c.JSON(http.StatusOK, photo)
}

func (h *PhotoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo

	if err := services.GetDB().First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// 删除文件
	if photo.FilePath != "" {
		filePath := "." + photo.FilePath
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			// 文件删除失败，但继续删除数据库记录
			fmt.Printf("Failed to delete file: %v\n", err)
		}
	}

	// 删除数据库记录
	if err := services.GetDB().Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func (h *PhotoHandler) IncrementView(c *gin.Context) {
	id := c.Param("id")

	if err := services.GetDB().Model(&models.Photo{}).Where("id = ?", id).
		UpdateColumn("view_count", services.GetDB().Raw("view_count + 1")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View count incremented"})
}

// UploadFile 单独的文件上传接口
func (h *PhotoHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	uploadPath := "./uploads/" + filename

	// 确保上传目录存在
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 保存文件
	dst, err := os.Create(uploadPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_path":      "/" + uploadPath,
		"original_name":  file.Filename,
		"size":          file.Size,
	})
}

func (h *PhotoHandler) ServeImage(c *gin.Context) {
	filePath := c.Param("filepath")
	fullPath := filepath.Join("./uploads", filePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.File(fullPath)
}
