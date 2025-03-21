package userfollowservicelogic

import (
	userservicelogic "coderhub/rpc/coderhub/internal/logic/userservice"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFansLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFansLogic {
	return &GetUserFansLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserFans 获取用户粉丝列表
func (l *GetUserFansLogic) GetUserFans(in *coderhub.GetUserFansReq) (*coderhub.GetUserFansResp, error) {
	// 获取粉丝列表
	userFollows, err := l.svcCtx.UserFollowRepository.GetUserFans(in.FollowedId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	// 构建粉丝IDs
	userIDs := make([]int64, 0, len(userFollows))
	for _, userFollow := range userFollows {
		userIDs = append(userIDs, userFollow.FollowerID)
	}
	// 批量获取用户信息
	batchGetUserByIdService := userservicelogic.NewBatchGetUserByIDLogic(l.ctx, l.svcCtx)
	users, err := batchGetUserByIdService.BatchGetUserByID(&coderhub.BatchGetUserByIDRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}

	// 查询 userID 是否关注了 userIDs 这批用户
	followedUserIDs, _ := l.svcCtx.UserFollowRepository.GetUserAndUserHasFollow(in.RequestUserId, userIDs)

	// 把 followedUserIDs 转换为 map，提高查询效率
	followedMap := make(map[int64]bool, len(followedUserIDs))
	for _, id := range followedUserIDs {
		followedMap[id] = true
	}

	userInfos := make([]*coderhub.UserInfo, 0, len(users.UserInfos))
	for _, userInfo := range users.UserInfos {
		userInfos = append(userInfos, &coderhub.UserInfo{
			UserId:     userInfo.UserId,
			UserName:   userInfo.UserName,
			Avatar:     userInfo.Avatar,
			Email:      userInfo.Email,
			NickName:   userInfo.NickName,
			IsAdmin:    userInfo.IsAdmin,
			Status:     userInfo.Status,
			IsFollowed: followedMap[userInfo.UserId], // 如果在 map 里，就是已关注,
			CreatedAt:  userInfo.CreatedAt,
			UpdatedAt:  userInfo.UpdatedAt,
		})
	}

	return &coderhub.GetUserFansResp{
		UserFans: userInfos,
		Total:    int64(len(userFollows)),
	}, nil
}
