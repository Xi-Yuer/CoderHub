package messageservicelogic

import (
	"coderhub/model"
	"coderhub/rpc/coderhub/coderhub"
	"coderhub/rpc/coderhub/internal/svc"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *coderhub.GetMessageListRequest) (*coderhub.GetMessageListResponse, error) {
	list, i, err := l.svcCtx.MessageRepository.List(l.ctx, &model.Message{
		ReceiverID: in.UserId,
		Type:       in.Type,
	}, int64(in.Page), int64(in.PageSize))
	if err != nil {
		return nil, err
	}
	var messages []*coderhub.MessageListItem
	for _, v := range list {
		messages = append(messages, &coderhub.MessageListItem{
			Id:         int64(v.ID),
			SenderId:   v.SenderID,
			ReceiverId: v.ReceiverID,
			Type:       v.Type,
			EntityId:   v.EntityID,
			Content:    v.Content,
			CreateTime: v.CreatedAt.Unix(),
			IsRead:     v.IsRead,
			IsDelete:   v.DeletedAt.Valid,
		})
	}

	return &coderhub.GetMessageListResponse{
		List:  messages,
		Total: int32(i),
	}, nil
}
