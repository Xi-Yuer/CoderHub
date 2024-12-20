package repository

import (
	"coderhub/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ImageRepository 图片仓储接口
type ImageRepository interface {
	// Create 创建图片元数据记录
	Create(ctx context.Context, image *model.Image) error
	// BatchCreate 批量创建图片元数据记录
	BatchCreate(ctx context.Context, images []*model.Image) error
	// GetByID 获取图片元数据
	GetByID(ctx context.Context, id int64) (*model.Image, error)
	// Delete 删除图片元数据（软删除）
	Delete(ctx context.Context, id int64) error
	// ListByUserID 获取用户上传的图片列表
	ListByUserID(ctx context.Context, userID int64, page, pageSize int32) ([]model.Image, int64, error)
	// BatchGetImagesByID 批量获取图片关联，根据实体ID列表、实体类型列表获取
	BatchGetImagesByID(ctx context.Context, ids []int64) ([]model.Image, error)
}

var (
	ErrImageNotFound = errors.New("图片未找到")
)

type imageRepository struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{
		DB: db,
	}
}

func (r *imageRepository) Create(ctx context.Context, image *model.Image) error {
	return r.DB.WithContext(ctx).Create(image).Error
}

func (r *imageRepository) BatchCreate(ctx context.Context, images []*model.Image) error {
	if len(images) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&images).Error
}

func (r *imageRepository) GetByID(ctx context.Context, id int64) (*model.Image, error) {
	var image model.Image
	err := r.DB.WithContext(ctx).First(&image, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *imageRepository) Delete(ctx context.Context, id int64) error {
	result := r.DB.WithContext(ctx).Delete(&model.Image{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrImageNotFound
	}
	return result.Error
}

func (r *imageRepository) ListByUserID(ctx context.Context, userID int64, page, pageSize int32) ([]model.Image, int64, error) {
	var images []model.Image
	var total int64
	err := r.DB.WithContext(ctx).
		Model(&model.Image{}).
		Where("user_id = ?", userID).
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&images).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return images, total, nil
}

func (r *imageRepository) BatchGetImagesByID(ctx context.Context, ids []int64) ([]model.Image, error) {
	var images []model.Image
	return images, r.DB.WithContext(ctx).Where("id IN (?)", ids).Find(&images).Error
}
