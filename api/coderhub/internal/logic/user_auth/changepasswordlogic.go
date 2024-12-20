package user_auth

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改密码
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.UpdatePasswordReq) (resp *types.UpdatePasswordResp, err error) {
	if err := utils.NewValidator().Password(req.OldPassword).Password(req.NewPassword).Check(); err != nil {
		return &types.UpdatePasswordResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: false,
		}, nil
	}

	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return &types.UpdatePasswordResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
		}, nil
	}
	ctx := utils.SetUserMetaData(l.ctx) // 设置元数据
	response, err := l.svcCtx.UserService.ChangePassword(ctx, &coderhub.ChangePasswordRequest{
		UserId:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		return &types.UpdatePasswordResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
		}, nil
	}

	// 返回成功响应
	return &types.UpdatePasswordResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: response.Success,
	}, nil
}
