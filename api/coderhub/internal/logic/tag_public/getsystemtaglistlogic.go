package tag_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSystemTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetSystemTagListLogic 获取系统分类标签
func NewGetSystemTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSystemTagListLogic {
	return &GetSystemTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSystemTagListLogic) GetSystemTagList(req *types.GetSystemTagReq) (resp *types.GetTagListResp, err error) {
	list, err := l.svcCtx.TagService.GetSystemProviderTagList(l.ctx, &coderhub.GetSystemProviderTagListRequest{
		Type: req.Type,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(list)
}

func (l *GetSystemTagListLogic) errorResp(err error) (resp *types.GetTagListResp, err1 error) {
	return &types.GetTagListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *GetSystemTagListLogic) successResp(list *coderhub.GetSystemProviderTagListResponse) (resp *types.GetTagListResp, err error) {
	var tagList *types.TagList
	if list != nil {
		tagList = &types.TagList{
			Total: 0,
			List:  nil,
		}
	}
	if list != nil {
		tagList.Total = list.Total
		for _, v := range list.ArticleTags {
			tagList.List = append(tagList.List, &types.Tag{
				ID:               utils.Int2String(v.Id),
				Name:             v.Name,
				Description:      v.Description,
				Type:             v.Type,
				IsSystemProvider: v.IsSystemProvider,
				Icon:             v.Icon,
				UsageCount:       v.UsageCount,
				CreatedAt:        v.CreatedAt,
				UpdatedAt:        v.UpdatedAt,
			})
		}
	}
	return &types.GetTagListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: tagList,
	}, nil
}
