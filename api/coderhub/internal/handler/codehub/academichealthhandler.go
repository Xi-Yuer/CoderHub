package codehub

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/codehub"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 健康检查
func AcademicHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := codehub.NewAcademicHealthLogic(r.Context(), svcCtx)
		resp, err := l.AcademicHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}