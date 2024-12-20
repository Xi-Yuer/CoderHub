package academic_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/academic_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 取消点赞学术导航
func CancelLikeAcademicNavigatorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CancelLikeAcademicNavigatorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := academic_auth.NewCancelLikeAcademicNavigatorLogic(r.Context(), svcCtx)
		resp, err := l.CancelLikeAcademicNavigator(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
