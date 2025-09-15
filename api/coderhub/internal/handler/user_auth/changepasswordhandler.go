package user_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/user_auth"
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 修改密码
func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdatePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user_auth.NewChangePasswordLogic(r.Context(), svcCtx)
		resp, err := l.ChangePassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
