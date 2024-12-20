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

type DeleteAcademicNavigatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteAcademicNavigatorLogic 删除学术导航
func NewDeleteAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAcademicNavigatorLogic {
	return &DeleteAcademicNavigatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAcademicNavigatorLogic) DeleteAcademicNavigator(req *types.DeleteAcademicNavigatorReq) (resp *types.DeleteAcademicNavigatorResp, err error) {
	// 获取用户ID
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.AcademicService.DeleteAcademicNavigator(l.ctx, &coderhub.DeleteAcademicNavigatorRequest{
		Id:     req.Id,
		UserId: userId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteAcademicNavigatorLogic) successResp() (*types.DeleteAcademicNavigatorResp, error) {
	return &types.DeleteAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteAcademicNavigatorLogic) errorResp(err error) (*types.DeleteAcademicNavigatorResp, error) {
	return &types.DeleteAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
	}, nil
}
