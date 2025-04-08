package model

import "time"

type ArticleTag struct {
	ID               int64      `gorm:"primaryKey" json:"id"`                    // 文章标签ID
	Name             string     `gorm:"size:255;not null;unique" json:"name"`    // 标签名称（唯一）
	Type             string     `gorm:"size:50;default:'article'" json:"type"`   // 标签类型（默认为文章标签）
	Description      string     `gorm:"size:500;default:''" json:"description"`  // 标签描述（可选）
	Icon             string     `gorm:"size:255;default:''" json:"icon"`         // 标签图标（可选）
	IsSystemProvider bool       `gorm:"default:false" json:"is_system_provider"` // 是否为系统标签
	UsageCount       int64      `gorm:"default:0" json:"usage_count"`            // 标签使用次数
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`        // 标签创建时间
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`        // 标签更新时间
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`       // 标签删除时间（软删除）
}
