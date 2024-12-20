package follow_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetFollowListLogic 获取关注列表
func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFollowListLogic) GetFollowList(req *types.GetFollowListReq) (resp *types.GetFollowListResp, err error) {
	userFollowsResp, err := l.svcCtx.UserFollowService.GetUserFollows(l.ctx, &coderhub.GetUserFollowsReq{
		FollowerId: req.UserId,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp(userFollowsResp)
}

func (l *GetFollowListLogic) successResp(userFollowsResp *coderhub.GetUserFollowsResp) (*types.GetFollowListResp, error) {
	userFollowsList := make([]types.UserInfo, 0, len(userFollowsResp.UserFollows))
	for _, userFollow := range userFollowsResp.UserFollows {
		userFollowsList = append(userFollowsList, types.UserInfo{
			Id:       userFollow.UserId,
			Username: userFollow.UserName,
			Nickname: userFollow.NickName,
			Email:    userFollow.Email,
			Phone:    userFollow.Phone,
			Avatar:   userFollow.Avatar,
			Gender:   userFollow.Gender,
			Age:      userFollow.Age,
			Status:   userFollow.Status,
			IsAdmin:  userFollow.IsAdmin,
			CreateAt: userFollow.CreatedAt,
			UpdateAt: userFollow.UpdatedAt,
		})
	}
	return &types.GetFollowListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.FollowList{
			List:  userFollowsList,
			Total: userFollowsResp.Total,
		},
	}, nil
}

func (l *GetFollowListLogic) errorResp(err error) (*types.GetFollowListResp, error) {
	return &types.GetFollowListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
