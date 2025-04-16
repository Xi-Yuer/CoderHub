package model

import (
	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	FollowerID int64 `gorm:"column:follower_id;not null"` // 关注者ID
	FollowedID int64 `gorm:"column:followed_id;not null"` // 被关注者ID

	// 优化后的联合索引，包含软删除字段
	_ struct{} `gorm:"index:idx_follower_followed_del,priority:1,columns:follower_id,followed_id,deleted_at"`
}
