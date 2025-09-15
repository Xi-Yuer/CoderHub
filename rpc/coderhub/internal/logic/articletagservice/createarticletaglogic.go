package articletagservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"coderhub/shared/utils"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateArticleTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateArticleTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateArticleTagLogic {
	return &CreateArticleTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateArticleTagLogic) CreateArticleTag(in *coderhub.CreateArticleTagRequest) (*coderhub.CreateArticleTagResponse, error) {
	err := l.svcCtx.ArticleTagRepository.Create(l.ctx, &model.ArticleTag{
		ID:               utils.GenID(),
		Name:             in.Name,
		Type:             in.Type,
		Description:      in.Description,
		Icon:             in.Icon,
		IsSystemProvider: in.IsSystemProvider,
		UsageCount:       0,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateArticleTagResponse{
		Success: true,
	}, nil
}
