package questionservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionTreeLogic {
	return &GetQuestionTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetQuestionTree 获取题库下的所有题目目录
func (l *GetQuestionTreeLogic) GetQuestionTree(in *coderhub.GetQuestionTreeRequest) (*coderhub.GetQuestionTreeResponse, error) {
	questions, total, err := l.svcCtx.QuestionRepository.GetQuestions(l.ctx, []int64{in.BankId}, in.Page, in.PageSize)
	if err != nil {
		return nil, err
	}

	list := make([]*coderhub.QuestionTreeNode, 0)
	for _, question := range questions {
		list = append(list, &coderhub.QuestionTreeNode{
			Id:    int64(question.ID),
			Title: question.Title,
		})
	}

	return &coderhub.GetQuestionTreeResponse{
		Nodes: list,
		Total: total,
	}, nil
}
