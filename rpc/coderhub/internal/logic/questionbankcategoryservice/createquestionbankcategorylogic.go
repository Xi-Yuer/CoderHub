package questionbankcategoryservicelogic

import (
	"coderhub/model"
	"coderhub/shared/utils"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionBankCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateQuestionBankCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionBankCategoryLogic {
	return &CreateQuestionBankCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateQuestionBankCategoryLogic) CreateQuestionBankCategory(in *coderhub.CreateQuestionBankCategoryRequest) (*coderhub.CreateQuestionBankCategoryResponse, error) {
	err := l.svcCtx.QuestionBankCategoryRepository.Create(l.ctx, &model.QuestionBankCategory{
		ID:          utils.GenID(),
		Name:        in.Name,
		Description: in.Description,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateQuestionBankCategoryResponse{
		Success: true,
	}, nil
}
