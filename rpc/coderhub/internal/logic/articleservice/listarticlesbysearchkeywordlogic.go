package articleservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticlesBySearchKeywordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticlesBySearchKeywordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticlesBySearchKeywordLogic {
	return &ListArticlesBySearchKeywordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticlesBySearchKeywordLogic) ListArticlesBySearchKeyword(in *coderhub.ListArticlesRequest) (*coderhub.ListRecommendedArticlesResponse, error) {
	articlesBySearchKeys, err := l.svcCtx.ArticleRepository.GetArticlesBySearchKeys(in.Keyword, in.Type, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	return &coderhub.ListRecommendedArticlesResponse{
		Ids:   articlesBySearchKeys,
		Total: int64(len(articlesBySearchKeys)),
	}, nil
}
