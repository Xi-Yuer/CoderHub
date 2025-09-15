package banner_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteBannerLogic 删除轮播图
func NewDeleteBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBannerLogic {
	return &DeleteBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBannerLogic) DeleteBanner() (resp *types.DeleteBannerReq, err error) {
	// todo: add your logic here and delete this line

	return
}
