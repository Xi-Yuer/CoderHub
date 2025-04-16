package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserFollowRepository interface {
	// CreateUserFollow 创建用户关注关系
	CreateUserFollow(userFollow *model.UserFollow) error
	// DeleteUserFollow 删除用户关注关系
	DeleteUserFollow(userFollow *model.UserFollow) error
	// GetUserFollows 查询用户关注的所有用户
	GetUserFollows(followerID int64, page int32, pageSize int32) ([]*model.UserFollow, error)
	// BatchGetUserFollows 批量查询用户关注的所有用户
	BatchGetUserFollows(followerID int64, page int32, pageSize int32) ([]*model.UserFollow, error)
	// GetUserFans 查询某用户的粉丝列表
	GetUserFans(followedID int64, page int32, pageSize int32) ([]*model.UserFollow, error)
	// GetUserFansCount 查询某用户的粉丝数量
	GetUserFansCount(followedID int64) (int64, error)
	// BatchGetUserFans 批量查询某用户的粉丝列表
	BatchGetUserFans(followedID int64, page int32, pageSize int32) ([]*model.UserFollow, error)
	// IsUserFollowed 判断两个用户是否存在关注关系
	IsUserFollowed(followerID int64, followedID int64) (bool, error)
	// GetMutualFollows 查询互相关注的用户
	GetMutualFollows(userID int64, page int32, pageSize int32) ([]*model.UserFollow, error)
	// GetUserAndUserHasFollow 查询用户和用户是否已经关注
	GetUserAndUserHasFollow(userID int64, userIDs []int64) ([]int64, error)
}

func NewUserFollowRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *UserFollowRepositoryImpl {
	return &UserFollowRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

type UserFollowRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

// CreateUserFollow 创建用户关注关系
func (r *UserFollowRepositoryImpl) CreateUserFollow(userFollow *model.UserFollow) error {
	isFollowed, err := r.IsUserFollowed(userFollow.FollowerID, userFollow.FollowedID)
	if err != nil {
		return err
	}
	if isFollowed {
		return r.DeleteUserFollow(userFollow)
	}
	if userFollow.FollowerID == userFollow.FollowedID {
		return errors.New("不能关注自己")
	}

	return r.DB.Model(&model.UserFollow{}).Create(userFollow).Error
}

// DeleteUserFollow 删除用户关注关系
func (r *UserFollowRepositoryImpl) DeleteUserFollow(userFollow *model.UserFollow) error {
	if userFollow.FollowerID == userFollow.FollowedID {
		return errors.New("不能取消关注自己")
	}

	// 开启事务
	tx := r.DB.Begin()
	err := tx.Model(&model.UserFollow{}).
		Where("follower_id = ? AND followed_id = ?", userFollow.FollowerID, userFollow.FollowedID).
		Unscoped().
		Delete(userFollow).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除缓存
	key := fmt.Sprintf("follow:%d:%d", userFollow.FollowerID, userFollow.FollowedID)
	if err := r.Redis.Del(key); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return tx.Commit().Error
}

// GetUserFollows 查询用户关注的所有用户
func (r *UserFollowRepositoryImpl) GetUserFollows(followerID int64, page int32, pageSize int32) ([]*model.UserFollow, error) {
	var userFollows []*model.UserFollow
	// 使用 idx_follower_followed 索引
	return userFollows, r.DB.Model(&model.UserFollow{}).
		Select("id, follower_id, followed_id, created_at, updated_at, deleted_at").
		Where("follower_id = ?", followerID).
		Order("created_at DESC").
		Offset((int(page) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&userFollows).Error
}

// GetUserFans 查询某用户的粉丝列表
func (r *UserFollowRepositoryImpl) GetUserFans(followedID int64, page int32, pageSize int32) ([]*model.UserFollow, error) {
	var userFollows []*model.UserFollow
	if page < 1 {
		page = 1
	}

	// 使用 idx_followed_follower 索引
	return userFollows, r.DB.Model(&model.UserFollow{}).
		Select("id, follower_id, followed_id, created_at, updated_at, deleted_at").
		Where("followed_id = ?", followedID).
		Order("created_at DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&userFollows).Error
}

// IsUserFollowed 判断两个用户是否存在关注关系
func (r *UserFollowRepositoryImpl) IsUserFollowed(followerID int64, followedID int64) (bool, error) {
	// 先查缓存
	key := fmt.Sprintf("follow:%d:%d", followerID, followedID)
	exists, err := r.Redis.Exists(key)
	if err == nil && exists {
		return true, nil
	}

	// 使用 uk_follower_followed 唯一索引
	var count int64
	err = r.DB.Model(&model.UserFollow{}).
		Select("1").
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).
		Count(&count).Error

	exists = count > 0
	// 如果存在关系，写入缓存并设置过期时间
	if exists {
		if err = r.Redis.Set(key, "1"); err != nil {
			return exists, fmt.Errorf("failed to set cache: %w", err)
		}
	}
	return exists, err
}

// GetMutualFollows 查询互相关注的用户
func (r *UserFollowRepositoryImpl) GetMutualFollows(userID int64, page int32, pageSize int32) ([]*model.UserFollow, error) {
	var userFollows []*model.UserFollow

	// 使用子查询和索引优化互关查询
	return userFollows, r.DB.Model(&model.UserFollow{}).
		Select("f1.*").
		Joins("AS f1").
		Joins("JOIN user_follows f2 ON f1.follower_id = f2.followed_id AND f1.followed_id = f2.follower_id").
		Where("f1.follower_id = ?", userID).
		Order("f1.created_at DESC").
		Offset((int(page) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&userFollows).Error
}

// GetUserAndUserHasFollow 查询用户和用户是否已经关注
func (r *UserFollowRepositoryImpl) GetUserAndUserHasFollow(userID int64, userIDs []int64) ([]int64, error) {
	if len(userIDs) == 0 {
		return []int64{}, nil
	}

	var followedUserIDs []int64
	// 使用 idx_follower_followed 索引
	err := r.DB.Model(&model.UserFollow{}).
		Select("followed_id").
		Where("follower_id = ? AND followed_id IN ?", userID, userIDs).
		Pluck("followed_id", &followedUserIDs).Error

	return followedUserIDs, err
}

// BatchGetUserFollows 批量查询用户关注的所有用户
func (r *UserFollowRepositoryImpl) BatchGetUserFollows(followerID int64, page int32, pageSize int32) ([]*model.UserFollow, error) {
	var userFollows []*model.UserFollow

	// 使用 idx_follower_followed 索引
	err := r.DB.Model(&model.UserFollow{}).
		Select("id, follower_id, followed_id, created_at, updated_at, deleted_at").
		Where("follower_id = ?", followerID).
		Order("created_at DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&userFollows).Error

	if err != nil {
		return nil, err
	}

	// 如果结果为空，返回空切片而不是 nil
	if len(userFollows) == 0 {
		return []*model.UserFollow{}, nil
	}

	return userFollows, nil
}

// BatchGetUserFans 批量查询某用户的粉丝列表
func (r *UserFollowRepositoryImpl) BatchGetUserFans(followedID int64, page int32, pageSize int32) ([]*model.UserFollow, error) {
	var userFollows []*model.UserFollow

	// 使用 idx_followed_follower 索引
	err := r.DB.Model(&model.UserFollow{}).
		Select("id, follower_id, followed_id, created_at, updated_at, deleted_at").
		Where("followed_id = ?", followedID).
		Order("created_at DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&userFollows).Error

	if err != nil {
		return nil, err
	}

	// 如果结果为空，返回空切片而不是 nil
	if len(userFollows) == 0 {
		return []*model.UserFollow{}, nil
	}

	return userFollows, nil
}

// GetUserFansCount 查询某用户的粉丝数量
func (r *UserFollowRepositoryImpl) GetUserFansCount(followedID int64) (int64, error) {
	var count int64
	// 使用 idx_followed_follower 索引
	err := r.DB.Model(&model.UserFollow{}).
		Where("followed_id =?", followedID).
		Count(&count).Error
	return count, err
}
