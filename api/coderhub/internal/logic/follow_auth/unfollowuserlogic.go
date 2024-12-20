package follow_auth

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnfollowUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUnfollowUserLogic 取消关注
func NewUnfollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnfollowUserLogic {
	return &UnfollowUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnfollowUserLogic) UnfollowUser(req *types.UnfollowUserReq) (resp *types.UnfollowUserResp, err error) {
	UserID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.UserFollowService.DeleteUserFollow(l.ctx, &coderhub.DeleteUserFollowReq{
		FollowerId: UserID,
		FollowedId: req.FollowUserId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *UnfollowUserLogic) successResp() (*types.UnfollowUserResp, error) {
	return &types.UnfollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
	}, nil
}

func (l *UnfollowUserLogic) errorResp(err error) (*types.UnfollowUserResp, error) {
	return &types.UnfollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
