package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"context"
	"gorm.io/gorm"
)

type QuestionBankRepository interface {
	CreateQuestionBank(ctx context.Context, questionBank *model.QuestionBank) error
	GetQuestionBankByID(ctx context.Context, id int64) (*model.QuestionBank, error)
	GetQuestionBanks(ctx context.Context, page, pageSize int32) ([]*model.QuestionBank, int64, error)
	UpdateQuestionBank(ctx context.Context, questionBank *model.QuestionBank) error
	DeleteQuestionBank(ctx context.Context, id int64) error
}

// NewQuestionRepositoryRepositoryImpl 实例
func NewQuestionRepositoryRepositoryImpl(db *gorm.DB, rdb storage.RedisDB) *QuestionRepositoryRepositoryImpl {
	return &QuestionRepositoryRepositoryImpl{
		DB:    db,
		Redis: rdb,
	}
}

type QuestionRepositoryRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

// CreateQuestionBank 创建题库
func (r *QuestionRepositoryRepositoryImpl) CreateQuestionBank(ctx context.Context, questionBank *model.QuestionBank) error {
	return r.DB.WithContext(ctx).Model(&model.QuestionBank{}).Create(questionBank).Error
}

func (r *QuestionRepositoryRepositoryImpl) GetQuestionBankByID(ctx context.Context, id int64) (*model.QuestionBank, error) {
	questionBank := &model.QuestionBank{}
	db := r.DB.WithContext(ctx).Where("id = ?", id).First(questionBank)
	return questionBank, db.Error
}

func (r *QuestionRepositoryRepositoryImpl) GetQuestionBanks(ctx context.Context, page, pageSize int32) ([]*model.QuestionBank, int64, error) {
	var questionBanks []*model.QuestionBank
	var total int64
	err := r.DB.WithContext(ctx).Model(&model.QuestionBank{}).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&questionBanks).Count(&total).Error
	return questionBanks, total, err
}

func (r *QuestionRepositoryRepositoryImpl) UpdateQuestionBank(ctx context.Context, questionBank *model.QuestionBank) error {
	return r.DB.WithContext(ctx).Model(&model.QuestionBank{}).Where("id = ?", questionBank.ID).Updates(questionBank).Error
}

func (r *QuestionRepositoryRepositoryImpl) DeleteQuestionBank(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.QuestionBank{}).Error
}
