package messageservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMessageLogic {
	return &CreateMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMessageLogic) CreateMessage(in *coderhub.CreateMessageRequest) (*coderhub.CreateMessageResponse, error) {
	if in.SenderId == in.ReceiverId {
		return &coderhub.CreateMessageResponse{
			Id: 0,
		}, nil
	}
	message, _ := l.svcCtx.MessageRepository.GetMessage(l.ctx, &model.Message{
		SenderID:   in.SenderId,
		ReceiverID: in.ReceiverId,
		Type:       in.Type,
		EntityID:   in.EntityId,
	})
	if message != nil {
		return &coderhub.CreateMessageResponse{
			Id: int64(message.ID),
		}, nil
	}

	err := l.svcCtx.MessageRepository.Create(l.ctx, &model.Message{
		SenderID:   in.SenderId,
		ReceiverID: in.ReceiverId,
		Type:       in.Type,
		EntityID:   in.EntityId,
		Content:    in.Content,
		IsRead:     false,
	})
	if err != nil {
		return nil, err
	}

	return &coderhub.CreateMessageResponse{
		Id: 0,
	}, nil
}
