package questionservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"coderhub/shared/utils"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type CreateQuestionBankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateQuestionBankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateQuestionBankLogic {
	return &CreateQuestionBankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateQuestionBank 创建题库
func (l *CreateQuestionBankLogic) CreateQuestionBank(in *coderhub.CreateQuestionBankRequest) (*coderhub.CreateQuestionBankResponse, error) {
	id := utils.GenID()
	err := l.svcCtx.QuestionBankRepository.CreateQuestionBank(l.ctx, &model.QuestionBank{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		Difficulty:  in.Difficulty,
		Tags:        strings.Join(in.Tags, ","),
		CreateUser:  in.CreateUser,
		CoverImage:  nil,
	})
	if err != nil {
		return nil, err
	}

	// 创建题库封面
	if in.CoverImage != nil {
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.ImageRelationRepository.Create(l.ctx, &model.ImageRelation{
			ID:         utils.GenID(),
			ImageID:    in.CoverImage.ImageId,
			EntityID:   id,
			EntityType: model.ImageRelationQuestionCover,
			Sort:       0,
		})
		if err != nil {
			return nil, err
		}
	}

	return &coderhub.CreateQuestionBankResponse{
		Success: true,
	}, nil
}
