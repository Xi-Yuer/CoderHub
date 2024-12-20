package userservicelogic

import (
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteUser 删除用户
func (l *DeleteUserLogic) DeleteUser(in *coderhub.DeleteUserRequest) (*coderhub.DeleteUserResponse, error) {
	UserInfo, err := NewGetUserInfoLogic(l.ctx, l.svcCtx).GetUserInfo(&coderhub.GetUserInfoRequest{
		UserId: in.UserId,
	})
	if err != nil {
		return nil, err
	}
	if UserInfo.UserId == 0 {
		return nil, errors.New("用户不存在")
	}
	if err := l.svcCtx.UserRepository.DeleteUser(in.UserId); err != nil {
		return nil, err
	}

	return &coderhub.DeleteUserResponse{
		Success: true,
	}, nil
}
