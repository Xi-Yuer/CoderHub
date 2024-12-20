package academic_auth

import (
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelLikeAcademicNavigatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCancelLikeAcademicNavigatorLogic 取消点赞学术导航
func NewCancelLikeAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelLikeAcademicNavigatorLogic {
	return &CancelLikeAcademicNavigatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelLikeAcademicNavigatorLogic) CancelLikeAcademicNavigator(req *types.CancelLikeAcademicNavigatorReq) (resp *types.CancelLikeAcademicNavigatorResp, err error) {
	// 获取用户ID
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.AcademicService.CancelLikeAcademicNavigator(l.ctx, &coderhub.CancelLikeAcademicNavigatorRequest{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *CancelLikeAcademicNavigatorLogic) successResp() (*types.CancelLikeAcademicNavigatorResp, error) {
	return &types.CancelLikeAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *CancelLikeAcademicNavigatorLogic) errorResp(err error) (*types.CancelLikeAcademicNavigatorResp, error) {
	return &types.CancelLikeAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
