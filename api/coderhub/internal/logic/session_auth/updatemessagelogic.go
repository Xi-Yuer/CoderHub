package session_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateMessageLogic 更新会话消息
func NewUpdateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageLogic {
	return &UpdateMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMessageLogic) UpdateMessage(req *types.UpdateChatSessionReq) (resp *types.UpdateChatSessionResp, err error) {
	session, err := l.svcCtx.UserSessionRepository.UpdateUserSession(l.ctx, &model.UserSession{
		SessionID:          req.SessionID,
		SessionName:        req.SessionName,
		UserID:             req.SenderID,
		PeerID:             req.PeerID,
		UnreadMessageCount: 0,
	})
	if err != nil {
		return l.errorResp(err)
	}

	return l.successResp(&types.Session{
		ID:                 session.SessionID,
		SessionName:        session.SessionName,
		UserID:             session.UserID,
		PeerID:             session.PeerID,
		LastMessageID:      session.LastMessageID,
		LastMessageContent: session.LastMessageContent,
		UnreadMessageCount: int64(session.UnreadMessageCount),
	})
}

func (l *UpdateMessageLogic) errorResp(err error) (resp *types.UpdateChatSessionResp, err1 error) {
	return &types.UpdateChatSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *UpdateMessageLogic) successResp(data *types.Session) (resp *types.UpdateChatSessionResp, err error) {
	return &types.UpdateChatSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
