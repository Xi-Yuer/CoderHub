package user_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if err := utils.NewValidator().Username(req.Username).Password(req.Password).Check(); err != nil {
		return &types.LoginResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: "",
		}, nil
	}

	var authorize *coderhub.AuthorizeResponse

	if authorize, err = l.svcCtx.UserService.Authorize(l.ctx, &coderhub.AuthorizeRequest{
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		return &types.LoginResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: "",
		}, nil
	}

	return &types.LoginResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: authorize.Token,
	}, nil
}
