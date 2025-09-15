package emotion_public

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEmotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListEmotionLogic 获取表情包列表
func NewListEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEmotionLogic {
	return &ListEmotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEmotionLogic) ListEmotion(req *types.GetEmojiListReq) (resp *types.GetEmojiListResp, err error) {
	emotionList, err := l.svcCtx.EmotionService.GetEmotionList(l.ctx, &coderhub.GetEmotionListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return l.errorResp(err)
	}
	list := make([]*types.Emoji, 0, len(emotionList.Emotions))
	for _, v := range emotionList.Emotions {
		list = append(list, &types.Emoji{
			ID:          utils.Int2String(v.Id),
			Code:        v.Code,
			Description: v.Description,
			URL:         v.Url,
			CreatedAt:   v.CreateTime,
			UpdatedAt:   v.CreateTime,
		})
	}
	return l.successResp(list, emotionList.Total)
}

func (l *ListEmotionLogic) successResp(list []*types.Emoji, total int64) (*types.GetEmojiListResp, error) {
	return &types.GetEmojiListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.EmojiList{
			Total: total,
			List:  list,
		},
	}, nil
}

func (l *ListEmotionLogic) errorResp(err error) (*types.GetEmojiListResp, error) {
	return &types.GetEmojiListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.EmojiList{
			Total: 0,
			List:  nil,
		},
	}, err
}
