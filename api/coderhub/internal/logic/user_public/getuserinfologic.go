package user_public

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

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
	UserInfo, err := l.svcCtx.UserService.GetUserInfo(l.ctx, &coderhub.GetUserInfoRequest{UserId: utils.String2Int(req.Id), RequestUserId: utils.String2Int(req.RequestUserID)})
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
			Id:           utils.Int2String(UserInfo.UserId),
			Username:     UserInfo.UserName,
			Nickname:     UserInfo.NickName,
			Email:        UserInfo.Email,
			Phone:        UserInfo.Phone,
			Avatar:       UserInfo.Avatar,
			Gender:       UserInfo.Gender,
			Level:        UserInfo.Level,
			Age:          UserInfo.Age,
			Status:       UserInfo.Status,
			IsAdmin:      UserInfo.IsAdmin,
			CreateAt:     UserInfo.CreatedAt,
			UpdateAt:     UserInfo.UpdatedAt,
			FollowCount:  UserInfo.FollowerCount,
			FansCount:    UserInfo.FollowCount,
			IsFollowed:   UserInfo.IsFollowed,
			ArticleCount: int64(UserInfo.ArticleCount),
		},
	}, nil
}
