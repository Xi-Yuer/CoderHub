package schoolexpservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSchoolExpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSchoolExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSchoolExpLogic {
	return &CreateSchoolExpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSchoolExpLogic) CreateSchoolExp(in *coderhub.CreateSchoolExpRequest) (*coderhub.CreateSchoolExpResponse, error) {
	err := l.svcCtx.SchoolExpRepository.Create(l.ctx, &model.SchoolExp{
		Education: in.Education,
		Major:     in.Major,
		School:    in.School,
		WorkExp:   in.WorkExp,
		Content:   in.Content,
		UserId:    in.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateSchoolExpResponse{
		Success: true,
	}, nil
}
