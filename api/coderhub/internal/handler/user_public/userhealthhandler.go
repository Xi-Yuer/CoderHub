package user_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/user_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 健康检查
func UserHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user_public.NewUserHealthLogic(r.Context(), svcCtx)
		resp, err := l.UserHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
