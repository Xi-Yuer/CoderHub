package model

import "gorm.io/gorm"

type CommentRelationLike struct {
	gorm.Model
	CommentID int64 `gorm:"column:comment_id;not null;index:idx_comment_id"`
	UserID    int64 `gorm:"column:user_id;not null;index:idx_user_id"`
	
	// 联合唯一索引
	_ struct {
		CommentID int64
		UserID    int64
	} `gorm:"uniqueIndex:idx_comment_id_user_id"`
}
