package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserFavorEntityRepository interface {
	Create(ctx context.Context, userFavorEntity *model.UserFavor) error
	Delete(ctx context.Context, userFavorEntity *model.UserFavor) error
	GetList(ctx context.Context, userFavorEntity *model.UserFavor, page, pageSize int) ([]*model.UserFavor, int64, error)
	BatchGetUserFavorEntity(ctx context.Context, entityIDs []int64, userID int64) (map[int64]bool, error)
}

type UserFavorEntityRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func NewUserFavorEntityRepository(db *gorm.DB, rdb storage.RedisDB) *UserFavorEntityRepositoryImpl {
	return &UserFavorEntityRepositoryImpl{DB: db, Redis: rdb}
}

func (r *UserFavorEntityRepositoryImpl) Create(ctx context.Context, userFavorEntity *model.UserFavor) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 更新收藏夹数量
		update := tx.Model(&model.UserFavorFolder{}).Where("id = ?", userFavorEntity.FavorFoldId).Update("favor_num", gorm.Expr("favor_num + ?", 1))
		if update.Error != nil {
			return update.Error
		}
		// 创建收藏记录
		return tx.WithContext(ctx).Create(userFavorEntity).Error
	})
}

func (r *UserFavorEntityRepositoryImpl) Delete(ctx context.Context, userFavorEntity *model.UserFavor) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.UserFavorFolder{}).Where("id = ?", userFavorEntity.FavorFoldId).Update("favor_num", gorm.Expr("favor_num - ?", 1)).Error
		if err != nil {
			return err
		}
		return r.DB.WithContext(ctx).Where("id = ? user_id = ? AND favor_fold_id = ? AND entity_id = ? AND entity_type = ?", userFavorEntity.ID, userFavorEntity.UserId, userFavorEntity.FavorFoldId, userFavorEntity.EntityId, userFavorEntity.EntityType).Delete(userFavorEntity).Error
	})
}

// GetList 获取收藏夹内容
func (r *UserFavorEntityRepositoryImpl) GetList(ctx context.Context, userFavorEntity *model.UserFavor, page, pageSize int) ([]*model.UserFavor, int64, error) {
	var userFavorEntities []*model.UserFavor
	var total int64
	fmt.Println("userFavorEntity.UserId", userFavorEntity.UserId)
	fmt.Println("userFavorEntity.FavorFoldId", userFavorEntity.FavorFoldId)
	fmt.Println("userFavorEntity.EntityType", userFavorEntity.EntityType)
	err := r.DB.WithContext(ctx).Model(&model.UserFavor{}).Where("user_id = ? AND favor_fold_id = ? AND entity_type = ?", userFavorEntity.UserId, userFavorEntity.FavorFoldId, userFavorEntity.EntityType).Limit(pageSize).Offset((page - 1) * pageSize).Count(&total).Find(&userFavorEntities).Error
	return userFavorEntities, total, err
}

// BatchGetUserFavorEntity 批量获取实体是否被用户收藏
func (r *UserFavorEntityRepositoryImpl) BatchGetUserFavorEntity(ctx context.Context, entityIDs []int64, userID int64) (map[int64]bool, error) {
	var userFavorEntities []*model.UserFavor
	err := r.DB.WithContext(ctx).Model(&model.UserFavor{}).Where("user_id = ? AND entity_id IN ?", userID, entityIDs).Find(&userFavorEntities).Error
	if err != nil {
		return nil, err
	}
	userFavorEntityMap := make(map[int64]bool)
	for _, userFavorEntity := range userFavorEntities {
		userFavorEntityMap[userFavorEntity.EntityId] = true
	}
	return userFavorEntityMap, nil
}
