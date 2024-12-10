package repository

import (
	"coderhub/model"
	"context"

	"gorm.io/gorm"
)

type AcademicRelationLikeRepository interface {
	AddAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error
	DeleteAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error
	GetAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) bool
	GetAcademicRelationLikeCount(ctx context.Context, academicRelationLike *model.AcademicRelationLike) (int64, error)
}

type AcademicRelationLikeRepositoryImpl struct {
	DB *gorm.DB
}

func NewAcademicRelationLikeRepositoryImpl(db *gorm.DB) *AcademicRelationLikeRepositoryImpl {
	return &AcademicRelationLikeRepositoryImpl{
		DB: db,
	}
}

// 创建学术关系点赞
func (r *AcademicRelationLikeRepositoryImpl) AddAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error {
	return r.DB.Create(academicRelationLike).Error
}

// 删除学术关系点赞
func (r *AcademicRelationLikeRepositoryImpl) DeleteAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error {
	return r.DB.Delete(academicRelationLike).Where("id = ? AND user_id = ?", academicRelationLike.AcademicNavigatorID, academicRelationLike.UserID).Error
}

// 是否点赞
func (r *AcademicRelationLikeRepositoryImpl) GetAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) bool {
	var count int64
	err := r.DB.Model(&model.AcademicRelationLike{}).Where("id = ? AND user_id = ?", academicRelationLike.AcademicNavigatorID, academicRelationLike.UserID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// 获取学术关系点赞数量
func (r *AcademicRelationLikeRepositoryImpl) GetAcademicRelationLikeCount(ctx context.Context, academicRelationLike *model.AcademicRelationLike) (int64, error) {
	var count int64
	err := r.DB.Model(&model.AcademicRelationLike{}).Where("id = ?", academicRelationLike.AcademicNavigatorID).Count(&count).Error
	return count, err
}
