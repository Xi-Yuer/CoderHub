package repository

import (
	"coderhub/model"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, message *model.Message, page, pageSize int64) ([]*model.Message, int64, error)
}

type MessageRepositoryImpl struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{
		DB: db,
	}
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	tx := r.DB.WithContext(ctx).Begin()
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create message: %w", err)
	}
	return tx.Commit().Error
}

func (r *MessageRepositoryImpl) GetMessage(ctx context.Context, message *model.Message) (*model.Message, error) {
	var msg *model.Message
	err := r.DB.WithContext(ctx).Model(&model.Message{}).First(&msg, message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("message not found: %w", err)
		}
		return nil, err
	}
	return msg, nil
}

func (r *MessageRepositoryImpl) Delete(ctx context.Context, id int64) error {
	var msg model.Message
	err := r.DB.WithContext(ctx).First(&msg, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("message not found: %w", err)
		}
		return err
	}
	return r.DB.WithContext(ctx).Delete(&msg).Error
}

func (r *MessageRepositoryImpl) List(ctx context.Context, message *model.Message, page, pageSize int64) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	err := r.DB.WithContext(ctx).Model(&model.Message{}).Where(message).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	if total > 0 {
		err = r.DB.WithContext(ctx).
			Where(message).
			Limit(int(pageSize)).
			Offset(int((page - 1) * pageSize)).
			Find(&messages).Error
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch messages: %w", err)
		}
	}
	return messages, total, nil
}
