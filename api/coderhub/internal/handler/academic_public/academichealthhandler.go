package academic_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/academic_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 健康检查
func AcademicHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := academic_public.NewAcademicHealthLogic(r.Context(), svcCtx)
		resp, err := l.AcademicHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
