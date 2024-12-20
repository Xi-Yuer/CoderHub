package user_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUserInfoLogic 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	UserInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{UserId: req.Id})
	if err != nil {
		return &types.GetUserInfoResp{
			Response: types.Response{
				Code:    conf.HttpCode.HttpBadRequest,
				Message: err.Error(),
			},
			Data: nil,
		}, nil
	}
	return &types.GetUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: &types.UserInfo{
			Id:       UserInfo.UserId,
			Username: UserInfo.UserName,
			Nickname: UserInfo.NickName,
			Email:    UserInfo.Email,
			Phone:    UserInfo.Phone,
			Avatar:   UserInfo.Avatar,
			Gender:   UserInfo.Gender,
			Age:      UserInfo.Age,
			Status:   UserInfo.Status,
			IsAdmin:  UserInfo.IsAdmin,
			CreateAt: UserInfo.CreatedAt,
			UpdateAt: UserInfo.UpdatedAt,
		},
	}, nil
}
