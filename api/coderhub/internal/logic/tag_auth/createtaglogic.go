package tag_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateTagLogic 创建分类标签
func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagReq) (resp *types.CreateTagResp, err error) {
	_, err = l.svcCtx.TagService.CreateArticleTag(l.ctx, &coderhub.CreateArticleTagRequest{
		Name:             req.Name,
		Description:      req.Description,
		IsSystemProvider: req.IsSystemProvider,
		Type:             req.Type,
		Icon:             req.Icon,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *CreateTagLogic) errorResp(err error) (resp *types.CreateTagResp, err1 error) {
	return &types.CreateTagResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateTagLogic) successResp() (resp *types.CreateTagResp, err error) {
	return &types.CreateTagResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
