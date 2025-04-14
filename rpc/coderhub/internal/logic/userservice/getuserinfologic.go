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
	// 定义两个 channel 用于接收查询结果和错误
	userChan := make(chan *model.User, 1)
	errChan := make(chan error, 2)
	isFollowedChan := make(chan bool, 1)

	// 并行查询用户信息
	go func() {
		user, err := l.svcCtx.UserRepository.GetUserByID(in.UserId)
		if err != nil {
			errChan <- fmt.Errorf("failed to get user by ID %d: %w", in.UserId, err)
			return
		}
		userChan <- user
	}()

	// 并行查询关注状态
	go func() {
		var user *model.User
		// 等待用户信息查询完成
		select {
		case user = <-userChan:
		case err := <-errChan:
			errChan <- err
			return
		}

		isFollowed, err := l.svcCtx.UserFollowRepository.IsUserFollowed(in.RequestUserId, user.ID)
		if err != nil {
			errChan <- fmt.Errorf("failed to check if user %d is followed by %d: %w", user.ID, in.RequestUserId, err)
			return
		}
		isFollowedChan <- isFollowed
	}()

	// 等待结果
	var user *model.User
	var isFollowed bool
	select {
	case err := <-errChan:
		l.Errorf("%v", err)
		return nil, err
	case user = <-userChan:
		select {
		case err := <-errChan:
			l.Errorf("%v", err)
			return nil, err
		case isFollowed = <-isFollowedChan:
		}
	}

	return buildUserInfoResponse(user, isFollowed), nil
}

func buildUserInfoResponse(user *model.User, isFollowed bool) *coderhub.UserInfo {
	return &coderhub.UserInfo{
		UserId:        user.ID,
		UserName:      user.UserName,
		Avatar:        user.Avatar.String,
		Email:         user.Email.String,
		Password:      user.Password,
		Gender:        user.Gender,
		Age:           user.Age,
		Phone:         user.Phone.String,
		NickName:      user.NickName.String,
		IsAdmin:       user.IsAdmin,
		Status:        user.Status,
		CreatedAt:     user.CreatedAt.Unix(),
		UpdatedAt:     user.UpdatedAt.Unix(),
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		ArticleCount:  int32(user.ArticleCount),
		IsFollowed:    isFollowed,
	}
}
