package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"gorm.io/gorm"
)

type UserFavorFolderRepository interface {
	Create(ctx context.Context, userFavorFolder *model.UserFavorFolder) error
	Delete(ctx context.Context, userFavorFolder *model.UserFavorFolder) error
	Update(ctx context.Context, userFavorFolder *model.UserFavorFolder) error
	GetFolderByID(ctx context.Context, id int64) (*model.UserFavorFolder, error)
	UpdateFolderNum(ctx context.Context, id int64, num int64) error
	GetList(ctx context.Context, userID int64, requestUserId, page, pageSize int64) ([]*model.UserFavorFolder, int64, error)
}

func NewUserFavorFolderRepository(db *gorm.DB, rdb storage.RedisDB) *UserFavorFolderRepositoryImpl {
	return &UserFavorFolderRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

type UserFavorFolderRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func (r *UserFavorFolderRepositoryImpl) Create(ctx context.Context, userFavorFolder *model.UserFavorFolder) error {
	return r.DB.WithContext(ctx).Create(userFavorFolder).Error
}

func (r *UserFavorFolderRepositoryImpl) Delete(ctx context.Context, userFavorFolder *model.UserFavorFolder) error {
	return r.DB.WithContext(ctx).Where("id = ?", userFavorFolder.ID).Delete(userFavorFolder).Error
}

func (r *UserFavorFolderRepositoryImpl) Update(ctx context.Context, userFavorFolder *model.UserFavorFolder) error {
	return r.DB.WithContext(ctx).Where("id = ?", userFavorFolder.ID).Updates(userFavorFolder).Error
}

func (r *UserFavorFolderRepositoryImpl) UpdateFolderNum(ctx context.Context, id int64, num int64) error {
	return r.DB.WithContext(ctx).Model(&model.UserFavorFolder{}).Where("id = ?", id).UpdateColumn("favor_num", gorm.Expr("favor_num + ?", num)).Error
}

func (r *UserFavorFolderRepositoryImpl) GetFolderByID(ctx context.Context, id int64) (*model.UserFavorFolder, error) {
	var userFavorFolder model.UserFavorFolder
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&userFavorFolder).Error
	return &userFavorFolder, err
}

func (r *UserFavorFolderRepositoryImpl) GetList(ctx context.Context, userID int64, requestUserId, page, pageSize int64) ([]*model.UserFavorFolder, int64, error) {
	var userFavorFolders []*model.UserFavorFolder
	var count int64
	// 请求的用户不是自己的话，只能查看公开收藏夹
	if requestUserId != userID {
		err := r.DB.WithContext(ctx).Where("user_id = ? AND is_public = ?", userID, true).Order("created_at desc").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&userFavorFolders).Count(&count).Error
		return userFavorFolders, count, err
	} else {
		// 请求用户是自己的话，可以查看所有收藏夹
		err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&userFavorFolders).Count(&count).Error
		return userFavorFolders, count, err
	}
}
