package userservicelogic

import (
	"coderhub/model"
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByUsernameLogic {
	return &GetUserInfoByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByUsernameLogic) GetUserInfoByUsername(in *coderhub.GetUserInfoByUsernameRequest) (*coderhub.UserInfo, error) {
	var User *model.User
	User, err := l.svcCtx.UserRepository.GetUserByName(in.Username)
	if err != nil {
		return nil, err
	}
	if User == nil {
		return nil, errors.New("用户不存在")
	}
	return &coderhub.UserInfo{
		UserId:    User.ID,
		UserName:  User.UserName,
		Avatar:    User.Avatar.String,
		Email:     User.Email.String,
		Password:  User.Password,
		Gender:    User.Gender,
		Age:       User.Age,
		Phone:     User.Phone.String,
		NickName:  User.NickName.String,
		IsAdmin:   User.IsAdmin,
		Status:    User.Status,
		CreatedAt: User.CreatedAt.Unix(),
		UpdatedAt: User.UpdatedAt.Unix(),
	}, nil
}
