package codehub

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLikeAcademicNavigatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消点赞学术导航
func NewCancelLikeAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLikeAcademicNavigatorLogic {
	return &CancelLikeAcademicNavigatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLikeAcademicNavigatorLogic) CancelLikeAcademicNavigator(req *types.CancelLikeAcademicNavigatorReq) (resp *types.CancelLikeAcademicNavigatorResp, err error) {
	// todo: add your logic here and delete this line

	return
}