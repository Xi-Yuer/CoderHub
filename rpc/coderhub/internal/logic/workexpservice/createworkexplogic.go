package workexpservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkExpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkExpLogic {
	return &CreateWorkExpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateWorkExpLogic) CreateWorkExp(in *coderhub.CreateWorkExpRequest) (*coderhub.CreateWorkExpResponse, error) {
	err := l.svcCtx.WorkExpRepository.Create(l.ctx, &model.WorkExp{
		Company:      in.Company,
		WorkDuration: in.WorkExp,
		Region:       in.Region,
		Content:      in.Content,
		UserId:       in.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateWorkExpResponse{
		Success: true,
	}, nil
}
