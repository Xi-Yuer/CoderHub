package userservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserExistsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserExistsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserExistsLogic {
	return &CheckUserExistsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckUserExists 检查用户是否存在
func (l *CheckUserExistsLogic) CheckUserExists(in *coderhub.CheckUserExistsRequest) (*coderhub.CheckUserExistsResponse, error) {
	UserInfo, err := NewGetUserInfoByUsernameLogic(l.ctx, l.svcCtx).GetUserInfoByUsername(&coderhub.GetUserInfoByUsernameRequest{Username: in.Username})

	if UserInfo == nil || err != nil {
		return &coderhub.CheckUserExistsResponse{
			Exists: false,
		}, nil
	}
	return &coderhub.CheckUserExistsResponse{
		Exists: true,
	}, nil
}
