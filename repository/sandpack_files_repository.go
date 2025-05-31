package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type SandpackFilesRepository interface {
	Create(ctx context.Context, sandpackFiles []*model.SandpackFiles) error
	Get(ctx context.Context, sandpackProjects int64) ([]*model.SandpackFiles, error)
}
type SandpackFilesRepositoryImpl struct {
	DB *gorm.DB
}

func NewSandpackFilesRepository(db *gorm.DB) SandpackFilesRepository {
	return &SandpackFilesRepositoryImpl{DB: db}
}

func (r *SandpackFilesRepositoryImpl) Create(ctx context.Context, sandpackFiles []*model.SandpackFiles) error {
	return r.DB.WithContext(ctx).Create(sandpackFiles).Error
}

func (r *SandpackFilesRepositoryImpl) Get(ctx context.Context, sandpackProjects int64) ([]*model.SandpackFiles, error) {
	var sandpackFiles []*model.SandpackFiles
	return sandpackFiles, r.DB.WithContext(ctx).Where("sandpack_id = ?", sandpackProjects).Find(&sandpackFiles).Error
}
