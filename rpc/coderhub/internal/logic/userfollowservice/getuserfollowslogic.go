package userfollowservicelogic

import (
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowsLogic {
	return &GetUserFollowsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserFollows 获取用户关注列表
func (l *GetUserFollowsLogic) GetUserFollows(in *coderhub.GetUserFollowsReq) (*coderhub.GetUserFollowsResp, error) {
	// 获取关注列表
	userFollows, err := l.svcCtx.UserFollowRepository.GetUserFollows(in.FollowerId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	// 构建用户IDs
	userIDs := make([]int64, 0, len(userFollows))
	for _, userFollow := range userFollows {
		userIDs = append(userIDs, userFollow.FollowedID)
	}
	// 批量获取用户信息
	batchGetUserByIdService := userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx)
	users, err := batchGetUserByIdService.BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}
	// 转换为user_follow.UserInfo
	userInfos := make([]*coderhub.UserInfo, 0, len(users.UserInfos))
	for _, userInfo := range users.UserInfos {
		userInfos = append(userInfos, &coderhub.UserInfo{
			UserId:    userInfo.UserId,
			UserName:  userInfo.UserName,
			Avatar:    userInfo.Avatar,
			Email:     userInfo.Email,
			Gender:    0,
			Age:       0,
			Phone:     "",
			NickName:  userInfo.NickName,
			IsAdmin:   userInfo.IsAdmin,
			Status:    userInfo.Status,
			CreatedAt: userInfo.CreatedAt,
			UpdatedAt: userInfo.UpdatedAt,
		})
	}
	l.Logger.Info("userInfos_length", len(userInfos))
	return &coderhub.GetUserFollowsResp{
		UserFollows: userInfos,
		Total:       int64(len(userFollows)),
	}, nil
}
