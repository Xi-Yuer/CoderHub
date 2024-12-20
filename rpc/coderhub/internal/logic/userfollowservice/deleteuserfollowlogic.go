package userfollowservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserFollowLogic {
	return &DeleteUserFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteUserFollow 删除用户关注关系
func (l *DeleteUserFollowLogic) DeleteUserFollow(in *coderhub.DeleteUserFollowReq) (*coderhub.DeleteUserFollowResp, error) {
	l.Logger.Info("DeleteUserFollow", in.FollowedId, in.FollowerId)
	err := l.svcCtx.UserFollowRepository.DeleteUserFollow(&model.UserFollow{
		FollowerID: in.FollowerId,
		FollowedID: in.FollowedId,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteUserFollowResp{
		Success: true,
	}, nil
}
