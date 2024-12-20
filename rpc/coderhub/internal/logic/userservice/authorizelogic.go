package userservicelogic

import (
	"coderhub/shared/security"
	"coderhub/shared/utils"
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Authorize 授权
func (l *AuthorizeLogic) Authorize(in *coderhub.AuthorizeRequest) (*coderhub.AuthorizeResponse, error) {
	if err := utils.NewValidator().Username(in.Username).Password(in.Password).Check(); err != nil {
		return nil, err
	}

	UserInfo, err := NewGetUserInfoByUsernameLogic(l.ctx, l.svcCtx).GetUserInfoByUsername(&coderhub.GetUserInfoByUsernameRequest{Username: in.Username})
	if err != nil {
		return nil, err
	}

	if !security.CompareHashAndPassword(UserInfo.Password, in.Password) {
		return nil, errors.New("密码错误")
	}

	if authorization, err := security.GenerateAuthorization(UserInfo.UserId); err != nil {
		return nil, err
	} else {
		return &coderhub.AuthorizeResponse{
			Token: authorization,
		}, nil
	}
}
