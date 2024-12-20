package user_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 健康检查
func NewUserHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHealthLogic {
	return &UserHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserHealthLogic) UserHealth() (resp *types.HealthResp, err error) {
	return &types.HealthResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
