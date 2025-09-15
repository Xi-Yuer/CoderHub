package academic_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/academic_public"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取学术导航
func GetAcademicNavigatorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAcademicNavigatorReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := academic_public.NewGetAcademicNavigatorLogic(r.Context(), svcCtx)
		resp, err := l.GetAcademicNavigator(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
