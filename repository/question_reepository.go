package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *model.Question) error
	GetQuestionByID(ctx context.Context, id int64) (*model.Question, error)
	GetQuestions(ctx context.Context, ids []int64, page, pageSize int32) ([]*model.Question, int64, error)
	UpdateQuestion(ctx context.Context, question *model.Question) error
	DeleteQuestion(ctx context.Context, id int64) error
}

func NewQuestionRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *QuestionRepositoryImpl {
	return &QuestionRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

type QuestionRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

func (r *QuestionRepositoryImpl) CreateQuestion(ctx context.Context, question *model.Question) error {
	return r.DB.WithContext(ctx).Create(question).Error
}

func (r *QuestionRepositoryImpl) GetQuestionByID(ctx context.Context, id int64) (*model.Question, error) {
	var question model.Question
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&question).Error
	return &question, err
}

func (r *QuestionRepositoryImpl) GetQuestions(ctx context.Context, ids []int64, page, pageSize int32) ([]*model.Question, int64, error) {
	var questions []*model.Question
	var total int64
	err := r.DB.WithContext(ctx).Where("bank_id IN (?)", ids).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&questions).Count(&total).Error
	return questions, total, err
}

func (r *QuestionRepositoryImpl) UpdateQuestion(ctx context.Context, question *model.Question) error {
	return r.DB.WithContext(ctx).Updates(question).Error
}

func (r *QuestionRepositoryImpl) DeleteQuestion(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&model.Question{}, id).Error
}
