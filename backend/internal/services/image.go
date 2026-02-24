package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const (
	ThumbnailWidth  = 400
	ThumbnailHeight = 0 // 0 means auto height based on aspect ratio
	ThumbnailQuality = 85
)

// GenerateThumbnail generates a thumbnail for the given image file
func GenerateThumbnail(srcPath, destPath string) error {
	// Open source image
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}
	defer srcFile.Close()

	// Decode image
	var img image.Image
	ext := strings.ToLower(filepath.Ext(srcPath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(srcFile)
	case ".png":
		img, err = png.Decode(srcFile)
	default:
		// Try to decode as generic image
		img, err = imaging.Decode(srcFile)
	}

	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Generate thumbnail
	thumbnail := imaging.Resize(img, ThumbnailWidth, ThumbnailHeight, imaging.Lanczos)

	// Ensure thumbnail directory exists
	thumbnailDir := filepath.Dir(destPath)
	if err := os.MkdirAll(thumbnailDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create thumbnail directory: %w", err)
	}

	// Save thumbnail
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail file: %w", err)
	}
	defer dstFile.Close()

	// Always save as JPEG for consistency
	err = jpeg.Encode(dstFile, thumbnail, &jpeg.Options{Quality: ThumbnailQuality})
	if err != nil {
		return fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return nil
}

// GenerateThumbnailFromUpload generates a thumbnail for uploaded file
func GenerateThumbnailFromUpload(uploadPath string) (string, error) {
	// Generate thumbnail path
	ext := filepath.Ext(uploadPath)
	thumbnailPath := strings.TrimSuffix(uploadPath, ext) + "_thumb.jpg"

	err := GenerateThumbnail(uploadPath, thumbnailPath)
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}
