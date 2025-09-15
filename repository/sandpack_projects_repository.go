package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type SandpackProjectsRepository interface {
	Create(ctx context.Context, sandpacks *model.SandpackProjects) (sandpaperProjectsId int64, err error)
	Get(ctx context.Context, userID, articleID int64) (*model.SandpackProjects, error)
}

type SandpackProjectsRepositoryImpl struct {
	DB *gorm.DB
}

func NewSandpackProjectsRepository(db *gorm.DB) SandpackProjectsRepository {
	return &SandpackProjectsRepositoryImpl{
		DB: db,
	}
}

func (s *SandpackProjectsRepositoryImpl) Create(ctx context.Context, sandpacks *model.SandpackProjects) (int64, error) {
	result := &model.SandpackProjects{
		Name:        sandpacks.Name,
		Description: sandpacks.Description,
		Template:    sandpacks.Template,
		UserID:      sandpacks.UserID,
		ArticleID:   sandpacks.ArticleID,
	}
	err := s.DB.WithContext(ctx).Create(result).Error
	return int64(result.ID), err
}

func (s *SandpackProjectsRepositoryImpl) Get(ctx context.Context, userID, articleID int64) (*model.SandpackProjects, error) {
	var sandpack model.SandpackProjects
	return &sandpack, s.DB.WithContext(ctx).Where("user_id = ? and article_id = ?", userID, articleID).First(&sandpack).Error
}
