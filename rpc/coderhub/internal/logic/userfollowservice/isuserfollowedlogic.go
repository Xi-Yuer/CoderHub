package userfollowservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

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

// IsUserFollowed 检查是否关注
func (l *IsUserFollowedLogic) IsUserFollowed(in *coderhub.IsUserFollowedReq) (*coderhub.IsUserFollowedResp, error) {
	// 检查是否关注
	isFollowed, err := l.svcCtx.UserFollowRepository.IsUserFollowed(in.FollowerId, in.FollowedId)
	if err != nil {
		return nil, err
	}

	return &coderhub.IsUserFollowedResp{IsFollowed: isFollowed}, nil
}
