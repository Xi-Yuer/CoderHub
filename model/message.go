package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SenderID   int64  `gorm:"column:sender_id;not null;index:idx_sender_id" json:"sender_id"`       // 发送者ID
	ReceiverID int64  `gorm:"column:receiver_id;not null;index:idx_receiver_id" json:"receiver_id"` // 接收者ID
	Type       int32  `gorm:"column:type;not null;default:0" json:"type"`                           // 消息类型（点赞、评论、系统消息）
	EntityID   int64  `gorm:"column:entity_id;index:idx_entity_id" json:"entity_id"`                // 实体ID（文章ID、评论ID）
	Content    string `gorm:"column:content;not null;type:text" json:"content"`                     // 消息内容
	IsRead     bool   `gorm:"column:is_read;not null;default:false" json:"is_read"`                 // 是否已读（0 未读，1 已读）
	_          struct {
		ReceiverID int64
		SenderID   int64
		Type       int32
		EntityID   int64
		EntityType int32
	} `gorm:"index:idx_receiver_id_entity_type_entity_id"`
}

const (
	MessageComment  = 1
	MessageFollow   = 2
	MessageLike     = 3
	MessageFavorite = 4
	MessageSystem   = 5
)
