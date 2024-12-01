package handler

import (
	"net/http"

	"coderhub/api/TechSphere/Comment/internal/logic"
	"coderhub/api/TechSphere/Comment/internal/svc"
	"coderhub/api/TechSphere/Comment/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取某条评论的子评论列表
func GetCommentRepliesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCommentRepliesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetCommentRepliesLogic(r.Context(), svcCtx)
		resp, err := l.GetCommentReplies(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}