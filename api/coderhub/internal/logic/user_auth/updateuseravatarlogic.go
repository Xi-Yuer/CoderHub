package user_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateUserAvatarLogic 更新用户头像
func NewUpdateUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserAvatarLogic {
	return &UpdateUserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserAvatarLogic) UpdateUserAvatar(req *types.UpdateUserAvatarReq) (resp *types.UpdateUserAvatarResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	ctx := utils.SetUserMetaData(l.ctx) // 设置元数据
	_, err = l.svcCtx.UserService.UploadAvatar(ctx, &coderhub.UploadAvatarRequest{
		UserId:  userID,
		ImageId: req.Avatar,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *UpdateUserAvatarLogic) errorResp(err error) (*types.UpdateUserAvatarResp, error) {
	return &types.UpdateUserAvatarResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *UpdateUserAvatarLogic) successResp() (*types.UpdateUserAvatarResp, error) {
	return &types.UpdateUserAvatarResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
