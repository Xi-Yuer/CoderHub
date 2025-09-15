package imagerelationservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRelationLogic {
	return &BatchDeleteRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchDeleteRelation 批量删除图片关系
func (l *BatchDeleteRelationLogic) BatchDeleteRelation(in *coderhub.BatchDeleteRelationRequest) (*coderhub.BatchDeleteRelationResponse, error) {
	if err := l.svcCtx.ImageRelationRepository.BatchDelete(l.ctx, in.Ids); err != nil {
		return nil, err
	}

	return &coderhub.BatchDeleteRelationResponse{
		Success: true,
	}, nil
}
