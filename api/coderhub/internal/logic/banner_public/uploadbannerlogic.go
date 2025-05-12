package banner_public

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUploadBannerLogic 上传轮播图
func NewUploadBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadBannerLogic {
	return &UploadBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadBannerLogic) UploadBanner() (resp *types.CreateBannerReq, err error) {
	// todo: add your logic here and delete this line

	return
}
