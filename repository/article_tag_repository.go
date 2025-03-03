package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

func NewArticleTagRepositoryImpl(db *gorm.DB) *ArticleTagRepositoryImpl {
	return &ArticleTagRepositoryImpl{DB: db}
}

type ArticleTagRepository interface {
	Create(ctx context.Context, articleTag *model.ArticleTag) error
	Delete(ctx context.Context, articleTag *model.ArticleTag) error
	GetList(ctx context.Context, page int64, pageSize int64) ([]*model.ArticleTag, int64, error)
	GetSystemTags(ctx context.Context) ([]*model.ArticleTag, int64, error)
	UpdateTagUsageCount(ctx context.Context, tagID int64) (err error)
}

type ArticleTagRepositoryImpl struct {
	DB *gorm.DB
}

func (r *ArticleTagRepositoryImpl) Create(ctx context.Context, articleTag *model.ArticleTag) error {
	return r.DB.WithContext(ctx).Create(articleTag).Error
}

func (r *ArticleTagRepositoryImpl) Delete(ctx context.Context, articleTag *model.ArticleTag) error {
	return r.DB.WithContext(ctx).Delete(articleTag).Error
}

func (r *ArticleTagRepositoryImpl) GetList(ctx context.Context, page int64, pageSize int64) ([]*model.ArticleTag, int64, error) {
	// 分页查询，且优先获取系统标签
	var articleTags []*model.ArticleTag
	var count int64
	if err := r.DB.WithContext(ctx).Where("is_system_provider = ?", true).Order("created_at desc").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&articleTags).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return articleTags, count, nil
}

func (r *ArticleTagRepositoryImpl) GetSystemTags(ctx context.Context) ([]*model.ArticleTag, int64, error) {
	// 获取系统标签
	var articleTags []*model.ArticleTag
	var count int64
	if err := r.DB.WithContext(ctx).Where("is_system_provider = ?", true).Order("created_at desc").Find(&articleTags).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return articleTags, count, nil
}

func (r *ArticleTagRepositoryImpl) UpdateTagUsageCount(ctx context.Context, tagID int64) (err error) {
	// 更新标签使用次数
	if err = r.DB.WithContext(ctx).Model(&model.ArticleTag{}).Where("id = ?", tagID).Update("usage_count", gorm.Expr("usage_count + 1")).Error; err != nil {
		return err
	}
	return nil
}
