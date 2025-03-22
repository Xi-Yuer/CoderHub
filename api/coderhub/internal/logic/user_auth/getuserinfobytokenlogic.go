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

type GetUserInfoByTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUserInfoByTokenLogic 根据用户的token获取用户信息
func NewGetUserInfoByTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByTokenLogic {
	return &GetUserInfoByTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByTokenLogic) GetUserInfoByToken() (resp *types.GetUserInfoResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	user, err := l.svcCtx.UserService.GetUserInfo(utils.SetUserMetaData(l.ctx), &coderhub.GetUserInfoRequest{
		UserId:        userID,
		RequestUserId: userID,
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp(&types.UserInfo{
		Id:           utils.Int2String(user.UserId),
		Username:     user.UserName,
		Nickname:     user.NickName,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Age:          user.Age,
		Status:       user.Status,
		IsAdmin:      user.IsAdmin,
		CreateAt:     user.CreatedAt,
		UpdateAt:     user.UpdatedAt,
		FollowCount:  user.FollowerCount,
		FansCount:    user.FollowCount,
		IsFollowed:   user.IsFollowed,
		ArticleCount: int64(user.ArticleCount),
	})
}

func (l *GetUserInfoByTokenLogic) successResp(user *types.UserInfo) (*types.GetUserInfoResp, error) {
	return &types.GetUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: user,
	}, nil
}

func (l *GetUserInfoByTokenLogic) errorResp(err error) (*types.GetUserInfoResp, error) {
	return &types.GetUserInfoResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
