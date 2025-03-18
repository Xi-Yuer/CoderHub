package school_exp_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSchoolExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSchoolExpLogic 创建学校经历
func NewCreateSchoolExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSchoolExpLogic {
	return &CreateSchoolExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSchoolExpLogic) CreateSchoolExp(req *types.CreateSchoolExpReq) (resp *types.CreateSchoolExpResp, err error) {
	_, err = l.svcCtx.SchoolExpService.CreateSchoolExp(l.ctx, &coderhub.CreateSchoolExpRequest{
		Education: req.Education,
		School:    req.School,
		Major:     req.Major,
		WorkExp:   req.WorkExp,
		Content:   req.Content,
		UserId:    utils.String2Int(req.RequestUserId),
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp()
}

func (l *CreateSchoolExpLogic) errorResp(err error) (resp *types.CreateSchoolExpResp, err1 error) {
	return &types.CreateSchoolExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateSchoolExpLogic) successResp() (resp *types.CreateSchoolExpResp, err error) {
	return &types.CreateSchoolExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
