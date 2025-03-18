package model

import "time"

type SchoolExp struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Education    string    `gorm:"column:education;not null;index:idx_education" json:"education"`           // 教育背景
	Major        string    `gorm:"column:major;not null;index:idx_major" json:"major"`                       // 专业
	School       string    `gorm:"column:school;not null;index:idx_school" json:"school"`                    // 学校
	WorkExp      string    `gorm:"column:work_exp;not null;index:idx_work_exp" json:"work_exp"`              // 工作经验
	Content      string    `gorm:"column:content;not null;type:text" json:"content"`                         // 内容
	UserId       int64     `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`                 // 用户ID
	ReviewStatus int32     `gorm:"column:review_status;not null;default:0" json:"review_status"`             // 审核状态
	ViewNum      int64     `gorm:"column:view_num;not null;default:0" json:"view_num"`                       // 浏览次数
	ThumbNum     int64     `gorm:"column:thumb_num;not null;default:0" json:"thumb_num"`                     // 点赞次数
	HasThumb     bool      `gorm:"column:has_thumb;not null;default:false" json:"has_thumb"`                 // 是否已点赞
	CreateTime   time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"` // 更新时间
	IsDelete     bool      `gorm:"column:is_delete;not null;default:false" json:"is_delete"`                 // 是否删除
	_            struct {
		Education string
		Major     string
		School    string
		WorkExp   string
		UserID    int64
	} `gorm:"index:idx_education_major_school"`
}

type WorkExp struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Company      string    `gorm:"column:company;not null;index:idx_company" json:"company"`                   // 公司名称
	WorkDuration string    `gorm:"column:work_duration;not null;index:idx_work_duration" json:"work_duration"` // 工作时长
	Region       string    `gorm:"column:region;not null;index:idx_region" json:"region"`                      // 地区
	Content      string    `gorm:"column:content;not null;type:text" json:"content"`                           // 内容
	UserId       int64     `gorm:"column:user_id;not null;index:idx_user_id" json:"user_id"`                   // 用户ID
	ReviewStatus int32     `gorm:"column:review_status;not null;default:0" json:"review_status"`               // 审核状态
	ViewNum      int64     `gorm:"column:view_num;not null;default:0" json:"view_num"`                         // 浏览次数
	ThumbNum     int64     `gorm:"column:thumb_num;not null;default:0" json:"thumb_num"`                       // 点赞次数
	HasThumb     bool      `gorm:"column:has_thumb;not null;default:false" json:"has_thumb"`                   // 是否已点赞
	CreateTime   time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`   // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"`   // 更新时间
	IsDelete     bool      `gorm:"column:is_delete;not null;default:false" json:"is_delete"`                   // 是否删除
	_            struct {
		Company      string
		Position     string
		Department   string
		WorkDuration string
		Region       string
		UserID       int64
	} `gorm:"index:idx_company_position_department_region"`
}
