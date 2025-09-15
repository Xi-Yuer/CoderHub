package emotionservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmotionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmotionLogic {
	return &DeleteEmotionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteEmotionLogic) DeleteEmotion(in *coderhub.DeleteEmotionRequest) (*coderhub.DeleteEmotionResponse, error) {
	err := l.svcCtx.EmotionRepository.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteEmotionResponse{
		Success: true,
	}, nil
}
