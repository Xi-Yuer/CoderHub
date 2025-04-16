package model

import "gorm.io/gorm"

// UserFavorFolder 用户收藏夹
type UserFavorFolder struct {
	gorm.Model
	UserId      int64  `gorm:"column:user_id;type:bigint(20);not null;comment:用户id" json:"user_id"`
	FavorName   string `gorm:"column:favor_name;type:varchar(255);not null;comment:收藏夹名称" json:"favor_name"`
	FavorNum    int64  `gorm:"column:favor_num;type:bigint(20);not null;comment:收藏夹内问题数量" json:"favor_num"`
	IsPublic    bool   `gorm:"column:is_public;type:tinyint(1);not null;default:0;comment:是否公开" json:"is_public"`
	Description string `gorm:"column:description;type:varchar(255);comment:收藏夹描述" json:"description"`
	_           struct {
		UserID  int64
		FavorID int64
	} `gorm:"uniqueIndex:idx_user_id_favor_id"`
	_ struct{} `gorm:"index:idx_user_public,priority:1,columns:user_id,is_public,created_at"` // 添加新的复合索引
}

type UserFavor struct {
	gorm.Model
	UserId      int64  `gorm:"column:user_id;type:bigint(20);not null;comment:用户id" json:"user_id"`
	FavorFoldId int64  `gorm:"column:favor_fold_id;type:bigint(20);not null;comment:收藏夹id" json:"favor_fold_id"`
	EntityId    int64  `gorm:"column:entity_id;type:bigint(20);not null;comment:实体id" json:"entity_id"`
	EntityType  string `gorm:"type:enum('question','article');not null;comment:收藏类型" json:"entity_type"`
	_           struct {
		UserID      int64
		FavorFoldId int64
		EntityID    int64
	} `gorm:"uniqueIndex:idx_user_id_favor_id_entity_id"`
	// 添加新的复合索引，包含软删除字段，优化计数查询
	_ struct{} `gorm:"index:idx_entity_type_del,priority:1,columns:entity_id,entity_type,deleted_at"`
}

var FavorEntityEnum = map[string]string{
	"question": "问题",
	"article":  "文章",
}

const (
	UserFavorRelationQuestion = "question"
	UserFavorRelationArticle  = "article"
)
