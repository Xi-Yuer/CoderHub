package favorfoldservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFavorFoldLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFavorFoldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavorFoldLogic {
	return &DeleteFavorFoldLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteFavorFold 删除
func (l *DeleteFavorFoldLogic) DeleteFavorFold(in *coderhub.DeleteFavorFoldRequest) (*coderhub.DeleteFavorFoldResponse, error) {
	err := l.svcCtx.UserFavorFolderRepository.Delete(l.ctx, &model.UserFavorFolder{
		Model: gorm.Model{ID: uint(in.Id)},
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteFavorFoldResponse{
		Success: true,
	}, nil
}
