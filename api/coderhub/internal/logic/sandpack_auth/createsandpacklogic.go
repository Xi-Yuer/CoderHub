package sandpack_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSandpackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSandpackLogic 创建一个sandpack
func NewCreateSandpackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSandpackLogic {
	return &CreateSandpackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSandpackLogic) CreateSandpack(req *types.CreateSandpackProject) (resp *types.CreateSandpackProjectResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	sandpackProjects, err := l.svcCtx.SandpackProjectService.CreateSandpackProjects(l.ctx, &coderhub.CreateSandpackProjectsRequest{
		Name:        req.Name,
		Description: req.Description,
		Template:    req.Template,
		UserId:      userID,
		ArticleId:   utils.String2Int(req.ArticleID),
	})
	if err != nil {
		return l.errorResp(err)
	}

	var sandpackProjectFiles []*coderhub.SandpackProjectFiles
	for _, file := range req.Files {
		sandpackProjectFiles = append(sandpackProjectFiles, &coderhub.SandpackProjectFiles{
			SandpackProjectsId: sandpackProjects.Id,
			FileName:           file.Name,
			Language:           file.Language,
			Content:            file.Code,
		})
	}

	if _, err = l.svcCtx.SandpackFileService.CreateSandpackProjectFiles(l.ctx, &coderhub.SandpackProjectFilesResponse{
		SandpackProjectFiles: sandpackProjectFiles,
	}); err != nil {
		return l.errorResp(err)
	}
	return l.successResp(nil)
}

func (l *CreateSandpackLogic) errorResp(err error) (resp *types.CreateSandpackProjectResp, err1 error) {
	return &types.CreateSandpackProjectResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *CreateSandpackLogic) successResp(data *types.CreateSandpackProject) (resp *types.CreateSandpackProjectResp, err error) {
	return &types.CreateSandpackProjectResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
