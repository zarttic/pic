package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"not null"`
	Description  string         `json:"description"`
	FilePath     string         `json:"file_path" gorm:"not null"`
	ThumbnailPath string        `json:"thumbnail_path"`
	Location     string         `json:"location"`
	ShotDate     *time.Time     `json:"shot_date"`
	Year         int            `json:"year"`
	CameraModel  string         `json:"camera_model"`
	Lens         string         `json:"lens"`
	Aperture     string         `json:"aperture"`
	ShutterSpeed string         `json:"shutter_speed"`
	ISO          int            `json:"iso"`
	Tags         string         `json:"tags"` // JSON array stored as string
	IsFeatured   bool           `json:"is_featured" gorm:"default:false"`
	ViewCount    int            `json:"view_count" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Album struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	CoverPhotoID *uint          `json:"cover_photo_id"`
	Password     string         `json:"-"` // 密码不返回给前端
	IsProtected  bool           `json:"is_protected" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Photos       []Photo        `json:"photos" gorm:"many2many:album_photos;"`
}

type AlbumPhoto struct {
	AlbumID   uint `gorm:"primaryKey"`
	PhotoID   uint `gorm:"primaryKey"`
	SortOrder int  `json:"sort_order"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
