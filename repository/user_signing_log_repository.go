package repository

import (
	"coderhub/model"
	"coderhub/shared/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// UserSigninLogRepository 定义用户签到记录仓库接口
type UserSigninLogRepository interface {
	// AddSigninLog 添加用户签到记录
	AddSigninLog(ctx context.Context, userID int64, signinDay time.Time) error
	// GetUserSigninStatusInMonth 获取用户某个月内的签到情况
	GetUserSigninStatusInMonth(ctx context.Context, userID int64, year int, month time.Month) ([]bool, error)
}

// UserSigninLogRepositoryImpl 实现用户签到记录仓库接口
type UserSigninLogRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserSigninLogRepository 创建用户签到记录仓库实例
func NewUserSigninLogRepository(db *gorm.DB) *UserSigninLogRepositoryImpl {
	return &UserSigninLogRepositoryImpl{
		DB: db,
	}
}

// AddSigninLog 添加用户签到记录
func (r *UserSigninLogRepositoryImpl) AddSigninLog(ctx context.Context, userID int64, signinDay time.Time) error {
	log := model.UserSigningLog{
		ID:        utils.GenID(),
		UserID:    userID,
		SignDate:  signinDay,
		CreatedAt: time.Now(),
	}
	result := r.DB.WithContext(ctx).Create(&log)
	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			return fmt.Errorf("用户 %d 在 %s 已签到", userID, signinDay.Format("2006-01-02"))
		}
		return result.Error
	}
	return nil
}

// GetUserSigninStatusInMonth 获取用户某个月内的签到情况
func (r *UserSigninLogRepositoryImpl) GetUserSigninStatusInMonth(ctx context.Context, userID int64, year int, month time.Month) ([]bool, error) {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	firstDayStr := firstDay.Format("2006-01-02")
	lastDayStr := lastDay.Format("2006-01-02")

	var logs []model.UserSigningLog
	result := r.DB.WithContext(ctx).Where("user_id = ? AND sign_date BETWEEN ? AND ?", userID, firstDayStr, lastDayStr).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	daysInMonth := lastDay.Day()
	signinStatus := make([]bool, daysInMonth)

	for _, log := range logs {
		day := log.SignDate.Day() - 1
		if day >= 0 && day < daysInMonth {
			signinStatus[day] = true
		}
	}

	return signinStatus, nil
}
