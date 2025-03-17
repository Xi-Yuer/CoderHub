package model

import (
	"gorm.io/gorm"
	"time"
)

// QuestionBank 题库模型
type QuestionBank struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(255);not null;unique;index:idx_name;comment:题库名称" json:"name"`
	Description string `gorm:"type:text;comment:题库描述" json:"description"`
	Difficulty  string `gorm:"type:enum('default','easy','medium','hard');comment:题库难度" json:"difficulty"`
	Tags        string `gorm:"type:text;comment:题库标签（JSON格式）" json:"tags"`
	CreateUser  int64  `gorm:"type:bigint;not null;comment:创建人" json:"create_user"`
	CoverImage  *Image `gorm:"-" json:"cover_image,omitempty"`
	CateGoryID  int64  `gorm:"index;not null;comment:题库分类" json:"cate_gory_id"`
}

// Question 题目模型
type Question struct {
	gorm.Model
	BankID       int64        `json:"bank_id" gorm:"type:bigint;index;comment:题库ID"`
	Title        string       `json:"title" gorm:"type:varchar(255);not null;comment:题目标题"`
	Content      string       `json:"content" gorm:"type:text;not null;comment:题目内容"`
	CreateUser   int64        `json:"create_user" gorm:"type:bigint;not null;comment:创建人"`
	Difficulty   string       `json:"difficulty" gorm:"type:enum('default','easy','medium','hard');comment:题库难度"`
	QuestionBank QuestionBank `gorm:"foreignKey:BankID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"question_bank,omitempty"`
}

type QuestionBanksPreviewWithCreateUser struct {
	ID             int64     `json:"id"`               // 题库ID
	Name           string    `json:"name"`             // 题库名称
	CoverImage     string    `json:"cover_image"`      // 题库封面图片
	Description    string    `json:"description"`      // 题库描述
	Difficulty     string    `json:"difficulty"`       // 题库难度
	Tags           string    `json:"tags"`             // 题库标签
	CreateTime     time.Time `json:"create_time"`      // 创建时间
	CreateUserID   int64     `json:"create_user_id"`   // 创建人ID
	CreateUserName string    `json:"create_user_name"` // 创建人
	Avatar         string    `json:"avatar"`           // 创建人头像
}
