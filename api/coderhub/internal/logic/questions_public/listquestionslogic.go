package questions_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListQuestionsLogic 获取题目列表
func NewListQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionsLogic {
	return &ListQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuestionsLogic) ListQuestions(req *types.GetQuestionListReq) (resp *types.GetQuestionListResp, err error) {
	tree, err := l.svcCtx.QuestionBankService.GetQuestionTree(l.ctx, &coderhub.GetQuestionTreeRequest{
		BankId:   req.BankId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}

	list := make([]*types.QuestionMenus, 0, len(tree.Nodes))
	for _, node := range tree.Nodes {
		list = append(list, &types.QuestionMenus{
			Id:    node.Id,
			Title: node.Title,
		})
	}

	return l.successResp(list, tree.Total)
}

func (l *ListQuestionsLogic) errorResp(err error) (*types.GetQuestionListResp, error) {
	return &types.GetQuestionListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.QuestionList{
			Total: 0,
			List:  nil,
		},
	}, nil
}
func (l *ListQuestionsLogic) successResp(list []*types.QuestionMenus, total int64) (*types.GetQuestionListResp, error) {
	return &types.GetQuestionListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.QuestionList{
			Total: total,
			List:  list,
		},
	}, nil
}
