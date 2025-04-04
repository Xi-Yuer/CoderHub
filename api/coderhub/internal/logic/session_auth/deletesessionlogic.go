package session_auth

import (
	"coderhub/conf"
	"coderhub/model"
	"context"

	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteSessionLogic 删除会话
func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionReq) (resp *types.DeleteSessionResp, err error) {
	userSession, err := l.svcCtx.UserSessionRepository.GetUserSession(l.ctx, &model.UserSession{
		SessionID: req.ID,
	})
	if err != nil {
		return l.errorResp(err)
	}
	err = l.svcCtx.UserSessionRepository.DeleteUserSession(l.ctx, &model.UserSession{
		UserID: userSession.UserID,
		PeerID: userSession.PeerID,
	})
	if err != nil {
		return l.errorResp(err)
	}
	err = l.svcCtx.UserSessionRepository.DeleteUserSession(l.ctx, &model.UserSession{
		UserID: userSession.PeerID,
		PeerID: userSession.UserID,
	})
	if err != nil {
		return l.errorResp(err)
	}
	return
}

func (l *DeleteSessionLogic) errorResp(err error) (resp *types.DeleteSessionResp, err1 error) {
	return &types.DeleteSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: false,
	}, nil
}

func (l *DeleteSessionLogic) successResp() (resp *types.DeleteSessionResp, err error) {
	return &types.DeleteSessionResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: true,
	}, nil
}
