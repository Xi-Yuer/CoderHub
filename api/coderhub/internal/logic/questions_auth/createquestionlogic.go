package questions_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateQuestionLogic 创建题目
func NewCreateQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionLogic {
	return &CreateQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateQuestionLogic) CreateQuestion(req *types.CreateQuestionReq) (resp *types.CreateQuestionResp, err error) {
	UserId, err := utils.GetUserID(l.ctx)
	if err != nil {
		return l.errorResp(err)
	}

	_, err = l.svcCtx.QuestionBankService.CreateQuestion(l.ctx, &coderhub.CreateQuestionRequest{
		Title:      req.Title,
		BankId:     req.BankId,
		Content:    req.Content,
		CreateUser: UserId,
		Difficulty: req.Difficulty,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *CreateQuestionLogic) errorResp(err error) (*types.CreateQuestionResp, error) {
	return &types.CreateQuestionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *CreateQuestionLogic) successResp() (*types.CreateQuestionResp, error) {
	return &types.CreateQuestionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
