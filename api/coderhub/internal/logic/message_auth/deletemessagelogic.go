package message_auth

import (
	"coderhub/conf"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteMessageLogic 删除消息
func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessageLogic) DeleteMessage(req *types.DeleteMessage) (resp *types.DeleteMessageResp, err error) {
	if _, err = l.svcCtx.MessageService.DeleteMessage(l.ctx, &coderhub.DeleteMessageRequest{
		Id:     utils.String2Int(req.ID),
		UserId: utils.String2Int(req.RequestUserId),
	}); err != nil {
		return l.errorResp(err)
	}

	return l.successResp()
}

func (l *DeleteMessageLogic) errorResp(err error) (resp *types.DeleteMessageResp, err1 error) {
	return &types.DeleteMessageResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteMessageLogic) successResp() (resp *types.DeleteMessageResp, err error) {
	return &types.DeleteMessageResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
