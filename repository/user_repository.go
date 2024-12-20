package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	FindOneByEmail(email string) (*model.User, error)
	BatchGetUserByID(ids []int64) ([]*model.User, error)
	UpdateUser(user *model.User) error
	ResetPassword(email string, password string) error
	DeleteUser(id int64) error
}

func NewUserRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

type UserRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	// 使用事务确保数据一致性
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// 创建后设置缓存
		return r.setCache(user.CacheKeyByID(user.ID), user)
	})
}

func (r *UserRepositoryImpl) GetUserByName(name string) (*model.User, error) {
	var user model.User
	key := user.CacheKeyByName(name)

	// 尝试从缓存获取
	if cached, err := r.getCache(key); err == nil {
		return cached, nil
	}

	// 简化数据库查询
	if err := r.DB.First(&user, "user_name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %s", name)
		}
		return nil, err
	}

	// 异步设置缓存
	go func() {
		_ = r.setCache(key, &user)
		// 同时设置ID缓存
		_ = r.setCache(user.CacheKeyByID(user.ID), &user)
	}()

	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	key := user.CacheKeyByID(id)

	// 尝试从缓存获取
	if cached, err := r.getCache(key); err == nil {
		return cached, nil
	}

	// 简化数据库查询
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %d", id)
		}
		return nil, err
	}

	// 异步设置缓存
	go func() {
		_ = r.setCache(key, &user)
		// 同时设置用户名缓存
		_ = r.setCache(user.CacheKeyByName(user.UserName), &user)
	}()

	return &user, nil
}

func (r *UserRepositoryImpl) FindOneByEmail(email string) (*model.User, error) {
	var user model.User
	return &user, r.DB.Where("email = ?", email).First(&user).Error
}

func (r *UserRepositoryImpl) BatchGetUserByID(ids []int64) ([]*model.User, error) {
	var users []*model.User
	fmt.Println("repository_ids", ids)
	err := r.DB.Where("id IN (?)", ids).Find(&users).Error
	fmt.Println("repository_users_length", len(users))
	return users, err
}

func (r *UserRepositoryImpl) UpdateUser(user *model.User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 获取旧数据用于清理缓存
		var oldUser model.User
		if err := tx.First(&oldUser, user.ID).Error; err != nil {
			return fmt.Errorf("获取用户失败: %w", err)
		}

		updates := tx.Model(&model.User{}).Where("id = ?", user.ID).Updates(user)
		if err := updates.Error; err != nil {
			return err
		}
		// 清理所有相关缓存
		keys := []string{
			user.CacheKeyByID(user.ID),
			user.CacheKeyByName(user.UserName),
		}

		for _, key := range keys {
			fmt.Println("delCache", key)
			_ = r.delCache(key)
		}

		return nil
	})
}

func (r *UserRepositoryImpl) ResetPassword(email string, password string) error {
	// 获取旧数据用于清理缓存
	var oldUser model.User
	if err := r.DB.First(&oldUser, email).Error; err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	// 清理所有相关缓存
	keys := []string{
		oldUser.CacheKeyByID(oldUser.ID),
		oldUser.CacheKeyByName(oldUser.UserName),
	}

	for _, key := range keys {
		fmt.Println("delCache", key)
		_ = r.delCache(key)
	}
	return r.DB.Model(&model.User{}).Where("email = ?", email).Update("password", password).Error
}

func (r *UserRepositoryImpl) DeleteUser(id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 获取用户信息用于清理缓存
		var user model.User
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		// 清理所有相关缓存
		keys := []string{
			user.CacheKeyByID(id),
			user.CacheKeyByName(user.UserName),
		}

		for _, key := range keys {
			if err := r.delCache(key); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *UserRepositoryImpl) getCache(key string) (*model.User, error) {
	var user model.User
	data, err := r.Redis.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) setCache(key string, user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.Redis.Set(key, string(data))
}

func (r *UserRepositoryImpl) delCache(key string) error {
	return r.Redis.Del(key)
}
