package ws_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/model"
	"coderhub/pkg/ws"
	"coderhub/shared/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWebSocketConnectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateWebSocketConnectionLogic 创建WebSocket连接
func NewCreateWebSocketConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWebSocketConnectionLogic {
	return &CreateWebSocketConnectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWebSocketConnectionLogic) CreateWebSocketConnection(w http.ResponseWriter, r *http.Request, req *types.WebSocketResponse) error {
	// 升级 HTTP 连接为 WebSocket
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return err
	}
	// 获取用户身份信息 (从 JWT 解析 UserID)
	userToken := r.URL.Query().Get("token")
	if userToken == "" {
		_ = conn.Close()
		return err
	}
	userID, err := utils.ParseUserIDFromToken(userToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return err
	}
	// 创建 WebSocket 连接实例
	connection := &ws.Connection{
		Conn:         conn,
		UserID:       userID,
		LastPingTime: time.Now().Unix(), // 初始化最后 ping 时间
	}

	// 注册到 WebSocket Hub
	l.svcCtx.WsHub.Register <- connection

	// 启动一个 goroutine 监听客户端消息
	go func() {
		defer func() {
			l.svcCtx.WsHub.Unregister <- connection
			_ = conn.Close()
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logx.Errorf("Failed to read message: %v", err)
				break
			}
			// 处理接收到的消息
			// 这里可以根据需要处理消息，例如解析消息内容
			logx.Infof("Received message: %s", msg)
			var message model.PrivateMessage
			err = json.Unmarshal(msg, &message)
			if err != nil {
				logx.Errorf("Failed to unmarshal message: %v", err)
				continue
			}
			// 数据校验
			if message.SessionID == "" || message.ReceiverID == "" || message.Content == "" {
				logx.Errorf("Invalid message: %v", err)
				continue
			}
			// 发送消息到 Hub
			l.svcCtx.WsHub.Messages <- model.PrivateMessage{
				MessageID:   utils.Int2String(utils.GenID()),
				SessionID:   message.SessionID,
				SenderID:    connection.UserID,
				ReceiverID:  message.ReceiverID,
				Content:     message.Content,
				ContentType: message.ContentType,
				Status:      model.Sent,
				IsRecalled:  message.IsRecalled,
			}
		}
	}()

	return nil
}
