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

type CreateWorkExpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateWorkExpLogic 创建工作经历
func NewCreateWorkExpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkExpLogic {
	return &CreateWorkExpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkExpLogic) CreateWorkExp(req *types.CreateWorkExpReq) (resp *types.CreateWorkExpResp, err error) {
	_, err = l.svcCtx.WorkExpService.CreateWorkExp(l.ctx, &coderhub.CreateWorkExpRequest{
		Company: req.Company,
		WorkExp: req.WorkExp,
		Region:  req.Region,
		Content: req.Content,
		UserId:  utils.String2Int(req.RequestUserId),
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp()
}
func (l *CreateWorkExpLogic) errorResp(err error) (resp *types.CreateWorkExpResp, err1 error) {
	return &types.CreateWorkExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateWorkExpLogic) successResp() (resp *types.CreateWorkExpResp, err error) {
	return &types.CreateWorkExpResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
