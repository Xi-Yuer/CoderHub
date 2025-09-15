package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"

	"gorm.io/gorm"
)

type CommentRelationLikeRepository interface {
	Create(ctx context.Context, commentRelationLike *model.CommentRelationLike) error
	Delete(ctx context.Context, commentRelationLike *model.CommentRelationLike) error
	Get(ctx context.Context, commentRelationLike *model.CommentRelationLike) bool
	List(ctx context.Context, commentID int64) (int64, error)
	BatchList(ctx context.Context, commentIDs []int64) (map[int64]int64, error)
	BatchGetCommentsHasBeenUserLiked(ctx context.Context, commentIDs []int64, userID int64) (map[int64]bool, error)
}

type commentRelationLikeRepository struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func NewCommentRelationLikeRepository(db *gorm.DB, redis storage.RedisDB) CommentRelationLikeRepository {
	return &commentRelationLikeRepository{
		DB:    db,
		Redis: redis,
	}
}

// Create 创建评论点赞
func (r *commentRelationLikeRepository) Create(ctx context.Context, commentRelationLike *model.CommentRelationLike) error {
	return r.DB.Create(commentRelationLike).Error
}

// Delete 删除评论点赞
func (r *commentRelationLikeRepository) Delete(ctx context.Context, commentRelationLike *model.CommentRelationLike) error {
	return r.DB.Delete(commentRelationLike, "comment_id = ? AND user_id = ?", commentRelationLike.CommentID, commentRelationLike.UserID).Error
}

// Get 获取评论是否被用户点赞
func (r *commentRelationLikeRepository) Get(ctx context.Context, commentRelationLike *model.CommentRelationLike) bool {
	var count int64
	r.DB.Model(commentRelationLike).Where("comment_id = ? AND user_id = ?", commentRelationLike.CommentID, commentRelationLike.UserID).Count(&count)
	return count > 0
}

// List 获取评论点赞数
func (r *commentRelationLikeRepository) List(ctx context.Context, commentID int64) (int64, error) {
	var count int64
	r.DB.Model(&model.CommentRelationLike{}).Where("comment_id = ?", commentID).Count(&count)
	return count, nil
}

// BatchList 批量获取评论点赞数
func (r *commentRelationLikeRepository) BatchList(ctx context.Context, commentIDs []int64) (map[int64]int64, error) {
	commentRelationLikes := make([]model.CommentRelationLike, 0)
	err := r.DB.Where("comment_id IN (?)", commentIDs).Find(&commentRelationLikes).Error
	if err != nil {
		return nil, err
	}
	commentLikeCountMap := make(map[int64]int64)
	for _, commentRelationLike := range commentRelationLikes {
		if _, ok := commentLikeCountMap[commentRelationLike.CommentID]; !ok {
			commentLikeCountMap[commentRelationLike.CommentID] = 0
		}
		commentLikeCountMap[commentRelationLike.CommentID]++
	}
	return commentLikeCountMap, nil
}

// BatchGetCommentsHasBeenUserLiked 批量获取评论是否被用户点赞
func (r *commentRelationLikeRepository) BatchGetCommentsHasBeenUserLiked(ctx context.Context, commentIDs []int64, userID int64) (map[int64]bool, error) {
	commentRelationLikes := make([]model.CommentRelationLike, 0)
	err := r.DB.Where("comment_id IN (?) AND user_id = ?", commentIDs, userID).Find(&commentRelationLikes).Error
	if err != nil {
		return nil, err
	}
	commentLikeMap := make(map[int64]bool)
	for _, commentRelationLike := range commentRelationLikes {
		commentLikeMap[commentRelationLike.CommentID] = true
	}
	return commentLikeMap, nil
}
