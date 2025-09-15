package favorites_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListFavoriteLogic 获取收藏夹列表
func NewListFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFavoriteLogic {
	return &ListFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFavoriteLogic) ListFavorite(req *types.GetFavorFoldListReq) (resp *types.GetFavorFoldListResp, err error) {
	list, err := l.svcCtx.FavoriteService.GetFavorFoldList(l.ctx, &coderhub.GetFavorFoldListRequest{
		Page:          req.Page,
		PageSize:      req.PageSize,
		UserId:        utils.String2Int(req.UserId),
		RequestUserId: utils.String2Int(req.RequestUserId),
	})
	if err != nil {
		return l.errorResp(err)
	}

	response := make([]*types.FavorFold, 0, len(list.FavorFolds))

	for _, v := range list.FavorFolds {
		response = append(response, &types.FavorFold{
			ID:          utils.Int2String(v.Id),
			Name:        v.Name,
			Description: v.Description,
			IsPublic:    v.IsPublic,
			CreateUser:  utils.Int2String(v.UserId),
			CreatedAt:   v.CreateTime,
			UpdatedAt:   v.UpdateTime,
		})
	}

	return l.successResp(types.FavorFoldList{
		List:  response,
		Total: list.Total,
	})
}

func (l *ListFavoriteLogic) successResp(list types.FavorFoldList) (*types.GetFavorFoldListResp, error) {
	return &types.GetFavorFoldListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: list,
	}, nil
}

func (l *ListFavoriteLogic) errorResp(err error) (*types.GetFavorFoldListResp, error) {
	return &types.GetFavorFoldListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.FavorFoldList{},
	}, nil
}
