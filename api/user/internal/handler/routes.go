// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package handler

import (
	"net/http"

	"coderhub/api/user/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/user/authenticate",
				Handler: AuthenticateUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/user/check-exists",
				Handler: CheckUserExistsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/user/info",
				Handler: GetUserInfoHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/user/change-password",
				Handler: ChangePasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/user/create",
				Handler: CreateUserHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/user/delete",
				Handler: DeleteUserHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/user/reset-password",
				Handler: ResetPasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/user/update",
				Handler: UpdateUserInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
