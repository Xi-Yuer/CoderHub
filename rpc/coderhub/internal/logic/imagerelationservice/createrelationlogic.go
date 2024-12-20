package imagerelationservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRelationLogic {
	return &CreateRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateRelation 创建图片关系
func (l *CreateRelationLogic) CreateRelation(in *coderhub.CreateRelationRequest) (*coderhub.CreateRelationResponse, error) {
	err := l.svcCtx.ImageRelationRepository.Create(l.ctx, &model.ImageRelation{
		ImageID:    in.ImageId,
		EntityID:   in.EntityId,
		EntityType: in.EntityType,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateRelationResponse{
		Relation: &coderhub.ImageRelation{
			ImageId:    in.ImageId,
			EntityId:   in.EntityId,
			EntityType: in.EntityType,
		},
	}, nil
}
