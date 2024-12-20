package model

import (
	"gorm.io/gorm"
)

// QuestionBank 题库模型
type QuestionBank struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `json:"name" gorm:"type:varchar(255);not null;unique;index:idx_name;comment:题库名称"`  // 题库名称，添加唯一约束和索引
	Description string `json:"description" gorm:"type:text;comment:题库描述"`                                  // 题库描述
	Difficulty  string `json:"difficulty" gorm:"type:enum('default','easy','medium','hard');comment:题库难度"` // 题库难度（例如：default, easy, medium, hard）
	Tags        string `json:"tags" gorm:"type:text;comment:题库标签（JSON格式）"`                                 // 标签（存储标签列表的JSON格式）
	CreateUser  int64  `json:"create_user" gorm:"type:bigint;not null;comment:创建人"`                        // 创建人
	CoverImage  *Image `gorm:"-" json:"cover_image,omitempty"`                                             // 封面图片
}

// Question 题目模型
type Question struct {
	gorm.Model
	BankID     int64  `json:"bank_id" gorm:"not null;index;comment:题库ID"`                                 // 外键，关联到题库
	Title      string `json:"title" gorm:"type:varchar(255);not null;comment:题目标题"`                       // 题目标题
	Content    string `json:"content" gorm:"type:text;not null;comment:题目内容"`                             // 题目内容
	CreateUser int64  `json:"create_user" gorm:"type:bigint;not null;comment:创建人"`                        // 创建人
	Difficulty string `json:"difficulty" gorm:"type:enum('default','easy','medium','hard');comment:题库难度"` // 题目难度（例如：default, easy, medium, hard）
}
