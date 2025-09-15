package articleservicelogic

import (
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRecommendedArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRecommendedArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRecommendedArticlesLogic {
	return &ListRecommendedArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListRecommendedArticlesLogic) ListRecommendedArticles(in *coderhub.ListRecommendedArticlesRequest) (*coderhub.ListRecommendedArticlesResponse, error) {
	articles, err := l.svcCtx.ArticleRepository.ListRecommendedArticles(in.Type, in.CategoryId, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}
	return &coderhub.ListRecommendedArticlesResponse{
		Ids: articles,
	}, nil
}
