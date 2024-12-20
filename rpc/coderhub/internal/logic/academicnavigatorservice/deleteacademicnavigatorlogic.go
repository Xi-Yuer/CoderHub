package academicnavigatorservicelogic

import (
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAcademicNavigatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAcademicNavigatorLogic {
	return &DeleteAcademicNavigatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteAcademicNavigator 删除学术导航
func (l *DeleteAcademicNavigatorLogic) DeleteAcademicNavigator(in *coderhub.DeleteAcademicNavigatorRequest) (*coderhub.Response, error) {
	academicNavigator, err := l.svcCtx.AcademicNavigatorRepository.GetAcademicNavigatorByID(in.Id)
	if err != nil {
		return nil, err
	}
	if academicNavigator == nil {
		return nil, errors.New("资源不存在")
	}

	if academicNavigator.UserId != in.UserId {
		return nil, errors.New("非法操作")
	}

	err = l.svcCtx.AcademicNavigatorRepository.DeleteAcademicNavigator(int64(academicNavigator.ID))
	if err != nil {
		return nil, err
	}

	return &coderhub.Response{
		Success: true,
	}, nil
}
