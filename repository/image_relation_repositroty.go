package repository

import (
	"coderhub/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ImageRelationRepository interface {
	Create(ctx context.Context, imageRelation *model.ImageRelation) error
	BatchCreate(ctx context.Context, imageRelations []*model.ImageRelation) error
	BatchDelete(ctx context.Context, ids []int64) error
	BatchGetImagesByEntity(ctx context.Context, entityIds []int64, entityType string) ([]model.ImageRelation, error)
	Delete(ctx context.Context, imageID int64, entityID int64, entityType string) error
	ListByEntityID(ctx context.Context, entityID int64, entityType string) ([]model.ImageRelation, error)
	ListByImageID(ctx context.Context, imageID int64) ([]model.ImageRelation, error)
}

type imageRelationRepository struct {
	DB *gorm.DB
}

func NewImageRelationRepository(db *gorm.DB) ImageRelationRepository {
	return &imageRelationRepository{DB: db}
}

func (r *imageRelationRepository) Create(ctx context.Context, imageRelation *model.ImageRelation) error {
	return r.DB.WithContext(ctx).Create(imageRelation).Error
}

func (r *imageRelationRepository) BatchCreate(ctx context.Context, imageRelations []*model.ImageRelation) error {
	if len(imageRelations) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&imageRelations).Error
}

// BatchDelete 批量删除图片关联, 根据ID列表删除
func (r *imageRelationRepository) BatchDelete(ctx context.Context, ids []int64) error {
	return r.DB.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.ImageRelation{}).Error
}

// BatchGetImagesByEntity 批量获取图片关联，根据实体ID列表、实体类型列表获取
func (r *imageRelationRepository) BatchGetImagesByEntity(ctx context.Context, entityIds []int64, entityType string) ([]model.ImageRelation, error) {
	var imageRelations []model.ImageRelation
	
	// 添加SQL查询日志
	query := r.DB.WithContext(ctx).
		Where("entity_id IN (?) AND entity_type = ?", entityIds, entityType)
	
	// 打印SQL语句
	sql := query.Statement.SQL.String()
	fmt.Printf("执行的SQL: %s\n参数: entityIds=%v, entityType=%s\n", sql, entityIds, entityType)
	
	err := query.Find(&imageRelations).Error
	if err != nil {
		return nil, err
	}
	
	// 打印查询结果
	fmt.Printf("查询到的记录数: %d\n", len(imageRelations))
	for i, rel := range imageRelations {
		fmt.Printf("记录 %d: entity_id=%d, image_id=%d\n", i+1, rel.EntityID, rel.ImageID)
	}
	
	return imageRelations, nil
}

// DeleteByEntityID 批量删除关联，根据实体ID、实体类型删除
func (r *imageRelationRepository) DeleteByEntityID(ctx context.Context, entityID int64, entityType string) error {
	return r.DB.WithContext(ctx).Where("entity_id = ? AND entity_type = ?", entityID, entityType).Delete(&model.ImageRelation{}).Error
}

// Delete 删除图片关联, 根据图片ID、实体ID、实体类型删除
func (r *imageRelationRepository) Delete(ctx context.Context, imageID int64, entityID int64, entityType string) error {
	return r.DB.WithContext(ctx).Where("image_id = ? AND entity_id = ? AND entity_type = ?", imageID, entityID, entityType).Delete(&model.ImageRelation{}).Error
}

func (r *imageRelationRepository) ListByEntityID(ctx context.Context, entityID int64, entityType string) ([]model.ImageRelation, error) {
	var imageRelations []model.ImageRelation
	return imageRelations, r.DB.WithContext(ctx).Where("entity_id = ? AND entity_type = ?", entityID, entityType).Find(&imageRelations).Error
}

func (r *imageRelationRepository) ListByImageID(ctx context.Context, imageID int64) ([]model.ImageRelation, error) {
	var imageRelations []model.ImageRelation
	return imageRelations, r.DB.WithContext(ctx).Where("image_id = ?", imageID).Find(&imageRelations).Error
}
