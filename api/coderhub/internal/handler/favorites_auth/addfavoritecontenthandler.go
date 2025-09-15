package favorites_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/favorites_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 添加收藏内容
func AddFavoriteContentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateFavorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorites_auth.NewAddFavoriteContentLogic(r.Context(), svcCtx)
		resp, err := l.AddFavoriteContent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
