package follow_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/follow_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 健康检查
func FollowHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := follow_public.NewFollowHealthLogic(r.Context(), svcCtx)
		resp, err := l.FollowHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
