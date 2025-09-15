package banner_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/banner_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除轮播图
func DeleteBannerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := banner_public.NewDeleteBannerLogic(r.Context(), svcCtx)
		resp, err := l.DeleteBanner()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
