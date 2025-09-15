package favorfoldservicelogic

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFavorFoldLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFavorFoldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFavorFoldLogic {
	return &CreateFavorFoldLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateFavorFold 创建
func (l *CreateFavorFoldLogic) CreateFavorFold(in *coderhub.CreateFavorFoldRequest) (*coderhub.CreateFavorFoldResponse, error) {
	err := l.svcCtx.UserFavorFolderRepository.Create(l.ctx, &model.UserFavorFolder{
		UserId:      in.UserId,
		FavorName:   in.Name,
		FavorNum:    0,
		IsPublic:    in.IsPublic,
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}
	// 增加用户经验
	go func() {
		_ = l.svcCtx.UserRepository.IncrUserLevel(in.UserId, conf.FavorFoldCreateLevel)
	}()
	return &coderhub.CreateFavorFoldResponse{
		Success: true,
	}, nil
}
