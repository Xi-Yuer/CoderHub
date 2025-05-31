package sandpackprojectsservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSandpackProjectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSandpackProjectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSandpackProjectsLogic {
	return &GetSandpackProjectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSandpackProjectsLogic) GetSandpackProjects(in *coderhub.GetSandpackProjectsRequest) (*coderhub.SandpackProjects, error) {
	SandpackProject, err := l.svcCtx.SandpackProjectsRepository.Get(l.ctx, in.UserId, in.ArticleId)
	if err != nil {
		return nil, err
	}

	return &coderhub.SandpackProjects{
		Id:          int64(SandpackProject.ID),
		Name:        SandpackProject.Name,
		Template:    SandpackProject.Template,
		UserId:      SandpackProject.UserID,
		ArticleId:   SandpackProject.ArticleID,
		Description: SandpackProject.Description,
		CreatedAt:   SandpackProject.CreatedAt.Unix(),
		UpdatedAt:   SandpackProject.UpdatedAt.Unix(),
	}, nil
}
