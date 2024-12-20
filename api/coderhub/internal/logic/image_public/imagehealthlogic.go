package image_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 健康检查
func NewImageHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageHealthLogic {
	return &ImageHealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageHealthLogic) ImageHealth() (resp *types.HealthResp, err error) {
	return &types.HealthResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
	}, nil
}
