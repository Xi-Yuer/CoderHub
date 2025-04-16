package model

import "gorm.io/gorm"

type ArticlesRelationLike struct {
	gorm.Model
	ArticleID int64 `gorm:"column:article_id;not null;index:idx_article_user,priority:1"`
	UserID    int64 `gorm:"column:user_id;not null;index:idx_article_user,priority:2"`

	// 联合唯一索引
	_ struct {
		ArticleID int64
		UserID    int64
	} `gorm:"uniqueIndex:uk_article_user"`
}
