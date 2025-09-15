package model

import "gorm.io/gorm"

type AcademicNavigator struct {
	gorm.Model
	UserId    int64  `gorm:"column:user_id;not null;index:idx_user_id"`     // 用户ID
	Content   string `gorm:"column:content;not null"`                       // 内容
	Education string `gorm:"column:education;not null;index:idx_education"` // 学历
	Major     string `gorm:"column:major;not null;index:idx_major"`         // 专业
	School    string `gorm:"column:school;not null;index:idx_school"`       // 学校
	WorkExp   string `gorm:"column:work_exp;not null;index:idx_work_exp"`   // 工作经验

	// 联合索引
	_ struct {
		UserID    int64
		Education string
		Major     string
		School    string
	} `gorm:"index:idx_user_id_education_major_school"`
}
