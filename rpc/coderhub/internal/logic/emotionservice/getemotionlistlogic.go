package emotionservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmotionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmotionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmotionListLogic {
	return &GetEmotionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmotionListLogic) GetEmotionList(in *coderhub.GetEmotionListRequest) (*coderhub.GetEmotionListResponse, error) {
	emoticons, i, err := l.svcCtx.EmotionRepository.List(l.ctx, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}

	list := make([]*coderhub.Emotion, 0)
	for _, v := range emoticons {
		list = append(list, &coderhub.Emotion{
			Id:          int64(v.ID),
			Code:        v.Code,
			Url:         v.URL,
			Description: v.Description,
			CreateTime:  v.CreatedAt.Unix(),
		})
	}

	return &coderhub.GetEmotionListResponse{
		Emotions: list,
		Total:    i,
	}, nil
}
