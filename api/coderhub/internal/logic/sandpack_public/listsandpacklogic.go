package sandpack_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSandpackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListSandpackLogic 获取sandpack列表
func NewListSandpackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSandpackLogic {
	return &ListSandpackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSandpackLogic) ListSandpack(req *types.GetSandpackProjectListReq) (resp *types.GetSandpackProjectListResp, err error) {
	sandpackProjects, err := l.svcCtx.SandpackProjectService.GetSandpackProjects(l.ctx, &coderhub.GetSandpackProjectsRequest{
		UserId:    utils.String2Int(req.UserID),
		ArticleId: utils.String2Int(req.ArticleID),
	})
	if err != nil {
		return l.errorResp(err)
	}
	projectFiles, err := l.svcCtx.SandpackFileService.GetSandpackProjectFiles(l.ctx, &coderhub.GetSandpackProjectFilesRequest{
		SandpackProjectsId: sandpackProjects.Id,
	})
	if err != nil {
		return l.errorResp(err)
	}
	var files []*types.SandpackProjectFile
	for _, file := range projectFiles.SandpackProjectFiles {
		files = append(files, &types.SandpackProjectFile{
			Name:     file.FileName,
			Code:     file.Content,
			Language: file.Language,
		})
	}

	sandpackProject := types.SandpackProject{
		ID:          utils.Int2String(sandpackProjects.Id),
		Name:        sandpackProjects.Name,
		Description: sandpackProjects.Description,
		Template:    sandpackProjects.Template,
		UserID:      utils.Int2String(sandpackProjects.UserId),
		ArticleID:   utils.Int2String(sandpackProjects.ArticleId),
		Files:       files,
		CreatedAt:   sandpackProjects.CreatedAt,
		UpdatedAt:   sandpackProjects.UpdatedAt,
	}

	return l.successResp(&sandpackProject)
}

func (l *ListSandpackLogic) errorResp(err error) (resp *types.GetSandpackProjectListResp, err1 error) {
	return &types.GetSandpackProjectListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListSandpackLogic) successResp(data *types.SandpackProject) (resp *types.GetSandpackProjectListResp, err error) {
	return &types.GetSandpackProjectListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
