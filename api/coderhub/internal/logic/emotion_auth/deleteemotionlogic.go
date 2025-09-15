package emotion_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmotionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除表情包
func NewDeleteEmotionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmotionLogic {
	return &DeleteEmotionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEmotionLogic) DeleteEmotion(req *types.DeleteEmojiReq) (resp *types.DeleteEmojiResp, err error) {
	_, err = l.svcCtx.EmotionService.DeleteEmotion(l.ctx, &coderhub.DeleteEmotionRequest{
		Id: utils.String2Int(req.Id),
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteEmotionLogic) successResp() (*types.DeleteEmojiResp, error) {
	return &types.DeleteEmojiResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}

func (l *DeleteEmotionLogic) errorResp(err error) (*types.DeleteEmojiResp, error) {
	return &types.DeleteEmojiResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, err
}
