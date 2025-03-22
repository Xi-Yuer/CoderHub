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

	// 不关心消息是否创建成功
	_ = l.svcCtx.MessageRepository.Create(l.ctx, &model.Message{
		SenderID:   in.SenderId,
		ReceiverID: in.ReceiverId,
		Type:       in.Type,
		EntityID:   in.EntityId,
		Content:    in.Content,
		IsRead:     false,
	})

	return &coderhub.CreateMessageResponse{
		Id: 0,
	}, nil
}
