package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"sync"

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

	// 构建基础查询
	query := r.DB.WithContext(ctx).Model(&model.UserFavorFolder{})

	// 根据请求用户设置查询条件
	if requestUserId != userID {
		query = query.Where("user_id = ? AND is_public = ?", userID, true)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	// 并发执行计数查询和数据查询
	var wg sync.WaitGroup
	var countErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 创建新的查询会话，避免影响主查询
		countErr = query.Session(&gorm.Session{}).
			Select("COUNT(1)").
			Count(&count).Error
	}()

	// 数据查询
	err := query.Order("created_at desc").
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&userFavorFolders).Error

	wg.Wait()

	if countErr != nil {
		return nil, 0, countErr
	}
	if err != nil {
		return nil, 0, err
	}

	return userFavorFolders, count, nil
}
