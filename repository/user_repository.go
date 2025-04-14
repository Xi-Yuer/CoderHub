package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	GetUserBirthday(id int64) (int64, error)
	FindOneByEmail(email string) (*model.User, error)
	BatchGetUserByID(ids []int64) ([]*model.User, error)
	UpdateUser(user *model.User) error
	IncrUserLevel(id int64, level int32) error
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
	// 查询用户粉丝数量
	var followerCount int64
	if err := r.DB.Table("user_follows").Where("followed_id = ?", user.ID).Count(&followerCount).Error; err != nil {
		return nil, err
	}

	// 查询用户关注数量
	var followCount int64
	if err := r.DB.Table("user_follows").Where("follower_id = ?", user.ID).Count(&followCount).Error; err != nil {
		return nil, err
	}

	// 查询用户文章数量
	var articleCount int64
	if err := r.DB.Model(&model.Articles{}).
		Where("author_id = ? AND deleted_at IS NULL AND status = ?", user.ID, "published").
		Count(&articleCount).Error; err != nil {
		return nil, err
	}
	user.FollowerCount = followerCount
	user.FollowCount = followCount
	user.ArticleCount = articleCount

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

	// 并行查询用户粉丝数量、关注数量和文章数量
	var (
		followerCount int64
		followCount   int64
		articleCount  int64
		errFollower   error
		errFollow     error
		errArticle    error
		done          = make(chan bool)
	)

	go func() {
		errFollower = r.DB.Table("user_follows").Where("followed_id = ?", id).Count(&followerCount).Error
		done <- true
	}()

	go func() {
		errFollow = r.DB.Table("user_follows").Where("follower_id = ?", id).Count(&followCount).Error
		done <- true
	}()

	go func() {
		errArticle = r.DB.Model(&model.Articles{}).
			Where("author_id = ? AND deleted_at IS NULL AND status = ?", id, "published").
			Count(&articleCount).Error
		done <- true
	}()

	// 等待所有 goroutine 完成
	for i := 0; i < 3; i++ {
		<-done
	}

	// 检查错误
	if errFollower != nil {
		return nil, errFollower
	}
	if errFollow != nil {
		return nil, errFollow
	}
	if errArticle != nil {
		return nil, errArticle
	}

	user.FollowerCount = followerCount
	user.FollowCount = followCount
	user.ArticleCount = articleCount

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
	err := r.DB.Where("id IN (?)", ids).Find(&users).Error
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
			_ = r.delCache(key)
		}

		return nil
	})
}

func (r *UserRepositoryImpl) IncrUserLevel(id int64, level int32) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 获取旧数据用于清理缓存
		var oldUser model.User
		if err := tx.First(&oldUser, id).Error; err != nil {
			return fmt.Errorf("获取用户失败: %w", err)
		}
		// 清理所有相关缓存
		keys := []string{
			oldUser.CacheKeyByID(oldUser.ID),
			oldUser.CacheKeyByName(oldUser.UserName),
		}
		for _, key := range keys {
			_ = r.delCache(key)
		}
		// 更新用户等级,原来的基础上加数值
		if err := tx.Model(&model.User{}).Where("id = ?", id).Update("level", gorm.Expr("level + ?", level)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *UserRepositoryImpl) ResetPassword(email string, password string) error {
	// 获取旧数据用于清理缓存
	var oldUser model.User
	if err := r.DB.Where("email = ?", email).First(&oldUser).Error; err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	// 清理所有相关缓存
	keys := []string{
		oldUser.CacheKeyByID(oldUser.ID),
		oldUser.CacheKeyByName(oldUser.UserName),
	}

	for _, key := range keys {
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

func (r *UserRepositoryImpl) GetUserBirthday(userID int64) (int64, error) {
	var createdAt time.Time
	now := time.Now()

	// 查询用户注册时间
	if err := r.DB.Table("users").Where("id = ?", userID).Pluck("created_at", &createdAt).Error; err != nil {
		return 0, err
	}

	// 计算时间差（单位：秒）
	return (now.Unix() - createdAt.Unix()) / 86400, nil
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
