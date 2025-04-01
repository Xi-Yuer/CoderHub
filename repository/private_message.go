package repository

import (
	"coderhub/model"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PrivateMessageRepository interface {
	Create(ctx context.Context, message *model.PrivateMessage) error
	GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error)
	UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error
}

type PrivateMessageRepositoryImpl struct {
	DB *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) *PrivateMessageRepositoryImpl {
	return &PrivateMessageRepositoryImpl{DB: db}
}

func (r *PrivateMessageRepositoryImpl) Create(ctx context.Context, message *model.PrivateMessage) error {
	// 新增消息并且更新会话，没有会话消息则创建会话
	tx := r.DB.WithContext(ctx).Begin()
	var err error
	// 创建消息
	if err = tx.Create(message).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create message: %w", err)
	}
	// 查询会话
	var userSession model.UserSession
	if err = tx.Where("receiver_id = ? AND sender_id = ?", message.ReceiverID, message.SenderID).First(&userSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有会话消息则创建会话
			userSession = model.UserSession{
				UserID:        message.SenderID,
				PeerID:        message.ReceiverID,
				LastMessageID: 0,
				UnreadCount:   0,
				UpdatedAt:     time.Now().Unix(),
			}
			if err = tx.Create(&userSession).Error; err != nil {
			}
		}
	}
	if err = tx.Model(&userSession).Where("receiver_id = ? AND sender_id = ?", message.ReceiverID, message.SenderID).Updates(model.UserSession{
		LastMessageID: message.MessageID,
		UpdatedAt:     time.Now().Unix(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user session: %w", err)
	}
	return tx.Commit().Error
}

func (r *PrivateMessageRepositoryImpl) GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error) {
	var messages []*model.PrivateMessage
	var total int64
	err := r.DB.WithContext(ctx).Where("receiver_id = ? AND sender_id = ?", message.ReceiverID, message.SenderID).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&messages).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get private messages: %w", err)
	}
	err = r.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("receiver_id = ? AND sender_id = ?", message.ReceiverID, message.SenderID).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count private messages: %w", err)
	}
	return messages, total, nil
}

func (r *PrivateMessageRepositoryImpl) UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error {
	return r.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("message_id = ?", message.MessageID).Updates(message).Error
}
