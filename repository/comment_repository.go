package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	Delete(ctx context.Context, id int64) error
	ListByArticleID(ctx context.Context, articleID int64, page int64, pageSize int64) ([]model.Comment, int64, error)
	ListReplies(ctx context.Context, rootID int64, page int64, pageSize int64) ([]model.Comment, int64, error)
	CountByArticleID(ctx context.Context, articleID int64) (int64, error)
}

var ErrConcurrentUpdate = errors.New("并发更新冲突，请重试")

// commentRepository 评论仓储实现
type commentRepository struct {
	DB       *gorm.DB
	Redis    storage.RedisDB
	minLikes int32
}

// NewCommentRepository 创建评论仓储实例
func NewCommentRepository(db *gorm.DB, rdb storage.RedisDB) CommentRepository {
	return &commentRepository{
		DB:       db,
		Redis:    rdb,
		minLikes: 10,
	}
}

// Create 创建评论
func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.DB.WithContext(ctx).Create(comment).Error
}

// GetByID 根据ID获取评论，优先从缓存获取
func (r *commentRepository) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	// 从数据库获取
	if err := r.DB.WithContext(ctx).First(comment, id).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// Delete 删除评论并清除缓存
func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	// 判断该评论是否有子级评论
	var count int64
	r.DB.WithContext(ctx).Model(&model.Comment{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		// 获取原始评论
		originalComment, err := r.GetByID(ctx, id)
		if err != nil {
			return err
		}
		// 更新评论内容
		originalComment.Content = "该评论已被删除"
		if err := r.UpdateByID(ctx, id, originalComment); err != nil {
			return err
		}
		return nil
	}
	err := r.DB.WithContext(ctx).Delete(&model.Comment{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

// ListByArticleID 获取文章的顶级评论列表(分页)
func (r *commentRepository) ListByArticleID(ctx context.Context, articleID int64, page int64, pageSize int64) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	// 查询评论列表 parent_id = 0 为顶级评论，此外还要额外查询顶级下三条回复评论,并且构建树形结构
	fmt.Println("文章ID: ", articleID)
	fmt.Println("页码: ", page)
	fmt.Println("每页大小: ", pageSize)
	// 查询顶级评论
	r.DB.WithContext(ctx).Where("entity_id = ? AND root_id = 0", articleID).Order("created_at ASC").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&comments)
	// 查询评论总数
	r.DB.WithContext(ctx).Model(&model.Comment{}).Where("entity_id = ?", articleID).Count(&total)
	// 构建树形结构并获取回复数量
	for i := range comments {
		if comments[i].RootID == 0 {
			// 获取回复总数
			var replyCount int64
			// 查询回复总数
			if err := r.DB.WithContext(ctx).Model(&model.Comment{}).
				Where("root_id = ?", comments[i].ID).
				Count(&replyCount).Error; err != nil {
				return nil, 0, err
			}
			comments[i].ReplyCount = replyCount // 假设 Comment 结构体中有 ReplyCount 字段

			// 获取前三条回复
			replies, _, err := r.ListReplies(ctx, comments[i].ID, 1, 3)
			if err != nil {
				return nil, 0, err
			}
			comments[i].Replies = replies
		}
	}

	return comments, total, nil
}

// ListReplies 获取评论的回复列表(分页)，这里评论一定是顶级评论的回复
func (r *commentRepository) ListReplies(ctx context.Context, rootID int64, page int64, pageSize int64) ([]model.Comment, int64, error) {
	var replies []model.Comment
	var total int64
	r.DB.WithContext(ctx).Where("root_id = ?", rootID).Order("created_at ASC").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&replies)
	// 查询回复总数
	r.DB.WithContext(ctx).Model(&model.Comment{}).Where("root_id = ?", rootID).Count(&total)
	return replies, total, nil
}

// UpdateByID 更新评论
func (r *commentRepository) UpdateByID(ctx context.Context, id int64, comment *model.Comment) error {
	return r.DB.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).Updates(comment).Error
}

// CountByArticleID 获取文章评论数
func (r *commentRepository) CountByArticleID(ctx context.Context, articleID int64) (int64, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&model.Comment{}).Where("entity_id = ?", articleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
