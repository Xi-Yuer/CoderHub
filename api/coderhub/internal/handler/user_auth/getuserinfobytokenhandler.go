package user_auth

import (
	"net/http"

	"coderhub/api/coderhub/internal/logic/user_auth"
	"coderhub/api/coderhub/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 根据用户的token获取用户信息
func GetUserInfoByTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user_auth.NewGetUserInfoByTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetUserInfoByToken()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
