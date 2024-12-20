package image_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/image_auth"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 上传图片
func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := image_auth.NewUploadLogic(r.Context(), svcCtx, r)
		resp, err := l.Upload()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
