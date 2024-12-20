package image_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/image_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 健康检查
func ImageHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := image_public.NewImageHealthLogic(r.Context(), svcCtx)
		resp, err := l.ImageHealth()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
