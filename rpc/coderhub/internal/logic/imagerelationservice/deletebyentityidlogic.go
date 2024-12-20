package imagerelationservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteByEntityIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteByEntityIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteByEntityIDLogic {
	return &DeleteByEntityIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteByEntityID 根据实体ID、实体类型删除图片关系
func (l *DeleteByEntityIDLogic) DeleteByEntityID(in *coderhub.DeleteByEntityIDRequest) (*coderhub.DeleteByEntityIDResponse, error) {
	// todo: add your logic here and delete this line

	return &coderhub.DeleteByEntityIDResponse{}, nil
}
