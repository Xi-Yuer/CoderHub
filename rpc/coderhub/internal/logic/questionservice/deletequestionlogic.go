package questionservicelogic

import (
	"context"
	"errors"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteQuestion 删除题目
func (l *DeleteQuestionLogic) DeleteQuestion(in *coderhub.DeleteQuestionRequest) (*coderhub.DeleteQuestionResponse, error) {
	questionBank, err := l.svcCtx.QuestionRepository.GetQuestionByID(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if questionBank.CreateUser != in.CreateUser {
		return nil, errors.New("非法操作")
	}

	if questionBank == nil {
		return nil, errors.New("题目不存在")
	}

	if err := l.svcCtx.QuestionRepository.DeleteQuestion(l.ctx, in.Id); err != nil {
		return nil, err
	}

	return &coderhub.DeleteQuestionResponse{
		Success: true,
	}, nil
}
