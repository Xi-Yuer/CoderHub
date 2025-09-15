package favorfoldservicelogic

import (
	"coderhub/model"
	"context"
	"gorm.io/gorm"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavorFoldLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFavorFoldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavorFoldLogic {
	return &UpdateFavorFoldLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateFavorFold 更新
func (l *UpdateFavorFoldLogic) UpdateFavorFold(in *coderhub.UpdateFavorFoldRequest) (*coderhub.UpdateFavorFoldResponse, error) {
	err := l.svcCtx.UserFavorFolderRepository.Update(l.ctx, &model.UserFavorFolder{
		Model: gorm.Model{
			ID: uint(in.Id),
		},
		FavorName:   in.Name,
		IsPublic:    in.IsPublic,
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.UpdateFavorFoldResponse{
		Success: true,
	}, nil
}
