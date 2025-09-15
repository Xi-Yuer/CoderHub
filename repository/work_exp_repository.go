package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type WorkExpRepository interface {
	Create(ctx context.Context, workExp *model.WorkExp) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, workExp *model.WorkExp, page, pageSize int64) ([]*model.WorkExp, int64, error)
}

type WorkExpRepositoryImpl struct {
	DB *gorm.DB
}

func NewWorkExpRepository(db *gorm.DB) *WorkExpRepositoryImpl {
	return &WorkExpRepositoryImpl{DB: db}
}

func (r *WorkExpRepositoryImpl) Create(ctx context.Context, workExp *model.WorkExp) error {
	return r.DB.WithContext(ctx).Create(workExp).Error
}

func (r *WorkExpRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.WorkExp{}, id).Error
}
func (r *WorkExpRepositoryImpl) List(ctx context.Context, workExp *model.WorkExp, page, pageSize int64) ([]*model.WorkExp, int64, error) {
	var schoolExps []*model.WorkExp
	var total int64
	err := r.DB.WithContext(ctx).Model(&model.WorkExp{}).Where(workExp).Count(&total).Error // 获取总数
	if err != nil {
		return nil, 0, err
	}
	if total > 0 {
		err = r.DB.WithContext(ctx).Model(&model.WorkExp{}).
			Where(workExp).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&schoolExps).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return schoolExps, total, nil
}
