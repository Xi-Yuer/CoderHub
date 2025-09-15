package questionservicelogic

import (
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionBankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionBankLogic {
	return &DeleteQuestionBankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteQuestionBank 删除题库
func (l *DeleteQuestionBankLogic) DeleteQuestionBank(in *coderhub.DeleteQuestionBankRequest) (*coderhub.DeleteQuestionBankResponse, error) {
	question, err := l.svcCtx.QuestionBankRepository.GetQuestionBankByID(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, errors.New("题库不存在")
	}

	if question.CreateUser != in.CreateUser {
		return nil, errors.New("非法操作")
	}

	if err := l.svcCtx.QuestionBankRepository.DeleteQuestionBank(l.ctx, in.Id); err != nil {
		return nil, err
	}

	return &coderhub.DeleteQuestionBankResponse{
		Success: true,
	}, nil
}
