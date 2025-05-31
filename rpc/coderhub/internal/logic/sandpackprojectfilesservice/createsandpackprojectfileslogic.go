package sandpackprojectfilesservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSandpackProjectFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSandpackProjectFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSandpackProjectFilesLogic {
	return &CreateSandpackProjectFilesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSandpackProjectFilesLogic) CreateSandpackProjectFiles(in *coderhub.SandpackProjectFilesResponse) (*coderhub.CreateSandpackProjectFilesResponse, error) {
	var sandpackFiles []*model.SandpackFiles
	for _, sandpackFile := range in.SandpackProjectFiles {
		fmt.Println("sandpackFile.SandpackProjectsId", sandpackFile.SandpackProjectsId)
		sandpackFiles = append(sandpackFiles, &model.SandpackFiles{
			SandpackID: sandpackFile.SandpackProjectsId,
			FileName:   sandpackFile.FileName,
			Language:   sandpackFile.Language,
			Content:    sandpackFile.Content,
		})
	}

	if err := l.svcCtx.SandpackFilesRepository.Create(l.ctx, sandpackFiles); err != nil {
		return nil, err
	}

	return &coderhub.CreateSandpackProjectFilesResponse{
		Success: true,
	}, nil
}
