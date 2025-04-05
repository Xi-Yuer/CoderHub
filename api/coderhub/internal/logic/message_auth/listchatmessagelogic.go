package message_auth

import (
	"coderhub/conf"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListChatMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListChatMessageLogic 获取聊天消息列表
func NewListChatMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListChatMessageLogic {
	return &ListChatMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListChatMessageLogic) ListChatMessage(req *types.GetChatMessageListReq) (resp *types.GetPrivateChatMessageListResp, err error) {
	privateMessages, i, err := l.svcCtx.PrivateMessageRepository.GetPrivateSessionHistoryMessage(l.ctx, req.ReceiverID, req.SenderID, req.Page, req.PageSize)
	if err != nil {
		return l.errorResp(err)
	}
	if len(privateMessages) == 0 {
		return l.successResp(nil)
	}

	var data types.ChatMessageList
	for _, v := range privateMessages {
		data.List = append(data.List, &types.PrivateMessage{
			MessageID:   v.MessageID,
			SessionID:   v.SessionID,
			SenderID:    v.SenderID,
			ReceiverID:  v.ReceiverID,
			Content:     v.Content,
			ContentType: v.ContentType,
			Status:      v.Status,
			Timestamp:   v.Timestamp,
			IsRecalled:  v.IsRecalled,
			CreatedAt:   v.CreatedAt.Unix(),
			UpdatedAt:   v.UpdatedAt.Unix(),
		})
	}
	data.Total = i

	return l.successResp(&data)
}

func (l *ListChatMessageLogic) errorResp(err error) (resp *types.GetPrivateChatMessageListResp, err1 error) {
	return &types.GetPrivateChatMessageListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListChatMessageLogic) successResp(data *types.ChatMessageList) (resp *types.GetPrivateChatMessageListResp, err error) {
	return &types.GetPrivateChatMessageListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
