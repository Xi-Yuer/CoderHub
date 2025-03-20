package messageservicelogic

import (
	"context"

	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnReadMessageCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUnReadMessageCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnReadMessageCountLogic {
	return &GetUnReadMessageCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUnReadMessageCountLogic) GetUnReadMessageCount(in *coderhub.GetUnReadMessageCountRequest) (*coderhub.GetUnReadMessageCountResponse, error) {
	count, err := l.svcCtx.MessageRepository.GetUnReadMessageCount(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return &coderhub.GetUnReadMessageCountResponse{
		Count: count,
	}, nil
}
