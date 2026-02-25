package services

import (
	"testing"
)

func TestExtractEXIFFromUpload(t *testing.T) {
	t.Run("extract EXIF from non-existent file", func(t *testing.T) {
		exifData := ExtractEXIFFromUpload("non_existent_file.jpg")

		// 应该返回空结构,不应该 panic
		if exifData == nil {
			t.Error("Expected non-nil EXIF data")
		}

		if exifData.CameraModel != "" {
			t.Error("Expected empty camera model for non-existent file")
		}
	})
}

func TestGenerateThumbnailFromUpload(t *testing.T) {
	t.Run("generate thumbnail from non-existent file", func(t *testing.T) {
		thumbnailPath, err := GenerateThumbnailFromUpload("non_existent_file.jpg")

		// 应该返回错误
		if err == nil {
			t.Error("Expected error for non-existent file")
		}

		if thumbnailPath != "" {
			t.Error("Expected empty thumbnail path for non-existent file")
		}
	})
}
