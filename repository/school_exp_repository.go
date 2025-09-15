package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type SchoolExpRepository interface {
	Create(ctx context.Context, schoolExp *model.SchoolExp) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, schoolExp *model.SchoolExp, page, pageSize int64) ([]*model.SchoolExp, int64, error)
}

type SchoolExpRepositoryImpl struct {
	DB *gorm.DB
}

func NewSchoolExpRepository(db *gorm.DB) *SchoolExpRepositoryImpl {
	return &SchoolExpRepositoryImpl{DB: db}
}

func (r *SchoolExpRepositoryImpl) Create(ctx context.Context, schoolExp *model.SchoolExp) error {
	return r.DB.WithContext(ctx).Create(schoolExp).Error
}

func (r *SchoolExpRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.SchoolExp{}, id).Error
}
func (r *SchoolExpRepositoryImpl) List(ctx context.Context, schoolExp *model.SchoolExp, page, pageSize int64) ([]*model.SchoolExp, int64, error) {
	var schoolExps []*model.SchoolExp
	var total int64
	err := r.DB.WithContext(ctx).Model(&model.SchoolExp{}).Where(schoolExp).Count(&total).Error // 查询总数
	if err != nil {
		return nil, 0, err
	}
	if total > 0 {
		err = r.DB.WithContext(ctx).Model(&model.SchoolExp{}).
			Where(schoolExp).
			Offset(int((page - 1) * pageSize)).
			Limit(int(pageSize)).
			Find(&schoolExps).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return schoolExps, total, nil
}
