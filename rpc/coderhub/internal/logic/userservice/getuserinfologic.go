package userservicelogic

import (
	"coderhub/model"
	"context"
	"fmt"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserInfo 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *coderhub.GetUserInfoRequest) (*coderhub.UserInfo, error) {
	var User *model.User
	User, err := l.svcCtx.UserRepository.GetUserByID(in.UserId)
	if err != nil {
		return nil, err
	}

	isUserFollowed, err := l.svcCtx.UserFollowRepository.IsUserFollowed(in.RequestUserId, User.ID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("isUserFollowed: %v\n", isUserFollowed)
	fmt.Printf("RequestUserId: %v\n", in.RequestUserId)
	fmt.Printf("User.ID: %v\n", User.ID)
	fmt.Printf("User: %v\n", User)

	return &coderhub.UserInfo{
		UserId:        User.ID,
		UserName:      User.UserName,
		Avatar:        User.Avatar.String,
		Email:         User.Email.String,
		Password:      User.Password,
		Gender:        User.Gender,
		Age:           User.Age,
		Phone:         User.Phone.String,
		NickName:      User.NickName.String,
		IsAdmin:       User.IsAdmin,
		Status:        User.Status,
		CreatedAt:     User.CreatedAt.Unix(),
		UpdatedAt:     User.UpdatedAt.Unix(),
		FollowCount:   User.FollowCount,
		FollowerCount: User.FollowerCount,
		ArticleCount:  int32(User.ArticleCount),
		IsFollowed:    isUserFollowed,
	}, nil
}
