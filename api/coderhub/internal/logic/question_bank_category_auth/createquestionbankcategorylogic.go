package question_bank_category_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionBankCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateQuestionBankCategoryLogic 创建题库分类
func NewCreateQuestionBankCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionBankCategoryLogic {
	return &CreateQuestionBankCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionBankCategoryLogic) CreateQuestionBankCategory(req *types.CreateQuestionBankCategoryReq) (resp *types.CreateQuestionBankCategoryResp, err error) {
	_, err = l.svcCtx.QuestionBankCategoryService.CreateQuestionBankCategory(l.ctx, &coderhub.CreateQuestionBankCategoryRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp()
}

func (l *CreateQuestionBankCategoryLogic) errorResp(err error) (resp *types.CreateQuestionBankCategoryResp, err1 error) {
	return &types.CreateQuestionBankCategoryResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateQuestionBankCategoryLogic) successResp() (resp *types.CreateQuestionBankCategoryResp, err1 error) {
	return &types.CreateQuestionBankCategoryResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "创建题库分类成功",
		},
		Data: true,
	}, nil
}
