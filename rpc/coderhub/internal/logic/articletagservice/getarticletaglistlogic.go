package articletagservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleTagListLogic {
	return &GetArticleTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetArticleTagListLogic) GetArticleTagList(in *coderhub.GetArticleTagListRequest) (*coderhub.GetArticleTagListResponse, error) {
	list, total, err := l.svcCtx.ArticleTagRepository.GetList(l.ctx, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}
	var articleTags []*coderhub.ArticleTag
	for _, v := range list {
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

	return &coderhub.GetArticleTagListResponse{
		ArticleTags: articleTags,
		Total:       total,
	}, nil
}
