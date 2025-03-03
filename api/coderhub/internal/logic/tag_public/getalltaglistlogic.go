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

type GetAllTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetAllTagListLogic 获取全部分类标签（包含系统和用户自定义的标签）
func NewGetAllTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTagListLogic {
	return &GetAllTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllTagListLogic) GetAllTagList(req *types.GetTagListReq) (resp *types.GetTagListResp, err error) {
	list, err := l.svcCtx.TagService.GetArticleTagList(l.ctx, &coderhub.GetArticleTagListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(list)
}

func (l *GetAllTagListLogic) errorResp(err error) (resp *types.GetTagListResp, err1 error) {
	return &types.GetTagListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *GetAllTagListLogic) successResp(list *coderhub.GetArticleTagListResponse) (resp *types.GetTagListResp, err error) {
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
