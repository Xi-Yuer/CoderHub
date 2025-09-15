package model

import "gorm.io/gorm"

type AcademicRelationLike struct {
	gorm.Model
	AcademicNavigatorID int64 `gorm:"column:academic_navigator_id;not null;index:idx_academic_navigator_id"` // 学术导航ID
	UserID              int64 `gorm:"column:user_id;not null;index:idx_user_id"`                             // 用户ID

	// 联合唯一索引
	_ struct {
		AcademicNavigatorID int64
		UserID              int64
	} `gorm:"uniqueIndex:idx_academic_navigator_id_user_id"`
}
