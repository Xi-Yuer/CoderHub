package banner_public

import (
	"coderhub/conf"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListBannerLogic 获取轮播图列表
func NewListBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBannerLogic {
	return &ListBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListBannerLogic) ListBanner() (resp *types.GetBannerListResp, err error) {

	var respData *types.GetBannerListResp

	respData = &types.GetBannerListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: []*types.Banner{
			{
				ID:        "1",
				Title:     "一分耕耘一分收获",
				ImageUrl:  "https://xiyuer.club/minio/coderhub/17470312824075",
				LinkUrl:   "https://xiyuer.club",
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			{
				ID:        "2",
				Title:     "今日处暑",
				ImageUrl:  "https://xiyuer.club/minio/coderhub/17470312971996",
				LinkUrl:   "https://xiyuer.club",
				CreatedAt: 0,
				UpdatedAt: 0,
			},
			{
				ID:        "3",
				Title:     "努力的你比谁都可爱",
				ImageUrl:  "https://xiyuer.club/minio/coderhub/17470313163037",
				LinkUrl:   "https://xiyuer.club",
				CreatedAt: 0,
				UpdatedAt: 0,
			},
		},
	}
	return respData, nil
}
