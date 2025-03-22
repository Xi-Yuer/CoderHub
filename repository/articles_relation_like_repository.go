package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"gorm.io/gorm"
)

type ArticlesRelationLikeRepository interface {
	Create(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) error
	Delete(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) error
	Get(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) bool
	List(ctx context.Context, articleID int64) (int64, error)
	BatchList(ctx context.Context, articleIDs []int64) (map[int64]int64, error)
	BatchArticlesHasBeenUserLiked(ctx context.Context, articleIDs []int64, userID int64) (map[int64]bool, error)
	GetArticleLikeCount(ctx context.Context, articleIDs []int64) (int64, error)
}
type articlesRelationLikeRepository struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func NewArticlesRelationLikeRepository(db *gorm.DB, redis storage.RedisDB) ArticlesRelationLikeRepository {
	return &articlesRelationLikeRepository{
		DB:    db,
		Redis: redis,
	}
}

// Create 创建文章点赞
func (r *articlesRelationLikeRepository) Create(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) error {
	return r.DB.Create(articleRelationLike).Error
}

// Delete 删除文章点赞
func (r *articlesRelationLikeRepository) Delete(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) error {
	return r.DB.Delete(articleRelationLike, "article_id = ? AND user_id = ?", articleRelationLike.ArticleID, articleRelationLike.UserID).Error
}

// Get 获取文章是否被用户点赞
func (r *articlesRelationLikeRepository) Get(ctx context.Context, articleRelationLike *model.ArticlesRelationLike) bool {
	var count int64
	r.DB.Model(articleRelationLike).Where("article_id = ? AND user_id = ?", articleRelationLike.ArticleID, articleRelationLike.UserID).Count(&count)
	return count > 0
}

// List 获取文章点赞数
func (r *articlesRelationLikeRepository) List(ctx context.Context, articleID int64) (int64, error) {
	var articlesRelationLikesCount int64
	r.DB.Model(&model.ArticlesRelationLike{}).Where("article_id = ?", articleID).Count(&articlesRelationLikesCount)
	return articlesRelationLikesCount, nil
}

// BatchList 批量获取文章点赞数
func (r *articlesRelationLikeRepository) BatchList(ctx context.Context, articleIDs []int64) (map[int64]int64, error) {
	articlesRelationLikes := make([]model.ArticlesRelationLike, 0)
	err := r.DB.Where("article_id IN (?)", articleIDs).Find(&articlesRelationLikes).Error
	if err != nil {
		return nil, err
	}
	articlesRelationLikeCountMap := make(map[int64]int64)
	for _, articlesRelationLike := range articlesRelationLikes {
		if _, ok := articlesRelationLikeCountMap[articlesRelationLike.ArticleID]; !ok {
			articlesRelationLikeCountMap[articlesRelationLike.ArticleID] = 1
		}
	}
	return articlesRelationLikeCountMap, nil
}

// BatchArticlesHasBeenUserLiked 批量获取文章是否被用户点赞
func (r *articlesRelationLikeRepository) BatchArticlesHasBeenUserLiked(ctx context.Context, articleIDs []int64, userID int64) (map[int64]bool, error) {
	articlesRelationLikes := make([]model.ArticlesRelationLike, 0)
	err := r.DB.Where("article_id IN (?) AND user_id = ?", articleIDs, userID).Find(&articlesRelationLikes).Error
	if err != nil {
		return nil, err
	}
	articlesRelationLikeMap := make(map[int64]bool)
	for _, articlesRelationLike := range articlesRelationLikes {
		articlesRelationLikeMap[articlesRelationLike.ArticleID] = true
	}
	return articlesRelationLikeMap, nil
}

func (r *articlesRelationLikeRepository) GetArticleLikeCount(ctx context.Context, articleIDs []int64) (int64, error) {
	var count int64
	err := r.DB.Model(&model.ArticlesRelationLike{}).Where("article_id IN (?)", articleIDs).Count(&count).Error
	return count, err
}
