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

type DeleteSchoolExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteSchoolExpLogic 删除学校经历
func NewDeleteSchoolExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSchoolExpLogic {
	return &DeleteSchoolExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSchoolExpLogic) DeleteSchoolExp(req *types.DeleteSchoolExpReq) (resp *types.DeleteSchoolExpResp, err error) {
	_, err = l.svcCtx.SchoolExpService.DeleteSchoolExp(l.ctx, &coderhub.DeleteSchoolExpRequest{
		Id: utils.String2Int(req.Id),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}
func (l *DeleteSchoolExpLogic) errorResp(err error) (resp *types.DeleteSchoolExpResp, err1 error) {
	return &types.DeleteSchoolExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteSchoolExpLogic) successResp() (resp *types.DeleteSchoolExpResp, err error) {
	return &types.DeleteSchoolExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
