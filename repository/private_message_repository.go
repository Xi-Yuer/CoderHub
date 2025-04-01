package repository

import (
	"coderhub/model"
	"context"

	"gorm.io/gorm"
)

/**
	type PrivateMessage struct {
	MessageID   uint64         `gorm:"primaryKey;column:message_id;comment:消息ID（雪花算法生成）"`
	SenderID    uint64         `gorm:"index:idx_sender_receiver;column:sender_id;not null;comment:发送者ID"`
	ReceiverID  uint64         `gorm:"index:idx_sender_receiver;index:idx_receiver_status;column:receiver_id;not null;comment:接收者ID"`
	Content     string         `gorm:"type:text;column:content;comment:消息内容"`
	ContentType string         `gorm:"type:varchar(20);column:content_type;comment:消息类型(text/image/file/audio)"`
	Status      string         `gorm:"type:ENUM('sent', 'delivered', 'read');column:status;default:'sent';comment:消息状态"`
	Timestamp   int64          `gorm:"column:timestamp;comment:消息时间戳（毫秒级）"`
	IsRecalled  bool           `gorm:"column:is_recalled;default:false;comment:是否已撤回"`
	CreatedAt   time.Time      `gorm:"->;column:created_at;autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"->;column:updated_at;autoUpdateTime;comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at;comment:软删除时间"`
}
**/

// PrivateMessageRepository 私聊消息相关操作
type PrivateMessageRepository interface {
	// Create 创建私聊消息
	Create(ctx context.Context, message *model.PrivateMessage) error
	// GetPrivateMessage 获取私聊消息
	GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error)
	// UpdatePrivateMessage 更新私聊消息
	UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error
}

type PrivateMessageRepositoryImpl struct {
	DB *gorm.DB
}

func NewPrivateMessageRepository(db *gorm.DB) *PrivateMessageRepositoryImpl {
	return &PrivateMessageRepositoryImpl{DB: db}
}

func (p *PrivateMessageRepositoryImpl) Create(ctx context.Context, message *model.PrivateMessage) error {
	if err := p.DB.WithContext(ctx).Create(message).Error; err != nil {
		return err
	}
	return nil
}
func (p *PrivateMessageRepositoryImpl) GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error) {
	var messages []*model.PrivateMessage
	var total int64
	if err := p.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("sender_id = ? AND receiver_id = ?", message.SenderID, message.ReceiverID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := p.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("sender_id =? AND receiver_id =?", message.SenderID, message.ReceiverID).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}
func (p *PrivateMessageRepositoryImpl) UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error {
	if err := p.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("message_id =?", message.MessageID).Updates(message).Error; err != nil {
		return err
	}
	return nil
}
