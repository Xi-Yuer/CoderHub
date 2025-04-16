package model

import (
	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	FollowerID int64 `gorm:"column:follower_id;not null"` // 关注者ID
	FollowedID int64 `gorm:"column:followed_id;not null"` // 被关注者ID

	// 优化索引结构
	_ struct{} `gorm:"index:idx_follower_followed,priority:1,columns:follower_id,followed_id"`
	_ struct{} `gorm:"index:idx_followed_follower,priority:1,columns:followed_id,follower_id"`
	_ struct{} `gorm:"uniqueIndex:uk_follower_followed,priority:1,columns:follower_id,followed_id,deleted_at"`
}
