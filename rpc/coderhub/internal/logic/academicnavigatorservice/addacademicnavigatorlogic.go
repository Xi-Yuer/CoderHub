package academicnavigatorservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAcademicNavigatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAcademicNavigatorLogic {
	return &AddAcademicNavigatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddAcademicNavigator 新增学术导航
func (l *AddAcademicNavigatorLogic) AddAcademicNavigator(in *coderhub.AddAcademicNavigatorRequest) (*coderhub.Response, error) {
	err := l.svcCtx.AcademicNavigatorRepository.AddAcademicNavigator(&model.AcademicNavigator{
		UserId:    in.UserId,
		Content:   in.Content,
		Education: in.Education,
		Major:     in.Major,
		School:    in.School,
		WorkExp:   in.WorkExp,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.Response{
		Success: true,
	}, nil
}
