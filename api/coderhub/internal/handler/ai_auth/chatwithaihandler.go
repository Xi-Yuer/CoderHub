package ai_auth

import (
	"coderhub/api/coderhub/internal/logic/ai_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// ChatWithAIHandler 与AI对话 - 流式返回
func ChatWithAIHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatWithAIReq
		w.Header().Set(`Content-Type`, `text/event-stream;charset=utf-8`)
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := ai_auth.NewChatWithAILogic(r.Context(), svcCtx)
		l.ChatWithAI(w, &req)
	}
}
