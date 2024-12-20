package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID         int64          `gorm:"<-:create;primaryKey" json:"id"`                  // 评论ID
	EntityID   int64          `gorm:"index" json:"entity_id"`                          // 文章ID
	Content    string         `gorm:"type:text;not null" json:"content"`               // 评论内容
	RootID     int64          `gorm:"default:0;index" json:"root_id"`                  // 根评论ID，0表示顶级评论
	ParentID   int64          `gorm:"default:0;index" json:"parent_id"`                // 父评论ID，0表示顶级评论
	UserID     int64          `gorm:"index" json:"user_id"`                            // 评论者ID
	CreatedAt  time.Time      `gorm:"<-:create" json:"created_at"`                     // 创建时间
	UpdatedAt  time.Time      `gorm:"autoCreateTime;autoUpdateTime" json:"updated_at"` // 更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`                         // 删除时间
	Version    int32          `gorm:"default:0" json:"version"`                        // 版本号:乐观锁
	ReplyToUID int64          `gorm:"index" json:"reply_to_uid"`                       // 回复目标用户ID
	// gorm:"-" 标签的含义是告诉 GORM 忽略这个字段，即这个字段不会被映射到数据库表中。
	// 这些图片ID可能存储在另一个关联表中，而不是直接存储在评论表里
	ImageIDs   []string  `gorm:"-" json:"image_ids"`   // 评论图片ID列表
	Replies    []Comment `gorm:"-" json:"replies"`     // 子评论列表
	ReplyCount int64     `gorm:"-" json:"reply_count"` // 回复数量
	LikeCount  int64     `gorm:"-" json:"like_count"`  // 点赞次数
}

// CacheKeyByID 生成评论缓存键
func (c *Comment) CacheKeyByID(id int64) string {
	return fmt.Sprintf("Comment:id:%d", id)
}
