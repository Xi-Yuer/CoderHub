package user_auth

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置密码
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp *types.ResetPasswordResp, err error) {
	if err := utils.NewValidator().ConfirmPassword(req.NewPassword, req.ConfirmPassword).Password(req.NewPassword).Check(); err != nil {
		return &types.ResetPasswordResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: false,
		}, nil
	}

	_, err = utils.GetUserID(l.ctx)
	if err != nil {
		return &types.ResetPasswordResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
		}, nil
	}
	_ = utils.SetUserMetaData(l.ctx) // 设置元数据

	return
}
