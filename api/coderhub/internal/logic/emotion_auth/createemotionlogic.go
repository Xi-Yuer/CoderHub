package emotion_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateEmotionLogic 创建表情包
func NewCreateEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmotionLogic {
	return &CreateEmotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEmotionLogic) CreateEmotion(req *types.CreateEmojiReq) (resp *types.CreateEmojiResp, err error) {
	_, err = l.svcCtx.EmotionService.CreateEmotion(l.ctx, &coderhub.CreateEmotionRequest{
		Code:        req.Code,
		Url:         req.URL,
		Description: req.Description,
	})
	if err != nil {
		return l.errorResp(err)
	}
	return l.successResp()
}

func (l *CreateEmotionLogic) successResp() (*types.CreateEmojiResp, error) {
	return &types.CreateEmojiResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *CreateEmotionLogic) errorResp(err error) (*types.CreateEmojiResp, error) {
	return &types.CreateEmojiResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, err
}
