package model

import (
	"time"

	"gorm.io/gorm"
)

// Image 通用图片模型
type Image struct {
	ID           int64          `gorm:"primaryKey;size:50" json:"image_id"`     // 图片ID，使用uuid或其他唯一标识
	BucketName   string         `gorm:"size:100;not null" json:"bucket_name"`   // MinIO bucket名称
	ObjectName   string         `gorm:"size:100;not null" json:"object_name"`   // MinIO中的对象名称
	URL          string         `gorm:"size:2083;not null" json:"url"`          // 完整的访问URL
	ThumbnailURL string         `gorm:"size:2083" json:"thumbnail_url"`         // 缩略图URL
	ContentType  string         `gorm:"size:100;not null" json:"content_type"`  // 文件MIME类型
	Size         int64          `json:"size"`                                   // 文件大小(bytes)
	Width        int32          `json:"width"`                                  // 图片宽度(px)
	Height       int32          `json:"height"`                                 // 图片高度(px)
	UploadIP     string         `gorm:"size:50" json:"upload_ip"`               // 上传者IP
	UserID       int64          `gorm:"index" json:"user_id"`                   // 上传者ID
	Status       string         `gorm:"size:20;default:'active'" json:"status"` // 状态：active,deleted
	CreatedAt    time.Time      `gorm:"<-:create" json:"created_at"`            // 创建时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`                // 删除时间
}

const (
	ImageStatusActive = "active"
	ImageStatusDeleted = "deleted"
)
