package model

import "gorm.io/gorm"

type SandpackFiles struct {
	gorm.Model
	SandpackID int64  `gorm:"column:sandpack_id;not null;index:sandpack_id"` // 项目ID
	FileName   string `gorm:"column:file_name;not null"`                     // 文件名称
	Language   string `gorm:"column:language;not null"`                      // 语言类型
	Content    string `gorm:"column:content;not null"`                       //  文件内容
}
