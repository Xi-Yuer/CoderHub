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

type AddAcademicNavigatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAddAcademicNavigatorLogic 新增学术导航
func NewAddAcademicNavigatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAcademicNavigatorLogic {
	return &AddAcademicNavigatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAcademicNavigatorLogic) AddAcademicNavigator(req *types.AddAcademicNavigatorReq) (resp *types.AddAcademicNavigatorResp, err error) {
	// 获取用户ID
	userId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.AcademicService.AddAcademicNavigator(l.ctx, &coderhub.AddAcademicNavigatorRequest{
		UserId:  userId,
		Content: req.Content,
		Major:   req.Major,
		School:  req.School,
		WorkExp: req.WorkExp,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *AddAcademicNavigatorLogic) successResp() (*types.AddAcademicNavigatorResp, error) {
	return &types.AddAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *AddAcademicNavigatorLogic) errorResp(err error) (*types.AddAcademicNavigatorResp, error) {
	return &types.AddAcademicNavigatorResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}
