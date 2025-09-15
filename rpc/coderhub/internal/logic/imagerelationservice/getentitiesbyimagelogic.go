package imagerelationservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEntitiesByImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEntitiesByImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEntitiesByImageLogic {
	return &GetEntitiesByImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetEntitiesByImage 获取图片关联的实体列表
func (l *GetEntitiesByImageLogic) GetEntitiesByImage(in *coderhub.GetEntitiesByImageRequest) (*coderhub.GetEntitiesByImageResponse, error) {
	imageRelations, err := l.svcCtx.ImageRelationRepository.ListByImageID(l.ctx, in.ImageId)
	if err != nil {
		return nil, err
	}

	entities := make([]*coderhub.ImageRelation, len(imageRelations))
	for i, relation := range imageRelations {
		entities[i] = &coderhub.ImageRelation{
			Id:         relation.ID,
			ImageId:    relation.ImageID,
			EntityId:   relation.EntityID,
			EntityType: relation.EntityType,
		}
	}
	entityInfos := make([]*coderhub.EntityInfo, 0)
	for _, relation := range imageRelations {
		entityInfos = append(entityInfos, &coderhub.EntityInfo{
			EntityType: relation.EntityType,
			EntityId:   relation.EntityID,
			CreatedAt:  relation.CreatedAt.Unix(),
		})
	}
	return &coderhub.GetEntitiesByImageResponse{
		Entities: entityInfos,
	}, nil
}
