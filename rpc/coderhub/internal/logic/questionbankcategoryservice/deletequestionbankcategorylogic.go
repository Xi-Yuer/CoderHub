package questionbankcategoryservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionBankCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionBankCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionBankCategoryLogic {
	return &DeleteQuestionBankCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteQuestionBankCategoryLogic) DeleteQuestionBankCategory(in *coderhub.DeleteQuestionBankCategoryRequest) (*coderhub.DeleteQuestionBankCategoryResponse, error) {
	err := l.svcCtx.QuestionBankRepository.DeleteQuestionBank(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteQuestionBankCategoryResponse{
		Success: true,
	}, nil
}
