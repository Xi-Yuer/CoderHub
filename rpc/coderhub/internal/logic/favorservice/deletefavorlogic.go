package favorservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFavorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFavorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavorLogic {
	return &DeleteFavorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteFavor 删除
func (l *DeleteFavorLogic) DeleteFavor(in *coderhub.DeleteFavorRequest) (*coderhub.DeleteFavorResponse, error) {
	err := l.svcCtx.UserFavorEntityRepository.Delete(l.ctx, &model.UserFavor{
		UserId:     in.UserId,
		EntityId:   in.EntityId,
		EntityType: in.EntityType,
	})
	if err != nil {
		return nil, err
	}
	return &coderhub.DeleteFavorResponse{
		Success: true,
	}, nil
}
