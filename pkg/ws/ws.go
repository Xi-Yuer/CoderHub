package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
	"time"
)

// Upgrader is used to upgrade the HTTP connection to a WebSocket connection
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Connection struct {
	Conn         *websocket.Conn
	UserID       string
	LastPingTime int64 // 记录最近一次心跳时间
	mu           sync.Mutex
}

type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

type Hub struct {
	Connections map[string]*Connection // 通过 UserID 存储连接
	Register    chan *Connection       // 注册新连接
	Unregister  chan *Connection       // 注销连接
	Messages    chan Message           // 存储私聊消息
	mu          sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Connections: make(map[string]*Connection),
		Register:    make(chan *Connection, 10),
		Unregister:  make(chan *Connection),
		Messages:    make(chan Message, 10),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.mu.Lock()
			h.Connections[conn.UserID] = conn
			h.mu.Unlock()
		case conn := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Connections[conn.UserID]; ok {
				delete(h.Connections, conn.UserID)
				_ = conn.Conn.Close()
			}
			h.mu.Unlock()
		case msg := <-h.Messages:
			h.sendMessage(msg)
		}
	}
}

func (h *Hub) sendMessage(msg Message) {
	h.mu.Lock()
	target, ok := h.Connections[msg.ReceiverID]
	h.mu.Unlock()
	if ok {
		msgBytes, err := json.Marshal(msg)
		if err == nil {
			target.Write(msgBytes)
		}
	} else {
		logx.Infof("User %s is not online", msg.ReceiverID)
	}
}

// StartHeartbeat 服务器定期发送 Ping
func (h *Hub) StartHeartbeat(interval time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(interval) // 每 30 秒触发一次
	defer ticker.Stop()
	for range ticker.C {
		h.mu.Lock()
		now := time.Now().Unix()
		for userID, conn := range h.Connections {
			// 超过 60 秒没收到 Pong，关闭连接
			if now-conn.LastPingTime > int64(timeout.Seconds()) {
				logx.Infof("User %s disconnected due to timeout", userID)
				delete(h.Connections, userID)
				_ = conn.Conn.Close()
			} else {
				// 发送 Ping
				err := conn.Conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					logx.Errorf("Failed to send ping to %s: %v", userID, err)
					delete(h.Connections, userID)
					_ = conn.Conn.Close()
				}
			}
		}
		h.mu.Unlock()
	}
}
func (c *Connection) Write(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		logx.Errorf("Error writing message: %v", err)
	}
}
