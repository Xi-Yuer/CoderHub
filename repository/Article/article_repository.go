package repository

import (
	"coderhub/model"
	"coderhub/shared/CacheDB"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	CreateArticle(article *model.Articles) error
	GetArticleByID(id int64) (*model.Articles, error)
	UpdateArticle(article *model.Articles) error
	DeleteArticle(id int64) error
}
type ArticleRepositoryImpl struct {
	DB    *gorm.DB
	Redis CacheDB.RedisDB
}

func NewArticleRepositoryImpl(db *gorm.DB, rdb CacheDB.RedisDB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

func (r *ArticleRepositoryImpl) CreateArticle(article *model.Articles) error {
	if err := r.DB.Create(article).Error; err != nil {
		return err
	}
	// 创建后设置缓存
	return r.setCache(article.CacheKeyByID(article.ID), article)
}

func (r *ArticleRepositoryImpl) GetArticleByID(id int64) (*model.Articles, error) {
	var article model.Articles
	key := article.CacheKeyByID(id)

	// 尝试从缓存获取
	if cached, err := r.getCache(key); err == nil {
		return cached, nil
	}

	// 简化数据库查询
	if err := r.DB.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文章不存在: %v", id)
		}
		return nil, err
	}

	// 异步设置缓存，避免影响主流程
	go func() {
		_ = r.setCache(key, &article)
	}()

	return &article, nil
}

func (r *ArticleRepositoryImpl) UpdateArticle(article *model.Articles) error {
	if article.ID <= 0 {
		return errors.New("无效的文章ID")
	}

	// 使用事务确保数据一致性
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 检查文章是否存在
		var count int64
		if err := tx.Model(&model.Articles{}).Select("id").Where("id = ?", article.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return errors.New("文章不存在")
		}

		// 更新数据库
		if err := tx.Model(article).Updates(article).Error; err != nil {
			return err
		}

		// 删除旧缓存并设置新缓存
		key := article.CacheKeyByID(article.ID)
		if err := r.delCache(key); err != nil {
			return err
		}
		return r.setCache(key, article)
	})
}

func (r *ArticleRepositoryImpl) DeleteArticle(id int64) error {
	article, err := r.GetArticleByID(id)
	if err != nil {
		return err
	}
	err = r.delCache(article.CacheKeyByID(id))
	if err != nil {
		return err
	}
	return r.DB.Delete(&model.Articles{}, id).Error
}

func (r *ArticleRepositoryImpl) getCache(key string) (*model.Articles, error) {
	var article model.Articles
	data, err := r.Redis.Get(key)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(data), &article); err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepositoryImpl) setCache(key string, article *model.Articles) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.Redis.Set(key, string(data))
}

func (r *ArticleRepositoryImpl) delCache(key string) error {
	return r.Redis.Del(key)
}