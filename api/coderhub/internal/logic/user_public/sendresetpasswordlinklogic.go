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

type SendResetPasswordLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送重置密码链接
func NewSendResetPasswordLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetPasswordLinkLogic {
	return &SendResetPasswordLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendResetPasswordLinkLogic) SendResetPasswordLink(req *types.SendResetPasswordLinkReq) (resp *types.SendResetPasswordLinkResp, err error) {
	err = utils.NewValidator().Email(req.Email).Check()
	if err != nil {
		return l.errorResp(err), nil
	}
	_, err = l.svcCtx.UserService.ResetPassword(l.ctx, &coderhub.ResetPasswordRequest{
		Email: req.Email,
	})
	if err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

func (l *SendResetPasswordLinkLogic) successResp() *types.SendResetPasswordLinkResp {
	return &types.SendResetPasswordLinkResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}
}

func (l *SendResetPasswordLinkLogic) errorResp(err error) *types.SendResetPasswordLinkResp {
	return &types.SendResetPasswordLinkResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}
