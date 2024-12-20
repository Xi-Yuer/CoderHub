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

type DeleteQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteQuestionLogic 删除题目
func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(req *types.DeleteQuestionReq) (resp *types.DeleteQuestionResp, err error) {
	UserId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}
	_, err = l.svcCtx.QuestionBankService.DeleteQuestion(l.ctx, &coderhub.DeleteQuestionRequest{
		Id:         req.Id,
		CreateUser: UserId,
	})

	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteQuestionLogic) errorResp(err error) (resp *types.DeleteQuestionResp, err1 error) {
	return &types.DeleteQuestionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteQuestionLogic) successResp() (resp *types.DeleteQuestionResp, err1 error) {
	return &types.DeleteQuestionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
