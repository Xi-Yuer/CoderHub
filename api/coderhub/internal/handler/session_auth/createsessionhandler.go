package session_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/session_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建一个会话
func CreateSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSessionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := session_auth.NewCreateSessionLogic(r.Context(), svcCtx)
		resp, err := l.CreateSession(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
