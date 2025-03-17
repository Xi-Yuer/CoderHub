package model

import (
	"gorm.io/gorm"
)

type QuestionBankCategory struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey;index;" json:"id"`
	Name        string `gorm:"type:varchar(255);not null;unique;index:idx_name" json:"name"`
	Description string `gorm:"type:text;comment:题库分类描述" json:"description"`
}
