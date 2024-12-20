package model

import (
	"time"

	"gorm.io/gorm"
)

// ImageRelation 图片关联模型（用于关联图片与其他业务实体）
type ImageRelation struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ImageID    int64          `gorm:"size:50;index" json:"image_id"`    // 图片ID
	EntityID   int64          `gorm:"index" json:"entity_id"`           // 关联实体ID
	EntityType string         `gorm:"size:50;index" json:"entity_type"` // 关联实体类型(comment/article等)
	Sort       int32          `gorm:"default:0" json:"sort"`            // 排序号
	CreatedAt  time.Time      `gorm:"<-:create" json:"created_at"`      // 创建时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`          // 删除时间
}

const (
	ImageRelationArticleCover   = "article_cover"
	ImageRelationArticleContent = "article_content"
	ImageRelationComment        = "comment"
	ImageRelationUserAvatar     = "user_avatar"
	ImageRelationQuestionCover  = "question_cover"
)
