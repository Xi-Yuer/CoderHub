package school_exp_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/school_exp_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除学校经历
func DeleteSchoolExpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteSchoolExpReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := school_exp_auth.NewDeleteSchoolExpLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSchoolExp(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
