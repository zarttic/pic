package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"picsite/internal/models"
	"picsite/internal/services"
	"strconv"
	"strings"
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

	// 搜索功能 - 支持多字段搜索
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where(
			"title LIKE ? OR description LIKE ? OR location LIKE ? OR tags LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}

	// 筛选
	if featured := c.Query("featured"); featured != "" {
		query = query.Where("is_featured = ?", featured == "true")
	}

	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	if location := c.Query("location"); location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	if year := c.Query("year"); year != "" {
		query = query.Where("year = ?", year)
	}

	if camera := c.Query("camera"); camera != "" {
		query = query.Where("camera_model LIKE ?", "%"+camera+"%")
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
	shotDate := c.PostForm("shot_date")
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

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅支持 jpg, jpeg, png, webp"})
		return
	}

	// 验证文件大小 (10MB)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制 (最大 10MB)"})
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

	// 提取 EXIF 信息
	exifData := services.ExtractEXIFFromUpload(uploadPath)

	// 生成缩略图
	thumbnailPath, err := services.GenerateThumbnailFromUpload(uploadPath)
	if err != nil {
		fmt.Printf("Failed to generate thumbnail: %v\n", err)
		// 不阻止上传，继续保存记录
	}

	// 使用 EXIF 数据填充表单数据（如果表单中没有提供）
	if cameraModel == "" && exifData.CameraModel != "" {
		cameraModel = exifData.CameraModel
	}
	if lens == "" && exifData.Lens != "" {
		lens = exifData.Lens
	}
	if aperture == "" && exifData.Aperture != "" {
		aperture = exifData.Aperture
	}
	if shutterSpeed == "" && exifData.ShutterSpeed != "" {
		shutterSpeed = exifData.ShutterSpeed
	}
	if iso == 0 && exifData.ISO != 0 {
		iso = exifData.ISO
	}
	var shotDateValue *time.Time
	if shotDate == "" && exifData.ShotDate != nil {
		shotDateValue = exifData.ShotDate
		if year == 0 {
			year = exifData.ShotDate.Year()
		}
	}

	// 创建照片记录
	photo := models.Photo{
		Title:         title,
		Description:   description,
		FilePath:      "/" + strings.TrimPrefix(uploadPath, "./"),
		ThumbnailPath: "/" + strings.TrimPrefix(thumbnailPath, "./"),
		Location:      location,
		ShotDate:      shotDateValue,
		Year:          year,
		CameraModel:   cameraModel,
		Lens:          lens,
		Aperture:      aperture,
		ShutterSpeed:  shutterSpeed,
		ISO:           iso,
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

// BatchDelete 批量删除照片
func (h *PhotoHandler) BatchDelete(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误，最多支持100张照片"})
		return
	}

	// 查询要删除的照片
	var photos []models.Photo
	if err := services.GetDB().Where("id IN ?", request.IDs).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询照片失败"})
		return
	}

	// 删除文件
	for _, photo := range photos {
		if photo.FilePath != "" {
			filePath := "." + photo.FilePath
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to delete file: %v\n", err)
			}
		}
		if photo.ThumbnailPath != "" {
			thumbnailPath := "." + photo.ThumbnailPath
			if err := os.Remove(thumbnailPath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to delete thumbnail: %v\n", err)
			}
		}
	}

	// 批量删除数据库记录
	if err := services.GetDB().Where("id IN ?", request.IDs).Delete(&models.Photo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量删除成功",
		"deleted": len(photos),
	})
}

// BatchUpdateTags 批量更新标签
func (h *PhotoHandler) BatchUpdateTags(c *gin.Context) {
	var request struct {
		IDs  []uint `json:"ids" binding:"required,min=1,max=100"`
		Tags string `json:"tags" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 批量更新
	if err := services.GetDB().Model(&models.Photo{}).Where("id IN ?", request.IDs).
		Update("tags", request.Tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量更新标签成功",
		"updated": len(request.IDs),
	})
}

// BatchUpdateFeatured 批量设置精选
func (h *PhotoHandler) BatchUpdateFeatured(c *gin.Context) {
	var request struct {
		IDs        []uint `json:"ids" binding:"required,min=1,max=100"`
		IsFeatured bool   `json:"is_featured"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 批量更新
	if err := services.GetDB().Model(&models.Photo{}).Where("id IN ?", request.IDs).
		Update("is_featured", request.IsFeatured).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	action := "取消精选"
	if request.IsFeatured {
		action = "设置精选"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量" + action + "成功",
		"updated": len(request.IDs),
	})
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

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型，仅支持 jpg, jpeg, png, webp"})
		return
	}

	// 验证文件大小 (10MB)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制 (最大 10MB)"})
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
		"file_path":     "/" + uploadPath,
		"original_name": file.Filename,
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
