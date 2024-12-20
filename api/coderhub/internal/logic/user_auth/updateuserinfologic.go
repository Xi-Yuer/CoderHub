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

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateUserInfoLogic 更新用户信息
func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	// 权限验证
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	ctx := utils.SetUserMetaData(l.ctx) // 设置元数据

	userInfo, err := l.svcCtx.UserService.UpdateUserInfo(ctx, &coderhub.UpdateUserInfoRequest{
		UserId:   userID,
		Email:    req.Email,
		Nickname: req.Nickname,
		Age:      req.Age,
		Gender:   req.Gender,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, err
	}
	return l.successResp(userInfo)
}

func (l *UpdateUserInfoLogic) errorResp(err error) (*types.UpdateUserInfoResp, error) {
	return &types.UpdateUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *UpdateUserInfoLogic) successResp(data *coderhub.UpdateUserInfoResponse) (*types.UpdateUserInfoResp, error) {
	return &types.UpdateUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
