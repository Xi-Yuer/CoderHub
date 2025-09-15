package questionservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateQuestion 创建题目
func (l *CreateQuestionLogic) CreateQuestion(in *coderhub.CreateQuestionRequest) (*coderhub.CreateQuestionResponse, error) {
	bank, err := l.svcCtx.QuestionBankRepository.GetQuestionBankByID(l.ctx, in.BankId)
	if err != nil {
		return nil, err
	}
	if bank == nil {
		return nil, errors.New("题库不存在")
	}
	err = l.svcCtx.QuestionRepository.CreateQuestion(l.ctx, &model.Question{
		BankID:     in.BankId,
		Title:      in.Title,
		Content:    in.Content,
		CreateUser: in.CreateUser,
		Difficulty: in.Difficulty,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateQuestionResponse{
		Success: true,
	}, nil
}
