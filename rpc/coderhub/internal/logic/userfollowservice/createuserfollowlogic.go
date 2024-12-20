package userfollowservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserFollowLogic {
	return &CreateUserFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateUserFollow 创建用户关注关系
func (l *CreateUserFollowLogic) CreateUserFollow(in *coderhub.CreateUserFollowReq) (*coderhub.CreateUserFollowResp, error) {
	if err := l.svcCtx.UserFollowRepository.CreateUserFollow(&model.UserFollow{
		FollowerID: in.FollowerId,
		FollowedID: in.FollowedId,
	}); err != nil {
		return nil, err
	}

	return &coderhub.CreateUserFollowResp{
		Success: true,
	}, nil
}
