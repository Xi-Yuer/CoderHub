package tag_public

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/tag_public"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取系统分类标签
func GetSystemTagListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := tag_public.NewGetSystemTagListLogic(r.Context(), svcCtx)
		resp, err := l.GetSystemTagList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
