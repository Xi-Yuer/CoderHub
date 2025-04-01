package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	Sent         string = "sent"
	Delivered    string = "delivered"
	Read         string = "read"
	TextMessage  string = "text"
	ImageMessage string = "image"
	FileMessage  string = "file"
	VideoMessage string = "video"
	AudioMessage string = "audio"
)

// PrivateMessage 私聊消息实体
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

// UserSession 用户会话实体（用于查询用户的会话信息）
type UserSession struct {
	UserID        uint64 `gorm:"primaryKey;column:user_id;comment:用户ID"`
	PeerID        uint64 `gorm:"primaryKey;column:peer_id;comment:对方用户ID"`
	LastMessageID uint64 `gorm:"column:last_message_id;comment:最后一条消息ID"`
	UnreadCount   int    `gorm:"column:unread_count;default:0;comment:未读消息数量"`
	UpdatedAt     int64  `gorm:"column:updated_at;comment:最后更新时间戳（毫秒）"`
}
