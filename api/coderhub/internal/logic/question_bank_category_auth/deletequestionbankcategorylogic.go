package question_bank_category_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionBankCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteQuestionBankCategoryLogic 删除题库分类
func NewDeleteQuestionBankCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionBankCategoryLogic {
	return &DeleteQuestionBankCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionBankCategoryLogic) DeleteQuestionBankCategory(req *types.DeleteQuestionBankCategoryReq) (resp *types.DeleteQuestionBankCategoryResp, err error) {
	_, err = l.svcCtx.QuestionBankCategoryService.DeleteQuestionBankCategory(l.ctx, &coderhub.DeleteQuestionBankCategoryRequest{
		Id: utils.String2Int(req.Id),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteQuestionBankCategoryLogic) errorResp(err error) (resp *types.DeleteQuestionBankCategoryResp, err1 error) {
	return &types.DeleteQuestionBankCategoryResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteQuestionBankCategoryLogic) successResp() (resp *types.DeleteQuestionBankCategoryResp, err1 error) {
	return &types.DeleteQuestionBankCategoryResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: "删除题库分类成功",
		},
		Data: true,
	}, nil
}
