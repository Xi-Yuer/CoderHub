package favorites_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateFavoriteLogic 更新收藏夹
func NewUpdateFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteLogic {
	return &UpdateFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFavoriteLogic) UpdateFavorite(req *types.UpdateFavorFoldReq) (resp *types.UpdateFavorFoldResp, err error) {
	if _, err := l.svcCtx.FavoriteService.UpdateFavorFold(l.ctx, &coderhub.UpdateFavorFoldRequest{
		Id:          utils.String2Int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
	}); err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *UpdateFavoriteLogic) successResp() (*types.UpdateFavorFoldResp, error) {
	return &types.UpdateFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *UpdateFavoriteLogic) errorResp(err error) (*types.UpdateFavorFoldResp, error) {
	return &types.UpdateFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
