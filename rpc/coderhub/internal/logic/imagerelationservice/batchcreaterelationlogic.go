package imagerelationservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateRelationLogic {
	return &BatchCreateRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchCreateRelation 批量创建图片关系
func (l *BatchCreateRelationLogic) BatchCreateRelation(in *coderhub.BatchCreateRelationRequest) (*coderhub.BatchCreateRelationResponse, error) {
	// 检查输入切片是否为空
	if len(in.Relations) == 0 {
		return &coderhub.BatchCreateRelationResponse{
			Relations: []*coderhub.ImageRelation{},
		}, nil
	}

	imageRelations := make([]*model.ImageRelation, len(in.Relations))
	for i, relation := range in.Relations {
		imageRelations[i] = &model.ImageRelation{
			ImageID:    relation.ImageId,
			EntityID:   relation.EntityId,
			EntityType: relation.EntityType,
		}
	}
	err := l.svcCtx.ImageRelationRepository.BatchCreate(l.ctx, imageRelations)
	if err != nil {
		return nil, err
	}
	relations := make([]*coderhub.ImageRelation, len(imageRelations))
	for i, relation := range imageRelations {
		relations[i] = &coderhub.ImageRelation{
			ImageId:    relation.ImageID,
			EntityId:   relation.EntityID,
			EntityType: relation.EntityType,
		}
	}
	return &coderhub.BatchCreateRelationResponse{
		Relations: relations,
	}, nil
}
