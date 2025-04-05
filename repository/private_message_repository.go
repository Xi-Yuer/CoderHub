package repository

import (
	"coderhub/model"
	"coderhub/shared/storage"
	"coderhub/shared/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// PrivateMessage 私聊消息结构体
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

// PrivateMessageRepository 私聊消息相关操作
type PrivateMessageRepository interface {
	// Create 创建私聊消息
	Create(ctx context.Context, message *model.PrivateMessage) error
	// GetPrivateMessage 获取私聊消息
	GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error)
	// GetPrivateSessionHistoryMessage 获取会话中的所有聊天记录，分页查询，加载更多是加载双方的历史消息记录,这里不能使用 session_id 来查询，因为 session_id 是由两个用户的 ID 拼接而成的，所以需要使用 sender_id 和 receiver_id 来查询
	GetPrivateSessionHistoryMessage(ctx context.Context, senderID string, receiverID string, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error)
	// UpdatePrivateMessage 更新私聊消息
	UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error
	// GetOfflineMessages 获取用户的离线消息
	GetOfflineMessages(ctx context.Context, userID string) ([]*model.PrivateMessage, error)
	// MarkUserOnline 标记用户上线
	MarkUserOnline(ctx context.Context, userID string) error
	// MarkUserOffline 标记用户离线
	MarkUserOffline(ctx context.Context, userID string) error
}

// PrivateMessageRepositoryImpl 私聊消息仓库实现
type PrivateMessageRepositoryImpl struct {
	DB    *gorm.DB
	Redis storage.RedisDB
}

// NewPrivateMessageRepository 创建私聊消息仓库实例
func NewPrivateMessageRepository(db *gorm.DB, redis storage.RedisDB) *PrivateMessageRepositoryImpl {
	return &PrivateMessageRepositoryImpl{DB: db, Redis: redis}
}

// Create 创建私聊消息
func (p *PrivateMessageRepositoryImpl) Create(ctx context.Context, message *model.PrivateMessage) error {
	if err := p.DB.WithContext(ctx).Create(message).Error; err != nil {
		return err
	}
	// 检查接收者是否在线，如果不在线则将消息保存到 Redis
	onlineKey := fmt.Sprintf("user_online:%s", message.ReceiverID)
	isOnline, err := p.Redis.Exists(onlineKey)
	if err != nil {
		return err
	}
	if !isOnline {
		offlineKey := fmt.Sprintf("offline_messages:%s", message.ReceiverID)
		if err := p.Redis.RPush(offlineKey, message.MessageID); err != nil {
			return err
		}
	}
	return nil
}

// GetPrivateMessage 获取私聊消息
func (p *PrivateMessageRepositoryImpl) GetPrivateMessage(ctx context.Context, message *model.PrivateMessage, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error) {
	var messages []*model.PrivateMessage
	var total int64
	if err := p.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("sender_id =? AND receiver_id =?", message.SenderID, message.ReceiverID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := p.DB.WithContext(ctx).Model(&model.PrivateMessage{}).Where("sender_id =? AND receiver_id =?", message.SenderID, message.ReceiverID).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}

// GetPrivateSessionHistoryMessage 获取会话中的所有聊天记录，分页查询，加载更多是加载双方的历史消息记录
func (p *PrivateMessageRepositoryImpl) GetPrivateSessionHistoryMessage(ctx context.Context, senderID string, receiverID string, page int64, pageSize int64) ([]*model.PrivateMessage, int64, error) {
	var messages []*PrivateMessage
	var total int64

	// 将 string 类型的 senderID 和 receiverID 转换为 uint64 类型
	senderIDUint64, err := strconv.ParseUint(senderID, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	receiverIDUint64, err := strconv.ParseUint(receiverID, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	// 构建查询条件
	query := p.DB.WithContext(ctx).Model(&PrivateMessage{}).
		Where("(sender_id =? AND receiver_id =?) OR (sender_id =? AND receiver_id =?)", senderIDUint64, receiverIDUint64, receiverIDUint64, senderIDUint64)

	// 统计消息总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询指定页的消息，按时间戳升序排序
	if err := query.Order("timestamp DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 model.PrivateMessage 类型
	var result []*model.PrivateMessage
	for _, msg := range messages {
		result = append(result, &model.PrivateMessage{
			MessageID:   utils.Int2String(int64(msg.MessageID)),
			SenderID:    utils.Int2String(int64(msg.SenderID)),
			ReceiverID:  utils.Int2String(int64(msg.ReceiverID)),
			Content:     msg.Content,
			ContentType: msg.ContentType,
			Status:      msg.Status,
			Timestamp:   msg.Timestamp,
			IsRecalled:  msg.IsRecalled,
			CreatedAt:   msg.CreatedAt,
			UpdatedAt:   msg.UpdatedAt,
			DeletedAt:   msg.DeletedAt,
		})
	}

	// 反转消息顺序，确保最新的消息在最前面
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result, total, nil
}

// UpdatePrivateMessage 更新私聊消息
func (p *PrivateMessageRepositoryImpl) UpdatePrivateMessage(ctx context.Context, message *model.PrivateMessage) error {
	if err := p.DB.WithContext(ctx).Save(message).Error; err != nil {
		return err
	}
	return nil
}

// GetOfflineMessages 获取用户的离线消息
func (p *PrivateMessageRepositoryImpl) GetOfflineMessages(ctx context.Context, userID string) ([]*model.PrivateMessage, error) {
	key := fmt.Sprintf("offline_messages:%s", userID)
	messageIDs, err := p.Redis.LRange(key, 0, -1)
	if err != nil {
		return nil, err
	}

	var messages []*model.PrivateMessage
	for _, idStr := range messageIDs {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		var message model.PrivateMessage
		if err := p.DB.WithContext(ctx).First(&message, id).Error; err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	// 清空离线消息列表
	if err := p.Redis.Del(key); err != nil {
		return nil, err
	}

	return messages, nil
}

// MarkUserOnline 标记用户上线
func (p *PrivateMessageRepositoryImpl) MarkUserOnline(ctx context.Context, userID string) error {
	// 标记用户为在线
	err := p.Redis.Set(fmt.Sprintf("user_online:%s", userID), "1")
	if err != nil {
		return err
	}
	return nil
}

// MarkUserOffline 标记用户离线
func (p *PrivateMessageRepositoryImpl) MarkUserOffline(ctx context.Context, userID string) error {
	// 标记用户为离线
	err := p.Redis.Del(fmt.Sprintf("user_online:%s", userID))
	if err != nil {
		return err
	}
	return nil
}
