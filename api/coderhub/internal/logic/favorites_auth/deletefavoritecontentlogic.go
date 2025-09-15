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

type DeleteFavoriteContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteFavoriteContentLogic 删除收藏夹内容
func NewDeleteFavoriteContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavoriteContentLogic {
	return &DeleteFavoriteContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFavoriteContentLogic) DeleteFavoriteContent(req *types.DeleteFavorReq) (resp *types.DeleteFavorResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	if _, err := l.svcCtx.FavoriteContentService.DeleteFavor(l.ctx, &coderhub.DeleteFavorRequest{
		UserId:     userID,
		EntityId:   utils.String2Int(req.EntityId),
		EntityType: req.EntityType,
	}); err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteFavoriteContentLogic) successResp() (*types.DeleteFavorResp, error) {
	return &types.DeleteFavorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteFavoriteContentLogic) errorResp(err error) (*types.DeleteFavorResp, error) {
	return &types.DeleteFavorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
