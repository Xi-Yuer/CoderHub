package work_exp_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWorkExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteWorkExpLogic 删除工作经历
func NewDeleteWorkExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkExpLogic {
	return &DeleteWorkExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWorkExpLogic) DeleteWorkExp(req *types.DeleteWorkExpReq) (resp *types.DeleteWorkExpResp, err error) {
	_, err = l.svcCtx.WorkExpService.DeleteWorkExp(l.ctx, &coderhub.DeleteWorkExpRequest{
		Id: utils.String2Int(req.Id),
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp()
}
func (l *DeleteWorkExpLogic) errorResp(err error) (resp *types.DeleteWorkExpResp, err1 error) {
	return &types.DeleteWorkExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteWorkExpLogic) successResp() (resp *types.DeleteWorkExpResp, err error) {
	return &types.DeleteWorkExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
