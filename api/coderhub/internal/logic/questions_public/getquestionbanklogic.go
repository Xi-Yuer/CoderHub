package questions_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionBankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetQuestionBankLogic 获取题目详情
func NewGetQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionBankLogic {
	return &GetQuestionBankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionBankLogic) GetQuestionBank(req *types.GetQuestionBankReq) (resp *types.GetQuestionBankResp, err error) {
	question, err := l.svcCtx.QuestionBankService.GetQuestion(l.ctx, &coderhub.GetQuestionRequest{
		QuestionId: req.Id,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(&types.Question{
		Id:        question.Id,
		Title:     question.Title,
		BankId:    question.BankId,
		Content:   question.Content,
		Difficult: question.Difficulty,
		CreatedAt: question.CreateTime,
		UpdatedAt: question.UpdateTime,
	})
}

func (l *GetQuestionBankLogic) errorResp(err error) (resp *types.GetQuestionBankResp, err1 error) {
	return &types.GetQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *GetQuestionBankLogic) successResp(data *types.Question) (resp *types.GetQuestionBankResp, err1 error) {
	return &types.GetQuestionBankResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
