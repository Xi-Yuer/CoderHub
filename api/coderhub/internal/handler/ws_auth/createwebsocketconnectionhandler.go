package ws_auth

import (
	"coderhub/api/coderhub/internal/logic/ws_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// CreateWebSocketConnectionHandler 创建WebSocket连接
func CreateWebSocketConnectionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebSocketResponse
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := ws_auth.NewCreateWebSocketConnectionLogic(r.Context(), svcCtx)
		err := l.CreateWebSocketConnection(w, r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
