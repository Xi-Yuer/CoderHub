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
	BatchGetQuestion(ctx context.Context, ids []int64) ([]*model.QuestionBanksPreviewWithCreateUser, error)
	GetQuestionBanks(ctx context.Context, categoryID int64, page, pageSize int32) ([]*model.QuestionBank, int64, error)
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
	err := r.DB.WithContext(ctx).Where("id = ?", id).Limit(1).Find(questionBank).Error
	if err != nil {
		return nil, err
	}
	return questionBank, nil
}

func (r *QuestionRepositoryRepositoryImpl) GetQuestionBanks(ctx context.Context, categoryID int64, page, pageSize int32) ([]*model.QuestionBank, int64, error) {
	var questionBanks []*model.QuestionBank
	var total int64
	err := r.DB.WithContext(ctx).Model(&model.QuestionBank{}).Where("cate_gory_id = ?", categoryID).Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&questionBanks).Count(&total).Error
	return questionBanks, total, err
}

func (r *QuestionRepositoryRepositoryImpl) BatchGetQuestion(ctx context.Context, ids []int64) ([]*model.QuestionBanksPreviewWithCreateUser, error) {
	var questionBanks []*model.QuestionBanksPreviewWithCreateUser
	err := r.DB.Table("question_banks AS qb").
		Select(`
        qb.id AS id, 
        qb.name, 
        img.url AS cover_image, 
        qb.description, 
        qb.difficulty, 
        qb.tags, 
        qb.created_at AS create_time, 
        u.id AS create_user_id, 
        u.user_name AS create_user_name, 
        u.avatar
    `).
		Joins("JOIN users u ON qb.create_user = u.id").
		Joins("JOIN image_relations ir ON qb.id = ir.entity_id AND ir.entity_type = ?", "question_cover").
		Joins("JOIN images img ON ir.image_id = img.id").
		Where("qb.id IN ?", ids).
		Scan(&questionBanks).Error
	return questionBanks, err
}

func (r *QuestionRepositoryRepositoryImpl) UpdateQuestionBank(ctx context.Context, questionBank *model.QuestionBank) error {
	return r.DB.WithContext(ctx).Model(&model.QuestionBank{}).Where("id = ?", questionBank.ID).Updates(questionBank).Error
}

func (r *QuestionRepositoryRepositoryImpl) DeleteQuestionBank(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.QuestionBank{}).Error
}
