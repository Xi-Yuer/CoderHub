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

type DeleteFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteFavoriteLogic 删除收藏夹
func NewDeleteFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavoriteLogic {
	return &DeleteFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFavoriteLogic) DeleteFavorite(req *types.DeleteFavorFoldReq) (resp *types.DeleteFavorFoldResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.FavoriteService.DeleteFavorFold(l.ctx, &coderhub.DeleteFavorFoldRequest{
		Id:     utils.String2Int(req.Id),
		UserId: userID,
	}); err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteFavoriteLogic) successResp() (resp *types.DeleteFavorFoldResp, err error) {
	return &types.DeleteFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteFavoriteLogic) errorResp(err error) (*types.DeleteFavorFoldResp, error) {
	return &types.DeleteFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
