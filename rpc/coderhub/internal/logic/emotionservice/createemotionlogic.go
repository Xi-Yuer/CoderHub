package emotionservicelogic

import (
	"coderhub/model"
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmotionLogic {
	return &CreateEmotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateEmotionLogic) CreateEmotion(in *coderhub.CreateEmotionRequest) (*coderhub.CreateEmotionResponse, error) {
	err := l.svcCtx.EmotionRepository.Create(l.ctx, &model.Emoticon{
		Code:        in.Code,
		Description: in.Description,
		URL:         in.Url,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateEmotionResponse{
		Success: true,
	}, nil
}
