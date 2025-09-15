package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, emotion *model.Emoticon) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, emotion *model.Emoticon) error
	GetByID(ctx context.Context, id int64) (*model.Emoticon, error)
	List(ctx context.Context, page, pageSize int64) ([]*model.Emoticon, int64, error)
}

type repository struct {
	DB *gorm.DB
}

func NewEmoticonRepository(db *gorm.DB) Repository {
	return &repository{DB: db}
}

func (r *repository) Create(ctx context.Context, emotion *model.Emoticon) error {
	return r.DB.WithContext(ctx).Create(emotion).Error
}
func (r *repository) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.Emoticon{}, id).Error
}

func (r *repository) Update(ctx context.Context, emotion *model.Emoticon) error {
	return r.DB.WithContext(ctx).Updates(emotion).Error
}

func (r *repository) GetByID(ctx context.Context, id int64) (*model.Emoticon, error) {
	var emoticon model.Emoticon
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&emoticon).Error
	if err != nil {
		return nil, err
	}
	return &emoticon, nil
}

func (r *repository) List(ctx context.Context, page, pageSize int64) ([]*model.Emoticon, int64, error) {
	var emoticons []*model.Emoticon
	var total int64
	err := r.DB.WithContext(ctx).Model(&model.Emoticon{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.DB.WithContext(ctx).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&emoticons).Error
	if err != nil {
		return nil, 0, err
	}
	return emoticons, total, nil
}
