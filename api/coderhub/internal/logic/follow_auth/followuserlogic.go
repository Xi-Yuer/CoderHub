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

type FollowUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFollowUserLogic 关注用户
func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowUserLogic) FollowUser(req *types.FollowUserReq) (resp *types.FollowUserResp, err error) {
	UserID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.UserFollowService.CreateUserFollow(l.ctx, &coderhub.CreateUserFollowReq{
		FollowerId: UserID,
		FollowedId: req.FollowUserId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *FollowUserLogic) successResp() (*types.FollowUserResp, error) {
	return &types.FollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *FollowUserLogic) errorResp(err error) (*types.FollowUserResp, error) {
	return &types.FollowUserResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
