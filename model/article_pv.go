package model

import "gorm.io/gorm"
import "time"

// ArticlePV 文章的浏览量
type ArticlePV struct {
	gorm.Model
	ArticleID  int64     `gorm:"column:article_id;not null;uniqueIndex:idx_article_id"`
	Count      int64     `gorm:"column:count;not null;default:0"`
	LastSyncAt time.Time `gorm:"column:last_sync_at;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (ArticlePV) TableName() string {
	return "article_pvs"
}
