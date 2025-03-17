package repository

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"
)

type QuestionBankCategoryRepository interface {
	Create(ctx context.Context, questionBankCategory *model.QuestionBankCategory) error
	Delete(ctx context.Context, id int64) error
	GetCategoryByID(ctx context.Context, id int64) (*model.QuestionBankCategory, error)
	List(ctx context.Context) ([]model.QuestionBankCategory, error)
}

type QuestionBankCategoryRepositoryImpl struct {
	DB *gorm.DB
}

func NewQuestionBankCategoryRepository(db *gorm.DB) QuestionBankCategoryRepository {
	return &QuestionBankCategoryRepositoryImpl{DB: db}
}

func (r *QuestionBankCategoryRepositoryImpl) Create(ctx context.Context, questionBankCategory *model.QuestionBankCategory) error {
	return r.DB.WithContext(ctx).Create(questionBankCategory).Error
}

func (r *QuestionBankCategoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.QuestionBankCategory{}, id).Error
}

func (r *QuestionBankCategoryRepositoryImpl) GetCategoryByID(ctx context.Context, id int64) (*model.QuestionBankCategory, error) {
	var questionBankCategory model.QuestionBankCategory
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&questionBankCategory).Error
	return &questionBankCategory, err
}

func (r *QuestionBankCategoryRepositoryImpl) List(ctx context.Context) ([]model.QuestionBankCategory, error) {
	var questionBankCategories []model.QuestionBankCategory
	err := r.DB.WithContext(ctx).Find(&questionBankCategories).Error
	return questionBankCategories, err
}
