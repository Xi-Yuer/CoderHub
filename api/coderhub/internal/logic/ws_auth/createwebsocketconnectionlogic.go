package ws_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/pkg/ws"
	"coderhub/shared/utils"
	"context"
	"fmt"
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
	fmt.Println("UserToken:", userToken)
	if userToken == "" {
		_ = conn.Close()
		return err
	}
	userID, err := utils.ParseUserIDFromToken(userToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return err
	}

	fmt.Println("UserID:", userID)

	// 创建 WebSocket 连接实例
	connection := &ws.Connection{
		Conn:         conn,
		UserID:       userID,
		LastPingTime: time.Now().Unix(), // 初始化最后 ping 时间
	}

	// 注册到 WebSocket Hub
	l.svcCtx.WsHub.Register <- connection

	return nil
}
