package favorites_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/favorites_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除收藏夹内容
func DeleteFavoriteContentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteFavorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorites_auth.NewDeleteFavoriteContentLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFavoriteContent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
