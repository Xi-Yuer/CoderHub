package articletagservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArticleTagUsageCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleTagUsageCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleTagUsageCountLogic {
	return &UpdateArticleTagUsageCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateArticleTagUsageCountLogic) UpdateArticleTagUsageCount(in *coderhub.UpdateArticleTagUsageCountRequest) (*coderhub.UpdateArticleTagUsageCountResponse, error) {
	err := l.svcCtx.ArticleTagRepository.UpdateTagUsageCount(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &coderhub.UpdateArticleTagUsageCountResponse{
		Success: true,
	}, nil
}
