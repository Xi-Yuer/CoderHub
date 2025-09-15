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
	BatchGetAcademicRelationLikeCount(ctx context.Context, IDs []int64) (map[int64]int64, error)
}

type AcademicRelationLikeRepositoryImpl struct {
	DB *gorm.DB
}

func NewAcademicRelationLikeRepositoryImpl(db *gorm.DB) *AcademicRelationLikeRepositoryImpl {
	return &AcademicRelationLikeRepositoryImpl{
		DB: db,
	}
}

// AddAcademicRelationLike 创建学术关系点赞
func (r *AcademicRelationLikeRepositoryImpl) AddAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error {
	like := r.GetAcademicRelationLike(ctx, academicRelationLike)
	if !like {
		return r.DB.Create(academicRelationLike).Error
	} else {
		return r.DB.Where("id = ? AND user_id = ?", academicRelationLike.AcademicNavigatorID, academicRelationLike.UserID).Delete(&model.AcademicRelationLike{}).Error
	}
}

// DeleteAcademicRelationLike 删除学术关系点赞
func (r *AcademicRelationLikeRepositoryImpl) DeleteAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) error {
	return r.DB.Delete(&model.AcademicRelationLike{}).Where("id = ? AND user_id = ?", academicRelationLike.AcademicNavigatorID, academicRelationLike.UserID).Error
}

// GetAcademicRelationLike 是否点赞
func (r *AcademicRelationLikeRepositoryImpl) GetAcademicRelationLike(ctx context.Context, academicRelationLike *model.AcademicRelationLike) bool {
	var count int64
	err := r.DB.Model(&model.AcademicRelationLike{}).Where("id = ? AND user_id = ?", academicRelationLike.AcademicNavigatorID, academicRelationLike.UserID).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// GetAcademicRelationLikeCount 获取学术关系点赞数量
func (r *AcademicRelationLikeRepositoryImpl) GetAcademicRelationLikeCount(ctx context.Context, academicRelationLike *model.AcademicRelationLike) (int64, error) {
	var count int64
	err := r.DB.Model(&model.AcademicRelationLike{}).Where("id = ?", academicRelationLike.AcademicNavigatorID).Count(&count).Error
	return count, err
}

// BatchGetAcademicRelationLikeCount 批量获取学术关系点赞数量
func (r *AcademicRelationLikeRepositoryImpl) BatchGetAcademicRelationLikeCount(ctx context.Context, IDs []int64) (map[int64]int64, error) {
	academicRelationLikes := make([]model.AcademicRelationLike, 0)
	err := r.DB.Where("id IN (?)", IDs).Find(&academicRelationLikes).Error
	if err != nil {
		return nil, err
	}
	academicRelationLikeCountMap := make(map[int64]int64)
	for _, academicRelationLike := range academicRelationLikes {
		if _, ok := academicRelationLikeCountMap[academicRelationLike.AcademicNavigatorID]; !ok {
			academicRelationLikeCountMap[academicRelationLike.AcademicNavigatorID] = 0
		}
		academicRelationLikeCountMap[academicRelationLike.AcademicNavigatorID]++
	}
	return academicRelationLikeCountMap, nil
}
