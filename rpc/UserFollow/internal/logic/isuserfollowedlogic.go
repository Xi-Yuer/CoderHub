package logic

import (
	"context"

	"coderhub/rpc/UserFollow/internal/svc"
	"coderhub/rpc/UserFollow/user_follow"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsUserFollowedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsUserFollowedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsUserFollowedLogic {
	return &IsUserFollowedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查是否关注
func (l *IsUserFollowedLogic) IsUserFollowed(in *user_follow.IsUserFollowedReq) (*user_follow.IsUserFollowedResp, error) {
	// 检查是否关注
	isFollowed, err := l.svcCtx.UserFollowRepository.IsUserFollowed(in.FollowerId, in.FollowedId)
	if err != nil {
		return nil, err
	}

	return &user_follow.IsUserFollowedResp{IsFollowed: isFollowed}, nil
}
