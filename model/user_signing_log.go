package model

import (
	"time"
)

type UserSigningLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	UserID    int64     `gorm:"not null;index:idx_user_sign_date,unique;column:user_id"`             // 用户 ID
	SignDate  time.Time `gorm:"type:date;not null;index:idx_user_sign_date,unique;column:sign_date"` // 签到日期
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`                                    // 签到时间
}
