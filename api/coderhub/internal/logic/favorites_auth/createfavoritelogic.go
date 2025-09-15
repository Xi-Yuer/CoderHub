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

type CreateFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateFavoriteLogic 创建收藏夹
func NewCreateFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFavoriteLogic {
	return &CreateFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateFavoriteLogic) CreateFavorite(req *types.CreateFavorFoldReq) (resp *types.CreateFavorFoldResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	fold, err := l.svcCtx.FavoriteService.CreateFavorFold(l.ctx, &coderhub.CreateFavorFoldRequest{
		Name:        req.Name,
		Description: req.Description,
		UserId:      userID,
		IsPublic:    req.IsPublic,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(fold)
}

func (l *CreateFavoriteLogic) successResp(response *coderhub.CreateFavorFoldResponse) (resp *types.CreateFavorFoldResp, err error) {
	return &types.CreateFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: response.Success,
	}, nil
}

func (l *CreateFavoriteLogic) errorResp(err error) (*types.CreateFavorFoldResp, error) {
	return &types.CreateFavorFoldResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
