package session_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"coderhub/shared/utils"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateSessionLogic 创建一个会话
func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq) (resp *types.CreateSessionResp, err error) {

	// 接受者信息
	peerUser, err := l.svcCtx.UserRepository.GetUserByID(utils.String2Int(req.PeerID))
	if err != nil {
		return l.errorResp(err)
	}
	if peerUser == nil {
		return l.errorResp(err)
	}
	UserID, err := utils.GetUserID(l.ctx)

	if err != nil {
		return l.errorResp(err)
	}

	// 发送者信息
	senderUser, err := l.svcCtx.UserRepository.GetUserByID(UserID)
	if err != nil {
		return l.errorResp(err)
	}

	// 创建会话
	// 查询会话是否存在
	userSession, _ := l.svcCtx.UserSessionRepository.GetUserSession(l.ctx, &model.UserSession{
		UserID: utils.Int2String(UserID),
		PeerID: req.PeerID,
	})

	if userSession != nil {
		return l.successResp(types.Session{
			ID:                 userSession.SessionID,
			SessionName:        userSession.SessionName,
			UserID:             userSession.UserID,
			PeerID:             userSession.PeerID,
			LastMessageID:      userSession.LastMessageID,
			LastMessageContent: userSession.LastMessageContent,
			UnreadMessageCount: int64(userSession.UnreadMessageCount),
		})
	}

	peerSession := &model.UserSession{
		SessionID:   utils.Int2String(utils.GenID()),
		SessionName: senderUser.UserName,
		PeerID:      utils.Int2String(UserID),
		UserID:      req.PeerID,
	}
	userNickName := senderUser.NickName
	if userNickName.String != "" {
		peerSession.SessionName = userNickName.String
	}
	senderSession := &model.UserSession{
		SessionID:   utils.Int2String(utils.GenID()),
		SessionName: peerUser.UserName,
		UserID:      utils.Int2String(UserID),
		PeerID:      req.PeerID,
	}
	peerUserNickName := peerUser.NickName
	if peerUserNickName.String != "" {
		senderSession.SessionName = peerUserNickName.String
	}
	session1, err := l.svcCtx.UserSessionRepository.Create(l.ctx, peerSession)

	if err != nil || session1 == nil {
		return l.errorResp(err)
	}

	session2, err := l.svcCtx.UserSessionRepository.Create(l.ctx, senderSession)

	if err != nil || session2 == nil {
		return l.errorResp(err)
	}
	return l.successResp(types.Session{
		ID:                 session2.SessionID,
		SessionName:        session2.SessionName,
		UserID:             session2.UserID,
		PeerID:             session2.PeerID,
		LastMessageID:      session2.LastMessageID,
		LastMessageContent: session2.LastMessageContent,
		UnreadMessageCount: int64(session2.UnreadMessageCount),
	})
}

func (l *CreateSessionLogic) errorResp(err error) (resp *types.CreateSessionResp, err1 error) {
	return &types.CreateSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: types.Session{},
	}, nil
}

func (l *CreateSessionLogic) successResp(data types.Session) (resp *types.CreateSessionResp, err error) {
	return &types.CreateSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
