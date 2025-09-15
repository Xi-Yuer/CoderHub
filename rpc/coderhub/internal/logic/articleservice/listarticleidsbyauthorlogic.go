package articleservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListArticleIDsByAuthorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListArticleIDsByAuthorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArticleIDsByAuthorLogic {
	return &ListArticleIDsByAuthorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListArticleIDsByAuthorLogic) ListArticleIDsByAuthor(in *coderhub.ListAuthorArticlesRequest) (*coderhub.ListRecommendedArticlesResponse, error) {
	listArticlesByAuthor, total, err := l.svcCtx.ArticleRepository.ListArticlesByAuthor(in.AuthorId, in.Type, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	return &coderhub.ListRecommendedArticlesResponse{
		Ids:   listArticlesByAuthor,
		Total: total,
	}, nil
}
