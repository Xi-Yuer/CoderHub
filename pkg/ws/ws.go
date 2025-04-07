package ws

import (
	"coderhub/model"
	"coderhub/repository"
	"coderhub/shared/storage"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// Upgrader 用于将 HTTP 连接升级为 WebSocket 连接
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Connection 表示一个 WebSocket 连接
type Connection struct {
	Conn         *websocket.Conn
	UserID       string
	LastPingTime int64 // 记录最近一次心跳时间
	mu           sync.Mutex
}

// Hub 管理所有的 WebSocket 连接和消息
type Hub struct {
	Connections              map[string]*Connection    // 通过 UserID 存储连接
	Register                 chan *Connection          // 注册新连接
	Unregister               chan *Connection          // 注销连接
	Messages                 chan model.PrivateMessage // 存储私聊消息
	mu                       sync.Mutex
	PrivateMessageRepository repository.PrivateMessageRepository
	UserSessionRepository    repository.UserSessionRepository
}

// NewHub 创建一个新的 Hub 实例
func NewHub() (*Hub, error) {
	redisDB, err := storage.NewRedisDB(storage.DefaultConfig())
	if err != nil {
		return nil, err
	}
	sql := storage.NewGorm()
	return &Hub{
		Connections:              make(map[string]*Connection),
		Register:                 make(chan *Connection, 10),
		Unregister:               make(chan *Connection),
		Messages:                 make(chan model.PrivateMessage, 10),
		PrivateMessageRepository: repository.NewPrivateMessageRepository(sql, redisDB),
		UserSessionRepository:    repository.NewUserSessionRepository(sql),
	}, nil
}

// Run 启动 Hub 并处理连接注册、注销和消息发送
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.handleRegister(conn)
		case conn := <-h.Unregister:
			h.handleUnregister(conn)
		case msg := <-h.Messages:
			h.sendMessage(msg)
		}
	}
}

// handleRegister 处理新连接的注册
func (h *Hub) handleRegister(conn *Connection) {
	h.mu.Lock()
	h.Connections[conn.UserID] = conn
	h.mu.Unlock()

	// 设置 WebSocket pong 处理器
	conn.Conn.SetPongHandler(func(appData string) error {
		conn.HandlePong()
		return nil
	})

	// 标记用户为在线
	if err := h.PrivateMessageRepository.MarkUserOnline(context.Background(), conn.UserID); err != nil {
		logx.Errorf("Failed to mark user %s online: %v", conn.UserID, err)
		return
	}

	h.sendOfflineMessages(conn, conn.UserID)
}

// handleUnregister 处理连接的注销
func (h *Hub) handleUnregister(conn *Connection) {
	h.mu.Lock()
	if _, ok := h.Connections[conn.UserID]; ok {
		delete(h.Connections, conn.UserID)
		_ = conn.Conn.Close()
	}
	h.mu.Unlock()

	if err := h.PrivateMessageRepository.MarkUserOffline(context.Background(), conn.UserID); err != nil {
		logx.Errorf("Failed to mark user %s offline: %v", conn.UserID, err)
	}
}

// sendOfflineMessages 发送用户的离线消息
func (h *Hub) sendOfflineMessages(conn *Connection, userID string) {
	// TODO:当用户前端点击会话之后，需要修改会话的未读消息数量为 0，这里需要在前端进行处理
	offlineMessages, err := h.PrivateMessageRepository.GetOfflineMessages(context.Background(), conn.UserID)
	if err != nil {
		logx.Errorf("Failed to get offline messages for user %s: %v", userID, err)
		return
	}
	for _, msg := range offlineMessages {
		msgBytes, err := json.Marshal(msg)
		if err == nil {
			conn.Write(msgBytes)
		}
		msg.Status = "read"
		if err := h.PrivateMessageRepository.UpdatePrivateMessage(context.Background(), msg); err != nil {
			logx.Errorf("Failed to update message status: %v", err)
		}
	}
}

