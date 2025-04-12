package repository

import (
	"coderhub/model"
	"context"
	"strconv"

	"gorm.io/gorm"
)

/**
	type UserSession struct {
	UserID        uint64 `gorm:"primaryKey;column:user_id;comment:用户ID"`
	PeerID        uint64 `gorm:"primaryKey;column:peer_id;comment:对方用户ID"`
	LastMessageID uint64 `gorm:"column:last_message_id;comment:最后一条消息ID"`
	UnreadCount   int    `gorm:"column:unread_count;default:0;comment:未读消息数量"`
	UpdatedAt     int64  `gorm:"column:updated_at;comment:最后更新时间戳（毫秒）"`
}
**/

// UserSessionRepository 用户会话相关操作
type UserSessionRepository interface {
	// Create 创建用户会话
	Create(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error)
	// GetUserSession 获取用户会话
	GetUserSession(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error)
	// GetUserSessionBySessionID 根据会话ID获取用户会话
	GetUserSessionBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error)
	// UpdateUserSession 更新用户会话
	UpdateUserSession(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error)
	// GetUserSessions 获取用户会话列表
	GetUserSessions(ctx context.Context, userID uint64, page, pageSize int64, sessionName string) ([]*model.UserSession, int64, error)
	// DeleteUserSession 删除用户会话
	DeleteUserSession(ctx context.Context, userSession *model.UserSession) error
}

type UserSessionRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepositoryImpl {
	return &UserSessionRepositoryImpl{DB: db}
}
func (u *UserSessionRepositoryImpl) Create(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error) {
	if err := u.DB.WithContext(ctx).Create(userSession).Error; err != nil {
		return nil, err
	}
	return userSession, nil
}
func (u *UserSessionRepositoryImpl) GetUserSession(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error) {
	var session model.UserSession
	if err := u.DB.WithContext(ctx).Where("user_id =? AND peer_id =?", userSession.UserID, userSession.PeerID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (u *UserSessionRepositoryImpl) GetUserSessionBySessionID(ctx context.Context, sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	if err := u.DB.WithContext(ctx).Where("session_id =?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (u *UserSessionRepositoryImpl) UpdateUserSession(ctx context.Context, userSession *model.UserSession) (*model.UserSession, error) {
	// 构造要更新的字段
	updateFields := map[string]interface{}{
		"unread_message_count": userSession.UnreadMessageCount,
		"last_message_id":      userSession.LastMessageID,
		"last_message_content": userSession.LastMessageContent,
	}

	// 如果 session_name 非空，则添加到更新字段中
	if userSession.SessionName != "" {
		updateFields["session_name"] = userSession.SessionName
	}

	// 执行更新
	if err := u.DB.WithContext(ctx).Model(&model.UserSession{}).
		Where("session_id = ?", userSession.SessionID).
		Updates(updateFields).Error; err != nil {
		return nil, err
	}

	return userSession, nil
}

func (u *UserSessionRepositoryImpl) GetUserSessions(ctx context.Context, userID uint64, page, pageSize int64, sessionName string) ([]*model.UserSession, int64, error) {
	var sessions []*model.UserSession

	// 构造查询条件
	query := u.DB.WithContext(ctx).Model(&model.UserSession{}).
		Where("user_id = ?", strconv.FormatUint(userID, 10))

	// 如果 sessionName 不为空，则添加模糊查询条件
	if sessionName != "" {
		query = query.Where("session_name LIKE ?", "%"+sessionName)
	}

	// 查询分页数据
	if err := query.Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	// 查询总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (u *UserSessionRepositoryImpl) DeleteUserSession(ctx context.Context, userSession *model.UserSession) error {
	if err := u.DB.WithContext(ctx).Where("user_id =? AND peer_id =?", userSession.UserID, userSession.PeerID).Delete(&model.UserSession{}).Error; err != nil {
		return err
	}
	return nil
}
