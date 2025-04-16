package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
	// 先查询缓存
	key := fmt.Sprintf("article:like:%d:%d", articleRelationLike.ArticleID, articleRelationLike.UserID)
	exists, err := r.Redis.Exists(key)
	if err == nil && exists {
		return true
	}

	// 缓存未命中，查询数据库
	var count int64
	r.DB.Model(articleRelationLike).
		Select("1").
		Where("article_id = ? AND user_id = ?", articleRelationLike.ArticleID, articleRelationLike.UserID).
		Limit(1).
		Count(&count)

	// 如果存在点赞关系，写入缓存
	if count > 0 {
		r.Redis.Set(key, "1")
	}
	return count > 0
}

// BatchArticlesHasBeenUserLiked 批量获取文章是否被用户点赞
func (r *articlesRelationLikeRepository) BatchArticlesHasBeenUserLiked(ctx context.Context, articleIDs []int64, userID int64) (map[int64]bool, error) {
	result := make(map[int64]bool)

	// 批量查询缓存
	pipeline := r.Redis.Pipeline()
	keys := make([]string, len(articleIDs))
	for i, articleID := range articleIDs {
		keys[i] = fmt.Sprintf("article:like:%d:%d", articleID, userID)
		pipeline.Exists(ctx, keys[i])
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 记录未命中缓存的文章ID
	missedArticleIDs := make([]int64, 0)
	for i, articleID := range articleIDs {
		if existsCmd := cmders[i].(*redis.IntCmd); existsCmd.Val() > 0 {
			result[articleID] = true
		} else {
			missedArticleIDs = append(missedArticleIDs, articleID)
		}
	}

	// 如果有未命中的，查询数据库
	if len(missedArticleIDs) > 0 {
		var likes []model.ArticlesRelationLike
		err := r.DB.Select("article_id").
			Where("article_id IN (?) AND user_id = ?", missedArticleIDs, userID).
			Find(&likes).Error
		if err != nil {
			return nil, err
		}

		// 更新结果和缓存
		pipeline = r.Redis.Pipeline()
		for _, like := range likes {
			result[like.ArticleID] = true
			key := fmt.Sprintf("article:like:%d:%d", like.ArticleID, userID)
			pipeline.Set(ctx, key, 1, time.Hour*24)
		}
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
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
// BatchArticlesHasBeenUserLikedV2 批量获取文章是否被用户点赞的第二个实现版本
func (r *articlesRelationLikeRepository) BatchArticlesHasBeenUserLikedV2(ctx context.Context, articleIDs []int64, userID int64) (map[int64]bool, error) {
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
