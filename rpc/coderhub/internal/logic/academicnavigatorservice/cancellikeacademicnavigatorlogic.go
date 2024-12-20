package academicnavigatorservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLikeAcademicNavigatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelLikeAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLikeAcademicNavigatorLogic {
	return &CancelLikeAcademicNavigatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CancelLikeAcademicNavigator 取消点赞学术导航
func (l *CancelLikeAcademicNavigatorLogic) CancelLikeAcademicNavigator(in *coderhub.CancelLikeAcademicNavigatorRequest) (*coderhub.Response, error) {
	err := l.svcCtx.AcademicRelationLikeRepository.DeleteAcademicRelationLike(l.ctx, &model.AcademicRelationLike{
		AcademicNavigatorID: in.Id,
		UserID:              in.UserId,
	})

	if err != nil {
		return nil, err
	}

	return &coderhub.Response{
		Success: true,
	}, nil
}
