package articletagservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteArticleTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleTagLogic {
	return &DeleteArticleTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteArticleTagLogic) DeleteArticleTag(in *coderhub.DeleteArticleTagRequest) (*coderhub.DeleteArticleTagResponse, error) {
	err := l.svcCtx.ArticleTagRepository.Delete(l.ctx, &model.ArticleTag{ID: in.Id})
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteArticleTagResponse{
		Success: true,
	}, nil
}
