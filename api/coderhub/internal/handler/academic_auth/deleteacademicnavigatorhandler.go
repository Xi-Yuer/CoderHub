package academic_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/academic_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除学术导航
func DeleteAcademicNavigatorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAcademicNavigatorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := academic_auth.NewDeleteAcademicNavigatorLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAcademicNavigator(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
