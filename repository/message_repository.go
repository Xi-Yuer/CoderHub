package repository

import (
	"coderhub/model"
	"context"
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
	GetMessage(ctx context.Context, message *model.Message) (*model.Message, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, message *model.Message, page, pageSize int64) ([]*model.Message, int64, error)
	GetUnReadMessageCount(ctx context.Context, receiverId int64) (int32, error)
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

	// 使用 Select 指定只计数需要的字段
	query := r.DB.WithContext(ctx).Model(&model.Message{}).Where(message)
	countQuery := query.Session(&gorm.Session{}) // 创建新会话避免影响主查询

	// 并发执行计数查询和数据查询
	var wg sync.WaitGroup
	var countErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		countErr = countQuery.Count(&total).Error
	}()

	var fetchErr error
	if page > 0 && pageSize > 0 {
		fetchErr = query.
			Order("is_read ASC, id DESC"). // 添加第二个排序条件提高性能
			Limit(int(pageSize)).
			Offset(int((page - 1) * pageSize)).
			Find(&messages).Error
	}

	wg.Wait()

	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", countErr)
	}
	if fetchErr != nil {
		return nil, 0, fmt.Errorf("failed to fetch messages: %w", fetchErr)
	}

	// 只有在有消息且需要更新状态时才执行更新
	if len(messages) > 0 {
		if err := r.markMessagesAsRead(ctx, messages); err != nil {
			return nil, 0, err
		}
	}

	return messages, total, nil
}

// 将更新逻辑抽取为单独的方法
func (r *MessageRepositoryImpl) markMessagesAsRead(ctx context.Context, messages []*model.Message) error {
	return r.DB.WithContext(ctx).Model(&model.Message{}).
		Where("id IN ? AND is_read = ?", getMessageIDs(messages), false).
		Update("is_read", true).Error
}

func (r *MessageRepositoryImpl) GetUnReadMessageCount(ctx context.Context, receiverId int64) (int32, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.Message{}).Where("receiver_id = ? AND is_read = ?", receiverId, false).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	// 查询用户所有会话信息未读数量
	var userSessions []model.UserSession
	err = r.DB.WithContext(ctx).Where("user_id = ?", receiverId).Find(&userSessions).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch user sessions: %w", err)
	}
	for _, userSession := range userSessions {
		count += int64(userSession.UnreadMessageCount)
	}
	return int32(count), nil
}

// getMessageIDs 提取消息 ID 列表
func getMessageIDs(messages []*model.Message) []int64 {
	ids := make([]int64, len(messages))
	for i, msg := range messages {
		ids[i] = int64(msg.ID)
	}
	return ids
}
