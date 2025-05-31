package sandpackprojectsservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSandpackProjectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSandpackProjectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSandpackProjectsLogic {
	return &CreateSandpackProjectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSandpackProjectsLogic) CreateSandpackProjects(in *coderhub.CreateSandpackProjectsRequest) (*coderhub.CreateSandpackProjectsResponse, error) {
	id, err := l.svcCtx.SandpackProjectsRepository.Create(l.ctx, &model.SandpackProjects{
		Name:        in.Name,
		Description: in.Description,
		Template:    in.Template,
		UserID:      in.UserId,
		ArticleID:   in.ArticleId,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateSandpackProjectsResponse{
		Id: id,
	}, nil
}
