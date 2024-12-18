package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Articles struct {
	ID           int64          `gorm:"<-:create;primaryKey" json:"id"`                               // 主键 ID
	Type         string         `gorm:"type:enum('article','micro_post');not null" json:"type"`       // 内容类型：长文(article) 或 短文(micro_post)
	Title        string         `gorm:"size:255" json:"title"`                                        // 长文标题，短文可为空
	Content      string         `gorm:"type:longtext;not null" json:"content"`                        // 内容（长文或短文）
	Summary      string         `gorm:"type:text" json:"summary"`                                     // 长文摘要，短文为空
	AuthorID     int64          `gorm:"not null;index" json:"author_id"`                              // 作者 ID
	Images       []Image        `gorm:"-" json:"images"`                                              // 文章图片列表
	CoverImage   *Image         `gorm:"-" json:"cover_image,omitempty"`                               // 封面图片
	Tags         string         `gorm:"size:255" json:"tags"`                                         // 标签，逗号分隔（适用于长文）
	CommentCount int64          `gorm:"default:0" json:"comment_count"`                               // 评论数
	Status       string         `gorm:"type:enum('draft','published');default:'draft'" json:"status"` // 内容状态
	Version      int64          `gorm:"default:0" json:"version"`                                     // 版本号
	CreatedAt    time.Time      `gorm:"<-:create" json:"created_at"`                                  // 创建时间
	UpdatedAt    time.Time      `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"`              // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`                                      // 删除时间

	ViewCount int64 `gorm:"-" json:"view_count"` // 阅读次数（长文专用）
	LikeCount int64 `gorm:"-" json:"like_count"` // 点赞次数

}

func (a *Articles) CacheKeyByID(id int64) string {
	return fmt.Sprintf("Articles:id:%d", id)
}
