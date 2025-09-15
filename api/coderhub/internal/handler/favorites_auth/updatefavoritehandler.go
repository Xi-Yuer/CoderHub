package favorites_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/favorites_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新收藏夹
func UpdateFavoriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFavorFoldReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := favorites_auth.NewUpdateFavoriteLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFavorite(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
