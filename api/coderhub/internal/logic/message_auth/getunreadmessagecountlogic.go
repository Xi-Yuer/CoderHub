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

type GetUnReadMessageCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUnReadMessageCountLogic 获取是否有未读消息
func NewGetUnReadMessageCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnReadMessageCountLogic {
	return &GetUnReadMessageCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnReadMessageCountLogic) GetUnReadMessageCount(req *types.GetUnreadMessageReq) (resp *types.GetUnreadMessageResp, err error) {
	userID, err := utils.GetUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	count, err := l.svcCtx.MessageService.GetUnReadMessageCount(l.ctx, &coderhub.GetUnReadMessageCountRequest{
		UserId: userID,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(int64(count.Count))
}

func (l *GetUnReadMessageCountLogic) errorResp(err error) (resp *types.GetUnreadMessageResp, err1 error) {
	return &types.GetUnreadMessageResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.UnreadMessage{
			Total: 0,
		},
	}, nil
}

func (l *GetUnReadMessageCountLogic) successResp(total int64) (resp *types.GetUnreadMessageResp, err error) {
	return &types.GetUnreadMessageResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: types.UnreadMessage{
			Total: total,
		},
	}, nil
}
