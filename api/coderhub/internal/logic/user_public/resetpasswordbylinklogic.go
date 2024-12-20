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

type ResetPasswordByLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通过链接重置密码
func NewResetPasswordByLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordByLinkLogic {
	return &ResetPasswordByLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordByLinkLogic) ResetPasswordByLink(req *types.ResetPasswordByLinkReq) (resp *types.ResetPasswordByLinkResp, err error) {
	err = utils.NewValidator().Email(req.Email).Password(req.Password).ConfirmPassword(req.Password, req.ConfirmPassword).Token(req.Token).Check()
	if err != nil {
		return l.errorResp(err), nil
	}
	_, err = l.svcCtx.UserService.ResetPasswordByLink(l.ctx, &coderhub.ResetPasswordByLinkRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Token:           req.Token,
	})
	if err != nil {
		return l.errorResp(err), nil
	}

	return l.successResp(), nil
}

func (l *ResetPasswordByLinkLogic) successResp() *types.ResetPasswordByLinkResp {
	return &types.ResetPasswordByLinkResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}
}

func (l *ResetPasswordByLinkLogic) errorResp(err error) *types.ResetPasswordByLinkResp {
	return &types.ResetPasswordByLinkResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}
}
