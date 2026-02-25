package services

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// EXIFData EXIF 数据结构
type EXIFData struct {
	CameraModel  string     `json:"camera_model"`
	Lens         string     `json:"lens"`
	Aperture     string     `json:"aperture"`
	ShutterSpeed string     `json:"shutter_speed"`
	ISO          int        `json:"iso"`
	ShotDate     *time.Time `json:"shot_date"`
}

// ExtractEXIF 从图片文件提取 EXIF 信息
func ExtractEXIF(filePath string) (*EXIFData, error) {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 检查是否包含 EXIF 数据
	entries, _, err := exif.GetFlatExifData(data, nil)
	if err != nil {
		// 没有 EXIF 数据不是错误，返回空结构
		return &EXIFData{}, nil
	}

	exifData := &EXIFData{}

	// 解析 EXIF 数据
	for _, entry := range entries {
		if entry.TagName == "" || entry.Value == nil {
			continue
		}

		switch entry.TagName {
		case "Model":
			if model, ok := entry.Value.(string); ok {
				exifData.CameraModel = strings.TrimSpace(model)
			}

		case "LensModel":
			if lens, ok := entry.Value.(string); ok {
				exifData.Lens = strings.TrimSpace(lens)
			}

		case "FNumber":
			exifData.Aperture = formatAperture(entry.Value)

		case "ExposureTime":
			exifData.ShutterSpeed = formatShutterSpeed(entry.Value)

		case "ISOSpeedRatings":
			if iso, ok := entry.Value.([]uint16); ok && len(iso) > 0 {
				exifData.ISO = int(iso[0])
			} else if iso, ok := entry.Value.(uint16); ok {
				exifData.ISO = int(iso)
			}

		case "DateTimeOriginal":
			if dateStr, ok := entry.Value.(string); ok {
				// EXIF 日期格式: "2006:01:02 15:04:05"
				t, err := time.Parse("2006:01:02 15:04:05", dateStr)
				if err == nil {
					exifData.ShotDate = &t
				}
			}
		}
	}

	return exifData, nil
}

// formatAperture 格式化光圈值
func formatAperture(value interface{}) string {
	switch v := value.(type) {
	case []exifcommon.Rational:
		if len(v) > 0 {
			aperture := float64(v[0].Numerator) / float64(v[0].Denominator)
			return fmt.Sprintf("f/%.1f", aperture)
		}
	case exifcommon.Rational:
		aperture := float64(v.Numerator) / float64(v.Denominator)
		return fmt.Sprintf("f/%.1f", aperture)
	}
	return ""
}

// formatShutterSpeed 格式化快门速度
func formatShutterSpeed(value interface{}) string {
	switch v := value.(type) {
	case []exifcommon.Rational:
		if len(v) > 0 {
			shutter := float64(v[0].Numerator) / float64(v[0].Denominator)
			if shutter >= 1 {
				return fmt.Sprintf("%.0fs", shutter)
			}
			return fmt.Sprintf("1/%.0fs", 1/shutter)
		}
	case exifcommon.Rational:
		shutter := float64(v.Numerator) / float64(v.Denominator)
		if shutter >= 1 {
			return fmt.Sprintf("%.0fs", shutter)
		}
		return fmt.Sprintf("1/%.0fs", 1/shutter)
	}
	return ""
}

// ExtractEXIFFromUpload 从上传的图片提取 EXIF
func ExtractEXIFFromUpload(filePath string) *EXIFData {
	exifData, err := ExtractEXIF(filePath)
	if err != nil {
		log.Printf("Failed to extract EXIF: %v", err)
		return &EXIFData{}
	}
	return exifData
}
