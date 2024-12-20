package academicnavigatorservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeAcademicNavigatorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeAcademicNavigatorLogic {
	return &LikeAcademicNavigatorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LikeAcademicNavigator 点赞学术导航
func (l *LikeAcademicNavigatorLogic) LikeAcademicNavigator(in *coderhub.LikeAcademicNavigatorRequest) (*coderhub.Response, error) {
	err := l.svcCtx.AcademicRelationLikeRepository.AddAcademicRelationLike(l.ctx, &model.AcademicRelationLike{
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
