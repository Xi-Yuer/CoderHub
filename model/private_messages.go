package model

import (
	"time"

	"gorm.io/gorm"
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
	MessageID   string         `gorm:"primaryKey;column:message_id;comment:消息ID（雪花算法生成）" json:"message_id"`
	SessionID   string         `gorm:"index:idx_session_id;column:session_id;not null;comment:会话ID" json:"session_id"`
	SenderID    string         `gorm:"index:idx_sender_receiver;column:sender_id;not null;comment:发送者ID" json:"sender_id"`
	ReceiverID  string         `gorm:"index:idx_sender_receiver;index:idx_receiver_status;column:receiver_id;not null;comment:接收者ID" json:"receiver_id"`
	Content     string         `gorm:"type:text;column:content;comment:消息内容" json:"content"`
	ContentType string         `gorm:"type:varchar(20);column:content_type;comment:消息类型(text/image/file/audio)" json:"content_type"`
	Status      string         `gorm:"type:ENUM('sent', 'delivered', 'read');column:status;default:'sent';comment:消息状态" json:"status"`
	Timestamp   int64          `gorm:"column:timestamp;comment:消息时间戳（毫秒级）" json:"timestamp"`
	IsRecalled  bool           `gorm:"column:is_recalled;default:false;comment:是否已撤回" json:"is_recalled"`
	CreatedAt   time.Time      `gorm:"->;column:created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"->;column:updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at;comment:软删除时间" json:"-"`
}

// UserSession 用户会话实体（用于查询用户的会话信息）
type UserSession struct {
	SessionID          string `gorm:"primaryKey;column:session_id;comment:会话ID" json:"session_id"`
	SessionName        string `gorm:"column:session_name;comment:会话名称" json:"session_name"`
	UserID             string `gorm:"column:user_id;comment:用户ID;index:idx_user_peer,unique" json:"user_id"`
	PeerID             string `gorm:"column:peer_id;comment:对方用户ID;index:idx_user_peer,unique" json:"peer_id"`
	LastMessageID      string `gorm:"column:last_message_id;comment:最后一条消息ID" json:"last_message_id"`
	LastMessageContent string `gorm:"column:last_message_content;comment:最后一条消息内容" json:"last_message_content"`
	UnreadMessageCount int    `gorm:"column:unread_message_count;default:0;comment:未读消息数量" json:"unread_message_count"`
	UnreadCount        int    `gorm:"column:unread_count;default:0;comment:未读消息数量" json:"unread_count"`
	CreatedAt          int64  `gorm:"column:created_at;comment:创建时间戳（毫秒）" json:"created_at"`
	UpdatedAt          int64  `gorm:"column:updated_at;comment:最后更新时间戳（毫秒）" json:"updated_at"`
}
