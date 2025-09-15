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

	user, err := l.svcCtx.UserService.UpdateUserInfo(utils.SetUserMetaData(l.ctx), &coderhub.UpdateUserInfoRequest{
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
	return l.successResp(user)
}

func (l *UpdateUserInfoLogic) errorResp(err error) (*types.UpdateUserInfoResp, error) {
	return &types.UpdateUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *UpdateUserInfoLogic) successResp(data *coderhub.UpdateUserInfoResponse) (*types.UpdateUserInfoResp, error) {
	return &types.UpdateUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.UserInfo{
			Id:       utils.Int2String(data.UserInfo.UserId),
			Username: data.UserInfo.UserName,
			Nickname: data.UserInfo.NickName,
			Email:    data.UserInfo.Email,
			Phone:    data.UserInfo.Phone,
			Avatar:   data.UserInfo.Avatar,
			Gender:   data.UserInfo.Gender,
			Age:      data.UserInfo.Age,
			Status:   data.UserInfo.Status,
			IsAdmin:  data.UserInfo.IsAdmin,
			CreateAt: data.UserInfo.CreatedAt,
			UpdateAt: data.UserInfo.UpdatedAt,
		},
	}, nil
}
