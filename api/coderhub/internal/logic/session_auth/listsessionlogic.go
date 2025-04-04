package session_auth

import (
	"coderhub/api/coderhub/internal/svc"
	"coderhub/api/coderhub/internal/types"
	"coderhub/conf"
	"coderhub/shared/utils"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListSessionLogic 获取会话列表
func NewListSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionLogic {
	return &ListSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSessionLogic) ListSession(req *types.GetSessionListReq) (resp *types.GetSessionListResp, err error) {

	fmt.Println("req.Page:", req.Page)
	fmt.Println("req.PageSize:", req.PageSize)
	sessions, i, err := l.svcCtx.UserSessionRepository.GetUserSessions(l.ctx, uint64(utils.String2Int(req.UserID)), req.Page, req.PageSize, req.SessionName)
	if err != nil {
		return l.errorResp(err)
	}

	var data *types.SessionList
	if len(sessions) > 0 {
		data = &types.SessionList{
			List:  make([]*types.Session, 0, len(sessions)),
			Total: i,
		}
		for _, session := range sessions {
			data.List = append(data.List, &types.Session{
				ID:                 session.SessionID,
				SessionName:        session.SessionName,
				UserID:             session.UserID,
				PeerID:             session.PeerID,
				LastMessageID:      session.LastMessageID,
				LastMessageContent: session.LastMessageContent,
				UnreadMessageCount: int64(session.UnreadMessageCount),
			})
		}
		return l.successResp(data)
	}

	return l.successResp(nil)
}

func (l *ListSessionLogic) errorResp(err error) (resp *types.GetSessionListResp, err1 error) {
	return &types.GetSessionListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpBadRequest,
			Message: err.Error(),
		},
		Data: nil,
	}, nil
}

func (l *ListSessionLogic) successResp(data *types.SessionList) (resp *types.GetSessionListResp, err error) {
	return &types.GetSessionListResp{
		Response: types.Response{
			Code:    conf.HttpCode.HttpStatusOK,
			Message: conf.HttpMessage.MsgOK,
		},
		Data: data,
	}, nil
}
