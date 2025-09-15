package favorservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFavorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFavorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFavorLogic {
	return &CreateFavorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateFavor 创建
func (l *CreateFavorLogic) CreateFavor(in *coderhub.CreateFavorRequest) (*coderhub.CreateFavorResponse, error) {
	folder, err := l.svcCtx.UserFavorFolderRepository.GetFolderByID(l.ctx, in.FavorFolderId)
	if err != nil {
		return nil, err
	}
	if folder.UserId != in.UserId {
		return nil, errors.New("非法操作")
	}
	if folder == nil {
		return nil, errors.New("收藏夹不存在")
	}

	err = l.svcCtx.UserFavorEntityRepository.Create(l.ctx, &model.UserFavor{
		UserId:      in.UserId,
		FavorFoldId: in.FavorFolderId,
		EntityId:    in.EntityId,
		EntityType:  in.EntityType,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateFavorResponse{
		Success: true,
	}, nil
}
