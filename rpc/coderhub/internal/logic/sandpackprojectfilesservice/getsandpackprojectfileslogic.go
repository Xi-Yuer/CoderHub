package sandpackprojectfilesservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSandpackProjectFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSandpackProjectFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSandpackProjectFilesLogic {
	return &GetSandpackProjectFilesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSandpackProjectFilesLogic) GetSandpackProjectFiles(in *coderhub.GetSandpackProjectFilesRequest) (*coderhub.SandpackProjectFilesResponse, error) {
	files, err := l.svcCtx.SandpackFilesRepository.Get(l.ctx, in.SandpackProjectsId)
	if err != nil {
		return nil, err
	}
	sandpackFiles := make([]*coderhub.SandpackProjectFiles, 0, len(files))
	for _, file := range files {
		sandpackFiles = append(sandpackFiles, &coderhub.SandpackProjectFiles{
			SandpackProjectsId: file.SandpackID,
			FileName:           file.FileName,
			Language:           file.Language,
			Content:            file.Content,
		})
	}

	return &coderhub.SandpackProjectFilesResponse{
		SandpackProjectFiles: sandpackFiles,
	}, nil
}
