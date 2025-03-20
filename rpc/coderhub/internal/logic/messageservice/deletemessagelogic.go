package messageservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageLogic {
	return &DeleteMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMessageLogic) DeleteMessage(in *coderhub.DeleteMessageRequest) (*coderhub.DeleteMessageResponse, error) {
	err := l.svcCtx.MessageRepository.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &coderhub.DeleteMessageResponse{
		Success: true,
	}, nil
}
