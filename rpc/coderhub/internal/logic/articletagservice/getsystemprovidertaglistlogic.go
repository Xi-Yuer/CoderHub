package articletagservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSystemProviderTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSystemProviderTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSystemProviderTagListLogic {
	return &GetSystemProviderTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSystemProviderTagListLogic) GetSystemProviderTagList(in *coderhub.GetSystemProviderTagListRequest) (*coderhub.GetSystemProviderTagListResponse, error) {
	tags, total, err := l.svcCtx.ArticleTagRepository.GetSystemTags(l.ctx)
	if err != nil {
		return nil, err
	}
	var articleTags []*coderhub.ArticleTag
	for _, v := range tags {
		articleTags = append(articleTags, &coderhub.ArticleTag{
			Id:               v.ID,
			Name:             v.Name,
			Description:      v.Description,
			IsSystemProvider: v.IsSystemProvider,
			Icon:             v.Icon,
			UsageCount:       v.UsageCount,
			CreatedAt:        v.CreatedAt.Unix(),
			UpdatedAt:        v.UpdatedAt.Unix(),
		})
	}

	return &coderhub.GetSystemProviderTagListResponse{
		ArticleTags: articleTags,
		Total:       total,
	}, nil
}
