package questionbankcategoryservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionBankCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionBankCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionBankCategoryListLogic {
	return &GetQuestionBankCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionBankCategoryListLogic) GetQuestionBankCategoryList(in *coderhub.GetQuestionBankCategoryListRequest) (*coderhub.GetQuestionBankCategoryListResponse, error) {
	list, err := l.svcCtx.QuestionBankCategoryRepository.List(l.ctx)
	if err != nil {
		return nil, err
	}
	var questionBankCategories []*coderhub.QuestionBankCategory
	for _, v := range list {
		questionBankCategories = append(questionBankCategories, &coderhub.QuestionBankCategory{
			Id:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Unix(),
			UpdatedAt:   v.UpdatedAt.Unix(),
		})
	}
	return &coderhub.GetQuestionBankCategoryListResponse{
		QuestionBankCategories: questionBankCategories,
	}, nil
}
