package userfollowservicelogic

import (
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMutualFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMutualFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMutualFollowsLogic {
	return &GetMutualFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMutualFollows 获取互相关注列表
func (l *GetMutualFollowsLogic) GetMutualFollows(in *coderhub.GetMutualFollowsReq) (*coderhub.GetMutualFollowsResp, error) {
	// 获取互相关注列表(ID)
	mutualFollows, err := l.svcCtx.UserFollowRepository.GetMutualFollows(in.UserId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	// 获取互相关注列表(User)
	userIDs := make([]int64, 0, len(mutualFollows))
	for _, follow := range mutualFollows {
		userIDs = append(userIDs, follow.FollowedID)
	}
	batchGetUserByIdService := userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx)
	users, err := batchGetUserByIdService.BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}
	userInfos := make([]*coderhub.UserFollowInfo, 0, len(users.UserInfos))
	for _, userInfo := range users.UserInfos {
		userInfos = append(userInfos, &coderhub.UserFollowInfo{
			Id:       userInfo.UserId,
			Username: userInfo.UserName,
			Avatar:   userInfo.Avatar,
		})
	}

	return &coderhub.GetMutualFollowsResp{
		MutualFollows: userInfos,
		Total:         int64(len(mutualFollows)),
	}, nil
}
