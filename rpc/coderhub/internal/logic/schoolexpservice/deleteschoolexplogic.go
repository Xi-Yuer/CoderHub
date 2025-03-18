package schoolexpservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSchoolExpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSchoolExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSchoolExpLogic {
	return &DeleteSchoolExpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSchoolExpLogic) DeleteSchoolExp(in *coderhub.DeleteSchoolExpRequest) (*coderhub.DeleteSchoolExpResponse, error) {
	err := l.svcCtx.SchoolExpRepository.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &coderhub.DeleteSchoolExpResponse{
		Success: true,
	}, nil
}
