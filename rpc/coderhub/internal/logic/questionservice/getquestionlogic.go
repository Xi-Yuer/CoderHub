package questionservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetQuestion 获取题目详情
func (l *GetQuestionLogic) GetQuestion(in *coderhub.GetQuestionRequest) (*coderhub.GetQuestionResponse, error) {
	question, err := l.svcCtx.QuestionRepository.GetQuestionByID(l.ctx, in.QuestionId)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserRepository.GetUserByID(question.CreateUser)
	if err != nil {
		return nil, err
	}

	return &coderhub.GetQuestionResponse{
		Id:      int64(question.ID),
		BankId:  question.BankID,
		Title:   question.Title,
		Content: question.Content,
		CreateUser: &coderhub.UserInfo{
			UserId:    user.ID,
			UserName:  user.UserName,
			Avatar:    user.Avatar.String,
			Email:     user.Email.String,
			Gender:    user.Gender,
			Age:       user.Age,
			Phone:     user.Phone.String,
			NickName:  user.NickName.String,
			IsAdmin:   user.IsAdmin,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
		Difficulty: question.Difficulty,
		CreateTime: question.CreatedAt.Unix(),
		UpdateTime: question.UpdatedAt.Unix(),
	}, nil
}
