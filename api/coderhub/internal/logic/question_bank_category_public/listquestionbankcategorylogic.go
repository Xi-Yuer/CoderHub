package question_bank_category_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuestionBankCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListQuestionBankCategoryLogic 获取题库分类列表
func NewListQuestionBankCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionBankCategoryLogic {
	return &ListQuestionBankCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuestionBankCategoryLogic) ListQuestionBankCategory(req *types.GetQuestionBankCategoryListReq) (resp *types.GetQuestionBankCategoryListResp, err error) {
	list, err := l.svcCtx.QuestionBankCategoryService.GetQuestionBankCategoryList(l.ctx, &coderhub.GetQuestionBankCategoryListRequest{})
	if err != nil {
		return l.errorResp(err)
	}
	var questionBankCategories []*types.QuestionBankCategory
	for _, v := range list.QuestionBankCategories {
		questionBankCategories = append(questionBankCategories, &types.QuestionBankCategory{
			ID:          utils.Int2String(v.Id),
			Name:        v.Name,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return l.successResp(questionBankCategories)
}

func (l *ListQuestionBankCategoryLogic) errorResp(err error) (resp *types.GetQuestionBankCategoryListResp, err1 error) {
	return &types.GetQuestionBankCategoryListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: &types.QuestionBankCategoryList{
			List: nil,
		},
	}, nil
}

func (l *ListQuestionBankCategoryLogic) successResp(list []*types.QuestionBankCategory) (resp *types.GetQuestionBankCategoryListResp, err1 error) {
	return &types.GetQuestionBankCategoryListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "获取题库分类成功",
		},
		Data: &types.QuestionBankCategoryList{List: list},
	}, nil
}
