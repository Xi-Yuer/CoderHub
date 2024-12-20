package questions_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteQuestionBankLogic 删除题库
func NewDeleteQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionBankLogic {
	return &DeleteQuestionBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionBankLogic) DeleteQuestionBank(req *types.DeleteQuestionBankReq) (resp *types.DeleteQuestionBankResp, err error) {
	UserId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.QuestionBankService.DeleteQuestionBank(l.ctx, &coderhub.DeleteQuestionBankRequest{
		Id:         req.Id,
		CreateUser: UserId,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteQuestionBankLogic) errorResp(err error) (resp *types.DeleteQuestionBankResp, err1 error) {
	return &types.DeleteQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteQuestionBankLogic) successResp() (resp *types.DeleteQuestionBankResp, err1 error) {
	return &types.DeleteQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
