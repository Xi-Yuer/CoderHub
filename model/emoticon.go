package model

import "gorm.io/gorm"

// Emoticon 表
type Emoticon struct {
	gorm.Model
	Code        string `gorm:"size:255;not null" json:"code"`        // 表情代码
	URL         string `gorm:"size:255;not null" json:"url"`         // 表情图片地址
	Description string `gorm:"size:255;not null" json:"description"` // 表情描述
}