// sendMessage 发送消息并处理消息持久化
func (h *Hub) sendMessage(msg model.PrivateMessage) {
	h.mu.Lock()
	target, ok := h.Connections[msg.ReceiverID]
	h.mu.Unlock()

	// 检查会话是否存在
	session, err := h.UserSessionRepository.GetUserSession(context.Background(), &model.UserSession{
		SessionID: msg.SessionID,
		UserID:    msg.SenderID,
		PeerID:    msg.ReceiverID,
	})
	if err != nil {
		logx.Errorf("Failed to get user session: %v", err)
	}
	if session == nil {
		// 会话信息不存在，返回错误
		logx.Errorf("User session not found for sender %s and receiver %s", msg.SenderID, msg.ReceiverID)
		return
	} else {
		// 更新发送者会话信息
		session.LastMessageID = msg.MessageID
		session.LastMessageContent = msg.Content
		session.UpdatedAt = time.Now()
		session.UnreadMessageCount = 0
		if _, err := h.UserSessionRepository.UpdateUserSession(context.Background(), session); err != nil {
			logx.Errorf("Failed to update sender's user session: %v", err)
		}
		// 更新接收者会话信息
		receiverSession, err := h.UserSessionRepository.GetUserSession(context.Background(), &model.UserSession{
			UserID: msg.ReceiverID,
			PeerID: msg.SenderID,
		})
		if err != nil {
			logx.Errorf("Failed to get receiver's user session: %v", err)
		} else {
			receiverSession.LastMessageID = msg.MessageID
			receiverSession.LastMessageContent = msg.Content
			if !ok { // 接收者不在线，未读消息数量加 1
				receiverSession.UnreadMessageCount++
			} else {
				receiverSession.UnreadMessageCount = 0
			}
			receiverSession.UpdatedAt = time.Now()
			if _, err := h.UserSessionRepository.UpdateUserSession(context.Background(), receiverSession); err != nil {
				logx.Errorf("Failed to update receiver's user session: %v", err)
			}
		}
	}

	if ok {
		msg.Status = "read"
		if err := h.sendAndSaveMessage(target, msg); err != nil {
			logx.Errorf("Failed to send and save message: %v", err)
		}
	} else {
		logx.Infof("User %s is not online", msg.ReceiverID)
		if err := h.saveMessage(msg); err != nil {
			logx.Errorf("Failed to save offline message: %v", err)
		}
	}
	if msg.Content == "pong" {
		h.handlePongMessage(msg.SenderID)
	}
}

// sendAndSaveMessage 发送消息并保存到数据库
func (h *Hub) sendAndSaveMessage(target *Connection, msg model.PrivateMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	target.Write(msgBytes)
	return h.saveMessage(msg)
}

// saveMessage 保存消息到数据库
func (h *Hub) saveMessage(msg model.PrivateMessage) error {
	// 更新会话信息
	_, err := h.UserSessionRepository.UpdateUserSession(context.Background(), &model.UserSession{
		SessionID:          msg.SessionID,
		SessionName:        msg.SessionID,
		UserID:             msg.SenderID,
		PeerID:             msg.ReceiverID,
		LastMessageID:      msg.MessageID,
		LastMessageContent: msg.Content,
	})
	if err != nil {
		return err
	}
	message := &model.PrivateMessage{
		MessageID:   msg.MessageID,
		SessionID:   msg.SessionID,
		SenderID:    msg.SenderID,
		ReceiverID:  msg.ReceiverID,
		Content:     msg.Content,
		ContentType: msg.ContentType,
		Status:      msg.Status,
		IsRecalled:  msg.IsRecalled,
	}
	return h.PrivateMessageRepository.Create(context.Background(), message)
}

// handlePongMessage 处理 pong 消息
func (h *Hub) handlePongMessage(senderID string) {
	h.mu.Lock()
	conn, ok := h.Connections[senderID]
	h.mu.Unlock()
	if ok {
		conn.HandlePong()
	}
}

// StartHeartbeat 服务器定期发送 Ping
func (h *Hub) StartHeartbeat(interval time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		h.checkAndSendPing(timeout)
	}
}

// checkAndSendPing 发送 Ping 并检查超时
func (h *Hub) checkAndSendPing(timeout time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now().Unix()
	for userID, conn := range h.Connections {
		if now-conn.LastPingTime > int64(timeout.Seconds()) {
			logx.Infof("User %s disconnected due to timeout", userID)
			delete(h.Connections, userID)
			_ = conn.Conn.Close()
		} else {
			if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logx.Errorf("Failed to send ping to %s: %v", userID, err)
				delete(h.Connections, userID)
				_ = conn.Conn.Close()
			}
		}
	}
}

// Write 向连接写入消息
func (c *Connection) Write(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
		logx.Errorf("Error writing message: %v", err)
	}
}

// HandlePong 处理 pong 消息，将用户的最后一次 ping 时间更新为当前时间
func (c *Connection) HandlePong() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastPingTime = time.Now().Unix()
	logx.Infof("Received pong from user %s", c.UserID)
}
