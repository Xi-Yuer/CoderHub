package workexpservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkExpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWorkExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkExpLogic {
	return &DeleteWorkExpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteWorkExpLogic) DeleteWorkExp(in *coderhub.DeleteWorkExpRequest) (*coderhub.DeleteWorkExpResponse, error) {
	err := l.svcCtx.WorkExpRepository.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteWorkExpResponse{
		Success: true,
	}, nil
}
